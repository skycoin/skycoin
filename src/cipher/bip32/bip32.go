package bip32

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"log"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

const (
	// FirstHardenedChild is the index of the firxt "harded" child key as per the
	// bip32 spec
	FirstHardenedChild = uint32(0x80000000)

	// publicKeyCompressedLength is the byte count of a compressed public key
	publicKeyCompressedLength = 33

	// masterKey is the "Key" value specified by bip32 for the master key
	masterKey = "Bitcoin seed"
)

var (
	// PrivateWalletVersion is the version flag for serialized private keys ("xpriv")
	PrivateWalletVersion = []byte{0x04, 0x88, 0xAD, 0xE4}

	// PublicWalletVersion is the version flag for serialized public keys ("xpub")
	PublicWalletVersion = []byte{0x04, 0x88, 0xB2, 0x1E}

	// ErrSerializedKeyWrongSize is returned when trying to deserialize a key that
	// has an incorrect length
	ErrSerializedKeyWrongSize = errors.New("Serialized keys should by exactly 82 bytes")

	// ErrHardenedChildPublicKey is returned when trying to create a harded child
	// of the public key
	ErrHardenedChildPublicKey = errors.New("Can't create hardened child for public key")

	// ErrInvalidChecksum is returned when deserializing a key with an incorrect
	// checksum
	ErrInvalidChecksum = errors.New("Checksum doesn't match")

	// ErrInvalidPrivateKey is returned when a derived private key is invalid
	ErrInvalidPrivateKey = errors.New("Invalid private key")

	// ErrInvalidPublicKey is returned when a derived public key is invalid
	ErrInvalidPublicKey = errors.New("Invalid public key")

	// ErrInvalidSeedLength is returned when generating a master key with an invalid number of seed bits
	ErrInvalidSeedLength = errors.New("Invalid master key seed length")

	// ErrInvalidPublicKeyBytesLength deserialize public key from bytes length was not 33
	ErrInvalidPublicKeyBytesLength = errors.New("Public keys have 33 bytes")

	// ErrInvalidPrivateKeyBytesLength deserialize private key from bytes length was not 32
	ErrInvalidPrivateKeyBytesLength = errors.New("Private keys have 32 bytes")

	// ErrDeserializePrivateFromPublic attempted to deserialize a private key from an encoded public key
	ErrDeserializePrivateFromPublic = errors.New("Cannot deserialize a private key from a public key")
)

// key represents a bip32 extended key
type key struct {
	Version     []byte // 4 bytes
	Depth       byte   // 1 bytes
	Fingerprint []byte // 4 bytes
	ChildNumber []byte // 4 bytes
	ChainCode   []byte // 32 bytes
	Key         []byte // 33 bytes for public keys; 32 bytes for private keys
}

// PrivateKey represents a bip32 extended private key
type PrivateKey struct {
	key
}

// PublicKey represents a bip32 extended public key
type PublicKey struct {
	key
}

// NewMasterKey creates a new master extended key from a seed.
// Seed should be between 128 and 512 bits; 256 bits are recommended.
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#master-key-generation
func NewMasterKey(seed []byte) (*PrivateKey, error) {
	if len(seed) < 16 || len(seed) > 64 {
		return nil, ErrInvalidSeedLength
	}

	return newMasterKey(seed)
}

func newMasterKey(seed []byte) (*PrivateKey, error) {
	// Generate key and chaincode
	hmac := hmac.New(sha512.New, []byte(masterKey))
	if _, err := hmac.Write(seed); err != nil {
		log.Panic(err)
	}
	intermediary := hmac.Sum(nil)

	// Split it into our key and chain code
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]

	// Validate key
	if err := validatePrivateKey(keyBytes); err != nil {
		return nil, err
	}

	// Create the key struct
	key := &PrivateKey{
		key: key{
			Version:     PrivateWalletVersion,
			ChainCode:   chainCode,
			Key:         keyBytes,
			Depth:       0x0,
			ChildNumber: []byte{0x00, 0x00, 0x00, 0x00},
			// Master key fingerprint specified to be 0x00000000 since it has no parent
			Fingerprint: []byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	return key, nil

}

// PublicKey returns the public version of key or return a copy
// The 'Neuter' function from the bip32 spec, N((k, c) -> (K, c).
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--public-child-key
func (k *PrivateKey) PublicKey() *PublicKey {
	pubKey, err := publicKeyForPrivateKey(k.Key)
	if err != nil {
		log.Panicf("PrivateKey.PublicKey failed: %v", err)
	}
	return &PublicKey{
		key: key{
			Version:     PublicWalletVersion,
			Key:         pubKey,
			Depth:       k.Depth,
			ChildNumber: k.ChildNumber,
			Fingerprint: k.Fingerprint,
			ChainCode:   k.ChainCode,
		},
	}
}

func (k *PrivateKey) fingerprint() []byte {
	// "Extended keys can be identified by the Hash160 (RIPEMD160 after SHA256)
	// of the serialized ECDSA public key K, ignoring the chain code."
	return k.PublicKey().fingerprint()
}

func (k *PublicKey) fingerprint() []byte {
	return fingerprint(k.Key)
}

func fingerprint(k []byte) []byte {
	// "Extended keys can be identified by the Hash160 (RIPEMD160 after SHA256)
	// of the serialized ECDSA public key K, ignoring the chain code."

	// ripemd160(sha256(key))
	fp := hash160(k)

	// "The first 32 bits of the identifier are called the key fingerprint."
	return fp[:4]
}

// NewPrivateChildKey derives a private child key from a given parent as outlined by bip32, CDKpriv().
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--private-child-key
func (k *PrivateKey) NewPrivateChildKey(childIdx uint32) (*PrivateKey, error) {
	intermediary := k.ckdPrivHMAC(childIdx)

	iL := intermediary[:32]        // used for computing the next key
	chainCode := intermediary[32:] // iR

	// ki = parse256(IL) + kpar (mod n)
	// In case parse256(IL) ≥ n or ki = 0, the resulting key is invalid
	// TODO -- bip32 says we should move to the next childIdx if this fails, need to return a value
	// that can be identified as this particular failure so the caller knows to retrys
	newKey, err := addPrivateKeys(iL, k.Key)
	if err != nil {
		return nil, err
	}

	// Precalculate the fingerprint and cache it.
	// It could be calculated on demand, instead of doing it here.
	fp := k.fingerprint()

	return &PrivateKey{
		key: key{
			Version:     PrivateWalletVersion,
			ChildNumber: uint32Bytes(childIdx),
			ChainCode:   chainCode,
			Depth:       k.Depth + 1,
			Fingerprint: fp,
			Key:         newKey,
		},
	}, nil
}

// NewPublicChildKey derives a public child key from an extended public key, N(CKDpriv()).
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--public-child-key
func (k *PrivateKey) NewPublicChildKey(childIdx uint32) (*PublicKey, error) {
	k2, err := k.NewPrivateChildKey(childIdx)
	if err != nil {
		return nil, err
	}
	return k2.PublicKey(), nil
}

// NewPublicChildKey derives a public child key from an extended public key, CKDpub().
// Hardened child keys cannot be derived; the value of childIdx must be less than 2^31.
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#public-parent-key--public-child-key
func (k *PublicKey) NewPublicChildKey(childIdx uint32) (*PublicKey, error) {
	// CKDPub step 1
	intermediary, err := k.ckdPubHMAC(childIdx)
	if err != nil {
		return nil, err
	}

	// CKDPub step 2
	iL := intermediary[:32] // used for computing the next key

	// CKDPub step 3: Ki = point(parse256(IL)) + Kpar
	// CKDPub step 5: In case parse256(IL) ≥ n or Ki is the point at infinity, the resulting key is invalid
	// TODO -- bip32 says we should move to the next childIdx if this fails, need to return a value
	// that can be identified as this particular failure so the caller knows to retrys
	keyBytes, err := publicKeyForPrivateKey(iL)
	if err != nil {
		return nil, err
	}

	// CKDPub step 4
	chainCode := intermediary[32:] // iR

	newKey, err := addPublicKeys(keyBytes[:], k.Key)
	if err != nil {
		return nil, err
	}

	// Precalculate the fingerprint and cache it.
	// It could be calculated on demand, instead of doing it here.
	fp := k.fingerprint()

	// Create child Key with data common to all both scenarios
	return &PublicKey{
		key: key{
			Version:     PublicWalletVersion,
			ChildNumber: uint32Bytes(childIdx),
			ChainCode:   chainCode,
			Depth:       k.Depth + 1,
			Fingerprint: fp,
			Key:         newKey,
		},
	}, nil
}

// ckdPrivHMAC computes the first step of the CKDPriv function, which computes an HMAC
func (k *key) ckdPrivHMAC(childIdx uint32) []byte {
	// Get intermediary to create key and chaincode.
	// Hardened children are based on the private key.
	// Non-hardened children are based on the public key.

	var data []byte
	if childIdx >= FirstHardenedChild {
		// Hardened child
		// I = HMAC-SHA512(Key = cpar, Data = 0x00 || ser256(kpar) || ser32(i))
		// 0x00 || ser256(kpar)
		data = append([]byte{0x0}, k.Key...)
	} else {
		// Non-hardened child
		// I = HMAC-SHA512(Key = cpar, Data = serP(point(kpar)) || ser32(i))
		// The equation below is "convert private key to compressed public key"
		// serP(point(kpar))
		var err error
		data, err = publicKeyForPrivateKey(k.Key)
		if err != nil {
			log.Panic(err)
		}
	}

	// Append the child index as big-endian serialized bytes
	// || ser32(i)
	childIndexBytes := uint32Bytes(childIdx)
	data = append(data, childIndexBytes...)

	// HMAC-SHA512(Key = cpar, Data)
	hmac := hmac.New(sha512.New, k.ChainCode)
	if _, err := hmac.Write(data); err != nil {
		log.Panic(err)
	}

	return hmac.Sum(nil)
}

// ckdPubHMAC computes the first step of the CKDPub function, which computes an HMAC
func (k *key) ckdPubHMAC(childIdx uint32) ([]byte, error) {
	// Get intermediary to create key and chaincode.
	// Hardened children are based on the private key.
	// Non-hardened children are based on the public key.

	if childIdx >= FirstHardenedChild {
		// Public keys can't derive hardened child keys
		return nil, ErrHardenedChildPublicKey
	}

	// I = HMAC-SHA512(Key = cpar, Data = serP(Kpar) || ser32(i))

	// serP(Kpar)
	data := make([]byte, len(k.Key))
	copy(data, k.Key)

	// ser32(i)
	childIndexBytes := uint32Bytes(childIdx)

	// serP(Kpar) || ser32(i)
	data = append(data, childIndexBytes...)

	// HMAC-SHA512(Key = cpar, Data = serP(Kpar) || ser32(i))
	hmac := hmac.New(sha512.New, k.ChainCode)
	if _, err := hmac.Write(data); err != nil {
		log.Panic(err)
	}

	return hmac.Sum(nil), nil
}

// Serialize a Key to a 78 byte byte slice
func (k *PrivateKey) Serialize() []byte {
	// Private keys should be prepended with a single null byte
	return k.serialize(append([]byte{0x0}, k.Key...))
}

// Serialize a Key to a 78 byte byte slice
func (k *PublicKey) Serialize() []byte {
	return k.serialize(k.Key)
}

// serialize a Key to a 78 byte byte slice
func (k *key) serialize(keyBytes []byte) []byte {
	n := len(k.Version)
	n++ // k.Depth
	n += len(k.Fingerprint)
	n += len(k.ChildNumber)
	n += len(k.ChainCode)
	n += len(keyBytes)

	buffer := &bytes.Buffer{}
	buffer.Grow(n)

	// Write fields to buffer in order
	buffer.Write(k.Version)
	buffer.WriteByte(k.Depth)
	buffer.Write(k.Fingerprint)
	buffer.Write(k.ChildNumber)
	buffer.Write(k.ChainCode)
	buffer.Write(keyBytes)

	// Append the standard doublesha256 checksum
	return addChecksumToBytes(buffer.Bytes())
}

// B58Serialize encodes the Key in the standard Bitcoin base58 encoding
func (k *PrivateKey) B58Serialize() string {
	return base58.Encode(k.Serialize())
}

// String encodes the Key in the standard Bitcoin base58 encoding
func (k *PrivateKey) String() string {
	return k.B58Serialize()
}

// B58Serialize encodes the Key in the standard Bitcoin base58 encoding
func (k *PublicKey) B58Serialize() string {
	return base58.Encode(k.Serialize())
}

// String encodes the Key in the standard Bitcoin base58 encoding
func (k *PublicKey) String() string {
	return k.B58Serialize()
}

// DeserializePrivate deserializes a byte slice into a PrivateKey
func DeserializePrivate(data []byte) (*PrivateKey, error) {
	k, err := deserialize(data)
	if err != nil {
		return nil, err
	}
	if len(k.Key) != 32 {
		return nil, ErrInvalidPrivateKeyBytesLength
	}
	return &PrivateKey{
		key: *k,
	}, nil
}

// DeserializePublic deserializes a byte slice into a PublicKey
func DeserializePublic(data []byte) (*PublicKey, error) {
	k, err := deserialize(data)
	if err != nil {
		return nil, err
	}
	if len(k.Key) != 33 {
		return nil, ErrInvalidPublicKeyBytesLength
	}
	return &PublicKey{
		key: *k,
	}, nil
}

// deserialize a byte slice into a Key.
// If the Key.Key length is 32 bytes it is a private key, otherwise it is a public key.
func deserialize(data []byte) (*key, error) {
	if len(data) != 82 {
		return nil, ErrSerializedKeyWrongSize
	}

	k := &key{}
	k.Version = data[0:4]
	k.Depth = data[4]
	k.Fingerprint = data[5:9]
	k.ChildNumber = data[9:13]
	k.ChainCode = data[13:45]

	if data[45] == byte(0) {
		k.Key = data[46:78]
	} else {
		k.Key = data[45:78]
	}

	// validate checksum
	cs1 := checksum(data[0 : len(data)-4])
	cs2 := data[len(data)-4:]
	if !bytes.Equal(cs1, cs2) {
		return nil, ErrInvalidChecksum
	}

	return k, nil
}
