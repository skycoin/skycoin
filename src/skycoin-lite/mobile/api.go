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
func GetAddresses(seed string, num int) (addr string, err error) {
	defer func() {
		r := recover()
		err, _ = r.(error)
	}()
	hexSeed := hex.EncodeToString([]byte(seed))
	addresses := liteclient.GenerateAddresses(hexSeed, num)

	byteaddr, err := json.Marshal(addresses) //nolint:gosec // this API deliberately returns key material to its caller
	if err != nil {
		return "", err
	}
	addr = string(byteaddr)

	return addr, nil
}

// PrepareTransaction receives inputs and outputs and returns a signed transaction
// inputsBody and outputsBody are JSONified arrays of TransactionInput and TransactionOutput.
func PrepareTransaction(inputsBody string, outputsBody string) (tx string, err error) {

	defer func() {
		r := recover()
		err, _ = r.(error)
	}()

	tx = liteclient.PrepareTransaction(inputsBody, outputsBody)

	return
}

// NewWordSeed wraps bip39.NewDefaultMnemonic
func NewWordSeed() (string, error) {
	seed, err := bip39.NewDefaultMnemonic()
	if err != nil {
		panic(err)
	}

	return seed, nil
}
