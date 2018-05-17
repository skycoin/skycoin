package webrpc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func TestClientGetUnspentOutputs(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	headTime := uint64(time.Now().UTC().Unix())
	uxouts := make([]coin.UxOut, 5)
	addrs := make([]cipher.Address, 5)
	rbOutputs := make(visor.ReadableOutputs, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxouts[i] = coin.UxOut{}
		uxouts[i].Body.Address = addrs[i]
		rbOut, err := visor.NewReadableOutput(headTime, uxouts[i])
		require.NoError(t, err)
		rbOutputs[i] = rbOut
	}

	s.Gateway = &fakeGateway{
		uxouts: uxouts,
	}

	cases := []struct {
		name   string
		params []string
		errMsg string
	}{
		{
			name:   "valid, multiple addresses",
			params: []string{addrs[0].String(), addrs[1].String()},
		},
		{
			name:   "invalid addresses",
			params: []string{"invalid-address-foo"},
			errMsg: "invalid address: invalid-address-foo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rpcReq, err := NewRequest("get_outputs", tc.params, "1")
			require.NoError(t, err)

			body, err := json.Marshal(rpcReq)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var resp Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.errMsg != "" {
				require.NotNil(t, resp.Error)
				require.NotEmpty(t, resp.Error.Code)
				require.Equal(t, tc.errMsg, resp.Error.Message)
				return
			}

			require.Nil(t, resp.Error)

			var outputs OutputsResult
			err = json.Unmarshal(resp.Result, &outputs)
			require.NoError(t, err)

			require.Len(t, outputs.Outputs.HeadOutputs, 2)
			require.Len(t, outputs.Outputs.IncomingOutputs, 0)
			require.Len(t, outputs.Outputs.OutgoingOutputs, 0)

			// GetUnspentOutputs sorts outputs by most recent time first, then by hash
			expectedOutputs := rbOutputs[:2]
			sort.Slice(expectedOutputs, func(i, j int) bool {
				if expectedOutputs[i].Time == expectedOutputs[j].Time {
					return strings.Compare(expectedOutputs[i].Hash, expectedOutputs[j].Hash) < 1
				}

				return expectedOutputs[i].Time > expectedOutputs[j].Time
			})

			require.Equal(t, rbOutputs[:2], outputs.Outputs.HeadOutputs)
		})
	}
}

func TestClientInjectTransaction(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	s.Gateway.(*fakeGateway).injectRawTxMap = map[string]bool{
		rawTxID: true,
	}
	require.Empty(t, s.Gateway.(*fakeGateway).injectedTransactions)

	rpcReq, err := NewRequest("inject_transaction", []string{rawTxStr}, "1")
	require.NoError(t, err)

	body, err := json.Marshal(rpcReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	var txidJSON TxIDJson
	err = json.Unmarshal(resp.Result, &txidJSON)
	require.NoError(t, err)
	require.NotEmpty(t, txidJSON.Txid)

	require.Len(t, s.Gateway.(*fakeGateway).injectedTransactions, 1)
	require.Contains(t, s.Gateway.(*fakeGateway).injectedTransactions, rawTxID)
}

func TestClientGetStatus(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	rpcReq, err := NewRequest("get_status", nil, "1")
	require.NoError(t, err)

	body, err := json.Marshal(rpcReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	var result StatusResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	// values derived from hardcoded `blockString`
	require.Equal(t, StatusResult{
		Running:            true,
		BlockNum:           455,
		LastBlockHash:      "",
		TimeSinceLastBlock: "18446744072232256374s",
	}, result)
}

func TestClientGetTransactionByID(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	cases := []struct {
		name                string
		txid                string
		errMsg              string
		gatewayTransactions map[string]string
	}{
		{
			name:   "invalid txn id",
			txid:   "foo",
			errMsg: "invalid transaction hash",
		},
		{
			name:   "valid txn id, but does not exist",
			txid:   rawTxID,
			errMsg: "transaction doesn't exist",
		},
		{
			name: "valid txn id exists",
			txid: rawTxID,
			gatewayTransactions: map[string]string{
				rawTxID: rawTxStr,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewayTransactions != nil {
				s.Gateway.(*fakeGateway).transactions = tc.gatewayTransactions
			}

			rpcReq, err := NewRequest("get_transaction", []string{tc.txid}, "1")
			require.NoError(t, err)

			body, err := json.Marshal(rpcReq)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var resp Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.errMsg != "" {
				require.NotNil(t, resp.Error)
				require.Equal(t, tc.errMsg, resp.Error.Message)
				return
			}

			var txn TxnResult
			err = json.Unmarshal(resp.Result, &txn)
			require.NoError(t, err)

			expectedTxn := decodeRawTransaction(rawTxStr)
			rbTx, err := visor.NewReadableTransaction(expectedTxn)
			require.NoError(t, err)
			require.Equal(t, &daemon.TransactionResult{
				Status:      expectedTxn.Status,
				Time:        0,
				Transaction: *rbTx,
			}, txn.Transaction)
		})
	}
}

func TestClientGetAddressUxOuts(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	cases := []struct {
		name   string
		addr   string
		errMsg string
	}{
		{
			name: "valid address",
			addr: "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
		},
		{
			name:   "invalid address",
			addr:   "foo",
			errMsg: "Invalid address length",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gatewayerMock, mockData := newUxOutMock()
			s.Gateway = gatewayerMock

			rpcReq, err := NewRequest("get_address_uxouts", []string{tc.addr}, "1")
			require.NoError(t, err)

			body, err := json.Marshal(rpcReq)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			var resp Response
			err = json.Unmarshal(rr.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.errMsg != "" {
				require.NotNil(t, resp.Error)
				require.Equal(t, tc.errMsg, resp.Error.Message)
				return
			}

			var uxouts []AddrUxoutResult
			err = json.Unmarshal(resp.Result, &uxouts)
			require.NoError(t, err)

			require.Equal(t, []AddrUxoutResult{{
				Address: tc.addr,
				UxOuts:  mockData(tc.addr),
			}}, uxouts)
		})
	}
}

func TestClientGetBlocks(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	// blockString borrowed from block_test.go
	rpcReq, err := NewRequest("get_blocks", []uint64{0, 1}, "1")
	require.NoError(t, err)

	body, err := json.Marshal(rpcReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	var blocks visor.ReadableBlocks
	err = json.Unmarshal(resp.Result, &blocks)
	require.NoError(t, err)

	require.NotNil(t, blocks.Blocks)
	require.Equal(t, decodeBlock(blockString), &blocks)
}

func TestClientGetBlocksBySeq(t *testing.T) {
	s := setupWebRPC(t)

	gatewayerMock := NewGatewayerMock()
	s.Gateway = gatewayerMock
	gatewayerMock.On("GetBlocksInDepth", []uint64{454}).Return(decodeBlock(blockString), nil)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	// blockString and seq borrowed from block_test.go
	var seq uint64 = 454
	rpcReq, err := NewRequest("get_blocks_by_seq", []uint64{seq}, "1")
	require.NoError(t, err)

	body, err := json.Marshal(rpcReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	var blocks visor.ReadableBlocks
	err = json.Unmarshal(resp.Result, &blocks)
	require.NoError(t, err)

	require.NotNil(t, blocks.Blocks)
	require.Equal(t, decodeBlock(blockString), &blocks)
}

func TestClientGetLastBlocks(t *testing.T) {
	s := setupWebRPC(t)

	mux := http.NewServeMux()
	mux.Handle("/api/v1/webrpc", http.HandlerFunc(s.Handler))

	var n uint64 = 1
	rpcReq, err := NewRequest("get_lastblocks", []uint64{n}, "1")
	require.NoError(t, err)

	body, err := json.Marshal(rpcReq)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/api/v1/webrpc", bytes.NewReader(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	var blocks visor.ReadableBlocks
	err = json.Unmarshal(resp.Result, &blocks)
	require.NoError(t, err)

	require.Len(t, blocks.Blocks, 1)
	require.Equal(t, decodeBlock(blockString), &blocks)
}
