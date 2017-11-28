package wallet

import (

	//"fmt"

	"errors"
	"fmt"

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
func NewReadableEntry(w Entry, v string) ReadableEntry {
	var secret string
	switch v {
	case "0.1":
		secret = w.Secret.Hex()
	case wltVersion:
		secret = w.EncryptedSeckey
	}
	return ReadableEntry{
		Address: w.Address.String(),
		Public:  w.Public.Hex(),
		Secret:  secret,
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
// converts base on the wallet version.
func (res ReadableEntries) toWalletEntries(v string) ([]Entry, error) {
	entries := make([]Entry, len(res))
	for i, re := range res {
		e, err := newEntryFromReadable(&re, v)
		if err != nil {
			return []Entry{}, err
		}

		entries[i] = *e
	}
	return entries, nil
}

// newEntryFromReadable creates WalletEntry base one ReadableWalletEntry
func newEntryFromReadable(w *ReadableEntry, v string) (*Entry, error) {
	if w.Secret == "" {
		return nil, errors.New("secret field is empty")
	}

	a, err := cipher.DecodeBase58Address(w.Address)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(w.Public)
	if err != nil {
		return nil, err
	}

	switch v {
	case "0.1":
		s, err := cipher.SecKeyFromHex(w.Secret)
		if err != nil {
			return nil, err
		}
		return &Entry{
			Address: a,
			Public:  p,
			Secret:  s,
		}, nil
	case wltVersion:
		return &Entry{
			Address:         a,
			Public:          p,
			EncryptedSeckey: w.Secret,
		}, nil
	}

	return nil, errors.New("invalid wallet version")
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
		readable[i] = NewReadableEntry(e, w.GetVersion())
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
func (rw *ReadableWallet) toWallet() (*Wallet, error) {
	ets, err := rw.Entries.toWalletEntries(rw.version())
	if err != nil {
		return nil, err
	}

	w := &Wallet{
		Meta:    rw.Meta,
		Entries: ets,
	}

	if err := w.Validate(); err != nil {
		return nil, fmt.Errorf("invalid wallet %s: %v", w.GetFilename(), err)
	}

	return w, nil
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

func (rw *ReadableWallet) version() string {
	return rw.Meta["version"]
}
