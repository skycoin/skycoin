package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/skycoin/skycoin/src/visor"
)

// ClientV2 provides an interface to a remote node's HTTP API
type ClientV2 struct {
	Client
}

// NewClientV2 creates a Client
func NewClientV2(addr string) *ClientV2 {
	return &ClientV2{*NewClient(addr)}
}

// Get Adds extra data to response
func (c *ClientV2) Get(endpoint string, obj interface{}) error {
	var resp ReceivedHTTPResponse
	var err error
	err = c.Client.Get(endpoint, &resp)
	if err != nil {
		return err
	} else if resp.Error != nil {
		err = ClientError{
			Status:     http.StatusText(resp.Error.Code),
			StatusCode: resp.Error.Code,
		}
	} else if resp.Data == nil {
		err = ClientError{
			Status:     http.StatusText(http.StatusInternalServerError),
			StatusCode: http.StatusInternalServerError,
		}
	} else {
		err = json.Unmarshal(resp.Data, obj)
	}
	return err
}

// BlockByHash makes a request to GET /api/v2/block?hash=xxx
func (c *ClientV2) BlockByHash(hash string) (*visor.ReadableBlockV2, error) {
	v := url.Values{}
	v.Add("hash", hash)
	endpoint := "/api/v2/block?" + v.Encode()

	var b visor.ReadableBlockV2
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// BlockBySeq makes a request to GET /api/v2/block?seq=xxx
func (c *ClientV2) BlockBySeq(seq uint64) (*visor.ReadableBlockV2, error) {
	v := url.Values{}
	v.Add("seq", fmt.Sprint(seq))
	endpoint := "/api/v2/block?" + v.Encode()

	var b visor.ReadableBlockV2
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Blocks makes a request to GET /api/v2/blocks
func (c *ClientV2) Blocks(start, end int) (*visor.ReadableBlocksV2, error) {
	v := url.Values{}
	v.Add("start", fmt.Sprint(start))
	v.Add("end", fmt.Sprint(end))
	endpoint := "/api/v2/blocks?" + v.Encode()

	var b visor.ReadableBlocksV2
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// LastBlocks makes a request to GET /api/v1/last_blocks
func (c *ClientV2) LastBlocks(n int) (*visor.ReadableBlocksV2, error) {
	v := url.Values{}
	v.Add("num", fmt.Sprint(n))
	endpoint := "/api/v2/last_blocks?" + v.Encode()

	var b visor.ReadableBlocksV2
	if err := c.Get(endpoint, &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// PendingTransactions makes a request to GET /api/v2/pendingTxs
func (c *ClientV2) PendingTransactions() (*visor.ReadableUnconfirmedTxnsV2, error) {
	var v visor.ReadableUnconfirmedTxnsV2
	if err := c.Get("/api/v2/pendingTxs", &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Transaction makes a request to GET /api/v2/transaction
func (c *ClientV2) Transaction(txid string) (*visor.TransactionResultV2, error) {
	v := url.Values{}
	v.Add("txid", txid)
	endpoint := "/api/v2/transaction?" + v.Encode()

	var r visor.TransactionResultV2
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Transactions makes a request to GET /api/v2/transactions
func (c *ClientV2) Transactions(addrs []string) (*visor.TransactionResultsV2, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	endpoint := "/api/v2/transactions?" + v.Encode()

	var r visor.TransactionResultsV2
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ConfirmedTransactions makes a request to GET /api/v2/transactions?confirmed=true
func (c *ClientV2) ConfirmedTransactions(addrs []string) (*visor.TransactionResultsV2, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "true")
	endpoint := "/api/v2/transactions?" + v.Encode()

	var r visor.TransactionResultsV2
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UnconfirmedTransactions makes a request to GET /api/v1/transactions?confirmed=false
func (c *ClientV2) UnconfirmedTransactions(addrs []string) (*visor.TransactionResultsV2, error) {
	v := url.Values{}
	v.Add("addrs", strings.Join(addrs, ","))
	v.Add("confirmed", "false")
	endpoint := "/api/v2/transactions?" + v.Encode()

	var r visor.TransactionResultsV2
	if err := c.Get(endpoint, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
