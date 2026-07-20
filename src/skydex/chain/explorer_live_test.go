//go:build livenet

// Live validation of the Esplora adapter against the real public APIs
// (mempool.space for BTC, litecoinspace.org for LTC). These are NOT part of the
// normal offline test run — they need outbound network and hit third-party
// services. Run explicitly:
//
//	go test -tags livenet -run TestLiveEsplora -v ./internal/skydex-market/chain/
//
// The test derives a real, currently-confirmed payment from each chain and
// asserts the adapter (a) parses the live response shape and (b) matches the
// exact amount with a positive confirmation count — the two things the
// documented-shape unit test can't prove.
package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// liveGetText fetches a plain-text endpoint (e.g. /api/block-height/{h}, which
// returns a bare block hash — not JSON).
func liveGetText(t *testing.T, base, path string) string {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, base+path, nil) //nolint:errcheck
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s%s: %v", base, path, err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s%s: status %d", base, path, resp.StatusCode)
	}
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) //nolint:errcheck
	return strings.TrimSpace(string(b))
}

func liveGet(t *testing.T, base, path string, out any) {
	t.Helper()
	req, _ := http.NewRequest(http.MethodGet, base+path, nil) //nolint:errcheck
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s%s: %v", base, path, err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET %s%s: status %d", base, path, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("decode %s%s: %v", base, path, err)
	}
}

func TestLiveEsplora(t *testing.T) {
	for _, tc := range []struct{ currency, base string }{
		{"BTC", "https://mempool.space"},
		{"LTC", "https://litecoinspace.org"},
	} {
		t.Run(tc.currency, func(t *testing.T) {
			hc := &http.Client{Timeout: 30 * time.Second}
			e := newEsplora(tc.currency, "", hc) // "" => default public host

			// Derive a live address with confirmed activity: take a block a few
			// deep (so its txs are confirmed) and pull a paid-to address from a
			// non-coinbase tx.
			var tipH uint64
			liveGet(t, tc.base, "/api/blocks/tip/height", &tipH)
			blockHash := liveGetText(t, tc.base, fmt.Sprintf("/api/block-height/%d", tipH-6))
			var txids []string
			liveGet(t, tc.base, "/api/block/"+blockHash+"/txids", &txids)
			if len(txids) < 2 {
				t.Skip("block has no non-coinbase tx to sample")
			}
			var sampleTx esploraTx
			liveGet(t, tc.base, "/api/tx/"+txids[1], &sampleTx)
			var addr string
			for _, o := range sampleTx.Vout {
				if o.Address != "" {
					addr = o.Address
					break
				}
			}
			if addr == "" {
				t.Skip("sampled tx has no address-bearing output")
			}

			// Fetch that address's txs the same way the adapter will, and pick a
			// confirmed tx to reproduce: amount = sum of its outputs to addr
			// (exactly what the adapter sums), so the adapter must match it.
			var txs []esploraTx
			liveGet(t, tc.base, "/api/address/"+addr+"/txs", &txs)
			var wantAmount float64
			var wantTxSeen bool
			for _, tx := range txs {
				if !tx.Status.Confirmed {
					continue
				}
				var sats uint64
				for _, o := range tx.Vout {
					if o.Address == addr {
						sats += o.Value
					}
				}
				if sats > 0 {
					wantAmount = float64(sats) / 1e8
					wantTxSeen = true
					break
				}
			}
			if !wantTxSeen {
				t.Skipf("no confirmed payment to %s in its recent txs", addr)
			}

			// (a) shape + (b) positive match: the adapter must confirm this payment.
			// This test matches an arbitrary real payment by amount only, so the
			// sender is left empty and the window is wide open.
			anySender := ""
			wideOpen, farFuture := time.Unix(0, 0), time.Now().Add(24*time.Hour)
			confs, txid, err := e.PaymentConfirmations(tc.currency, addr, anySender, wantAmount, wideOpen, farFuture)
			if err != nil {
				t.Fatalf("PaymentConfirmations(%s, %s, %.8f): %v", tc.currency, addr, wantAmount, err)
			}
			if confs <= 0 || txid == "" {
				t.Fatalf("live payment of %.8f to %s: got confs=%d txid=%q, want a positive confirmed match",
					wantAmount, addr, confs, txid)
			}
			t.Logf("%s: matched %.8f to %s at %d confs (tx %s)", tc.currency, wantAmount, addr, confs, txid)

			// A clearly-impossible amount must not match (no false positive, and
			// proves the full parse path ran without error).
			confs, _, err = e.PaymentConfirmations(tc.currency, addr, anySender, 987654.32109876, wideOpen, farFuture)
			if err != nil {
				t.Fatalf("PaymentConfirmations(no-match): %v", err)
			}
			if confs != 0 {
				t.Fatalf("impossible amount matched at %d confs — false positive", confs)
			}
		})
	}
}
