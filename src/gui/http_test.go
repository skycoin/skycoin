package gui

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestGetOutputsHandler(t *testing.T) {
	validAddr := "2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf"
	invalidAddr := "invalidAddr"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"

	type httpBody struct {
		addrs   string
		hashStr string
	}
	tt := []struct {
		name                      string
		method                    string
		url                       string
		status                    int
		err                       string
		httpBody                  *httpBody
		uxid                      string
		getUnspentOutputsResponse *visor.ReadableOutputSet
		getUnspentOutputsError    error
		httpResponse              *visor.ReadableOutputSet
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - addrs and hashes together",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - addrs and hashes cannot be specified together",
			httpBody: &httpBody{
				addrs:   validAddr,
				hashStr: validHash,
			},
		},
		{
			name:   "400 - invalid address",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - addrs contains invalid address",
			httpBody: &httpBody{
				addrs: invalidAddr,
			},
		},
		{
			name:   "500 - getUnspentOutputsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			getUnspentOutputsResponse: nil,
			getUnspentOutputsError:    errors.New("getUnspentOutputsError"),
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			status: http.StatusOK,
			getUnspentOutputsResponse: &visor.ReadableOutputSet{},
			httpResponse:              &visor.ReadableOutputSet{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			endpoint := "/outputs"
			gateway.On("GetUnspentOutputs", mock.Anything).Return(tc.getUnspentOutputsResponse, tc.getUnspentOutputsError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.hashStr != "" {
					v.Add("hashes", tc.httpBody.hashStr)
				}
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{})
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *visor.ReadableOutputSet
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestGetBalanceHandler(t *testing.T) {
	type httpBody struct {
		addrs string
	}
	invalidAddr := "invalidAddr"
	validAddr := "2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf"
	address, err := cipher.DecodeBase58Address(validAddr)
	require.NoError(t, err)
	tt := []struct {
		name                      string
		method                    string
		url                       string
		status                    int
		err                       string
		httpBody                  *httpBody
		uxid                      string
		getBalanceOfAddrsArg      []cipher.Address
		getBalanceOfAddrsResponse []wallet.BalancePair
		getBalanceOfAddrsError    error
		httpResponse              wallet.BalancePair
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - invalid address",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - address invalidAddr is invalid: Invalid base58 character",
			httpBody: &httpBody{
				addrs: invalidAddr,
			},
		},
		{
			name:   "500 - GetBalanceOfAddrsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - Get balance failed: GetBalanceOfAddrsError",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg:   []cipher.Address{address},
			getBalanceOfAddrsError: errors.New("GetBalanceOfAddrsError"),
		},
		{
			name:   "500 - balance Confirmed coins uint64 addition overflow",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - uint64 addition overflow",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
			},
		},
		{
			name:   "500 - balance Predicted coins uint64 addition overflow",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - uint64 addition overflow",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
				},
			},
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "200 - OK",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
			},
			httpResponse: wallet.BalancePair{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			endpoint := "/balance"
			gateway.On("GetBalanceOfAddrs", tc.getBalanceOfAddrsArg).Return(tc.getBalanceOfAddrsResponse, tc.getBalanceOfAddrsError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{})
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg wallet.BalancePair
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
