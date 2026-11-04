package ui2

import (
	"errors"
	"net/url"
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
	}
}

// Send is the send page.
type Send struct {
	ID        string
	Label     string
	Type      string
	Encrypted bool
	Balance   Amount
	From      string
	FromCoins string
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
	if from != "" {
		ab, ok := b.Addresses[from]
		if !ok {
			s.Err = "not an address of this wallet: " + from
			return s
		}
		s.FromCoins = coins(ab.Confirmed.Coins)
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

// Review is a signed transaction waiting for the user to send it.
type Review struct {
	ID       string
	To       string
	Coins    string
	Hours    string
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
	errPending   = errors.New("send.err-pending")
	errNoCoins   = errors.New("send.err-no-coins")
)

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

	req := api.WalletCreateTransactionRequest{WalletID: f.ID, Password: f.Password}
	req.HoursSelection = api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.5"}

	to, amount := f.To, f.Coins
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
		if to, err = unusedAddress(n, w, use, f.Password, wallet.ExternalChain, f.From); err != nil {
			return fail(err)
		}
		amount = coins(ab.Confirmed.Coins)
		req.Addresses = []string{f.From}
	} else {
		if _, err := cipher.DecodeBase58Address(to); err != nil {
			return fail(errBadTo)
		}
		if c, err := droplet.FromString(amount); err != nil || c == 0 {
			return fail(errBadCoins)
		}
		if own[to] {
			r.Warnings = append(r.Warnings, warnToOwn)
		} else {
			tu, err := addressUse(n, []string{to})
			if err != nil {
				return fail(err)
			}
			switch {
			case tu[to].Spent:
				r.Warnings = append(r.Warnings, warnToSpent)
			case tu[to].Seen:
				r.Warnings = append(r.Warnings, warnToReused)
			}
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
		if change, err = unusedAddress(n, w, use, f.Password, wallet.ChangeChain, to); err != nil {
			return fail(err)
		}
	}
	if change != "" {
		req.ChangeAddress = &change
	}
	req.To = []api.Receiver{{Address: to, Coins: amount}}

	txn, err := n.WalletCreateTransaction(req)
	if err != nil {
		return fail(err)
	}
	r.To, r.Coins = to, amount
	for _, o := range txn.Transaction.Out {
		if o.Address == to && r.Hours == "" {
			r.Hours = o.Hours
			continue
		}
		r.Change = append(r.Change, o)
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
