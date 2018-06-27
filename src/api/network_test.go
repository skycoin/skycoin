package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/daemon"
)

func TestConnection(t *testing.T) {
	bp := daemon.BlockchainProgress{
		Current: 35,
		Highest: 39,
		Peers:   nil}
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.3.5.1",
		Height:  39,
	})
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.0.0.1",
		Height:  12,
	})
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.0.5.1",
		Height:  13,
	})

	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		addr                               string
		gatewayGetConnectionResult         *daemon.Connection
		gatewayGetBlockchainProgressResult *daemon.BlockchainProgress
		gatewayGetBlockchainProgressError  error
		result                             *daemon.Connection
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - empty addr",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - addr is required",
			addr:   "",
			gatewayGetConnectionResult: nil,
			result: nil,
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			addr:   "addr",
			gatewayGetBlockchainProgressResult: &bp,
			gatewayGetBlockchainProgressError:  nil,
			gatewayGetConnectionResult: &daemon.Connection{
				ID:           1,
				Addr:         "127.0.0.1",
				LastSent:     99999,
				LastReceived: 1111111,
				Outgoing:     true,
				Introduced:   true,
				Mirror:       9876,
				ListenPort:   9877,
			},
			result: &daemon.Connection{
				ID:           1,
				Addr:         "127.0.0.1",
				LastSent:     99999,
				LastReceived: 1111111,
				Outgoing:     true,
				Introduced:   true,
				Mirror:       9876,
				ListenPort:   9877,
			},
		},
		{
			name:   "500 - blockchain progress failed",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - some error",
			addr:   "addr",
			gatewayGetBlockchainProgressResult: nil,
			gatewayGetBlockchainProgressError:  errors.New("some error"),
			gatewayGetConnectionResult: &daemon.Connection{
				ID:           1,
				Addr:         "127.0.0.1",
				LastSent:     99999,
				LastReceived: 1111111,
				Outgoing:     true,
				Introduced:   true,
				Mirror:       9876,
				ListenPort:   9877,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connection"
			gateway := NewGatewayerMock()
			gateway.On("GetConnection", tc.addr).Return(tc.gatewayGetConnectionResult)
			gateway.On("GetBlockchainProgress").Return(
				tc.gatewayGetBlockchainProgressResult,
				tc.gatewayGetBlockchainProgressError,
			)
			v := url.Values{}
			if tc.addr != "" {
				v.Add("addr", tc.addr)
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
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *daemon.Connection
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestConnections(t *testing.T) {
	bp := daemon.BlockchainProgress{
		Current: 35,
		Highest: 39,
		Peers:   nil}
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.3.5.1",
		Height:  39,
	})
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.0.0.1",
		Height:  12,
	})
	bp.Peers = append(bp.Peers, struct {
		Address string `json:"address"`
		Height  uint64 `json:"height"`
	}{
		Address: "127.0.5.1",
		Height:  13,
	})

	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		gatewayGetConnectionsResult        *daemon.Connections
		gatewayGetBlockchainProgressResult *daemon.BlockchainProgress
		gatewayGetBlockchainProgressError  error
		result                             *daemon.Connections
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetBlockchainProgressResult: &bp,
			gatewayGetBlockchainProgressError:  nil,
			gatewayGetConnectionsResult: &daemon.Connections{
				Connections: []*daemon.Connection{
					&daemon.Connection{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     true,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
					},
				},
			},
			result: &daemon.Connections{
				Connections: []*daemon.Connection{
					&daemon.Connection{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     true,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
					},
				},
			},
		},
		{
			name:   "500 - blockchain progress failed",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - some error",
			gatewayGetBlockchainProgressResult: nil,
			gatewayGetBlockchainProgressError:  errors.New("some error"),
			gatewayGetConnectionsResult: &daemon.Connections{
				Connections: []*daemon.Connection{
					&daemon.Connection{
						ID:           1,
						Addr:         "127.0.0.1",
						LastSent:     99999,
						LastReceived: 1111111,
						Outgoing:     true,
						Introduced:   true,
						Mirror:       9876,
						ListenPort:   9877,
					},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections"
			gateway := NewGatewayerMock()
			gateway.On("GetConnections").Return(tc.gatewayGetConnectionsResult)
			gateway.On("GetBlockchainProgress").Return(
				tc.gatewayGetBlockchainProgressResult,
				tc.gatewayGetBlockchainProgressError,
			)
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *daemon.Connections
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestDefaultConnections(t *testing.T) {
	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		gatewayGetDefaultConnectionsResult []string
		result                             []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetDefaultConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result: []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/defaultConnections"
			gateway := NewGatewayerMock()
			gateway.On("GetDefaultConnections").Return(tc.gatewayGetDefaultConnectionsResult)
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetTrustConnections(t *testing.T) {
	tt := []struct {
		name                             string
		method                           string
		status                           int
		err                              string
		gatewayGetTrustConnectionsResult []string
		result                           []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetTrustConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result: []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/trust"
			gateway := NewGatewayerMock()
			gateway.On("GetTrustConnections").Return(tc.gatewayGetTrustConnectionsResult)
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetExchgConnection(t *testing.T) {
	tt := []struct {
		name                            string
		method                          string
		status                          int
		err                             string
		gatewayGetExchgConnectionResult []string
		result                          []string
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			gatewayGetExchgConnectionResult: []string{"44.33.22.11", "11.44.66.88"},
			result: []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/exchange"
			gateway := NewGatewayerMock()
			gateway.On("GetExchgConnection").Return(tc.gatewayGetExchgConnectionResult)
			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
