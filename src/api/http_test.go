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
	EndpointsRead:                  struct{}{},
	EndpointsStatus:                struct{}{},
	EndpointsWallet:                struct{}{},
	EndpointsInsecureWalletSeed:    struct{}{},
	EndpointsDeprecatedWalletSpend: struct{}{},
}

func defaultMuxConfig() muxConfig {
	return muxConfig{
		host:           configuredHost,
		appLoc:         ".",
		disableCSP:     true,
		enabledAPISets: allAPISetsEnabled,
	}
}

var endpoints = []string{
	"/address_uxouts",
	"/addresscount",
	"/balance",
	"/block",
	"/blockchain/metadata",
	"/blockchain/progress",
	"/blocks",
	"/coinSupply",
	"/explorer/address",
	"/health",
	"/injectTransaction",
	"/last_blocks",
	"/version",
	"/network/connection",
	"/network/connections",
	"/network/connections/exchange",
	"/network/connections/trust",
	"/network/defaultConnections",
	"/outputs",
	"/pendingTxs",
	"/rawtx",
	"/richlist",
	"/resendUnconfirmedTxns",
	"/transaction",
	"/transactions",
	"/uxout",
	"/wallet",
	"/wallet/balance",
	"/wallet/create",
	"/wallet/newAddress",
	"/wallet/newSeed",
	"/wallet/seed",
	"/wallet/spend",
	"/wallet/transaction",
	"/wallet/transactions",
	"/wallet/unload",
	"/wallet/update",
	"/wallets",
	"/wallets/folderName",
	"/webrpc",

	"/api/v1/address_uxouts",
	"/api/v1/addresscount",
	"/api/v1/balance",
	"/api/v1/block",
	"/api/v1/blockchain/metadata",
	"/api/v1/blockchain/progress",
	"/api/v1/blocks",
	"/api/v1/coinSupply",
	"/api/v1/explorer/address",
	"/api/v1/health",
	"/api/v1/injectTransaction",
	"/api/v1/last_blocks",
	"/api/v1/version",
	"/api/v1/network/connection",
	"/api/v1/network/connections",
	"/api/v1/network/connections/exchange",
	"/api/v1/network/connections/trust",
	"/api/v1/network/defaultConnections",
	"/api/v1/outputs",
	"/api/v1/pendingTxs",
	"/api/v1/rawtx",
	"/api/v1/richlist",
	"/api/v1/resendUnconfirmedTxns",
	"/api/v1/transaction",
	"/api/v1/transactions",
	"/api/v1/uxout",
	"/api/v1/wallet",
	"/api/v1/wallet/balance",
	"/api/v1/wallet/create",
	"/api/v1/wallet/newAddress",
	"/api/v1/wallet/newSeed",
	"/api/v1/wallet/seed",
	"/api/v1/wallet/spend",
	"/api/v1/wallet/transaction",
	"/api/v1/wallet/transactions",
	"/api/v1/wallet/unload",
	"/api/v1/wallet/update",
	"/api/v1/wallets",
	"/api/v1/wallets/folderName",
	"/api/v1/webrpc",

	"/api/v2/transaction/verify",
	"/api/v2/address/verify",
	"/api/v2/wallet/recover",
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
			handler := newServerMux(muxConfig{
				host:       configuredHost,
				appLoc:     tc.appLoc,
				disableCSP: true,
			}, gateway, &CSRFStore{}, nil)
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
			rsp, err := http.Get(url)
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
	for _, e := range append(endpoints, []string{"/csrf", "/api/v1/csrf"}...) {
		switch e {
		case "/webrpc", "/api/v1/webrpc":
			continue
		}

		t.Run(e, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, e, nil)
			require.NoError(t, err)

			cfg := defaultMuxConfig()
			cfg.enableUnversionedAPI = true
			cfg.enableJSON20RPC = false
			cfg.enabledAPISets = map[string]struct{}{} // disable all API sets

			handler := newServerMux(cfg, &MockGatewayer{}, &CSRFStore{
				Enabled: true,
			}, nil)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			switch e {
			case "/csrf", "/api/v1/csrf", "/version", "/api/v1/version": // always enabled
				require.Equal(t, http.StatusOK, rr.Code)
			default:
				require.Equal(t, http.StatusForbidden, rr.Code)
				require.Equal(t, "403 Forbidden - Endpoint is disabled", strings.TrimSpace(rr.Body.String()))
			}
		})
	}
}

func TestCORS(t *testing.T) {
	// Make sure cross origin requests are blocked by default
	// Make sure cross origin requests are blocked when using a whitelist
	// Make sure regular requests work
	//		-- GET request
	//		-- POST request with CSRF token
	// Make sure cross origin requests work if whitelisted
	//		-- GET request
	//		-- POST request with CSRF token
	// Make sure OPTIONS responds for all endpoints

	// getCSRFTokenFromOrigin := func(t *testing.T, origin string, cfg muxConfig, csrfStore *CSRFStore) string {
	// 	rr := makeCSRFTokenRequestFromOrigin(t, origin, cfg, csrfStore)

	// 	require.Equal(t, http.StatusOK, rr.Code)

	// 	var msg map[string]string
	// 	err := json.Unmarshal(rr.Body.Bytes(), &msg)
	// 	require.NoError(t, err)

	// 	token := msg["csrf_token"]
	// 	require.NotEmpty(t, token)

	// 	return token
	// }

	// makePOSTRequestFromOrigin := func(t *testing.T, origin string, cfg muxConfig, csrfStore *CSRFStore, token string) {
	// 	// Make a POST request with the CSRF token
	// 	// The POST request is sent to /api/v2/address/verify since this has no side effects to handle
	// 	body := `{"address":"7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD"}`
	// 	req, err := http.NewRequest(http.MethodPost, "/api/v2/address/verify", bytes.NewBufferString(body))
	// 	require.NoError(t, err)

	// 	req.Header.Set("Origin", fmt.Sprintf("http://%s", origin))
	// 	req.Header.Set(CSRFHeaderName, token)
	// 	req.Header.Set("Content-Type", "application/json")

	// 	handler := newServerMux(cfg, &MockGatewayer{}, csrfStore, nil)

	// 	rr := httptest.NewRecorder()
	// 	handler.ServeHTTP(rr, req)

	// 	require.Equal(t, http.StatusOK, rr.Code)
	// }

	// testSameOriginRequestOK := func(t *testing.T) {
	// 	// Tests that same origin requests are fine,
	// 	// GET /csrf, then use CSRF to POST something
	// 	cfg := defaultMuxConfig()
	// 	csrfStore := &CSRFStore{
	// 		Enabled: true,
	// 	}
	// 	token := getCSRFTokenFromOrigin(t, cfg.host, cfg, csrfStore)
	// 	makePOSTRequestFromOrigin(t, cfg.host, cfg, csrfStore, token)
	// }

	// testCrossOriginWhitelistedRequestOK := func(t *testing.T) {
	// 	// Tests that whitelist cross origin requests are fine,
	// 	// GET /csrf, then use CSRF to POST something
	// 	cfg := defaultMuxConfig()
	// 	cfg.hostWhitelist = []string{"example.com"}
	// 	csrfStore := &CSRFStore{
	// 		Enabled: true,
	// 	}
	// 	token := getCSRFTokenFromOrigin(t, "example.com", cfg, csrfStore)
	// 	makePOSTRequestFromOrigin(t, "example.com", cfg, csrfStore, token)
	// }

	// testCrossOriginNotAllowed := func(t *testing.T) {
	// 	// Tests that cross origin requests are not allowed by default
	// 	cfg := defaultMuxConfig()
	// 	csrfStore := &CSRFStore{
	// 		Enabled: true,
	// 	}
	// 	rr := makeCSRFTokenRequestFromOrigin(t, "example.com", cfg, csrfStore)
	// 	require.Equal(t, http.StatusForbidden, rr.Code)
	// 	require.Equal(t, "foo", rr.Body.String())

	// 	// makePOSTRequestFromOrigin(t, "example.com", cfg, csrfStore, token)
	// }

	// makeRequestFromOrigin := func(t *testing.T, method, endpoint, origin string, cfg muxConfig) *httptest.ResponseRecorder {
	// 	req, err := http.NewRequest(method, endpoint, nil)
	// 	require.NoError(t, err)

	// 	csrfStore := &CSRFStore{
	// 		Enabled: true,
	// 	}
	// 	setCSRFParameters(csrfStore, tokenValid, req)

	// 	req.Header.Set("Origin", fmt.Sprintf("http://%s", origin))

	// 	handler := newServerMux(cfg, &MockGatewayer{}, csrfStore, nil)

	// 	rr := httptest.NewRecorder()
	// 	handler.ServeHTTP(rr, req)

	// 	return rr
	// }

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

	for _, e := range append(endpoints, "/api/v1/csrf") {
		if !strings.HasPrefix(e, "/api/v") {
			continue
		}

		for _, tc := range cases {
			for _, m := range []string{http.MethodPost, http.MethodGet} {
				name := fmt.Sprintf("%s %s %s", tc.name, m, e)
				t.Run(name, func(t *testing.T) {
					cfg := defaultMuxConfig()
					cfg.hostWhitelist = tc.hostWhitelist

					req, err := http.NewRequest(http.MethodOptions, e, nil)
					require.NoError(t, err)

					csrfStore := &CSRFStore{
						Enabled: true,
					}
					setCSRFParameters(csrfStore, tokenValid, req)

					req.Header.Set("Origin", fmt.Sprintf("http://%s", tc.origin))
					req.Header.Set("Access-Control-Request-Method", m)

					requestHeaders := strings.ToLower(fmt.Sprintf("%s, Content-Type", CSRFHeaderName))
					req.Header.Set("Access-Control-Request-Headers", requestHeaders)

					handler := newServerMux(cfg, &MockGatewayer{}, csrfStore, nil)

					rr := httptest.NewRecorder()
					handler.ServeHTTP(rr, req)

					resp := rr.Result()

					fmt.Println(resp.Header)
					fmt.Println(rr.Body.String())

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
