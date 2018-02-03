package gui

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"/lastTxs",
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
					handler := NewServerMux(configuredHost, ".", gateway, csrfStore)

					handler.ServeHTTP(rr, req)

					status := rr.Code
					require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)
					require.Equal(t, "403 Forbidden - invalid CSRF token\n", rr.Body.String())
				})
			}
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
			handler := NewServerMux(configuredHost, ".", gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, http.StatusForbidden, status, "wrong status code: got `%v` want `%v`", status, http.StatusForbidden)
			require.Equal(t, "403 Forbidden\n", rr.Body.String())
		})
	}
}
