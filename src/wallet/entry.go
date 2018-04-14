package wallet

import (
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
)

// Entry represents the wallet entry
type Entry struct {
	Address cipher.Address
	Public  cipher.PubKey
	Secret  cipher.SecKey
}

// Verify checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (we *Entry) Verify() error {
	if cipher.PubKeyFromSecKey(we.Secret) != we.Public {
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
