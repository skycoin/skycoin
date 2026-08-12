// Package mobile provides mobile API bindings for Skycoin
package mobile

import (
	"encoding/hex"
	"encoding/json"

	"github.com/skycoin/skycoin/src/cipher/bip39"

	"github.com/skycoin/skycoin/src/skycoin-lite/liteclient"
)

// GetAddresses generates num addresses from seed using
// the Skycoin deterministic address generator
//
// This used to recover a panic from liteclient into err. The deferred function
// ran on the way out of every call, so a successful one whose json.Marshal had
// failed left with err set back to nil and an empty string reported as a
// result. liteclient returns errors now, so there is nothing to recover.
func GetAddresses(seed string, num int) (string, error) {
	hexSeed := hex.EncodeToString([]byte(seed))

	addresses, err := liteclient.GenerateAddresses(hexSeed, num)
	if err != nil {
		return "", err
	}

	byteaddr, err := json.Marshal(addresses) //nolint:gosec // this API deliberately returns key material to its caller
	if err != nil {
		return "", err
	}

	return string(byteaddr), nil
}

// PrepareTransaction receives inputs and outputs and returns a signed transaction
// inputsBody and outputsBody are JSONified arrays of TransactionInput and TransactionOutput.
func PrepareTransaction(inputsBody string, outputsBody string) (string, error) {
	return liteclient.PrepareTransaction(inputsBody, outputsBody)
}

// NewWordSeed wraps bip39.NewDefaultMnemonic
func NewWordSeed() (string, error) {
	// This panicked on error despite returning one, so a caller that checked
	// err — the whole point of the signature — was killed instead.
	return bip39.NewDefaultMnemonic()
}
