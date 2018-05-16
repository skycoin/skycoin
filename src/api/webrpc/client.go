package webrpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

const (
	dialTimeout         = 60 * time.Second
	httpClientTimeout   = 120 * time.Second
	tlsHandshakeTimeout = 60 * time.Second
)

// ErrJSONUnmarshal is returned if JSON unmarshal fails
var ErrJSONUnmarshal = errors.New("JSON unmarshal failed")

// ClientError is used for non-200 API responses
type ClientError struct {
	Status     string
	StatusCode int
	Message    string
}

func (e ClientError) Error() string {
	return e.Message
}

// Client is an RPC client
type Client struct {
	Addr       string
	HTTPClient *http.Client
	UseCSRF    bool
	reqIDCtr   int
}

// NewClient creates a Client
func NewClient(addr string) (*Client, error) {
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

	if _, err := url.Parse(addr); err != nil {
		return nil, err
	}

	return &Client{
		Addr:       addr,
		HTTPClient: httpClient,
	}, nil
}

// Do makes an RPC request
func (c *Client) Do(obj interface{}, method string, params interface{}) error {
	c.reqIDCtr++

	var csrf string
	if c.UseCSRF {
		var err error
		csrf, err = c.CSRF()
		if err != nil {
			return err
		}

		if csrf == "" {
			return errors.New("Remote node has CSRF disabled")
		}
	}

	req, err := NewRequest(method, params, strconv.Itoa(c.reqIDCtr))
	if err != nil {
		return err
	}

	rsp, err := do(c.HTTPClient, req, c.Addr, csrf)
	if err != nil {
		return err
	}

	if rsp.Error != nil {
		return rsp.Error
	}

	return decodeJSON(rsp.Result, obj)
}

// CSRF returns a CSRF token. If CSRF is disabled on the node, returns an empty string and nil error.
func (c *Client) CSRF() (string, error) {
	endpoint := c.Addr + "csrf"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)

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

		return "", fmt.Errorf("%d %s: %s", resp.StatusCode, resp.Status, string(body))
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

// GetUnspentOutputs returns unspent outputs for a set of addresses
// TODO -- what is the difference between this and GetAddressUxOuts?
func (c *Client) GetUnspentOutputs(addrs []string) (*OutputsResult, error) {
	outputs := OutputsResult{}
	if err := c.Do(&outputs, "get_outputs", addrs); err != nil {
		return nil, err
	}

	return &outputs, nil
}

// InjectTransactionString injects a hex-encoded transaction string to the network
func (c *Client) InjectTransactionString(rawtx string) (string, error) {
	params := []string{rawtx}
	rlt := TxIDJson{}

	if err := c.Do(&rlt, "inject_transaction", params); err != nil {
		return "", err
	}

	return rlt.Txid, nil
}

// InjectTransaction injects a *coin.Transaction to the network
func (c *Client) InjectTransaction(tx *coin.Transaction) (string, error) {
	d := tx.Serialize()
	rawTx := hex.EncodeToString(d)
	return c.InjectTransactionString(rawTx)
}

// GetStatus returns status info for a skycoin node
func (c *Client) GetStatus() (*StatusResult, error) {
	status := StatusResult{}
	if err := c.Do(&status, "get_status", nil); err != nil {
		return nil, err
	}

	return &status, nil
}

// GetTransactionByID returns a transaction given a txid
func (c *Client) GetTransactionByID(txid string) (*TxnResult, error) {
	txn := TxnResult{}
	if err := c.Do(&txn, "get_transaction", []string{txid}); err != nil {
		return nil, err
	}

	return &txn, nil
}

// GetAddressUxOuts returns unspent outputs for a set of addresses
// TODO -- what is the difference between this and GetUnspentOutputs?
func (c *Client) GetAddressUxOuts(addrs []string) ([]AddrUxoutResult, error) {
	uxouts := []AddrUxoutResult{}
	if err := c.Do(&uxouts, "get_address_uxouts", addrs); err != nil {
		return nil, err
	}

	return uxouts, nil
}

// GetBlocks returns a range of blocks
func (c *Client) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{start, end}
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks", param); err != nil {
		return nil, err
	}

	return &blocks, nil
}

// GetBlocksBySeq returns blocks for a set of block sequences (heights)
func (c *Client) GetBlocksBySeq(ss []uint64) (*visor.ReadableBlocks, error) {
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks_by_seq", ss); err != nil {
		return nil, err
	}

	return &blocks, nil
}

// GetLastBlocks returns the last n blocks
func (c *Client) GetLastBlocks(n uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{n}
	blocks := visor.ReadableBlocks{}
	if err := c.Do(&blocks, "get_lastblocks", param); err != nil {
		return nil, err
	}

	return &blocks, nil
}

// do send request to web. rpcAddress should have forward slash appended
func do(httpClient *http.Client, rpcReq *Request, rpcAddress, csrf string) (*Response, error) {
	d, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	url := rpcAddress + "api/v1/webrpc"
	body := bytes.NewBuffer(d)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if csrf != "" {
		req.Header.Set("X-CSRF-Token", csrf)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, ClientError{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    strings.TrimSpace(string(body)),
		}
	}

	res := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeJSON(data []byte, obj interface{}) error {
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(obj); err != nil {
		return ErrJSONUnmarshal
	}
	return nil
}
