package blockdb

import (
	"errors"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	logger = logging.MustGetLogger("blockdb")
	// ErrNoHeadBlock is returned when calling Blockchain.Head() when no head block exists
	ErrNoHeadBlock = fmt.Errorf("found no head block")
)

// ErrSignatureLost is returned if a block in the db does not have a corresponding signature in the db
type ErrSignatureLost struct {
	b *coin.Block
}

// NewErrSignatureLost creates ErrSignatureLost from *coin.Block
func NewErrSignatureLost(b *coin.Block) error {
	return ErrSignatureLost{
		b: b,
	}
}

func (e ErrSignatureLost) Error() string {
	return fmt.Sprintf("Signature not found for block seq=%d hash=%s", e.b.Head.BkSeq, e.b.HashHeader().Hex())
}

// BlockTree block storage
type BlockTree interface {
	AddBlock(*bolt.Tx, *coin.Block) error
	GetBlock(*bolt.Tx, cipher.SHA256) (*coin.Block, error)
	GetBlockInDepth(*bolt.Tx, uint64, Walker) (*coin.Block, error)
	ForEachBlock(*bolt.Tx, func(*coin.Block) error) error
}

// BlockSigs block signature storage
type BlockSigs interface {
	Add(*bolt.Tx, cipher.SHA256, cipher.Sig) error
	Get(*bolt.Tx, cipher.SHA256) (cipher.Sig, bool, error)
	ForEach(*bolt.Tx, func(cipher.SHA256, cipher.Sig) error) error
}

// UnspentPool unspent outputs pool
type UnspentPool interface {
	Len(*bolt.Tx) (uint64, error)
	Contains(*bolt.Tx, cipher.SHA256) (bool, error)
	Get(*bolt.Tx, cipher.SHA256) (*coin.UxOut, error)
	GetAll(*bolt.Tx) (coin.UxArray, error)
	GetArray(*bolt.Tx, []cipher.SHA256) (coin.UxArray, error)
	GetUxHash(*bolt.Tx) (cipher.SHA256, error)
	GetUnspentsOfAddrs(*bolt.Tx, []cipher.Address) (coin.AddressUxOuts, error)
	ProcessBlock(*bolt.Tx, *coin.SignedBlock) error
	// GetForTransactionInputs(*bolt.Tx, coin.Transactions) (coin.TransactionUnspents, error)
}

// Blockchain maintain the buckets for blockchain
type Blockchain struct {
	db      *dbutil.DB
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
func NewBlockchain(db *dbutil.DB, walker Walker) (*Blockchain, error) {
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

func createBlockchain(db *dbutil.DB, walker Walker, tree BlockTree, sigs BlockSigs, unspent UnspentPool) (*Blockchain, error) {
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
	if err := bc.updateWithTx(tx, bc.updateHeadSeq(b)); err != nil {
		return err
	}

	if err := bc.unspent.ProcessBlock(tx, b); err != nil {
		return err
	}

	return bc.updateWithTx(tx, bc.cacheGenesisBlock(b))
}

// Head returns head block, returns error if no block does exist
func (bc *Blockchain) Head(tx *bolt.Tx) (*coin.SignedBlock, error) {
	b, err := bc.GetSignedBlockBySeq(tx, bc.HeadSeq())
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, ErrNoHeadBlock
	}

	return b, nil
}

// UnspentPool returns the unspent pool
func (bc *Blockchain) UnspentPool() UnspentPool {
	return bc.unspent
}

// HeadSeq returns the head block sequence
func (bc *Blockchain) HeadSeq() uint64 {
	bc.RLock()
	defer bc.RUnlock()
	return bc.cache.headSeq
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

// GetBlockSignature returns the signature of a block
func (bc *Blockchain) GetBlockSignature(tx *bolt.Tx, b *coin.Block) (cipher.Sig, bool, error) {
	return bc.sigs.Get(tx, b.HashHeader())
}

// GetBlockByHash returns block of given hash
func (bc *Blockchain) GetBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.Block, error) {
	b, err := bc.tree.GetBlock(tx, hash)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetSignedBlockByHash returns signed block of given hash
func (bc *Blockchain) GetSignedBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
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
		return nil, NewErrSignatureLost(b)
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// GetSignedBlockBySeq returns signed block of given seq
func (bc *Blockchain) GetSignedBlockBySeq(tx *bolt.Tx, seq uint64) (*coin.SignedBlock, error) {
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
		return nil, NewErrSignatureLost(b)
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// ForEachSignature iterates all signatures and calls f on them
func (bc *Blockchain) ForEachSignature(tx *bolt.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
	return bc.sigs.ForEach(tx, f)
}

// ForEachBlock iterates all blocks and calls f on them
func (bc *Blockchain) ForEachBlock(tx *bolt.Tx, f func(b *coin.Block) error) error {
	return bc.tree.ForEachBlock(tx, f)
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
		b, err := bc.GetSignedBlockBySeq(tx, 0)
		if err != nil {
			return err
		}

		bc.cache.genesisBlock = b
	}
	return nil
}

func (bc *Blockchain) updateWithTx(tx *bolt.Tx, ps ...dbutil.TxHandler) error {
	rollbackFuncs := []dbutil.Rollback{}
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

func (bc *Blockchain) updateHeadSeq(b *coin.SignedBlock) dbutil.TxHandler {
	return func(tx *bolt.Tx) (dbutil.Rollback, error) {
		// meta := chainMeta{tx.Bucket(bc.meta.Name)}
		if err := bc.meta.setHeadSeq(tx, b.Seq()); err != nil {
			return func() {}, err
		}

		bc.Lock()
		defer bc.Unlock()
		// get current head seq
		seq := bc.cache.headSeq

		// update the cache head seq
		bc.cache.headSeq = b.Seq()

		return func() {
			// reset the cache head seq
			bc.Lock()
			defer bc.Unlock()
			bc.cache.headSeq = seq
		}, nil
	}
}

// cacheGenesisBlock will cache genesis block if the current block is genesis
func (bc *Blockchain) cacheGenesisBlock(b *coin.SignedBlock) dbutil.TxHandler {
	return func(tx *bolt.Tx) (dbutil.Rollback, error) {
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
