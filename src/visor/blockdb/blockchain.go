/*
Package blockdb is the core blockchain database wrapper
*/
package blockdb

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

var (
	logger = logging.MustGetLogger("blockdb")

	// ErrNoHeadBlock is returned when calling Blockchain.Head() when no head block exists
	ErrNoHeadBlock = fmt.Errorf("found no head block")
)

//go:generate skyencoder -unexported -struct Block -output-path . -package blockdb github.com/SkycoinProject/skycoin/src/coin
//go:generate skyencoder -unexported -struct UxOut -output-path . -package blockdb github.com/SkycoinProject/skycoin/src/coin
//go:generate skyencoder -unexported -struct hashPairsWrapper
//go:generate skyencoder -unexported -struct hashesWrapper
//go:generate skyencoder -unexported -struct sigWrapper

// hashesWrapper wraps []cipher.SHA256 so it can be used by skyencoder
type hashesWrapper struct {
	Hashes []cipher.SHA256
}

// sigWrapper wraps cipher.Sig in struct so it can be used by skyencoder
type sigWrapper struct {
	Sig cipher.Sig
}

// hashPairsWrapper wraps []coin.HashPair so it can be used by skyencoder
type hashPairsWrapper struct {
	HashPairs []coin.HashPair
}

// ErrMissingSignature is returned if a block in the db does not have a corresponding signature in the db
type ErrMissingSignature struct {
	b *coin.Block
}

// NewErrMissingSignature creates ErrMissingSignature from *coin.Block
func NewErrMissingSignature(b *coin.Block) error {
	return ErrMissingSignature{
		b: b,
	}
}

func (e ErrMissingSignature) Error() string {
	return fmt.Sprintf("Signature not found for block seq=%d hash=%s", e.b.Head.BkSeq, e.b.HashHeader().Hex())
}

// CreateBuckets creates bolt.DB buckets used by the blockdb
func CreateBuckets(tx *dbutil.Tx) error {
	return dbutil.CreateBuckets(tx, [][]byte{
		BlockSigsBkt,
		BlocksBkt,
		TreeBkt,
		BlockchainMetaBkt,
		UnspentPoolBkt,
		UnspentPoolAddrIndexBkt,
		UnspentMetaBkt,
	})
}

// BlockTree block storage
type BlockTree interface {
	AddBlock(*dbutil.Tx, *coin.Block) error
	GetBlock(*dbutil.Tx, cipher.SHA256) (*coin.Block, error)
	GetBlockInDepth(*dbutil.Tx, uint64, Walker) (*coin.Block, error)
	ForEachBlock(*dbutil.Tx, func(*coin.Block) error) error
}

// BlockSigs block signature storage
type BlockSigs interface {
	Add(*dbutil.Tx, cipher.SHA256, cipher.Sig) error
	Get(*dbutil.Tx, cipher.SHA256) (cipher.Sig, bool, error)
	ForEach(*dbutil.Tx, func(cipher.SHA256, cipher.Sig) error) error
}

//go:generate mockery -name UnspentPooler -case underscore -testonly -inpkg

// UnspentPooler unspent outputs pool
type UnspentPooler interface {
	MaybeBuildIndexes(*dbutil.Tx, uint64) error
	Len(*dbutil.Tx) (uint64, error)
	Contains(*dbutil.Tx, cipher.SHA256) (bool, error)
	Get(*dbutil.Tx, cipher.SHA256) (*coin.UxOut, error)
	GetAll(*dbutil.Tx) (coin.UxArray, error)
	GetArray(*dbutil.Tx, []cipher.SHA256) (coin.UxArray, error)
	GetUxHash(*dbutil.Tx) (cipher.SHA256, error)
	GetUnspentsOfAddrs(*dbutil.Tx, []cipher.Address) (coin.AddressUxOuts, error)
	GetUnspentHashesOfAddrs(*dbutil.Tx, []cipher.Address) (AddressHashes, error)
	ProcessBlock(*dbutil.Tx, *coin.SignedBlock) error
	AddressCount(*dbutil.Tx) (uint64, error)
}

// ChainMeta blockchain metadata
type ChainMeta interface {
	GetHeadSeq(*dbutil.Tx) (uint64, bool, error)
	SetHeadSeq(*dbutil.Tx, uint64) error
}

// Blockchain maintain the buckets for blockchain
type Blockchain struct {
	db      *dbutil.DB
	meta    ChainMeta
	unspent UnspentPooler
	tree    BlockTree
	sigs    BlockSigs
	walker  Walker
}

// NewBlockchain creates a new blockchain instance
func NewBlockchain(db *dbutil.DB, walker Walker) (*Blockchain, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	if walker == nil {
		return nil, errors.New("blockchain walker is nil")
	}

	return &Blockchain{
		db:      db,
		unspent: NewUnspentPool(),
		meta:    &chainMeta{},
		tree:    &blockTree{},
		sigs:    &blockSigs{},
		walker:  walker,
	}, nil
}

// UnspentPool returns the unspent pool
func (bc *Blockchain) UnspentPool() UnspentPooler {
	return bc.unspent
}

// AddBlock adds signed block
func (bc *Blockchain) AddBlock(tx *dbutil.Tx, sb *coin.SignedBlock) error {
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
func (bc *Blockchain) processBlock(tx *dbutil.Tx, b *coin.SignedBlock) error {
	if err := bc.unspent.ProcessBlock(tx, b); err != nil {
		return err
	}

	return bc.meta.SetHeadSeq(tx, b.Seq())
}

// Head returns head block, returns error if no head block exists
func (bc *Blockchain) Head(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	seq, ok, err := bc.HeadSeq(tx)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrNoHeadBlock
	}

	b, err := bc.GetSignedBlockBySeq(tx, seq)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, ErrNoHeadBlock
	}

	return b, nil
}

// HeadSeq returns the head block sequence
func (bc *Blockchain) HeadSeq(tx *dbutil.Tx) (uint64, bool, error) {
	return bc.meta.GetHeadSeq(tx)
}

// Len returns blockchain length
func (bc *Blockchain) Len(tx *dbutil.Tx) (uint64, error) {
	seq, ok, err := bc.meta.GetHeadSeq(tx)
	if err != nil {
		return 0, err
	} else if !ok {
		return 0, nil
	}

	return seq + 1, nil
}

// GetBlockSignature returns the signature of a block
func (bc *Blockchain) GetBlockSignature(tx *dbutil.Tx, b *coin.Block) (cipher.Sig, bool, error) {
	return bc.sigs.Get(tx, b.HashHeader())
}

// GetBlockByHash returns block of given hash
func (bc *Blockchain) GetBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.Block, error) {
	b, err := bc.tree.GetBlock(tx, hash)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GetSignedBlockByHash returns signed block of given hash
func (bc *Blockchain) GetSignedBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
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
		return nil, NewErrMissingSignature(b)
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// GetSignedBlockBySeq returns signed block of given seq
func (bc *Blockchain) GetSignedBlockBySeq(tx *dbutil.Tx, seq uint64) (*coin.SignedBlock, error) {
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
		return nil, NewErrMissingSignature(b)
	}

	return &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}, nil
}

// GetGenesisBlock returns genesis block
func (bc *Blockchain) GetGenesisBlock(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	return bc.GetSignedBlockBySeq(tx, 0)
}

// ForEachBlock iterates all blocks and calls f on them
func (bc *Blockchain) ForEachBlock(tx *dbutil.Tx, f func(b *coin.Block) error) error {
	return bc.tree.ForEachBlock(tx, f)
}
