package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/readable"
)

func TestUI2Mount(t *testing.T) {
	cfg := defaultMuxConfig()
	cfg.health.BuildInfo = readable.BuildInfo{Version: "9.9.9"}
	cfg.username, cfg.password = "user", "pass"
	cfg.ui2 = func(c *Client) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			v, err := c.Version()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_, _ = w.Write([]byte(r.URL.Path + " " + v.Version)) //nolint:errcheck,gosec
		})
	}
	mux := newServerMux(cfg, &MockGatewayer{})

	get := func(host, user, referer string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		req.Host = host
		if user != "" {
			req.SetBasicAuth(user, "pass")
		}
		if referer != "" {
			req.Header.Set("Referer", referer)
		}
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		return rr
	}

	rr := get(configuredHost, "user", "")
	require.Equal(t, http.StatusOK, rr.Code, rr.Body.String())
	require.Equal(t, "/status 9.9.9", rr.Body.String(), "the client reaches the API in process")

	require.Equal(t, http.StatusUnauthorized, get(configuredHost, "", "").Code)
	require.Equal(t, http.StatusUnauthorized, get(configuredHost, "other", "").Code)
	require.Equal(t, http.StatusForbidden, get("evil.example:6420", "user", "").Code, "DNS rebinding")
	require.Equal(t, http.StatusOK, get(configuredHost, "user", "https://example.org/").Code, "a link from another site opens a page")

	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	req.Host = configuredHost
	req.SetBasicAuth("user", "pass")
	req.Header.Set("Origin", "https://example.org")
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code, "a write from another origin")
}

func TestUI2MountCSRF(t *testing.T) {
	cfg := defaultMuxConfig()
	cfg.disableCSRF = false
	cfg.ui2 = func(*Client) http.Handler {
		return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	}
	mux := newServerMux(cfg, &MockGatewayer{})

	req := httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Host = configuredHost
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code, "a write without a CSRF token")

	token, err := newCSRFToken()
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/x", nil)
	req.Host = configuredHost
	req.Header.Set(CSRFHeaderName, token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}

func TestUI2MountPaths(t *testing.T) {
	cfg := defaultMuxConfig()
	cfg.webWallet = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("angular " + r.URL.Path)) //nolint:errcheck,gosec
	})
	cfg.ui2 = func(*Client) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ui2 " + r.URL.Path)) //nolint:errcheck,gosec
		})
	}
	mux := newServerMux(cfg, &MockGatewayer{})

	for _, tc := range []struct{ path, body, location string }{
		{path: "/", body: "ui2 /"},
		{path: "/wallet", body: "ui2 /wallet"},
		{path: "/angular/", body: "angular /"},
		{path: "/angular/main.js", body: "angular /main.js"},
		{path: "/assets/scripts/wasm_exec.js", body: "angular /assets/scripts/wasm_exec.js"},
		{path: "/favicon.ico", body: "angular /favicon.ico"},
		{path: "/ngsw-worker.js", body: "angular /safety-worker.js"},
		{path: "/ui2/send?id=w", body: `<a href="/send?id=w">Moved Permanently</a>.`, location: "/send?id=w"},
	} {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		req.Host = configuredHost
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		require.Equal(t, tc.body, strings.TrimSpace(rr.Body.String()), tc.path)
		require.Equal(t, tc.location, rr.Header().Get("Location"), tc.path)
	}

	cfg.ui2 = nil
	mux = newServerMux(cfg, &MockGatewayer{})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = configuredHost
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, "/angular/", rr.Header().Get("Location"), "without ui2 the root goes to the Angular wallet")
}
