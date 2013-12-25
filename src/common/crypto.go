package sb

import (
	"log"
)
import "lib/secp256k1-go"

/* TODO
   - optimize hashing function
*/

func SHA256(data []byte) [32]byte {
	sha256_hash.Reset()
	sha256_hash.Write(data)
	sum := sha256_hash.Sum(nil)
	var out [32]byte
	copy(out[0:32], sum[0:32])
	return out
}

func HashFunc(data []byte) [32]byte {
	sha256_hash.Reset()
	sha256_hash.Write(data)
	sum := sha256_hash.Sum(nil)
	var out [32]byte
	copy(out[0:32], sum[0:32])
	return out
}

func HashBuffer(data []byte, out []byte) {
	h := HashFunc(data)
	copy(out[0:32], h[0:32])
}

func GenerateSignature(seckey []byte, msg []byte) []byte {
	if secp256.VerifySeckey(seckey) != 1 {
		log.Panic()
	}
	sig := secp256.Sign(msg, seckey) //test that signature is valid
	return sig
}

func VerifySignature(pubkey []byte, msg []byte, sig []byte) int {
	if secp256.VerifyPubkey(pubkey) != 1 {
		log.Panic()
	}
	if secp256.VerifySignatureValidity(sig) != 1 {
		log.Panic()
	}
	return secp256.VerifySignature(msg, sig, pubkey)
}

func GenerateKeyPair() ([]byte, []byte) {
	seckey, pubkey := secp256.GenerateKeyPair()
	return seckey, pubkey
}
