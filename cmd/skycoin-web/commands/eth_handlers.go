package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/eth"
	"github.com/skycoin/skycoin/src/wallet"
)

// ethToken is an ERC-20 token held at Ethereum wallet addresses.
type ethToken struct {
	Symbol   string
	Contract eth.Address
	Decimals int
}

// ethHandler serves /v1/eth/... for Ether, or for one token when token is set.
type ethHandler struct {
	rpc        *eth.RPC
	indexer    *eth.Indexer
	wltService *wallet.Service
	token      *ethToken
	endpoints  *ethEndpoints
}

const ethRequestTimeout = 60 * time.Second

func (h *ethHandler) handleEthAPI(c *webCtx, apiPath string) bool {
	path := strings.TrimSuffix(apiPath, "/")
	method := c.Request.Method
	switch {
	case path == "/v1/eth/balance" && method == http.MethodGet:
		h.balance(c)
	case path == "/v1/eth/history" && method == http.MethodGet:
		h.history(c)
	case path == "/v1/eth/fee" && method == http.MethodGet:
		h.fee(c)
	case path == "/v1/eth/send" && method == http.MethodPost:
		h.send(c)
	case path == "/v1/eth/endpoints" && method == http.MethodGet:
		h.endpoints.get(c)
	case path == "/v1/eth/endpoints" && method == http.MethodPut:
		h.endpoints.put(c)
	case path == "/v1/eth/health" && method == http.MethodGet:
		h.health(c)
	default:
		return false
	}
	return true
}

func (h *ethHandler) ctx(c *webCtx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), ethRequestTimeout)
}

func parseEthAddrs(c *webCtx) ([]eth.Address, bool) {
	param := c.Query("addrs")
	if param == "" {
		c.JSON(http.StatusBadRequest, H{"error": "addrs parameter required"})
		return nil, false
	}
	var addrs []eth.Address
	for _, s := range strings.Split(param, ",") {
		a, err := eth.ParseAddress(strings.TrimSpace(s))
		if err != nil {
			c.JSON(http.StatusBadRequest, H{"error": fmt.Sprintf("%s: %v", s, err)})
			return nil, false
		}
		addrs = append(addrs, a)
	}
	return addrs, true
}

// addressBalance is the balance h serves, Ether in wei or the token in its base units.
func (h *ethHandler) addressBalance(ctx context.Context, a eth.Address) (*big.Int, error) {
	if h.token != nil {
		return h.rpc.TokenBalance(ctx, h.token.Contract, a)
	}
	return h.rpc.Balance(ctx, a)
}

// balance answers in the shape of /v1/btc/balance, amounts in base units as strings.
func (h *ethHandler) balance(c *webCtx) {
	addrs, ok := parseEthAddrs(c)
	if !ok {
		return
	}
	ctx, cancel := h.ctx(c)
	defer cancel()

	total := new(big.Int)
	perAddress := make(map[string]H, len(addrs))
	for _, a := range addrs {
		b, err := h.addressBalance(ctx, a)
		if err != nil {
			log.Printf("[ETH] balance %s: %v", a, err)
			c.JSON(http.StatusBadGateway, H{"error": err.Error()})
			return
		}
		total.Add(total, b)
		perAddress[a.String()] = H{"confirmed": H{"coins": b.String()}}
	}
	c.JSON(http.StatusOK, H{
		"confirmed": H{"coins": total.String()},
		"predicted": H{"coins": total.String()},
		"addresses": perAddress,
	})
}

// history answers in the shape of /v1/btc/history: the sender as the one
// input and the recipient and amount as the one output.
func (h *ethHandler) history(c *webCtx) {
	addrs, ok := parseEthAddrs(c)
	if !ok {
		return
	}
	ctx, cancel := h.ctx(c)
	defer cancel()

	seen := make(map[string]bool)
	var txs []H
	for _, a := range addrs {
		var list []eth.IndexedTx
		var err error
		if h.token != nil {
			list, err = h.indexer.TokenTransfers(ctx, a.String(), h.token.Contract.String())
		} else {
			list, err = h.indexer.Transactions(ctx, a.String())
		}
		if err != nil {
			log.Printf("[ETH] history %s: %v", a, err)
			c.JSON(http.StatusBadGateway, H{"error": err.Error()})
			return
		}
		for _, tx := range list {
			// A transfer between two of the wallet's addresses is listed for both.
			if seen[tx.Hash] {
				continue
			}
			seen[tx.Hash] = true
			value := tx.Value
			if tx.Failed {
				value = new(big.Int)
			}
			entry := H{
				"txid":          tx.Hash,
				"timestamp":     tx.Timestamp,
				"block_height":  tx.Block,
				"confirmations": confirmations(tx.Block),
				"failed":        tx.Failed,
				"inputs":        []H{{"address": checksummed(tx.From)}},
				"outputs":       []H{{"address": checksummed(tx.To), "value": value.String()}},
			}
			if h.token == nil {
				entry["fee"] = tx.Fee.String()
			}
			txs = append(txs, entry)
		}
	}
	sort.SliceStable(txs, func(i, j int) bool { return txs[i]["timestamp"].(int64) > txs[j]["timestamp"].(int64) })
	if txs == nil {
		txs = []H{}
	}
	c.JSON(http.StatusOK, txs)
}

func confirmations(block uint64) int {
	if block > 0 {
		return 1
	}
	return 0
}

// checksummed is the EIP-55 form of an indexer's lowercase address, or s itself.
func checksummed(s string) string {
	a, err := eth.ParseAddress(s)
	if err != nil {
		return s
	}
	return a.String()
}

// tokenTransferGas is the fee estimate's gas for a token, before the node
// estimates the real transfer.
const tokenTransferGas = 65000

func (h *ethHandler) fee(c *webCtx) {
	ctx, cancel := h.ctx(c)
	defer cancel()
	fees, base, err := eth.SuggestFees(ctx, h.rpc)
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": err.Error()})
		return
	}
	gas := int64(21000)
	if h.token != nil {
		gas = tokenTransferGas
	}
	presets := H{}
	for name, f := range eth.FeePresets(base, fees) {
		presets[name] = H{
			"max_priority_fee_per_gas": f.MaxPriorityFeePerGas.String(),
			"max_fee_per_gas":          f.MaxFeePerGas.String(),
			"max_fee":                  new(big.Int).Mul(big.NewInt(gas), f.MaxFeePerGas).String(),
		}
	}
	c.JSON(http.StatusOK, H{
		"base_fee":                 base.String(),
		"presets":                  presets,
		"max_priority_fee_per_gas": fees.MaxPriorityFeePerGas.String(),
		"max_fee_per_gas":          fees.MaxFeePerGas.String(),
		"gas":                      gas,
		"max_fee":                  new(big.Int).Mul(big.NewInt(gas), fees.MaxFeePerGas).String(),
	})
}

func (h *ethHandler) health(c *webCtx) {
	ctx, cancel := h.ctx(c)
	defer cancel()
	id, err := h.rpc.ChainID(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": err.Error()})
		return
	}
	block, err := h.rpc.BlockNumber(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, H{"chain_id": id.String(), "tip_height": block})
}

type ethSendRequest struct {
	WalletID string `json:"wallet_id"`
	Password string `json:"password"`
	// From is the address to send from, or empty for the first that can pay.
	From   string `json:"from"`
	To     string `json:"to"`
	Amount string `json:"amount"`
	// MaxFeePerGas and MaxPriorityFeePerGas, in wei, replace the suggested fees.
	MaxFeePerGas         string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas string `json:"max_priority_fee_per_gas"`
}

type ethKey struct {
	address eth.Address
	secret  cipher.SecKey
}

func (h *ethHandler) send(c *webCtx) {
	var req ethSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, H{"error": fmt.Sprintf("invalid request: %v", err)})
		return
	}
	to, err := eth.ParseAddress(req.To)
	if err != nil {
		c.JSON(http.StatusBadRequest, H{"error": "invalid destination address"})
		return
	}
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		c.JSON(http.StatusBadRequest, H{"error": "amount must be a positive integer in base units"})
		return
	}

	keys, err := h.walletKeys(req.WalletID, req.Password)
	if err != nil {
		handleWalletError(c, err)
		return
	}
	defer func() {
		for i := range keys {
			keys[i].secret = cipher.SecKey{}
		}
	}()

	ctx, cancel := h.ctx(c)
	defer cancel()
	fees, err := h.sendFees(ctx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	key, err := h.chooseSender(ctx, keys, req.From, amount, fees)
	if err != nil {
		c.JSON(http.StatusBadRequest, H{"error": err.Error()})
		return
	}

	var token *eth.Address
	if h.token != nil {
		token = &h.token.Contract
	}
	tx, signed, err := eth.Transfer(ctx, h.rpc, eth.TransferRequest{From: key.secret, To: to, Amount: amount, Token: token, Fees: &fees})
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": err.Error()})
		return
	}
	ether, err := h.rpc.Balance(ctx, key.address)
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": err.Error()})
		return
	}
	need := tx.MaxFee()
	if h.token == nil {
		need.Add(need, amount)
	}
	if ether.Cmp(need) < 0 {
		c.JSON(http.StatusBadRequest, H{"error": "not enough ETH to pay the amount and the fee"})
		return
	}

	txid, err := h.rpc.SendRaw(ctx, signed.Raw)
	if err != nil {
		c.JSON(http.StatusBadGateway, H{"error": fmt.Sprintf("broadcast failed: %v", err)})
		return
	}
	c.JSON(http.StatusOK, H{"txid": txid, "from": key.address.String(), "max_fee": tx.MaxFee().String()})
}

// walletKeys are the wallet's external addresses and keys, unlocking it with
// password when it is encrypted.
func (h *ethHandler) walletKeys(walletID, password string) ([]ethKey, error) {
	if walletID == "" {
		return nil, wallet.NewError(errors.New("wallet_id required"))
	}
	var keys []ethKey
	err := h.wltService.ViewSecrets(walletID, []byte(password), func(w wallet.Wallet) error {
		entries, err := w.GetEntries(wallet.OptionExternal())
		if err != nil {
			return err
		}
		for _, e := range entries {
			a, err := eth.ParseAddress(e.Address.String())
			if err != nil {
				return err
			}
			keys = append(keys, ethKey{address: a, secret: e.Secret})
		}
		return nil
	})
	return keys, err
}

func (h *ethHandler) sendFees(ctx context.Context, req ethSendRequest) (eth.Fees, error) {
	if req.MaxFeePerGas == "" {
		f, _, err := eth.SuggestFees(ctx, h.rpc)
		return f, err
	}
	feeCap, ok1 := new(big.Int).SetString(req.MaxFeePerGas, 10)
	tip, ok2 := new(big.Int).SetString(req.MaxPriorityFeePerGas, 10)
	if !ok1 || !ok2 || feeCap.Sign() <= 0 || tip.Sign() < 0 || tip.Cmp(feeCap) > 0 {
		return eth.Fees{}, errors.New("give both fee values in wei, the tip no more than the cap")
	}
	return eth.Fees{MaxFeePerGas: feeCap, MaxPriorityFeePerGas: tip}, nil
}

// chooseSender is the requested address, or the first whose balance covers
// the amount, plus the most the fee can be when sending Ether.
func (h *ethHandler) chooseSender(ctx context.Context, keys []ethKey, from string, amount *big.Int, fees eth.Fees) (ethKey, error) {
	if from != "" {
		a, err := eth.ParseAddress(from)
		if err != nil {
			return ethKey{}, errors.New("invalid from address")
		}
		for _, k := range keys {
			if k.address == a {
				return k, nil
			}
		}
		return ethKey{}, errors.New("the from address is not in this wallet")
	}
	need := new(big.Int).Set(amount)
	if h.token == nil {
		need.Add(need, new(big.Int).Mul(big.NewInt(21000), fees.MaxFeePerGas))
	}
	for _, k := range keys {
		b, err := h.addressBalance(ctx, k.address)
		if err != nil {
			return ethKey{}, err
		}
		if b.Cmp(need) >= 0 {
			return k, nil
		}
	}
	return ethKey{}, errors.New("no address in this wallet holds enough to send this")
}
