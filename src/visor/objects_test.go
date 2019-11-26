package visor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/transaction"
)

func TestNewTransactionInputsFromUxBalance(t *testing.T) {
	headTime := uint64(time.Now().Unix())

	uxa := []coin.UxOut{
		{
			Head: coin.UxHead{
				Time:  headTime / 2,
				BkSeq: 60,
			},
			Body: coin.UxBody{
				SrcTransaction: testutil.RandSHA256(t),
				Address:        testutil.MakeAddress(),
				Coins:          11e6,
				Hours:          12345,
			},
		},
		{
			Head: coin.UxHead{
				Time:  headTime/2 + headTime/4,
				BkSeq: 120,
			},
			Body: coin.UxBody{
				SrcTransaction: testutil.RandSHA256(t),
				Address:        testutil.MakeAddress(),
				Coins:          12345678,
				Hours:          987654,
			},
		},
	}

	uxb, err := transaction.NewUxBalances(uxa, headTime)
	require.NoError(t, err)

	inputsFromUxa, err := NewTransactionInputs(uxa, headTime)
	require.NoError(t, err)

	inputsFromUxb := NewTransactionInputsFromUxBalance(uxb)
	require.Equal(t, inputsFromUxa, inputsFromUxb)

	require.Nil(t, NewTransactionInputsFromUxBalance([]transaction.UxBalance{}))
}
