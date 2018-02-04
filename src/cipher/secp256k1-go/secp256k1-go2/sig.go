package secp256k1go

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

// Signature represents the signature
type Signature struct {
	R, S Number
}

// Print prints the signature
func (sig *Signature) Print(lab string) {
	fmt.Println(lab+".R:", hex.EncodeToString(sig.R.Bytes()))
	fmt.Println(lab+".S:", hex.EncodeToString(sig.S.Bytes()))
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

// Recover TODO: return type, or nil on failure
func (sig *Signature) Recover(pubkey *XY, m *Number, recid int) (ret bool) {
	var rx, rn, u1, u2 Number
	var fx Field
	var X XY
	var xj, qj XYZ

	rx.Set(&sig.R.Int)
	if (recid & 2) != 0 {
		rx.Add(&rx.Int, &TheCurve.Order.Int)
		if rx.Cmp(&TheCurve.p.Int) >= 0 {
			return false
		}
	}

	fx.SetB32(rx.getBin(32))

	X.SetXO(&fx, (recid&1) != 0)
	if !X.IsValid() {
		return false
	}

	xj.SetXY(&X)
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

/*
//uncompressed Signature Parsing in DER
func (r *Signature) ParseBytes(sig []byte) int {
	if sig[0] != 0x30 || len(sig) < 5 {
		return -1
	}

	lenr := int(sig[3])
	if lenr == 0 || 5+lenr >= len(sig) || sig[lenr+4] != 0x02 {
		return -1
	}

	lens := int(sig[lenr+5])
	if lens == 0 || int(sig[1]) != lenr+lens+4 || lenr+lens+6 > len(sig) || sig[2] != 0x02 {
		return -1
	}

	r.R.SetBytes(sig[4 : 4+lenr])
	r.S.SetBytes(sig[6+lenr : 6+lenr+lens])
	return 6 + lenr + lens
}
*/

/*
//uncompressed Signature parsing in DER
func (sig *Signature) Bytes() []byte {
	r := sig.R.Bytes()
	if r[0] >= 0x80 {
		r = append([]byte{0}, r...)
	}
	s := sig.S.Bytes()
	if s[0] >= 0x80 {
		s = append([]byte{0}, s...)
	}
	res := new(bytes.Buffer)
	res.WriteByte(0x30)
	res.WriteByte(byte(4 + len(r) + len(s)))
	res.WriteByte(0x02)
	res.WriteByte(byte(len(r)))
	res.Write(r)
	res.WriteByte(0x02)
	res.WriteByte(byte(len(s)))
	res.Write(s)
	return res.Bytes()
}
*/

// ParseBytes compressed signature parsing
func (sig *Signature) ParseBytes(v []byte) {
	if len(v) != 64 {
		log.Panic()
	}
	sig.R.SetBytes(v[0:32])
	sig.S.SetBytes(v[32:64])
}

//secp256k1_num_get_bin(sig64, 32, &sig.r);
//secp256k1_num_get_bin(sig64 + 32, 32, &sig.s);

// Bytes compressed signature parsing
func (sig *Signature) Bytes() []byte {
	r := sig.R.Bytes() //endianess
	s := sig.S.Bytes() //endianess

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
	res.Write(r)
	res.Write(s)

	//test
	if true {
		ret := res.Bytes()
		var sig2 Signature
		sig2.ParseBytes(ret)
		if bytes.Equal(sig.R.Bytes(), sig2.R.Bytes()) == false {
			log.Panic("serialization failed 1")
		}
		if bytes.Equal(sig.S.Bytes(), sig2.S.Bytes()) == false {
			log.Panic("serialization failed 2")
		}
	}

	if len(res.Bytes()) != 64 {
		log.Panic()
	}
	return res.Bytes()
}
