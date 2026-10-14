package blockdb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BlockPQSigsBkt holds the ML-DSA-44 signatures of blocks, keyed by header hash.
	// Only blocks from the post-quantum fork height have one.
	BlockPQSigsBkt = []byte("block_pq_sigs")
)

// blockPQSigs stores post-quantum block signatures apart from the secp256k1 ones,
// so the block_sigs encoding and the nodes that read it are unchanged.
type blockPQSigs struct{}

// Get returns the post-quantum signature of a block, if it has one.
func (bs *blockPQSigs) Get(tx *dbutil.Tx, hash cipher.SHA256) ([]byte, error) {
	return dbutil.GetBucketValue(tx, BlockPQSigsBkt, hash[:])
}

// Add stores the post-quantum signature of a block.
func (bs *blockPQSigs) Add(tx *dbutil.Tx, hash cipher.SHA256, sig []byte) error {
	return dbutil.PutBucketValue(tx, BlockPQSigsBkt, hash[:], sig)
}
