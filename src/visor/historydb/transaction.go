package historydb

// transaction.go mainly provides transaction corresponding buckets and apis,
// The transactions bucket, tx hash as key, and tx as value, it's the main bucket that stores the
// transaction value. All other buckets that index different field of transaction will only records the
// transaction hash, and get the tx value from transactions bucket.

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// TransactionsBkt holds Transactions
var TransactionsBkt = []byte("transactions")

// Transactions transaction bucket instance.
type transactions struct{}

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
	return &transactions{}, nil
}

// Add transaction to the db.
func (txs *transactions) Add(tx *dbutil.Tx, txn *Transaction) error {
	hash := txn.Hash()
	return dbutil.PutBucketValue(tx, TransactionsBkt, hash[:], encoder.Serialize(txn))
}

// Get gets transaction by tx hash, return nil on not found.
func (txs *transactions) Get(tx *dbutil.Tx, hash cipher.SHA256) (*Transaction, error) {
	var txn Transaction

	if ok, err := dbutil.GetBucketObjectDecoded(tx, TransactionsBkt, hash[:], &txn); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &txn, nil
}

// GetSlice returns transactions slice of given hashes
func (txs *transactions) GetSlice(tx *dbutil.Tx, hashes []cipher.SHA256) ([]Transaction, error) {
	var txns []Transaction
	for _, h := range hashes {
		var txn Transaction

		if ok, err := dbutil.GetBucketObjectDecoded(tx, TransactionsBkt, h[:], &txn); err != nil {
			return nil, err
		} else if !ok {
			continue
		}

		txns = append(txns, txn)
	}

	return txns, nil
}

// IsEmpty checks if transaction bucket is empty
func (txs *transactions) IsEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, TransactionsBkt)
}

// Reset resets the bucket
func (txs *transactions) Reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, TransactionsBkt)
}

// ForEach traverses the transactions in db
func (txs *transactions) ForEach(tx *dbutil.Tx, f func(cipher.SHA256, *Transaction) error) error {
	return dbutil.ForEach(tx, TransactionsBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromBytes(k)
		if err != nil {
			return err
		}

		var txn Transaction
		if err := encoder.DeserializeRaw(v, &txn); err != nil {
			return err
		}

		return f(hash, &txn)
	})
}
