package wallet

import (
	"encoding/hex"
	"errors"
	"log"
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

//Filename
//Seed
//Type - wallet type
//Coin - coin type
type Wallet struct {
	Meta  map[string]string
	Entry WalletEntry
}

//Generate Deterministic Wallet
//generates a random seed if seed is ""
func NewWallet(seed string, wltName string) Wallet {
	//if seed is blank, generate a new seed
	if seed == "" {
		seed_raw := cipher.SumSHA256(secp256k1.RandByte(64))
		seed = hex.EncodeToString(seed_raw[:])
	}

	pub, sec := cipher.GenerateDeterministicKeyPair([]byte(seed[:]))
	return Wallet{
		Meta: map[string]string{
			"filename": wltName,
			"seed":     seed,
			"type":     "deterministic",
			"coin":     "sky"},
		Entry: NewWalletEntryFromKeypair(pub, sec),
	}
}

func NewWalletFromReadable(r *ReadableWallet) Wallet {
	if len(r.Entries) != 1 {
		log.Panic("Deterministic wallets have exactly 1 entry")
	}

	w := Wallet{
		Meta:  r.Meta,
		Entry: r.Entries.ToWalletEntries().ToArray()[0],
	}

	err := w.Validate()
	if err != nil {
		log.Panic("Wallet %s invalid: %v", w.GetFilename, err)
	}
	return w

}

func (self *Wallet) Validate() error {

	if _, ok := self.Meta["filename"]; !ok {
		return errors.New("filename not set")
	}
	if _, ok := self.Meta["seed"]; !ok {
		return errors.New("seed not set")
	}
	wallet_type, ok := self.Meta["type"]
	if !ok {
		return errors.New("type not set")
	}
	if wallet_type != "deterministic" {
		return errors.New("wallet type invalid")
	}

	coin_type, ok := self.Meta["coin"]
	if !ok {
		return errors.New("coin field not set")
	}
	if coin_type != "sky" {
		return errors.New("coin type invalid")
	}

	return nil

}

func (self *Wallet) GetType() string {
	return self.Meta["type"]
}

func (self *Wallet) GetFilename() string {
	return self.Meta["filename"]
}

func (self *Wallet) SetFilename(fn string) {
	self.Meta["filename"] = fn
}

func (self *Wallet) GetID() WalletID {
	return WalletID(self.Meta["filename"])
}

/*
Refactor
- should be list of entries
- should deterministically generate new entries
*/

func (self *Wallet) NumEntries() int {
	return 1
}

func (self *Wallet) GetEntries() WalletEntries {
	m := make(WalletEntries, 1)
	m[self.Entry.Address] = self.Entry
	return m
}

func (self *Wallet) GetAddressSet() AddressSet {
	m := make(AddressSet, 1)
	m[self.Entry.Address] = byte(1)
	return m
}

func (self *Wallet) CreateEntry() WalletEntry {
	log.Panic("Multiple entries not implemented for deterministic wallet")
	return WalletEntry{}
}

func (self *Wallet) GetAddresses() []cipher.Address {
	return []cipher.Address{self.Entry.Address}
}

func (self *Wallet) GetEntry(a cipher.Address) (WalletEntry, bool) {
	if a == self.Entry.Address {
		return self.Entry, true
	} else {
		return WalletEntry{}, false
	}
}

func (self *Wallet) AddEntry(e WalletEntry) error {
	return errors.New("Adding entries to deterministic wallet not allowed")
}

/*
	End Refactor
*/

func (self *Wallet) Save(dir string) error {
	r := NewReadableWallet(*self)
	return r.Save(filepath.Join(dir, self.GetFilename()))
}

func (self *Wallet) Load(dir string) error {
	r := &ReadableWallet{}
	if err := r.Load(filepath.Join(dir, self.GetFilename())); err != nil {
		return err
	}
	*self = NewWalletFromReadable(r)
	return nil
}
