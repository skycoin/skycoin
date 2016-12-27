package cipher

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"log"
	"time"

	"github.com/skycoin/skycoin/src/cipher/ripemd160"
	"github.com/skycoin/skycoin/src/util"

	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

var (
	logger      = util.MustGetLogger("crypto")
	DebugLevel1 = true //checks for extremely unlikely conditions (10e-40)
	DebugLevel2 = true //enable checks for impossible conditions
)

type PubKey [33]byte
type PubKeySlice []PubKey

func (slice PubKeySlice) Len() int {
	return len(slice)
}

func (slice PubKeySlice) Less(i, j int) bool {
	return bytes.Compare(slice[i][:], slice[j][:]) < 0
}

func (slice PubKeySlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func RandByte(n int) []byte {
	return secp256k1.RandByte(n)
}

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

func PubKeyFromHex(s string) (PubKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return PubKey{}, errors.New("Invalid public key")
	}
	return NewPubKey(b), nil
}

// Recovers the public key for a secret key
func PubKeyFromSecKey(seckey SecKey) PubKey {
	if seckey == (SecKey{}) {
		log.Panic("PubKeyFromSecKey, attempt to load null seckey, unsafe")
	}
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

// Decodes a hex encoded SecKey, or panics
func SecKeyFromHex(s string) (SecKey, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return SecKey{}, errors.New("Invalid SecKey: not valid hex")
	}
	if len(b) != 32 {
		return SecKey{}, errors.New("Invalid SecKey: invalid length")
	}
	return NewSecKey(b), nil
}

// Verify attempts to determine if SecKey is valid. Returns nil on success.
// If DebugLevel2, will do additional sanity checking
func (self SecKey) Verify() error {
	if secp256k1.VerifySeckey(self[:]) != 1 {
		return errors.New("Invalid SecKey")
	}
	if DebugLevel2 {
		err := TestSecKey(self)
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

//Generates a shared secret
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
		log.Panic("ECDH invalid pubkey input")
	}

	if err := sec.Verify(); err != nil {
		log.Panic("ECDH invalid seckey input")
	}

	buff := secp256k1.ECDH(pub[:], sec[:])
	ret := SumSHA256(buff) //hash this so they cant screw up
	return ret[:]

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
	if len(b) != 65 {
		log.Panic("Signature Length is Invalid")
	}
	return NewSig(b)
}

func SigFromHex(s string) (Sig, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Sig{}, err
	}
	if len(b) != 65 {
		return Sig{}, errors.New("Signature Length is Invalid")
	}
	return NewSig(b), nil
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
		if TestSecKey(NewSecKey(secret)) != nil {
			log.Panic("DebugLevel1, GenerateKeyPair, generated private key " +
				"failed TestSecKey")
		}
	}

	return NewPubKey(public), NewSecKey(secret)
}

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

//DeterministicKeyPairIterator takes SHA256 value, returns a new
//SHA256 value and publickey and private key. Apply multiple times
//feeding the SHA256 value back into generate sequence of keys
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

//Returns sequence of n private keys from intial seed, and return the new seed
func GenerateDeterministicKeyPairsSeed(seed []byte, n int) ([]byte, []SecKey) {
	var keys []SecKey
	var seckey SecKey
	for i := 0; i < n; i++ {
		seed, _, seckey = DeterministicKeyPairIterator(seed)
		keys = append(keys, seckey)
	}
	return seed, keys
}

func TestSecKey(seckey SecKey) error {
	hash := SumSHA256([]byte(time.Now().String()))
	return TestSecKeyHash(seckey, hash)
}

// TestPrivKey performs a series of tests to determine if a seckey is valid.
// All generated keys and keys loaded from disc must pass the TestSecKey suite.
// TestPrivKey returns error if a key fails any test in the test suite.
func TestSecKeyHash(seckey SecKey, hash SHA256) error {
	//check seckey with verify
	if secp256k1.VerifySeckey(seckey[:]) != 1 {
		return errors.New("Seckey verification failed")
	}

	//check pubkey recovery
	pubkey := PubKeyFromSecKey(seckey)
	if pubkey == (PubKey{}) {
		errors.New("impossible error, TestSecKey, nil pubkey recovered")
	}
	//verify recovered pubkey
	if secp256k1.VerifyPubkey(pubkey[:]) != 1 {
		return errors.New("impossible error, TestSecKey, Derived Pubkey " +
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
		return errors.New("impossible error, TestSecKey, pubkey recovery " +
			"from signature failed")
	}
	if pubkey != recoveredPubkey {
		return errors.New("impossible error TestSecKey, pubkey does not " +
			"match recovered pubkey")
	}

	//verify produced signature
	err = VerifySignature(pubkey, sig, hash)
	if err != nil {
		errors.New("impossible error, TestSecKey, verify signature failed " +
			"for sig")
	}

	//verify ChkSig
	addr := AddressFromPubKey(pubkey)
	err = ChkSig(addr, hash, sig)
	if err != nil {
		return errors.New("impossible error TestSecKey, ChkSig Failed, " +
			"should not get this far")
	}

	return nil
}

//do not allow program to start if crypto tests fail
func init() {
	// init the reuse hash pool.
	sha256HashChan = make(chan hash.Hash, poolsize)
	ripemd160HashChan = make(chan hash.Hash, poolsize)
	for i := 0; i < poolsize; i++ {
		sha256HashChan <- sha256.New()
		ripemd160HashChan <- ripemd160.New()
	}

	_, seckey := GenerateKeyPair()
	if TestSecKey(seckey) != nil {
		log.Fatal("CRYPTOGRAPHIC INTEGRITY CHECK FAILED: TERMINATING " +
			"PROGRAM TO PROTECT COINS")
	}
}
