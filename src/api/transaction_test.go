package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
)

func createUnconfirmedTxn(t *testing.T) visor.UnconfirmedTransaction {
	ut := visor.UnconfirmedTransaction{}
	ut.Transaction = coin.Transaction{}
	ut.Transaction.InnerHash = testutil.RandSHA256(t)
	ut.Transaction.In = []cipher.SHA256{testutil.RandSHA256(t)}
	ut.Received = time.Now().UTC().UnixNano()
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

	err := txn.PushInput(ux.Hash())
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 5e6, 50)
	require.NoError(t, err)
	txn.SignInputs([]cipher.SecKey{s})
	err = txn.UpdateHeader()
	require.NoError(t, err)
	return txn
}

func TestGetPendingTxs(t *testing.T) {
	invalidTxn := createUnconfirmedTxn(t)
	invalidTxn.Transaction.Out = append(invalidTxn.Transaction.Out, coin.TransactionOutput{
		Coins: math.MaxInt64 + 1,
	})

	type verboseResult struct {
		Transactions []visor.UnconfirmedTransaction
		Inputs       [][]visor.TransactionInput
	}

	tt := []struct {
		name                                 string
		method                               string
		status                               int
		err                                  string
		verbose                              bool
		verboseStr                           string
		getAllUnconfirmedTxnsResponse        []visor.UnconfirmedTransaction
		getAllUnconfirmedTxnsErr             error
		getAllUnconfirmedTxnsVerboseResponse verboseResult
		getAllUnconfirmedTxnsVerboseErr      error
		httpResponse                         interface{}
	}{
		{
			name:                          "405",
			method:                        http.MethodPost,
			status:                        http.StatusMethodNotAllowed,
			err:                           "405 Method Not Allowed",
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTransaction{},
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
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTransaction{
				invalidTxn,
			},
		},
		{
			name:                     "500 - get unconfirmedTxn error",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - GetAllUnconfirmedTransactions failed",
			getAllUnconfirmedTxnsErr: errors.New("GetAllUnconfirmedTransactions failed"),
		},
		{
			name:                            "500 - get unconfirmedTxnVerbose error",
			method:                          http.MethodGet,
			status:                          http.StatusInternalServerError,
			verboseStr:                      "1",
			verbose:                         true,
			err:                             "500 Internal Server Error - GetAllUnconfirmedTransactionsVerbose failed",
			getAllUnconfirmedTxnsVerboseErr: errors.New("GetAllUnconfirmedTransactionsVerbose failed"),
		},
		{
			name:                          "200",
			method:                        http.MethodGet,
			status:                        http.StatusOK,
			getAllUnconfirmedTxnsResponse: []visor.UnconfirmedTransaction{},
			httpResponse:                  []readable.UnconfirmedTransactions{},
		},
		{
			name:       "200 verbose",
			method:     http.MethodGet,
			status:     http.StatusOK,
			verboseStr: "1",
			verbose:    true,
			getAllUnconfirmedTxnsVerboseResponse: verboseResult{
				Transactions: []visor.UnconfirmedTransaction{},
				Inputs:       [][]visor.TransactionInput{},
			},
			httpResponse: []readable.UnconfirmedTransactionVerbose{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/pendingTxs"
			gateway := &MockGatewayer{}
			gateway.On("GetAllUnconfirmedTransactions").Return(tc.getAllUnconfirmedTxnsResponse, tc.getAllUnconfirmedTxnsErr)
			gateway.On("GetAllUnconfirmedTransactionsVerbose").Return(tc.getAllUnconfirmedTxnsVerboseResponse.Transactions,
				tc.getAllUnconfirmedTxnsVerboseResponse.Inputs, tc.getAllUnconfirmedTxnsVerboseErr)

			v := url.Values{}
			if tc.verboseStr != "" {
				v.Add("verbose", tc.verboseStr)
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg []readable.UnconfirmedTransactionVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				} else {
					var msg []readable.UnconfirmedTransactions
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
	validHashRaw, err := cipher.SHA256FromHex(validHash)
	require.NoError(t, err)

	validAddr := "28ATuZGXJm6dJDGyJbdknFWgv8kbBX9hAdN"
	validAddrRaw, err := cipher.DecodeBase58Address(validAddr)
	require.NoError(t, err)

	validSig := "cca1595fb27375789da47bb1cf78e14febc2be6f3c3034247fea6f700b853cddbab5d16f4ffc1912fca8373f10e468b745d6a1d686cb73ade1e3c3b3653b2f9d7f"
	validSigRaw, err := cipher.SigFromHex(validSig)
	require.NoError(t, err)

	type httpBody struct {
		txid    string
		verbose string
		encoded string
	}

	type verboseResult struct {
		Transaction *visor.Transaction
		Inputs      []visor.TransactionInput
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
		getTransactionResultVerboseReponse verboseResult
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
			name:   "500 - getTransactionError encoded",
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
			name:   "500 - getTransactionError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - getTransactionError",
			httpBody: &httpBody{
				txid: validHash,
			},
			txid:                testutil.SHA256FromHex(t, validHash),
			getTransactionError: errors.New("getTransactionError"),
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
			txid: testutil.SHA256FromHex(t, validHash),
			getTransactionReponse: &visor.Transaction{
				Transaction: coin.Transaction{
					Sigs: []cipher.Sig{validSigRaw},
					In:   []cipher.SHA256{validHashRaw},
					Out: []coin.TransactionOutput{
						{
							Coins:   9999,
							Hours:   1111,
							Address: validAddrRaw,
						},
					},
				},
				Status: visor.TransactionStatus{
					Confirmed: true,
					BlockSeq:  100,
					Height:    9,
				},
			},
			httpResponse: &readable.TransactionWithStatus{
				Status: readable.TransactionStatus{
					Confirmed: true,
					BlockSeq:  100,
					Height:    9,
				},
				Transaction: readable.Transaction{
					Hash:      "b64525bc14edb3c838ff3ef4f01bd74712432b32c18463dbda59b431959b2e52",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Sigs:      []string{validSig},
					In:        []string{validHash},
					Out: []readable.TransactionOutput{
						{
							Hash:    "87ec4d440fd64bb4c26839d58684e567e499265ca396649c03304b928378720b",
							Coins:   "0.009999",
							Hours:   1111,
							Address: validAddr,
						},
					},
				},
			},
		},

		{
			name:   "200 verbose",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid:    validHash,
				verbose: "1",
			},
			verbose: true,
			txid:    testutil.SHA256FromHex(t, validHash),
			getTransactionResultVerboseReponse: verboseResult{
				Transaction: &visor.Transaction{
					Transaction: coin.Transaction{
						Sigs: []cipher.Sig{validSigRaw},
						In:   []cipher.SHA256{validHashRaw},
						Out: []coin.TransactionOutput{
							{
								Coins:   9999,
								Hours:   1111,
								Address: validAddrRaw,
							},
						},
					},
					Status: visor.TransactionStatus{
						Confirmed: true,
						BlockSeq:  100,
						Height:    9,
					},
				},
				Inputs: []visor.TransactionInput{
					{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   9999,
								Hours:   1111,
								Address: validAddrRaw,
							},
						},
						CalculatedHours: 3333,
					},
				},
			},
			httpResponse: &readable.TransactionWithStatusVerbose{
				Status: readable.TransactionStatus{
					Confirmed: true,
					BlockSeq:  100,
					Height:    9,
				},
				Transaction: readable.TransactionVerbose{
					BlockTransactionVerbose: readable.BlockTransactionVerbose{
						Fee:       2222,
						Hash:      "b64525bc14edb3c838ff3ef4f01bd74712432b32c18463dbda59b431959b2e52",
						InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
						Sigs:      []string{validSig},
						In: []readable.TransactionInput{
							{
								Hash:            "50e8ad459e29a051d969f221f1fb9775e26248e8b443982fef0cfaa117ee6c0c",
								Coins:           "0.009999",
								Hours:           1111,
								CalculatedHours: 3333,
								Address:         validAddr,
							},
						},
						Out: []readable.TransactionOutput{
							{
								Hash:    "87ec4d440fd64bb4c26839d58684e567e499265ca396649c03304b928378720b",
								Coins:   "0.009999",
								Hours:   1111,
								Address: validAddr,
							},
						},
					},
				},
			},
		},

		{
			name:   "200 encoded",
			method: http.MethodGet,
			status: http.StatusOK,
			httpBody: &httpBody{
				txid:    validHash,
				encoded: "1",
			},
			encoded: true,
			txid:    testutil.SHA256FromHex(t, validHash),
			getTransactionReponse: &visor.Transaction{
				Transaction: coin.Transaction{
					Sigs: []cipher.Sig{validSigRaw},
					In:   []cipher.SHA256{validHashRaw},
					Out: []coin.TransactionOutput{
						{
							Coins:   9999,
							Hours:   1111,
							Address: validAddrRaw,
						},
					},
				},
				Status: visor.TransactionStatus{
					Confirmed: true,
					BlockSeq:  100,
					Height:    9,
				},
			},
			httpResponse: &TransactionEncodedResponse{
				Status: readable.TransactionStatus{
					Confirmed: true,
					BlockSeq:  100,
					Height:    9,
				},
				EncodedTransaction: "0000000000000000000000000000000000000000000000000000000000000000000000000001000000cca1595fb27375789da47bb1cf78e14febc2be6f3c3034247fea6f700b853cddbab5d16f4ffc1912fca8373f10e468b745d6a1d686cb73ade1e3c3b3653b2f9d7f0100000079216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b0100000000a1f1da0612c870cbb2d88fb3d7f95ba7118d6efb0f270000000000005704000000000000",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/transaction"
			gateway := &MockGatewayer{}
			gateway.On("GetTransaction", tc.txid).Return(tc.getTransactionReponse, tc.getTransactionError)
			gateway.On("GetTransactionWithInputs", tc.txid).Return(tc.getTransactionResultVerboseReponse.Transaction,
				tc.getTransactionResultVerboseReponse.Inputs, tc.getTransactionResultVerboseError)

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

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg readable.TransactionWithStatusVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				} else if tc.encoded {
					var msg TransactionEncodedResponse
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				} else {
					var msg readable.TransactionWithStatus
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, &msg, tc.name)
				}
			}

			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)
		})
	}
}

func TestInjectTransaction(t *testing.T) {
	validTransaction := makeTransaction(t)

	validTxnBody := &InjectTransactionRequest{
		RawTxn: validTransaction.MustSerializeHex(),
	}
	validTxnBodyJSON, err := json.Marshal(validTxnBody)
	require.NoError(t, err)

	validTxnBodyNoBroadcast := &InjectTransactionRequest{
		RawTxn:      validTransaction.MustSerializeHex(),
		NoBroadcast: true,
	}
	validTxnBodyNoBroadcastJSON, err := json.Marshal(validTxnBodyNoBroadcast)
	require.NoError(t, err)

	b := &InjectTransactionRequest{
		RawTxn: hex.EncodeToString(testutil.RandBytes(t, 128)),
	}
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
			name:     "400 - rawtx required",
			method:   http.MethodPost,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - rawtx is required",
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
			err:      "400 Bad Request - Invalid transaction: Not enough buffer data to deserialize",
			httpBody: string(invalidTxnBodyJSON),
		},
		{
			name:                   "503 - daemon.ErrNetworkingDisabled",
			method:                 http.MethodPost,
			status:                 http.StatusServiceUnavailable,
			err:                    "503 Service Unavailable - Networking is disabled",
			httpBody:               string(validTxnBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: daemon.ErrNetworkingDisabled,
		},
		{
			name:                   "503 - gnet.ErrNoReachableConnections",
			method:                 http.MethodPost,
			status:                 http.StatusServiceUnavailable,
			err:                    "503 Service Unavailable - All pool connections are unreachable at this time",
			httpBody:               string(validTxnBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: gnet.ErrNoReachableConnections,
		},
		{
			name:                   "503 - gnet.ErrPoolEmpty",
			method:                 http.MethodPost,
			status:                 http.StatusServiceUnavailable,
			err:                    "503 Service Unavailable - Connection pool is empty after filtering connections",
			httpBody:               string(validTxnBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: gnet.ErrPoolEmpty,
		},
		{
			name:                   "500 - other injectBroadcastTransactionError",
			method:                 http.MethodPost,
			status:                 http.StatusInternalServerError,
			err:                    "500 Internal Server Error - injectBroadcastTransactionError",
			httpBody:               string(validTxnBodyJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: errors.New("injectBroadcastTransactionError"),
		},
		{
			name:                   "500 - no broadcast other injectTransactionError",
			method:                 http.MethodPost,
			status:                 http.StatusInternalServerError,
			err:                    "500 Internal Server Error - injectTransactionError",
			httpBody:               string(validTxnBodyNoBroadcastJSON),
			injectTransactionArg:   validTransaction,
			injectTransactionError: errors.New("injectTransactionError"),
		},
		{
			name:                 "400 - txn constraint violation",
			method:               http.MethodPost,
			status:               http.StatusBadRequest,
			err:                  "400 Bad Request - Transaction violates hard constraint: bad transaction",
			httpBody:             string(validTxnBodyJSON),
			injectTransactionArg: validTransaction,
			injectTransactionError: visor.ErrTxnViolatesHardConstraint{
				Err: errors.New("bad transaction"),
			},
		},
		{
			name:                 "400 - no broadcast txn constraint violation",
			method:               http.MethodPost,
			status:               http.StatusBadRequest,
			err:                  "400 Bad Request - Transaction violates hard constraint: bad transaction",
			httpBody:             string(validTxnBodyNoBroadcastJSON),
			injectTransactionArg: validTransaction,
			injectTransactionError: visor.ErrTxnViolatesHardConstraint{
				Err: errors.New("bad transaction"),
			},
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
			name:                 "200 no broadcast",
			method:               http.MethodPost,
			status:               http.StatusOK,
			httpBody:             string(validTxnBodyNoBroadcastJSON),
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
			gateway.On("InjectTransaction", tc.injectTransactionArg).Return(tc.injectTransactionError)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)

			}

			rr := httptest.NewRecorder()

			handler := newServerMux(defaultMuxConfig(), gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				expectedResponse, err := json.MarshalIndent(tc.httpResponse, "", "    ")
				require.NoError(t, err)
				require.Equal(t, string(expectedResponse), rr.Body.String(), tc.name)
			}
		})
	}
}

func TestResendUnconfirmedTxns(t *testing.T) {
	validHash1 := testutil.RandSHA256(t)
	validHash2 := testutil.RandSHA256(t)

	tt := []struct {
		name                          string
		method                        string
		status                        int
		err                           string
		httpBody                      string
		resendUnconfirmedTxnsResponse []cipher.SHA256
		resendUnconfirmedTxnsErr      error
		httpResponse                  ResendResult
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:                     "500 resend failed network error",
			method:                   http.MethodPost,
			status:                   http.StatusServiceUnavailable,
			err:                      "503 Service Unavailable - All pool connections are unreachable at this time",
			resendUnconfirmedTxnsErr: gnet.ErrNoReachableConnections,
		},

		{
			name:                     "500 resend failed unknown error",
			method:                   http.MethodPost,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error - ResendUnconfirmedTxns failed",
			resendUnconfirmedTxnsErr: errors.New("ResendUnconfirmedTxns failed"),
		},

		{
			name:                          "200",
			method:                        http.MethodPost,
			status:                        http.StatusOK,
			resendUnconfirmedTxnsResponse: nil,
			httpResponse: ResendResult{
				Txids: []string{},
			},
		},

		{
			name:                          "200 with hashes",
			method:                        http.MethodPost,
			status:                        http.StatusOK,
			resendUnconfirmedTxnsResponse: []cipher.SHA256{validHash1, validHash2},
			httpResponse: ResendResult{
				Txids: []string{validHash1.Hex(), validHash2.Hex()},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/resendUnconfirmedTxns"
			gateway := &MockGatewayer{}
			gateway.On("ResendUnconfirmedTxns").Return(tc.resendUnconfirmedTxnsResponse, tc.resendUnconfirmedTxnsErr)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg ResendResult
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestGetRawTxn(t *testing.T) {
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

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
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

	type verboseResult struct {
		Transactions []visor.Transaction
		Inputs       [][]visor.TransactionInput
	}

	tt := []struct {
		name                           string
		method                         string
		status                         int
		err                            string
		httpBody                       *httpBody
		verbose                        bool
		getTransactionsArg             []visor.TxFilter
		getTransactionsResponse        []visor.Transaction
		getTransactionsError           error
		getTransactionsVerboseResponse verboseResult
		getTransactionsVerboseError    error
		httpResponse                   interface{}
	}{
		{
			name:   "405",
			method: http.MethodDelete,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:   "400 - invalid `addrs` param",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - parse parameter: 'addrs' failed: address \"invalid\" is invalid: Invalid base58 character",
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
			getTransactionsResponse: []visor.Transaction{},
			httpResponse:            []readable.TransactionWithStatus{},
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
			getTransactionsVerboseResponse: verboseResult{
				Transactions: []visor.Transaction{},
				Inputs:       [][]visor.TransactionInput{},
			},
			httpResponse: []readable.TransactionWithStatusVerbose{},
		},

		{
			name:   "200 POST",
			method: http.MethodPost,
			status: http.StatusOK,
			httpBody: &httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			getTransactionsArg: []visor.TxFilter{
				visor.NewAddrsFilter(addrs),
				visor.NewConfirmedTxFilter(true),
			},
			getTransactionsResponse: []visor.Transaction{},
			httpResponse:            []readable.TransactionWithStatus{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/transactions"
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

			gateway.On("GetTransactions", matchFunc).Return(tc.getTransactionsResponse, tc.getTransactionsError)
			gateway.On("GetTransactionsWithInputs", matchFunc).Return(tc.getTransactionsVerboseResponse.Transactions,
				tc.getTransactionsVerboseResponse.Inputs, tc.getTransactionsVerboseError)

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

			setCSRFParameters(t, tokenValid, req)

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

			handler := newServerMux(cfg, gateway)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				if tc.verbose {
					var msg []readable.TransactionWithStatusVerbose
					err = json.Unmarshal(rr.Body.Bytes(), &msg)
					require.NoError(t, err)
					require.Equal(t, tc.httpResponse, msg, tc.name)
				} else {
					var msg []readable.TransactionWithStatus
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
	inputs []visor.TransactionInput
}

func newVerifyTxnResponseJSON(t *testing.T, txn *coin.Transaction, inputs []visor.TransactionInput, isTxnConfirmed, isUnsigned bool) VerifyTransactionResponse {
	ctxn, err := newCreatedTransactionFuzzy(txn, inputs)
	require.NoError(t, err)
	return VerifyTransactionResponse{
		Transaction: *ctxn,
		Confirmed:   isTxnConfirmed,
		Unsigned:    isUnsigned,
	}
}

func prepareTxnAndInputs(t *testing.T) transactionAndInputs {
	txn := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	err := txn.PushInput(ux.Hash())
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 5e6, 50)
	require.NoError(t, err)
	txn.SignInputs([]cipher.SecKey{s})
	err = txn.UpdateHeader()
	require.NoError(t, err)

	input, err := visor.NewTransactionInput(ux, uint64(time.Now().UTC().Unix()))
	require.NoError(t, err)

	return transactionAndInputs{
		txn:    txn,
		inputs: []visor.TransactionInput{input},
	}
}

func makeTransactionWithEmptyAddressOutput(t *testing.T) transactionAndInputs {
	txn := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	err := txn.PushInput(ux.Hash())
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.PushOutput(cipher.Address{}, 5e6, 50)
	require.NoError(t, err)
	txn.SignInputs([]cipher.SecKey{s})
	err = txn.UpdateHeader()
	require.NoError(t, err)

	input, err := visor.NewTransactionInput(ux, uint64(time.Now().UTC().Unix()))
	require.NoError(t, err)

	return transactionAndInputs{
		txn:    txn,
		inputs: []visor.TransactionInput{input},
	}
}

func TestVerifyTransaction(t *testing.T) {
	txnAndInputs := prepareTxnAndInputs(t)
	type httpBody struct {
		Unsigned           bool   `json:"unsigned"`
		EncodedTransaction string `json:"encoded_transaction"`
	}

	validTxnBody := &httpBody{
		EncodedTransaction: txnAndInputs.txn.MustSerializeHex(),
	}
	validTxnBodyJSON, err := json.Marshal(validTxnBody)
	require.NoError(t, err)

	b := &httpBody{
		EncodedTransaction: hex.EncodeToString(testutil.RandBytes(t, 128)),
	}
	invalidTxnBodyJSON, err := json.Marshal(b)
	require.NoError(t, err)

	invalidTxnEmptyAddress := makeTransactionWithEmptyAddressOutput(t)
	invalidTxnEmptyAddressBody := &httpBody{
		EncodedTransaction: invalidTxnEmptyAddress.txn.MustSerializeHex(),
	}
	invalidTxnEmptyAddressBodyJSON, err := json.Marshal(invalidTxnEmptyAddressBody)
	require.NoError(t, err)

	unsignedTxnAndInputs := prepareTxnAndInputs(t)
	unsignedTxnAndInputs.txn.Sigs = make([]cipher.Sig, len(unsignedTxnAndInputs.txn.Sigs))
	err = unsignedTxnAndInputs.txn.UpdateHeader()
	require.NoError(t, err)
	unsignedTxnBody := &httpBody{
		EncodedTransaction: unsignedTxnAndInputs.txn.MustSerializeHex(),
	}
	unsignedTxnBodyJSON, err := json.Marshal(unsignedTxnBody)
	require.NoError(t, err)

	unsignedTxnBodyUnsigned := &httpBody{
		Unsigned:           true,
		EncodedTransaction: unsignedTxnAndInputs.txn.MustSerializeHex(),
	}
	unsignedTxnBodyUnsignedJSON, err := json.Marshal(unsignedTxnBodyUnsigned)
	require.NoError(t, err)

	type verifyTxnVerboseResult struct {
		Uxouts         []visor.TransactionInput
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
		gatewayVerifyTxnVerboseSigned visor.TxnSignedFlag
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
			contentType:  ContentTypeJSON,
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
			name:         "400 - encoded_transaction is required",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpBody:     `{"wrongKey":"wrongValue"}`,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "encoded_transaction is required"),
		},
		{
			name:         "400 - encoding/hex: odd length hex string",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpBody:     `{"encoded_transaction":"aab"}`,
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "decode transaction failed: encoding/hex: odd length hex string"),
		},
		{
			name:         "400 - deserialization error",
			method:       http.MethodPost,
			contentType:  ContentTypeJSON,
			status:       http.StatusBadRequest,
			httpBody:     string(invalidTxnBodyJSON),
			httpResponse: NewHTTPErrorResponse(http.StatusBadRequest, "decode transaction failed: Invalid transaction: Not enough buffer data to deserialize"),
		},
		{
			name:                          "422 - txn sends to empty address",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusUnprocessableEntity,
			httpBody:                      string(invalidTxnEmptyAddressBodyJSON),
			gatewayVerifyTxnVerboseArg:    invalidTxnEmptyAddress.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnSigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: invalidTxnEmptyAddress.inputs,
				Err:    visor.NewErrTxnViolatesUserConstraint(errors.New("Transaction.Out contains an output sending to an empty address")),
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &invalidTxnEmptyAddress.txn, invalidTxnEmptyAddress.inputs, false, false),
				Error: &HTTPError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Transaction violates user constraint: Transaction.Out contains an output sending to an empty address",
				},
			},
		},
		{
			name:                          "422 - txn is unsigned",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusUnprocessableEntity,
			httpBody:                      string(unsignedTxnBodyJSON),
			gatewayVerifyTxnVerboseArg:    unsignedTxnAndInputs.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnSigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: unsignedTxnAndInputs.inputs,
				Err:    visor.NewErrTxnViolatesUserConstraint(errors.New("Transaction.Out contains an output sending to an empty address")),
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &unsignedTxnAndInputs.txn, unsignedTxnAndInputs.inputs, false, true),
				Error: &HTTPError{
					Code:    http.StatusUnprocessableEntity,
					Message: "Transaction violates user constraint: Transaction.Out contains an output sending to an empty address",
				},
			},
		},
		{
			name:                          "500 - internal server error",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusInternalServerError,
			httpBody:                      string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg:    txnAndInputs.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnSigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Err: errors.New("verify transaction failed"),
			},
			httpResponse: NewHTTPErrorResponse(http.StatusInternalServerError, "verify transaction failed"),
		},
		{
			name:                          "422 - txn is confirmed",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusUnprocessableEntity,
			httpBody:                      string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg:    txnAndInputs.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnSigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts:         txnAndInputs.inputs,
				IsTxnConfirmed: true,
			},
			httpResponse: HTTPResponse{
				Error: &HTTPError{
					Message: "transaction has been spent",
					Code:    http.StatusUnprocessableEntity,
				},
				Data: newVerifyTxnResponseJSON(t, &txnAndInputs.txn, txnAndInputs.inputs, true, false),
			},
		},
		{
			name:                          "200 - unsigned",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusOK,
			httpBody:                      string(unsignedTxnBodyUnsignedJSON),
			gatewayVerifyTxnVerboseArg:    unsignedTxnAndInputs.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnUnsigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: unsignedTxnAndInputs.inputs,
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &unsignedTxnAndInputs.txn, unsignedTxnAndInputs.inputs, false, true),
			},
		},
		{
			name:                          "200",
			method:                        http.MethodPost,
			contentType:                   ContentTypeJSON,
			status:                        http.StatusOK,
			httpBody:                      string(validTxnBodyJSON),
			gatewayVerifyTxnVerboseArg:    txnAndInputs.txn,
			gatewayVerifyTxnVerboseSigned: visor.TxnSigned,
			gatewayVerifyTxnVerboseResult: verifyTxnVerboseResult{
				Uxouts: txnAndInputs.inputs,
			},
			httpResponse: HTTPResponse{
				Data: newVerifyTxnResponseJSON(t, &txnAndInputs.txn, txnAndInputs.inputs, false, false),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v2/transaction/verify"
			gateway := &MockGatewayer{}
			gateway.On("VerifyTxnVerbose", &tc.gatewayVerifyTxnVerboseArg, tc.gatewayVerifyTxnVerboseSigned).Return(tc.gatewayVerifyTxnVerboseResult.Uxouts,
				tc.gatewayVerifyTxnVerboseResult.IsTxnConfirmed, tc.gatewayVerifyTxnVerboseResult.Err)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(tc.httpBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", tc.contentType)

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			rr := httptest.NewRecorder()

			cfg := defaultMuxConfig()
			cfg.disableCSRF = false

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

				var txnRsp VerifyTransactionResponse
				err := json.Unmarshal(rsp.Data, &txnRsp)
				require.NoError(t, err)

				require.Equal(t, tc.httpResponse.Data.(VerifyTransactionResponse), txnRsp)
			}
		})
	}
}
