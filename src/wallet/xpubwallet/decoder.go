package xpubwallet

import (
	"encoding/json"
	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

// TODO: test this

// JSONDecoder implements the the WalletDecoder interface,
// which provides methods for encoding and decoding a XPub wallet in JSON format.
type JSONDecoder struct{}

// Encode encodes the XPub wallet to []byte, and error if any
func (d JSONDecoder) Encode(w wallet.Wallet) ([]byte, error) {
	return json.MarshalIndent(newReadableWallet(w.(*Wallet)), "", "    ")
}

// Decode decodes the XPub wallet from byte slice
func (d JSONDecoder) Decode(b []byte) (wallet.Wallet, error) {
	rw := readableWallet{}
	if err := json.Unmarshal(b, &rw); err != nil {
		return nil, err
	}

	return rw.toWallet()
}

type readableWallet struct {
	wallet.Meta `json:"meta"`
	Entries     readableXPubEntries `json:"entries"`
}

func (w readableWallet) toWallet() (*Wallet, error) {
	ad := wallet.ResolveAddressDecoder(w.Coin())
	entries, err := w.Entries.toXPubEntries(ad)
	if err != nil {
		return nil, err
	}

	xpubStr := w.Meta[wallet.MetaXPub]
	if xpubStr == "" {
		return nil, errors.New("missing xpub meta field")
	}

	xPub, err := parseXPub(xpubStr)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Meta:    w.Meta.Clone(),
		entries: entries,
		xpub:    xPub,
		decoder: &JSONDecoder{},
	}, nil
}

func newReadableWallet(w *Wallet) *readableWallet {
	return &readableWallet{
		Meta:    w.Meta.Clone(),
		Entries: newReadableEntries(w.entries),
	}
}

type readableXPubEntries []readableXPubEntry

func (es readableXPubEntries) toXPubEntries(ad wallet.AddressDecoder) (wallet.Entries, error) {
	entries := make(wallet.Entries, len(es))
	for i, e := range es {
		addr, err := ad.DecodeBase58Address(e.Address)
		if err != nil {
			return nil, err
		}

		p, err := cipher.PubKeyFromHex(e.Public)
		if err != nil {
			return nil, err
		}

		entries[i] = wallet.Entry{
			Address:     addr,
			Public:      p,
			ChildNumber: e.ChildNumber,
		}
	}

	return entries, nil
}

func newReadableEntries(entries wallet.Entries) readableXPubEntries {
	var res readableXPubEntries
	res = make([]readableXPubEntry, len(entries))
	for i, e := range entries {
		res[i] = readableXPubEntry{
			Address:     e.Address.String(),
			Public:      e.Public.Hex(),
			ChildNumber: e.ChildNumber,
		}
	}

	return res
}

type readableXPubEntry struct {
	Address     string `json:"address"`
	Public      string `json:"public"`
	ChildNumber uint32 `json:"child_number"` // For bip32/bip44
}
