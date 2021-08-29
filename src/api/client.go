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

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/kvstorage"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second

	// ContentTypeJSON json content type header
	ContentTypeJSON = "application/json"
	// ContentTypeForm form data content type header
	ContentTypeForm = "application/x-www-form-urlencoded"
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
	Username   string
	Password   string
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

// SetAuth configures the Client's request authentication
func (c *Client) SetAuth(username, password string) {
	c.Username = username
	c.Password = password
}

func (c *Client) applyAuth(req *http.Request) {
	if c.Username == "" && c.Password == "" {
		return
	}

	req.SetBasicAuth(c.Username, c.Password)
}

// GetV2 makes a GET request to an endpoint and unmarshals the response to respObj.
// If the response is not 200 OK, returns an error
func (c *Client) GetV2(endpoint string, respObj interface{}) (bool, error) {
	return c.requestV2(http.MethodGet, endpoint, nil, respObj)
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
	return c.makeRequestWithoutBody(endpoint, http.MethodGet)
}

// makeRequestWithoutBody makes a `method` request to an endpoint. Caller must close response body.
func (c *Client) makeRequestWithoutBody(endpoint, method string) (*http.Response, error) {
	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return nil, err
	}

	c.applyAuth(req)

	return c.HTTPClient.Do(req)
}

// DeleteV2 makes a DELETE request to an endpoint with body of json data,
// and parses the standard JSON response.
func (c *Client) DeleteV2(endpoint string, respObj interface{}) (bool, error) {
	return c.requestV2(http.MethodDelete, endpoint, nil, respObj)
}

// PostForm makes a POST request to an endpoint with body of ContentTypeForm formated data.
func (c *Client) PostForm(endpoint string, body io.Reader, obj interface{}) error {
	return c.Post(endpoint, ContentTypeForm, body, obj)
}

// PostJSON makes a POST request to an endpoint with body of json data.
func (c *Client) PostJSON(endpoint string, reqObj, respObj interface{}) error {
	body, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	return c.Post(endpoint, ContentTypeJSON, bytes.NewReader(body), respObj)
}

// Post makes a POST request to an endpoint.
func (c *Client) Post(endpoint string, contentType string, body io.Reader, obj interface{}) error {
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

	c.applyAuth(req)

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

	return c.requestV2(http.MethodPost, endpoint, bytes.NewReader(body), respObj)
}

func (c *Client) requestV2(method, endpoint string, body io.Reader, respObj interface{}) (bool, error) {
	csrf, err := c.CSRF()
	if err != nil {
		return false, err
	}

	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return false, err
	}

	c.applyAuth(req)

	if csrf != "" {
		req.Header.Set(CSRFHeaderName, csrf)
	}

	switch method {
	case http.MethodPost:
		req.Header.Set("Content-Type", ContentTypeJSON)
	}

	req.Header.Set("Accept", ContentTypeJSON)

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
			return false, NewClientError(resp.Status, resp.StatusCode, string(respBody))
		}

		return false, err
	}

	// The JSON decoder stops at the end of the first valid JSON object.
	// Check that there is no trailing data after the end of the first valid JSON object.
	// This could occur if an endpoint mistakenly wrote an object twice, for example.
	// This line returns the decoder's underlying read buffer. Read(nil) will return io.EOF
	// if the buffer was completely consumed.
	if _, err := decoder.Buffered().Read(nil); err != io.EOF {
		return false, NewClientError(resp.Status, resp.StatusCode, "Response has additional bytes after the first JSON object: "+string(respBody))
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

// OutputsForAddresses makes a request to POST /api/v1/outputs?addrs=xxx
func (c *Client) OutputsForAddresses(addrs []string) (*readable.UnspentOutputsSummary, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))

	endpoint := "/api/v1/outputs"

	var o readable.UnspentOutputsSummary
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &o); err != nil {
		return nil, err
	}
	return &o, nil
}

// OutputsForHashes makes a request to POST /api/v1/outputs?hashes=zzz
func (c *Client) OutputsForHashes(hashes []string) (*readable.UnspentOutputsSummary, error) {
	v := url.Values{}
	v.Add("hashes", strings.Join(hashes, ","))
	endpoint := "/api/v1/outputs"

	var o readable.UnspentOutputsSummary
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &o); err != nil {
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

// Blocks makes a request to POST /api/v1/blocks?seqs=
func (c *Client) Blocks(seqs []uint64) (*readable.Blocks, error) {
	sSeqs := make([]string, len(seqs))
	for i, x := range seqs {
		sSeqs[i] = fmt.Sprint(x)
	}

	v := url.Values{}
	v.Add("seqs", strings.Join(sSeqs, ","))
	endpoint := "/api/v1/blocks"

	var b readable.Blocks
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlocksVerbose makes a request to POST /api/v1/blocks?verbose=1&seqs=
func (c *Client) BlocksVerbose(seqs []uint64) (*readable.BlocksVerbose, error) {
	sSeqs := make([]string, len(seqs))
	for i, x := range seqs {
		sSeqs[i] = fmt.Sprint(x)
	}

	v := url.Values{}
	v.Add("seqs", strings.Join(sSeqs, ","))
	v.Add("verbose", "1")
	endpoint := "/api/v1/blocks"

	var b readable.BlocksVerbose
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlocksInRange makes a request to GET /api/v1/blocks?start=&end=
func (c *Client) BlocksInRange(start, end uint64) (*readable.Blocks, error) {
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

// BlocksInRangeVerbose makes a request to GET /api/v1/blocks?verbose=1&start=&end=
func (c *Client) BlocksInRangeVerbose(start, end uint64) (*readable.BlocksVerbose, error) {
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

// Balance makes a request to POST /api/v1/balance?addrs=xxx
func (c *Client) Balance(addrs []string) (*BalanceResponse, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/balance"

	var b BalanceResponse
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &b); err != nil {
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

// CreateWalletOptions are the options for creating a wallet
type CreateWalletOptions struct {
	Type                  string
	Seed                  string
	SeedPassphrase        string
	Label                 string
	Password              string
	ScanN                 uint64
	XPub                  string
	Encrypt               bool
	Bip44Coin             *bip44.CoinType
	CollectionPrivateKeys string
}

// CreateWallet makes a request to POST /api/v1/wallet/create and creates a wallet.
// If scanN is <= 0, the scan number defaults to 1
func (c *Client) CreateWallet(o CreateWalletOptions) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("type", o.Type)
	v.Add("seed", o.Seed)
	v.Add("label", o.Label)
	v.Add("encrypt", fmt.Sprint(o.Encrypt))

	if o.Password != "" {
		v.Add("password", o.Password)
	}

	if o.ScanN > 0 {
		v.Add("scan", fmt.Sprint(o.ScanN))
	}

	if o.SeedPassphrase != "" {
		v.Add("seed-passphrase", o.SeedPassphrase)
	}

	if o.Bip44Coin != nil {
		v.Add("bip44-coin", fmt.Sprintf("%d", *o.Bip44Coin))
	}

	if o.XPub != "" {
		v.Add("xpub", o.XPub)
	}

	if o.CollectionPrivateKeys != "" {
		v.Add("private-keys", o.CollectionPrivateKeys)
	}

	var w WalletResponse
	if err := c.PostForm("/api/v1/wallet/create", strings.NewReader(v.Encode()), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// CreateWalletTemp makes a request to POST /api/v1/wallet/createTemp and creates a
// temporary wallet.
// If scanN is <= 0, the scan number defaults to 1
func (c *Client) CreateWalletTemp(o CreateWalletOptions) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("type", o.Type)
	v.Add("seed", o.Seed)
	v.Add("label", o.Label)

	if o.ScanN > 0 {
		v.Add("scan", fmt.Sprint(o.ScanN))
	}

	if o.Bip44Coin != nil {
		v.Add("bip44-coin", fmt.Sprintf("%d", *o.Bip44Coin))
	}

	if o.XPub != "" {
		v.Add("xpub", o.XPub)
	}

	if o.CollectionPrivateKeys != "" {
		v.Add("private-keys", o.CollectionPrivateKeys)
	}

	var w WalletResponse
	if err := c.PostForm("/api/v1/wallet/createTemp", strings.NewReader(v.Encode()), &w); err != nil {
		return nil, err
	}
	return &w, nil
}

// NewWalletAddress makes a request to POST /api/v1/wallet/newAddress
// if n is <= 0, defaults to 1
func (c *Client) NewWalletAddress(id string, password string, options ...wallet.Option) ([]string, error) {
	v := url.Values{}
	v.Add("id", id)
	if len(password) > 0 {
		v.Add("password", password)
	}

	var opts wallet.AdvancedOptions
	for _, f := range options {
		f(&opts)
	}

	if opts.GenerateN > 0 {
		v.Add("num", fmt.Sprint(opts.GenerateN))
	}

	if len(opts.PrivateKeys) > 0 {
		keys := make([]string, 0, len(opts.PrivateKeys))
		for _, k := range opts.PrivateKeys {
			keys = append(keys, k.Hex())
		}
		v.Add("private-keys", strings.Join(keys, ","))
	}

	var obj struct {
		Addresses []string `json:"addresses"`
	}
	if err := c.PostForm("/api/v1/wallet/newAddress", strings.NewReader(v.Encode()), &obj); err != nil {
		return nil, err
	}
	return obj.Addresses, nil
}

// ScanWalletAddresses makes a request to POST /api/v1/wallet/scan
// if n is <= 0, defaults to 20
func (c *Client) ScanWalletAddresses(id string, n int, password string) ([]string, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)
	if n > 0 {
		v.Add("num", fmt.Sprint(n))
	}

	var obj struct {
		Addresses []string `json:"addresses"`
	}
	if err := c.PostForm("/api/v1/wallet/scan", strings.NewReader(v.Encode()), &obj); err != nil {
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

// CreateTransactionRequest is sent to /api/v2/transaction
type CreateTransactionRequest struct {
	IgnoreUnconfirmed bool           `json:"ignore_unconfirmed"`
	HoursSelection    HoursSelection `json:"hours_selection"`
	ChangeAddress     *string        `json:"change_address,omitempty"`
	To                []Receiver     `json:"to"`
	UxOuts            []string       `json:"unspents,omitempty"`
	Addresses         []string       `json:"addresses,omitempty"`
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

// WalletCreateTransactionRequest is sent to /api/v1/wallet/transaction
type WalletCreateTransactionRequest struct {
	Unsigned bool   `json:"unsigned"`
	WalletID string `json:"wallet_id"`
	Password string `json:"password"`
	CreateTransactionRequest
}

// WalletCreateTransaction makes a request to POST /api/v1/wallet/transaction
func (c *Client) WalletCreateTransaction(req WalletCreateTransactionRequest) (*CreateTransactionResponse, error) {
	var r CreateTransactionResponse
	endpoint := "/api/v1/wallet/transaction"
	if err := c.PostJSON(endpoint, req, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// WalletSignTransaction makes a request to POST /api/v2/wallet/transaction/sign
func (c *Client) WalletSignTransaction(req WalletSignTransactionRequest) (*CreateTransactionResponse, error) {
	var r CreateTransactionResponse
	endpoint := "/api/v2/wallet/transaction/sign"
	ok, err := c.PostJSONV2(endpoint, req, &r)
	if ok {
		return &r, err
	}
	return nil, err
}

// CreateTransaction makes a request to POST /api/v2/transaction
func (c *Client) CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error) {
	var r CreateTransactionResponse
	endpoint := "/api/v2/transaction"
	ok, err := c.PostJSONV2(endpoint, req, &r)
	if ok {
		return &r, err
	}
	return nil, err
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

// VerifySeed verifies whether the given seed is a valid bip39 mnemonic or not
func (c *Client) VerifySeed(seed string) (bool, error) {
	ok, err := c.PostJSONV2("/api/v2/wallet/seed/verify", VerifySeedRequest{
		Seed: seed,
	}, &struct{}{})
	if err != nil {
		return false, err
	}
	return ok, nil
}

// WalletSeed makes a request to POST /api/v1/wallet/seed
func (c *Client) WalletSeed(id string, password string) (*WalletSeedResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)

	var r WalletSeedResponse
	if err := c.PostForm("/api/v1/wallet/seed", strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}

	return &r, nil
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

// NetworkConnectionsFilter filters for network connections
type NetworkConnectionsFilter struct {
	States    []daemon.ConnectionState // "pending", "connected" and "introduced"
	Direction string                   // "incoming" or "outgoing"
}

// NetworkConnections makes a request to GET /api/v1/network/connections.
// Connections can be filtered by state and direction. By default, "connected" and "introduced" connections
// of both directions are returned.
func (c *Client) NetworkConnections(filters *NetworkConnectionsFilter) (*Connections, error) {
	v := url.Values{}
	if filters != nil {
		if len(filters.States) != 0 {
			states := make([]string, len(filters.States))
			for i, s := range filters.States {
				states[i] = string(s)
			}
			v.Add("states", strings.Join(states, ","))
		}
		if filters.Direction != "" {
			v.Add("direction", filters.Direction)
		}
	}
	endpoint := "/api/v1/network/connections?" + v.Encode()

	var dc Connections
	if err := c.Get(endpoint, &dc); err != nil {
		return nil, err
	}
	return &dc, nil
}

// NetworkDefaultPeers makes a request to GET /api/v1/network/defaultConnections
func (c *Client) NetworkDefaultPeers() ([]string, error) {
	var dc []string
	if err := c.Get("/api/v1/network/defaultConnections", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkTrustedPeers makes a request to GET /api/v1/network/connections/trust
func (c *Client) NetworkTrustedPeers() ([]string, error) {
	var dc []string
	if err := c.Get("/api/v1/network/connections/trust", &dc); err != nil {
		return nil, err
	}
	return dc, nil
}

// NetworkExchangedPeers makes a request to GET /api/v1/network/connections/exchange
func (c *Client) NetworkExchangedPeers() ([]string, error) {
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

// Transactions makes a request to POST /api/v1/transactions
func (c *Client) Transactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatus
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// ConfirmedTransactions makes a request to POST /api/v1/transactions?confirmed=true
func (c *Client) ConfirmedTransactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatus
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// UnconfirmedTransactions makes a request to POST /api/v1/transactions?confirmed=false
func (c *Client) UnconfirmedTransactions(addrs []string) ([]readable.TransactionWithStatus, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatus
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// TransactionsVerbose makes a request to POST /api/v1/transactions?verbose=1
func (c *Client) TransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatusVerbose
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// ConfirmedTransactionsVerbose makes a request to POST /api/v1/transactions?confirmed=true&verbose=1
func (c *Client) ConfirmedTransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatusVerbose
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// UnconfirmedTransactionsVerbose makes a request to POST /api/v1/transactions?confirmed=false&verbose=1
func (c *Client) UnconfirmedTransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	v.Add("verbose", "1")
	endpoint := "/api/v1/transactions"

	var r []readable.TransactionWithStatusVerbose
	if err := c.PostForm(endpoint, strings.NewReader(v.Encode()), &r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetTransactionsNum makes a GET request to /api/v1/transactions/num
func (c *Client) GetTransactionsNum() (uint64, error) {
	var r struct {
		TxnsNum uint64 `json:"txns_num"`
	}
	if err := c.Get("api/v1/transactions/num", &r); err != nil {
		return 0, err
	}

	return r.TxnsNum, nil
}

// InjectTransaction makes a request to POST /api/v1/injectTransaction.
func (c *Client) InjectTransaction(txn *coin.Transaction) (string, error) {
	rawTxn, err := txn.SerializeHex()
	if err != nil {
		return "", err
	}
	return c.InjectEncodedTransaction(rawTxn)
}

// InjectTransactionNoBroadcast makes a request to POST /api/v1/injectTransaction
// but does not broadcast the transaction.
func (c *Client) InjectTransactionNoBroadcast(txn *coin.Transaction) (string, error) {
	rawTxn, err := txn.SerializeHex()
	if err != nil {
		return "", err
	}
	return c.InjectEncodedTransactionNoBroadcast(rawTxn)
}

// InjectEncodedTransaction makes a request to POST /api/v1/injectTransaction.
// rawTxn is a hex-encoded, serialized transaction
func (c *Client) InjectEncodedTransaction(rawTxn string) (string, error) {
	return c.injectEncodedTransaction(rawTxn, false)
}

// InjectEncodedTransactionNoBroadcast makes a request to POST /api/v1/injectTransaction
// but does not broadcast the transaction.
// rawTxn is a hex-encoded, serialized transaction
func (c *Client) InjectEncodedTransactionNoBroadcast(rawTxn string) (string, error) {
	return c.injectEncodedTransaction(rawTxn, true)
}

func (c *Client) injectEncodedTransaction(rawTxn string, noBroadcast bool) (string, error) {
	v := InjectTransactionRequest{
		RawTxn:      rawTxn,
		NoBroadcast: noBroadcast,
	}

	var txid string
	if err := c.PostJSON("/api/v1/injectTransaction", v, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// ResendUnconfirmedTransactions makes a request to POST /api/v1/resendUnconfirmedTxns
func (c *Client) ResendUnconfirmedTransactions() (*ResendResult, error) {
	endpoint := "/api/v1/resendUnconfirmedTxns"
	var r ResendResult
	if err := c.PostForm(endpoint, strings.NewReader(""), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// RawTransaction makes a request to GET /api/v1/rawtx
func (c *Client) RawTransaction(txid string) (string, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/api/v1/rawtx?" + v.Encode()

	var rawTxn string
	if err := c.Get(endpoint, &rawTxn); err != nil {
		return "", err
	}
	return rawTxn, nil
}

// VerifyTransaction makes a request to POST /api/v2/transaction/verify.
func (c *Client) VerifyTransaction(req VerifyTransactionRequest) (*VerifyTransactionResponse, error) {
	var rsp VerifyTransactionResponse
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
func (c *Client) EncryptWallet(id, password string) (*WalletResponse, error) {
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
func (c *Client) DecryptWallet(id, password string) (*WalletResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	v.Add("password", password)
	var wlt WalletResponse
	if err := c.PostForm("/api/v1/wallet/decrypt", strings.NewReader(v.Encode()), &wlt); err != nil {
		return nil, err
	}

	return &wlt, nil
}

// RecoverWallet makes a request to POST /api/v2/wallet/recover to recover an encrypted wallet by seed.
// The password argument is optional, if provided, the recovered wallet will be encrypted with this password,
// otherwise the recovered wallet will be unencrypted.
func (c *Client) RecoverWallet(req WalletRecoverRequest) (*WalletResponse, error) {
	var rsp WalletResponse
	ok, err := c.PostJSONV2("/api/v2/wallet/recover", req, &rsp)
	if ok {
		return &rsp, err
	}

	return nil, err
}

// Disconnect disconnect a connections by ID
func (c *Client) Disconnect(id uint64) error {
	v := url.Values{}
	v.Add("id", fmt.Sprint(id))

	var obj struct{}
	return c.PostForm("/api/v1/network/connection/disconnect", strings.NewReader(v.Encode()), &obj)
}

// GetAllStorageValues makes a GET request to /api/v2/data to get all the values from the storage of
// `storageType` type
func (c *Client) GetAllStorageValues(storageType kvstorage.Type) (map[string]string, error) {
	var values map[string]string
	ok, err := c.GetV2(fmt.Sprintf("/api/v2/data?type=%s", storageType), &values)
	if !ok {
		return nil, err
	}

	return values, err
}

// GetStorageValue makes a GET request to /api/v2/data to get the value associated with `key` from storage
// of `storageType` type
func (c *Client) GetStorageValue(storageType kvstorage.Type, key string) (string, error) {
	var value string
	ok, err := c.GetV2(fmt.Sprintf("/api/v2/data?type=%s&key=%s", storageType, key), &value)
	if !ok {
		return "", err
	}

	return value, err
}

// AddStorageValue make a POST request to /api/v2/data to add a value with the key to the storage
// of `storageType` type
func (c *Client) AddStorageValue(storageType kvstorage.Type, key, val string) error {
	_, err := c.PostJSONV2("/api/v2/data", StorageRequest{
		StorageType: storageType,
		Key:         key,
		Val:         val,
	}, nil)

	return err
}

// RemoveStorageValue makes a DELETE request to /api/v2/data to remove a value associated with the `key`
// from the storage of `storageType` type
func (c *Client) RemoveStorageValue(storageType kvstorage.Type, key string) error {
	_, err := c.DeleteV2(fmt.Sprintf("/api/v2/data?type=%s&key=%s", storageType, key), nil)

	return err
}

// RequestArg is the general data type for sending request
type RequestArg struct {
	Key   string
	Value string
}

// TransactionsWithStatusV2 represents transactions result with page info
type TransactionsWithStatusV2 struct {
	PageInfo readable.PageInfo                `json:"page_info"`
	Txns     []readable.TransactionWithStatus `json:"txns"`
}

// TransactionsWithStatusVerboseV2 represents verbose transactions result with page info
type TransactionsWithStatusVerboseV2 struct {
	PageInfo readable.PageInfo                       `json:"page_info"`
	Txns     []readable.TransactionWithStatusVerbose `json:"txns"`
}

// TransactionsV2 make a GET request to /api/v2/transaction to get transactions with no verbose.
func (c *Client) TransactionsV2(args ...RequestArg) (*TransactionsWithStatusV2, error) {
	kvs := make([]string, len(args))
	for i, arg := range args {
		kvs[i] = fmt.Sprintf("%s=%s", arg.Key, arg.Value)
		if strings.Contains(arg.Key, "verbose") {
			return nil, errors.New("arguments should not include 'verbose'")
		}
	}

	endpoint := fmt.Sprintf("/api/v2/transactions?%s", strings.Join(kvs, "&"))

	var obj TransactionsWithStatusV2
	_, err := c.GetV2(endpoint, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

// TransactionsVerboseV2 make a GET request to /api/v2/transaction to get transactions with no verbose.
func (c *Client) TransactionsVerboseV2(args ...RequestArg) (*TransactionsWithStatusVerboseV2, error) {
	kvs := make([]string, len(args))
	for i, arg := range args {
		kvs[i] = fmt.Sprintf("%s=%s", arg.Key, arg.Value)
	}

	endpoint := fmt.Sprintf("/api/v2/transactions?verbose=true&%s", strings.Join(kvs, "&"))

	var obj TransactionsWithStatusVerboseV2
	_, err := c.GetV2(endpoint, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
