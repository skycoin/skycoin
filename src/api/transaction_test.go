package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"math"

	"time"

	"net/url"

	"errors"

	"bytes"
	"encoding/hex"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
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

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeTransaction(t *testing.T) coin.Transaction {
	txn := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	txn.PushInput(ux.Hash())
	txn.SignInputs([]cipher.SecKey{s})
	txn.PushOutput(makeAddress(), 1e6, 50)
	txn.PushOutput(makeAddress(), 5e6, 50)
	txn.UpdateHeader()
	return txn
}

func TestGetPendingTxs(t *testing.T) {
	invalidTxn := createUnconfirmedTxn(t)
	invalidTxn.Txn.Out = append(invalidTxn.Txn.Out, coin.TransactionOutput{
		Coins: math.MaxInt64 + 1,
	})

	tt := []struct {
		name                                 string
		method                               string
		status                               int
		err                                  string
		verbose                              bool
		verboseStr                           string
		getAllUnconfirmedTxnsResponse        []visor.UnconfirmedTxn
		getAllUnconfirmedTxnsErr             error
		getAllUnconfirmedTxnsVerboseResponse []visor.ReadableUnconfirmedTxnVerbose
		getAllUnconfirmedTxnsVerboseErr      error
		httpResponse                         interface{}
	}{
		{
			name:                          "405",
			method:                        http.MethodPost,
			status:                        http.StatusMethodNotAllowed,
			err:                           "405 Method Not Allowed",
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{},
		},
		{
			name:       "400 - bad verbose",
			method:     http.MethodGet,
			status:     http.StatusBadRequest,
			err:        "400 Bad Request - Invalid value for verbose",
			verboseStr: "foo",
		},
		{
			name:   "500 - bad unconfirmedTxn",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - Droplet string conversion failed: Value is too large",
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{
				invalidTxn,
			},
		},
		{
			name:                     "500 - get unconfirmedTxn error",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - GetAllUnconfirmedTxns failed",
			getAllUnconfirmedTxnsErr: errors.New("GetAllUnconfirmedTxns failed"),
		},
		{
			name:                            "500 - get unconfirmedTxnVerbose error",
			method:                          http.MethodGet,
			status:                          http.StatusInternalServerError,
			verboseStr:                      "1",
			verbose:                         true,
			err:                             "500 Internal Server Error - GetAllUnconfirmedTxnsVerbose failed",
			getAllUnconfirmedTxnsVerboseErr: errors.New("GetAllUnconfirmedTxnsVerbose failed"),
		},
		{
			name:                          "200",
			method:                        http.MethodGet,
			status:                        http.StatusOK,
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTxn{},
			httpResponse:                  []visor.ReadableUnconfirmedTxn{},
		},
		{
			name:                                 "200 verbose",
			method:                               http.MethodGet,
			status:                               http.StatusOK,
			verboseStr:                           "1",
			verbose:                              true,
			getAllUnconfirmedTxnsVerboseResponse: []visor.ReadableUnconfirmedTxnVerbose{},
			httpResponse:                         []visor.ReadableUnconfirmedTxnVerbose{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/pendingTxs"
			gateway := &MockGatewayer{}
			gateway.On("GetAllUnconfirmedTxns").Return(tc.getAllUnconfirmedTxnsResponse, tc.getAllUnconfirmedTxnsErr)
			gateway.On("GetAllUnconfirmedTxnsVerbose").Return(tc.getAllUnconfirmedTxnsVerboseResponse, tc.getAllUnconfirmedTxnsVerboseErr)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

			v := url.Values{}
			if tc.verboseStr != "" {
				v.Add("verbose", tc.verboseStr)
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

			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg []visor.ReadableUnconfirmedTxnVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				} else {
					var msg []visor.ReadableUnconfirmedTxn
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				}
			}
		})
	}
}

func TestGetTransactionByID(t *testing.T) {
	oddHash := "cafcb"
	invalidHash := "cabrca"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	type httpBody struct {
		txid    string
		verbose string
		encoded string
	}

	tt := []struct {
		name                               string
		method                             string
		status                             int
		err                                string
		httpBody                           *httpBody
		verbose                            bool
		encoded                            bool
		txid                               cipher.SHA256
		getTransactionReponse              *visor.Transaction
		getTransactionError                error
		getTransactionResultReponse        *daemon.TransactionResult
		getTransactionResultError          error
		getTransactionResultVerboseReponse *daemon.TransactionResultVerbose
		getTransactionResultVerboseError   error
		httpResponse                       interface{}
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
			txid:   testutil.RandSHA256(t),
		},

		{
			name:   "400 - empty txid",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - txid is empty",
			httpBody: &httpBody{
				txid: "",
			},
			txid: testutil.RandSHA256(t),
		},

		{
			name:   "400 - invalid hash: odd length hex string",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: odd length hex string",
			httpBody: &httpBody{
				txid: oddHash,
			},
			txid: testutil.RandSHA256(t),
		},

		{
			name:   "400 - invalid hash: invalid byte: U+0072 'r'",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			httpBody: &httpBody{
				txid: invalidHash,
			},
			txid: testutil.RandSHA256(t),
		},

		{
			name:   "400 - invalid verbose",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid value for verbose",
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "foo",
			},
		},

		{
			name:   "400 - invalid encoded",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid value for encoded",
			httpBody: &httpBody{
				txid:    validHash,
				encoded: "foo",
			},
		},

		{
			name:   "400 - verbose and encoded combined",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - verbose and encoded cannot be combined",
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "1",
				encoded: "1",
			},
		},

		{
			name:   "500 - getTransactionError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionError",
			httpBody: &httpBody{
				txid:    validHash,
				encoded: "1",
			},
			encoded:             true,
			txid:                testutil.SHA256FromHex(t, validHash),
			getTransactionError: errors.New("getTransactionError"),
		},

		{
			name:   "500 - getTransactionResultError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionResultError",
			httpBody: &httpBody{
				txid: validHash,
			},
			txid:                      testutil.SHA256FromHex(t, validHash),
			getTransactionResultError: errors.New("getTransactionResultError"),
		},

		{
			name:   "500 - getTransactionResultVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionResultVerboseError",
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "1",
			},
			verbose:                          true,
			txid:                             testutil.SHA256FromHex(t, validHash),
			getTransactionResultVerboseError: errors.New("getTransactionResultVerboseError"),
		},

		{
			name:   "404",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			httpBody: &httpBody{
				txid: validHash,
			},
			txid: testutil.SHA256FromHex(t, validHash),
		},

		{
			name:   "404 verbose",
			method: http.MethodGet,
			status: http.StatusNotFound,
			err:    "404 Not Found",
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "1",
			},
			verbose: true,
			txid:    testutil.SHA256FromHex(t, validHash),
		},

		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid: validHash,
			},
			txid:                        testutil.SHA256FromHex(t, validHash),
			getTransactionResultReponse: &daemon.TransactionResult{},
			httpResponse:                &daemon.TransactionResult{},
		},

		{
			name:   "200 verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "1",
			},
			verbose:                            true,
			txid:                               testutil.SHA256FromHex(t, validHash),
			getTransactionResultVerboseReponse: &daemon.TransactionResultVerbose{},
			httpResponse:                       &daemon.TransactionResultVerbose{},
		},

		{
			name:   "200 encoded",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid:    validHash,
				encoded: "1",
			},
			encoded:               true,
			txid:                  testutil.SHA256FromHex(t, validHash),
			getTransactionReponse: &visor.Transaction{},
			httpResponse: &TransactionEncodedResponse{
				EncodedTransaction: "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/transaction"
			gateway := &MockGatewayer{}
			gateway.On("GetTransaction", tc.txid).Return(tc.getTransactionReponse, tc.getTransactionError)
			gateway.On("GetTransactionResult", tc.txid).Return(tc.getTransactionResultReponse, tc.getTransactionResultError)
			gateway.On("GetTransactionResultVerbose", tc.txid).Return(tc.getTransactionResultVerboseReponse, tc.getTransactionResultVerboseError)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.txid != "" {
					v.Add("txid", tc.httpBody.txid)
				}
				if tc.httpBody.verbose != "" {
					v.Add("verbose", tc.httpBody.verbose)
				}
				if tc.httpBody.encoded != "" {
					v.Add("encoded", tc.httpBody.encoded)
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg daemon.TransactionResultVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				} else if tc.encoded {
					var msg TransactionEncodedResponse
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				} else {
					var msg daemon.TransactionResult
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				}
			}

			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)
		})
	}
}

func TestInjectTransaction(t *testing.T) {
	validTransaction := makeTransaction(t)
	type httpBody struct {
		Rawtx string `json:"rawtx"`
	}

	validTxnBody := &httpBody{Rawtx: hex.EncodeToString(validTransaction.Serialize())}
	validTxnBodyJSON, err := json.Marshal(validTxnBody)
	require.NoError(t, err)

	b := &httpBody{Rawtx: hex.EncodeToString(testutil.RandBytes(t, 128))}
	invalidTxnBodyJSON, err := json.Marshal(b)
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
			httpBody: string(invalidTxnBodyJSON),
		},
		{
			name:                   "503 - injectTransactionError",
			method:                 http.MethodPost,
			status:                 http.StatusServiceUnavailable,
			err:                    "503 Service Unavailable - inject tx failed: injectTransactionError",
			httpBody:               string(validTxnBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: errors.New("injectTransactionError"),
		},
		{
			name:                 "200",
			method:               http.MethodPost,
			status:               http.StatusOK,
			httpBody:             string(validTxnBodyJSON),
			injectTransactionArg: validTransaction,
			httpResponse:         validTransaction.Hash().Hex(),
		},
		{
			name:                 "200 - csrf disabled",
			method:               http.MethodPost,
			status:               http.StatusOK,
			httpBody:             string(validTxnBodyJSON),
			injectTransactionArg: validTransaction,
			httpResponse:         validTransaction.Hash().Hex(),
			csrfDisabled:         true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/injectTransaction"
			gateway := &MockGatewayer{}
			gateway.On("InjectBroadcastTransaction", tc.injectTransactionArg).Return(tc.injectTransactionError)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
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
		resendUnconfirmedTxnsErr      error
		httpResponse                  *daemon.ResendResult
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                     "500 resend failed",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - ResendUnconfirmedTxns failed",
			resendUnconfirmedTxnsErr: errors.New("ResendUnconfirmedTxns failed"),
		},
		{
			name:                          "200",
			method:                        http.MethodGet,
			status:                        http.StatusOK,
			resendUnconfirmedTxnsResponse: &daemon.ResendResult{},
			httpResponse:                  &daemon.ResendResult{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/resendUnconfirmedTxns"
			gateway := &MockGatewayer{}
			gateway.On("ResendUnconfirmedTxns").Return(tc.resendUnconfirmedTxnsResponse, tc.resendUnconfirmedTxnsErr)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
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
			endpoint := "/api/v1/rawtx"
			gateway := &MockGatewayer{}
			gateway.On("GetTransaction", tc.getTransactionArg).Return(tc.getTransactionResponse, tc.getTransactionError)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
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

	type httpBody struct {
		addrs     string
		confirmed string
		verbose   string
	}

	tt := []struct {
		name                           string
		method                         string
		status                         int
		err                            string
		httpBody                       *httpBody
		verbose                        bool
		getTransactionsArg             []visor.TxFilter
		getTransactionsResponse        *daemon.TransactionResults
		getTransactionsError           error
		getTransactionsVerboseResponse *daemon.TransactionResultsVerbose
		getTransactionsVerboseError    error
		httpResponse                   interface{}
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
				visor.NewAddrsFilter(addrs),
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
				visor.NewAddrsFilter(addrs),
			},
		},

		{
			name:   "400 - invalid verbose",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - Invalid value for verbose",
			httpBody: &httpBody{
				addrs:   addrsStr,
				verbose: "foo",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.NewAddrsFilter(addrs),
			},
		},

		{
			name:   "500 - getTransactionsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionsError",
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.NewAddrsFilter(addrs),
				visor.NewConfirmedTxFilter(true),
			},
			getTransactionsError: errors.New("getTransactionsError"),
		},

		{
			name:   "500 - getTransactionsVerboseError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionsVerboseError",
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
				verbose:   "1",
			},
			verbose: true,
			getTransactionsArg: []visor.TxFilter{
				visor.NewAddrsFilter(addrs),
				visor.NewConfirmedTxFilter(true),
			},
			getTransactionsVerboseError: errors.New("getTransactionsVerboseError"),
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
				visor.NewAddrsFilter(addrs),
				visor.NewConfirmedTxFilter(true),
			},
			getTransactionsResponse: &daemon.TransactionResults{},
			httpResponse:            []daemon.TransactionResult{},
		},

		{
			name:   "200 verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
				verbose:   "1",
			},
			verbose: true,
			getTransactionsArg: []visor.TxFilter{
				visor.NewAddrsFilter(addrs),
				visor.NewConfirmedTxFilter(true),
			},
			getTransactionsVerboseResponse: &daemon.TransactionResultsVerbose{},
			httpResponse:                   []daemon.TransactionResultVerbose{},
		},
	}

	for _, tc := range tt {
		endpoint := "/api/v1/transactions"
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			// Custom argument matching function for matching TxFilter args
			matchFunc := mock.MatchedBy(func(flts []visor.TxFilter) bool {
				if len(flts) != len(tc.getTransactionsArg) {
					return false
				}

				for i, f := range flts {
					switch f.(type) {
					case visor.AddrsFilter:
						flt, ok := tc.getTransactionsArg[i].(visor.AddrsFilter)
						if !ok {
							return false
						}

						if len(flt.Addrs) != len(f.(visor.AddrsFilter).Addrs) {
							return false
						}

						for j, a := range flt.Addrs {
							ab := a.Bytes()
							bb := f.(visor.AddrsFilter).Addrs[j].Bytes()
							if !bytes.Equal(ab[:], bb[:]) {
								return false
							}
						}

					case visor.BaseFilter:
						// This part assumes that the filter is a ConfirmedTxFilter
						flt, ok := tc.getTransactionsArg[i].(visor.BaseFilter)
						if !ok {
							return false
						}

						dummyTxn := &visor.Transaction{
							Status: visor.TransactionStatus{
								Confirmed: true,
							},
						}

						if flt.F(dummyTxn) != f.(visor.BaseFilter).F(dummyTxn) {
							return false
						}

					default:
						return false
					}
				}

				return true
			})

			gateway.On("GetTransactionResults", matchFunc).Return(tc.getTransactionsResponse, tc.getTransactionsError)
			gateway.On("GetTransactionResultsVerbose", matchFunc).Return(tc.getTransactionsVerboseResponse, tc.getTransactionsVerboseError)
			gateway.On("IsAPISetEnabled", "READ_ONLY", []string(nil)).Return(true)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
				if tc.httpBody.confirmed != "" {
					v.Add("confirmed", tc.httpBody.confirmed)
				}
				if tc.httpBody.verbose != "" {
					v.Add("verbose", tc.httpBody.verbose)
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg []daemon.TransactionResultVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				} else {
					var msg []daemon.TransactionResult
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				}
			}
		})
	}
}

type transactionAndInputs struct {
	txn    coin.Transaction
	inputs []wallet.UxBalance
}

func newVerifyTxnResponseJSON(t *testing.T, txn *coin.Transaction, inputs []wallet.UxBalance, isTxnConfirmed bool) VerifyTxnResponse {
	ctxn, err := newCreatedTransactionFuzzy(txn, inputs)
	require.NoError(t, err)
	return VerifyTxnResponse{
		Transaction: *ctxn,
		Confirmed:   isTxnConfirmed,
	}
}

func prepareTxnAndInputs(t *testing.T) transactionAndInputs {
	txn := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	txn.PushInput(ux.Hash())
	txn.SignInputs([]cipher.SecKey{s})
	txn.PushOutput(makeAddress(), 1e6, 50)
	txn.PushOutput(makeAddress(), 5e6, 50)
	txn.UpdateHeader()

	input, err := wallet.NewUxBalance(uint64(utc.UnixNow()), ux)
	require.NoError(t, err)

	return transactionAndInputs{txn: txn, inputs: []wallet.UxBalance{input}}
}

func makeTransactionWithEmptyAddressOutput(t *testing.T) transactionAndInputs {
	txn := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	txn.PushInput(ux.Hash())
	txn.SignInputs([]cipher.SecKey{s})
	txn.PushOutput(makeAddress(), 1e6, 50)
	txn.PushOutput(cipher.Address{}, 5e6, 50)
	txn.UpdateHeader()

	input, err := wallet.NewUxBalance(uint64(utc.UnixNow()), ux)
	require.NoError(t, err)

	return transactionAndInputs{txn: txn, inputs: []wallet.UxBalance{input}}
}

func TestVerifyTransaction(t *testing.T) {
	txnAndInputs := prepareTxnAndInputs(t)
	type httpBody struct {
		EncodedTransaction string `json:"encoded_transaction"`
	}

	validTxnBody := &httpBody{EncodedTransaction: hex.EncodeToString(txnAndInputs.txn.Serialize())}
	validTxnBodyJSON, err := json.Marshal(validTxnBody)
	require.NoError(t, err)

	b := &httpBody{EncodedTransaction: hex.EncodeToString(testutil.RandBytes(t, 128))}
	invalidTxnBodyJSON, err := json.Marshal(b)
	require.NoError(t, err)

	invalidTxnEmptyAddress := makeTransactionWithEmptyAddressOutput(t)
	invalidTxnEmptyAddressBody := &httpBody{
		EncodedTransaction: hex.EncodeToString(invalidTxnEmptyAddress.txn.Serialize()),
	}
	invalidTxnEmptyAddressBodyJSON, err := json.Marshal(invalidTxnEmptyAddressBody)
	require.NoError(t, err)

	type verifyTxnVerboseResult struct {
		Uxouts         []wallet.UxBalance
		IsTxnConfirmed bool
		Err            error
	}

	tt := []struct {
		name                          string
		method                        string
		contentType                   string
		status                        int
		httpBody                      string
		gatewayVerifyTxnVerboseArg    coin.Transaction
		gatewayVerifyTxnVerboseResult verifyTxnVerboseResult
		httpResponse                  HTTPResponse
		csrfDisabled                  bool
	}{
		{
			name:                       "405",
			method:                     http.MethodGet,
			status:                     http.StatusMethodNotAllowed,
			gatewayVerifyTxnVerboseArg: txnAndInputs.txn,
			httpResponse:               NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""),
		},
		{
			name:         "400 - EOF",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "EOF"),
		},
		{
			name:         "415 - Unsupported Media Type",
			method:       http.MethodPost,
			contentType:  "",
			status:       http.StatusUnsupportedMediaType,
			httpResponse: NewHTTPErrorResponse(http.StatusUnsupportedMediaType, ""),
		},
		{
			name:         "400 - Invalid transaction: Deserialization failed",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpBody:     `{"wrongKey":"wrongValue"}`,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "decode transaction failed: Invalid transaction: Deserialization failed"),
		},
		{
			name:         "400 - encoding/hex: odd length hex string",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpBody:     `{"encoded_transaction":"aab"}`,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "decode transaction failed: encoding/hex: odd length hex string"),
		},
		{
			name:         "400 - deserialization error",
			method:       http.MethodPost,
			contentType:  "application/json",
			status:       http.StatusBadRequest,
			httpBody:     string(invalidTxnBodyJSON),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "decode transaction failed: Invalid transaction: Deserialization failed"),
		},
		{
			name:                       "422 - txn sends to empty address",
			method:                     http.MethodPost,
			contentType:                "application/json",
			status:                     http.StatusUnprocessableEntity,
			httpBody:                   string(invalidTxnEmptyAddressBodyJSON),
			gatewayVerifyTxnVerboseArg: invalidTxnEmptyAddress.txn,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: invalidTxnEmptyAddress.inputs,
				Err:    visor.NewErrTxnViolatesUserConstraint(errors.New("Transaction.Out contains an output sending to an empty address")),
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &invalidTxnEmptyAddress.txn, invalidTxnEmptyAddress.inputs, false),
				Error: &HTTPError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Transaction violates user constraint: Transaction.Out contains an output sending to an empty address",
				},
			},
		},
		{
			name:                       "500 - internal server error",
			method:                     http.MethodPost,
			contentType:                "application/json",
			status:                     http.StatusInternalServerError,
			httpBody:                   string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg: txnAndInputs.txn,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Err: errors.New("verify transaction failed"),
			},
			httpResponse: NewHTTPErrorResponse(http.StatusInternalServerError, "verify transaction failed"),
		},
		{
			name:                       "422 - txn is confirmed",
			method:                     http.MethodPost,
			contentType:                "application/json",
			status:                     http.StatusUnprocessableEntity,
			httpBody:                   string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg: txnAndInputs.txn,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts:         txnAndInputs.inputs,
				IsTxnConfirmed: true,
			},
			httpResponse: HTTPResponse{
				Error: &HTTPError{
					Message: "transaction has been spent",
					Code:    http.StatusUnprocessableEntity,
				},
				Data: newVerifyTxnResponseJSON(t, &txnAndInputs.txn, txnAndInputs.inputs, true),
			},
		},
		{
			name:                       "200",
			method:                     http.MethodPost,
			contentType:                "application/json",
			status:                     http.StatusOK,
			httpBody:                   string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg: txnAndInputs.txn,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: txnAndInputs.inputs,
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &txnAndInputs.txn, txnAndInputs.inputs, false),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/transaction/verify"
			gateway := &MockGatewayer{}
			gateway.On("VerifyTxnVerbose", &tc.gatewayVerifyTxnVerboseArg).Return(tc.gatewayVerifyTxnVerboseResult.Uxouts,
				tc.gatewayVerifyTxnVerboseResult.IsTxnConfirmed, tc.gatewayVerifyTxnVerboseResult.Err)

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tc.contentType)

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "handler returned wrong status code: got `%v` want `%v`", status, tc.status)

			var rsp ReceivedHTTPResponse
			err = json.NewDecoder(rr.Body).Decode(&rsp)
			require.NoError(t, err)

			require.Equal(t, tc.httpResponse.Error, rsp.Error)

			if rsp.Data == nil {
				require.Nil(t, tc.httpResponse.Data)
			} else {
				require.NotNil(t, tc.httpResponse.Data)

				var txnRsp VerifyTxnResponse
				err := json.Unmarshal(rsp.Data, &txnRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(VerifyTxnResponse), txnRsp)
			}
		})
	}
}
