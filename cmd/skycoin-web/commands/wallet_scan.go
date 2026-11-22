// Package commands cmd/skycoin-web/commands/wallet_scan.go c4-ui-wallet
package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

// scanGap is how many unused addresses in a row end a scan, the BIP44 gap limit.
const scanGap = 20

// scanEvery is how often the bip44 and xpub wallets kept here are scanned.
const scanEvery = time.Minute

// nodeFinder reports which addresses have received coins, from a coin node.
type nodeFinder struct {
	nodeURL string
}

// AddressesActivity implements wallet.TransactionsFinder.
func (f nodeFinder) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	if len(addrs) == 0 {
		return nil, nil
	}
	strs := make([]string, len(addrs))
	for i, a := range addrs {
		strs[i] = a.String()
	}
	q := url.Values{"addrs": {strings.Join(strs, ",")}, "verbose": {"0"}}
	resp, err := (&http.Client{Transport: nodeTransport, Timeout: nodeTimeout}).Get(f.nodeURL + "/api/v1/transactions?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node %s returned %d: %s", f.nodeURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var txns []readable.TransactionWithStatus
	if err := json.Unmarshal(body, &txns); err != nil {
		return nil, fmt.Errorf("decode transactions: %w", err)
	}
	received := make(map[string]bool)
	for _, t := range txns {
		for _, o := range t.Transaction.Out {
			received[o.Address] = true
		}
	}
	active := make([]bool, len(strs))
	for i, s := range strs {
		active[i] = received[s]
	}
	return active, nil
}

// POST /api/v1/wallet/scan
func handleScanAddresses(c *webCtx, s *wallet.Service, nodeURL string) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}
	n := uint64(scanGap)
	if num := c.Request.FormValue("num"); num != "" {
		v, err := strconv.ParseUint(num, 10, 64)
		if err != nil || v == 0 || v > maxScanNum {
			errBadRequest(c, fmt.Sprintf("num must be between 1 and %d", maxScanNum))
			return
		}
		n = v
	}
	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	addrs, err := s.ScanAddresses(wltID, []byte(password), n, nodeFinder{nodeURL})
	if err != nil {
		handleWalletError(c, err)
		return
	}
	rlt := struct {
		Addresses []string `json:"addresses"`
	}{Addresses: []string{}}
	for _, a := range addrs {
		rlt.Addresses = append(rlt.Addresses, a.String())
	}
	c.JSON(http.StatusOK, rlt)
}

// maxScanNum bounds one scan request, as the node's /api/v1/wallet/scan does.
const maxScanNum = 65536

// scanAhead scans the bip44 and xpub wallets of every skycoin type coin past
// their last address at start and then every scanEvery, until ctx is done.
func scanAhead(ctx context.Context, table *coinTable) {
	t := time.NewTicker(scanEvery)
	defer t.Stop()
	for {
		type job struct {
			nodeURL  string
			services []*wallet.Service
		}
		var jobs []job
		coinsMu.RLock()
		for i, c := range table.coins {
			if c.CoinType == "skycoin" && !c.hidden && c.remoteNodeURL != "" {
				jobs = append(jobs, job{c.remoteNodeURL, table.wlt[i]})
			}
		}
		coinsMu.RUnlock()
		for _, j := range jobs {
			for _, svc := range j.services {
				if err := svc.ScanAhead(scanGap, nodeFinder{j.nodeURL}); err != nil {
					log.Printf("[WALLET] scan ahead with %s: %v", j.nodeURL, err) //nolint:gosec // the configured node URL
				}
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}
