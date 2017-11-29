package fee

import (
	"fmt"
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
	testutil.RequireError(t, err, ErrTxnNoFee.Error())

	// A txn with no outputs hours but with a coinhours burn fee is valid
	err = VerifyTransactionFee(emptyTxn, 100)
	require.NoError(t, err)

	txn := &coin.Transaction{}
	txn.Out = append(txn.Out, coin.TransactionOutput{
		Hours: 1e6,
	})
	txn.Out = append(txn.Out, coin.TransactionOutput{
		Hours: 3e6,
	})
	require.Equal(t, uint64(4e6), txn.OutputHours())

	// A txn with insufficient net coinhours burn fee is invalid
	err = VerifyTransactionFee(txn, 0)
	testutil.RequireError(t, err, ErrTxnNoFee.Error())

	err = VerifyTransactionFee(txn, 1)
	testutil.RequireError(t, err, ErrTxnInsufficientFee.Error())

	err = VerifyTransactionFee(txn, txn.OutputHours()-2)
	testutil.RequireError(t, err, ErrTxnInsufficientFee.Error())

	// A txn with sufficient net coinhours burn fee is valid
	err = VerifyTransactionFee(txn, txn.OutputHours())
	require.NoError(t, err)
	err = VerifyTransactionFee(txn, txn.OutputHours()*10)
	require.NoError(t, err)

	cases := []struct {
		inputHours  uint64
		outputHours uint64
		err         error
	}{
		{0, 0, ErrTxnNoFee},
		{1, 0, nil},
		{1, 1, ErrTxnNoFee},
		{2, 0, nil},
		{2, 1, nil},
		{2, 2, ErrTxnNoFee},
		{3, 0, nil},
		{3, 1, nil},
		{3, 2, ErrTxnInsufficientFee},
		{3, 3, ErrTxnNoFee},
		{4, 0, nil},
		{4, 1, nil},
		{4, 2, nil},
		{4, 3, ErrTxnInsufficientFee},
		{4, 4, ErrTxnNoFee},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("input=%d output=%d", tc.inputHours, tc.outputHours)
		t.Run(name, func(t *testing.T) {
			txn := &coin.Transaction{}
			txn.Out = append(txn.Out, coin.TransactionOutput{
				Hours: tc.outputHours,
			})

			require.True(t, tc.inputHours >= tc.outputHours)
			err := VerifyTransactionFee(txn, tc.inputHours-tc.outputHours)
			require.Equal(t, tc.err, err)
		})
	}
}
