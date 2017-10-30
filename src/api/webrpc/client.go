package webrpc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

// ErrJSONUnmarshal is returned if JSON unmarshal fails
var ErrJSONUnmarshal = errors.New("JSON unmarshal failed")

// Client is an RPC client
type Client struct {
	Addr     string
	reqIDCtr int
}

// Do makes an RPC request
func (c *Client) Do(obj interface{}, method string, params interface{}) error {
	c.reqIDCtr++
	req, err := NewRequest(method, params, strconv.Itoa(c.reqIDCtr))
	if err != nil {
		return err
	}

	rsp, err := Do(req, c.Addr)
	if err != nil {
		return err
	}

	if rsp.Error != nil {
		return rsp.Error
	}

	return decodeJSON(rsp.Result, obj)
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

// Do send request to web
func Do(req *Request, rpcAddress string) (*Response, error) {
	d, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	rsp, err := http.Post(fmt.Sprintf("http://%s/webrpc", rpcAddress), "application/json", bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	res := Response{}
	if err := json.NewDecoder(rsp.Body).Decode(&res); err != nil {
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
