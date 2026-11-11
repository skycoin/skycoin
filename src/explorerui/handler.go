package explorerui

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/skycoin/skycoin/src/ui2"
)

// Options configure Handler.
type Options struct {
	// Root leads from where Handler is mounted to the explorer's root, which
	// serves api.html and search.xml. Empty means "./".
	Root string
	// Classic, when set, is the Angular explorer relative to the mount
	// point. Every page then links to the same page there.
	Classic string
	// Origin is the public scheme://host for canonical links. Empty takes it
	// from each request.
	Origin string
}

// Handler serves the explorer pages, reading the node through n. It expects
// the mount prefix to be stripped, so mounted at /ui2/ it sees
// /ui2/app/blocks/1 as /app/blocks/1. The paths below the mount point are the
// Angular explorer's, so its links keep working.
func Handler(n Node, o Options) http.Handler {
	if o.Root == "" {
		o.Root = "./"
	}
	info := &nodeInfo{}
	mux := http.NewServeMux()
	mux.Handle("/static/", assets)

	page := func(f func(r *http.Request) (templ.Component, Status)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			if l := r.URL.Query().Get("lang"); l != "" {
				setLang(w, r, l)
				return
			}
			c, s := f(r)
			render(w, r, n, info, o, c, s)
		}
	}
	redirect := func(to string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			relRedirect(w, r, to)
		}
	}

	mux.HandleFunc("/{$}", redirect("app/blocks/1"))
	mux.HandleFunc("/app/{$}", redirect("app/blocks/1"))
	mux.HandleFunc("/app/blocks", redirect("app/blocks/1"))
	mux.HandleFunc("/app/blocks/{page}", page(func(r *http.Request) (templ.Component, Status) {
		b := LoadBlocks(n, atoi(r.PathValue("page")))
		return BlocksPage(b), b.Status
	}))
	mux.HandleFunc("/app/block/{id}", page(func(r *http.Request) (templ.Component, Status) {
		b := LoadBlock(n, r.PathValue("id"))
		return BlockPage(b), b.Status
	}))
	mux.HandleFunc("/app/transaction/{txid}", page(func(r *http.Request) (templ.Component, Status) {
		t := LoadTx(n, r.PathValue("txid"))
		return TxPageView(t), t.Status
	}))
	search := func(w http.ResponseWriter, r *http.Request, term string) {
		if to := Search(n, term); to != "" {
			relRedirect(w, r, to)
			return
		}
		render(w, r, n, info, o, SearchPage(term), NotFound)
	}
	mux.HandleFunc("/app/search", func(w http.ResponseWriter, r *http.Request) {
		search(w, r, r.URL.Query().Get("q"))
	})
	mux.HandleFunc("/app/search/{term}", func(w http.ResponseWriter, r *http.Request) {
		search(w, r, r.PathValue("term"))
	})
	return mux
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s) //nolint:errcheck
	return n
}

// relRedirect sends the browser to to, a path below the mount point. The
// Location is relative so the redirect works wherever the handler is mounted.
func relRedirect(w http.ResponseWriter, r *http.Request, to string) {
	w.Header().Set("Location", ui2.BaseFor(r.URL.Path)+to)
	w.WriteHeader(http.StatusFound)
}

// setLang keeps the picked language in a cookie and reloads the page without
// the query, so the address stays clean and shareable.
func setLang(w http.ResponseWriter, r *http.Request, lang string) {
	if catalog.Known(lang) {
		// Secure only over TLS, so a local explorer on plain HTTP keeps the choice.
		http.SetCookie(w, &http.Cookie{Name: LangCookie, Value: lang, Path: "/", MaxAge: 365 * 24 * 3600, HttpOnly: true, Secure: secure(r), SameSite: http.SameSiteLaxMode}) //nolint:gosec
	}
	to := r.URL.Path
	to = to[strings.LastIndex(to, "/")+1:]
	if to == "" {
		to = "./"
	}
	w.Header().Set("Location", to)
	w.WriteHeader(http.StatusFound)
}

func secure(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func origin(r *http.Request, fixed string) string {
	if fixed != "" {
		return strings.TrimRight(fixed, "/")
	}
	scheme := "http"
	if secure(r) {
		scheme = "https"
	}
	host := r.Host
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		host = h
	}
	return scheme + "://" + host
}

func render(w http.ResponseWriter, r *http.Request, n Node, info *nodeInfo, o Options, c templ.Component, s Status) {
	lang := ""
	if ck, err := r.Cookie(LangCookie); err == nil {
		lang = ck.Value
	}
	coin, syncing := info.get(n)
	p := &Page{
		Base:    ui2.BaseFor(r.URL.Path),
		Root:    o.Root,
		Path:    r.URL.Path,
		Origin:  origin(r, o.Origin),
		Lang:    catalog.Pick(lang, r.Header.Get("Accept-Language")),
		Coin:    coin,
		Syncing: syncing,
	}
	if o.Classic != "" {
		p.Classic = o.Classic + strings.TrimPrefix(r.URL.Path, "/")
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Vary", "Cookie, Accept-Language")
	switch s {
	case NotFound:
		w.WriteHeader(http.StatusNotFound)
	case Failed:
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	if err := c.Render(with(r.Context(), p), w); err != nil && s == OK {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
