package commands

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/cmd/skydex-client/skydexui"
	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/server"
)

// newReadyMarket is a market that accepts trades: SKY is an enabled sell
// coin and BTC a payment currency, as in src/skydex/server's tests.
func newReadyMarket(t *testing.T) (*server.Server, *db.Database) {
	t.Helper()
	database, err := db.New(filepath.Join(t.TempDir(), "m.db"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = database.Close() }) //nolint
	if err := database.Migrate(); err != nil {
		t.Fatal(err)
	}
	if err := database.InitDefaultConfig(); err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{
		"explorer_btc_provider": "esplora",
		"explorer_ltc_provider": "",
		"sky_wallet_seed":       "test escrow seed words for the market",
		"wallet_sky":            "sky-market-wallet",
	} {
		if err := database.SetConfig(k, v); err != nil {
			t.Fatal(err)
		}
	}
	if err := database.SetSellCoinEnabled("SKY", true); err != nil {
		t.Fatal(err)
	}
	return server.New(database, nil, 0), database
}

// newUI2Server serves the client's mux (API and /ui2/) against srv. The
// client never follows redirects, so tests see them.
func newUI2Server(t *testing.T, srv *server.Server, defaultPK string) (*httptest.Server, *http.Client) {
	t.Helper()
	sess := newSession(fakeDialer{srv: srv}, defaultPK)
	mux := http.NewServeMux()
	registerAPI(mux, sess)
	mux.Handle("/ui2/", http.StripPrefix("/ui2", skydexui.Handler(uiMarket{sess})))
	ts := httptest.NewServer(mux)
	t.Cleanup(func() { ts.Close(); sess.close() })
	c := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return ts, c
}

func send(t *testing.T, c *http.Client, method, u string, form url.Values, hdr ...string) (*http.Response, string) {
	t.Helper()
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, u, body) //nolint:noctx
	if err != nil {
		t.Fatal(err)
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for i := 0; i+1 < len(hdr); i += 2 {
		req.Header.Set(hdr[i], hdr[i+1])
	}
	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, u, err)
	}
	defer res.Body.Close() //nolint
	b, _ := io.ReadAll(res.Body) //nolint
	return res, string(b)
}

func mustContain(t *testing.T, body string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(body, w) {
			t.Errorf("page lacks %q", w)
		}
	}
}

func connectUI2(t *testing.T, ts *httptest.Server, c *http.Client) {
	t.Helper()
	res, _ := send(t, c, "POST", ts.URL+"/ui2/connect", url.Values{"market_pk": {testMarketRef}})
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "./market" {
		t.Fatalf("connect = %d to %q, want 303 to ./market", res.StatusCode, res.Header.Get("Location"))
	}
}

// TestUI2Gate covers the connect screen: the default key is pre-filled but
// not dialed, a connect lands on the market tab and a disconnect goes back.
func TestUI2Gate(t *testing.T) {
	ts, c := newUI2Server(t, newTestMarket(t), "the-default-pk")

	res, body := send(t, c, "GET", ts.URL+"/ui2/", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("gate = %d", res.StatusCode)
	}
	mustContain(t, body, `value="the-default-pk"`, `action="connect"`, `data-local="exchange:markets"`, `href="../"`)

	res, body = send(t, c, "POST", ts.URL+"/ui2/connect", url.Values{"market_pk": {" "}})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("empty connect = %d", res.StatusCode)
	}
	mustContain(t, body, "Please enter the market public key.")

	// Tabs need a connection.
	res, _ = send(t, c, "GET", ts.URL+"/ui2/settings", nil)
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "./" {
		t.Fatalf("settings while disconnected = %d to %q", res.StatusCode, res.Header.Get("Location"))
	}

	connectUI2(t, ts, c)
	res, _ = send(t, c, "GET", ts.URL+"/ui2/", nil)
	if res.Header.Get("Location") != "./market" {
		t.Fatalf("gate while connected goes to %q, want ./market", res.Header.Get("Location"))
	}
	_, body = send(t, c, "GET", ts.URL+"/ui2/market", nil)
	mustContain(t, body, "Connected", `data-market-pk="`+testMarketRef+`"`, `class="tab-link active" href="market"`, "Disconnect")
	// The default market has no sell coin, so it warns.
	mustContain(t, body, "This market isn’t accepting trades yet")

	// htmx gets a full page load instead of a redirect.
	res, _ = send(t, c, "POST", ts.URL+"/ui2/disconnect", url.Values{}, "HX-Request", "true")
	if res.Header.Get("HX-Redirect") != "./" {
		t.Fatalf("htmx disconnect HX-Redirect = %q", res.Header.Get("HX-Redirect"))
	}
	var st struct {
		Connected bool `json:"connected"`
	}
	getJSON(t, ts.URL+"/api/status", &st)
	if st.Connected {
		t.Fatal("still connected after disconnect")
	}
}

// TestUI2Recent renders the recent markets the browser posts.
func TestUI2Recent(t *testing.T) {
	ts, c := newUI2Server(t, newTestMarket(t), "")
	saved := `[{"pk":"` + testMarketRef + `","name":"Test market","ts":1},{"pk":""}]`
	_, body := send(t, c, "POST", ts.URL+"/ui2/recent", url.Values{"saved": {saved}})
	mustContain(t, body, "Recent markets", "Test market", `value="`+testMarketRef+`"`, "03aaaabbbb…8899aa", "d ago", `data-remove-market="`+testMarketRef+`"`)
	_, body = send(t, c, "POST", ts.URL+"/ui2/recent", url.Values{"saved": {"[]"}})
	if strings.Contains(body, "Recent markets") {
		t.Error("empty list still renders")
	}
}

// TestUI2Settings registers wallet addresses through the settings form.
func TestUI2Settings(t *testing.T) {
	srv, _ := newReadyMarket(t)
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/ui2/settings", nil)
	mustContain(t, body, `name="wallet_sky"`, `name="wallet_BTC"`, `value="`+testMarketRef+`" readonly`)
	if strings.Contains(body, "data-register-banner") {
		t.Error("settings shows the register banner")
	}

	res, body := send(t, c, "POST", ts.URL+"/ui2/settings", url.Values{"wallet_sky": {""}}, "HX-Request", "true")
	mustContain(t, body, "A Skycoin address is required.")
	if res.Header.Get("HX-Trigger") != "" {
		t.Error("a refused save triggers wallets-saved")
	}

	res, body = send(t, c, "POST", ts.URL+"/ui2/settings", url.Values{
		"wallet_sky": {"xPHfK6xn5LZvWAPB5EV6A62WyzQLywci1J"},
		"wallet_BTC": {"1CLdkSYgwRA4Hm7EUPqaqpAtGG9zgT7dFm"},
	}, "HX-Request", "true")
	mustContain(t, body, "Addresses saved. You can now trade.")
	if res.Header.Get("HX-Trigger") != "wallets-saved" {
		t.Errorf("HX-Trigger = %q, want wallets-saved", res.Header.Get("HX-Trigger"))
	}

	_, body = send(t, c, "POST", ts.URL+"/ui2/settings", url.Values{
		"wallet_sky": {"not-an-address"},
	}, "HX-Request", "true")
	mustContain(t, body, "alert-danger")
}

// TestUI2Readiness polls the readiness banner and the market name.
func TestUI2Readiness(t *testing.T) {
	srv, database := newReadyMarket(t)
	if err := database.SetConfig("market_name", "Test market"); err != nil {
		t.Fatal(err)
	}
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)
	_, body := send(t, c, "GET", ts.URL+"/ui2/readiness?tab=market", nil)
	mustContain(t, body, `id="readiness"`, `hx-swap-oob="true"`, "Test market")
	if strings.Contains(body, "banner-warn") {
		t.Error("a ready market warns")
	}
}

// TestUI2RefusesCrossSiteWrites keeps other sites from driving the client.
func TestUI2RefusesCrossSiteWrites(t *testing.T) {
	ts, c := newUI2Server(t, newTestMarket(t), "")
	res, _ := send(t, c, "POST", ts.URL+"/ui2/connect", url.Values{"market_pk": {testMarketRef}}, "Sec-Fetch-Site", "cross-site")
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-site connect = %d, want 403", res.StatusCode)
	}
	res, _ = send(t, c, "POST", ts.URL+"/ui2/disconnect", url.Values{}, "Origin", "http://evil.example")
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("foreign origin disconnect = %d, want 403", res.StatusCode)
	}
	res, _ = send(t, c, "POST", ts.URL+"/ui2/connect", url.Values{"market_pk": {testMarketRef}}, "Sec-Fetch-Site", "same-origin")
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("same-origin connect = %d, want 303", res.StatusCode)
	}
}
