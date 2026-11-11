package explorerui

import (
	"errors"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

// Node is the part of the node's REST API the pages read. Client provides it
// over HTTP. Every call is a GET.
type Node interface {
	Health() (*api.HealthResponse, error)
	BlockchainProgress() (*readable.BlockchainProgress, error)
	BlockchainMetadata() (*readable.BlockchainMetadata, error)
	CoinSupply() (*api.CoinSupply, error)
	BlocksInRangeVerbose(start, end uint64) (*readable.BlocksVerbose, error)
	BlockBySeqVerbose(seq uint64) (*readable.BlockVerbose, error)
	BlockByHashVerbose(hash string) (*readable.BlockVerbose, error)
	TransactionVerbose(txid string) (*readable.TransactionWithStatusVerbose, error)
}

// Client is a Node reached through the REST API client.
type Client struct {
	*api.Client
}

// NewClient reads the node at addr, such as http://127.0.0.1:6420.
func NewClient(addr string) Client {
	return Client{api.NewClient(addr)}
}

// notFound reports whether err is the node saying there is no such thing, or
// that the request named something invalid.
func notFound(err error) bool {
	var ce api.ClientError
	return errors.As(err, &ce) && ce.StatusCode >= 400 && ce.StatusCode < 500
}

// nodeInfo keeps the node's coin names and sync state for a few seconds, so
// a page does not ask for them on every request.
type nodeInfo struct {
	mu      sync.Mutex
	at      time.Time
	coin    Coin
	syncing bool
}

const nodeInfoTTL = 15 * time.Second

func (ni *nodeInfo) get(n Node) (Coin, bool) {
	ni.mu.Lock()
	defer ni.mu.Unlock()
	if !ni.at.IsZero() && time.Since(ni.at) < nodeInfoTTL {
		return ni.coin, ni.syncing
	}
	c := defaultCoin
	h, err := n.Health()
	if err != nil {
		return c, false
	}
	f := h.Fiber
	if f.DisplayName != "" {
		c.FullName = f.DisplayName
	}
	if f.Ticker != "" {
		c.Ticker = f.Ticker
	}
	if f.CoinHoursName != "" {
		c.Hours = f.CoinHoursName
	}
	if f.CoinHoursNameSingular != "" {
		c.HoursSingular = f.CoinHoursNameSingular
	}
	if f.QrURIPrefix != "" {
		c.QRPrefix = f.QrURIPrefix
	}
	if d := h.UserVerifyTxn.MaxDropletPrecision; d > 0 {
		c.MaxDecimals = int(d)
	}
	syncing := false
	if p, err := n.BlockchainProgress(); err == nil {
		syncing = p.Current < p.Highest
	}
	ni.coin, ni.syncing, ni.at = c, syncing, time.Now()
	return c, syncing
}
