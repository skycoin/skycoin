package webrpc

import (
	"reflect"
	"testing"
)

func Test_getAddrUxOutsHandler(t *testing.T) {
	m := NewGatewayerMock()
	// GetRecvUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error)
	// GetSpentUxOutOfAddr(addr cipher.Address) ([]*historydb.UxOut, error)
	mockData := []struct {
		method string
		addr   string
		ret    string
		err    error
	}{
		{
			"GetRecvUxOutOfAddr",
			"2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
			`[
                {
                    "uxid": "cc816392cef53a5b75f91bc3fb8155f133907c8ce7f6540507ab30e0456aec3e",
                    "time": 1482042899,
                    "src_block_seq": 562,
                    "src_tx": "ec9e876d4bb33beec203de769b0d3b23de21052de0e4df06b1444bcfec773c46",
                    "owner_address": "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
                    "coins": 1000000,
                    "hours": 0,
                    "spent_block_seq": 563,
                    "spent_tx": "31a21a4dd8331ce68756ddbb21f2c66279d5f5526e936f550e49e29b840ac1ff"
                }
            ]`,
			nil,
		},
        {
            "GetSpentUxOutOfAddr",
            "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
            `[
                {
                    "uxid": "cc816392cef53a5b75f91bc3fb8155f133907c8ce7f6540507ab30e0456aec3e",
                    "time": 1482042899,
                    "src_block_seq": 562,
                    "src_tx": "ec9e876d4bb33beec203de769b0d3b23de21052de0e4df06b1444bcfec773c46",
                    "owner_address": "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
                    "coins": 1000000,
                    "hours": 0,
                    "spent_block_seq": 563,
                    "spent_tx": "31a21a4dd8331ce68756ddbb21f2c66279d5f5526e936f550e49e29b840ac1ff"
                }
            ]`,
            nil,
        },
	}

	for _, d := range mockData {
		m.On(d.method, d.addr).Return(d.ret, d.err)
	}

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
				gateway: m,
			},
			makeSuccessResponse("1", )
		},
	}
	for _, tt := range tests {
		if got := getAddrUxOutsHandler(tt.args.req, tt.args.gateway); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getAddrUxOutsHandler() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
