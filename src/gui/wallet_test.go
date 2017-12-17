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

	"bytes"

	"github.com/pkg/errors"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

type httpBody struct {
	Id    string `url:"id,omitempty"`
	Dst   string `url:"dst,omitempty"`
	Coins string `url:"coins,omitempty"`
	Seed  string `url:"seed,omitempty"`
	Label string `url:"label,omitempty"`
	ScanN string `url:"scan,omitempty"`
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

func TestWalletHandler(t *testing.T) {
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
			http.MethodPut,
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
			wltName: tc.wltname,
			scanN:   tc.scnN,
			t:       t,
		}
		gateway.On("CreateWallet", "", tc.options).Return(tc.gatewayCreateWalletResult, tc.gatewayCreateWalletErr)
		gateway.On("ScanAheadWalletAddresses", tc.wltname, tc.scnN-1).Return(tc.scanWalletAddressesResult, tc.scanWalletAddressesError)
		body, _ := query.Values(tc.body)

		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(body.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WalletCreate(gateway))

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg wallet.ReadableWallet
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			if err != nil {
				t.Errorf("fail unmarshal json response while 200 OK. body: %s, err: %s", rr.Body.String(), err)
			}
			require.Equal(t, tc.responseBody, msg, tc.name)
		}
	}
}
