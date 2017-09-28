package blockdb

import (
	"errors"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	// blockchain meta info bucket
	blockchainMetaBkt = []byte("blockchain_meta")
	// blockchain head sequence number
	headSeqKey = []byte("head_seq")
)

type chainMeta struct {
	bucket.Bucket
}

func newChainMeta(db *bolt.DB) (*chainMeta, error) {
	bkt, err := bucket.New(blockchainMetaBkt, db)
	if err != nil {
		return nil, err
	}

	return &chainMeta{
		Bucket: *bkt,
	}, nil
}

func (m chainMeta) setHeadSeqWithTx(tx *bolt.Tx, seq uint64) error {
	return m.PutWithTx(tx, headSeqKey, bucket.Itob(seq))
}

// BlockTree block storage
type BlockTree interface {
	AddBlockWithTx(tx *bolt.Tx, b *coin.Block) error
	GetBlock(hash cipher.SHA256) *coin.Block
	GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block
}

// BlockSigs block signature storage
type BlockSigs interface {
	AddWithTx(*bolt.Tx, cipher.SHA256, cipher.Sig) error
	Get(hash cipher.SHA256) (cipher.Sig, bool, error)
}

// UnspentPool unspent outputs pool
type UnspentPool interface {
	Len() uint64                          // Len returns the length of unspent outputs pool
	Get(cipher.SHA256) (coin.UxOut, bool) // Get returns outpus
	GetAll() (coin.UxArray, error)
	GetArray(hashes []cipher.SHA256) (coin.UxArray, error)
	GetUxHash() cipher.SHA256
	GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts
	ProcessBlock(*coin.SignedBlock) bucket.TxHandler
	Contains(cipher.SHA256) bool
}

// Walker function for go through blockchain
type Walker func(hps []coin.HashPair) cipher.SHA256

// Blockchain maintain the buckets for blockchain
type Blockchain struct {
	db      *bolt.DB
	meta    *chainMeta
	unspent UnspentPool
	tree    BlockTree
	sigs    BlockSigs
	walker  Walker
	cache   struct {
		headSeq      uint64 // head block seq
		genesisBlock *coin.SignedBlock
	}
	sync.RWMutex // cache lock
}

// NewBlockchain creates a new blockchain instance
func NewBlockchain(db *bolt.DB, walker Walker) (*Blockchain, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if walker == nil {
		return nil, errors.New("blockchain walker is nil")
	}

	unspent, err := NewUnspentPool(db)
	if err != nil {
		return nil, err
	}

	tree, err := newBlockTree(db)
	if err != nil {
		return nil, err
	}

	sigs, err := NewBlockSigs(db)
	if err != nil {
		return nil, err
	}

	return createBlockchain(db, walker, tree, sigs, unspent)
}

func createBlockchain(db *bolt.DB,
	walker Walker,
	tree BlockTree,
	sigs BlockSigs,
	unspent UnspentPool,
) (*Blockchain, error) {
	meta, err := newChainMeta(db)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db:      db,
		unspent: unspent,
		meta:    meta,
		tree:    tree,
		sigs:    sigs,
		walker:  walker,
	}

	if err := bc.syncCache(); err != nil {
		return nil, err
	}

	return bc, nil
}

// AddBlockWithTx adds signed block
func (bc *Blockchain) AddBlockWithTx(tx *bolt.Tx, sb *coin.SignedBlock) error {
	if err := bc.sigs.AddWithTx(tx, sb.HashHeader(), sb.Sig); err != nil {
		return fmt.Errorf("save signature failed: %v", err)
	}

	if err := bc.tree.AddBlockWithTx(tx, &sb.Block); err != nil {
		return fmt.Errorf("save block failed: %v", err)
	}

	// update block head seq and unspent pool
	if err := bc.processBlockWithTx(tx, sb); err != nil {
		return err
	}

	return nil
}

// processBlockWithTx process block with *bolt.Tx
func (bc *Blockchain) processBlockWithTx(tx *bolt.Tx, b *coin.SignedBlock) error {
	return bc.updateWithTx(tx,
		bc.updateHeadSeq(b),
		bc.unspent.ProcessBlock(b),
		bc.cacheGenesisBlock(b))
}

// Head returns head block, returns error if no block does exist
func (bc *Blockchain) Head() (*coin.SignedBlock, error) {
	b, err := bc.GetBlockBySeq(bc.HeadSeq())
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, fmt.Errorf("found no head block: %v", bc.HeadSeq())
	}

	return b, nil
}

// HeadSeq returns the head block sequence
func (bc *Blockchain) HeadSeq() uint64 {
	bc.RLock()
	defer bc.RUnlock()
	return bc.cache.headSeq
}

// UnspentPool returns the unspent pool
func (bc *Blockchain) UnspentPool() UnspentPool {
	return bc.unspent
}

// Len returns blockchain length
func (bc *Blockchain) Len() uint64 {
	bc.RLock()
	defer bc.RUnlock()
	if bc.cache.genesisBlock == nil {
		return 0
	}
	return uint64(bc.cache.headSeq + 1)
}

// GetBlockByHash returns signed block of given hash
func (bc *Blockchain) GetBlockByHash(hash cipher.SHA256) (*coin.SignedBlock, error) {
	b := bc.tree.GetBlock(hash)
	if b == nil {
		return nil, nil
	}

	// get signature
	sig, ok, err := bc.sigs.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("find signature of block: %v failed: %v", hash.Hex(), err)
	}

	if !ok {
		return nil, fmt.Errorf("find no signature of block: %v", hash.Hex())
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// GetBlockBySeq returns signed block of given seq
func (bc *Blockchain) GetBlockBySeq(seq uint64) (*coin.SignedBlock, error) {
	b := bc.tree.GetBlockInDepth(seq, bc.walker)
	if b == nil {
		return nil, nil
	}

	sig, ok, err := bc.sigs.Get(b.HashHeader())
	if err != nil {
		return nil, fmt.Errorf("find signature of block: %v failed: %v", seq, err)
	}

	if !ok {
		return nil, fmt.Errorf("find no signature of block: %v", seq)
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// GetGenesisBlock returns genesis block
func (bc *Blockchain) GetGenesisBlock() *coin.SignedBlock {
	bc.RLock()
	defer bc.RUnlock()
	return bc.cache.genesisBlock
}

func (bc *Blockchain) syncCache() error {
	// update head seq cache
	bc.Lock()
	defer bc.Unlock()
	bc.cache.headSeq = bc.getHeadSeqFromDB()

	// load genesis block
	if bc.cache.genesisBlock == nil {
		b, err := bc.GetBlockBySeq(0)
		if err != nil {
			return err
		}

		bc.cache.genesisBlock = b
	}
	return nil
}

func (bc *Blockchain) getHeadSeqFromDB() uint64 {
	if v := bc.meta.Get(headSeqKey); v != nil {
		return bucket.Btoi(v)
	}

	return 0
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

func (bc *Blockchain) updateHeadSeq(b *coin.SignedBlock) bucket.TxHandler {
	return func(tx *bolt.Tx) (bucket.Rollback, error) {
		// meta := chainMeta{tx.Bucket(bc.meta.Name)}
		if err := bc.meta.setHeadSeqWithTx(tx, b.Seq()); err != nil {
			return func() {}, err
		}

		bc.Lock()
		// get current head seq
		seq := bc.cache.headSeq

		// update the cache head seq
		bc.cache.headSeq = b.Seq()
		bc.Unlock()

		return func() {
			// reset the cache head seq
			bc.Lock()
			bc.cache.headSeq = seq
			bc.Unlock()
		}, nil
	}
}

// cacheGenesisBlock will cache genesis block if the current block is genesis
func (bc *Blockchain) cacheGenesisBlock(b *coin.SignedBlock) bucket.TxHandler {
	return func(tx *bolt.Tx) (bucket.Rollback, error) {
		bc.Lock()
		defer bc.Unlock()

		seq := bc.cache.headSeq
		originGenesisBlock := bc.cache.genesisBlock
		if seq == 0 {
			bc.cache.genesisBlock = b
		}

		return func() {
			bc.Lock()
			bc.cache.genesisBlock = originGenesisBlock
			bc.Unlock()
		}, nil
	}
}
