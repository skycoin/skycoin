package skydexui

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/skydex/protocol"
	"github.com/skycoin/skycoin/src/ui2"
)

// Page is what every page needs to render.
type Page struct {
	// Base leads from the page back to where the handler is mounted.
	Base string
	// Tab is the open tab, empty on the connect screen.
	Tab       string
	Connected bool
	Ref       string
	Name      string
	// Coins is the market's get_currencies reply, nil when it failed.
	Coins *protocol.GetCurrenciesResponse
}

// Ready reports whether the market accepts trades, and whether that is known.
func (p *Page) Ready() (ready, known bool) {
	if p.Coins == nil {
		return false, false
	}
	return len(p.Coins.SellCoins) > 0, true
}

// Tabs are the client's tabs in the React client's order.
var Tabs = []struct{ ID, Label string }{
	{"market", "Market"},
	{"orders", "My Orders"},
	{"listings", "My Listings"},
	{"history", "History"},
	{"settings", "Settings"},
}

// Handler serves the client pages, talking to the market through m. It
// expects the mount prefix to be stripped, so mounted at /ui2/ it sees
// /ui2/market as /market.
func Handler(m Market) http.Handler {
	h := &handler{m: m}
	mux := http.NewServeMux()
	mux.Handle("/static/", assets)
	mux.HandleFunc("GET /{$}", h.home)
	mux.HandleFunc("POST /connect", h.connect)
	mux.HandleFunc("POST /disconnect", h.disconnect)
	mux.HandleFunc("POST /recent", h.recent)
	mux.HandleFunc("GET /readiness", h.readiness)
	mux.HandleFunc("GET /settings", h.tab("settings"))
	mux.HandleFunc("GET /market/products", h.productGrid)
	mux.HandleFunc("POST /market/buy", h.buy)
	mux.HandleFunc("POST /market/listing", h.createListing)
	mux.HandleFunc("POST /settings", h.register)
	for _, t := range Tabs {
		if t.ID != "settings" {
			mux.HandleFunc("GET /"+t.ID, h.tab(t.ID))
		}
	}
	return sameOriginWrites(mux)
}

type handler struct{ m Market }

func (h *handler) page(r *http.Request, tab string) *Page {
	connected, ref, name := h.m.Status()
	p := &Page{Base: ui2.BaseFor(r.URL.Path), Tab: tab, Connected: connected, Ref: ref, Name: name}
	if connected && tab != "" {
		var c protocol.GetCurrenciesResponse
		if h.m.Call(protocol.TypeGetCurrencies, nil, &c) == nil {
			p.Coins = &c
			if n := strings.TrimSpace(c.MarketName); n != "" {
				p.Name = n
			}
		}
	}
	return p
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	if connected, _, _ := h.m.Status(); connected {
		redirect(w, r, "market")
		return
	}
	render(w, r, GatePage(h.page(r, ""), h.m.Default(), ""))
}

func (h *handler) tab(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := h.page(r, id)
		if !p.Connected {
			redirect(w, r, "")
			return
		}
		render(w, r, TabPage(p, h.tabBody(r, p)))
	}
}

// tabBody is the content of a tab below the banners.
func (h *handler) tabBody(r *http.Request, p *Page) templ.Component {
	switch p.Tab {
	case "settings":
		return SettingsTab(p, nil)
	case "market":
		v := h.products()
		v.Open = r.URL.Query().Get("new") != ""
		return MarketTab(p, v)
	}
	return ClassicTab(p)
}

func (h *handler) connect(w http.ResponseWriter, r *http.Request) {
	ref := strings.TrimSpace(r.FormValue("market_pk"))
	if ref == "" {
		render(w, r, GatePage(h.page(r, ""), ref, "Please enter the market public key."))
		return
	}
	if _, err := h.m.Connect(ref); err != nil {
		render(w, r, GatePage(h.page(r, ""), ref, err.Error()))
		return
	}
	redirect(w, r, "market")
}

func (h *handler) disconnect(w http.ResponseWriter, r *http.Request) {
	h.m.Disconnect()
	redirect(w, r, "")
}

// RecentMarket is an entry of the browser's recent markets list, kept in
// localStorage under exchange:markets as the React client keeps it.
type RecentMarket struct {
	PK   string `json:"pk"`
	Name string `json:"name"`
	TS   int64  `json:"ts"`
}

// recent renders the recent markets the browser posts from its storage.
func (h *handler) recent(w http.ResponseWriter, r *http.Request) {
	var list []RecentMarket
	_ = json.Unmarshal([]byte(r.FormValue("saved")), &list) //nolint:errcheck
	out := list[:0]
	for _, m := range list {
		if m.PK != "" {
			out = append(out, m)
		}
	}
	render(w, r, RecentMarkets(out))
}

func (h *handler) readiness(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, r.URL.Query().Get("tab"))
	render(w, r, Readiness(p, true))
}

func (h *handler) register(w http.ResponseWriter, r *http.Request) {
	p := h.page(r, "settings")
	if !p.Connected {
		redirect(w, r, "")
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sky := strings.TrimSpace(r.PostForm.Get("wallet_sky"))
	var msg *Flash
	if sky == "" {
		msg = &Flash{Kind: "danger", Text: "A Skycoin address is required."}
	} else {
		req := protocol.RegisterRequest{WalletSKY: sky, Wallets: map[string]string{}}
		for k, v := range r.PostForm {
			if cur, ok := strings.CutPrefix(k, "wallet_"); ok && cur != "sky" && len(v) > 0 && strings.TrimSpace(v[0]) != "" {
				req.Wallets[cur] = strings.TrimSpace(v[0])
			}
		}
		if err := h.m.Call(protocol.TypeRegister, req, nil); err != nil {
			msg = &Flash{Kind: "danger", Text: err.Error()}
		} else {
			msg = &Flash{Kind: "success", Text: "Addresses saved. You can now trade."}
			w.Header().Set("HX-Trigger", "wallets-saved")
		}
	}
	if isHTMX(r) {
		render(w, r, FlashBox(msg))
		return
	}
	render(w, r, TabPage(p, SettingsTab(p, msg)))
}

// Flash is a message shown above a form.
type Flash struct {
	Kind string // success, danger or info
	Text string
}

func isHTMX(r *http.Request) bool { return r.Header.Get("HX-Request") == "true" }

// redirect sends the browser to to, a path below the mount point. htmx
// requests get a full page load, as a connect or disconnect changes the whole
// page. htmx resolves its target against the page's <base>, the mount
// point, where a Location resolves against the request path.
func redirect(w http.ResponseWriter, r *http.Request, to string) {
	if isHTMX(r) {
		w.Header().Set("HX-Redirect", "./"+to)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	code := http.StatusFound
	if r.Method == http.MethodPost {
		code = http.StatusSeeOther
	}
	w.Header().Set("Location", ui2.BaseFor(r.URL.Path)+to)
	w.WriteHeader(code)
}

func render(w http.ResponseWriter, r *http.Request, c templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	if err := c.Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// sameOriginWrites refuses requests that can change state when a browser says
// they come from another site, so no other page can place or cancel orders.
// It follows net/http's CrossOriginProtection, which TinyGo's net/http lacks.
func sameOriginWrites(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if crossOriginWrite(r) {
			http.Error(w, "cross-origin request refused", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func crossOriginWrite(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	switch r.Header.Get("Sec-Fetch-Site") {
	case "same-origin", "none":
		return false
	case "":
	default:
		return true
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	return err != nil || u.Host != r.Host
}
