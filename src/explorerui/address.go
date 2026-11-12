// Package explorerui src/explorerui/address.go c3-ui-explorer
package explorerui

import (
	"sort"
	"strings"
	"sync"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

// AddressPageSize and manyTxns follow the Angular address page: up to
// manyTxns transactions are read at once, which gives the running balance
// and the totals, and more are read a page at a time.
const (
	AddressPageSize = 25
	manyTxns        = 100
)

// Address is an address's page.
type Address struct {
	Status
	Addr string
	// Count is how many transactions the address has. Many is more than
	// manyTxns, which leaves out the totals and running balances.
	Count int
	Many  bool
	Coins int64
	Hours uint64
	// Pending is what the unconfirmed transactions will change, PendingIn
	// and PendingOut the coins they bring and take.
	PendingCoins, PendingHours int64
	PendingIn, PendingOut      int64
	Received, Sent             int64
	// SentKnown says whether the node answered whether the address has
	// spent an output, Sent being the answer.
	SentKnown, HasSent bool
	Txns               []Tx
	Page, Pages        int
}

// Exposed is an address that has sent coins, so its public key is on the
// chain, and still holds coins.
func (a Address) Exposed() bool { return a.SentKnown && a.HasSent && a.Coins > 0 }

// addrTx is a transaction seen from addr, Balance being what it changed
// there.
func addrTx(t readable.TransactionWithStatusVerbose, addr string) Tx {
	tx := txFrom(t.Transaction.BlockTransactionVerbose)
	tx.Status = t.Status
	tx.Time = unixWhen(t.Transaction.Timestamp)
	tx.HasBalance = true
	for _, in := range tx.In {
		if strings.EqualFold(in.Addr, addr) {
			tx.Balance -= in.Coins
		}
	}
	for _, out := range tx.Out {
		if strings.EqualFold(out.Addr, addr) {
			tx.Balance += out.Coins
		}
	}
	return tx
}

// LoadAddress reads page (1 based) of addr's transactions with its balance,
// following explorer/src/app/components/pages/address-detail.
func LoadAddress(n Node, addr string, page int) Address {
	a := Address{Addr: addr, Page: 1}
	if !isAlnum(addr) {
		a.Status = NotFound
		return a
	}
	// The balance, the pending transactions and the spent outputs do not
	// depend on the list, so they are read while it is.
	var (
		wg      sync.WaitGroup
		bal     *api.BalanceResponse
		balErr  error
		pend    []readable.UnconfirmedTransactionVerbose
		pendErr error
		sent    bool
		sentErr error
	)
	wg.Add(3)
	go func() { defer wg.Done(); bal, balErr = n.AddressBalance(addr) }()
	go func() { defer wg.Done(); pend, pendErr = n.PendingTransactionsVerbose() }()
	go func() { defer wg.Done(); sent, sentErr = n.AddressHasSent(addr) }()
	defer wg.Wait()

	first, err := n.AddressTransactionsPage(addr, 1, 1)
	if err != nil {
		a.Status = Failed
		if notFound(err) {
			a.Status = NotFound
		}
		return a
	}
	a.Count = int(first.PageInfo.TotalPages) //nolint:gosec
	a.Pages = max((a.Count+AddressPageSize-1)/AddressPageSize, 1)
	a.Page = min(max(page, 1), a.Pages)
	a.Many = a.Count > manyTxns

	if a.Many {
		p, err := n.AddressTransactionsPage(addr, a.Page, AddressPageSize)
		if err != nil {
			a.Status = Failed
			return a
		}
		for _, t := range p.Txns {
			a.Txns = append(a.Txns, addrTx(t, addr))
		}
	} else {
		all, err := n.AddressTransactions(addr)
		if err != nil {
			a.Status = Failed
			return a
		}
		sort.SliceStable(all, func(i, j int) bool { return all[i].Transaction.Timestamp < all[j].Transaction.Timestamp })
		var run int64
		txs := make([]Tx, len(all))
		for i, t := range all {
			tx := addrTx(t, addr)
			tx.HasRunning, tx.Before = true, run
			run += tx.Balance
			tx.After = run
			if tx.Status.Confirmed {
				if tx.Balance > 0 {
					a.Received += tx.Balance
				} else {
					a.Sent -= tx.Balance
				}
			}
			txs[len(all)-1-i] = tx
		}
		from := (a.Page - 1) * AddressPageSize
		a.Txns = txs[min(from, len(txs)):min(from+AddressPageSize, len(txs))]
	}

	wg.Wait()
	if balErr != nil {
		a.Status = Failed
		return a
	}
	a.Coins, a.Hours = int64(bal.Confirmed.Coins), bal.Confirmed.Hours //nolint:gosec

	if pendErr == nil {
		for _, u := range pend {
			var coins, hours int64
			for _, in := range u.Transaction.In {
				if in.Address == addr {
					coins -= Droplets(in.Coins)
					hours -= int64(in.CalculatedHours) //nolint:gosec
				}
			}
			for _, out := range u.Transaction.Out {
				if out.Address == addr {
					coins += Droplets(out.Coins)
					hours += int64(out.Hours) //nolint:gosec
				}
			}
			a.PendingCoins += coins
			a.PendingHours += hours
			if coins > 0 {
				a.PendingIn += coins
			} else {
				a.PendingOut -= coins
			}
		}
	}

	a.SentKnown, a.HasSent = sentErr == nil, sent
	return a
}

// UnspentMax is how many outputs show before the rest fold away, as in the
// Angular unspent outputs page.
const UnspentMax = 100

// Unspent is an address's unspent outputs page.
type Unspent struct {
	Status
	Addr    string
	Outputs []readable.UnspentOutput
	Coins   int64
	Hours   uint64
}

// LoadUnspent reads addr's unspent outputs.
func LoadUnspent(n Node, addr string) Unspent {
	u := Unspent{Addr: addr}
	if !isAlnum(addr) {
		u.Status = NotFound
		return u
	}
	o, err := n.AddressOutputs(addr)
	if err != nil {
		u.Status = Failed
		if notFound(err) {
			u.Status = NotFound
		}
		return u
	}
	u.Outputs = o.HeadOutputs
	for _, out := range u.Outputs {
		u.Coins += Droplets(out.Coins)
		u.Hours += out.CalculatedHours
	}
	return u
}

// Unconfirmed is the unconfirmed transactions page.
type Unconfirmed struct {
	Status
	Txns           []Tx
	Size           uint64
	Newest, Oldest When
}

// LoadUnconfirmed reads the node's unconfirmed transactions, newest first.
func LoadUnconfirmed(n Node) Unconfirmed {
	var u Unconfirmed
	pend, err := n.PendingTransactionsVerbose()
	if err != nil {
		u.Status = Failed
		return u
	}
	sort.SliceStable(pend, func(i, j int) bool { return pend[i].Received.After(pend[j].Received) })
	for _, p := range pend {
		tx := txFrom(p.Transaction)
		tx.Time = When{p.Received.UTC()}
		tx.Status.Confirmed = p.IsValid
		u.Size += uint64(p.Transaction.Length)
		u.Txns = append(u.Txns, tx)
	}
	if len(u.Txns) > 0 {
		u.Newest, u.Oldest = u.Txns[0].Time, u.Txns[len(u.Txns)-1].Time
	}
	return u
}

// RichlistPageSize follows the Angular rich list.
const RichlistPageSize = 50

// RichRow is one address of the rich list.
type RichRow struct {
	Rank  int
	Addr  string
	Coins int64
}

// Richlist is a page of the rich list. Last is 0 when the node does not
// say how long the list is.
type Richlist struct {
	Status
	Rows    []RichRow
	Page    int
	Last    int
	HasNext bool
}

// LoadRichlist reads page (1 based) of the rich list. A node without
// offset and total paging is asked for the whole list up to the page.
func LoadRichlist(n Node, page int) Richlist {
	r := Richlist{Page: max(page, 1)}
	start := (r.Page - 1) * RichlistPageSize
	res, err := n.Richlist(&api.RichlistParams{N: RichlistPageSize + 1, Offset: start})
	if err != nil {
		r.Status = Failed
		return r
	}
	entries := res.Richlist
	if res.Total > 0 {
		r.Last = (res.Total + RichlistPageSize - 1) / RichlistPageSize
		r.HasNext = start+min(len(entries), RichlistPageSize) < res.Total
	} else {
		if start > 0 {
			if res, err = n.Richlist(&api.RichlistParams{N: start + RichlistPageSize + 1}); err != nil {
				r.Status = Failed
				return r
			}
			entries = res.Richlist[min(start, len(res.Richlist)):]
		}
		r.HasNext = len(entries) > RichlistPageSize
	}
	for i, e := range entries[:min(len(entries), RichlistPageSize)] {
		r.Rows = append(r.Rows, RichRow{Rank: start + i + 1, Addr: e.Address, Coins: Droplets(e.Coins)})
	}
	return r
}
