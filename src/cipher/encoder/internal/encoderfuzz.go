package encoderfuzz

// To use the fuzzer:
// Follow the install instructions from https://github.com/dvyukov/go-fuzz
// Then, from the repo root,
// $ go-fuzz-build github.com/skycoin/skycoin/src/cipher/base58/internal
// This creates a file encoderfuzz-fuzz.zip
// Then,
// $ go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/base58/internal
// New corpus and crash objects will be put in src/cipher/base58/internal

// Fuzz is the entrypoint for go-fuzz
func Fuzz(b []byte) int {
	buf := SerializeString(string(b))
	if buf == nil {
		panic("SerializeString buf == nil")
	}

	s, x, err := DeserializeString(b, 8)
	if len(s) != x+4 {
		panic("DeserializeString len(s) != x + 4")
	}
	if err != nil {
		return 0
	}

	return 1
}
