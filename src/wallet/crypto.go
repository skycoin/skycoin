package wallet

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/encrypt"
)

// secrets key name
const (
	secretSeed     = "seed"
	secretLastSeed = "lastSeed"
)

type cryptor interface {
	Encrypt(data, password []byte) ([]byte, error)
	Decrypt(data, password []byte) ([]byte, error)
}

// CryptoType represents the type of crypto name
type CryptoType string

// StrToCryptoType converts string to CryptoType
func StrToCryptoType(s string) (CryptoType, error) {
	switch CryptoType(s) {
	case CryptoTypeSha256Xor:
		return CryptoTypeSha256Xor, nil
	case CryptoTypeScryptChacha20poly1305:
		return CryptoTypeScryptChacha20poly1305, nil
	default:
		return "", errors.New("unknow crypto type")
	}
}

// Crypto types
const (
	CryptoTypeSha256Xor              = CryptoType("sha256-xor")
	CryptoTypeScryptChacha20poly1305 = CryptoType("scrypt-chacha20poly1305")
)

// Scrypt paramenters
// scryptN: scrypt N paramenter. 1<<20 is the recommended value for file encryption, it takes about 3 seconds in 2.9 GHz Intel core i7.
// scryptR: scrypt r paramenter. Cache line size have not significantly increased since 2009, 8 should still be optimal for r.
// scryptP: scrypt p paramenter. The parallel difficulty, 1 is still optimal.
// scryptKeyLen: The length of returned byte slice that can be used as cryptographic key.
var (
	scryptN      = 1 << 20
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
)

// cryptoTable records all supported wallet crypto methods
// If want to support new crypto methods, register here.
var cryptoTable = map[CryptoType]cryptor{
	CryptoTypeSha256Xor:              encrypt.NewSha256Xor(),
	CryptoTypeScryptChacha20poly1305: encrypt.NewScryptChacha20poly1305(scryptN, scryptR, scryptP, scryptKeyLen),
}

// ErrAuthenticationFailed wraps the error of decryption.
type ErrAuthenticationFailed struct {
	err error
}

func (e ErrAuthenticationFailed) Error() string {
	return e.err.Error()
}

// getCrypto gets crypto of given type
func getCrypto(cryptoType CryptoType) (cryptor, error) {
	c, ok := cryptoTable[cryptoType]
	if !ok {
		return nil, fmt.Errorf("can not find crypto %v in crypto table", cryptoType)
	}

	return c, nil
}

type secrets map[string]string

func (s secrets) get(key string, v interface{}) error {
	d, ok := s[key]
	if !ok {
		return fmt.Errorf("secret %v doesn't exist", key)
	}

	return json.Unmarshal([]byte(d), v)
}

func (s secrets) set(key string, v interface{}) error {
	d, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s[key] = string(d)
	return nil
}

func (s secrets) serialize() ([]byte, error) {
	return json.Marshal(s)
}

func (s secrets) deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}

func (s secrets) erase() {
	for k := range s {
		s[k] = ""
		delete(s, k)
	}
}
