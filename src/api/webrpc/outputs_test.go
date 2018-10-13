package webrpc

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func filterOut(t *testing.T, headTime uint64, outs []coin.UxOut, f func(out coin.UxOut) bool) readable.UnspentOutputsSummary {
	os := []coin.UxOut{}
	for _, o := range outs {
		if f(o) {
			os = append(os, o)
		}
	}

	vOuts, err := visor.NewUnspentOutputs(os, headTime)
	require.NoError(t, err)

	headOuts, err := readable.NewUnspentOutputs(vOuts)
	require.NoError(t, err)

	return readable.UnspentOutputsSummary{
		Head: readable.BlockHeader{
			Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
			PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
			BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
			UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
		},
		HeadOutputs:     headOuts,
		IncomingOutputs: readable.UnspentOutputs{},
		OutgoingOutputs: readable.UnspentOutputs{},
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
		{
			name: "invalid address",
			args: args{
				addrs: []string{"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"},
			},
			want: MakeErrorResponse(ErrCodeInvalidParams, "invalid address: fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"),
		},
		{
			name: "invalid params: empty addresses",
			args: args{},
			want: MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
		{
			name: "single address",
			args: args{
				addrs:   []string{addrs[0].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			want: makeSuccessResponse("1", OutputsResult{filterOut(t, headTime, uxouts[:], func(out coin.UxOut) bool {
				return out.Body.Address == addrs[0]
			})}),
		},
		{
			name: "multiple addresses",
			args: args{
				addrs:   []string{addrs[0].String(), addrs[1].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			want: makeSuccessResponse("1", OutputsResult{filterOut(t, headTime, uxouts, func(out coin.UxOut) bool {
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
