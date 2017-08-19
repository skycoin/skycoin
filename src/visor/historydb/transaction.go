package historydb

// transaction.go mainly provides transaction corresponding buckets and apis,
// The transactions bucket, tx hash as key, and tx as value, it's the main bucket that stores the
// transaction value. All other buckets that index different field of transaction will only records the
// transaction hash, and get the tx value from transactions bucket.

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// lastTxNum reprsents the number of transactions that the GetLastTxs function will return.
const lastTxNum = 20

// Transactions transaction bucket instance.
type transactions struct {
	bkt     *bucket.Bucket
	lastTxs []cipher.SHA256 // records the latest transactions
}

// Transaction contains transaction info and the seq of block which executed this block.
type Transaction struct {
	Tx       coin.Transaction
	BlockSeq uint64
}

// Hash return the Tx hash.
func (tx *Transaction) Hash() cipher.SHA256 {
	return tx.Tx.Hash()
}

// New create a transaction db instance.
func newTransactionsBkt(db *bolt.DB) (*transactions, error) {
	txBkt, err := bucket.New([]byte("transactions"), db)
	if err != nil {
		return nil, nil
	}

	return &transactions{bkt: txBkt}, nil
}

func addTransaction(b *bolt.Bucket, tx *Transaction) error {
	hash := tx.Hash()
	return b.Put(hash[:], encoder.Serialize(tx))
}

// Add transaction to the db.
func (txs *transactions) Add(t *Transaction) error {
	txs.lastTxs = append(txs.lastTxs, t.Hash())
	if len(txs.lastTxs) > lastTxNum {
		txs.lastTxs = txs.lastTxs[1:]
	}

	key := t.Hash()
	v := encoder.Serialize(t)
	return txs.bkt.Put(key[:], v)
}

// Get get transaction by tx hash, return nil on not found.
func (txs *transactions) Get(hash cipher.SHA256) (*Transaction, error) {
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

// GetSlice returns transactions slice of given hashes
func (txs *transactions) GetSlice(hashes []cipher.SHA256) ([]Transaction, error) {
	keys := [][]byte{}
	for i := range hashes {
		keys = append(keys, hashes[i][:])
	}

	vs := txs.bkt.GetSlice(keys)
	txns := make([]Transaction, 0, len(vs))
	for i := range vs {
		var tx Transaction
		if err := encoder.DeserializeRaw(vs[i], &tx); err != nil {
			return []Transaction{}, err
		}
		txns = append(txns, tx)
	}

	return txns, nil
}

// IsEmpty checks if transaction bucket is empty
func (txs *transactions) IsEmpty() bool {
	return txs.bkt.IsEmpty()
}

// Reset resets the bucket
func (txs *transactions) Reset() error {
	return txs.bkt.Reset()
}

// GetLastTxs get latest tx hash set.
func (txs *transactions) GetLastTxs() []cipher.SHA256 {
	return txs.lastTxs
}

func (txs *transactions) updateLastTxs(hash cipher.SHA256) {
	txs.lastTxs = append(txs.lastTxs, hash)
	if len(txs.lastTxs) > lastTxNum {
		txs.lastTxs = txs.lastTxs[1:]
	}
}
