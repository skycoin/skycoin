package chain

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"
)

// fakeStore is a chain.Store returning a fixed explorer config for one currency.
// SellCoinConfig is unused by these explorer tests (returns disabled).
type fakeStore struct {
	currency, provider, url, key string
}

func (f fakeStore) ExplorerConfig(currency string) (string, string, string, error) {
	if currency == f.currency {
		return f.provider, f.url, f.key, nil
	}
	return "", "", "", nil
}

func (f fakeStore) SellCoinConfig(string) (string, string, string, int, bool, error) {
	return "", "", "", 0, false, nil
}

func TestProvidersFor(t *testing.T) {
	// The esplora adapter is universal: it covers the built-in coins and any
	// operator-added custom coin (the operator supplies the explorer URL).
	for _, c := range []string{"BTC", "LTC", "BCH", "DOGE", "DASH", "GRS", "ZEC"} {
		if got := ProvidersFor(c); !slices.Contains(got, "esplora") {
			t.Fatalf("ProvidersFor(%s) = %v, want esplora present", c, got)
		}
		if !SupportsProvider(c, "esplora") {
			t.Fatalf("SupportsProvider(%s, esplora) = false, want true", c)
		}
	}
	// An unknown provider name is still rejected.
	if SupportsProvider("BTC", "bogus") {
		t.Fatal("unknown provider must not be supported")
	}
}

// TestEsploraPaymentConfirmations serves an Esplora-shaped address txs response
// and a tip height, then checks the router matches the exact amount + confs.
func TestEsploraPaymentConfirmations(t *testing.T) {
	const addr = "bc1qbuyerpays"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/blocks/tip/height":
			_, _ = w.Write([]byte("800010")) //nolint:errcheck
		case "/api/address/" + addr + "/txs":
			txs := []map[string]any{
				{ // pays exactly 0.05 to addr from the buyer, confirmed at height 800000 (11 confs), block time 1000
					"txid":   "tx-pay",
					"status": map[string]any{"confirmed": true, "block_height": 800000, "block_time": 1000},
					"vin":    []map[string]any{{"prevout": map[string]any{"scriptpubkey_address": "bc1qbuyer"}}},
					"vout": []map[string]any{
						{"scriptpubkey_address": addr, "value": 5000000},
						{"scriptpubkey_address": "other", "value": 999},
					},
				},
			}
			_ = json.NewEncoder(w).Encode(txs) //nolint:errcheck
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	// Route BTC to an esplora adapter pointed at the mock server.
	c := New(fakeStore{currency: "BTC", provider: "esplora", url: srv.URL})
	c.exp.(*router).hc = srv.Client()

	// Window [500, 2000] covers the payment's block time (1000); buyer corroborates.
	nb, na := time.Unix(500, 0), time.Unix(2000, 0)
	confs, txid, err := c.PaymentConfirmations("BTC", addr, "bc1qbuyer", 0.05, nb, na)
	if err != nil {
		t.Fatalf("PaymentConfirmations: %v", err)
	}
	if confs != 11 || txid != "tx-pay" {
		t.Fatalf("got confs=%d txid=%q, want 11 / tx-pay", confs, txid)
	}

	// A different amount is not matched.
	if confs, _, err := c.PaymentConfirmations("BTC", addr, "bc1qbuyer", 0.04, nb, na); err != nil || confs != 0 {
		t.Fatalf("expected no match for a different amount, got confs=%d err=%v", confs, err)
	}

	// A payment confirmed before the window (anti-replay) is not matched.
	if confs, _, err := c.PaymentConfirmations("BTC", addr, "bc1qbuyer", 0.05, time.Unix(1500, 0), time.Unix(2000, 0)); err != nil || confs != 0 {
		t.Fatalf("expected no match for a payment before the window, got confs=%d err=%v", confs, err)
	}
}

// TestRouterUnconfiguredCurrency returns 0 confirmations (not an error) when the
// currency has no explorer configured.
func TestRouterUnconfiguredCurrency(t *testing.T) {
	c := New(fakeStore{}) // store returns empty provider for everything
	confs, _, err := c.PaymentConfirmations("BTC", "addr", "", 1.0, time.Unix(0, 0), time.Now())
	if err != nil || confs != 0 {
		t.Fatalf("unconfigured currency should be 0/nil, got confs=%d err=%v", confs, err)
	}
}
