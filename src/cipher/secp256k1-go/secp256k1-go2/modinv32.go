package secp256k1go

// modinv32.go - Port of bitcoin-core/secp256k1 modular inverse using safegcd
//
// This implements modular inversion based on the paper "Fast constant-time gcd computation and
// modular inversion" by Daniel J. Bernstein and Bo-Yin Yang.
//
// For an explanation, see doc/safegcd_implementation.md in bitcoin-core/secp256k1
// This file implements the algorithm for N=30, using 30-bit signed limbs as int32.

// modInv32Signed30 represents a signed integer in 30-bit limbs
// Value is sum(v[i] * 2^(30*i), i=0..8)
type modInv32Signed30 struct {
	v [9]int32
}

// modInv32Trans2x2 is a 2x2 transition matrix
// t = [ u  v ]
//     [ q  r ]
type modInv32Trans2x2 struct {
	u, v, q, r int32
}

// modInv32ModInfo contains modulus information
type modInv32ModInfo struct {
	modulus      modInv32Signed30
	modulusInv30 uint32 // modulus^{-1} mod 2^30
}

// fieldModInfo is the modulus info for the secp256k1 field prime
// p = 2^256 - 2^32 - 977 = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F
var fieldModInfo modInv32ModInfo

func init() {
	// Initialize the secp256k1 field prime in 30-bit signed limbs
	// p = 2^256 - 2^32 - 977
	// The C library represents this as: {{-0x3D1, -4, 0, 0, 0, 0, 0, 0, 65536}}
	fieldModInfo.modulus.v[0] = -0x3D1   // -(2^32 + 977) mod 2^30 = -977
	fieldModInfo.modulus.v[1] = -4       // -2^32 / 2^30 = -4
	fieldModInfo.modulus.v[2] = 0
	fieldModInfo.modulus.v[3] = 0
	fieldModInfo.modulus.v[4] = 0
	fieldModInfo.modulus.v[5] = 0
	fieldModInfo.modulus.v[6] = 0
	fieldModInfo.modulus.v[7] = 0
	fieldModInfo.modulus.v[8] = 65536    // 2^256 / 2^240 = 2^16
	
	// modulus^{-1} mod 2^30 (from C library)
	fieldModInfo.modulusInv30 = 0x2DDACACF
}

// modInv32Inv256 table: modInv32Inv256[i] = -(2*i+1)^-1 (mod 256)
var modInv32Inv256 = [128]uint8{
	0xFF, 0x55, 0x33, 0x49, 0xC7, 0x5D, 0x3B, 0x11, 0x0F, 0xE5, 0xC3, 0x59,
	0xD7, 0xED, 0xCB, 0x21, 0x1F, 0x75, 0x53, 0x69, 0xE7, 0x7D, 0x5B, 0x31,
	0x2F, 0x05, 0xE3, 0x79, 0xF7, 0x0D, 0xEB, 0x41, 0x3F, 0x95, 0x73, 0x89,
	0x07, 0x9D, 0x7B, 0x51, 0x4F, 0x25, 0x03, 0x99, 0x17, 0x2D, 0x0B, 0x61,
	0x5F, 0xB5, 0x93, 0xA9, 0x27, 0xBD, 0x9B, 0x71, 0x6F, 0x45, 0x23, 0xB9,
	0x37, 0x4D, 0x2B, 0x81, 0x7F, 0xD5, 0xB3, 0xC9, 0x47, 0xDD, 0xBB, 0x91,
	0x8F, 0x65, 0x43, 0xD9, 0x57, 0x6D, 0x4B, 0xA1, 0x9F, 0xF5, 0xD3, 0xE9,
	0x67, 0xFD, 0xDB, 0xB1, 0xAF, 0x85, 0x63, 0xF9, 0x77, 0x8D, 0x6B, 0xC1,
	0xBF, 0x15, 0xF3, 0x09, 0x87, 0x1D, 0xFB, 0xD1, 0xCF, 0xA5, 0x83, 0x19,
	0x97, 0xAD, 0x8B, 0xE1, 0xDF, 0x35, 0x13, 0x29, 0xA7, 0x3D, 0x1B, 0xF1,
	0xEF, 0xC5, 0xA3, 0x39, 0xB7, 0xCD, 0xAB, 0x01,
}

// ctz32Var counts trailing zeros (variable time)
func ctz32Var(x uint32) int {
	if x == 0 {
		return 32
	}
	var count int
	for (x & 1) == 0 {
		x >>= 1
		count++
	}
	return count
}

// modInv32Normalize30 takes a signed30 number in range (-2*modulus, modulus), and adds a multiple
// of the modulus to it to bring it to range [0, modulus). If sign < 0, the input will also be
// negated in the process.
func modInv32Normalize30(r *modInv32Signed30, sign int32, modinfo *modInv32ModInfo) {
	const M30 = int32(0x3FFFFFFF) // 2^30 - 1

	r0, r1, r2, r3, r4 := r.v[0], r.v[1], r.v[2], r.v[3], r.v[4]
	r5, r6, r7, r8 := r.v[5], r.v[6], r.v[7], r.v[8]

	var condAdd, condNegate int32

	// In a first step, add the modulus if the input is negative, and then negate if requested.
	// This brings r from range (-2*modulus, modulus) to range (-modulus, modulus).
	condAdd = r8 >> 31 // Sign extend: -1 if negative, 0 if positive
	r0 += modinfo.modulus.v[0] & condAdd
	r1 += modinfo.modulus.v[1] & condAdd
	r2 += modinfo.modulus.v[2] & condAdd
	r3 += modinfo.modulus.v[3] & condAdd
	r4 += modinfo.modulus.v[4] & condAdd
	r5 += modinfo.modulus.v[5] & condAdd
	r6 += modinfo.modulus.v[6] & condAdd
	r7 += modinfo.modulus.v[7] & condAdd
	r8 += modinfo.modulus.v[8] & condAdd

	condNegate = sign >> 31
	r0 = (r0 ^ condNegate) - condNegate
	r1 = (r1 ^ condNegate) - condNegate
	r2 = (r2 ^ condNegate) - condNegate
	r3 = (r3 ^ condNegate) - condNegate
	r4 = (r4 ^ condNegate) - condNegate
	r5 = (r5 ^ condNegate) - condNegate
	r6 = (r6 ^ condNegate) - condNegate
	r7 = (r7 ^ condNegate) - condNegate
	r8 = (r8 ^ condNegate) - condNegate

	// Propagate the top bits, to bring limbs back to range (-2^30, 2^30)
	r1 += r0 >> 30
	r0 &= M30
	r2 += r1 >> 30
	r1 &= M30
	r3 += r2 >> 30
	r2 &= M30
	r4 += r3 >> 30
	r3 &= M30
	r5 += r4 >> 30
	r4 &= M30
	r6 += r5 >> 30
	r5 &= M30
	r7 += r6 >> 30
	r6 &= M30
	r8 += r7 >> 30
	r7 &= M30

	// In a second step add the modulus again if the result is still negative, bringing r to range [0, modulus)
	condAdd = r8 >> 31
	r0 += modinfo.modulus.v[0] & condAdd
	r1 += modinfo.modulus.v[1] & condAdd
	r2 += modinfo.modulus.v[2] & condAdd
	r3 += modinfo.modulus.v[3] & condAdd
	r4 += modinfo.modulus.v[4] & condAdd
	r5 += modinfo.modulus.v[5] & condAdd
	r6 += modinfo.modulus.v[6] & condAdd
	r7 += modinfo.modulus.v[7] & condAdd
	r8 += modinfo.modulus.v[8] & condAdd

	// And propagate again
	r1 += r0 >> 30
	r0 &= M30
	r2 += r1 >> 30
	r1 &= M30
	r3 += r2 >> 30
	r2 &= M30
	r4 += r3 >> 30
	r3 &= M30
	r5 += r4 >> 30
	r4 &= M30
	r6 += r5 >> 30
	r5 &= M30
	r7 += r6 >> 30
	r6 &= M30
	r8 += r7 >> 30
	r7 &= M30

	r.v[0], r.v[1], r.v[2], r.v[3], r.v[4] = r0, r1, r2, r3, r4
	r.v[5], r.v[6], r.v[7], r.v[8] = r5, r6, r7, r8
}

// modInv32Divsteps30Var computes the transition matrix and eta for 30 divsteps (variable time).
// Implements the divsteps_n_matrix_var function from the safegcd paper.
func modInv32Divsteps30Var(eta int32, f0, g0 uint32) (newEta int32, t modInv32Trans2x2) {
	// Transformation matrix; starts with identity
	u, v, q, r := uint32(1), uint32(0), uint32(0), uint32(1)
	f, g := f0, g0
	var w uint16
	i := 30

	for {
		// Use a sentinel bit to count zeros only up to i
		zeros := ctz32Var(g | (0xFFFFFFFF << uint(i)))
		// Perform zeros divsteps at once; they all just divide g by two
		g >>= uint(zeros)
		u <<= uint(zeros)
		v <<= uint(zeros)
		eta -= int32(zeros)
		i -= zeros

		// We're done once we've done 30 divsteps
		if i == 0 {
			break
		}

		// If eta is negative, negate it and replace f,g with g,-f
		if eta < 0 {
			eta = -eta
			f, g = g, -f
			u, q = q, -u
			v, r = r, -v
		}

		// eta is now >= 0. In what follows we're going to cancel out the bottom bits of g.
		// No more than i can be cancelled out, and no more than eta+1
		limit := int(eta) + 1
		if limit > i {
			limit = i
		}

		// m is a mask for the bottom min(limit, 8) bits (our table only supports 8 bits)
		m := (uint32(0xFFFFFFFF) >> uint(32-limit)) & 255

		// Find what multiple of f must be added to g to cancel its bottom min(limit, 8) bits
		w = uint16((g * uint32(modInv32Inv256[(f>>1)&127])) & m)

		// Do so
		g += f * uint32(w)
		q += u * uint32(w)
		r += v * uint32(w)
	}

	// Return data in t and return value
	t.u = int32(u)
	t.v = int32(v)
	t.q = int32(q)
	t.r = int32(r)

	return eta, t
}

// fieldToModInv32 converts a Field to modInv32Signed30 format
// Port of secp256k1_fe_to_signed30 from field_10x26_impl.h
func fieldToModInv32(r *modInv32Signed30, a *Field) {
	const M30 = uint64(0x3FFFFFFF) // 2^30 - 1
	a0, a1, a2, a3, a4 := uint64(a.n[0]), uint64(a.n[1]), uint64(a.n[2]), uint64(a.n[3]), uint64(a.n[4])
	a5, a6, a7, a8, a9 := uint64(a.n[5]), uint64(a.n[6]), uint64(a.n[7]), uint64(a.n[8]), uint64(a.n[9])

	r.v[0] = int32((a0 | a1<<26) & M30)
	r.v[1] = int32((a1>>4 | a2<<22) & M30)
	r.v[2] = int32((a2>>8 | a3<<18) & M30)
	r.v[3] = int32((a3>>12 | a4<<14) & M30)
	r.v[4] = int32((a4>>16 | a5<<10) & M30)
	r.v[5] = int32((a5>>20 | a6<<6) & M30)
	r.v[6] = int32((a6>>24 | a7<<2 | a8<<28) & M30)
	r.v[7] = int32((a8>>2 | a9<<24) & M30)
	r.v[8] = int32(a9 >> 6)
}

// modInv32ToField converts modInv32Signed30 back to Field format
// Port of secp256k1_fe_from_signed30 from field_10x26_impl.h
func modInv32ToField(r *Field, a *modInv32Signed30) {
	const M26 = uint32(0x3FFFFFF) // 2^26 - 1
	a0, a1, a2, a3, a4 := uint32(a.v[0]), uint32(a.v[1]), uint32(a.v[2]), uint32(a.v[3]), uint32(a.v[4])
	a5, a6, a7, a8 := uint32(a.v[5]), uint32(a.v[6]), uint32(a.v[7]), uint32(a.v[8])

	r.n[0] = a0 & M26
	r.n[1] = (a0>>26 | a1<<4) & M26
	r.n[2] = (a1>>22 | a2<<8) & M26
	r.n[3] = (a2>>18 | a3<<12) & M26
	r.n[4] = (a3>>14 | a4<<16) & M26
	r.n[5] = (a4>>10 | a5<<20) & M26
	r.n[6] = (a5>>6 | a6<<24) & M26
	r.n[7] = (a6 >> 2) & M26
	r.n[8] = (a6>>28 | a7<<2) & M26
	r.n[9] = (a7>>24 | a8<<6)
}

// modInv32UpdateDE30 updates d and e using the transition matrix t, with modular reduction
// This computes (t/2^30) * [d, e] + modulus * [md, me], choosing md, me such that the bottom
// 30 bits of the result are zero
func modInv32UpdateDE30(d, e *modInv32Signed30, t *modInv32Trans2x2, modinfo *modInv32ModInfo) {
	const M30 = int32(0x3FFFFFFF) // 2^30 - 1
	u, v, q, r := t.u, t.v, t.q, t.r
	var di, ei, md, me, sd, se int32
	var cd, ce int64

	// [md,me] start as zero; plus [u,q] if d is negative; plus [v,r] if e is negative
	sd = d.v[8] >> 31
	se = e.v[8] >> 31
	md = (u & sd) + (v & se)
	me = (q & sd) + (r & se)

	// Begin computing t*[d,e]
	di = d.v[0]
	ei = e.v[0]
	cd = int64(u)*int64(di) + int64(v)*int64(ei)
	ce = int64(q)*int64(di) + int64(r)*int64(ei)

	// Correct md,me so that t*[d,e]+modulus*[md,me] has 30 zero bottom bits
	md -= int32((modinfo.modulusInv30*uint32(cd) + uint32(md)) & uint32(M30))
	me -= int32((modinfo.modulusInv30*uint32(ce) + uint32(me)) & uint32(M30))

	// Update the beginning of computation for t*[d,e]+modulus*[md,me] now md,me are known
	cd += int64(modinfo.modulus.v[0]) * int64(md)
	ce += int64(modinfo.modulus.v[0]) * int64(me)

	// Verify that the low 30 bits are zero, then throw them away
	cd >>= 30
	ce >>= 30

	// Now iteratively compute limb i=1..8 of t*[d,e]+modulus*[md,me],
	// and store them in output limb i-1 (shifting down by 30 bits)
	for i := 1; i < 9; i++ {
		di = d.v[i]
		ei = e.v[i]
		cd += int64(u)*int64(di) + int64(v)*int64(ei)
		ce += int64(q)*int64(di) + int64(r)*int64(ei)
		cd += int64(modinfo.modulus.v[i]) * int64(md)
		ce += int64(modinfo.modulus.v[i]) * int64(me)
		d.v[i-1] = int32(cd) & M30
		e.v[i-1] = int32(ce) & M30
		cd >>= 30
		ce >>= 30
	}

	// What remains is limb 9 of t*[d,e]+modulus*[md,me]; store it as output limb 8
	d.v[8] = int32(cd)
	e.v[8] = int32(ce)
}

// modInv32UpdateFG30Var updates f and g using the transition matrix t (variable time)
// This computes (t/2^30) * [f, g] for a variable number of limbs
func modInv32UpdateFG30Var(length int, f, g *modInv32Signed30, t *modInv32Trans2x2) {
	const M30 = int32(0x3FFFFFFF) // 2^30 - 1
	u, v, q, r := t.u, t.v, t.q, t.r
	var fi, gi int32
	var cf, cg int64

	// Start computing t*[f,g]
	fi = f.v[0]
	gi = g.v[0]
	cf = int64(u)*int64(fi) + int64(v)*int64(gi)
	cg = int64(q)*int64(fi) + int64(r)*int64(gi)

	// Verify that the bottom 30 bits are zero, then throw them away
	cf >>= 30
	cg >>= 30

	// Now iteratively compute limb i=1..length of t*[f,g],
	// and store them in output limb i-1 (shifting down by 30 bits)
	for i := 1; i < length; i++ {
		fi = f.v[i]
		gi = g.v[i]
		cf += int64(u)*int64(fi) + int64(v)*int64(gi)
		cg += int64(q)*int64(fi) + int64(r)*int64(gi)
		f.v[i-1] = int32(cf) & M30
		g.v[i-1] = int32(cg) & M30
		cf >>= 30
		cg >>= 30
	}

	// What remains is limb length of t*[f,g]; store it as output limb length-1
	f.v[length-1] = int32(cf)
	g.v[length-1] = int32(cg)
}

// modInv32Var computes the modular inverse of x modulo modinfo.modulus
// This is the variable-time main loop that performs iterations of divsteps
func modInv32Var(x *modInv32Signed30, modinfo *modInv32ModInfo) {
	// Start with d=0, e=1, f=modulus, g=x, eta=-1
	d := modInv32Signed30{v: [9]int32{0, 0, 0, 0, 0, 0, 0, 0, 0}}
	e := modInv32Signed30{v: [9]int32{1, 0, 0, 0, 0, 0, 0, 0, 0}}
	f := modinfo.modulus
	g := *x

	length := 9
	eta := int32(-1) // eta = -delta; delta is initially 1

	// Do iterations of 30 divsteps each until g=0
	for {
		// Compute transition matrix and new eta after 30 divsteps
		var t modInv32Trans2x2
		eta, t = modInv32Divsteps30Var(eta, uint32(f.v[0]), uint32(g.v[0]))

		// Update d,e using that transition matrix
		modInv32UpdateDE30(&d, &e, &t, modinfo)

		// Update f,g using that transition matrix
		modInv32UpdateFG30Var(length, &f, &g, &t)

		// If the bottom limb of g is 0, there is a chance g=0
		if g.v[0] == 0 {
			cond := int32(0)
			// Check if all other limbs are also 0
			for j := 1; j < length; j++ {
				cond |= g.v[j]
			}
			// If so, we're done
			if cond == 0 {
				break
			}
		}

		// Determine if length>1 and limb (length-1) of both f and g is 0 or -1
		fn := f.v[length-1]
		gn := g.v[length-1]
		cond := (int32(length) - 2) >> 31
		cond |= fn ^ (fn >> 31)
		cond |= gn ^ (gn >> 31)

		// If so, reduce length, propagating the sign of f and g's top limb into the one below
		if cond == 0 {
			f.v[length-2] |= int32(uint32(fn) << 30)
			g.v[length-2] |= int32(uint32(gn) << 30)
			length--
		}
	}

	// At this point g is 0 and d contains +/- the modular inverse
	// Optionally negate d, normalize to [0,modulus), and return it
	modInv32Normalize30(&d, f.v[length-1], modinfo)
	*x = d
}
