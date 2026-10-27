package secp256k1go

import (
	"math/big"
	"math/rand"
	"testing"
)

// The big.Int forms rshX and maskBits replaced, kept as the reference.
func rshXRef(num *Number, bits uint) int {
	res := int(new(big.Int).And(&num.Int, new(big.Int).SetUint64((1<<bits)-1)).Uint64()) //nolint:gosec
	num.Rsh(&num.Int, bits)
	return res
}

func maskBitsRef(num *Number, bits uint) {
	mask := new(big.Int).Lsh(BigInt1, bits)
	mask.Sub(mask, BigInt1)
	num.Int.And(&num.Int, mask)
}

func randNumber(r *rand.Rand) *Number {
	b := make([]byte, r.Intn(40))
	r.Read(b) //nolint:errcheck,gosec // math/rand never fails
	var n Number
	n.SetBytes(b)
	if r.Intn(2) == 0 {
		n.Neg(&n.Int)
	}
	return &n
}

func TestRshXAndMaskBitsMatchBigInt(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	edges := []int64{0, 1, -1, 2, -2, 15, -15, 16, -16, 1 << 31, -(1 << 31), 1<<32 - 1, -(1<<32 - 1)}
	for i := 0; i < 200000; i++ {
		var a Number
		if i < len(edges) {
			a.SetInt64(edges[i])
		} else {
			a = *randNumber(r)
		}
		for _, bits := range []uint{1, 4, winA, 13, winG, 31, 32, 40} {
			var x, y Number
			x.Set(&a.Int)
			y.Set(&a.Int)
			if got, want := x.rshX(bits), rshXRef(&y, bits); got != want || x.Cmp(&y.Int) != 0 {
				t.Fatalf("rshX(%s, %d) = %d, %s; want %d, %s", a.String(), bits, got, x.String(), want, y.String())
			}
		}
		for _, bits := range []uint{0, 1, 31, 32, 63, 64, 65, 128, 200, 320, uint(r.Intn(330))} {
			var x, y Number
			x.Set(&a.Int)
			y.Set(&a.Int)
			x.maskBits(bits)
			maskBitsRef(&y, bits)
			if x.Cmp(&y.Int) != 0 {
				t.Fatalf("maskBits(%s, %d) = %s; want %s", a.String(), bits, x.String(), y.String())
			}
		}
	}
}
