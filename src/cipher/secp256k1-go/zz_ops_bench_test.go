package secp256k1

import (
	"crypto/sha256"
	"testing"
)

// BenchmarkOpsAllocs measures the operations a skywire deployment host runs
// most: CXO root checks (VerifySignature, RecoverPubkey), noise handshakes
// (ECDH) and noise key generation (GenerateKeyPair).
func BenchmarkOpsAllocs(b *testing.B) {
	pub, sec := GenerateDeterministicKeyPair([]byte("bench"))
	msg := sha256.Sum256([]byte("message"))
	sig := Sign(msg[:], sec)
	b.Run("VerifySignature", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if VerifySignature(msg[:], sig, pub) != 1 {
				b.Fatal("verify failed")
			}
		}
	})
	b.Run("RecoverPubkey", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if RecoverPubkey(msg[:], sig) == nil {
				b.Fatal("recover failed")
			}
		}
	})
	b.Run("ECDH", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = ECDH(pub, sec)
		}
	})
	b.Run("Sign", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_ = Sign(msg[:], sec)
		}
	})
	b.Run("GenerateKeyPair", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			_, _ = GenerateKeyPair()
		}
	})
}
