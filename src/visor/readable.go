package visor

import "github.com/skycoin/skycoin/src/coin"

// Transaction wraps around coin.Transaction, tagged with its status.  This allows us
// to include unconfirmed txns
type Transaction struct {
	Txn    coin.Transaction
	Status TransactionStatus
	Time   uint64
}

// TransactionInput includes the UxOut spent in a transaction and the calculated hours of the output at spending time
type TransactionInput struct {
	UxOut           coin.UxOut
	CalculatedHours uint64
}

// NewTransactionInput creates a TransactionInput.
// calculateHoursTime is the time against which the CalculatedHours should be computed
func NewTransactionInput(ux coin.UxOut, calculateHoursTime uint64) (TransactionInput, error) {
	// The overflow bug causes this to fail for some transactions, allow it to pass
	calculatedHours, err := ux.CoinHours(calculateHoursTime)
	if err != nil {
		logger.Critical().Warningf("Ignoring NewTransactionInput ux.CoinHours failed: %v", err)
		calculatedHours = 0
	}

	return TransactionInput{
		UxOut:           ux,
		CalculatedHours: hours,
	}, nil
}
