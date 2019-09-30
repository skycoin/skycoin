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

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/visor/historydb"
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
		status                  int
		err                     string
		httpBody                *httpBody
		uxid                    string
		getGetUxOutByIDArg      cipher.SHA256
		getGetUxOutByIDResponse *historydb.UxOut
		getGetUxOutByIDError    error
		httpResponse            readable.SpentOutput
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
			httpResponse:            readable.NewSpentOutput(&historydb.UxOut{}),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
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

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg readable.SpentOutput
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
		name                                string
		method                              string
		status                              int
		err                                 string
		httpBody                            *httpBody
		getSpentOutputsForAddressesArg      []cipher.Address
		getSpentOutputsForAddressesResponse [][]historydb.UxOut
		getSpentOutputsForAddressesError    error
		httpResponse                        []readable.SpentOutput
		csrfDisabled                        bool
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
			name:   "400 - gateway.GetSpentOutputsForAddresses error",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - getSpentOutputsForAddressesError",
			httpBody: &httpBody{
				address: addressForGwError.String(),
			},
			getSpentOutputsForAddressesArg:   []cipher.Address{addressForGwError},
			getSpentOutputsForAddressesError: errors.New("getSpentOutputsForAddressesError"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				address: addressForGwResponse.String(),
			},
			getSpentOutputsForAddressesArg:      []cipher.Address{addressForGwResponse},
			getSpentOutputsForAddressesResponse: [][]historydb.UxOut{{}},
			httpResponse:                        []readable.SpentOutput{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/address_uxouts"
			gateway := &MockGatewayer{}
			gateway.On("GetSpentOutputsForAddresses", tc.getSpentOutputsForAddressesArg).Return(tc.getSpentOutputsForAddressesResponse, tc.getSpentOutputsForAddressesError)

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

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []readable.SpentOutput
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
