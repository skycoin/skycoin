package webrpc

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newRPCHandler(t *testing.T) {
	rpc1 := newRPCHandler(1, 1, &fakeGateway{}, make(chan struct{}))
	assert.NotNil(t, rpc1.mux)
	assert.NotNil(t, rpc1.handlers)
	assert.NotNil(t, rpc1.gateway)

	assert.Panics(t, func() {
		newRPCHandler(1, 0, &fakeGateway{}, make(chan struct{}))
	})
}

func Test_makeJob(t *testing.T) {
	job := makeJob(Request{})
	assert.NotNil(t, job.ResC)
}

func Test_rpcHandler_HandlerFunc(t *testing.T) {
	rpc := newRPCHandler(1, 1, &fakeGateway{}, make(chan struct{}))
	rpc.HandlerFunc("get_status", getStatusHandler)
	assert.Panics(t, func() {
		rpc.HandlerFunc("get_status", getStatusHandler)
	})
}

func Test_rpcHandler_Handler(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()

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
