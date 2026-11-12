package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const configuredHost = "127.0.0.1:6420"

var allAPISetsEnabled = map[string]struct{}{
	EndpointsRead:               {},
	EndpointsTransaction:        {},
	EndpointsStatus:             {},
	EndpointsWallet:             {},
	EndpointsInsecureWalletSeed: {},
	EndpointsNetCtrl:            {},
	EndpointsStorage:            {},
	EndpointsExplorer:           {},
}

func defaultMuxConfig() muxConfig {
	return muxConfig{
		host:           configuredHost,
		disableCSRF:    true,
		disableCSP:     true,
		enabledAPISets: allAPISetsEnabled,
	}
}

var endpointsMethods = map[string][]string{
	"/api/v1/address_uxouts": {
		http.MethodGet,
	},
	"/api/v1/addresscount": {
		http.MethodGet,
	},
	"/api/v1/balance": {
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/block": {
		http.MethodGet,
	},
	"/api/v1/blockchain/metadata": {
		http.MethodGet,
	},
	"/api/v1/blockchain/progress": {
		http.MethodGet,
	},
	"/api/v1/blocks": {
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/coinSupply": {
		http.MethodGet,
	},
	"/api/v1/explorer/address": {
		http.MethodGet,
	},
	"/api/v1/health": {
		http.MethodGet,
	},
	"/api/v1/injectTransaction": {
		http.MethodPost,
	},
	"/api/v1/last_blocks": {
		http.MethodGet,
	},
	"/api/v1/version": {
		http.MethodGet,
	},
	"/api/v1/network/connection": {
		http.MethodGet,
	},
	"/api/v1/network/connections": {
		http.MethodGet,
	},
	"/api/v1/network/connections/exchange": {
		http.MethodGet,
	},
	"/api/v1/network/connections/trust": {
		http.MethodGet,
	},
	"/api/v1/network/defaultConnections": {
		http.MethodGet,
	},
	"/api/v1/network/connection/disconnect": {
		http.MethodPost,
	},
	"/api/v1/outputs": {
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/pendingTxs": {
		http.MethodGet,
	},
	"/api/v1/rawtx": {
		http.MethodGet,
	},
	"/api/v1/richlist": {
		http.MethodGet,
	},
	"/api/v1/resendUnconfirmedTxns": {
		http.MethodPost,
	},
	"/api/v1/transaction": {
		http.MethodGet,
	},
	"/api/v1/transactions": {
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/uxout": {
		http.MethodGet,
	},
	"/api/v1/wallet": {
		http.MethodGet,
	},
	"/api/v1/wallet/balance": {
		http.MethodGet,
	},
	"/api/v1/wallet/create": {
		http.MethodPost,
	},
	"/api/v1/wallet/newAddress": {
		http.MethodPost,
	},
	"/api/v1/wallet/newSeed": {
		http.MethodGet,
	},
	"/api/v1/wallet/seed": {
		http.MethodPost,
	},
	"/api/v1/wallet/transaction": {
		http.MethodPost,
	},
	"/api/v1/wallet/transactions": {
		http.MethodGet,
	},
	"/api/v1/wallet/unload": {
		http.MethodPost,
	},
	"/api/v1/wallet/newAccount": {
		http.MethodPost,
	},
	"/api/v1/wallet/update": {
		http.MethodPost,
	},
	"/api/v1/wallets": {
		http.MethodGet,
	},
	"/api/v1/wallets/folderName": {
		http.MethodGet,
	},

	"/api/v2/transaction/verify": {
		http.MethodPost,
	},
	"/api/v2/address/verify": {
		http.MethodPost,
	},
	"/api/v2/wallet/recover": {
		http.MethodPost,
	},
	"/api/v2/wallet/seed/verify": {
		http.MethodPost,
	},
	"/api/v2/wallet/transaction/sign": {
		http.MethodPost,
	},
	"/api/v2/transaction": {
		http.MethodPost,
	},

	"/api/v2/data": {
		http.MethodGet,
		http.MethodPost,
		http.MethodDelete,
	},
}

func allEndpoints() []string {
	endpoints := make([]string, len(endpointsMethods))
	i := 0
	for e := range endpointsMethods {
		endpoints[i] = e
		i++
	}
	return endpoints
}

func TestAPISetDisabled(t *testing.T) {
	tf := func(t *testing.T, endpoint, method string, disableCSRF bool) {
		req, err := http.NewRequest(method, endpoint, nil)
		require.NoError(t, err)

		isAPIV2 := strings.HasPrefix(endpoint, "/api/v2/")
		if isAPIV2 {
			req.Header.Set("Content-Type", ContentTypeJSON)
		}

		cfg := defaultMuxConfig()
		cfg.disableCSRF = disableCSRF
		cfg.enabledAPISets = map[string]struct{}{} // disable all API sets

		handler := newServerMux(cfg, &MockGatewayer{})

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		switch endpoint {
		case "/api/v1/csrf", "/api/v1/version": // always enabled
			require.Equal(t, http.StatusOK, rr.Code)
		default:
			require.Equal(t, http.StatusForbidden, rr.Code)
			if isAPIV2 {
				require.Equal(t, "{\n    \"error\": {\n        \"message\": \"Endpoint is disabled\",\n        \"code\": 403\n    }\n}", rr.Body.String())
			} else {
				require.Equal(t, "403 Forbidden - Endpoint is disabled", strings.TrimSpace(rr.Body.String()))
			}
		}
	}

	for e, methods := range endpointsMethods {
		for _, m := range methods {
			t.Run(fmt.Sprintf("%s %s", m, e), func(t *testing.T) {
				tf(t, e, m, true)
			})
		}
	}

	t.Run("GET /api/v1/csrf", func(t *testing.T) {
		tf(t, "/api/v1/csrf", http.MethodGet, false)
	})
}

func TestCORS(t *testing.T) {
	cases := []struct {
		name          string
		origin        string
		hostWhitelist []string
		valid         bool
	}{
		{
			name:   "options no whitelist",
			origin: configuredHost,
			valid:  true,
		},
		{
			name:          "options whitelist",
			origin:        "example.com",
			hostWhitelist: []string{"example.com"},
			valid:         true,
		},
		{
			name:   "options no whitelist not whitelisted",
			origin: "example.com",
			valid:  false,
		},
	}

	for _, e := range append(allEndpoints(), "/api/v1/csrf") {
		for _, tc := range cases {
			for _, m := range []string{http.MethodPost, http.MethodGet} {
				name := fmt.Sprintf("%s %s %s", tc.name, m, e)
				t.Run(name, func(t *testing.T) {
					cfg := defaultMuxConfig()
					cfg.disableCSRF = false
					cfg.hostWhitelist = tc.hostWhitelist

					req, err := http.NewRequest(http.MethodOptions, e, nil)
					require.NoError(t, err)

					setCSRFParameters(t, tokenValid, req)

					isAPIV2 := strings.HasPrefix(e, "/api/v2/")
					if isAPIV2 {
						req.Header.Set("Content-Type", ContentTypeJSON)
					}

					req.Header.Set("Origin", fmt.Sprintf("http://%s", tc.origin))
					req.Header.Set("Access-Control-Request-Method", m)

					requestHeaders := strings.ToLower(fmt.Sprintf("%s, Content-Type", CSRFHeaderName))
					req.Header.Set("Access-Control-Request-Headers", requestHeaders)

					handler := newServerMux(cfg, &MockGatewayer{})

					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, req)

					resp := rr.Result()

					allowOrigins := resp.Header.Get("Access-Control-Allow-Origin")
					allowHeaders := resp.Header.Get("Access-Control-Allow-Headers")
					allowMethods := resp.Header.Get("Access-Control-Allow-Methods")

					if tc.valid {
						require.Equal(t, fmt.Sprintf("http://%s", tc.origin), allowOrigins)
						require.Equal(t, requestHeaders, strings.ToLower(allowHeaders))
						require.Equal(t, m, allowMethods)
					} else {
						require.Empty(t, allowOrigins)
						require.Empty(t, allowHeaders)
						require.Empty(t, allowMethods)
					}

					allowCreds := resp.Header.Get("Access-Control-Allow-Credentials")
					require.Empty(t, allowCreds)
				})
			}
		}
	}
}

func TestHTTPBasicAuthInvalid(t *testing.T) {
	username := "foo"
	badUsername := "foof"
	password := "bar"
	badPassword := "barb"

	userPassCombos := []struct {
		u, p string
	}{
		{},
		{
			u: username,
		},
		{
			p: password,
		},
		{
			u: username,
			p: password,
		},
	}

	reqUserPassCombos := []struct {
		u, p string
	}{
		{},
		{
			u: username,
		},
		{
			u: badUsername,
		},
		{
			p: password,
		},
		{
			p: badPassword,
		},
		{
			u: username,
			p: password,
		},
		{
			u: badUsername,
			p: badPassword,
		},
		{
			u: username,
			p: badPassword,
		},
		{
			u: badUsername,
			p: password,
		},
	}

	type testCase struct {
		username    string
		password    string
		reqUsername string
		reqPassword string
		authorized  bool
	}

	cases := []testCase{}

	for _, a := range userPassCombos {
		for _, b := range reqUserPassCombos {
			cases = append(cases, testCase{
				username:    a.u,
				password:    a.p,
				reqUsername: b.u,
				reqPassword: b.p,
				authorized:  a.u == b.u && a.p == b.p,
			})
		}
	}

	for _, e := range append(allEndpoints(), []string{"/api/v1/csrf"}...) {
		for _, tc := range cases {
			name := fmt.Sprintf("u=%s p=%s ru=%s rp=%s auth=%v e=%s", tc.username, tc.password, tc.reqUsername, tc.reqPassword, tc.authorized, e)
			t.Run(name, func(t *testing.T) {
				// Use a made-up request method so that any authorized request
				// is guaranteed to fail before it reaches the mock gateway,
				// which will panic without the mocks configured
				req, err := http.NewRequest("FOOBAR", e, nil)
				require.NoError(t, err)

				req.SetBasicAuth(tc.reqUsername, tc.reqPassword)

				cfg := defaultMuxConfig()
				cfg.disableCSRF = false
				cfg.username = tc.username
				cfg.password = tc.password

				handler := newServerMux(cfg, &MockGatewayer{})

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if !tc.authorized {
					require.Equal(t, http.StatusUnauthorized, rr.Code)
					if strings.HasPrefix(e, "/api/v2") {
						require.Equal(t, "{\n    \"error\": {\n        \"message\": \"Unauthorized\",\n        \"code\": 401\n    }\n}", rr.Body.String())
					} else {
						require.Equal(t, "401 Unauthorized", strings.TrimSpace(rr.Body.String()))
					}
				} else {
					require.NotEqual(t, http.StatusUnauthorized, rr.Code)
				}
			})
		}
	}
}

// //////////////////////////////////////////////////////////////
// Test helper tools
// //////////////////////////////////////////////////////////////
type httpMockClient struct {
	gateway     *MockGatewayer
	contentType string
	disableCSRF bool
}

func newHTTPMockClient(gateway *MockGatewayer, contentType string, disableCSRF bool) *httpMockClient {
	return &httpMockClient{
		gateway:     gateway,
		contentType: contentType,
		disableCSRF: disableCSRF,
	}
}

func (c httpMockClient) Do(t *testing.T, method, endpoint, value string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, endpoint, strings.NewReader(value))
	require.NoError(t, err)
	req.Header.Add("Content-Type", c.contentType)
	if c.disableCSRF {
		setCSRFParameters(t, tokenInvalid, req)
	} else {
		setCSRFParameters(t, tokenValid, req)
	}

	rr := httptest.NewRecorder()

	cfg := defaultMuxConfig()
	cfg.disableCSRF = c.disableCSRF

	handler := newServerMux(cfg, c.gateway)
	handler.ServeHTTP(rr, req)
	return rr
}
