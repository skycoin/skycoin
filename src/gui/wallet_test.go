package gui

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

const configuredHost = "127.0.0.1:6420"

var mxConfig = muxConfig{host: configuredHost, appLoc: "."}

func TestWalletSpendHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
		Password string
	}

	tt := []struct {
		name                          string
		method                        string
		body                          *httpBody
		status                        int
		err                           string
		walletID                      string
		coins                         uint64
		dst                           string
		password                      string
		gatewaySpendResult            *coin.Transaction
		gatewaySpendErr               error
		gatewayGetWalletBalanceResult wallet.BalancePair
		gatewayBalanceErr             error
		spendResult                   *SpendResult
		csrfDisabled                  bool
	}{
		{
			name:     "405",
			method:   http.MethodGet,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodPost,
			body:     &httpBody{},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "0",
		},
		{
			name:   "400 - no dst",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing destination address \"dst\"",
			walletID: "0",
		},
		{
			name:   "400 - bad dst addr",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      " 2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid destination address: Invalid base58 character",
			walletID: "0",
		},
		{
			name:   "400 - no coins",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - coins is string",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "foo",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - coins is negative value",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "-123",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - zero coins",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "0",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value, must > 0",
			walletID: "0",
		},
		{
			name:   "400 - gw spend error txn no fee",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - Transaction has zero coinhour fee",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: fee.ErrTxnNoFee,
			spendResult: &SpendResult{
				Error: fee.ErrTxnNoFee.Error(),
			},
		},
		{
			name:   "400 - gw spend error spending unconfirmed",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - please spend after your pending transaction is confirmed",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrSpendingUnconfirmed,
			spendResult: &SpendResult{
				Error: wallet.ErrSpendingUnconfirmed.Error(),
			},
		},
		{
			name:   "400 - gw spend error insufficient balance",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - balance is not sufficient",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrInsufficientBalance,
			spendResult: &SpendResult{
				Error: wallet.ErrInsufficientBalance.Error(),
			},
		},
		{
			name:   "404 - gw spend error wallet not exist",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusNotFound,
			err:             "404 Not Found",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrWalletNotExist,
			spendResult: &SpendResult{
				Error: wallet.ErrWalletNotExist.Error(),
			},
		},
		{
			name:   "500 - gw spend error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusInternalServerError,
			err:             "500 Internal Server Error - Spend error",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: errors.New("Spend error"),
			spendResult: &SpendResult{
				Error: "Spend error",
			},
		},
		{
			name:   "200 - gw GetWalletBalance error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:             http.StatusOK,
			walletID:           "1234",
			coins:              12,
			dst:                "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{},
			gatewayBalanceErr:  errors.New("GetWalletBalance error"),
			spendResult: &SpendResult{
				Error: "Get wallet balance failed: GetWalletBalance error",
				Transaction: &visor.ReadableTransaction{
					Sigs:      []string{},
					In:        []string{},
					Out:       []visor.ReadableTransactionOutput{},
					Hash:      "78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrWalletAPIDisabled,
			spendResult: &SpendResult{
				Error: wallet.ErrWalletAPIDisabled.Error(),
			},
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:             http.StatusOK,
			walletID:           "1234",
			coins:              12,
			dst:                "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{},
			spendResult: &SpendResult{
				Balance: &wallet.BalancePair{},
				Transaction: &visor.ReadableTransaction{
					Length:    0,
					Type:      0,
					Hash:      "78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Timestamp: 0,
					Sigs:      []string{},
					In:        []string{},
					Out:       []visor.ReadableTransactionOutput{},
				},
			},
		},
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:             http.StatusOK,
			walletID:           "1234",
			coins:              12,
			dst:                "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{},
			spendResult: &SpendResult{
				Balance: &wallet.BalancePair{},
				Transaction: &visor.ReadableTransaction{
					Length:    0,
					Type:      0,
					Hash:      "78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Timestamp: 0,
					Sigs:      []string{},
					In:        []string{},
					Out:       []visor.ReadableTransactionOutput{},
				},
			},
			csrfDisabled: true,
		},
		{
			name:   "400 - missing password",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
			},
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrMissingPassword,
			err:             "400 Bad Request - missing password",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrMissingPassword.Error(),
			},
		},
		{
			name:   "400 - invalid password",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
				Password: "pwd",
			},
			password:        "pwd",
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrInvalidPassword,
			err:             "400 Bad Request - invalid password",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrInvalidPassword.Error(),
			},
		},
		{
			name:   "400 - wallet is encrypted",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
			},
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrWalletEncrypted,
			err:             "400 Bad Request - wallet is encrypted",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrWalletEncrypted.Error(),
			},
		},
		{
			name:   "400 - wallet is not encrypted",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
				Password: "pwd",
			},
			password:        "pwd",
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrWalletNotEncrypted,
			err:             "400 Bad Request - wallet is not encrypted",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrWalletNotEncrypted.Error(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewaySpendResult == nil {
				tc.gatewaySpendResult = &coin.Transaction{}
			}

			gateway := &GatewayerMock{}
			addr, _ := cipher.DecodeBase58Address(tc.dst)
			gateway.On("Spend", tc.walletID, []byte(tc.password), tc.coins, addr).Return(tc.gatewaySpendResult, tc.gatewaySpendErr)
			gateway.On("GetWalletBalance", tc.walletID).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)

			endpoint := "/wallet/spend"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
				if tc.body.Dst != "" {
					v.Add("dst", tc.body.Dst)
				}
				if tc.body.Coins != "" {
					v.Add("coins", tc.body.Coins)
				}
				if tc.body.Password != "" {
					v.Add("password", tc.body.Password)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
			} else {
				var msg SpendResult
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, *tc.spendResult, msg)
			}
		})
	}
}

func TestWalletGet(t *testing.T) {
	entries, resEntries := makeEntries([]byte("seed"), 5)
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
	}

	tt := []struct {
		name                   string
		method                 string
		url                    string
		body                   *httpBody
		status                 int
		err                    string
		walletID               string
		gatewayGetWalletResult wallet.Wallet
		responseBody           WalletResponse
		gatewayGetWalletErr    error
	}{
		{
			name:     "405",
			method:   http.MethodPost,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "",
		},
		{
			name:   "400 - error from the `gateway.GetWallet(wltID)`",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "123",
			},
			status:              http.StatusBadRequest,
			err:                 "400 Bad Request - wallet 123 doesn't exist",
			walletID:            "123",
			gatewayGetWalletErr: errors.New("wallet 123 doesn't exist"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "1234",
			},
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "1234",
			gatewayGetWalletResult: wallet.Wallet{
				Meta:    map[string]string{"seed": "seed", "lastSeed": "seed"},
				Entries: []wallet.Entry{},
			},
			gatewayGetWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "1234",
			},
			status:   http.StatusOK,
			walletID: "1234",
			gatewayGetWalletResult: wallet.Wallet{
				Meta:    map[string]string{"seed": "seed", "lastSeed": "seed"},
				Entries: cloneEntries(entries),
			},
			responseBody: WalletResponse{Entries: resEntries[:]},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("GetWallet", tc.walletID).Return(&tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)
			v := url.Values{}

			endpoint := "/wallet"

			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
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
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt WalletResponse
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestWalletBalanceHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
	}

	tt := []struct {
		name                          string
		method                        string
		body                          *httpBody
		status                        int
		err                           string
		walletID                      string
		gatewayGetWalletBalanceResult wallet.BalancePair
		gatewayBalanceErr             error
		result                        *wallet.BalancePair
	}{
		{
			name:     "405",
			method:   http.MethodPost,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "0",
		},
		{
			name:   "404 - gw `wallet doesn't exist` error",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "notFoundId",
			},
			status:   http.StatusNotFound,
			err:      "404 Not Found",
			walletID: "notFoundId",
			gatewayGetWalletBalanceResult: wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
			gatewayBalanceErr: wallet.ErrWalletNotExist,
			result: &wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			name:   "500 - gw other error",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "someId",
			},
			status:   http.StatusInternalServerError,
			err:      "500 Internal Server Error - gatewayBalanceError",
			walletID: "someId",
			gatewayGetWalletBalanceResult: wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
			gatewayBalanceErr: errors.New("gatewayBalanceError"),
			result: &wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:                        http.StatusForbidden,
			err:                           "403 Forbidden",
			walletID:                      "foo",
			gatewayGetWalletBalanceResult: wallet.BalancePair{},
			gatewayBalanceErr:             wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusOK,
			err:      "",
			walletID: "foo",
			result:   &wallet.BalancePair{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("GetWalletBalance", tc.walletID).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)

			endpoint := "/wallet/balance"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}
			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)
			if status != tc.status {
				t.Errorf("case: %s, handler returned wrong status code: got `%v` want `%v`",
					tc.name, status, tc.status)
			}
			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg wallet.BalancePair
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, &msg, tc.name)
			}
		})
	}
}

func TestUpdateWalletLabelHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Label    string
	}

	tt := []struct {
		name                        string
		method                      string
		url                         string
		body                        *httpBody
		status                      int
		err                         string
		walletID                    string
		label                       string
		gatewayUpdateWalletLabelErr error
		responseBody                string
	}{
		{
			name:   "405",
			method: http.MethodGet,
			body:   &httpBody{},
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			body:   &httpBody{},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:   "400 - missing label",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing label",
			walletID: "foo",
		},
		{
			name:   "404 - gateway.UpdateWalletLabel ErrWalletNotExist",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:   http.StatusNotFound,
			err:      "404 Not Found",
			walletID: "foo",
			label:    "label",
			gatewayUpdateWalletLabelErr: wallet.ErrWalletNotExist,
		},
		{
			name:   "500 - gateway.UpdateWalletLabel error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:   http.StatusInternalServerError,
			err:      "500 Internal Server Error - gateway.UpdateWalletLabel error",
			walletID: "foo",
			label:    "label",
			gatewayUpdateWalletLabelErr: errors.New("gateway.UpdateWalletLabel error"),
		},
		{
			name:   "403 Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "foo",
			label:    "label",
			gatewayUpdateWalletLabelErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 OK",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:   http.StatusOK,
			err:      "",
			walletID: "foo",
			label:    "label",
			gatewayUpdateWalletLabelErr: nil,
			responseBody:                "\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("UpdateWalletLabel", tc.walletID, tc.label).Return(tc.gatewayUpdateWalletLabelErr)

			endpoint := "/wallet/update"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
				if tc.body.Label != "" {
					v.Add("label", tc.body.Label)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				require.Equal(t, tc.responseBody, rr.Body.String(), tc.name)
			}
		})
	}
}

func TestWalletTransactionsHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
	}

	unconfirmedTxn, _ := visor.NewReadableUnconfirmedTxn(&visor.UnconfirmedTxn{})
	tt := []struct {
		name                                  string
		method                                string
		body                                  *httpBody
		status                                int
		err                                   string
		walletID                              string
		gatewayGetWalletUnconfirmedTxnsResult []visor.UnconfirmedTxn
		gatewayGetWalletUnconfirmedTxnsErr    error
		responseBody                          UnconfirmedTxnsResponse
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:   "500 - gateway.GetWalletUnconfirmedTxns error",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusInternalServerError,
			err:      "500 Internal Server Error - gateway.GetWalletUnconfirmedTxns error",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsErr: errors.New("gateway.GetWalletUnconfirmedTxns error"),
		},
		{
			name:   "404 - wallet doesn't exist",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusNotFound,
			err:      "404 Not Found",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsErr: wallet.ErrWalletNotExist,
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusOK,
			err:      "",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsResult: make([]visor.UnconfirmedTxn, 1),
			responseBody:                          UnconfirmedTxnsResponse{Transactions: []visor.ReadableUnconfirmedTxn{*unconfirmedTxn}},
		},
	}

	for _, tc := range tt {
		gateway := &GatewayerMock{}
		gateway.On("GetWalletUnconfirmedTxns", tc.walletID).Return(tc.gatewayGetWalletUnconfirmedTxnsResult, tc.gatewayGetWalletUnconfirmedTxnsErr)

		endpoint := "/wallet/transactions"

		v := url.Values{}
		if tc.body != nil {
			if tc.body.WalletID != "" {
				v.Add("id", tc.body.WalletID)
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
		handler := newServerMux(mxConfig, gateway, csrfStore)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg UnconfirmedTxnsResponse
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			// require.Equal on whole response might result in flaky tests as there is a time field attached to unconfirmed txn response
			require.IsType(t, msg, tc.responseBody)
			require.Len(t, msg.Transactions, 1)
			require.Equal(t, msg.Transactions[0].Txn, tc.responseBody.Transactions[0].Txn)
		}
	}
}

func TestWalletCreateHandler(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type httpBody struct {
		Seed     string
		Label    string
		ScanN    string
		Encrypt  bool
		Password string
	}
	tt := []struct {
		name                      string
		method                    string
		body                      *httpBody
		status                    int
		err                       string
		wltName                   string
		options                   wallet.Options
		gatewayCreateWalletResult wallet.Wallet
		gatewayCreateWalletErr    error
		scanWalletAddressesResult wallet.Wallet
		scanWalletAddressesError  error
		responseBody              WalletResponse
		csrfDisabled              bool
	}{
		{
			name:    "405",
			method:  http.MethodGet,
			status:  http.StatusMethodNotAllowed,
			err:     "405 Method Not Allowed",
			wltName: "foo",
		},
		{
			name:    "400 - missing seed",
			method:  http.MethodPost,
			body:    &httpBody{},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - missing seed",
			wltName: "foo",
		},
		{
			name:   "400 - missing label",
			method: http.MethodPost,
			body: &httpBody{
				Seed: "foo",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - missing label",
			wltName: "foo",
		},
		{
			name:   "400 - invalid scan value",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "bad scanN",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - invalid scan value",
			wltName: "foo",
		},
		{
			name:   "400 - scan must be > 0",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "0",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - scan must be > 0",
			wltName: "foo",
		},
		{
			name:   "400 - gateway.CreateWallet error",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "1",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - gateway.CreateWallet error",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
			},
			gatewayCreateWalletErr: errors.New("gateway.CreateWallet error"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusForbidden,
			err:     "403 Forbidden",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
				Entries: cloneEntries(entries),
			},
			scanWalletAddressesResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
				Entries: cloneEntries(entries),
			},
			responseBody: WalletResponse{
				Meta: WalletMeta{
					Filename: "filename",
				},
				Entries: responseEntries[:],
			},
		},
		// CSRF Tests
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			scanWalletAddressesResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			responseBody: WalletResponse{
				Meta: WalletMeta{
					Filename: "filename",
				},
			},
			csrfDisabled: true,
		},
		{
			name:   "200 - OK - Encrypted",
			method: http.MethodPost,
			body: &httpBody{
				Seed:     "foo",
				Label:    "bar",
				Encrypt:  true,
				Password: "pwd",
				ScanN:    "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename":  "filename",
					"label":     "bar",
					"encrypted": "true",
					"secrets":   "secrets",
				},
			},
			scanWalletAddressesResult: wallet.Wallet{
				Meta: map[string]string{
					"filename":  "filename",
					"label":     "bar",
					"encrypted": "true",
					"secrets":   "secrets",
				},
			},
			responseBody: WalletResponse{
				Meta: WalletMeta{
					Filename:  "filename",
					Label:     "bar",
					Encrypted: true,
				},
			},
		},
		{
			name:   "400 Bad request - encrypt without password",
			method: http.MethodPost,
			body: &httpBody{
				Seed:    "foo",
				Label:   "bar",
				Encrypt: true,
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing password",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			if tc.options.ScanN == 0 {
				tc.options.ScanN = 1
			}
			gateway.On("CreateWallet", "", tc.options).Return(&tc.gatewayCreateWalletResult, tc.gatewayCreateWalletErr)
			// gateway.On("ScanAheadWalletAddresses", tc.wltName, tc.options.Password, tc.scnN-1).Return(&tc.scanWalletAddressesResult, tc.scanWalletAddressesError)

			endpoint := "/wallet/create"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Seed != "" {
					v.Add("seed", tc.body.Seed)
				}
				if tc.body.Label != "" {
					v.Add("label", tc.body.Label)
				}
				if tc.body.ScanN != "" {
					v.Add("scan", tc.body.ScanN)
				}

				if tc.body.Encrypt {
					v.Add("encrypt", strconv.FormatBool(tc.body.Encrypt))
				}

				if tc.body.Password != "" {
					v.Add("password", tc.body.Password)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
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
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg WalletResponse
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, msg, tc.name)
			}

		})
	}
}

func TestWalletNewSeed(t *testing.T) {
	type httpBody struct {
		Entropy string
	}
	tt := []struct {
		name      string
		method    string
		body      *httpBody
		status    int
		err       string
		entropy   string
		resultLen int
	}{
		{
			name:   "405",
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - invalid entropy type",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "xx",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - invalid entropy",
			entropy: "xx",
		},
		{
			name:   "400 - `wrong entropy length` error",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "200",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - entropy length must be 128 or 256",
			entropy: "200",
		},
		{
			name:      "200 - OK with no entropy",
			method:    http.MethodGet,
			body:      &httpBody{},
			status:    http.StatusOK,
			entropy:   "128",
			resultLen: 12,
		},
		{
			name:   "200 - OK | 12 word seed",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "128",
			},
			status:    http.StatusOK,
			entropy:   "128",
			resultLen: 12,
		},
		{
			name:   "200 - OK | 24 word seed",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "256",
			},
			status:    http.StatusOK,
			entropy:   "256",
			resultLen: 24,
		},
	}

	// Loop over each test case
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("IsWalletAPIEnabled").Return(true)

			endpoint := "/wallet/newSeed"

			// Add request parameters to url
			v := url.Values{}
			if tc.body != nil {
				if tc.body.Entropy != "" {
					v.Add("entropy", tc.body.Entropy)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` expected `%v`", tc.name, status, tc.status)
			if status != tc.status {
				t.Errorf("case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)
			}
			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, expected `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg struct {
					Seed string `json:"seed"`
				}
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				// check that expected length is equal to response length
				require.Equal(t, tc.resultLen, len(strings.Fields(msg.Seed)), tc.name)
			}
		})
	}
}

func TestGetWalletSeed(t *testing.T) {
	type gatewayReturnPair struct {
		seed string
		err  error
	}

	tt := []struct {
		name              string
		method            string
		wltID             string
		password          string
		gatewayReturnArgs []interface{}
		expectStatus      int
		expectSeed        string
		expectErr         string
	}{
		{
			name:     "200 - OK",
			method:   http.MethodGet,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"seed",
				nil,
			},
			expectStatus: http.StatusOK,
			expectSeed:   "seed",
		},
		{
			name:              "400 - missing wallet id ",
			method:            http.MethodGet,
			wltID:             "",
			password:          "pwd",
			gatewayReturnArgs: []interface{}{},
			expectStatus:      http.StatusBadRequest,
			expectErr:         "400 Bad Request - missing wallet id",
		},
		{
			name:              "400 - missing password",
			method:            http.MethodGet,
			wltID:             "wallet.wlt",
			password:          "",
			gatewayReturnArgs: []interface{}{},
			expectStatus:      http.StatusBadRequest,
			expectErr:         "400 Bad Request - missing password",
		},
		{
			name:     "400 - invalid password",
			method:   http.MethodGet,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				nil,
				wallet.ErrInvalidPassword,
			},
			expectStatus: http.StatusBadRequest,
			expectErr:    "400 Bad Request - invalid password",
		},
		{
			name:     "403 - wallet not encrypted",
			method:   http.MethodGet,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				nil,
				wallet.ErrWalletNotEncrypted,
			},
			expectStatus: http.StatusForbidden,
			expectErr:    "403 Forbidden",
		},
		{
			name:     "404 - wallet does not exist",
			method:   http.MethodGet,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				nil,
				wallet.ErrWalletNotExist,
			},
			expectStatus: http.StatusNotFound,
			expectErr:    "404 Not Found",
		},
		{
			name:         "405 - Method Not Allowed",
			method:       http.MethodPost,
			expectStatus: http.StatusMethodNotAllowed,
			expectErr:    "405 Method Not Allowed",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			gateway.On("GetWalletSeed", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturnArgs...)

			v := url.Values{}
			v.Add("id", tc.wltID)
			if len(tc.password) > 0 {
				v.Add("password", tc.password)
			}
			endpoint := "/wallet/seed?" + v.Encode()

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.expectStatus, status)

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
			} else {
				var r struct {
					Seed string `json:"seed"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &r)
				require.NoError(t, err)
				require.Equal(t, tc.expectSeed, r.Seed)
			}
		})
	}
}

func TestWalletNewAddressesHandler(t *testing.T) {
	type httpBody struct {
		ID       string
		Num      string
		Password string
	}
	type Addresses struct {
		Address []string `json:"addresses"`
	}

	var responseAddresses = Addresses{}
	var responseEmptyAddresses = Addresses{}

	var emptyAddrs = make([]cipher.Address, 0)
	var addrs = make([]cipher.Address, 3)

	for i := 0; i < 3; i++ {
		pub, _ := cipher.GenerateDeterministicKeyPair(cipher.RandByte(32))
		addrs[i] = cipher.AddressFromPubKey(pub)
		responseAddresses.Address = append(responseAddresses.Address, addrs[i].String())
	}

	tt := []struct {
		name                      string
		method                    string
		body                      *httpBody
		status                    int
		err                       string
		walletID                  string
		n                         uint64
		password                  string
		gatewayNewAddressesResult []cipher.Address
		gatewayNewAddressesErr    error
		responseBody              Addresses
		csrfDisabled              bool
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:   "400 - invalid num value",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "bar",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid num value",
		},
		{
			name:   "400 - gateway.NewAddresses error",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - gateway.NewAddresses error",
			walletID: "foo",
			n:        1,
			gatewayNewAddressesErr: errors.New("gateway.NewAddresses error"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "foo",
			n:        1,
			gatewayNewAddressesErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "400 Bad Request - missing password",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing password",
			walletID: "foo",
			n:        1,
			gatewayNewAddressesErr: wallet.ErrMissingPassword,
		},
		{
			name:   "400 Bad Request - wallet invalid password",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid password",
			walletID: "foo",
			n:        1,
			gatewayNewAddressesErr: wallet.ErrInvalidPassword,
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
		},
		{
			name:   "200 - OK with password",
			method: http.MethodPost,
			body: &httpBody{
				ID:       "foo",
				Num:      "1",
				Password: "pwd",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
		},
		{
			name:   "200 - OK empty addresses",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "0",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        0,
			gatewayNewAddressesResult: emptyAddrs,
			responseBody:              responseEmptyAddresses,
		},
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:   http.StatusOK,
			walletID: "foo",
			n:        1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
			csrfDisabled:              true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("NewAddresses", tc.walletID, []byte(tc.password), tc.n).Return(tc.gatewayNewAddressesResult, tc.gatewayNewAddressesErr)

			endpoint := "/wallet/newAddress"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.ID != "" {
					v.Add("id", tc.body.ID)
				}
				if tc.body.Num != "" {
					v.Add("num", tc.body.Num)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Addresses
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, msg, tc.name)
			}
		})
	}
}

func TestGetWalletFolderHandler(t *testing.T) {
	tt := []struct {
		name                 string
		method               string
		status               int
		err                  string
		getWalletDirResponse string
		getWalletDirErr      error
		httpResponse         WalletFolder
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                 "200",
			method:               http.MethodGet,
			status:               http.StatusOK,
			getWalletDirResponse: "/wallet/folder/address",
			httpResponse: WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
		{
			name:            "403 - wallet API disabled",
			method:          http.MethodGet,
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			getWalletDirErr: wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		gateway := &GatewayerMock{}
		gateway.On("GetWalletDir").Return(tc.getWalletDirResponse, tc.getWalletDirErr)

		endpoint := "/wallets/folderName"

		req, err := http.NewRequest(tc.method, endpoint, nil)
		require.NoError(t, err)

		csrfStore := &CSRFStore{
			Enabled: true,
		}
		setCSRFParameters(csrfStore, tokenValid, req)

		rr := httptest.NewRecorder()
		handler := newServerMux(mxConfig, gateway, csrfStore)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg WalletFolder
			json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}

func TestWalletUnloadHandler(t *testing.T) {
	tt := []struct {
		name            string
		method          string
		status          int
		err             string
		walletID        string
		unloadWalletErr error
		csrfDisabled    bool
	}{
		{
			name:     "405",
			method:   http.MethodGet,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "wallet.wlt",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:            "403 - Forbidden - wallet API disabled",
			method:          http.MethodPost,
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			walletID:        "wallet.wlt",
			unloadWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:     "200 - ok",
			method:   http.MethodPost,
			status:   http.StatusOK,
			walletID: "wallet.wlt",
		},
		{
			name:         "200 - ok, csrf disabled",
			method:       http.MethodPost,
			status:       http.StatusOK,
			walletID:     "wallet.wlt",
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &GatewayerMock{}
			gateway.On("UnloadWallet", tc.walletID).Return(tc.unloadWalletErr)

			endpoint := "/wallet/unload"
			v := url.Values{}
			v.Add("id", tc.walletID)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %d, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			}
		})
	}
}

func TestEncryptWallet(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type gatewayReturnPair struct {
		w   *wallet.Wallet
		err error
	}
	tt := []struct {
		name          string
		method        string
		wltID         string
		password      string
		gatewayReturn gatewayReturnPair
		status        int
		expectWallet  WalletResponse
		expectErr     string
	}{
		{
			name:     "200 - OK",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				w: &wallet.Wallet{
					Meta: map[string]string{
						"filename":  "wallet.wlt",
						"seed":      "seed",
						"lastSeed":  "lastSeed",
						"secrets":   "secrets",
						"encrypted": "true",
					},
					Entries: cloneEntries(entries),
				},
			},
			status: http.StatusOK,
			expectWallet: WalletResponse{
				Meta: WalletMeta{
					Filename:  "wallet.wlt",
					Encrypted: true,
				},
				Entries: responseEntries,
			},
		},
		{
			name:     "403 Forbidden",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletAPIDisabled,
			},
			status:    http.StatusForbidden,
			expectErr: "403 Forbidden",
		},
		{
			name:      "405 Method Not Allowed",
			method:    http.MethodGet,
			wltID:     "wallet.wlt",
			password:  "pwd",
			status:    http.StatusMethodNotAllowed,
			expectErr: "405 Method Not Allowed",
		},
		{
			name:      "400 - Missing Password",
			method:    http.MethodPost,
			wltID:     "wallet.wlt",
			password:  "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing password",
		},
		{
			name:      "400 - Missing Wallet Id",
			method:    http.MethodPost,
			wltID:     "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing wallet id",
		},
		{
			name:     "400 - Invalid Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrInvalidPassword,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - invalid password",
		},
		{
			name:     "404 - Wallet Not Found",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotExist,
			},
			status:    http.StatusNotFound,
			expectErr: "404 Not Found",
		},
		{
			name:     "400 - Wallet Is Already Encrypted",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletEncrypted,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - wallet is already encrypted",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			gateway.On("EncryptWallet", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturn.w, tc.gatewayReturn.err)

			endpoint := "/wallet/encrypt"
			v := url.Values{}
			v.Add("id", tc.wltID)
			v.Add("password", tc.password)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`, body: %v", status, tc.status, rr.Body.String())

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
				return
			}

			var rlt WalletResponse
			err = json.NewDecoder(rr.Body).Decode(&rlt)
			require.NoError(t, err)
			require.Equal(t, tc.expectWallet, rlt)
		})
	}
}

func TestDecryptWallet(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type gatewayReturnPair struct {
		w   *wallet.Wallet
		err error
	}

	tt := []struct {
		name          string
		method        string
		wltID         string
		password      string
		gatewayReturn gatewayReturnPair
		status        int
		expectWallet  WalletResponse
		expectErr     string
	}{
		{
			name:     "200 OK",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				w: &wallet.Wallet{
					Meta: map[string]string{
						"filename":  "wallet",
						"seed":      "seed",
						"lastSeed":  "lastSeed",
						"secrets":   "",
						"encrypted": "false",
					},
					Entries: cloneEntries(entries),
				},
			},
			status: http.StatusOK,
			expectWallet: WalletResponse{
				Meta: WalletMeta{
					Filename:  "wallet",
					Encrypted: false,
				},
				Entries: responseEntries,
			},
		},
		{
			name:     "403 Forbidden",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletAPIDisabled,
			},
			status:    http.StatusForbidden,
			expectErr: "403 Forbidden",
		},
		{
			name:      "405 Method Not Allowed",
			method:    http.MethodGet,
			status:    http.StatusMethodNotAllowed,
			expectErr: "405 Method Not Allowed",
		},
		{
			name:      "400 - Missing Wallet ID",
			method:    http.MethodPost,
			wltID:     "",
			password:  "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing wallet id",
		},
		{
			name:      "400 - Missing Password",
			method:    http.MethodPost,
			wltID:     "wallet.wlt",
			password:  "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing password",
		},
		{
			name:     "400 - Wallet IS Not Encrypted",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotEncrypted,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - wallet is not encrypted",
		},
		{
			name:     "400 - Invalid Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrInvalidPassword,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - invalid password",
		},
		{
			name:     "404 - Wallet Does Not Exist",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotExist,
			},
			status:    http.StatusNotFound,
			expectErr: "404 Not Found",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := NewGatewayerMock()
			gateway.On("DecryptWallet", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturn.w, tc.gatewayReturn.err)

			endpoint := "/wallet/decrypt"
			v := url.Values{}
			v.Add("id", tc.wltID)
			v.Add("password", tc.password)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(mxConfig, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
				return
			}

			var r WalletResponse
			err = json.NewDecoder(rr.Body).Decode(&r)
			require.NoError(t, err)
			require.Equal(t, tc.expectWallet, r)
		})
	}
}

// makeEntries derives N wallet address entries from given seed
// Returns set of wallet.Entry and wallet.ReadableEntry, the readable
// entries' secrets are removed.
func makeEntries(seed []byte, n int) ([]wallet.Entry, []WalletEntry) {
	seckeys := cipher.GenerateDeterministicKeyPairs(seed, n)
	var entries []wallet.Entry
	var responseEntries []WalletEntry
	for i, seckey := range seckeys {
		pubkey := cipher.PubKeyFromSecKey(seckey)
		entries = append(entries, wallet.Entry{
			Address: cipher.AddressFromPubKey(pubkey),
			Public:  pubkey,
			Secret:  seckey,
		})
		responseEntries = append(responseEntries, WalletEntry{
			Address: entries[i].Address.String(),
			Public:  entries[i].Public.Hex(),
		})
	}
	return entries, responseEntries
}

func cloneEntries(es []wallet.Entry) []wallet.Entry {
	var entries []wallet.Entry
	for _, e := range es {
		entries = append(entries, e)
	}
	return entries
}
