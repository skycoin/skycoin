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

	for _, e := range append(endpoints, []string{"/csrf", "/api/v1/csrf"}...) {
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
				cfg.enableUnversionedAPI = true
				cfg.enableJSON20RPC = false
				cfg.username = tc.username
				cfg.password = tc.password

				handler := newServerMux(cfg, &MockGatewayer{}, &CSRFStore{
					Enabled: true,
				}, nil)

				rr := httptest.NewRecorder()
				handler.ServeHTTP(rr, req)

				if !tc.authorized {
					require.Equal(t, http.StatusUnauthorized, rr.Code)
					require.Equal(t, "401 Unauthorized", strings.TrimSpace(rr.Body.String()))
				} else {
					require.NotEqual(t, http.StatusUnauthorized, rr.Code)
				}
			})
		}
	}
}
