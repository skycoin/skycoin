package api

import (
	"bytes"
	"encoding/hex"
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

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second
)

// ClientError is used for non-200 API responses
type ClientError struct {
	Status     string
	StatusCode int
	Message    string
}

// NewClientError creates a ClientError
func NewClientError(status string, statusCode int, message string) ClientError {
	return ClientError{
		Status:     status,
		StatusCode: statusCode,
		Message:    strings.TrimRight(message, "\n"),
	}
}

func (e ClientError) Error() string {
	return e.Message
}

// ReceivedHTTPResponse parsed a HTTPResponse received by the Client, for the V2 API
type ReceivedHTTPResponse struct {
	Error *HTTPError      `json:"error,omitempty"`
	Data  json.RawMessage `json:"data"`
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

		return NewClientError(resp.Status, resp.StatusCode, string(body))
	}

	if obj == nil {
		return nil
	}

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()
	return d.Decode(obj)
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
func (c *Client) PostForm(endpoint string, body io.Reader, obj interface{}) error {
	return c.post(endpoint, "application/x-www-form-urlencoded", body, obj)
}

// PostJSON makes a POST request to an endpoint with body of json data.
func (c *Client) PostJSON(endpoint string, reqObj, respObj interface{}) error {
	body, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	return c.post(endpoint, "application/json", bytes.NewReader(body), respObj)
}

// post makes a POST request to an endpoint.
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

	if csrf != "" {
		req.Header.Set(CSRFHeaderName, csrf)
	}

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

		return NewClientError(resp.Status, resp.StatusCode, string(body))
	}

	if obj == nil {
		return nil
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(obj)
}

// PostJSONV2 makes a POST request to an endpoint with body of json data,
// and parses the standard JSON response.
func (c *Client) PostJSONV2(endpoint string, reqObj, respObj interface{}) (bool, error) {
	body, err := json.Marshal(reqObj)
	if err != nil {
		return false, err
	}

	csrf, err := c.CSRF()
	if err != nil {
		return false, err
	}

	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return false, err
	}

	if csrf != "" {
		req.Header.Set(CSRFHeaderName, csrf)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	decoder := json.NewDecoder(bytes.NewReader(respBody))
	decoder.DisallowUnknownFields()

	var wrapObj ReceivedHTTPResponse
	if err := decoder.Decode(&wrapObj); err != nil {
		// In some cases, the server can send an error response in a non-JSON format,
		// such as a 404 when the endpoint is not registered, or if a 500 error
		// occurs in the go HTTP stack, outside of the application's control.
		// If this happens, treat the entire response body as the error message.
		if resp.StatusCode != http.StatusOK {
			return false, NewClientError(resp.Status, resp.StatusCode, string(body))
		}

		return false, err
	}

	var rspErr error
	if resp.StatusCode != http.StatusOK {
		rspErr = NewClientError(resp.Status, resp.StatusCode, wrapObj.Error.Message)
	}

	if wrapObj.Data == nil {
		return false, rspErr
	}

	decoder = json.NewDecoder(bytes.NewReader(wrapObj.Data))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(respObj); err != nil {
		return false, err
	}

	return true, rspErr
}

// CSRF returns a CSRF token. If CSRF is disabled on the node, returns an empty string and nil error.
func (c *Client) CSRF() (string, error) {
	resp, err := c.get("/api/v1/csrf")
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

		return "", NewClientError(resp.Status, resp.StatusCode, string(body))
	}

	d := json.NewDecoder(resp.Body)
	d.DisallowUnknownFields()

	var m map[string]string
	if err := d.Decode(&m); err != nil {
		return "", err
	}

	token, ok := m["csrf_token"]
	if !ok {
		return "", errors.New("csrf_token not found in response")
	}

	return token, nil
}

// Version makes a request to GET /api/v1/version
func (c *Client) Version() (*readable.BuildInfo, error) {
	var bi readable.BuildInfo
	if err := c.Get("/api/v1/version", &bi); err != nil {
		return nil, err
	}
	return &bi, nil
}

// Outputs makes a request to GET /api/v1/outputs
func (c *Client) Outputs() (*readable.UnspentOutputsSummary, error) {
	var o readable.UnspentOutputsSummary
	if err := c.Get("/api/v1/outputs", &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForAddresses makes a request to GET /api/v1/outputs?addrs=xxx
func (c *Client) OutputsForAddresses(addrs []string) (*readable.UnspentOutputsSummary, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/outputs?" + v.Encode()

	var o readable.UnspentOutputsSummary
	if err := c.Get(endpoint, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForHashes makes a request to GET /api/v1/outputs?hashes=zzz
func (c *Client) OutputsForHashes(hashes []string) (*readable.UnspentOutputsSummary, error) {
	v := url.Values{}
	v.Add("hashes", strings.Join(hashes, ","))
	endpoint := "/api/v1/outputs?" + v.Encode()

	var o readable.UnspentOutputsSummary
	if err := c.Get(endpoint, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// CoinSupply makes a request to GET /api/v1/coinSupply
func (c *Client) CoinSupply() (*CoinSupply, error) {
	var cs CoinSupply
	if err := c.Get("/api/v1/coinSupply", &cs); err != nil {
		return nil, err
	}
	return &cs, nil
}

// BlockByHash makes a request to GET /api/v1/block?hash=xxx
func (c *Client) BlockByHash(hash string) (*readable.Block, error) {
	v := url.Values{}
	v.Add("hash", hash)
	endpoint := "/api/v1/block?" + v.Encode()

	var b readable.Block
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockByHashVerbose makes a request to GET /api/v1/block?hash=xxx&verbose=1
func (c *Client) BlockByHashVerbose(hash string) (*readable.BlockVerbose, error) {
	v := url.Values{}
	v.Add("hash", hash)
	v.Add("verbose", "1")
	endpoint := "/api/v1/block?" + v.Encode()

	var b readable.BlockVerbose
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeq makes a request to GET /api/v1/block?seq=xxx
func (c *Client) BlockBySeq(seq uint64) (*readable.Block, error) {
	v := url.Values{}
	v.Add("seq", fmt.Sprint(seq))
	endpoint := "/api/v1/block?" + v.Encode()

	var b readable.Block
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeqVerbose makes a request to GET /api/v1/block?seq=xxx&verbose=1
func (c *Client) BlockBySeqVerbose(seq uint64) (*readable.BlockVerbose, error) {
	v := url.Values{}
	v.Add("seq", fmt.Sprint(seq))
	v.Add("verbose", "1")
	endpoint := "/api/v1/block?" + v.Encode()

	var b readable.BlockVerbose
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Blocks makes a request to GET /api/v1/blocks
func (c *Client) Blocks(start, end uint64) (*readable.Blocks, error) {
	v := url.Values{}
	v.Add("start", fmt.Sprint(start))
	v.Add("end", fmt.Sprint(end))
	endpoint := "/api/v1/blocks?" + v.Encode()

	var b readable.Blocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlocksVerbose makes a request to GET /api/v1/blocks?verbose=1
func (c *Client) BlocksVerbose(start, end uint64) (*readable.BlocksVerbose, error) {
	v := url.Values{}
	v.Add("start", fmt.Sprint(start))
	v.Add("end", fmt.Sprint(end))
	v.Add("verbose", "1")
	endpoint := "/api/v1/blocks?" + v.Encode()

	var b readable.BlocksVerbose
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// LastBlocks makes a request to GET /api/v1/last_blocks
func (c *Client) LastBlocks(n uint64) (*readable.Blocks, error) {
	v := url.Values{}
	v.Add("num", fmt.Sprint(n))
	endpoint := "/api/v1/last_blocks?" + v.Encode()

	var b readable.Blocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// LastBlocksVerbose makes a request to GET /api/v1/last_blocks?verbose=1
func (c *Client) LastBlocksVerbose(n uint64) (*readable.BlocksVerbose, error) {
	v := url.Values{}
	v.Add("num", fmt.Sprint(n))
	v.Add("verbose", "1")
	endpoint := "/api/v1/last_blocks?" + v.Encode()

	var b readable.BlocksVerbose
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainMetadata makes a request to GET /api/v1/blockchain/metadata
func (c *Client) BlockchainMetadata() (*readable.BlockchainMetadata, error) {
	var b readable.BlockchainMetadata
	if err := c.Get("/api/v1/blockchain/metadata", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainProgress makes a request to GET /api/v1/blockchain/progress
func (c *Client) BlockchainProgress() (*readable.BlockchainProgress, error) {
	var b readable.BlockchainProgress
	if err := c.Get("/api/v1/blockchain/progress", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Balance makes a request to GET /api/v1/balance?addrs=xxx
func (c *Client) Balance(addrs []string) (*BalanceResponse, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/balance?" + v.Encode()

	var b BalanceResponse
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// UxOut makes a request to GET /api/v1/uxout?uxid=xxx
func (c *Client) UxOut(uxID string) (*readable.SpentOutput, error) {
	v := url.Values{}
	v.Add("uxid", uxID)
	endpoint := "/api/v1/uxout?" + v.Encode()

	var b readable.SpentOutput
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// AddressUxOuts makes a request to GET /api/v1/address_uxouts
func (c *Client) AddressUxOuts(addr string) ([]readable.SpentOutput, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/api/v1/address_uxouts?" + v.Encode()

	var b []readable.SpentOutput
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return b, nil
}

// Wallet makes a request to GET /api/v1/wallet
func (c *Client) Wallet(id string) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/api/v1/wallet?" + v.Encode()

	var wr WalletResponse
	if err := c.Get(endpoint, &wr); err != nil {
		return nil, err
	}

	return &wr, nil
}

// Wallets makes a request to GET /api/v1/wallets
func (c *Client) Wallets() ([]WalletResponse, error) {
	var wrs []WalletResponse
	if err := c.Get("/api/v1/wallets", &wrs); err != nil {
		return nil, err
	}

	return wrs, nil
}

// CreateUnencryptedWallet makes a request to POST /api/v1/wallet/create and creates
// a wallet without encryption.
// If scanN is <= 0, the scan number defaults to 1
func (c *Client) CreateUnencryptedWallet(seed, label string, scanN int) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("seed", seed)
	v.Add("label", label)
	v.Add("encrypt", "false")

	if scanN > 0 {
		v.Add("scan", fmt.Sprint(scanN))
	}

	var w WalletResponse
	if err := c.PostForm("/api/v1/wallet/create", strings.NewReader(v.Encode()), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// CreateEncryptedWallet makes a request to POST /api/v1/wallet/create and try to create
// a wallet with encryption.
// If scanN is <= 0, the scan number defaults to 1
func (c *Client) CreateEncryptedWallet(seed, label, password string, scanN int) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("seed", seed)
	v.Add("label", label)
	v.Add("encrypt", "true")
	v.Add("password", password)

	if scanN > 0 {
		v.Add("scan", fmt.Sprint(scanN))
	}

	var w WalletResponse
	if err := c.PostForm("/api/v1/wallet/create", strings.NewReader(v.Encode()), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// NewWalletAddress makes a request to POST /api/v1/wallet/newAddress
// if n is <= 0, defaults to 1
func (c *Client) NewWalletAddress(id string, n int, password string) ([]string, error) {
	v := url.Values{}
	v.Add("id", id)
	if n > 0 {
		v.Add("num", fmt.Sprint(n))
	}

	v.Add("password", password)

	var obj struct {
		Addresses []string `json:"addresses"`
	}
	if err := c.PostForm("/api/v1/wallet/newAddress", strings.NewReader(v.Encode()), &obj); err != nil {
		return nil, err
	}
	return obj.Addresses, nil
}

// WalletBalance makes a request to GET /api/v1/wallet/balance
func (c *Client) WalletBalance(id string) (*BalanceResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/api/v1/wallet/balance?" + v.Encode()

	var b BalanceResponse
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Spend makes a request to POST /api/v1/wallet/spend
func (c *Client) Spend(id, dst string, coins uint64, password string) (*SpendResult, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("dst", dst)
	v.Add("coins", fmt.Sprint(coins))
	v.Add("password", password)

	var r SpendResult
	endpoint := "/api/v1/wallet/spend"
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// CreateTransactionRequest is sent to /api/v1/wallet/transaction
type CreateTransactionRequest struct {
	IgnoreUnconfirmed bool                           `json:"ignore_unconfirmed"`
	HoursSelection    HoursSelection                 `json:"hours_selection"`
	Wallet            CreateTransactionRequestWallet `json:"wallet"`
	ChangeAddress     *string                        `json:"change_address,omitempty"`
	To                []Receiver                     `json:"to"`
}

// CreateTransactionRequestWallet defines a wallet to spend from and optionally which addresses in the wallet
type CreateTransactionRequestWallet struct {
	ID        string   `json:"id"`
	UxOuts    []string `json:"unspents,omitempty"`
	Addresses []string `json:"addresses,omitempty"`
	Password  string   `json:"password"`
}

// HoursSelection defines options for hours distribution
type HoursSelection struct {
	Type        string `json:"type"`
	Mode        string `json:"mode"`
	ShareFactor string `json:"share_factor,omitempty"`
}

// Receiver specifies a spend destination
type Receiver struct {
	Address string `json:"address"`
	Coins   string `json:"coins"`
	Hours   string `json:"hours,omitempty"`
}

// CreateTransaction makes a request to POST /api/v1/wallet/transaction
func (c *Client) CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error) {
	var r CreateTransactionResponse
	endpoint := "/api/v1/wallet/transaction"
	if err := c.PostJSON(endpoint, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// WalletUnconfirmedTransactions makes a request to GET /api/v1/wallet/transactions
func (c *Client) WalletUnconfirmedTransactions(id string) (*UnconfirmedTxnsResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/api/v1/wallet/transactions?" + v.Encode()

	var utx *UnconfirmedTxnsResponse
	if err := c.Get(endpoint, &utx); err != nil {
		return nil, err
	}
	return utx, nil
}

// WalletUnconfirmedTransactionsVerbose makes a request to GET /api/v1/wallet/transactions&verbose=1
func (c *Client) WalletUnconfirmedTransactionsVerbose(id string) (*UnconfirmedTxnsVerboseResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("verbose", "1")
	endpoint := "/api/v1/wallet/transactions?" + v.Encode()

	var utx *UnconfirmedTxnsVerboseResponse
	if err := c.Get(endpoint, &utx); err != nil {
		return nil, err
	}
	return utx, nil
}

// UpdateWallet makes a request to POST /api/v1/wallet/update
func (c *Client) UpdateWallet(id, label string) error {
	v := url.Values{}
	v.Add("id", id)
	v.Add("label", label)

	return c.PostForm("/api/v1/wallet/update", strings.NewReader(v.Encode()), nil)
}

// WalletFolderName makes a request to GET /api/v1/wallets/folderName
func (c *Client) WalletFolderName() (*WalletFolder, error) {
	var w WalletFolder
	if err := c.Get("/api/v1/wallets/folderName", &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// NewSeed makes a request to GET /api/v1/wallet/newSeed
// entropy must be 128 or 256
func (c *Client) NewSeed(entropy int) (string, error) {
	v := url.Values{}
	v.Add("entropy", fmt.Sprint(entropy))
	endpoint := "/api/v1/wallet/newSeed?" + v.Encode()

	var r struct {
		Seed string `json:"seed"`
	}
	if err := c.Get(endpoint, &r); err != nil {
		return "", err
	}
	return r.Seed, nil
}

// WalletSeed makes a request to POST /api/v1/wallet/seed
func (c *Client) WalletSeed(id string, password string) (string, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)

	var r struct {
		Seed string `json:"seed"`
	}
	if err := c.PostForm("/api/v1/wallet/seed", strings.NewReader(v.Encode()), &r); err != nil {
		return "", err
	}

	return r.Seed, nil
}

// NetworkConnection makes a request to GET /api/v1/network/connection
func (c *Client) NetworkConnection(addr string) (*readable.Connection, error) {
	v := url.Values{}
	v.Add("addr", addr)
	endpoint := "/api/v1/network/connection?" + v.Encode()

	var dc readable.Connection
	if err := c.Get(endpoint, &dc); err != nil {
		return nil, err
	}
	return &dc, nil
}

// NetworkConnections makes a request to GET /api/v1/network/connections
func (c *Client) NetworkConnections() (*Connections, error) {
	var dc Connections
	if err := c.Get("/api/v1/network/connections", &dc); err != nil {
		return nil, err
	}
	return &dc, nil
}

// NetworkDefaultConnections makes a request to GET /api/v1/network/defaultConnections
func (c *Client) NetworkDefaultConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/api/v1/network/defaultConnections", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkTrustedConnections makes a request to GET /api/v1/network/connections/trust
func (c *Client) NetworkTrustedConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/api/v1/network/connections/trust", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkExchangeableConnections makes a request to GET /api/v1/network/connections/exchange
func (c *Client) NetworkExchangeableConnections() ([]string, error) {
	var dc []string
	if err := c.Get("/api/v1/network/connections/exchange", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// PendingTransactions makes a request to GET /api/v1/pendingTxs
func (c *Client) PendingTransactions() ([]readable.UnconfirmedTransactions, error) {
	var v []readable.UnconfirmedTransactions
	if err := c.Get("/api/v1/pendingTxs", &v); err != nil {
		return nil, err
	}
	return v, nil
}

// PendingTransactionsVerbose makes a request to GET /api/v1/pendingTxs?verbose=1
func (c *Client) PendingTransactionsVerbose() ([]readable.UnconfirmedTransactionVerbose, error) {
	var v []readable.UnconfirmedTransactionVerbose
	if err := c.Get("/api/v1/pendingTxs?verbose=1", &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Transaction makes a request to GET /api/v1/transaction
func (c *Client) Transaction(txid string) (*readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/api/v1/transaction?" + v.Encode()

	var r readable.TransactionWithStatus
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// TransactionVerbose makes a request to GET /api/v1/transaction?verbose=1
func (c *Client) TransactionVerbose(txid string) (*readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("txid", txid)
	v.Add("verbose", "1")
	endpoint := "/api/v1/transaction?" + v.Encode()

	var r readable.TransactionWithStatusVerbose
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// TransactionEncoded makes a request to GET /api/v1/transaction?encoded=1
func (c *Client) TransactionEncoded(txid string) (*TransactionEncodedResponse, error) {
	v := url.Values{}
	v.Add("txid", txid)
	v.Add("encoded", "1")
	endpoint := "/api/v1/transaction?" + v.Encode()

	var r TransactionEncodedResponse
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Transactions makes a request to GET /api/v1/transactions
func (c *Client) Transactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatus
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// ConfirmedTransactions makes a request to GET /api/v1/transactions?confirmed=true
func (c *Client) ConfirmedTransactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatus
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// UnconfirmedTransactions makes a request to GET /api/v1/transactions?confirmed=false
func (c *Client) UnconfirmedTransactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatus
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// TransactionsVerbose makes a request to GET /api/v1/transactions?verbose=1
func (c *Client) TransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatusVerbose
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// ConfirmedTransactionsVerbose makes a request to GET /api/v1/transactions?confirmed=true&verbose=1
func (c *Client) ConfirmedTransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatusVerbose
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// UnconfirmedTransactionsVerbose makes a request to GET /api/v1/transactions?confirmed=false&verbose=1
func (c *Client) UnconfirmedTransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []readable.TransactionWithStatusVerbose
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// InjectTransaction makes a request to POST /api/v1/injectTransaction.
func (c *Client) InjectTransaction(txn *coin.Transaction) (string, error) {
	d := txn.Serialize()
	rawTx := hex.EncodeToString(d)
	return c.InjectEncodedTransaction(rawTx)
}

// InjectEncodedTransaction makes a request to POST /api/v1/injectTransaction.
// rawTx is a hex-encoded, serialized transaction
func (c *Client) InjectEncodedTransaction(rawTx string) (string, error) {
	v := struct {
		Rawtx string `json:"rawtx"`
	}{
		Rawtx: rawTx,
	}

	var txid string
	if err := c.PostJSON("/api/v1/injectTransaction", v, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// ResendUnconfirmedTransactions makes a request to GET /api/v1/resendUnconfirmedTxns
func (c *Client) ResendUnconfirmedTransactions() (*ResendResult, error) {
	var r ResendResult
	if err := c.Get("/api/v1/resendUnconfirmedTxns", &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// RawTransaction makes a request to GET /api/v1/rawtx
func (c *Client) RawTransaction(txid string) (string, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/api/v1/rawtx?" + v.Encode()

	var rawTx string
	if err := c.Get(endpoint, &rawTx); err != nil {
		return "", err
	}
	return rawTx, nil
}

// VerifyTransaction makes a request to POST /api/v2/transaction/verify.
func (c *Client) VerifyTransaction(encodedTxn string) (*VerifyTxnResponse, error) {
	req := VerifyTxnRequest{
		EncodedTransaction: encodedTxn,
	}

	var rsp VerifyTxnResponse
	ok, err := c.PostJSONV2("/api/v2/transaction/verify", req, &rsp)
	if ok {
		return &rsp, err
	}

	return nil, err
}

// VerifyAddress makes a request to POST /api/v2/address/verify
// The API may respond with an error but include data useful for processing,
// so both return values may be non-nil.
func (c *Client) VerifyAddress(addr string) (*VerifyAddressResponse, error) {
	req := VerifyAddressRequest{
		Address: addr,
	}

	var rsp VerifyAddressResponse
	ok, err := c.PostJSONV2("/api/v2/address/verify", req, &rsp)
	if ok {
		return &rsp, err
	}

	return nil, err
}

// AddressTransactions makes a request to GET /api/v1/explorer/address
func (c *Client) AddressTransactions(addr string) ([]readable.TransactionVerbose, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/api/v1/explorer/address?" + v.Encode()

	var b []readable.TransactionVerbose
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

// Richlist makes a request to GET /api/v1/richlist
func (c *Client) Richlist(params *RichlistParams) (*Richlist, error) {
	endpoint := "/api/v1/richlist"

	if params != nil {
		v := url.Values{}
		v.Add("n", fmt.Sprint(params.N))
		v.Add("include-distribution", fmt.Sprint(params.IncludeDistribution))
		endpoint = "/api/v1/richlist?" + v.Encode()
	}

	var r Richlist
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// AddressCount makes a request to GET /api/v1/addresscount
func (c *Client) AddressCount() (uint64, error) {
	var r struct {
		Count uint64 `json:"count"`
	}
	if err := c.Get("/api/v1/addresscount", &r); err != nil {
		return 0, err
	}
	return r.Count, nil

}

// UnloadWallet makes a request to POST /api/v1/wallet/unload
func (c *Client) UnloadWallet(id string) error {
	v := url.Values{}
	v.Add("id", id)
	return c.PostForm("/api/v1/wallet/unload", strings.NewReader(v.Encode()), nil)
}

// Health makes a request to GET /api/v1/health
func (c *Client) Health() (*HealthResponse, error) {
	var r HealthResponse
	if err := c.Get("/api/v1/health", &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// EncryptWallet makes a request to POST /api/v1/wallet/encrypt to encrypt a specific wallet with the given password
func (c *Client) EncryptWallet(id string, password string) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)
	var wlt WalletResponse
	if err := c.PostForm("/api/v1/wallet/encrypt", strings.NewReader(v.Encode()), &wlt); err != nil {
		return nil, err
	}

	return &wlt, nil
}

// DecryptWallet makes a request to POST /api/v1/wallet/decrypt to decrypt a wallet
func (c *Client) DecryptWallet(id string, password string) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)
	var wlt WalletResponse
	if err := c.PostForm("/api/v1/wallet/decrypt", strings.NewReader(v.Encode()), &wlt); err != nil {
		return nil, err
	}

	return &wlt, nil
}
