package ui2

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
)

func (f *fakeNode) NewWalletAddress(_, _ string, _ ...wallet.Option) ([]string, error) {
	a := newAddr()
	f.newAddrs = append(f.newAddrs, a)
	return []string{a}, nil
}

func (f *fakeNode) WalletCreateTransaction(r api.WalletCreateTransactionRequest) (*api.CreateTransactionResponse, error) {
	f.txReqs = append(f.txReqs, r)
	change := "node-picked-change"
	if r.ChangeAddress != nil {
		change = *r.ChangeAddress
	}
	var t api.CreateTransactionResponse
	t.Transaction.TxID, t.Transaction.Fee, t.EncodedTransaction = "tx1", "3", "RAWTX"
	t.Transaction.Out = []api.CreatedTransactionOutput{
		{Address: r.To[0].Address, Coins: r.To[0].Coins, Hours: "5"},
		{Address: change, Coins: "1.000000", Hours: "4"},
	}
	return &t, nil
}

func (f *fakeNode) InjectEncodedTransaction(raw string) (string, error) {
	f.injected = append(f.injected, raw)
	return "tx1", nil
}

func newAddr() string {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p).String()
}

type sendFixture struct {
	n                          *fakeNode
	spent, unused, change      string
	toSpent, toReused, toFresh string
}

// sendNode is a bip44 wallet whose first receive address has spent and
// still holds coins, a second receive address never used, and a change
// address that received change.
func sendNode(t *testing.T) sendFixture {
	t.Helper()
	x := sendFixture{spent: newAddr(), unused: newAddr(), change: newAddr(), toSpent: newAddr(), toReused: newAddr(), toFresh: newAddr()}
	w := api.WalletResponse{
		Meta: readable.WalletMeta{Filename: "s.wlt", Label: "S", Type: "bip44"},
		Entries: []readable.WalletEntry{
			{Address: x.spent, ChildNumber: u32(0), Change: u32(0)},
			{Address: x.unused, ChildNumber: u32(1), Change: u32(0)},
			{Address: x.change, ChildNumber: u32(0), Change: u32(1)},
		},
	}
	bal := func(c uint64) readable.BalancePair {
		return readable.BalancePair{Confirmed: readable.Balance{Coins: c, Hours: 1}, Predicted: readable.Balance{Coins: c, Hours: 1}}
	}
	b := &api.BalanceResponse{
		BalancePair: bal(5_000_000),
		Addresses:   readable.AddressBalances{x.spent: bal(2_000_000), x.change: bal(3_000_000), x.unused: bal(0)},
	}
	txn := func(in, out []string) readable.TransactionWithStatusVerbose {
		var t readable.TransactionWithStatusVerbose
		t.Status.Confirmed = true
		for _, a := range in {
			t.Transaction.In = append(t.Transaction.In, readable.TransactionInput{Address: a, Coins: "1"})
		}
		for _, a := range out {
			t.Transaction.Out = append(t.Transaction.Out, readable.TransactionOutput{Address: a, Coins: "1"})
		}
		return t
	}
	x.n = &fakeNode{
		wallets:  []api.WalletResponse{w},
		balances: map[string]*api.BalanceResponse{"s.wlt": b},
		txns: []readable.TransactionWithStatusVerbose{
			txn([]string{x.toSpent}, []string{x.spent}),
			txn([]string{x.spent}, []string{x.toReused, x.change}),
		},
	}
	return x
}

func TestExposedAddresses(t *testing.T) {
	x := sendNode(t)
	w := LoadWallet(x.n, "s.wlt", ListQuery{})
	got := map[string]bool{}
	for _, a := range w.Addrs {
		got[a.Address] = a.Exposed
	}
	require.Equal(t, map[string]bool{x.spent: true, x.unused: false, x.change: false}, got)

	body := get(t, Handler(x.n, Options{}), "/wallet/addresses?id=s.wlt").Body.String()
	require.Equal(t, 1, strings.Count(body, "Exposed key"))
	require.Contains(t, body, `href="send?id=s.wlt&amp;from=`+x.spent+`"`)
}

func TestPrepareWarnings(t *testing.T) {
	x := sendNode(t)
	for to, want := range map[string][]string{
		x.toSpent:  {warnToSpent},
		x.toReused: {warnToReused},
		x.toFresh:  nil,
		x.unused:   {warnToOwn},
	} {
		r := Prepare(x.n, SendForm{ID: "s.wlt", To: to, Coins: "1.5", Password: "pw"})
		require.Empty(t, r.Err)
		require.Equal(t, want, r.Warnings, to)
	}
	req := x.n.txReqs[0]
	require.Nil(t, req.ChangeAddress, "a bip44 wallet lets the node pick an unused change address")
	require.Equal(t, api.HoursSelection{Type: "auto", Mode: "share", ShareFactor: "0.5"}, req.HoursSelection)
	require.Equal(t, "pw", req.Password)
	require.Empty(t, req.Addresses)

	r := Prepare(x.n, SendForm{ID: "s.wlt", To: x.toFresh, Coins: "1", Change: x.change})
	require.Equal(t, []string{warnChangeReused}, r.Warnings)
	r = Prepare(x.n, SendForm{ID: "s.wlt", To: x.toFresh, Coins: "1", Change: x.toFresh})
	require.Equal(t, []string{warnChangeForeign}, r.Warnings)
	r = Prepare(x.n, SendForm{ID: "s.wlt", To: x.toFresh, Coins: "1", Change: x.unused})
	require.Empty(t, r.Warnings)

	require.Equal(t, errBadTo.Error(), Prepare(x.n, SendForm{ID: "s.wlt", To: "nope", Coins: "1"}).Err)
	require.Equal(t, errBadCoins.Error(), Prepare(x.n, SendForm{ID: "s.wlt", To: x.toFresh, Coins: "0"}).Err)
}

func TestPrepareEncryptedDeterministic(t *testing.T) {
	x := sendNode(t)
	w := &x.n.wallets[0]
	w.Meta.Type, w.Meta.Encrypted = "deterministic", true
	for i := range w.Entries {
		w.Entries[i].Change, w.Entries[i].ChildNumber = nil, nil
	}
	r := Prepare(x.n, SendForm{ID: "s.wlt", To: x.toFresh, Coins: "1", Password: "pw"})
	require.Empty(t, r.Err)
	req := x.n.txReqs[len(x.n.txReqs)-1]
	require.NotNil(t, req.ChangeAddress)
	require.Equal(t, x.unused, *req.ChangeAddress, "the unused address of the wallet")
	require.Empty(t, x.n.newAddrs)
}

func TestPrepareMove(t *testing.T) {
	x := sendNode(t)
	r := Prepare(x.n, SendForm{ID: "s.wlt", From: x.spent})
	require.Empty(t, r.Err)
	req := x.n.txReqs[0]
	require.Equal(t, []string{x.spent}, req.Addresses, "only the exposed address is spent")
	require.Equal(t, x.unused, req.To[0].Address, "to an unused receive address of the wallet")
	require.Equal(t, "2.000000", req.To[0].Coins, "every coin it holds")

	x.n.wallets[0].Entries = x.n.wallets[0].Entries[:1]
	x.n.wallets[0].Entries = append(x.n.wallets[0].Entries, readable.WalletEntry{Address: x.change, Change: u32(1)})
	r = Prepare(x.n, SendForm{ID: "s.wlt", From: x.spent})
	require.Empty(t, r.Err)
	require.Len(t, x.n.newAddrs, 1, "a new address when none is unused")
	require.Equal(t, x.n.newAddrs[0], x.n.txReqs[1].To[0].Address)

	b := x.n.balances["s.wlt"]
	p := b.Addresses[x.spent]
	p.Predicted.Coins = 0
	b.Addresses[x.spent] = p
	require.Equal(t, errPending.Error(), Prepare(x.n, SendForm{ID: "s.wlt", From: x.spent}).Err)
	require.Equal(t, errBadTo.Error(), Prepare(x.n, SendForm{ID: "s.wlt", From: x.toFresh}).Err)
}

func TestSendHandlers(t *testing.T) {
	x := sendNode(t)
	x.n.wallets[0].Meta.Encrypted = true
	h := Handler(x.n, Options{})

	body := get(t, h, "/send?id=s.wlt").Body.String()
	require.Contains(t, body, `name="password"`)
	require.Contains(t, body, `hx-post="send/review"`)

	rr := postForm(h, "/send/review", url.Values{"id": {"s.wlt"}, "to": {x.toSpent}, "coins": {"1"}, "password": {"secret-pw"}})
	body = rr.Body.String()
	require.Contains(t, body, `name="rawtx" value="RAWTX"`)
	require.Contains(t, body, "public key is public")
	require.NotContains(t, body, "secret-pw")
	require.NotContains(t, body, `name="password"`)
	require.Equal(t, "secret-pw", x.n.txReqs[0].Password)

	rr = postForm(h, "/send/inject", url.Values{"id": {"s.wlt"}, "rawtx": {"RAWTX"}})
	require.Contains(t, rr.Body.String(), "Sent.")
	require.Equal(t, []string{"RAWTX"}, x.n.injected)

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/send/inject?rawtx=RAWTX", nil))
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Len(t, x.n.injected, 1)

	body = get(t, h, "/send?id=s.wlt&from="+x.spent).Body.String()
	require.Contains(t, body, "Move coins to a new address")
	require.Contains(t, body, "2.000000")
}
