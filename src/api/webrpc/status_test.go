package webrpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getStatusHandler(t *testing.T) {
	b := makeTestReadableBlocks(t)
	m := &MockGatewayer{}
	m.On("GetLastBlocks", uint64(1)).Return(makeTestBlocks(t), nil)

	type args struct {
		req Request
		in1 Gatewayer
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
					Method:  "get_status",
					Jsonrpc: jsonRPC,
				},
				in1: m,
			},
			makeSuccessResponse("1", StatusResult{
				Running:            true,
				BlockNum:           b.Blocks[0].Head.BkSeq + 1,
				LastBlockHash:      b.Blocks[0].Head.Hash,
				TimeSinceLastBlock: "", // can't check this
			}),
		},
		{
			"invalid params",
			args{
				req: Request{
					ID:      "1",
					Method:  "get_status",
					Jsonrpc: jsonRPC,
					Params:  []byte(`{"abc": "123"}`),
				},
				in1: m,
			},
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatusHandler(tt.args.req, tt.args.in1)

			if got.Error == nil {
				// Patch out TimeSinceLastBlock since it can increment during the test
				var gotStatus StatusResult
				err := json.Unmarshal(got.Result, &gotStatus)
				require.NoError(t, err)
				require.NotEmpty(t, gotStatus.TimeSinceLastBlock)
				require.NotEqual(t, "s", gotStatus.TimeSinceLastBlock)
				gotStatus.TimeSinceLastBlock = ""
				d, err := json.Marshal(gotStatus)
				require.NoError(t, err)
				got.Result = d
			}

			require.Equal(t, tt.want, got)
		})
	}
}
