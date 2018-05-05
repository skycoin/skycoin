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
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// lastTxNum reprsents the number of transactions that the GetLastTxs function will return.
const lastTxNum = 20

var transactionsBkt = []byte("transactions")

// Transactions transaction bucket instance.
type transactions struct {
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
func newTransactions(db *dbutil.DB) (*transactions, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			transactionsBkt,
		})
	}); err != nil {
		return nil, err
	}

	return &transactions{}, nil
}

// Add transaction to the db.
func (txs *transactions) Add(tx *bolt.Tx, t *Transaction) error {
	// TODO -- cached data does not rollback on error, remove it
	// Use a sequence counter to store each hash in a bucket
	// And iterate this bucket backwards (using tx.Counter()) up to lastTxNum
	txs.lastTxs = append(txs.lastTxs, t.Hash())
	if len(txs.lastTxs) > lastTxNum {
		txs.lastTxs = txs.lastTxs[1:]
	}

	hash := t.Hash()
	return dbutil.PutBucketValue(tx, transactionsBkt, hash[:], encoder.Serialize(t))
}

// Get gets transaction by tx hash, return nil on not found.
func (txs *transactions) Get(tx *bolt.Tx, hash cipher.SHA256) (*Transaction, error) {
	var txn Transaction

	if ok, err := dbutil.GetBucketObjectDecoded(tx, transactionsBkt, hash[:], &txn); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &txn, nil
}

// GetSlice returns transactions slice of given hashes
func (txs *transactions) GetSlice(tx *bolt.Tx, hashes []cipher.SHA256) ([]Transaction, error) {
	var txns []Transaction
	for _, h := range hashes {
		var txn Transaction

		if ok, err := dbutil.GetBucketObjectDecoded(tx, transactionsBkt, h[:], &txn); err != nil {
			return nil, err
		} else if !ok {
			continue
		}

		txns = append(txns, txn)
	}

	return txns, nil
}

// IsEmpty checks if transaction bucket is empty
func (txs *transactions) IsEmpty(tx *bolt.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, transactionsBkt)
}

// Reset resets the bucket
func (txs *transactions) Reset(tx *bolt.Tx) error {
	return dbutil.Reset(tx, transactionsBkt)
}

// GetLastTxs get latest tx hash set.
func (txs *transactions) GetLastTxs() []cipher.SHA256 {
	return txs.lastTxs
}
