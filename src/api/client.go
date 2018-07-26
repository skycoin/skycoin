package api

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

// ClientError is used for non-200 API responses
type ClientError struct {
	Status     string
	StatusCode int
	Message    string
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

		return ClientError{
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

		return ClientError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
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
			return false, ClientError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}

		return false, err
	}

	var rspErr error
	if resp.StatusCode != http.StatusOK {
		rspErr = ClientError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    wrapObj.Error.Message,
		}
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

		return "", ClientError{
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

// Version makes a request to GET /api/v1/version
func (c *Client) Version() (*visor.BuildInfo, error) {
	var bi visor.BuildInfo
	if err := c.Get("/api/v1/version", &bi); err != nil {
		return nil, err
	}
	return &bi, nil
}

// Outputs makes a request to GET /api/v1/outputs
func (c *Client) Outputs() (*visor.ReadableOutputSet, error) {
	var o visor.ReadableOutputSet
	if err := c.Get("/api/v1/outputs", &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForAddresses makes a request to GET /api/v1/outputs?addrs=xxx
func (c *Client) OutputsForAddresses(addrs []string) (*visor.ReadableOutputSet, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/outputs?" + v.Encode()

	var o visor.ReadableOutputSet
	if err := c.Get(endpoint, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForHashes makes a request to GET /api/v1/outputs?hashes=zzz
func (c *Client) OutputsForHashes(hashes []string) (*visor.ReadableOutputSet, error) {
	v := url.Values{}
	v.Add("hashes", strings.Join(hashes, ","))
	endpoint := "/api/v1/outputs?" + v.Encode()

	var o visor.ReadableOutputSet
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
func (c *Client) BlockByHash(hash string) (*visor.ReadableBlock, error) {
	v := url.Values{}
	v.Add("hash", hash)
	endpoint := "/api/v1/block?" + v.Encode()

	var b visor.ReadableBlock
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeq makes a request to GET /api/v1/block?seq=xxx
func (c *Client) BlockBySeq(seq uint64) (*visor.ReadableBlock, error) {
	v := url.Values{}
	v.Add("seq", fmt.Sprint(seq))
	endpoint := "/api/v1/block?" + v.Encode()

	var b visor.ReadableBlock
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Blocks makes a request to GET /api/v1/blocks
func (c *Client) Blocks(start, end int) (*visor.ReadableBlocks, error) {
	v := url.Values{}
	v.Add("start", fmt.Sprint(start))
	v.Add("end", fmt.Sprint(end))
	endpoint := "/api/v1/blocks?" + v.Encode()

	var b visor.ReadableBlocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// LastBlocks makes a request to GET /api/v1/last_blocks
func (c *Client) LastBlocks(n int) (*visor.ReadableBlocks, error) {
	v := url.Values{}
	v.Add("num", fmt.Sprint(n))
	endpoint := "/api/v1/last_blocks?" + v.Encode()

	var b visor.ReadableBlocks
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainMetadata makes a request to GET /api/v1/blockchain/metadata
func (c *Client) BlockchainMetadata() (*visor.BlockchainMetadata, error) {
	var b visor.BlockchainMetadata
	if err := c.Get("/api/v1/blockchain/metadata", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockchainProgress makes a request to GET /api/v1/blockchain/progress
func (c *Client) BlockchainProgress() (*daemon.BlockchainProgress, error) {
	var b daemon.BlockchainProgress
	if err := c.Get("/api/v1/blockchain/progress", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Balance makes a request to GET /api/v1/balance?addrs=xxx
func (c *Client) Balance(addrs []string) (*wallet.BalancePair, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/balance?" + v.Encode()

	var b wallet.BalancePair
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// UxOut makes a request to GET /api/v1/uxout?uxid=xxx
func (c *Client) UxOut(uxID string) (*historydb.UxOutJSON, error) {
	v := url.Values{}
	v.Add("uxid", uxID)
	endpoint := "/api/v1/uxout?" + v.Encode()

	var b historydb.UxOutJSON
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// AddressUxOuts makes a request to GET /api/v1/address_uxouts
func (c *Client) AddressUxOuts(addr string) ([]*historydb.UxOutJSON, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/api/v1/address_uxouts?" + v.Encode()

	var b []*historydb.UxOutJSON
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
func (c *Client) Wallets() ([]*WalletResponse, error) {
	var wrs []*WalletResponse
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

// CreateTransactionRequest is sent to /wallet/transaction
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

// WalletTransactions makes a request to GET /api/v1/wallet/transactions
func (c *Client) WalletTransactions(id string) (*UnconfirmedTxnsResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	endpoint := "/api/v1/wallet/transactions?" + v.Encode()

	var utx *UnconfirmedTxnsResponse
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

// GetWalletSeed makes a request to POST /api/v1/wallet/seed
func (c *Client) GetWalletSeed(id string, password string) (string, error) {
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
func (c *Client) NetworkConnection(addr string) (*daemon.Connection, error) {
	v := url.Values{}
	v.Add("addr", addr)
	endpoint := "/api/v1/network/connection?" + v.Encode()

	var dc daemon.Connection
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
func (c *Client) PendingTransactions() ([]*visor.ReadableUnconfirmedTxn, error) {
	var v []*visor.ReadableUnconfirmedTxn
	if err := c.Get("/api/v1/pendingTxs", &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Transaction makes a request to GET /api/v1/transaction
func (c *Client) Transaction(txid string) (*daemon.TransactionResult, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/api/v1/transaction?" + v.Encode()

	var r daemon.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Transactions makes a request to GET /api/v1/transactions
func (c *Client) Transactions(addrs []string) (*[]daemon.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []daemon.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ConfirmedTransactions makes a request to GET /api/v1/transactions?confirmed=true
func (c *Client) ConfirmedTransactions(addrs []string) (*[]daemon.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []daemon.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UnconfirmedTransactions makes a request to GET /api/v1/transactions?confirmed=false
func (c *Client) UnconfirmedTransactions(addrs []string) (*[]daemon.TransactionResult, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	endpoint := "/api/v1/transactions?" + v.Encode()

	var r []daemon.TransactionResult
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// InjectTransaction makes a request to POST /api/v1/injectTransaction
func (c *Client) InjectTransaction(rawTx string) (string, error) {
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
func (c *Client) ResendUnconfirmedTransactions() (*daemon.ResendResult, error) {
	var r daemon.ResendResult
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
func (c *Client) AddressTransactions(addr string) ([]daemon.ReadableTransaction, error) {
	v := url.Values{}
	v.Add("address", addr)
	endpoint := "/api/v1/explorer/address?" + v.Encode()

	var b []daemon.ReadableTransaction
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
