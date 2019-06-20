/*
Package secp256k1 provides private and public key generation with the secp256k1 elliptic curve.
*/
package secp256k1

import (
	"bytes"
	"encoding/hex"
	"log"

	secp "github.com/skycoin/skycoin/src/cipher/secp256k1-go/secp256k1-go2"
)

// DebugPrint enable debug print statements
var DebugPrint = false

// intenal, may fail
// may return nil
func pubkeyFromSeckey(seckey []byte) []byte {
	if len(seckey) != 32 {
		log.Panic("seckey length invalid")
	}

	if secp.SeckeyIsValid(seckey) != 1 {
		log.Panic("always ensure seckey is valid")
		return nil
	}

	pubkey := secp.GeneratePublicKey(seckey) // always returns true
	if pubkey == nil {
		log.Panic("ERROR: impossible, secp.BaseMultiply always returns true")
		return nil
	}
	if len(pubkey) != 33 {
		log.Panic("ERROR: impossible, invalid pubkey length")
	}

	if ret := secp.PubkeyIsValid(pubkey); ret != 1 {
		log.Panicf("ERROR: pubkey invald, ret=%d", ret)
		return nil
	}

	if ret := VerifyPubkey(pubkey); ret != 1 {
		log.Printf("seckey= %s", hex.EncodeToString(seckey))
		log.Printf("pubkey= %s", hex.EncodeToString(pubkey))
		log.Panicf("ERROR: pubkey verification failed, for deterministic. ret=%d", ret)
		return nil
	}

	return pubkey
}

// GenerateKeyPair generates public and private key pairs
func GenerateKeyPair() ([]byte, []byte) {
	const seckeyLen = 32
	var seckey []byte
	var pubkey []byte

new_seckey:
	seckey = RandByte(seckeyLen)
	if secp.SeckeyIsValid(seckey) != 1 {
		goto new_seckey // regen
	}

	pubkey = pubkeyFromSeckey(seckey)
	if pubkey == nil {
		log.Panic("IMPOSSIBLE: pubkey invalid from valid seckey")
		goto new_seckey
	}
	if ret := secp.PubkeyIsValid(pubkey); ret != 1 {
		log.Panicf("ERROR: Pubkey invalid, ret=%d", ret)
		goto new_seckey
	}

	return pubkey, seckey
}

// PubkeyFromSeckey generates a compressed public key from a secret key
func PubkeyFromSeckey(seckey []byte) []byte {
	// This method must succeed
	// TODO; hash on fail
	// TODO: must match, result of private key from deterministic gen?
	// deterministic gen will always return a valid private key
	if len(seckey) != 32 {
		log.Panic("PubkeyFromSeckey: invalid length")
	}

	pubkey := pubkeyFromSeckey(seckey)
	if pubkey == nil {
		log.Panic("ERRROR: impossible, pubkey generation failed")
		return nil
	}
	if ret := secp.PubkeyIsValid(pubkey); ret != 1 {
		log.Panicf("ERROR: Pubkey invalid, ret=%d", ret)
		return nil
	}

	return pubkey
}

// UncompressPubkey uncompresses pubkey
func UncompressPubkey(pubkey []byte) []byte {
	if VerifyPubkey(pubkey) != 1 {
		log.Panic("cannot uncompress invalid pubkey")
		return nil
	}

	var pubXY secp.XY
	if err := pubXY.ParsePubkey(pubkey); err != nil {
		log.Panicf("ERROR: impossible, pubkey parse fail: %v", err)
	}

	var pubkey2 = pubXY.BytesUncompressed()
	if pubkey2 == nil {
		log.Panic("ERROR: pubkey, uncompression fail")
		return nil
	}

	return pubkey2
}

// UncompressedPubkeyFromSeckey returns nil on error
// should only need pubkey, not private key
func UncompressedPubkeyFromSeckey(seckey []byte) []byte {
	if len(seckey) != 32 {
		log.Panic("UncompressedPubkeyFromSeckey: invalid length")
	}

	pubkey := PubkeyFromSeckey(seckey)
	if pubkey == nil {
		log.Panic("Generating seckey from pubkey, failed")
		return nil
	}

	if VerifyPubkey(pubkey) != 1 {
		log.Panic("ERROR: impossible, Pubkey generation succeeded but pubkey validation failed")
	}

	var uncompressedPubkey = UncompressPubkey(pubkey)
	if uncompressedPubkey == nil {
		log.Panic("decompression failed")
		return nil
	}

	return uncompressedPubkey
}

// deterministicKeyPairIteratorStep generates deterministic keypair with weak SHA256 hash of seed.
// internal use only
func deterministicKeyPairIteratorStep(seed []byte) ([]byte, []byte) {
	if len(seed) != 32 {
		log.Panic("ERROR: deterministicKeyPairIteratorStep: seed must be 32 bytes")
	}

	const seckeyLen = 32
	seckey := make([]byte, seckeyLen)

new_seckey:
	seed = SumSHA256(seed)
	copy(seckey, seed)

	if secp.SeckeyIsValid(seckey) != 1 {
		if DebugPrint {
			log.Printf("deterministicKeyPairIteratorStep, secp.SeckeyIsValid fail")
		}
		goto new_seckey //regen
	}

	pubkey := secp.GeneratePublicKey(seckey)
	if pubkey == nil {
		log.Panic("ERROR: deterministicKeyPairIteratorStep: GeneratePublicKey failed, impossible, secp.BaseMultiply always returns true")
		goto new_seckey
	}

	if len(pubkey) != 33 {
		log.Panic("ERROR: deterministicKeyPairIteratorStep: impossible, pubkey length wrong")
	}

	if ret := secp.PubkeyIsValid(pubkey); ret != 1 {
		log.Panicf("ERROR: deterministicKeyPairIteratorStep: PubkeyIsValid failed, ret=%d", ret)
	}

	if ret := VerifyPubkey(pubkey); ret != 1 {
		log.Printf("seckey= %s", hex.EncodeToString(seckey))
		log.Printf("pubkey= %s", hex.EncodeToString(pubkey))

		log.Panicf("ERROR: deterministicKeyPairIteratorStep: VerifyPubkey failed, ret=%d", ret)
		goto new_seckey
	}

	return pubkey, seckey
}

// Secp256k1Hash double SHA256, salted with ECDH operation in curve
func Secp256k1Hash(seed []byte) []byte { //nolint:golint
	hash := SumSHA256(seed)
	_, seckey := deterministicKeyPairIteratorStep(hash) // seckey1 is usually sha256 of hash
	pubkeySeed := SumSHA256(hash)
	pubkey, _ := deterministicKeyPairIteratorStep(pubkeySeed) // SumSHA256(hash) usually equals seckey
	ecdh := ECDH(pubkey, seckey)                              // raise pubkey to power of seckey in curve
	out := SumSHA256(append(hash, ecdh...))                   // append signature to sha256(seed) and hash
	return out
}

// GenerateDeterministicKeyPair generate a single secure key
func GenerateDeterministicKeyPair(seed []byte) ([]byte, []byte) {
	_, pubkey, seckey := DeterministicKeyPairIterator(seed)
	return pubkey, seckey
}

// DeterministicKeyPairIterator iterator for deterministic keypair generation.
// Returns SHA256, PubKey, SecKey as bytes
// Feeds SHA256 back into function to generate sequence of seckeys
// If private key is disclosed, should not be able to compute future or past keys in sequence
func DeterministicKeyPairIterator(seedIn []byte) ([]byte, []byte, []byte) {
	seed1 := Secp256k1Hash(seedIn) // make it difficult to derive future seckeys from previous seckeys
	seed2 := SumSHA256(append(seedIn, seed1...))
	pubkey, seckey := deterministicKeyPairIteratorStep(seed2) // this is our seckey
	return seed1, pubkey, seckey
}

func newRandomNonceNumber() secp.Number {
	nonce := RandByte(32)
	var n secp.Number
	n.SetBytes(nonce)
	return n
}

// newSigningNonce creates a nonce for signing. This is the `k` parameter in
// ECDSA signing. `k` must be 0 < k < n, where `n` is the order of the curve
func newSigningNonce() secp.Number {
	nonce := newRandomNonceNumber()
	for nonce.Sign() == 0 || nonce.Cmp(&secp.TheCurve.Order.Int) >= 0 {
		nonce = newRandomNonceNumber()
	}
	return nonce
}

// Sign sign hash, returns a compact recoverable signature
func Sign(msg []byte, seckey []byte) []byte {
	if len(seckey) != 32 {
		log.Panic("Sign, Invalid seckey length")
	}
	if secp.SeckeyIsValid(seckey) != 1 {
		log.Panic("Attempting to sign with invalid seckey")
	}
	if len(msg) == 0 {
		log.Panic("Sign, message nil")
	}
	if len(msg) != 32 {
		log.Panic("Sign, message must be 32 bytes")
	}

	nonce := newSigningNonce()
	sig := make([]byte, 65)
	var recid int // recovery byte, used to recover pubkey from sig

	var cSig secp.Signature

	var seckey1 secp.Number
	var msg1 secp.Number

	seckey1.SetBytes(seckey)
	msg1.SetBytes(msg)

	if msg1.Sign() == 0 {
		log.Panic("Sign: message is 0")
	}

	ret := cSig.Sign(&seckey1, &msg1, &nonce, &recid)

	if ret != 1 {
		log.Panic("Secp25k1-go, Sign, signature operation failed")
	}

	sigBytes := cSig.Bytes()
	for i := 0; i < 64; i++ {
		sig[i] = sigBytes[i]
	}
	if len(sigBytes) != 64 {
		log.Panicf("Invalid signature byte count: %d", len(sigBytes))
	}
	sig[64] = byte(recid)

	if recid > 4 {
		log.Panic("invalid recovery id")
	}

	return sig
}

// VerifySeckey verifies a secret key
// Returns 1 on success
func VerifySeckey(seckey []byte) int {
	if len(seckey) != 32 {
		return -1
	}

	//does conversion internally if less than order of curve
	if secp.SeckeyIsValid(seckey) != 1 {
		return -2
	}

	//seckey is just 32 bit integer
	//assume all seckey are valid
	//no. must be less than order of curve
	//note: converts internally
	return 1
}

// VerifyPubkey verifies a public key
// Returns 1 on success
func VerifyPubkey(pubkey []byte) int {
	if len(pubkey) != 33 {
		return -2
	}

	if secp.PubkeyIsValid(pubkey) != 1 {
		return -1 // tests parse and validity
	}

	return 1 //valid
}

// VerifySignatureValidity verifies a signature is well formed and not malleable
// Returns 1 on success
func VerifySignatureValidity(sig []byte) int {
	//64+1
	if len(sig) != 65 {
		log.Panic("VerifySignatureValidity: sig len is not 65 bytes")
		return 0
	}
	//malleability check:
	//highest bit of 32nd byte must be 1
	//0x7f is 126 or 0b01111111
	if (sig[32] >> 7) == 1 {
		return 0 // signature is malleable
	}
	//recovery id check
	if sig[64] >= 4 {
		return 0 // recovery id invalid
	}
	return 1
}

// VerifySignature for compressed signatures, does not need pubkey
// Returns 1 on success
func VerifySignature(msg []byte, sig []byte, pubkey1 []byte) int {
	if msg == nil || len(sig) == 0 || len(pubkey1) == 0 {
		log.Panic("VerifySignature, ERROR: invalid input, empty slices")
	}
	if len(sig) != 65 {
		log.Panic("VerifySignature, invalid signature length")
	}
	if len(pubkey1) != 33 {
		log.Panic("VerifySignature, invalid pubkey length")
	}

	if len(msg) == 0 {
		return 0 // empty message
	}

	// malleability check:
	// to enforce malleability, highest bit of S must be 1
	// S starts at 32nd byte
	// 0x80 is 0b10000000 or 128 and masks highest bit
	if (sig[32] >> 7) == 1 {
		return 0 // valid signature, but fails malleability
	}

	if sig[64] >= 4 {
		return 0 // recovery byte invalid
	}

	pubkey2 := RecoverPubkey(msg, sig)
	if pubkey2 == nil {
		return 0 // pubkey could not be recovered, signature is invalid
	}

	if len(pubkey2) != 33 {
		log.Panic("recovered pubkey length invalid") // sanity check
	}

	if !bytes.Equal(pubkey1, pubkey2) {
		return 0 // pubkeys do not match
	}

	return 1 // valid signature
}

// RecoverPubkey recovers the public key from the signature
func RecoverPubkey(msg []byte, sig []byte) []byte {
	if len(sig) != 65 {
		log.Panic("sig length must be 65 bytes")
	}

	var recid = int(sig[64])

	pubkey, ret := secp.RecoverPublicKey(sig[0:64], msg, recid)

	if ret != 1 {
		if DebugPrint {
			log.Printf("RecoverPubkey: code %d", ret)
		}
		return nil
	}

	if pubkey == nil {
		log.Panic("ERROR: impossible, pubkey nil and ret == 1")
	}
	if len(pubkey) != 33 {
		log.Panic("pubkey length wrong")
	}

	return pubkey
}

// ECDH raise a pubkey to the power of a seckey
func ECDH(pub, sec []byte) []byte {
	if len(sec) != 32 {
		log.Panic("secret key must be 32 bytes")
	}

	if len(pub) != 33 {
		log.Panic("public key must be 33 bytes")
	}

	if VerifySeckey(sec) != 1 {
		if DebugPrint {
			log.Printf("Invalid Seckey")
		}
		return nil
	}

	if ret := VerifyPubkey(pub); ret != 1 {
		if DebugPrint {
			log.Printf("Invalid Pubkey, %d", ret)
		}
		return nil
	}

	pubkeyOut := secp.Multiply(pub, sec)
	if pubkeyOut == nil {
		return nil
	}
	if len(pubkeyOut) != 33 {
		log.Panic("ERROR: impossible, invalid pubkey length")
	}
	return pubkeyOut
}
