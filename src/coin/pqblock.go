package coin

import (
	"github.com/skycoin/skycoin/src/cipher/pqsig"
)

// PQBlockRule says which blocks must also carry an ML-DSA-44 signature by the
// post-quantum block publisher key, and checks it.
type PQBlockRule struct {
	// Pubkey is the post-quantum block publisher key. Empty turns the rule off.
	Pubkey pqsig.PubKey
	// ForkSeq is the first block that needs the signature. Genesis never does.
	ForkSeq uint64
}

// Applies reports whether the block at seq needs a post-quantum signature.
func (r PQBlockRule) Applies(seq uint64) bool {
	return len(r.Pubkey) != 0 && seq > 0 && seq >= r.ForkSeq
}

// Verify checks the post-quantum signature of b when the rule applies to it.
func (r PQBlockRule) Verify(b SignedBlock) error {
	if !r.Applies(b.Seq()) {
		return nil
	}
	return b.VerifyPQSignature(r.Pubkey)
}
