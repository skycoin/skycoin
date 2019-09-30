package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/daemon"
	"github.com/SkycoinProject/skycoin/src/daemon/pex"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/util/useragent"
)

func TestConnection(t *testing.T) {
	tt := []struct {
		name                       string
		method                     string
		status                     int
		err                        string
		addr                       string
		gatewayGetConnectionResult *daemon.Connection
		gatewayGetConnectionError  error
		result                     *readable.Connection
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                       "400 - empty addr",
			method:                     http.MethodGet,
			status:                     http.StatusBadRequest,
			err:                        "400 Bad Request - addr is required",
			addr:                       "",
			gatewayGetConnectionResult: nil,
			result:                     nil,
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "",
			addr:   "addr",
			gatewayGetConnectionResult: &daemon.Connection{
				Addr: "127.0.0.1:6061",
				Gnet: daemon.GnetConnectionDetails{
					ID:           1,
					LastSent:     time.Unix(99999, 0),
					LastReceived: time.Unix(1111111, 0),
				},
				ConnectionDetails: daemon.ConnectionDetails{
					Outgoing:    true,
					ConnectedAt: time.Unix(222222, 0),
					State:       daemon.ConnectionStateIntroduced,
					Mirror:      6789,
					ListenPort:  9877,
					Height:      1234,
					UserAgent:   useragent.MustParse("skycoin:0.25.1(foo)"),
				},
				Pex: pex.Peer{
					Trusted: false,
				},
			},
			result: &readable.Connection{
				Addr:          "127.0.0.1:6061",
				GnetID:        1,
				LastSent:      99999,
				LastReceived:  1111111,
				ConnectedAt:   222222,
				Outgoing:      true,
				State:         daemon.ConnectionStateIntroduced,
				Mirror:        6789,
				ListenPort:    9877,
				Height:        1234,
				UserAgent:     useragent.MustParse("skycoin:0.25.1(foo)"),
				IsTrustedPeer: false,
			},
		},

		{
			name:                      "500 - GetConnection failed",
			method:                    http.MethodGet,
			status:                    http.StatusInternalServerError,
			err:                       "500 Internal Server Error - GetConnection failed",
			addr:                      "addr",
			gatewayGetConnectionError: errors.New("GetConnection failed"),
		},

		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			addr:   "addr",
			err:    "404 Not Found",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connection"
			gateway := &MockGatewayer{}
			gateway.On("GetConnection", tc.addr).Return(tc.gatewayGetConnectionResult, tc.gatewayGetConnectionError)

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
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *readable.Connection
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestConnections(t *testing.T) {
	intrOut := daemon.Connection{
		Addr: "127.0.0.1:6061",
		Gnet: daemon.GnetConnectionDetails{
			ID:           1,
			LastSent:     time.Unix(99999, 0),
			LastReceived: time.Unix(1111111, 0),
		},
		ConnectionDetails: daemon.ConnectionDetails{
			Outgoing:    true,
			State:       daemon.ConnectionStateIntroduced,
			ConnectedAt: time.Unix(222222, 0),
			Mirror:      9876,
			ListenPort:  9877,
			Height:      1234,
			UserAgent:   useragent.MustParse("skycoin:0.25.1(foo)"),
		},
		Pex: pex.Peer{
			Trusted: true,
		},
	}

	intrIn := daemon.Connection{
		Addr: "127.0.0.2:6062",
		Gnet: daemon.GnetConnectionDetails{
			ID:           2,
			LastSent:     time.Unix(99999, 0),
			LastReceived: time.Unix(1111111, 0),
		},
		ConnectionDetails: daemon.ConnectionDetails{
			Outgoing:    false,
			State:       daemon.ConnectionStateIntroduced,
			ConnectedAt: time.Unix(222222, 0),
			Mirror:      9877,
			ListenPort:  9879,
			Height:      1234,
			UserAgent:   useragent.MustParse("skycoin:0.25.1(foo)"),
		},
	}

	readIntrOut := readable.Connection{
		Addr:          "127.0.0.1:6061",
		GnetID:        1,
		LastSent:      99999,
		LastReceived:  1111111,
		ConnectedAt:   222222,
		Outgoing:      true,
		State:         daemon.ConnectionStateIntroduced,
		Mirror:        9876,
		ListenPort:    9877,
		Height:        1234,
		UserAgent:     useragent.MustParse("skycoin:0.25.1(foo)"),
		IsTrustedPeer: true,
	}

	readIntrIn := readable.Connection{
		Addr:          "127.0.0.2:6062",
		GnetID:        2,
		LastSent:      99999,
		LastReceived:  1111111,
		ConnectedAt:   222222,
		Outgoing:      false,
		State:         daemon.ConnectionStateIntroduced,
		Mirror:        9877,
		ListenPort:    9879,
		Height:        1234,
		UserAgent:     useragent.MustParse("skycoin:0.25.1(foo)"),
		IsTrustedPeer: false,
	}

	conns := []daemon.Connection{intrOut, intrIn}
	readConns := []readable.Connection{readIntrOut, readIntrIn}

	tt := []struct {
		name                                 string
		method                               string
		status                               int
		states                               string
		direction                            string
		err                                  string
		gatewayGetSolicitedConnectionsResult []daemon.Connection
		gatewayGetSolicitedConnectionsError  error
		result                               Connections
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:                                 "200 defaults",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:                                 "200 incoming",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			direction:                            "incoming",
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:                                 "200 outgoing",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			direction:                            "outgoing",
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:                                 "200 pending,connected",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			states:                               "pending,connected",
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:                                 "200 pending,connected outgoing",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			states:                               "pending,connected",
			direction:                            "outgoing",
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:                                 "200 pending,introduced,connected",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			states:                               "pending,introduced,connected",
			err:                                  "",
			gatewayGetSolicitedConnectionsResult: conns,
			result: Connections{
				Connections: readConns,
			},
		},

		{
			name:   "400 - bad state",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			states: "pending,foo",
			err:    "400 Bad Request - Invalid state in states. Valid states are \"pending\", \"connected\" or \"introduced\"",
		},

		{
			name:      "400 - bad direction",
			method:    http.MethodGet,
			status:    http.StatusBadRequest,
			direction: "foo",
			err:       "400 Bad Request - Invalid direction. Valid directions are \"outgoing\" or \"incoming\"",
		},

		{
			name:                                "500 - GetConnections failed",
			method:                              http.MethodGet,
			status:                              http.StatusInternalServerError,
			err:                                 "500 Internal Server Error - GetConnections failed",
			gatewayGetSolicitedConnectionsError: errors.New("GetConnections failed"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections"
			gateway := &MockGatewayer{}
			gateway.On("GetConnections", mock.Anything).Return(tc.gatewayGetSolicitedConnectionsResult, tc.gatewayGetSolicitedConnectionsError)

			v := url.Values{}
			if tc.states != "" {
				v.Add("states", tc.states)
			}
			if tc.direction != "" {
				v.Add("direction", tc.direction)
			}

			endpoint += "?" + v.Encode()

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Connections
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
			name:                               "200",
			method:                             http.MethodGet,
			status:                             http.StatusOK,
			err:                                "",
			gatewayGetDefaultConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                             []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/defaultConnections"
			gateway := &MockGatewayer{}
			gateway.On("GetDefaultConnections").Return(tc.gatewayGetDefaultConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
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
			name:                             "200",
			method:                           http.MethodGet,
			status:                           http.StatusOK,
			err:                              "",
			gatewayGetTrustConnectionsResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                           []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/trust"
			gateway := &MockGatewayer{}
			gateway.On("GetTrustConnections").Return(tc.gatewayGetTrustConnectionsResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
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
			name:                            "200",
			method:                          http.MethodGet,
			status:                          http.StatusOK,
			err:                             "",
			gatewayGetExchgConnectionResult: []string{"44.33.22.11", "11.44.66.88"},
			result:                          []string{"11.44.66.88", "44.33.22.11"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/network/connections/exchange"
			gateway := &MockGatewayer{}
			gateway.On("GetExchgConnection").Return(tc.gatewayGetExchgConnectionResult)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []string
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestDisconnect(t *testing.T) {
	tt := []struct {
		name          string
		method        string
		status        int
		err           string
		disconnectErr error
		id            string
		gnetID        uint64
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:   "400 missing ID",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - id is required",
		},

		{
			name:   "400 invalid ID 0",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid id",
			id:     "0",
		},

		{
			name:   "400 invalid ID negative",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid id",
			id:     "-100",
		},

		{
			name:          "404 Disconnect connection not found",
			method:        http.MethodPost,
			status:        http.StatusNotFound,
			err:           "404 Not Found",
			disconnectErr: daemon.ErrConnectionNotExist,
			id:            "100",
			gnetID:        100,
		},

		{
			name:          "500 Disconnect error",
			method:        http.MethodPost,
			status:        http.StatusInternalServerError,
			err:           "500 Internal Server Error - foo",
			disconnectErr: errors.New("foo"),
			id:            "100",
			gnetID:        100,
		},

		{
			name:   "200",
			method: http.MethodPost,
			status: http.StatusOK,
			id:     "100",
			gnetID: 100,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("DisconnectByGnetID", tc.gnetID).Return(tc.disconnectErr)

			endpoint := "/api/v1/network/connection/disconnect"
			v := url.Values{}
			if tc.id != "" {
				v.Add("id", tc.id)
			}

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var obj struct{}
				err = json.Unmarshal(rr.Body.Bytes(), &obj)
				require.NoError(t, err)
			}
		})
	}
}
