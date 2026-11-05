package ui2

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// BrowserSendForm is the send form of a browser wallet. FreshChange and
// FreshReceive are addresses bw.js derived for the wallet and has not stored
// yet, the default change address and the destination of a move.
type BrowserSendForm struct {
	Wallet       BrowserWallet
	To           string
	Coins        string
	Change       string
	From         string
	FreshChange  string
	FreshReceive string
}

// ParseBrowserSendForm reads the form through post, a request's
// PostFormValue.
func ParseBrowserSendForm(post func(string) string) (BrowserSendForm, error) {
	w, err := ParseBrowserWallet(post("wallet"))
	if err != nil {
		return BrowserSendForm{}, err
	}
	return BrowserSendForm{
		Wallet:       w,
		To:           strings.TrimSpace(post("to")),
		Coins:        strings.TrimSpace(post("coins")),
		Change:       strings.TrimSpace(post("change")),
		From:         strings.TrimSpace(post("from")),
		FreshChange:  strings.TrimSpace(post("fresh_change")),
		FreshReceive: strings.TrimSpace(post("fresh_receive")),
	}, nil
}

// SignInput and SignOutput are what the wasm cipher's prepareTransaction
// takes, less the secret keys bw.js adds.
type SignInput struct {
	Hash    string `json:"hash"`
	Address string `json:"address"`
}

// SignOutput is an output in droplets and hours, as the cipher takes it.
type SignOutput struct {
	Address string `json:"address"`
	Coins   uint64 `json:"coins"`
	Hours   uint64 `json:"hours"`
}

// SignRequest is the unsigned transaction bw.js signs. To and Coins are the
// destination and droplets it checks the outputs against first.
type SignRequest struct {
	Inputs  []SignInput  `json:"inputs"`
	Outputs []SignOutput `json:"outputs"`
	To      string       `json:"to"`
	Coins   uint64       `json:"coins"`
}

// JSON is the request as the review page carries it.
func (s SignRequest) JSON() string {
	b, err := json.Marshal(s)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// PrepareBrowser checks a browser wallet's send form and has the node build
// the unsigned transaction, which bw.js signs in the page.
func PrepareBrowser(n Node, f BrowserSendForm) (Review, SignRequest) {
	w := f.Wallet
	r := Review{ID: w.ID}
	fail := func(err error) (Review, SignRequest) {
		r.Err = err.Error()
		return r, SignRequest{}
	}
	addrs := w.addrs()
	own := make(map[string]bool, len(addrs))
	for _, a := range addrs {
		own[a] = true
	}
	b, err := n.Balance(addrs)
	if err != nil {
		return fail(err)
	}
	use, err := addressUse(n, addrs)
	if err != nil {
		return fail(err)
	}

	req := api.CreateTransactionRequest{HoursSelection: api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.5"}}
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
		if _, err := cipher.DecodeBase58Address(f.FreshReceive); err != nil {
			return fail(errBadTo)
		}
		to, amount = f.FreshReceive, coins(ab.Confirmed.Coins)
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
		req.Addresses = addrs
	}

	change := f.Change
	if change == "" {
		change = f.FreshChange
	} else if !own[change] {
		r.Warnings = append(r.Warnings, warnChangeForeign)
	}
	if _, err := cipher.DecodeBase58Address(change); err != nil {
		return fail(errBadChange)
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
	req.ChangeAddress = &change
	req.To = []api.Receiver{{Address: to, Coins: amount}}

	txn, err := n.CreateTransaction(req)
	if err != nil {
		return fail(err)
	}
	sr := SignRequest{To: to}
	if sr.Coins, err = droplet.FromString(amount); err != nil {
		return fail(errBadCoins)
	}
	for _, i := range txn.Transaction.In {
		sr.Inputs = append(sr.Inputs, SignInput{Hash: i.UxID, Address: i.Address})
	}
	r.To, r.Coins = to, amount
	for _, o := range txn.Transaction.Out {
		c, err := droplet.FromString(o.Coins)
		if err != nil {
			return fail(err)
		}
		h, err := strconv.ParseUint(o.Hours, 10, 64)
		if err != nil {
			return fail(err)
		}
		sr.Outputs = append(sr.Outputs, SignOutput{Address: o.Address, Coins: c, Hours: h})
		if o.Address == to && r.Hours == "" {
			r.Hours = o.Hours
			continue
		}
		r.Change = append(r.Change, o)
	}
	// The id of the signed transaction differs, so none is shown before signing.
	r.Fee = txn.Transaction.Fee
	return r, sr
}
