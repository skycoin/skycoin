package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/api/webrpc"
	"github.com/skycoin/skycoin/src/visor"
)

type RpcClient struct {
	Addr string
}

func (c RpcClient) Do(obj interface{}, method string, params interface{}, id string) error {
	req, err := webrpc.NewRequest("get_status", nil, "1")
	if err != nil {
		return fmt.Errorf("create rpc request failed: %v", err)
	}

	rsp, err := webrpc.Do(req, c.Addr)
	if err != nil {
		return fmt.Errorf("do rpc request failed: %v", err)
	}

	if rsp.Error != nil {
		return fmt.Errorf("rpc response error: %+v", *rsp.Error)
	}

	return decodeJson(rsp.Result, obj)
}

func decodeJson(data []byte, obj interface{}) error {
	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(obj); err != nil {
		return ErrJSONUnmarshal
	}
	return nil
}

func (c *RpcClient) GetUnspent(addrs []string) (*UnspentOutSet, error) {
	outputs := webrpc.OutputsResult{}
	if err := c.Do(&outputs, "get_outputs", addrs, "1"); err != nil {
		return nil, err
	}

	return &UnspentOutSet{outputs.Outputs}, nil
}

func (c *RpcClient) GetAddressOutputs(addrs []string) (*webrpc.OutputsResult, error) {
	outputs := webrpc.OutputsResult{}
	if err := c.Do(&outputs, "get_outputs", addrs, "1"); err != nil {
		return nil, err
	}

	return &outputs, nil
}

// Returns TxId
func (c *RpcClient) BroadcastTx(rawtx string) (string, error) {
	params := []string{rawtx}
	rlt := webrpc.TxIDJson{}

	if err := c.Do(&rlt, "inject_transaction", params, "1"); err != nil {
		return "", err
	}

	return rlt.Txid, nil
}

func (c *RpcClient) GetStatus() (*webrpc.StatusResult, error) {
	status := webrpc.StatusResult{}
	if err := c.Do(&status, "get_status", nil, "1"); err != nil {
		return nil, err
	}

	return &status, nil
}

func (c *RpcClient) GetTransactionByID(txid string) (*webrpc.TxnResult, error) {
	txn := webrpc.TxnResult{}
	if err := c.Do(&txn, "get_transaction", []string{txid}, "1"); err != nil {
		return nil, err
	}

	return &txn, nil
}

func (c *RpcClient) GetBlocks(start, end uint64) (*visor.ReadableBlocks, error) {
	param := []uint64{start, end}
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks", param, "1"); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (c *RpcClient) GetBlocksBySeq(ss []uint64) (*visor.ReadableBlocks, error) {
	blocks := visor.ReadableBlocks{}

	if err := c.Do(&blocks, "get_blocks_by_seq", ss, "1"); err != nil {
		return nil, err
	}

	return &blocks, nil
}

func (c *RpcClient) GetAddressUxOuts(addrs []string) ([]webrpc.AddrUxoutResult, error) {
	uxouts := []webrpc.AddrUxoutResult{}
	if err := c.Do(&uxouts, "get_address_uxouts", addrs, "1"); err != nil {
		return nil, err
	}

	return uxouts, nil
}

func (c *RpcClient) GetLastBlocks(n uint64) (*visor.ReadableBlocks, error) {
	if n <= 0 {
		return nil, errors.New("block number must >= 0")
	}

	param := []uint64{n}
	blocks := visor.ReadableBlocks{}
	if err := c.Do(&blocks, "get_lastblocks", param, "1"); err != nil {
		return nil, err
	}

	return &blocks, nil
}
