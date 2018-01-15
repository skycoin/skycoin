package gui

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

// GetUxOutByID gets UxOut by hash id.
func (gw *FakeGateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	args := gw.Called(id)
	return args.Get(0).(*historydb.UxOut), args.Error(1)
}

// GetAddrUxOuts gets all the address affected UxOuts.
func (gw *FakeGateway) GetAddrUxOuts(addr cipher.Address) ([]*historydb.UxOutJSON, error) {
	args := gw.Called(addr)
	return args.Get(0).([]*historydb.UxOutJSON), args.Error(1)
}

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
	}{
		{
			"405",
			http.MethodPost,
			"/uxout",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			nil,
			"uxid",
			testutil.RandSHA256(t),
			nil,
			nil,
			nil,
		},
		{
			"400 - empty uxin value",
			http.MethodGet,
			"/uxout",
			http.StatusBadRequest,
			"400 Bad Request - uxid is empty",
			&httpBody{
				uxid: "",
			},
			"",
			testutil.RandSHA256(t),
			nil,
			nil,
			nil,
		},
		{
			"400 - odd length uxin value",
			http.MethodGet,
			"/uxout",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			&httpBody{
				uxid: oddHash,
			},
			oddHash,
			testutil.RandSHA256(t),
			nil,
			nil,
			nil,
		},
		{
			"400 - invalid uxin value",
			http.MethodGet,
			"/uxout",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			&httpBody{
				uxid: invalidHash,
			},
			invalidHash,
			testutil.RandSHA256(t),
			nil,
			nil,
			nil,
		},
		{
			"400 - getGetUxOutByIDError",
			http.MethodGet,
			"/uxout",
			http.StatusBadRequest,
			"400 Bad Request - getGetUxOutByIDError",
			&httpBody{
				uxid: validHash,
			},
			validHash,
			testutil.SHA256FromHex(t, validHash),
			nil,
			errors.New("getGetUxOutByIDError"),
			nil,
		},
		{
			"404 - uxout == nil",
			http.MethodGet,
			"/uxout",
			http.StatusNotFound,
			"404 Not Found",
			&httpBody{
				uxid: validHash,
			},
			validHash,
			testutil.SHA256FromHex(t, validHash),
			nil,
			nil,
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/uxout",
			http.StatusOK,
			"404 Not Found",
			&httpBody{
				uxid: validHash,
			},
			validHash,
			testutil.SHA256FromHex(t, validHash),
			&historydb.UxOut{},
			nil,
			historydb.NewUxOutJSON(&historydb.UxOut{}),
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("GetUxOutByID", tc.getGetUxOutByIDArg).Return(tc.getGetUxOutByIDResponse, tc.getGetUxOutByIDError)

		var urlFull = tc.url
		v := url.Values{}
		if tc.httpBody != nil {
			if tc.httpBody.uxid != "" {
				v.Add("uxid", tc.httpBody.uxid)
			}
		}

		if len(v) > 0 {
			urlFull += "?" + v.Encode()
		}

		req, err := http.NewRequest(tc.method, urlFull, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getUxOutByID(gateway))

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
		getAddrUxOutsArg      cipher.Address
		getAddrUxOutsResponse []*historydb.UxOutJSON
		getAddrUxOutsError    error
		httpResponse          []*historydb.UxOutJSON
	}{
		{
			"405",
			http.MethodPost,
			"/address_uxouts",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			nil,
			testutil.MakeAddress(),
			nil,
			nil,
			nil,
		},
		{
			"400 - address is empty",
			http.MethodGet,
			"/address_uxouts",
			http.StatusBadRequest,
			"400 Bad Request - address is empty",
			&httpBody{
				address: "",
			},
			testutil.MakeAddress(),
			nil,
			nil,
			nil,
		},
		{
			"400 - cipher.DecodeBase58Address error",
			http.MethodGet,
			"/address_uxouts",
			http.StatusBadRequest,
			"400 Bad Request - Invalid address length",
			&httpBody{
				address: "abcd",
			},
			testutil.MakeAddress(),
			nil,
			nil,
			nil,
		},
		{
			"400 - gateway.GetAddrUxOuts error",
			http.MethodGet,
			"/address_uxouts",
			http.StatusBadRequest,
			"400 Bad Request - getAddrUxOutsError",
			&httpBody{
				address: addressForGwError.String(),
			},
			addressForGwError,
			nil,
			errors.New("getAddrUxOutsError"),
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/address_uxouts",
			http.StatusOK,
			"",
			&httpBody{
				address: addressForGwResponse.String(),
			},
			addressForGwResponse,
			[]*historydb.UxOutJSON{},
			nil,
			[]*historydb.UxOutJSON{},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("GetAddrUxOuts", tc.getAddrUxOutsArg).Return(tc.getAddrUxOutsResponse, tc.getAddrUxOutsError)

		var urlFull = tc.url
		v := url.Values{}
		if tc.httpBody != nil {
			if tc.httpBody.address != "" {
				v.Add("address", tc.httpBody.address)
			}
		}

		if len(v) > 0 {
			urlFull += "?" + v.Encode()
		}

		req, err := http.NewRequest(tc.method, urlFull, nil)
		require.NoError(t, err)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getAddrUxOuts(gateway))

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
	}
}
