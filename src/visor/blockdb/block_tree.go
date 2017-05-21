package blockdb

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	emptyHash      cipher.SHA256
	errBlockExist  = errors.New("block already exist")
	errNoParent    = errors.New("block is not genesis and have no parent")
	errWrongParent = errors.New("wrong parent")
	errHasChild    = errors.New("remove block failed, it has children")
)

// BlockTree use the blockdb store all blocks and maintains the block tree struct.
type BlockTree struct {
	db     *bolt.DB
	blocks *bucket.Bucket
	tree   *bucket.Bucket
}

// NewBlockTree create buckets in blockdb if does not exist.
func NewBlockTree(db *bolt.DB) *BlockTree {
	blocks, err := bucket.New([]byte("blocks"), db)
	if err != nil {
		panic(err)
	}

	tree, err := bucket.New([]byte("block_tree"), db)
	if err != nil {
		panic(err)
	}

	return &BlockTree{
		blocks: blocks,
		tree:   tree,
		db:     db,
	}
}

// AddBlock write the block into blocks bucket, add the pair of block hash and pre block hash into
// tree in the block depth.
func (bt *BlockTree) AddBlock(b *coin.Block) error {
	return bt.db.Update(func(tx *bolt.Tx) error {
		blocks := tx.Bucket(bt.blocks.Name)

		// can't store block if it's not genesis block and has no parent.
		if b.Seq() > 0 && b.PreHashHeader() == emptyHash {
			return errNoParent
		}

		// check if the block already exist.
		hash := b.HashHeader()
		if blk := blocks.Get(hash[:]); blk != nil {
			return errBlockExist
		}

		// write block into blocks bucket.
		if err := setBlock(blocks, b); err != nil {
			return err
		}

		// get tree bucket.
		tree := tx.Bucket(bt.tree.Name)

		// the pre hash must be in depth - 1.
		if b.Seq() > 0 {
			preHash := b.PreHashHeader()
			parentHashPair, err := getHashPairInDepth(tree, b.Seq()-1, func(hp coin.HashPair) bool {
				return hp.Hash == preHash
			})
			if err != nil {
				return err
			}
			if len(parentHashPair) == 0 {
				return errWrongParent
			}
		}

		hp := coin.HashPair{Hash: hash, PreHash: b.Head.PrevHash}

		// get block pairs in the depth
		hashPairs, err := getHashPairInDepth(tree, b.Seq(), allPairs)
		if err != nil {
			return err
		}

		if len(hashPairs) == 0 {
			// no hash pair exist in the depth.
			// write the hash pair into tree.
			return setHashPairInDepth(tree, b.Seq(), []coin.HashPair{hp})
		}

		// check dup block
		if containHash(hashPairs, hp) {
			return errBlockExist
		}

		hashPairs = append(hashPairs, hp)
		return setHashPairInDepth(tree, b.Seq(), hashPairs)
	})
}

// RemoveBlock remove block from blocks bucket and tree bucket.
// can't remove block if it has children.
func (bt *BlockTree) RemoveBlock(b *coin.Block) error {
	return bt.db.Update(func(tx *bolt.Tx) error {
		// delete block in blocks bucket.
		blocks := tx.Bucket(bt.blocks.Name)
		hash := b.HashHeader()
		if err := blocks.Delete(hash[:]); err != nil {
			return err
		}

		// get tree bucket.
		tree := tx.Bucket(bt.tree.Name)

		// check if this block has children
		has, err := hasChild(tree, *b)
		if err != nil {
			return err
		}
		if has {
			return errHasChild
		}

		// get block hash pairs in depth
		hashPairs, err := getHashPairInDepth(tree, b.Seq(), func(hp coin.HashPair) bool {
			return true
		})
		if err != nil {
			return err
		}

		// remove block hash pair in tree.
		ps := removePairs(hashPairs, coin.HashPair{Hash: hash, PreHash: b.PreHashHeader()})
		if len(ps) == 0 {
			tree.Delete(itob(b.Seq()))
			return nil
		}

		// update the hash pairs in tree.
		return setHashPairInDepth(tree, b.Seq(), ps)
	})
}

// GetBlock get block by hash, return nil on not found
func (bt *BlockTree) GetBlock(hash cipher.SHA256) *coin.Block {
	return bt.getBlock(hash)
}

// GetBlockInDepth get block in depth, return nil on not found,
// the filter is used to choose the appropriate block.
func (bt *BlockTree) GetBlockInDepth(depth uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block {
	hash, err := bt.getHashInDepth(depth, filter)
	if err != nil {
		return nil
	}

	return bt.getBlock(hash)
}

// GetAllBlockHashInDepth returns all block hash of N depth in the tree.
func (bt *BlockTree) GetAllBlockHashInDepth(depth uint64) ([]cipher.SHA256, error) {
	key := itob(depth)
	pairsBin := bt.tree.Get(key)
	pairs := []coin.HashPair{}
	if err := encoder.DeserializeRaw(pairsBin, &pairs); err != nil {
		return []cipher.SHA256{}, err
	}
	hashes := make([]cipher.SHA256, len(pairs))
	for i, hp := range pairs {
		hashes[i] = hp.Hash
	}
	return hashes, nil
}

func (bt *BlockTree) getBlock(hash cipher.SHA256) *coin.Block {
	bin := bt.blocks.Get(hash[:])
	if bin == nil {
		return nil
	}
	block := coin.Block{}
	if err := encoder.DeserializeRaw(bin, &block); err != nil {
		return nil
	}
	return &block
}

func (bt *BlockTree) getHashInDepth(depth uint64, filter func(ps []coin.HashPair) cipher.SHA256) (cipher.SHA256, error) {
	key := itob(depth)
	pairsBin := bt.tree.Get(key)
	pairs := []coin.HashPair{}
	if err := encoder.DeserializeRaw(pairsBin, &pairs); err != nil {
		return cipher.SHA256{}, err
	}

	hash := filter(pairs)
	return hash, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func containHash(hashPairs []coin.HashPair, pair coin.HashPair) bool {
	for _, p := range hashPairs {
		if p.Hash == pair.Hash {
			return true
		}
	}
	return false
}

func removePairs(hps []coin.HashPair, pair coin.HashPair) []coin.HashPair {
	pairs := []coin.HashPair{}
	for _, p := range hps {
		if p.Hash == pair.Hash && p.PreHash == pair.PreHash {
			continue
		}
		pairs = append(pairs, p)
	}
	return pairs
}

func getHashPairInDepth(tree *bolt.Bucket, dep uint64, fn func(hp coin.HashPair) bool) ([]coin.HashPair, error) {
	v := tree.Get(itob(dep))
	if v == nil {
		return []coin.HashPair{}, nil
	}

	hps := []coin.HashPair{}
	if err := encoder.DeserializeRaw(v, &hps); err != nil {
		return nil, err
	}
	pairs := []coin.HashPair{}
	for _, ps := range hps {
		if fn(ps) {
			pairs = append(pairs, ps)
		}
	}
	return pairs, nil
}

func setBlock(bkt *bolt.Bucket, b *coin.Block) error {
	bin := encoder.Serialize(b)
	key := b.HashHeader()
	return bkt.Put(key[:], bin)
}

// check if this block has children
func hasChild(bkt *bolt.Bucket, b coin.Block) (bool, error) {
	// get the child block hash pair, whose pre hash point to current block.
	childHashPair, err := getHashPairInDepth(bkt, b.Head.BkSeq+1, func(hp coin.HashPair) bool {
		return hp.PreHash == b.HashHeader()
	})

	if err != nil {
		return false, nil
	}

	return len(childHashPair) > 0, nil
}

func setHashPairInDepth(bkt *bolt.Bucket, dep uint64, hps []coin.HashPair) error {
	hpsBin := encoder.Serialize(hps)
	key := itob(dep)
	return bkt.Put(key, hpsBin)
}

func allPairs(hp coin.HashPair) bool {
	return true
}
