package webrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/uxotutil"
	"github.com/stretchr/testify/require"
)

func Test_getRichlistHandler(t *testing.T) {
	uxouts := make([]coin.UxOut, 3)
	addrs := []string{"fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B", "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B"}
	coins := []uint64{20 * 1e6, 30 * 1e6, 40 * 1e6}
	for i := 0; i < 3; i++ {
		addr, err := cipher.DecodeBase58Address(addrs[i])
		require.NoError(t, err)
		uxouts[i] = coin.UxOut{}
		uxouts[i].Body.Address = addr
		uxouts[i].Body.Coins = coins[i]
	}

	account1 := uxotutil.AccountJSON{Addr: "fyqX5YuwXMUs4GEUE3LjLyhrqvNztFHQ4B", Coins: "60.000000", Locked: false}
	account2 := uxotutil.AccountJSON{Addr: "cBnu9sUvv12dovBmjQKTtfE4rbjMmf3fzW", Coins: "30.000000", Locked: false}
	type args struct {
		Topn           int
		IsDistribution bool
		gateway        Gatewayer
	}
	tests := []struct {
		name string
		args args
		want Response
	}{
		{
			"richest one",
			args{
				Topn:           1,
				IsDistribution: false,
				gateway:        &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsTopn{[]uxotutil.AccountJSON{account1}}),
		},
		{
			"richest two",
			args{
				Topn:           2,
				IsDistribution: false,
				gateway:        &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsTopn{[]uxotutil.AccountJSON{account1, account2}}),
		},
		{
			"richest all",
			args{
				Topn:           -1,
				IsDistribution: true,
				gateway:        &fakeGateway{uxouts: uxouts},
			},
			makeSuccessResponse("1", OutputsTopn{[]uxotutil.AccountJSON{account1, account2}}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, err := json.Marshal(tt.args)
			fmt.Println("param:", string(params))
			require.NoError(t, err)
			req := Request{
				ID:      "1",
				Jsonrpc: jsonRPC,
				Method:  "get_richlist",
				Params:  params,
			}

			got := getRichlistHandler(req, tt.args.gateway)
			require.Equal(t, tt.want, got)
		})
	}
}
