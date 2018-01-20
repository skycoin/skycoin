package wallet

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/scryptChacha20poly1305"
	"github.com/skycoin/skycoin/src/cipher/sha256xor"
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

// Scrypt paraments
//
const (
	scryptN      = 1 << 20
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32
)

// cryptoTable records all supported wallet crypto methods
// If want to support new crypto methods, register here.
var cryptoTable = map[CryptoType]cryptor{
	CryptoTypeSha256Xor:              sha256xor.New(),
	CryptoTypeScryptChacha20poly1305: scryptChacha20poly1305.New(scryptN, scryptR, scryptP, scryptKeyLen),
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
