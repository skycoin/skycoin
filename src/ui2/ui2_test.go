package ui2

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

func (f *fakeNode) TransactionsVerbose([]string) ([]readable.TransactionWithStatusVerbose, error) {
	return f.txns, nil
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
