// Package ui2 src/ui2/gateway.go c3-ui-wallet
package ui2

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

// GwWalletRow is a server wallet of a gateway coin in the wallet list.
type GwWalletRow struct {
	ID        string
	Label     string
	Encrypted bool
	Addresses int
	Balance   string
	Err       string
}

// GwWallets is the wallet list of a gateway coin.
type GwWallets struct {
	Coin  Coin
	Rows  []GwWalletRow
	Total string
	Err   string
}

// GwAddress is one address of a gateway wallet with its balance.
type GwAddress struct {
	Address string
	Balance string
	Empty   bool
}

// GwWallet is the page of one gateway wallet.
type GwWallet struct {
	Coin      Coin
	ID        string
	Label     string
	Encrypted bool
	Balance   string
	Addrs     []GwAddress
	Notice    Notice
	Err       string
}

// GwTxRow is one transaction of a gateway wallet, from the wallet's side.
type GwTxRow struct {
	Txid      string
	Time      string
	Confirmed bool
	Failed    bool
	Kind      string
	Amount    string
	Fee       string
}

// GwTxns is a gateway wallet's history.
type GwTxns struct {
	Coin Coin
	Rows []GwTxRow
	Err  string
}

func walletAddrs(w *api.WalletResponse) []string {
	out := make([]string, 0, len(w.Entries))
	for _, e := range w.Entries {
		out = append(out, e.Address)
	}
	return out
}

// LoadGwWallets reads the coin's server wallets and their balances.
func LoadGwWallets(h *CoinHost, c Coin) GwWallets {
	out := GwWallets{Coin: c}
	ws, err := h.Node(c.ID).Wallets()
	if err != nil {
		out.Err = err.Error()
		return out
	}
	total := new(big.Int)
	for _, w := range ws {
		r := GwWalletRow{ID: w.Meta.Filename, Label: w.Meta.Label, Encrypted: w.Meta.Encrypted, Addresses: len(w.Entries)}
		if b, err := h.Balance(c, walletAddrs(&w)); err != nil {
			r.Err = err.Error()
		} else {
			r.Balance = b.Confirmed.Coins.Format(c.Places())
			total.Add(total, &b.Confirmed.Coins.Int)
		}
		out.Rows = append(out.Rows, r)
	}
	sort.SliceStable(out.Rows, func(i, j int) bool { return strings.ToLower(out.Rows[i].Label) < strings.ToLower(out.Rows[j].Label) })
	out.Total = formatUnits(total, c.Places())
	return out
}

// LoadGwWallet reads gateway wallet id with its address balances.
func LoadGwWallet(h *CoinHost, c Coin, id string) GwWallet {
	out := GwWallet{Coin: c, ID: id}
	w, err := h.Node(c.ID).Wallet(id)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	out.Label, out.Encrypted = w.Meta.Label, w.Meta.Encrypted
	addrs := walletAddrs(w)
	b, err := h.Balance(c, addrs)
	if err != nil {
		out.Err = err.Error()
		b = &GwBalance{}
	}
	out.Balance = b.Confirmed.Coins.Format(c.Places())
	for _, a := range addrs {
		ab := b.Addresses[a].Confirmed.Coins
		out.Addrs = append(out.Addrs, GwAddress{Address: a, Balance: ab.Format(c.Places()), Empty: ab.Sign() == 0})
	}
	return out
}

// LoadGwTxns reads the history of gateway wallet id. A transaction that
// spends from the wallet sends what goes to other addresses, and one that
// sends nothing outside moved coins between the wallet's own addresses.
func LoadGwTxns(h *CoinHost, c Coin, id string) GwTxns {
	out := GwTxns{Coin: c}
	w, err := h.Node(c.ID).Wallet(id)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	addrs := walletAddrs(w)
	own := map[string]bool{}
	for _, a := range addrs {
		own[strings.ToLower(a)] = true
	}
	txs, err := h.History(c, addrs)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	for _, t := range txs {
		r := GwTxRow{Txid: t.TxID, Confirmed: t.Confirmations > 0 || t.BlockHeight > 0, Failed: t.Failed}
		if t.Timestamp > 0 {
			r.Time = time.Unix(t.Timestamp, 0).UTC().Format("2006-01-02 15:04 UTC")
		}
		outgoing := false
		for _, i := range t.Inputs {
			outgoing = outgoing || own[strings.ToLower(i.Address)]
		}
		sum := new(big.Int)
		for _, o := range t.Outputs {
			if own[strings.ToLower(o.Address)] != outgoing {
				sum.Add(sum, &o.Value.Int)
			}
		}
		switch {
		case !outgoing:
			r.Kind = "in"
		case sum.Sign() == 0:
			r.Kind = "self"
			for _, o := range t.Outputs {
				sum.Add(sum, &o.Value.Int)
			}
		default:
			r.Kind = "out"
		}
		r.Amount = formatUnits(sum, c.Places())
		if outgoing && t.Fee.Sign() > 0 {
			r.Fee = t.Fee.Format(gwFeePlaces(c))
		}
		out.Rows = append(out.Rows, r)
	}
	return out
}

// gwFeePlaces are the decimal places of a fee: an ERC-20 token's fee is
// paid in Ether.
func gwFeePlaces(c Coin) int {
	if c.Gateway() == "eth" {
		return 18
	}
	return c.Places()
}

// GwSendView is the send page or its review.
type GwSendView struct {
	Coin      Coin
	ID        string
	Label     string
	Encrypted bool
	Balance   string
	To        string
	Amount    string
	Preset    string
	Fees      map[string]FeePreset
	FeeText   string
	Err       string
}

// presetNames are the fee presets in the order the form offers them.
var presetNames = []string{"slow", "normal", "fast"}

// feeText describes a preset: sat/byte for Bitcoin, the most the fee can
// be for an Ethereum coin.
func feeText(c Coin, p FeePreset) string {
	if c.Gateway() == "btc" {
		return itoa(int(p.SatPerByte)) + " sat/byte"
	}
	return "≤ " + p.MaxFee.Format(18) + " ETH"
}

// LoadGwSend reads what the send page of wallet id shows.
func LoadGwSend(h *CoinHost, c Coin, id string) GwSendView {
	v := GwSendView{Coin: c, ID: id, Preset: "normal"}
	w := LoadGwWallet(h, c, id)
	v.Label, v.Encrypted, v.Balance, v.Err = w.Label, w.Encrypted, w.Balance, w.Err
	if f, err := h.Fee(c); err != nil {
		if v.Err == "" {
			v.Err = err.Error()
		}
	} else {
		v.Fees = f.Presets
	}
	return v
}

var (
	errGwPreset = errors.New("gw.err-preset")
	errGwTo     = errors.New("send.err-to")
)

// prepareGwSend checks the form and finds its fee preset. It reads the
// presets again rather than taking fee values from the form.
func prepareGwSend(h *CoinHost, c Coin, post func(string) string) (GwSendView, GwSend, error) {
	v := LoadGwSend(h, c, post("id"))
	v.To, v.Amount, v.Preset = strings.TrimSpace(post("to")), strings.TrimSpace(post("amount")), post("preset")
	if v.Err != "" {
		return v, GwSend{}, errors.New(v.Err)
	}
	if v.To == "" || strings.ContainsAny(v.To, " \t") {
		return v, GwSend{}, errGwTo
	}
	amount, err := parseUnits(v.Amount, c.Places())
	if err != nil {
		return v, GwSend{}, err
	}
	p, ok := v.Fees[v.Preset]
	if !ok {
		return v, GwSend{}, errGwPreset
	}
	v.FeeText = feeText(c, p)
	return v, GwSend{WalletID: v.ID, To: v.To, Amount: amount, Fee: p}, nil
}

// gatewayHandler serves the pages of a Bitcoin or Ethereum coin, whose
// server wallets the host keeps and whose sends it signs and broadcasts.
func gatewayHandler(h *CoinHost, c Coin, o Options) http.Handler {
	node := h.Node(c.ID)
	mux := http.NewServeMux()
	mux.Handle("/static/", assets)
	page := func(f func(r *http.Request) templ.Component) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			render(w, r, o, f(r))
		}
	}
	post := func(f func(w http.ResponseWriter, r *http.Request) templ.Component) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
			if cmp := f(w, r); cmp != nil {
				render(w, r, o, cmp)
			}
		}
	}
	mux.HandleFunc("/{$}", page(func(*http.Request) templ.Component { return GwWalletsPage(LoadGwWallets(h, c)) }))
	mux.HandleFunc("/wallet", page(func(r *http.Request) templ.Component {
		w := LoadGwWallet(h, c, r.FormValue("id"))
		w.Notice = ParseNotice(r.FormValue)
		return GwWalletPage(w)
	}))
	mux.HandleFunc("/gw/txns", page(func(r *http.Request) templ.Component { return GwTxnList(LoadGwTxns(h, c, r.FormValue("id"))) }))
	mux.HandleFunc("/wallet/do/{op}", post(walletOpHandler(node)))
	mux.HandleFunc("/new", page(func(*http.Request) templ.Component { return CreatePage(false) }))
	mux.HandleFunc("/load", page(func(*http.Request) templ.Component { return CreatePage(true) }))
	mux.HandleFunc("/create/seed", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		entropy := 128
		if r.PostFormValue("words") == "24" {
			entropy = 256
		}
		seed, err := node.NewSeed(entropy)
		if err != nil {
			return ErrorBox(err.Error())
		}
		return SeedStep(seed)
	}))
	mux.HandleFunc("/create", post(func(w http.ResponseWriter, r *http.Request) templ.Component {
		f := ParseCreateForm(r.PostFormValue)
		f.Type = "bip44"
		if err := f.Check(node); err != nil {
			return formError(err)
		}
		id, err := h.CreateGwWallet(c, f, r.PostFormValue("legacy") != "")
		if err != nil {
			return formError(err)
		}
		w.Header().Set("HX-Redirect", walletURL(id))
		return nil
	}))
	mux.HandleFunc("/send", page(func(r *http.Request) templ.Component { return GwSendPage(LoadGwSend(h, c, r.FormValue("id"))) }))
	mux.HandleFunc("/send/review", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		v, _, err := prepareGwSend(h, c, r.PostFormValue)
		if err != nil {
			return formError(err)
		}
		return GwReview(v)
	}))
	mux.HandleFunc("/send/do", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		v, s, err := prepareGwSend(h, c, r.PostFormValue)
		if err != nil {
			return formError(err)
		}
		s.Password = r.PostFormValue("password")
		txid, err := h.Send(c, s)
		if err != nil {
			return formError(err)
		}
		return SentView(walletURL(v.ID), txid, "")
	}))
	mux.HandleFunc("/statusbar", page(func(*http.Request) templ.Component { return templ.NopComponent }))
	return mux
}

func gwMeta(w GwWallet) readable.WalletMeta {
	return readable.WalletMeta{Filename: w.ID, Label: w.Label, Type: "bip44", Encrypted: w.Encrypted}
}

// qrScheme is the URI scheme of a gateway coin's payment QR codes.
func qrScheme(c Coin) string {
	if c.Gateway() == "btc" {
		return "bitcoin"
	}
	return "ethereum"
}

func removeVals(id string) string {
	b, _ := json.Marshal(map[string]string{"remove": id}) //nolint:errcheck
	return string(b)
}

func joinLines(s []string) string { return strings.Join(s, "\n") }

func joinComma(s []string) string { return strings.Join(s, ", ") }
