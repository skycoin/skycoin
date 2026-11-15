// Package ui2 src/ui2/sweep.go c3-ui-wallet
package ui2

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

// Destination choices of the sweep form. "next" is the first account after
// the shared ones and "self" the wallet's own next unused address.
const (
	sweepToNext    = "next"
	sweepToAccount = "account"
	sweepToSelf    = "self"
	sweepToWallet  = "wallet"
	sweepToAddress = "address"
)

// SweepTTLs are the unlock durations offered, in seconds. 0 keeps the
// password until the node stops.
var SweepTTLs = []struct {
	Seconds uint64
	Key     string
}{
	{15 * 60, "sweep.ttl-15m"},
	{60 * 60, "sweep.ttl-1h"},
	{24 * 60 * 60, "sweep.ttl-24h"},
	{7 * 24 * 60 * 60, "sweep.ttl-7d"},
	{0, "sweep.ttl-restart"},
}

// SweepAccount is one bip44 account in the sweep form.
type SweepAccount struct {
	Index uint32
	Name  string
}

// SweepView is the auto-sweep panel of one wallet.
type SweepView struct {
	ID        string
	Bip44     bool
	Encrypted bool
	Accounts  []SweepAccount
	// Wallets are the node's other wallets, for a wallet destination.
	Wallets []WalletRow
	Rule    sweep.Rule
	// To is the destination choice the form shows.
	To     string
	Status api.SweepStatus
	// Saved says the last change went through.
	Saved string
	Err   string
}

// Active says the node sweeps this wallet now.
func (v SweepView) Active() bool { return v.Status.Enabled && v.Status.Paused == "" }

// statusOnly is v for the status block inside the panel, which shows any
// error itself.
func (v SweepView) statusOnly() SweepView {
	v.Err = ""
	return v
}

// Shared says account a is marked as shared.
func (v SweepView) Shared(a uint32) bool { return slices.Contains(v.Rule.SharedAccounts, a) }

// sweepable says the panel applies to a node wallet: one that holds keys and
// is kept on disk.
func sweepable(w Wallet) bool {
	return !w.Meta.Temp && w.Meta.Type != "xpub" && w.Meta.Filename != ""
}

func destinationChoice(d sweep.Destination, bip44 bool) string {
	switch d.Kind {
	case sweep.DestAccount:
		if d.Account != nil {
			return sweepToAccount
		}
		return sweepToNext
	case sweep.DestWallet:
		if d.Wallet != "" {
			return sweepToWallet
		}
		return sweepToSelf
	case sweep.DestAddress:
		return sweepToAddress
	}
	if bip44 {
		return sweepToNext
	}
	return sweepToSelf
}

// LoadSweep reads the wallet, the node's other wallets and the sweep rule and
// status of wallet id.
func LoadSweep(n Node, id string) SweepView {
	v := SweepView{ID: id}
	w, err := n.Wallet(id)
	if err != nil {
		v.Err = err.Error()
		return v
	}
	v.Bip44, v.Encrypted = w.Meta.Type == "bip44", w.Meta.Encrypted
	for _, a := range w.Accounts {
		v.Accounts = append(v.Accounts, SweepAccount{Index: a.Index, Name: a.Name})
	}
	if v.Bip44 && len(v.Accounts) == 0 {
		v.Accounts = []SweepAccount{{Index: 0}}
	}
	if ws, err := n.Wallets(); err == nil {
		for _, o := range ws {
			if o.Meta.Filename != id && !o.Meta.Temp && o.Meta.Type != "xpub" {
				v.Wallets = append(v.Wallets, WalletRow{ID: o.Meta.Filename, Label: o.Meta.Label})
			}
		}
	}
	resp, err := n.WalletSweep(id)
	if err != nil {
		v.Err = err.Error()
		return v
	}
	v.apply(resp)
	return v
}

func (v *SweepView) apply(resp *api.SweepResponse) {
	if resp.Rule != nil {
		v.Rule = *resp.Rule
	} else {
		v.Rule = sweep.Rule{Wallet: v.ID}
	}
	v.To = destinationChoice(v.Rule.Destination, v.Bip44)
	v.Status = resp.Status
}

// Errors of the sweep form, as catalog keys.
var (
	errSweepBurn    = errors.New("sweep.err-burn")
	errSweepAccount = errors.New("sweep.err-account")
	errSweepTo      = errors.New("sweep.err-destination")
)

// ParseSweepRule builds the rule a sweep form posts, on top of the rule the
// node has, so the gap and any field the form leaves out are kept.
func ParseSweepRule(base sweep.Rule, bip44 bool, post func(string) string, posts []string, addrs string) (sweep.Rule, error) {
	r := base
	r.SharedAccounts = nil
	if bip44 {
		for _, s := range posts {
			a, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return r, errSweepAccount
			}
			r.SharedAccounts = append(r.SharedAccounts, uint32(a))
		}
	}
	r.Addresses = strings.Fields(addrs)
	switch post("to") {
	case sweepToNext:
		r.Destination = sweep.Destination{Kind: sweep.DestAccount}
	case sweepToAccount:
		a, err := strconv.ParseUint(post("to_account"), 10, 32)
		if err != nil {
			return r, errSweepAccount
		}
		acct := uint32(a)
		r.Destination = sweep.Destination{Kind: sweep.DestAccount, Account: &acct}
	case sweepToSelf:
		r.Destination = sweep.Destination{Kind: sweep.DestWallet}
	case sweepToWallet:
		if post("to_wallet") == "" {
			return r, errSweepTo
		}
		r.Destination = sweep.Destination{Kind: sweep.DestWallet, Wallet: post("to_wallet")}
	case sweepToAddress:
		addr := strings.TrimSpace(post("to_address"))
		if addr == "" {
			return r, errSweepTo
		}
		r.Destination = sweep.Destination{Kind: sweep.DestAddress, Address: addr}
	default:
		return r, errSweepTo
	}
	r.Enabled = post("enabled") != ""
	if r.Enabled && !base.Enabled && post("burn_ok") == "" {
		return r, errSweepBurn
	}
	return r, nil
}

// sweepHandlers serve the auto-sweep panel under wallet/sweep.
func sweepHandlers(mux *http.ServeMux, n Node, page func(func(*http.Request) templ.Component) http.HandlerFunc,
	post func(func(http.ResponseWriter, *http.Request) templ.Component) http.HandlerFunc) {
	mux.HandleFunc("/wallet/sweep", page(func(r *http.Request) templ.Component {
		return SweepPanel(LoadSweep(n, r.FormValue("id")))
	}))
	mux.HandleFunc("/wallet/sweep/status", page(func(r *http.Request) templ.Component {
		v := SweepView{ID: r.FormValue("id"), Encrypted: r.FormValue("enc") != ""}
		if resp, err := n.WalletSweep(v.ID); err != nil {
			v.Err = err.Error()
		} else {
			v.apply(resp)
		}
		return SweepStatusBlock(v)
	}))
	// Each write renders the whole panel again, with the node's answer or its
	// refusal, and the form keeps what was posted when the node refuses it.
	mux.HandleFunc("/wallet/sweep/rule", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		v := LoadSweep(n, r.PostFormValue("id"))
		if v.Err != "" {
			return SweepPanel(v)
		}
		rule, err := ParseSweepRule(v.Rule, v.Bip44, r.PostFormValue, r.PostForm["shared"], r.PostFormValue("addresses"))
		if err == nil {
			var resp *api.SweepResponse
			if resp, err = n.SetWalletSweep(rule); err == nil {
				v.apply(resp)
				v.Saved = "sweep.done-saved"
				return SweepPanel(v)
			}
		}
		v.Rule, v.To, v.Err = rule, r.PostFormValue("to"), err.Error()
		return SweepPanel(v)
	}))
	mux.HandleFunc("/wallet/sweep/unlock", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		v := LoadSweep(n, r.PostFormValue("id"))
		if v.Err != "" {
			return SweepPanel(v)
		}
		if r.PostFormValue("password") == "" {
			v.Err = errNoPassword.Error()
			return SweepPanel(v)
		}
		ttl, err := strconv.ParseUint(r.PostFormValue("ttl"), 10, 64)
		if err != nil {
			ttl = 0
		}
		resp, err := n.WalletSweepUnlock(api.SweepUnlockRequest{ID: v.ID, Password: r.PostFormValue("password"), TTLSeconds: ttl})
		if err != nil {
			v.Err = err.Error()
			return SweepPanel(v)
		}
		v.apply(resp)
		v.Saved = "sweep.done-unlocked"
		return SweepPanel(v)
	}))
	mux.HandleFunc("/wallet/sweep/lock", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		v := LoadSweep(n, r.PostFormValue("id"))
		if v.Err != "" {
			return SweepPanel(v)
		}
		resp, err := n.WalletSweepLock(v.ID)
		if err != nil {
			v.Err = err.Error()
			return SweepPanel(v)
		}
		v.apply(resp)
		v.Saved = "sweep.done-locked"
		return SweepPanel(v)
	}))
}

func sweepTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02 15:04:05 UTC")
}
