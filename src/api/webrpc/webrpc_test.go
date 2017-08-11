package webrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

const testWebRPCAddr = "127.0.0.1:8081"

func setupWebRPC(t *testing.T) *WebRPC {
	rpc, err := New(testWebRPCAddr, &fakeGateway{})
	require.NoError(t, err)
	rpc.WorkerNum = 1
	rpc.ChanBuffSize = 2
	return rpc
}

type fakeGateway struct {
	transactions    map[string]string
	injectRawTxMap  map[string]bool // key: transaction hash, value indicates whether the injectTransaction should return error.
	addrRecvUxOuts  []*historydb.UxOut
	addrSpentUxOUts []*historydb.UxOut
}

func (fg fakeGateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	var blocks visor.ReadableBlocks
	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		panic(err)
	}

	return &blocks
}

func (fg fakeGateway) GetBlocks(start, end uint64) *visor.ReadableBlocks {
	var blocks visor.ReadableBlocks
	if start > end {
		return &blocks
	}

	if err := json.Unmarshal([]byte(blockString), &blocks); err != nil {
		panic(err)
	}

	return &blocks
}

func (fg fakeGateway) GetBlocksInDepth(vs []uint64) *visor.ReadableBlocks {
	return nil
}

func (fg fakeGateway) GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error) {
	v := decodeOutputStr(outputStr)
	for _, f := range filters {
		v.HeadOutputs = f(v.HeadOutputs)
		v.OutgoingOutputs = f(v.OutgoingOutputs)
		v.IncomingOutputs = f(v.IncomingOutputs)
	}
	return v, nil
}

func (fg fakeGateway) GetTransaction(txid cipher.SHA256) (*visor.Transaction, error) {
	str, ok := fg.transactions[txid.Hex()]
	if ok {
		return decodeRawTransaction(str), nil
	}
	return nil, nil
}

func (fg fakeGateway) InjectTransaction(txn coin.Transaction) error {
	if _, v := fg.injectRawTxMap[txn.Hash().Hex()]; v {
		return nil
	}

	return errors.New("inject transaction failed")
}

func (fg fakeGateway) GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error) {
	return nil, nil
}

func (fg fakeGateway) GetTimeNow() uint64 {
	return 0
}

func Test_rpcHandler_HandlerFunc(t *testing.T) {
	rpc := setupWebRPC(t)
	rpc.HandleFunc("get_status", getStatusHandler)
	err := rpc.HandleFunc("get_status", getStatusHandler)
	require.Error(t, err)
}

func Test_rpcHandler_Handler(t *testing.T) {
	rpc := setupWebRPC(t)
	errC := make(chan error, 1)
	go func() {
		errC <- rpc.Run()
	}()
	defer func() {
		rpc.Shutdown()
		require.NoError(t, <-errC)
	}()

	time.Sleep(50 * time.Millisecond)

	type args struct {
		httpMethod string
		req        Request
	}

	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			"http GET",
			args{
				httpMethod: "GET",
				req:        Request{},
			},
			Response{
				Jsonrpc: jsonRPC,
				Error: &RPCError{
					Code:    errCodeInvalidRequest,
					Message: errMsgNotPost,
				},
			},
		},
		{
			"invalid jsonrpc",
			args{
				httpMethod: "POST",
				req: Request{
					ID:      "1",
					Jsonrpc: "1.0",
					Method:  "get_status",
				},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidJsonrpc),
		},
	}

	for _, tt := range tests {
		d, err := json.Marshal(tt.args.req)
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest(tt.args.httpMethod, "/webrpc", bytes.NewBuffer(d))
		w := httptest.NewRecorder()
		rpc.Handler(w, r)
		var res Response
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, res, tt.want)
	}
}
