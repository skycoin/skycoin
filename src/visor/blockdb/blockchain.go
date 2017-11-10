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

// BlockTree block storage
type BlockTree interface {
	AddBlock(*bolt.Tx, *coin.Block) error
	GetBlock(*bolt.Tx, cipher.SHA256) (*coin.Block, error)
	GetBlockInDepth(*bolt.Tx, uint64, Walker) (*coin.Block, error)
}

// BlockSigs block signature storage
type BlockSigs interface {
	Add(*bolt.Tx, cipher.SHA256, cipher.Sig) error
	Get(*bolt.Tx, cipher.SHA256) (cipher.Sig, bool, error)
}

// UnspentPool unspent outputs pool
type UnspentPool interface {
	Len() uint64                          // Len returns the length of unspent outputs pool
	Get(cipher.SHA256) (coin.UxOut, bool) // Get returns outpus
	GetAll() (coin.UxArray, error)
	GetArray([]cipher.SHA256) (coin.UxArray, error)
	GetUxHash() cipher.SHA256
	GetUnspentsOfAddrs([]cipher.Address) coin.AddressUxOuts
	ProcessBlock(*coin.SignedBlock) bucket.TxHandler
	Contains(cipher.SHA256) bool
}

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
		return nil, fmt.Errorf("NewUnspentPool failed: %v", err)
	}

	tree, err := newBlockTree(db)
	if err != nil {
		return nil, fmt.Errorf("newBlockTree failed: %v", err)
	}

	sigs, err := newBlockSigs(db)
	if err != nil {
		return nil, fmt.Errorf("newBlockSigs failed: %v", err)
	}

	return createBlockchain(db, walker, tree, sigs, unspent)
}

func createBlockchain(db *bolt.DB, walker Walker, tree BlockTree, sigs BlockSigs, unspent UnspentPool) (*Blockchain, error) {
	meta, err := newChainMeta(db)
	if err != nil {
		return nil, fmt.Errorf("newChainMeta failed: %v", err)
	}

	bc := &Blockchain{
		db:      db,
		unspent: unspent,
		meta:    meta,
		tree:    tree,
		sigs:    sigs,
		walker:  walker,
	}

	if err := db.View(func(tx *bolt.Tx) error {
		return bc.syncCache(tx)
	}); err != nil {
		return nil, err
	}

	return bc, nil
}

// AddBlock adds signed block
func (bc *Blockchain) AddBlock(tx *bolt.Tx, sb *coin.SignedBlock) error {
	if err := bc.sigs.Add(tx, sb.HashHeader(), sb.Sig); err != nil {
		return fmt.Errorf("save signature failed: %v", err)
	}

	if err := bc.tree.AddBlock(tx, &sb.Block); err != nil {
		return fmt.Errorf("save block failed: %v", err)
	}

	// update block head seq and unspent pool
	if err := bc.processBlock(tx, sb); err != nil {
		return err
	}

	return nil
}

// processBlock processes a block and updates the db
func (bc *Blockchain) processBlock(tx *bolt.Tx, b *coin.SignedBlock) error {
	return bc.updateWithTx(tx, bc.updateHeadSeq(b), bc.unspent.ProcessBlock(b), bc.cacheGenesisBlock(b))
}

// Head returns head block, returns error if no block does exist
func (bc *Blockchain) Head(tx *bolt.Tx) (*coin.SignedBlock, error) {
	b, err := bc.GetBlockBySeq(tx, bc.HeadSeq())
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
func (bc *Blockchain) GetBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
	b, err := bc.tree.GetBlock(tx, hash)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, nil
	}

	// get signature
	sig, ok, err := bc.sigs.Get(tx, hash)
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
func (bc *Blockchain) GetBlockBySeq(tx *bolt.Tx, seq uint64) (*coin.SignedBlock, error) {
	b, err := bc.tree.GetBlockInDepth(tx, seq, bc.walker)
	if err != nil {
		return nil, fmt.Errorf("bc.tree.GetBlockInDepth failed: %v", err)
	}
	if b == nil {
		return nil, nil
	}

	sig, ok, err := bc.sigs.Get(tx, b.HashHeader())
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

func (bc *Blockchain) syncCache(tx *bolt.Tx) error {
	// update head seq cache
	bc.Lock()
	defer bc.Unlock()
	headSeq, err := bc.meta.getHeadSeq(tx)
	if err != nil {
		return err
	}

	bc.cache.headSeq = headSeq

	// load genesis block
	if bc.cache.genesisBlock == nil {
		b, err := bc.GetBlockBySeq(tx, 0)
		if err != nil {
			return err
		}

		bc.cache.genesisBlock = b
	}
	return nil
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
		if err := bc.meta.setHeadSeq(tx, b.Seq()); err != nil {
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
			defer bc.Unlock()
			bc.cache.genesisBlock = originGenesisBlock
		}, nil
	}
}
