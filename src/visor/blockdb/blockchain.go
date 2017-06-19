package blockdb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	headSeqKey = []byte("head_seq")
)

type chainMeta struct {
	*bolt.Bucket
}

func (m chainMeta) setHeadSeq(seq uint64) error {
	return m.Put(headSeqKey, bucket.Itob(seq))
}

// Blockchain maintain the buckets for blockchain
type Blockchain struct {
	db      *bolt.DB
	meta    *bucket.Bucket
	Unspent *UnspentPool
}

// NewBlockchain creates a new blockchain instance
func NewBlockchain(db *bolt.DB) (*Blockchain, error) {
	unspent, err := NewUnspentPool(db)
	if err != nil {
		return nil, err
	}

	meta, err := bucket.New([]byte("blockchain_meta"), db)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db:      db,
		Unspent: unspent,
		meta:    meta,
	}
	return bc, nil
}

// ProcessBlock processes block
func (bc *Blockchain) ProcessBlock(b *coin.Block) error {
	txns := b.Body.Transactions
	return bc.db.Update(func(tx *bolt.Tx) error {
		for _, txn := range txns {
			// Remove spent outputs
			if err := bc.Unspent.deleteWithTx(tx, txn.In); err != nil {
				return err
			}

			// Create new outputs
			txUxs := coin.CreateUnspents(b.Head, txn)
			for i := range txUxs {
				if err := bc.Unspent.addWithTx(tx, txUxs[i]); err != nil {
					return err
				}
			}
		}

		// update block head seq
		return chainMeta{tx.Bucket(bc.meta.Name)}.setHeadSeq(b.Head.BkSeq)
	})
}

// AddUxOut adds a UxOut to pool
func (bc *Blockchain) AddUxOut(ux coin.UxOut) error {
	return bc.db.Update(func(tx *bolt.Tx) error {
		meta := chainMeta{tx.Bucket(bc.meta.Name)}
		if err := bc.Unspent.addWithTx(tx, ux); err != nil {
			return err
		}

		return meta.setHeadSeq(ux.Head.BkSeq)
	})
}

// Reset resets the chain
func (bc *Blockchain) Reset() error {
	return bc.db.Update(func(tx *bolt.Tx) error {
		if err := bc.Unspent.resetWithTx(tx); err != nil {
			return err
		}

		if err := tx.DeleteBucket(bc.meta.Name); err != nil {
			return err
		}

		_, err := tx.CreateBucket(bc.meta.Name)
		return err
	})
}

// HeadSeq returns the head block sequence
func (bc *Blockchain) HeadSeq() int64 {
	if v := bc.meta.Get(headSeqKey); v != nil {
		return int64(bucket.Btoi(v))
	}

	return -1
}
