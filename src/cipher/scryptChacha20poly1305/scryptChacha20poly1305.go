package scryptChacha20poly1305

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"

	"github.com/skycoin/skycoin/src/cipher"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
)

const (
	metaLengthSize = 2 // meta data length field size in bytes
	saltSize       = 8 // salt bytes number
)

// ScryptChacha20poly1305 provides methods for encryption/decryption with scrypt and chacha20poly1305
type ScryptChacha20poly1305 struct {
	n      int
	r      int
	p      int
	keyLen int
}

// New creates a ScryptChacha20poly1305 instance
func New(N, r, p, keyLen int) *ScryptChacha20poly1305 {
	return &ScryptChacha20poly1305{
		n:      N,
		r:      r,
		p:      p,
		keyLen: keyLen,
	}
}

type meta struct {
	N      int
	R      int
	P      int
	KeyLen int
	Salt   []byte
	Nonce  []byte
}

// Encrypt encrypts data with password,
// 1. Scrypt derives the key from password
// 2. Chacha20poly1305 generates AEAD from the derived key
// 4. Puts scrypt paramenters, salt and nonce into metadata, json serialize it and get the serialized metadata length
// 5. AEAD.Seal encrypts the data, and use [length][metadata] as additional data
// 6. Final format: base64([[length][metadata]][ciphertext]), length is 2 bytes.
func (s *ScryptChacha20poly1305) Encrypt(data, password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, errors.New("missing password")
	}

	// Scyrpt derives key from password
	salt := cipher.RandByte(saltSize)
	dk, err := scrypt.Key(password, salt, s.n, s.r, s.p, s.keyLen)
	if err != nil {
		return nil, err
	}

	// Prepare metadata
	m := meta{s.n, s.r, s.p, s.keyLen, salt, cipher.RandByte(chacha20poly1305.NonceSize)}
	// json serialize the metadata
	ms, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	if len(ms) > math.MaxUint16 {
		return nil, errors.New("metadata length beyond the math.MaxUint16")
	}

	length := make([]byte, metaLengthSize)
	binary.LittleEndian.PutUint16(length, uint16(len(ms)))

	// Additional data for AEAD
	ad := append(length, ms...)
	aead, err := chacha20poly1305.New(dk)
	if err != nil {
		return nil, err
	}

	ciphertext := aead.Seal(nil, m.Nonce, data, ad)

	// Base64 encode the [[length][metadata]][ciphertext]
	rawData := append(ad, ciphertext...)
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(rawData)))
	enc.Encode(buf, rawData)
	return buf, nil
}

// Decrypt decrypts the data with password
// 1. Base64 decodes the data
// 2. Reads the first [metaLengthSize] bytes data to get the metadata length, and reads out the metadata.
// 3. Scrypt derives key from password and paramenters in metadata
// 4. Chacha20poly1305 geneates AEAD
// 5. AEAD decrypts ciphertext with nonce in metadata and [length][metadata] as additional data.
func (s *ScryptChacha20poly1305) Decrypt(data, password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, errors.New("missing password")
	}

	enc := base64.StdEncoding
	encData := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(encData, data)
	if err != nil {
		return nil, err
	}
	encData = encData[:n]

	length := binary.LittleEndian.Uint16(encData[:metaLengthSize])
	if length > math.MaxUint16 {
		return nil, errors.New("metadata length beyond math.MaxUint64")
	}

	if int(metaLengthSize+length) > len(encData) {
		return nil, errors.New("invalid metadata length")
	}

	var m meta
	if err := json.Unmarshal(encData[metaLengthSize:metaLengthSize+length], &m); err != nil {
		return nil, err
	}

	ad := encData[:metaLengthSize+length]
	// Scrypt derives key
	dk, err := scrypt.Key(password, m.Salt, m.N, m.R, m.P, m.KeyLen)
	if err != nil {
		return nil, err
	}

	// Geneates AEAD
	aead, err := chacha20poly1305.New(dk)
	if err != nil {
		return nil, err
	}

	return aead.Open(nil, m.Nonce, encData[metaLengthSize+length:], ad)
}
