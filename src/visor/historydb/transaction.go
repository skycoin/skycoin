package historydb

// transaction.go mainly provides transaction corresponding buckets and apis,
// The transactions bucket, tx hash as key, and tx as value, it's the main bucket that stores the
// transaction value. All other buckets that index different field of transaction will only records the
// transaction hash, and get the tx value from transactions bucket.

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// Transactions transaction bucket instance.
type Transactions struct {
	bkt *bucket.Bucket
}

// Transaction contains transaction info and the seq of block which executed this block.
type Transaction struct {
	coin.Transaction
	BlockSeq uint64
}

// New create a transaction db instance.
func newTransactions(db *bolt.DB) (*Transactions, error) {
	txBkt, err := bucket.New([]byte("transactions"), db)
	if err != nil {
		return nil, nil
	}

	return &Transactions{txBkt}, nil
}

// Add transaction to the db.
func (txs *Transactions) Add(t *Transaction) error {
	key := t.Hash()
	v := encoder.Serialize(t)
	return txs.bkt.Put(key[:], v)
}

// Get get transaction by tx hash, return nil on not found.
func (txs Transactions) Get(hash cipher.SHA256) (*Transaction, error) {
	bin := txs.bkt.Get(hash[:])
	if bin == nil {
		return nil, nil
	}

	// deserialize tx
	var tx Transaction
	if err := encoder.DeserializeRaw(bin, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}
