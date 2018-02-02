package gui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

// Uses /wallet/newAddress to test CSRF wrapper

// Scenarios
// Methods: POST, GET, PUT, DELETE
// CSRF State: Enabled, Disabled
// Token: Valid, Invalid(Expired, Wrong), Empty

const (
	TokenValid = iota
	TokenInvalid
	TokenExpired
	TokenEmpty
)

func TestCSRFWrapper(t *testing.T) {
	type httpBody struct {
		ID  string
		Num string
	}
	type Addresses struct {
		Address []string `json:"addresses"`
	}
	var responseAddresses = Addresses{}
	var addrs = make([]cipher.Address, 3)
	var csrfStore = &CSRFStore{}

	for i := 0; i < 3; i++ {
		pub, _ := cipher.GenerateDeterministicKeyPair(cipher.RandByte(32))

		addrs[i] = cipher.AddressFromPubKey(pub)
		responseAddresses.Address = append(responseAddresses.Address, addrs[i].String())
	}
	tt := []struct {
		name                      string
		method                    string
		body                      *httpBody
		status                    int
		err                       string
		walletID                  string
		n                         uint64
		gatewayNewAddressesResult []cipher.Address
		gatewayNewAddressesErr    error
		responseCode              int
		csrfDisabled              bool
		csrfTokenType             int
	}{
		{
			name:   "200 - OK - Valid CSRF Token",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
			csrfTokenType:             TokenValid,
		},
		{
			name:     "403 - Forbidden - Invalid CSRF Token",
			method:   http.MethodPost,
			body:     &httpBody{},
			status:   http.StatusForbidden,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
			csrfTokenType:             TokenInvalid,
		},
		{
			name:     "403 - Forbidden - Expired CSRF Token",
			method:   http.MethodPost,
			body:     &httpBody{},
			status:   http.StatusForbidden,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
			csrfTokenType:             TokenExpired,
		},
		{
			name:     "403 - Forbidden - Empty CSRF Token",
			method:   http.MethodPost,
			body:     &httpBody{},
			status:   http.StatusForbidden,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
			csrfTokenType:             TokenEmpty,
		},
		{
			name:     "405 - GET Method",
			method:   http.MethodGet,
			body:     &httpBody{},
			status:   http.StatusMethodNotAllowed,
			walletID: "",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
		},
		{
			name:     "405 - PUT Method",
			method:   http.MethodPut,
			body:     &httpBody{},
			status:   http.StatusMethodNotAllowed,
			walletID: "",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
		},
		{
			name:     "405 - DELETE Method",
			method:   http.MethodDelete,
			body:     &httpBody{},
			status:   http.StatusMethodNotAllowed,
			walletID: "",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              false,
		},
		{
			name:   "200 - OK - CSRF Disabled",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			csrfDisabled:              true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("NewAddresses", tc.walletID, tc.n).Return(tc.gatewayNewAddressesResult, tc.gatewayNewAddressesErr)

			endpoint := "/wallet/newAddress"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.ID != "" {
					v.Add("id", tc.body.ID)
				}
				if tc.body.Num != "" {
					v.Add("num", tc.body.Num)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			(*csrfStore).Enabled = !tc.csrfDisabled
			if csrfStore.Enabled {
				csrfStore, req = setCSRFParameters(csrfStore, tc.csrfTokenType, req)
			}
			rr := httptest.NewRecorder()
			handler := NewServerMux(configuredHost, ".", gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)
		})
	}

}

func setCSRFParameters(csrfStore *CSRFStore, tokenType int, req *http.Request) (*CSRFStore, *http.Request) {
	csrfStore.setToken(generateToken())
	// token check
	switch tokenType {
	case TokenValid:
		req.Header.Add("X-CSRF-Token", csrfStore.getTokenValue())
	case TokenInvalid:
		// set invalid token value
		req.Header.Add("X-CSRF-Token", "xcasadsadsa")
	case TokenExpired:
		req.Header.Add("X-CSRF-Token", csrfStore.getTokenValue())
		// set some old unix time
		csrfStore.token.ExpiresAt = time.Unix(1517509381, 10)
	case TokenEmpty:
		// set empty token
		req.Header.Add("X-CSRF-Token", "")
	}

	return csrfStore, req
}
