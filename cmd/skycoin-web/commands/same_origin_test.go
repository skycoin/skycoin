package commands

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSameOriginWrites(t *testing.T) {
	h := sameOriginWrites(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	for _, tc := range []struct {
		method, site, origin string
		code                 int
	}{
		{http.MethodGet, "cross-site", "https://evil.example", http.StatusNoContent},
		{http.MethodPost, "same-origin", "http://127.0.0.1:8000", http.StatusNoContent},
		{http.MethodPost, "none", "", http.StatusNoContent},
		{http.MethodPost, "", "", http.StatusNoContent},
		{http.MethodPost, "", "http://127.0.0.1:8000", http.StatusNoContent},
		{http.MethodPost, "cross-site", "https://evil.example", http.StatusForbidden},
		{http.MethodPut, "same-site", "http://localhost:9000", http.StatusForbidden},
		{http.MethodDelete, "", "https://evil.example", http.StatusForbidden},
		{http.MethodPost, "", "null", http.StatusForbidden},
	} {
		req := httptest.NewRequest(tc.method, "http://127.0.0.1:8000/coin/1/api/v1/eth/send", nil)
		if tc.site != "" {
			req.Header.Set("Sec-Fetch-Site", tc.site)
		}
		if tc.origin != "" {
			req.Header.Set("Origin", tc.origin)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		require.Equal(t, tc.code, rec.Code, "%s %q %q", tc.method, tc.site, tc.origin)
	}
}
