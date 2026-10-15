package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebWallet(t *testing.T) {
	cfg := defaultMuxConfig()
	cfg.webWallet = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("wallet")) //nolint:errcheck
	})
	cfg.webWalletUsername = "user"
	cfg.webWalletPassword = "pass"
	mux := newServerMux(cfg, &MockGatewayer{})

	serve := func(path, user, pass string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Host = configuredHost
		if user != "" {
			req.SetBasicAuth(user, pass)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	for _, path := range []string{"/", "/main.js", "/assets/scripts/skycoin-lite.wasm"} {
		require.Equal(t, http.StatusUnauthorized, serve(path, "", "").Code, path)
		require.Equal(t, http.StatusUnauthorized, serve(path, "user", "wrong").Code, path)

		rr := serve(path, "user", "pass")
		require.Equal(t, http.StatusOK, rr.Code, path)
		require.Equal(t, "wallet", rr.Body.String(), path)
	}

	// The wallet's password does not reach the API, which thin clients use.
	rr := serve("/api/v1/version", "", "")
	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEqual(t, "wallet", rr.Body.String())
}

func TestWebWalletAndGUIConflict(t *testing.T) {
	_, err := Create("127.0.0.1:0", Config{
		EnableGUI: true,
		WebWallet: http.NotFoundHandler(),
	}, &MockGatewayer{})
	require.Error(t, err)
}

func TestWebWalletCredentialsReachTheAPI(t *testing.T) {
	cfg := defaultMuxConfig()
	cfg.webWallet = http.NotFoundHandler()
	cfg.webWalletUsername = "user"
	cfg.webWalletPassword = "pass"
	mux := newServerMux(cfg, &MockGatewayer{})

	call := func(user, pass string) int {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/version", nil)
		req.Host = configuredHost
		req.SetBasicAuth(user, pass)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr.Code
	}

	// The browser sends the wallet's credentials with its API calls too.
	require.Equal(t, http.StatusOK, call("user", "pass"))
	// Any other credentials still tell the client the API has no auth.
	require.Equal(t, http.StatusUnauthorized, call("user", "wrong"))
}
