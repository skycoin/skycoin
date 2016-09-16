package coin

import (
	"encoding/binary"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/blockdb"
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
	blocks *blockdb.Bucket
	tree   *blockdb.Bucket
}

// NewBlockTree create buckets in blockdb if does not exist.
func NewBlockTree() *BlockTree {
	blocks, err := blockdb.NewBucket([]byte("blocks"))
	if err != nil {
		panic(err)
	}

	tree, err := blockdb.NewBucket([]byte("block_tree"))
	if err != nil {
		panic(err)
	}

	return &BlockTree{
		blocks: blocks,
		tree:   tree,
	}
}

// AddBlock write the block into blocks bucket, add the pair of block hash and pre block hash into
// tree in the block depth.
func (bt *BlockTree) AddBlock(b Block) error {
	return blockdb.UpdateTx(func(tx *bolt.Tx) error {
		blocks := tx.Bucket(bt.blocks.Name)

		// can't store block if it's not genesis block and has no parent.
		if b.Seq() > 0 && b.PreHashHeader() == emptyHash {
			return errNoParent
		}

		// write block into blocks bucket.
		if err := setBlock(blocks, &b); err != nil {
			return err
		}

		// get tree bucket.
		tree := tx.Bucket(bt.tree.Name)

		// the pre hash must be in depth - 1.
		if b.Seq() > 0 {
			preHash := b.PreHashHeader()
			parentHashPair, err := getHashPairInDepth(tree, b.Seq()-1, func(hp HashPair) bool {
				return hp.Hash == preHash
			})
			if err != nil {
				return err
			}
			if len(parentHashPair) == 0 {
				return errWrongParent
			}
		}

		hash := b.HashHeader()
		hp := HashPair{hash, b.Head.PrevHash}

		// get block pairs in the depth
		hashPairs, err := getHashPairInDepth(tree, b.Seq(), allPairs)
		if err != nil {
			return err
		}

		if len(hashPairs) == 0 {
			// no hash pair exist in the depth.
			// write the hash pair into tree.
			return setHashPairInDepth(tree, b.Seq(), []HashPair{hp})
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
func (bt *BlockTree) RemoveBlock(b Block) error {
	return blockdb.UpdateTx(func(tx *bolt.Tx) error {
		// delete block in blocks bucket.
		blocks := tx.Bucket(bt.blocks.Name)
		hash := b.HashHeader()
		if err := blocks.Delete(hash[:]); err != nil {
			return err
		}

		// get tree bucket.
		tree := tx.Bucket(bt.tree.Name)

		// check if this block has children
		has, err := hasChild(tree, b)
		if err != nil {
			return err
		}
		if has {
			return errHasChild
		}

		// get block hash pairs in depth
		hashPairs, err := getHashPairInDepth(tree, b.Seq(), func(hp HashPair) bool {
			return true
		})
		if err != nil {
			return err
		}

		// remove block hash pair in tree.
		ps := removePairs(hashPairs, HashPair{hash, b.PreHashHeader()})
		if len(ps) == 0 {
			tree.Delete(itob(b.Seq()))
			return nil
		}

		// update the hash pairs in tree.
		return setHashPairInDepth(tree, b.Seq(), ps)
	})
}

// GetBlock get block by hash, return nil on not found
func (bt *BlockTree) GetBlock(hash cipher.SHA256) *Block {
	return bt.getBlock(hash)
}

// GetBlockInDepth get block in depth, return nil on not found,
// the filter is used to choose the appropriate block.
func (bt *BlockTree) GetBlockInDepth(depth uint64, filter func(hps []HashPair) cipher.SHA256) *Block {
	hash, err := bt.getHashInDepth(depth, filter)
	if err != nil {
		logger.Error(err.Error())
		return nil
	}

	return bt.getBlock(hash)
}

func (bt *BlockTree) getBlock(hash cipher.SHA256) *Block {
	bin := bt.blocks.Get(hash[:])
	if bin == nil {
		return nil
	}
	block := Block{}
	if err := encoder.DeserializeRaw(bin, &block); err != nil {
		return nil
	}
	return &block
}

func (bt *BlockTree) getHashInDepth(depth uint64, filter func(ps []HashPair) cipher.SHA256) (cipher.SHA256, error) {
	key := itob(depth)
	pairsBin := bt.tree.Get(key)
	pairs := []HashPair{}
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

func containHash(hashPairs []HashPair, pair HashPair) bool {
	for _, p := range hashPairs {
		if p.Hash == pair.Hash {
			return true
		}
	}
	return false
}

func removePairs(HashPairs []HashPair, pair HashPair) []HashPair {
	pairs := []HashPair{}
	for _, p := range HashPairs {
		if p.Hash == pair.Hash && p.PreHash == pair.PreHash {
			continue
		}
		pairs = append(pairs, p)
	}
	return pairs
}

func getHashPairInDepth(tree *bolt.Bucket, dep uint64, fn func(hp HashPair) bool) ([]HashPair, error) {
	v := tree.Get(itob(dep))
	if v == nil {
		return []HashPair{}, nil
	}

	hps := []HashPair{}
	if err := encoder.DeserializeRaw(v, &hps); err != nil {
		return nil, err
	}
	pairs := []HashPair{}
	for _, ps := range hps {
		if fn(ps) {
			pairs = append(pairs, ps)
		}
	}
	return pairs, nil
}

func setBlock(bkt *bolt.Bucket, b *Block) error {
	bin := encoder.Serialize(b)
	key := b.HashHeader()
	return bkt.Put(key[:], bin)
}

// check if this block has children
func hasChild(bkt *bolt.Bucket, b Block) (bool, error) {
	// get the child block hash pair, whose pre hash point to current block.
	childHashPair, err := getHashPairInDepth(bkt, b.Head.BkSeq+1, func(hp HashPair) bool {
		return hp.PreHash == b.HashHeader()
	})

	if err != nil {
		return false, nil
	}

	return len(childHashPair) > 0, nil
}

func setHashPairInDepth(bkt *bolt.Bucket, dep uint64, hps []HashPair) error {
	hpsBin := encoder.Serialize(hps)
	key := itob(dep)
	return bkt.Put(key, hpsBin)
}

func allPairs(hp HashPair) bool {
	return true
}
