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
