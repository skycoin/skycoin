package secp256k1go

// This file contains micro-optimizations to existing field operations
// that provide immediate performance improvements without risky rewrites.

// MulInt2 multiplies a field element by 2 (optimized version)
func (fd *Field) MulInt2() {
	// Multiply by 2 is just a left shift by 1 bit
	// This is faster than the generic MulInt(2)
	fd.n[0] <<= 1
	fd.n[1] <<= 1
	fd.n[2] <<= 1
	fd.n[3] <<= 1
	fd.n[4] <<= 1
	fd.n[5] <<= 1
	fd.n[6] <<= 1
	fd.n[7] <<= 1
	fd.n[8] <<= 1
	fd.n[9] <<= 1
}

// MulInt4 multiplies a field element by 4 (optimized version)
func (fd *Field) MulInt4() {
	// Multiply by 4 is a left shift by 2 bits
	fd.n[0] <<= 2
	fd.n[1] <<= 2
	fd.n[2] <<= 2
	fd.n[3] <<= 2
	fd.n[4] <<= 2
	fd.n[5] <<= 2
	fd.n[6] <<= 2
	fd.n[7] <<= 2
	fd.n[8] <<= 2
	fd.n[9] <<= 2
}

// MulInt8 multiplies a field element by 8 (optimized version)
func (fd *Field) MulInt8() {
	// Multiply by 8 is a left shift by 3 bits
	fd.n[0] <<= 3
	fd.n[1] <<= 3
	fd.n[2] <<= 3
	fd.n[3] <<= 3
	fd.n[4] <<= 3
	fd.n[5] <<= 3
	fd.n[6] <<= 3
	fd.n[7] <<= 3
	fd.n[8] <<= 3
	fd.n[9] <<= 3
}

// NormalizeWeak performs weak normalization (magnitude 1 but not fully reduced)
// This is faster than full Normalize when full reduction isn't required.
func (fd *Field) NormalizeWeak() {
	t0, t1, t2, t3, t4 := fd.n[0], fd.n[1], fd.n[2], fd.n[3], fd.n[4]
	t5, t6, t7, t8, t9 := fd.n[5], fd.n[6], fd.n[7], fd.n[8], fd.n[9]

	// Reduce t9 at the start so there will be at most a single carry
	x := t9 >> 22
	t9 &= 0x03FFFFF

	// First pass ensures magnitude is 1
	t0 += x * 0x3D1
	t1 += x << 6
	t1 += t0 >> 26
	t0 &= 0x3FFFFFF
	t2 += t1 >> 26
	t1 &= 0x3FFFFFF
	t3 += t2 >> 26
	t2 &= 0x3FFFFFF
	t4 += t3 >> 26
	t3 &= 0x3FFFFFF
	t5 += t4 >> 26
	t4 &= 0x3FFFFFF
	t6 += t5 >> 26
	t5 &= 0x3FFFFFF
	t7 += t6 >> 26
	t6 &= 0x3FFFFFF
	t8 += t7 >> 26
	t7 &= 0x3FFFFFF
	t9 += t8 >> 26
	t8 &= 0x3FFFFFF

	fd.n[0], fd.n[1], fd.n[2], fd.n[3], fd.n[4] = t0, t1, t2, t3, t4
	fd.n[5], fd.n[6], fd.n[7], fd.n[8], fd.n[9] = t5, t6, t7, t8, t9
}

// Cmov performs constant-time conditional move: r = flag ? a : r
// This is critical for constant-time operations with secret keys.
func (fd *Field) Cmov(a *Field, flag bool) {
	// Convert flag to mask (all 1s if true, all 0s if false)
	var mask uint32
	if flag {
		mask = 0xFFFFFFFF
	}

	// Constant-time select: r = (r & ~mask) | (a & mask)
	fd.n[0] = (fd.n[0] &^ mask) | (a.n[0] & mask)
	fd.n[1] = (fd.n[1] &^ mask) | (a.n[1] & mask)
	fd.n[2] = (fd.n[2] &^ mask) | (a.n[2] & mask)
	fd.n[3] = (fd.n[3] &^ mask) | (a.n[3] & mask)
	fd.n[4] = (fd.n[4] &^ mask) | (a.n[4] & mask)
	fd.n[5] = (fd.n[5] &^ mask) | (a.n[5] & mask)
	fd.n[6] = (fd.n[6] &^ mask) | (a.n[6] & mask)
	fd.n[7] = (fd.n[7] &^ mask) | (a.n[7] & mask)
	fd.n[8] = (fd.n[8] &^ mask) | (a.n[8] & mask)
	fd.n[9] = (fd.n[9] &^ mask) | (a.n[9] & mask)
}

// NegateUnchecked computes r = -a (mod p) with given magnitude
func (fd *Field) NegateUnchecked(r *Field, m uint32) {
	// From C library: r->n[i] = constant * 2 * (m + 1) - a->n[i]
	r.n[0] = 0x3FFFC2F*2*(m+1) - fd.n[0]
	r.n[1] = 0x3FFFFBF*2*(m+1) - fd.n[1]
	r.n[2] = 0x3FFFFFF*2*(m+1) - fd.n[2]
	r.n[3] = 0x3FFFFFF*2*(m+1) - fd.n[3]
	r.n[4] = 0x3FFFFFF*2*(m+1) - fd.n[4]
	r.n[5] = 0x3FFFFFF*2*(m+1) - fd.n[5]
	r.n[6] = 0x3FFFFFF*2*(m+1) - fd.n[6]
	r.n[7] = 0x3FFFFFF*2*(m+1) - fd.n[7]
	r.n[8] = 0x3FFFFFF*2*(m+1) - fd.n[8]
	r.n[9] = 0x03FFFFF*2*(m+1) - fd.n[9]
}
