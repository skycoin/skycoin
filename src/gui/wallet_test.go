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

	"github.com/pkg/errors"

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
	var w wallet.Wallet
	args := gw.Called(wltID)
	if args.Get(0) != nil {
		return args.Get(0).(wallet.Wallet), args.Error(1)
	} else {
		return w, args.Error(1)
	}
}

// GetWalletBalance returns balance pair of specific wallet
func (gw *FakeGateway) CreateWallet(wltID string) (wallet.Wallet, error) {
	var w wallet.Wallet
	args := gw.Called(wltID)
	if args.Get(0) != nil {
		return args.Get(0).(wallet.Wallet), args.Error(1)
	} else {
		return w, args.Error(1)
	}
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
		gatewayGetWalletResult wallet.Wallet
		gatewayGetWalletErr    error
	}{
		{
			"405",
			http.MethodPut,
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
				Id: "123",
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
				Id: "1234",
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
			walletId: tc.walletId,
			t:        t,
		}
		gateway.On("GetWallet", tc.walletId).Return(tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)
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
		handler := http.HandlerFunc(WalletGet(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, strings.TrimSpace(rr.Body.String()), tc.err, "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg wallet.Wallet
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			if err != nil {
				t.Errorf("fail unmarshal json response while 200 OK. body: %s, err: %s", rr.Body.String(), err)
			}
			require.Equal(t, tc.gatewayGetWalletResult, msg, tc.name)
		}
	}
}
