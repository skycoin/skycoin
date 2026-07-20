package chain

import (
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// Two distinct, syntactically valid SHA256 hex source-transaction ids. They must
// differ so the two crafted outputs hash to distinct uxids (the uxid is derived
// from the output body, which includes the source-transaction id).
const (
	srcTxA = "0101010101010101010101010101010101010101010101010101010101010101"
	srcTxB = "0202020202020202020202020202020202020202020202020202020202020202"
)

// deriveAddr returns the first deterministic address for a seed (same derivation
// the escrow signer uses).
func deriveAddr(t *testing.T, seed string) string {
	t.Helper()
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 1)
	addr, err := cipher.AddressFromSecKey(seckeys[0])
	if err != nil {
		t.Fatalf("derive addr: %v", err)
	}
	return addr.String()
}

func escrowOutputs(addr string, headTime uint64) readable.UnspentOutputs {
	// Two 10 SKY outputs with ample coin-hours to cover the burn fee.
	return readable.UnspentOutputs{
		{Hash: srcTxA, Time: headTime, BkSeq: 1, SourceTransaction: srcTxA, Address: addr, Coins: "10.000000", Hours: 1000, CalculatedHours: 1000},
		{Hash: srcTxB, Time: headTime, BkSeq: 1, SourceTransaction: srcTxB, Address: addr, Coins: "10.000000", Hours: 1000, CalculatedHours: 1000},
	}
}

func TestNewSigner_MatchAndMismatch(t *testing.T) {
	seed := "market escrow hot wallet seed one"
	addr := deriveAddr(t, seed)

	// Matching wallet_sky is accepted.
	if _, err := newSigner(seed, addr); err != nil {
		t.Fatalf("newSigner with matching addr: %v", err)
	}
	// Empty wallet_sky (no cross-check) is accepted.
	if _, err := newSigner(seed, ""); err != nil {
		t.Fatalf("newSigner with empty addr: %v", err)
	}
	// A wrong wallet_sky is rejected, guarding against a seed/address mismatch.
	if _, err := newSigner(seed, "2bvxvirrH5ttckcHqLo5Po4yeyq6KsZoc8j"); err == nil {
		t.Fatal("newSigner with mismatched addr: expected error, got nil")
	}
}

func TestSendSKY_DisabledWithoutSeed(t *testing.T) {
	n := NewSkyNode("http://node.invalid", "", "", 2, nil)
	if n.EscrowAddress() != "" {
		t.Fatalf("expected no escrow address without seed, got %q", n.EscrowAddress())
	}
	if _, err := n.SendCoin(deriveAddr(t, "buyer"), 1.0); err == nil {
		t.Fatal("expected SendSKY to fail with no seed configured")
	}
}

func TestBuildSignedTxn_DeliversAndReturnsChange(t *testing.T) {
	seed := "market escrow hot wallet seed two"
	escrow := deriveAddr(t, seed)
	n := NewSkyNode("http://node.invalid", seed, escrow, 2, nil)
	if n.signer == nil {
		t.Fatalf("signer not initialized: %v", n.spendErr)
	}

	var headTime uint64 = 1_600_000_000
	dst, err := cipher.DecodeBase58Address(deriveAddr(t, "the buyer wallet"))
	if err != nil {
		t.Fatalf("decode dst: %v", err)
	}
	const sendCoins uint64 = 3_000_000 // 3 SKY

	txn, err := n.buildSignedTxn(dst, sendCoins, escrowOutputs(escrow, headTime), headTime)
	if err != nil {
		t.Fatalf("buildSignedTxn: %v", err)
	}

	if !txn.IsFullySigned() {
		t.Fatal("transaction is not fully signed")
	}
	if err := txn.Verify(); err != nil {
		t.Fatalf("txn.Verify: %v", err)
	}

	// Expect exactly two outputs: the delivery and the change.
	if len(txn.Out) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(txn.Out))
	}
	if txn.Out[0].Address != dst {
		t.Errorf("output 0 address = %s, want dst %s", txn.Out[0].Address, dst)
	}
	if txn.Out[0].Coins != sendCoins {
		t.Errorf("output 0 coins = %d, want %d", txn.Out[0].Coins, sendCoins)
	}

	change := txn.Out[1]
	if change.Address.String() != escrow {
		t.Errorf("change address = %s, want escrow %s", change.Address, escrow)
	}
	// Spend selection minimizes inputs: a single 10-SKY output covers the 3 SKY
	// delivery, leaving 7 SKY change.
	wantChange, _ := droplet.FromString("7.000000") //nolint
	if change.Coins != wantChange {
		t.Errorf("change coins = %d, want %d", change.Coins, wantChange)
	}
}

func TestSendSKY_RejectsBadPrecision(t *testing.T) {
	seed := "market escrow hot wallet seed three"
	n := NewSkyNode("http://node.invalid", seed, "", 2, nil)
	// 1.0001 SKY has 4 decimals; mainnet allows at most 3 (MaxDropletPrecision).
	_, err := n.SendCoin(deriveAddr(t, "buyer"), 1.0001)
	if err == nil || !strings.Contains(err.Error(), "precision") {
		t.Fatalf("expected precision error, got %v", err)
	}
}
