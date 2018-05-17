package webrpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/historydb"
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
			MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError),
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
			MakeErrorResponse(ErrCodeInvalidParams, "Invalid address length"),
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
			MakeErrorResponse(ErrCodeInvalidParams, "Invalid address length"),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAddrUxOutsHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}

func newUxOutMock() (*GatewayerMock, func(addr string) []*historydb.UxOutJSON) {
	m := NewGatewayerMock()

	hash, _ := cipher.SHA256FromHex("31a21a4dd8331ce68756ddbb21f2c66279d5f5526e936f550e49e29b840ac1ff")
	address, _ := cipher.DecodeBase58Address("2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT")
	srcTxHash, _ := cipher.SHA256FromHex("ec9e876d4bb33beec203de769b0d3b23de21052de0e4df06b1444bcfec773c46")

	mockData := map[string]struct {
		ret []*historydb.UxOut
		err error
	}{
		"2kmKohJrwURrdcVtDNaWK6hLCNsWWbJhTqT": {
			[]*historydb.UxOut{
				{
					Out: coin.UxOut{
						Head: coin.UxHead{
							Time:  1482042899,
							BkSeq: 562,
						},
						Body: coin.UxBody{
							SrcTransaction: srcTxHash,
							Address:        address,
							Coins:          1000000,
							Hours:          0,
						},
					},
					SpentTxID:     hash,
					SpentBlockSeq: 563,
				},
			},
			nil,
		},
		"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B": {
			nil,
			errors.New("internal server error"),
		},
	}

	for addr, d := range mockData {
		a := cipher.MustDecodeBase58Address(addr)
		m.On("GetAddrUxOuts", []cipher.Address{a}).Return(d.ret, d.err)
	}

	f := func(addr string) []*historydb.UxOutJSON {
		// Convert UxOut to UxOutJson for handler test
		uxouts := mockData[addr].ret
		uxs := make([]*historydb.UxOutJSON, len(uxouts))

		for i, ux := range uxouts {
			uxs[i] = historydb.NewUxOutJSON(ux)
		}

		return uxs
	}
	return m, f
}
