package blockdb

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	// blockchain head sequence number
	headSeqKey = []byte("head_seq")
)

type chainMeta struct {
	*bolt.Bucket
}

func (m chainMeta) setHeadSeq(seq uint64) error {
	return m.Put(headSeqKey, bucket.Itob(seq))
}

func (m chainMeta) getHeadSeq() uint64 {
	return bucket.Btoi(m.Get(headSeqKey))
}

// Blockchain maintain the buckets for blockchain
type Blockchain struct {
	db      *bolt.DB
	meta    *bucket.Bucket
	Unspent *UnspentPool

	cache struct {
		headSeq        int64  // head block seq
		verifiedSigSeq uint64 // verified block seq
	}
	sync.Mutex // cache lock
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
	bc.cache.headSeq = -1

	bc.syncCache()

	return bc, nil
}

func (bc *Blockchain) syncCache() {
	// update head seq cache
	bc.Lock()
	bc.cache.headSeq = bc.getHeadSeqFromDB()
	bc.Unlock()
}

func (bc *Blockchain) getHeadSeqFromDB() int64 {
	if v := bc.meta.Get(headSeqKey); v != nil {
		return int64(bucket.Btoi(v))
	}
	return -1
}

// ProcessBlock processes block
func (bc *Blockchain) ProcessBlock(b *coin.Block) error {
	if err := bc.dbUpdate(
		bc.updateHeadSeq(b),
		bc.Unspent.processBlock(b)); err != nil {
		return err
	}

	return nil
}

// ProcessBlockWithTx process block with *bolt.Tx
func (bc *Blockchain) ProcessBlockWithTx(tx *bolt.Tx, b *coin.Block) error {
	return bc.updateWithTx(tx,
		bc.updateHeadSeq(b),
		bc.Unspent.processBlock(b))
}

// dbUpdate will execute all processors in sequence, return error will rollback all
// updates to the db
func (bc *Blockchain) dbUpdate(ps ...bucket.TxHandler) error {
	return bc.db.Update(func(tx *bolt.Tx) error {
		return bc.updateWithTx(tx, ps...)
	})
}

func (bc *Blockchain) updateWithTx(tx *bolt.Tx, ps ...bucket.TxHandler) error {
	rollbackFuncs := []bucket.Rollback{}
	for _, p := range ps {
		rb, err := p(tx)
		if err != nil {
			// rollback previous updates if any
			for _, r := range rollbackFuncs {
				r()
			}
			return err
		}
		rollbackFuncs = append(rollbackFuncs, rb)
	}

	return nil
}

func (bc *Blockchain) updateHeadSeq(b *coin.Block) bucket.TxHandler {
	return func(tx *bolt.Tx) (bucket.Rollback, error) {
		meta := chainMeta{tx.Bucket(bc.meta.Name)}

		if err := meta.setHeadSeq(b.Seq()); err != nil {
			return func() {}, err
		}

		bc.Lock()
		// get current head seq
		seq := bc.cache.headSeq

		// update the cache head seq
		bc.cache.headSeq = int64(b.Seq())
		bc.Unlock()

		return func() {
			// reset the cache head seq
			bc.Lock()
			bc.cache.headSeq = int64(seq)
			bc.Unlock()
		}, nil
	}
}

// HeadSeq returns the head block sequence
func (bc *Blockchain) HeadSeq() int64 {
	bc.Lock()
	defer bc.Unlock()
	return bc.cache.headSeq
}
