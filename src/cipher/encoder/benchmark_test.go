package encoder

import "testing"

// benchmarkExample is the same struct used in https://github.com/gz-c/gosercomp benchmarks
type benchmarkExample struct {
	ID     int32
	Name   string
	Colors []string
}

var benchmarkExampleObj = benchmarkExample{
	ID:     1,
	Name:   "Reds",
	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
}

func BenchmarkDeserializeRaw(b *testing.B) {
	byt := Serialize(benchmarkExampleObj)
	result := &benchmarkExample{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeserializeRaw(byt, result) //nolint:errcheck
	}
}

func BenchmarkSerialize(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Serialize(&benchmarkExampleObj)
	}
}
