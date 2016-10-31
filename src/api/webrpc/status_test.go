package webrpc

import (
	"reflect"
	"testing"
)

func Test_getStatusHandler(t *testing.T) {
	type args struct {
		req Request
		in1 Gatewayer
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		// TODO: Add test cases.
		{
			"normal",
			args{
				req: Request{
					ID:      "1",
					Method:  "get_status",
					Jsonrpc: jsonRPC,
				},
				in1: &fakeGateway{},
			},
			Response{
				ID:      ptrString("1"),
				Jsonrpc: jsonRPC,
				Result:  ptrString(`{"running": true}`),
			},
		},
	}
	for _, tt := range tests {
		if got := getStatusHandler(tt.args.req, tt.args.in1); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getStatusHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// func TestGetStatus(t *testing.T) {
// 	rpc, teardown := setup()
// 	defer teardown()

// 	tests := []struct {
// 		Req          Request
// 		WantResponse Response
// 	}{
// 		{
// 			Request{
// 				ID:      "1",
// 				Method:  "get_status",
// 				Jsonrpc: jsonRPC,
// 			},
// 			Response{
// 				ID:      "1",
// 				Jsonrpc: jsonRPC,
// 				Result:  `{"running": true}`,
// 			},
// 		},
// 		{
// 			Request{
// 				ID:      "1",
// 				Method:  "invalid_method",
// 				Jsonrpc: jsonRPC,
// 			},
// 			makeErrorResponse("", &RPCError{
// 				Code:    errCodeMethodNotFound,
// 				Message: errMsgMethodNotFound,
// 			}),
// 		},
// 	}

// 	for _, tt := range tests {
// 		d, err := json.Marshal(tt.Req)
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		r := httptest.NewRequest("POST", "/webrpc", bytes.NewBuffer(d))
// 		w := httptest.NewRecorder()
// 		rpc.Handler(w, r)
// 		var res Response
// 		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
// 			t.Fatal(err)
// 		}
// 		assert.EqualValues(t, res, tt.WantResponse)
// 	}
// }
