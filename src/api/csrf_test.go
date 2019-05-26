package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"encoding/json"
	"sync"

	"github.com/stretchr/testify/require"
)

const (
	tokenValid            = "token_valid"
	tokenInvalid          = "token_invalid"
	tokenInvalidSignature = "token_invalid_signature"
	tokenExpired          = "token_expired"
	tokenEmpty            = "token_empty"
)

func setCSRFParameters(t *testing.T, tokenType string, req *http.Request) {
	token, err := newCSRFToken()
	require.NoError(t, err)
	// token check
	switch tokenType {
	case tokenValid:
		req.Header.Set("X-CSRF-Token", token)
	case tokenInvalid:
		// add invalid token value
		req.Header.Set("X-CSRF-Token", "xcasadsadsa")
	case tokenInvalidSignature:
		req.Header.Set("X-CSRF-Token", "YXNkc2Fkcw.YXNkc2Fkcw")
	case tokenExpired:
		// set some old unix time
		expiredToken, err := newCSRFTokenWithTime(time.Unix(1517509381, 10))
		require.NoError(t, err)
		req.Header.Set("X-CSRF-Token", expiredToken)
	case tokenEmpty:
		// add empty token
		req.Header.Set("X-CSRF-Token", "")
	}
}

func TestCSRFWrapper(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	cases := []string{tokenInvalid, tokenExpired, tokenEmpty, tokenInvalidSignature}

	for endpoint := range endpointsMethods {
		for _, method := range methods {
			for _, c := range cases {
				name := fmt.Sprintf("%s %s %s", method, endpoint, c)
				t.Run(name, func(t *testing.T) {
					gateway := &MockGatewayer{}

					req, err := http.NewRequest(method, endpoint, nil)
					require.NoError(t, err)

					setCSRFParameters(t, c, req)

					isAPIV2 := strings.HasPrefix(endpoint, "/api/v2")
					if isAPIV2 {
						req.Header.Set("Content-Type", ContentTypeJSON)
					}

					rr := httptest.NewRecorder()
					handler := newServerMux(muxConfig{
						host:           configuredHost,
						appLoc:         ".",
						disableCSRF:    false,
						disableCSP:     true,
						enabledAPISets: allAPISetsEnabled,
					}, gateway)

					handler.ServeHTTP(rr, req)

					status := rr.Code
					require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)

					var errMsg error
					switch c {
					case tokenInvalid, tokenEmpty:
						errMsg = ErrCSRFInvalid
					case tokenInvalidSignature:
						errMsg = ErrCSRFInvalidSignature
					case tokenExpired:
						errMsg = ErrCSRFExpired
					}

					if isAPIV2 {
						require.Equal(t, fmt.Sprintf("{\n    \"error\": {\n        \"message\": \"%s\",\n        \"code\": 403\n    }\n}", errMsg), rr.Body.String())
					} else {
						require.Equal(t, fmt.Sprintf("403 Forbidden - %s\n", errMsg), rr.Body.String())
					}
				})
			}
		}
	}
}

func TestCSRFWrapperConcurrent(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	cases := []string{tokenInvalid, tokenExpired, tokenEmpty, tokenInvalidSignature}

	gateway := &MockGatewayer{}

	handler := newServerMux(muxConfig{
		host:           configuredHost,
		appLoc:         ".",
		disableCSRF:    false,
		disableCSP:     true,
		enabledAPISets: allAPISetsEnabled,
	}, gateway)

	var wg sync.WaitGroup

	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for endpoint := range endpointsMethods {
				for _, method := range methods {
					for _, c := range cases {
						name := fmt.Sprintf("%s %s %s", method, endpoint, c)
						t.Run(name, func(t *testing.T) {

							req, err := http.NewRequest(method, endpoint, nil)
							require.NoError(t, err)

							setCSRFParameters(t, c, req)

							isAPIV2 := strings.HasPrefix(endpoint, "/api/v2")
							if isAPIV2 {
								req.Header.Set("Content-Type", ContentTypeJSON)
							}

							rr := httptest.NewRecorder()

							handler.ServeHTTP(rr, req)

							status := rr.Code
							require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)

							var errMsg error
							switch c {
							case tokenInvalid, tokenEmpty:
								errMsg = ErrCSRFInvalid
							case tokenInvalidSignature:
								errMsg = ErrCSRFInvalidSignature
							case tokenExpired:
								errMsg = ErrCSRFExpired
							}

							if isAPIV2 {
								require.Equal(t, fmt.Sprintf("{\n    \"error\": {\n        \"message\": \"%s\",\n        \"code\": 403\n    }\n}", errMsg), rr.Body.String())
							} else {
								require.Equal(t, fmt.Sprintf("403 Forbidden - %s\n", errMsg), rr.Body.String())
							}
						})
					}
				}
			}
		}()
	}
	wg.Wait()

}

func TestCSRF(t *testing.T) {
	updateWalletLabel := func(csrfToken string) *httptest.ResponseRecorder {
		gateway := &MockGatewayer{}
		gateway.On("UpdateWalletLabel", "fooid", "foolabel").Return(nil)

		endpoint := "/api/v1/wallet/update"

		v := url.Values{}
		v.Add("id", "fooid")
		v.Add("label", "foolabel")

		req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(v.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", ContentTypeForm)

		if csrfToken != "" {
			req.Header.Set("X-CSRF-Token", csrfToken)
		}

		rr := httptest.NewRecorder()
		handler := newServerMux(muxConfig{
			host:           configuredHost,
			appLoc:         ".",
			disableCSRF:    false,
			disableCSP:     true,
			enabledAPISets: allAPISetsEnabled,
		}, gateway)

		handler.ServeHTTP(rr, req)

		return rr
	}

	// First request to POST /wallet/update is rejected because of missing CSRF
	rr := updateWalletLabel("")
	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())

	// Make a request to /csrf to get a token
	gateway := &MockGatewayer{}
	cfg := defaultMuxConfig()
	cfg.disableCSRF = false
	handler := newServerMux(cfg, gateway)

	// non-GET request to /csrf is invalid
	req, err := http.NewRequest(http.MethodPost, "/api/v1/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// CSRF disabled 404s
	cfg.disableCSRF = true
	handler = newServerMux(cfg, gateway)

	req, err = http.NewRequest(http.MethodGet, "/api/v1/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	cfg.disableCSRF = false
	handler = newServerMux(cfg, gateway)

	// Request a CSRF token, use it in a request
	req, err = http.NewRequest(http.MethodGet, "/api/v1/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var msg map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &msg)
	require.NoError(t, err)

	token := msg["csrf_token"]
	require.NotEmpty(t, token)

	req, err = http.NewRequest(http.MethodPost, "/api/v1/version", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Make a request to POST /wallet/update again, using the CSRF token
	rr = updateWalletLabel(token)
	require.Equal(t, http.StatusOK, rr.Code)
}
