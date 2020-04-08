/*
Package bip32 implements the bip32 spec https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
*/
package bip32

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"github.com/SkycoinProject/skycoin/src/cipher/base58"
)

const (
	// FirstHardenedChild is the index of the firxt "hardened" child key as per the
	// bip32 spec
	FirstHardenedChild = uint32(0x80000000)

	// masterKey is the "Key" value specified by bip32 for the master key
	masterKey = "Bitcoin seed"
)

// Error wraps bip32 errors
type Error struct {
	error
	impossibleChild bool
}

// NewError creates an Error
func NewError(err error) error {
	if err == nil {
		log.Panic("NewError called with nil error")
	}
	return Error{
		error: err,
	}
}

// ImpossibleChild returns true if this error indicates
// that a given child number cannot produce a valid key.
// If the caller receives this error, they should skip this child and go to the next number.
// The probability of this happening is less than 1 in 2^127.
func (err Error) ImpossibleChild() bool {
	return err.impossibleChild
}

// NewImpossibleChildError creates an Error flagged as an ImpossibleChild error
func NewImpossibleChildError(err error, childNumber uint32) error {
	if err == nil {
		log.Panic("NewImpossibleChildError called with nil error")
	}
	return Error{
		error:           fmt.Errorf("childNumber=%d %v", childNumber, err),
		impossibleChild: true,
	}
}

// IsImpossibleChildError returns true if the error is an ImpossibleChild Error
func IsImpossibleChildError(err error) bool {
	if err == nil {
		return false
	}

	switch t := err.(type) {
	case Error:
		return t.ImpossibleChild()
	default:
		return false
	}
}

var (
	// PrivateWalletVersion is the version flag for serialized private keys ("xpriv")
	PrivateWalletVersion = []byte{0x04, 0x88, 0xAD, 0xE4}

	// PublicWalletVersion is the version flag for serialized public keys ("xpub")
	PublicWalletVersion = []byte{0x04, 0x88, 0xB2, 0x1E}

	// ErrSerializedKeyWrongSize is returned when trying to deserialize a key that
	// has an incorrect length
	ErrSerializedKeyWrongSize = NewError(errors.New("Serialized keys should be exactly 82 bytes"))

	// ErrHardenedChildPublicKey is returned when trying to create a harded child
	// of the public key
	ErrHardenedChildPublicKey = NewError(errors.New("Can't create hardened child for public key"))

	// ErrInvalidChecksum is returned when deserializing a key with an incorrect checksum
	ErrInvalidChecksum = NewError(errors.New("Checksum doesn't match"))

	// ErrDerivedInvalidPrivateKey is returned when an invalid private key was derived
	ErrDerivedInvalidPrivateKey = NewError(errors.New("Derived invalid private key"))

	// ErrDerivedInvalidPublicKey is returned when an invalid public key was derived
	ErrDerivedInvalidPublicKey = NewError(errors.New("Derived invalid public key"))

	// ErrInvalidPrivateKeyVersion is returned when a deserializing a private key without the 'xprv' prefix
	ErrInvalidPrivateKeyVersion = NewError(errors.New("Invalid private key version"))

	// ErrInvalidPublicKeyVersion is returned when a deserializing a public key without the 'xpub' prefix
	ErrInvalidPublicKeyVersion = NewError(errors.New("Invalid public key version"))

	// ErrInvalidSeedLength is returned when generating a master key with an invalid number of seed bits
	ErrInvalidSeedLength = NewError(errors.New("Invalid master key seed length"))

	// ErrDeserializePrivateFromPublic attempted to deserialize a private key from an encoded public key
	ErrDeserializePrivateFromPublic = NewError(errors.New("Cannot deserialize a private key from a public key"))

	// ErrInvalidKeyVersion is returned if the key version is not 'xpub' or 'xprv'
	ErrInvalidKeyVersion = NewError(errors.New("Invalid key version"))

	// ErrInvalidFingerprint is returned if a deserialized key has an invalid fingerprint
	ErrInvalidFingerprint = NewError(errors.New("Invalid key fingerprint"))

	// ErrInvalidChildNumber is returned if a deserialized key has an invalid child number
	ErrInvalidChildNumber = NewError(errors.New("Invalid key child number"))

	// ErrInvalidPrivateKey is returned if a deserialized xprv key's private key is invalid
	ErrInvalidPrivateKey = NewError(errors.New("Invalid private key"))

	// ErrInvalidPublicKey is returned if a deserialized xpub key's public key is invalid
	ErrInvalidPublicKey = NewError(errors.New("Invalid public key"))

	// ErrMaxDepthReached maximum allowed depth (255) reached for child key
	ErrMaxDepthReached = NewError(errors.New("Maximum child depth reached"))
)

// key represents a bip32 extended key
type key struct {
	Version           []byte // 4 bytes
	Depth             byte   // 1 bytes
	ParentFingerprint []byte // 4 bytes
	childNumber       []byte // 4 bytes
	ChainCode         []byte // 32 bytes
	Key               []byte // 33 bytes for public keys; 32 bytes for private keys
}

func (k key) ChildNumber() uint32 {
	return binary.BigEndian.Uint32(k.childNumber)
}

func (k key) clone() key {
	newK := key{}
	newK.Depth = k.Depth
	newK.Version = append(newK.Version, k.Version...)
	newK.ParentFingerprint = append(newK.ParentFingerprint, k.ParentFingerprint...)
	newK.childNumber = append(newK.childNumber, k.childNumber...)
	newK.ChainCode = append(newK.ChainCode, k.ChainCode...)
	newK.Key = append(newK.Key, k.Key...)
	return newK
}

// PrivateKey represents a bip32 extended private key
type PrivateKey struct {
	key
}

// Clone returns a copy of the private key
func (k PrivateKey) Clone() PrivateKey {
	return PrivateKey{key: k.clone()}
}

// PublicKey represents a bip32 extended public key
type PublicKey struct {
	key
}

// Clone returns a copy of the public key
func (k PublicKey) Clone() PublicKey {
	return PublicKey{key: k.clone()}
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
			Version:           PrivateWalletVersion,
			ChainCode:         chainCode,
			Key:               keyBytes,
			Depth:             0x0,
			childNumber:       []byte{0x00, 0x00, 0x00, 0x00},
			ParentFingerprint: []byte{0x00, 0x00, 0x00, 0x00},
		},
	}

	return key, nil
}

// NewPrivateKeyFromPath returns a private key at a given bip32 path.
// The path must be a full path starting with m/, and the initial seed
// must be provided.
// This method can return an ImpossibleChild error.
func NewPrivateKeyFromPath(seed []byte, p string) (*PrivateKey, error) {
	path, err := ParsePath(p)
	if err != nil {
		return nil, err
	}

	k, err := NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	if len(path.Elements) > 1 {
		return k.DeriveSubpath(path.Elements[1:])
	}

	return k, nil
}

// DeriveSubpath derives a PrivateKey at at bip32 subpath, e.g. `0'/1'/0`.
// The nodes argument must not be empty.
// This method can return an ImpossibleChild error.
func (k *PrivateKey) DeriveSubpath(nodes []PathNode) (*PrivateKey, error) {
	if len(nodes) == 0 {
		return nil, errors.New("Path nodes array empty when deriving a bip32 subpath")
	}

	ck, err := k.newPrivateChildKeyFromPathNode(nodes[0])
	if err != nil {
		return nil, err
	}

	for _, e := range nodes[1:] {
		ck, err = ck.newPrivateChildKeyFromPathNode(e)
		if err != nil {
			return nil, err
		}
	}

	return ck, nil
}

func (k *PrivateKey) newPrivateChildKeyFromPathNode(n PathNode) (*PrivateKey, error) {
	if n.Master {
		return nil, errors.New("PathNode is Master at a non-zero depth")
	}

	return k.NewPrivateChildKey(n.ChildNumber)
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
			Version:           PublicWalletVersion,
			Key:               pubKey,
			Depth:             k.Depth,
			childNumber:       k.childNumber,
			ChainCode:         k.ChainCode,
			ParentFingerprint: k.ParentFingerprint,
		},
	}
}

// Fingerprint returns the key fingerprint
func (k *PrivateKey) Fingerprint() []byte {
	// "Extended keys can be identified by the Hash160 (RIPEMD160 after SHA256)
	// of the serialized ECDSA public key K, ignoring the chain code."
	return k.PublicKey().Fingerprint()
}

// Identifier returns the key ID
func (k *PrivateKey) Identifier() []byte {
	return k.PublicKey().Identifier()
}

// Fingerprint returns the key fingerprint
func (k *PublicKey) Fingerprint() []byte {
	return fingerprint(k.Key)
}

// Identifier returns the key ID
func (k *PublicKey) Identifier() []byte {
	return identifier(k.Key)
}

func identifier(key []byte) []byte {
	// "Extended keys can be identified by the Hash160 (RIPEMD160 after SHA256)
	// of the serialized ECDSA public key K, ignoring the chain code."

	// ripemd160(sha256(key))
	return hash160(key)
}

func fingerprint(key []byte) []byte {
	id := identifier(key)
	// "The first 32 bits of the identifier are called the key fingerprint."
	return id[:4]
}

// NewPrivateChildKey derives a private child key from a given parent as outlined by bip32, CDKpriv().
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--private-child-key
// This method can return an ImpossibleChild error.
func (k *PrivateKey) NewPrivateChildKey(childIdx uint32) (*PrivateKey, error) {
	if k.Depth == 0xFF {
		return nil, ErrMaxDepthReached
	}

	intermediary := k.ckdPrivHMAC(childIdx)

	iL := intermediary[:32]        // used for computing the next key
	chainCode := intermediary[32:] // iR

	// ki = parse256(IL) + kpar (mod n)
	// In case parse256(IL) ≥ n or ki = 0, the resulting key is invalid
	newKey, err := addPrivateKeys(iL, k.Key)
	if err != nil {
		return nil, NewImpossibleChildError(err, childIdx)
	}

	return &PrivateKey{
		key: key{
			Version:           PrivateWalletVersion,
			childNumber:       uint32Bytes(childIdx),
			ChainCode:         chainCode,
			Depth:             k.Depth + 1,
			Key:               newKey,
			ParentFingerprint: k.Fingerprint(),
		},
	}, nil
}

// NewPublicChildKey derives a public child key from an extended public key, N(CKDpriv()).
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#private-parent-key--public-child-key
// This method can return an ImpossibleChild error.
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
// This method can return an ImpossibleChild error.
func (k *PublicKey) NewPublicChildKey(childIdx uint32) (*PublicKey, error) {
	if k.Depth == 0xFF {
		return nil, ErrMaxDepthReached
	}

	// CKDPub step 1
	intermediary, err := k.ckdPubHMAC(childIdx)
	if err != nil {
		return nil, err
	}

	// CKDPub step 2
	iL := intermediary[:32] // used for computing the next key

	// CKDPub step 3: Ki = point(parse256(IL)) + Kpar
	// CKDPub step 5: In case parse256(IL) ≥ n or Ki is the point at infinity, the resulting key is invalid
	keyBytes, err := publicKeyForPrivateKey(iL)
	if err != nil {
		return nil, NewImpossibleChildError(err, childIdx)
	}

	// CKDPub step 4
	chainCode := intermediary[32:] // iR

	newKey, err := addPublicKeys(keyBytes[:], k.Key)
	if err != nil {
		return nil, err
	}

	// Create child Key with data common to all both scenarios
	return &PublicKey{
		key: key{
			Version:           PublicWalletVersion,
			childNumber:       uint32Bytes(childIdx),
			ChainCode:         chainCode,
			Depth:             k.Depth + 1,
			Key:               newKey,
			ParentFingerprint: k.Fingerprint(),
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
func (k *key) serialize(key []byte) []byte {
	n := len(k.Version)
	n++ // k.Depth
	n += len(k.ParentFingerprint)
	n += len(k.childNumber)
	n += len(k.ChainCode)
	n += len(key)

	buffer := &bytes.Buffer{}
	buffer.Grow(n)

	// Write fields to buffer in order
	buffer.Write(k.Version)
	buffer.WriteByte(k.Depth)
	buffer.Write(k.ParentFingerprint)
	buffer.Write(k.childNumber)
	buffer.Write(k.ChainCode)
	buffer.Write(key)

	// Append the standard doublesha256 checksum
	return addChecksumToBytes(buffer.Bytes())
}

// String encodes the Key in the standard Bitcoin base58 encoding
func (k *PrivateKey) String() string {
	return base58.Encode(k.Serialize())
}

// String encodes the Key in the standard Bitcoin base58 encoding
func (k *PublicKey) String() string {
	return base58.Encode(k.Serialize())
}

// DeserializeEncodedPrivateKey deserializes a base58 xprv key to a PrivateKey
func DeserializeEncodedPrivateKey(xprv string) (*PrivateKey, error) {
	b, err := base58.Decode(xprv)
	if err != nil {
		return nil, err
	}
	return DeserializePrivateKey(b)
}

// DeserializePrivateKey deserializes the []byte serialization of a PrivateKey
func DeserializePrivateKey(data []byte) (*PrivateKey, error) {
	k, err := deserialize(data, true)
	if err != nil {
		return nil, err
	}

	if len(k.Key) != 32 {
		log.Panic("DeserializePrivateKey expected 32 bytes key length")
	}
	if !bytes.Equal(k.Version, PrivateWalletVersion) {
		log.Panic("DeserializePrivateKey expected xprv prefix")
	}

	return &PrivateKey{
		key: *k,
	}, nil
}

// DeserializeEncodedPublicKey deserializes a base58 xpub key to a PublicKey
func DeserializeEncodedPublicKey(xpub string) (*PublicKey, error) {
	b, err := base58.Decode(xpub)
	if err != nil {
		return nil, err
	}
	return DeserializePublicKey(b)
}

// DeserializePublicKey deserializes the []byte serialization of a PublicKey
func DeserializePublicKey(data []byte) (*PublicKey, error) {
	k, err := deserialize(data, false)
	if err != nil {
		return nil, err
	}

	if len(k.Key) != 33 {
		log.Panic("DeserializePublicKey expected 33 bytes key length")
	}
	if !bytes.Equal(k.Version, PublicWalletVersion) {
		log.Panic("DeserializePublicKey expected xpub prefix")
	}

	return &PublicKey{
		key: *k,
	}, nil
}

// deserialize a byte slice into a Key.
// If the Key.Key length is 32 bytes it is a private key, otherwise it is a public key.
func deserialize(data []byte, wantPrivate bool) (*key, error) {
	if len(data) != 82 {
		return nil, ErrSerializedKeyWrongSize
	}

	// Validate checksum
	cs1 := checksum(data[0 : len(data)-4])
	cs2 := data[len(data)-4:]
	if !bytes.Equal(cs1, cs2) {
		return nil, ErrInvalidChecksum
	}

	k := &key{}
	k.Version = data[0:4]
	k.Depth = data[4]
	k.ParentFingerprint = data[5:9]
	k.childNumber = data[9:13]
	k.ChainCode = data[13:45]

	isPrivate := bytes.Equal(k.Version, PrivateWalletVersion)
	isPublic := bytes.Equal(k.Version, PublicWalletVersion)

	if !isPrivate && !isPublic {
		return nil, ErrInvalidKeyVersion
	}

	if wantPrivate && !isPrivate {
		return nil, ErrInvalidPrivateKeyVersion
	}

	if !wantPrivate && !isPublic {
		return nil, ErrInvalidPublicKeyVersion
	}

	// Master keys (depth=0) have an empty fingerprint and a ChildNumber of 0
	if k.Depth == 0 {
		var emptyBytes [4]byte
		if !bytes.Equal(k.ParentFingerprint, emptyBytes[:]) {
			return nil, ErrInvalidFingerprint
		}

		if k.ChildNumber() != 0 {
			return nil, ErrInvalidChildNumber
		}
	}

	// Private keys must have a 0 byte prefix padding
	if isPrivate {
		if data[45] != 0 {
			return nil, ErrInvalidPrivateKey
		}
		k.Key = data[46:78]
		if err := validatePrivateKey(k.Key); err != nil {
			return nil, ErrInvalidPrivateKey
		}
	} else {
		k.Key = data[45:78]
		if err := validatePublicKey(k.Key); err != nil {
			return nil, ErrInvalidPublicKey
		}
	}

	return k, nil
}
