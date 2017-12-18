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

func TestWalletTransactionsHandler(t *testing.T) {

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
			http.MethodPut,
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
			"400 - gateway.GetWalletUnconfirmedTxns error",
			http.MethodGet,
			"/wallet/transactions",
			&httpBody{
				Id: "foo",
			},
			http.StatusBadRequest,
			"400 Bad Request - get wallet unconfirmed transactions failed: gateway.GetWalletUnconfirmedTxns error",
			"foo",
			make([]visor.UnconfirmedTxn, 0),
			errors.New("gateway.GetWalletUnconfirmedTxns error"),
			[]visor.UnconfirmedTxn{},
		},
		{
			"200 - OK",
			http.MethodGet,
			"/wallet/transactions",
			&httpBody{
				Id: "foo",
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
		handler := http.HandlerFunc(WalletTransactionsHandler(gateway))

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
			if err != nil {
				t.Errorf("fail unmarshal json response while 200 OK. body: %s, err: %s", rr.Body.String(), err)
			}
			require.Equal(t, tc.responseBody, msg, tc.name)
		}
	}
}
