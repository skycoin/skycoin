package readable

import (
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/util/timeutil"
	"github.com/skycoin/skycoin/src/visor"
)

// BlockBodyVerbose represents a verbose readable block body
type BlockBodyVerbose struct {
	Transactions []BlockTransactionVerbose `json:"txns"`
}

// BlockVerbose represents a readable block with verbose data
type BlockVerbose struct {
	Head BlockHeader      `json:"header"`
	Body BlockBodyVerbose `json:"body"`
	Size uint32           `json:"size"`
}

// NewBlockBodyVerbose creates a verbose readable block body
func NewBlockBodyVerbose(b coin.Block, inputs [][]visor.TransactionInput) (*BlockBodyVerbose, error) {
	if len(inputs) != len(b.Body.Transactions) {
		return nil, fmt.Errorf("NewBlockBodyVerbose: len(inputs) != len(b.Body.Transactions) (seq=%d)", b.Head.BkSeq)
	}

	txns := make([]BlockTransactionVerbose, len(b.Body.Transactions))
	for i := range b.Body.Transactions {
		t := b.Body.Transactions[i]

		txn, err := NewBlockTransactionVerbose(t, inputs[i], b.Head.BkSeq == 0)
		if err != nil {
			return nil, err
		}
		txns[i] = txn
	}

	return &BlockBodyVerbose{
		Transactions: txns,
	}, nil
}

// NewBlockVerbose creates a verbose readable block
func NewBlockVerbose(b coin.Block, inputs [][]visor.TransactionInput) (*BlockVerbose, error) {
	body, err := NewBlockBodyVerbose(b, inputs)
	if err != nil {
		return nil, err
	}

	size, err := b.Size()
	if err != nil {
		return nil, err
	}

	return &BlockVerbose{
		Head: NewBlockHeader(b.Head),
		Body: *body,
		Size: size,
	}, nil
}

// BlocksVerbose an array of verbose readable blocks.
type BlocksVerbose struct {
	Blocks []BlockVerbose `json:"blocks"`
}

// NewBlocksVerbose creates BlocksVerbose from []BlockVerbose
func NewBlocksVerbose(blocks []coin.SignedBlock, inputs [][][]visor.TransactionInput) (*BlocksVerbose, error) {
	bs := make([]BlockVerbose, len(blocks))
	for i := range blocks {
		if i >= len(inputs) {
			return nil, errors.New("NewBlocksVerbose: not enough inputs for blocks")
		}

		b, err := NewBlockVerbose(blocks[i].Block, inputs[i])
		if err != nil {
			return nil, err
		}

		bs[i] = *b
	}

	return &BlocksVerbose{
		Blocks: bs,
	}, nil
}

// BlockTransactionVerbose has readable transaction data for transactions inside a block. It differs from Transaction
// in that it includes metadata for transaction inputs and the calculated coinhour fee spent by the block
type BlockTransactionVerbose struct {
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	Hash      string `json:"txid"`
	InnerHash string `json:"inner_hash"`
	Fee       uint64 `json:"fee"`

	Sigs []string            `json:"sigs"`
	In   []TransactionInput  `json:"inputs"`
	Out  []TransactionOutput `json:"outputs"`
}

// NewBlockTransactionVerbose creates BlockTransactionVerbose
func NewBlockTransactionVerbose(txn coin.Transaction, inputs []visor.TransactionInput, isGenesis bool) (BlockTransactionVerbose, error) {
	if len(inputs) != len(txn.In) {
		return BlockTransactionVerbose{}, errors.New("NewBlockTransactionVerbose: len(inputs) != len(txn.In)")
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

	out := make([]TransactionOutput, len(txn.Out))
	for i := range txn.Out {
		o, err := NewTransactionOutput(&txn.Out[i], txID)
		if err != nil {
			return BlockTransactionVerbose{}, err
		}

		out[i] = *o
	}

	var hoursIn uint64
	for _, i := range inputs {
		if _, err := mathutil.AddUint64(hoursIn, i.CalculatedHours); err != nil {
			logger.Critical().Warningf("Ignoring NewBlockTransactionVerbose summing txn %s input hours error: %v", txID.Hex(), err)
		}
		hoursIn += i.CalculatedHours
	}

	var hoursOut uint64
	for _, o := range txn.Out {
		if _, err := mathutil.AddUint64(hoursOut, o.Hours); err != nil {
			logger.Critical().Warningf("Ignoring NewBlockTransactionVerbose summing txn %s outputs hours error: %v", txID.Hex(), err)
		}

		hoursOut += o.Hours
	}

	var fee uint64
	if isGenesis {
		if hoursIn != 0 {
			err := errors.New("NewBlockTransactionVerbose genesis block should have 0 input hours")
			return BlockTransactionVerbose{}, err
		}

		fee = 0
	} else {
		if hoursIn < hoursOut {
			err := fmt.Errorf("NewBlockTransactionVerbose input hours is less than output hours, txid=%s", txID.Hex())
			return BlockTransactionVerbose{}, err
		}

		fee = hoursIn - hoursOut
	}

	txnInputs := make([]TransactionInput, len(inputs))
	for i, input := range inputs {
		var err error
		txnInputs[i], err = NewTransactionInput(input)
		if err != nil {
			return BlockTransactionVerbose{}, err
		}
	}

	return BlockTransactionVerbose{
		Length:    txn.Length,
		Type:      txn.Type,
		Hash:      txn.Hash().Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fee,

		Sigs: sigs,
		In:   txnInputs,
		Out:  out,
	}, nil
}

// TransactionVerbose has readable transaction data. It adds TransactionStatus to a BlockTransactionVerbose
type TransactionVerbose struct {
	Status    *TransactionStatus `json:"status,omitempty"`
	Timestamp uint64             `json:"timestamp,omitempty"`
	BlockTransactionVerbose
}

// NewTransactionVerbose creates TransactionVerbose
func NewTransactionVerbose(txn visor.Transaction, inputs []visor.TransactionInput) (TransactionVerbose, error) {
	rb, err := NewBlockTransactionVerbose(txn.Transaction, inputs, txn.Status.BlockSeq == 0 && txn.Status.Confirmed)
	if err != nil {
		return TransactionVerbose{}, err
	}

	status := NewTransactionStatus(txn.Status)

	return TransactionVerbose{
		Status:                  &status,
		Timestamp:               txn.Time,
		BlockTransactionVerbose: rb,
	}, nil
}

// UnconfirmedTransactionVerbose represents a verbose readable unconfirmed transaction
type UnconfirmedTransactionVerbose struct {
	Transaction BlockTransactionVerbose `json:"transaction"`
	Received    time.Time               `json:"received"`
	Checked     time.Time               `json:"checked"`
	Announced   time.Time               `json:"announced"`
	IsValid     bool                    `json:"is_valid"`
}

// NewUnconfirmedTransactionVerbose creates a verbose readable unconfirmed transaction
func NewUnconfirmedTransactionVerbose(unconfirmed *visor.UnconfirmedTransaction, inputs []visor.TransactionInput) (*UnconfirmedTransactionVerbose, error) {
	isGenesis := false // The genesis transaction is never unconfirmed
	txn, err := NewBlockTransactionVerbose(unconfirmed.Transaction, inputs, isGenesis)
	if err != nil {
		return nil, err
	}

	return &UnconfirmedTransactionVerbose{
		Transaction: txn,
		Received:    timeutil.NanoToTime(unconfirmed.Received),
		Checked:     timeutil.NanoToTime(unconfirmed.Checked),
		Announced:   timeutil.NanoToTime(unconfirmed.Announced),
		IsValid:     unconfirmed.IsValid == 1,
	}, nil
}

// NewUnconfirmedTransactionsVerbose creates []UnconfirmedTransactions from []UnconfirmedTransaction and their readable transaction inputs
func NewUnconfirmedTransactionsVerbose(txns []visor.UnconfirmedTransaction, inputs [][]visor.TransactionInput) ([]UnconfirmedTransactionVerbose, error) {
	if len(inputs) != len(txns) {
		return nil, fmt.Errorf("NewUnconfirmedTransactionsVerbose: len(inputs) != len(txns)")
	}

	rTxns := make([]UnconfirmedTransactionVerbose, len(txns))
	for i, txn := range txns {
		rTxn, err := NewUnconfirmedTransactionVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}

		rTxns[i] = *rTxn
	}

	return rTxns, nil
}
