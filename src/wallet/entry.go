package wallet

import (
	"errors"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
)

// Entry represents the wallet entry
type Entry struct {
	Address cipher.Address
	Public  cipher.PubKey
	Secret  cipher.SecKey
}

// NewEntryFromKeypair creates wallet entry base on key pairs
func NewEntryFromKeypair(pub cipher.PubKey, sec cipher.SecKey) Entry {
	return Entry{
		Address: cipher.AddressFromPubKey(pub),
		Public:  pub,
		Secret:  sec,
	}
}

// NewEntry creates wallet entry
func NewEntry() Entry {
	pub, sec := cipher.GenerateKeyPair()
	return NewEntryFromKeypair(pub, sec)
}

// NewEntryFromReadable creates WalletEntry base one ReadableWalletEntry
func NewEntryFromReadable(w *ReadableEntry) Entry {
	// SimpleWallet entries are shared as a form of identification, the secret key
	// is not required
	// TODO -- fix lib/base58 to not panic on invalid input -- should
	// return error, so we can detect a broken wallet.

	if w.Address == "" {
		//log.Panic("ReadableWalletEntry has no Address")
	}
	var s cipher.SecKey
	if w.Secret != "" {
		s = cipher.MustSecKeyFromHex(w.Secret)
	}

	//regen from the private key
	//redundant/
	if w.Address == "" {
		addr := cipher.AddressFromSecKey(s)
		pub := cipher.PubKeyFromSecKey(s)

		return Entry{
			Address: addr,
			Public:  pub,
			Secret:  s,
		}
	}

	return Entry{
		Address: cipher.MustDecodeBase58Address(w.Address),
		Public:  cipher.MustPubKeyFromHex(w.Public),
		Secret:  s,
	}
}

// LoadEntry loads a WalletEntry from filename, where the file contains a
// ReadableWalletEntry
func LoadEntry(filename string) (Entry, error) {
	w, err := LoadReadableEntry(filename)
	if err != nil {
		return Entry{}, err
	}

	return NewEntryFromReadable(&w), nil
}

// MustLoadEntry loads a WalletEntry from filename but panics is unable to load or contents
// are invalid
func MustLoadEntry(filename string) Entry {
	keys, err := LoadEntry(filename)
	if err != nil {
		log.Panicf("Failed to load wallet entry: %v", err)
	}
	if err := keys.Verify(); err != nil {
		log.Panicf("Invalid wallet entry: %v", err)
	}
	return keys
}

// Verify checks that the public key is derivable from the secret key,
// and that the public key is associated with the address
func (we *Entry) Verify() error {
	if cipher.PubKeyFromSecKey(we.Secret) != we.Public {
		return errors.New("Invalid public key for secret key")
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
