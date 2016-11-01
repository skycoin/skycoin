package webrpc

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
)

var outputStr = `{
       "outputs": [
            {
                "txid": "0ccc33d0f771b60c95e660745a1b1a92f9dc29f50bdf28d81755cdd3466caf41",
                "src_tx": "5d4f0397d0f7278d70a81623cc5e3899e65a8b9063d61e15a908657189855082",
                "address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
                "coins": "6",
                "hours": 430
            },
            {
                "txid": "574d7e5afaefe4ee7e0adf6ce1971d979f038adc8ebbd35771b2c19b0bad7e3d",
                "src_tx": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a",
                "address": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
                "coins": "1",
                "hours": 3455
            },
            {
                "txid": "6d8a9c89177ce5e9d3b4b59fff67c00f0471fdebdfbb368377841b03fc7d688b",
                "src_tx": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a",
                "address": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                "coins": "5",
                "hours": 3455
            }
        ]
    }`

func decodeOutputStr(str string) []visor.ReadableOutput {
	outs := OutputsResult{}
	if err := json.NewDecoder(strings.NewReader(outputStr)).Decode(&outs); err != nil {
		panic(err)
	}
	return outs.Outputs
}

func filterOut(outs []visor.ReadableOutput, f func(out visor.ReadableOutput) bool) []visor.ReadableOutput {
	ret := []visor.ReadableOutput{}
	for _, o := range outs {
		if f(o) {
			ret = append(ret, o)
		}
	}
	return ret
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
					Params: map[string]string{
						"addresses": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C",
					},
				},
			},
			makeErrorResponse(errCodeInvalidParams, "invalid address: fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4C"),
		},
		{
			"normal, single address",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_outputs",
					Params: map[string]string{
						"addresses": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
					},
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
					Params: map[string]string{
						"addresses": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B, cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
					},
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
