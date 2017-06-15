package cipher

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher/chacha20"
)

// Chacha20Encrypt encrypt the data in chacha20
func Chacha20Encrypt(data []byte, key []byte, nonce []byte) (d []byte, err error) {
	if len(key) != 32 {
		return []byte{}, errors.New("Key is 32 bytes")
	}
	e := make([]byte, len(data))
	c, err := chacha20.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}

// Chacha20Decrypt decrypt the data in chacha20
func Chacha20Decrypt(data []byte, key []byte, nonce []byte) (d []byte, err error) {
	if len(key) != 32 {
		return []byte{}, errors.New("Key is 32 bytes")
	}
	e := make([]byte, len(data))
	c, err := chacha20.New(key, nonce)
	if err != nil {
		return []byte{}, err
	}
	c.XORKeyStream(e, data)
	return e, nil
}
