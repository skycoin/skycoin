package wallet

import (
	//"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

// ReadableEntry wallet entry with json tags
type ReadableEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

// CoinSupply records the coin supply info
type CoinSupply struct {
	CurrentSupply                           int      `json:"coinSupply"`
	CoinCap                                 int      `json:"coinCap"`
	UndistributedLockedCoinBalance          int      `json:"UndistributedLockedCoinBalance"`
	UndistributedLockedCoinHoldingAddresses []string `json:"UndistributedLockedCoinHoldingAddresses"`
}

// NewReadableEntry creates readable wallet entry
func NewReadableEntry(w *Entry) ReadableEntry {
	return ReadableEntry{
		Address: w.Address.String(),
		Public:  w.Public.Hex(),
		Secret:  w.Secret.Hex(),
	}
}

// LoadReadableEntry load readable wallet entry from given file
func LoadReadableEntry(filename string) (ReadableEntry, error) {
	w := ReadableEntry{}
	err := util.LoadJSON(filename, &w)
	return w, err
}

// NewReadableEntryFromPubkey creates a ReadableWalletEntry given a pubkey hex string.
// The Secret field is left empty.
func NewReadableEntryFromPubkey(pub string) ReadableEntry {
	pubkey := cipher.MustPubKeyFromHex(pub)
	addr := cipher.AddressFromPubKey(pubkey)
	return ReadableEntry{
		Address: addr.String(),
		Public:  pub,
	}
}

// Save persists to disk
func (re *ReadableEntry) Save(filename string) error {
	return util.SaveJSONSafe(filename, re, 0600)
}

// ReadableEntries array of ReadableEntry
type ReadableEntries []ReadableEntry

// ToWalletEntries convert readable entries to entries
func (res ReadableEntries) ToWalletEntries() []Entry {
	entries := make([]Entry, len(res))
	for i, re := range res {
		we := NewEntryFromReadable(&re)
		if err := we.Verify(); err != nil {
			log.Panicf("Invalid wallet entry loaded. Address: %s", re.Address)
		}
		entries[i] = we
	}
	return entries
}

// ReadableWallet used for [de]serialization of a Wallet
type ReadableWallet struct {
	Meta    map[string]string `json:"meta"`
	Entries ReadableEntries   `json:"entries"`
}

// ByTm for sort ReadableWallets
type ByTm []*ReadableWallet

func (bt ByTm) Len() int {
	return len(bt)
}

func (bt ByTm) Less(i, j int) bool {
	return bt[i].Meta["tm"] < bt[j].Meta["tm"]
}

func (bt ByTm) Swap(i, j int) {
	bt[i], bt[j] = bt[j], bt[i]
}

// ReadableWalletCtor readable wallet creator
type ReadableWalletCtor func(w Wallet) *ReadableWallet

// NewReadableWallet creates readable wallet
func NewReadableWallet(w Wallet) *ReadableWallet {
	//return newReadableWallet(w, NewReadableWalletEntry)
	readable := make(ReadableEntries, len(w.Entries))
	i := 0
	for _, e := range w.Entries {
		readable[i] = NewReadableEntry(&e)
		i++
	}
	return &ReadableWallet{
		Meta:    w.Meta,
		Entries: readable,
	}
}

// LoadReadableWallet loads a ReadableWallet from disk
func LoadReadableWallet(filename string) (*ReadableWallet, error) {
	w := &ReadableWallet{}
	err := w.Load(filename)
	return w, err
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableWallet) ToWallet() (Wallet, error) {
	return NewWalletFromReadable(rw), nil
}

// Save saves to filename
func (rw *ReadableWallet) Save(filename string) error {
	// logger.Info("Saving readable wallet to %s with filename %s", filename,
	// 	self.Meta["filename"])
	return util.SaveJSON(filename, rw, 0600)
}

// SaveSafe saves to filename, but won't overwrite existing
func (rw *ReadableWallet) SaveSafe(filename string) error {
	return util.SaveJSONSafe(filename, rw, 0600)
}

// Load loads from filename
func (rw *ReadableWallet) Load(filename string) error {
	return util.LoadJSON(filename, rw)
}
