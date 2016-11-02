package webrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/visor"
)

var transactionStr = `{
        "transaction": {
            "txn": {
                "length": 220,
                "type": 0,
                "txid": "bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43",
                "inner_hash": "a8558b814926ed0062cd720a572bd67367aa0d01c0769ea4800adcc89cdee524",
                "sigs": [
                    "8756e4bde4ee1c725510a6a9a308c6a90d949de7785978599a87faba601d119f27e1be695cbb32a1e346e5dd88653a97006bf1a93c9673ac59cf7b5db7e0790100"
                ],
                "inputs": [
                    "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
                ],
                "outputs": [
                    {
                        "uxid": "1252a942d721187781147ce4936c46cb7e66c8f9356ca917226dd5f02b4c3695",
                        "dst": "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B",
                        "coins": "1",
                        "hours": 0
                    },
                    {
                        "uxid": "63b7b7ec191d3ece5b8b22a32d78fb74e38cf83858e2ffaac497264a873270d2",
                        "dst": "hFfJrygn4ux3outFEnb1oJ7pqeEqWdAX2M",
                        "coins": "989",
                        "hours": 0
                    }
                ]
            },
            "status": {
                "unconfirmed": false,
                "unknown": false,
                "confirmed": true,
                "height": 103
            }
        }
    }`

var emptyTransactionStr = `{
        "transaction": null
    }`

func decodeTransaction(str string) *visor.TransactionResult {
	t := visor.TransactionResult{}
	if err := json.NewDecoder(strings.NewReader(str)).Decode(&t); err != nil {
		panic(err)
	}
	return &t
}

func Test_getTransactionHandler(t *testing.T) {
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
					Method:  "get_transaction",
					Params:  []byte(`["bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43"]`),
				},
				gateway: &fakeGateway{transactions: map[string]string{
					"bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43": transactionStr,
				}},
			},
			makeSuccessResponse("1", TxnResult{decodeTransaction(transactionStr)}),
		},
		{
			"transaction hash not exist",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_transaction",
					Params:  []byte(`["bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d44"]`),
				},
				gateway: &fakeGateway{},
			},
			makeSuccessResponse("1", TxnResult{}),
		},
		{
			"invalid params: invalid transaction hash",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_transaction",
					Params:  []byte(`["bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d4h"]`),
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, "invalid transaction hash"),
		},
		{
			"invalid params: decode failed",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "get_transaction",
					Params:  []byte("aoo"),
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
					Method:  "get_transaction",
				},
				gateway: &fakeGateway{},
			},
			makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams),
		},
	}
	for _, tt := range tests {
		if got := getTransactionHandler(tt.args.req, tt.args.gateway); !reflect.DeepEqual(got, tt.want) {
			fmt.Println(string(got.Result))
			t.Errorf("%q. getTransactionHandler() = %+v, want %+v", tt.name, got, tt.want)
		}
	}
}
