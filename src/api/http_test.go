package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/visor"
)

const configuredHost = "127.0.0.1:6420"

func defaultMuxConfig() muxConfig {
	return muxConfig{
		host:       configuredHost,
		appLoc:     ".",
		disableCSP: true,
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
	}{
		{
			name:            "enable CSP GET /",
			endpoint:        "/",
			enableCSP:       true,
			appLoc:          "../gui/static/dist",
			expectCSPHeader: "script-src 'self' 127.0.0.1",
		},
		{
			name:            "disable CSP GET /",
			endpoint:        "/",
			enableCSP:       false,
			appLoc:          "../gui/static/dist",
			expectCSPHeader: "",
		},
		{
			// Confirms that the /csrf api won't be affected by the csp setting
			name:            "enable CSP GET /csrf",
			endpoint:        "/api/v1/csrf",
			enableCSP:       true,
			appLoc:          "",
			expectCSPHeader: "",
		},
		{
			// Confirms that the /version api won't be affected by the csp setting
			name:            "enable CSP GET /version",
			endpoint:        "/api/v1/version",
			enableCSP:       true,
			appLoc:          "",
			expectCSPHeader: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.endpoint, nil)
			require.NoError(t, err)

			gateway := &MockGatewayer{}
			gateway.On("GetBuildInfo").Return(visor.BuildInfo{})

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{
				host:       configuredHost,
				appLoc:     tc.appLoc,
				enableGUI:  true,
				disableCSP: !tc.enableCSP,
			}, gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			csp := rr.Header().Get("Content-Security-Policy")
			require.Equal(t, tc.expectCSPHeader, csp)
		})
	}
}
