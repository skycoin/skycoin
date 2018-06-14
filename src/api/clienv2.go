package api

import (
	"fmt"
	"net/url"

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
