package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"errors"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

func TestGetUxOutByID(t *testing.T) {
	invalidHash := "carccb"
	oddHash := "caccb"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"

	type httpBody struct {
		uxid string
	}

	tt := []struct {
		name                    string
		method                  string
		url                     string
		status                  int
		err                     string
		httpBody                *httpBody
		uxid                    string
		getGetUxOutByIDArg      cipher.SHA256
		getGetUxOutByIDResponse *historydb.UxOut
		getGetUxOutByIDError    error
		httpResponse            *historydb.UxOutJSON
		csrfDisabled            bool
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - empty uxin value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - uxid is empty",
			httpBody: &httpBody{
				uxid: "",
			},
		},
		{
			name:   "400 - odd length uxin value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: odd length hex string",
			httpBody: &httpBody{
				uxid: oddHash,
			},
			uxid: oddHash,
		},
		{
			name:   "400 - invalid uxin value",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			httpBody: &httpBody{
				uxid: invalidHash,
			},
			uxid: invalidHash,
		},
		{
			name:   "400 - getGetUxOutByIDError",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - getGetUxOutByIDError",
			httpBody: &httpBody{
				uxid: validHash,
			},
			uxid:                 validHash,
			getGetUxOutByIDArg:   testutil.SHA256FromHex(t, validHash),
			getGetUxOutByIDError: errors.New("getGetUxOutByIDError"),
		},
		{
			name:   "404 - uxout == nil",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			httpBody: &httpBody{
				uxid: validHash,
			},
			uxid:               validHash,
			getGetUxOutByIDArg: testutil.SHA256FromHex(t, validHash),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "404 Not Found",
			httpBody: &httpBody{
				uxid: validHash,
			},
			uxid:                    validHash,
			getGetUxOutByIDArg:      testutil.SHA256FromHex(t, validHash),
			getGetUxOutByIDResponse: &historydb.UxOut{},
			httpResponse:            historydb.NewUxOutJSON(&historydb.UxOut{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			endpoint := "/api/v1/uxout"
			gateway.On("GetUxOutByID", tc.getGetUxOutByIDArg).Return(tc.getGetUxOutByIDResponse, tc.getGetUxOutByIDError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.uxid != "" {
					v.Add("uxid", tc.httpBody.uxid)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *historydb.UxOutJSON
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestGetAddrUxOuts(t *testing.T) {
	addressForGwError := testutil.MakeAddress()
	addressForGwResponse := testutil.MakeAddress()
	type httpBody struct {
		address string
	}

	tt := []struct {
		name                  string
		method                string
		url                   string
		status                int
		err                   string
		httpBody              *httpBody
		getAddrUxOutsArg      []cipher.Address
		getAddrUxOutsResponse []*historydb.UxOut
		getAddrUxOutsError    error
		httpResponse          []*historydb.UxOutJSON
		csrfDisabled          bool
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - address is empty",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - address is empty",
			httpBody: &httpBody{
				address: "",
			},
		},
		{
			name:   "400 - cipher.DecodeBase58Address error",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid address length",
			httpBody: &httpBody{
				address: "abcd",
			},
		},
		{
			name:   "400 - gateway.GetAddrUxOuts error",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - getAddrUxOutsError",
			httpBody: &httpBody{
				address: addressForGwError.String(),
			},
			getAddrUxOutsArg:   []cipher.Address{addressForGwError},
			getAddrUxOutsError: errors.New("getAddrUxOutsError"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				address: addressForGwResponse.String(),
			},
			getAddrUxOutsArg:      []cipher.Address{addressForGwResponse},
			getAddrUxOutsResponse: []*historydb.UxOut{},
			httpResponse:          []*historydb.UxOutJSON{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/address_uxouts"
			gateway := NewGatewayerMock()
			gateway.On("GetAddrUxOuts", tc.getAddrUxOutsArg).Return(tc.getAddrUxOutsResponse, tc.getAddrUxOutsError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.address != "" {
					v.Add("address", tc.httpBody.address)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)
			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}
			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []*historydb.UxOutJSON
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
