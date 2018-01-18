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

	"github.com/pkg/errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor"
	"bytes"
)

// GetAllUnconfirmedTxns returns all unconfirmed transactions
func (gw FakeGateway) GetAllUnconfirmedTxns() []visor.UnconfirmedTxn {
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

func createUnconfirmedTxn(t *testing.T) visor.UnconfirmedTxn {
	ut := visor.UnconfirmedTxn{}
	ut.Txn = coin.Transaction{}
	ut.Txn.InnerHash = testutil.RandSHA256(t)
	ut.Received = utc.Now().UnixNano()
	ut.Checked = ut.Received
	ut.Announced = time.Time{}.UnixNano()
	return ut
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
	}
}

func TestInjectTransaction(t *testing.T) {
	validTransaction := testutil.MakeTransaction(t)
	type httpBody struct {
		Rawtx string `json:"rawtx"`
	}
	b := &httpBody{Rawtx: string(validTransaction.Serialize())}
	body, err := json.Marshal(b)
	require.NoError(t, err)
	tt := []struct {
		name                   string
		method                 string
		url                    string
		status                 int
		err                    string
		httpBody               string
		injectTransactionArg   cipher.SHA256
		injectTransactionError error
		httpResponse           visor.TransactionResult
	}{
		{
			"405",
			http.MethodGet,
			"/injectTransaction",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			testutil.RandSHA256(t),
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - EOF",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - EOF",
			"",
			testutil.RandSHA256(t),
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - Invalid transaction: Deserialization failed",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - Invalid transaction: Deserialization failed",
			`{"wrongKey":"wrongValue"}`,
			testutil.RandSHA256(t),
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - encoding/hex: odd length hex string",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			`{"rawtx":"aab"}`,
			testutil.RandSHA256(t),
			nil,
			visor.TransactionResult{},
		},
		{
			"400 - unknown",
			http.MethodPost,
			"/injectTransaction",
			http.StatusBadRequest,
			"400 Bad Request - encoding/hex: odd length hex string",
			string(body),
			testutil.RandSHA256(t),
			nil,
			visor.TransactionResult{},
		},
	}

	for _, tc := range tt {
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
			var msg visor.TransactionResult
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}
