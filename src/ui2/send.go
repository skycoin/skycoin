package ui2

import (
	"errors"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/wallet"
)

// Use is what the blockchain shows about an address.
type Use struct {
	// Seen is an address in any transaction, confirmed or pending.
	Seen bool
	// Spent is an address that signed an input, which put its public key on
	// the blockchain.
	Spent bool
}

// addressUse reads the history of addrs.
func addressUse(n Node, addrs []string) (map[string]Use, error) {
	out := make(map[string]Use, len(addrs))
	if len(addrs) == 0 {
		return out, nil
	}
	want := make(map[string]bool, len(addrs))
	for _, a := range addrs {
		want[a] = true
	}
	txns, err := n.TransactionsVerbose(addrs)
	if err != nil {
		return out, err
	}
	for _, t := range txns {
		for _, i := range t.Transaction.In {
			if want[i.Address] {
				out[i.Address] = Use{Seen: true, Spent: true}
			}
		}
		for _, o := range t.Transaction.Out {
			if want[o.Address] {
				u := out[o.Address]
				u.Seen = true
				out[o.Address] = u
			}
		}
	}
	return out, nil
}

// SendForm is the send form, read from a POST body only.
type SendForm struct {
	ID       string
	To       string
	Coins    string
	Change   string
	Password string
	// From moves every coin of this one address of the wallet to an unused
	// receive address of the same wallet.
	From string
	// Unsigned makes the transaction without signing it, to sign elsewhere.
	Unsigned bool
	// Dests are every destination. Empty means the one in To and Coins.
	Dests []Dest
	// Manual sets the hours of each destination by hand. Otherwise Share
	// is the factor of the hours that go to the destinations.
	Manual bool
	Share  string
	// FromAddrs and FromOuts limit what the transaction may spend.
	FromAddrs []string
	FromOuts  []string
	// Bulk is pasted lines of address, coins and optionally hours, which
	// replace the destination rows when set.
	Bulk string
}

// Dest is one destination of a send.
type Dest struct {
	To    string
	Coins string
	Hours string
}

// ParseSendForm reads the form through post, a request's PostFormValue.
func ParseSendForm(post func(string) string) SendForm {
	return SendForm{
		ID:       post("id"),
		To:       strings.TrimSpace(post("to")),
		Coins:    strings.TrimSpace(post("coins")),
		Change:   strings.TrimSpace(post("change")),
		Password: post("password"),
		From:     strings.TrimSpace(post("from")),
		Manual:   post("hours_mode") == "manual",
		Share:    strings.TrimSpace(post("share")),
		Bulk:     post("bulk"),
	}
}

// ReadRows reads the destination rows and the sources through all, a
// request's multi value form fields, after the form was parsed.
func (f *SendForm) ReadRows(all func(string) []string) {
	to, cs, hs := all("to"), all("coins"), all("hours")
	for i := range to {
		d := Dest{To: strings.TrimSpace(to[i])}
		if i < len(cs) {
			d.Coins = strings.TrimSpace(cs[i])
		}
		if i < len(hs) {
			d.Hours = strings.TrimSpace(hs[i])
		}
		if d.To != "" || d.Coins != "" || d.Hours != "" {
			f.Dests = append(f.Dests, d)
		}
	}
	for _, a := range all("from_addr") {
		f.FromAddrs = append(f.FromAddrs, strings.TrimSpace(a))
	}
	for _, o := range all("from_out") {
		f.FromOuts = append(f.FromOuts, strings.TrimSpace(o))
	}
}

// Errors of the bulk send lines, as catalog keys, as the Angular wallet's
// bulk send dialog reports them.
var (
	errBulkEmpty        = errors.New("send.bulk-empty")
	errBulkInvalid      = errors.New("send.bulk-invalid")
	errBulkInconsistent = errors.New("send.bulk-inconsistent")
)

// parseBulk reads lines of address,coins[,hours]. Every line sets hours or
// none does.
func parseBulk(text string) ([]Dest, bool, error) {
	var out []Dest
	parts := 0
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		if (len(fields) != 2 && len(fields) != 3) || slices.Contains(fields, "") {
			return nil, false, errBulkInvalid
		}
		if parts != 0 && len(fields) != parts {
			return nil, false, errBulkInconsistent
		}
		parts = len(fields)
		d := Dest{To: fields[0], Coins: fields[1]}
		if parts == 3 {
			d.Hours = fields[2]
		}
		out = append(out, d)
	}
	if len(out) == 0 {
		return nil, false, errBulkEmpty
	}
	return out, parts == 3, nil
}

// dests are the destinations the form asks for.
func (f SendForm) dests() ([]Dest, bool, error) {
	if strings.TrimSpace(f.Bulk) != "" {
		return parseBulk(f.Bulk)
	}
	if len(f.Dests) > 0 {
		return f.Dests, f.Manual, nil
	}
	return []Dest{{To: f.To, Coins: f.Coins}}, f.Manual, nil
}

// SendOutput is one unspent output of the send page's source list.
type SendOutput struct {
	Hash    string
	Address string
	Coins   string
	Hours   string
}

// Send is the send page.
type Send struct {
	ID        string
	Label     string
	Type      string
	Encrypted bool
	Balance   Amount
	// Available is the confirmed coins, which "send all" fills in.
	Available string
	From      string
	FromCoins string
	Addrs     []AddressRow
	Outputs   []SendOutput
	Err       string
}

// LoadSend reads what the send page shows about wallet id. from, when set,
// is an address of the wallet whose coins move to a fresh address.
func LoadSend(n Node, id, from string) Send {
	s := Send{ID: id, From: from}
	w, err := n.Wallet(id)
	if err != nil {
		s.Err = err.Error()
		return s
	}
	s.Label, s.Type, s.Encrypted = w.Meta.Label, w.Meta.Type, w.Meta.Encrypted
	b, err := n.WalletBalance(id)
	if err != nil {
		s.Err = err.Error()
		return s
	}
	s.Balance = amountOf(b.BalancePair)
	s.Available = coins(b.Confirmed.Coins)
	if from != "" {
		ab, ok := b.Addresses[from]
		if !ok {
			s.Err = "not an address of this wallet: " + from
			return s
		}
		s.FromCoins = coins(ab.Confirmed.Coins)
		return s
	}
	var addrs []string
	for _, e := range w.Entries {
		ab := b.Addresses[e.Address]
		if ab.Confirmed.Coins == 0 {
			continue
		}
		addrs = append(addrs, e.Address)
		s.Addrs = append(s.Addrs, AddressRow{Address: e.Address, Coins: coins(ab.Confirmed.Coins), Hours: utoa(ab.Confirmed.Hours)})
	}
	if len(addrs) > 0 {
		if o, err := n.OutputsForAddresses(addrs); err == nil {
			for _, x := range o.HeadOutputs {
				s.Outputs = append(s.Outputs, SendOutput{Hash: x.Hash, Address: x.Address, Coins: x.Coins, Hours: utoa(x.CalculatedHours)})
			}
		}
	}
	return s
}

// Warnings the review shows, as catalog keys.
const (
	warnToSpent       = "send.warn-to-spent"
	warnToReused      = "send.warn-to-reused"
	warnToOwn         = "send.warn-to-own"
	warnChangeReused  = "send.warn-change-reused"
	warnChangeForeign = "send.warn-change-foreign"
)

// Review is a signed transaction waiting for the user to send it. To,
// Coins and Hours are the first destination.
type Review struct {
	ID       string
	To       string
	Coins    string
	Hours    string
	Dests    []Dest
	Change   []api.CreatedTransactionOutput
	Fee      string
	Txid     string
	Raw      string
	Warnings []string
	Err      string
}

var (
	errBadTo     = errors.New("send.err-to")
	errBadChange = errors.New("send.err-change")
	errBadCoins  = errors.New("send.err-coins")
	errBadHours  = errors.New("send.err-hours")
	errBadShare  = errors.New("send.err-share")
	errBadSource = errors.New("send.err-source")
	errPending   = errors.New("send.err-pending")
	errNoCoins   = errors.New("send.err-no-coins")
)

func addWarning(ws []string, w string) []string {
	if slices.Contains(ws, w) {
		return ws
	}
	return append(ws, w)
}

// Prepare checks the form, builds and signs the transaction through n
// without sending it, and lists what the user should know first.
func Prepare(n Node, f SendForm) Review {
	r := Review{ID: f.ID}
	fail := func(err error) Review {
		r.Err = err.Error()
		return r
	}
	w, err := n.Wallet(f.ID)
	if err != nil {
		return fail(err)
	}
	b, err := n.WalletBalance(f.ID)
	if err != nil {
		return fail(err)
	}
	own := make(map[string]bool, len(w.Entries))
	all := make([]string, 0, len(w.Entries))
	for _, e := range w.Entries {
		own[e.Address] = true
		all = append(all, e.Address)
	}
	use, err := addressUse(n, all)
	if err != nil {
		return fail(err)
	}

	req := api.WalletCreateTransactionRequest{WalletID: f.ID, Password: f.Password, Unsigned: f.Unsigned}
	req.HoursSelection = api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.5"}

	var dests []Dest
	if f.From != "" {
		if !own[f.From] {
			return fail(errBadTo)
		}
		ab := b.Addresses[f.From]
		if ab.Confirmed.Coins == 0 {
			return fail(errNoCoins)
		}
		if ab.Confirmed != ab.Predicted {
			return fail(errPending)
		}
		to, err := unusedAddress(n, w, use, f.Password, wallet.ExternalChain, f.From)
		if err != nil {
			return fail(err)
		}
		dests = []Dest{{To: to, Coins: coins(ab.Confirmed.Coins)}}
		req.Addresses = []string{f.From}
	} else {
		var manual bool
		if dests, manual, err = f.dests(); err != nil {
			return fail(err)
		}
		if manual {
			req.HoursSelection = api.HoursSelection{Type: "manual"}
		} else if f.Share != "" {
			v, err := strconv.ParseFloat(f.Share, 64)
			if err != nil || v < 0 || v > 1 {
				return fail(errBadShare)
			}
			req.HoursSelection.ShareFactor = f.Share
		}
		for _, d := range dests {
			if _, err := cipher.DecodeBase58Address(d.To); err != nil {
				return fail(errBadTo)
			}
			if c, err := droplet.FromString(d.Coins); err != nil || c == 0 {
				return fail(errBadCoins)
			}
			if manual {
				if _, err := strconv.ParseUint(d.Hours, 10, 64); err != nil {
					return fail(errBadHours)
				}
			}
			if own[d.To] {
				r.Warnings = addWarning(r.Warnings, warnToOwn)
				continue
			}
			tu, err := addressUse(n, []string{d.To})
			if err != nil {
				return fail(err)
			}
			switch {
			case tu[d.To].Spent:
				r.Warnings = addWarning(r.Warnings, warnToSpent)
			case tu[d.To].Seen:
				r.Warnings = addWarning(r.Warnings, warnToReused)
			}
		}
		if !manual {
			for i := range dests {
				dests[i].Hours = ""
			}
		}
		for _, a := range f.FromAddrs {
			if !own[a] {
				return fail(errBadSource)
			}
		}
		// The node takes either addresses or outputs to spend from.
		if len(f.FromOuts) > 0 {
			req.UxOuts = f.FromOuts
		} else {
			req.Addresses = f.FromAddrs
		}
	}

	change := f.Change
	switch {
	case change != "":
		if _, err := cipher.DecodeBase58Address(change); err != nil {
			return fail(errBadChange)
		}
		if !own[change] {
			r.Warnings = append(r.Warnings, warnChangeForeign)
		}
		cu := use[change]
		if !own[change] {
			if m, err := addressUse(n, []string{change}); err == nil {
				cu = m[change]
			}
		}
		if cu.Seen {
			r.Warnings = append(r.Warnings, warnChangeReused)
		}
	case w.Meta.Type != "bip44" && w.Meta.Encrypted:
		// The node picks an unused change address itself except for an
		// encrypted deterministic wallet, which needs one named.
		if change, err = unusedAddress(n, w, use, f.Password, wallet.ChangeChain, dests[0].To); err != nil {
			return fail(err)
		}
	}
	if change != "" {
		req.ChangeAddress = &change
	}
	for _, d := range dests {
		req.To = append(req.To, api.Receiver{Address: d.To, Coins: d.Coins, Hours: d.Hours})
	}

	txn, err := n.WalletCreateTransaction(req)
	if err != nil {
		return fail(err)
	}
	// The node writes the destinations first, in order, and the change after.
	for i, o := range txn.Transaction.Out {
		if i < len(dests) {
			r.Dests = append(r.Dests, Dest{To: o.Address, Coins: o.Coins, Hours: o.Hours})
			continue
		}
		r.Change = append(r.Change, o)
	}
	if len(r.Dests) > 0 {
		r.To, r.Coins, r.Hours = r.Dests[0].To, dests[0].Coins, r.Dests[0].Hours
	}
	r.Fee, r.Txid, r.Raw = txn.Transaction.Fee, txn.Transaction.TxID, txn.EncodedTransaction
	return r
}

// unusedAddress is an address of w that has never been in a transaction,
// other than skip, or a new one on chain when there is none.
func unusedAddress(n Node, w *api.WalletResponse, use map[string]Use, password string, chain wallet.ChainMode, skip string) (string, error) {
	for _, e := range w.Entries {
		if e.Address == skip || use[e.Address].Seen {
			continue
		}
		if e.Change != nil && (*e.Change == 1) != (chain == wallet.ChangeChain) {
			continue
		}
		return e.Address, nil
	}
	opts := []wallet.Option{wallet.OptionGenerateN(1)}
	switch {
	case w.Meta.Type != "bip44":
	case chain == wallet.ChangeChain:
		opts = append(opts, wallet.OptionChange())
	default:
		opts = append(opts, wallet.OptionExternal())
	}
	addrs, err := n.NewWalletAddress(w.Meta.Filename, password, opts...)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", errors.New("the node made no new address")
	}
	return addrs[0], nil
}

func sendURL(id, from string) string {
	u := "send?id=" + url.QueryEscape(id)
	if from != "" {
		u += "&from=" + url.QueryEscape(from)
	}
	return u
}
