package wallet

import (
	"encoding/base64"

	"github.com/skycoin/skycoin/src/cipher"
)

// cipher.Encrypt the data, then encode the result into base64 string
func encrypt(data []byte, password []byte) (string, error) {
	encData, err := cipher.Encrypt(data, password)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encData), nil
}

// base64 decode the string, then cipher.Decrypt the data
func decrypt(data string, password []byte) ([]byte, error) {
	base64DecodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return cipher.Decrypt(base64DecodedData, password)
}
