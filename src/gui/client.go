package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second
)

// APIError is used for non-200 API responses
type APIError struct {
	Status     string
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return e.Message
}

// Client provides an interface to a remote node's HTTP API
type Client struct {
	HTTPClient *http.Client
	Addr       string
}

// NewClient creates a Client
func NewClient(addr string) *Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dialTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}
	httpClient := &http.Client{
		Transport: transport,
		Timeout:   httpClientTimeout,
	}
	addr = strings.TrimRight(addr, "/")
	addr += "/"

	return &Client{
		Addr:       addr,
		HTTPClient: httpClient,
	}
}

// Get makes a GET request to an endpoint and unmarshals the response to obj.
// If the response is not 200 OK, returns an error
func (c *Client) Get(endpoint string, obj interface{}) error {
	resp, err := c.get(endpoint)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return APIError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if obj == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(obj)
}

// get makes a GET request to an endpoint. Caller must close response body.
func (c *Client) get(endpoint string) (*http.Response, error) {
	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return c.HTTPClient.Do(req)
}

// PostForm makes a POST request to an endpoint with body of "application/x-www-form-urlencoded" formated data.
func (c *Client) PostForm(endpoints string, body io.Reader, obj interface{}) error {
	return c.post(endpoints, "application/x-www-form-urlencoded", body, obj)
}

// PostJSON makes a POST request to an endpoint with body of json data.
func (c *Client) PostJSON(endpoints string, body io.Reader, obj interface{}) error {
	return c.post(endpoints, "application/json", body, obj)
}

// Post makes a POST request to an endpoint. Caller must close response body.
func (c *Client) post(endpoint string, contentType string, body io.Reader, obj interface{}) error {
	csrf, err := c.CSRF()
	if err != nil {
		return err
	}

	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("X-CSRF-Token", csrf)
	req.Header.Set("Content-Type", contentType)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return APIError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if obj == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(obj)
}

// CSRF returns a CSRF token. If CSRF is disabled on the node, returns an empty string and nil error.
func (c *Client) CSRF() (string, error) {
	resp, err := c.get("/csrf")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		// CSRF is disabled on the node
		return "", nil
	default:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return "", APIError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	var m map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return "", err
	}

	token, ok := m["csrf_token"]
	if !ok {
		return "", errors.New("csrf_token not found in response")
	}

	return token, nil
}

// Version makes a request to /version
func (c *Client) Version() (*visor.BuildInfo, error) {
	var bi visor.BuildInfo
	if err := c.Get("/version", &bi); err != nil {
		return nil, err
	}
	return &bi, nil
}

// Outputs makes a request to /outputs
func (c *Client) Outputs() (*visor.ReadableOutputSet, error) {
	var o visor.ReadableOutputSet
	if err := c.Get("/outputs", &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForAddresses makes a request to /outputs?addrs=xxx
func (c *Client) OutputsForAddresses(addrs []string) (*visor.ReadableOutputSet, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/outputs?" + v.Encode()

	var o visor.ReadableOutputSet
	if err := c.Get(endpoint, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForHashes makes a request to /outputs?hashes=zzz
func (c *Client) OutputsForHashes(hashes []string) (*visor.ReadableOutputSet, error) {
	v := url.Values{}
	v.Add("hashes", strings.Join(hashes, ","))
	endpoint := "/outputs?" + v.Encode()

	var o visor.ReadableOutputSet
	if err := c.Get(endpoint, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// CoinSupply makes a request to /coinSupply
func (c *Client) CoinSupply() (*CoinSupply, error) {
	var cs CoinSupply
	if err := c.Get("/coinSupply", &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

// BlockByHash makes a request to /block?hash=xxx
func (c *Client) BlockByHash(hash string) (*visor.ReadableBlock, error) {
	v := url.Values{}
	v.Add("hash", hash)
	endpoint := "/block?" + v.Encode()

	var b visor.ReadableBlock
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeq makes a request to /block?seq=xxx
func (c *Client) BlockBySeq(seq uint64) (*visor.ReadableBlock, error) {
	v := url.Values{}
	v.Add("seq", fmt.Sprint(seq))
	endpoint := "/block?" + v.Encode()

	var b visor.ReadableBlock
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Blocks makes a request to /blocks
func (c *Client) Blocks(start, end int) (*visor.ReadableBlocks, error) {
	v := url.Values{}
	v.Add("start", fmt.Sprint(start))
	v.Add("end", fmt.Sprint(end))
	endpoint := "/blocks?" + v.Encode()

	var b visor.ReadableBlocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// LastBlocks makes a request to /last_blocks
func (c *Client) LastBlocks(n int) (*visor.ReadableBlocks, error) {
	v := url.Values{}
	v.Add("num", fmt.Sprint(n))
	endpoint := "/last_blocks?" + v.Encode()

	var b visor.ReadableBlocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainMetadata makes a request to /blockchain/metadata
func (c *Client) BlockchainMetadata() (*visor.BlockchainMetadata, error) {
	var b visor.BlockchainMetadata
	if err := c.Get("/blockchain/metadata", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainProgress makes a request to /blockchain/progress
func (c *Client) BlockchainProgress() (*daemon.BlockchainProgress, error) {
	var b daemon.BlockchainProgress
	if err := c.Get("/blockchain/progress", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Balance makes a request to /balance?addrs=xxx
func (c *Client) Balance(addrs []string) (*wallet.BalancePair, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/balance?" + v.Encode()

	var b wallet.BalancePair
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// UxOut makes a request to /uxout?uxid=xxx
func (c *Client) UxOut(uxID string) (*historydb.UxOutJSON, error) {
	v := url.Values{}
	v.Add("uxid", uxID)
	endpoint := "/uxout?" + v.Encode()

	var b historydb.UxOutJSON
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// AddressUxOuts makes a request to /address_uxouts
func (c *Client) AddressUxOuts(addr string) ([]*historydb.UxOutJSON, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/address_uxouts?" + v.Encode()

	var b []*historydb.UxOutJSON
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// Wallet makes a request to /wallet
func (c *Client) Wallet(id string) (*wallet.Wallet, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/wallet?" + v.Encode()

	var w wallet.Wallet
	if err := c.Get(endpoint, &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// Wallets makes a request to /wallets
func (c *Client) Wallets() ([]*wallet.ReadableWallet, error) {
	var w []*wallet.ReadableWallet
	if err := c.Get("/wallets", &w); err != nil {
		return nil, err
	}
	return w, nil
}

// CreateWallet makes a request to /wallet/create
// If scanN is <= 0, the scan number defaults to 1
func (c *Client) CreateWallet(seed, label string, scanN int) (*wallet.ReadableWallet, error) {
	v := url.Values{}
	v.Add("seed", seed)
	v.Add("label", label)
	if scanN > 0 {
		v.Add("scan", fmt.Sprint(scanN))
	}

	var w wallet.ReadableWallet
	if err := c.PostForm("/wallet/create", strings.NewReader(v.Encode()), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// NewWalletAddress makes a request to /wallet/newAddress
// if n is <= 0, defaults to 1
func (c *Client) NewWalletAddress(id string, n int) ([]string, error) {
	v := url.Values{}
	v.Add("id", id)
	if n > 0 {
		v.Add("num", fmt.Sprint(n))
	}

	var obj struct {
		Addresses []string `json:"addresses"`
	}
	if err := c.PostForm("/wallet/newAddress", strings.NewReader(v.Encode()), &obj); err != nil {
		return nil, err
	}
	return obj.Addresses, nil
}

// WalletBalance makes a request to /wallet/balance
func (c *Client) WalletBalance(id string) (*wallet.BalancePair, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/wallet/balance?" + v.Encode()

	var b wallet.BalancePair
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Spend makes a request to /wallet/spend
func (c *Client) Spend(id, dst string, coins uint64) (*SpendResult, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("dst", dst)
	v.Add("coins", fmt.Sprint(coins))

	var r SpendResult
	endpoint := "/wallet/spend"
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// WalletTransactions makes a request to /wallet/transactions
func (c *Client) WalletTransactions(id string) (*UnconfirmedTxnsResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/wallet/transactions?" + v.Encode()

	var utx *UnconfirmedTxnsResponse
	if err := c.Get(endpoint, &utx); err != nil {
		return nil, err
	}
	return utx, nil
}

// UpdateWallet makes a request to /wallet/update
func (c *Client) UpdateWallet(id, label string) error {
	v := url.Values{}
	v.Add("id", id)
	v.Add("label", label)

	return c.PostForm("/wallet/update", strings.NewReader(v.Encode()), nil)
}

// WalletFolderName makes a request to /wallets/folderName
func (c *Client) WalletFolderName() (*WalletFolder, error) {
	var w WalletFolder
	if err := c.Get("/wallets/folderName", &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// NewSeed makes a request to /wallet/newSeed
// entropy must be 128 or 256
func (c *Client) NewSeed(entropy int) (string, error) {
	v := url.Values{}
	v.Add("entropy", fmt.Sprint(entropy))
	endpoint := "/wallet/newSeed?" + v.Encode()

	var r struct {
		Seed string `json:"seed"`
	}
	if err := c.Get(endpoint, &r); err != nil {
		return "", err
	}
	return r.Seed, nil
}

// NetworkConnection makes a request to /network/connection
func (c *Client) NetworkConnection(addr string) (*daemon.Connection, error) {
	v := url.Values{}
	v.Add("addr", addr)
	endpoint := "/network/connection?" + v.Encode()

	var dc daemon.Connection
	if err := c.Get(endpoint, &dc); err != nil {
		return nil, err
	}
	return &dc, nil
}

// NetworkConnections makes a request to /network/connections
func (c *Client) NetworkConnections() (*daemon.Connections, error) {
	var dc daemon.Connections
	if err := c.Get("/network/connections", &dc); err != nil {
		return nil, err
	}
	return &dc, nil
}

// NetworkDefaultConnections makes a request to /network/defaultConnections
func (c *Client) NetworkDefaultConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/network/defaultConnections", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkTrustedConnections makes a request to /network/connections/trust
func (c *Client) NetworkTrustedConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/network/connections/trust", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkExchangeableConnections makes a request to /network/connections/exchange
func (c *Client) NetworkExchangeableConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/network/connections/exchange", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// PendingTransactions makes a request to /pendingTxs
func (c *Client) PendingTransactions() ([]*visor.ReadableUnconfirmedTxn, error) {
	var v []*visor.ReadableUnconfirmedTxn
	if err := c.Get("/pendingTxs", &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Transaction makes a request to /transaction
func (c *Client) Transaction(txid string) (*visor.TransactionResult, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/transaction?" + v.Encode()

	var r visor.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Transactions makes a request to /transactions
func (c *Client) Transactions(addrs []string) (*[]visor.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/transactions?" + v.Encode()

	var r []visor.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ConfirmedTransactions makes a request to /transactions?confirmed=true
func (c *Client) ConfirmedTransactions(addrs []string) (*[]visor.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	endpoint := "/transactions?" + v.Encode()

	var r []visor.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UnconfirmedTransactions makes a request to /transactions?confirmed=false
func (c *Client) UnconfirmedTransactions(addrs []string) (*[]visor.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	endpoint := "/transactions?" + v.Encode()

	var r []visor.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// InjectTransaction makes a request to /injectTransaction
func (c *Client) InjectTransaction(rawTx string) (string, error) {
	v := struct {
		Rawtx string `json:"rawtx"`
	}{
		Rawtx: rawTx,
	}

	d, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	var txid string
	if err := c.PostJSON("/injectTransaction", bytes.NewReader(d), &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// ResendUnconfirmedTransactions makes a request to /resendUnconfirmedTxns
func (c *Client) ResendUnconfirmedTransactions() (*daemon.ResendResult, error) {
	var r daemon.ResendResult
	if err := c.Get("/resendUnconfirmedTxns", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// RawTransaction makes a request to /rawtx
func (c *Client) RawTransaction(txid string) (string, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/rawtx?" + v.Encode()

	var rawTx string
	if err := c.Get(endpoint, &rawTx); err != nil {
		return "", err
	}
	return rawTx, nil
}

// AddressTransactions makes a request to /explorer/address
func (c *Client) AddressTransactions(addr string) ([]ReadableTransaction, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/explorer/address?" + v.Encode()

	var b []ReadableTransaction
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// RichlistParams are arguments to the /richlist endpoint
type RichlistParams struct {
	N                   int
	IncludeDistribution bool
}

// Richlist makes a request to /richlist
func (c *Client) Richlist(params *RichlistParams) (*Richlist, error) {
	endpoint := "/richlist"

	if params != nil {
		v := url.Values{}
		v.Add("n", fmt.Sprint(params.N))
		v.Add("include-distribution", fmt.Sprint(params.IncludeDistribution))
		endpoint = "/richlist?" + v.Encode()
	}

	var r Richlist
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// AddressCount makes a request to /addresscount
func (c *Client) AddressCount() (uint64, error) {
	var r struct {
		Count uint64 `json:"count"`
	}
	if err := c.Get("/addresscount", &r); err != nil {
		return 0, err
	}
	return r.Count, nil

}

// UnloadWallet make a request to /wallet/unload
func (c *Client) UnloadWallet(id string) error {
	v := url.Values{}
	v.Add("id", id)
	return c.PostForm("/wallet/unload", strings.NewReader(v.Encode()), nil)
}
