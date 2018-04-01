// uses scrypt for password key derivation, and chacha20poly1305 for
// encryption/decryption

package encrypt

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/chacha20poly1305"
	"github.com/skycoin/skycoin/src/cipher/scrypt"
)

const (
	scryptChacha20MetaLengthSize = 2  // meta data length field size in bytes
	scryptChacha20SaltSize       = 32 // salt bytes number
)

// Default scrypt paramenters
const (
	// ScryptN: scrypt N paramenter. 1<<20 is the recommended value for file encryption, it takes about 3 seconds in 2.9 GHz Intel core i7.
	ScryptN = 1 << 20
	// ScryptR: scrypt r paramenter. Cache line size have not significantly increased since 2009, 8 should still be optimal for r.
	ScryptR = 8
	// ScryptP: scrypt p paramenter. The parallel difficulty, 1 is still optimal.
	ScryptP = 1
	// ScryptKeyLen: The length of returned byte slice that can be used as cryptographic key.
	ScryptKeyLen = 32
)

// DefaultScryptChacha20poly1305 default ScryptChacha20poly1305 encryptor
var DefaultScryptChacha20poly1305 = ScryptChacha20poly1305{
	N:      ScryptN,
	R:      ScryptR,
	P:      ScryptP,
	KeyLen: ScryptKeyLen,
}

// ScryptChacha20poly1305 provides methods for encryption/decryption with scrypt and chacha20poly1305
type ScryptChacha20poly1305 struct {
	N      int
	R      int
	P      int
	KeyLen int
}

type meta struct {
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	KeyLen int    `json:"keyLen"`
	Salt   []byte `json:"salt"`
	Nonce  []byte `json:"nonce"`
}

// Encrypt encrypts data with password,
// 1. Scrypt derives the key from password
// 2. Chacha20poly1305 generates AEAD from the derived key
// 4. Puts scrypt paramenters, salt and nonce into metadata, json serialize it and get the serialized metadata length
// 5. AEAD.Seal encrypts the data, and use [length][metadata] as additional data
// 6. Final format: base64([[length][metadata]][ciphertext]), length is 2 bytes.
func (s ScryptChacha20poly1305) Encrypt(data, password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, errors.New("missing password")
	}

	// Scyrpt derives key from password
	salt := cipher.RandByte(scryptChacha20SaltSize)
	dk, err := scrypt.Key(password, salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return nil, err
	}

	// Prepare metadata
	m := meta{
		N:      s.N,
		R:      s.R,
		P:      s.P,
		KeyLen: s.KeyLen,
		Salt:   salt,
		Nonce:  cipher.RandByte(chacha20poly1305.NonceSize),
	}
	// json serialize the metadata
	ms, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	if len(ms) > math.MaxUint16 {
		return nil, errors.New("metadata length beyond the math.MaxUint16")
	}

	length := make([]byte, scryptChacha20MetaLengthSize)
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
func (s ScryptChacha20poly1305) Decrypt(data, password []byte) ([]byte, error) {
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

	length := binary.LittleEndian.Uint16(encData[:scryptChacha20MetaLengthSize])
	if int(scryptChacha20MetaLengthSize+length) > len(encData) {
		return nil, errors.New("invalid metadata length")
	}

	var m meta
	if err := json.Unmarshal(encData[scryptChacha20MetaLengthSize:scryptChacha20MetaLengthSize+length], &m); err != nil {
		return nil, err
	}

	ad := encData[:scryptChacha20MetaLengthSize+length]
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

	return aead.Open(nil, m.Nonce, encData[scryptChacha20MetaLengthSize+length:], ad)
}
