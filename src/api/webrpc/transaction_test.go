package webrpc

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

const (
	rawTxStr    = "dc00000000a8558b814926ed0062cd720a572bd67367aa0d01c0769ea4800adcc89cdee524010000008756e4bde4ee1c725510a6a9a308c6a90d949de7785978599a87faba601d119f27e1be695cbb32a1e346e5dd88653a97006bf1a93c9673ac59cf7b5db7e07901000100000079216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b020000000060dfa95881cdc827b45a6d49b11dbc152ecd4de640420f00000000000000000000000000006409744bcacb181bf98b1f02a11e112d7e4fa9f940f1f23a000000000000000000000000"
	rawTxID     = "bdc4a85a3e9d17a8fe00aa7430d0347c7f1dd6480a16da7147b6e43905057d43"
	txHeight    = uint64(103)
	txConfirmed = true
)

func decodeRawTransaction(rawTxStr string) *visor.Transaction {
	rawTx, err := hex.DecodeString(rawTxStr)
	if err != nil {
		panic(fmt.Sprintf("invalid raw transaction:%v", err))
	}

	tx := coin.MustTransactionDeserialize(rawTx)
	return &visor.Transaction{
		Txn: tx,
		Status: visor.TransactionStatus{
			Confirmed: txConfirmed,
			Height:    txHeight,
		},
	}
}

func Test_getTransactionHandler(t *testing.T) {
	type args struct {
		req     Request
		gateway Gatewayer
	}

	tx := decodeRawTransaction(rawTxStr)
	rbTx, err := visor.NewReadableTransaction(tx)
	require.NoError(t, err)
	txRlt := daemon.TransactionResult{
		Status: visor.TransactionStatus{
			Confirmed: true,
			Height:    103,
		},
		Transaction: *rbTx,
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
					Params:  []byte(fmt.Sprintf(`["%s"]`, rawTxID)),
				},
				gateway: &fakeGateway{transactions: map[string]string{
					rawTxID: rawTxStr,
				}},
			},
			makeSuccessResponse("1", TxnResult{&txRlt}),
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
			MakeErrorResponse(ErrCodeInvalidRequest, "transaction doesn't exist"),
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
			MakeErrorResponse(ErrCodeInvalidParams, "invalid transaction hash"),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
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
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getTransactionHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_injectTransactionHandler(t *testing.T) {
	var rawTx = `dc0000000010e05181fd4023f865a84359bf72a304e687b6f00e42f93ad9a4b8ee5a64aabc01000000dcb5b236eecd97a36c7d0a0b8ed68bb5df6274433a51fddf911f02f3926d20bf6eaabdc21529b7696f498545b06cc7e69f2f08b4dc5fa823c5b3f03da06794a300010000006d8a9c89177ce5e9d3b4b59fff67c00f0471fdebdfbb368377841b03fc7d688b02000000005771eeda2e253697cf5368f16fe05210d5cd319040420f0000000000af010000000000000060dfa95881cdc827b45a6d49b11dbc152ecd4de600093d0000000000af01000000000000`
	var txid = "3e52703a21bf9462799f52ab0cedb314efcf7c43aadb815429cd79f35f040954"

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
					Method:  "inject_transaction",
					Params:  []byte(fmt.Sprintf("[%q]", rawTx)),
				},
				gateway: &fakeGateway{
					injectRawTxMap: map[string]bool{
						txid: true,
					},
				},
			},
			makeSuccessResponse("1", TxIDJson{txid}),
		},
		{
			"invalid params: invalid raw transaction",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "inject_transaction",
					Params:  []byte(`["abcdeffedca"]`),
				},
			},
			MakeErrorResponse(ErrCodeInvalidParams, "invalid raw transaction: encoding/hex: odd length hex string"),
		},
		{
			"invalid params type",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "inject_transaction",
					Params:  []byte("abcdeffedca"),
				},
			},
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
		{
			"invalid params: more than one raw transaction",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "inject_transaction",
					Params:  []byte(fmt.Sprintf("[%q,%q]", rawTx, rawTx)),
				},
			},
			MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams),
		},
		{
			"internal error",
			args{
				req: Request{
					ID:      "1",
					Jsonrpc: jsonRPC,
					Method:  "inject_transaction",
					Params:  []byte(fmt.Sprintf("[%q]", rawTx)),
				},
				gateway: &fakeGateway{},
			},
			MakeErrorResponse(ErrCodeInternalError, "inject transaction failed: fake gateway inject transaction failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := injectTransactionHandler(tt.args.req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}
