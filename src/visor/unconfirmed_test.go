package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func TestUnconfirmedVerifyTransactionFee(t *testing.T) {
	emptyTxn := &coin.Transaction{}
	require.Equal(t, uint64(0), emptyTxn.OutputHours())

	// A txn with no outputs hours and no coinhours burn fee is valid
	err := VerifyTransactionFee(emptyTxn, 0)
	testutil.RequireError(t, err, ErrTxnNoCoinHours.Error())

	// A txn with no outputs hours but with a coinhours burn fee is valid
	err = VerifyTransactionFee(emptyTxn, 100)
	require.NoError(t, err)

	txn := &coin.Transaction{}
	txn.Out = append(txn.Out, coin.TransactionOutput{
		Coins: 1e6,
		Hours: 1e6,
	})
	txn.Out = append(txn.Out, coin.TransactionOutput{
		Coins: 1e6,
		Hours: 3e6,
	})
	require.Equal(t, uint64(4e6), txn.OutputHours())

	// A txn with insufficient net coinhours burn fee is invalid
	err = VerifyTransactionFee(txn, 0)
	testutil.RequireError(t, err, ErrTxnInsufficientCoinHourFee.Error())

	err = VerifyTransactionFee(txn, txn.OutputHours()-2)
	testutil.RequireError(t, err, ErrTxnInsufficientCoinHourFee.Error())

	// A txn with sufficient net coinhours burn fee is valid
	err = VerifyTransactionFee(txn, txn.OutputHours())
	require.NoError(t, err)
	err = VerifyTransactionFee(txn, txn.OutputHours()*10)
	require.NoError(t, err)
}
