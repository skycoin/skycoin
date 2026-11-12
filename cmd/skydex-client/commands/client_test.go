package commands

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

// The server rendered client has the root, the React one /angular/, and the
// /ui2/ paths the server rendered one was tried at move to the root.
func TestUIHandlerPaths(t *testing.T) {
	sess := newSession(fakeDialer{srv: newTestMarket(t)}, testMarketRef)
	t.Cleanup(sess.close)
	h, err := uiHandler(sess)
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		path, want, location string
		code                 int
	}{
		{path: "/", want: `href="angular/"`, code: http.StatusOK},
		{path: "/angular/", want: `src="/angular/assets/index.js"`, code: http.StatusOK},
		{path: "/angular/orders", want: `src="/angular/assets/index.js"`, code: http.StatusOK},
		{path: "/angular/assets/", want: `src="/angular/assets/index.js"`, code: http.StatusOK},
		{path: "/angular/assets/index.js", code: http.StatusOK},
		{path: "/ui2/market?new=1", location: "/market?new=1", code: http.StatusMovedPermanently},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Code != tc.code || !strings.Contains(rec.Body.String(), tc.want) || rec.Header().Get("Location") != tc.location {
			t.Errorf("GET %s: %d Location %q, want %d Location %q and a body with %q", tc.path, rec.Code, rec.Header().Get("Location"), tc.code, tc.location, tc.want)
		}
	}
}
