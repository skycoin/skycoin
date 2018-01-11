package wallet

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/sha256xor"
)

// Crypto is the interface that provides methods for encryption and decryption
type Crypto interface {
	Encrypt(data, password []byte) ([]byte, error)
	Decrypt(data, password []byte) ([]byte, error)
}

// CryptoType is the data type of crypto name
type CryptoType string

// Crypto types
const (
	CryptoTypeSha256Xor = CryptoType("sha256-xor")
)

// cryptoTable records all supported wallet crypto methods
// If want to support new crypto methods, register them here.
var cryptoTable = map[CryptoType]Crypto{
	CryptoTypeSha256Xor: sha256xor.New(),
}

// getCrypto gets crypto of givn type
func getCrypto(cryptoType CryptoType) (Crypto, error) {
	c, ok := cryptoTable[cryptoType]
	if !ok {
		return nil, fmt.Errorf("could not find crypto %v in crypto table", cryptoType)
	}

	return c, nil
}
