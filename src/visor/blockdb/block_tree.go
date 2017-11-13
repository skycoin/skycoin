package blockdb

import (
	"errors"
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	emptyHash      cipher.SHA256
	errBlockExist  = errors.New("block already exists")
	errNoParent    = errors.New("block is not genesis and has no parent")
	errWrongParent = errors.New("wrong parent")
	errHasChild    = errors.New("remove block failed, it has children")

	blocksBkt = []byte("blocks")
	treeBkt   = []byte("block_tree")
)

// Walker function for go through blockchain
type Walker func(*bolt.Tx, []coin.HashPair) (cipher.SHA256, bool)

// blockTree use the blockdb store all blocks and maintains the block tree struct.
type blockTree struct {
}

// newBlockTree create buckets in blockdb if does not exist.
func newBlockTree(db *dbutil.DB) (*blockTree, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			blocksBkt,
			treeBkt,
		})
	}); err != nil {
		return nil, err
	}

	return &blockTree{}, nil
}

// AddBlock adds block with *bolt.Tx
func (bt *blockTree) AddBlock(tx *bolt.Tx, b *coin.Block) error {
	// can't store block if it's not genesis block and has no parent.
	if b.Seq() > 0 && b.PreHashHeader() == emptyHash {
		return errNoParent
	}

	// check if the block already exist.
	hash := b.HashHeader()
	if ok, err := dbutil.BucketHasKey(tx, blocksBkt, hash[:]); err != nil {
		return err
	} else if ok {
		return errBlockExist
	}

	// write block into blocks bucket.
	if err := dbutil.PutBucketValue(tx, blocksBkt, hash[:], encoder.Serialize(b)); err != nil {
		return err
	}

	// the pre hash must be in depth - 1.
	if b.Seq() > 0 {
		preHash := b.PreHashHeader()
		parentHashPair, err := getHashPairInDepth(tx, b.Seq()-1, func(hp coin.HashPair) bool {
			return hp.Hash == preHash
		})
		if err != nil {
			return err
		}
		if len(parentHashPair) == 0 {
			return errWrongParent
		}
	}

	hp := coin.HashPair{
		Hash:    hash,
		PreHash: b.Head.PrevHash,
	}

	// get block pairs in the depth
	hashPairs, err := getHashPairInDepth(tx, b.Seq(), allPairs)
	if err != nil {
		return err
	}

	if len(hashPairs) == 0 {
		// no hash pair exist in the depth.
		// write the hash pair into tree.
		return setHashPairInDepth(tx, b.Seq(), []coin.HashPair{hp})
	}

	// check dup block
	if containHash(hashPairs, hp) {
		return errBlockExist
	}

	hashPairs = append(hashPairs, hp)
	return setHashPairInDepth(tx, b.Seq(), hashPairs)
}

// RemoveBlock remove block from blocks bucket and tree bucket.
// can't remove block if it has children.
func (bt *blockTree) RemoveBlock(tx *bolt.Tx, b *coin.Block) error {
	// delete block in blocks bucket.
	hash := b.HashHeader()
	if err := dbutil.Delete(tx, blocksBkt, hash[:]); err != nil {
		return err
	}

	// check if this block has children
	if has, err := hasChild(tx, *b); err != nil {
		return err
	} else if has {
		return errHasChild
	}

	// get block hash pairs in depth
	hashPairs, err := getHashPairInDepth(tx, b.Seq(), allPairs)
	if err != nil {
		return err
	}

	// remove block hash pair in tree.
	ps := removePairs(hashPairs, coin.HashPair{
		Hash:    hash,
		PreHash: b.PreHashHeader(),
	})

	if len(ps) == 0 {
		return dbutil.Delete(tx, treeBkt, dbutil.Itob(b.Seq()))
	}

	// update the hash pairs in tree.
	return setHashPairInDepth(tx, b.Seq(), ps)
}

// GetBlock get block by hash, return nil on not found
func (bt *blockTree) GetBlock(tx *bolt.Tx, hash cipher.SHA256) (*coin.Block, error) {
	var b coin.Block

	if ok, err := dbutil.GetBucketObjectDecoded(tx, blocksBkt, hash[:], &b); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	if hash != b.HashHeader() {
		return nil, fmt.Errorf("DB key %s does not match block hash header %s", hash, b.HashHeader())
	}

	return &b, nil
}

// GetBlockInDepth get block in depth, return nil on not found,
// the filter is used to choose the appropriate block.
func (bt *blockTree) GetBlockInDepth(tx *bolt.Tx, depth uint64, filter Walker) (*coin.Block, error) {
	hash, ok, err := bt.getHashInDepth(tx, depth, filter)
	if err != nil {
		return nil, fmt.Errorf("BlockTree.getHashInDepth failed: %v", err)
	} else if !ok {
		return nil, nil
	}

	return bt.GetBlock(tx, hash)
}

// ForEachBlock iterates all blocks and calls f on them
func (bt *blockTree) ForEachBlock(tx *bolt.Tx, f func(b *coin.Block) error) error {
	return dbutil.ForEach(tx, blocksBkt, func(_, v []byte) error {
		var b coin.Block
		if err := encoder.DeserializeRaw(v, &b); err != nil {
			return err
		}

		return f(&b)
	})
}

func (bt *blockTree) getHashInDepth(tx *bolt.Tx, depth uint64, filter Walker) (cipher.SHA256, bool, error) {
	var pairs []coin.HashPair
	if ok, err := dbutil.GetBucketObjectDecoded(tx, treeBkt, dbutil.Itob(depth), &pairs); err != nil {
		return cipher.SHA256{}, false, err
	} else if !ok {
		return cipher.SHA256{}, false, nil
	}

	hash, ok := filter(tx, pairs)
	if !ok {
		return cipher.SHA256{}, false, errors.New("No hash found in depth")
	}

	return hash, true, nil
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

func getHashPairInDepth(tx *bolt.Tx, dep uint64, fn func(hp coin.HashPair) bool) ([]coin.HashPair, error) {
	var hps []coin.HashPair
	if ok, err := dbutil.GetBucketObjectDecoded(tx, treeBkt, dbutil.Itob(dep), &hps); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	var pairs []coin.HashPair
	for _, ps := range hps {
		if fn(ps) {
			pairs = append(pairs, ps)
		}
	}
	return pairs, nil
}

// check if this block has children
func hasChild(tx *bolt.Tx, b coin.Block) (bool, error) {
	// get the child block hash pair, whose pre hash point to current block.
	childHashPair, err := getHashPairInDepth(tx, b.Head.BkSeq+1, func(hp coin.HashPair) bool {
		return hp.PreHash == b.HashHeader()
	})

	if err != nil {
		return false, err
	}

	return len(childHashPair) > 0, nil
}

func setHashPairInDepth(tx *bolt.Tx, dep uint64, hps []coin.HashPair) error {
	return dbutil.PutBucketValue(tx, treeBkt, dbutil.Itob(dep), encoder.Serialize(hps))
}

func allPairs(hp coin.HashPair) bool {
	return true
}
