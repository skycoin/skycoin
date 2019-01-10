package base58fuzz

import (
	"bytes"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

// To use the fuzzer:
// Follow the install instructions from https://github.com/dvyukov/go-fuzz
// Then, from the repo root,
// $ go-fuzz-build github.com/skycoin/skycoin/src/cipher/base58/internal
// This creates a file base58fuzz-fuzz.zip
// Then,
// $ go-fuzz -bin=base58fuzz-fuzz.zip -workdir=src/cipher/base58/internal
// New corpus and crash objects will be put in src/cipher/base58/internal

// Fuzz is the entrypoint for go-fuzz
func Fuzz(b []byte) int {
	s := base58.Encode(b)

	x, err := base58.Decode(s)
	if err != nil {
		if x != nil {
			panic("x != nil on error")
		}
		return 0
	}

	if !bytes.Equal(b, x) {
		panic("decoded bytes are not equal")
	}

	x, err = base58.Decode(string(b))
	if err != nil {
		if x != nil {
			panic("x != nil on error number 2")
		}
		return 0
	}

	s = base58.Encode(x)

	if s != string(b) {
		panic("encoded strings are not equal")
	}

	return 1
}
