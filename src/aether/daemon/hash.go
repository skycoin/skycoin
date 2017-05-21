package daemon

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
)

var (
	sha256Hash = sha256.New()
)

// SHA256

// SHA256 sha256 type
type SHA256 [32]byte

// Set sets value
func (g *SHA256) Set(b []byte) {
	if len(b) != 32 {
		log.Panic("Invalid sha256 length")
	}
	copy(g[:], b[:])
}

// SumSHA256 sum sha256
func SumSHA256(b []byte) SHA256 {
	sha256Hash.Reset()
	sha256Hash.Write(b)
	sum := sha256Hash.Sum(nil)
	h := SHA256{}
	h.Set(sum)
	return h
}

// Hex returns sha256 hex string
func (g SHA256) Hex() string {
	return hex.EncodeToString(g[:])
}
