package webrpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var lastBlockStr = `
{
        "blocks": [
            {
                "header": {
                    "seq": 571,
                    "block_hash": "8b677839c52d04f3d33373d2c39ac8f2c9ec2cbaf15cbd5d12e5d9557c9705bf",
                    "previous_block_hash": "e2074ed9e090a8770dc3c53f27262762ce006e0ddec3c81f823b6b2790aa256e",
                    "timestamp": 1482970140,
                    "fee": 13004,
                    "version": 0,
                    "tx_body_hash": "dc486954796df209bda87947e5c877445c3d4508485e83ddf407f4505366bf74"
                },
                "body": {
                    "txns": [
                        {
                            "length": 317,
                            "type": 0,
                            "txid": "dc486954796df209bda87947e5c877445c3d4508485e83ddf407f4505366bf74",
                            "inner_hash": "8fc6b76d6e861e719142507baccd3c75d4eb09ee76856084be20af31b3db6220",
                            "sigs": [
                                "6df7b96b463bdc354054ef5ac3a8b4fe7586e8cc55ffa50f05d973bd9d228a0c4566f101f29f5e66765e9ecc0893e18a114e8b75503174c193206dc8f0267ece01",
                                "e79c4a98cce98dac4d9aca7bda7f98b202521ff61f4abd521673c1ca0394b87f75dc55eded93c8f87ebf71eb2b69195b6a8a99410753da7200dda43844b199d701"
                            ],
                            "inputs": [
                                "93b907eb6c22475322b19f9812d065cde8635986267b85960d8bd7521f7622d6",
                                "37469c7964ae16f2f29a2afdc94fd6b8db217a00161ceeac5c0b59d92f652d41"
                            ],
                            "outputs": [
                                {
                                    "uxid": "a4131b92b3629fe9614f840986d56ed0d85c9d26839430468fc005fe689b0523",
                                    "dst": "Yk1kPnXauZDLBS4ZoA3SFv6avR2HoqAaHP",
                                    "coins": "5",
                                    "hours": 0
                                },
                                {
                                    "uxid": "ea9c516eb27f8ee1793cbc9ef6b0a9a18e521ab974887a3e45f679da81a6fa6c",
                                    "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                                    "coins": "55",
                                    "hours": 0
                                }
                            ]
                        }
                    ]
                }
            }
        ]
    }`

func Test_getStatusHandler(t *testing.T) {
	b := decodeBlock(lastBlockStr)
	now := time.Now().Unix()
	m := NewGatewayerMock()
	m.On("GetLastBlocks", uint64(1)).Return(b, nil)
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
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStatusHandler(tt.args.req, tt.args.in1)
			require.Equal(t, tt.want, got)
		})
	}
}
