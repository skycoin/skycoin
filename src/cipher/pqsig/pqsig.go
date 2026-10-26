// Package pqsig signs and verifies with ML-DSA-44 (FIPS 204), the post-quantum
// signature the block publisher adds beside its secp256k1 signature.
package pqsig

import (
	"crypto/mldsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
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

// The context strings keep a block signature from passing as an input signature
// and the other way round (FIPS 204 section 5.2).
var (
	blockContext = &mldsa.Options{Context: "skycoin block"}
	inputContext = &mldsa.Options{Context: "skycoin txn input"}
)

var (
	// ErrInvalidSig is returned for a signature that does not verify.
	ErrInvalidSig = errors.New("invalid ML-DSA-44 signature")
	// ErrMissingSig is returned when a block that needs a signature has none.
	ErrMissingSig = errors.New("missing ML-DSA-44 signature")
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
	return sign(sk, headerHash, blockContext)
}

// VerifyBlock checks sig over a block header hash against pk.
func VerifyBlock(pk PubKey, headerHash, sig []byte) error {
	return verify(pk, headerHash, sig, blockContext)
}

// SignInput signs the hash a transaction input commits to. The signature is deterministic.
func SignInput(sk SecKey, inputHash []byte) ([]byte, error) {
	return sign(sk, inputHash, inputContext)
}

// VerifyInput checks sig over a transaction input hash against pk.
func VerifyInput(pk PubKey, inputHash, sig []byte) error {
	return verify(pk, inputHash, sig, inputContext)
}

func sign(sk SecKey, msg []byte, ctx *mldsa.Options) ([]byte, error) {
	priv, err := mldsa.NewPrivateKey(mldsa.MLDSA44(), sk)
	if err != nil {
		return nil, err
	}
	return priv.SignDeterministic(msg, ctx)
}

func verify(pk PubKey, msg, sig []byte, ctx *mldsa.Options) error {
	if len(sig) == 0 {
		return ErrMissingSig
	}
	pub, err := mldsa.NewPublicKey(mldsa.MLDSA44(), pk)
	if err != nil {
		return err
	}
	if err := mldsa.Verify(pub, msg, sig, ctx); err != nil {
		return ErrInvalidSig
	}
	return nil
}

// DeriveSecKey derives the secret key at index from a wallet seed, as SHA-256 over
// a domain tag, the seed and the big-endian index. One seed backs every key.
func DeriveSecKey(seed []byte, index uint32) SecKey {
	h := sha256.New()
	h.Write([]byte("skycoin pq key")) //nolint:errcheck // hash writes cannot fail
	h.Write(seed)                     //nolint:errcheck // hash writes cannot fail
	var i [4]byte
	binary.BigEndian.PutUint32(i[:], index)
	h.Write(i[:]) //nolint:errcheck // hash writes cannot fail
	return h.Sum(nil)
}
