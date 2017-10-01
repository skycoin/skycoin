package webrpc

import (
	"encoding/json"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
	"github.com/stretchr/testify/require"
)

var blockString = `{
    "blocks": [
        {
            "header": {
                "version": 0,
                "timestamp": 1477295242,
                "seq": 454,
                "fee": 20732,
                "prev_hash": "f680fe1f068a1cd5c3ef9194f91a9bc3cacffbcae4a32359a3c014da4ef7516f",
                "hash": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a"
            },
            "body": {
                "txns": [
                    {
                        "length": 608,
                        "type": 0,
                        "txid": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a",
                        "inner_hash": "37f1111bd83d9c995b9e48511bd52de3b0e440dccbf6d2cfd41dee31a10f1aa4",
                        "sigs": [
                            "ef0b8e1465557e6f21cb2bfad17136188f0b9bd54bba3db76c3488eb8bc900bc7662e3fe162dd6c236d9e52a7051a2133855081a91f6c1a63e1fce2ae9e3820e00",
                            "800323c8c22a2c078cecdfad35210902f91af6f97f0c63fe324e0a9c2159e9356f2fbbfff589edea5a5c24453ef5fc0cd5929f24bebee28e37057acd6d42f3d700",
                            "ca6a6ef5f5fb67490d88ddeeee5e5d11055246613b03e7ed2ad5cc82d01077d262e2da56560083928f5389580ae29500644719cf0e82a5bf065cecbed857598400",
                            "78ddc117607159c7b4c76fc91deace72425f21f2df5918d44d19a377da68cc610668c335c84e2bb7a8f16cd4f9431e900585fc0a3f1024b722b974fcef59dfd500",
                            "4c484d44072e23e97a437deb03a85e3f6eca0bd8875031efe833e3c700fc17f91491969b9864b56c280ef8a68d18dd728b211ce1d46fe477fe3104d73d55ad6501"
                        ],
                        "inputs": [
                            "4bd7c68ecf3039c2b2d8c26a5e2983e20cf53b6d62b099e7786546b3c3f600f9",
                            "f9e39908677cae43832e1ead2514e01eaae48c9a3614a97970f381187ee6c4b1",
                            "7e8ac23a2422b4666ff45192fe36b1bd05f1285cf74e077ac92cabf5a7c1100e",
                            "b3606a4f115d4161e1c8206f4fb5ac0e91551c40d0ee6fe40c86040d2faacac0",
                            "305f1983f5b630bba27e2777c229c725b6b57f37a6ddee138d1d82ae56311909"
                        ],
                        "outputs": [
                            {
                                "uxid": "574d7e5afaefe4ee7e0adf6ce1971d979f038adc8ebbd35771b2c19b0bad7e3d",
                                "dst": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
                                "coins": "1",
                                "hours": 3455
                            },
                            {
                                "uxid": "6d8a9c89177ce5e9d3b4b59fff67c00f0471fdebdfbb368377841b03fc7d688b",
                                "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                                "coins": "5",
                                "hours": 3455
                            }
                        ]
                    }
                ]
            }
        }
    ]
}`

var emptyBlockString = `{
							"blocks":[]
						}`

func decodeBlock(str string) *visor.ReadableBlocks {
	var blocks visor.ReadableBlocks
	if err := json.Unmarshal([]byte(str), &blocks); err != nil {
		panic(err)
	}
	return &blocks
}

func Test_getLastBlocksHandler(t *testing.T) {
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
					Method:  "get_lastblocks",
					Params:  []byte("[1]"),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", decodeBlock(blockString)),
		},
		{
			"invalid params: num value",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_lastblocks",
					Params:  []byte(`[1a]`), // invalid params
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"invalid params: no num value",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_lastblocks",
					Params:  []byte(`{"foo": 1}`), // invalid params
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"invalid params: empty params",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_lastblocks",
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"invalid params: more than one param",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_lastblocks",
					Params:  []byte("[1,2]"),
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getLastBlocksHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_getBlocksHandler(t *testing.T) {
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
					Method:  "get_blocks",
					Params:  []byte("[0, 1]"),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", decodeBlock(blockString)),
		},
		{
			"invalid params: lost end",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks",
					Params:  []byte("[0]"),
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"invalid params:lost start",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks",
					Params:  []byte("[1]"),
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"invalid params: start = abc",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks",
					Params:  []byte(`{ "start": "abc"}`),
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"empty params",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks",
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"start > end",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks",
					Params:  []byte(`[2, 1]`),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlocksHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_getBlocksBySeqHandler(t *testing.T) {
	m := NewGatewayerMock()
	m.On("GetBlocksInDepth", []uint64{454}).Return(decodeBlock(blockString), nil)
	m.On("GetBlocksInDepth", []uint64{1000}).Return(decodeBlock(emptyBlockString), nil)

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
					Method:  "get_blocks_by_seq",
					Params:  []byte(`[454]`),
				},
				gateway: m,
			},
			makeSuccessResponse("1", decodeBlock(blockString)),
		},
		{
			"none exist seq",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks_by_seq",
					Params:  []byte(`[1000]`),
				},
				gateway: m,
			},
			makeSuccessResponse("1", decodeBlock(emptyBlockString)),
		},
		{
			"invalid request param",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks_by_seq",
					Params:  []byte(`["454"]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
		{
			"empty param",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_blocks_by_seq",
					Params:  []byte(`[]`),
				},
				gateway: m,
			},
			makeErrorResponse(errCodeInvalidParams, "empty params"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlocksBySeqHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}
