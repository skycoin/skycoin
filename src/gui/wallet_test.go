package gui

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"encoding/json"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/require"
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
			"",
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
			"",
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
			"",
			"1234",
			12,
			"2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			&coin.Transaction{},
			nil,
			wallet.BalancePair{},
			nil,
			&SpendResult{
				Error:   "",
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
		require.Equal(t, status, tc.status, "case: %s, handler returned wrong status code: got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, strings.TrimSpace(rr.Body.String()), tc.err, "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
				tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg SpendResult
			json.Unmarshal(rr.Body.Bytes(), &msg)
			require.Equal(t, *tc.spendResult, msg, "test")
		}
	}

}
