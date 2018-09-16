package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/visor"
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
		status                    int
		err                       string
		httpBody                  *httpBody
		getUnspentOutputsResponse *visor.UnspentOutputsSummary
		getUnspentOutputsError    error
		httpResponse              *readable.UnspentOutputsSummary
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
			name:                      "500 - getUnspentOutputsError",
			method:                    http.MethodGet,
			status:                    http.StatusInternalServerError,
			err:                       "500 Internal Server Error - gateway.GetUnspentOutputsSummary failed: getUnspentOutputsError",
			getUnspentOutputsResponse: nil,
			getUnspentOutputsError:    errors.New("getUnspentOutputsError"),
		},
		{
			name:                      "200 - OK",
			method:                    http.MethodGet,
			status:                    http.StatusOK,
			getUnspentOutputsResponse: &visor.UnspentOutputsSummary{},
			httpResponse: &readable.UnspentOutputsSummary{
				HeadOutputs:     readable.UnspentOutputs{},
				OutgoingOutputs: readable.UnspentOutputs{},
				IncomingOutputs: readable.UnspentOutputs{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			endpoint := "/api/v1/outputs"
			gateway.On("GetUnspentOutputsSummary", mock.Anything).Return(tc.getUnspentOutputsResponse, tc.getUnspentOutputsError)

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
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *readable.UnspentOutputsSummary
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
