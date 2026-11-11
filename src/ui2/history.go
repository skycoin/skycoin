package ui2

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// HistWallet is a wallet the history covers.
type HistWallet struct {
	ID    string
	Label string
	Addrs []string
}

// TxnRow is one transaction, seen from the addresses the history shows.
type TxnRow struct {
	Txid      string
	Time      string
	Confirmed bool
	// Block is the sequence of the block that holds it.
	Block uint64
	// Kind is "in", "out" or "self", coins moved between the wallets' own
	// addresses.
	Kind  string
	Coins string
	Hours string
	// Addrs are the addresses the row is about: the receiving ones of an
	// incoming transaction, the sending ones of an outgoing one.
	Addrs  []string
	Burned string
	In     []TxnIO
	Out    []TxnIO
}

// TxnIO is an input or output of a transaction.
type TxnIO struct {
	Address string
	Coins   string
	Hours   string
	Own     bool
}

// FilterOption is a wallet or one of its addresses in the history filter.
type FilterOption struct {
	Value string
	Label string
	Coins string
	Addrs []FilterOption
}

// Txns is a page of transaction history.
type Txns struct {
	ID     string
	Filter string
	// Options are the wallets and addresses the history can be narrowed to.
	// Empty when the list belongs to one wallet's page.
	Options []FilterOption
	Rows    []TxnRow
	Pager   Pager
	Err     string
	// URL is where the list asks for another page. Post sends the request
	// as a POST carrying the browser wallets' addresses.
	URL  string
	Post bool
}

// TxnPageSize is how many transactions a history page shows.
const TxnPageSize = 20

// A filter is "w:" and a wallet id, or "a:" and an address.
func walletFilter(id string) string { return "w:" + id }

func addrFilter(a string) string { return "a:" + a }

// selected are the addresses filter names among ws, every address when it
// names none of them.
func selected(ws []HistWallet, filter string) []string {
	var all []string
	for _, w := range ws {
		if filter == walletFilter(w.ID) {
			return w.Addrs
		}
		for _, a := range w.Addrs {
			if filter == addrFilter(a) {
				return []string{a}
			}
		}
		all = append(all, w.Addrs...)
	}
	return all
}

// LoadHistory is the history of every wallet the node keeps.
func LoadHistory(n Node, filter string, q ListQuery) Txns {
	out := Txns{Filter: filter, URL: "history/list"}
	ws, err := n.Wallets()
	if err != nil {
		out.Err = err.Error()
		return out
	}
	hw := make([]HistWallet, 0, len(ws))
	for _, w := range ws {
		h := HistWallet{ID: w.Meta.Filename, Label: w.Meta.Label}
		for _, e := range w.Entries {
			h.Addrs = append(h.Addrs, e.Address)
		}
		hw = append(hw, h)
	}
	return history(n, hw, filter, q, out)
}

// history fills out with the page of ws's history that filter and q ask for,
// and the filter options with their balances.
func history(n Node, ws []HistWallet, filter string, q ListQuery, out Txns) Txns {
	sort.SliceStable(ws, func(i, j int) bool { return strings.ToLower(ws[i].Label) < strings.ToLower(ws[j].Label) })
	var all []string
	for _, w := range ws {
		all = append(all, w.Addrs...)
	}
	var bal map[string]readable.BalancePair
	if len(all) > 0 {
		if b, err := n.Balance(all); err != nil {
			out.Err = err.Error()
		} else {
			bal = b.Addresses
		}
	}
	for _, w := range ws {
		o := FilterOption{Value: walletFilter(w.ID), Label: w.Label}
		if o.Label == "" {
			o.Label = w.ID
		}
		var sum uint64
		for _, a := range w.Addrs {
			c := bal[a].Confirmed.Coins
			sum += c
			o.Addrs = append(o.Addrs, FilterOption{Value: addrFilter(a), Label: a, Coins: coins(c)})
		}
		o.Coins = coins(sum)
		out.Options = append(out.Options, o)
	}
	page := historyOf(n, ws, selected(ws, filter), q)
	out.Rows, out.Pager = page.Rows, page.Pager
	if page.Err != "" {
		out.Err = page.Err
	}
	return out
}

// LoadTxns is the history of node wallet id, for its page.
func LoadTxns(n Node, id string, q ListQuery) Txns {
	out := Txns{ID: id, URL: walletPart("txns", id)}
	w, err := n.Wallet(id)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	hw := HistWallet{ID: id}
	for _, e := range w.Entries {
		hw.Addrs = append(hw.Addrs, e.Address)
	}
	page := historyOf(n, []HistWallet{hw}, hw.Addrs, q)
	out.Rows, out.Pager, out.Err = page.Rows, page.Pager, page.Err
	return out
}

// historyOf is a page of the history of addrs, newest first, pending ones
// on top. ws are every wallet, which tell change from a payment.
func historyOf(n Node, ws []HistWallet, addrs []string, q ListQuery) Txns {
	var out Txns
	var txns []readable.TransactionWithStatusVerbose
	if len(addrs) > 0 {
		var err error
		if txns, err = n.TransactionsVerbose(addrs); err != nil {
			out.Err = err.Error()
			return out
		}
	}
	sort.SliceStable(txns, func(i, j int) bool {
		a, b := txns[i], txns[j]
		if a.Status.Confirmed != b.Status.Confirmed {
			return !a.Status.Confirmed
		}
		return a.Time > b.Time
	})
	sel := make(map[string]bool, len(addrs))
	for _, a := range addrs {
		sel[a] = true
	}
	walletOf := map[string]int{}
	for i, w := range ws {
		for _, a := range w.Addrs {
			walletOf[a] = i
		}
	}
	rows := make([]TxnRow, 0, len(txns))
	for _, t := range txns {
		r := txnRow(t, sel, ws, walletOf)
		if q.Match(append([]string{r.Txid}, r.Addrs...)...) {
			rows = append(rows, r)
		}
	}
	out.Rows, out.Pager = Paginate(rows, q, TxnPageSize)
	return out
}

func droplets(s string) uint64 {
	c, _ := droplet.FromString(s) //nolint:errcheck
	return c
}

// txnRow reads t from the side of the sel addresses as the Angular wallet's
// history does (src/skycoin-web/src/app/services/wallet/history.service.ts).
// An outgoing transaction sends what goes outside every wallet that signed
// it, and one that sends nothing outside moved coins between own addresses.
func txnRow(t readable.TransactionWithStatusVerbose, sel map[string]bool, ws []HistWallet, walletOf map[string]int) TxnRow {
	tx := t.Transaction
	r := TxnRow{
		Txid:      tx.Hash,
		Time:      unixTime(t.Time),
		Confirmed: t.Status.Confirmed,
		Block:     t.Status.BlockSeq,
	}
	own := func(a string) bool { _, ok := walletOf[a]; return ok }
	var inHours, outHours uint64
	for _, i := range tx.In {
		inHours += i.CalculatedHours
		r.In = append(r.In, TxnIO{Address: i.Address, Coins: i.Coins, Hours: utoa(i.CalculatedHours), Own: own(i.Address)})
	}
	for _, o := range tx.Out {
		outHours += o.Hours
		r.Out = append(r.Out, TxnIO{Address: o.Address, Coins: o.Coins, Hours: utoa(o.Hours), Own: own(o.Address)})
	}
	if inHours > outHours {
		r.Burned = utoa(inHours - outHours)
	} else {
		r.Burned = "0"
	}

	seen := map[string]bool{}
	note := func(a string) {
		if !seen[a] {
			seen[a] = true
			r.Addrs = append(r.Addrs, a)
		}
	}
	var amount, hours uint64
	outgoing := false
	for _, i := range tx.In {
		if sel[i.Address] {
			outgoing = true
		}
	}
	if !outgoing {
		r.Kind = "in"
		for _, o := range tx.Out {
			if sel[o.Address] {
				note(o.Address)
				amount += droplets(o.Coins)
				hours += o.Hours
			}
		}
	} else {
		back := map[string]bool{}
		for _, i := range tx.In {
			if !sel[i.Address] {
				continue
			}
			note(i.Address)
			if k, ok := walletOf[i.Address]; ok {
				for _, a := range ws[k].Addrs {
					back[a] = true
				}
			} else {
				back[i.Address] = true
			}
		}
		r.Kind = "out"
		for _, o := range tx.Out {
			if !back[o.Address] {
				amount += droplets(o.Coins)
				hours += o.Hours
			}
		}
		if amount == 0 {
			r.Kind = "self"
			signer := map[string]bool{}
			for _, i := range tx.In {
				signer[i.Address] = true
			}
			for _, o := range tx.Out {
				if !signer[o.Address] {
					note(o.Address)
					amount += droplets(o.Coins)
					hours += o.Hours
				}
			}
		}
	}
	r.Coins, r.Hours = coins(amount), strconv.FormatUint(hours, 10)
	return r
}

func historyURL(filter string) string {
	if filter == "" {
		return "history"
	}
	return "history?filter=" + url.QueryEscape(filter)
}
