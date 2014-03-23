package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"log"
)

var (
	sha256Hash hash.Hash = sha256.New()
)

// SHA256

type SHA256 [32]byte

func (g *SHA256) Set(b []byte) {
	if len(b) != 32 {
		log.Panic("Invalid sha256 length")
	}
	copy(g[:], b[:])
}

func SumSHA256(b []byte) SHA256 {
	sha256Hash.Reset()
	sha256Hash.Write(b)
	sum := sha256Hash.Sum(nil)
	h := SHA256{}
	h.Set(sum)
	return h
}

func (g SHA256) Hex() string {
	return hex.EncodeToString(g[:])
}
