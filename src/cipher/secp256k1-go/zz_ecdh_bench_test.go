package secp256k1

import "testing"

// BenchmarkECDHAllocs measures per-ECDH allocations — the exact operation the
// dmsg Noise handshake runs per stream. Option A (stack precomp tables in
// ECmult) should drop allocs/op sharply.
func BenchmarkECDHAllocs(b *testing.B) {
	pub, sec := GenerateKeyPair()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ECDH(pub, sec)
	}
}
