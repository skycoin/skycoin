package wallet

import (
	"errors"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
)

type WalletEntries map[cipher.Address]WalletEntry

func (self WalletEntries) ToArray() []WalletEntry {
	e := make([]WalletEntry, len(self))
	i := 0
	for _, we := range self {
		e[i] = we
		i++
	}
	return e
}

type WalletEntry struct {
	Address cipher.Address
	Public  cipher.PubKey
	Secret  cipher.SecKey
}

func NewWalletEntryFromKeypair(pub cipher.PubKey, sec cipher.SecKey) WalletEntry {
	return WalletEntry{
		Public:  pub,
		Secret:  sec,
		Address: cipher.AddressFromPubKey(pub),
	}
}

func NewWalletEntry() WalletEntry {
	pub, sec := cipher.GenerateKeyPair()
	return NewWalletEntryFromKeypair(pub, sec)
}

func WalletEntryFromReadable(w *ReadableWalletEntry) WalletEntry {
	// SimpleWallet entries are shared as a form of identification, the secret key
	// is not required
	// TODO -- fix lib/base58 to not panic on invalid input -- should
	// return error, so we can detect a broken wallet.
	if w.Address == "" {
		log.Panic("ReadableWalletEntry has no Address")
	}
	var s cipher.SecKey
	if w.Secret != "" {
		s = cipher.MustSecKeyFromHex(w.Secret)
	}
	return WalletEntry{
		Address: cipher.MustDecodeBase58Address(w.Address),
		Public:  cipher.MustPubKeyFromHex(w.Public),
		Secret:  s,
	}
}

// Loads a WalletEntry from filename, where the file contains a
// ReadableWalletEntry
func LoadWalletEntry(filename string) (WalletEntry, error) {
	w, err := LoadReadableWalletEntry(filename)
	if err != nil {
		return WalletEntry{}, err
	} else {
		return WalletEntryFromReadable(&w), nil
	}
}

// Loads a WalletEntry from filename but panics is unable to load or contents
// are invalid
func MustLoadWalletEntry(filename string) WalletEntry {
	keys, err := LoadWalletEntry(filename)
	if err != nil {
		log.Panicf("Failed to load wallet entry: %v", err)
	}
	if err := keys.Verify(); err != nil {
		log.Panicf("Invalid wallet entry: %v", err)
	}
	return keys
}

// Checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (self *WalletEntry) Verify() error {
	if cipher.PubKeyFromSecKey(self.Secret) != self.Public {
		return errors.New("Invalid public key for secret key")
	}
	return self.VerifyPublic()
}

// Checks that the public key is associated with the address
func (self *WalletEntry) VerifyPublic() error {
	if err := self.Public.Verify(); err != nil {
		return err
	} else {
		return self.Address.Verify(self.Public)
	}
}
