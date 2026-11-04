package ui2

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

// Options configure Handler.
type Options struct {
	// Classic is the Angular wallet, relative to where Handler is mounted.
	// The pages also take the favicon and the CSRF endpoint from there.
	// Empty means "../".
	Classic string
}

// Handler serves the pages, reading the node through n. It expects the mount
// prefix to be stripped, so mounted at /ui2/ it sees /ui2/status as /status.
// The host puts it behind the same auth, host and CSRF checks as the API.
func Handler(n Node, o Options) http.Handler {
	if o.Classic == "" {
		o.Classic = "../"
	}
	files := http.StripPrefix("/static/", http.FileServer(http.FS(staticFS)))

	mux := http.NewServeMux()
	mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		if v := r.URL.Query().Get("v"); v != "" && v == assetVer[strings.TrimPrefix(r.URL.Path, "/static/")] {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		files.ServeHTTP(w, r)
	})
	page := func(f func(r *http.Request) templ.Component) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			render(w, r, o, f(r))
		}
	}
	mux.HandleFunc("/{$}", page(func(*http.Request) templ.Component {
		return WalletsPage(LoadWallets(n))
	}))
	mux.HandleFunc("/wallet", page(func(r *http.Request) templ.Component {
		return WalletPage(LoadWallet(n, r.FormValue("id"), ParseListQuery(r.FormValue)))
	}))
	mux.HandleFunc("/wallet/addresses", page(func(r *http.Request) templ.Component {
		return AddrList(LoadWallet(n, r.FormValue("id"), ParseListQuery(r.FormValue)))
	}))
	mux.HandleFunc("/wallet/txns", page(func(r *http.Request) templ.Component {
		return TxnList(LoadTxns(n, r.FormValue("id"), ParseListQuery(r.FormValue)))
	}))
	mux.HandleFunc("/status", page(func(*http.Request) templ.Component {
		return StatusPage(LoadStatus(n))
	}))
	mux.HandleFunc("/status/live", page(func(*http.Request) templ.Component {
		return StatusLive(LoadStatus(n))
	}))
	mux.HandleFunc("/wallet/xpub", page(func(r *http.Request) templ.Component {
		accts, err := LoadXPubs(n, r.FormValue("id"))
		if err != nil {
			return XPubCard(nil, err.Error())
		}
		return XPubCard(accts, "")
	}))
	mux.HandleFunc("/new", page(func(*http.Request) templ.Component { return CreatePage(false) }))
	mux.HandleFunc("/load", page(func(*http.Request) templ.Component { return CreatePage(true) }))

	// The writes read their fields with PostFormValue, from the body only, so
	// a seed or password in a URL is ignored.
	post := func(f func(w http.ResponseWriter, r *http.Request) templ.Component) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
			if c := f(w, r); c != nil {
				render(w, r, o, c)
			}
		}
	}
	mux.HandleFunc("/create/seed", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		entropy := 128
		if r.PostFormValue("words") == "24" {
			entropy = 256
		}
		seed, err := n.NewSeed(entropy)
		if err != nil {
			return ErrorBox(err.Error())
		}
		return SeedStep(seed)
	}))
	mux.HandleFunc("/create", post(func(w http.ResponseWriter, r *http.Request) templ.Component {
		f := ParseCreateForm(r.PostFormValue)
		if err := f.Check(n); err != nil {
			return formError(err)
		}
		wlt, err := n.CreateWallet(f.Options())
		if err != nil {
			return formError(err)
		}
		w.Header().Set("HX-Redirect", walletURL(wlt.Meta.Filename))
		return nil
	}))
	return mux
}

const maxFormBytes = 64 << 10

func render(w http.ResponseWriter, r *http.Request, o Options, c templ.Component) {
	lang := ""
	if ck, err := r.Cookie(LangCookie); err == nil {
		lang = ck.Value
	}
	ctx := with(r.Context(), Page{
		Base:    BaseFor(r.URL.Path),
		Lang:    PickLang(lang, r.Header.Get("Accept-Language")),
		Classic: o.Classic,
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// The pages show balances and addresses, so no cache may keep them.
	w.Header().Set("Cache-Control", "no-store")
	if err := c.Render(ctx, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
