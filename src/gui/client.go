package gui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/skycoin/skycoin/src/wallet"
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
	Addr string
}

// NewClient creates a Client
func NewClient(addr string) *Client {
	addr = strings.TrimRight(addr, "/")
	addr += "/"

	return &Client{
		Addr: addr,
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

	return json.NewDecoder(resp.Body).Decode(&obj)
}

// get makes a GET request to an endpoint. Caller must close response body.
func (c *Client) get(endpoint string) (*http.Response, error) {
	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}

// Post makes a POST request to an endpoint. Caller must close response body.
func (c *Client) Post(endpoint string, body io.Reader) (*http.Response, error) {
	csrf, err := c.CSRF()
	if err != nil {
		return nil, err
	}

	endpoint = strings.TrimLeft(endpoint, "/")
	endpoint = c.Addr + endpoint

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-CSRF-Token", csrf)
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
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

// Outputs makes a request to /outputs?addrs=xxx&hashes=zzz
func (c *Client) Outputs(addrs, hashes []string) (*visor.ReadableOutputSet, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("hashes", strings.Join(hashes, ","))

	var o visor.ReadableOutputSet
	if err := c.Get("/outputs?"+v.Encode(), &o); err != nil {
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
	var b visor.ReadableBlock
	if err := c.Get(fmt.Sprintf("/block?hash=%s", hash), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeq makes a request to /block?seq=xxx
func (c *Client) BlockBySeq(seq uint64) (*visor.ReadableBlock, error) {
	var b visor.ReadableBlock
	if err := c.Get(fmt.Sprintf("/block?seq=%d", seq), &b); err != nil {
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
	var b wallet.BalancePair
	endpoint := fmt.Sprintf("/balance?addrs=%s", strings.Join(addrs, ","))
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// UxOut makes a request to /uxout?uxid=xxx
func (c *Client) UxOut(uxID string) (*historydb.UxOutJSON, error) {
	var b historydb.UxOutJSON
	endpoint := fmt.Sprintf("/uxout?uxid=%s", uxID)
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}
