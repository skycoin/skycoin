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

// NewEntryFromReadable creates WalletEntry base one ReadableWalletEntry
func NewEntryFromReadable(w *ReadableEntry) (*Entry, error) {
	if w.Secret == "" {
		return nil, errors.New("secret field is empty")
	}

	s, err := cipher.SecKeyFromHex(w.Secret)
	if err != nil {
		return nil, err
	}

	a := cipher.AddressFromSecKey(s)
	if w.Address != "" {
		if a.String() != w.Address {
			return nil, errors.New("address does not match the secret")
		}
	}

	return &Entry{
		Address: a,
		Public:  cipher.PubKeyFromSecKey(s),
		Secret:  s,
	}, nil
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
