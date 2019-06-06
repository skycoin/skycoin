package wallet

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
)

// Entry represents the wallet entry
type Entry struct {
	Address cipher.Addresser
	Public  cipher.PubKey
	Secret  cipher.SecKey
}

// SkycoinAddress returns the Skycoin address of an entry. Panics if Address is not a Skycoin address
func (we *Entry) SkycoinAddress() cipher.Address {
	return we.Address.(cipher.Address)
}

// BitcoinAddress returns the Skycoin address of an entry. Panics if Address is not a Bitcoin address
func (we *Entry) BitcoinAddress() cipher.BitcoinAddress {
	return we.Address.(cipher.BitcoinAddress)
}

// Verify checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (we *Entry) Verify() error {
	pk, err := cipher.PubKeyFromSecKey(we.Secret)
	if err != nil {
		return err
	}

	if pk != we.Public {
		return errors.New("invalid public key for secret key")
	}

	return we.VerifyPublic()
}

// VerifyPublic checks that the public key is associated with the address
func (we *Entry) VerifyPublic() error {
	if err := we.Public.Verify(); err != nil {
		return err
	}
	return we.Address.Verify(we.Public)
}
