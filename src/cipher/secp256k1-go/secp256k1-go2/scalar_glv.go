package secp256k1go

// scalar_glv.go - Optimized GLV scalar decomposition
// Port of secp256k1_scalar_split_lambda from bitcoin-core/secp256k1
//
// This optimizes the scalar split to use precomputed constants g1, g2
// with mul_shift instead of slow big.Int division

import (
	"math/big"
)

// GLV constants for scalar decomposition
// These are ported from bitcoin-core/secp256k1/src/scalar_impl.h
var (
	// lambda: k1 + k2*lambda = k (mod n)
	scalarLambda = mustParseBig("5363AD4CC05C30E0A5261C028812645A122E22EA20816678DF02967C1B23BD72")

	// minus_b1, minus_b2: used in scalar split
	scalarMinusB1 = mustParseBig("00000000000000000000000000000000E4437ED6010E88286F547FA90ABFE4C3")
	scalarMinusB2 = mustParseBig("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE8A280AC50774346DD765CDA83DB1562C")

	// g1, g2: precomputed for fast rounding
	// g1 = round(2^384 * b2/n)
	// g2 = round(2^384 * (-b1)/n)
	scalarG1 = mustParseBig("3086D221A7D46BCDE86C90E49284EB153DAA8A1471E8CA7FE893209A45DBB031")
	scalarG2 = mustParseBig("E4437ED6010E88286F547FA90ABFE4C4221208AC9DF506C61571B4AE8AC47F71")
)

func mustParseBig(s string) *big.Int {
	n := new(big.Int)
	n, ok := n.SetString(s, 16) // Always use base 16 for hex strings
	if !ok {
		panic("failed to parse big int: " + s)
	}
	return n
}

// scalarMulShift computes (a * b) >> shift
// This is much faster than division for our use case
func scalarMulShift(result, a, b *big.Int, shift uint) {
	result.Mul(a, b)
	result.Rsh(result, shift)
}

// SplitScalarFast splits scalar k into k1, k2 such that k = k1 + k2*lambda (mod n)
// This uses the optimized algorithm from bitcoin-core/secp256k1
// It's much faster than the original splitExp which uses division
func SplitScalarFast(k1, k2, k *Number) {
	var c1, c2, tmp big.Int

	// Compute c1 = round(k * g1 / 2^384) = (k * g1) >> 384
	scalarMulShift(&c1, &k.Int, scalarG1, 384)

	// Compute c2 = round(k * g2 / 2^384) = (k * g2) >> 384
	scalarMulShift(&c2, &k.Int, scalarG2, 384)

	// r2 = -(c1 * minus_b1 + c2 * minus_b2)
	// Note: In C this is: r2 = c1*b1 + c2*b2
	// But we have minus_b1, minus_b2, so: r2 = -(c1*(-b1) + c2*(-b2)) = c1*b1 + c2*b2
	tmp.Mul(&c1, scalarMinusB1)
	k2.Mul(&c2, scalarMinusB2)
	k2.Add(&k2.Int, &tmp)

	// r2 = r2 mod n
	k2.Mod(&k2.Int, &TheCurve.Order.Int)

	// r1 = k - r2 * lambda (mod n)
	tmp.Mul(&k2.Int, scalarLambda)
	k1.Sub(&k.Int, &tmp)
	k1.Mod(&k1.Int, &TheCurve.Order.Int)
}
