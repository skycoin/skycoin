package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/visor"
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
			method: http.MethodDelete,
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
			err:    "400 Bad Request - address \"invalidAddr\" is invalid: Invalid base58 character",
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
			name:   "200 - OK",
			method: http.MethodGet,
			status: http.StatusOK,
			getUnspentOutputsResponse: &visor.UnspentOutputsSummary{
				HeadBlock: &coin.SignedBlock{},
			},
			httpResponse: &readable.UnspentOutputsSummary{
				Head: readable.BlockHeader{
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
				HeadOutputs:     readable.UnspentOutputs{},
				OutgoingOutputs: readable.UnspentOutputs{},
				IncomingOutputs: readable.UnspentOutputs{},
			},
		},
		{
			name:   "200 - OK POST",
			method: http.MethodPost,
			status: http.StatusOK,
			getUnspentOutputsResponse: &visor.UnspentOutputsSummary{
				HeadBlock: &coin.SignedBlock{},
			},
			httpResponse: &readable.UnspentOutputsSummary{
				Head: readable.BlockHeader{
					Hash:         "7b8ec8dd836b564f0c85ad088fc744de820345204e154bc1503e04e9d6fdd9f1",
					PreviousHash: "0000000000000000000000000000000000000000000000000000000000000000",
					BodyHash:     "0000000000000000000000000000000000000000000000000000000000000000",
					UxHash:       "0000000000000000000000000000000000000000000000000000000000000000",
				},
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

			var reqBody io.Reader
			if len(v) > 0 {
				if tc.method == http.MethodPost {
					reqBody = strings.NewReader(v.Encode())
				} else {
					endpoint += "?" + v.Encode()
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, reqBody)
			require.NoError(t, err)

			if tc.method == http.MethodPost {
				req.Header.Set("Content-Type", ContentTypeForm)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				var msg *readable.UnspentOutputsSummary
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
