package wallet

import (
	"github.com/skycoin/skycoin/src/cipher"
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

func newReadableEntries(entries Entries, coinType CoinType) ReadableEntries {
	re := make(ReadableEntries, len(entries))
	for i, e := range entries {
		re[i] = NewReadableEntry(coinType, e)
	}
	return re
}

// GetEntries returns this array
func (res ReadableEntries) GetEntries() ReadableEntries {
	return res
}

// toWalletEntries convert readable entries to entries
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

// Readable defines the readable wallet API.
// A readable wallet is the on-disk representation of a wallet.
type Readable interface {
	ToWallet() (Wallet, error)
	Timestamp() int64
	SetFilename(string)
	Filename() string
	GetEntries() ReadableEntries
}
