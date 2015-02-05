package wallet

import (
	"encoding/hex"
	"errors"
	"log"
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

//const DeterministicSeedLength = 1024

//type WalletSeed cipher.SHA256

/*
func NewWalletSeed() string {
	seed := cipher.SumSHA256(secp256k1.RandByte(DeterministicSeedLength))
	return hex.EncodeToString(seed[:])
}
*/

/*
func (self *WalletSeed) toWalletID() WalletID {
	// Uses the first 16 bytes of SHA256(seed) as id
	shaid := cipher.SumSHA256(self[:])
	return WalletID(hex.EncodeToString(shaid[:16]))
}
*/
//func (self WalletSeed) Hex() string {
//	return cipher.SHA256(self).Hex()
//}

type Wallet struct {
	//Name     string //deprecate
	//Filename string //deprecate
	//Seed     string
	// Only holds one entry for now, and is assumed to be the first
	// entry generated from seed.
	Meta  map[string]string
	Entry WalletEntry
}

//Generate Deterministic Wallet
//generates a random seed if seed is ""
func NewWallet(seed string) Wallet {

	//if seed is blank, generate a new seed
	if seed == "" {
		seed_raw := cipher.SumSHA256(secp256k1.RandByte(64))
		seed = hex.EncodeToString(seed_raw[:])
	}

	pub, sec := cipher.GenerateDeterministicKeyPair([]byte(seed[:]))
	return Wallet{
		//Filename: NewWalletFilename(""),
		Meta: map[string]string{
			"filename": NewWalletFilename(),
			"seed":     seed,
			"type":     "deterministic",
			"coin":     "sky"},

		//Seed:  seed,
		Entry: NewWalletEntryFromKeypair(pub, sec),
	}
}

func NewWalletFromReadable(r *ReadableWallet) Wallet {
	//if r.Type != WalletType {
	//	log.Panic("ReadableWallet type must be Deterministic")
	//}
	if len(r.Entries) != 1 {
		log.Panic("Deterministic wallets have exactly 1 entry")
	}
	//should be string
	//seed := r.Extra["seed"].(string)
	w := Wallet{
		/*
			Meta: map[string]string{
				"filename": r.Meta["filename"],
				"seed":     r.Meta["seed"],
				"type":     r.Meta["type"],
				"coin":     r.Meta["coin"],
			},
		*/
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

//deprecate
//func (self *Wallet) GetName() string {
//	return self.Meta["filename"]
//}

//deprecate
//func (self *Wallet) SetName(name string) {
//	self.Meta["filename"] = name
//}

/*
	Refactor
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

//

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

/*
func (self *Wallet) GetExtraSerializerData() map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["seed"] = self.Seed
	return m
}
*/
