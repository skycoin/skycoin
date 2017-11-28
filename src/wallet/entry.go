package wallet

import (
	"github.com/skycoin/skycoin/src/cipher"
)

// Entry represents the wallet entry
type Entry struct {
	Address         cipher.Address
	Public          cipher.PubKey
	Secret          cipher.SecKey
	EncryptedSeckey string
}

// Verify checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
// func (we *Entry) Verify(password []byte) error {
// 	var seckey cipher.SecKey
// 	if password == nil {

// 	}
// 	ds, err := decrypt(we.Secret, password)
// 	if err != nil {
// 		return err
// 	}

// 	copy(seckey[:], ds[:])

// 	if cipher.PubKeyFromSecKey(seckey) != we.Public {
// 		return errors.New("invalid public key for secret key")
// 	}
// 	return we.VerifyPublic()
// }

// VerifyPublic checks that the public key is associated with the address
// func (we *Entry) VerifyPublic() error {
// 	if err := we.Public.Verify(); err != nil {
// 		return err
// 	}
// 	return we.Address.Verify(we.Public)
// }
