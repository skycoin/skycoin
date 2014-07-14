package visor

import (
	"log"

	"github.com/skycoin/skycoin/src/coin"
)

// Encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block's header
	Head ReadableBlockHeader `json:"head"`
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64 `json:"unspents"`
	// Number of known unconfirmed txns
	Unconfirmed uint64 `json:"unconfirmed"`
}

func NewBlockchainMetadata(v *Visor) BlockchainMetadata {
	head := v.blockchain.Head().Head
	return BlockchainMetadata{
		Head:        NewReadableBlockHeader(&head),
		Unspents:    uint64(len(v.blockchain.Unspent.Pool)),
		Unconfirmed: uint64(len(v.Unconfirmed.Txns)),
	}
}

// Wrapper around coin.Transaction, tagged with its status.  This allows us
// to include unconfirmed txns
type Transaction struct {
	Txn    coin.Transaction  `json:"txn"`
	Status TransactionStatus `json:"status"`
}

type TransactionStatus struct {
	// This txn is in the unconfirmed pool
	Unconfirmed bool `json:"unconfirmed"`
	// We can't find anything about this txn.  Be aware that the txn may be
	// in someone else's unconfirmed pool, and if valid, it may become a
	// confirmed txn in the future
	Unknown   bool `json:"unknown"`
	Confirmed bool `json:"confirmed"`
	// If confirmed, how many blocks deep in the chain it is. Will be at least
	// 1 if confirmed.
	Height uint64 `json:"height"`
}

func NewUnconfirmedTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Unconfirmed: true,
		Unknown:     false,
		Confirmed:   false,
		Height:      0,
	}
}

func NewUnknownTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Unconfirmed: false,
		Unknown:     true,
		Confirmed:   false,
		Height:      0,
	}
}

func NewConfirmedTransactionStatus(height uint64) TransactionStatus {
	if height == 0 {
		log.Panic("Invalid confirmed transaction height")
	}
	return TransactionStatus{
		Unconfirmed: false,
		Unknown:     false,
		Confirmed:   true,
		Height:      height,
	}
}

type ReadableTransactionHeader struct {
	Hash string   `json:"hash"`
	Sigs []string `json:"sigs"`
}

func NewReadableTransactionHeader(t *coin.TransactionHeader) ReadableTransactionHeader {
	sigs := make([]string, len(t.Sigs))
	for i, _ := range t.Sigs {
		sigs[i] = t.Sigs[i].Hex()
	}
	return ReadableTransactionHeader{
		Hash: t.Hash.Hex(),
		Sigs: sigs,
	}
}

type ReadableTransactionOutput struct {
	Address string `json:"dst"`
	Coins   uint64 `json:"coins"`
	Hours   uint64 `json:"hours"`
}

func NewReadableTransactionOutput(t *coin.TransactionOutput) ReadableTransactionOutput {
	return ReadableTransactionOutput{
		Address: t.Address.String(),
		Coins:   t.Coins,
		Hours:   t.Hours,
	}
}

type ReadableTransaction struct {
	Head ReadableTransactionHeader   `json:"header"`
	In   []string                    `json:"inputs"`
	Out  []ReadableTransactionOutput `json:"outputs"`
}

func NewReadableTransaction(t *coin.Transaction) ReadableTransaction {
	in := make([]string, len(t.In))
	for i, _ := range t.In {
		in[i] = t.In[i].Hex()
	}
	out := make([]ReadableTransactionOutput, len(t.Out))
	for i, _ := range t.Out {
		out[i] = NewReadableTransactionOutput(&t.Out[i])
	}
	return ReadableTransaction{
		Head: NewReadableTransactionHeader(&t.Head),
		In:   in,
		Out:  out,
	}
}

type ReadableBlockHeader struct {
	Version  uint32 `json:"version"`
	Time     uint64 `json:"timestamp"`
	BkSeq    uint64 `json:"seq"`
	Fee      uint64 `json:"fee"`
	PrevHash string `json:"prev_hash"`
	BodyHash string `json:"hash"`
}

func NewReadableBlockHeader(b *coin.BlockHeader) ReadableBlockHeader {
	return ReadableBlockHeader{
		Version:  b.Version,
		Time:     b.Time,
		BkSeq:    b.BkSeq,
		Fee:      b.Fee,
		PrevHash: b.PrevHash.Hex(),
		BodyHash: b.BodyHash.Hex(),
	}
}

type ReadableBlockBody struct {
	Transactions []ReadableTransaction `json:"txns"`
}

func NewReadableBlockBody(b *coin.BlockBody) ReadableBlockBody {
	txns := make([]ReadableTransaction, len(b.Transactions))
	for i, _ := range b.Transactions {
		txns[i] = NewReadableTransaction(&b.Transactions[i])
	}
	return ReadableBlockBody{
		Transactions: txns,
	}
}

type ReadableBlock struct {
	Head ReadableBlockHeader `json:"header"`
	Body ReadableBlockBody   `json:"body"`
}

func NewReadableBlock(b *coin.Block) ReadableBlock {
	return ReadableBlock{
		Head: NewReadableBlockHeader(&b.Head),
		Body: NewReadableBlockBody(&b.Body),
	}
}
