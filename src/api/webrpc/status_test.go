package webrpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_getStatusHandler(t *testing.T) {
	b := makeTestReadableBlocks(t)
	now := time.Now().Unix()
	m := &MockGatewayer{}
	m.On("GetLastBlocks", uint64(1)).Return(makeTestBlocks(t), nil)
	m.On("GetTimeNow").Return(uint64(now))

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
				LastBlockHash:      b.Blocks[0].Head.BlockHash,
				TimeSinceLastBlock: fmt.Sprintf("%vs", uint64(now)-b.Blocks[0].Head.Time),
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
			require.Equal(t, tt.want, got)
		})
	}
}
