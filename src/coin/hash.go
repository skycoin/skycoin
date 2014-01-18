package sb_coin

import (
	"crypto/sha256"
	"encoding/hex"
	//"fmt"
	"hash"
	"log"
)

import "lib/encoder"
import "lib/ripemd160"

/*
	SHA256
*/
var sha256_hash hash.Hash = sha256.New()

func Sha256_func(data []byte) [32]byte {
	sha256_hash.Reset()
	sha256_hash.Write(data)
	sum := sha256_hash.Sum(nil)
	var out [32]byte
	copy(out[0:32], sum[0:32])
	return out
}

type SHA256 struct {
	Value [32]byte
}

func SHA256sum(b []byte) SHA256 {
	var ret SHA256
	ret.Value = Sha256_func(b)
	return ret
}

//sums to two hashes together
func SHA256add(a1 SHA256, b1 SHA256) SHA256 {
	b := append(a1.Value[:], b1.Value[:]...)
	return SHA256sumn(b, 32*2)
}

//double SHA256
func DSHA256sum(b []byte) SHA256 {
	var ret SHA256 = SHA256{Value: Sha256_func(b)}
	return SHA256add(ret, ret)
}

//assert number of bytes as input
func SHA256sumn(b []byte, n int) SHA256 {
	if len(b) != n {
		log.Panic()
	}
	return SHA256{Value: Sha256_func(b)}
}

func (h1 SHA256) Xor(h2 SHA256) SHA256 {
	var h3 SHA256
	for i := 0; i < 32; i++ {
		h3.Value[i] = h1.Value[i] ^ h2.Value[i]
	}
	return h3
}

//compute root merkle tree hash of hash list
//pad input to power of 16
//group inputs hashes into groups of 16 and hash them down to single hash
//repeat until there is single hash in list
func Merkle(h0 []SHA256) SHA256 {
	//fmt.Printf("Merkle 0: len= %v \n", len(h0))
	if len(h0) == 0 {
		return SHA256{} //zero hash
	}
	np := 0
	for np = 1; np < len(h0); np *= 16 {
	}
	h1 := make([]SHA256, np)

	var th SHA256 = h0[0]

	var lh0 = len(h0)
	for i := 0; i < np-lh0; i++ {
		th = th.Xor(h0[1+i])
		h1[i+lh0] = th //pad to power of 16
	}

	for len(h1) != 1 {
		//fmt.Printf("Merkle 1: len= %v \n", len(h1))
		h2 := make([]SHA256, len(h1)/16)
		var h3 [16]SHA256
		for i := 0; i < len(h2); i++ {
			for j := 0; j < 16; j++ {
				h3[j] = h1[16*i+j]
			}
			h2[i] = SHA256sumn(encoder.Serialize(h3), 16*32)
		}
		h1 = h2
	}
	return h1[0]
}

func (g *SHA256) Set(b []byte) {
	if len(b) != 32 {
		log.Panic()
	}
	copy(g.Value[0:32], b[0:32])
}

/*
	Ripmd160
*/

var ripemd160_hash hash.Hash = ripemd160.New()

func Ripmd160_func(data []byte) [20]byte {
	ripemd160_hash.Reset()
	ripemd160_hash.Write(data)
	sum := ripemd160_hash.Sum(nil)
	var out [20]byte
	copy(out[0:20], sum[0:20])
	return out
}

func Hex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Panic()
		return nil
	}
	return b
}
