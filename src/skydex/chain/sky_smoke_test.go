//go:build livenet

// Live end-to-end smoke test for the SKY escrow spend path against a real,
// synced Skycoin fullnode. It derives the escrow wallet from a seed, reads the
// wallet's live unspent outputs, and builds + locally signs a real delivery
// transaction. It is DRY-RUN by default: it does NOT broadcast unless explicitly
// told to, so it is safe to run repeatedly.
//
// The seed is read from the environment and never hardcoded or logged.
//
//	# dry-run (build + sign only, no broadcast):
//	SKY_NODE_URL=http://45.56.79.245:6420 \
//	SKY_WALLET_SEED='your escrow seed' \
//	SKY_DEST=<destination SKY address> \
//	SKY_SEND=0.001 \
//	go test -tags livenet -run TestLiveSkySpend -v ./internal/skydex-market/chain/
//
//	# real broadcast (adds SKY_BROADCAST=1):
//	... SKY_BROADCAST=1 go test -tags livenet -run TestLiveSkySpend -v ./...
package chain

import (
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLiveSkySpend(t *testing.T) {
	nodeURL := os.Getenv("SKY_NODE_URL")
	seed := os.Getenv("SKY_WALLET_SEED")
	dest := os.Getenv("SKY_DEST")
	if nodeURL == "" || seed == "" || dest == "" {
		t.Skip("set SKY_NODE_URL, SKY_WALLET_SEED and SKY_DEST to run the live spend smoke test")
	}
	sendStr := os.Getenv("SKY_SEND")
	if sendStr == "" {
		sendStr = "0.001"
	}
	amount, err := strconv.ParseFloat(sendStr, 64)
	if err != nil {
		t.Fatalf("SKY_SEND is not a number: %v", err)
	}

	// wallet_sky (SKY_WALLET) is optional; if given, the seed is cross-checked.
	node := NewSkyNode(nodeURL, seed, os.Getenv("SKY_WALLET"), 2, &http.Client{Timeout: 60 * time.Second})
	escrow := node.EscrowAddress()
	if escrow == "" {
		t.Fatalf("escrow wallet not initialized: %v", node.spendErr)
	}
	t.Logf("escrow address (fund this): %s", escrow)

	// Confirm the escrow wallet actually holds spendable SKY on the live chain.
	bal, err := node.Balance(escrow)
	if err != nil {
		t.Fatalf("read escrow balance: %v", err)
	}
	t.Logf("escrow live balance: %.6f SKY, %d coin-hours (across %d confirmed outputs)", bal.CoinsSKY, bal.Hours, bal.Outputs)
	if bal.CoinsSKY < amount {
		t.Skipf("escrow holds %.6f SKY, need at least %.6f — fund %s and retry", bal.CoinsSKY, amount, escrow)
	}

	// Build + locally sign a real delivery transaction against live UTXOs.
	txn, cl, err := node.prepareSpend(dest, amount)
	if err != nil {
		t.Fatalf("build+sign spend: %v", err)
	}
	rawtx, err := txn.SerializeHex()
	if err != nil {
		t.Fatalf("serialize txn: %v", err)
	}
	t.Logf("built + signed delivery txn: id=%s inputs=%d outputs=%d fully_signed=%v",
		txn.Hash().Hex(), len(txn.In), len(txn.Out), txn.IsFullySigned())
	for i, o := range txn.Out {
		t.Logf("  out[%d] -> %s : %d droplets, %d hours", i, o.Address, o.Coins, o.Hours)
	}
	t.Logf("raw tx (%d bytes): %s", len(rawtx)/2, rawtx)

	if os.Getenv("SKY_BROADCAST") != "1" {
		t.Log("DRY-RUN: not broadcasting (set SKY_BROADCAST=1 to inject this transaction)")
		return
	}

	txid, err := cl.InjectTransaction(txn)
	if err != nil {
		t.Fatalf("broadcast: %v", err)
	}
	t.Logf("BROADCAST OK — network accepted txn %s (%.6f SKY -> %s)", txid, amount, dest)
}
