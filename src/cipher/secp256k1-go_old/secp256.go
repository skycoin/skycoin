package secp256k1

/*
#cgo CFLAGS: -std=gnu99 -Wno-error
#cgo LDFLAGS: -lgmp
#cgo CFLAGS: -Wno-error
#define USE_FIELD_10X26
#define USE_NUM_GMP
#define USE_FIELD_INV_BUILTIN
#include "./secp256k1/src/secp256k1.c"
*/
import "C"

//#cgo pkg-config: gmp
//#cgo pkg-config: secp256

//for osx 'xcode-select --install'
//in 10.9 (mavericks) it fails. See README.md

import (
	"unsafe"
	//"fmt"
	//"errors"
	"bytes"
	"log"
)

//#define USE_FIELD_5X64

/*
   Todo:
   > Centralize key management in module
   > add pubkey/private key struct
   > Dont let keys leave module; address keys as ints

   > store private keys in buffer and shuffle (deters persistance on swap disc)
   > Byte permutation (changing)
   > xor with chaning random block (to deter scanning memory for 0x63) (stream cipher?)

   On Disk
   > Store keys in wallets
   > use slow key derivation function for wallet encryption key (2 seconds)
*/

func init() {
	C.secp256k1_start() //takes 10ms to 100ms
}

func Stop() {
	C.secp256k1_stop()
}

/*
int secp256k1_ecdsa_pubkey_create(
    unsigned char *pubkey, int *pubkeylen,
    const unsigned char *seckey, int compressed);
*/

/** Compute the public key for a secret key.
 *  In:     compressed: whether the computed public key should be compressed
 *          seckey:     pointer to a 32-byte private key.
 *  Out:    pubkey:     pointer to a 33-byte (if compressed) or 65-byte (if uncompressed)
 *                      area to store the public key.
 *          pubkeylen:  pointer to int that will be updated to contains the pubkey's
 *                      length.
 *  Returns: 1: secret was valid, public key stores
 *           0: secret was invalid, try again.
 */

//pubkey, seckey

func GenerateKeyPair() ([]byte, []byte) {

	pubkey_len := C.int(33)
	const seckey_len = 32

	var pubkey []byte = make([]byte, pubkey_len)
	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))

	var ret C.int

new_seckey:
	var seckey []byte = RandByte(seckey_len)
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))

	ret = C.secp256k1_ecdsa_seckey_verify(seckey_ptr)

	if ret != 1 {
		goto new_seckey
	}

	//for C.secp256k1_ecdsa_seckey_verify(seckey_ptr) != 1 {
	//	seckey = RandByte(seckey_len)
	//	seckey_ptr = (*C.uchar)(unsafe.Pointer(&seckey[0]))
	//}

	ret = C.secp256k1_ecdsa_pubkey_create(
		pubkey_ptr, &pubkey_len,
		seckey_ptr, 1)

	if ret != 1 {
		goto new_seckey
	}

	return pubkey, seckey
}

//returns nil on error
func PubkeyFromSeckey(SecKey []byte) []byte {
	if len(SecKey) != 32 {
		log.Panic("PubkeyFromSeckey: invalid length")
	}

	pubkey_len := C.int(33)
	const seckey_len = 32

	var pubkey []byte = make([]byte, pubkey_len)
	var seckey []byte = make([]byte, seckey_len)
	copy(seckey, SecKey)

	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))

	ret := C.secp256k1_ecdsa_pubkey_create(
		pubkey_ptr, &pubkey_len,
		seckey_ptr, 1)

	if ret != 1 {
		return nil
	}
	return pubkey
}

//returns nil on error
func UncompressedPubkeyFromSeckey(SecKey []byte) []byte {
	if len(SecKey) != 32 {
		log.Panic("PubkeyFromSeckey: invalid length")
	}

	pubkey_len := C.int(65)
	const seckey_len = 32

	var pubkey []byte = make([]byte, pubkey_len)
	var seckey []byte = make([]byte, seckey_len)
	copy(seckey, SecKey)

	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))

	ret := C.secp256k1_ecdsa_pubkey_create(
		pubkey_ptr, &pubkey_len,
		seckey_ptr, 0)

	if ret != 1 {
		return nil
	}
	return pubkey
}

//generates deterministic keypair with weak SHA256 hash of seed
//internal use only
//be extremely careful with golang slice semantics
func generateDeterministicKeyPair(seed []byte) ([]byte, []byte) {
	if seed == nil {
		log.Panic()
	}

	pubkey_len := C.int(33)
	const seckey_len = 32

	var pubkey []byte = make([]byte, pubkey_len)
	var seckey []byte = make([]byte, seckey_len)

	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))

new_seckey:
	seed = SumSHA256(seed[0:32])
	copy(seckey[0:32], seed[0:32])
	if C.secp256k1_ecdsa_seckey_verify(seckey_ptr) != 1 {
		goto new_seckey //rehash seckey until it succeeds, almost never happens
	}

	ret := C.secp256k1_ecdsa_pubkey_create(
		pubkey_ptr, &pubkey_len,
		seckey_ptr, 1)

	if ret != 1 {
		log.Panic("secp256k1-g0, generateDeterministicKeyPair, pubkey generation failing for valid seckey")
	}

	return pubkey, seckey
}

//this is a GPU and ASIC resistant hash function that combines SHA256 with operations on
// elliptic curve through  slow secp256k1 signature operations. designed to protect
// brainwallet seeds against GPU brute forcing

/*
func Secp256k1Hash(hash []byte) []byte {
	hash = SumSHA256(hash)                            //sha256
	_, seckey := generateDeterministicKeyPair(hash)   //generate key
	sig := SignDeterministic(hash, seckey, hash)      //sign with key
	return SumSHA256(append(SumSHA256(hash), sig...)) //append signature to sha256(seed) and hash
}
*/

//double SHA256, salted with ECDH operation in curve
func Secp256k1Hash(hash []byte) []byte {
	hash = SumSHA256(hash)
	_, seckey := generateDeterministicKeyPair(hash)            //seckey1 is usually sha256 of hash
	pubkey, _ := generateDeterministicKeyPair(SumSHA256(hash)) //SumSHA256(hash) equals seckey usually
	ecdh := ECDH(pubkey, seckey)                               //raise pubkey to power of seckey in curve
	return SumSHA256(append(hash, ecdh...))                    //append signature to sha256(seed) and hash
}

//generate a single secure key
func GenerateDeterministicKeyPair(seed []byte) ([]byte, []byte) {
	_, pubkey, seckey := DeterministicKeyPairIterator(seed)
	return pubkey, seckey
}

//Iterator for deterministic keypair generation. Returns SHA256, Pubkey, Seckey
//Feed SHA256 back into function to generate sequence of seckeys
func DeterministicKeyPairIterator(seed_in []byte) ([]byte, []byte, []byte) {
	seed1 := Secp256k1Hash(seed_in) //make it difficult to derive future seckeys from previous seckeys
	seed2 := SumSHA256(append(seed_in, seed1...))
	pubkey, seckey := generateDeterministicKeyPair(seed2) //this is our seckey
	return seed1, pubkey, seckey
}

/*
*  Create a compact ECDSA signature (64 byte + recovery id).
*  Returns: 1: signature created
*           0: nonce invalid, try another one
*  In:      msg:    the message being signed
*           msglen: the length of the message being signed
*           seckey: pointer to a 32-byte secret key (assumed to be valid)
*           nonce:  pointer to a 32-byte nonce (generated with a cryptographic PRNG)
*  Out:     sig:    pointer to a 64-byte array where the signature will be placed.
*           recid:  pointer to an int, which will be updated to contain the recovery id.
 */

/*
int secp256k1_ecdsa_sign_compact(const unsigned char *msg, int msglen,
                                 unsigned char *sig64,
                                 const unsigned char *seckey,
                                 const unsigned char *nonce,
                                 int *recid);
*/

//Rename SignHash
func Sign(msg []byte, seckey []byte) []byte {

	if len(seckey) != 32 {
		log.Panic("Sign, Invalid seckey length")
	}
	if msg == nil {
		log.Panic("Sign, message nil")
	}
	var nonce []byte = RandByte(32) //going to get bitcoins stolen!

	var sig []byte = make([]byte, 65)
	var recid C.int

	var msg_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&msg[0]))
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))
	var nonce_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&nonce[0]))
	var sig_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&sig[0]))

	if C.secp256k1_ecdsa_seckey_verify(seckey_ptr) != C.int(1) {
		log.Panic() //invalid seckey
	}

	ret := C.secp256k1_ecdsa_sign_compact(
		msg_ptr, C.int(len(msg)),
		sig_ptr,
		seckey_ptr,
		nonce_ptr,
		&recid)

	sig[64] = byte(int(recid))

	if int(recid) > 4 {
		log.Panic()
	}

	if ret != 1 {
		return Sign(msg, seckey) //nonce invalid,retry
	}

	return sig
}

//generate signature in repeatable way
func SignDeterministic(msg []byte, seckey []byte, nonce_seed []byte) []byte {
	nonce := SumSHA256(nonce_seed) //deterministicly generate nonce

	var sig []byte = make([]byte, 65)
	var recid C.int

	var msg_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&msg[0]))
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))
	var nonce_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&nonce[0]))
	var sig_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&sig[0]))

	if C.secp256k1_ecdsa_seckey_verify(seckey_ptr) != C.int(1) {
		log.Panic("Invalid secret key")
	}

	ret := C.secp256k1_ecdsa_sign_compact(
		msg_ptr, C.int(len(msg)),
		sig_ptr,
		seckey_ptr,
		nonce_ptr,
		&recid)

	sig[64] = byte(int(recid))

	if int(recid) > 4 {
		log.Panic()
	}

	if ret != 1 {
		return SignDeterministic(msg, seckey, nonce_seed) //nonce invalid,retry
	}

	return sig
}

/*
* Verify an ECDSA secret key.
*  Returns: 1: secret key is valid
*           0: secret key is invalid
*  In:      seckey: pointer to a 32-byte secret key
 */

//Rename ChkSeckeyValidity
func VerifySeckey(seckey []byte) int {
	if len(seckey) != 32 {
		return 0
	}
	var seckey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&seckey[0]))
	ret := C.secp256k1_ecdsa_seckey_verify(seckey_ptr)
	return int(ret)
}

/*
* Validate a public key.
*  Returns: 1: valid public key
*           0: invalid public key
 */

//Rename ChkPubkeyValidity
func VerifyPubkey(pubkey []byte) int {
	if len(pubkey) != 33 {
		return 0
	}
	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))
	ret := C.secp256k1_ecdsa_pubkey_verify(pubkey_ptr, 33)
	return int(ret)
}

//Rename ChkSignatureValidity
func VerifySignatureValidity(sig []byte) int {
	//64+1
	if len(sig) != 65 {
		return 0
	}
	//malleability check:
	//highest bit of 32nd byte must be 1
	//0x7f us 126 or 0b01111111
	if (sig[32] >> 7) == 1 {
		return 0
	}
	//recovery id check
	if sig[64] >= 4 {
		return 0
	}
	return 1
}

//for compressed signatures, does not need pubkey
//Rename SignatureChk
func VerifySignature(msg []byte, sig []byte, pubkey1 []byte) int {
	if msg == nil || sig == nil || pubkey1 == nil {
		log.Panic("VerifySignature, ERROR: invalid input, nils")
	}
	if len(sig) != 65 {
		log.Panic("VerifySignature, invalid signature length")
	}
	if len(pubkey1) != 33 {
		log.Panic("VerifySignature, invalid pubkey length")
	}

	//malleability check:
	//to enforce malleability, highest bit of S must be 1
	//S starts at 32nd byte
	//0x80 is 0b10000000 or 128 and masks highest bit
	if (sig[32] >> 7) == 1 {
		return 0 //valid signature, but fails malleability
	}

	if sig[64] >= 4 {
		return 0 //recover byte invalid
	}

	pubkey2 := RecoverPubkey(msg, sig) //if pubkey recovered, signature valid

	if pubkey2 == nil {
		return 0
	}

	if len(pubkey2) != 33 {
		log.Panic("recovered pubkey length invalid")
	}

	if bytes.Equal(pubkey1, pubkey2) != true {
		return 0 //pubkeys do not match
	}

	return 1 //valid signature
}

//SignatureErrorString returns error string for signature failure
func SignatureErrorString(msg []byte, sig []byte, pubkey1 []byte) string {

	if msg == nil || len(sig) != 65 || len(pubkey1) != 33 {
		log.Panic()
	}

	if (sig[32] >> 7) == 1 {
		return "signature fails malleability requirement"
	}

	if sig[64] >= 4 {
		return "signature recovery byte is invalid, must be 0 to 3"
	}

	pubkey2 := RecoverPubkey(msg, sig) //if pubkey recovered, signature valid
	if pubkey2 == nil {
		return "pubkey from signature failed"
	}

	if bytes.Equal(pubkey1, pubkey2) == false {
		return "input pubkey and recovered pubkey do not match"
	}

	return "No Error!"
}

/*
int secp256k1_ecdsa_recover_compact(const unsigned char *msg, int msglen,
                                    const unsigned char *sig64,
                                    unsigned char *pubkey, int *pubkeylen,
                                    int compressed, int recid);
*/

/*
 * Recover an ECDSA public key from a compact signature.
 *  Returns: 1: public key succesfully recovered (which guarantees a correct signature).
 *           0: otherwise.
 *  In:      msg:        the message assumed to be signed
 *           msglen:     the length of the message
 *           compressed: whether to recover a compressed or uncompressed pubkey
 *           recid:      the recovery id (as returned by ecdsa_sign_compact)
 *  Out:     pubkey:     pointer to a 33 or 65 byte array to put the pubkey.
 *           pubkeylen:  pointer to an int that will contain the pubkey length.
 */

//recovers the public key from the signature
//recovery of pubkey means correct signature
func RecoverPubkey(msg []byte, sig []byte) []byte {
	if len(sig) != 65 {
		log.Panic()
	}

	var pubkey []byte = make([]byte, 33)

	var msg_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&msg[0]))
	var sig_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&sig[0]))
	var pubkey_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pubkey[0]))

	var pubkeylen C.int

	ret := C.secp256k1_ecdsa_recover_compact(
		msg_ptr, C.int(len(msg)),
		sig_ptr,
		pubkey_ptr, &pubkeylen,
		C.int(1), C.int(sig[64]),
	)

	if ret == 0 || int(pubkeylen) != 33 {
		return nil
	}

	return pubkey
}

//int secp256k1_ecdsa_privkey_tweak_mul(unsigned char *seckey, const unsigned char *tweak);

//int secp256k1_ecdsa_pubkey_tweak_mul(unsigned char *pubkey, int pubkeylen, const unsigned char *tweak);

//raise a pubkey to the power of a seckey
func ECDH(pub []byte, sec []byte) []byte {
	if len(sec) != 32 {
		log.Panic()
	}

	if len(pub) != 33 {
		log.Panic()
	}

	if VerifyPubkey(pub) != 1 {
		log.Printf("Invalid Pubkey")
		return nil
	}

	if VerifySeckey(sec) != 1 {
		log.Printf("Invalid Seckey")
	}

	var pub2 []byte = make([]byte, 33)
	copy(pub2[0:33], pub[0:33])

	var pub_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&pub2[0]))
	var sec_ptr *C.uchar = (*C.uchar)(unsafe.Pointer(&sec[0]))

	ret := C.secp256k1_ecdsa_pubkey_tweak_mul(
		pub_ptr, C.int(len(pub2)),
		sec_ptr,
	)

	if ret != 1 {
		return nil
	}

	return pub2
}
