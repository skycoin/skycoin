// Package commands cmd/apps/skydex-market/commands/otp_flow_test.go
package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/skydex/db"
)

// TestOTPFlow_EndToEnd walks the whole operator path the way the browser does:
// the market mints an OTP and publishes it (the hypervisor app list is where the
// operator reads it), the OTP buys a session token, the token unlocks the API,
// and the OTP is dead afterwards.
func TestOTPFlow_EndToEnd(t *testing.T) {
	const seed = "escrow hot wallet seed words that must never be served"

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
	if err := database.SetConfig("sky_wallet_seed", seed); err != nil {
		t.Fatalf("set seed: %v", err)
	}

	// Capture what the market publishes to the visor — this stands in for the
	// hypervisor app list the operator reads the code from.
	var publishedOTP string
	gate, err := newAuthGate(func(otp string) { publishedOTP = otp })
	if err != nil {
		t.Fatalf("newAuthGate: %v", err)
	}

	apiMux := http.NewServeMux()
	registerOperatorAPI(apiMux, database, stdHost{log: logrus.New()})

	mux := http.NewServeMux()
	mux.Handle("/api/login", gate.loginHandler())
	mux.Handle("/api/", gate.requireAuth(apiMux))

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	get := func(path, token string) (int, string) {
		t.Helper()
		req, err := http.NewRequest(http.MethodGet, ts.URL+path, nil) //nolint:noctx
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("GET %s: %v", path, err)
		}
		defer res.Body.Close() //nolint
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}

	login := func(otp string) (int, string) {
		t.Helper()
		body, _ := json.Marshal(map[string]string{"otp": otp})                                //nolint
		res, err := http.Post(ts.URL+"/api/login", "application/json", bytes.NewReader(body)) //nolint:noctx
		if err != nil {
			t.Fatalf("POST /api/login: %v", err)
		}
		defer res.Body.Close() //nolint
		var out struct {
			Token string `json:"token"`
		}
		_ = json.NewDecoder(res.Body).Decode(&out) //nolint
		return res.StatusCode, out.Token
	}

	// 1. The market published a code at startup.
	if publishedOTP == "" {
		t.Fatal("market did not publish an OTP to the visor")
	}

	// 2. Nothing is reachable without a token.
	if code, _ := get("/api/config", ""); code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated /api/config = %d, want 401", code)
	}

	// Snapshot before logging in: a successful login rotates and republishes,
	// which overwrites publishedOTP via the callback above.
	firstOTP := publishedOTP

	// 3. The published code buys a session token.
	code, token := login(firstOTP)
	if code != http.StatusOK || token == "" {
		t.Fatalf("login with published OTP = %d (token %q), want 200 + token", code, token)
	}

	// 4. That token unlocks the API...
	status, body := get("/api/config", token)
	if status != http.StatusOK {
		t.Fatalf("authenticated /api/config = %d, want 200", status)
	}

	// 5. ...but the escrow seed is still not in the response, auth or no auth.
	if strings.Contains(body, seed) {
		t.Fatal("escrow seed served to an authenticated operator; it must be write-only")
	}
	var cfgResp struct {
		Config     map[string]string `json:"config"`
		SecretsSet map[string]bool   `json:"secrets_set"`
	}
	if err := json.Unmarshal([]byte(body), &cfgResp); err != nil {
		t.Fatalf("decode config: %v", err)
	}
	if !cfgResp.SecretsSet["sky_wallet_seed"] {
		t.Fatal("secrets_set should report the seed as configured")
	}

	// 6. The consumed code is dead, and a fresh one was published in its place.
	if code, _ := login(firstOTP); code == http.StatusOK {
		t.Fatal("a consumed OTP was accepted a second time")
	}
	if firstOTP == gate.currentOTP() {
		t.Fatal("OTP was not rotated after a successful login")
	}
	if publishedOTP != gate.currentOTP() {
		t.Fatal("the replacement OTP was not published to the visor")
	}

	// 7. Saving config with a blank seed must not wipe the stored one — the UI
	// can never read it back, so blank has to mean "keep".
	upd, _ := json.Marshal(map[string]string{"sky_wallet_seed": "", "market_name": "Renamed"}) //nolint
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"/api/config", bytes.NewReader(upd))     //nolint:noctx,errcheck
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST /api/config: %v", err)
	}
	res.Body.Close() //nolint
	if res.StatusCode != http.StatusOK {
		t.Fatalf("POST /api/config = %d, want 200", res.StatusCode)
	}

	stored, err := database.GetConfig("sky_wallet_seed")
	if err != nil {
		t.Fatalf("read back seed: %v", err)
	}
	if stored != seed {
		t.Fatalf("blank submit wiped the escrow seed (got %q)", stored)
	}
}
