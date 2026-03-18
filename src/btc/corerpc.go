package btc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// CoreRPCClient is a JSON-RPC client for Bitcoin Core
type CoreRPCClient struct {
	url    string
	client *http.Client
	nextID atomic.Int64
}

type coreRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type coreRPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *coreRPCError   `json:"error"`
	ID     int64           `json:"id"`
}

type coreRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewCoreRPCClient creates a new Bitcoin Core RPC client.
// url should include auth credentials: http://user:pass@host:port
func NewCoreRPCClient(url string) *CoreRPCClient {
	return &CoreRPCClient{
		url: url,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CoreRPCClient) call(method string, params []any, result any) error {
	id := c.nextID.Add(1)
	req := coreRPCRequest{
		JSONRPC: "1.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck,gosec

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var rpcResp coreRPCResponse
	if err := json.Unmarshal(respBody, &rpcResp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("RPC error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if result != nil {
		if err := json.Unmarshal(rpcResp.Result, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}

// GetBlockchainInfo returns basic blockchain information to verify connectivity
func (c *CoreRPCClient) GetBlockchainInfo() (json.RawMessage, error) {
	var result json.RawMessage
	if err := c.call("getblockchaininfo", nil, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// ScanTxOutSet scans the UTXO set for outputs matching the given descriptors.
// descriptors should be like ["addr(1A1zP1...)", "addr(1BvBM...)"]
func (c *CoreRPCClient) ScanTxOutSet(descriptors []string) (*ScanTxOutSetResult, error) {
	scanObjects := make([]map[string]string, len(descriptors))
	for i, desc := range descriptors {
		scanObjects[i] = map[string]string{"desc": desc}
	}

	var result ScanTxOutSetResult
	if err := c.call("scantxoutset", []any{"start", scanObjects}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ScanTxOutSetResult is the result of scantxoutset
type ScanTxOutSetResult struct {
	Success  bool             `json:"success"`
	TxOuts   int              `json:"txouts"`
	Height   int64            `json:"height"`
	Unspents []CoreRPCUnspent `json:"unspents"`
}

// CoreRPCUnspent is a single unspent output from scantxoutset
type CoreRPCUnspent struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Desc         string  `json:"desc"`
	Amount       float64 `json:"amount"`
	Height       int64   `json:"height"`
}

// SendRawTransaction broadcasts a hex-encoded transaction
func (c *CoreRPCClient) SendRawTransaction(hexTx string) (string, error) {
	var txid string
	if err := c.call("sendrawtransaction", []any{hexTx}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// EstimateSmartFee returns the estimated fee rate in BTC/kB
func (c *CoreRPCClient) EstimateSmartFee(confTarget int) (*EstimateSmartFeeResult, error) {
	var result EstimateSmartFeeResult
	if err := c.call("estimatesmartfee", []any{confTarget}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// EstimateSmartFeeResult is the result of estimatesmartfee
type EstimateSmartFeeResult struct {
	FeeRate float64  `json:"feerate"` // BTC/kB
	Errors  []string `json:"errors"`
	Blocks  int      `json:"blocks"`
}

// GetRawTransaction returns a decoded transaction
func (c *CoreRPCClient) GetRawTransaction(txid string) (*RawTransactionResult, error) {
	var result RawTransactionResult
	if err := c.call("getrawtransaction", []any{txid, true}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// RawTransactionResult is the result of getrawtransaction with verbose=true
type RawTransactionResult struct {
	TxID          string        `json:"txid"`
	Size          int           `json:"size"`
	VSize         int           `json:"vsize"`
	Confirmations int64         `json:"confirmations"`
	Time          int64         `json:"time"`
	BlockTime     int64         `json:"blocktime"`
	Vin           []CoreRPCVin  `json:"vin"`
	Vout          []CoreRPCVout `json:"vout"`
}

// CoreRPCVin is a transaction input from getrawtransaction
type CoreRPCVin struct {
	TxID string `json:"txid"`
	Vout uint32 `json:"vout"`
}

// CoreRPCVout is a transaction output from getrawtransaction
type CoreRPCVout struct {
	Value        float64             `json:"value"`
	N            uint32              `json:"n"`
	ScriptPubKey CoreRPCScriptPubKey `json:"scriptPubKey"`
}

// CoreRPCScriptPubKey is the scriptPubKey info from a vout
type CoreRPCScriptPubKey struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}
