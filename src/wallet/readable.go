package wallet

import (
	"fmt"
	"strconv"

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
func NewReadableEntry(coinType CoinType, w Entry) ReadableEntry {
	re := ReadableEntry{}
	if !w.Address.Null() {
		re.Address = w.Address.String()
	}

	if !w.Public.Null() {
		re.Public = w.Public.Hex()
	}

	if !w.Secret.Null() {
		switch coinType {
		case CoinTypeSkycoin:
			re.Secret = w.Secret.Hex()
		case CoinTypeBitcoin:
			re.Secret = cipher.BitcoinWalletImportFormatFromSeckey(w.Secret)
		default:
			logger.Panicf("Invalid coin type %q", coinType)
		}
	}

	return re
}

// ReadableEntries array of ReadableEntry
type ReadableEntries []ReadableEntry

// ToWalletEntries convert readable entries to entries
// converts base on the wallet version.
func (res ReadableEntries) toWalletEntries(coinType CoinType, isEncrypted bool) ([]Entry, error) {
	entries := make([]Entry, len(res))
	for i, re := range res {
		e, err := newEntryFromReadable(coinType, &re)
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
func newEntryFromReadable(coinType CoinType, w *ReadableEntry) (*Entry, error) {
	var a cipher.Addresser
	var err error

	switch coinType {
	case CoinTypeSkycoin:
		a, err = cipher.DecodeBase58Address(w.Address)
	case CoinTypeBitcoin:
		a, err = cipher.DecodeBase58BitcoinAddress(w.Address)
	default:
		logger.Panicf("Invalid coin type %q", coinType)
	}

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
		switch coinType {
		case CoinTypeSkycoin:
			secret, err = cipher.SecKeyFromHex(w.Secret)
		case CoinTypeBitcoin:
			secret, err = cipher.SecKeyFromBitcoinWalletImportFormat(w.Secret)
		default:
			logger.Panicf("Invalid coin type %q", coinType)
		}
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
		readable[i] = NewReadableEntry(w.coin(), e)
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

	ets, err := rw.Entries.toWalletEntries(w.coin(), w.IsEncrypted())
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

func (rw *ReadableWallet) timestamp() int64 {
	// Intentionally ignore the error when parsing the timestamp,
	// if it isn't valid or is missing it will be set to 0
	x, _ := strconv.ParseInt(rw.Meta[metaTimestamp], 10, 64) //nolint:errcheck
	return x
}

func (rw *ReadableWallet) filename() string {
	return rw.Meta[metaFilename]
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
