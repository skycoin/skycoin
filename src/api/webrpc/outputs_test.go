package webrpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func filterOut(headTime uint64, outs []coin.UxOut, f func(out coin.UxOut) bool) visor.ReadableOutputSet {
	os := []coin.UxOut{}
	for _, o := range outs {
		if f(o) {
			os = append(os, o)
		}
	}

	headOuts, err := visor.NewReadableOutputs(headTime, os)
	if err != nil {
		panic(err)
	}
	return visor.ReadableOutputSet{
		HeadOutputs: headOuts,
	}
}

func Test_getOutputsHandler(t *testing.T) {
	uxouts := make([]coin.UxOut, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxouts[i] = coin.UxOut{}
		uxouts[i].Body.Address = addrs[i]
	}

	headTime := uint64(time.Now().UTC().Unix())

	type args struct {
		addrs   []string
		gateway Gatewayer
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		// TODO: Add test cases.
		{
			"invalid address",
			args{
				addrs: []string{"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"},
			},
			MakeErrorResponse(ErrCodeInvalidParams, "invalid address: fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"),
		},
		{
			"invalid params: empty addresses",
			args{},
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
		{
			"single address",
			args{
				addrs:   []string{addrs[0].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(headTime, uxouts[:], func(out coin.UxOut) bool {
				return out.Body.Address == addrs[0]
			})}),
		},
		{
			"multiple addresses",
			args{
				addrs:   []string{addrs[0].String(), addrs[1].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(headTime, uxouts, func(out coin.UxOut) bool {
				return out.Body.Address == addrs[0] || out.Body.Address == addrs[1]
			})}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := json.Marshal(tt.args.addrs)
			require.NoError(t, err)
			req := Request{
				ID:      "1",
				Jsonrpc: jsonRPC,
				Method:  "get_outputs",
				Params:  params,
			}

			got := getOutputsHandler(req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}
