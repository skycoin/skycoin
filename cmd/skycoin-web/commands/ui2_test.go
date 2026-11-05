package commands

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
)

func TestMountUI2(t *testing.T) {
	var asked []string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		asked = append(asked, r.URL.Path)
		var h api.HealthResponse
		h.CoinName, h.Version.Version = "skycoin", "0.28.8"
		h.BlockchainMetadata.Head.BkSeq = 4321
		switch r.URL.Path {
		case "/api/v1/health":
			_ = json.NewEncoder(w).Encode(h) //nolint:errcheck
		case "/api/v1/blockchain/progress":
			_, _ = w.Write([]byte(`{"current":4321,"highest":4321,"peers":[]}`)) //nolint:errcheck
		case "/api/v1/network/connections":
			_, _ = w.Write([]byte(`{"connections":[]}`)) //nolint:errcheck
		default:
			http.NotFound(w, r)
		}
	})
	mountUI2(mux)
	h := sameOriginWrites(mux)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/ui2/status", nil))
	require.Equal(t, http.StatusOK, rr.Code)
	require.Contains(t, rr.Body.String(), "4321", "the page reads the node through the /api proxy")
	require.Contains(t, asked, "/api/v1/health")

	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/ui2/", nil))
	require.Equal(t, http.StatusOK, rr.Code)
	require.True(t, strings.Contains(rr.Body.String(), `data-bw-page="list"`), "a node without the wallet API gets the browser wallets")

	req := httptest.NewRequest(http.MethodPost, "/ui2/bw/send/inject", strings.NewReader("rawtx=00"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	require.Equal(t, http.StatusForbidden, rr.Code, "the same-origin check covers the pages")
}
