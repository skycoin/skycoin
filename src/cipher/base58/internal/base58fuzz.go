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
	decodeErr := encodeDecode(b)
	encodeErr := decodeEncode(b)

	if decodeErr == nil || encodeErr == nil {
		return 1
	}

	return 0
}

func decodeEncode(b []byte) error {
	x, err := base58.Decode(string(b))
	if err != nil {
		if x != nil {
			panic("x != nil on error number 2")
		}
	} else {
		s := base58.Encode(x)

		if s != string(b) {
			panic("encoded strings are not equal")
		}
	}

	return err
}

func encodeDecode(b []byte) error {
	s := base58.Encode(b)

	x, err := base58.Decode(s)
	if err != nil {
		if x != nil {
			panic("x != nil on error")
		}
	} else {
		if !bytes.Equal(b, x) {
			panic("decoded bytes are not equal")
		}
	}

	return err
}
