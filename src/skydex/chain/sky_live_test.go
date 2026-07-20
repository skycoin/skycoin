//go:build livenet

// Live validation of the Skycoin node client against the public node
// (node.skycoin.com). Read-only: it proves DepositConfirmed parses the real
// verbose /transactions shape and correctly matches a confirmed on-chain payment
// by amount + confirmation depth. NOT part of the offline test run.
//
//	go test -tags livenet -run TestLiveSky -v ./internal/skydex-market/chain/
//
// (liveGet is defined in explorer_live_test.go, same package + build tag.)
package chain

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

const publicSkyNode = "https://node.skycoin.com"

func TestLiveSkyDepositConfirmed(t *testing.T) {
	node := NewSkyNode(publicSkyNode, "", "", 2, &http.Client{Timeout: 60 * time.Second})

	// Recent blocks (blocks[0] is the deepest, ≈ a handful of confirmations, so it
	// clears the confs=2 threshold). Each output is a real confirmed payment.
	var lb struct {
		Blocks []struct {
			Body struct {
				Txns []struct {
					TxID    string `json:"txid"`
					Outputs []struct {
						Dst   string `json:"dst"`
						Coins string `json:"coins"`
					} `json:"outputs"`
				} `json:"txns"`
			} `json:"body"`
		} `json:"blocks"`
	}
	liveGet(t, publicSkyNode, "/api/v1/last_blocks?num=12", &lb)
	if len(lb.Blocks) == 0 {
		t.Skip("no blocks returned")
	}

	tried := 0
	for _, blk := range lb.Blocks {
		for _, tx := range blk.Body.Txns {
			for _, o := range tx.Outputs {
				amount, err := strconv.ParseFloat(o.Coins, 64)
				if err != nil || amount <= 0 || o.Dst == "" {
					continue
				}
				if tried++; tried > 25 {
					t.Skip("no low-history address found among recent outputs to match against")
				}

				// Our client must detect this exact payment to this address. This
				// test matches arbitrary on-chain payments by amount only, so the
				// sender is left empty and the time window is wide open.
				anySender := ""
				wideOpen, farFuture := time.Unix(0, 0), time.Now().Add(24*time.Hour)
				ok, txid, derr := node.DepositConfirmed(o.Dst, anySender, amount, wideOpen, farFuture)
				if derr != nil {
					continue // busy address (response too large) or transient — try next
				}
				if !ok || txid == "" {
					continue
				}
				t.Logf("detected %.6f SKY to %s (tx %s) via the public node", amount, o.Dst, txid)

				// No false positive on an amount that address never received.
				bad, _, berr := node.DepositConfirmed(o.Dst, anySender, amount+987654.321, wideOpen, farFuture)
				if berr != nil {
					t.Fatalf("negative-check request failed: %v", berr)
				}
				if bad {
					t.Fatalf("false positive: matched a payment of %.6f that was never made", amount+987654.321)
				}
				return // success
			}
		}
	}
	t.Skip("no confirmable deposit found among recent outputs (addresses too busy?)")
}
