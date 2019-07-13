package wallet

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

// ReadableEntry wallet entry with json tags
type ReadableEntry struct {
	Address     string  `json:"address"`
	Public      string  `json:"public_key"`
	Secret      string  `json:"secret_key"`
	ChildNumber *uint32 `json:"child_number,omitempty"` // For bip32/bip44
	Change      *uint32 `json:"change,omitempty"`       // For bip44
}

// NewReadableEntry creates readable wallet entry
func NewReadableEntry(coinType CoinType, walletType string, e Entry) ReadableEntry {
	re := ReadableEntry{}
	if !e.Address.Null() {
		re.Address = e.Address.String()
	}

	if !e.Public.Null() {
		re.Public = e.Public.Hex()
	}

	if !e.Secret.Null() {
		switch coinType {
		case CoinTypeSkycoin:
			re.Secret = e.Secret.Hex()
		case CoinTypeBitcoin:
			re.Secret = cipher.BitcoinWalletImportFormatFromSeckey(e.Secret)
		default:
			logger.Panicf("Invalid coin type %q", coinType)
		}
	}

	switch walletType {
	case WalletTypeBip44:
		cn := e.ChildNumber
		re.ChildNumber = &cn
		change := e.Change
		re.Change = &change
	case WalletTypeXPub:
		cn := e.ChildNumber
		re.ChildNumber = &cn
		if e.Change != 0 {
			logger.Panicf("wallet.Entry.Change is not 0 but wallet type is %q", walletType)
		}
	default:
		if e.ChildNumber != 0 {
			logger.Panicf("wallet.Entry.ChildNumber is not 0 but wallet type is %q", walletType)
		}
		if e.Change != 0 {
			logger.Panicf("wallet.Entry.Change is not 0 but wallet type is %q", walletType)
		}
	}

	return re
}

// ReadableEntries array of ReadableEntry
type ReadableEntries []ReadableEntry

func newReadableEntries(entries Entries, coinType CoinType, walletType string) ReadableEntries {
	re := make(ReadableEntries, len(entries))
	for i, e := range entries {
		re[i] = NewReadableEntry(coinType, walletType, e)
	}
	return re
}

// GetEntries returns this array
func (res ReadableEntries) GetEntries() ReadableEntries {
	return res
}

// toWalletEntries convert readable entries to entries
// converts base on the wallet version.
func (res ReadableEntries) toWalletEntries(coinType CoinType, walletType string, isEncrypted bool) ([]Entry, error) {
	entries := make([]Entry, len(res))
	for i, re := range res {
		e, err := newEntryFromReadable(coinType, walletType, &re)
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
func newEntryFromReadable(coinType CoinType, walletType string, re *ReadableEntry) (*Entry, error) {
	var a cipher.Addresser
	var err error

	switch coinType {
	case CoinTypeSkycoin:
		a, err = cipher.DecodeBase58Address(re.Address)
	case CoinTypeBitcoin:
		a, err = cipher.DecodeBase58BitcoinAddress(re.Address)
	default:
		logger.Panicf("Invalid coin type %q", coinType)
	}

	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(re.Public)
	if err != nil {
		return nil, err
	}

	// Decodes the secret hex string if any
	var secret cipher.SecKey
	if re.Secret != "" {
		switch coinType {
		case CoinTypeSkycoin:
			secret, err = cipher.SecKeyFromHex(re.Secret)
		case CoinTypeBitcoin:
			secret, err = cipher.SecKeyFromBitcoinWalletImportFormat(re.Secret)
		default:
			logger.Panicf("Invalid coin type %q", coinType)
		}
		if err != nil {
			return nil, err
		}
	}

	var childNumber uint32
	var change uint32
	switch walletType {
	case WalletTypeBip44:
		if re.ChildNumber == nil {
			return nil, fmt.Errorf("child_number required for %q wallet type", walletType)
		}
		if re.Change == nil {
			return nil, fmt.Errorf("change required for %q wallet type", walletType)
		}

		childNumber = *re.ChildNumber
		change = *re.Change

		switch change {
		case 0, 1:
		default:
			return nil, errors.New("change must be either 0 or 1")
		}

	default:
		if re.ChildNumber != nil {
			return nil, fmt.Errorf("child_number should not be set for %q wallet type", walletType)
		}
		if re.Change != nil {
			return nil, fmt.Errorf("change should not be set for %q wallet type", walletType)
		}
	}

	return &Entry{
		Address:     a,
		Public:      p,
		Secret:      secret,
		ChildNumber: childNumber,
		Change:      change,
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
