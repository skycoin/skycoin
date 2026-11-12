// Package ui2 src/ui2/browser.go c3-ui-wallet
package ui2

import (
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
)

// A browser wallet is kept in the browser's localStorage, in the Angular
// wallet's format, and only static/bw.js reads it. The pages are sent the
// public part, BrowserWallet, and never a seed, key or password.

// BrowserAddress is one address of a browser wallet. Chain is "external" or
// "change" for a bip44 wallet and empty otherwise.
type BrowserAddress struct {
	Address string `json:"address"`
	Chain   string `json:"chain,omitempty"`
	Index   int    `json:"index"`
}

// BrowserWallet is the public part of a browser wallet. ID is its first
// address, which is how the Angular wallet tells wallets apart.
type BrowserWallet struct {
	ID        string           `json:"id"`
	Label     string           `json:"label"`
	Type      string           `json:"type"`
	Encrypted bool             `json:"encrypted"`
	Temporary bool             `json:"temporary,omitempty"`
	XPub      string           `json:"xpub,omitempty"`
	Addresses []BrowserAddress `json:"addresses"`
}

var errBadWallet = errors.New("bw.err-data")

// ParseBrowserWallets reads the summaries bw.js posts, dropping none: a
// summary with an address that does not decode fails the whole request.
func ParseBrowserWallets(s string) ([]BrowserWallet, error) {
	var ws []BrowserWallet
	if err := json.Unmarshal([]byte(s), &ws); err != nil {
		return nil, errBadWallet
	}
	for _, w := range ws {
		if len(w.Addresses) == 0 || w.ID != w.Addresses[0].Address {
			return nil, errBadWallet
		}
		for _, a := range w.Addresses {
			if _, err := cipher.DecodeBase58Address(a.Address); err != nil {
				return nil, errBadWallet
			}
		}
	}
	return ws, nil
}

// ParseBrowserWallet reads one summary.
func ParseBrowserWallet(s string) (BrowserWallet, error) {
	ws, err := ParseBrowserWallets("[" + s + "]")
	if err != nil || len(ws) != 1 {
		return BrowserWallet{}, errBadWallet
	}
	return ws[0], nil
}

func (w BrowserWallet) addrs() []string {
	out := make([]string, len(w.Addresses))
	for i, a := range w.Addresses {
		out[i] = a.Address
	}
	return out
}

func (w BrowserWallet) entries() []readable.WalletEntry {
	out := make([]readable.WalletEntry, len(w.Addresses))
	for i, a := range w.Addresses {
		out[i] = readable.WalletEntry{Address: a.Address}
		if a.Chain != "" {
			change := uint32(0)
			if a.Chain == "change" {
				change = 1
			}
			idx := uint32(a.Index) //nolint:gosec
			out[i].Change, out[i].ChildNumber = &change, &idx
		}
	}
	return out
}

func bwWalletURL(id string) string { return "bw/wallet?id=" + url.QueryEscape(id) }

func bwSendURL(id, from string) string {
	u := "bw/send?id=" + url.QueryEscape(id)
	if from != "" {
		u += "&from=" + url.QueryEscape(from)
	}
	return u
}

// LoadBrowserWallets reads the balance of each browser wallet.
func LoadBrowserWallets(n Node, ws []BrowserWallet) Wallets {
	var out Wallets
	var total readable.BalancePair
	for _, w := range ws {
		r := WalletRow{
			ID:        w.ID,
			Label:     w.Label,
			Type:      w.Type,
			Encrypted: w.Encrypted,
			Temp:      w.Temporary,
			Addresses: len(w.Addresses),
			Href:      bwWalletURL(w.ID),
		}
		if b, err := n.Balance(w.addrs()); err != nil {
			r.Err = err.Error()
		} else {
			r.Balance = amountOf(b.BalancePair)
			total.Confirmed.Coins += b.Confirmed.Coins
			total.Confirmed.Hours += b.Confirmed.Hours
			total.Predicted.Coins += b.Predicted.Coins
			total.Predicted.Hours += b.Predicted.Hours
		}
		out.Rows = append(out.Rows, r)
	}
	sort.SliceStable(out.Rows, func(i, j int) bool {
		return strings.ToLower(out.Rows[i].Label) < strings.ToLower(out.Rows[j].Label)
	})
	out.Total = amountOf(total)
	return out
}

// LoadBrowserWallet reads the balance and addresses of browser wallet w.
func LoadBrowserWallet(n Node, w BrowserWallet) Wallet {
	out := Wallet{ID: w.ID, Meta: readable.WalletMeta{Label: w.Label, Type: w.Type, Encrypted: w.Encrypted, XPub: w.XPub}}
	out.BrowserAddrs = w.addrs()
	b, err := n.Balance(w.addrs())
	if err != nil {
		out.Err = err.Error()
		b = &api.BalanceResponse{}
	}
	out.Balance = amountOf(b.BalancePair)
	rows, exposed, err := addressRows(n, w.entries(), b, ListQuery{}, func(a string) string { return bwSendURL(w.ID, a) })
	out.Exposed = exposed
	if err != nil {
		out.Err = err.Error()
	}
	out.Addrs, out.Pager = Paginate(rows, ListQuery{}, len(rows))
	return out
}

// LoadBrowserTxns is a page of the history of browser wallet w.
func LoadBrowserTxns(n Node, w BrowserWallet, q ListQuery) Txns {
	out := historyOf(n, []HistWallet{w.hist()}, w.addrs(), q)
	out.ID, out.URL, out.Post = w.ID, "bw/txns", true
	return out
}

func (w BrowserWallet) hist() HistWallet {
	return HistWallet{ID: w.ID, Label: w.Label, Addrs: w.addrs()}
}

// LoadBrowserHistory is the history of every browser wallet in ws.
func LoadBrowserHistory(n Node, ws []BrowserWallet, filter string, q ListQuery) Txns {
	hw := make([]HistWallet, 0, len(ws))
	for _, w := range ws {
		hw = append(hw, w.hist())
	}
	return history(n, hw, filter, q, Txns{Filter: filter, URL: "bw/history/list", Post: true, History: true})
}
