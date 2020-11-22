package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOriginRefererCheck(t *testing.T) {
	cases := []struct {
		name              string
		origin            string
		referer           string
		enableHeaderCheck bool
		status            int
		errV1             string
		errV2             string
		hostWhitelist     []string
	}{
		{
			name:              "unparseable origin header",
			origin:            ":?4foo",
			enableHeaderCheck: true,
			status:            http.StatusForbidden,
			errV1:             "403 Forbidden - Invalid URL in Origin or Referer header\n",
			errV2:             "{\n    \"error\": {\n        \"message\": \"Invalid URL in Origin or Referer header\",\n        \"code\": 403\n    }\n}",
		},
		{
			name:              "mismatched origin header",
			origin:            "http://example.com/",
			enableHeaderCheck: true,
			status:            http.StatusForbidden,
			errV1:             "403 Forbidden - Invalid Origin or Referer\n",
			errV2:             "{\n    \"error\": {\n        \"message\": \"Invalid Origin or Referer\",\n        \"code\": 403\n    }\n}",
		},
		{
			name:              "mismatched referer header",
			referer:           "http://example.com/",
			enableHeaderCheck: true,
			status:            http.StatusForbidden,
			errV1:             "403 Forbidden - Invalid Origin or Referer\n",
			errV2:             "{\n    \"error\": {\n        \"message\": \"Invalid Origin or Referer\",\n        \"code\": 403\n    }\n}",
		},
		{
			name:              "whitelisted referer header",
			referer:           "http://example.com/",
			enableHeaderCheck: true,
			hostWhitelist:     []string{"example.com"},
		},
		{
			name:              "whitelisted origin header",
			referer:           "http://example.com/",
			enableHeaderCheck: true,
			hostWhitelist:     []string{"example.com"},
		},
		{
			name:              "mismatched referer header",
			referer:           "http://example.com/",
			enableHeaderCheck: false,
		},
		{
			name:              "mismatched origin header",
			origin:            "http://example.com/",
			enableHeaderCheck: false,
		},
	}

	for _, endpoint := range allEndpoints() {
		for _, tc := range cases {
			name := fmt.Sprintf("%s %s", tc.name, endpoint)
			t.Run(name, func(t *testing.T) {
				gateway := &MockGatewayer{}

				req, err := http.NewRequest(http.MethodGet, endpoint, nil)
				require.NoError(t, err)

				setCSRFParameters(t, tokenValid, req)

				isAPIV2 := strings.HasPrefix(endpoint, "/api/v2")
				if isAPIV2 {
					req.Header.Set("Content-Type", ContentTypeJSON)
				}

				if tc.origin != "" {
					req.Header.Set("Origin", tc.origin)
				}
				if tc.referer != "" {
					req.Header.Set("Referer", tc.referer)
				}

				rr := httptest.NewRecorder()

				cfg := defaultMuxConfig()
				cfg.disableCSRF = false
				cfg.disableHeaderCheck = !tc.enableHeaderCheck
				// disable all api sets to avoid mocking gateway methods
				cfg.enabledAPISets = map[string]struct{}{}

				handler := newServerMux(cfg, gateway)
				handler.ServeHTTP(rr, req)

				switch tc.status {
				case http.StatusForbidden:
					require.Equal(t, tc.status, rr.Code)

					if isAPIV2 {
						require.Equal(t, tc.errV2, rr.Body.String())
					} else {
						require.Equal(t, tc.errV1, rr.Body.String())
					}
				default:
					if tc.enableHeaderCheck || tc.hostWhitelist == nil {
						// Arbitrary endpoints could return any status, since we don't customize the request per endpoint
						// Make sure that the request only didn't return the origin check error
						require.False(t, strings.Contains("Invalid URL in Origin or Referer header", rr.Body.String()))
						require.False(t, strings.Contains("Invalid Origin or Referer", rr.Body.String()))
					}
				}

			})
		}
	}
}

func TestHostCheck(t *testing.T) {
	cases := []struct {
		name              string
		host              string
		status            int
		enableHeaderCheck bool
		errV1             string
		errV2             string
		hostWhitelist     []string
	}{
		{
			name:              "invalid host",
			host:              "example.com",
			status:            http.StatusForbidden,
			enableHeaderCheck: true,
			errV1:             "403 Forbidden - Invalid Host\n",
			errV2:             "{\n    \"error\": {\n        \"message\": \"Invalid Host\",\n        \"code\": 403\n    }\n}",
		},
		{
			name:              "invalid host is whitelisted",
			host:              "example.com",
			hostWhitelist:     []string{"example.com"},
			enableHeaderCheck: true,
		},
		{
			name:              "invalid host - header check disabled",
			host:              "example.com",
			enableHeaderCheck: false,
		},
	}

	for endpoint, methods := range endpointsMethods {
		for _, m := range methods {
			for _, tc := range cases {
				name := fmt.Sprintf("%s %s %s", tc.name, m, endpoint)
				t.Run(name, func(t *testing.T) {
					gateway := &MockGatewayer{}

					req, err := http.NewRequest(m, endpoint, nil)
					require.NoError(t, err)

					setCSRFParameters(t, tokenValid, req)

					isAPIV2 := strings.HasPrefix(endpoint, "/api/v2")
					if isAPIV2 {
						req.Header.Set("Content-Type", ContentTypeJSON)
					}

					req.Host = "example.com"

					rr := httptest.NewRecorder()
					handler := newServerMux(muxConfig{
						host:               configuredHost,
						appLoc:             ".",
						disableCSRF:        false,
						disableHeaderCheck: !tc.enableHeaderCheck,
						disableCSP:         true,
						hostWhitelist:      tc.hostWhitelist,
					}, gateway)

					handler.ServeHTTP(rr, req)

					switch tc.status {
					case http.StatusForbidden:
						require.Equal(t, http.StatusForbidden, rr.Code)
						if isAPIV2 {
							require.Equal(t, tc.errV2, rr.Body.String())
						} else {
							require.Equal(t, tc.errV1, rr.Body.String())
						}
					default:
						if tc.enableHeaderCheck || tc.hostWhitelist == nil {
							// Arbitrary endpoints could return any status, since we don't customize the request per endpoint
							// Make sure that the request only didn't return the invalid host error
							require.False(t, strings.Contains("Invalid Host", rr.Body.String()))
						}
					}
				})
			}
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
			name:      "enable CSP GET /",
			endpoint:  "/",
			enableCSP: true,
			appLoc:    "../gui/static/dist",
			expectCSPHeader: "default-src 'self'; connect-src 'self' https://api.coinpaprika.com https://swaplab.cc https://version.skycoin.com https://downloads.skycoin.com http://127.0.0.1:9510; img-src 'self' 'unsafe-inline' data:; style-src 'self' 'unsafe-inline'; object-src	'none'; form-action 'none'; frame-ancestors 'none'; block-all-mixed-content; base-uri 'self'",
			enableGUI: true,
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
				host:        configuredHost,
				appLoc:      tc.appLoc,
				enableGUI:   tc.enableGUI,
				disableCSP:  !tc.enableCSP,
				disableCSRF: true,
			}, &MockGatewayer{})
			handler.ServeHTTP(rr, req)

			csp := rr.Header().Get("Content-Security-Policy")
			require.Equal(t, tc.expectCSPHeader, csp)
		})
	}
}

func TestIsContentTypeJSON(t *testing.T) {
	require.True(t, isContentTypeJSON(ContentTypeJSON))
	require.True(t, isContentTypeJSON("application/json"))
	require.True(t, isContentTypeJSON("application/json; charset=utf-8"))
	require.False(t, isContentTypeJSON("application/x-www-form-urlencoded"))
	require.False(t, isContentTypeJSON(ContentTypeForm))
}
