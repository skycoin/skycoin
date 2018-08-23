package visor

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// ReadableBlockBodyVerbose represents a verbose readable block body
type ReadableBlockBodyVerbose struct {
	Transactions []ReadableTransactionVerbose `json:"txns"`
}

// ReadableBlockVerbose represents a readable block with verbose data
type ReadableBlockVerbose struct {
	Head ReadableBlockHeader      `json:"header"`
	Body ReadableBlockBodyVerbose `json:"body"`
	Size int                      `json:"size'`
}

// ReadableTransactionVerbose has readable transaction data. It differs from ReadableTransaction
// in that it includes metadata for transaction inputs
type ReadableTransactionVerbose struct {
	Status    TransactionStatus `json:"status"`
	Length    uint32            `json:"length"`
	Type      uint8             `json:"type"`
	Hash      string            `json:"txid"`
	InnerHash string            `json:"inner_hash"`
	Timestamp uint64            `json:"timestamp,omitempty"`
	Fee       uint64            `json:"fee"`

	Sigs []string                    `json:"sigs"`
	In   []ReadableTransactionInput  `json:"inputs"`
	Out  []ReadableTransactionOutput `json:"outputs"`
}

// NewReadableTransactionVerbose creates ReadableTransactionVerbose
func NewReadableTransactionVerbose(t Transaction, inputs []ReadableTransactionInput) (ReadableTransactionVerbose, error) {
	// Genesis transaction use empty SHA256 as txid
	txID := cipher.SHA256{}
	if t.Status.BlockSeq != 0 {
		txID = t.Txn.Hash()
	}

	sigs := make([]string, len(t.Txn.Sigs))
	for i, s := range t.Txn.Sigs {
		sigs[i] = s.Hex()
	}

	out := make([]ReadableTransactionOutput, len(t.Txn.Out))
	for i := range t.Txn.Out {
		o, err := NewReadableTransactionOutput(&t.Txn.Out[i], txID)
		if err != nil {
			return ReadableTransactionVerbose{}, err
		}

		out[i] = *o
	}

	var hoursIn uint64
	for _, i := range inputs {
		if _, err := coin.AddUint64(hoursIn, i.CalculatedHours); err != nil {
			logger.Critical().Warningf("Ignoring NewReadableTransactionVerbose summing txn %s input hours error: %v", txID.Hex(), err)
		}
		hoursIn += i.CalculatedHours
	}

	var hoursOut uint64
	for _, o := range t.Txn.Out {
		if _, err := coin.AddUint64(hoursOut, o.Hours); err != nil {
			logger.Critical().Warningf("Ignoring NewReadableTransactionVerbose summing txn %s outputs hours error: %v", txID.Hex(), err)
		}

		hoursOut += o.Hours
	}

	if hoursIn < hoursOut {
		err := fmt.Errorf("NewReadableTransactionVerbose input hours is less than output hours, txid=%s", txID.Hex())
		return ReadableTransactionVerbose{}, err
	}

	fee := hoursIn - hoursOut

	return ReadableTransactionVerbose{
		Status:    t.Status,
		Length:    t.Txn.Length,
		Type:      t.Txn.Type,
		Hash:      t.Txn.Hash().Hex(),
		InnerHash: t.Txn.InnerHash.Hex(),
		Timestamp: t.Time,
		Fee:       fee,

		Sigs: sigs,
		In:   inputs,
		Out:  out,
	}, nil
}
