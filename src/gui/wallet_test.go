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

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

// Gateway RPC interface wrapper for daemon state
type FakeGateway struct {
	mock.Mock
	walletID string
	coins    uint64
	dst      cipher.Address
	t        *testing.T
}

func (gw *FakeGateway) Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	args := gw.Called(wltID, coins, dest)
	return args.Get(0).(*coin.Transaction), args.Error(1)
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) GetWalletBalance(wltID string) (wallet.BalancePair, error) {
	args := gw.Called(wltID)
	return args.Get(0).(wallet.BalancePair), args.Error(1)
}

// GetWallet returns the wallet by wltID
func (gw *FakeGateway) GetWallet(wltID string) (wallet.Wallet, error) {
	args := gw.Called(wltID)
	return args.Get(0).(wallet.Wallet), args.Error(1)
}

func (gw *FakeGateway) UpdateWalletLabel(wltID, label string) error {
	args := gw.Called(wltID, label)
	return args.Error(0)
}

// NewAddresses generate addresses in given wallet
func (gw *FakeGateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error) {
	args := gw.Called(wltID)
	return args.Get(0).([]visor.UnconfirmedTxn), args.Error(1)
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) CreateWallet(wltName string, options wallet.Options) (wallet.Wallet, error) {
	args := gw.Called(wltName, options)
	return args.Get(0).(wallet.Wallet), args.Error(1)
}

func (gw *FakeGateway) ScanAheadWalletAddresses(wltName string, scanN uint64) (wallet.Wallet, error) {
	args := gw.Called(wltName, scanN)
	return args.Get(0).(wallet.Wallet), args.Error(1)
}

// NewAddresses generate addresses in given wallet
func (gw *FakeGateway) NewAddresses(wltID string, n uint64) ([]cipher.Address, error) {
	args := gw.Called(wltID, n)
	return args.Get(0).([]cipher.Address), args.Error(1)
}

func (gw *FakeGateway) GetWalletDir() string {
	args := gw.Called()
	return args.String(0)
}

func TestWalletSpendHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
	}

	tt := []struct {
		name                          string
		method                        string
		url                           string
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
	}{
		{
			"405",
			http.MethodGet,
			"/wallet/spend",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - no walletID",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{},
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - no dst",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
			},
			http.StatusBadRequest,
			"400 Bad Request - missing destination address \"dst\"",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - bad dst addr",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      " 2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid destination address: Invalid base58 character",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - no coins",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - coins is string",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "foo",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - coins is negative value",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "-123",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - zero coins",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "0",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid \"coins\" value, must > 0",
			"0",
			0,
			"",
			nil,
			nil,
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - gw spend error txn no fee",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusBadRequest,
			"400 Bad Request - Transaction has zero coinhour fee",
			"123",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			fee.ErrTxnNoFee,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: fee.ErrTxnNoFee.Error(),
			},
		},
		{
			"400 - gw spend error spending unconfirmed",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusBadRequest,
			"400 Bad Request - please spend after your pending transaction is confirmed",
			"123",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			wallet.ErrSpendingUnconfirmed,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: wallet.ErrSpendingUnconfirmed.Error(),
			},
		},
		{
			"400 - gw spend error insufficient balance",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusBadRequest,
			"400 Bad Request - balance is not sufficient",
			"123",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			wallet.ErrInsufficientBalance,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: wallet.ErrInsufficientBalance.Error(),
			},
		},
		{
			"404 - gw spend error wallet not exist",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusNotFound,
			"404 Not Found",
			"123",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			wallet.ErrWalletNotExist,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: wallet.ErrWalletNotExist.Error(),
			},
		},
		{
			"500 - gw spend error",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusInternalServerError,
			"500 Internal Server Error - Spend error",
			"123",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			errors.New("Spend error"),
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: "Spend error",
			},
		},
		{
			"200 - gw GetWalletBalance error",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusOK,
			"",
			"1234",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			nil,
			wallet.BalancePair{},
			errors.New("GetWalletBalance error"),
			&SpendResult{
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
			"200 - OK",
			http.MethodPost,
			"/wallet/spend",
			&httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			http.StatusOK,
			"",
			"1234",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			nil,
			wallet.BalancePair{},
			nil,
			&SpendResult{
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
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				walletID: tc.walletID,
				t:        t,
			}
			addr, _ := cipher.DecodeBase58Address(tc.dst)
			gateway.On("Spend", tc.walletID, tc.coins, addr).Return(tc.gatewaySpendResult, tc.gatewaySpendErr)
			gateway.On("GetWalletBalance", tc.walletID).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)

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

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(walletSpendHandler(gateway))

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
		walletId               string
		gatewayGetWalletResult wallet.Wallet
		gatewayGetWalletErr    error
	}{
		{
			"405",
			http.MethodPost,
			"/wallet",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			wallet.Wallet{},
			nil,
		},
		{
			"400 - no walletId",
			http.MethodGet,
			"/wallet",
			nil,
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"",
			wallet.Wallet{},
			nil,
		},
		{
			"400 - error from the `gateway.GetWallet(wltID)`",
			http.MethodGet,
			"/wallet",
			&httpBody{
				WalletID: "123",
			},
			http.StatusBadRequest,
			"400 Bad Request - wallet 123 doesn't exist",
			"123",
			wallet.Wallet{},
			errors.New("wallet 123 doesn't exist"),
		},
		{
			"200 - OK",
			http.MethodGet,
			"/wallet",
			&httpBody{
				WalletID: "1234",
			},
			http.StatusOK,
			"",
			"1234",
			wallet.Wallet{
				Meta:    map[string]string{},
				Entries: []wallet.Entry{},
			},
			nil,
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			walletID: tc.walletId,
			t:        t,
		}
		gateway.On("GetWallet", tc.walletId).Return(tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)
		v := url.Values{}
		var url = tc.url
		if tc.body != nil {
			if tc.body.WalletID != "" {
				v.Add("id", tc.body.WalletID)
			}
		}

		if len(v) > 0 {
			url += "?" + v.Encode()
		}

		req, err := http.NewRequest(tc.method, url, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(walletGet(gateway))

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
		url                           string
		body                          *httpBody
		status                        int
		err                           string
		walletId                      string
		gatewayGetWalletBalanceResult wallet.BalancePair
		gatewayBalanceErr             error
		result                        *wallet.BalancePair
	}{
		{
			"405",
			http.MethodPost,
			"/wallet/balance",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - no walletId",
			http.MethodGet,
			"/wallet/balance",
			nil,
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"0",
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"404 - gw `wallet doesn't exist` error",
			http.MethodGet,
			"/wallet/balance",
			&httpBody{
				WalletID: "notFoundId",
			},
			http.StatusNotFound,
			"404 Not Found",
			"notFoundId",
			wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
			wallet.ErrWalletNotExist,
			&wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			"500 - gw other error",
			http.MethodGet,
			"/wallet/balance",
			&httpBody{
				WalletID: "someId",
			},
			http.StatusInternalServerError,
			"500 Internal Server Error - gatewayBalanceError",
			"someId",
			wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
			errors.New("gatewayBalanceError"),
			&wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			"200 - OK",
			http.MethodGet,
			"/wallet/balance",
			&httpBody{
				WalletID: "foo",
			},
			http.StatusOK,
			"",
			"foo",
			wallet.BalancePair{},
			nil,
			&wallet.BalancePair{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				walletID: tc.walletId,
				t:        t,
			}
			gateway.On("GetWalletBalance", tc.walletId).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)

			v := url.Values{}
			var url = tc.url
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
			}
			if len(v) > 0 {
				url += "?" + v.Encode()
			}
			req, err := http.NewRequest(tc.method, url, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(walletBalanceHandler(gateway))

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
		walletId                    string
		label                       string
		gatewayUpdateWalletLabelErr error
		responseBody                string
	}{
		{
			"405",
			http.MethodGet,
			"/wallet/update",
			&httpBody{},
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			"",
			nil,
			"",
		},
		{
			"400 - missing wallet id",
			http.MethodPost,
			"/wallet/update",
			&httpBody{},
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"",
			"",
			nil,
			"",
		},
		{
			"400 - missing label",
			http.MethodPost,
			"/wallet/update",
			&httpBody{
				WalletID: "foo",
			},
			http.StatusBadRequest,
			"400 Bad Request - missing label",
			"foo",
			"",
			nil,
			"",
		},
		{
			"404 - gateway.UpdateWalletLabel ErrWalletNotExist",
			http.MethodPost,
			"/wallet/update",
			&httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			http.StatusNotFound,
			"404 Not Found",
			"foo",
			"label",
			wallet.ErrWalletNotExist,
			"",
		},
		{
			"500 - gateway.UpdateWalletLabel error",
			http.MethodPost,
			"/wallet/update",
			&httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			http.StatusInternalServerError,
			"500 Internal Server Error - gateway.UpdateWalletLabel error",
			"foo",
			"label",
			errors.New("gateway.UpdateWalletLabel error"),
			"",
		},
		{
			"200 OK",
			http.MethodPost,
			"/wallet/update",
			&httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			http.StatusOK,
			"",
			"foo",
			"label",
			nil,
			"\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("UpdateWalletLabel", tc.walletId, tc.label).Return(tc.gatewayUpdateWalletLabelErr)

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
				if tc.body.Label != "" {
					v.Add("label", tc.body.Label)
				}
			}
			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(walletUpdateHandler(gateway))

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

	tt := []struct {
		name                                  string
		method                                string
		url                                   string
		body                                  *httpBody
		status                                int
		err                                   string
		walletId                              string
		gatewayGetWalletUnconfirmedTxnsResult []visor.UnconfirmedTxn
		gatewayGetWalletUnconfirmedTxnsErr    error
		responseBody                          []visor.UnconfirmedTxn
	}{
		{
			"405",
			http.MethodPost,
			"/wallet/transactions",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			make([]visor.UnconfirmedTxn, 0),
			nil,
			[]visor.UnconfirmedTxn{},
		},
		{
			"400 - missing wallet id",
			http.MethodGet,
			"/wallet/transactions",
			nil,
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"",
			make([]visor.UnconfirmedTxn, 0),
			nil,
			[]visor.UnconfirmedTxn{},
		},
		{
			"500 - gateway.GetWalletUnconfirmedTxns error",
			http.MethodGet,
			"/wallet/transactions",
			&httpBody{
				WalletID: "foo",
			},
			http.StatusInternalServerError,
			"500 Internal Server Error - gateway.GetWalletUnconfirmedTxns error",
			"foo",
			make([]visor.UnconfirmedTxn, 0),
			errors.New("gateway.GetWalletUnconfirmedTxns error"),
			[]visor.UnconfirmedTxn{},
		},
		{
			"404 - wallet doesn't exist",
			http.MethodGet,
			"/wallet/transactions",
			&httpBody{
				WalletID: "foo",
			},
			http.StatusNotFound,
			"404 Not Found",
			"foo",
			make([]visor.UnconfirmedTxn, 0),
			wallet.ErrWalletNotExist,
			[]visor.UnconfirmedTxn{},
		},
		{
			"200 - OK",
			http.MethodGet,
			"/wallet/transactions",
			&httpBody{
				WalletID: "foo",
			},
			http.StatusOK,
			"",
			"foo",
			make([]visor.UnconfirmedTxn, 0),
			nil,
			[]visor.UnconfirmedTxn{},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("GetWalletUnconfirmedTxns", tc.walletId).Return(tc.gatewayGetWalletUnconfirmedTxnsResult, tc.gatewayGetWalletUnconfirmedTxnsErr)
		v := url.Values{}
		var urlFull = tc.url
		if tc.body != nil {
			if tc.body.WalletID != "" {
				v.Add("id", tc.body.WalletID)
			}
		}
		if len(v) > 0 {
			urlFull = urlFull + "?" + v.Encode()
		}
		req, err := http.NewRequest(tc.method, urlFull, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(walletTransactionsHandler(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg []visor.UnconfirmedTxn
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.responseBody, msg, tc.name)
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
		url                       string
		body                      *httpBody
		status                    int
		err                       string
		wltname                   string
		scnN                      uint64
		options                   wallet.Options
		gatewayCreateWalletResult wallet.Wallet
		gatewayCreateWalletErr    error
		scanWalletAddressesResult wallet.Wallet
		scanWalletAddressesError  error
		responseBody              wallet.ReadableWallet
	}{
		{
			"405",
			http.MethodGet,
			"/wallet/create",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"foo",
			0,
			wallet.Options{},
			wallet.Wallet{},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"400 - missing seed",
			http.MethodPost,
			"/wallet/create",
			&httpBody{},
			http.StatusBadRequest,
			"400 Bad Request - missing seed",
			"foo",
			0,
			wallet.Options{},
			wallet.Wallet{},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"400 - missing label",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed: "foo",
			},
			http.StatusBadRequest,
			"400 Bad Request - missing label",
			"foo",
			0,
			wallet.Options{},
			wallet.Wallet{},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"400 - invalid scan value",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "bad scanN",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid scan value",
			"foo",
			0,
			wallet.Options{},
			wallet.Wallet{},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"400 - scan must be > 0",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "0",
			},
			http.StatusBadRequest,
			"400 Bad Request - scan must be > 0",
			"foo",
			0,
			wallet.Options{},
			wallet.Wallet{},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"400 - gateway.CreateWallet error",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "1",
			},
			http.StatusBadRequest,
			"400 Bad Request - gateway.CreateWallet error",
			"",
			0,
			wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			wallet.Wallet{},
			errors.New("gateway.CreateWallet error"),
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{},
		},
		{
			"500 - gateway.ScanAheadWalletAddresses error",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			http.StatusInternalServerError,
			"500 Internal Server Error",
			"filename",
			2,
			wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			nil,
			wallet.Wallet{},
			errors.New("gateway.ScanAheadWalletAddresses error"),
			wallet.ReadableWallet{},
		},
		{
			"200 - OK",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			http.StatusOK,
			"",
			"filename",
			2,
			wallet.Options{
				Label: "bar",
				Seed:  "foo",
			},
			wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			nil,
			wallet.Wallet{},
			nil,
			wallet.ReadableWallet{
				Meta:    map[string]string{},
				Entries: wallet.ReadableEntries{},
			},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("CreateWallet", "", tc.options).Return(tc.gatewayCreateWalletResult, tc.gatewayCreateWalletErr)
		gateway.On("ScanAheadWalletAddresses", tc.wltname, tc.scnN-1).Return(tc.scanWalletAddressesResult, tc.scanWalletAddressesError)
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

		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(v.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(walletCreate(gateway))

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
		name       string
		method     string
		url        string
		body       *httpBody
		statusCode int
		err        string
		entropy    string
		resultLen  int
	}{
		{
			"405",
			http.MethodPut,
			"/wallet/newSeed",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			0,
		},
		{
			"400 - invalid entropy type",
			http.MethodGet,
			"/wallet/newSeed",
			&httpBody{
				Entropy: "xx",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid entropy",
			"xx",
			0,
		},
		{
			"400 - `wrong entropy length` error",
			http.MethodGet,
			"/wallet/balance",
			&httpBody{
				Entropy: "200",
			},
			http.StatusBadRequest,
			"400 Bad Request - entropy length must be 128 or 256",
			"200",
			0,
		},
		{
			"200 - OK with no entropy",
			http.MethodGet,
			"/wallet/newSeed",
			&httpBody{
				Entropy: "",
			},
			http.StatusOK,
			"",
			"128",
			12,
		},
		{
			"200 - OK | 12 word seed",
			http.MethodGet,
			"/wallet/newSeed",
			&httpBody{
				Entropy: "128",
			},
			http.StatusOK,
			"",
			"128",
			12,
		},
		{
			"200 - OK | 24 word seed",
			http.MethodGet,
			"/wallet/newSeed",
			&httpBody{
				Entropy: "256",
			},
			http.StatusOK,
			"",
			"256",
			24,
		},
	}

	// Loop over each test case
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{t: t}
			// Add request parameters to url
			v := url.Values{}
			var url = tc.url
			if tc.body != nil {
				if tc.body.Entropy != "" {
					v.Add("entropy", tc.body.Entropy)
				}
			}
			if len(v) > 0 {
				url += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, url, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(newWalletSeed(gateway))
			handler.ServeHTTP(rr, req)

			statusCode := rr.Code
			require.Equal(t, tc.statusCode, statusCode, "case: %s, handler returned wrong status code: got `%v` expected `%v`", tc.name, statusCode, tc.statusCode)
			if statusCode != tc.statusCode {
				t.Errorf("case: %s, handler returned wrong status code: got `%v` want `%v`",
					tc.name, statusCode, tc.statusCode)
			}
			if statusCode != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, expected `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), statusCode, tc.err)
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
		Id  string
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
		url                       string
		body                      *httpBody
		status                    int
		err                       string
		walletId                  string
		n                         uint64
		gatewayNewAddressesResult []cipher.Address
		gatewayNewAddressesErr    error
		responseBody              Addresses
	}{
		{
			"405",
			http.MethodGet,
			"/wallet/create",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"",
			1,
			make([]cipher.Address, 0),
			nil,
			Addresses{},
		},
		{
			"400 - missing wallet id",
			http.MethodPost,
			"/wallet/create",
			nil,
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"foo",
			1,
			make([]cipher.Address, 0),
			nil,
			Addresses{},
		},
		{
			"400 - invalid num value",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Id:  "foo",
				Num: "bar",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid num value",
			"foo",
			1,
			make([]cipher.Address, 0),
			nil,
			Addresses{},
		},
		{
			"400 - gateway.NewAddresses error",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Id:  "foo",
				Num: "1",
			},
			http.StatusBadRequest,
			"400 Bad Request - gateway.NewAddresses error",
			"foo",
			1,
			make([]cipher.Address, 0),
			errors.New("gateway.NewAddresses error"),
			Addresses{},
		},
		{
			"200 - OK",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Id:  "foo",
				Num: "1",
			},
			http.StatusOK,
			"",
			"foo",
			1,
			addrs,
			nil,
			responseAddresses,
		},
		{
			"200 - OK empty addresses",
			http.MethodPost,
			"/wallet/create",
			&httpBody{
				Id:  "foo",
				Num: "0",
			},
			http.StatusOK,
			"",
			"foo",
			0,
			emptyAddrs,
			nil,
			responseEmptyAddresses,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("NewAddresses", tc.walletId, tc.n).Return(tc.gatewayNewAddressesResult, tc.gatewayNewAddressesErr)
			v := url.Values{}
			if tc.body != nil {
				if tc.body.Id != "" {
					v.Add("id", tc.body.Id)
				}
				if tc.body.Num != "" {
					v.Add("num", tc.body.Num)
				}
			}

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(walletNewAddresses(gateway))

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
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
		url                  string
		status               int
		err                  string
		getWalletDirResponse string
		httpResponse         WalletFolder
	}{
		{
			"200",
			http.MethodGet,
			"/wallets/folderName",
			http.StatusOK,
			"",
			"/wallet/folder/address",
			WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
		{
			"200 - POST",
			http.MethodPost,
			"/wallets/folderName",
			http.StatusOK,
			"",
			"/wallet/folder/address",
			WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("GetWalletDir").Return(tc.getWalletDirResponse)

		req, err := http.NewRequest(tc.method, tc.url, nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getWalletFolder(gateway))

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
