package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

// pqFakeOutputs serves one confirmed output of 10 coins at addr.
type pqFakeOutputs struct {
	t    *testing.T
	addr cipher.Address
	ux   coin.UxOut
}

func newPQFakeOutputs(t *testing.T, addr cipher.Address) *pqFakeOutputs {
	return &pqFakeOutputs{t: t, addr: addr, ux: coin.UxOut{
		Head: coin.UxHead{BkSeq: 9, Time: 1700000000 - 3600},
		Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(t), Address: addr, Coins: 10e6, Hours: 1000},
	}}
}

func (f *pqFakeOutputs) OutputsForAddresses(_ []string) (*readable.UnspentOutputsSummary, error) {
	out, err := readable.NewUnspentOutput(visor.UnspentOutput{UxOut: f.ux, CalculatedHours: 1000})
	require.NoError(f.t, err)
	return &readable.UnspentOutputsSummary{
		Head:        readable.NewBlockHeader(coin.BlockHeader{BkSeq: 10, Time: 1700000000}),
		HeadOutputs: readable.UnspentOutputs{out},
	}, nil
}

func TestCreatePQTxn(t *testing.T) {
	sk := pqsig.DeriveSecKey([]byte("seed"), 0)
	pk, err := sk.PubKey()
	require.NoError(t, err)
	from := cipher.AddressFromPQPubKey(pk)
	f := newPQFakeOutputs(t, from)

	to := testutil.MakeAddress()
	txn, err := createPQTxn(f, sk, to, 4e6)
	require.NoError(t, err)
	require.Equal(t, coin.TxnTypePQ, txn.Type)
	require.Len(t, txn.Out, 2)
	require.Equal(t, to, txn.Out[0].Address)
	require.Equal(t, from, txn.Out[1].Address, "change returns to the post-quantum address")
	require.NoError(t, txn.VerifyInputSignatures(coin.UxArray{f.ux}))

	raw, err := txn.SerializeHexWithPQ()
	require.NoError(t, err)
	back, err := coin.DeserializeTransactionHexWithPQ(raw)
	require.NoError(t, err)
	require.NoError(t, back.VerifyInputSignatures(coin.UxArray{f.ux}))
}
