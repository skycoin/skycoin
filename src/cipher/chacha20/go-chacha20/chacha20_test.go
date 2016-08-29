package chacha20

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"io"
	"testing"
)

func fromHex(s string) []byte {
	ret, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return ret
}

// chacha2020TestVectors was taken from the agl's chacha20-poly1305 proposal:
// https://tools.ietf.org/html/draft-agl-tls-chacha20poly1305-04
var chacha2020TestVectors = []struct {
	key        []byte
	iv         []byte
	ciphertext []byte
}{
	{
		fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		fromHex("0000000000000000"),
		fromHex("76b8e0ada0f13d90405d6ae55386bd28bdd219b8a08ded1aa836efcc8b770dc7da41597c5157488d7724e03fb8d84a376a43b8f41518a11cc387b669"),
	},
	{
		fromHex("0000000000000000000000000000000000000000000000000000000000000001"),
		fromHex("0000000000000000"),
		fromHex("4540f05a9f1fb296d7736e7b208e3c96eb4fe1834688d2604f450952ed432d41bbe2a0b6ea7566d2a5d1e7e20d42af2c53d792b1c43fea817e9ad275"),
	},
	{
		fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		fromHex("0000000000000001"),
		fromHex("de9cba7bf3d69ef5e786dc63973f653a0b49e015adbff7134fcb7df137821031e85a050278a7084527214f73efc7fa5b5277062eb7a0433e445f41e3"),
	},
	{
		fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
		fromHex("0100000000000000"),
		fromHex("ef3fdfd6c61578fbf5cf35bd3dd33b8009631634d21e42ac33960bd138e50d32111e4caf237ee53ca8ad6426194a88545ddc497a0b466e7d6bbdb004"),
	},
	{
		fromHex("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"),
		fromHex("0001020304050607"),
		fromHex("f798a189f195e66982105ffb640bb7757f579da31602fc93ec01ac56f85ac3c134a4547b733b46413042c9440049176905d3be59ea1c53f15916155c2be8241a38008b9a26bc35941e2444177c8ade6689de95264986d95889fb60e84629c9bd9a5acb1cc118be563eb9b3a4a472f82e09a7e778492b562ef7130e88dfe031c79db9d4f7c7a899151b9a475032b63fc385245fe054e3dd5a97a5f576fe064025d3ce042c566ab2c507b138db853e3d6959660996546cc9c4a6eafdc777c040d70eaf46f76dad3979e5c5360c3317166a1c894c94a371876a94df7628fe4eaaf2ccb27d5aaae0ad7ad0f9d4b6ad3b54098746d4524d38407a6deb"),
	},

	// These are taken from
	/*{
		fromHex("80000000000000000000000000000000"),
		fromHex("0000000000000000"),
		fromHex("beb1e81e0f747e43ee51922b3e87fb38d0163907b4ed49336032ab78b67c24579fe28f751bd3703e51d876c017faa43589e63593e03355a7d57b2366f30047c5"),
	},
	{
		fromHex("0053a6f94c9ff24598eb3e91e4378add"),
		fromHex("0d74db42a91077de"),
		fromHex("509b267e7266355fa2dc0a25c023fce47922d03dd9275423d7cb7118b2aedf220568854bf47920d6fc0fd10526cfe7f9de472835afc73c916b849e91eee1f529"),
	},
	{
		fromHex("0000200000000000000000000000000000000000000000000000000000000000"),
		fromHex("0000000000000000"),
		fromHex("653f4a18e3d27daf51f841a00b6c1a2bd2489852d4ae0711e1a4a32ad166fa6f881a2843238c7e17786ba5162bc019d573849c167668510ada2f62b4ff31ad04"),
	},
	{
		fromHex("0f62b5085bae0154a7fa4da0f34699ec3f92e5388bde3184d72a7dd02376c91c"),
		fromHex("288ff65dc42b92f9"),
		fromHex("db165814f66733b7a8e34d1ffc1234271256d3bf8d8da2166922e598acac70f412b3fe35a94190ad0ae2e8ec62134819ab61addcccfe99d867ca3d73183fa3fd"),
	},*/
}

func TestChacha2020(t *testing.T) {
	var inBuf, outBuf []byte
	var numBytes int

	for i, test := range chacha2020TestVectors {
		numBytes = len(test.ciphertext)
		if numBytes > len(inBuf) {
			inBuf = make([]byte, numBytes)
			outBuf = make([]byte, numBytes)
		}
		in := inBuf[:numBytes]
		out := outBuf[:numBytes]
		XORKeyStream(out, in, test.iv, test.key)

		if !bytes.Equal(out, test.ciphertext) {
			t.Errorf("#%d: bad result. Got: %v Expected: %v", i, out, test.ciphertext)
		}
	}
}

func BenchmarkChacha2020_1K(b *testing.B) {
	benchmarkChacha(b, XORKeyStream, 1024)
}

func BenchmarkChacha2012_1K(b *testing.B) {
	benchmarkChacha(b, XORKeyStream12, 1024)
}

func BenchmarkChacha208_1K(b *testing.B) {
	benchmarkChacha(b, XORKeyStream8, 1024)
}

func benchmarkChacha(b *testing.B, xor func(out, in []byte, nonce []byte, key []byte), dataLen int64) {
	b.StopTimer()
	var buff []byte = make([]byte, dataLen)
	var key []byte = make([]byte, 32)
	var iv []byte = make([]byte, 8)
	io.ReadFull(rand.Reader, key)
	io.ReadFull(rand.Reader, iv)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xor(buff, buff, iv, key)
	}
	b.SetBytes(dataLen)
}
