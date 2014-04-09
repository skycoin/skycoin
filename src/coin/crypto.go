package coin

import (
	"encoding/hex"
	"errors"
	"github.com/skycoin/skycoin/src/lib/secp256k1-go"
	"log"
	"time"
)

type PubKey [33]byte

// Converts []byte to a PubKey. Panics is []byte is not the exact size
func NewPubKey(b []byte) PubKey {
	p := PubKey{}
	if len(b) != len(p) {
		log.Panic("Invalid public key length")
	}
	copy(p[:], b[:])
	return p
}

// Decodes a hex encoded PubKey, or panics
func MustPubKeyFromHex(s string) PubKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return NewPubKey(b)
}

// Recovers the public key for a secret key
func PubKeyFromSecKey(seckey SecKey) PubKey {
	b := secp256k1.PubkeyFromSeckey(seckey[:])
	if b == nil {
		log.Panic("PubKeyFromSecKey, pubkey recovery failed. Function " +
			"assumes seckey is valid. Check seckey")
	}
	return NewPubKey(b)
}

// Recovers the public key from a signed hash
func PubKeyFromSig(sig Sig, hash SHA256) (PubKey, error) {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return PubKey{}, errors.New("Invalig sig: PubKey recovery failed")
	}
	return NewPubKey(rawPubKey), nil
}

//Verify attempts to determine if pubkey is valid. Returns nil on success
func (self PubKey) Verify() error {
	if secp256k1.VerifyPubkey(self[:]) != 1 {
		return errors.New("Invalid public key")
	}
	return nil
}

// Returns a hex encoded PubKey string
func (self PubKey) Hex() string {
	return hex.EncodeToString(self[:])
}

// Returns the public key as ripemd160(sha256(sha256(key)))
func (self *PubKey) ToAddressHash() Ripemd160 {
	r1 := SumSHA256(self[:])
	r2 := SumSHA256(r1[:])
	return HashRipemd160(r2[:])
}

type SecKey [32]byte

// Converts []byte to a SecKey. Panics is []byte is not the exact size
func NewSecKey(b []byte) SecKey {
	p := SecKey{}
	if len(b) != len(p) {
		log.Panic("Invalid secret key length")
	}
	copy(p[:], b[:])
	return p
}

// Decodes a hex encoded SecKey, or panics
func MustSecKeyFromHex(s string) SecKey {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return NewSecKey(b)
}

// Verify attempts to determine if SecKey is valid. Returns nil on success.
// If DebugLevel2, will do additional sanity checking
func (self SecKey) Verify() error {
	if secp256k1.VerifySeckey(self[:]) != 1 {
		return errors.New("Invalid SecKey")
	}
	if DebugLevel2 {
		err := testSecKey(self)
		if err != nil {
			log.Panic("DebugLevel2, WARNING CRYPTO ARMAGEDDON")
		}
	}
	return nil
}

// Returns a hex encoded SecKey string
func (s SecKey) Hex() string {
	return hex.EncodeToString(s[:])
}

type Sig [64 + 1]byte //64 byte signature with 1 byte for key recovery

// Converts []byte to a Sig. Panics is []byte is not the exact size
func NewSig(b []byte) Sig {
	s := Sig{}
	if len(b) != len(s) {
		log.Panic("Invalid secret key length")
	}
	copy(s[:], b[:])
	return s
}

// Decodes a hex-encoded Sig, panicing if invalid
func MustSigFromHex(s string) Sig {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic(err)
	}
	return NewSig(b)
}

func (s Sig) Hex() string {
	return hex.EncodeToString(s[:])
}

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

/*
	Checks whether PubKey corresponding to address hash signed hash
	- recovers the PubKey from sig and hash
	- fail if PubKey cannot be be recovered
	- computes the address from the PubKey
	- fail if recovered address does not match PubKey hash
	- verify that signature is valid for hash for PubKey
*/
func ChkSig(address Address, hash SHA256, sig Sig) error {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return errors.New("Invalig sig: PubKey recovery failed")
	}
	if address != AddressFromPubKey(NewPubKey(rawPubKey)) {
		return errors.New("Invalid sig: address does not match output address")
	}
	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey[:]) != 1 {
		return errors.New("Invalid sig: invalid for hash")
	}
	return nil
}

// This only checks that the signature can be converted to a public key
// Since there is no pubkey or address argument, it cannot check that the
// signature is valid in that context.
func VerifySignedHash(sig Sig, hash SHA256) error {
	rawPubKey := secp256k1.RecoverPubkey(hash[:], sig[:])
	if rawPubKey == nil {
		return errors.New("Failed to recover public key")
	}
	if secp256k1.VerifySignature(hash[:], sig[:], rawPubKey) != 1 {
		// If this occurs, secp256k1 is bugged
		logger.Critical("Recovered public key is not valid for signed hash")
		return errors.New("Signature invalid for hash")
	}
	return nil
}

// Verifies that hash was signed by PubKey
func VerifySignature(pubkey PubKey, sig Sig, hash SHA256) error {
	pubkeyRec, err := PubKeyFromSig(sig, hash) //recovered pubkey
	if err != nil {
		return errors.New("Invalig sig: PubKey recovery failed")
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

func GenerateKeyPair() (PubKey, SecKey) {
	public, secret := secp256k1.GenerateKeyPair()

	if DebugLevel1 {
		if testSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateKeyPair, generated private key " +
				"failed testSecKey")
		}
	}

	return NewPubKey(public), NewSecKey(secret)
}

func GenerateDeterministicKeyPair(seed []byte) (PubKey, SecKey) {
	public, secret := secp256k1.GenerateDeterministicKeyPair(seed)

	if DebugLevel1 {

		if testSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"seckey invalid, failed testSecKey")
		}
		if testSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"generated private key failed testSecKey")
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

//DeterministicKeyPairIterator takes SHA256 value, returns a new
//SHA256 value and publickey and private key. Apply multiple times
//feeding the SHA256 value back into generate sequence of keys
func DeterministicKeyPairIterator(seed []byte) ([]byte, PubKey, SecKey) {
	hash, public, secret := secp256k1.DeterministicKeyPairIterator(seed)
	if DebugLevel1 {
		if testSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"generated private key failed testSecKey")
		}
		if PubKeyFromSecKey(NewSecKey(secret)) != NewPubKey(public) {
			log.Panic("DebugLevel1, GenerateDeterministicKeyPair, " +
				"public key does not match private key")
		}
	}
	return hash, NewPubKey(public), NewSecKey(secret)
}

//Returns sequence of n private keys from intial seed
func GenerateDeterministicKeyPairs(seed []byte, n int) []SecKey {
	var keys []SecKey
	var seckey SecKey
	for i := 0; i < n; i++ {
		seed, _, seckey = DeterministicKeyPairIterator(seed)
		keys = append(keys, seckey)
	}
	return keys
}

func testSecKey(seckey SecKey) error {
	hash := SumSHA256([]byte(time.Now().String()))
	return testSecKeyHash(seckey, hash)
}

// TestPrivKey performs a series of tests to determine if a seckey is valid.
// All generated keys and keys loaded from disc must pass the testSecKey suite.
// TestPrivKey returns error if a key fails any test in the test suite.
func testSecKeyHash(seckey SecKey, hash SHA256) error {
	//check seckey with verify
	if secp256k1.VerifySeckey(seckey[:]) != 1 {
		return errors.New("Seckey verification failed")
	}

	//check pubkey recovery
	pubkey := PubKeyFromSecKey(seckey)
	if pubkey == (PubKey{}) {
		errors.New("impossible error, testSecKey, nil pubkey recovered")
	}
	//verify recovered pubkey
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		return errors.New("impossible error, testSecKey, Derived Pubkey " +
			"verification failed")
	}

	//check signature production
	sig := SignHash(hash, seckey)
	pubkey2, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return err
	}
	if pubkey != pubkey2 {
		return errors.New("Recovered pubkey does not match signed hash")
	}

	//check pubkey recovered from sig
	recoveredPubkey, err := PubKeyFromSig(sig, hash)
	if err != nil {
		return errors.New("impossible error, testSecKey, pubkey recovery " +
			"from signature failed")
	}
	if pubkey != recoveredPubkey {
		return errors.New("impossible error testSecKey, pubkey does not " +
			"match recovered pubkey")
	}

	//verify produced signature
	err = VerifySignature(pubkey, sig, hash)
	if err != nil {
		errors.New("impossible error, testSecKey, verify signature failed " +
			"for sig")
	}

	//verify ChkSig
	addr := AddressFromPubKey(pubkey)
	err = ChkSig(addr, hash, sig)
	if err != nil {
		return errors.New("impossible error testSecKey, ChkSig Failed, " +
			"should not get this far")
	}

	return nil
}

//do not allow program to start if crypto tests fail
func init() {
	_, seckey := GenerateKeyPair()
	if testSecKey(seckey) != nil {
		log.Fatal("CRYPTOGRAPHIC INTEGRITY CHECK FAILED: TERMINATING " +
			"PROGRAM TO PROTECT COINS")
	}
}
