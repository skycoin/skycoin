package chain

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Explorer verifies a buyer's payment on an external chain. Concrete adapters
// (esplora, …) implement it; the router below dispatches each currency to the
// adapter the operator configured for it.
type Explorer interface {
	// PaymentConfirmations returns the confirmation count of a payment of
	// expectedAmount received at addr in the given currency (0 if not observed or
	// unconfirmed), plus its transaction hash. A confirmed payment is only counted
	// if its block time is within [notBefore, notAfter] — the order's window —
	// which blocks an old, unrelated payment of a recycled amount from matching.
	// senderAddr, if non-empty, is the buyer's registered payment address; it is
	// used as a soft corroborating signal (a tiebreaker), not a hard requirement,
	// since wallets may spend from a different address than the one on file.
	PaymentConfirmations(currency, addr, senderAddr string, expectedAmount float64, notBefore, notAfter time.Time) (int, string, error)
}

// noExplorer is used when no config store is available; it never confirms.
type noExplorer struct{}

func (noExplorer) PaymentConfirmations(string, string, string, float64, time.Time, time.Time) (int, string, error) {
	return 0, "", nil
}

// ExplorerConfigStore returns the operator's configured explorer for a currency:
// the provider name and its optional base-URL override / API key. An empty
// provider means the currency is disabled.
type ExplorerConfigStore interface {
	ExplorerConfig(currency string) (provider, baseURL, apiKey string, err error)
}

// providerFactory describes an explorer adapter: which currencies it can verify
// and how to build one from the operator's per-coin config. A nil `only` set
// means the adapter is universal (it can verify any coin the operator points it
// at with a URL).
type providerFactory struct {
	only  map[string]bool // nil => universal
	build func(currency, baseURL, apiKey string, hc *http.Client) Explorer
}

// supports reports whether the adapter can verify currency.
func (f providerFactory) supports(currency string) bool {
	return f.only == nil || f.only[currency]
}

// registry of available explorer adapters, keyed by provider name. New providers
// (blockcypher, etherscan, trongrid, 3xpl, …) are added here.
var registry = map[string]providerFactory{
	"esplora": {
		// The Esplora adapter is universal: it can verify ANY UTXO coin exposed by
		// an Esplora-compatible API. BTC/LTC ship with a default public endpoint;
		// every other coin (BCH/DOGE/DASH and any custom coin) just needs the
		// operator to supply an explorer URL.
		only: nil,
		build: func(currency, baseURL, _ string, hc *http.Client) Explorer {
			return newEsplora(currency, baseURL, hc)
		},
	},
}

// DefaultExplorerURL returns the built-in Esplora endpoint for a currency, or ""
// if it has none (the operator must configure a URL to enable it).
func DefaultExplorerURL(currency string) string {
	return esploraDefaultBase[currency]
}

// ProvidersFor returns the explorer provider names that can verify currency, so
// the operator UI can offer a per-coin dropdown. Sorted for stable output.
func ProvidersFor(currency string) []string {
	var out []string
	for name, f := range registry {
		if f.supports(currency) {
			out = append(out, name)
		}
	}
	slices.Sort(out)
	return out
}

// SupportsProvider reports whether provider can verify currency (used to
// validate operator config before saving).
func SupportsProvider(currency, provider string) bool {
	f, ok := registry[provider]
	return ok && f.supports(currency)
}

// router dispatches each currency to its configured adapter, caching built
// adapters. It reads config lazily so operator changes take effect without a
// restart.
type router struct {
	store ExplorerConfigStore
	hc    *http.Client

	mu    sync.Mutex
	cache map[string]Explorer
}

func newRouter(store ExplorerConfigStore, hc *http.Client) *router {
	return &router{store: store, hc: hc, cache: make(map[string]Explorer)}
}

func (r *router) PaymentConfirmations(currency, addr, senderAddr string, amount float64, notBefore, notAfter time.Time) (int, string, error) {
	provider, baseURL, apiKey, err := r.store.ExplorerConfig(currency)
	if err != nil {
		return 0, "", err
	}
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return 0, "", nil // currency not configured => never confirms
	}
	f, ok := registry[provider]
	if !ok || !f.supports(currency) {
		return 0, "", fmt.Errorf("no %q explorer adapter for %s", provider, currency)
	}

	cacheKey := currency + "|" + provider + "|" + baseURL
	r.mu.Lock()
	exp, ok := r.cache[cacheKey]
	if !ok {
		exp = f.build(currency, baseURL, apiKey, r.hc)
		r.cache[cacheKey] = exp
	}
	r.mu.Unlock()

	return exp.PaymentConfirmations(currency, addr, senderAddr, amount, notBefore, notAfter)
}

// --- Esplora adapter (mempool.space / litecoinspace / Blockstream) ---

// esploraDefaultBase is the default API host per coin; overridable via config.
var esploraDefaultBase = map[string]string{
	"BTC": "https://mempool.space",
	"LTC": "https://litecoinspace.org",
}

type esplora struct {
	base string
	hc   *http.Client
}

func newEsplora(currency, baseURL string, hc *http.Client) *esplora {
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if base == "" {
		base = esploraDefaultBase[currency]
	}
	return &esplora{base: base, hc: hc}
}

type esploraTx struct {
	TxID   string `json:"txid"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight uint64 `json:"block_height"`
		BlockTime   int64  `json:"block_time"`
	} `json:"status"`
	Vin []struct {
		Prevout struct {
			Address string `json:"scriptpubkey_address"`
		} `json:"prevout"`
	} `json:"vin"`
	Vout []struct {
		Address string `json:"scriptpubkey_address"`
		Value   uint64 `json:"value"` // satoshis/litoshis
	} `json:"vout"`
}

// hasSender reports whether any input to the tx was funded by addr.
func (tx esploraTx) hasSender(addr string) bool {
	for _, in := range tx.Vin {
		if in.Prevout.Address == addr {
			return true
		}
	}
	return false
}

// PaymentConfirmations scans the address's recent transactions for one that pays
// exactly amount to addr within the order window, returning its confirmation
// count. A confirmed tx must have a block time in [notBefore, notAfter] (this
// blocks an old payment of a recycled amount); unconfirmed txs are necessarily
// recent and bypass the window. When senderAddr is set, a tx that was funded by
// it is preferred over one that was not (a soft corroborating tiebreaker).
func (e *esplora) PaymentConfirmations(_ /*currency*/, addr, senderAddr string, amount float64, notBefore, notAfter time.Time) (int, string, error) {
	var txs []esploraTx
	if err := e.getJSON("/api/address/"+url.PathEscape(addr)+"/txs", &txs); err != nil {
		return 0, "", err
	}
	tip, err := e.tipHeight()
	if err != nil {
		return 0, "", err
	}

	const eps = 0.5 / 1e8 // half a base unit (8 decimals)
	nb, na := notBefore.Unix(), notAfter.Unix()
	best := -1
	var bestTx string
	bestSender := false
	for _, tx := range txs {
		var received uint64
		for _, o := range tx.Vout {
			if o.Address == addr {
				received += o.Value
			}
		}
		if received == 0 {
			continue
		}
		if math.Abs(float64(received)/1e8-amount) > eps {
			continue
		}
		// A confirmed payment must fall inside the order window; unconfirmed ones
		// are in the mempool now, so they can't be an old replay and are accepted.
		if tx.Status.Confirmed && (tx.Status.BlockTime < nb || tx.Status.BlockTime > na) {
			continue
		}
		confs := 0
		if tx.Status.Confirmed && tip >= tx.Status.BlockHeight {
			confs = int(tip-tx.Status.BlockHeight) + 1 //nolint
		}
		senderMatch := senderAddr != "" && tx.hasSender(senderAddr)
		// Prefer a sender-corroborated tx; otherwise the most-confirmed one.
		if best < 0 || (senderMatch != bestSender && senderMatch) || (senderMatch == bestSender && confs > best) {
			best, bestTx, bestSender = confs, tx.TxID, senderMatch
		}
	}
	if best < 0 {
		return 0, "", nil // payment not seen yet
	}
	return best, bestTx, nil
}

func (e *esplora) tipHeight() (uint64, error) {
	body, err := e.get("/api/blocks/tip/height")
	if err != nil {
		return 0, err
	}
	h, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("esplora tip height: %w", err)
	}
	return h, nil
}

func (e *esplora) getJSON(path string, out any) error {
	body, err := e.get(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (e *esplora) get(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, e.base+path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := e.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("esplora %s: %w", path, err)
	}
	defer resp.Body.Close()                                 //nolint:errcheck
	data, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<20)) //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("esplora %s: status %d", path, resp.StatusCode)
	}
	return data, nil
}
