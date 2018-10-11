package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOriginRefererCheck(t *testing.T) {
	cases := []struct {
		name          string
		origin        string
		referer       string
		status        int
		err           string
		hostWhitelist []string
	}{
		{
			name:   "unparseable origin header",
			origin: ":?4foo",
			status: http.StatusForbidden,
			err:    "403 Forbidden - Invalid URL in Origin or Referer header\n",
		},
		{
			name:   "mismatched origin header",
			origin: "http://example.com/",
			status: http.StatusForbidden,
			err:    "403 Forbidden - Invalid Origin or Referer\n",
		},
		{
			name:    "mismatched referer header",
			referer: "http://example.com/",
			status:  http.StatusForbidden,
			err:     "403 Forbidden - Invalid Origin or Referer\n",
		},
		{
			name:          "whitelisted referer header",
			referer:       "http://example.com/",
			hostWhitelist: []string{"example.com"},
		},
		{
			name:          "whitelisted origin header",
			referer:       "http://example.com/",
			hostWhitelist: []string{"example.com"},
		},
	}

	for _, endpoint := range endpoints {
		for _, tc := range cases {
			name := fmt.Sprintf("%s %s", tc.name, endpoint)
			t.Run(name, func(t *testing.T) {
				gateway := &MockGatewayer{}

				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				require.NoError(t, err)

				csrfStore := &CSRFStore{
					Enabled: true,
				}
				setCSRFParameters(csrfStore, tokenValid, req)

				if tc.origin != "" {
					req.Header.Set("Origin", tc.origin)
				}
				if tc.referer != "" {
					req.Header.Set("Referer", tc.referer)
				}

				rr := httptest.NewRecorder()
				handler := newServerMux(muxConfig{
					host:            configuredHost,
					appLoc:          ".",
					enableJSON20RPC: true,
					disableCSP:      true,
					hostWhitelist:   tc.hostWhitelist,
				}, gateway, csrfStore, nil)

				handler.ServeHTTP(rr, req)

				switch tc.status {
				case http.StatusForbidden:
					require.Equal(t, tc.status, rr.Code)
					require.Equal(t, tc.err, rr.Body.String())
				default:
					// Arbitrary endpoints could return any status, since we don't customize the request per endpoint
					// Make sure that the request only didn't return the origin check error
					require.NotEqual(t, "403 Forbidden - Invalid URL in Origin or Referer header\n", rr.Body.String())
					require.NotEqual(t, "403 Forbidden - Invalid Origin or Referer\n", rr.Body.String())
				}
			})
		}
	}
}

func TestHostCheck(t *testing.T) {
	cases := []struct {
		name          string
		host          string
		status        int
		err           string
		hostWhitelist []string
	}{
		{
			name:   "invalid host",
			host:   "example.com",
			status: http.StatusForbidden,
			err:    "403 Forbidden - Invalid Host\n",
		},
		{
			name:          "invalid host is whitelisted",
			host:          "example.com",
			hostWhitelist: []string{"example.com"},
		},
	}

	for _, endpoint := range endpoints {
		for _, tc := range cases {
			name := fmt.Sprintf("%s %s", tc.name, endpoint)
			t.Run(name, func(t *testing.T) {
				gateway := &MockGatewayer{}

				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				require.NoError(t, err)

				csrfStore := &CSRFStore{
					Enabled: true,
				}
				setCSRFParameters(csrfStore, tokenValid, req)

				req.Host = "example.com"

				rr := httptest.NewRecorder()
				handler := newServerMux(muxConfig{
					host:            configuredHost,
					appLoc:          ".",
					enableJSON20RPC: true,
					disableCSP:      true,
					hostWhitelist:   tc.hostWhitelist,
				}, gateway, csrfStore, nil)

				handler.ServeHTTP(rr, req)

				switch tc.status {
				case http.StatusForbidden:
					require.Equal(t, http.StatusForbidden, rr.Code)
					require.Equal(t, tc.err, rr.Body.String())
				default:
					// Arbitrary endpoints could return any status, since we don't customize the request per endpoint
					// Make sure that the request only didn't return the invalid host error
					require.NotEqual(t, "403 Forbidden - Invalid Host\n", rr.Body.String())
				}
			})
		}
	}
}

func TestContentSecurityPolicy(t *testing.T) {
	tt := []struct {
		name            string
		endpoint        string
		enableCSP       bool
		appLoc          string
		expectCSPHeader string
		enableGUI       bool
	}{
		{
			name:            "enable CSP GET /",
			endpoint:        "/",
			enableCSP:       true,
			appLoc:          "../gui/static/dist",
			expectCSPHeader: "script-src 'self' 127.0.0.1",
			enableGUI:       true,
		},
		{
			name:            "disable CSP GET /",
			endpoint:        "/",
			enableCSP:       false,
			appLoc:          "../gui/static/dist",
			expectCSPHeader: "",
			enableGUI:       true,
		},
		{
			// Confirms that the /csrf api won't be affected by the csp setting
			name:            "enable CSP GET /csrf",
			endpoint:        "/api/v1/csrf",
			enableCSP:       true,
			appLoc:          "",
			expectCSPHeader: "",
			enableGUI:       false,
		},
		{
			// Confirms that the /version api won't be affected by the csp setting
			name:            "enable CSP GET /version",
			endpoint:        "/api/v1/version",
			enableCSP:       true,
			appLoc:          "",
			expectCSPHeader: "",
			enableGUI:       false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{
				host:       configuredHost,
				appLoc:     tc.appLoc,
				enableGUI:  tc.enableGUI,
				disableCSP: !tc.enableCSP,
			}, &MockGatewayer{}, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			csp := rr.Header().Get("Content-Security-Policy")
			require.Equal(t, tc.expectCSPHeader, csp)
		})
	}
}
