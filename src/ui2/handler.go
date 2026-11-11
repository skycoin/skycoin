package ui2

import (
	"net/http"

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
	mux.HandleFunc("/{$}", page(func(*http.Request) templ.Component {
		if browserMode(n) {
			return BrowserWalletsPage()
		}
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
	mux.HandleFunc("/history", page(func(r *http.Request) templ.Component {
		if browserMode(n) {
			return BrowserHistoryPage(r.FormValue("filter"))
		}
		return HistoryPage(LoadHistory(n, r.FormValue("filter"), ParseListQuery(r.FormValue)))
	}))
	mux.HandleFunc("/history/list", page(func(r *http.Request) templ.Component {
		return TxnList(LoadHistory(n, r.FormValue("filter"), ParseListQuery(r.FormValue)))
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
	mux.HandleFunc("/send", page(func(r *http.Request) templ.Component {
		return SendPage(LoadSend(n, r.FormValue("id"), r.FormValue("from")))
	}))
	mux.HandleFunc("/send/review", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		return ReviewView(Prepare(n, ParseSendForm(r.PostFormValue)))
	}))
	mux.HandleFunc("/send/inject", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		id := r.PostFormValue("id")
		txid, err := n.InjectEncodedTransaction(r.PostFormValue("rawtx"))
		if err != nil {
			return SentView(walletURL(id), "", err.Error())
		}
		return SentView(walletURL(id), txid, "")
	}))

	mux.HandleFunc("/bw/wallet", page(func(r *http.Request) templ.Component { return BrowserWalletPage(r.FormValue("id")) }))
	mux.HandleFunc("/bw/new", page(func(*http.Request) templ.Component { return BrowserCreatePage(false) }))
	mux.HandleFunc("/bw/load", page(func(*http.Request) templ.Component { return BrowserCreatePage(true) }))
	mux.HandleFunc("/bw/list", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		ws, err := ParseBrowserWallets(r.PostFormValue("wallets"))
		if err != nil {
			return ErrorBox(T(r.Context(), err.Error()))
		}
		return BrowserList(LoadBrowserWallets(n, ws))
	}))
	mux.HandleFunc("/bw/wallet/view", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		w, err := ParseBrowserWallet(r.PostFormValue("wallet"))
		if err != nil {
			return bwError(err)
		}
		return BrowserWalletView(LoadBrowserWallet(n, w))
	}))
	mux.HandleFunc("/bw/txns", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		w, err := ParseBrowserWallet(r.PostFormValue("wallet"))
		if err != nil {
			return bwError(err)
		}
		return TxnList(LoadBrowserTxns(n, w, ParseListQuery(r.PostFormValue)))
	}))
	bwHistory := func(r *http.Request) (Txns, error) {
		ws, err := ParseBrowserWallets(r.PostFormValue("wallets"))
		if err != nil {
			return Txns{}, err
		}
		return LoadBrowserHistory(n, ws, r.PostFormValue("filter"), ParseListQuery(r.PostFormValue)), nil
	}
	mux.HandleFunc("/bw/history", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		t, err := bwHistory(r)
		if err != nil {
			return bwError(err)
		}
		return HistoryBody(t)
	}))
	mux.HandleFunc("/bw/history/list", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		t, err := bwHistory(r)
		if err != nil {
			return bwError(err)
		}
		return TxnList(t)
	}))
	mux.HandleFunc("/bw/send", page(func(r *http.Request) templ.Component {
		return BrowserSendPage(r.FormValue("id"), r.FormValue("from"))
	}))
	mux.HandleFunc("/bw/send/review", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		f, err := ParseBrowserSendForm(r.PostFormValue)
		if err != nil {
			return bwError(err)
		}
		return BrowserReview(PrepareBrowser(n, f))
	}))
	mux.HandleFunc("/bw/send/inject", post(func(_ http.ResponseWriter, r *http.Request) templ.Component {
		back := bwWalletURL(r.PostFormValue("id"))
		txid, err := n.InjectEncodedTransaction(r.PostFormValue("rawtx"))
		if err != nil {
			return SentView(back, "", err.Error())
		}
		return SentView(back, txid, "")
	}))
	return mux
}

// maxFormBytes bounds a form, which for a browser wallet carries every
// address it has.
const maxFormBytes = 1 << 20

// browserMode is a node that keeps no wallets, as with -web-wallet, so the
// wallets are kept in the browser.
func browserMode(n Node) bool {
	h, err := n.Health()
	return err == nil && !h.WalletAPIEnabled
}

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
