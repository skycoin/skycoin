package cipher

import (
	"github.com/tang0th/go-chacha20"
)

//32 byte input key
func ChaCha20Encrypt(in []byte, Key []byte) []byte {

	if len(Key != 32) {
		log.Panic("Key is 32 bytes")
	}
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20.XORKeyStream(out, in, []byte("nonce123"), Key[:])
	return out
}

//32 byte input key
func ChaCha20Decrypt(in []byte, Key []byte) []byte {

	if len(Key != 32) {
		log.Panic("Key is 32 bytes")
	}
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20.XORKeyStream(out, in, []byte("nonce123"), self.key[:])
	return out
}
