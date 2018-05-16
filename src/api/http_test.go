package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
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
			err:    "500 Internal Server Error - get unspent outputs failed: getUnspentOutputsError",
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
			endpoint := "/api/v1/outputs"
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
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
			name:     "400 - no addresses",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - addrs is required",
			httpBody: &httpBody{},
		},
		{
			name:   "500 - GetBalanceOfAddrsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.GetBalanceOfAddrs failed: GetBalanceOfAddrsError",
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
			endpoint := "/api/v1/balance"
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
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

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: tc.appLoc}, nil, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			c := Config{
				EnableGUI:   tc.enableGUI,
				DisableCSRF: true,
				StaticDir:   tc.appLoc,
			}

			host := "127.0.0.1:6423"
			s, err := Create(host, c, nil)
			require.NoError(t, err)

			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.Serve()
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
