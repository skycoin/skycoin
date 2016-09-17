package blockdb

import (
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
)

type Transactions struct {
	txns *Bucket
}

type TransactionBlock struct {
	TxHash    cipher.SHA256
	BlockHash cipher.SHA256
}

func NewTransactions() *Transactions {
	txns, err := NewBucket([]byte("transactions"))
	if err != nil {
		panic(err)
	}
	return &Transactions{
		txns: txns,
	}
}

func (t *Transactions) Add(tb *TransactionBlock) error {
	return t.txns.Put(tb.TxHash[:], encoder.Serialize(tb.BlockHash))
}

func (t Transactions) Get(txHash cipher.SHA256) *TransactionBlock {
	bin := t.txns.Get(txHash[:])
	if bin == nil {
		return nil
	}
	blockHash := cipher.SHA256{}
	if err := encoder.DeserializeRaw(bin, &blockHash); err != nil {
		return nil
	}

	return &TransactionBlock{
		TxHash:    txHash,
		BlockHash: blockHash,
	}
}
