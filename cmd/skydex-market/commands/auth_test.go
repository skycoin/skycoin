// Package commands cmd/apps/skydex-market/commands/auth_test.go
package commands

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestGate(t *testing.T) (*authGate, *[]string) {
	t.Helper()

	published := &[]string{}
	g, err := newAuthGate(func(otp string) { *published = append(*published, otp) })
	if err != nil {
		t.Fatalf("newAuthGate: %v", err)
	}

	return g, published
}

func TestAuthGate_PublishesOTPOnStart(t *testing.T) {
	g, published := newTestGate(t)

	if len(*published) != 1 {
		t.Fatalf("expected 1 published OTP, got %d", len(*published))
	}
	if (*published)[0] != g.currentOTP() {
		t.Fatal("published OTP does not match current")
	}
	if len(g.currentOTP()) != otpLen {
		t.Fatalf("OTP length = %d, want %d", len(g.currentOTP()), otpLen)
	}
	// The alphabet is all digits, so 0 and 1 are legitimate here — there are no
	// letters left for them to be confused with, which is part of why the
	// alphabet is numeric. Assert membership instead of banning glyphs.
	for _, c := range g.currentOTP() {
		if !strings.ContainsRune(otpAlphabet, c) {
			t.Fatalf("OTP %q contains %q, outside otpAlphabet %q", g.currentOTP(), c, otpAlphabet)
		}
	}
}

func TestAuthGate_OTPIsSingleUse(t *testing.T) {
	g, published := newTestGate(t)
	otp := g.currentOTP()

	token, ok := g.login(otp)
	if !ok || token == "" {
		t.Fatal("first login with a valid OTP should succeed")
	}

	// The whole point: replaying the same code must fail.
	if _, ok := g.login(otp); ok {
		t.Fatal("OTP was accepted twice")
	}

	if g.currentOTP() == otp {
		t.Fatal("OTP was not rotated after use")
	}
	// Three publishes: the initial code, the replacement minted by the successful
	// login, and a third from the rejected replay above — under
	// otpFailureRotate == 1 a wrong guess burns the live code as well.
	if len(*published) != 3 {
		t.Fatalf("expected 3 publishes (initial, post-login, post-failed-replay), got %d", len(*published))
	}
	if (*published)[len(*published)-1] != g.currentOTP() {
		t.Fatal("published replacement does not match current OTP")
	}
}

// wrongCode returns a well-formed code guaranteed not to equal cur, by stepping
// its first character one place through otpAlphabet. Negative tests need this:
// with a five-digit numeric alphabet a hardcoded guess would collide with the
// live code once in 100_000 runs, which is exactly the kind of rare CI flake
// that gets dismissed as noise.
func wrongCode(cur string) string {
	b := []byte(cur)
	i := strings.IndexByte(otpAlphabet, b[0])
	b[0] = otpAlphabet[(i+1)%len(otpAlphabet)]

	return string(b)
}

func TestAuthGate_LoginRejectsWrongOTP(t *testing.T) {
	g, _ := newTestGate(t)

	// Every wrong guess burns the live code (otpFailureRotate == 1), so re-read it
	// each iteration instead of assuming it survived the previous attempt. The
	// literals below cannot collide with a real code — none of them is five
	// digits — so only the last case needs wrongCode.
	bad := []string{"", "WRONG", strings.Repeat("A", otpLen), g.currentOTP() + "X"}
	for _, b := range append(bad, wrongCode(g.currentOTP())) {
		before := g.currentOTP()
		if _, ok := g.login(b); ok {
			t.Fatalf("login accepted invalid OTP %q", b)
		}
		// The point of the change: a single wrong guess must move the target, so
		// an attacker can never enumerate the space.
		if g.currentOTP() == before {
			t.Fatalf("wrong OTP %q did not burn the code", b)
		}
	}

	// The operator is not locked out — the code current at this moment still works.
	if _, ok := g.login(g.currentOTP()); !ok {
		t.Fatal("valid OTP stopped working after failed attempts")
	}
}

func TestAuthGate_LoginNormalizesInput(t *testing.T) {
	g, _ := newTestGate(t)
	otp := g.currentOTP()

	// Operators retype these by hand; browsers add whitespace.
	if _, ok := g.login("  " + strings.ToLower(otp) + "\n"); !ok {
		t.Fatal("login should tolerate surrounding whitespace and lowercase")
	}
}

func TestAuthGate_SessionTokens(t *testing.T) {
	g, _ := newTestGate(t)

	token, ok := g.login(g.currentOTP())
	if !ok {
		t.Fatal("login failed")
	}

	if !g.valid(token) {
		t.Fatal("freshly issued token should be valid")
	}
	if g.valid("") || g.valid("not-a-real-token") {
		t.Fatal("bogus token accepted")
	}

	// Two logins must not collide.
	token2, ok := g.login(g.currentOTP())
	if !ok {
		t.Fatal("second login failed")
	}
	if token == token2 {
		t.Fatal("two logins produced the same token")
	}
	if !g.valid(token) || !g.valid(token2) {
		t.Fatal("both tokens should remain valid")
	}
}

func TestAuthGate_RequireAuthMiddleware(t *testing.T) {
	g, _ := newTestGate(t)

	protected := g.requireAuth(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`)) //nolint
	}))

	// No credential at all.
	rec := httptest.NewRecorder()
	protected.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/config", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated request: got %d, want 401", rec.Code)
	}

	token, _ := g.login(g.currentOTP())

	// Correct token.
	rec = httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	protected.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("authenticated request: got %d, want 200", rec.Code)
	}

	// HTTP Basic must NOT work — browsers resend it ambiently, which is exactly
	// the CSRF exposure this design avoids.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/config", nil)
	req.SetBasicAuth("operator", token)
	protected.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("basic auth: got %d, want 401", rec.Code)
	}
}

func TestAuthGate_LoginHandler(t *testing.T) {
	g, _ := newTestGate(t)
	h := g.loginHandler()

	post := func(body string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(body))
		req.RemoteAddr = "10.0.0.9:1234"
		h.ServeHTTP(rec, req)
		return rec
	}

	rec := post(`{"otp":"` + g.currentOTP() + `"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("valid login: got %d, want 200 (%s)", rec.Code, rec.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if !g.valid(resp.Token) {
		t.Fatal("token returned by loginHandler is not valid")
	}

	// GET must not be a way in.
	rec = httptest.NewRecorder()
	getReq := httptest.NewRequest(http.MethodGet, "/api/login", nil)
	getReq.RemoteAddr = "10.0.0.9:1234"
	h.ServeHTTP(rec, getReq)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("GET /api/login: got %d, want 405", rec.Code)
	}
}

func TestAuthGate_LoginRateLimited(t *testing.T) {
	g, _ := newTestGate(t)
	h := g.loginHandler()

	attempt := func(ip string) int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBufferString(`{"otp":"BADCODE1"}`))
		req.RemoteAddr = ip + ":5555"
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	// Burst is consumed by wrong guesses, then the limiter takes over.
	var throttled bool
	for range loginBurst + 3 {
		if attempt("203.0.113.7") == http.StatusTooManyRequests {
			throttled = true
			break
		}
	}
	if !throttled {
		t.Fatal("brute-force attempts were never rate limited")
	}

	// The limiter is per-IP, so one attacker must not lock out the operator.
	if code := attempt("198.51.100.4"); code == http.StatusTooManyRequests {
		t.Fatal("a different IP was throttled by another client's attempts")
	}
}

// The wiring in serveUI relies on net/http's mux preferring the exact
// "/api/login" pattern over the "/api/" prefix. If that ever stopped holding,
// login would be shadowed by the auth wrapper and the panel would be
// permanently unreachable — so pin it.
func TestServeUIRouting_LoginIsReachableButRestIsGated(t *testing.T) {
	g, _ := newTestGate(t)

	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/api/config", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"secret":"present"}`)) //nolint
	})

	mux := http.NewServeMux()
	mux.Handle("/api/login", g.loginHandler())
	mux.Handle("/api/", g.requireAuth(apiMux))

	// Login is reachable with no token.
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login",
		bytes.NewBufferString(`{"otp":"`+g.currentOTP()+`"}`))
	req.RemoteAddr = "10.0.0.1:1111"
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login through mux: got %d, want 200 (%s)", rec.Code, rec.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Everything else is gated.
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/config", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ungated /api/config: got %d, want 401", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "present") {
		t.Fatal("unauthenticated request leaked handler output")
	}

	// ...and reachable once authenticated.
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/config", nil)
	req.Header.Set("Authorization", "Bearer "+resp.Token)
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("authenticated /api/config: got %d, want 200", rec.Code)
	}
}

// The code is only ~25 bits, so these two controls — not entropy — are what
// make it safe to brute-force. If either regresses, a short OTP becomes
// guessable, so pin both.
func TestAuthGate_RotatesOTPAfterRepeatedFailures(t *testing.T) {
	g, published := newTestGate(t)
	original := g.currentOTP()

	// Drive failures directly, bypassing the rate limiter (which is exercised
	// separately) so this test isolates the rotation behavior.
	for range otpFailureRotate {
		if _, ok := g.login("AAAAA"); ok {
			t.Fatal("a junk code was accepted")
		}
	}

	if g.currentOTP() == original {
		t.Fatalf("OTP did not rotate after %d failed attempts", otpFailureRotate)
	}
	// The operator must be able to find the replacement.
	if (*published)[len(*published)-1] != g.currentOTP() {
		t.Fatal("the rotated OTP was not published to the visor")
	}
	// The old code must be dead even though it was never successfully used.
	if _, ok := g.login(original); ok {
		t.Fatal("the pre-rotation OTP still works")
	}
}

func TestAuthGate_FailureCountResetsOnRotation(t *testing.T) {
	g, _ := newTestGate(t)

	// A failure rotates the code and must leave the counter clean, so the
	// operator's freshly published replacement is not rotated out from under them
	// by carried-over state. "AAAAA" is guaranteed wrong against a digit alphabet.
	g.login("AAAAA") //nolint:errcheck
	if _, ok := g.login(g.currentOTP()); !ok {
		t.Fatal("valid login failed")
	}

	g.mu.Lock()
	fails := g.fails
	g.mu.Unlock()
	if fails != 0 {
		t.Fatalf("failure count = %d after a successful login, want 0", fails)
	}
}

func TestAuthGate_GlobalRateLimitCapsDistributedAttack(t *testing.T) {
	g, _ := newTestGate(t)

	// Every attempt from a distinct source address, so the per-IP limiter never
	// fires. Without a global cap this loop would never be throttled — which is
	// exactly how a botnet defeats per-IP limiting.
	blocked := 0
	for i := range 60 {
		ip := "203.0.113." + string(rune('0'+i%10)) + "." + string(rune('0'+i/10))
		if !g.allowLogin(ip) {
			blocked++
		}
	}

	if blocked == 0 {
		t.Fatal("distributed attempts from unique IPs were never globally throttled")
	}
}

func TestOTPLengthIsReadableButGuarded(t *testing.T) {
	g, _ := newTestGate(t)

	if len(g.currentOTP()) != otpLen {
		t.Fatalf("OTP length = %d, want %d", len(g.currentOTP()), otpLen)
	}
	for _, c := range g.currentOTP() {
		if !strings.ContainsRune(otpAlphabet, c) {
			t.Fatalf("OTP contains %q, which is outside otpAlphabet %q", c, otpAlphabet)
		}
	}

	// The code space is deliberately small (10^5) for typing ergonomics, which
	// only holds up because two controls compensate. Pin both, so shrinking the
	// space further or loosening a control trips here rather than in production.
	space := 1
	for range otpLen {
		space *= len(otpAlphabet)
	}

	// Enumeration defense: below ~1e6 the space is walkable, so every wrong guess
	// must burn the code and force the attacker to guess independently.
	if space < 1_000_000 && otpFailureRotate != 1 {
		t.Fatalf("code space %d needs otpFailureRotate == 1, got %d", space, otpFailureRotate)
	}

	// Rate defense: the global limiter is what actually bounds an attacker, since
	// per-IP budget is multiplied away by anyone with several addresses.
	//
	// As configured — 10^5 codes, one attempt per 6m — sustained global guessing
	// needs ~417 days to cover the space, so the floor is a real bar rather than a
	// rubber stamp: a full year. Shrinking otpLen or speeding up globalLoginRefill
	// past that point trips here instead of in production. Each extra digit is a
	// 10x multiplier if more headroom is ever needed.
	const floorDays = 365.0

	perDay := float64(24*time.Hour) / float64(globalLoginRefill)
	if days := float64(space) / perDay; days < floorDays {
		t.Fatalf("sustained global guessing exhausts the %d-code space in %.1f days "+
			"(floor %.0f); lengthen otpLen or slow globalLoginRefill", space, days, floorDays)
	}
}

func TestBearerToken(t *testing.T) {
	cases := map[string]string{
		"Bearer abc123": "abc123",
		"bearer abc123": "abc123", // scheme is case-insensitive per RFC 7235
		"Bearer   xyz ": "xyz",
		"Basic abc123":  "",
		"abc123":        "",
		"":              "",
		"Bearer ":       "",
	}

	for header, want := range cases {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		if got := bearerToken(req); got != want {
			t.Fatalf("bearerToken(%q) = %q, want %q", header, got, want)
		}
	}
}

func TestClientIP_IgnoresForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.0.2.10:4321"
	// Attacker-controlled: honoring it would let one client evade the limiter.
	req.Header.Set("X-Forwarded-For", "1.2.3.4")

	if got := clientIP(req); got != "192.0.2.10" {
		t.Fatalf("clientIP = %q, want 192.0.2.10", got)
	}
}
