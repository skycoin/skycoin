package encoderfuzz

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// To use the fuzzer:
// Follow the install instructions from https://github.com/dvyukov/go-fuzz
// Then, from the repo root,
// $ go-fuzz-build github.com/skycoin/skycoin/src/cipher/encoder/internal
// This creates a file encoderfuzz-fuzz.zip
// Then,
// $ go-fuzz -bin=encoderfuzz-fuzz.zip -workdir=src/cipher/encoder/internal
// New corpus and crash objects will be put in src/cipher/encoder/internal

type thing struct {
	X uint8 `enc:"-"`
	A int8
	B uint32
	T innerThing
	M map[uint64]int64
	F uint16
	G int16
	H []innerThing
	Z []byte `enc:",omitempty"`
}

type innerThing struct {
	A [2]byte
	C string `enc:",maxlen=128"`
}

// Fuzz is the entrypoint for go-fuzz
func Fuzz(b []byte) int {
	buf := encoder.SerializeString(string(b))
	if buf == nil {
		panic("SerializeString buf == nil")
	}

	s, x, errA := encoder.DeserializeString(b, 8)
	if errA == nil {
		if x != uint64(len(s)+4) {
			panic(fmt.Sprintf("DeserializeString x != len(s) + 4 (%d = %d, %q)", x, len(s)+4, s))
		}
	}

	var v uint32
	n, errB := encoder.DeserializeAtomic(b, &v)
	if errB == nil {
		if n != 4 {
			panic("DeserializeAtomic uint32 n bytes read is not 4")
		}
	}

	var t thing
	errC := encoder.DeserializeRawExact(b, &t)

	if errA == nil || errB == nil || errC == nil {
		return 1
	}

	return 0
}

// Uncomment and change package name to "main", then "go run" this to write a serialized "thing"
// func main() {
// 	x := thing{
// 		A: 12,
// 		B: 0xFF33AA01,
// 		T: innerThing{
// 			A: [2]byte{0x40, 0xF7},
// 			C: "foo",
// 		},
// 		// M: map[uint64]int64{0x00FF134444: -1234567},
// 		F: 0xAA11,
// 		G: -1000,
// 		H: []innerThing{{
// 			A: [2]byte{0x01, 0x02},
// 			C: "",
// 		}},
// 	}

// 	byt := encoder.Serialize(x)

// 	f, err := os.Create("thing3.serialized")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	f.Write(byt)
// }
