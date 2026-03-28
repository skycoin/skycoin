package readable

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/wallet"
)

// Balance has coins and hours
type Balance struct {
	Coins uint64 `json:"coins"`
	Hours uint64 `json:"hours"`
}

// NewBalance copies from wallet.Balance
func NewBalance(b wallet.Balance) Balance {
	return Balance{
		Coins: b.Coins,
		Hours: b.Hours,
	}
}

// BalancePair records the confirmed and predicted balance of an address
type BalancePair struct {
	Confirmed Balance `json:"confirmed"`
	Predicted Balance `json:"predicted"` // TODO rename "pending"
}

// NewBalancePair copies from wallet.BalancePair
func NewBalancePair(bp wallet.BalancePair) BalancePair {
	return BalancePair{
		Confirmed: NewBalance(bp.Confirmed),
		Predicted: NewBalance(bp.Predicted),
	}
}

// AddressBalances represents a map of address balances
type AddressBalances map[string]BalancePair

// NewAddressBalances copies from wallet.AddressBalances
func NewAddressBalances(wab wallet.AddressBalances) AddressBalances {
	ab := make(AddressBalances, len(wab))
	for k, v := range wab {
		ab[k] = NewBalancePair(v)
	}
	return ab
}

// WalletEntry the wallet entry struct
type WalletEntry struct {
	Address     string  `json:"address"`
	Public      string  `json:"public_key"`
	ChildNumber *uint32 `json:"child_number,omitempty"` // For bip32/44
	Change      *uint32 `json:"change,omitempty"`       // For bip44
}

// WalletAccount represents a BIP44 account with entries split by chain
type WalletAccount struct {
	Name            string        `json:"name"`
	Index           uint32        `json:"index"`
	ExternalEntries []WalletEntry `json:"external_entries"`
	ChangeEntries   []WalletEntry `json:"change_entries"`
}

// WalletMeta the wallet meta struct
type WalletMeta struct {
	Coin       wallet.CoinType   `json:"coin"`
	Filename   string            `json:"filename"`
	Label      string            `json:"label"`
	Type       string            `json:"type"`
	Version    string            `json:"version"`
	CryptoType crypto.CryptoType `json:"crypto_type"`
	Timestamp  int64             `json:"timestamp"`
	Temp       bool              `json:"temp"`
	Encrypted  bool              `json:"encrypted"`
	Bip44Coin  *bip44.CoinType   `json:"bip44_coin,omitempty"` // For bip44
	XPub       string            `json:"xpub,omitempty"`       // For xpub
}

// WalletResponse is the common wallet response structure used by both
// the daemon API and the skycoin-web thin client.
type WalletResponse struct {
	Meta     WalletMeta      `json:"meta"`
	Entries  []WalletEntry   `json:"entries"`
	Accounts []WalletAccount `json:"accounts,omitempty"`
}

// NewWalletResponse creates a WalletResponse from a wallet.Wallet interface.
func NewWalletResponse(w wallet.Wallet) (*WalletResponse, error) {
	var wr WalletResponse

	wr.Meta.Coin = w.Coin()
	wr.Meta.Filename = w.Filename()
	wr.Meta.Label = w.Label()
	wr.Meta.Type = w.Type()
	wr.Meta.Version = w.Version()
	wr.Meta.CryptoType = w.CryptoType()
	wr.Meta.Encrypted = w.IsEncrypted()
	wr.Meta.Timestamp = w.Timestamp()
	wr.Meta.Temp = w.IsTemp()

	var options []wallet.Option
	switch w.Type() {
	case wallet.WalletTypeBip44:
		bip44Coin := w.Bip44Coin()
		if bip44Coin == nil {
			return nil, errors.New("wallet has no Bip44Coin meta data")
		}
		wr.Meta.Bip44Coin = bip44Coin
		options = append(options, wallet.OptionExternal(), wallet.OptionChange())

		accounts := w.Accounts()
		wr.Accounts = make([]WalletAccount, len(accounts))
		for ai, acct := range accounts {
			wa := WalletAccount{
				Name:  acct.Name,
				Index: acct.Index,
			}

			extEntries, err := w.GetEntries(wallet.OptionAccount(acct.Index), wallet.OptionExternal())
			if err != nil {
				return nil, fmt.Errorf("failed to get external entries for account %d: %v", acct.Index, err)
			}
			wa.ExternalEntries = entriesToReadable(extEntries)

			chgEntries, err := w.GetEntries(wallet.OptionAccount(acct.Index), wallet.OptionChange())
			if err != nil {
				return nil, fmt.Errorf("failed to get change entries for account %d: %v", acct.Index, err)
			}
			wa.ChangeEntries = entriesToReadable(chgEntries)

			wr.Accounts[ai] = wa
		}
	case wallet.WalletTypeXPub:
		wr.Meta.XPub = w.XPub()
	}

	entries, err := w.GetEntries(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet entries: %v", err)
	}
	wr.Entries = make([]WalletEntry, len(entries))

	for i, e := range entries {
		wr.Entries[i] = WalletEntry{
			Address: e.Address.String(),
			Public:  e.Public.Hex(),
		}

		switch w.Type() {
		case wallet.WalletTypeBip44:
			childNumber := e.ChildNumber
			wr.Entries[i].ChildNumber = &childNumber
			change := e.Change
			wr.Entries[i].Change = &change
		case wallet.WalletTypeXPub:
			childNumber := e.ChildNumber
			wr.Entries[i].ChildNumber = &childNumber
		}
	}

	return &wr, nil
}

// entriesToReadable converts wallet entries to readable format with child number info.
func entriesToReadable(entries wallet.Entries) []WalletEntry {
	result := make([]WalletEntry, len(entries))
	for i, e := range entries {
		childNumber := e.ChildNumber
		change := e.Change
		result[i] = WalletEntry{
			Address:     e.Address.String(),
			Public:      e.Public.Hex(),
			ChildNumber: &childNumber,
			Change:      &change,
		}
	}
	return result
}
