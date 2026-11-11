package ui2

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

type fakeNode struct {
	wallets  []api.WalletResponse
	balances map[string]*api.BalanceResponse
	txns     []readable.TransactionWithStatusVerbose
	err      error
	created  []api.CreateWalletOptions
	txReqs   []api.WalletCreateTransactionRequest
	injected []string
	newAddrs []string
	browser  bool
	created2 []api.CreateTransactionRequest
	ops      []string
}

func (f *fakeNode) Health() (*api.HealthResponse, error) {
	h := &api.HealthResponse{CoinName: "skycoin", OpenConnections: 2, IncomingConnections: 1, OutgoingConnections: 1, WalletAPIEnabled: !f.browser}
	h.Version.Version = "0.28.7"
	h.BlockchainMetadata.Head.BkSeq = 123456
	h.Uptime = wh.FromDuration(3600e9)
	return h, nil
}

func (f *fakeNode) BlockchainProgress() (*readable.BlockchainProgress, error) {
	return &readable.BlockchainProgress{Current: 50, Highest: 200}, nil
}

func (f *fakeNode) NetworkConnections(*api.NetworkConnectionsFilter) (*api.Connections, error) {
	c := readable.Connection{Addr: "10.0.0.1:6000", Outgoing: true, Height: 199}
	c.UserAgent.Coin, c.UserAgent.Version = "skycoin", "0.28.7"
	return &api.Connections{Connections: []readable.Connection{c}}, nil
}

func (f *fakeNode) Wallets() ([]api.WalletResponse, error) { return f.wallets, f.err }

func (f *fakeNode) Wallet(id string) (*api.WalletResponse, error) {
	if f.err != nil {
		return nil, f.err
	}
	for i := range f.wallets {
		if f.wallets[i].Meta.Filename == id {
			return &f.wallets[i], nil
		}
	}
	return nil, api.NewClientError("404 Not Found", http.StatusNotFound, "wallet not found")
}

func (f *fakeNode) WalletBalance(id string) (*api.BalanceResponse, error) {
	return f.balances[id], nil
}

func (f *fakeNode) TransactionsVerbose(addrs []string) ([]readable.TransactionWithStatusVerbose, error) {
	want := map[string]bool{}
	for _, a := range addrs {
		want[a] = true
	}
	var out []readable.TransactionWithStatusVerbose
	for _, t := range f.txns {
		hit := false
		for _, i := range t.Transaction.In {
			hit = hit || want[i.Address]
		}
		for _, o := range t.Transaction.Out {
			hit = hit || want[o.Address]
		}
		if hit {
			out = append(out, t)
		}
	}
	return out, nil
}

func u32(n uint32) *uint32 { return &n }

func testNode() *fakeNode {
	w := api.WalletResponse{
		Meta: readable.WalletMeta{Filename: "a.wlt", Label: "Savings", Type: "bip44", Encrypted: true},
		Entries: []readable.WalletEntry{
			{Address: "addrRecv0", ChildNumber: u32(0), Change: u32(0)},
			{Address: "addrRecv1", ChildNumber: u32(1), Change: u32(0)},
			{Address: "addrChange0", ChildNumber: u32(0), Change: u32(1)},
		},
	}
	b := &api.BalanceResponse{
		BalancePair: readable.BalancePair{
			Confirmed: readable.Balance{Coins: 3_500_000, Hours: 10},
			Predicted: readable.Balance{Coins: 3_500_000, Hours: 10},
		},
		Addresses: readable.AddressBalances{
			"addrRecv0": {Confirmed: readable.Balance{Coins: 3_500_000, Hours: 10}, Predicted: readable.Balance{Coins: 3_500_000, Hours: 10}},
		},
	}
	txn := func(id string, time uint64, in, out []string) readable.TransactionWithStatusVerbose {
		var t readable.TransactionWithStatusVerbose
		t.Time = time
		t.Status.Confirmed = true
		t.Transaction.Hash = id
		for _, a := range in {
			t.Transaction.In = append(t.Transaction.In, readable.TransactionInput{Address: a, Coins: "2", CalculatedHours: 4})
		}
		for _, a := range out {
			t.Transaction.Out = append(t.Transaction.Out, readable.TransactionOutput{Address: a, Coins: "1", Hours: 1})
		}
		return t
	}
	return &fakeNode{
		wallets:  []api.WalletResponse{w},
		balances: map[string]*api.BalanceResponse{"a.wlt": b},
		txns: []readable.TransactionWithStatusVerbose{
			txn("received", 100, []string{"stranger"}, []string{"addrRecv0", "stranger"}),
			txn("sent", 300, []string{"addrRecv1"}, []string{"stranger", "addrChange0"}),
			txn("moved", 200, []string{"addrRecv0"}, []string{"addrRecv1", "addrChange0"}),
		},
	}
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, path, nil))
	return rr
}

func TestWalletsPage(t *testing.T) {
	h := Handler(testNode(), Options{})
	rr := get(t, h, "/")
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "no-store", rr.Header().Get("Cache-Control"))
	body := rr.Body.String()
	require.Contains(t, body, "Savings")
	require.Contains(t, body, `href="wallet?id=a.wlt"`)
	require.Contains(t, body, "3.5")
	require.Contains(t, body, `<base href="./">`)
	require.NotContains(t, body, "<script>", "the CSP allows no inline script")
}

func TestWalletAPIOff(t *testing.T) {
	n := testNode()
	n.err = api.NewClientError("403 Forbidden", http.StatusForbidden, "Endpoint is disabled")
	body := get(t, Handler(n, Options{}), "/").Body.String()
	require.Contains(t, body, "does not keep wallets")

	n.err = errors.New("boom")
	body = get(t, Handler(n, Options{}), "/").Body.String()
	require.Contains(t, body, "boom")
	require.NotContains(t, body, "does not keep wallets")
}

func TestWalletPage(t *testing.T) {
	h := Handler(testNode(), Options{})
	body := get(t, h, "/wallet?id=a.wlt").Body.String()
	require.Contains(t, body, "addrChange0")
	require.Contains(t, body, "Change")
	require.Contains(t, body, `hx-get="wallet/txns?id=a.wlt"`)

	body = get(t, h, "/wallet/addresses?id=a.wlt&hide=1").Body.String()
	require.Contains(t, body, "addrRecv0")
	require.NotContains(t, body, "addrRecv1")

	body = get(t, h, "/wallet/addresses?id=a.wlt&q=CHANGE").Body.String()
	require.Contains(t, body, "addrChange0")
	require.NotContains(t, body, "addrRecv0")

	body = get(t, h, "/wallet?id=missing.wlt").Body.String()
	require.Contains(t, body, "wallet not found")
}

func TestTxns(t *testing.T) {
	tx := LoadTxns(testNode(), "a.wlt", ListQuery{})
	require.Empty(t, tx.Err)
	require.Len(t, tx.Rows, 3)
	got := map[string]TxnRow{}
	for _, r := range tx.Rows {
		got[r.Txid] = r
	}
	require.Equal(t, "sent", tx.Rows[0].Txid, "newest first")
	require.Equal(t, "in", got["received"].Kind)
	require.Equal(t, "1.000000", got["received"].Coins)
	require.Equal(t, "out", got["sent"].Kind)
	require.Equal(t, "1.000000", got["sent"].Coins, "the change does not count as sent")
	require.Equal(t, "self", got["moved"].Kind)
	require.Equal(t, "2.000000", got["moved"].Coins)

	body := get(t, Handler(testNode(), Options{}), "/wallet/txns?id=a.wlt").Body.String()
	require.Contains(t, body, "Received")
	require.Contains(t, body, `id="txns"`)
}

func TestStatusPage(t *testing.T) {
	h := Handler(testNode(), Options{})
	body := get(t, h, "/status").Body.String()
	require.Contains(t, body, "123456")
	require.Contains(t, body, "25.0%")
	require.Contains(t, body, "10.0.0.1:6000")
	require.Contains(t, body, "skycoin:0.28.7")
	require.Contains(t, body, `<base href="./">`)

	body = get(t, h, "/status/live").Body.String()
	require.True(t, strings.HasPrefix(body, `<div id="status-live"`), body)
}

func TestStaticAndLang(t *testing.T) {
	h := Handler(testNode(), Options{})
	rr := get(t, h, "/"+Asset("ui2.css"))
	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Header().Get("Cache-Control"), "immutable")

	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	req.Header.Set("Accept-Language", "es-AR,es;q=0.9")
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Contains(t, rr.Body.String(), "Estado del nodo")

	require.Equal(t, http.StatusMethodNotAllowed, func() int {
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/status", nil))
		return rr.Code
	}())
}

func TestBaseFor(t *testing.T) {
	require.Equal(t, "./", BaseFor("/"))
	require.Equal(t, "./", BaseFor("/wallet"))
	require.Equal(t, "../", BaseFor("/wallet/txns"))
}

func TestCatalogsMatch(t *testing.T) {
	en := catalog.Load("en")
	for _, l := range Langs {
		m := catalog.Load(l.Code)
		require.Len(t, m, len(en), l.Code)
		for k := range en {
			require.NotEmpty(t, m[k], "%s %s", l.Code, k)
		}
	}
}

func TestHistory(t *testing.T) {
	n := testNode()
	other := api.WalletResponse{
		Meta:    readable.WalletMeta{Filename: "b.wlt", Label: "Spending", Type: "deterministic"},
		Entries: []readable.WalletEntry{{Address: "stranger"}},
	}
	n.wallets = append(n.wallets, other)

	all := LoadHistory(n, "", ListQuery{})
	require.Empty(t, all.Err)
	require.Len(t, all.Options, 2)
	require.Equal(t, "Savings", all.Options[0].Label, "sorted by label")
	require.Equal(t, "w:a.wlt", all.Options[0].Value)
	require.Equal(t, "a:addrRecv0", all.Options[0].Addrs[0].Value)
	require.Equal(t, "3.500000", all.Options[0].Coins)
	kinds := map[string]TxnRow{}
	for _, r := range all.Rows {
		kinds[r.Txid] = r
	}
	require.Equal(t, "out", kinds["sent"].Kind, "as in the Angular wallet, a payment to another own wallet is sent")

	one := LoadHistory(n, "w:a.wlt", ListQuery{})
	for _, r := range one.Rows {
		kinds[r.Txid] = r
	}
	require.Equal(t, "out", kinds["sent"].Kind)
	require.Equal(t, []string{"addrRecv1"}, kinds["sent"].Addrs)
	require.Equal(t, "1.000000", kinds["sent"].Coins)
	require.Equal(t, "2", kinds["sent"].Burned, "4 calculated hours in, 2 out")
	require.Len(t, kinds["sent"].In, 1)
	require.Len(t, kinds["sent"].Out, 2)
	require.Equal(t, "in", kinds["received"].Kind)
	require.Equal(t, []string{"addrRecv0"}, kinds["received"].Addrs)

	h := Handler(n, Options{})
	body := get(t, h, "/history?filter=w%3Aa.wlt").Body.String()
	require.Contains(t, body, `<option value="w:a.wlt" selected>`)
	require.Contains(t, body, `data-note="sent"`)
	require.Contains(t, body, `hx-get="history/list"`)
	require.Contains(t, body, `class="active">History`)

	body = get(t, h, "/history/list?filter=w%3Aa.wlt&q=addrRecv0").Body.String()
	require.Contains(t, body, `id="txns"`)
	require.Contains(t, body, `data-note="received"`)
	require.NotContains(t, body, `data-note="sent"`)

	body = get(t, h, "/history/list?q=nothing").Body.String()
	require.Contains(t, body, "matching the current filter")
}

func TestBrowserHistory(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})
	body := get(t, h, "/history?filter=w%3A"+x.spent).Body.String()
	require.Contains(t, body, `data-bw-page="list"`)
	require.Contains(t, body, `hx-post="bw/history"`)
	require.Contains(t, body, `value="w:`+x.spent+`"`)

	ws := url.Values{"wallets": {jsonOf(t, []BrowserWallet{browserWallet(x)})}}
	body = postForm(h, "/bw/history", ws).Body.String()
	require.Contains(t, body, `id="hist-form"`)
	require.Contains(t, body, `hx-post="bw/history/list"`)
	require.Contains(t, body, `hx-include="#hist-form, #bw-data"`)
	require.Equal(t, 1, strings.Count(body, ">Sent<"))

	ws.Set("filter", "a:"+x.change)
	body = postForm(h, "/bw/history/list", ws).Body.String()
	require.Equal(t, 1, strings.Count(body, ">Received<"), "the change address received its change")
	require.NotContains(t, body, ">Sent<")

	require.Contains(t, postForm(h, "/bw/history", url.Values{"wallets": {"x"}}).Body.String(), "could not read")
}

func (f *fakeNode) ScanWalletAddresses(id string, n int, _ string) ([]string, error) {
	f.ops = append(f.ops, "scan "+id+" "+itoa(n))
	w, err := f.Wallet(id)
	if err != nil {
		return nil, err
	}
	w.Entries = append(w.Entries, readable.WalletEntry{Address: "scanned"})
	return []string{"scanned"}, nil
}

func (f *fakeNode) UpdateWallet(id, label string) error {
	f.ops = append(f.ops, "rename "+id+" "+label)
	return nil
}

func (f *fakeNode) EncryptWallet(id, password string) (*api.WalletResponse, error) {
	f.ops = append(f.ops, "encrypt "+id+" "+password)
	return f.Wallet(id)
}

func (f *fakeNode) DecryptWallet(id, password string) (*api.WalletResponse, error) {
	f.ops = append(f.ops, "decrypt "+id+" "+password)
	return f.Wallet(id)
}

func TestWalletOptions(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})
	body := get(t, h, "/wallet?id=a.wlt").Body.String()
	require.Contains(t, body, `hx-post="wallet/do/rename"`)
	require.Contains(t, body, `hx-post="wallet/do/decrypt"`, "an encrypted wallet offers decrypt")
	require.NotContains(t, body, `hx-post="wallet/do/encrypt"`)
	require.Contains(t, body, `value="Savings"`)

	do := func(op string, v url.Values) *httptest.ResponseRecorder {
		v.Set("id", "a.wlt")
		return postForm(h, "/wallet/do/"+op, v)
	}
	rr := do("rename", url.Values{"label": {" Rainy day "}})
	require.Equal(t, "wallet?id=a.wlt&done=renamed", rr.Header().Get("HX-Redirect"))
	require.Equal(t, "rename a.wlt Rainy day", n.ops[len(n.ops)-1])

	rr = do("rename", url.Values{"label": {"  "}})
	require.Empty(t, rr.Header().Get("HX-Redirect"))
	require.Contains(t, rr.Body.String(), "Enter a name.")

	rr = do("address", url.Values{"num": {"3"}, "password": {"pw"}})
	require.Equal(t, "wallet?id=a.wlt&done=added&n=1", rr.Header().Get("HX-Redirect"), "the fake makes one address")
	require.Contains(t, do("address", url.Values{"num": {"101"}}).Body.String(), "from 1 to 100")

	rr = do("scan", url.Values{"num": {""}})
	require.Equal(t, "wallet?id=a.wlt&done=scanned&n=1", rr.Header().Get("HX-Redirect"))
	require.Equal(t, "scan a.wlt 20", n.ops[len(n.ops)-1])

	require.Contains(t, do("encrypt", url.Values{"password": {"a"}, "password2": {"b"}}).Body.String(), "do not match")
	rr = do("encrypt", url.Values{"password": {"a"}, "password2": {"a"}})
	require.Equal(t, "wallet?id=a.wlt&done=encrypted", rr.Header().Get("HX-Redirect"))
	require.Contains(t, do("decrypt", url.Values{}).Body.String(), "Enter the password.")
	require.Equal(t, http.StatusNotFound, do("delete", url.Values{}).Code)
	require.Equal(t, http.StatusMethodNotAllowed, get(t, h, "/wallet/do/rename?id=a.wlt&label=x").Code)

	body = get(t, h, "/wallet?id=a.wlt&done=scanned&n=4").Body.String()
	require.Contains(t, body, "Found 4 new used address(es).")
	body = get(t, h, "/wallet?id=a.wlt&done=%3Cb%3E").Body.String()
	require.NotContains(t, body, `class="alert ok notice"`)
}
