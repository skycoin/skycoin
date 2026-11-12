// Package ui2 src/ui2/multicoin.go c3-ui-wallet
package ui2

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// CoinCookie holds the coin the pages are about. ui2.js keeps it in step
// with the picker.
const CoinCookie = "ui2-coin"

type coinsKey struct{}

// coinState is the coin a request is about and every coin the host offers.
type coinState struct {
	cur *Coin
	all []Coin
}

// multiCoin serves the pages of every coin a host offers. A Skycoin or fiber
// coin gets the usual pages over its own node API, a Bitcoin or Ethereum
// coin the gateway pages.
type multiCoin struct {
	base Node
	o    Options

	mu    sync.Mutex
	pages map[int]http.Handler
	own   *http.ServeMux
}

func newMultiCoin(base Node, o Options) *multiCoin {
	m := &multiCoin{base: base, o: o, pages: map[int]http.Handler{}}
	m.own = http.NewServeMux()
	m.own.Handle("/static/", assets)
	m.own.HandleFunc("/coins", m.coinsPage)
	m.own.HandleFunc("/coins/endpoints", m.endpointsPage)
	return m
}

// pick is the coin the request is about: the cookie's, else the first one.
func pick(r *http.Request, all []Coin) *Coin {
	if len(all) == 0 {
		return nil
	}
	if ck, err := r.Cookie(CoinCookie); err == nil {
		if id, err := strconv.Atoi(ck.Value); err == nil {
			for i := range all {
				if all[i].ID == id {
					return &all[i]
				}
			}
		}
	}
	return &all[0]
}

func (m *multiCoin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	all, err := m.o.Coins.Coins()
	if err != nil {
		all = nil
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	cur := pick(r, all)
	r = r.WithContext(context.WithValue(r.Context(), coinsKey{}, &coinState{cur: cur, all: all}))
	if r.URL.Path == "/coins" || strings.HasPrefix(r.URL.Path, "/coins/") || strings.HasPrefix(r.URL.Path, "/static/") {
		m.own.ServeHTTP(w, r)
		return
	}
	if cur == nil {
		coinHandler(m.base, m.o).ServeHTTP(w, r)
		return
	}
	m.pagesOf(*cur).ServeHTTP(w, r)
}

// pagesOf are the pages of coin c, made once per coin.
func (m *multiCoin) pagesOf(c Coin) http.Handler {
	m.mu.Lock()
	defer m.mu.Unlock()
	if h, ok := m.pages[c.ID]; ok {
		return h
	}
	var h http.Handler
	switch {
	case c.Gateway() != "":
		h = gatewayHandler(m.o.Coins, c, m.o)
	case c.ID == 0:
		// The host's own /api is coin 0's node.
		h = coinHandler(m.base, m.o)
	default:
		h = coinHandler(m.o.Coins.Node(c.ID), m.o)
	}
	m.pages[c.ID] = h
	return h
}

// CoinsView is the coins page: every coin, the added ones and the Ethereum
// endpoints.
type CoinsView struct {
	Coins  []Coin
	Custom []CustomCoin
	Err    string
	Notice Notice
}

func (m *multiCoin) coinsPage(w http.ResponseWriter, r *http.Request) {
	h := m.o.Coins
	switch r.Method {
	case http.MethodGet:
		v := CoinsView{Notice: ParseNotice(r.FormValue)}
		var err error
		if v.Coins, err = h.Coins(); err != nil {
			v.Err = err.Error()
		}
		if v.Custom, err = h.CustomCoins(); err != nil && v.Err == "" {
			v.Err = err.Error()
		}
		render(w, r, m.o, CoinsPage(v))
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
		var err error
		done := ""
		if id := r.PostFormValue("remove"); id != "" {
			err, done = h.RemoveCoin(id), "coin-removed"
		} else {
			cc := CustomCoin{
				Kind:     r.PostFormValue("kind"),
				Name:     strings.TrimSpace(r.PostFormValue("name")),
				Ticker:   strings.TrimSpace(r.PostFormValue("ticker")),
				NodeURL:  strings.TrimSpace(r.PostFormValue("node_url")),
				Contract: strings.TrimSpace(r.PostFormValue("contract")),
			}
			if d := strings.TrimSpace(r.PostFormValue("decimals")); d != "" {
				if cc.Decimals, err = strconv.Atoi(d); err != nil {
					err = errDecimals
				}
			}
			if err == nil {
				err, done = h.AddCoin(cc), "coin-added"
			}
		}
		if err != nil {
			render(w, r, m.o, formError(err))
			return
		}
		w.Header().Set("HX-Redirect", "coins?done="+done)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// ethCoin is a coin that answers the Ethereum endpoint routes, if any.
func ethCoin(all []Coin) *Coin {
	for i := range all {
		if all[i].Gateway() == "eth" {
			return &all[i]
		}
	}
	return nil
}

// EndpointsView is the Ethereum endpoints page.
type EndpointsView struct {
	Endpoints *Endpoints
	None      bool
	Err       string
	Notice    Notice
}

func (m *multiCoin) endpointsPage(w http.ResponseWriter, r *http.Request) {
	h := m.o.Coins
	all, err := h.Coins()
	if err != nil {
		render(w, r, m.o, EndpointsPage(EndpointsView{Err: err.Error()}))
		return
	}
	c := ethCoin(all)
	switch r.Method {
	case http.MethodGet:
		v := EndpointsView{Notice: ParseNotice(r.FormValue), None: c == nil}
		if c != nil {
			if v.Endpoints, err = h.Endpoints(*c); err != nil {
				v.Err = err.Error()
			}
		}
		render(w, r, m.o, EndpointsPage(v))
	case http.MethodPost:
		if c == nil {
			http.NotFound(w, r)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
		var e Endpoints
		for _, line := range strings.Split(r.PostFormValue("rpc_urls"), "\n") {
			if s := strings.TrimSpace(line); s != "" {
				e.RPCURLs = append(e.RPCURLs, s)
			}
		}
		e.IndexerURL = strings.TrimSpace(r.PostFormValue("indexer_url"))
		if err := h.SetEndpoints(*c, e); err != nil {
			render(w, r, m.o, formError(err))
			return
		}
		w.Header().Set("HX-Redirect", "endpoints?done=endpoints")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

var errDecimals = errors.New("coins.err-decimals")
