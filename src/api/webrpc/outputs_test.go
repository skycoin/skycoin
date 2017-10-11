package webrpc

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/require"
)

const outputStr = `{
       "outputs":
			{
				"head_outputs": [
					{
						"hash": "ca02361ef6d658cac5b5aadcb502b4b6046d1403e0f4b1f16b35c06a3f27e3df",
						"src_tx": "e00196267e879c76215ccb93d046bd248e2bc5accad93d246ba43c71c42ff44a",
						"address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
						"coins": "4",
						"hours": 0
					},
					{
						"hash": "22f489be1a2f87ed826c183b516bd10f1703c6591643796f48630ba97db3b16c",
						"src_tx": "fe50714012b29b3ffe5bc2f8e12a95af35004513d61e329e33b9b2a964ae2924",
						"address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
						"coins": "1",
						"hours": 0
					},
					{
						"hash": "f34f2f08c0a9bab56920b4ef946c0cb3ce31bbd641e44b23c5c1c39a14c86c86",
						"src_tx": "059197c06b3a236c550bec377e26401c50ee6480b51206b6f3899ece55209b50",
						"address": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
						"coins": "53",
						"hours": 1
					},
					{
						"hash": "86c43aeaa420e17843fee51ec28275726c6422f6bb0f844e70c552d65dd63df8",
						"src_tx": "bb35c6b277f432c6cf13d4a6b36d64f75cc405bc2b864aad718e53a6cbbd9105",
						"address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
						"coins": "1",
						"hours": 0
					}
				],
				"outgoing_outputs": [],
				"incoming_outputs": []
			}
    }`

func decodeOutputStr(str string) visor.ReadableOutputSet {
	outs := OutputsResult{}
	if err := json.NewDecoder(strings.NewReader(outputStr)).Decode(&outs); err != nil {
		panic(err)
	}
	return outs.Outputs
}

func filterOut(outs []coin.UxOut, f func(out coin.UxOut) bool) visor.ReadableOutputSet {
	os := []coin.UxOut{}
	for _, o := range outs {
		if f(o) {
			os = append(os, o)
		}
	}

	headOuts, err := visor.NewReadableOutputs(os)
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
			makeErrorResponse(errCodeInvalidParams, "invalid address: fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"),
		},
		{
			"invalid params: empty addresses",
			args{},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"single address",
			args{
				addrs:   []string{addrs[0].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(uxouts[:], func(out coin.UxOut) bool {
				return out.Body.Address == addrs[0]
			})}),
		},
		{
			"multiple addresses",
			args{
				addrs:   []string{addrs[0].String(), addrs[1].String()},
				gateway: &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(uxouts, func(out coin.UxOut) bool {
				return out.Body.Address == addrs[0] || out.Body.Address == addrs[1]
			})}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := json.Marshal(tt.args.addrs)
			fmt.Println("param:", string(params))
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
