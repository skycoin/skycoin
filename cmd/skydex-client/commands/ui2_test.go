package commands

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/market"
	"github.com/skycoin/skycoin/src/skydex/protocol"
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

// newUI2Server serves the client's whole UI and API against srv. The
// client never follows redirects, so tests see them.
func newUI2Server(t *testing.T, srv *server.Server, defaultPK string) (*httptest.Server, *http.Client) {
	t.Helper()
	sess := newSession(fakeDialer{srv: srv}, defaultPK)
	h, err := uiHandler(sess)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewServer(h)
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
	defer res.Body.Close()       //nolint
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
	res, _ := send(t, c, "POST", ts.URL+"/connect", url.Values{"market_pk": {testMarketRef}})
	if res.StatusCode != http.StatusSeeOther || res.Header.Get("Location") != "./market" {
		t.Fatalf("connect = %d to %q, want 303 to ./market", res.StatusCode, res.Header.Get("Location"))
	}
}

// TestUI2Gate covers the connect screen: the default key is pre-filled but
// not dialed, a connect lands on the market tab and a disconnect goes back.
func TestUI2Gate(t *testing.T) {
	ts, c := newUI2Server(t, newTestMarket(t), "the-default-pk")

	res, body := send(t, c, "GET", ts.URL+"/", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("gate = %d", res.StatusCode)
	}
	mustContain(t, body, `value="the-default-pk"`, `action="connect"`, `data-local="exchange:markets"`, `href="angular/"`)

	res, body = send(t, c, "POST", ts.URL+"/connect", url.Values{"market_pk": {" "}})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("empty connect = %d", res.StatusCode)
	}
	mustContain(t, body, "Please enter the market public key.")

	// Tabs need a connection.
	res, _ = send(t, c, "GET", ts.URL+"/settings", nil)
	if res.StatusCode != http.StatusFound || res.Header.Get("Location") != "./" {
		t.Fatalf("settings while disconnected = %d to %q", res.StatusCode, res.Header.Get("Location"))
	}

	connectUI2(t, ts, c)
	res, _ = send(t, c, "GET", ts.URL+"/", nil)
	if res.Header.Get("Location") != "./market" {
		t.Fatalf("gate while connected goes to %q, want ./market", res.Header.Get("Location"))
	}
	_, body = send(t, c, "GET", ts.URL+"/market", nil)
	mustContain(t, body, "Connected", `data-market-pk="`+testMarketRef+`"`, `class="tab-link active" href="market"`, "Disconnect")
	// The default market has no sell coin, so it warns.
	mustContain(t, body, "This market isn’t accepting trades yet")

	// htmx gets a full page load instead of a redirect.
	res, _ = send(t, c, "POST", ts.URL+"/disconnect", url.Values{}, "HX-Request", "true")
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
	_, body := send(t, c, "POST", ts.URL+"/recent", url.Values{"saved": {saved}})
	mustContain(t, body, "Recent markets", "Test market", `value="`+testMarketRef+`"`, "03aaaabbbb…8899aa", "d ago", `data-remove-market="`+testMarketRef+`"`)
	_, body = send(t, c, "POST", ts.URL+"/recent", url.Values{"saved": {"[]"}})
	if strings.Contains(body, "Recent markets") {
		t.Error("empty list still renders")
	}
}

// TestUI2Settings registers wallet addresses through the settings form.
func TestUI2Settings(t *testing.T) {
	srv, _ := newReadyMarket(t)
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/settings", nil)
	mustContain(t, body, `name="wallet_sky"`, `name="wallet_BTC"`, `value="`+testMarketRef+`" readonly`)
	if strings.Contains(body, "data-register-banner") {
		t.Error("settings shows the register banner")
	}

	res, body := send(t, c, "POST", ts.URL+"/settings", url.Values{"wallet_sky": {""}}, "HX-Request", "true")
	mustContain(t, body, "A Skycoin address is required.")
	if res.Header.Get("HX-Trigger") != "" {
		t.Error("a refused save triggers wallets-saved")
	}

	res, body = send(t, c, "POST", ts.URL+"/settings", url.Values{
		"wallet_sky": {"xPHfK6xn5LZvWAPB5EV6A62WyzQLywci1J"},
		"wallet_BTC": {"1CLdkSYgwRA4Hm7EUPqaqpAtGG9zgT7dFm"},
	}, "HX-Request", "true")
	mustContain(t, body, "Addresses saved. You can now trade.")
	if res.Header.Get("HX-Trigger") != "wallets-saved" {
		t.Errorf("HX-Trigger = %q, want wallets-saved", res.Header.Get("HX-Trigger"))
	}

	_, body = send(t, c, "POST", ts.URL+"/settings", url.Values{
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
	_, body := send(t, c, "GET", ts.URL+"/readiness?tab=market", nil)
	mustContain(t, body, `id="readiness"`, `hx-swap-oob="true"`, "Test market")
	if strings.Contains(body, "banner-warn") {
		t.Error("a ready market warns")
	}
}

// TestUI2RefusesCrossSiteWrites keeps other sites from driving the client.
func TestUI2RefusesCrossSiteWrites(t *testing.T) {
	ts, c := newUI2Server(t, newTestMarket(t), "")
	res, _ := send(t, c, "POST", ts.URL+"/connect", url.Values{"market_pk": {testMarketRef}}, "Sec-Fetch-Site", "cross-site")
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("cross-site connect = %d, want 403", res.StatusCode)
	}
	res, _ = send(t, c, "POST", ts.URL+"/disconnect", url.Values{}, "Origin", "http://evil.example")
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("foreign origin disconnect = %d, want 403", res.StatusCode)
	}
	res, _ = send(t, c, "POST", ts.URL+"/connect", url.Values{"market_pk": {testMarketRef}}, "Sec-Fetch-Site", "same-origin")
	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("same-origin connect = %d, want 303", res.StatusCode)
	}
}

const ui2SellerPK = "0311111111111111111111111111111111111111111111111111111111111111aa"

// seedProduct lists an active product from another seller.
func seedProduct(t *testing.T, srv *server.Server, database *db.Database, id string) {
	t.Helper()
	ce, se := net.Pipe()
	go srv.Serve(se, ui2SellerPK)
	conn := market.NewConn(ce)
	defer conn.Close() //nolint
	resp, err := conn.Do(protocol.TypeRegister, protocol.RegisterRequest{WalletSKY: "FbmPhy5bhsMX8JNeAQdQ4W3DeCYFH97FBg", WalletBTC: "19wiPpvWPxHEmcANKzwUjiRgABQk8wZzCp"})
	if err != nil || resp.IsError() {
		t.Fatalf("register seller: %v %s", err, resp.Data)
	}
	if err := database.CreateProduct(&db.Product{
		ID: id, SellerPubKey: ui2SellerPK, SellCoin: "SKY", Amount: 10, Price: 0.0002, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}
}

func registerUI2(t *testing.T, ts *httptest.Server, c *http.Client) {
	t.Helper()
	_, body := send(t, c, "POST", ts.URL+"/settings", url.Values{
		"wallet_sky": {"xPHfK6xn5LZvWAPB5EV6A62WyzQLywci1J"},
		"wallet_BTC": {"1CLdkSYgwRA4Hm7EUPqaqpAtGG9zgT7dFm"},
	}, "HX-Request", "true")
	mustContain(t, body, "Addresses saved.")
}

// TestUI2MarketBuy lists products and buys one through the market tab.
func TestUI2MarketBuy(t *testing.T) {
	srv, database := newReadyMarket(t)
	seedProduct(t, srv, database, "prod-ui2")
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/market", nil)
	mustContain(t, body, "Available products", "10 SKY", "0.0002 BTC", "03111111…11aa", `value="prod-ui2"`, `href="market?new=1"`, `id="qr-dialog"`)

	// The market refuses an unregistered buyer, and the client sends them to Settings.
	res, _ := send(t, c, "POST", ts.URL+"/market/buy", url.Values{"product_id": {"prod-ui2"}}, "HX-Request", "true")
	if res.Header.Get("HX-Redirect") != "./settings" {
		t.Fatalf("unregistered buy HX-Redirect = %q, want ./settings", res.Header.Get("HX-Redirect"))
	}

	registerUI2(t, ts, c)
	res, body = send(t, c, "POST", ts.URL+"/market/buy", url.Values{"product_id": {"prod-ui2"}}, "HX-Request", "true")
	if res.Header.Get("HX-Retarget") != "#market-payment" {
		t.Fatalf("buy HX-Retarget = %q, body %s", res.Header.Get("HX-Retarget"), body)
	}
	mustContain(t, body, "Pay the seller to complete your purchase", "19wiPpvWPxHEmcANKzwUjiRgABQk8wZzCp",
		`data-qr-link="bitcoin:19wiPpvWPxHEmcANKzwUjiRgABQk8wZzCp?amount=`, `id="market-error" hx-swap-oob="true"`, `id="products"`)
	// The frozen product left the grid.
	if strings.Contains(body, `value="prod-ui2"`) {
		t.Error("the bought product is still offered")
	}

	// A second buy of the same product fails and shows why.
	res, body = send(t, c, "POST", ts.URL+"/market/buy", url.Values{"product_id": {"prod-ui2"}}, "HX-Request", "true")
	if res.Header.Get("HX-Retarget") != "" {
		t.Fatal("a failed buy retargets")
	}
	mustContain(t, body, `id="market-error"`, "alert-danger")
}

// TestUI2MarketListing creates a sell order, with and without htmx.
func TestUI2MarketListing(t *testing.T) {
	srv, _ := newReadyMarket(t)
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)
	registerUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/market?new=1", nil)
	mustContain(t, body, `id="create-form"`, `<option value="SKY">SKY</option>`, `<option value="BTC" selected>BTC</option>`)
	if strings.Contains(body, `id="create-panel" hidden`) {
		t.Error("?new=1 leaves the form closed")
	}

	res, body := send(t, c, "POST", ts.URL+"/market/listing", url.Values{
		"sell_coin": {"SKY"}, "amount": {"0"}, "price": {"1"}, "payment_currency": {"BTC"},
	}, "HX-Request", "true")
	if res.Header.Get("HX-Trigger") != "" {
		t.Fatal("a refused listing triggers listing-created")
	}
	mustContain(t, body, "alert-danger")

	res, body = send(t, c, "POST", ts.URL+"/market/listing", url.Values{
		"sell_coin": {"SKY"}, "amount": {"10"}, "price": {"0.001"}, "payment_currency": {"BTC"},
	}, "HX-Request", "true")
	if res.Header.Get("HX-Trigger") != "listing-created" || res.Header.Get("HX-Retarget") != "#market-deposit" {
		t.Fatalf("listing headers = %v, body %s", res.Header, body)
	}
	mustContain(t, body, "Fund your sell order", "sky-market-wallet", `data-copy="sky-market-wallet"`, "data-local-time")

	// Without htmx the whole page comes back with the instructions.
	res, body = send(t, c, "POST", ts.URL+"/market/listing", url.Values{
		"sell_coin": {"SKY"}, "amount": {"10"}, "price": {"0.001"}, "payment_currency": {"BTC"},
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("plain listing = %d", res.StatusCode)
	}
	mustContain(t, body, "<!doctype html>", `class="tab-link active" href="market"`)
}

// TestUI2Orders shows a buy in flight and cancels it.
func TestUI2Orders(t *testing.T) {
	srv, database := newReadyMarket(t)
	seedProduct(t, srv, database, "prod-ui2")
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)
	registerUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/orders", nil)
	mustContain(t, body, "My Orders", "have any orders yet.", `hx-get="orders/table"`)

	send(t, c, "POST", ts.URL+"/market/buy", url.Values{"product_id": {"prod-ui2"}}, "HX-Request", "true")
	_, body = send(t, c, "GET", ts.URL+"/orders/table", nil)
	mustContain(t, body, `id="orders"`, "badge buy", "Pending Payment", "Awaiting payment…", "19wiPpvWPxHEmcANKzwUjiRgABQk8wZzCp",
		`hx-post="orders/cancel"`, "hx-confirm=", "data-local-time=\"datetime\"")
	i := strings.Index(body, `name="id" value="`)
	if i < 0 {
		t.Fatal("no cancel form")
	}
	id := body[i+len(`name="id" value="`):]
	id = id[:strings.Index(id, `"`)]

	_, body = send(t, c, "POST", ts.URL+"/orders/cancel", url.Values{"id": {id}}, "HX-Request", "true")
	mustContain(t, body, `id="orders"`, `id="tab-notice" hx-swap-oob="true"`, "alert-info", "Cancelled") //nolint:misspell // the label React shows
	if strings.Contains(body, `hx-post="orders/cancel"`) {
		t.Error("a canceled buy can still be canceled")
	}

	_, body = send(t, c, "POST", ts.URL+"/orders/cancel", url.Values{"id": {"nope"}}, "HX-Request", "true")
	mustContain(t, body, "alert-danger")
}

// TestUI2Listings shows a pending listing and cancels it.
func TestUI2Listings(t *testing.T) {
	srv, _ := newReadyMarket(t)
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)
	registerUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/listings", nil)
	mustContain(t, body, "My Listings", "You have no active sell listings.")

	send(t, c, "POST", ts.URL+"/market/listing", url.Values{
		"sell_coin": {"SKY"}, "amount": {"10"}, "price": {"0.001"}, "payment_currency": {"BTC"},
	}, "HX-Request", "true")
	_, body = send(t, c, "GET", ts.URL+"/listings/table", nil)
	mustContain(t, body, `id="listings"`, "sky-market-wallet", "Pending deposit", `class="badge bg-info">Pending deposit`,
		"fee", "Cancel this listing?", "For each pending listing")
	i := strings.Index(body, `name="id" value="`)
	if i < 0 {
		t.Fatal("no cancel form")
	}
	id := body[i+len(`name="id" value="`):]
	id = id[:strings.Index(id, `"`)]

	// Without htmx the whole tab comes back with the notice.
	res, body := send(t, c, "POST", ts.URL+"/listings/cancel", url.Values{"id": {id}})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("plain cancel = %d", res.StatusCode)
	}
	mustContain(t, body, "<!doctype html>", `id="tab-notice"`, "alert-info")
}

// TestUI2History merges the browser's trades with the market's finished
// orders.
func TestUI2History(t *testing.T) {
	srv, database := newReadyMarket(t)
	seedProduct(t, srv, database, "prod-ui2")
	ts, c := newUI2Server(t, srv, "")
	connectUI2(t, ts, c)
	registerUI2(t, ts, c)

	_, body := send(t, c, "GET", ts.URL+"/history", nil)
	mustContain(t, body, `hx-post="history/table"`, `data-local="exchange:history"`)

	_, body = send(t, c, "POST", ts.URL+"/history/table", url.Values{"saved": {"[]"}})
	mustContain(t, body, "No completed trades yet.", `id="history-data"`, `>[]`)

	send(t, c, "POST", ts.URL+"/market/buy", url.Values{"product_id": {"prod-ui2"}}, "HX-Request", "true")
	_, body = send(t, c, "GET", ts.URL+"/orders/table", nil)
	i := strings.Index(body, `name="id" value="`)
	id := body[i+len(`name="id" value="`):]
	id = id[:strings.Index(id, `"`)]
	send(t, c, "POST", ts.URL+"/orders/cancel", url.Values{"id": {id}}, "HX-Request", "true")

	saved := `[{"id":"purged","type":"sell","sell_coin":"SKY","amount":1,"price":2,"payment_currency":"BTC","status":"completed","created_at":"2020-01-01T00:00:00Z","market_name":"Gone market</script>"}]`
	_, body = send(t, c, "POST", ts.URL+"/history/table", url.Values{"saved": {saved}})
	mustContain(t, body, "Clear history", `data-remove-trade="purged"`, `data-remove-trade="`+id+`"`, "badge enabled", "Gone market", `"market_pk":"`+testMarketRef+`"`)
	if strings.Count(body, "</script>") != 1 {
		t.Error("a saved trade can close the JSON script")
	}
}
