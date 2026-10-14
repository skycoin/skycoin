// Package pqsig signs and verifies with ML-DSA-44 (FIPS 204), the post-quantum
// signature the block publisher adds beside its secp256k1 signature.
package pqsig

import (
	"crypto/mldsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	// SeedSize is the size of a secret key, which is the ML-DSA seed.
	SeedSize = mldsa.PrivateKeySize
	// PubKeySize is the size of an ML-DSA-44 public key.
	PubKeySize = mldsa.MLDSA44PublicKeySize
	// SigSize is the size of an ML-DSA-44 signature.
	SigSize = mldsa.MLDSA44SignatureSize
)

// blockContext separates block signatures from any other use of the same key
// (FIPS 204 section 5.2, the context string).
var blockContext = &mldsa.Options{Context: "skycoin block"}

var (
	// ErrInvalidSig is returned for a signature that does not verify.
	ErrInvalidSig = errors.New("invalid ML-DSA-44 block signature")
	// ErrMissingSig is returned when a block that needs a signature has none.
	ErrMissingSig = errors.New("missing ML-DSA-44 block signature")
)

// PubKey is an ML-DSA-44 public key.
type PubKey []byte

// SecKey is an ML-DSA-44 secret key, in its 32-byte seed form.
type SecKey []byte

// NewSecKey makes a random secret key.
func NewSecKey() (SecKey, error) {
	seed := make([]byte, SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}
	return seed, nil
}

// PubKeyFromHex decodes and checks a hex public key.
func PubKeyFromHex(s string) (PubKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	pk := PubKey(b)
	if err := pk.Verify(); err != nil {
		return nil, err
	}
	return pk, nil
}

// SecKeyFromHex decodes and checks a hex secret key.
func SecKeyFromHex(s string) (SecKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(b) != SeedSize {
		return nil, fmt.Errorf("ML-DSA-44 secret key must be %d bytes, got %d", SeedSize, len(b))
	}
	return b, nil
}

// Verify checks that pk is a well-formed ML-DSA-44 public key.
func (pk PubKey) Verify() error {
	_, err := mldsa.NewPublicKey(mldsa.MLDSA44(), pk)
	return err
}

// Hex returns the key in hex.
func (pk PubKey) Hex() string { return hex.EncodeToString(pk) }

// PubKey returns the public key of sk.
func (sk SecKey) PubKey() (PubKey, error) {
	priv, err := mldsa.NewPrivateKey(mldsa.MLDSA44(), sk)
	if err != nil {
		return nil, err
	}
	return priv.PublicKey().Bytes(), nil
}

// Hex returns the key in hex.
func (sk SecKey) Hex() string { return hex.EncodeToString(sk) }

// SignBlock signs a block header hash. The signature is deterministic.
func SignBlock(sk SecKey, headerHash []byte) ([]byte, error) {
	priv, err := mldsa.NewPrivateKey(mldsa.MLDSA44(), sk)
	if err != nil {
		return nil, err
	}
	return priv.SignDeterministic(headerHash, blockContext)
}

// VerifyBlock checks sig over a block header hash against pk.
func VerifyBlock(pk PubKey, headerHash, sig []byte) error {
	if len(sig) == 0 {
		return ErrMissingSig
	}
	pub, err := mldsa.NewPublicKey(mldsa.MLDSA44(), pk)
	if err != nil {
		return err
	}
	if err := mldsa.Verify(pub, headerHash, sig, blockContext); err != nil {
		return ErrInvalidSig
	}
	return nil
}
