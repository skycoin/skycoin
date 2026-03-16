package commands

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/skycoin/skycoin/src/btc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

// btcHandler holds the Bitcoin backend and wallet service for handling HTTP requests
type btcHandler struct {
	backend    btc.Backend
	wltService *wallet.Service
}

// handleBtcAPI routes Bitcoin-specific API requests
func (h *btcHandler) handleBtcAPI(c *gin.Context, apiPath string) bool {
	path := strings.TrimSuffix(apiPath, "/")
	method := c.Request.Method

	switch {
	case path == "/v1/btc/balance" && method == http.MethodGet:
		h.getBalance(c)
		return true
	case path == "/v1/btc/utxos" && method == http.MethodGet:
		h.getUTXOs(c)
		return true
	case path == "/v1/btc/history" && method == http.MethodGet:
		h.getHistory(c)
		return true
	case path == "/v1/btc/send" && method == http.MethodPost:
		h.sendTransaction(c)
		return true
	case path == "/v1/btc/fee" && method == http.MethodGet:
		h.estimateFee(c)
		return true
	}

	return false
}

// getBalance returns balances for comma-separated addresses
func (h *btcHandler) getBalance(c *gin.Context) {
	addrsParam := c.Query("addrs")
	if addrsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addrs parameter required"})
		return
	}

	addresses := strings.Split(addrsParam, ",")
	balances, err := h.backend.GetBalance(addresses)
	if err != nil {
		log.Printf("[BTC] GetBalance error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Compute totals
	var totalConfirmed int64
	addrBalances := make(map[string]map[string]map[string]int64, len(balances))
	for addr, bal := range balances {
		totalConfirmed += bal.Confirmed
		addrBalances[addr] = map[string]map[string]int64{
			"confirmed": {"coins": bal.Confirmed},
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"confirmed": gin.H{"coins": totalConfirmed},
		"predicted": gin.H{"coins": totalConfirmed},
		"addresses": addrBalances,
	})
}

// getUTXOs returns unspent outputs for comma-separated addresses
func (h *btcHandler) getUTXOs(c *gin.Context) {
	addrsParam := c.Query("addrs")
	if addrsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addrs parameter required"})
		return
	}

	addresses := strings.Split(addrsParam, ",")
	utxos, err := h.backend.ListUnspent(addresses)
	if err != nil {
		log.Printf("[BTC] ListUnspent error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utxos)
}

// getHistory returns transaction history for comma-separated addresses
func (h *btcHandler) getHistory(c *gin.Context) {
	addrsParam := c.Query("addrs")
	if addrsParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "addrs parameter required"})
		return
	}

	addresses := strings.Split(addrsParam, ",")
	history, err := h.backend.GetHistory(addresses)
	if err != nil {
		log.Printf("[BTC] GetHistory error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// sendTransaction builds, signs, and broadcasts a Bitcoin transaction
func (h *btcHandler) sendTransaction(c *gin.Context) {
	var req btc.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}

	if req.WalletID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_id required"})
		return
	}
	if len(req.Destinations) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "at least one destination required"})
		return
	}
	if req.FeeRate <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fee_rate must be positive"})
		return
	}

	// Get wallet
	wlt, err := h.wltService.GetWallet(req.WalletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("wallet not found: %v", err)})
		return
	}

	// Collect addresses and keys from wallet
	entries, err := wlt.GetEntries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("get wallet entries: %v", err)})
		return
	}

	addresses := make([]string, 0, len(entries))
	keys := make(map[string]cipher.SecKey, len(entries))
	for _, entry := range entries {
		addr := entry.Address.String()
		addresses = append(addresses, addr)
		keys[addr] = entry.Secret
	}

	// Calculate total amount needed
	var totalNeeded int64
	destinations := make([]btc.TxDestination, len(req.Destinations))
	for i, dest := range req.Destinations {
		totalNeeded += dest.Coins
		destinations[i] = btc.TxDestination{
			Address: dest.Address,
			Value:   dest.Coins,
		}
	}

	// Get UTXOs
	utxos, err := h.backend.ListUnspent(addresses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("list unspent: %v", err)})
		return
	}

	// Select UTXOs (need to account for fee)
	estimatedFee := int64(btc.EstimatedTxSize(1, len(destinations)+1)) * req.FeeRate
	selectedUTXOs, _, err := btc.SelectUTXOs(utxos, totalNeeded+estimatedFee)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}

	// Use first address as change address
	changeAddr := addresses[0]

	// Build and sign transaction
	rawTx, err := btc.BuildTransaction(selectedUTXOs, destinations, changeAddr, req.FeeRate, keys)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("build transaction: %v", err)})
		return
	}

	// Broadcast
	txid, err := h.backend.BroadcastTransaction(rawTx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("broadcast failed: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"txid": txid})
}

// estimateFee returns fee estimation for a given confirmation target
func (h *btcHandler) estimateFee(c *gin.Context) {
	blocksStr := c.DefaultQuery("blocks", "6")
	blocks, err := strconv.Atoi(blocksStr)
	if err != nil || blocks < 1 {
		blocks = 6
	}

	satPerByte, err := h.backend.EstimateFee(blocks)
	if err != nil {
		log.Printf("[BTC] EstimateFee error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sat_per_byte": satPerByte,
		"blocks":       blocks,
	})
}
