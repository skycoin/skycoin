package gui

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestWalletSpendHandler(t *testing.T) {
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
		coins                         uint64
		dst                           string
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
			gatewaySpendErr: wallet.ErrWalletApiDisabled,
			spendResult: &SpendResult{
				Error: wallet.ErrWalletApiDisabled.Error(),
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewaySpendResult == nil {
				tc.gatewaySpendResult = &coin.Transaction{}
			}

			gateway := &GatewayerMock{}
			addr, _ := cipher.DecodeBase58Address(tc.dst)
			gateway.On("Spend", tc.walletID, tc.coins, addr).Return(tc.gatewaySpendResult, tc.gatewaySpendErr)
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
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
				Meta:    map[string]string{},
				Entries: []wallet.Entry{},
			},
			gatewayGetWalletErr: wallet.ErrWalletApiDisabled,
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
				Meta:    map[string]string{},
				Entries: []wallet.Entry{},
			},
		},
	}

	for _, tc := range tt {
		gateway := &GatewayerMock{}
		gateway.On("GetWallet", tc.walletID).Return(tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)
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
		handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
				"case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg wallet.Wallet
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.gatewayGetWalletResult, msg, tc.name)
		}
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
			gatewayBalanceErr:             wallet.ErrWalletApiDisabled,
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
			gatewayUpdateWalletLabelErr: wallet.ErrWalletApiDisabled,
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
			gatewayGetWalletUnconfirmedTxnsErr: wallet.ErrWalletApiDisabled,
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
		handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
	type httpBody struct {
		Seed  string
		Label string
		ScanN string
	}
	tt := []struct {
		name                      string
		method                    string
		body                      *httpBody
		status                    int
		err                       string
		wltName                   string
		scnN                      uint64
		options                   wallet.Options
		gatewayCreateWalletResult wallet.Wallet
		gatewayCreateWalletErr    error
		scanWalletAddressesResult wallet.Wallet
		scanWalletAddressesError  error
		responseBody              wallet.ReadableWallet
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
				Label: "bar",
				Seed:  "foo",
			},
			gatewayCreateWalletErr: errors.New("gateway.CreateWallet error"),
		},
		{
			name:   "500 - gateway.ScanAheadWalletAddresses error",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusInternalServerError,
			err:     "500 Internal Server Error",
			wltName: "filename",
			scnN:    2,
			options: wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			scanWalletAddressesError: errors.New("gateway.ScanAheadWalletAddresses error"),
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
			scnN:    2,
			options: wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			gatewayCreateWalletErr: wallet.ErrWalletApiDisabled,
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
			scnN:    2,
			options: wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			responseBody: wallet.ReadableWallet{
				Meta:    map[string]string{},
				Entries: wallet.ReadableEntries{},
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
			scnN:    2,
			options: wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			responseBody: wallet.ReadableWallet{
				Meta:    map[string]string{},
				Entries: wallet.ReadableEntries{},
			},
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		gateway := &GatewayerMock{}
		gateway.On("CreateWallet", "", tc.options).Return(tc.gatewayCreateWalletResult, tc.gatewayCreateWalletErr)
		gateway.On("ScanAheadWalletAddresses", tc.wltName, tc.scnN-1).Return(tc.scanWalletAddressesResult, tc.scanWalletAddressesError)

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
		handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
				"case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg wallet.ReadableWallet
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.responseBody, msg, tc.name)
		}
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
			gateway.On("IsWalletAPIDisabled").Return(false)

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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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

func TestWalletNewAddressesHandler(t *testing.T) {
	type httpBody struct {
		ID  string
		Num string
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
			gatewayNewAddressesErr: wallet.ErrWalletApiDisabled,
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
			gateway.On("NewAddresses", tc.walletID, tc.n).Return(tc.gatewayNewAddressesResult, tc.gatewayNewAddressesErr)

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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
			getWalletDirErr: wallet.ErrWalletApiDisabled,
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
		handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
			unloadWalletErr: wallet.ErrWalletApiDisabled,
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
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)

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
