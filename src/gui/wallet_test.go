package gui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

type httpBody struct {
	Id    string `url:"id,omitempty"`
	Dst   string `url:"dst,omitempty"`
	Coins string `url:"coins,omitempty"`
}

// Gateway RPC interface wrapper for daemon state
type FakeGateway struct {
	mock.Mock
	walletId string
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

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) GetWallet(wltID string) (wallet.Wallet, error) {
	args := gw.Called(wltID)
	return args.Get(0).(wallet.Wallet), args.Error(1)
}

func TestWalletHandler(t *testing.T) {
	tt := []struct {
		name                   string
		method                 string
		url                    string
		body                   *httpBody
		status                 int
		err                    string
		walletId               string
		coins                  uint64
		dst                    string
		gatewayGetWalletResult wallet.Wallet
		gatewayGetWalletErr    error
		result                 *wallet.Wallet
	}{
		{
			"405",
			"PUT",
			"/wallet",
			nil,
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			0,
			"",
			wallet.Wallet{},
			nil,
			nil,
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			walletId: tc.walletId,
			t:        t,
		}
		gateway.On("GetWalletBalance", tc.walletId).Return(tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)
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
