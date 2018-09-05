package cipher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	skyerrors "github.com/skycoin/skycoin/src/util/errors"
)

var (
	// DebugLevel1 debug level one
	DebugLevel1 = true //checks for extremely unlikely conditions (10e-40)
	// DebugLevel2 debug level two
	DebugLevel2 = true //enable checks for impossible conditions

	// ErrInvalidLengthPubKey  Invalid public key length
	ErrInvalidLengthPubKey = errors.New("Invalid public key length")
	// ErrPubKeyFromNullSecKey PubKeyFromSecKey, attempt to load null seckey, unsafe
	ErrPubKeyFromNullSecKey = errors.New("PubKeyFromSecKey, attempt to load null seckey, unsafe")
	// ErrPubKeyFromBadSecKey  PubKeyFromSecKey, pubkey recovery failed. Function
	ErrPubKeyFromBadSecKey = errors.New("PubKeyFromSecKey, pubkey recovery failed. Function " +
		"assumes seckey is valid. Check seckey")
	// ErrInvalidLengthSecKey Invalid secret key length
	ErrInvalidLengthSecKey = errors.New("Invalid secret key length")
	// ErrECHDInvalidPubKey   ECDH invalid pubkey input
	ErrECHDInvalidPubKey = errors.New("ECDH invalid pubkey input")
	// ErrECHDInvalidSecKey   ECDH invalid seckey input
	ErrECHDInvalidSecKey = errors.New("ECDH invalid seckey input")
	// ErrInvalidLengthSig    Invalid signature length
	ErrInvalidLengthSig = errors.New("Invalid signature length")
	// ErrInvalidPubKey       Invalid public key
	ErrInvalidPubKey = errors.New("Invalid public key")
	// ErrInvalidSecKey       Invalid public key
	ErrInvalidSecKey = errors.New("Invalid secret key")
	// ErrInvalidSigForPubKey Invalig sig: PubKey recovery failed
	ErrInvalidSigForPubKey = errors.New("Invalig sig: PubKey recovery failed")
	// ErrInvalidSecKeyHex    Invalid SecKey: not valid hex
	ErrInvalidSecKeyHex = errors.New("Invalid SecKey: not valid hex")
	// ErrInvalidAddressForSig Invalid sig: address does not match output address
	ErrInvalidAddressForSig = errors.New("Invalid sig: address does not match output address")
	// ErrInvalidHashForSig   Signature invalid for hash
	ErrInvalidHashForSig = errors.New("Signature invalid for hash")
	// ErrPubKeyRecoverMismatch Recovered pubkey does not match pubkey
	ErrPubKeyRecoverMismatch = errors.New("Recovered pubkey does not match pubkey")
	// ErrInvalidSigInvalidPubKey VerifySignature, secp256k1.VerifyPubkey failed
	ErrInvalidSigInvalidPubKey = errors.New("VerifySignature, secp256k1.VerifyPubkey failed")
	// ErrInvalidSigValidity  VerifySignature, VerifySignatureValidity failed
	ErrInvalidSigValidity = errors.New("VerifySignature, VerifySignatureValidity failed")
	// ErrInvalidSigForMessage Invalid signature for this message
	ErrInvalidSigForMessage = errors.New("Invalid signature for this message")
	// ErrInvalidSecKyVerification Seckey secp256k1 verification failed
	ErrInvalidSecKyVerification = errors.New("Seckey verification failed")
	// ErrNullPubKeyFromSecKey Impossible error, TestSecKey, nil pubkey recovered
	ErrNullPubKeyFromSecKey = errors.New("impossible error, TestSecKey, nil pubkey recovered")
	// ErrInvalidDerivedPubKeyFromSecKey impossible error, TestSecKey, Derived Pubkey verification failed
	ErrInvalidDerivedPubKeyFromSecKey = errors.New("impossible error, TestSecKey, Derived Pubkey verification failed")
	// ErrInvalidPubKeyFromHash Recovered pubkey does not match signed hash
	ErrInvalidPubKeyFromHash = errors.New("Recovered pubkey does not match signed hash")
	// ErrPubKeyFromSecKeyMissmatch impossible error TestSecKey, pubkey does not match recovered pubkey
	ErrPubKeyFromSecKeyMissmatch = errors.New("impossible error TestSecKey, pubkey does not match recovered pubkey")
)

// PubKey public key
type PubKey [33]byte

// PubKeySlice PubKey slice
type PubKeySlice []PubKey

// Len returns length for sorting
func (slice PubKeySlice) Len() int {
	return len(slice)
}

// Less for sorting
func (slice PubKeySlice) Less(i, j int) bool {
	return bytes.Compare(slice[i][:], slice[j][:]) < 0
}

// Swap for sorting
func (slice PubKeySlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// RandByte returns rand N bytes
func RandByte(n int) []byte {
	return secp256k1.RandByte(n)
}

// NewPubKey converts []byte to a PubKey. Panics is []byte is not the exact size
func NewPubKey(b []byte) PubKey {
	p := PubKey{}
	if len(b) != len(p) {
		err := skyerrors.NewValueError(ErrInvalidLengthPubKey, "b", b)
		log.Print(err)
		panic(err)
	}
	copy(p[:], b[:])
	return p
}

// MustPubKeyFromHex decodes a hex encoded PubKey, or panics
func MustPubKeyFromHex(s string) PubKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return NewPubKey(b)
}

// PubKeyFromHex generates PubKey from hex string
func PubKeyFromHex(s string) (PubKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PubKey{}, ErrInvalidPubKey
	}
	if len(b) != len(PubKey{}) {
		return PubKey{}, ErrInvalidLengthPubKey
	}
	return NewPubKey(b), nil
}

// PubKeyFromSecKey recovers the public key for a secret key
func PubKeyFromSecKey(seckey SecKey) PubKey {
	if seckey == (SecKey{}) {
		err := skyerrors.NewValueError(ErrPubKeyFromNullSecKey, "seckey", seckey)
		log.Print(err)
		panic(err)
	}
	b := secp256k1.PubkeyFromSeckey(seckey[:])
	if b == nil {
		err := skyerrors.NewValueError(ErrPubKeyFromBadSecKey, "seckey", seckey)
		log.Print(err)
		panic(err)
	}
	return NewPubKey(b)
}

// PubKeyFromSig recovers the public key from a signed hash
func PubKeyFromSig(sig Sig, hash SHA256) (PubKey, error) {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return PubKey{}, ErrInvalidSigForPubKey
	}
	return NewPubKey(rawPubKey), nil
}

// Verify attempts to determine if pubkey is valid. Returns nil on success
func (pk PubKey) Verify() error {
	if secp256k1.VerifyPubkey(pk[:]) != 1 {
		return ErrInvalidPubKey
	}
	return nil
}

// Hex returns a hex encoded PubKey string
func (pk PubKey) Hex() string {
	return hex.EncodeToString(pk[:])
}

// ToAddressHash returns the public key as ripemd160(sha256(sha256(key)))
func (pk *PubKey) ToAddressHash() Ripemd160 {
	r1 := SumSHA256(pk[:])
	r2 := SumSHA256(r1[:])
	return HashRipemd160(r2[:])
}

// SecKey secret key
type SecKey [32]byte

// NewSecKey converts []byte to a SecKey. Panics is []byte is not the exact size
func NewSecKey(b []byte) SecKey {
	p := SecKey{}
	if len(b) != len(p) {
		err := skyerrors.NewValueError(ErrInvalidLengthSecKey, "b", b)
		log.Print(err)
		panic(err)
	}
	copy(p[:], b[:])
	return p
}

// MustSecKeyFromHex decodes a hex encoded SecKey, or panics
func MustSecKeyFromHex(s string) SecKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return NewSecKey(b)
}

// SecKeyFromHex decodes a hex encoded SecKey, or panics
func SecKeyFromHex(s string) (SecKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return SecKey{}, ErrInvalidSecKeyHex
	}
	if len(b) != 32 {
		return SecKey{}, ErrInvalidLengthSecKey
	}
	return NewSecKey(b), nil
}

// Verify attempts to determine if SecKey is valid. Returns nil on success.
// If DebugLevel2, will do additional sanity checking
func (sk SecKey) Verify() error {
	if secp256k1.VerifySeckey(sk[:]) != 1 {
		return ErrInvalidSecKey
	}
	if DebugLevel2 {
		err := TestSecKey(sk)
		if err != nil {
			log.Panic("DebugLevel2, WARNING CRYPTO ARMAGEDDON")
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
func ECDH(pub PubKey, sec SecKey) []byte {

	if err := pub.Verify(); err != nil {
		err := skyerrors.NewValueError(ErrECHDInvalidPubKey, "pub", pub)
		log.Print(err)
		panic(err)
	}

	// WARNING: This calls TestSecKey if DebugLevel2 is set to true.
	// TestSecKey is extremely slow and will kill performance if ECDH is called frequently
	if err := sec.Verify(); err != nil {
		err := skyerrors.NewValueError(ErrECHDInvalidSecKey, "sec", sec)
		log.Print(err)
		panic(err)
	}

	buff := secp256k1.ECDH(pub[:], sec[:])
	ret := SumSHA256(buff) //hash this so they cant screw up
	return ret[:]

}

// Sig signature
type Sig [64 + 1]byte //64 byte signature with 1 byte for key recovery

// NewSig converts []byte to a Sig. Panics is []byte is not the exact size
func NewSig(b []byte) Sig {
	s := Sig{}
	if len(b) != len(s) {
		err := skyerrors.NewValueError(ErrInvalidLengthSig, "b", b)
		log.Print(err)
		panic(err)
	}
	copy(s[:], b[:])
	return s
}

// MustSigFromHex decodes a hex-encoded Sig, panicing if invalid
func MustSigFromHex(s string) Sig {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	if len(b) != 65 {
		err := skyerrors.NewValueError(ErrInvalidLengthSig, "s", s)
		log.Print(err)
		panic(err)
	}
	return NewSig(b)
}

// SigFromHex generates signature from hex string
func SigFromHex(s string) (Sig, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Sig{}, err
	}
	if len(b) != 65 {
		return Sig{}, ErrInvalidLengthSig
	}
	return NewSig(b), nil
}

// Hex converts signature to hex string
func (s Sig) Hex() string {
	return hex.EncodeToString(s[:])
}

// SignHash sign hash
func SignHash(hash SHA256, sec SecKey) Sig {
	sig := NewSig(secp256k1.Sign(hash[:], sec[:]))

	if DebugLevel2 || DebugLevel1 { //!!! Guard against coin loss
		pubkey, err := PubKeyFromSig(sig, hash)
		if err != nil {
			log.Panic("SignHash, error: pubkey from sig recovery failure")
		}
		if VerifySignature(pubkey, sig, hash) != nil {
			log.Panic("SignHash, error: secp256k1.Sign returned non-null " +
				"invalid non-null signature")
		}
		if ChkSig(AddressFromPubKey(pubkey), hash, sig) != nil {
			log.Panic("SignHash error: ChkSig failed for signature")
		}
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
		return ErrInvalidSigForPubKey
	}
	if address != AddressFromPubKey(NewPubKey(rawPubKey)) {
		return ErrInvalidAddressForSig
	}
	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey[:]) != 1 {
		return ErrInvalidHashForSig
	}
	return nil
}

// VerifySignedHash this only checks that the signature can be converted to a public key
// Since there is no pubkey or address argument, it cannot check that the
// signature is valid in that context.
func VerifySignedHash(sig Sig, hash SHA256) error {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return ErrInvalidSigForPubKey
	}
	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey) != 1 {
		// If this occurs, secp256k1 is bugged
		log.Printf("Recovered public key is not valid for signed hash")
		return ErrInvalidHashForSig
	}
	return nil
}

// VerifySignature verifies that hash was signed by PubKey
func VerifySignature(pubkey PubKey, sig Sig, hash SHA256) error {
	pubkeyRec, err := PubKeyFromSig(sig, hash) //recovered pubkey
	if err != nil {
		return ErrInvalidSigForPubKey
	}
	if pubkeyRec != pubkey {
		return ErrPubKeyRecoverMismatch
	}
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		if DebugLevel2 {
			if secp256k1.VerifySignature(hash[:], sig[:], pubkey[:]) == 1 {
				log.Panic("VerifySignature warning, ")
			}
		}
		return ErrInvalidSigInvalidPubKey
	}
	if secp256k1.VerifySignatureValidity(sig[:]) != 1 {
		return ErrInvalidSigValidity
	}
	if secp256k1.VerifySignature(hash[:], sig[:], pubkey[:]) != 1 {
		return ErrInvalidSigForMessage
	}
	return nil
}

// GenerateKeyPair creates key pair
func GenerateKeyPair() (PubKey, SecKey) {
	public, secret := secp256k1.GenerateKeyPair()

	if DebugLevel1 {
		if TestSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateKeyPair, generated private key " +
				"failed TestSecKey")
		}
	}

	return NewPubKey(public), NewSecKey(secret)
}

// GenerateDeterministicKeyPair generates deterministic key pair
func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
	public, secret := secp256k1.GenerateDeterministicKeyPair(seed)

	if DebugLevel1 {

		if TestSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"seckey invalid, failed TestSecKey")
		}
		if TestSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"generated private key failed TestSecKey")
		}
		if PubKeyFromSecKey(NewSecKey(secret)) != NewPubKey(public) {
			//s1 := NewSecKey(secret).Hex()
			//s2 := NewPubKey(public).Hex()
			//s3 := PubKeyFromSecKey(NewSecKey(secret)).Hex()
			//log.Printf("sec= %s, pub= %s recpub= %s \n", s1,s2, s3 )
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"public key does not match private key")
		}
	}
	return NewPubKey(public), NewSecKey(secret)
}

// DeterministicKeyPairIterator takes SHA256 value, returns a new
// SHA256 value and publickey and private key. Apply multiple times
// feeding the SHA256 value back into generate sequence of keys
func DeterministicKeyPairIterator(seed []byte) ([]byte, PubKey, SecKey) {
	hash, public, secret := secp256k1.DeterministicKeyPairIterator(seed)
	if DebugLevel1 {
		if TestSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"generated private key failed TestSecKey")
		}
		if PubKeyFromSecKey(NewSecKey(secret)) != NewPubKey(public) {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"public key does not match private key")
		}
	}
	return hash, NewPubKey(public), NewSecKey(secret)
}

// GenerateDeterministicKeyPairs returns sequence of n private keys from initial seed
func GenerateDeterministicKeyPairs(seed []byte, n int) []SecKey {
	var keys []SecKey
	var seckey SecKey
	for i := 0; i < n; i++ {
		seed, _, seckey = DeterministicKeyPairIterator(seed)
		keys = append(keys, seckey)
	}
	return keys
}

// GenerateDeterministicKeyPairsSeed returns sequence of n private keys from initial seed, and return the new seed
func GenerateDeterministicKeyPairsSeed(seed []byte, n int) ([]byte, []SecKey) {
	var keys []SecKey
	var seckey SecKey
	for i := 0; i < n; i++ {
		seed, _, seckey = DeterministicKeyPairIterator(seed)
		keys = append(keys, seckey)
	}
	return seed, keys
}

// TestSecKey test seckey hash
func TestSecKey(seckey SecKey) error {
	hash := SumSHA256([]byte(time.Now().String()))
	return TestSecKeyHash(seckey, hash)
}

// TestSecKeyHash performs a series of tests to determine if a seckey is valid.
// All generated keys and keys loaded from disc must pass the TestSecKey suite.
// TestPrivKey returns error if a key fails any test in the test suite.
func TestSecKeyHash(seckey SecKey, hash SHA256) error {
	//check seckey with verify
	if secp256k1.VerifySeckey(seckey[:]) != 1 {
		return ErrInvalidSecKyVerification
	}

	//check pubkey recovery
	pubkey := PubKeyFromSecKey(seckey)
	if pubkey == (PubKey{}) {
		return ErrNullPubKeyFromSecKey
	}
	//verify recovered pubkey
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		return ErrInvalidDerivedPubKeyFromSecKey
	}

	//check signature production
	sig := SignHash(hash, seckey)
	pubkey2, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return fmt.Errorf("PubKeyFromSig failed: %v", err)
	}
	if pubkey != pubkey2 {
		return ErrInvalidPubKeyFromHash
	}

	//check pubkey recovered from sig
	recoveredPubkey, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return fmt.Errorf("impossible error, TestSecKey, pubkey recovery from signature failed: %v", err)
	}
	if pubkey != recoveredPubkey {
		return ErrPubKeyFromSecKeyMissmatch
	}

	//verify produced signature
	err = VerifySignature(pubkey, sig, hash)
	if err != nil {
		return fmt.Errorf("impossible error, TestSecKey, verify signature failed for sig: %v", err)
	}

	//verify ChkSig
	addr := AddressFromPubKey(pubkey)
	err = ChkSig(addr, hash, sig)
	if err != nil {
		return fmt.Errorf("impossible error TestSecKey, ChkSig Failed, should not get this far: %v", err)
	}

	//verify VerifySignedHash
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
	_, seckey := GenerateKeyPair()
	if err := TestSecKey(seckey); err != nil {
		log.Fatalf("CRYPTOGRAPHIC INTEGRITY CHECK FAILED: TERMINATING PROGRAM TO PROTECT COINS: %v", err)
	}
}
