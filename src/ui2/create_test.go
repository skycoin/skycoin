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

func (f *fakeNode) CreateWalletTemp(o api.CreateWalletOptions) (*api.WalletResponse, error) {
	f.ops = append(f.ops, "temp "+o.Label)
	f.created = append(f.created, o)
	return &api.WalletResponse{Meta: readable.WalletMeta{Filename: "temp.wlt", Temp: true}}, nil
}

func (f *fakeNode) RecoverWallet(r api.WalletRecoverRequest) (*api.WalletResponse, error) {
	if r.Seed != testSeed {
		return nil, api.NewClientError("400 Bad Request", http.StatusBadRequest, "wallet recovery seed is wrong")
	}
	f.ops = append(f.ops, "recover "+r.ID+" "+r.SeedPassphrase+" "+r.Password)
	return f.Wallet(r.ID)
}

func TestLoadTemporaryAndPassphrase(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})
	body := get(t, h, "/load").Body.String()
	require.Contains(t, body, `name="temp"`)
	require.Contains(t, body, `<option value="pq">`)
	require.NotContains(t, get(t, h, "/new").Body.String(), `name="temp"`)

	form := url.Values{"mode": {"load"}, "label": {"T"}, "type": {"bip44"}, "seed": {testSeed}, "passphrase": {" extra "}, "temp": {"1"}, "password": {"ignored"}}
	rr := postForm(h, "/create", form)
	require.Equal(t, "wallet?id=temp.wlt", rr.Header().Get("HX-Redirect"))
	o := n.created[len(n.created)-1]
	require.Equal(t, " extra ", o.SeedPassphrase, "used exactly as typed")
	require.False(t, o.Encrypt, "a temporary wallet has no password")
	require.Empty(t, o.Password)

	form.Set("type", "deterministic")
	form.Del("temp")
	form.Set("password2", "ignored")
	postForm(h, "/create", form)
	o = n.created[len(n.created)-1]
	require.Empty(t, o.SeedPassphrase, "only a bip44 wallet takes a passphrase")
	require.True(t, o.Encrypt)

	form.Set("type", "pq")
	require.NotEmpty(t, postForm(h, "/create", form).Header().Get("HX-Redirect"))

	n.wallets[0].Meta.Temp = true
	require.Contains(t, get(t, h, "/").Body.String(), ">Temporary<")
}

func TestResetPassword(t *testing.T) {
	n := testNode()
	h := Handler(n, Options{})
	require.Contains(t, get(t, h, "/wallet?id=a.wlt").Body.String(), `href="reset?id=a.wlt"`)
	body := get(t, h, "/reset?id=a.wlt").Body.String()
	require.Contains(t, body, `name="passphrase"`, "a bip44 wallet may have one")
	require.Contains(t, body, "Savings")

	v := url.Values{"id": {"a.wlt"}, "seed": {"  " + testSeed + " "}, "passphrase": {"p"}, "password": {"new"}, "password2": {"new"}}
	rr := postForm(h, "/reset", v)
	require.Equal(t, "wallet?id=a.wlt&done=reset", rr.Header().Get("HX-Redirect"))
	require.Equal(t, "recover a.wlt p new", n.ops[len(n.ops)-1])

	v.Set("password", "")
	v.Set("password2", "")
	require.Equal(t, "wallet?id=a.wlt&done=decrypted", postForm(h, "/reset", v).Header().Get("HX-Redirect"))

	v.Set("seed", "wrong")
	require.Contains(t, postForm(h, "/reset", v).Body.String(), "seed is wrong")
	v.Set("password2", "x")
	require.Contains(t, postForm(h, "/reset", v).Body.String(), "do not match")
	require.Contains(t, get(t, h, "/reset?id=missing.wlt").Body.String(), "wallet not found")
}
