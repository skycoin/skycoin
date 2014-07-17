package wallet

import (
	"encoding/hex"
	"errors"
	"log"
	"path/filepath"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

const DeterministicSeedLength = 1024

type DeterministicWalletSeed cipher.SHA256

func NewDeterministicWalletSeed() DeterministicWalletSeed {
	seed := cipher.SumSHA256(secp256k1.RandByte(DeterministicSeedLength))
	return DeterministicWalletSeed(seed)
}

func (self *DeterministicWalletSeed) toWalletID() WalletID {
	// Uses the first 16 bytes of SHA256(seed) as id
	shaid := cipher.SumSHA256(self[:])
	return WalletID(hex.EncodeToString(shaid[:16]))
}

func (self DeterministicWalletSeed) Hex() string {
	return cipher.SHA256(self).Hex()
}

type DeterministicWallet struct {
	Name     string
	Filename string
	Seed     DeterministicWalletSeed
	// Only holds one entry for now, and is assumed to be the first
	// entry generated from seed.
	Entry WalletEntry
}

func NewDeterministicWallet() Wallet {
	seed := NewDeterministicWalletSeed()
	pub, sec := cipher.GenerateDeterministicKeyPair(seed[:])
	return &DeterministicWallet{
		Filename: NewWalletFilename(seed.toWalletID()),
		Seed:     seed,
		Entry:    NewWalletEntryFromKeypair(pub, sec),
	}
}

func NewDeterministicWalletFromReadable(r *ReadableWallet) Wallet {
	if r.Type != DeterministicWalletType {
		log.Panic("ReadableWallet type must be Deterministic")
	}
	if len(r.Entries) != 1 {
		log.Panic("Deterministic wallets have exactly 1 entry")
	}
	seed := cipher.MustSHA256FromHex(r.Extra["seed"].(string))
	return &DeterministicWallet{
		Filename: r.Filename,
		Name:     r.Name,
		Entry:    r.Entries.ToWalletEntries().ToArray()[0],
		Seed:     DeterministicWalletSeed(seed),
	}
}

func (self *DeterministicWallet) GetType() WalletType {
	return DeterministicWalletType
}

func (self *DeterministicWallet) GetFilename() string {
	return self.Filename
}

func (self *DeterministicWallet) SetFilename(fn string) {
	self.Filename = fn
}

func (self *DeterministicWallet) GetID() WalletID {
	return self.Seed.toWalletID()
}

func (self *DeterministicWallet) GetName() string {
	return self.Name
}

func (self *DeterministicWallet) SetName(name string) {
	self.Name = name
}

func (self *DeterministicWallet) NumEntries() int {
	return 1
}

func (self *DeterministicWallet) GetEntries() WalletEntries {
	m := make(WalletEntries, 1)
	m[self.Entry.Address] = self.Entry
	return m
}

func (self *DeterministicWallet) GetAddressSet() AddressSet {
	m := make(AddressSet, 1)
	m[self.Entry.Address] = byte(1)
	return m
}

func (self *DeterministicWallet) CreateEntry() WalletEntry {
	log.Panic("Multiple entries not implemented for deterministic wallet")
	return WalletEntry{}
}

func (self *DeterministicWallet) GetAddresses() []cipher.Address {
	return []cipher.Address{self.Entry.Address}
}

func (self *DeterministicWallet) GetEntry(a cipher.Address) (WalletEntry, bool) {
	if a == self.Entry.Address {
		return self.Entry, true
	} else {
		return WalletEntry{}, false
	}
}

func (self *DeterministicWallet) AddEntry(e WalletEntry) error {
	return errors.New("Adding entries to deterministic wallet not allowed")
}

func (self *DeterministicWallet) Save(dir string) error {
	r := NewReadableWallet(self)
	return r.Save(filepath.Join(dir, self.Filename))
}

func (self *DeterministicWallet) Load(dir string) error {
	r := &ReadableWallet{}
	if err := r.Load(filepath.Join(dir, self.Filename)); err != nil {
		return err
	}
	*self = *(NewDeterministicWalletFromReadable(r)).(*DeterministicWallet)
	return nil
}

func (self *DeterministicWallet) GetExtraSerializerData() map[string]interface{} {
	m := make(map[string]interface{}, 1)
	m["seed"] = self.Seed.Hex()
	return m
}
