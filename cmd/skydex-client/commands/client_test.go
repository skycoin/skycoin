package commands

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// A page on another site must not be able to place or cancel orders through
// the API the React client uses.
func TestUIHandlerRefusesCrossSiteWrites(t *testing.T) {
	sess := newSession(fakeDialer{srv: newTestMarket(t)}, testMarketRef)
	t.Cleanup(sess.close)
	h, err := uiHandler(sess)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		site   string
		origin string
		want   bool // refused
	}{
		{site: "cross-site", want: true},
		{origin: "http://evil.example", want: true},
		{site: "same-origin", want: false},
		{origin: "http://example.com", want: false}, // httptest's default Host
		{want: false}, // no browser headers: curl, the CLI
	} {
		for _, path := range []string{"/api/buy", "/api/orders/cancel", "/api/disconnect"} {
			req := httptest.NewRequest(http.MethodPost, path, nil)
			if tc.site != "" {
				req.Header.Set("Sec-Fetch-Site", tc.site)
			}
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if refused := rec.Code == http.StatusForbidden; refused != tc.want {
				t.Errorf("POST %s site=%q origin=%q: status %d, refused %v, want refused %v", path, tc.site, tc.origin, rec.Code, refused, tc.want)
			}
		}
	}

	// Reads stay open to any page, as before.
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code == http.StatusForbidden {
		t.Errorf("GET /api/status cross-site refused, reads should pass")
	}
}
