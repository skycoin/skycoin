package webrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/stretchr/testify/assert"
)

func setup() (*WebRPC, func()) {
	c := make(chan struct{})

	rpc, err := New(
		"0.0.0.0:8081",
		ChanBuffSize(1),
		ThreadNum(1),
		Gateway(&fakeGateway{}),
		Quit(c))
	if err != nil {
		panic(err)
	}

	return rpc, func() {
		rpc.Shutdown()
	}
}

type fakeGateway struct {
	transactions    map[string]string
	injectRawTxMap  map[string]bool // key: transacion hash, value indicates whether the injectTransaction should return error.
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
		v.IncommingOutputs = f(v.IncommingOutputs)
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

func TestNewWebRPC(t *testing.T) {
	rpc1, err := New("0.0.0.0:8080", ChanBuffSize(1), ThreadNum(1), Gateway(&fakeGateway{}), Quit(make(chan struct{})))
	assert.Nil(t, err)
	assert.NotNil(t, rpc1.mux)
	assert.NotNil(t, rpc1.handlers)
	assert.NotNil(t, rpc1.gateway)
}

func Test_rpcHandler_HandlerFunc(t *testing.T) {
	rpc, err := New("0.0.0.0:8080", ChanBuffSize(1), ThreadNum(1), Gateway(&fakeGateway{}), Quit(make(chan struct{})))
	assert.Nil(t, err)
	rpc.HandleFunc("get_status", getStatusHandler)
	err = rpc.HandleFunc("get_status", getStatusHandler)
	assert.NotNil(t, err)
}

func Test_rpcHandler_Handler(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()
	go rpc.Run()

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
		assert.Equal(t, res, tt.want)
	}

}
