package secp256k1go

import (
	"bytes"
	"log"
)

// Signature represents the signature
type Signature struct {
	R, S Number
}

// Verify verify the signature
func (sig *Signature) Verify(pubkey *XY, message *Number) (ret bool) {
	var r2 Number
	ret = sig.recompute(&r2, pubkey, message) && sig.R.Cmp(&r2.Int) == 0
	return
}

func (sig *Signature) recompute(r2 *Number, pubkey *XY, message *Number) (ret bool) {
	var sn, u1, u2 Number

	sn.modInv(&sig.S, &TheCurve.Order)
	u1.modMul(&sn, message, &TheCurve.Order)
	u2.modMul(&sn, &sig.R, &TheCurve.Order)

	var pr, pubkeyj XYZ
	pubkeyj.SetXY(pubkey)

	pubkeyj.ECmult(&pr, &u2, &u1)
	if !pr.IsInfinity() {
		var xr Field
		pr.getX(&xr)
		xr.Normalize()
		var xrb [32]byte
		xr.GetB32(xrb[:])
		r2.SetBytes(xrb[:])
		r2.Mod(&r2.Int, &TheCurve.Order.Int)
		ret = true
	}

	return
}

// Recover recovers a pubkey XY point given the message that was signed to create
// this signature.
// TODO: return type, or nil on failure
func (sig *Signature) Recover(pubkey *XY, m *Number, recid int) (ret bool) {
	var rx, rn, u1, u2 Number
	var fx Field
	var x XY
	var xj, qj XYZ

	rx.Set(&sig.R.Int)
	if (recid & 2) != 0 {
		rx.Add(&rx.Int, &TheCurve.Order.Int)
		if rx.Cmp(&TheCurve.p.Int) >= 0 {
			return false
		}
	}

	fx.SetB32(LeftPadBytes(rx.Bytes(), 32))

	x.SetXO(&fx, (recid&1) != 0)
	if !x.IsValid() {
		return false
	}

	xj.SetXY(&x)
	rn.modInv(&sig.R, &TheCurve.Order)
	u1.modMul(&rn, m, &TheCurve.Order)
	u1.Sub(&TheCurve.Order.Int, &u1.Int)
	u2.modMul(&rn, &sig.S, &TheCurve.Order)
	xj.ECmult(&qj, &u2, &u1)
	pubkey.SetXYZ(&qj)

	return true
}

// Sign signs the signature
func (sig *Signature) Sign(seckey, message, nonce *Number, recid *int) int {
	var r XY
	var rp XYZ
	var n Number
	var b [32]byte

	ECmultGen(&rp, nonce)
	r.SetXYZ(&rp)
	r.X.Normalize()
	r.Y.Normalize()
	r.X.GetB32(b[:])
	sig.R.SetBytes(b[:])
	if recid != nil {
		*recid = 0
		if sig.R.Cmp(&TheCurve.Order.Int) >= 0 {
			*recid |= 2
		}
		if r.Y.IsOdd() {
			*recid |= 1
		}
	}
	sig.R.mod(&TheCurve.Order)
	n.modMul(&sig.R, seckey, &TheCurve.Order)
	n.Add(&n.Int, &message.Int)
	n.mod(&TheCurve.Order)
	sig.S.modInv(nonce, &TheCurve.Order)
	sig.S.modMul(&sig.S, &n, &TheCurve.Order)
	if sig.S.Sign() == 0 {
		return 0
	}
	if sig.S.IsOdd() {
		sig.S.Sub(&TheCurve.Order.Int, &sig.S.Int)
		if recid != nil {
			*recid ^= 1
		}
	}

	if forceLowS && sig.S.Cmp(&TheCurve.halfOrder.Int) == 1 {
		sig.S.Sub(&TheCurve.Order.Int, &sig.S.Int)
		if recid != nil {
			*recid ^= 1
		}
	}

	return 1
}

// ParseBytes parses a serialized R||S pair to Signature.
// R and S should be in big-endian encoding.
func (sig *Signature) ParseBytes(v []byte) {
	if len(v) != 64 {
		log.Panic("Signature.ParseBytes requires 64 bytes")
	}
	sig.R.SetBytes(v[0:32])
	sig.S.SetBytes(v[32:64])
}

//secp256k1_num_get_bin(sig64, 32, &sig.r);
//secp256k1_num_get_bin(sig64 + 32, 32, &sig.s);

// Bytes serializes compressed signatures as bytes.
// The serialization format is R||S. R and S bytes are big-endian.
// R and S are left-padded with the NUL byte to ensure their length is 32 bytes.
func (sig *Signature) Bytes() []byte {
	r := sig.R.Bytes() // big-endian
	s := sig.S.Bytes() // big-endian

	for len(r) < 32 {
		r = append([]byte{0}, r...)
	}
	for len(s) < 32 {
		s = append([]byte{0}, s...)
	}

	if len(r) != 32 || len(s) != 32 {
		log.Panicf("signature size invalid: %d, %d", len(r), len(s))
	}

	res := new(bytes.Buffer)
	if _, err := res.Write(r); err != nil {
		panic(err)
	}
	if _, err := res.Write(s); err != nil {
		panic(err)
	}

	//test
	if true {
		ret := res.Bytes()
		var sig2 Signature
		sig2.ParseBytes(ret)
		if !bytes.Equal(sig.R.Bytes(), sig2.R.Bytes()) {
			log.Panic("serialization failed 1")
		}
		if !bytes.Equal(sig.S.Bytes(), sig2.S.Bytes()) {
			log.Panic("serialization failed 2")
		}
	}

	if len(res.Bytes()) != 64 {
		log.Panic()
	}
	return res.Bytes()
}
