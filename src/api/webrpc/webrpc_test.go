package webrpc

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/assert"
)

type fakeGateway struct {
}

func (fg fakeGateway) GetLastBlocks(num uint64) *visor.ReadableBlocks {
	return nil
}

func setup() (*rpcHandler, func()) {
	c := make(chan struct{})
	f := func() {
		close(c)
	}

	return makeRPC(1, 1, &fakeGateway{}, c), f
}

func TestHTTPMethod(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()
	d, err := json.Marshal(Request{})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/webrpc", bytes.NewBuffer(d))
	w := httptest.NewRecorder()
	rpc.Handler(w, r)

	var res Response
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.Error.Code, errCodeInvalidRequest)
	assert.Equal(t, res.Error.Message, "only support http POST")
}

func TestInvalidJsonRpc(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()

	d, err := json.Marshal(Request{
		ID:      "1",
		Jsonrpc: "1.0",
		Method:  "get_status",
	})

	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/webrpc", bytes.NewBuffer(d))
	w := httptest.NewRecorder()
	rpc.Handler(w, r)

	var res Response
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.Error, &RPCError{
		Code:    errCodeInvalidParams,
		Message: errMsgInvalidJsonrpc,
	})
}

func TestGetStatus(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()

	tests := []struct {
		Req          Request
		WantResponse Response
	}{
		{
			Request{
				ID:      "1",
				Method:  "get_status",
				Jsonrpc: jsonRPC,
			},
			Response{
				ID:      "1",
				Jsonrpc: jsonRPC,
				Result:  `{"running": true}`,
			},
		},
	}

	for _, tt := range tests {
		d, err := json.Marshal(tt.Req)
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest("POST", "/webrpc", bytes.NewBuffer(d))
		w := httptest.NewRecorder()
		rpc.Handler(w, r)
		var res Response
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}
		assert.EqualValues(t, res, tt.WantResponse)
	}
}
