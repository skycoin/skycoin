package liteclient

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

/*
Functions used mainly during test procedures.
*/

// VerifyPubKeySignedHash verifies that hash was signed by PubKey
func VerifyPubKeySignedHash(pubkey, sig, hash string) {
	p := cipher.MustPubKeyFromHex(pubkey)
	s := cipher.MustSigFromHex(sig)
	h := cipher.MustSHA256FromHex(hash)

	err := cipher.VerifyPubKeySignedHash(p, s, h)
	if err != nil {
		panic(err)
	}
}

// VerifyAddressSignedHash checks whether PubKey corresponding to address hash signed hash
// - recovers the PubKey from sig and hash
// - fail if PubKey cannot be be recovered
// - computes the address from the PubKey
// - fail if recovered address does not match PubKey hash
// - verify that signature is valid for hash for PubKey
func VerifyAddressSignedHash(address, sig, hash string) {
	a := cipher.MustDecodeBase58Address(address)
	h := cipher.MustSHA256FromHex(hash)
	s := cipher.MustSigFromHex(sig)

	err := cipher.VerifyAddressSignedHash(a, s, h)
	if err != nil {
		panic(err)
	}
}

// VerifySignatureRecoverPubKey this only checks that the signature can be converted to a public key.
// It does not check that the signature signed the hash.
// The original public key or address is required to verify that the signature signed the hash.
func VerifySignatureRecoverPubKey(sig, hash string) {
	s := cipher.MustSigFromHex(sig)
	h := cipher.MustSHA256FromHex(hash)

	err := cipher.VerifySignatureRecoverPubKey(s, h)
	if err != nil {
		panic(err)
	}
}

// VerifySeckey validate a private key
func VerifySeckey(seckey string) int {
	s := cipher.MustSecKeyFromHex(seckey)
	return secp256k1.VerifySeckey(s[:])
}

// VerifyPubkey validate a public key
func VerifyPubkey(pubkey string) int {
	p := cipher.MustPubKeyFromHex(pubkey)
	return secp256k1.VerifyPubkey(p[:])
}

// AddressFromPubKey creates Address from PubKey as ripemd160(sha256(sha256(pubkey)))
func AddressFromPubKey(pubkey string) string {
	p := cipher.MustPubKeyFromHex(pubkey)
	return cipher.AddressFromPubKey(p).String()
}

// AddressFromSecKey generates address from secret key
func AddressFromSecKey(seckey string) string {
	s := cipher.MustSecKeyFromHex(seckey)
	return cipher.MustAddressFromSecKey(s).String()
}

// PubKeyFromSig recovers the public key from a signed hash
func PubKeyFromSig(sig string, hash string) string {
	s := cipher.MustSigFromHex(sig)
	h := cipher.MustSHA256FromHex(hash)

	pubKey, err := cipher.PubKeyFromSig(s, h)
	if err != nil {
		panic(err)
	}

	return pubKey.Hex()
}

// SignHash sign hash
func SignHash(hash string, seckey string) string {
	h := cipher.MustSHA256FromHex(hash)
	s := cipher.MustSecKeyFromHex(seckey)

	sig := cipher.MustSignHash(h, s)
	return sig.Hex()
}
