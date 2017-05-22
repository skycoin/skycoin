package blockdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// BlockSigs manages known BlockSigs as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// BlockSigs per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type BlockSigs struct {
	Sigs *bucket.Bucket
}

// NewBlockSigs create block signature buckets
func NewBlockSigs(db *bolt.DB) *BlockSigs {
	sigs, err := bucket.New([]byte("block_sigs"), db)
	if err != nil {
		panic(err)
	}

	return &BlockSigs{
		Sigs: sigs,
	}
}

// Verify Checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
// func (self *BlockSigs) Verify(masterPublic cipher.PubKey, bc *visor.Blockchain) error {
// 	for i := uint64(0); i <= bc.Head().Seq(); i++ {
// 		b := bc.GetBlockInDepth(i)
// 		if b == nil {
// 			return fmt.Errorf("no block in depth %v", i)
// 		}
// 		// get sig
// 		sig, err := self.Get(b.HashHeader())
// 		if err != nil {
// 			return err
// 		}

// 		if err := cipher.VerifySignature(masterPublic, sig, bc.GetBlockInDepth(i).HashHeader()); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// Get returns signature of specific block
func (bs BlockSigs) Get(hash cipher.SHA256) (cipher.Sig, error) {
	bin := bs.Sigs.Get(hash[:])
	if bin == nil {
		return cipher.Sig{}, fmt.Errorf("no sig for %v", hash.Hex())
	}
	var sig cipher.Sig
	if err := encoder.DeserializeRaw(bin, &sig); err != nil {
		return cipher.Sig{}, err
	}
	return sig, nil
}

// Add stores the signed block into db.
func (bs *BlockSigs) Add(sb *coin.SignedBlock) error {
	hash := sb.Block.HashHeader()
	return bs.Sigs.Put(hash[:], encoder.Serialize(sb.Sig))
}
