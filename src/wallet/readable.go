package wallet

import (
	"fmt"
	//"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

// ReadableEntry wallet entry with json tags
type ReadableEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

// NewReadableEntry creates readable wallet entry
func NewReadableEntry(w Entry) ReadableEntry {
	return ReadableEntry{
		Address: w.Address.String(),
		Public:  w.Public.Hex(),
		Secret:  w.Secret.Hex(),
	}
}

// LoadReadableEntry load readable wallet entry from given file
func LoadReadableEntry(filename string) (ReadableEntry, error) {
	w := ReadableEntry{}
	err := file.LoadJSON(filename, &w)
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
	return file.SaveJSONSafe(filename, re, 0600)
}

// ReadableEntries array of ReadableEntry
type ReadableEntries []ReadableEntry

// ToWalletEntries convert readable entries to entries
func (res ReadableEntries) ToWalletEntries() ([]Entry, error) {
	entries := make([]Entry, len(res))
	for i, re := range res {
		e, err := NewEntryFromReadable(&re)
		if err != nil {
			return []Entry{}, err
		}

		if err := e.Verify(); err != nil {
			return []Entry{}, fmt.Errorf("convert readable wallet entry failed: %v", err)
		}

		entries[i] = *e
	}
	return entries, nil
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
	readable := make(ReadableEntries, len(w.Entries))
	for i, e := range w.Entries {
		readable[i] = NewReadableEntry(e)
	}

	meta := make(map[string]string, len(w.Meta))
	for k, v := range w.Meta {
		meta[k] = v
	}

	return &ReadableWallet{
		Meta:    meta,
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
	w, err := newWalletFromReadable(rw)
	if err != nil {
		return Wallet{}, err
	}
	return *w, nil
}

// Save saves to filename
func (rw *ReadableWallet) Save(filename string) error {
	// logger.Info("Saving readable wallet to %s with filename %s", filename,
	// 	self.Meta["filename"])
	return file.SaveJSON(filename, rw, 0600)
}

// SaveSafe saves to filename, but won't overwrite existing
func (rw *ReadableWallet) SaveSafe(filename string) error {
	return file.SaveJSONSafe(filename, rw, 0600)
}

// Load loads from filename
func (rw *ReadableWallet) Load(filename string) error {
	return file.LoadJSON(filename, rw)
}
