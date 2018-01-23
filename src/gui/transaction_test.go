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

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw *FakeGateway) GetAllUnconfirmedTxns() []visor.UnconfirmedTxn {
	args := gw.Called()
	return args.Get(0).([]visor.UnconfirmedTxn)
}

// GetTransaction returns transaction by txid
func (gw *FakeGateway) GetTransaction(txid cipher.SHA256) (tx *visor.Transaction, err error) {
	args := gw.Called(txid)
	return args.Get(0).(*visor.Transaction), args.Error(1)
}

// InjectTransaction injects transaction
func (gw *FakeGateway) InjectTransaction(txn coin.Transaction) error {
	args := gw.Called(txn)
	return args.Error(0)
}

// ResendUnconfirmedTxns resents all unconfirmed transactions
func (gw *FakeGateway) ResendUnconfirmedTxns() (rlt *daemon.ResendResult) {
	args := gw.Called()
	return args.Get(0).(*daemon.ResendResult)
}

func (gw *FakeGateway) GetTransactions(flts ...visor.TxFilter) ([]visor.Transaction, error) {
	args := gw.Called(flts)
	return args.Get(0).([]visor.Transaction), args.Error(1)
}

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
			"405",
			http.MethodPost,
			"/pendingTxs",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			[]visor.UnconfirmedTxn{},
			nil,
		},
		{
			"500 - bad unconfirmedTxn",
			http.MethodGet,
			"/pendingTxs",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			[]visor.UnconfirmedTxn{
				invalidTxn,
			},
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/pendingTxs",
			http.StatusOK,
			"",
			[]visor.UnconfirmedTxn{},
			[]*visor.ReadableUnconfirmedTxn{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetAllUnconfirmedTxns").Return(tc.getAllUnconfirmedTxnsResponse)

			req, err := http.NewRequest(tc.method, tc.url, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getPendingTxs(gateway))

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
	oddHash := "caicb"
	invalidHash := "cabrca"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	type httpBody struct {
		txid string
	}

	tt := []struct {
		name                  string
		method                string
		url                   string
		status                int
		err                   string
		httpBody              *httpBody
		getTransactionArg     cipher.SHA256
		getTransactionReponse *visor.Transaction
		getTransactionError   error
		httpResponse          visor.TransactionResult
	}{
		{
			"405",
			http.MethodPost,
			"/transaction",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			nil,
			testutil.RandSHA256(t),
			nil,
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - empty txid",
			http.MethodGet,
			"/transaction",
			http.StatusBadRequest,
			"400 Bad Request - txid is empty",
			&httpBody{
				txid: "",
			},
			testutil.RandSHA256(t),
			nil,
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - invalid hash: odd length hex string",
			http.MethodGet,
			"/transaction",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			&httpBody{
				txid: oddHash,
			},
			testutil.RandSHA256(t),
			nil,
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - invalid hash: invalid byte: U+0072 'r'",
			http.MethodGet,
			"/transaction",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			&httpBody{
				txid: invalidHash,
			},
			testutil.RandSHA256(t),
			nil,
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - getTransactionError",
			http.MethodGet,
			"/transaction",
			http.StatusBadRequest,
			"400 Bad Request - getTransactionError",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			nil,
			errors.New("getTransactionError"),
			visor.TransactionResult{},
		},
		{
			"404",
			http.MethodGet,
			"/transaction",
			http.StatusNotFound,
			"404 Not Found",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			nil,
			nil,
			visor.TransactionResult{},
		},
		{
			"200",
			http.MethodGet,
			"/transaction",
			http.StatusOK,
			"",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			&visor.Transaction{},
			nil,
			visor.TransactionResult{
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
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetTransaction", tc.getTransactionArg).Return(tc.getTransactionReponse, tc.getTransactionError)

			v := url.Values{}
			urlFull := tc.url
			if tc.httpBody != nil {
				if tc.httpBody.txid != "" {
					v.Add("txid", tc.httpBody.txid)
				}
			}
			if len(v) > 0 {
				urlFull += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, urlFull, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getTransactionByID(gateway))

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
	validTxBodyJson, err := json.Marshal(validTxBody)
	require.NoError(t, err)
	b := &httpBody{Rawtx: hex.EncodeToString(testutil.RandBytes(t, 128))}
	invalidTxBodyJson, err := json.Marshal(b)
	require.NoError(t, err)
	tt := []struct {
		name                   string
		method                 string
		url                    string
		status                 int
		err                    string
		httpBody               string
		injectTransactionArg   coin.Transaction
		injectTransactionError error
		httpResponse           string
	}{
		{
			"405",
			http.MethodGet,
			"/injectTransaction",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			validTransaction,
			nil,
			"",
		},
		{
			"400 - EOF",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - EOF",
			"",
			validTransaction,
			nil,
			"",
		},
		{
			"400 - Invalid transaction: Deserialization failed",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - Invalid transaction: Deserialization failed",
			`{"wrongKey":"wrongValue"}`,
			validTransaction,
			nil,
			"",
		},
		{
			"400 - encoding/hex: odd length hex string",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			`{"rawtx":"aab"}`,
			validTransaction,
			nil,
			"",
		},
		{
			"400 - rawtx deserialization error",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - Invalid transaction: Deserialization failed",
			string(invalidTxBodyJson),
			validTransaction,
			nil,
			"",
		},
		{
			"400 - injectTransactionError",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - inject tx failed:injectTransactionError",
			string(validTxBodyJson),
			validTransaction,
			errors.New("injectTransactionError"),
			"",
		},
		{
			"200",
			http.MethodPost,
			"/injectTransaction",
			http.StatusOK,
			"",
			string(validTxBodyJson),
			validTransaction,
			nil,
			validTransaction.Hash().Hex(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("InjectTransaction", tc.injectTransactionArg).Return(tc.injectTransactionError)

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(injectTransaction(gateway))
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
		url                           string
		status                        int
		err                           string
		httpBody                      string
		resendUnconfirmedTxnsResponse *daemon.ResendResult
		httpResponse                  *daemon.ResendResult
	}{
		{
			"405",
			http.MethodPost,
			"/resendUnconfirmedTxns",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			&daemon.ResendResult{},
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/resendUnconfirmedTxns",
			http.StatusOK,
			"",
			"",
			&daemon.ResendResult{},
			&daemon.ResendResult{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("ResendUnconfirmedTxns").Return(tc.resendUnconfirmedTxnsResponse)

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(tc.httpBody))
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resendUnconfirmedTxns(gateway))
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
	oddHash := "caicb"
	invalidHash := "cabrca"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	type httpBody struct {
		txid string
	}
	tt := []struct {
		name             string
		method           string
		url              string
		status           int
		err              string
		httpBody         *httpBody
		getRawTxArg      cipher.SHA256
		getRawTxResponse *visor.Transaction
		getRawTxError    error
		httpResponse     string
	}{
		{
			"405",
			http.MethodPost,
			"/rawtx",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			nil,
			testutil.RandSHA256(t),
			nil,
			nil,
			"",
		},
		{
			"400 - txid is empty",
			http.MethodGet,
			"/rawtx",
			http.StatusBadRequest,
			"400 Bad Request - txid is empty",
			&httpBody{},
			testutil.RandSHA256(t),
			nil,
			nil,
			"",
		},
		{
			"400 - invalid hash: odd length hex string",
			http.MethodGet,
			"/rawtx",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			&httpBody{
				txid: oddHash,
			},
			testutil.RandSHA256(t),
			nil,
			nil,
			"",
		},
		{
			"400 - invalid hash: invalid byte: U+0072 'r'",
			http.MethodGet,
			"/rawtx",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: invalid byte: U+0072 'r'",
			&httpBody{
				txid: invalidHash,
			},
			testutil.RandSHA256(t),
			nil,
			nil,
			"",
		},
		{
			"400 - getTransactionError",
			http.MethodGet,
			"/rawtx",
			http.StatusBadRequest,
			"400 Bad Request - getTransactionError",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			nil,
			errors.New("getTransactionError"),
			"",
		},
		{
			"400 - getTransactionError",
			http.MethodGet,
			"/rawtx",
			http.StatusBadRequest,
			"400 Bad Request - getTransactionError",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			nil,
			errors.New("getTransactionError"),
			"",
		},
		{
			"404",
			http.MethodGet,
			"/rawtx",
			http.StatusNotFound,
			"404 Not Found",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			nil,
			nil,
			"",
		},
		{
			"200",
			http.MethodGet,
			"/rawtx",
			http.StatusOK,
			"",
			&httpBody{
				txid: validHash,
			},
			testutil.SHA256FromHex(t, validHash),
			&visor.Transaction{},
			nil,
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetTransaction", tc.getRawTxArg).Return(tc.getRawTxResponse, tc.getRawTxError)
			v := url.Values{}
			urlFull := tc.url
			if tc.httpBody != nil {
				if tc.httpBody.txid != "" {
					v.Add("txid", tc.httpBody.txid)
				}
			}
			if len(v) > 0 {
				urlFull += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, urlFull, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getRawTx(gateway))
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
		url                     string
		status                  int
		err                     string
		httpBody                *httpBody
		getTransactionsArg      []visor.TxFilter
		getTransactionsResponse []visor.Transaction
		getTransactionsError    error
		httpResponse            []visor.Transaction
	}{
		{
			"405",
			http.MethodPost,
			"/transactions",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			nil,
			[]visor.TxFilter{},
			nil,
			nil,
			[]visor.Transaction{},
		},
		{
			"400 - invalid `addrs` param",
			http.MethodGet,
			"/transactions",
			http.StatusBadRequest,
			"400 Bad Request - parse parament: 'addrs' failed: Invalid base58 character",
			&httpBody{
				addrs: invalidAddrsStr,
			},
			[]visor.TxFilter{
				visor.AddrsFilter(addrs),
			},
			nil,
			nil,
			[]visor.Transaction{},
		},
		{
			"400 - invalid `confirmed` param",
			http.MethodGet,
			"/transactions",
			http.StatusBadRequest,
			"400 Bad Request - invalid 'confirmed' value: strconv.ParseBool: parsing \"invalidConfirmed\": invalid syntax",
			&httpBody{
				addrs:     addrsStr,
				confirmed: "invalidConfirmed",
			},
			[]visor.TxFilter{
				visor.AddrsFilter(addrs),
			},
			nil,
			nil,
			[]visor.Transaction{},
		},
		{
			"500 - getTransactionsError",
			http.MethodGet,
			"/transactions",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			&httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			[]visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			[]visor.Transaction{},
			errors.New("getTransactionsError"),
			[]visor.Transaction{},
		},
		{
			"500 - visor.NewTransactionResults error",
			http.MethodGet,
			"/transactions",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			&httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			[]visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			[]visor.Transaction{
				{
					Txn: invalidTxn,
					Status: visor.TransactionStatus{
						Confirmed: true,
						Height:    103,
					},
				},
			},
			nil,
			[]visor.Transaction{},
		},
		{
			"200",
			http.MethodGet,
			"/transactions",
			http.StatusOK,
			"",
			&httpBody{
				addrs:     addrsStr,
				confirmed: "true",
			},
			[]visor.TxFilter{
				visor.AddrsFilter(addrs),
				visor.ConfirmedTxFilter(true),
			},
			[]visor.Transaction{},
			nil,
			[]visor.Transaction{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetTransactions", mock.Anything).Return(tc.getTransactionsResponse, tc.getTransactionsError)

			v := url.Values{}
			urlFull := tc.url
			if tc.httpBody != nil {
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
				if tc.httpBody.confirmed != "" {
					v.Add("confirmed", tc.httpBody.confirmed)
				}
			}
			if len(v) > 0 {
				urlFull += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, urlFull, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getTransactions(gateway))

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
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
