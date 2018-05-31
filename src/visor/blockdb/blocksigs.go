package blockdb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BlockSigsBkt holds block signatures
	BlockSigsBkt = []byte("block_sigs")
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

// Get returns the signature of a specific block
func (bs blockSigs) Get(tx *dbutil.Tx, hash cipher.SHA256) (cipher.Sig, bool, error) {
	var sig cipher.Sig

	if ok, err := dbutil.GetBucketObjectDecoded(tx, BlockSigsBkt, hash[:], &sig); err != nil {
		return cipher.Sig{}, false, err
	} else if !ok {
		return cipher.Sig{}, false, nil
	}

	return sig, true, nil
}

// Add adds a signed block to the db
func (bs *blockSigs) Add(tx *dbutil.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	return dbutil.PutBucketValue(tx, BlockSigsBkt, hash[:], encoder.Serialize(sig))
}

// ForEach iterates all signatures and calls f on them
func (bs *blockSigs) ForEach(tx *dbutil.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
	return dbutil.ForEach(tx, BlockSigsBkt, func(k, v []byte) error {
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
