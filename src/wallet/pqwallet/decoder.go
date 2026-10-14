package pqwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

// JSONDecoder implements the Decoder interface for post-quantum wallets
type JSONDecoder struct{}

// Encode encodes the wallet to []byte
func (d JSONDecoder) Encode(w wallet.Wallet) ([]byte, error) {
	return json.MarshalIndent(newReadableWallet(w.(*Wallet)), "", "    ")
}

// Decode decodes the wallet from a byte slice
func (d JSONDecoder) Decode(b []byte) (wallet.Wallet, error) {
	var rw readableWallet
	if err := json.Unmarshal(b, &rw); err != nil {
		return nil, err
	}

	return rw.toWallet()
}

// readableEntry has no public key. The ML-DSA-44 key is 1312 bytes and is
// derived from the secret key when it is needed.
type readableEntry struct {
	Address string `json:"address"`
	Secret  string `json:"secret_key"`
}

type readableWallet struct {
	wallet.Meta `json:"meta"`
	Entries     []readableEntry `json:"entries"`
}

func newReadableWallet(w *Wallet) *readableWallet {
	rw := &readableWallet{
		Meta:    w.Meta.Clone(),
		Entries: make([]readableEntry, len(w.entries)),
	}
	for i, e := range w.entries {
		rw.Entries[i].Address = e.Address.String()
		if !e.Secret.Null() {
			rw.Entries[i].Secret = e.Secret.Hex()
		}
	}
	return rw
}

func (rw *readableWallet) toWallet() (wallet.Wallet, error) {
	w := &Wallet{
		Meta: rw.Meta.Clone(),
	}

	if err := w.Validate(); err != nil {
		return nil, fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
	}

	if err := validateMeta(w.Meta); err != nil {
		return nil, fmt.Errorf("invalid wallet %q: %v", w.Filename(), err)
	}

	w.entries = make(wallet.Entries, len(rw.Entries))
	for i, re := range rw.Entries {
		a, err := cipher.DecodeBase58Address(re.Address)
		if err != nil {
			return nil, err
		}
		if a.Version != cipher.AddressVersionPQ {
			return nil, fmt.Errorf("address %s is not a post-quantum address", re.Address)
		}

		e := wallet.Entry{
			Address:     a,
			ChildNumber: uint32(i), //nolint:gosec // entry count fits
		}
		if re.Secret != "" {
			b, err := hex.DecodeString(re.Secret)
			if err != nil || len(b) != len(e.Secret) {
				return nil, fmt.Errorf("invalid secret key for address %s", re.Address)
			}
			copy(e.Secret[:], b)
			if !w.IsEncrypted() {
				if err := e.Verify(); err != nil {
					return nil, err
				}
			}
		}
		w.entries[i] = e
	}

	return w, nil
}
