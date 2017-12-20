package gui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

type httpBody struct {
	Id    string `url:"id,omitempty"`
	Dst   string `url:"dst,omitempty"`
	Coins string `url:"coins,omitempty"`
	Seed  string `url:"seed,omitempty"`
	Label string `url:"label,omitempty"`
	ScanN string `url:"scan,omitempty"`
	Num   string `url:"num,omitempty"`
}

// Gateway RPC interface wrapper for daemon state
type FakeGateway struct {
	mock.Mock
	walletId string
	label    string
	coins    uint64
	wltName  string
	scanN    uint64
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

// NewAddresses generate addresses in given wallet
func (gw *FakeGateway) GetWalletUnconfirmedTxns(wltID string) ([]visor.UnconfirmedTxn, error) {
	args := gw.Called(wltID)
	return args.Get(0).([]visor.UnconfirmedTxn), args.Error(1)
}

func (gw *FakeGateway) UpdateWalletLabel(wltID, label string) error {
	args := gw.Called(wltID, label)
	return args.Error(0)
}

func (gw *FakeGateway) ReloadWallets() error {
	args := gw.Called()
	return args.Error(0)
}

func TestWalletsReloadHandler(t *testing.T) {
	tt := []struct {
		name                    string
		method                  string
		url                     string
		body                    *httpBody
		status                  int
		err                     string
		gatewayWalletsReloadErr error
		responseBody            string
	}{
		{
			"500 - gateway.ReloadWallets error",
			http.MethodGet,
			"/wallets/reload",
			&httpBody{},
			http.StatusInternalServerError,
			"500 Internal Server Error",
			errors.New("gateway.ReloadWallets error"),
			"",
		},
		{
			"200 - OK",
			http.MethodGet,
			"/wallets/reload",
			nil,
			http.StatusOK,
			"",
			nil,
			"\"success\"",
		},
		{
			"200 - OK POST",
			http.MethodPost,
			"/wallets/reload",
			nil,
			http.StatusOK,
			"",
			nil,
			"\"success\"",
		},
		{
			"200 - OK trailing backslash",
			http.MethodGet,
			"/wallets/reload/",
			nil,
			http.StatusOK,
			"",
			nil,
			"\"success\"",
		},
	}

	for _, tc := range tt {
		gateway := &FakeGateway{
			t: t,
		}
		gateway.On("ReloadWallets").Return(tc.gatewayWalletsReloadErr)
		params, _ := query.Values(tc.body)
		paramsEncoded := params.Encode()
		var url = tc.url
		if paramsEncoded != "" {
			url = url + "?" + paramsEncoded
		}
		req, err := http.NewRequest(tc.method, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WalletsReloadHandler(gateway))

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
	}
}
