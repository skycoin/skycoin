package webrpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
)

// var blockString = `{
//     "blocks": [
//         {
//             "header": {
//                 "version": 0,
//                 "timestamp": 1477295242,
//                 "seq": 454,
//                 "fee": 20732,
//                 "previous_block_hash": "f680fe1f068a1cd5c3ef9194f91a9bc3cacffbcae4a32359a3c014da4ef7516f",
//                 "tx_body_hash": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a"
//             },
//             "body": {
//                 "txns": [
//                     {
//                         "length": 608,
//                         "type": 0,
//                         "txid": "662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a",
//                         "inner_hash": "37f1111bd83d9c995b9e48511bd52de3b0e440dccbf6d2cfd41dee31a10f1aa4",
//                         "sigs": [
//                             "ef0b8e1465557e6f21cb2bfad17136188f0b9bd54bba3db76c3488eb8bc900bc7662e3fe162dd6c236d9e52a7051a2133855081a91f6c1a63e1fce2ae9e3820e00",
//                             "800323c8c22a2c078cecdfad35210902f91af6f97f0c63fe324e0a9c2159e9356f2fbbfff589edea5a5c24453ef5fc0cd5929f24bebee28e37057acd6d42f3d700",
//                             "ca6a6ef5f5fb67490d88ddeeee5e5d11055246613b03e7ed2ad5cc82d01077d262e2da56560083928f5389580ae29500644719cf0e82a5bf065cecbed857598400",
//                             "78ddc117607159c7b4c76fc91deace72425f21f2df5918d44d19a377da68cc610668c335c84e2bb7a8f16cd4f9431e900585fc0a3f1024b722b974fcef59dfd500",
//                             "4c484d44072e23e97a437deb03a85e3f6eca0bd8875031efe833e3c700fc17f91491969b9864b56c280ef8a68d18dd728b211ce1d46fe477fe3104d73d55ad6501"
//                         ],
//                         "inputs": [
//                             "4bd7c68ecf3039c2b2d8c26a5e2983e20cf53b6d62b099e7786546b3c3f600f9",
//                             "f9e39908677cae43832e1ead2514e01eaae48c9a3614a97970f381187ee6c4b1",
//                             "7e8ac23a2422b4666ff45192fe36b1bd05f1285cf74e077ac92cabf5a7c1100e",
//                             "b3606a4f115d4161e1c8206f4fb5ac0e91551c40d0ee6fe40c86040d2faacac0",
//                             "305f1983f5b630bba27e2777c229c725b6b57f37a6ddee138d1d82ae56311909"
//                         ],
//                         "outputs": [
//                             {
//                                 "uxid": "574d7e5afaefe4ee7e0adf6ce1971d979f038adc8ebbd35771b2c19b0bad7e3d",
//                                 "dst": "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW",
//                                 "coins": "1",
//                                 "hours": 3455
//                             },
//                             {
//                                 "uxid": "6d8a9c89177ce5e9d3b4b59fff67c00f0471fdebdfbb368377841b03fc7d688b",
//                                 "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
//                                 "coins": "5",
//                                 "hours": 3455
//                             }
//                         ]
//                     }
//                 ]
//             }
//         }
//     ]
// }`

var emptyBlockString = `{
							"blocks":[]
						}`

// func decodeBlock(str string) *readable.Blocks {
// 	var blocks readable.Blocks
// 	if err := json.Unmarshal([]byte(str), &blocks); err != nil {
// 		panic(err)
// 	}
// 	return &blocks
// }

// func mustReadableBlocksToSignedBlocks(t *testing.T, rBlocks *readable.Blocks) []coin.SignedBlock {
// 	blocks, err := readableBlocksToSignedBlocks(rBlocks)
// 	require.NoError(t, err)
// 	return blocks
// }

func makeTestBlocks(t *testing.T) []coin.SignedBlock {
	blocks, err := makeTestBlocksWithErr()
	require.NoError(t, err)
	return blocks
}

func makeTestBlocksWithErr() ([]coin.SignedBlock, error) {
	prevHash, err := cipher.SHA256FromHex("f680fe1f068a1cd5c3ef9194f91a9bc3cacffbcae4a32359a3c014da4ef7516f")
	if err != nil {
		return nil, err
	}

	bodyHash, err := cipher.SHA256FromHex("662835cc081e037561e1fe05860fdc4b426f6be562565bfaa8ec91be5675064a")
	if err != nil {
		return nil, err
	}

	innerHash, err := cipher.SHA256FromHex("37f1111bd83d9c995b9e48511bd52de3b0e440dccbf6d2cfd41dee31a10f1aa4")
	if err != nil {
		return nil, err
	}

	sig, err := cipher.SigFromHex("ef0b8e1465557e6f21cb2bfad17136188f0b9bd54bba3db76c3488eb8bc900bc7662e3fe162dd6c236d9e52a7051a2133855081a91f6c1a63e1fce2ae9e3820e00")
	if err != nil {
		return nil, err
	}

	input, err := cipher.SHA256FromHex("4bd7c68ecf3039c2b2d8c26a5e2983e20cf53b6d62b099e7786546b3c3f600f9")
	if err != nil {
		return nil, err
	}

	addr, err := cipher.DecodeBase58Address("cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW")
	if err != nil {
		return nil, err
	}

	transactions := coin.Transactions{
		{
			Length:    608,
			Type:      0,
			InnerHash: innerHash,
			Sigs:      []cipher.Sig{sig},
			In:        []cipher.SHA256{input},
			Out: []coin.TransactionOutput{
				{
					Address: addr,
					Coins:   1e6,
					Hours:   3455,
				},
			},
		},
	}

	return []coin.SignedBlock{
		{
			Block: coin.Block{
				Head: coin.BlockHeader{
					Version:  0,
					Time:     1477295242,
					BkSeq:    454,
					Fee:      20732,
					PrevHash: prevHash,
					BodyHash: bodyHash,
				},
				Body: coin.BlockBody{
					Transactions: transactions,
				},
			},
		},
	}, nil
}

func makeTestReadableBlocks(t *testing.T) *readable.Blocks {
	blocks := makeTestBlocks(t)
	rb, err := readable.NewBlocks(blocks)
	require.NoError(t, err)
	return rb
}

func blockString(t *testing.T) string {
	rb := makeTestReadableBlocks(t)
	data, err := json.Marshal(rb)
	require.NoError(t, err)
	return string(data)
}

// func readableBlocksToSignedBlocks(rBlocks *readable.Blocks) ([]coin.SignedBlock, error) {
// 	blocks := make([]coin.SignedBlock, len(rBlocks.Blocks))
// 	for i, r := range rBlocks.Blocks {
// 		prevHash, err := cipher.SHA256FromHex(r.Head.PreviousBlockHash)
// 		if err != nil {
// 			return nil, err
// 		}

// 		bodyHash, err := cipher.SHA256FromHex(r.Head.BodyHash)
// 		if err != nil {
// 			return nil, err
// 		}

// 		transactions := make(coin.Transactions, len(r.Body.Transactions))
// 		for j, rTxn := range r.Body.Transactions {
// 			innerHash, err := cipher.SHA256FromHex(rTxn.InnerHash)
// 			if err != nil {
// 				return nil, err
// 			}

// 			sigs := make([]cipher.Sig, len(rTxn.Sigs))
// 			for k, sigStr := range rTxn.Sigs {
// 				sig, err := cipher.SigFromHex(sigStr)
// 				if err != nil {
// 					return nil, err
// 				}
// 				sigs[k] = sig
// 			}

// 			inputs := make([]cipher.SHA256, len(rTxn.In))
// 			for k, rIn := range rTxn.In {
// 				input, err := cipher.SHA256FromHex(rIn)
// 				if err != nil {
// 					return nil, err
// 				}
// 				inputs[k] = input
// 			}

// 			outputs := make([]coin.TransactionOutput, len(rTxn.Out))
// 			for k, rOut := range rTxn.Out {
// 				coins, err := droplet.FromString(rOut.Coins)
// 				if err != nil {
// 					return nil, err
// 				}

// 				addr, err := cipher.DecodeBase58Address(rOut.Address)
// 				if err != nil {
// 					return nil, err
// 				}

// 				out := coin.TransactionOutput{
// 					Address: addr,
// 					Coins:   coins,
// 					Hours:   rOut.Hours,
// 				}
// 				outputs[k] = out
// 			}

// 			txn := coin.Transaction{
// 				Length:    rTxn.Length,
// 				Type:      rTxn.Type,
// 				InnerHash: innerHash,
// 				Sigs:      sigs,
// 				In:        inputs,
// 				Out:       outputs,
// 			}

// 			transactions[j] = txn
// 		}

// 		b := coin.SignedBlock{
// 			Block: coin.Block{
// 				Head: coin.BlockHeader{
// 					Version:  r.Head.Version,
// 					Time:     r.Head.Time,
// 					BkSeq:    r.Head.BkSeq,
// 					Fee:      r.Head.Fee,
// 					PrevHash: prevHash,
// 					BodyHash: bodyHash,
// 				},
// 				Body: coin.BlockBody{
// 					Transactions: transactions,
// 				},
// 			},
// 		}

// 		blocks[i] = b
// 	}

// 	return blocks, nil
// }

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
			makeSuccessResponse("1", makeTestReadableBlocks(t)),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getLastBlocksHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got, "%v", got)
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
			makeSuccessResponse("1", makeTestReadableBlocks(t)),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			makeSuccessResponse("1", &readable.Blocks{
				Blocks: []readable.Block{},
			}),
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
	m := &MockGatewayer{}
	m.On("GetBlocks", []uint64{454}).Return(makeTestBlocks(t), nil)
	m.On("GetBlocks", []uint64{1000}).Return([]coin.SignedBlock{}, nil)

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
			makeSuccessResponse("1", makeTestReadableBlocks(t)),
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
			makeSuccessResponse("1", &readable.Blocks{
				Blocks: []readable.Block{},
			}),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, "empty params"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getBlocksBySeqHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}
