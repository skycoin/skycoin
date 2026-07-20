package commands

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/market"
	"github.com/skycoin/skycoin/src/skydex/server"
)

// testMarketRef is the market reference the UI dials in these tests. fakeDialer
// ignores it (it always serves the in-process market), so any non-empty value
// works; the tests just assert it round-trips through the connect API.
const testMarketRef = "03aaaabbbbccccddddeeeeffff00001111222233334444555566667777888899aa"

// fakeDialer is the test MarketDialer: every Dial returns a *market.Conn whose
// other end is served by a real in-process market Server over net.Pipe.
type fakeDialer struct{ srv *server.Server }

func (f fakeDialer) Dial(_ string) (*market.Conn, error) {
	clientEnd, serverEnd := net.Pipe()
	go f.srv.Serve(serverEnd, testMarketRef)
	return market.NewConn(clientEnd), nil
}

func newTestMarket(t *testing.T) *server.Server {
	t.Helper()
	database, err := db.New(filepath.Join(t.TempDir(), "m.db"), "")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() }) //nolint
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := database.InitDefaultConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	return server.New(database, nil, 0)
}

func newTestServer(t *testing.T, defaultPK string) *httptest.Server {
	t.Helper()
	sess := newSession(fakeDialer{srv: newTestMarket(t)}, defaultPK)
	mux := http.NewServeMux()
	registerAPI(mux, sess)
	ts := httptest.NewServer(mux)
	t.Cleanup(func() { ts.Close(); sess.close() })
	return ts
}

// TestConfigExposesDefaultPK verifies the --market-pk value reaches the UI.
func TestConfigExposesDefaultPK(t *testing.T) {
	ts := newTestServer(t, "the-default-pk")
	var out struct {
		MarketPK string `json:"market_pk"`
	}
	getJSON(t, ts.URL+"/api/config", &out)
	if out.MarketPK != "the-default-pk" {
		t.Fatalf("market_pk = %q, want the default", out.MarketPK)
	}
}

// TestConnectFlow drives status -> connect -> status -> disconnect through the
// HTTP API, exercising the real dmsg-style handshake against a market.
func TestConnectFlow(t *testing.T) {
	ts := newTestServer(t, "")

	// Initially disconnected.
	var st struct {
		Connected bool `json:"connected"`
	}
	getJSON(t, ts.URL+"/api/status", &st)
	if st.Connected {
		t.Fatal("expected disconnected initially")
	}

	// Connect with a market reference.
	var connResp struct {
		Connected bool     `json:"connected"`
		MarketPK  string   `json:"market_pk"`
		Curr      []string `json:"currencies"`
	}
	body, _ := json.Marshal(map[string]string{"market_pk": testMarketRef}) //nolint
	res, err := http.Post(ts.URL+"/api/connect", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST connect: %v", err)
	}
	decode(t, res, &connResp)
	if res.StatusCode != http.StatusOK || !connResp.Connected || connResp.MarketPK != testMarketRef {
		t.Fatalf("connect resp = %+v (status %d), want connected", connResp, res.StatusCode)
	}

	// Status now reflects the connection.
	getJSON(t, ts.URL+"/api/status", &st)
	if !st.Connected {
		t.Fatal("expected connected after /api/connect")
	}

	// Disconnect.
	res, err = http.Post(ts.URL+"/api/disconnect", "application/json", nil)
	if err != nil {
		t.Fatalf("POST disconnect: %v", err)
	}
	res.Body.Close() //nolint
	getJSON(t, ts.URL+"/api/status", &st)
	if st.Connected {
		t.Fatal("expected disconnected after /api/disconnect")
	}
}

// TestProxyRoutes verifies the UI-facing proxy routes forward to the connected
// market and return its data (register -> products -> currencies).
func TestProxyRoutes(t *testing.T) {
	ts := newTestServer(t, "")

	// Proxying before connecting must be rejected with 409.
	res, err := http.Get(ts.URL + "/api/products") //nolint:noctx
	if err != nil {
		t.Fatalf("GET products: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusConflict {
		t.Fatalf("products before connect = %d, want 409", res.StatusCode)
	}

	// Connect.
	body, _ := json.Marshal(map[string]string{"market_pk": testMarketRef}) //nolint
	res, err = http.Post(ts.URL+"/api/connect", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST connect: %v", err)
	}
	res.Body.Close() //nolint

	// Register wallets.
	// Real mainnet addresses — registration validates address format.
	reg, _ := json.Marshal(map[string]string{ //nolint
		"wallet_sky": "xPHfK6xn5LZvWAPB5EV6A62WyzQLywci1J",
		"wallet_btc": "1CLdkSYgwRA4Hm7EUPqaqpAtGG9zgT7dFm",
	})
	res, err = http.Post(ts.URL+"/api/register", "application/json", bytes.NewReader(reg))
	if err != nil {
		t.Fatalf("POST register: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusOK {
		t.Fatalf("register = %d, want 200", res.StatusCode)
	}

	// Products returns an empty list (nothing listed yet).
	var prods struct {
		Products []any `json:"products"`
	}
	getJSON(t, ts.URL+"/api/products", &prods)
	if prods.Products == nil {
		t.Fatalf("products payload missing 'products' field")
	}

	// Currencies returns a (possibly empty) list.
	var cur struct {
		Currencies []string `json:"currencies"`
	}
	getJSON(t, ts.URL+"/api/currencies", &cur)

	// The caller's own listings proxy returns a (possibly empty) list.
	var mine struct {
		Listings []any `json:"listings"`
	}
	getJSON(t, ts.URL+"/api/listings/mine", &mine)
	if mine.Listings == nil {
		t.Fatalf("listings payload missing 'listings' field")
	}
}

// TestConnectRejectsEmptyPK verifies connect fails cleanly with no key.
func TestConnectRejectsEmptyPK(t *testing.T) {
	ts := newTestServer(t, "")
	body, _ := json.Marshal(map[string]string{"market_pk": ""}) //nolint
	res, err := http.Post(ts.URL+"/api/connect", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST connect: %v", err)
	}
	defer res.Body.Close() //nolint
	if res.StatusCode == http.StatusOK {
		t.Fatal("expected error status for empty market_pk")
	}
}

func getJSON(t *testing.T, url string, out any) {
	t.Helper()
	res, err := http.Get(url) //nolint:noctx,gosec
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	decode(t, res, out)
}

func decode(t *testing.T, res *http.Response, out any) {
	t.Helper()
	defer res.Body.Close() //nolint
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		t.Fatalf("decode: %v", err)
	}
}
