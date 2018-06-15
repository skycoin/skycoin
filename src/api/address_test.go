package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyAddress(t *testing.T) {
	toJSON := func(r VerifyAddressRequest) string {
		b, err := json.Marshal(r)
		require.NoError(t, err)
		return string(b)
	}

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
			contentType:  "application/x-www-form-urlencoded",
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},

		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},

		{
			name:         "400 - Missing address",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpBody:     "{}",
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "address is required"),
		},

		{
			name:   "422 - Invalid checksum",
			method: http.MethodPost,
			status: http.StatusUnprocessableEntity,
			httpBody: toJSON(VerifyAddressRequest{
				Address: "7apQ7t3PZZXvjTst8G7Uvs7XH4LeM8fBPD",
			}),
			httpResponse: NewHTTPErrorResponse(http.StatusUnprocessableEntity, "Invalid checksum"),
		},

		{
			name:   "200",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: toJSON(VerifyAddressRequest{
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
			httpBody: toJSON(VerifyAddressRequest{
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
			gateway := NewGatewayerMock()

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)

			contentType := tc.contentType
			if contentType == "" {
				contentType = "application/json"
			}

			req.Header.Set("Content-Type", contentType)

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			cfg := muxConfig{host: configuredHost, appLoc: "."}
			handler := newServerMux(cfg, gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
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
