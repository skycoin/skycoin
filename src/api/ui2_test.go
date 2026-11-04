package api

import (
	"net/http"
	"net/http/httptest"
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
		req := httptest.NewRequest(http.MethodGet, "/ui2/status", nil)
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
	require.Equal(t, "/status 9.9.9", rr.Body.String(), "the prefix is stripped and the client reaches the API in process")

	require.Equal(t, http.StatusUnauthorized, get(configuredHost, "", "").Code)
	require.Equal(t, http.StatusUnauthorized, get(configuredHost, "other", "").Code)
	require.Equal(t, http.StatusForbidden, get("evil.example:6420", "user", "").Code, "DNS rebinding")
	require.Equal(t, http.StatusOK, get(configuredHost, "user", "https://example.org/").Code, "a link from another site opens a page")

	req := httptest.NewRequest(http.MethodPost, "/ui2/status", nil)
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

	req := httptest.NewRequest(http.MethodPost, "/ui2/x", nil)
	req.Host = configuredHost
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code, "a write without a CSRF token")

	token, err := newCSRFToken()
	require.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/ui2/x", nil)
	req.Host = configuredHost
	req.Header.Set(CSRFHeaderName, token)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
}
