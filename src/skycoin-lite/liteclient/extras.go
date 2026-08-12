package liteclient

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

/*
Functions used mainly during test procedures.

Like the rest of the package these report failure by returning an error. They
used to panic, both directly and through the cipher Must* constructors, which
made a malformed hex string from a caller indistinguishable from a real
verification failure and killed the wasm module outright under TinyGo.
*/

var (
	// ErrInvalidSecKey is returned when a secret key does not verify.
	ErrInvalidSecKey = errors.New("invalid secret key")

	// ErrInvalidPubKey is returned when a public key does not verify.
	ErrInvalidPubKey = errors.New("invalid public key")
)

// VerifyPubKeySignedHash verifies that hash was signed by PubKey
func VerifyPubKeySignedHash(pubkey, sig, hash string) error {
	p, err := cipher.PubKeyFromHex(pubkey)
	if err != nil {
		return err
	}

	s, err := cipher.SigFromHex(sig)
	if err != nil {
		return err
	}

	h, err := cipher.SHA256FromHex(hash)
	if err != nil {
		return err
	}

	return cipher.VerifyPubKeySignedHash(p, s, h)
}

// VerifyAddressSignedHash checks whether PubKey corresponding to address hash signed hash
// - recovers the PubKey from sig and hash
// - fail if PubKey cannot be be recovered
// - computes the address from the PubKey
// - fail if recovered address does not match PubKey hash
// - verify that signature is valid for hash for PubKey
func VerifyAddressSignedHash(address, sig, hash string) error {
	a, err := cipher.DecodeBase58Address(address)
	if err != nil {
		return err
	}

	h, err := cipher.SHA256FromHex(hash)
	if err != nil {
		return err
	}

	s, err := cipher.SigFromHex(sig)
	if err != nil {
		return err
	}

	return cipher.VerifyAddressSignedHash(a, s, h)
}

// VerifySignatureRecoverPubKey this only checks that the signature can be converted to a public key.
// It does not check that the signature signed the hash.
// The original public key or address is required to verify that the signature signed the hash.
func VerifySignatureRecoverPubKey(sig, hash string) error {
	s, err := cipher.SigFromHex(sig)
	if err != nil {
		return err
	}

	h, err := cipher.SHA256FromHex(hash)
	if err != nil {
		return err
	}

	return cipher.VerifySignatureRecoverPubKey(s, h)
}

// VerifySeckey validate a private key
func VerifySeckey(seckey string) error {
	s, err := cipher.SecKeyFromHex(seckey)
	if err != nil {
		return err
	}

	if secp256k1.VerifySeckey(s[:]) != 1 {
		return ErrInvalidSecKey
	}

	return nil
}

// VerifyPubkey validate a public key
func VerifyPubkey(pubkey string) error {
	p, err := cipher.PubKeyFromHex(pubkey)
	if err != nil {
		return err
	}

	if secp256k1.VerifyPubkey(p[:]) != 1 {
		return ErrInvalidPubKey
	}

	return nil
}

// AddressFromPubKey creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubkey string) (string, error) {
	p, err := cipher.PubKeyFromHex(pubkey)
	if err != nil {
		return "", err
	}

	return cipher.AddressFromPubKey(p).String(), nil
}

// AddressFromSecKey generates address from secret key
func AddressFromSecKey(seckey string) (string, error) {
	s, err := cipher.SecKeyFromHex(seckey)
	if err != nil {
		return "", err
	}

	address, err := cipher.AddressFromSecKey(s)
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

// PubKeyFromSig recovers the public key from a signed hash
func PubKeyFromSig(sig string, hash string) (string, error) {
	s, err := cipher.SigFromHex(sig)
	if err != nil {
		return "", err
	}

	h, err := cipher.SHA256FromHex(hash)
	if err != nil {
		return "", err
	}

	pubKey, err := cipher.PubKeyFromSig(s, h)
	if err != nil {
		return "", err
	}

	return pubKey.Hex(), nil
}

// SignHash sign hash
func SignHash(hash string, seckey string) (string, error) {
	h, err := cipher.SHA256FromHex(hash)
	if err != nil {
		return "", err
	}

	s, err := cipher.SecKeyFromHex(seckey)
	if err != nil {
		return "", err
	}

	sig, err := cipher.SignHash(h, s)
	if err != nil {
		return "", err
	}

	return sig.Hex(), nil
}
