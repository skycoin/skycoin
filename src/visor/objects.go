package visor

import (
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// Transaction wraps around coin.Transaction, tagged with its status.  This allows us
// to include unconfirmed txns
type Transaction struct {
	Transaction coin.Transaction
	Status      TransactionStatus
	Time        uint64
}

// TransactionStatus represents the transaction status
type TransactionStatus struct {
	Confirmed bool
	// If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed.
	Height uint64
	// If confirmed, the sequence of the block in which the transaction was executed
	BlockSeq uint64
}

// NewUnconfirmedTransactionStatus creates unconfirmed transaction status
func NewUnconfirmedTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Confirmed: false,
		Height:    0,
		BlockSeq:  0,
	}
}

// NewConfirmedTransactionStatus creates confirmed transaction status
func NewConfirmedTransactionStatus(height, blockSeq uint64) TransactionStatus {
	// Height starts at 1
	// TODO -- height should start at 0?
	if height == 0 {
		logger.Panic("Invalid confirmed transaction height")
	}
	return TransactionStatus{
		Confirmed: true,
		Height:    height,
		BlockSeq:  blockSeq,
	}
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
		CalculatedHours: calculatedHours,
	}, nil
}

// BlockchainMetadata encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block
	HeadBlock coin.SignedBlock
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64
	// Number of known unconfirmed txns
	Unconfirmed uint64
}

// NewBlockchainMetadata creates blockchain meta data
func NewBlockchainMetadata(head coin.SignedBlock, unconfirmedLen, unspentsLen uint64) (*BlockchainMetadata, error) {
	return &BlockchainMetadata{
		HeadBlock:   head,
		Unspents:    unspentsLen,
		Unconfirmed: unconfirmedLen,
	}, nil
}

// UnconfirmedTransaction unconfirmed transaction
type UnconfirmedTransaction struct {
	Transaction coin.Transaction
	// Time the txn was last received
	Received int64
	// Time the txn was last checked against the blockchain
	Checked int64
	// Last time we announced this txn
	Announced int64
	// If this txn is valid
	IsValid int8
}

// Hash returns the coin.Transaction's hash
func (ut *UnconfirmedTransaction) Hash() cipher.SHA256 {
	return ut.Transaction.Hash()
}

// NewUnconfirmedTransaction creates an UnconfirmedTransaction
func NewUnconfirmedTransaction(txn coin.Transaction) UnconfirmedTransaction {
	now := time.Now().UTC()
	return UnconfirmedTransaction{
		Transaction: txn,
		Received:    now.UnixNano(),
		Checked:     now.UnixNano(),
		Announced:   time.Time{}.UnixNano(),
		IsValid:     0,
	}
}

// UnspentOutput includes coin.UxOut and adds CalculatedHours
type UnspentOutput struct {
	coin.UxOut
	CalculatedHours uint64
}

// NewUnspentOutput creates an UnspentOutput
func NewUnspentOutput(uxOut coin.UxOut, calculateHoursTime uint64) (UnspentOutput, error) {
	calculatedHours, err := uxOut.CoinHours(calculateHoursTime)

	// Treat overflowing coin hours calculations as a non-error and force hours to 0
	// This affects one bad spent output which had overflowed hours, spent in block 13277.
	switch err {
	case nil:
	case coin.ErrAddEarnedCoinHoursAdditionOverflow:
		calculatedHours = 0
	default:
		return UnspentOutput{}, err
	}

	return UnspentOutput{
		UxOut:           uxOut,
		CalculatedHours: calculatedHours,
	}, nil
}

// NewUnspentOutputs creates []UnspentOutput
func NewUnspentOutputs(uxOuts []coin.UxOut, calculateHoursTime uint64) ([]UnspentOutput, error) {
	outs := make([]UnspentOutput, len(uxOuts))
	for i, ux := range uxOuts {
		u, err := NewUnspentOutput(ux, calculateHoursTime)
		if err != nil {
			return nil, err
		}
		outs[i] = u
	}

	return outs, nil
}

// UnspentOutputsSummary includes current unspent outputs and incoming and outgoing unspent outputs
type UnspentOutputsSummary struct {
	Confirmed []UnspentOutput
	Outgoing  []UnspentOutput
	Incoming  []UnspentOutput
}
