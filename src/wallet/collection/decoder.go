package collection

import (
	"encoding/json"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/wallet"
)

// JSONDecoder implements the Decoder interface for deterministic wallet
type JSONDecoder struct{}

// Encode encodes the deterministic wallet to []byte, and error if any
func (d JSONDecoder) Encode(w wallet.Wallet) ([]byte, error) {
	rw := newReadableDeterministicWallet(w.(*Wallet))
	return json.MarshalIndent(rw, "", "    ")
}

// Decode decodes the deterministic wallet to []byte, and error if any
func (d JSONDecoder) Decode(b []byte) (wallet.Wallet, error) {
	var rw readableDeterministicWallet
	if err := json.Unmarshal(b, &rw); err != nil {
		return nil, err
	}

	return rw.toWallet()
}

// readableEntry wallet entry with json tags
type readableEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

// newReadableEntry creates readable wallet entry
func newReadableEntry(coinType wallet.CoinType, e wallet.Entry) readableEntry {
	re := readableEntry{}
	if !e.Address.Null() {
		re.Address = e.Address.String()
	}

	if !e.Public.Null() {
		re.Public = e.Public.Hex()
	}

	if !e.Secret.Null() {
		d := wallet.ResolveSecKeyDecoder(coinType)
		re.Secret = d.SecKeyToHex(e.Secret)
	}

	return re
}

// readableEntries array of readableEntry
type readableEntries []readableEntry

func newReadableEntries(entries wallet.Entries, coinType wallet.CoinType) readableEntries {
	re := make(readableEntries, len(entries))
	for i, e := range entries {
		re[i] = newReadableEntry(coinType, e)
	}
	return re
}

// toWalletEntries convert readable entries to entries
// converts base on the wallet version.
func (res readableEntries) toWalletEntries(coinType wallet.CoinType, isEncrypted bool) ([]wallet.Entry, error) {
	entries := make([]wallet.Entry, len(res))
	for i, re := range res {
		e, err := newEntryFromReadable(coinType, &re)
		if err != nil {
			return []wallet.Entry{}, err
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
func newEntryFromReadable(coinType wallet.CoinType, re *readableEntry) (*wallet.Entry, error) {
	var a cipher.Addresser
	var err error

	switch coinType {
	case wallet.CoinTypeSkycoin:
		a, err = cipher.DecodeBase58Address(re.Address)
	case wallet.CoinTypeBitcoin:
		a, err = cipher.DecodeBase58BitcoinAddress(re.Address)
	default:
		panic(fmt.Errorf("invalid coin type %q", coinType))
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
		case wallet.CoinTypeSkycoin:
			secret, err = cipher.SecKeyFromHex(re.Secret)
		case wallet.CoinTypeBitcoin:
			secret, err = cipher.SecKeyFromBitcoinWalletImportFormat(re.Secret)
		default:
			panic(fmt.Errorf("invalid coin type %q", coinType))
		}
		if err != nil {
			return nil, err
		}
	}

	return &wallet.Entry{
		Address: a,
		Public:  p,
		Secret:  secret,
	}, nil
}

// readableDeterministicWallet used for [de]serialization of a deterministic wallet
type readableDeterministicWallet struct {
	wallet.Meta `json:"meta"`
	Entries     readableEntries `json:"entries"`
}

// newReadableDeterministicWallet creates readable wallet
func newReadableDeterministicWallet(w *Wallet) *readableDeterministicWallet {
	return &readableDeterministicWallet{
		Meta:    w.Meta.Clone(),
		Entries: newReadableEntries(w.entries, w.Meta.Coin()),
	}
}

// toWallet convert readable wallet to Wallet
func (rw *readableDeterministicWallet) toWallet() (wallet.Wallet, error) {
	w := &Wallet{
		Meta: rw.Meta.Clone(),
	}

	// make sure "sky", "btc" normalize to "skycoin", "bitcoin"
	ct, err := wallet.ResolveCoinType(string(w.Meta.Coin()))
	if err != nil {
		return nil, err
	}

	w.SetCoin(ct)

	if err := w.Validate(); err != nil {
		err := fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
		return nil, err
	}

	ets, err := rw.Entries.toWalletEntries(w.Meta.Coin(), w.Meta.IsEncrypted())
	if err != nil {
		return nil, err
	}

	w.entries = ets

	return w, nil
}
