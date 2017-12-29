package gui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"encoding/json"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

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

func (gw *FakeGateway) GetWalletDir() string {
	args := gw.Called()
	return args.String(0)
}

func TestGetWalletFolderHandler(t *testing.T) {
	type httpBody struct{}

	tt := []struct {
		name                 string
		method               string
		url                  string
		body                 *httpBody
		status               int
		err                  string
		getWalletDirResponse string
		httpResponse         WalletFolder
	}{
		{
			"200 - OK",
			http.MethodGet,
			"/wallets/folderName",
			&httpBody{},
			http.StatusOK,
			"",
			"/wallet/folder/address",
			WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
		{
			"200 - OK. trailed backslash",
			http.MethodGet,
			"/wallets/folderName/",
			&httpBody{},
			http.StatusOK,
			"",
			"/wallet/folder/address",
			WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
		{
			"200 -OK. POST",
			http.MethodPost,
			"/wallets/folderName",
			&httpBody{}, http.StatusOK,
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
		handler := http.HandlerFunc(GetWalletFolder(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg WalletFolder
			if err := json.Unmarshal(rr.Body.Bytes(), &msg); err != nil {
				t.Fatal("Failed unmarshal responseBidy `%s`: %v", rr.Body.String(), err)
			}
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}
