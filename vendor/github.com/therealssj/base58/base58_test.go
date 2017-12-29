package base58

import (
	"crypto/rand"
	"testing"
	"fmt"
)

type testValues struct {
	dec []byte
	enc string
}

var n = 5000000
var testPairs = make([]testValues, 0, n)

func initTestPairs() {
	if len(testPairs) > 0 {
		return
	}
	// pre-make the test pairs, so it doesn't take up benchmark time...
	data := make([]byte, 32)
	for i := 0; i < n; i++ {
		rand.Read(data)
		b58string, err := Encode(data)
		if err != nil {
			fmt.Sprintf("error in test setup: %v", err)

		}
		testPairs = append(testPairs, testValues{dec: data, enc: b58string})
	}
}

func BenchmarkBase58Encoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Encode([]byte(testPairs[i].dec))
	}
}

func BenchmarkBase58Decoding(b *testing.B) {
	initTestPairs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Decode(testPairs[i].enc)
	}
}

