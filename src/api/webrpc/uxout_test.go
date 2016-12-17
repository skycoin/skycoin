package webrpc

import (
	"reflect"
	"testing"
)

func Test_getAddrUxOutsHandler(t *testing.T) {
	type args struct {
		req     Request
		gateway Gatewayer
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
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"]`),
				},
				gateway: &fakeGateway{},
			},
			// makeSuccessResponse("1", )
		},
	}
	for _, tt := range tests {
		if got := getAddrUxOutsHandler(tt.args.req, tt.args.gateway); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getAddrUxOutsHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
