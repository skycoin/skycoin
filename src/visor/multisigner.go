package visor

import (
	"errors"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// SignerVotesBkt records, by height, the one block this node signed as a block
// signer of a multi-signer chain.
var SignerVotesBkt = []byte("block_signer_votes")

// ErrAlreadyVoted is returned when a signer is asked to sign a second block at a height.
type ErrAlreadyVoted struct {
	Seq  uint64
	Hash cipher.SHA256
}

func (e ErrAlreadyVoted) Error() string {
	return fmt.Sprintf("already signed block %s at height %d", e.Hash.Hex(), e.Seq)
}

// verifyBlockSignatures checks the signatures a block needs: k of the listed
// signers where the multi-signer rule applies, otherwise the publisher key and
// the post-quantum rule.
func verifyBlockSignatures(b coin.SignedBlock, pubkey cipher.PubKey, pq coin.PQBlockRule, ms coin.MultiSignerRule) error {
	if ms.Applies(b.Seq()) {
		_, err := ms.Verify(b)
		return err
	}
	if err := b.VerifySignature(pubkey); err != nil {
		return err
	}
	return pq.Verify(b)
}

// SignerIndex returns this node's index in the block signer list, or -1 if it
// is not a signer of a multi-signer chain.
func (vs *Visor) SignerIndex() int {
	if !vs.Config.IsBlockPublisher || !vs.Config.MultiSigner.Enabled() {
		return -1
	}
	pk, err := cipher.PubKeyFromSecKey(vs.Config.BlockchainSeckey)
	if err != nil {
		return -1
	}
	return vs.Config.MultiSigner.SignerIndex(pk)
}

// BlockVote returns the block this signer signed at seq, if any.
func (vs *Visor) BlockVote(seq uint64) (cipher.SHA256, bool, error) {
	var h cipher.SHA256
	var ok bool
	err := vs.db.View("BlockVote", func(tx *dbutil.Tx) error {
		v, err := dbutil.GetBucketValue(tx, SignerVotesBkt, dbutil.Itob(seq))
		if err != nil || v == nil {
			return err
		}
		if len(v) != len(h) {
			return errors.New("malformed block signer vote")
		}
		copy(h[:], v)
		ok = true
		return nil
	})
	return h, ok, err
}

// RecordBlockVote records that this signer signs the block hash at seq. It is
// written before the signature leaves the node, and fails with ErrAlreadyVoted
// if another block was signed at seq (docs/multi-signer-fiber.md, Safety).
func (vs *Visor) RecordBlockVote(seq uint64, hash cipher.SHA256) error {
	return vs.db.Update("RecordBlockVote", func(tx *dbutil.Tx) error {
		v, err := dbutil.GetBucketValue(tx, SignerVotesBkt, dbutil.Itob(seq))
		if err != nil {
			return err
		}
		if v != nil {
			var prev cipher.SHA256
			copy(prev[:], v)
			if prev != hash {
				return ErrAlreadyVoted{Seq: seq, Hash: prev}
			}
			return nil
		}
		return dbutil.PutBucketValue(tx, SignerVotesBkt, dbutil.Itob(seq), hash[:])
	})
}

// CreateBlockProposal builds and signs a block from the unconfirmed pool without
// executing it. Only a block signer may call it.
func (vs *Visor) CreateBlockProposal(when time.Time) (coin.SignedBlock, error) {
	if vs.SignerIndex() < 0 {
		return coin.SignedBlock{}, errors.New("only a block signer can propose blocks")
	}
	var sb coin.SignedBlock
	err := vs.db.View("CreateBlockProposal", func(tx *dbutil.Tx) error {
		var err error
		sb, err = vs.createBlock(tx, uint64(when.UTC().Unix())) //nolint:gosec
		return err
	})
	return sb, err
}

// SignBlockCosig signs the header of b as this block signer.
func (vs *Visor) SignBlockCosig(b coin.Block) (coin.BlockCosig, error) {
	i := vs.SignerIndex()
	if i < 0 {
		return coin.BlockCosig{}, errors.New("only a block signer can cosign blocks")
	}
	return coin.SignBlockCosig(b.HashHeader(), i, vs.Config.BlockchainSeckey)
}
