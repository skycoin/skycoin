package ui2

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

const testSeed = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

func (f *fakeNode) WalletXPub(id, path string) (*api.WalletXPub, error) {
	return &api.WalletXPub{XPubKey: "xpub-" + id + "-" + path}, nil
}

func (f *fakeNode) NewSeed(entropy int) (string, error) {
	if entropy == 256 {
		return testSeed + " " + testSeed, nil
	}
	return testSeed, nil
}

func (f *fakeNode) VerifySeed(seed string) (bool, error) {
	n := len(strings.Fields(seed))
	return n == 12 || n == 24, nil
}

func (f *fakeNode) CreateWallet(o api.CreateWalletOptions) (*api.WalletResponse, error) {
	f.created = append(f.created, o)
	return &api.WalletResponse{Meta: readable.WalletMeta{Filename: "new.wlt"}}, nil
}

func postForm(h http.Handler, path string, v url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}

func TestCreateWallet(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})

	body := get(t, h, "/new").Body.String()
	require.Contains(t, body, `hx-post="create/seed"`)
	require.Contains(t, body, `<option value="bip44" selected>`)

	rr := postForm(h, "/create/seed", url.Values{"words": {"24"}})
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, 24, strings.Count(rr.Body.String(), "<li class=\"mono\">"))
	require.Contains(t, rr.Body.String(), `name="seed_confirm"`)

	form := url.Values{
		"mode": {"new"}, "label": {"Mine"}, "type": {"bip44"},
		"seed": {testSeed}, "seed_confirm": {"abandon abandon"},
		"password": {"hunter2-secret"}, "password2": {"hunter2-secret"},
	}
	rr = postForm(h, "/create", form)
	require.Contains(t, rr.Body.String(), "does not match")
	require.NotContains(t, rr.Body.String(), "hunter2-secret", "the password is never sent back")
	require.NotContains(t, rr.Body.String(), "abandon", "the seed is never sent back")
	require.Empty(t, n.created)

	form.Set("seed_confirm", "  "+testSeed+"\n")
	rr = postForm(h, "/create", form)
	require.Equal(t, "wallet?id=new.wlt", rr.Header().Get("HX-Redirect"), rr.Body.String())
	require.Len(t, n.created, 1)
	o := n.created[0]
	require.Equal(t, "bip44", o.Type)
	require.Equal(t, testSeed, o.Seed)
	require.True(t, o.Encrypt)
	require.Equal(t, "hunter2-secret", o.Password)
	require.EqualValues(t, 100, o.ScanN)

	form.Set("password2", "other")
	require.Contains(t, postForm(h, "/create", form).Body.String(), "passwords do not match")

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/create", nil))
	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestCreateIgnoresURL(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})
	q := url.Values{"seed": {testSeed}, "seed_confirm": {testSeed}}
	rr := postForm(h, "/create?"+q.Encode(), url.Values{"mode": {"new"}, "label": {"x"}, "type": {"bip44"}})
	require.Contains(t, rr.Body.String(), "There is no seed")
	require.Empty(t, n.created)
}

func TestLoadWallet(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})
	require.Contains(t, get(t, h, "/load").Body.String(), `name="nonstandard"`)

	odd := url.Values{"mode": {"load"}, "label": {"Old"}, "type": {"deterministic"}, "seed": {"my own words"}}
	require.Contains(t, postForm(h, "/create", odd).Body.String(), "not a standard")

	odd.Set("type", "bip44")
	odd.Set("nonstandard", "1")
	require.Contains(t, postForm(h, "/create", odd).Body.String(), "not a standard", "bip44 needs a standard seed")

	odd.Set("type", "deterministic")
	rr := postForm(h, "/create", odd)
	require.Equal(t, "wallet?id=new.wlt", rr.Header().Get("HX-Redirect"))
	require.False(t, n.created[0].Encrypt)
	require.Zero(t, n.created[0].ScanN)
}

func TestXPubCard(t *testing.T) {
	n := testNode()
	n.wallets[0].Accounts = []readable.WalletAccount{{Index: 0, Name: "main"}, {Index: 1}}
	h := Handler(n, Options{})
	require.Contains(t, get(t, h, "/wallet?id=a.wlt").Body.String(), `hx-get="wallet/xpub?id=a.wlt"`)

	body := get(t, h, "/wallet/xpub?id=a.wlt").Body.String()
	require.Contains(t, body, "Account xPub")
	require.Contains(t, body, `data-copy="xpub-a.wlt-0"`)
	require.Contains(t, body, `data-copy="xpub-a.wlt-1"`)
	require.Contains(t, body, "m/44&#39;/8000&#39;/0&#39;")
	require.Contains(t, body, "xpub-a.wlt-0/1")
	require.Less(t, strings.Index(body, "xpub-a.wlt-0\""), strings.Index(body, "Chain xPubs"), "the account key comes before the chain keys")
}
