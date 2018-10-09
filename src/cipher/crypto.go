/*
Package cipher implements cryptographic methods.

These methods include:

* Public and private key generation
* Address generation
* Signing

Private keys are secp256k1 keys. Addresses are base58 encoded.

All dependencies are either from the go stdlib, or are manually vendored
below this package. This manual vendoring ensures that the exact same dependencies
are used by any user of this package, regardless of their gopath.
*/
package cipher

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

var (
	// DebugLevel1 debug level one
	DebugLevel1 = true //checks for extremely unlikely conditions (10e-40)
	// DebugLevel2 debug level two
	DebugLevel2 = true //enable checks for impossible conditions
)

// PubKey public key
type PubKey [33]byte

// RandByte returns rand N bytes
func RandByte(n int) []byte {
	return secp256k1.RandByte(n)
}

// NewPubKey converts []byte to a PubKey
func NewPubKey(b []byte) (PubKey, error) {
	p := PubKey{}
	if len(b) != len(p) {
		return PubKey{}, errors.New("Invalid public key length")
	}
	copy(p[:], b[:])

	if err := p.Verify(); err != nil {
		return PubKey{}, err
	}

	return p, nil
}

// MustNewPubKey converts []byte to a PubKey, panics on error
func MustNewPubKey(b []byte) PubKey {
	p, err := NewPubKey(b)
	if err != nil {
		log.Panic(err)
	}
	return p
}

// PubKeyFromHex generates PubKey from hex string
func PubKeyFromHex(s string) (PubKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PubKey{}, errors.New("Invalid public key")
	}
	return NewPubKey(b)
}

// MustPubKeyFromHex decodes a hex encoded PubKey, panics on error
func MustPubKeyFromHex(s string) PubKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return MustNewPubKey(b)
}

// PubKeyFromSecKey recovers the public key for a secret key
func PubKeyFromSecKey(seckey SecKey) (PubKey, error) {
	if seckey == (SecKey{}) {
		return PubKey{}, errors.New("Cannot convert null SecKey to PubKey")
	}

	b := secp256k1.PubkeyFromSeckey(seckey[:])
	if b == nil {
		return PubKey{}, errors.New("PubKey recovery from SecKey failed. The recovery function assumes SecKey is valid, check SecKey")
	}

	return NewPubKey(b)
}

// MustPubKeyFromSecKey recovers the public key for a secret key. Panics on error.
func MustPubKeyFromSecKey(seckey SecKey) PubKey {
	pk, err := PubKeyFromSecKey(seckey)
	if err != nil {
		log.Panic(err)
	}
	return pk
}

// PubKeyFromSig recovers the public key from a signed hash
func PubKeyFromSig(sig Sig, hash SHA256) (PubKey, error) {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return PubKey{}, errors.New("Invalig sig: PubKey recovery failed")
	}
	return NewPubKey(rawPubKey)
}

// MustPubKeyFromSig recovers the public key from a signed hash, panics on error
func MustPubKeyFromSig(sig Sig, hash SHA256) PubKey {
	pk, err := PubKeyFromSig(sig, hash)
	if err != nil {
		log.Panic(err)
	}
	return pk
}

// Verify attempts to determine if pubkey is valid. Returns nil on success
func (pk PubKey) Verify() error {
	if secp256k1.VerifyPubkey(pk[:]) != 1 {
		return errors.New("Invalid public key")
	}
	return nil
}

// Hex returns a hex encoded PubKey string
func (pk PubKey) Hex() string {
	return hex.EncodeToString(pk[:])
}

// SecKey secret key
type SecKey [32]byte

// NewSecKey converts []byte to a SecKey
func NewSecKey(b []byte) (SecKey, error) {
	p := SecKey{}
	if len(b) != len(p) {
		return SecKey{}, errors.New("Invalid secret key length")
	}
	copy(p[:], b[:])

	// Disable the DebugLevel2 check here because it is too slow.
	// If desired, perform the full Verify() check after using this method
	if err := p.verify(false); err != nil {
		return SecKey{}, err
	}

	return p, nil
}

// MustNewSecKey converts []byte to a SecKey. Panics is []byte is not the exact size
func MustNewSecKey(b []byte) SecKey {
	p, err := NewSecKey(b)
	if err != nil {
		log.Panic(err)
	}
	return p
}

// MustSecKeyFromHex decodes a hex encoded SecKey, or panics
func MustSecKeyFromHex(s string) SecKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return MustNewSecKey(b)
}

// SecKeyFromHex decodes a hex encoded SecKey, or panics
func SecKeyFromHex(s string) (SecKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return SecKey{}, errors.New("Invalid secret key")
	}
	return NewSecKey(b)
}

// Verify attempts to determine if SecKey is valid. Returns nil on success.
// If DebugLevel2, will do additional sanity checking
func (sk SecKey) Verify() error {
	return sk.verify(DebugLevel2)
}

func (sk SecKey) verify(debugLevel2Check bool) error {
	if secp256k1.VerifySeckey(sk[:]) != 1 {
		return errors.New("Invalid secret key")
	}

	if debugLevel2Check {
		if err := CheckSecKey(sk); err != nil {
			log.Panicf("DebugLevel2, WARNING CRYPTO ARMAGEDDON: %v", err)
		}
	}

	return nil
}

// Hex returns a hex encoded SecKey string
func (sk SecKey) Hex() string {
	return hex.EncodeToString(sk[:])
}

//ECDH generates a shared secret
// A: pub1,sec1
// B: pub2,sec2
// person A sends their public key pub1
// person B sends an emphameral pubkey pub2
// person A computes cipher.ECDH(pub2, sec1)
// person B computes cipher.ECDH(pub1, sec2)
// cipher.ECDH(pub2, sec1) equals cipher.ECDH(pub1, sec2)
// This is their shared secret
func ECDH(pub PubKey, sec SecKey) ([]byte, error) {
	if err := pub.Verify(); err != nil {
		return nil, errors.New("ECDH invalid pubkey input")
	}

	// Don't perform the DebugLevel2 verification check for the secret key,
	// it is too slow to use in an ECDH context and is not important for that use case
	if err := sec.verify(false); err != nil {
		return nil, errors.New("ECDH invalid seckey input")
	}

	buff := secp256k1.ECDH(pub[:], sec[:])
	ret := SumSHA256(buff) // hash this so they cant screw up
	return ret[:], nil
}

// MustECDH calls ECDH and panics on error
func MustECDH(pub PubKey, sec SecKey) []byte {
	r, err := ECDH(pub, sec)
	if err != nil {
		log.Panic(err)
	}
	return r
}

// Sig signature
type Sig [64 + 1]byte //64 byte signature with 1 byte for key recovery

// NewSig converts []byte to a Sig
func NewSig(b []byte) (Sig, error) {
	s := Sig{}
	if len(b) != len(s) {
		return Sig{}, errors.New("Invalid signature length")
	}
	copy(s[:], b[:])
	return s, nil
}

// MustNewSig converts []byte to a Sig. Panics is []byte is not the exact size
func MustNewSig(b []byte) Sig {
	s := Sig{}
	if len(b) != len(s) {
		log.Panic("Invalid signature length")
	}
	copy(s[:], b[:])
	return s
}

// SigFromHex converts a hex string to a signature
func SigFromHex(s string) (Sig, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Sig{}, errors.New("Invalid signature")
	}
	return NewSig(b)
}

// MustSigFromHex converts a hex string to a signature, panics on error
func MustSigFromHex(s string) Sig {
	sig, err := SigFromHex(s)
	if err != nil {
		log.Panic(err)
	}
	return sig
}

// Hex converts signature to hex string
func (s Sig) Hex() string {
	return hex.EncodeToString(s[:])
}

// SignHash sign hash
func SignHash(hash SHA256, sec SecKey) (Sig, error) {
	if secp256k1.VerifySeckey(sec[:]) != 1 {
		// can't use sec.Verify() because that calls SignHash again, with DebugLevel2 set
		return Sig{}, errors.New("Invalid secret key")
	}

	s := secp256k1.Sign(hash[:], sec[:])

	sig, err := NewSig(s)
	if err != nil {
		return Sig{}, err
	}

	if DebugLevel2 || DebugLevel1 {
		// Guard against coin loss;
		// if the generated signature is somehow invalid, coins would be lost,
		// make sure that the signature is valid
		pubkey, err := PubKeyFromSig(sig, hash)
		if err != nil {
			log.Panic("MustSignHash error: pubkey from sig recovery failure")
		}
		if VerifySignature(pubkey, sig, hash) != nil {
			log.Panic("MustSignHash error: secp256k1.Sign returned non-null invalid non-null signature")
		}
		if ChkSig(AddressFromPubKey(pubkey), hash, sig) != nil {
			log.Panic("MustSignHash error: ChkSig failed for signature")
		}
	}

	return sig, nil
}

// MustSignHash sign hash, panics on error
func MustSignHash(hash SHA256, sec SecKey) Sig {
	sig, err := SignHash(hash, sec)
	if err != nil {
		log.Panic(err)
	}
	return sig
}

// ChkSig checks whether PubKey corresponding to address hash signed hash
// - recovers the PubKey from sig and hash
// - fail if PubKey cannot be be recovered
// - computes the address from the PubKey
// - fail if recovered address does not match PubKey hash
// - verify that signature is valid for hash for PubKey
func ChkSig(address Address, hash SHA256, sig Sig) error {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return errors.New("Invalig sig: PubKey recovery failed")
	}

	pubKey, err := NewPubKey(rawPubKey)
	if err != nil {
		return err
	}

	if address != AddressFromPubKey(pubKey) {
		return errors.New("Invalid sig: address does not match output address")
	}

	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey[:]) != 1 {
		return errors.New("Invalid sig: invalid for hash")
	}

	return nil
}

// VerifySignedHash this only checks that the signature can be converted to a public key
// Since there is no pubkey or address argument, it cannot check that the
// signature is valid in that context.
func VerifySignedHash(sig Sig, hash SHA256) error {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return errors.New("Failed to recover public key")
	}
	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey) != 1 {
		// If this occurs, secp256k1 is bugged
		log.Printf("Recovered public key is not valid for signed hash")
		return errors.New("Signature invalid for hash")
	}
	return nil
}

// VerifySignature verifies that hash was signed by PubKey
func VerifySignature(pubkey PubKey, sig Sig, hash SHA256) error {
	pubkeyRec, err := PubKeyFromSig(sig, hash) //recovered pubkey
	if err != nil {
		return errors.New("Invalid sig: PubKey recovery failed")
	}
	if pubkeyRec != pubkey {
		return errors.New("Recovered pubkey does not match pubkey")
	}
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		if DebugLevel2 {
			if secp256k1.VerifySignature(hash[:], sig[:], pubkey[:]) == 1 {
				log.Panic("VerifySignature warning, ")
			}
		}
		return errors.New("VerifySignature, secp256k1.VerifyPubkey failed")
	}
	if secp256k1.VerifySignatureValidity(sig[:]) != 1 {
		return errors.New("VerifySignature, VerifySignatureValidity failed")
	}
	if secp256k1.VerifySignature(hash[:], sig[:], pubkey[:]) != 1 {
		return errors.New("Invalid signature for this message")
	}
	return nil
}

// GenerateKeyPair creates key pair
func GenerateKeyPair() (PubKey, SecKey) {
	public, secret := secp256k1.GenerateKeyPair()

	secKey, err := NewSecKey(secret)
	if err != nil {
		log.Panicf("GenerateKeyPair: secp256k1.GenerateKeyPair returned invalid secKey: %v", err)
	}

	pubKey, err := NewPubKey(public)
	if err != nil {
		log.Panicf("GenerateKeyPair: secp256k1.GenerateKeyPair returned invalid pubKey: %v", err)
	}

	if DebugLevel1 {
		if err := CheckSecKey(secKey); err != nil {
			log.Panicf("DebugLevel1, GenerateKeyPair, generated private key failed CheckSecKey: %v", err)
		}

		if MustPubKeyFromSecKey(secKey) != pubKey {
			log.Panic("DebugLevel1, GenerateKeyPair, public key does not match private key")
		}
	}

	return pubKey, secKey
}

// GenerateDeterministicKeyPair generates deterministic key pair
func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey, error) {
	if len(seed) == 0 {
		return PubKey{}, SecKey{}, errors.New("seed input is empty")
	}

	public, secret := secp256k1.GenerateDeterministicKeyPair(seed)

	secKey, err := NewSecKey(secret)
	if err != nil {
		log.Panicf("GenerateDeterministicKeyPair: secp256k1.GenerateDeterministicKeyPair returned invalid secKey: %v", err)
	}

	pubKey, err := NewPubKey(public)
	if err != nil {
		log.Panicf("GenerateDeterministicKeyPair: secp256k1.GenerateDeterministicKeyPair returned invalid pubKey: %v", err)
	}

	if DebugLevel1 {
		if err := CheckSecKey(secKey); err != nil {
			log.Panicf("DebugLevel1, GenerateDeterministicKeyPair, CheckSecKey failed: %v", err)
		}

		if MustPubKeyFromSecKey(secKey) != pubKey {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, public key does not match private key")
		}
	}

	return pubKey, secKey, nil
}

// MustGenerateDeterministicKeyPair generates deterministic key pair, panics on error
func MustGenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
	p, s, err := GenerateDeterministicKeyPair(seed)
	if err != nil {
		log.Panic(err)
	}
	return p, s
}

// DeterministicKeyPairIterator takes SHA256 value, returns a new
// SHA256 value and publickey and private key. Apply multiple times
// feeding the SHA256 value back into generate sequence of keys
func DeterministicKeyPairIterator(seed []byte) ([]byte, PubKey, SecKey, error) {
	if len(seed) == 0 {
		return nil, PubKey{}, SecKey{}, errors.New("seed input is empty")
	}

	hash, public, secret := secp256k1.DeterministicKeyPairIterator(seed)

	secKey := MustNewSecKey(secret)
	pubKey := MustNewPubKey(public)

	if DebugLevel1 {
		if err := CheckSecKey(secKey); err != nil {
			log.Panicf("DebugLevel1, DeterministicKeyPairIterator, CheckSecKey failed: %v", err)
		}

		if MustPubKeyFromSecKey(secKey) != pubKey {
			log.Panic("DebugLevel1, DeterministicKeyPairIterator, public key does not match private key")
		}
	}

	return hash, pubKey, secKey, nil
}

// MustDeterministicKeyPairIterator takes SHA256 value, returns a new
// SHA256 value and publickey and private key. Apply multiple times
// feeding the SHA256 value back into generate sequence of keys, panics on error
func MustDeterministicKeyPairIterator(seed []byte) ([]byte, PubKey, SecKey) {
	hash, p, s, err := DeterministicKeyPairIterator(seed)
	if err != nil {
		log.Panic(err)
	}
	return hash, p, s
}

// GenerateDeterministicKeyPairs returns sequence of n private keys from initial seed
func GenerateDeterministicKeyPairs(seed []byte, n int) ([]SecKey, error) {
	_, keys, err := GenerateDeterministicKeyPairsSeed(seed, n)
	return keys, err
}

// MustGenerateDeterministicKeyPairs returns sequence of n private keys from initial seed, panics on error
func MustGenerateDeterministicKeyPairs(seed []byte, n int) []SecKey {
	keys, err := GenerateDeterministicKeyPairs(seed, n)
	if err != nil {
		log.Panic(err)
	}
	return keys
}

// GenerateDeterministicKeyPairsSeed returns sequence of n private keys from initial seed, and return the new seed
func GenerateDeterministicKeyPairsSeed(seed []byte, n int) ([]byte, []SecKey, error) {
	var keys []SecKey
	var seckey SecKey
	for i := 0; i < n; i++ {
		var err error
		seed, _, seckey, err = DeterministicKeyPairIterator(seed)
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, seckey)
	}
	return seed, keys, nil
}

// MustGenerateDeterministicKeyPairsSeed returns sequence of n private keys from initial seed, and return the new seed
func MustGenerateDeterministicKeyPairsSeed(seed []byte, n int) ([]byte, []SecKey) {
	newSeed, keys, err := GenerateDeterministicKeyPairsSeed(seed, n)
	if err != nil {
		log.Panic(err)
	}
	return newSeed, keys
}

// CheckSecKey test seckey hash
func CheckSecKey(seckey SecKey) error {
	hash := SumSHA256([]byte(time.Now().String()))
	return CheckSecKeyHash(seckey, hash)
}

// CheckSecKeyHash performs a series of tests to determine if a seckey is valid.
// All generated keys and keys loaded from disc must pass the CheckSecKey suite.
// TestPrivKey returns error if a key fails any test in the test suite.
func CheckSecKeyHash(seckey SecKey, hash SHA256) error {
	// check seckey with verify
	if secp256k1.VerifySeckey(seckey[:]) != 1 {
		return errors.New("Seckey verification failed")
	}

	// check pubkey recovery
	pubkey, err := PubKeyFromSecKey(seckey)
	if err != nil {
		return fmt.Errorf("PubKeyFromSecKey failed: %v", err)
	}
	if pubkey == (PubKey{}) {
		return errors.New("impossible error, CheckSecKey, nil pubkey recovered")
	}
	// verify recovered pubkey
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		return errors.New("impossible error, CheckSecKey, Derived Pubkey verification failed")
	}

	// check signature production
	sig, err := SignHash(hash, seckey)
	if err != nil {
		return fmt.Errorf("SignHash failed: %v", err)
	}

	pubkey2, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return fmt.Errorf("PubKeyFromSig failed: %v", err)
	}
	if pubkey != pubkey2 {
		return errors.New("Recovered pubkey does not match signed hash")
	}

	// check pubkey recovered from sig
	recoveredPubkey, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return fmt.Errorf("impossible error, CheckSecKey, pubkey recovery from signature failed: %v", err)
	}
	if pubkey != recoveredPubkey {
		return errors.New("impossible error CheckSecKey, pubkey does not match recovered pubkey")
	}

	// verify produced signature
	err = VerifySignature(pubkey, sig, hash)
	if err != nil {
		return fmt.Errorf("impossible error, CheckSecKey, verify signature failed for sig: %v", err)
	}

	// verify ChkSig
	addr := AddressFromPubKey(pubkey)
	err = ChkSig(addr, hash, sig)
	if err != nil {
		return fmt.Errorf("impossible error CheckSecKey, ChkSig Failed, should not get this far: %v", err)
	}

	// verify VerifySignedHash
	err = VerifySignedHash(sig, hash)
	if err != nil {
		return fmt.Errorf("VerifySignedHash failed: %v", err)
	}

	return nil
}

func init() {
	ripemd160HashPool = make(chan hash.Hash, ripemd160HashPoolSize)
	for i := 0; i < ripemd160HashPoolSize; i++ {
		ripemd160HashPool <- ripemd160.New()
	}

	sha256HashPool = make(chan hash.Hash, sha256HashPoolSize)
	for i := 0; i < sha256HashPoolSize; i++ {
		sha256HashPool <- sha256.New()
	}

	// Do not allow program to start if crypto tests fail
	pubkey, seckey := GenerateKeyPair()
	if err := CheckSecKey(seckey); err != nil {
		log.Fatalf("CRYPTOGRAPHIC INTEGRITY CHECK FAILED: TERMINATING PROGRAM TO PROTECT COINS: %v", err)
	}
	if MustPubKeyFromSecKey(seckey) != pubkey {
		log.Fatal("DebugLevel1, GenerateKeyPair, public key does not match private key")
	}
}
