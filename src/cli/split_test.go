package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// manyOutputs serves n confirmed outputs of one coin at addr.
type manyOutputs struct {
	t    *testing.T
	addr cipher.Address
	n    int
}

func (m manyOutputs) OutputsForAddresses(_ []string) (*readable.UnspentOutputsSummary, error) {
	head := coin.BlockHeader{BkSeq: 10, Time: 1700000000}
	outs := make(readable.UnspentOutputs, m.n)
	for i := range outs {
		ux := coin.UxOut{
			Head: coin.UxHead{BkSeq: 9, Time: head.Time - 3600},
			Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(m.t), Address: m.addr, Coins: 1e6, Hours: 100},
		}
		out, err := readable.NewUnspentOutput(visor.UnspentOutput{UxOut: ux, CalculatedHours: 100})
		require.NoError(m.t, err)
		outs[i] = out
	}
	return &readable.UnspentOutputsSummary{Head: readable.NewBlockHeader(head), HeadOutputs: outs}, nil
}

func TestCreateRawTxnsSplit(t *testing.T) {
	for _, password := range [][]byte{nil, []byte("pwd")} {
		t.Run("encrypted="+string(password), func(t *testing.T) {
			w, err := wallet.NewWallet("test.wlt", "label", "split seed", wallet.Options{
				Type:       wallet.WalletTypeDeterministic,
				GenerateN:  1,
				Encrypt:    len(password) != 0,
				Password:   password,
				CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure,
			})
			require.NoError(t, err)
			e, err := w.GetEntryAt(0)
			require.NoError(t, err)
			from := e.Address.(cipher.Address)
			change := testutil.MakeAddress()
			to := testutil.MakeAddress()
			c := manyOutputs{t, from, 800}
			send := []SendAmount{{Addr: to.String(), Coins: 700e6}}

			_, err = CreateRawTxn(c, w, []string{from.String()}, change.String(), send, password, params.MainNetDistribution)
			require.ErrorAs(t, err, new(transaction.TooLargeError))

			txns, err := CreateRawTxns(c, w, []string{from.String()}, change.String(), send, password, params.MainNetDistribution)
			require.NoError(t, err)
			require.Len(t, txns, 3)

			seen := make(map[cipher.SHA256]bool)
			var sent uint64
			for _, txn := range txns {
				require.True(t, txn.FullySigned())
				size, err := txn.Size()
				require.NoError(t, err)
				require.LessOrEqual(t, size, params.UserVerifyTxn.MaxTransactionSize)
				for _, h := range txn.In {
					require.False(t, seen[h], "an input is spent by two parts")
					seen[h] = true
				}
				var outCoins uint64
				for _, o := range txn.Out {
					outCoins += o.Coins
					if o.Address == to {
						sent += o.Coins
					} else {
						require.Equal(t, change, o.Address)
					}
				}
				require.Equal(t, uint64(len(txn.In))*1e6, outCoins) //nolint:gosec
			}
			require.Equal(t, uint64(700e6), sent)
		})
	}
}
