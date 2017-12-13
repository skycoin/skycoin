package gui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/mock"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/skycoin/skycoin/src/visor"
)

var (
	genPublic, genSecret         = cipher.GenerateKeyPair()
	genAddress                   = cipher.AddressFromPubKey(genPublic)
	testMaxSize                  = 1024 * 1024
	GenesisPublic, GenesisSecret = cipher.GenerateKeyPair()
	GenesisAddress               = cipher.AddressFromPubKey(GenesisPublic)
)

const (
	TimeIncrement    uint64 = 3600 * 1000
	GenesisTime      uint64 = 1000
	GenesisCoins     uint64 = 1000e6
	GenesisCoinHours uint64 = 1000 * 1000
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
	//logger.Info("arg[0]: %s, arg[1]: %s", args.Get(0), args.Get(1))
	return args.Get(0).(wallet.BalancePair), args.Error(1)
}

func TestWalletSpendHandler(t *testing.T) {
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
			"400 - no walletId",
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
				Id: "123",
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
				Id:  "123",
				Dst: "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ======",
			},
			http.StatusBadRequest,
			"400 Bad Request - invalid destination address: Invalid address length",
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
				Id:  "123",
				Dst: "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
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
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "foo",
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
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "-123",
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
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "0",
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
			"200 - gw spend error",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "123",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "12",
			},
			http.StatusOK,
			"400 Bad Request - invalid \"coins\" value, must > 0",
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
				Id:    "1234",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "12",
			},
			http.StatusOK,
			"400 Bad Request - Get wallet balance failed: GetWalletBalance error",
			"1234",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			errors.New("GetWalletBalance error"),
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error: "GetWalletBalance error",
			},
		},
		{
			"200 - OK",
			"POST",
			"/wallet/spend",
			&httpBody{
				Id:    "1234",
				Dst:   "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins: "12",
			},
			http.StatusOK,
			"400 Bad Request - Get wallet balance failed: GetWalletBalance error",
			"1234",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			nil,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error:"",
				Balance: &wallet.BalancePair{},
				Transaction: &visor.ReadableTransaction{
					Length:0,
					Type:0,
					Hash: "78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Timestamp:0,
					Sigs:[]string{},
					In: []string{},
					Out:[]visor.ReadableTransactionOutput{},
				},
			},
		},
	}
	for _, tc := range tt {
		gateway := &FakeGateway{
			walletId: tc.walletId,
			t:        t,
		}
		addr, _ := cipher.DecodeBase58Address(tc.dst)
		gateway.On("Spend", tc.walletId, tc.coins, addr).Return(tc.gatewaySpendResult, tc.gatewaySpendErr)
		gateway.On("GetWalletBalance", tc.walletId).Return(tc.gatewayGetWalletBalanceResult, tc.gatewayBalanceErr)
		body, _ := query.Values(tc.body)
		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBufferString(body.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(WalletSpendHandler(gateway))

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
			var msg SpendResult
			json.Unmarshal(rr.Body.Bytes(), &msg)
			assert.Equal(t, *tc.spendResult, msg, "test")
		}
	}

}
