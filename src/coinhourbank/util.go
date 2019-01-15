package coinhourbank

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

// Out represents a transaction destination
type Out struct {
	Address string `json:"address"` // Destination address
	Coins   uint64 `json:"coins"`   // Amount of coins
	Hours   uint64 `json:"hours"`   // Amount of coin hours
}

// Transaction represents a transaction to join
type Transaction struct {
	FromAddress string   `json:"fromAddress"` // Address that UxOuts belongs to
	UxOuts      []string `json:"uxOuts"`      // UxOuts of transaction sender
	Outs        []Out    `json:"outs"`        // Transaction outs
}

// Input represents a transaction input that must be signed by input owner
type Input struct {
	FromAddress string     // Address of transaction sender
	UxOut       string     // UxOut HEX of transaction sender
	Sign        cipher.Sig // Signature of input
	InputIndex  int        // Index of this input in the transaction inputs array
}

// signAddressInputs signs all UxOuts from tx inputs (inputs) that corresponds to specified address by its secret key.
func signAddressInputs(address string, txInnerHash cipher.SHA256, secKeyHex string, inputs []CoinjoinInput) error {
	secKey, err := cipher.SecKeyFromHex(secKeyHex)
	if err != nil {
		return err
	}

	for i, in := range inputs {
		if in.FromAddress == address {
			h := cipher.AddSHA256(txInnerHash, cipher.MustSHA256FromHex(in.UxOut)) // hash to sign
			inputs[i].Sign, err = cipher.SignHash(h, secKey)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func modelToBytes(i interface{}) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	return b, json.NewEncoder(b).Encode(i)
}


// ArgumentValidationErr is the error of the wrong value of function argument
type ArgumentValidationErr struct {
	message string
}

// NewArgumentValidationErr returns ArgumentValidationErr
func NewArgumentValidationErr(argument string, message string) *ArgumentValidationErr {
	return &ArgumentValidationErr{
		message: fmt.Sprintf("[%s] %s", argument, message),
	}
}

func (e *ArgumentValidationErr) Error() string {
	return e.message
}
