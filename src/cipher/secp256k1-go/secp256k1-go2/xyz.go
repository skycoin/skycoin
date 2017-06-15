package secp256k1go

import (
	"fmt"
	//	"encoding/hex"
)

// XYZ contains xyz fields
type XYZ struct {
	X, Y, Z  Field
	Infinity bool
}

// Print prints xyz
func (xyz XYZ) Print(lab string) {
	if xyz.Infinity {
		fmt.Println(lab + " - INFINITY")
		return
	}
	fmt.Println(lab+".X", xyz.X.String())
	fmt.Println(lab+".Y", xyz.Y.String())
	fmt.Println(lab+".Z", xyz.Z.String())
}

// SetXY sets xy
func (xyz *XYZ) SetXY(a *XY) {
	xyz.Infinity = a.Infinity
	xyz.X = a.X
	xyz.Y = a.Y
	xyz.Z.SetInt(1)
}

// IsInfinity check if xyz is infinity
func (xyz *XYZ) IsInfinity() bool {
	return xyz.Infinity
}

// IsValid check if xyz is valid
func (xyz *XYZ) IsValid() bool {
	if xyz.Infinity {
		return false
	}
	var y2, x3, z2, z6 Field
	xyz.Y.Sqr(&y2)
	xyz.X.Sqr(&x3)
	x3.Mul(&x3, &xyz.X)
	xyz.Z.Sqr(&z2)
	z2.Sqr(&z6)
	z6.Mul(&z6, &z2)
	z6.MulInt(7)
	x3.SetAdd(&z6)
	y2.Normalize()
	x3.Normalize()
	return y2.Equals(&x3)
}

func (xyz *XYZ) getX(r *Field) {
	var zi2 Field
	xyz.Z.InvVar(&zi2)
	zi2.Sqr(&zi2)
	xyz.X.Mul(r, &zi2)
}

// Normalize normalize all fields
func (xyz *XYZ) Normalize() {
	xyz.X.Normalize()
	xyz.Y.Normalize()
	xyz.Z.Normalize()
}

// Equals checks if equal
func (xyz *XYZ) Equals(b *XYZ) bool {
	if xyz.Infinity != b.Infinity {
		return false
	}
	// TODO: is the normalize really needed here?
	xyz.Normalize()
	b.Normalize()
	return xyz.X.Equals(&b.X) && xyz.Y.Equals(&b.Y) && xyz.Z.Equals(&b.Z)
}

func (xyz *XYZ) precomp(w int) (pre []XYZ) {
	var d XYZ
	pre = make([]XYZ, (1 << (uint(w) - 2)))
	pre[0] = *xyz
	pre[0].Double(&d)
	for i := 1; i < len(pre); i++ {
		d.Add(&pre[i], &pre[i-1])
	}
	return
}

func ecmultWnaf(wnaf []int, a *Number, w uint) (ret int) {
	var zeroes uint
	var X Number
	X.Set(&a.Int)

	for X.Sign() != 0 {
		for X.Bit(0) == 0 {
			zeroes++
			X.rsh(1)
		}
		word := X.rshX(w)
		for zeroes > 0 {
			wnaf[ret] = 0
			ret++
			zeroes--
		}
		if (word & (1 << (w - 1))) != 0 {
			X.inc()
			wnaf[ret] = (word - (1 << w))
		} else {
			wnaf[ret] = word
		}
		zeroes = w - 1
		ret++
	}
	return
}

// ECmult  r = na*a + ng*G
func (xyz *XYZ) ECmult(r *XYZ, na, ng *Number) {
	var na1, naLam, ng1, ng128 Number

	// split na into na_1 and na_lam (where na = na_1 + na_lam*lambda, and na_1 and na_lam are ~128 bit)
	na.splitExp(&na1, &naLam)

	// split ng into ng_1 and ng_128 (where gn = gn_1 + gn_128*2^128, and gn_1 and gn_128 are ~128 bit)
	ng.split(&ng1, &ng128, 128)

	// build wnaf representation for na_1, na_lam, ng_1, ng_128
	var wnafNa1, wnafNaLam, wnafNg1, wnafNg128 [129]int
	bitsNa1 := ecmultWnaf(wnafNa1[:], &na1, winA)
	bitsNaLam := ecmultWnaf(wnafNaLam[:], &naLam, winA)
	bitsNg1 := ecmultWnaf(wnafNg1[:], &ng1, winG)
	bitsNg128 := ecmultWnaf(wnafNg128[:], &ng128, winG)

	// calculate a_lam = a*lambda
	var aLam XYZ
	xyz.mulLambda(&aLam)

	// calculate odd multiples of a and a_lam
	preA1 := xyz.precomp(winA)
	preALam := aLam.precomp(winA)

	bits := bitsNa1
	if bitsNaLam > bits {
		bits = bitsNaLam
	}
	if bitsNg1 > bits {
		bits = bitsNg1
	}
	if bitsNg128 > bits {
		bits = bitsNg128
	}

	r.Infinity = true

	var tmpj XYZ
	var tmpa XY
	var n int

	for i := bits - 1; i >= 0; i-- {
		r.Double(r)

		if i < bitsNa1 {
			n = wnafNa1[i]
			if n > 0 {
				r.Add(r, &preA1[((n)-1)/2])
			} else if n != 0 {
				preA1[(-(n)-1)/2].Neg(&tmpj)
				r.Add(r, &tmpj)
			}
		}

		if i < bitsNaLam {
			n = wnafNaLam[i]
			if n > 0 {
				r.Add(r, &preALam[((n)-1)/2])
			} else if n != 0 {
				preALam[(-(n)-1)/2].Neg(&tmpj)
				r.Add(r, &tmpj)
			}
		}

		if i < bitsNg1 {
			n = wnafNg1[i]
			if n > 0 {
				r.AddXY(r, &preG[((n)-1)/2])
			} else if n != 0 {
				preG[(-(n)-1)/2].Neg(&tmpa)
				r.AddXY(r, &tmpa)
			}
		}

		if i < bitsNg128 {
			n = wnafNg128[i]
			if n > 0 {
				r.AddXY(r, &preG128[((n)-1)/2])
			} else if n != 0 {
				preG128[(-(n)-1)/2].Neg(&tmpa)
				r.AddXY(r, &tmpa)
			}
		}
	}
}

// Neg caculate neg
func (xyz *XYZ) Neg(r *XYZ) {
	r.Infinity = xyz.Infinity
	r.X = xyz.X
	r.Y = xyz.Y
	r.Z = xyz.Z
	r.Y.Normalize()
	r.Y.Negate(&r.Y, 1)
}

func (xyz *XYZ) mulLambda(r *XYZ) {
	*r = *xyz
	r.X.Mul(&r.X, &TheCurve.beta)
}

// Double cacule double
func (xyz *XYZ) Double(r *XYZ) {
	var t1, t2, t3, t4, t5 Field

	t5 = xyz.Y
	t5.Normalize()
	if xyz.Infinity || t5.IsZero() {
		r.Infinity = true
		return
	}

	t5.Mul(&r.Z, &xyz.Z)
	r.Z.MulInt(2)
	xyz.X.Sqr(&t1)
	t1.MulInt(3)
	t1.Sqr(&t2)
	t5.Sqr(&t3)
	t3.MulInt(2)
	t3.Sqr(&t4)
	t4.MulInt(2)
	xyz.X.Mul(&t3, &t3)
	r.X = t3
	r.X.MulInt(4)
	r.X.Negate(&r.X, 4)
	r.X.SetAdd(&t2)
	t2.Negate(&t2, 1)
	t3.MulInt(6)
	t3.SetAdd(&t2)
	t1.Mul(&r.Y, &t3)
	t4.Negate(&t2, 2)
	r.Y.SetAdd(&t2)
	r.Infinity = false
}

// AddXY adds XY
func (xyz *XYZ) AddXY(r *XYZ, b *XY) {
	if xyz.Infinity {
		r.Infinity = b.Infinity
		r.X = b.X
		r.Y = b.Y
		r.Z.SetInt(1)
		return
	}
	if b.Infinity {
		*r = *xyz
		return
	}
	r.Infinity = false
	var z12, u1, u2, s1, s2 Field
	xyz.Z.Sqr(&z12)
	u1 = xyz.X
	u1.Normalize()
	b.X.Mul(&u2, &z12)
	s1 = xyz.Y
	s1.Normalize()
	b.Y.Mul(&s2, &z12)
	s2.Mul(&s2, &xyz.Z)
	u1.Normalize()
	u2.Normalize()

	if u1.Equals(&u2) {
		s1.Normalize()
		s2.Normalize()
		if s1.Equals(&s2) {
			xyz.Double(r)
		} else {
			r.Infinity = true
		}
		return
	}

	var h, i, i2, h2, h3, t Field
	u1.Negate(&h, 1)
	h.SetAdd(&u2)
	s1.Negate(&i, 1)
	i.SetAdd(&s2)
	i.Sqr(&i2)
	h.Sqr(&h2)
	h.Mul(&h3, &h2)
	r.Z = xyz.Z
	r.Z.Mul(&r.Z, &h)
	u1.Mul(&t, &h2)
	r.X = t
	r.X.MulInt(2)
	r.X.SetAdd(&h3)
	r.X.Negate(&r.X, 3)
	r.X.SetAdd(&i2)
	r.X.Negate(&r.Y, 5)
	r.Y.SetAdd(&t)
	r.Y.Mul(&r.Y, &i)
	h3.Mul(&h3, &s1)
	h3.Negate(&h3, 1)
	r.Y.SetAdd(&h3)
}

// Add adds value
func (xyz *XYZ) Add(r, b *XYZ) {
	if xyz.Infinity {
		*r = *b
		return
	}
	if b.Infinity {
		*r = *xyz
		return
	}
	r.Infinity = false
	var z22, z12, u1, u2, s1, s2 Field

	b.Z.Sqr(&z22)
	xyz.Z.Sqr(&z12)
	xyz.X.Mul(&u1, &z22)
	b.X.Mul(&u2, &z12)
	xyz.Y.Mul(&s1, &z22)
	s1.Mul(&s1, &b.Z)
	b.Y.Mul(&s2, &z12)
	s2.Mul(&s2, &xyz.Z)
	u1.Normalize()
	u2.Normalize()
	if u1.Equals(&u2) {
		s1.Normalize()
		s2.Normalize()
		if s1.Equals(&s2) {
			xyz.Double(r)
		} else {
			r.Infinity = true
		}
		return
	}
	var h, i, i2, h2, h3, t Field

	u1.Negate(&h, 1)
	h.SetAdd(&u2)
	s1.Negate(&i, 1)
	i.SetAdd(&s2)
	i.Sqr(&i2)
	h.Sqr(&h2)
	h.Mul(&h3, &h2)
	xyz.Z.Mul(&r.Z, &b.Z)
	r.Z.Mul(&r.Z, &h)
	u1.Mul(&t, &h2)
	r.X = t
	r.X.MulInt(2)
	r.X.SetAdd(&h3)
	r.X.Negate(&r.X, 3)
	r.X.SetAdd(&i2)
	r.X.Negate(&r.Y, 5)
	r.Y.SetAdd(&t)
	r.Y.Mul(&r.Y, &i)
	h3.Mul(&h3, &s1)
	h3.Negate(&h3, 1)
	r.Y.SetAdd(&h3)
}

// ECmultGen r = a*G
//TODO: Change to returning result
//TODO: input should not be pointer
func ECmultGen(r *XYZ, a *Number) {
	var n Number
	n.Set(&a.Int)
	r.SetXY(&prec[0][n.rshX(4)])
	for j := 1; j < 64; j++ {
		r.AddXY(r, &prec[j][n.rshX(4)])
	}
	r.AddXY(r, &fin)
}
