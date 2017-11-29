package fee

import (
	"errors"

	"github.com/skycoin/skycoin/src/coin"
)

const (
	// BurnFactor half of coinhours must be burnt
	BurnFactor uint64 = 2
)

var (
	// ErrTxnNoFee is returned if a transaction has no coinhour fee
	ErrTxnNoFee = errors.New("Transaction has zero coinhour fee")

	// ErrTxnInsufficientFee is returned if a transaction's coinhour burn fee is not enough
	ErrTxnInsufficientFee = errors.New("Transaction coinhour fee minimum not met")
)

// VerifyTransactionFee performs additional transaction verification at the unconfirmed pool level.
// This checks tunable parameters that should prevent the transaction from
// entering the blockchain, but cannot be done at the blockchain level because
// they may be changed.
func VerifyTransactionFee(t *coin.Transaction, fee uint64) error {
	return VerifyTransactionFeeForHours(t.OutputHours(), fee)
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
	total := hours + fee

	// Make sure at least half (BurnFactor=2) the coinhours are destroyed
	// Compare hours to the burned coins, because total/BurnFactor
	// rounds down.
	// Examples for BurnFactor=2:
	// If inputs are 3, one coinhour can be spent and the fee is 2
	// If inputs are 2, one coinhour can be spent and the fee is 1
	// If inputs are 1, no coinhour can be spent, and the fee is 1.
	// If inputs are 0, the fee is 0 and ErrTxnNoFee was returned.
	if hours > total/BurnFactor {
		return ErrTxnInsufficientFee
	}

	return nil
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction
func TransactionFee(t *coin.Transaction, headTime uint64, inUxs coin.UxArray) (uint64, error) {
	// Compute input hours
	inHours := uint64(0)
	for _, ux := range inUxs {
		inHours += ux.CoinHours(headTime)
	}

	// Compute output hours
	outHours := uint64(0)
	for i := range t.Out {
		outHours += t.Out[i].Hours
	}

	if inHours < outHours {
		return 0, errors.New("Insufficient coinhours for transaction outputs")
	}

	return inHours - outHours, nil
}
