package gui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"math"

	"time"

	"net/url"

	"errors"

	"bytes"
	"encoding/hex"

	"github.com/stretchr/testify/mock"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
)

func createUnconfirmedTxn(t *testing.T) visor.UnconfirmedTxn {
	ut := visor.UnconfirmedTxn{}
	ut.Txn = coin.Transaction{}
	ut.Txn.InnerHash = testutil.RandSHA256(t)
	ut.Received = utc.Now().UnixNano()
	ut.Checked = ut.Received
	ut.Announced = time.Time{}.UnixNano()
	return ut
}

func makeTransaction(t *testing.T) coin.Transaction {
	tx, _ := makeTransactionWithSecret(t)
	return tx
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeTransactionWithSecret(t *testing.T) (coin.Transaction, cipher.SecKey) {
	tx := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{s})
	tx.PushOutput(makeAddress(), 1e6, 50)
	tx.PushOutput(makeAddress(), 5e6, 50)
	tx.UpdateHeader()
	return tx, s
}

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func TestGetPendingTxs(t *testing.T) {
	invalidTxn := createUnconfirmedTxn(t)
	invalidTxn.Txn.Out = append(invalidTxn.Txn.Out, coin.TransactionOutput{
		Coins: math.MaxInt64 + 1,
	})

	tt := []struct {
		name                          string
		method                        string
		url                           string
		status                        int
		err                           string
		getAllUnconfirmedTxnsResponse []visor.UnconfirmedTxn
		httpResponse                  []*visor.ReadableUnconfirmedTxn
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{},
		},
		{
			name:   "500 - bad unconfirmedTxn",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{
				invalidTxn,
			},
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{},
			httpResponse:                  []*visor.ReadableUnconfirmedTxn{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/pendingTxs"
			gateway := NewGatewayerMock()
			gateway.On("GetAllUnconfirmedTxns").Return(tc.getAllUnconfirmedTxnsResponse)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []*visor.ReadableUnconfirmedTxn
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestGetTransactionByID(t *testing.T) {
	oddHash := "cafcb"
	invalidHash := "cabrca"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	type httpBody struct {
		txid string
	}

	tt := []struct {
		name                  string
		method                string
		status                int
		err                   string
		httpBody              *httpBody
		getTransactionArg     cipher.SHA256
		getTransactionReponse *visor.Transaction
		getTransactionError   error
		httpResponse          visor.TransactionResult
	}{
		{
			name:              "405",
			method:            http.MethodPost,
			status:            http.StatusMethodNotAllowed,
			err:               "405 Method Not Allowed",
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - empty txid",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - txid is empty",
			httpBody: &httpBody{
				txid: "",
			},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - invalid hash: odd length hex string",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: odd length hex string",
			httpBody: &httpBody{
				txid: oddHash,
			},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - invalid hash: invalid byte: U+0072 'r'",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			httpBody: &httpBody{
				txid: invalidHash,
			},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - getTransactionError",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - getTransactionError",
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg:   testutil.SHA256FromHex(t, validHash),
			getTransactionError: errors.New("getTransactionError"),
		},
		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg: testutil.SHA256FromHex(t, validHash),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg:     testutil.SHA256FromHex(t, validHash),
			getTransactionReponse: &visor.Transaction{},
			httpResponse: visor.TransactionResult{
				Transaction: visor.ReadableTransaction{
					Sigs:      []string{},
					In:        []string{},
					Out:       []visor.ReadableTransactionOutput{},
					Hash:      "78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/transaction"
			gateway := NewGatewayerMock()
			gateway.On("GetTransaction", tc.getTransactionArg).Return(tc.getTransactionReponse, tc.getTransactionError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.txid != "" {
					v.Add("txid", tc.httpBody.txid)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)
			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg visor.TransactionResult
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestInjectTransaction(t *testing.T) {
	validTransaction := makeTransaction(t)
	type httpBody struct {
		Rawtx string `json:"rawtx"`
	}

	validTxBody := &httpBody{Rawtx: hex.EncodeToString(validTransaction.Serialize())}
	validTxBodyJSON, err := json.Marshal(validTxBody)
	require.NoError(t, err)
	b := &httpBody{Rawtx: hex.EncodeToString(testutil.RandBytes(t, 128))}
	invalidTxBodyJSON, err := json.Marshal(b)
	require.NoError(t, err)

	tt := []struct {
		name                   string
		method                 string
		status                 int
		err                    string
		httpBody               string
		injectTransactionArg   coin.Transaction
		injectTransactionError error
		httpResponse           string
		csrfDisabled           bool
	}{
		{
			name:                 "405",
			method:               http.MethodGet,
			status:               http.StatusMethodNotAllowed,
			err:                  "405 Method Not Allowed",
			injectTransactionArg: validTransaction,
		},
		{
			name:   "400 - EOF",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - EOF",
		},
		{
			name:     "400 - Invalid transaction: Deserialization failed",
			method:   http.MethodPost,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - Invalid transaction: Deserialization failed",
			httpBody: `{"wrongKey":"wrongValue"}`,
		},
		{
			name:     "400 - encoding/hex: odd length hex string",
			method:   http.MethodPost,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - encoding/hex: odd length hex string",
			httpBody: `{"rawtx":"aab"}`,
		},
		{
			name:     "400 - rawtx deserialization error",
			method:   http.MethodPost,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - Invalid transaction: Deserialization failed",
			httpBody: string(invalidTxBodyJSON),
		},
		{
			name:                   "400 - injectTransactionError",
			method:                 http.MethodPost,
			status:                 http.StatusBadRequest,
			err:                    "400 Bad Request - inject tx failed: injectTransactionError",
			httpBody:               string(validTxBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: errors.New("injectTransactionError"),
		},
		{
			name:                 "200",
			method:               http.MethodPost,
			status:               http.StatusOK,
			httpBody:             string(validTxBodyJSON),
			injectTransactionArg: validTransaction,
			httpResponse:         validTransaction.Hash().Hex(),
		},
		{
			name:                 "200 - csrf disabled",
			method:               http.MethodPost,
			status:               http.StatusOK,
			httpBody:             string(validTxBodyJSON),
			injectTransactionArg: validTransaction,
			httpResponse:         validTransaction.Hash().Hex(),
			csrfDisabled:         true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/injectTransaction"
			gateway := NewGatewayerMock()
			gateway.On("InjectBroadcastTransaction", tc.injectTransactionArg).Return(tc.injectTransactionError)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(tc.httpBody))
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				expectedResponse, err := json.MarshalIndent(tc.httpResponse, "", "    ")
				require.NoError(t, err)
				require.Equal(t, string(expectedResponse), rr.Body.String(), tc.name)
			}
		})
	}
}

func TestResendUnconfirmedTxns(t *testing.T) {
	tt := []struct {
		name                          string
		method                        string
		status                        int
		err                           string
		httpBody                      string
		resendUnconfirmedTxnsResponse *daemon.ResendResult
		httpResponse                  *daemon.ResendResult
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
			resendUnconfirmedTxnsResponse: &daemon.ResendResult{},
			httpResponse:                  &daemon.ResendResult{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/resendUnconfirmedTxns"
			gateway := NewGatewayerMock()
			gateway.On("ResendUnconfirmedTxns").Return(tc.resendUnconfirmedTxnsResponse)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *daemon.ResendResult
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestGetRawTx(t *testing.T) {
	oddHash := "cafcb"
	invalidHash := "cabrca"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	type httpBody struct {
		txid string
	}

	tt := []struct {
		name                   string
		method                 string
		url                    string
		status                 int
		err                    string
		httpBody               *httpBody
		getTransactionArg      cipher.SHA256
		getTransactionResponse *visor.Transaction
		getTransactionError    error
		httpResponse           string
	}{
		{
			name:              "405",
			method:            http.MethodPost,
			status:            http.StatusMethodNotAllowed,
			err:               "405 Method Not Allowed",
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:              "400 - txid is empty",
			method:            http.MethodGet,
			status:            http.StatusBadRequest,
			err:               "400 Bad Request - txid is empty",
			httpBody:          &httpBody{},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - invalid hash: odd length hex string",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: odd length hex string",
			httpBody: &httpBody{
				txid: oddHash,
			},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - invalid hash: invalid byte: U+0072 'r'",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			httpBody: &httpBody{
				txid: invalidHash,
			},
			getTransactionArg: testutil.RandSHA256(t),
		},
		{
			name:   "400 - getTransactionError",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - getTransactionError",
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg:   testutil.SHA256FromHex(t, validHash),
			getTransactionError: errors.New("getTransactionError"),
		},
		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg: testutil.SHA256FromHex(t, validHash),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid: validHash,
			},
			getTransactionArg:      testutil.SHA256FromHex(t, validHash),
			getTransactionResponse: &visor.Transaction{},
			httpResponse:           "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/rawtx"
			gateway := NewGatewayerMock()
			gateway.On("GetTransaction", tc.getTransactionArg).Return(tc.getTransactionResponse, tc.getTransactionError)
			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.txid != "" {
					v.Add("txid", tc.httpBody.txid)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				expectedResponse, err := json.MarshalIndent(tc.httpResponse, "", "    ")
				require.NoError(t, err)
				require.Equal(t, string(expectedResponse), rr.Body.String(), tc.name)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	invalidAddrsStr := "invalid,addrs"
	addrsStr := "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ,2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE"
	var addrs []cipher.Address
	for _, item := range []string{"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ", "2PBmUva7J8WFsyWg979cREZkU3z2pkYjNkE"} {
		addr, err := cipher.DecodeBase58Address(item)
		require.NoError(t, err)
		addrs = append(addrs, addr)
	}
	invalidTxn := makeTransaction(t)
	invalidTxn.Out = append(invalidTxn.Out, coin.TransactionOutput{
		Coins: math.MaxInt64 + 1,
	})
	type httpBody struct {
		addrs     string
		confirmed string
	}

	tt := []struct {
		name                    string
		method                  string
		status                  int
		err                     string
		httpBody                *httpBody
		getTransactionsArg      []visor.TxFilter
		getTransactionsResponse []visor.Transaction
		getTransactionsError    error
		httpResponse            []visor.Transaction
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - invalid `addrs` param",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - parse parameter: 'addrs' failed: Invalid base58 character",
			httpBody: &httpBody{
				addrs: invalidAddrsStr,
			},
			getTransactionsArg: []visor.TxFilter{
				visor.AddrsFilter(addrs),
			},
		},
		{
			name:   "400 - invalid `confirmed` param",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid 'confirmed' value: strconv.ParseBool: parsing \"invalidConfirmed\": invalid syntax",
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "invalidConfirmed",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.AddrsFilter(addrs),
			},
		},
		{
			name:   "500 - getTransactionsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			getTransactionsError: errors.New("getTransactionsError"),
		},
		{
			name:   "500 - visor.NewTransactionResults error",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			getTransactionsResponse: []visor.Transaction{
				{
					Txn: invalidTxn,
					Status: visor.TransactionStatus{
						Confirmed: true,
						Height:    103,
					},
				},
			},
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			getTransactionsResponse: []visor.Transaction{},
			httpResponse:            []visor.Transaction{},
		},
	}

	for _, tc := range tt {
		endpoint := "/transactions"
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			gateway.On("GetTransactions", mock.Anything).Return(tc.getTransactionsResponse, tc.getTransactionsError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
				if tc.httpBody.confirmed != "" {
					v.Add("confirmed", tc.httpBody.confirmed)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []visor.Transaction
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}
