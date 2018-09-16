package secp256k1go

import (
	//"encoding/hex"
	"bytes"
	"log"
)

// func ecdsaVerify(pubkey, sig, msg []byte) int {
// 	var m Number
// 	var s Signature
// 	m.SetBytes(msg)

// 	log.Println("pubkey len is", len(pubkey))

// 	var q XY
// 	if !q.ParsePubkey(pubkey) {
// 		return -1
// 	}

// 	//if s.ParseBytes(sig) < 0 {
// 	//	return -2
// 	//}
// 	if len(pubkey) != 32 {
// 		return -2
// 	}
// 	if len(sig) != 64 {
// 		return -3
// 	}

// 	if !s.Verify(&q, &m) {
// 		return 0
// 	}
// 	return 1
// }

// // Verify verifies ecdsa
// func Verify(k, s, m []byte) bool {
// 	return ecdsaVerify(k, s, m) == 1
// }

// DecompressPoint decompresses point
func DecompressPoint(X []byte, off bool, Y []byte) {
	var rx, ry, c, x2, x3 Field
	rx.SetB32(X)
	rx.Sqr(&x2)
	rx.Mul(&x3, &x2)
	c.SetInt(7)
	c.SetAdd(&x3)
	c.Sqrt(&ry)
	ry.Normalize()
	if ry.IsOdd() != off {
		ry.Negate(&ry, 1)
	}
	ry.Normalize()
	ry.GetB32(Y)
}

//TODO: change signature to []byte type
/*
func RecoverPublicKey2(sig Signature, h []byte, recid int, pubkey *XY) int {
	//var sig Signature
	var msg Number

	if sig.R.Sign() <= 0 || sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
		if sig.R.Sign() == 0 {
			return -10
		}
		if sig.R.Sign() <= 0 {
			return -11
		}
		if sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
			return -12
		}
		return -1
	}
	if sig.S.Sign() <= 0 || sig.S.Cmp(&TheCurve.Order.Int) >= 0 {
		return -2
	}

	msg.SetBytes(h)
	if !sig.Recover(pubkey, &msg, recid) {
		return -3
	}
	return 1
}
*/
//TODO: deprecate
/*
func RecoverPublicKey(r, s, h []byte, recid int, pubkey *XY) bool {
	var sig Signature
	var msg Number
	sig.R.SetBytes(r)
	if sig.R.Sign() <= 0 || sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
		return false
	}
	sig.S.SetBytes(s)
	if sig.S.Sign() <= 0 || sig.S.Cmp(&TheCurve.Order.Int) >= 0 {
		return false
	}
	msg.SetBytes(h)
	if !sig.Recover(pubkey, &msg, recid) {
		return false
	}
	return true
}
*/

// RecoverPublicKey nil on error
// returns error code
func RecoverPublicKey(sigByte []byte, h []byte, recid int) ([]byte, int) {

	var pubkey XY

	if len(sigByte) != 64 {
		log.Panic("must pass in 64 byte pubkey")
	}

	var sig Signature
	sig.ParseBytes(sigByte[0:64])

	//var sig Signature
	var msg Number

	if sig.R.Sign() <= 0 || sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
		if sig.R.Sign() == 0 {
			return nil, -1
		}
		if sig.R.Sign() <= 0 {
			return nil, -2
		}
		if sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
			return nil, -3
		}
		return nil, -4
	}
	if sig.S.Sign() <= 0 || sig.S.Cmp(&TheCurve.Order.Int) >= 0 {
		return nil, -5
	}

	msg.SetBytes(h)
	if !sig.Recover(&pubkey, &msg, recid) {
		return nil, -6
	}

	return pubkey.Bytes(), 1
}

// Multiply standard EC multiplacation k(xy)
// xy - is the standarized public key format (33 or 65 bytes long)
// out - should be the buffer for 33 bytes (1st byte will be set to either 02 or 03)
// TODO: change out to return type
func Multiply(xy, k []byte) []byte {
	var pk XY
	var xyz XYZ
	var na, nzero Number
	if !pk.ParsePubkey(xy) {
		return nil
	}
	xyz.SetXY(&pk)
	na.SetBytes(k)
	xyz.ECmult(&xyz, &na, &nzero)
	pk.SetXYZ(&xyz)

	if !pk.IsValid() {
		log.Panic()
	}
	return pk.GetPublicKey()
}

// Multiply k by G
// returns public key
// return nil on error, but never returns nil
// 33 bytes out

/*
func BaseMultiply2(k []byte) []byte {
	var r XYZ
	var n Number
	var pk XY
	n.SetBytes(k)
	ECmultGen(&r, &n)
	pk.SetXYZ(&r)
	if !pk.IsValid() {
		log.Panic()
	}

	return pk.GetPublicKey()
}
*/

//test assumptions
func _pubkeyTest(pk XY) {

	if !pk.IsValid() {
		log.Panic("IMPOSSIBLE3: pubkey invalid")
	}
	var pk2 XY
	retb := pk2.ParsePubkey(pk.Bytes())
	if !retb {
		log.Panic("IMPOSSIBLE2: parse failed")
	}
	if !pk2.IsValid() {
		log.Panic("IMPOSSIBLE3: parse failed non valid key")
	}
	if PubkeyIsValid(pk2.Bytes()) != 1 {
		log.Panic("IMPOSSIBLE4: pubkey failed")
	}
}

// BaseMultiply base multiply
func BaseMultiply(k []byte) []byte {
	var r XYZ
	var n Number
	var pk XY
	n.SetBytes(k)
	ECmultGen(&r, &n)
	pk.SetXYZ(&r)
	if !pk.IsValid() {
		log.Panic() //should not occur
	}

	_pubkeyTest(pk)

	return pk.Bytes()
}

// BaseMultiplyAdd out = G*k + xy
// TODO: switch to returning output as []byte
// nil on error
// 33 byte out
func BaseMultiplyAdd(xy, k []byte) []byte {
	var r XYZ
	var n Number
	var pk XY
	if !pk.ParsePubkey(xy) {
		return nil
	}
	n.SetBytes(k)
	ECmultGen(&r, &n)
	r.AddXY(&r, &pk)
	pk.SetXYZ(&r)

	_pubkeyTest(pk)
	return pk.Bytes()
}

// GeneratePublicKey returns nil on failure
//crash rather than fail
func GeneratePublicKey(k []byte) []byte {

	//log.Panic()
	if len(k) != 32 {
		log.Panic()
	}
	var r XYZ
	var n Number
	var pk XY

	//must not be zero
	//must not be negative
	//must be less than order of curve
	n.SetBytes(k)
	if n.Sign() <= 0 || n.Cmp(&TheCurve.Order.Int) >= 0 {
		log.Panic("only call for valid seckey, check that seckey is valid first")
		return nil
	}
	ECmultGen(&r, &n)
	pk.SetXYZ(&r)
	if !pk.IsValid() {
		log.Panic() //should not occur
	}
	_pubkeyTest(pk)
	return pk.Bytes()
}

// SeckeyIsValid 1 on success
// must not be zero
// must not be negative
// must be less than order of curve
func SeckeyIsValid(seckey []byte) int {
	if len(seckey) != 32 {
		log.Panic()
	}
	var n Number
	n.SetBytes(seckey)
	//must not be zero
	//must not be negative
	//must be less than order of curve
	if n.Sign() <= 0 {
		return -1
	}
	if n.Cmp(&TheCurve.Order.Int) >= 0 {
		return -2
	}
	return 1
}

// PubkeyIsValid returns 1 on success
func PubkeyIsValid(pubkey []byte) int {
	if len(pubkey) != 33 {
		log.Panic() //do not permit invalid length inputs
		return -2
	}
	var pubTest XY
	ok := pubTest.ParsePubkey(pubkey)
	if !ok {
		//log.Panic("PubkeyIsValid, ERROR: pubkey parse fail, bad pubkey from private key")
		return -1
	}
	if !bytes.Equal(pubTest.Bytes(), pubkey) {
		log.Panic("pubkey parses but serialize/deserialize roundtrip fails")
	}
	//this fails
	//if !pub_test.IsValid() {
	//	return -2
	//}
	return 1
}

/*
Note:
- choose random private key
- generate public key
- call "IsValid()" on the public key


*/
