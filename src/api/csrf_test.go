package api

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
					handler := newServerMux(muxConfig{
						host:                 configuredHost,
						appLoc:               ".",
						enableJSON20RPC:      true,
						enableUnversionedAPI: true,
					}, gateway, csrfStore, nil)

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
				handler := newServerMux(muxConfig{
					host:            configuredHost,
					appLoc:          ".",
					enableJSON20RPC: true,
				}, gateway, csrfStore, nil)

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
			handler := newServerMux(muxConfig{
				host:            configuredHost,
				appLoc:          ".",
				enableJSON20RPC: true,
			}, gateway, csrfStore, nil)

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

		endpoint := "/api/v1/wallet/update"

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
		handler := newServerMux(muxConfig{
			host:            configuredHost,
			appLoc:          ".",
			enableJSON20RPC: true,
		}, gateway, csrfStore, nil)

		handler.ServeHTTP(rr, req)

		return rr
	}

	// First request to POST /wallet/update is rejected because of missing CSRF
	rr := updateWalletLabel("")
	require.Equal(t, http.StatusForbidden, rr.Code)
	require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())

	// Make a request to /csrf to get a token
	gateway := &GatewayerMock{}
	handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore, nil)

	// non-GET request to /csrf is invalid
	req, err := http.NewRequest(http.MethodPost, "/api/v1/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusMethodNotAllowed, rr.Code)
	require.Nil(t, csrfStore.token, "csrfStore.token should not be set yet")

	// CSRF disabled 404s
	csrfStore.Enabled = false

	req, err = http.NewRequest(http.MethodGet, "/api/v1/csrf", nil)
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Nil(t, csrfStore.token, "csrfStore.token should not be set yet")

	csrfStore.Enabled = true

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

	// Make another call to /csrf, this will invalidate the first token
	// Request a CSRF token, use it in a request
	req, err = http.NewRequest(http.MethodGet, "/api/v1/csrf", nil)
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
