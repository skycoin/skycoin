package webrpc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
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

func filterOut(outs visor.ReadableOutputSet, f func(out visor.ReadableOutput) bool) visor.ReadableOutputSet {
	headOuts := []visor.ReadableOutput{}
	outgoingOuts := []visor.ReadableOutput{}
	incomingOuts := []visor.ReadableOutput{}

	for _, o := range outs.HeadOutputs {
		if f(o) {
			headOuts = append(headOuts, o)
		}
	}

	for _, o := range outs.OutgoingOutputs {
		if f(o) {
			outgoingOuts = append(outgoingOuts, o)
		}
	}

	for _, o := range outs.IncomingOutputs {
		if f(o) {
			incomingOuts = append(incomingOuts, o)
		}
	}

	return visor.ReadableOutputSet{
		HeadOutputs:     headOuts,
		OutgoingOutputs: outgoingOuts,
		IncomingOutputs: incomingOuts,
	}
}

func Test_getOutputsHandler(t *testing.T) {
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
			"invalid address",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_outputs",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"]`),
				},
			},
			makeErrorResponse(errCodeInvalidParams, "invalid address: fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"),
		},
		{
			"invalid params: empty addresses",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_outputs",
				},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"normal, single address",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_outputs",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"]`),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(decodeOutputStr(outputStr), func(out visor.ReadableOutput) bool {
				return out.Address == "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"
			})}),
		},
		{
			"normal, multiple addresses",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_outputs",
					Params:  []byte(`["fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW"]`),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", OutputsResult{filterOut(decodeOutputStr(outputStr), func(out visor.ReadableOutput) bool {
				return out.Address == "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B" || out.Address == "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW"
			})}),
		},
	}
	for _, tt := range tests {
		if got := getOutputsHandler(tt.args.req, tt.args.gateway); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. getOutputsHandler() = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}
