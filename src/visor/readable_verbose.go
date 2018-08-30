package visor

import (
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// ReadableBlockBodyVerbose represents a verbose readable block body
type ReadableBlockBodyVerbose struct {
	Transactions []ReadableBlockTransactionVerbose `json:"txns"`
}

// ReadableBlockVerbose represents a readable block with verbose data
type ReadableBlockVerbose struct {
	Head ReadableBlockHeader      `json:"header"`
	Body ReadableBlockBodyVerbose `json:"body"`
	Size int                      `json:"size"`
}

// NewReadableBlockBodyVerbose creates a verbose readable block body
func NewReadableBlockBodyVerbose(b *coin.Block, inputs [][]ReadableTransactionInput) (*ReadableBlockBodyVerbose, error) {
	if len(inputs) != len(b.Body.Transactions) {
		return nil, fmt.Errorf("NewReadableBlockBodyVerbose: len(inputs) != len(b.Body.Transactions) (seq=%d)", b.Head.BkSeq)
	}

	txns := make([]ReadableBlockTransactionVerbose, len(b.Body.Transactions))
	for i := range b.Body.Transactions {
		t := b.Body.Transactions[i]

		tx, err := NewReadableBlockTransactionVerbose(t, inputs[i], b.Head.BkSeq == 0)
		if err != nil {
			return nil, err
		}
		txns[i] = tx
	}

	return &ReadableBlockBodyVerbose{
		Transactions: txns,
	}, nil
}

// NewReadableBlockVerbose creates a verbose readable block
func NewReadableBlockVerbose(b *coin.Block, inputs [][]ReadableTransactionInput) (*ReadableBlockVerbose, error) {
	body, err := NewReadableBlockBodyVerbose(b, inputs)
	if err != nil {
		return nil, err
	}

	return &ReadableBlockVerbose{
		Head: NewReadableBlockHeader(&b.Head),
		Body: *body,
		Size: b.Size(),
	}, nil
}

// ReadableBlocksVerbose an array of verbose readable blocks.
type ReadableBlocksVerbose struct {
	Blocks []ReadableBlockVerbose `json:"blocks"`
}

// NewReadableBlocksVerbose creates ReadableBlocksVerbose from []ReadableBlockVerbose
func NewReadableBlocksVerbose(blocks []ReadableBlockVerbose) *ReadableBlocksVerbose {
	if blocks == nil {
		blocks = []ReadableBlockVerbose{}
	}
	return &ReadableBlocksVerbose{
		Blocks: blocks,
	}
}

// ReadableBlockTransactionVerbose has readable transaction data for transactions inside a block. It differs from ReadableTransaction
// in that it includes metadata for transaction inputs and the calculated coinhour fee spent by the block
type ReadableBlockTransactionVerbose struct {
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	Hash      string `json:"txid"`
	InnerHash string `json:"inner_hash"`
	Fee       uint64 `json:"fee"`

	Sigs []string                    `json:"sigs"`
	In   []ReadableTransactionInput  `json:"inputs"`
	Out  []ReadableTransactionOutput `json:"outputs"`
}

// NewReadableBlockTransactionVerbose creates ReadableBlockTransactionVerbose
func NewReadableBlockTransactionVerbose(txn coin.Transaction, inputs []ReadableTransactionInput, isGenesis bool) (ReadableBlockTransactionVerbose, error) {
	if len(inputs) != len(txn.In) {
		return ReadableBlockTransactionVerbose{}, errors.New("NewReadableTransactionVerbose: len(inputs) != len(txn.In)")
	}

	// Genesis transaction uses empty SHA256 as txid
	// FIXME: If/when the blockchain is regenerated, use a real hash as the txID for the genesis block. The bkSeq argument can be removed then.
	txID := cipher.SHA256{}
	if !isGenesis {
		txID = txn.Hash()
	}

	sigs := make([]string, len(txn.Sigs))
	for i, s := range txn.Sigs {
		sigs[i] = s.Hex()
	}

	out := make([]ReadableTransactionOutput, len(txn.Out))
	for i := range txn.Out {
		o, err := NewReadableTransactionOutput(&txn.Out[i], txID)
		if err != nil {
			return ReadableBlockTransactionVerbose{}, err
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
	for _, o := range txn.Out {
		if _, err := coin.AddUint64(hoursOut, o.Hours); err != nil {
			logger.Critical().Warningf("Ignoring NewReadableTransactionVerbose summing txn %s outputs hours error: %v", txID.Hex(), err)
		}

		hoursOut += o.Hours
	}

	var fee uint64
	if isGenesis {
		if hoursIn != 0 {
			err := errors.New("NewReadableTransactionVerbose genesis block should have 0 input hours")
			return ReadableBlockTransactionVerbose{}, err
		}

		fee = 0
	} else {
		if hoursIn < hoursOut {
			err := fmt.Errorf("NewReadableTransactionVerbose input hours is less than output hours, txid=%s", txID.Hex())
			return ReadableBlockTransactionVerbose{}, err
		}

		fee = hoursIn - hoursOut
	}

	if inputs == nil {
		inputs = make([]ReadableTransactionInput, 0)
	}

	return ReadableBlockTransactionVerbose{
		Length:    txn.Length,
		Type:      txn.Type,
		Hash:      txn.Hash().Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fee,

		Sigs: sigs,
		In:   inputs,
		Out:  out,
	}, nil
}

// ReadableTransactionVerbose has readable transaction data. It adds TransactionStatus to a ReadableBlockTransactionVerbose
type ReadableTransactionVerbose struct {
	Status    *TransactionStatus `json:"status,omitempty"`
	Timestamp uint64             `json:"timestamp,omitempty"`
	ReadableBlockTransactionVerbose
}

// NewReadableTransactionVerbose creates ReadableTransactionVerbose
func NewReadableTransactionVerbose(txn Transaction, inputs []ReadableTransactionInput) (ReadableTransactionVerbose, error) {
	rb, err := NewReadableBlockTransactionVerbose(txn.Txn, inputs, txn.Status.BlockSeq == 0 && txn.Status.Confirmed)
	if err != nil {
		return ReadableTransactionVerbose{}, err
	}

	return ReadableTransactionVerbose{
		Status:                          &txn.Status,
		Timestamp:                       txn.Time,
		ReadableBlockTransactionVerbose: rb,
	}, nil
}

// ReadableUnconfirmedTxnVerbose represents a verbose readable unconfirmed transaction
type ReadableUnconfirmedTxnVerbose struct {
	Txn       ReadableBlockTransactionVerbose `json:"transaction"`
	Received  time.Time                       `json:"received"`
	Checked   time.Time                       `json:"checked"`
	Announced time.Time                       `json:"announced"`
	IsValid   bool                            `json:"is_valid"`
}

// NewReadableUnconfirmedTxnVerbose creates a verbose readable unconfirmed transaction
func NewReadableUnconfirmedTxnVerbose(unconfirmed *UnconfirmedTxn, inputs []ReadableTransactionInput) (*ReadableUnconfirmedTxnVerbose, error) {
	isGenesis := false // The genesis transaction is never unconfirmed
	txn, err := NewReadableBlockTransactionVerbose(unconfirmed.Txn, inputs, isGenesis)
	if err != nil {
		return nil, err
	}

	return &ReadableUnconfirmedTxnVerbose{
		Txn:       txn,
		Received:  nanoToTime(unconfirmed.Received),
		Checked:   nanoToTime(unconfirmed.Checked),
		Announced: nanoToTime(unconfirmed.Announced),
		IsValid:   unconfirmed.IsValid == 1,
	}, nil
}

// NewReadableUnconfirmedTxnsVerbose creates []ReadableUnconfirmedTxn from []UnconfirmedTxn and their readable transaction inputs
func NewReadableUnconfirmedTxnsVerbose(txns []UnconfirmedTxn, inputs [][]ReadableTransactionInput) ([]ReadableUnconfirmedTxnVerbose, error) {
	if len(inputs) != len(txns) {
		return nil, fmt.Errorf("NewReadableUnconfirmedTxnsVerbose: len(inputs) != len(txns)")
	}

	rTxns := make([]ReadableUnconfirmedTxnVerbose, len(txns))
	for i, txn := range txns {
		rTxn, err := NewReadableUnconfirmedTxnVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}

		rTxns[i] = *rTxn
	}

	return rTxns, nil
}
