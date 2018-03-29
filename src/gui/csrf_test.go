package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	tokenValid        = "token_valid"
	tokenInvalid      = "token_invalid"
	tokenExpired      = "token_expired"
	tokenEmpty        = "token_empty"
	tokenNotGenerated = "token_not_generated"
)

func setCSRFParameters(csrfStore *CSRFStore, tokenType string, req *http.Request) {
	token := newCSRFToken()
	// token check
	switch tokenType {
	case tokenValid:
		csrfStore.setToken(token)
		req.Header.Set("X-CSRF-Token", token.String())
	case tokenInvalid:
		// set invalid token value
		csrfStore.setToken(token)
		req.Header.Set("X-CSRF-Token", "xcasadsadsa")
	case tokenExpired:
		csrfStore.setToken(token)
		csrfStore.token.ExpiresAt = time.Unix(1517509381, 10)
		req.Header.Set("X-CSRF-Token", token.String())
		// set some old unix time
	case tokenEmpty:
		// set empty token
		csrfStore.setToken(token)
		req.Header.Set("X-CSRF-Token", "")
	case tokenNotGenerated:
		// don't set token
		csrfStore.token = nil
		req.Header.Set("X-CSRF-Token", token.String())
	}
}

var endpoints = []string{
	"/version",
	"/outputs",
	"/balance",
	"/wallet",
	"/wallet/create",
	"/wallet/newAddress",
	"/wallet/balance",
	"/wallet/spend",
	"/wallet/transactions",
	"/wallet/update",
	"/wallets",
	"/wallets/folderName",
	"/wallet/newSeed",
	"/wallet/unload",
	"/blockchain/metadata",
	"/blockchain/progress",
	"/block",
	"/blocks",
	"/last_blocks",
	"/network/connection",
	"/network/connections",
	"/network/defaultConnections",
	"/network/connections/trust",
	"/network/connections/exchange",
	"/pendingTxs",
	"/transaction",
	"/transactions",
	"/injectTransaction",
	"/resendUnconfirmedTxns",
	"/rawtx",
	"/uxout",
	"/address_uxouts",
	"/explorer/address",
	"/coinSupply",
	"/richlist",
	"/addresscount",
}

func TestCSRFWrapper(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete}
	cases := []string{tokenInvalid, tokenExpired, tokenEmpty, tokenNotGenerated}

	for _, endpoint := range endpoints {
		for _, method := range methods {
			for _, c := range cases {
				name := fmt.Sprintf("%s %s %s", method, endpoint, c)
				t.Run(name, func(t *testing.T) {
					gateway := &GatewayerMock{}

					req, err := http.NewRequest(method, endpoint, nil)
					require.NoError(t, err)

					csrfStore := &CSRFStore{
						Enabled: true,
					}
					setCSRFParameters(csrfStore, c, req)

					rr := httptest.NewRecorder()
					handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

					handler.ServeHTTP(rr, req)

					status := rr.Code
					require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)
					require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())
				})
			}
		}
	}
}

func TestOriginRefererCheck(t *testing.T) {
	cases := []struct {
		name    string
		origin  string
		referer string
	}{
		{
			name:   "mismatched origin header",
			origin: "http://example.com/",
		},
		{
			name:    "mismatched referer header",
			referer: "http://example.com/",
		},
	}

	for _, endpoint := range endpoints {
		for _, tc := range cases {
			name := fmt.Sprintf("%s %s", tc.name, endpoint)
			t.Run(name, func(t *testing.T) {
				gateway := &GatewayerMock{}

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
				handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

				handler.ServeHTTP(rr, req)

				status := rr.Code
				require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)
				require.Equal(t, "403 Forbidden\n", rr.Body.String())
			})
		}
	}
}

func TestHostCheck(t *testing.T) {
	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			gateway := &GatewayerMock{}

			req, err := http.NewRequest(http.MethodGet, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			req.Host = "example.com"

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)
			require.Equal(t, "403 Forbidden\n", rr.Body.String())
		})
	}
}

func TestCSRF(t *testing.T) {
	csrfStore := &CSRFStore{
		Enabled: true,
	}

	updateWalletLabel := func(csrfToken string) *httptest.ResponseRecorder {
		gateway := &GatewayerMock{}
		gateway.On("UpdateWalletLabel", "fooid", "foolabel").Return(nil)

		endpoint := "/wallet/update"

		v := url.Values{}
		v.Add("id", "fooid")
		v.Add("label", "foolabel")

		req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(v.Encode()))
		require.NoError(t, err)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		if csrfToken != "" {
			req.Header.Set("X-CSRF-Token", csrfToken)
		}

		rr := httptest.NewRecorder()
		handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

		handler.ServeHTTP(rr, req)

		return rr
	}

	// First request to POST /wallet/update is rejected because of missing CSRF
	rr := updateWalletLabel("")
	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())

	// Make a request to /csrf to get a token
	gateway := &GatewayerMock{}
	handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

	// non-GET request to /csrf is invalid
	req, err := http.NewRequest(http.MethodPost, "/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Nil(t, csrfStore.token, "csrfStore.token should not be set yet")

	// CSRF disabled 404s
	csrfStore.Enabled = false

	req, err = http.NewRequest(http.MethodGet, "/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Nil(t, csrfStore.token, "csrfStore.token should not be set yet")

	csrfStore.Enabled = true

	// Request a CSRF token, use it in a request
	req, err = http.NewRequest(http.MethodGet, "/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var msg map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &msg)
	require.NoError(t, err)

	token := msg["csrf_token"]
	require.NotEmpty(t, token)

	req, err = http.NewRequest(http.MethodPost, "/version", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Make a request to POST /wallet/update again, using the CSRF token
	rr = updateWalletLabel(token)
	require.Equal(t, http.StatusOK, rr.Code)

	// Make another call to /csrf, this will invalidate the first token
	// Request a CSRF token, use it in a request
	req, err = http.NewRequest(http.MethodGet, "/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var msg2 map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &msg2)
	require.NoError(t, err)
	require.NotEmpty(t, msg2["csrf_token"])
	require.NotEqual(t, msg["csrf_token"], msg2["csrf_token"])

	rr = updateWalletLabel(token)
	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())
}
