package bip32

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"

	"github.com/SkycoinProject/skycoin/src/cipher/ripemd160"
	secp256k1 "github.com/SkycoinProject/skycoin/src/cipher/secp256k1-go"
	secp256k1go "github.com/SkycoinProject/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

//
// Hashes
//

func hashSHA256(data []byte) []byte {
	h := sha256.New()
	h.Write(data) //nolint:errcheck
	return h.Sum(nil)
}

func hashDoubleSHA256(data []byte) []byte {
	hash1 := hashSHA256(data)
	hash2 := hashSHA256(hash1)
	return hash2
}

func hashRipemd160(data []byte) []byte {
	h := ripemd160.New()
	h.Write(data) //nolint:errcheck
	return h.Sum(nil)
}

func hash160(data []byte) []byte {
	hash1 := hashSHA256(data)
	hash2 := hashRipemd160(hash1)
	return hash2
}

//
// Checksum
//

func checksum(data []byte) []byte {
	hash := hashDoubleSHA256(data)
	return hash[:4]
}

func addChecksumToBytes(data []byte) []byte {
	checksum := checksum(data)
	return append(data, checksum...)
}

//
// Keys
//

// publicKeyForPrivateKey converts a private key to a public key.
// Equivalent to `serP(point(k))` in the bip32 spec
func publicKeyForPrivateKey(key []byte) ([]byte, error) {
	// From bip32: If parse256(IL) ≥ n, fail
	if err := validatePrivateKey(key); err != nil {
		return nil, err
	}

	b := secp256k1.PubkeyFromSeckey(key)
	if b == nil {
		log.Panic("publicKeyForPrivateKey: invalid private key")
	}

	// From bip32: If Ki == 0 fail
	if err := validatePublicKey(b); err != nil {
		return nil, err
	}

	return b, nil
}

func addPublicKeys(key, keyPar []byte) ([]byte, error) {
	if code := secp256k1.VerifyPubkey(key); code != 1 {
		return nil, fmt.Errorf("addPublicKeys: key is invalid, secp256k1.VerifyPubkey error code %d", code)
	}
	if code := secp256k1.VerifyPubkey(keyPar); code != 1 {
		return nil, fmt.Errorf("addPublicKeys: keyPar is invalid, secp256k1.VerifyPubkey error code %d", code)
	}

	// expandPublicKey
	var pk1, pk2 secp256k1go.XY
	if err := pk1.ParsePubkey(key); err != nil {
		log.Panicf("addPublicKeys: invalid pubkey1: %v", err)
	}
	if err := pk2.ParsePubkey(keyPar); err != nil {
		log.Panicf("addPublicKeys: invalid pubkey1: %v", err)
	}

	// add public keys
	pk1.AddXY(&pk2)

	// compress
	newKey := pk1.Bytes()

	if code := secp256k1.VerifyPubkey(newKey); code != 1 {
		return nil, fmt.Errorf("addPublicKeys: newKey is invalid, secp256k1.VerifyPubkey error code %d", code)
	}

	return newKey, nil
}

// addPrivateKeys computes the CKDPriv equation `parse256(IL) + kpar (mod n)`
// and verifies the result
func addPrivateKeys(key, keyPar []byte) ([]byte, error) {
	// From bip32: If parse256(IL) ≥ n, fail
	if code := secp256k1.VerifySeckey(key); code != 1 {
		return nil, fmt.Errorf("addPrivateKeys: key is invalid, secp256k1.VerifySeckey error code %d", code)
	}
	if code := secp256k1.VerifySeckey(keyPar); code != 1 {
		return nil, fmt.Errorf("addPrivateKeys: keyPar is invalid, secp256k1.VerifySeckey error code %d", code)
	}

	var keyInt big.Int
	var keyParInt big.Int
	keyInt.SetBytes(key)
	keyParInt.SetBytes(keyPar)

	// Computes this CKDPriv equation:
	// parse256(IL) + kpar (mod n)
	keyInt.Add(&keyInt, &keyParInt)
	keyInt.Mod(&keyInt, &secp256k1go.TheCurve.Order.Int)

	newKey := secp256k1go.LeftPadBytes(keyInt.Bytes(), 32)

	// From bip32: If ki == 0 fail
	if code := secp256k1.VerifySeckey(newKey); code != 1 {
		return nil, fmt.Errorf("addPrivateKeys: newKey is invalid, secp256k1.VerifySeckey error code %d", code)
	}

	return newKey, nil
}

// validatePrivateKey verifies that the secret key is not zero and that it is inside the curve
// Corresponds to bip32 spec constraints `parse256(IL) < n && ki != 0`
func validatePrivateKey(key []byte) error {
	// VerifySeckey checks that the key is > 0 and inside the curve
	if secp256k1.VerifySeckey(key) != 1 {
		return ErrDerivedInvalidPrivateKey
	}

	return nil
}

func validatePublicKey(key []byte) error {
	if secp256k1.VerifyPubkey(key) != 1 {
		return ErrDerivedInvalidPublicKey
	}

	return nil
}

//
// Numerical
//

// uint32Bytes serializes a uint32 as bytes in big-endian form
func uint32Bytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}
