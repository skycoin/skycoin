package webrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/visor"
)

var ErrJSONUnmarshal = errors.New("json unmarshal failed")

type Client struct {
	Addr     string
	reqIdCtr int
}

func (c *Client) Do(obj interface{}, method string, params interface{}) error {
	c.reqIdCtr++
	req, err := NewRequest(method, params, strconv.Itoa(c.reqIdCtr))
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

	return decodeJson(rsp.Result, obj)
}

func (c *Client) GetUnspentOutputs(addrs []string) (*OutputsResult, error) {
	outputs := OutputsResult{}
	if err := c.Do(&outputs, "get_outputs", addrs); err != nil {
		return nil, err
	}

	return &outputs, nil
}

// Returns TxId
func (c *Client) InjectTransaction(rawtx string) (string, error) {
	params := []string{rawtx}
	rlt := TxIDJson{}

	if err := c.Do(&rlt, "inject_transaction", params); err != nil {
		return "", err
	}

	return rlt.Txid, nil
}

func (c *Client) GetStatus() (*StatusResult, error) {
	status := StatusResult{}
	if err := c.Do(&status, "get_status", nil); err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *Client) GetTransactionByID(txid string) (*TxnResult, error) {
	txn := TxnResult{}
	if err := c.Do(&txn, "get_transaction", []string{txid}); err != nil {
		return nil, err
	}

	return &txn, nil
}

func (c *Client) GetAddressUxOuts(addrs []string) ([]AddrUxoutResult, error) {
	uxouts := []AddrUxoutResult{}
	if err := c.Do(&uxouts, "get_address_uxouts", addrs); err != nil {
		return nil, err
	}

	return uxouts, nil
}

func (c *Client) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{start, end}
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks", param); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (c *Client) GetBlocksBySeq(ss []uint64) (*visor.ReadableBlocks, error) {
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks_by_seq", ss); err != nil {
		return nil, err
	}

	return &blocks, nil
}

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

func decodeJson(data []byte, obj interface{}) error {
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(obj); err != nil {
		return ErrJSONUnmarshal
	}
	return nil
}
