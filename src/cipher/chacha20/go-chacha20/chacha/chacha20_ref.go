// +build !cgo

package chacha

// core applies the ChaCha20 core function to 16-uint32 input matrix in, and puts
// the result into 64-byte array out.
//
// TODO: Use regs rather than an array
func core(in *[16]uint32, out *[64]byte, rounds int) {
	var x [16]uint32
	copy(x[:], in[:])

	// Applies 8 quarterrounds round / 2 times.
	for i := rounds; i > 0; i -= 2 {
		quarterround(&x, 0, 4, 8, 12)
		quarterround(&x, 1, 5, 9, 13)
		quarterround(&x, 2, 6, 10, 14)
		quarterround(&x, 3, 7, 11, 15)
		quarterround(&x, 0, 5, 10, 15)
		quarterround(&x, 1, 6, 11, 12)
		quarterround(&x, 2, 7, 8, 13)
		quarterround(&x, 3, 4, 9, 14)
	}

	// Add the block and input together. Then convert the 16-byte uint32 block to the
	// 64-byte output array.
	for i, j := 0, 0; i < 16; i++ {
		x[i] += in[i]
		out[j+0] = byte(x[i])
		out[j+1] = byte(x[i] >> 8)
		out[j+2] = byte(x[i] >> 16)
		out[j+3] = byte(x[i] >> 24)
		j += 4
	}
}

// quarterround applies the ChaCha20 quarterround on the given a, b, c and d
// elements of the array.
//
// The ChaCha quarterround is:
// a += b; d ^= a; d <<<= 16;
// c += d; b ^= c; b <<<= 12;
// a += b; d ^= a; d <<<= 8;
// c += d; b ^= c; b <<<= 7;
//
// TODO: Inline this function into the ChaCha core
// TODO: Use regs instead of arrays
func quarterround(x *[16]uint32, a, b, c, d int) {
	x[a] = x[a] + x[b]
	x[d] ^= x[a]
	x[d] = x[d]<<16 | x[d]>>(32-16)

	x[c] = x[c] + x[d]
	x[b] ^= x[c]
	x[b] = x[b]<<12 | x[b]>>(32-12)

	x[a] = x[a] + x[b]
	x[d] ^= x[a]
	x[d] = x[d]<<8 | x[d]>>(32-8)

	x[c] = x[c] + x[d]
	x[b] ^= x[c]
	x[b] = x[b]<<7 | x[b]>>(32-7)
}

// matrixSetup sets up the 16-uint32 ChaCha20 matrix with a 32-byte key k, an 8-byte
// iv i and a 16-byte constant
func matrixSetup(out *[16]uint32, k *[32]byte, i *[8]byte, c *[16]byte) {
	out[0] = uint32(c[0]) | uint32(c[1])<<8 | uint32(c[2])<<16 | uint32(c[3])<<24
	out[1] = uint32(c[4]) | uint32(c[5])<<8 | uint32(c[6])<<16 | uint32(c[7])<<24
	out[2] = uint32(c[8]) | uint32(c[9])<<8 | uint32(c[10])<<16 | uint32(c[11])<<24
	out[3] = uint32(c[12]) | uint32(c[13])<<8 | uint32(c[14])<<16 | uint32(c[15])<<24

	out[4] = uint32(k[0]) | uint32(k[1])<<8 | uint32(k[2])<<16 | uint32(k[3])<<24
	out[5] = uint32(k[4]) | uint32(k[5])<<8 | uint32(k[6])<<16 | uint32(k[7])<<24
	out[6] = uint32(k[8]) | uint32(k[9])<<8 | uint32(k[10])<<16 | uint32(k[11])<<24
	out[7] = uint32(k[12]) | uint32(k[13])<<8 | uint32(k[14])<<16 | uint32(k[15])<<24
	out[8] = uint32(k[16]) | uint32(k[17])<<8 | uint32(k[18])<<16 | uint32(k[19])<<24
	out[9] = uint32(k[20]) | uint32(k[21])<<8 | uint32(k[22])<<16 | uint32(k[23])<<24
	out[10] = uint32(k[24]) | uint32(k[25])<<8 | uint32(k[26])<<16 | uint32(k[27])<<24
	out[11] = uint32(k[28]) | uint32(k[29])<<8 | uint32(k[30])<<16 | uint32(k[31])<<24

	out[12] = 0
	out[13] = 0
	out[14] = uint32(i[0]) | uint32(i[1])<<8 | uint32(i[2])<<16 | uint32(i[3])<<24
	out[15] = uint32(i[4]) | uint32(i[5])<<8 | uint32(i[6])<<16 | uint32(i[7])<<24
}

// XORKeyStream crypts bytes from in to out using the given key, initialisation vector,
// constant and number of ChaCha20 rounds to perform.
//
// In and out may be the same slice but otherwise should not overlap. Counter
// contains the raw salsa20 counter bytes (both nonce and block counter).
func XORKeyStream(out, in []byte, iv *[8]byte, constant *[16]byte, key *[32]byte, rounds int) {
	var block [64]byte
	var matrix [16]uint32

	// Set up the matrix
	matrixSetup(&matrix, key, iv, constant)

	// For whole blocks
	for len(in) >= 64 {
		core(&matrix, &block, rounds)
		for i, x := range block {
			out[i] = in[i] ^ x
		}
		matrix[12] += 1
		if matrix[12] == 0 {
			matrix[13] += 1
		}
		in = in[64:]
		out = out[64:]
	}

	// Last remaining block
	if len(in) > 0 {
		core(&matrix, &block, rounds)
		for i, v := range in {
			out[i] = v ^ block[i]
		}
	}
}
