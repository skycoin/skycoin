package blockdb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BlockCosigsBkt holds the cosignatures of blocks of a multi-signer chain,
	// keyed by header hash.
	BlockCosigsBkt = []byte("block_cosigs")
)

// blockCosigs stores block cosignatures apart from the blocks, so the block
// encoding and the nodes that read it are unchanged.
type blockCosigs struct{}

// Get returns the cosignatures of a block, nil if it has none. A database from
// before the multi-signer code, opened read-only, has no bucket yet.
func (bs *blockCosigs) Get(tx *dbutil.Tx, hash cipher.SHA256) ([]coin.BlockCosig, error) {
	if !dbutil.Exists(tx, BlockCosigsBkt) {
		return nil, nil
	}
	v, err := dbutil.GetBucketValue(tx, BlockCosigsBkt, hash[:])
	if err != nil || v == nil {
		return nil, err
	}
	cs, _, err := coin.DecodeBlockCosigs(v)
	return cs, err
}

// Add stores the cosignatures of a block.
func (bs *blockCosigs) Add(tx *dbutil.Tx, hash cipher.SHA256, cs []coin.BlockCosig) error {
	return dbutil.PutBucketValue(tx, BlockCosigsBkt, hash[:], coin.EncodeBlockCosigs(cs))
}

// GetBlockCosigs returns the cosignatures of the block with the given header hash.
func (bc *Blockchain) GetBlockCosigs(tx *dbutil.Tx, hash cipher.SHA256) ([]coin.BlockCosig, error) {
	return bc.cosigs.Get(tx, hash)
}
