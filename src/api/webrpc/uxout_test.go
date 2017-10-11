package webrpc

import (
	"errors"
	"testing"

	"encoding/json"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/stretchr/testify/require"
)

func Test_getAddrUxOutsHandler(t *testing.T) {
	m, mockData := newUxOutMock()
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
					Params:  []byte(`["2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT"]`),
				},
				gateway: m,
			},
			makeSuccessResponse("1", []AddrUxoutResult{{
				Address: "2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT",
				UxOuts:  mockData("2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT")}}),
		},
		{
			"internal server error",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInternalError, errMsgInternalError),
		},
		{
			"invalid address length",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4BBB"]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, "Invalid address length"),
		},
		{
			"invalid address version",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`["111X5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, "Invalid address length"),
		},
		{
			"invalid params",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`[]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"decode params error",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_address_uxouts",
					Params:  []byte(`[invalid params]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAddrUxOutsHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}

func decodeUxJSON(s string) []*historydb.UxOutJSON {
	v := []*historydb.UxOutJSON{}
	if err := json.NewDecoder(strings.NewReader(s)).Decode(&v); err != nil {
		panic(err)
	}
	return v
}

func newUxOutMock() (*GatewayerMock, func(addr string) []*historydb.UxOutJSON) {
	m := NewGatewayerMock()
	uxoutJSON := `[
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
            ]`
	mockData := map[string]struct {
		ret []*historydb.UxOutJSON
		err error
	}{
		"2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT": {
			decodeUxJSON(uxoutJSON),
			nil,
		},
		"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B": {
			nil,
			errors.New("internal server error"),
		},
	}

	for addr, d := range mockData {
		a := cipher.MustDecodeBase58Address(addr)
		m.On("GetAddrUxOuts", a).Return(d.ret, d.err)
	}

	f := func(addr string) []*historydb.UxOutJSON {
		return mockData[addr].ret
	}
	return m, f
}
