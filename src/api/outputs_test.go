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
		getUnspentOutputsResponse *readable.OutputSet
		getUnspentOutputsError    error
		httpResponse              *readable.OutputSet
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
			err:                       "500 Internal Server Error - get unspent outputs failed: getUnspentOutputsError",
			getUnspentOutputsResponse: nil,
			getUnspentOutputsError:    errors.New("getUnspentOutputsError"),
		},
		{
			name:                      "200 - OK",
			method:                    http.MethodGet,
			status:                    http.StatusOK,
			getUnspentOutputsResponse: &readable.OutputSet{},
			httpResponse:              &readable.OutputSet{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
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
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *readable.OutputSet
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
