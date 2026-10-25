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
	mux := newServerMux(cfg, &MockGatewayer{})

	serve := func(path string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Host = configuredHost
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	for _, path := range []string{"/", "/main.js", "/assets/scripts/skycoin-lite.wasm"} {
		rr := serve(path)
		require.Equal(t, http.StatusOK, rr.Code, path)
		require.Equal(t, "wallet", rr.Body.String(), path)
	}

	// The API keeps its own routes beside the wallet.
	rr := serve("/api/v1/version")
	require.Equal(t, http.StatusOK, rr.Code)
	require.NotEqual(t, "wallet", rr.Body.String())
}
