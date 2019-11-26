package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func toJSON(t *testing.T, r interface{}) string {
	b, err := json.Marshal(r)
	require.NoError(t, err)
	return string(b)
}

func TestVerifyAddress(t *testing.T) {
	cases := []struct {
		name         string
		method       string
		status       int
		contentType  string
		csrfDisabled bool
		httpBody     string
		httpResponse HTTPResponse
	}{
		{
			name:         "405",
			method:       http.MethodGet,
			status:       http.StatusMethodNotAllowed,
			httpResponse: NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},

		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  ContentTypeForm,
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:         "400 - Missing address",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpBody:     "{}",
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "address is required"),
		},

		{
			name:   "422 - Invalid checksum",
			method: http.MethodPost,
			status: http.StatusUnprocessableEntity,
			httpBody: toJSON(t, VerifyAddressRequest{
				Address: "7apQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "Invalid checksum"),
		},
		{
			name:   "200",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: toJSON(t, VerifyAddressRequest{
				Address: "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
			}),
			httpResponse: HTTPResponse{
				Data: VerifyAddressResponse{
					Version: 0,
				},
			},
		},
		{
			name:   "200 - csrf disabled",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: toJSON(t, VerifyAddressRequest{
				Address: "7cpQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
			}),
			httpResponse: HTTPResponse{
				Data: VerifyAddressResponse{
					Version: 0,
				},
			},
			csrfDisabled: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/address/verify"
			gateway := &MockGatewayer{}

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = ContentTypeJSON
			}

			req.Header.Set("Content-Type", contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()
			cfg := defaultMuxConfig()
			cfg.disableCSRF = tc.csrfDisabled
			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.Unmarshal(rr.Body.Bytes(), &rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var addrRsp VerifyAddressResponse
				err := json.Unmarshal(rsp.Data, &addrRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(VerifyAddressResponse), addrRsp)
			}
		})
	}
}
