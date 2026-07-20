package chain

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// SellCoinConfigStore returns the operator's escrow-node config for a sell coin
// (SKY or a fibercoin): its fullnode URL, escrow hot-wallet seed + address, the
// required confirmation depth, and whether it is enabled. A coin that is unknown
// or disabled is not spendable/verifiable. The DB implements this.
type SellCoinConfigStore interface {
	SellCoinConfig(symbol string) (nodeURL, seed, wallet string, confs int, enabled bool, err error)
}

// Store is the full backing config the Chain reads: per-sell-coin escrow config
// and per-payment-currency explorer config. The market DB satisfies both.
type Store interface {
	SellCoinConfigStore
	ExplorerConfigStore
}

// Chain composes a per-sell-coin Skycoin-node router with a per-currency explorer
// router to satisfy jobs.Chain: sell-coin deposits/spends go to that coin's node;
// external-coin payment verification is dispatched to whichever explorer the
// operator configured.
type Chain struct {
	sky *skyRouter
	exp Explorer
}

// New builds a Chain from store. store provides both the per-sell-coin escrow
// config and the per-currency explorer config; a nil store disables everything.
func New(store Store) *Chain {
	hc := &http.Client{Timeout: 20 * time.Second}
	var exp Explorer = noExplorer{}
	var sky *skyRouter
	if store != nil {
		exp = newRouter(store, hc)
		sky = newSkyRouter(store, hc)
	}
	return &Chain{sky: sky, exp: exp}
}

// DepositConfirmed implements jobs.Chain for the given sell coin.
func (c *Chain) DepositConfirmed(coin, marketWallet, senderAddr string, amount float64, notBefore, notAfter time.Time) (bool, string, error) {
	n, err := c.node(coin)
	if err != nil {
		return false, "", err
	}
	return n.DepositConfirmed(marketWallet, senderAddr, amount, notBefore, notAfter)
}

// PaymentConfirmations implements jobs.Chain. When currency is an enabled sell
// coin (a fibercoin, incl. SKY), the payment is a peer-to-peer transfer on that
// coin's own chain: it is verified via the coin's fullnode by sender + window +
// FIXED amount (addr is the seller's receive address), so no unique non-round
// amount is needed. Otherwise it is an external UTXO coin verified via its
// explorer (unique amount + window + soft sender).
func (c *Chain) PaymentConfirmations(currency, addr, senderAddr string, expectedAmount float64, notBefore, notAfter time.Time) (int, string, error) {
	if c.sky != nil && c.sky.enabled(currency) {
		n, err := c.sky.node(currency)
		if err != nil {
			return 0, "", err
		}
		return n.DepositConfirmations(addr, senderAddr, expectedAmount, notBefore, notAfter)
	}
	return c.exp.PaymentConfirmations(currency, addr, senderAddr, expectedAmount, notBefore, notAfter)
}

// SendCoin implements jobs.Chain: spend the given sell coin from its escrow
// wallet to addr.
func (c *Chain) SendCoin(coin, toAddr string, amount float64) (string, error) {
	n, err := c.node(coin)
	if err != nil {
		return "", err
	}
	return n.SendCoin(toAddr, amount)
}

// EscrowBalance implements jobs.Chain: the live spendable balance of the sell
// coin's escrow wallet, for the escrow-audit consistency check.
func (c *Chain) EscrowBalance(coin, addr string) (float64, error) {
	n, err := c.node(coin)
	if err != nil {
		return 0, err
	}
	b, err := n.Balance(addr)
	if err != nil {
		return 0, err
	}
	return b.Coins, nil
}

func (c *Chain) node(coin string) (*SkyNode, error) {
	if c.sky == nil {
		return nil, fmt.Errorf("no sell-coin backend configured")
	}
	return c.sky.node(coin)
}

// skyRouter resolves and caches one SkyNode per sell coin from the config store.
// It reads config lazily so operator changes (new node URL, seed, confirmations)
// take effect without a restart — a changed config yields a new cache key.
type skyRouter struct {
	store SellCoinConfigStore
	hc    *http.Client

	mu    sync.Mutex
	cache map[string]*SkyNode
}

func newSkyRouter(store SellCoinConfigStore, hc *http.Client) *skyRouter {
	return &skyRouter{store: store, hc: hc, cache: make(map[string]*SkyNode)}
}

// enabled reports whether symbol is a configured, enabled sell coin — i.e. a
// fibercoin the market can verify on its own node (used to decide whether a
// payment currency is a fibercoin or an external explorer coin).
func (r *skyRouter) enabled(symbol string) bool {
	_, _, _, _, enabled, err := r.store.SellCoinConfig(symbol)
	return err == nil && enabled
}

func (r *skyRouter) node(symbol string) (*SkyNode, error) {
	nodeURL, seed, wallet, confs, enabled, err := r.store.SellCoinConfig(symbol)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, fmt.Errorf("sell coin %q is not available at this market", symbol)
	}
	if nodeURL == "" {
		return nil, fmt.Errorf("sell coin %q has no fullnode URL configured", symbol)
	}
	// Key on the full config so any change (URL, seed, wallet, confs) rebuilds the
	// node and its local signer.
	key := symbol + "|" + nodeURL + "|" + wallet + "|" + strconv.Itoa(confs) + "|" + seed

	r.mu.Lock()
	defer r.mu.Unlock()
	if n, ok := r.cache[key]; ok {
		return n, nil
	}
	n := NewSkyNode(nodeURL, seed, wallet, confs, r.hc)
	r.cache[key] = n
	return n, nil
}
