package wallet

import (

	//"fmt"

	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
)

var (
	emptyAddress = cipher.Address{}
	emptyPubkey  = cipher.PubKey{}
	emptySeckey  = cipher.SecKey{}
)

// ReadableEntry wallet entry with json tags
type ReadableEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

// NewReadableEntry creates readable wallet entry
func NewReadableEntry(w Entry) ReadableEntry {
	re := ReadableEntry{}
	if w.Address != emptyAddress {
		re.Address = w.Address.String()
	}

	if w.Public != emptyPubkey {
		re.Public = w.Public.Hex()
	}

	if w.Secret != emptySeckey {
		re.Secret = w.Secret.Hex()
	}

	return re
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
func (res ReadableEntries) toWalletEntries(isEncrypted bool) ([]Entry, error) {
	entries := make([]Entry, len(res))
	for i, re := range res {
		e, err := newEntryFromReadable(&re)
		if err != nil {
			return []Entry{}, err
		}

		// Verify the wallet if it's not encrypted
		if !isEncrypted && re.Secret != "" {
			if err := e.Verify(); err != nil {
				return nil, err
			}
		}

		entries[i] = *e
	}
	return entries, nil
}

// newEntryFromReadable creates WalletEntry base one ReadableWalletEntry
func newEntryFromReadable(w *ReadableEntry) (*Entry, error) {
	a, err := cipher.DecodeBase58Address(w.Address)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(w.Public)
	if err != nil {
		return nil, err
	}

	// Decodes the secret hex string if any
	var secret cipher.SecKey
	if w.Secret != "" {
		var err error
		secret, err = cipher.SecKeyFromHex(w.Secret)
		if err != nil {
			return nil, err
		}
	}

	return &Entry{
		Address: a,
		Public:  p,
		Secret:  secret,
	}, nil
}

// ReadableWallet used for [de]serialization of a Wallet
type ReadableWallet struct {
	Meta    map[string]string `json:"meta"`
	Entries ReadableEntries   `json:"entries"`
}

// NewReadableWallet creates readable wallet
func NewReadableWallet(w *Wallet) *ReadableWallet {
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
	if err := w.Load(filename); err != nil {
		return nil, fmt.Errorf("load wallet %s failed: %v", filename, err)
	}
	return w, nil
}

// ToWallet convert readable wallet to Wallet
func (rw *ReadableWallet) ToWallet() (*Wallet, error) {
	w := &Wallet{
		Meta: rw.Meta,
	}

	if err := w.Validate(); err != nil {
		return nil, fmt.Errorf("invalid wallet %s: %v", w.Filename(), err)
	}

	ets, err := rw.Entries.toWalletEntries(w.IsEncrypted())
	if err != nil {
		return nil, err
	}

	w.Entries = ets

	return w, nil
}

// Save saves to filename
func (rw *ReadableWallet) Save(filename string) error {
	return file.SaveJSON(filename, rw, 0600)
}

// Load loads from filename
func (rw *ReadableWallet) Load(filename string) error {
	return file.LoadJSON(filename, rw)
}

func (rw *ReadableWallet) version() string {
	return rw.Meta[metaVersion]
}

func (rw *ReadableWallet) time() string {
	return rw.Meta[metaTm]
}

// Erase remove sensitive data
func (rw *ReadableWallet) Erase() {
	delete(rw.Meta, metaSeed)
	delete(rw.Meta, metaLastSeed)
	delete(rw.Meta, metaSecrets)
	for i := range rw.Entries {
		rw.Entries[i].Secret = ""
	}
}
