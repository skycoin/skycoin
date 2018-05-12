package fee

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

type verifyTxFeeTestCase struct {
	inputHours  uint64
	outputHours uint64
	err         error
}

var burnFactor2verifyTxFeeTestCase = []verifyTxFeeTestCase{
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

var burnFactor3verifyTxFeeTestCase = []verifyTxFeeTestCase{
	{0, 0, ErrTxnNoFee},
	{1, 0, nil},
	{1, 1, ErrTxnNoFee},
	{2, 0, nil},
	{2, 1, nil},
	{2, 2, ErrTxnNoFee},
	{3, 0, nil},
	{3, 1, nil},
	{3, 2, nil},
	{3, 3, ErrTxnNoFee},
	{4, 0, nil},
	{4, 1, nil},
	{4, 2, nil},
	{4, 3, ErrTxnInsufficientFee},
	{4, 4, ErrTxnNoFee},
	{5, 0, nil},
	{5, 1, nil},
	{5, 2, nil},
	{5, 3, nil},
	{5, 4, ErrTxnInsufficientFee},
	{5, 5, ErrTxnNoFee},
}

func TestVerifyTransactionFee(t *testing.T) {
	emptyTxn := &coin.Transaction{}
	hours, err := emptyTxn.OutputHours()
	require.NoError(t, err)
	require.Equal(t, uint64(0), hours)

	// A txn with no outputs hours and no coinhours burn fee is valid
	err = VerifyTransactionFee(emptyTxn, 0)
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

	hours, err = txn.OutputHours()
	require.NoError(t, err)
	require.Equal(t, uint64(4e6), hours)

	// A txn with insufficient net coinhours burn fee is invalid
	err = VerifyTransactionFee(txn, 0)
	testutil.RequireError(t, err, ErrTxnNoFee.Error())

	err = VerifyTransactionFee(txn, 1)
	testutil.RequireError(t, err, ErrTxnInsufficientFee.Error())

	// A txn with sufficient net coinhours burn fee is valid
	hours, err = txn.OutputHours()
	require.NoError(t, err)
	err = VerifyTransactionFee(txn, hours)
	require.NoError(t, err)
	hours, err = txn.OutputHours()
	err = VerifyTransactionFee(txn, hours*10)
	require.NoError(t, err)

	// fee + hours overflows
	err = VerifyTransactionFee(txn, math.MaxUint64-3e6)
	testutil.RequireError(t, err, "Hours and fee overflow")

	// txn has overflowing output hours
	txn.Out = append(txn.Out, coin.TransactionOutput{
		Hours: math.MaxUint64 - 1e6 - 3e6 + 1,
	})
	err = VerifyTransactionFee(txn, 10)
	testutil.RequireError(t, err, "Transaction output hours overflow")

	var cases []verifyTxFeeTestCase
	switch BurnFactor {
	case 2:
		cases = burnFactor2verifyTxFeeTestCase
	case 3:
		cases = burnFactor3verifyTxFeeTestCase
	default:
		t.Fatalf("No test cases for BurnFactor=%d", BurnFactor)
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

type requiredFeeTestCase struct {
	hours uint64
	fee   uint64
}

var burnFactor2RequiredFeeTestCases = []requiredFeeTestCase{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 2},
	{5, 3},
	{6, 3},
	{7, 4},
	{998, 499},
	{999, 500},
	{1000, 500},
	{1001, 501},
}

var burnFactor3RequiredFeeTestCases = []requiredFeeTestCase{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 1},
	{4, 2},
	{5, 2},
	{6, 2},
	{7, 3},
	{999, 333},
	{1000, 334},
	{1001, 334},
	{1002, 334},
	{1003, 335},
}

func TestRequiredFee(t *testing.T) {
	var cases []requiredFeeTestCase
	switch BurnFactor {
	case 2:
		cases = burnFactor2RequiredFeeTestCases
	case 3:
		cases = burnFactor3RequiredFeeTestCases
	default:
		t.Fatalf("No test cases for BurnFactor=%d", BurnFactor)
	}

	for _, tc := range cases {
		name := fmt.Sprintf("hours=%d fee=%d", tc.hours, tc.fee)
		t.Run(name, func(t *testing.T) {
			fee := RequiredFee(tc.hours)
			require.Equal(t, tc.fee, fee)

			remainingHours := RemainingHours(tc.hours)
			require.Equal(t, tc.hours-fee, remainingHours)
		})
	}
}

func TestTransactionFee(t *testing.T) {
	var headTime uint64 = 1000
	nextTime := headTime + 3600 // 1 hour later

	type uxInput struct {
		time  uint64
		coins uint64
		hours uint64
	}

	cases := []struct {
		name     string
		out      []uint64
		in       []uxInput
		headTime uint64
		fee      uint64
		err      error
	}{
		// Test case with one output, one input
		{
			fee: 5,
			out: []uint64{5},
			in: []uxInput{
				{time: headTime, coins: 10e6, hours: 10},
			},
			headTime: headTime,
		},

		// Test case with multiple outputs, multiple inputs
		{
			fee: 0,
			out: []uint64{5, 7, 3},
			in: []uxInput{
				{time: headTime, coins: 10e6, hours: 10},
				{time: headTime, coins: 10e6, hours: 5},
			},
			headTime: headTime,
		},

		// Test case with multiple outputs, multiple inputs, and some inputs have more CoinHours once adjusted for HeadTime
		{
			fee: 8,
			out: []uint64{5, 10},
			in: []uxInput{
				{time: nextTime, coins: 10e6, hours: 10},
				{time: headTime, coins: 8e6, hours: 5},
			},
			headTime: nextTime,
		},

		// Test case with insufficient coin hours
		{
			err: ErrTxnInsufficientCoinHours,
			out: []uint64{5, 10, 1},
			in: []uxInput{
				{time: headTime, coins: 10e6, hours: 10},
				{time: headTime, coins: 8e6, hours: 5},
			},
			headTime: headTime,
		},

		// Test case with overflowing input hours
		{
			err: errors.New("UxArray.CoinHours addition overflow"),
			out: []uint64{0},
			in: []uxInput{
				{time: headTime, coins: 10e6, hours: 10},
				{time: headTime, coins: 10e6, hours: math.MaxUint64 - 9},
			},
			headTime: headTime,
		},

		// Test case with overflowing output hours
		{
			err: errors.New("Transaction output hours overflow"),
			out: []uint64{0, 10, math.MaxUint64 - 9},
			in: []uxInput{
				{time: headTime, coins: 10e6, hours: 10},
				{time: headTime, coins: 10e6, hours: 100},
			},
			headTime: headTime,
		},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("fee=%d headTime=%d", tc.fee, tc.headTime)
		t.Run(name, func(t *testing.T) {
			tx := &coin.Transaction{}
			for _, h := range tc.out {
				tx.Out = append(tx.Out, coin.TransactionOutput{
					Hours: h,
				})
			}

			inUxs := make(coin.UxArray, len(tc.in))
			for i, b := range tc.in {
				inUxs[i] = coin.UxOut{
					Head: coin.UxHead{
						Time: b.time,
					},
					Body: coin.UxBody{
						Coins: b.coins,
						Hours: b.hours,
					},
				}
			}

			fee, err := TransactionFee(tx, tc.headTime, inUxs)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.fee, fee)
		})
	}
}
