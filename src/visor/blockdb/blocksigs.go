package blockdb

import (
	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// blockSigs manages known blockSigs as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// blockSigs per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type blockSigs struct{}

var (
	blockSigsBkt = []byte("block_sigs")
)

// newBlockSigs create block signature bucket
func newBlockSigs(db *dbutil.DB) (*blockSigs, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(blockSigsBkt)
		return err
	}); err != nil {
		return nil, err
	}

	return &blockSigs{}, nil
}

// Get returns the signature of a specific block
func (bs blockSigs) Get(tx *bolt.Tx, hash cipher.SHA256) (cipher.Sig, bool, error) {
	var sig cipher.Sig

	if err := dbutil.GetBucketObjectDecoded(tx, blockSigsBkt, hash[:], &sig); err != nil {
		switch err.(type) {
		case dbutil.ObjectNotExistErr:
			return cipher.Sig{}, false, nil
		default:
			return cipher.Sig{}, false, err
		}
	}

	return sig, true, nil
}

// Add adds a signed block to the db
func (bs *blockSigs) Add(tx *bolt.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	return dbutil.PutBucketValue(tx, blockSigsBkt, hash[:], encoder.Serialize(sig))
}

// ForEach iterates all signatures and calls f on them
func (bs *blockSigs) ForEach(tx *bolt.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
	return dbutil.ForEach(tx, blocksBkt, func(k, v []byte) error {
		hash, err := cipher.SHA256FromBytes(k)
		if err != nil {
			return err
		}

		var sig cipher.Sig
		if err := encoder.DeserializeRaw(v, &sig); err != nil {
			return err
		}

		return f(hash, sig)
	})
}
