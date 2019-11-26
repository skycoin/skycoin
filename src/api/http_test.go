package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

const configuredHost = "127.0.0.1:6420"

var allAPISetsEnabled = map[string]struct{}{
	EndpointsRead:               struct{}{},
	EndpointsTransaction:        struct{}{},
	EndpointsStatus:             struct{}{},
	EndpointsWallet:             struct{}{},
	EndpointsInsecureWalletSeed: struct{}{},
	EndpointsPrometheus:         struct{}{},
	EndpointsNetCtrl:            struct{}{},
	EndpointsStorage:            struct{}{},
}

func defaultMuxConfig() muxConfig {
	return muxConfig{
		host:           configuredHost,
		appLoc:         ".",
		disableCSRF:    true,
		disableCSP:     true,
		enabledAPISets: allAPISetsEnabled,
	}
}

var endpointsMethods = map[string][]string{
	"/api/v1/address_uxouts": []string{
		http.MethodGet,
	},
	"/api/v1/addresscount": []string{
		http.MethodGet,
	},
	"/api/v1/balance": []string{
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/block": []string{
		http.MethodGet,
	},
	"/api/v1/blockchain/metadata": []string{
		http.MethodGet,
	},
	"/api/v1/blockchain/progress": []string{
		http.MethodGet,
	},
	"/api/v1/blocks": []string{
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/coinSupply": []string{
		http.MethodGet,
	},
	"/api/v1/health": []string{
		http.MethodGet,
	},
	"/api/v1/injectTransaction": []string{
		http.MethodPost,
	},
	"/api/v1/last_blocks": []string{
		http.MethodGet,
	},
	"/api/v1/version": []string{
		http.MethodGet,
	},
	"/api/v1/network/connection": []string{
		http.MethodGet,
	},
	"/api/v1/network/connections": []string{
		http.MethodGet,
	},
	"/api/v1/network/connections/exchange": []string{
		http.MethodGet,
	},
	"/api/v1/network/connections/trust": []string{
		http.MethodGet,
	},
	"/api/v1/network/defaultConnections": []string{
		http.MethodGet,
	},
	"/api/v1/network/connection/disconnect": []string{
		http.MethodPost,
	},
	"/api/v1/outputs": []string{
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/pendingTxs": []string{
		http.MethodGet,
	},
	"/api/v1/rawtx": []string{
		http.MethodGet,
	},
	"/api/v1/richlist": []string{
		http.MethodGet,
	},
	"/api/v1/resendUnconfirmedTxns": []string{
		http.MethodPost,
	},
	"/api/v1/transaction": []string{
		http.MethodGet,
	},
	"/api/v1/transactions": []string{
		http.MethodGet,
		http.MethodPost,
	},
	"/api/v1/uxout": []string{
		http.MethodGet,
	},
	"/api/v1/wallet": []string{
		http.MethodGet,
	},
	"/api/v1/wallet/balance": []string{
		http.MethodGet,
	},
	"/api/v1/wallet/create": []string{
		http.MethodPost,
	},
	"/api/v1/wallet/newAddress": []string{
		http.MethodPost,
	},
	"/api/v1/wallet/newSeed": []string{
		http.MethodGet,
	},
	"/api/v1/wallet/seed": []string{
		http.MethodPost,
	},
	"/api/v1/wallet/transaction": []string{
		http.MethodPost,
	},
	"/api/v1/wallet/transactions": []string{
		http.MethodGet,
	},
	"/api/v1/wallet/unload": []string{
		http.MethodPost,
	},
	"/api/v1/wallet/update": []string{
		http.MethodPost,
	},
	"/api/v1/wallets": []string{
		http.MethodGet,
	},
	"/api/v1/wallets/folderName": []string{
		http.MethodGet,
	},

	"/api/v2/transaction/verify": []string{
		http.MethodPost,
	},
	"/api/v2/address/verify": []string{
		http.MethodPost,
	},
	"/api/v2/wallet/recover": []string{
		http.MethodPost,
	},
	"/api/v2/wallet/seed/verify": []string{
		http.MethodPost,
	},
	"/api/v2/wallet/transaction/sign": []string{
		http.MethodPost,
	},
	"/api/v2/transaction": []string{
		http.MethodPost,
	},

	"/api/v2/data": []string{
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

// TestEnableGUI tests enable gui option, EnableGUI isn't part of Gateway API,
// we can't control the output by mocking the Gateway like other tests. Instead,
// we create a full webserver for each test case.
func TestEnableGUI(t *testing.T) {
	tt := []struct {
		name       string
		enableGUI  bool
		endpoint   string
		appLoc     string
		expectCode int
		expectBody string
	}{
		{
			name:       "disable gui GET /",
			enableGUI:  false,
			endpoint:   "/",
			appLoc:     "",
			expectCode: http.StatusNotFound,
			expectBody: "404 Not Found\n",
		},
		{
			name:       "disable gui GET /invalid-path",
			enableGUI:  false,
			endpoint:   "/invalid-path",
			appLoc:     "",
			expectCode: http.StatusNotFound,
			expectBody: "404 Not Found\n",
		},
		{
			name:       "enable gui GET /",
			enableGUI:  true,
			endpoint:   "/",
			appLoc:     "../gui/static",
			expectCode: http.StatusOK,
			expectBody: "",
		},
		{
			name:       "enable gui GET /invalid-path",
			enableGUI:  true,
			endpoint:   "/invalid-path",
			appLoc:     "../gui/static",
			expectCode: http.StatusNotFound,
			expectBody: "404 Not Found\n",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.endpoint, nil)
			require.NoError(t, err)

			gateway := &MockGatewayer{}

			rr := httptest.NewRecorder()
			cfg := defaultMuxConfig()
			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			c := Config{
				EnableGUI:   tc.enableGUI,
				DisableCSRF: true,
				StaticDir:   tc.appLoc,
			}

			host := "127.0.0.1:6423"
			s, err := Create(host, c, gateway)
			require.NoError(t, err)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := s.Serve()
				if err != nil && err.Error() != fmt.Sprintf("accept tcp %s: use of closed network connection", host) {
					require.NoError(t, err)
				}
			}()

			defer func() {
				s.listener.Close()
				wg.Wait()
			}()

			url := fmt.Sprintf("http://%s/%s", host, tc.endpoint)
			rsp, err := http.Get(url) //nolint:gosec
			require.NoError(t, err)

			defer rsp.Body.Close()
			require.Equal(t, tc.expectCode, rsp.StatusCode)

			body, err := ioutil.ReadAll(rr.Body)
			require.NoError(t, err)

			if rsp.StatusCode != http.StatusOK {
				require.Equal(t, tc.expectBody, string(body))
			}
		})
	}
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
