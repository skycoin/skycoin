// Package commands cmd/apps/skydex-market/commands/auth.go
package commands

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"math/big"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// The operator panel is authenticated with a one-time code the market itself
// mints and publishes to the visor, where it surfaces on the hypervisor's app
// list (which is behind hypervisor auth). The operator reads it there and types
// it into the market's login. On success the code is consumed and replaced, and
// the browser receives a short-lived token it keeps in memory only.
//
// Why no cookie: without an ambient credential a hostile page cannot make the
// browser replay it, so CSRF is structurally impossible rather than something
// we defend against. The cost is that a reload means a fresh code — deliberate.
const (
	// otpLen is the number of characters in a generated code — short enough to
	// read off a dashboard and retype, which is the whole ergonomic point.
	//
	// WARNING — the code space here is SMALL. Five digits is 100_000 codes, about
	// 16.6 bits. Entropy does not make this safe on its own; the rate limiters
	// below are the only thing that does, so treat them as load-bearing and do
	// not loosen them. For reference, at the sustained global budget
	// (globalLoginRefill, one attempt per 6m ≈ 240/day) an attacker who keeps
	// guessing needs ~417 days to cover the whole space. If this gate ever
	// guards more value, lengthen the code before touching anything else: each
	// extra digit multiplies the attacker's work by 10, and 8 digits (~26.6 bits)
	// restores the margin the previous 5-character alphanumeric code had.
	otpLen = 5

	// otpAlphabet is digits only: the operator retypes the code by eye from a
	// dashboard, and digits remove glyph ambiguity outright (no I/O vs 1/0) while
	// letting phones offer a numeric keypad. The trade is steep — dropping from a
	// 32-symbol alphabet to 10 shrinks the space 335x at the same length — and is
	// accepted deliberately in favor of entry ergonomics. See otpLen above.
	otpAlphabet = "0123456789"

	// tokenAlphabet is used for session tokens, which are generated, stored, and
	// replayed by the browser and never read or typed by a person. It is kept
	// separate from otpAlphabet on purpose: the readability constraints that
	// shrink the OTP alphabet buy nothing here, and sharing one alphabet meant an
	// ergonomic change to the OTP silently weakened tokens too.
	tokenAlphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// tokenLen gives ~190 bits over tokenAlphabet — far past guessing, and the
	// reason session tokens need no rotation or rate limiting of their own.
	tokenLen = 32

	// sessionTTL bounds how long a logged-in page may keep calling the API
	// before it must be re-authenticated with a fresh code.
	sessionTTL = 12 * time.Hour

	// loginBurst / loginRefill throttle login attempts per client IP. Budget is
	// consumed by every attempt, not just failures — the check necessarily runs
	// before the outcome is known. The burst is sized so a human operator
	// (one attempt, occasionally a typo) never notices.
	loginBurst  = 5
	loginRefill = time.Minute

	// globalLoginBurst / globalLoginRefill cap login attempts across ALL source
	// addresses. The per-IP limiter alone is useless against an attacker with a
	// botnet or an IPv6 range, which is exactly the threat a ~16.6-bit code faces,
	// so this is THE control that bounds brute force — size it against the code
	// space, not against what feels responsive.
	//
	// One attempt per 6 minutes is 240/day, so covering the 100_000-code space
	// costs an attacker ~417 days. A real operator never feels it: a login spends
	// one token, and globalLoginBurst absorbs a run of typos before the sustained
	// rate applies at all.
	globalLoginBurst  = 10
	globalLoginRefill = 6 * time.Minute

	// otpFailureRotate replaces the code after this many failed attempts. It is 1:
	// a single wrong guess burns the code. That is the strongest form of this
	// control and matters most now that otpLen is only ~16.6 bits — it denies an
	// attacker the ability to ENUMERATE. Walking 00000, 00001, 00002 … is
	// guaranteed to succeed within the space and takes half of it on average;
	// rotating every failure makes each guess independent, so there is no
	// guarantee at any number of attempts and the expected work doubles.
	//
	// Be clear about the size of that win: it is a factor of ~2, not a rescue. It
	// does NOT make a five-digit code strong, because it does not change the
	// per-guess probability of 1/100_000. The rate limiters remain the control
	// that actually bounds an attacker.
	//
	// Cost, accepted deliberately: one typo sends the operator back to the
	// hypervisor app list for a fresh code, and an attacker who spends his budget
	// on wrong guesses keeps the code churning, which is a denial-of-service
	// against login. Both are bounded by the global limiter (one attempt per 6m).
	otpFailureRotate = 1
)

// authGate owns the market's current OTP and the set of live page sessions.
//
// The zero value is not usable; construct with newAuthGate.
type authGate struct {
	mu       sync.Mutex
	otp      string
	fails    int                  // consecutive failed attempts since the last rotation
	sessions map[string]time.Time // token -> expiry
	limiters map[string]*rate.Limiter
	global   *rate.Limiter
	// publish hands a freshly minted OTP to the visor. Injected so tests can
	// observe rotation without an app client.
	publish func(otp string)
}

// newAuthGate mints the first OTP and publishes it. publish may be nil.
func newAuthGate(publish func(otp string)) (*authGate, error) {
	if publish == nil {
		publish = func(string) {}
	}
	g := &authGate{
		sessions: make(map[string]time.Time),
		limiters: make(map[string]*rate.Limiter),
		global:   rate.NewLimiter(rate.Every(globalLoginRefill), globalLoginBurst),
		publish:  publish,
	}
	if err := g.rotate(); err != nil {
		return nil, err
	}
	return g, nil
}

// rotate replaces the current OTP and publishes the replacement. Callers must
// not hold g.mu.
func (g *authGate) rotate() error {
	otp, err := randomCode(otpAlphabet, otpLen)
	if err != nil {
		return err
	}

	g.mu.Lock()
	g.otp = otp
	g.fails = 0
	g.mu.Unlock()

	g.publish(otp)

	return nil
}

// login consumes the current OTP and returns a session token. The OTP is
// rotated on success, so a code works exactly once.
func (g *authGate) login(otp string) (string, bool) {
	// Normalize the way an operator retypes it: browsers happily contribute stray
	// whitespace. ToUpper is a no-op for the all-digit alphabet but is kept so a
	// future alphabet change cannot silently reintroduce case sensitivity.
	otp = strings.ToUpper(strings.TrimSpace(otp))

	g.mu.Lock()
	current := g.otp
	// Constant-time compare so a wrong code can't be narrowed by timing. The
	// length check is deliberately outside it — length is not a secret, and
	// subtle.ConstantTimeCompare returns 0 for mismatched lengths anyway.
	ok := len(otp) == len(current) &&
		subtle.ConstantTimeCompare([]byte(otp), []byte(current)) == 1
	if !ok {
		g.fails++
		exhausted := g.fails >= otpFailureRotate
		g.mu.Unlock()

		// Move the target before a search can converge. rotate() takes g.mu
		// itself, so this must happen after the unlock above.
		if exhausted {
			_ = g.rotate() //nolint:errcheck // keep the old code rather than lock the operator out
		}

		return "", false
	}
	g.mu.Unlock()

	token, err := randomCode(tokenAlphabet, tokenLen)
	if err != nil {
		return "", false
	}

	g.mu.Lock()
	g.sessions[token] = time.Now().Add(sessionTTL)
	g.mu.Unlock()

	// Burn the used code. If minting a replacement fails the old one stays put
	// rather than leaving the panel unreachable; the operator can restart.
	if err := g.rotate(); err != nil {
		return token, true
	}

	return token, true
}

// valid reports whether token names a live session, extending it on use.
func (g *authGate) valid(token string) bool {
	if token == "" {
		return false
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	exp, ok := g.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		delete(g.sessions, token)
		return false
	}
	g.sessions[token] = time.Now().Add(sessionTTL)

	return true
}

// allowLogin reports whether ip may attempt a login now, consuming budget from
// both limiters. It runs before the OTP is checked, so a successful attempt
// costs the same as a failed one.
//
// Two limiters apply: one per source IP, and one global. The global one is what
// actually protects a short code, since per-IP budget is trivially multiplied by
// an attacker with many addresses.
func (g *authGate) allowLogin(ip string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.global.Allow() {
		return false
	}

	lim, ok := g.limiters[ip]
	if !ok {
		// Cap the table so a spoofed-source flood can't grow it without bound.
		if len(g.limiters) > 1024 {
			g.limiters = make(map[string]*rate.Limiter)
		}
		lim = rate.NewLimiter(rate.Every(loginRefill), loginBurst)
		g.limiters[ip] = lim
	}

	return lim.Allow()
}

// currentOTP is for tests and diagnostics. It is never served over HTTP.
func (g *authGate) currentOTP() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.otp
}

// randomCode returns n characters drawn uniformly from alphabet using the
// CSPRNG. Neither alphabet in use has a power-of-two size, so a naive modulo
// would bias the draw; crypto/rand.Int rejection-samples internally and stays
// uniform for any alphabet size.
func randomCode(alphabet string, n int) (string, error) {
	max := big.NewInt(int64(len(alphabet))) //nolint
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[idx.Int64()]
	}

	return string(b), nil
}

// clientIP extracts a rate-limiting key. X-Forwarded-For is deliberately
// ignored: it is attacker-controlled unless a trusted proxy is known to rewrite
// it, and trusting it would let one client evade the login limiter entirely.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return host
}

// bearerToken pulls the session token from the Authorization header.
//
// A custom header, never HTTP Basic: browsers cache Basic credentials and
// resend them automatically, which would reintroduce the ambient credential
// this design exists to avoid.
func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if len(h) > 7 && strings.EqualFold(h[:7], "Bearer ") {
		return strings.TrimSpace(h[7:])
	}

	return ""
}

// loginHandler exchanges a valid OTP for a session token.
func (g *authGate) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			mWriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		if !g.allowLogin(clientIP(r)) {
			mWriteError(w, http.StatusTooManyRequests, "too many attempts, wait a minute")
			return
		}

		var req struct {
			OTP string `json:"otp"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			mWriteError(w, http.StatusBadRequest, "invalid login payload")
			return
		}

		token, ok := g.login(req.OTP)
		if !ok {
			// Say that the code is gone, not just wrong: the failure itself rotated
			// it, so an operator who retries what they just typed — or what is still
			// on their screen — cannot succeed. This leaks nothing; that a wrong
			// guess burns the code is documented behavior, not a secret.
			mWriteError(w, http.StatusUnauthorized,
				"invalid OTP — the code has been replaced, re-read it in the hypervisor app list")
			return
		}

		mWriteJSON(w, http.StatusOK, map[string]any{"token": token})
	}
}

// requireAuth wraps a handler so only requests carrying a live session token
// reach it.
func (g *authGate) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !g.valid(bearerToken(r)) {
			mWriteError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
