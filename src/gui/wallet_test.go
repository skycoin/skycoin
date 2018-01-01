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

func TestWalletNewAddressesHandler(t *testing.T) {
	type httpBody struct {
		Id  string `url:"id,omitempty"`
		Num string `url:"num,omitempty"`
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
			http.MethodPut,
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
				Num: "1",
			},
			http.StatusOK,
			"",
			"foo",
			1,
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
			body, _ := query.Values(tc.body)

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(body.Encode()))
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
