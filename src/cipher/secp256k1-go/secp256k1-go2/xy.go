package secp256k1go

import (
	"fmt"
	"log"
)

// XY TODO...
type XY struct {
	X, Y     Field
	Infinity bool
}

// Print prints the xy
func (xy *XY) Print(lab string) {
	if xy.Infinity {
		fmt.Println(lab + " - Infinity")
		return
	}
	fmt.Println(lab+".X:", xy.X.String())
	fmt.Println(lab+".Y:", xy.Y.String())
}

//edited

/*
   if (size == 33 && (pub[0] == 0x02 || pub[0] == 0x03)) {
       secp256k1_fe_t x;
       secp256k1_fe_set_b32(&x, pub+1);
       return secp256k1_ge_set_xo(elem, &x, pub[0] == 0x03);
   } else if (size == 65 && (pub[0] == 0x04 || pub[0] == 0x06 || pub[0] == 0x07)) {
       secp256k1_fe_t x, y;
       secp256k1_fe_set_b32(&x, pub+1);
       secp256k1_fe_set_b32(&y, pub+33);
       secp256k1_ge_set_xy(elem, &x, &y);
       if ((pub[0] == 0x06 || pub[0] == 0x07) && secp256k1_fe_is_odd(&y) != (pub[0] == 0x07))
           return 0;
       return secp256k1_ge_is_valid(elem);
   }
*/
//All compact keys appear to be valid by construction, but may fail
//is valid check

// ParsePubkey WARNING: for compact signatures, will succeed unconditionally
//however, elem.IsValid will fail
func (xy *XY) ParsePubkey(pub []byte) bool {
	if len(pub) != 33 {
		log.Panic("pubkey len must be 33, len is ", len(pub)) // do not permit invalid length inputs
		return false
	}
	if len(pub) == 33 && (pub[0] == 0x02 || pub[0] == 0x03) {
		xy.X.SetB32(pub[1:33])
		xy.SetXO(&xy.X, pub[0] == 0x03)
	} else {
		return false
	}
	//THIS FAILS
	//reenable later
	//if elem.IsValid() == false {
	//	return false
	//}

	/*
		 else if len(pub) == 65 && (pub[0] == 0x04 || pub[0] == 0x06 || pub[0] == 0x07) {
			elem.X.SetB32(pub[1:33])
			elem.Y.SetB32(pub[33:65])
			if (pub[0] == 0x06 || pub[0] == 0x07) && elem.Y.IsOdd() != (pub[0] == 0x07) {
				return false
			}
		}
	*/
	return true
}

// Bytes Returns serialized key in in compressed format: "<02> <X>",
// eventually "<03> <X>"
//33 bytes
func (xy XY) Bytes() []byte {
	xy.X.Normalize() // See GitHub issue #15

	raw := make([]byte, 33)
	if xy.Y.IsOdd() {
		raw[0] = 0x03
	} else {
		raw[0] = 0x02
	}
	xy.X.GetB32(raw[1:])
	return raw
}

// BytesUncompressed returns serialized key in uncompressed format "<04> <X> <Y>"
//65 bytes
func (xy *XY) BytesUncompressed() (raw []byte) {
	xy.X.Normalize() // See GitHub issue #15
	xy.Y.Normalize() // See GitHub issue #15

	raw = make([]byte, 65)
	raw[0] = 0x04
	xy.X.GetB32(raw[1:33])
	xy.Y.GetB32(raw[33:65])
	return
}

// SetXY sets x y fields
func (xy *XY) SetXY(X, Y *Field) {
	xy.Infinity = false
	xy.X = *X
	xy.Y = *Y
}

/*
int static secp256k1_ecdsa_pubkey_parse(secp256k1_ge_t *elem, const unsigned char *pub, int size) {
    if (size == 33 && (pub[0] == 0x02 || pub[0] == 0x03)) {
        secp256k1_fe_t x;
        secp256k1_fe_set_b32(&x, pub+1);
        return secp256k1_ge_set_xo(elem, &x, pub[0] == 0x03);
    } else if (size == 65 && (pub[0] == 0x04 || pub[0] == 0x06 || pub[0] == 0x07)) {
        secp256k1_fe_t x, y;
        secp256k1_fe_set_b32(&x, pub+1);
        secp256k1_fe_set_b32(&y, pub+33);
        secp256k1_ge_set_xy(elem, &x, &y);
        if ((pub[0] == 0x06 || pub[0] == 0x07) && secp256k1_fe_is_odd(&y) != (pub[0] == 0x07))
            return 0;
        return secp256k1_ge_is_valid(elem);
    } else {
        return 0;
    }
}
*/

//    if (size == 33 && (pub[0] == 0x02 || pub[0] == 0x03)) {
//        secp256k1_fe_t x;
//        secp256k1_fe_set_b32(&x, pub+1);
//        return secp256k1_ge_set_xo(elem, &x, pub[0] == 0x03);

// IsValid checks if valid
func (xy *XY) IsValid() bool {
	if xy.Infinity {
		return false
	}
	var y2, x3, c Field
	xy.Y.Sqr(&y2)
	xy.X.Sqr(&x3)
	x3.Mul(&x3, &xy.X)
	c.SetInt(7)
	x3.SetAdd(&c)
	y2.Normalize()
	x3.Normalize()
	return y2.Equals(&x3)
}

// SetXYZ sets X Y Z fields
func (xy *XY) SetXYZ(a *XYZ) {
	var z2, z3 Field
	a.Z.InvVar(&a.Z)
	a.Z.Sqr(&z2)
	a.Z.Mul(&z3, &z2)
	a.X.Mul(&a.X, &z2)
	a.Y.Mul(&a.Y, &z3)
	a.Z.SetInt(1)
	xy.Infinity = a.Infinity
	xy.X = a.X
	xy.Y = a.Y
}

func (xy *XY) precomp(w int) (pre []XY) {
	pre = make([]XY, (1 << (uint(w) - 2)))
	pre[0] = *xy
	var X, d, tmp XYZ
	X.SetXY(xy)
	X.Double(&d)
	for i := 1; i < len(pre); i++ {
		d.AddXY(&tmp, &pre[i-1])
		pre[i].SetXYZ(&tmp)
	}
	return
}

// Neg caculates negate
func (xy *XY) Neg(r *XY) {
	r.Infinity = xy.Infinity
	r.X = xy.X
	r.Y = xy.Y
	r.Y.Normalize()
	r.Y.Negate(&r.Y, 1)
}

/*
int static secp256k1_ge_set_xo(secp256k1_ge_t *r, const secp256k1_fe_t *x, int odd) {
    r->x = *x;
    secp256k1_fe_t x2; secp256k1_fe_sqr(&x2, x);
    secp256k1_fe_t x3; secp256k1_fe_mul(&x3, x, &x2);
    r->infinity = 0;
    secp256k1_fe_t c; secp256k1_fe_set_int(&c, 7);
    secp256k1_fe_add(&c, &x3);
    if (!secp256k1_fe_sqrt(&r->y, &c))
        return 0;
    secp256k1_fe_normalize(&r->y);
    if (secp256k1_fe_is_odd(&r->y) != odd)
        secp256k1_fe_negate(&r->y, &r->y, 1);
    return 1;
}
*/

// SetXO sets
func (xy *XY) SetXO(X *Field, odd bool) {
	var c, x2, x3 Field
	xy.X = *X
	X.Sqr(&x2)
	X.Mul(&x3, &x2)
	xy.Infinity = false
	c.SetInt(7)
	c.SetAdd(&x3)
	c.Sqrt(&xy.Y) //does not return, can fail
	if xy.Y.IsOdd() != odd {
		xy.Y.Negate(&xy.Y, 1)
	}

	//r.X.Normalize() // See GitHub issue #15
	xy.Y.Normalize()
}

// AddXY adds xy
func (xy *XY) AddXY(a *XY) {
	var xyz XYZ
	xyz.SetXY(xy)
	xyz.AddXY(&xyz, a)
	xy.SetXYZ(&xyz)
}

/*
func (pk *XY) GetPublicKey() []byte {
	var out []byte = make([]byte, 65, 65)
	pk.X.GetB32(out[1:33])
	if len(out) == 65 {
		out[0] = 0x04
		pk.Y.GetB32(out[33:65])
	} else {
		if pk.Y.IsOdd() {
			out[0] = 0x03
		} else {
			out[0] = 0x02
		}
	}
	return out
}
*/

// GetPublicKey use compact format
//returns only 33 bytes
//same as bytes()
//TODO: deprecate, replace with .Bytes()
func (xy *XY) GetPublicKey() []byte {
	return xy.Bytes()
	/*
		var out []byte = make([]byte, 33, 33)
		pk.X.GetB32(out[1:33])
		if pk.Y.IsOdd() {
			out[0] = 0x03
		} else {
			out[0] = 0x02
		}
		return out
	*/
}
