package encoder

import "testing"

// benchmarkExample is the same struct used in https://github.com/gz-c/gosercomp benchmarks
type benchmarkExample struct {
	ID     int32
	Name   string
	Colors []string
}

func BenchmarkDeserializeRaw(b *testing.B) {
	obj := benchmarkExample{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	byt := Serialize(obj)
	result := &benchmarkExample{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DeserializeRaw(byt, result)
	}
}

func BenchmarkSerialize(b *testing.B) {
	obj := benchmarkExample{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Serialize(&obj)
	}
}
