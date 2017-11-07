package blockdb

import (
	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// BlockSigs manages known blockSigs as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// blockSigs per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type BlockSigs struct {
}

const (
	blockSigsBkt = "block_sigs"
)

// NewBlockSigs create block signature buckets
func NewBlockSigs(db *bolt.DB) (*BlockSigs, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(blockSigsBkt))
		return err
	}); err != nil {
		return nil, err
	}

	return &BlockSigs{}, nil
}

func (bs BlockSigs) sigs(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(blockSigsBkt))
}

// Get returns signature of specific block
func (bs BlockSigs) Get(tx *bolt.Tx, hash cipher.SHA256) (cipher.Sig, bool, error) {
	bin := bs.sigs(tx).Get(hash[:])
	if bin == nil {
		return cipher.Sig{}, false, nil
	}

	var sig cipher.Sig
	if err := encoder.DeserializeRaw(bin, &sig); err != nil {
		return cipher.Sig{}, false, err
	}

	return sig, true, nil
}

// Add add signed block with bolt.Tx
func (bs *BlockSigs) Add(tx *bolt.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	return bs.sigs(tx).Put(hash[:], encoder.Serialize(sig))
}
