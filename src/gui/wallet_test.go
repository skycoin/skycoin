package gui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"net/url"
	"bytes"
)

type httpBody struct {
	Id    string `url:"id,omitempty"`
	Dst   string `url:"dst,omitempty"`
	Coins string `url:"coins,omitempty"`
}

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
			"GET",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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
			"POST",
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

func TestWalletBalanceHandler(t *testing.T) {
	tt := []struct {
		name                          string
		method                        string
		url                           string
		body                          *httpBody
		status                        int
		err                           string
		walletId                      string
		coins                         uint64
		dst                           string
		gatewayGetWalletBalanceResult wallet.BalancePair
		gatewayBalanceErr             error
		result                        *wallet.BalancePair
	}{
		{
			"405",
			"PUT",
			"/wallet/balance",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			0,
			"",
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"400 - no walletId",
			"GET",
			"/wallet/balance",
			nil,
			http.StatusBadRequest,
			"400 Bad Request - missing wallet id",
			"0",
			0,
			"",
			wallet.BalancePair{},
			nil,
			nil,
		},
		{
			"200 - OK",
			"GET",
			"/wallet/balance",
			&httpBody{
				Id: "foo",
			},
			http.StatusOK,
			"",
			"foo",
			0,
			"",
			wallet.BalancePair{},
			nil,
			&wallet.BalancePair{},
		},
		{
			"200 - but with err from `b, err := gateway.GetWalletBalance(wltID)`",
			"GET",
			"/wallet/balance",
			&httpBody{
				Id: "someId",
			},
			http.StatusOK,
			"",
			"someId",
			0,
			"",
			wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
			errors.New("200 - but with err from `b, err := gateway.GetWalletBalance(wltID)`"),
			&wallet.BalancePair{
				Confirmed: wallet.Balance{Coins: 0, Hours: 0},
				Predicted: wallet.Balance{Coins: 0, Hours: 0},
			},
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			walletId: tc.walletId,
			t:        t,
		}
		gateway.On("GetWalletBalance", tc.walletId).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)
		query, _ := query.Values(tc.body)
		params := query.Encode()
		var url = tc.url
		if params != "" {
			url = url + "?" + params
		}
		req, err := http.NewRequest(tc.method, url, nil)

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WalletBalanceHandler(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		if status != tc.status {
			t.Errorf("case: %s, handler returned wrong status code: got `%v` want `%v`",
				tc.name, status, tc.status)
		}
		if status != http.StatusOK {
			if errMsg := rr.Body.String(); strings.TrimSpace(errMsg) != tc.err {
				t.Errorf("case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, errMsg, status, tc.err)
			}
		} else {
			var msg wallet.BalancePair
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			if err != nil {
				t.Errorf("fail unmarshal json response while 200 OK. body: %s, err: %s", rr.Body.String(), err)
			}
			require.Equal(t, tc.result, &msg, tc.name)
		}
	}
}
