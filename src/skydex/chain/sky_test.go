package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestDepositConfirmed verifies the SKY deposit check matches an exact-amount
// deposit that is sufficiently confirmed, sent from the seller's address, and
// within the listing window — and rejects on amount, sender, or time mismatch.
func TestDepositConfirmed(t *testing.T) {
	const wallet = "2sky...market"
	const seller = "2seller...addr"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/transactions" {
			http.NotFound(w, r)
			return
		}
		// One confirmed tx from the seller paying exactly 10 SKY to the market
		// wallet at block time 1000, 3 confs.
		resp := []map[string]any{
			{
				"status": map[string]any{"confirmed": true, "height": 3, "block_seq": 100},
				"txn": map[string]any{
					"txid":      "tx-good",
					"timestamp": 1000,
					"inputs":    []map[string]any{{"owner": seller}},
					"outputs":   []map[string]any{{"dst": wallet, "coins": "10.000000"}},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer srv.Close()

	node := NewSkyNode(srv.URL, "", "", 2, srv.Client())
	nb, na := time.Unix(500, 0), time.Unix(2000, 0) // window covers block time 1000

	// Exact amount, right sender, within window.
	ok, txid, err := node.DepositConfirmed(wallet, seller, 10.0, nb, na)
	if err != nil {
		t.Fatalf("DepositConfirmed: %v", err)
	}
	if !ok || txid != "tx-good" {
		t.Fatalf("expected confirmed deposit tx-good, got ok=%v txid=%q", ok, txid)
	}

	// Wrong amount → not matched.
	if ok, _, _ := node.DepositConfirmed(wallet, seller, 11.0, nb, na); ok { //nolint
		t.Fatal("expected no match for a different amount")
	}

	// Wrong sender → not matched.
	if ok, _, _ := node.DepositConfirmed(wallet, "2someone...else", 10.0, nb, na); ok { //nolint
		t.Fatal("expected no match for a different sender")
	}

	// Deposit predates the window (anti-replay of an old transaction) → no match.
	if ok, _, _ := node.DepositConfirmed(wallet, seller, 10.0, time.Unix(1500, 0), time.Unix(2000, 0)); ok { //nolint
		t.Fatal("expected no match for a deposit before the listing window")
	}
}

// fibStore routes one fibercoin symbol to a node URL and disables everything
// else, satisfying chain.Store for the fibercoin-payment routing test.
type fibStore struct {
	symbol, nodeURL string
}

func (f fibStore) ExplorerConfig(string) (string, string, string, error) { return "", "", "", nil }
func (f fibStore) SellCoinConfig(sym string) (string, string, string, int, bool, error) {
	if sym == f.symbol {
		return f.nodeURL, "", "", 1, true, nil
	}
	return "", "", "", 0, false, nil
}

// TestFiberPaymentRouting: a payment in an enabled sell coin is verified on that
// coin's node by sender + window + FIXED amount, returning the live confirmation
// count — no unique non-round amount involved.
func TestFiberPaymentRouting(t *testing.T) {
	const seller = "2seller...recv"
	const buyer = "2buyer...addr"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Buyer paid the seller EXACTLY 100 FIB (a fixed, round amount) at time 1000.
		resp := []map[string]any{{
			"status": map[string]any{"confirmed": true, "height": 4, "block_seq": 200},
			"txn": map[string]any{
				"txid":      "pay-tx",
				"timestamp": 1000,
				"inputs":    []map[string]any{{"owner": buyer}},
				"outputs":   []map[string]any{{"dst": seller, "coins": "100.000000"}},
			},
		}}
		_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer srv.Close()

	c := New(fibStore{symbol: "FIB", nodeURL: srv.URL})
	c.sky.hc = srv.Client()
	nb, na := time.Unix(500, 0), time.Unix(2000, 0)

	confs, txid, err := c.PaymentConfirmations("FIB", seller, buyer, 100.0, nb, na)
	if err != nil {
		t.Fatalf("PaymentConfirmations(FIB): %v", err)
	}
	if confs != 4 || txid != "pay-tx" {
		t.Fatalf("fiber payment = (%d,%q), want (4,pay-tx)", confs, txid)
	}
	// Wrong buyer (sender) → no match, even though the fixed amount is identical.
	if confs, _, _ := c.PaymentConfirmations("FIB", seller, "2other...buyer", 100.0, nb, na); confs != 0 { //nolint
		t.Fatalf("expected 0 confs for a different sender, got %d", confs)
	}
	// An unconfigured currency falls through to the (empty) explorer → 0.
	if confs, _, _ := c.PaymentConfirmations("BTC", seller, buyer, 100.0, nb, na); confs != 0 { //nolint
		t.Fatalf("expected 0 confs for a non-fiber unconfigured currency, got %d", confs)
	}
}

// TestDepositBelowConfirmations rejects a deposit that has too few confirmations.
func TestDepositBelowConfirmations(t *testing.T) {
	const wallet = "2sky...market"
	const seller = "2seller...addr"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := []map[string]any{{
			"status": map[string]any{"confirmed": true, "height": 1, "block_seq": 100},
			"txn": map[string]any{
				"txid":      "tx-shallow",
				"timestamp": 1000,
				"inputs":    []map[string]any{{"owner": seller}},
				"outputs":   []map[string]any{{"dst": wallet, "coins": "5.000000"}},
			},
		}}
		_ = json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}))
	defer srv.Close()

	node := NewSkyNode(srv.URL, "", "", 2, srv.Client()) // needs 2, tx has 1
	ok, _, err := node.DepositConfirmed(wallet, seller, 5.0, time.Unix(500, 0), time.Unix(2000, 0))
	if err != nil {
		t.Fatalf("DepositConfirmed: %v", err)
	}
	if ok {
		t.Fatal("expected rejection: only 1 confirmation, need 2")
	}
}

// Note: the SKY spend path (local build + sign + inject) is covered end-to-end in
// sky_signer_test.go. TestSendSKYNoWallet here only asserts the no-seed guard.

// TestSendSKYNoWallet errors clearly when no escrow seed is configured.
func TestSendSKYNoWallet(t *testing.T) {
	node := NewSkyNode("http://127.0.0.1:1", "", "", 2, nil)
	if _, err := node.SendCoin("addr", 1); err == nil {
		t.Fatal("expected an error when no escrow seed is configured")
	}
}

// TestNoExplorer confirms the default explorer never confirms a payment.
func TestNoExplorer(t *testing.T) {
	c := New(nil)
	confs, _, err := c.PaymentConfirmations("BTC", "addr", "", 1.0, time.Unix(0, 0), time.Now())
	if err != nil || confs != 0 {
		t.Fatalf("noExplorer should report 0 confirmations, got confs=%d err=%v", confs, err)
	}
}
