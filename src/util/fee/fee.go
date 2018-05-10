package fee

import (
	"errors"

	"github.com/skycoin/skycoin/src/coin"
)

const (
	// BurnFactor inverse fraction of coinhours that must be burned
	BurnFactor uint64 = 2
)

var (
	// ErrTxnNoFee is returned if a transaction has no coinhour fee
	ErrTxnNoFee = errors.New("Transaction has zero coinhour fee")

	// ErrTxnInsufficientFee is returned if a transaction's coinhour burn fee is not enough
	ErrTxnInsufficientFee = errors.New("Transaction coinhour fee minimum not met")

	// ErrTxnInsufficientCoinHours is returned if a transaction has more coinhours in its outputs than its inputs
	ErrTxnInsufficientCoinHours = errors.New("Insufficient coinhours for transaction outputs")
)

// VerifyTransactionFee performs additional transaction verification at the unconfirmed pool level.
// This checks tunable parameters that should prevent the transaction from
// entering the blockchain, but cannot be done at the blockchain level because
// they may be changed.
func VerifyTransactionFee(t *coin.Transaction, fee uint64) error {
	hours, err := t.OutputHours()
	if err != nil {
		return err
	}
	return VerifyTransactionFeeForHours(hours, fee)
}

// VerifyTransactionFeeForHours verifies the fee given fee and hours,
// where hours is the number of hours in a transaction's outputs,
// and hours+fee is the number of hours in a transaction's inputs
func VerifyTransactionFeeForHours(hours, fee uint64) error {
	// Require non-zero coinhour fee
	if fee == 0 {
		return ErrTxnNoFee
	}

	// Calculate total number of coinhours
	total, err := coin.AddUint64(hours, fee)
	if err != nil {
		return errors.New("Hours and fee overflow")
	}

	// Calculate the required fee
	requiredFee := RequiredFee(total)

	// Ensure that the required fee is met
	if fee < requiredFee {
		return ErrTxnInsufficientFee
	}

	return nil
}

// RequiredFee returns the coinhours fee required for an amount of hours
// The required fee is calculated as hours/BurnFactor, rounded up.
func RequiredFee(hours uint64) uint64 {
	feeHours := hours / BurnFactor
	if hours%BurnFactor != 0 {
		feeHours++
	}

	return feeHours
}

// RemainingHours returns the amount of coinhours leftover after paying the fee for the input.
func RemainingHours(hours uint64) uint64 {
	fee := RequiredFee(hours)
	return hours - fee
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction.
// Returns ErrTxnInsufficientCoinHours if input hours is less than output hours.
func TransactionFee(tx *coin.Transaction, headTime uint64, inUxs coin.UxArray) (uint64, error) {
	// Compute input hours
	inHours, err := inUxs.CoinHours(headTime)
	if err != nil {
		return 0, err
	}

	// Compute output hours
	outHours, err := tx.OutputHours()
	if err != nil {
		return 0, err
	}

	if inHours < outHours {
		return 0, ErrTxnInsufficientCoinHours
	}

	return inHours - outHours, nil
}
