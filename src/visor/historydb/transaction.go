package historydb

// transaction.go mainly provides transaction corresponding buckets and apis,
// The transactions bucket, tx hash as key, and tx as value, it's the main bucket that stores the
// transaction value. All other buckets that index different field of transaction will only record the
// transaction hash, and get the tx value from transactions bucket.

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

//go:generate skyencoder -unexported -struct Transaction

// Transaction contains transaction info and the seq of block which executed this block.
type Transaction struct {
	Txn      coin.Transaction
	BlockSeq uint64
}

// Hash return the Txn hash.
func (txn *Transaction) Hash() cipher.SHA256 {
	return txn.Txn.Hash()
}

// TransactionsBkt holds Transactions
var TransactionsBkt = []byte("transactions")

// Transactions transaction bucket instance.
type transactions struct{}

// put transaction in the db
func (txs *transactions) put(tx *dbutil.Tx, txn *Transaction) error {
	hash := txn.Hash()
	buf, err := encodeTransaction(txn)
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, TransactionsBkt, hash[:], buf)
}

// get gets transaction by transaction hash, return nil on not found
func (txs *transactions) get(tx *dbutil.Tx, hash cipher.SHA256) (*Transaction, error) {
	var txn Transaction

	v, err := dbutil.GetBucketValueNoCopy(tx, TransactionsBkt, hash[:])
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	if err := decodeTransactionExact(v, &txn); err != nil {
		return nil, err
	}

	return &txn, nil
}

// getArray returns transactions slice of given hashes
func (txs *transactions) getArray(tx *dbutil.Tx, hashes []cipher.SHA256) ([]Transaction, error) {
	txns := make([]Transaction, 0, len(hashes))
	for _, h := range hashes {
		txn, err := txs.get(tx, h)
		if err != nil {
			return nil, err
		}
		if txn == nil {
			return nil, errors.New("Transaction not found")
		}

		txns = append(txns, *txn)
	}

	return txns, nil
}

// isEmpty checks if transaction bucket is empty
func (txs *transactions) isEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, TransactionsBkt)
}

// reset resets the bucket
func (txs *transactions) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, TransactionsBkt)
}

// forEach traverses the transactions in db
func (txs *transactions) forEach(tx *dbutil.Tx, f func(cipher.SHA256, *Transaction) error) error {
	return dbutil.ForEach(tx, TransactionsBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromBytes(k)
		if err != nil {
			return err
		}

		var txn Transaction
		if err := decodeTransactionExact(v, &txn); err != nil {
			return err
		}

		return f(hash, &txn)
	})
}

// len returns the total number of all transactions
func (txs *transactions) len(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, TransactionsBkt)
}
