// Package explorerui src/explorerui/node.go c3-ui-explorer
package explorerui

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	PendingTransactionsVerbose() ([]readable.UnconfirmedTransactionVerbose, error)
	Richlist(*api.RichlistParams) (*api.Richlist, error)
	AddressHasSent(addr string) (bool, error)
	AddressBalance(addr string) (*api.BalanceResponse, error)
	AddressOutputs(addr string) (*readable.UnspentOutputsSummary, error)
	AddressTransactions(addr string) ([]readable.TransactionWithStatusVerbose, error)
	AddressTransactionsPage(addr string, page, limit int) (*api.TransactionsWithStatusVerboseV2, error)
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
	if p := f.QrURIPrefix; p != "" {
		c.QRPrefix = strings.TrimSuffix(p, ":") + ":"
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

func query(path string, kv ...string) string {
	v := url.Values{}
	for i := 0; i+1 < len(kv); i += 2 {
		v.Set(kv[i], kv[i+1])
	}
	return path + "?" + v.Encode()
}

// The query helpers below keep every request a GET, as the explorer's JSON
// proxy does, where the REST client would POST.

// AddressBalance is GET /api/v1/balance for one address.
func (c Client) AddressBalance(addr string) (*api.BalanceResponse, error) {
	var b api.BalanceResponse
	if err := c.Get(query("/api/v1/balance", "addrs", addr), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// AddressOutputs is GET /api/v1/outputs for one address.
func (c Client) AddressOutputs(addr string) (*readable.UnspentOutputsSummary, error) {
	var o readable.UnspentOutputsSummary
	if err := c.Get(query("/api/v1/outputs", "addrs", addr), &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// AddressTransactions is GET /api/v1/transactions, every transaction of one
// address, confirmed or not.
func (c Client) AddressTransactions(addr string) ([]readable.TransactionWithStatusVerbose, error) {
	var t []readable.TransactionWithStatusVerbose
	if err := c.Get(query("/api/v1/transactions", "addrs", addr, "verbose", "1"), &t); err != nil {
		return nil, err
	}
	return t, nil
}

// AddressTransactionsPage is GET /api/v2/transactions, one page of an
// address's transactions, newest first.
func (c Client) AddressTransactionsPage(addr string, page, limit int) (*api.TransactionsWithStatusVerboseV2, error) {
	return c.TransactionsVerboseV2(
		api.RequestArg{Key: "addrs", Value: addr},
		api.RequestArg{Key: "page", Value: strconv.Itoa(page)},
		api.RequestArg{Key: "limit", Value: strconv.Itoa(limit)},
		api.RequestArg{Key: "sort", Value: "desc"},
	)
}

// AddressHasSent reports whether addr has spent an output, reading GET
// /api/v1/address_uxouts only up to the first spent one. The whole list of
// a busy address runs to tens of megabytes.
func (c Client) AddressHasSent(addr string) (bool, error) {
	resp, err := c.HTTPClient.Get(c.Addr + strings.TrimPrefix(query("/api/v1/address_uxouts", "address", addr), "/"))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096)) //nolint:errcheck
		return false, api.NewClientError(resp.Status, resp.StatusCode, string(body))
	}
	d := json.NewDecoder(resp.Body)
	if tok, err := d.Token(); err != nil || tok == nil {
		return false, err
	}
	for d.More() {
		var o struct {
			SpentBlockSeq uint64 `json:"spent_block_seq"`
		}
		if err := d.Decode(&o); err != nil {
			return false, err
		}
		if o.SpentBlockSeq != 0 {
			return true, nil
		}
	}
	return false, nil
}
