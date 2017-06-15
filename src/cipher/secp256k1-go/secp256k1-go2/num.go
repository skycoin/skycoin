package secp256k1go

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

var (
	// BigInt1 represents big int with value 1
	BigInt1 = new(big.Int).SetInt64(1)
)

// Number wraps the big.Int
type Number struct {
	big.Int
}

// Print prints the label with hex number string
func (num *Number) Print(label string) {
	fmt.Println(label, hex.EncodeToString(num.Bytes()))
}

func (num *Number) modMul(a, b, m *Number) {
	num.Mul(&a.Int, &b.Int)
	num.Mod(&num.Int, &m.Int)
	return
}

func (num *Number) modInv(a, b *Number) {
	num.ModInverse(&a.Int, &b.Int)
	return
}

func (num *Number) mod(a *Number) {
	num.Mod(&num.Int, &a.Int)
	return
}

// SetHex sets number from string
func (num *Number) SetHex(s string) {
	num.SetString(s, 16)
}

//SetBytes and GetBytes are inherited by default
//added
//func (a *Number) SetBytes(b []byte) {
//	a.SetBytes(b)
//}

func (num *Number) maskBits(bits uint) {
	mask := new(big.Int).Lsh(BigInt1, bits)
	mask.Sub(mask, BigInt1)
	num.Int.And(&num.Int, mask)
}

func (num *Number) splitExp(r1, r2 *Number) {
	var bnc1, bnc2, bnn2, bnt1, bnt2 Number

	bnn2.Int.Rsh(&TheCurve.Order.Int, 1)

	bnc1.Mul(&num.Int, &TheCurve.a1b2.Int)
	bnc1.Add(&bnc1.Int, &bnn2.Int)
	bnc1.Div(&bnc1.Int, &TheCurve.Order.Int)

	bnc2.Mul(&num.Int, &TheCurve.b1.Int)
	bnc2.Add(&bnc2.Int, &bnn2.Int)
	bnc2.Div(&bnc2.Int, &TheCurve.Order.Int)

	bnt1.Mul(&bnc1.Int, &TheCurve.a1b2.Int)
	bnt2.Mul(&bnc2.Int, &TheCurve.a2.Int)
	bnt1.Add(&bnt1.Int, &bnt2.Int)
	r1.Sub(&num.Int, &bnt1.Int)

	bnt1.Mul(&bnc1.Int, &TheCurve.b1.Int)
	bnt2.Mul(&bnc2.Int, &TheCurve.a1b2.Int)
	r2.Sub(&bnt1.Int, &bnt2.Int)
}

func (num *Number) split(rl, rh *Number, bits uint) {
	rl.Int.Set(&num.Int)
	rh.Int.Rsh(&rl.Int, bits)
	rl.maskBits(bits)
}

func (num *Number) rsh(bits uint) {
	num.Rsh(&num.Int, bits)
}

func (num *Number) inc() {
	num.Add(&num.Int, BigInt1)
}

func (num *Number) rshX(bits uint) (res int) {
	res = int(new(big.Int).And(&num.Int, new(big.Int).SetUint64((1<<bits)-1)).Uint64())
	num.Rsh(&num.Int, bits)
	return
}

// IsOdd checks if is odd
func (num *Number) IsOdd() bool {
	return num.Bit(0) != 0
}

func (num *Number) getBin(le int) []byte {
	bts := num.Bytes()
	if len(bts) > le {
		panic("buffer too small")
	}
	if len(bts) == le {
		return bts
	}
	return append(make([]byte, le-len(bts)), bts...)
}
