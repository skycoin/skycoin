// Wrapper around secp256k1-go
package util

import (
	"errors"
	"lib/secp256k1-go"
	"log"
)

func GenerateSignature(seckey []byte, msg []byte) []byte {
	if secp256.VerifySeckey(seckey) != 1 {
		log.Panic("Invalid secret key")
		return nil
	}
	return secp256.Sign(msg, seckey) // test that signature is valid
}

func VerifySignature(pubkey []byte, msg []byte, sig []byte) error {
	if secp256.VerifyPubkey(pubkey) != 1 {
		log.Panic("Invalid public key")
		return errors.New("Invalid public key")
	}
	if secp256.VerifySignatureValidity(sig) != 1 {
		log.Panic("Invalid signature")
		return errors.New("Invalid signature")
	}
	if secp256.VerifySignature(msg, sig, pubkey) != 1 {
		return errors.New("Invalid signature for this message")
	}
	return nil
}

func GenerateKeyPair() (public, secret []byte) {
	public, secret = secp256.GenerateKeyPair()
	return
}
