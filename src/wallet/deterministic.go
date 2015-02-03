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
	Name     string //deprecate
	Filename string //deprecate
	Seed     string
	// Only holds one entry for now, and is assumed to be the first
	// entry generated from seed.
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
		Filename: NewWalletFilename(""),
		Seed:     seed,
		Entry:    NewWalletEntryFromKeypair(pub, sec),
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
	seed := r.Extra["seed"].(string)
	return Wallet{
		Filename: r.Filename,
		Name:     r.Name,
		Entry:    r.Entries.ToWalletEntries().ToArray()[0],
		Seed:     seed,
	}
}

func (self *Wallet) GetType() string {
	return "deterministic"
}

func (self *Wallet) GetFilename() string {
	return self.Filename
}

func (self *Wallet) SetFilename(fn string) {
	self.Filename = fn
}

func (self *Wallet) GetID() WalletID {
	return WalletID(self.Seed[0:4])
}

func (self *Wallet) GetName() string {
	return self.Name
}

func (self *Wallet) SetName(name string) {
	self.Name = name
}

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

func (self *Wallet) Save(dir string) error {
	r := NewReadableWallet(*self)
	return r.Save(filepath.Join(dir, self.Filename))
}

func (self *Wallet) Load(dir string) error {
	r := &ReadableWallet{}
	if err := r.Load(filepath.Join(dir, self.Filename)); err != nil {
		return err
	}
	*self = NewWalletFromReadable(r)
	return nil
}

func (self *Wallet) GetExtraSerializerData() map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["seed"] = self.Seed
	return m
}
