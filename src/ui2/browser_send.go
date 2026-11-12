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
// yet, the default change address and the destination of a move. Opts
// holds the destinations and advanced options, read as for a node wallet.
type BrowserSendForm struct {
	Wallet       BrowserWallet
	Opts         SendForm
	FreshChange  string
	FreshReceive string
}

// ParseBrowserSendForm reads the form through post, a request's
// PostFormValue, and all, its multi value fields.
func ParseBrowserSendForm(post func(string) string, all func(string) []string) (BrowserSendForm, error) {
	w, err := ParseBrowserWallet(post("wallet"))
	if err != nil {
		return BrowserSendForm{}, err
	}
	f := BrowserSendForm{
		Wallet:       w,
		Opts:         ParseSendForm(post),
		FreshChange:  strings.TrimSpace(post("fresh_change")),
		FreshReceive: strings.TrimSpace(post("fresh_receive")),
	}
	f.Opts.ID = w.ID
	f.Opts.ReadRows(all)
	return f, nil
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

// SignDest is a destination in droplets, which bw.js checks the outputs
// against before it signs.
type SignDest struct {
	To    string `json:"to"`
	Coins uint64 `json:"coins"`
}

// SignRequest is the unsigned transaction bw.js signs.
type SignRequest struct {
	Inputs  []SignInput  `json:"inputs"`
	Outputs []SignOutput `json:"outputs"`
	Dests   []SignDest   `json:"dests"`
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
	o := f.Opts
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
	var dests []Dest
	if o.From != "" {
		if !own[o.From] {
			return fail(errBadTo)
		}
		ab := b.Addresses[o.From]
		if ab.Confirmed.Coins == 0 {
			return fail(errNoCoins)
		}
		if ab.Confirmed != ab.Predicted {
			return fail(errPending)
		}
		if _, err := cipher.DecodeBase58Address(f.FreshReceive); err != nil {
			return fail(errBadTo)
		}
		dests = []Dest{{To: f.FreshReceive, Coins: coins(ab.Confirmed.Coins)}}
		req.Addresses = []string{o.From}
	} else {
		var manual bool
		if dests, manual, err = checkDests(n, &o, own, &r.Warnings, &req.HoursSelection); err != nil {
			return fail(err)
		}
		if !manual {
			for i := range dests {
				dests[i].Hours = ""
			}
		}
		for _, a := range o.FromAddrs {
			if !own[a] {
				return fail(errBadSource)
			}
		}
		switch {
		case len(o.FromOuts) > 0:
			req.UxOuts = o.FromOuts
		case len(o.FromAddrs) > 0:
			req.Addresses = o.FromAddrs
		default:
			req.Addresses = addrs
		}
	}

	change := o.Change
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
	var sr SignRequest
	for _, d := range dests {
		req.To = append(req.To, api.Receiver{Address: d.To, Coins: d.Coins, Hours: d.Hours})
		c, err := droplet.FromString(d.Coins)
		if err != nil {
			return fail(errBadCoins)
		}
		sr.Dests = append(sr.Dests, SignDest{To: d.To, Coins: c})
	}

	txn, err := n.CreateTransaction(req)
	if err != nil {
		return fail(err)
	}
	for _, i := range txn.Transaction.In {
		sr.Inputs = append(sr.Inputs, SignInput{Hash: i.UxID, Address: i.Address})
	}
	for i, out := range txn.Transaction.Out {
		c, err := droplet.FromString(out.Coins)
		if err != nil {
			return fail(err)
		}
		h, err := strconv.ParseUint(out.Hours, 10, 64)
		if err != nil {
			return fail(err)
		}
		sr.Outputs = append(sr.Outputs, SignOutput{Address: out.Address, Coins: c, Hours: h})
		if i < len(dests) {
			r.Dests = append(r.Dests, Dest{To: out.Address, Coins: out.Coins, Hours: out.Hours})
			continue
		}
		r.Change = append(r.Change, out)
	}
	r.To, r.Coins, r.Hours = r.Dests[0].To, dests[0].Coins, r.Dests[0].Hours
	r.Fee, r.Raw = txn.Transaction.Fee, txn.EncodedTransaction
	// The id of the signed transaction differs, so none is shown before signing.
	return r, sr
}
