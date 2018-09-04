package webrpc

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
)

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
