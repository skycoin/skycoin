package coin

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// MaxBlockSigners bounds the signer list of a multi-signer chain.
const MaxBlockSigners = 1024

var (
	// ErrBlockProposerNotSigner is returned when a block's own signature is not by a listed signer
	ErrBlockProposerNotSigner = errors.New("block is not signed by a listed block signer")
	// ErrBlockCosigsMalformed is returned when encoded cosignatures cannot be decoded
	ErrBlockCosigsMalformed = errors.New("malformed block cosignatures")
)

// ErrNotEnoughBlockSigners is returned when a block has fewer valid signatures than the threshold.
type ErrNotEnoughBlockSigners struct {
	Have, Need int
}

func (e ErrNotEnoughBlockSigners) Error() string {
	return fmt.Sprintf("block has %d valid signatures from listed signers, needs %d", e.Have, e.Need)
}

// BlockCosig is a signature over a block header hash by one of the block signers.
type BlockCosig struct {
	// Signer is the index of the signer in MultiSignerRule.Signers
	Signer uint32
	Sig    cipher.Sig
}

// MultiSignerRule says which blocks need signatures from k of N listed signers,
// and checks them. See docs/multi-signer-fiber.md.
type MultiSignerRule struct {
	// Signers are the block signer keys. Empty turns the rule off.
	Signers []cipher.PubKey
	// Threshold is k, the number of signatures that makes a block final.
	Threshold int
	// ForkSeq is the first block that needs k signatures. Genesis never does.
	ForkSeq uint64
}

// DefaultBlockSignerThreshold returns floor(2n/3)+1, the smallest k above 2n/3.
func DefaultBlockSignerThreshold(n int) int {
	return 2*n/3 + 1
}

// NewMultiSignerRule parses hex signer keys into a rule and checks it. A zero
// threshold takes the default. No signers gives the rule that is off.
func NewMultiSignerRule(signers []string, threshold int, forkSeq uint64) (MultiSignerRule, error) {
	if len(signers) == 0 {
		return MultiSignerRule{}, nil
	}
	r := MultiSignerRule{Signers: make([]cipher.PubKey, len(signers)), Threshold: threshold, ForkSeq: forkSeq}
	for i, s := range signers {
		pk, err := cipher.PubKeyFromHex(s)
		if err != nil {
			return MultiSignerRule{}, fmt.Errorf("invalid block signer %q: %v", s, err)
		}
		r.Signers[i] = pk
	}
	if r.Threshold == 0 {
		r.Threshold = DefaultBlockSignerThreshold(len(signers))
	}
	return r, r.Validate()
}

// Enabled reports whether the chain has a signer list
func (r MultiSignerRule) Enabled() bool {
	return len(r.Signers) != 0
}

// Applies reports whether the block at seq needs k signatures.
func (r MultiSignerRule) Applies(seq uint64) bool {
	return r.Enabled() && seq > 0 && seq >= r.ForkSeq
}

// Validate checks the signer list and the threshold.
func (r MultiSignerRule) Validate() error {
	if !r.Enabled() {
		return nil
	}
	n := len(r.Signers)
	if n > MaxBlockSigners {
		return fmt.Errorf("too many block signers: %d, at most %d", n, MaxBlockSigners)
	}
	seen := make(map[cipher.PubKey]struct{}, n)
	for i, pk := range r.Signers {
		if err := pk.Verify(); err != nil {
			return fmt.Errorf("block signer %d is not a valid public key: %v", i, err)
		}
		if _, ok := seen[pk]; ok {
			return fmt.Errorf("block signer %s is listed twice", pk.Hex())
		}
		seen[pk] = struct{}{}
	}
	if r.Threshold > n || 3*r.Threshold <= 2*n {
		return fmt.Errorf("block signer threshold %d must be more than 2/3 of %d signers and at most %d", r.Threshold, n, n)
	}
	return nil
}

// Proposer returns the index of the signer that proposes the block at seq in round.
func (r MultiSignerRule) Proposer(seq, round uint64) int {
	return int((seq + round) % uint64(len(r.Signers))) //nolint:gosec // bounded by MaxBlockSigners
}

// SignerIndex returns the index of pk in the signer list, or -1.
func (r MultiSignerRule) SignerIndex(pk cipher.PubKey) int {
	for i, s := range r.Signers {
		if s == pk {
			return i
		}
	}
	return -1
}

// ProposerIndex returns the index of the signer whose key made b.Sig, or -1.
func (r MultiSignerRule) ProposerIndex(b SignedBlock) int {
	h := b.HashHeader()
	pk, err := cipher.PubKeyFromSig(b.Sig, h)
	if err != nil {
		return -1
	}
	i := r.SignerIndex(pk)
	if i < 0 || cipher.VerifyPubKeySignedHash(pk, b.Sig, h) != nil {
		return -1
	}
	return i
}

// VerifyCosig checks that c is a valid signature by its signer over hash.
func (r MultiSignerRule) VerifyCosig(hash cipher.SHA256, c BlockCosig) error {
	if int(c.Signer) >= len(r.Signers) {
		return fmt.Errorf("cosignature signer index %d is not in the signer list", c.Signer)
	}
	return cipher.VerifyPubKeySignedHash(r.Signers[c.Signer], c.Sig, hash)
}

// ValidCosigs returns the cosignatures of b that are valid and by distinct listed
// signers other than the proposer, with the proposer's index. Others are dropped.
func (r MultiSignerRule) ValidCosigs(b SignedBlock) (int, []BlockCosig, error) {
	p := r.ProposerIndex(b)
	if p < 0 {
		return -1, nil, ErrBlockProposerNotSigner
	}
	h := b.HashHeader()
	signed := map[uint32]struct{}{uint32(p): {}} //nolint:gosec // p < MaxBlockSigners
	var valid []BlockCosig
	for _, c := range b.Cosigs {
		if _, ok := signed[c.Signer]; ok {
			continue
		}
		if r.VerifyCosig(h, c) != nil {
			continue
		}
		signed[c.Signer] = struct{}{}
		valid = append(valid, c)
	}
	return p, valid, nil
}

// Verify checks that b is signed by a listed signer and has at least k valid
// signatures from distinct listed signers, counting the proposer's. It returns
// the valid cosignatures. A rule that does not apply to b accepts it unchanged.
func (r MultiSignerRule) Verify(b SignedBlock) ([]BlockCosig, error) {
	if !r.Applies(b.Seq()) {
		return b.Cosigs, nil
	}
	_, valid, err := r.ValidCosigs(b)
	if err != nil {
		return nil, err
	}
	if have := len(valid) + 1; have < r.Threshold {
		return nil, ErrNotEnoughBlockSigners{Have: have, Need: r.Threshold}
	}
	return valid, nil
}

// BlockSigners returns the keys that validly signed b, the proposer first.
func (r MultiSignerRule) BlockSigners(b SignedBlock) []cipher.PubKey {
	if !r.Applies(b.Seq()) {
		return nil
	}
	p, valid, err := r.ValidCosigs(b)
	if err != nil {
		return nil
	}
	pks := make([]cipher.PubKey, 0, len(valid)+1)
	pks = append(pks, r.Signers[p])
	for _, c := range valid {
		pks = append(pks, r.Signers[c.Signer])
	}
	return pks
}

// SignBlockCosig signs a block header hash as signer index i.
func SignBlockCosig(hash cipher.SHA256, i int, sec cipher.SecKey) (BlockCosig, error) {
	sig, err := cipher.SignHash(hash, sec)
	if err != nil {
		return BlockCosig{}, err
	}
	return BlockCosig{Signer: uint32(i), Sig: sig}, nil //nolint:gosec // i < MaxBlockSigners
}

const blockCosigSize = 4 + len(cipher.Sig{})

// EncodeBlockCosigs encodes cosignatures for storage and transport.
func EncodeBlockCosigs(cs []BlockCosig) []byte {
	buf := make([]byte, 4+len(cs)*blockCosigSize)
	e := &encoder.Encoder{Buffer: buf}
	e.Uint32(uint32(len(cs))) //nolint:gosec // bounded by MaxBlockSigners when decoded
	for _, c := range cs {
		e.Uint32(c.Signer)
		e.CopyBytes(c.Sig[:])
	}
	return buf
}

// DecodeBlockCosigs decodes cosignatures. It returns how many bytes it read.
func DecodeBlockCosigs(b []byte) ([]BlockCosig, int, error) {
	d := &encoder.Decoder{Buffer: b}
	n, err := d.Uint32()
	if err != nil || n > MaxBlockSigners || int(n)*blockCosigSize > len(d.Buffer) {
		return nil, 0, ErrBlockCosigsMalformed
	}
	cs := make([]BlockCosig, n)
	for i := range cs {
		if cs[i].Signer, err = d.Uint32(); err != nil {
			return nil, 0, ErrBlockCosigsMalformed
		}
		copy(cs[i].Sig[:], d.Buffer[:len(cipher.Sig{})])
		d.Buffer = d.Buffer[len(cipher.Sig{}):]
	}
	return cs, len(b) - len(d.Buffer), nil
}
