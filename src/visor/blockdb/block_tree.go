package blockdb

import (
	"errors"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

const (
	blocksBkt = "blocks"
	treeBkt   = "block_tree"
)

var (
	emptyHash      cipher.SHA256
	errBlockExist  = errors.New("block already exist")
	errNoParent    = errors.New("block is not genesis and have no parent")
	errWrongParent = errors.New("wrong parent")
	errHasChild    = errors.New("remove block failed, it has children")
)

type hasBucket interface {
	Bucket(name []byte) *bolt.Bucket
}

// blockTree use the blockdb store all blocks and maintains the block tree struct.
type blockTree struct {
	db *bolt.DB
}

// newBlockTree create buckets in blockdb if does not exist.
func newBlockTree(db *bolt.DB) (*blockTree, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(blocksBkt)); err != nil {
			return err
		}

		_, err := tx.CreateBucketIfNotExists([]byte(treeBkt))
		return err
	}); err != nil {
		return nil, err
	}

	return &blockTree{
		db: db,
	}, nil
}

func (bt *blockTree) blocks(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(blocksBkt))
}

func (bt *blockTree) tree(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(treeBkt))
}

// AddBlock adds block with *bolt.Tx
func (bt *blockTree) AddBlock(tx *bolt.Tx, b *coin.Block) error {
	bkt := bt.blocks(tx)

	// can't store block if it's not genesis block and has no parent.
	if b.Seq() > 0 && b.PreHashHeader() == emptyHash {
		return errNoParent
	}

	// check if the block already exist.
	hash := b.HashHeader()
	if blk := bkt.Get(hash[:]); blk != nil {
		return errBlockExist
	}

	// write block into blocks bucket.
	if err := setBlock(bkt, b); err != nil {
		return err
	}

	// get tree bucket.
	tree := bt.tree(tx)

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
}

// RemoveBlock remove block from blocks bucket and tree bucket.
// can't remove block if it has children.
func (bt *blockTree) RemoveBlock(b *coin.Block) error {
	return bt.db.Update(func(tx *bolt.Tx) error {
		// delete block in blocks bucket.
		blocks := bt.blocks(tx)
		hash := b.HashHeader()
		if err := blocks.Delete(hash[:]); err != nil {
			return err
		}

		// get tree bucket.
		tree := bt.tree(tx)

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
			tree.Delete(bucket.Itob(b.Seq()))
			return nil
		}

		// update the hash pairs in tree.
		return setHashPairInDepth(tree, b.Seq(), ps)
	})
}

// GetBlockInDepth get block in depth, return nil on not found,
// the filter is used to choose the appropriate block.
func (bt *blockTree) GetBlockInDepth(tx *bolt.Tx, depth uint64, filter func(hps []coin.HashPair) cipher.SHA256) (*coin.Block, error) {
	hash, err := bt.getHashInDepth(tx, depth, filter)
	if err != nil {
		return nil, err
	}

	return bt.GetBlock(tx, hash)
}

// GetBlock get block by hash, return nil on not found
func (bt *blockTree) GetBlock(tx *bolt.Tx, hash cipher.SHA256) (*coin.Block, error) {
	bin := bt.blocks(tx).Get(hash[:])
	if bin == nil {
		return nil, nil
	}

	block := coin.Block{}
	if err := encoder.DeserializeRaw(bin, &block); err != nil {
		return nil, err
	}

	return &block, nil
}

func (bt *blockTree) getHashInDepth(tx *bolt.Tx, depth uint64, filter func(ps []coin.HashPair) cipher.SHA256) (cipher.SHA256, error) {
	key := bucket.Itob(depth)

	pairsBin := bt.tree(tx).Get(key)
	pairs := []coin.HashPair{}
	if err := encoder.DeserializeRaw(pairsBin, &pairs); err != nil {
		return cipher.SHA256{}, err
	}

	hash := filter(pairs)
	return hash, nil
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
	v := tree.Get(bucket.Itob(dep))
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
		return false, err
	}

	return len(childHashPair) > 0, nil
}

func setHashPairInDepth(bkt *bolt.Bucket, dep uint64, hps []coin.HashPair) error {
	hpsBin := encoder.Serialize(hps)
	key := bucket.Itob(dep)
	return bkt.Put(key, hpsBin)
}

func allPairs(hp coin.HashPair) bool {
	return true
}
