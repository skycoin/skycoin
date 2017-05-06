package wallet

import (
	//"fmt"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

type ReadableWalletEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

type CoinSupply struct {
	CurrentSupply int `json:"coinSupply"`
	CoinCap  int `json:"coinCap"`
	UndistributedLockedCoinBalance int `json:"UndistributedLockedCoinBalance"`
	UndistributedLockedCoinHoldingAddresses []string `json:"UndistributedLockedCoinHoldingAddresses"`
}

type ReadableWalletEntryCtor func(w *WalletEntry) ReadableWalletEntry

func NewReadableWalletEntry(w *WalletEntry) ReadableWalletEntry {
	return ReadableWalletEntry{
		Address: w.Address.String(),
		Public:  w.Public.Hex(),
		Secret:  w.Secret.Hex(),
	}
}

func LoadReadableWalletEntry(filename string) (ReadableWalletEntry, error) {
	w := ReadableWalletEntry{}
	err := util.LoadJSON(filename, &w)
	return w, err
}

// Creates a ReadableWalletEntry given a pubkey hex string.  The Secret field
// is left empty.
func ReadableWalletEntryFromPubkey(pub string) ReadableWalletEntry {
	pubkey := cipher.MustPubKeyFromHex(pub)
	addr := cipher.AddressFromPubKey(pubkey)
	return ReadableWalletEntry{
		Address: addr.String(),
		Public:  pub,
	}
}

func (self *ReadableWalletEntry) Save(filename string) error {
	return util.SaveJSONSafe(filename, self, 0600)
}

type ReadableWalletEntries []ReadableWalletEntry

func (self ReadableWalletEntries) ToWalletEntries() []WalletEntry {
	entries := make([]WalletEntry, len(self))
	for i, re := range self {
		we := WalletEntryFromReadable(&re)
		if err := we.Verify(); err != nil {
			log.Panicf("Invalid wallet entry loaded. Address: %s", re.Address)
		}
		entries[i] = we
	}
	return entries
}

// Used for [de]serialization of a Wallet
type ReadableWallet struct {
	Meta    map[string]string     `json:"meta"`
	Entries ReadableWalletEntries `json:"entries"`
}

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

type ReadableWalletCtor func(w Wallet) *ReadableWallet

func NewReadableWallet(w Wallet) *ReadableWallet {
	//return newReadableWallet(w, NewReadableWalletEntry)
	readable := make(ReadableWalletEntries, len(w.Entries))
	i := 0
	for _, e := range w.Entries {
		readable[i] = NewReadableWalletEntry(&e)
		i++
	}
	return &ReadableWallet{
		Meta:    w.Meta,
		Entries: readable,
	}
}

// Loads a ReadableWallet from disk
func LoadReadableWallet(filename string) (*ReadableWallet, error) {
	w := &ReadableWallet{}
	err := w.Load(filename)
	return w, err
}

func (self *ReadableWallet) ToWallet() (Wallet, error) {
	return NewWalletFromReadable(self), nil
}

// Saves to filename
func (self *ReadableWallet) Save(filename string) error {
	// logger.Info("Saving readable wallet to %s with filename %s", filename,
	// 	self.Meta["filename"])
	return util.SaveJSON(filename, self, 0600)
}

// Saves to filename, but won't overwrite existing
func (self *ReadableWallet) SaveSafe(filename string) error {
	return util.SaveJSONSafe(filename, self, 0600)
}

// Loads from filename
func (self *ReadableWallet) Load(filename string) error {
	return util.LoadJSON(filename, self)
}
