package ui2

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

func (f *fakeNode) Balance(addrs []string) (*api.BalanceResponse, error) {
	out := &api.BalanceResponse{Addresses: readable.AddressBalances{}}
	for _, b := range f.balances {
		for _, a := range addrs {
			if p, ok := b.Addresses[a]; ok {
				out.Addresses[a] = p
				out.Confirmed.Coins += p.Confirmed.Coins
				out.Confirmed.Hours += p.Confirmed.Hours
				out.Predicted.Coins += p.Predicted.Coins
				out.Predicted.Hours += p.Predicted.Hours
			}
		}
	}
	return out, nil
}

// browserWallet is the public part bw.js would post for the send fixture's
// wallet, kept in a browser.
func browserWallet(x sendFixture) BrowserWallet {
	return BrowserWallet{
		ID: x.spent, Label: "B", Type: "bip44", Encrypted: true, XPub: "xpub-test",
		Addresses: []BrowserAddress{
			{Address: x.spent, Chain: "external", Index: 0},
			{Address: x.unused, Chain: "external", Index: 1},
			{Address: x.change, Chain: "change", Index: 0},
		},
	}
}

func jsonOf(t *testing.T, v any) string {
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return string(b)
}

func TestParseBrowserWallets(t *testing.T) {
	x := sendNode(t)
	w := browserWallet(x)
	ws, err := ParseBrowserWallets(jsonOf(t, []BrowserWallet{w}))
	require.NoError(t, err)
	require.Equal(t, []BrowserWallet{w}, ws)

	bad := w
	bad.Addresses = append([]BrowserAddress{}, w.Addresses...)
	bad.Addresses[1].Address = "nope"
	_, err = ParseBrowserWallets(jsonOf(t, []BrowserWallet{bad}))
	require.Error(t, err)
	bad = w
	bad.ID = x.unused
	_, err = ParseBrowserWallets(jsonOf(t, []BrowserWallet{bad}))
	require.Error(t, err, "the id is the first address")
	_, err = ParseBrowserWallet("{")
	require.Error(t, err)
}

func TestBrowserPages(t *testing.T) {
	x := sendNode(t)
	h := Handler(x.n, Options{})
	require.Contains(t, get(t, h, "/").Body.String(), `href="wallet?id=s.wlt"`, "a node that keeps wallets lists them")

	x.n.browser = true
	body := get(t, h, "/").Body.String()
	require.Contains(t, body, `data-bw-page="list"`)
	require.Contains(t, body, `hx-post="bw/list"`)
	require.Contains(t, get(t, h, "/bw/wallet?id="+x.spent).Body.String(), `data-bw-id="`+x.spent+`"`)

	w := browserWallet(x)
	body = postForm(h, "/bw/list", url.Values{"wallets": {jsonOf(t, []BrowserWallet{w})}}).Body.String()
	require.Contains(t, body, `href="bw/wallet?id=`+x.spent+`"`)
	require.Contains(t, body, "5.000000")

	body = postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.Contains(t, body, `data-copy="xpub-test"`)
	require.Contains(t, body, "Account xPub")
	require.Equal(t, 1, strings.Count(body, "Exposed key"))
	require.Contains(t, body, "Move coins")
	require.Contains(t, body, `hx-post="bw/txns"`)

	body = postForm(h, "/bw/txns", url.Values{"wallet": {jsonOf(t, w)}}).Body.String()
	require.Contains(t, body, `id="txns"`)
	require.Equal(t, 1, strings.Count(body, ">Received<"))
	require.Equal(t, 1, strings.Count(body, ">Sent<"), "the change does not make a send a receipt")

	require.Contains(t, postForm(h, "/bw/wallet/view", url.Values{"wallet": {"{}"}}).Body.String(), "could not read")
}

func TestBrowserCreateHasNoNames(t *testing.T) {
	h := Handler(sendNode(t).n, Options{})
	for _, path := range []string{"/bw/new", "/bw/load"} {
		body := get(t, h, path).Body.String()
		i, j := strings.Index(body, `<form id="bw-create"`), strings.Index(body, "</form>")
		form := body[i:]
		form = form[:strings.Index(form, "</form>")]
		require.Positive(t, j)
		require.NotRegexp(t, regexp.MustCompile(`\sname=`), form, "%s: a named field could be submitted", path)
		require.Contains(t, form, `id="bw-password"`)
		require.Contains(t, body, "assets/scripts/wasm_exec.js")
	}
}

func TestBrowserOptionsHaveNoNames(t *testing.T) {
	x := sendNode(t)
	x.n.browser = true
	h := Handler(x.n, Options{})
	body := postForm(h, "/bw/wallet/view", url.Values{"wallet": {jsonOf(t, browserWallet(x))}}).Body.String()
	start := strings.Index(body, `id="wallet-options"`)
	require.Positive(t, start)
	opts := body[start : start+strings.Index(body[start:], "</section>")]
	for _, op := range []string{"rename", "address", "seed", "password", "reset", "delete"} {
		require.Contains(t, opts, `data-bw-op="`+op+`"`)
	}
	require.NotRegexp(t, regexp.MustCompile(`\sname=`), opts, "no option field can be posted")
	require.NotContains(t, opts, "<form class=\"form\" hx-", "bw.js handles them")

	page := get(t, h, "/bw/wallet?id="+x.spent+"&done=reset").Body.String()
	require.Contains(t, page, "Password successfully changed.")
	require.Contains(t, page, `data-msg-wrong-seed=`)
}
