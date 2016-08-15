package cipher

import (
	"errors"
	chacha20_1 "github.com/codahale/chacha20"
	chacha20_0 "github.com/tang0th/go-chacha20"
	"log"
)

//32 byte input key
func ChaCha20Encrypt(in []byte, Key []byte) []byte {

	if len(Key) != 32 {
		log.Panic("Key is 32 bytes")
	}
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20_0.XORKeyStream(out, in, []byte("nonce123"), Key[:])
	return out
}

//32 byte input key
func ChaCha20Decrypt(in []byte, Key []byte) []byte {

	if len(Key) != 32 {
		log.Panic("Key is 32 bytes")
	}
	out := make([]byte, len(in))
	//TODO, using fixed nonce
	chacha20_0.XORKeyStream(out, in, []byte("nonce123"), Key)
	return out
}

/*
Duplicate
*/

func Chacha20Encrypt(data []byte, pubkey PubKey, seckey SecKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("encrypt faild")
		}
	}()

	key := ECDH(pubkey, seckey)
	e := make([]byte, len(data))
	c, err := chacha20_1.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}

func Chacha20Decrypt(data []byte, pubkey PubKey, seckey SecKey, nonce []byte) (d []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("decrypt faild")
		}
	}()

	key := ECDH(pubkey, seckey)
	e := make([]byte, len(data))
	c, err := chacha20_1.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}
