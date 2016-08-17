package chacha20

import "github.com/tang0th/go-chacha20/chacha"

// XORKeyStream crypts bytes from in to out using the given key and nonce. It
// performs the 20 round chacha cipher operation. In and out may be the same
// slice but otherwise should not overlap. Nonce must be 8 bytes long. Key
// must be either 10, 16 or 32 byte long (32 bytes is recommended for security
// purposes).
func XORKeyStream(out, in []byte, nonce []byte, key []byte) {
	xorkeystream(out, in, nonce, key, 20)
}

// XORKeyStream12 crypts bytes from in to out using the given key and nonce. It
// performs the 12 round chacha cipher operation. In and out may be the same
// slice but otherwise should not overlap. Nonce must be 8 bytes long. Key
// must be either 10, 16 or 32 byte long (32 bytes is recommended for security
// purposes).
func XORKeyStream12(out, in []byte, nonce []byte, key []byte) {
	xorkeystream(out, in, nonce, key, 12)
}

// XORKeyStream8 crypts bytes from in to out using the given key and nonce. It
// performs the 8 round chacha cipher operation. In and out may be the same
// slice but otherwise should not overlap. Nonce must be 8 bytes long. Key
// must be either 10, 16 or 32 byte long (32 bytes is recommended for security
// purposes).
func XORKeyStream8(out, in []byte, nonce []byte, key []byte) {
	xorkeystream(out, in, nonce, key, 8)
}

// xorkeystream crypts bytes from in to out using the given key and nonce and
// number of rounds. It selects the right constants, elongates the key (if
// necessary) and in the future may even call HChacha20 (once it is proved secure).
// We then call the underlying chahcha.XORKeyStream function.
func xorkeystream(out, in []byte, nonce []byte, k []byte, rounds int) {
	if len(out) < len(in) {
		in = in[:len(out)]
	}

	var subNonce [8]byte
	var key [32]byte
	var constant *[16]byte

	switch len(k) {
	case 32:
		copy(key[:], k)
		constant = &chacha.Sigma
	case 16:
		copy(key[:16], k)
		copy(key[16:], k)
		constant = &chacha.Tau
	case 10:
		copy(key[:16], k)
		copy(key[16:], k)
		constant = &chacha.Upsilon
	default:
		panic("chacha20: key must be 32, 16 or 10 bytes")
	}

	if len(nonce) == 8 {
		copy(subNonce[:], nonce[:])
	} else {
		panic("chacha20: nonce must be 8 bytes")
	}

	chacha.XORKeyStream(out, in, &subNonce, constant, &key, rounds)
}
