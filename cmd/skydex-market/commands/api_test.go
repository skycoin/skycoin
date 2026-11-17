package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/skydex/db"
)

func newTestOperatorServer(t *testing.T) (*httptest.Server, *db.Database) {
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

	mux := http.NewServeMux()
	// A stdHost with empty PubKey stands in for the visor-backed host here.
	registerOperatorAPI(mux, database, stdHost{log: logrus.New()})
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, database
}

// TestOperatorConfigReadWrite verifies config GET/POST, including the
// editable-key whitelist.
func TestOperatorConfigReadWrite(t *testing.T) {
	ts, _ := newTestOperatorServer(t)

	// GET returns the seeded defaults. Response is {config, secrets_set} —
	// secret values are blanked server side (see redactConfig).
	readConfig := func() map[string]string {
		t.Helper()
		var resp struct {
			Config     map[string]string `json:"config"`
			SecretsSet map[string]bool   `json:"secrets_set"`
		}
		getJSONInto(t, ts.URL+"/api/config", &resp)
		return resp.Config
	}

	cfg := readConfig()
	if _, ok := cfg["sky_fullnode_url"]; !ok {
		t.Fatalf("config missing sky_fullnode_url: %v", cfg)
	}

	// POST a valid update: enable BTC via its explorer provider.
	body, _ := json.Marshal(map[string]string{"explorer_btc_provider": "esplora", "wallet_sky": "sky-w"}) //nolint
	res, err := http.Post(ts.URL+"/api/config", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST config: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusOK {
		t.Fatalf("POST config = %d, want 200", res.StatusCode)
	}
	cfg = readConfig()
	if cfg["explorer_btc_provider"] != "esplora" || cfg["wallet_sky"] != "sky-w" {
		t.Fatalf("config not persisted: %v", cfg)
	}

	// POST an unknown key is rejected.
	bad, _ := json.Marshal(map[string]string{"evil_key": "x"}) //nolint
	res, err = http.Post(ts.URL+"/api/config", "application/json", bytes.NewReader(bad))
	if err != nil {
		t.Fatalf("POST bad config: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("unknown key = %d, want 400", res.StatusCode)
	}

	// POST an explorer provider that doesn't support the coin is rejected.
	badProv, _ := json.Marshal(map[string]string{"explorer_btc_provider": "bogus"}) //nolint
	res, err = http.Post(ts.URL+"/api/config", "application/json", bytes.NewReader(badProv))
	if err != nil {
		t.Fatalf("POST bad provider: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad provider = %d, want 400", res.StatusCode)
	}
}

// TestOperatorMonitoring verifies the monitoring endpoints return their arrays.
func TestOperatorMonitoring(t *testing.T) {
	ts, database := newTestOperatorServer(t)

	// Seed one product so /api/products is non-empty.
	if err := database.CreateUser(&db.User{PubKey: "03seller", WalletSKY: "sky"}); err != nil {
		t.Fatal(err)
	}
	if err := database.CreateProduct(&db.Product{
		ID: "p1", SellerPubKey: "03seller", Amount: 5, Price: 1, PaymentCurrency: "BTC", Status: "active",
	}); err != nil {
		t.Fatal(err)
	}

	var prods struct {
		Products []map[string]any `json:"products"`
	}
	getJSONInto(t, ts.URL+"/api/products", &prods)
	if len(prods.Products) != 1 || prods.Products[0]["id"] != "p1" {
		t.Fatalf("products = %+v, want p1", prods.Products)
	}

	for _, path := range []string{"/api/orders", "/api/bans"} {
		res, err := http.Get(ts.URL + path) //nolint:noctx
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		res.Body.Close() //nolint
		if res.StatusCode != http.StatusOK {
			t.Fatalf("GET %s = %d, want 200", path, res.StatusCode)
		}
	}
}

// TestServeUIServesHTMLAndAPI runs the real serveUI (embedded HTML + operator
// API) and checks both surfaces respond.
func TestServeUIServesHTMLAndAPI(t *testing.T) {
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	const addr = "127.0.0.1:18050"
	go serveUI(ctx, stdHost{log: logrus.New()}, database, addr)

	base := "http://" + addr
	var body []byte
	for i := 0; i < 40; i++ {
		res, err := http.Get(base + "/") //nolint:noctx
		if err == nil {
			body, _ = io.ReadAll(res.Body) //nolint
			res.Body.Close()               //nolint
			break
		}
		time.Sleep(25 * time.Millisecond)
	}
	if !bytes.Contains(body, []byte("Operator")) {
		t.Fatalf("index.html not served (got %d bytes)", len(body))
	}

	// The API is mounted alongside the HTML, but behind the OTP gate: the
	// login page must be reachable while everything else is refused.
	res, err := http.Get(base + "/api/config") //nolint:noctx
	if err != nil {
		t.Fatalf("GET /api/config: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("/api/config without a token = %d, want 401", res.StatusCode)
	}

	// A wrong OTP is refused, so the gate isn't merely decorative.
	res, err = http.Post(base+"/api/login", "application/json", //nolint:noctx
		bytes.NewReader([]byte(`{"otp":"NOPENOPE"}`)))
	if err != nil {
		t.Fatalf("POST /api/login: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("login with a wrong OTP = %d, want 401", res.StatusCode)
	}
}

func getJSONInto(t *testing.T, url string, out any) {
	t.Helper()
	res, err := http.Get(url) //nolint:noctx,gosec
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer res.Body.Close() //nolint
	if err := json.NewDecoder(res.Body).Decode(out); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
}
