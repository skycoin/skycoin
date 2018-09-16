package api

import (
	"bytes"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestGetBalanceHandler(t *testing.T) {
	type httpBody struct {
		addrs string
	}
	invalidAddr := "invalidAddr"
	validAddr := "2eZYSbzBKJ7QCL4kd5LSqV478rJQGb4UNkf"
	address, err := cipher.DecodeBase58Address(validAddr)
	require.NoError(t, err)
	tt := []struct {
		name                      string
		method                    string
		status                    int
		err                       string
		httpBody                  *httpBody
		getBalanceOfAddrsArg      []cipher.Address
		getBalanceOfAddrsResponse []wallet.BalancePair
		getBalanceOfAddrsError    error
		httpResponse              readable.BalancePair
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - invalid address",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - address invalidAddr is invalid: Invalid base58 character",
			httpBody: &httpBody{
				addrs: invalidAddr,
			},
		},
		{
			name:     "400 - no addresses",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - addrs is required",
			httpBody: &httpBody{},
		},
		{
			name:   "500 - GetBalanceOfAddrsError",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.GetBalanceOfAddrs failed: GetBalanceOfAddrsError",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg:   []cipher.Address{address},
			getBalanceOfAddrsError: errors.New("GetBalanceOfAddrsError"),
		},
		{
			name:   "500 - balance Confirmed coins uint64 addition overflow",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - uint64 addition overflow",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
			},
		},
		{
			name:   "500 - balance Predicted coins uint64 addition overflow",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - uint64 addition overflow",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: math.MaxInt64 + 1, Hours: 0},
				},
			},
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			status: http.StatusOK,
			err:    "200 - OK",
			httpBody: &httpBody{
				addrs: validAddr,
			},
			getBalanceOfAddrsArg: []cipher.Address{address},
			getBalanceOfAddrsResponse: []wallet.BalancePair{
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
				{
					Confirmed: wallet.Balance{Coins: 0, Hours: 0},
					Predicted: wallet.Balance{Coins: 0, Hours: 0},
				},
			},
			httpResponse: readable.BalancePair{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			endpoint := "/api/v1/balance"
			gateway.On("GetBalanceOfAddrs", tc.getBalanceOfAddrsArg).Return(tc.getBalanceOfAddrsResponse, tc.getBalanceOfAddrsError)

			v := url.Values{}
			if tc.httpBody != nil {
				if tc.httpBody.addrs != "" {
					v.Add("addrs", tc.httpBody.addrs)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, &CSRFStore{}, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg readable.BalancePair
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.httpResponse, msg, tc.name)
			}
		})
	}
}

func TestWalletSpendHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
		Password string
	}

	type balanceResult struct {
		BalancePair wallet.BalancePair
		Addresses   wallet.AddressBalances
	}

	tt := []struct {
		name                          string
		method                        string
		body                          *httpBody
		status                        int
		err                           string
		walletID                      string
		coins                         uint64
		dst                           string
		password                      string
		gatewaySpendResult            *coin.Transaction
		gatewaySpendErr               error
		gatewayGetWalletBalanceResult balanceResult
		gatewayBalanceErr             error
		spendResult                   *SpendResult
		csrfDisabled                  bool
	}{
		{
			name:     "405",
			method:   http.MethodGet,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodPost,
			body:     &httpBody{},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "0",
		},
		{
			name:   "400 - no dst",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing destination address \"dst\"",
			walletID: "0",
		},
		{
			name:   "400 - bad dst addr",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      " 2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid destination address: Invalid base58 character",
			walletID: "0",
		},
		{
			name:   "400 - no coins",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - coins is string",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "foo",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - coins is negative value",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "-123",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value",
			walletID: "0",
		},
		{
			name:   "400 - zero coins",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "0",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - invalid \"coins\" value, must > 0",
			walletID: "0",
		},
		{
			name:   "400 - gw spend error txn no fee",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - Transaction has zero coinhour fee",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: fee.ErrTxnNoFee,
			spendResult: &SpendResult{
				Error: fee.ErrTxnNoFee.Error(),
			},
		},
		{
			name:   "400 - gw spend error spending unconfirmed",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - please spend after your pending transaction is confirmed",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrSpendingUnconfirmed,
			spendResult: &SpendResult{
				Error: wallet.ErrSpendingUnconfirmed.Error(),
			},
		},
		{
			name:   "400 - gw spend error insufficient balance",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusBadRequest,
			err:             "400 Bad Request - balance is not sufficient",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrInsufficientBalance,
			spendResult: &SpendResult{
				Error: wallet.ErrInsufficientBalance.Error(),
			},
		},
		{
			name:   "404 - gw spend error wallet not exist",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusNotFound,
			err:             "404 Not Found",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrWalletNotExist,
			spendResult: &SpendResult{
				Error: wallet.ErrWalletNotExist.Error(),
			},
		},
		{
			name:   "500 - gw spend error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusInternalServerError,
			err:             "500 Internal Server Error - Spend error",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: errors.New("Spend error"),
			spendResult: &SpendResult{
				Error: "Spend error",
			},
		},
		{
			name:   "200 - gw GetWalletBalance error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:   http.StatusOK,
			walletID: "1234",
			coins:    12,
			dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{
				In: []cipher.SHA256{cipher.MustSHA256FromHex("78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e")},
			},
			gatewayBalanceErr: errors.New("GetWalletBalance error"),
			spendResult: &SpendResult{
				Error: "gateway.GetWalletBalance failed: GetWalletBalance error",
				Transaction: &readable.Transaction{
					Sigs:      []string{},
					In:        []string{"78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e"},
					Out:       []readable.TransactionOutput{},
					Hash:      "110d27c6a0917ec3e3741a7fc5732996542d68a4c61b593335e1f0f1c071ba95",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "123",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			walletID:        "123",
			coins:           12,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendErr: wallet.ErrWalletAPIDisabled,
			spendResult: &SpendResult{
				Error: wallet.ErrWalletAPIDisabled.Error(),
			},
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:   http.StatusOK,
			walletID: "1234",
			coins:    12,
			dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{
				In: []cipher.SHA256{cipher.MustSHA256FromHex("78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e")},
			},
			spendResult: &SpendResult{
				Balance: &readable.BalancePair{},
				Transaction: &readable.Transaction{
					Length:    0,
					Type:      0,
					Hash:      "110d27c6a0917ec3e3741a7fc5732996542d68a4c61b593335e1f0f1c071ba95",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Timestamp: 0,
					Sigs:      []string{},
					In:        []string{"78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e"},
					Out:       []readable.TransactionOutput{},
				},
			},
		},
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "1234",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "12",
			},
			status:   http.StatusOK,
			walletID: "1234",
			coins:    12,
			dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			gatewaySpendResult: &coin.Transaction{
				In: []cipher.SHA256{cipher.MustSHA256FromHex("78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e")},
			},
			spendResult: &SpendResult{
				Balance: &readable.BalancePair{},
				Transaction: &readable.Transaction{
					Length:    0,
					Type:      0,
					Hash:      "110d27c6a0917ec3e3741a7fc5732996542d68a4c61b593335e1f0f1c071ba95",
					InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
					Timestamp: 0,
					Sigs:      []string{},
					In:        []string{"78877fa898f0b4c45c9c33ae941e40617ad7c8657a307db62bc5691f92f4f60e"},
					Out:       []readable.TransactionOutput{},
				},
			},
			csrfDisabled: true,
		},
		{
			name:   "400 - missing password",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
			},
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrMissingPassword,
			err:             "400 Bad Request - missing password",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrMissingPassword.Error(),
			},
		},
		{
			name:   "401 Unauthorized - invalid password",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
				Password: "pwd",
			},
			password:        "pwd",
			status:          http.StatusUnauthorized,
			gatewaySpendErr: wallet.ErrInvalidPassword,
			err:             "401 Unauthorized - invalid password",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrInvalidPassword.Error(),
			},
		},
		{
			name:   "400 - wallet is encrypted",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
			},
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrWalletEncrypted,
			err:             "400 Bad Request - wallet is encrypted",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrWalletEncrypted.Error(),
			},
		},
		{
			name:   "400 - wallet is not encrypted",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "wallet.wlt",
				Dst:      "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
				Coins:    "1",
				Password: "pwd",
			},
			password:        "pwd",
			status:          http.StatusBadRequest,
			gatewaySpendErr: wallet.ErrWalletNotEncrypted,
			err:             "400 Bad Request - wallet is not encrypted",
			walletID:        "wallet.wlt",
			coins:           1,
			dst:             "2konv5no3DZvSMxf2GPVtAfZinfwqCGhfVQ",
			spendResult: &SpendResult{
				Error: wallet.ErrWalletNotEncrypted.Error(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.gatewaySpendResult == nil {
				tc.gatewaySpendResult = &coin.Transaction{}
			}

			gateway := &MockGatewayer{}

			if tc.dst != "" {
				addr, err := cipher.DecodeBase58Address(tc.dst)
				require.NoError(t, err)
				gateway.On("Spend", tc.walletID, []byte(tc.password), tc.coins, addr).Return(tc.gatewaySpendResult, tc.gatewaySpendErr)
			}

			gateway.On("GetWalletBalance", tc.walletID).Return(tc.gatewayGetWalletBalanceResult.BalancePair,
				tc.gatewayGetWalletBalanceResult.Addresses, tc.gatewayBalanceErr)

			endpoint := "/api/v1/wallet/spend"

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
				if tc.body.Password != "" {
					v.Add("password", tc.body.Password)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()))
				if status == http.StatusUnauthorized {
					require.Equal(t, HTTP401AuthHeader, rr.Header().Get("WWW-Authenticate"))
				}
			} else {
				var msg SpendResult
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, *tc.spendResult, msg)
			}
		})
	}
}

func TestWalletGet(t *testing.T) {
	entries, resEntries := makeEntries([]byte("seed"), 5)
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
	}

	tt := []struct {
		name                   string
		method                 string
		body                   *httpBody
		status                 int
		err                    string
		walletID               string
		gatewayGetWalletResult wallet.Wallet
		responseBody           WalletResponse
		gatewayGetWalletErr    error
	}{
		{
			name:     "405",
			method:   http.MethodPost,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "",
		},
		{
			name:   "400 - error from the `gateway.GetWallet(wltID)`",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "123",
			},
			status:              http.StatusBadRequest,
			err:                 "400 Bad Request - wallet 123 doesn't exist",
			walletID:            "123",
			gatewayGetWalletErr: errors.New("wallet 123 doesn't exist"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "1234",
			},
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "1234",
			gatewayGetWalletResult: wallet.Wallet{
				Meta:    map[string]string{"seed": "seed", "lastSeed": "seed"},
				Entries: []wallet.Entry{},
			},
			gatewayGetWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "1234",
			},
			status:   http.StatusOK,
			walletID: "1234",
			gatewayGetWalletResult: wallet.Wallet{
				Meta:    map[string]string{"seed": "seed", "lastSeed": "seed"},
				Entries: cloneEntries(entries),
			},
			responseBody: WalletResponse{Entries: resEntries[:]},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetWallet", tc.walletID).Return(&tc.gatewayGetWalletResult, tc.gatewayGetWalletErr)

			v := url.Values{}

			endpoint := "/api/v1/wallet"

			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()),
					"got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var rlt WalletResponse
				err = json.Unmarshal(rr.Body.Bytes(), &rlt)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, rlt)
			}
		})
	}
}

func TestWalletBalanceHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Dst      string
		Coins    string
	}

	type balanceResult struct {
		BalancePair wallet.BalancePair
		Addresses   wallet.AddressBalances
	}

	tt := []struct {
		name                          string
		method                        string
		body                          *httpBody
		status                        int
		err                           string
		walletID                      string
		gatewayGetWalletBalanceResult balanceResult
		gatewayBalanceErr             error
		result                        *readable.BalancePair
	}{
		{
			name:     "405",
			method:   http.MethodPost,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "0",
		},
		{
			name:     "400 - no walletID",
			method:   http.MethodGet,
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing wallet id",
			walletID: "0",
		},
		{
			name:   "404 - gw `wallet doesn't exist` error",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "notFoundId",
			},
			status:                        http.StatusNotFound,
			err:                           "404 Not Found",
			walletID:                      "notFoundId",
			gatewayGetWalletBalanceResult: balanceResult{},
			gatewayBalanceErr:             wallet.ErrWalletNotExist,
			result: &readable.BalancePair{
				Confirmed: readable.Balance{Coins: 0, Hours: 0},
				Predicted: readable.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			name:   "500 - gw other error",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "someId",
			},
			status:                        http.StatusInternalServerError,
			err:                           "500 Internal Server Error - gatewayBalanceError",
			walletID:                      "someId",
			gatewayGetWalletBalanceResult: balanceResult{},
			gatewayBalanceErr:             errors.New("gatewayBalanceError"),
			result: &readable.BalancePair{
				Confirmed: readable.Balance{Coins: 0, Hours: 0},
				Predicted: readable.Balance{Coins: 0, Hours: 0},
			},
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:                        http.StatusForbidden,
			err:                           "403 Forbidden",
			walletID:                      "foo",
			gatewayGetWalletBalanceResult: balanceResult{},
			gatewayBalanceErr:             wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusOK,
			err:      "",
			walletID: "foo",
			result:   &readable.BalancePair{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetWalletBalance", tc.walletID).Return(tc.gatewayGetWalletBalanceResult.BalancePair,
				tc.gatewayGetWalletBalanceResult.Addresses, tc.gatewayBalanceErr)

			endpoint := "/api/v1/wallet/balance"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}
			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)
			if status != tc.status {
				t.Errorf("got `%v` want `%v`", status, tc.status)
			}
			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg readable.BalancePair
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, &msg, tc.name)
			}
		})
	}
}

func TestUpdateWalletLabelHandler(t *testing.T) {
	type httpBody struct {
		WalletID string
		Label    string
	}

	tt := []struct {
		name                        string
		method                      string
		body                        *httpBody
		status                      int
		err                         string
		walletID                    string
		label                       string
		gatewayUpdateWalletLabelErr error
		responseBody                string
	}{
		{
			name:   "405",
			method: http.MethodGet,
			body:   &httpBody{},
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			body:   &httpBody{},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:   "400 - missing label",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
			},
			status:   http.StatusBadRequest,
			err:      "400 Bad Request - missing label",
			walletID: "foo",
		},
		{
			name:   "404 - gateway.UpdateWalletLabel ErrWalletNotExist",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:                      http.StatusNotFound,
			err:                         "404 Not Found",
			walletID:                    "foo",
			label:                       "label",
			gatewayUpdateWalletLabelErr: wallet.ErrWalletNotExist,
		},
		{
			name:   "500 - gateway.UpdateWalletLabel error",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:                      http.StatusInternalServerError,
			err:                         "500 Internal Server Error - gateway.UpdateWalletLabel error",
			walletID:                    "foo",
			label:                       "label",
			gatewayUpdateWalletLabelErr: errors.New("gateway.UpdateWalletLabel error"),
		},
		{
			name:   "403 Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:                      http.StatusForbidden,
			err:                         "403 Forbidden",
			walletID:                    "foo",
			label:                       "label",
			gatewayUpdateWalletLabelErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 OK",
			method: http.MethodPost,
			body: &httpBody{
				WalletID: "foo",
				Label:    "label",
			},
			status:                      http.StatusOK,
			err:                         "",
			walletID:                    "foo",
			label:                       "label",
			gatewayUpdateWalletLabelErr: nil,
			responseBody:                "\"success\"",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("UpdateWalletLabel", tc.walletID, tc.label).Return(tc.gatewayUpdateWalletLabelErr)

			endpoint := "/api/v1/wallet/update"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.WalletID != "" {
					v.Add("id", tc.body.WalletID)
				}
				if tc.body.Label != "" {
					v.Add("label", tc.body.Label)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				require.Equal(t, tc.responseBody, rr.Body.String(), tc.name)
			}
		})
	}
}

func TestWalletTransactionsHandler(t *testing.T) {
	type httpBody struct {
		walletID string
		verbose  string
	}

	uTxn := &visor.UnconfirmedTransaction{
		Transaction: coin.Transaction{
			In: []cipher.SHA256{testutil.RandSHA256(t)},
		},
	}

	unconfirmedTxn, err := readable.NewUnconfirmedTransaction(uTxn)
	require.NoError(t, err)

	unconfirmedTxnVerbose, err := readable.NewUnconfirmedTransactionVerbose(uTxn, []visor.TransactionInput{
		visor.TransactionInput{},
	})
	require.NoError(t, err)

	tt := []struct {
		name                                         string
		method                                       string
		body                                         *httpBody
		status                                       int
		err                                          string
		walletID                                     string
		verbose                                      bool
		gatewayGetWalletUnconfirmedTxnsResult        []visor.UnconfirmedTransaction
		gatewayGetWalletUnconfirmedTxnsErr           error
		gatewayGetWalletUnconfirmedTxnsVerboseResult []readable.UnconfirmedTransactionVerbose
		gatewayGetWalletUnconfirmedTxnsVerboseErr    error
		responseBody                                 interface{}
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},

		{
			name:   "400 - missing wallet id",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},

		{
			name:   "400 - invalid verbose",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			body: &httpBody{
				verbose: "foo",
			},
			err: "400 Bad Request - Invalid value for verbose",
		},

		{
			name:   "500 - gateway.GetWalletUnconfirmedTransactions error",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
			},
			status:                             http.StatusInternalServerError,
			err:                                "500 Internal Server Error - gateway.GetWalletUnconfirmedTransactions error",
			walletID:                           "foo",
			gatewayGetWalletUnconfirmedTxnsErr: errors.New("gateway.GetWalletUnconfirmedTransactions error"),
		},

		{
			name:   "500 - gateway.GetWalletUnconfirmedTransactionsVerbose error",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
				verbose:  "1",
			},
			verbose:  true,
			status:   http.StatusInternalServerError,
			err:      "500 Internal Server Error - gateway.GetWalletUnconfirmedTransactionsVerbose error",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsVerboseErr: errors.New("gateway.GetWalletUnconfirmedTransactionsVerbose error"),
		},

		{
			name:   "404 - wallet doesn't exist",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
			},
			status:                             http.StatusNotFound,
			err:                                "404 Not Found",
			walletID:                           "foo",
			gatewayGetWalletUnconfirmedTxnsErr: wallet.ErrWalletNotExist,
		},

		{
			name:   "404 - wallet doesn't exist verbose",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
				verbose:  "1",
			},
			verbose:  true,
			status:   http.StatusNotFound,
			err:      "404 Not Found",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsVerboseErr: wallet.ErrWalletNotExist,
		},

		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
			},
			status:                             http.StatusForbidden,
			err:                                "403 Forbidden",
			walletID:                           "foo",
			gatewayGetWalletUnconfirmedTxnsErr: wallet.ErrWalletAPIDisabled,
		},

		{
			name:   "403 - Forbidden - wallet API disabled verbose",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
				verbose:  "1",
			},
			verbose:  true,
			status:   http.StatusForbidden,
			err:      "403 Forbidden",
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsVerboseErr: wallet.ErrWalletAPIDisabled,
		},

		{
			name:   "200 - OK",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
			},
			status:                                http.StatusOK,
			walletID:                              "foo",
			gatewayGetWalletUnconfirmedTxnsResult: make([]visor.UnconfirmedTransaction, 1),
			responseBody: UnconfirmedTxnsResponse{
				Transactions: []readable.UnconfirmedTransactions{
					*unconfirmedTxn,
				},
			},
		},

		{
			name:   "200 - OK verbose",
			method: http.MethodGet,
			body: &httpBody{
				walletID: "foo",
				verbose:  "1",
			},
			verbose:  true,
			status:   http.StatusOK,
			walletID: "foo",
			gatewayGetWalletUnconfirmedTxnsVerboseResult: make([]readable.UnconfirmedTransactionVerbose, 1),
			responseBody: UnconfirmedTxnsVerboseResponse{
				Transactions: []readable.UnconfirmedTransactionVerbose{
					*unconfirmedTxnVerbose,
				},
			},
		},
	}

	for _, tc := range tt {
		gateway := &MockGatewayer{}
		gateway.On("GetWalletUnconfirmedTransactions", tc.walletID).Return(tc.gatewayGetWalletUnconfirmedTxnsResult, tc.gatewayGetWalletUnconfirmedTxnsErr)
		gateway.On("GetWalletUnconfirmedTransactionsVerbose", tc.walletID).Return(tc.gatewayGetWalletUnconfirmedTxnsVerboseResult, tc.gatewayGetWalletUnconfirmedTxnsVerboseErr)

		endpoint := "/api/v1/wallet/transactions"

		v := url.Values{}
		if tc.body != nil {
			if tc.body.walletID != "" {
				v.Add("id", tc.body.walletID)
			}
		}
		if len(v) > 0 {
			endpoint += "?" + v.Encode()
		}
		req, err := http.NewRequest(tc.method, endpoint, nil)
		require.NoError(t, err)

		csrfStore := &CSRFStore{
			Enabled: true,
		}
		setCSRFParameters(csrfStore, tokenValid, req)

		rr := httptest.NewRecorder()
		handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
				strings.TrimSpace(rr.Body.String()), status, tc.err)
			return
		}

		if tc.verbose {
			var msg UnconfirmedTxnsVerboseResponse
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			// require.Equal on whole response might result in flaky tests as there is a time field attached to unconfirmed txn response
			require.IsType(t, msg, tc.responseBody)
			require.Len(t, msg.Transactions, 1)
			require.Equal(t, msg.Transactions[0].Transaction, tc.responseBody.(UnconfirmedTxnsVerboseResponse).Transactions[0].Transaction)
		} else {
			var msg UnconfirmedTxnsResponse
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			// require.Equal on whole response might result in flaky tests as there is a time field attached to unconfirmed txn response
			require.IsType(t, msg, tc.responseBody)
			require.Len(t, msg.Transactions, 1)
			require.Equal(t, msg.Transactions[0].Transaction, tc.responseBody.(UnconfirmedTxnsResponse).Transactions[0].Transaction)
		}
	}
}

func TestWalletCreateHandler(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type httpBody struct {
		Seed     string
		Label    string
		ScanN    string
		Encrypt  bool
		Password string
	}
	tt := []struct {
		name                      string
		method                    string
		body                      *httpBody
		status                    int
		err                       string
		wltName                   string
		options                   wallet.Options
		gatewayCreateWalletResult wallet.Wallet
		gatewayCreateWalletErr    error
		responseBody              WalletResponse
		csrfDisabled              bool
	}{
		{
			name:    "405",
			method:  http.MethodGet,
			status:  http.StatusMethodNotAllowed,
			err:     "405 Method Not Allowed",
			wltName: "foo",
		},
		{
			name:    "400 - missing seed",
			method:  http.MethodPost,
			body:    &httpBody{},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - missing seed",
			wltName: "foo",
		},
		{
			name:   "400 - missing label",
			method: http.MethodPost,
			body: &httpBody{
				Seed: "foo",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - missing label",
			wltName: "foo",
		},
		{
			name:   "400 - invalid scan value",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "bad scanN",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - invalid scan value",
			wltName: "foo",
		},
		{
			name:   "400 - scan must be > 0",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "0",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - scan must be > 0",
			wltName: "foo",
		},
		{
			name:   "400 - seed in use",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "1",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - a wallet already exists with this seed",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
			},
			gatewayCreateWalletErr: wallet.ErrSeedUsed,
		},
		{
			name:   "500 - gateway.CreateWallet error",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "1",
			},
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.CreateWallet error",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
			},
			gatewayCreateWalletErr: errors.New("gateway.CreateWallet error"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusForbidden,
			err:     "403 Forbidden",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
				Entries: cloneEntries(entries),
			},
			responseBody: WalletResponse{
				Meta: readable.WalletMeta{
					Filename: "filename",
				},
				Entries: responseEntries[:],
			},
		},
		// CSRF Tests
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				Seed:  "foo",
				Label: "bar",
				ScanN: "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Password: []byte{},
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename": "filename",
				},
			},
			responseBody: WalletResponse{
				Meta: readable.WalletMeta{
					Filename: "filename",
				},
			},
			csrfDisabled: true,
		},
		{
			name:   "200 - OK - Encrypted",
			method: http.MethodPost,
			body: &httpBody{
				Seed:     "foo",
				Label:    "bar",
				Encrypt:  true,
				Password: "pwd",
				ScanN:    "2",
			},
			status:  http.StatusOK,
			err:     "",
			wltName: "filename",
			options: wallet.Options{
				Label:    "bar",
				Seed:     "foo",
				Encrypt:  true,
				Password: []byte("pwd"),
				ScanN:    2,
			},
			gatewayCreateWalletResult: wallet.Wallet{
				Meta: map[string]string{
					"filename":  "filename",
					"label":     "bar",
					"encrypted": "true",
					"secrets":   "secrets",
				},
			},
			responseBody: WalletResponse{
				Meta: readable.WalletMeta{
					Filename:  "filename",
					Label:     "bar",
					Encrypted: true,
				},
			},
		},
		{
			name:   "400 Bad request - encrypt without password",
			method: http.MethodPost,
			body: &httpBody{
				Seed:    "foo",
				Label:   "bar",
				Encrypt: true,
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing password",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			if tc.options.ScanN == 0 {
				tc.options.ScanN = 1
			}
			gateway.On("CreateWallet", "", tc.options).Return(&tc.gatewayCreateWalletResult, tc.gatewayCreateWalletErr)

			endpoint := "/api/v1/wallet/create"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.Seed != "" {
					v.Add("seed", tc.body.Seed)
				}
				if tc.body.Label != "" {
					v.Add("label", tc.body.Label)
				}
				if tc.body.ScanN != "" {
					v.Add("scan", tc.body.ScanN)
				}

				if tc.body.Encrypt {
					v.Add("encrypt", strconv.FormatBool(tc.body.Encrypt))
				}

				if tc.body.Password != "" {
					v.Add("password", tc.body.Password)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				body := strings.TrimSpace(rr.Body.String())
				require.Equal(t, tc.err, body, "got `%v`| %d, want `%v`", body, status, tc.err)
				return
			}

			var msg WalletResponse
			err = json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.responseBody, msg, tc.name)
		})
	}
}

func TestWalletNewSeed(t *testing.T) {
	type httpBody struct {
		Entropy string
	}
	tt := []struct {
		name      string
		method    string
		body      *httpBody
		status    int
		err       string
		entropy   string
		resultLen int
	}{
		{
			name:   "405",
			method: http.MethodPut,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - invalid entropy type",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "xx",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - invalid entropy",
			entropy: "xx",
		},
		{
			name:   "400 - `wrong entropy length` error",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "200",
			},
			status:  http.StatusBadRequest,
			err:     "400 Bad Request - entropy length must be 128 or 256",
			entropy: "200",
		},
		{
			name:      "200 - OK with no entropy",
			method:    http.MethodGet,
			body:      &httpBody{},
			status:    http.StatusOK,
			entropy:   "128",
			resultLen: 12,
		},
		{
			name:   "200 - OK | 12 word seed",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "128",
			},
			status:    http.StatusOK,
			entropy:   "128",
			resultLen: 12,
		},
		{
			name:   "200 - OK | 24 word seed",
			method: http.MethodGet,
			body: &httpBody{
				Entropy: "256",
			},
			status:    http.StatusOK,
			entropy:   "256",
			resultLen: 24,
		},
	}

	// Loop over each test case
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}

			endpoint := "/api/v1/wallet/newSeed"

			// Add request parameters to url
			v := url.Values{}
			if tc.body != nil {
				if tc.body.Entropy != "" {
					v.Add("entropy", tc.body.Entropy)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` expected `%v`", tc.name, status, tc.status)
			if status != tc.status {
				t.Errorf("got `%v` want `%v`", status, tc.status)
			}
			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, expected `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg struct {
					Seed string `json:"seed"`
				}
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				// check that expected length is equal to response length
				require.Equal(t, tc.resultLen, len(strings.Fields(msg.Seed)), tc.name)
			}
		})
	}
}

func TestGetWalletSeed(t *testing.T) {

	tt := []struct {
		name              string
		method            string
		wltID             string
		password          string
		gatewayReturnArgs []interface{}
		expectStatus      int
		expectSeed        string
		expectErr         string
		csrfDisabled      bool
	}{
		{
			name:     "200 - OK",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"seed",
				nil,
			},
			expectStatus: http.StatusOK,
			expectSeed:   "seed",
		},
		{
			name:     "200 - OK - CSRF disabled",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"seed",
				nil,
			},
			expectStatus: http.StatusOK,
			expectSeed:   "seed",
			csrfDisabled: true,
		},
		{
			name:     "400 - missing wallet id ",
			method:   http.MethodPost,
			wltID:    "",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"seed",
				nil,
			},
			expectStatus: http.StatusBadRequest,
			expectErr:    "400 Bad Request - missing wallet id",
		},
		{
			name:     "400 - missing password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "",
			gatewayReturnArgs: []interface{}{
				"",
				wallet.ErrMissingPassword,
			},
			expectStatus: http.StatusBadRequest,
			expectErr:    "400 Bad Request - missing password",
		},
		{
			name:     "401 Unauthorized - Invalid password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"",
				wallet.ErrInvalidPassword,
			},
			expectStatus: http.StatusUnauthorized,
			expectErr:    "401 Unauthorized - invalid password",
		},
		{
			name:     "400 - wallet not encrypted",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"",
				wallet.ErrWalletNotEncrypted,
			},
			expectStatus: http.StatusBadRequest,
			expectErr:    "400 Bad Request - wallet is not encrypted",
		},
		{
			name:     "404 - wallet does not exist",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturnArgs: []interface{}{
				"",
				wallet.ErrWalletNotExist,
			},
			expectStatus: http.StatusNotFound,
			expectErr:    "404 Not Found",
		},
		{
			name:         "405 - Method Not Allowed",
			method:       http.MethodGet,
			expectStatus: http.StatusMethodNotAllowed,
			expectErr:    "405 Method Not Allowed",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("GetWalletSeed", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturnArgs...)

			endpoint := "/api/v1/wallet/seed"

			v := url.Values{}
			v.Add("id", tc.wltID)
			if len(tc.password) > 0 {
				v.Add("password", tc.password)
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.expectStatus, status)

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
				if status == http.StatusUnauthorized {
					require.Equal(t, HTTP401AuthHeader, rr.Header().Get("WWW-Authenticate"))
				}
			} else {
				var r struct {
					Seed string `json:"seed"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &r)
				require.NoError(t, err)
				require.Equal(t, tc.expectSeed, r.Seed)
			}
		})
	}
}

func TestWalletNewAddressesHandler(t *testing.T) {
	type httpBody struct {
		ID       string
		Num      string
		Password string
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
		body                      *httpBody
		status                    int
		err                       string
		walletID                  string
		n                         uint64
		password                  string
		gatewayNewAddressesResult []cipher.Address
		gatewayNewAddressesErr    error
		responseBody              Addresses
		csrfDisabled              bool
	}{
		{
			name:   "405",
			method: http.MethodGet,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:   "400 - invalid num value",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "bar",
			},
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid num value",
		},
		{
			name:   "400 - gateway.NewAddresses error",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                 http.StatusBadRequest,
			err:                    "400 Bad Request - gateway.NewAddresses error",
			walletID:               "foo",
			n:                      1,
			gatewayNewAddressesErr: errors.New("gateway.NewAddresses error"),
		},
		{
			name:   "403 - Forbidden - wallet API disabled",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                 http.StatusForbidden,
			err:                    "403 Forbidden",
			walletID:               "foo",
			n:                      1,
			gatewayNewAddressesErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:   "400 Bad Request - missing password",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                 http.StatusBadRequest,
			err:                    "400 Bad Request - missing password",
			walletID:               "foo",
			n:                      1,
			gatewayNewAddressesErr: wallet.ErrMissingPassword,
		},
		{
			name:   "401 Unauthorized - Invalid password",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                 http.StatusUnauthorized,
			err:                    "401 Unauthorized - invalid password",
			walletID:               "foo",
			n:                      1,
			gatewayNewAddressesErr: wallet.ErrInvalidPassword,
		},
		{
			name:   "200 - OK",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                    http.StatusOK,
			walletID:                  "foo",
			n:                         1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
		},
		{
			name:   "200 - OK with password",
			method: http.MethodPost,
			body: &httpBody{
				ID:       "foo",
				Num:      "1",
				Password: "pwd",
			},
			status:                    http.StatusOK,
			walletID:                  "foo",
			n:                         1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
		},
		{
			name:   "200 - OK empty addresses",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "0",
			},
			status:                    http.StatusOK,
			walletID:                  "foo",
			n:                         0,
			gatewayNewAddressesResult: emptyAddrs,
			responseBody:              responseEmptyAddresses,
		},
		{
			name:   "200 - OK - CSRF disabled",
			method: http.MethodPost,
			body: &httpBody{
				ID:  "foo",
				Num: "1",
			},
			status:                    http.StatusOK,
			walletID:                  "foo",
			n:                         1,
			gatewayNewAddressesResult: addrs,
			responseBody:              responseAddresses,
			csrfDisabled:              true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("NewAddresses", tc.walletID, []byte(tc.password), tc.n).Return(tc.gatewayNewAddressesResult, tc.gatewayNewAddressesErr)

			endpoint := "/api/v1/wallet/newAddress"

			v := url.Values{}
			if tc.body != nil {
				if tc.body.ID != "" {
					v.Add("id", tc.body.ID)
				}
				if tc.body.Num != "" {
					v.Add("num", tc.body.Num)
				}
			}

			req, err := http.NewRequest(tc.method, endpoint, bytes.NewBufferString(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
				if status == http.StatusUnauthorized {
					require.Equal(t, HTTP401AuthHeader, rr.Header().Get("WWW-Authenticate"))
				}
			} else {
				var msg Addresses
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.responseBody, msg, tc.name)
			}
		})
	}
}

func TestGetWalletFolderHandler(t *testing.T) {
	tt := []struct {
		name                 string
		method               string
		status               int
		err                  string
		getWalletDirResponse string
		getWalletDirErr      error
		httpResponse         WalletFolder
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:                 "200",
			method:               http.MethodGet,
			status:               http.StatusOK,
			getWalletDirResponse: "/wallet/folder/address",
			httpResponse: WalletFolder{
				Address: "/wallet/folder/address",
			},
		},
		{
			name:            "403 - wallet API disabled",
			method:          http.MethodGet,
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			getWalletDirErr: wallet.ErrWalletAPIDisabled,
		},
	}

	for _, tc := range tt {
		gateway := &MockGatewayer{}
		gateway.On("GetWalletDir").Return(tc.getWalletDirResponse, tc.getWalletDirErr)

		endpoint := "/api/v1/wallets/folderName"

		req, err := http.NewRequest(tc.method, endpoint, nil)
		require.NoError(t, err)

		csrfStore := &CSRFStore{
			Enabled: true,
		}
		setCSRFParameters(csrfStore, tokenValid, req)

		rr := httptest.NewRecorder()
		handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
				strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg WalletFolder
			err := json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}

func TestGetWallets(t *testing.T) {
	var pubkeys []cipher.PubKey
	var seckeys []cipher.SecKey
	var addrs []cipher.Address

	for i := 0; i < 4; i++ {
		pubkey, seckey := cipher.GenerateKeyPair()
		addr := cipher.AddressFromPubKey(pubkey)
		pubkeys = append(pubkeys, pubkey)
		seckeys = append(seckeys, seckey)
		addrs = append(addrs, addr)
	}

	cases := []struct {
		name               string
		method             string
		status             int
		err                string
		getWalletsResponse wallet.Wallets
		getWalletsErr      error
		httpResponse       []*WalletResponse
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:          "403 - wallet API disabled",
			method:        http.MethodGet,
			status:        http.StatusForbidden,
			err:           "403 Forbidden",
			getWalletsErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:               "200 no wallets",
			method:             http.MethodGet,
			status:             http.StatusOK,
			getWalletsResponse: nil,
			httpResponse:       []*WalletResponse{},
		},
		{
			name:               "200 no wallets 2",
			method:             http.MethodGet,
			status:             http.StatusOK,
			getWalletsResponse: wallet.Wallets{},
			httpResponse:       []*WalletResponse{},
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			getWalletsResponse: wallet.Wallets{
				"foofilename": {
					Meta: map[string]string{
						"foo":        "bar",
						"seed":       "fooseed",
						"lastSeed":   "foolastseed",
						"coin":       "foocoin",
						"filename":   "foofilename",
						"label":      "foolabel",
						"type":       "footype",
						"version":    "fooversion",
						"cryptoType": "foocryptotype",
						"tm":         "345678",
						"encrypted":  "true",
					},
					Entries: []wallet.Entry{
						{
							Address: addrs[0],
							Public:  pubkeys[0],
							Secret:  seckeys[0],
						},
					},
				},
				"foofilename2": {
					Meta: map[string]string{
						"foo":        "bar2",
						"seed":       "fooseed2",
						"lastSeed":   "foolastseed2",
						"coin":       "foocoin",
						"filename":   "foofilename2",
						"label":      "foolabel2",
						"type":       "footype",
						"version":    "fooversion",
						"cryptoType": "foocryptotype",
						"tm":         "123456",
						"encrypted":  "false",
					},
					Entries: []wallet.Entry{
						{
							Address: addrs[1],
							Public:  pubkeys[1],
							Secret:  seckeys[1],
						},
					},
				},
				"foofilename3": {
					Meta: map[string]string{
						"foo":        "bar3",
						"seed":       "fooseed3",
						"lastSeed":   "foolastseed3",
						"coin":       "foocoin",
						"filename":   "foofilename3",
						"label":      "foolabel3",
						"type":       "footype",
						"version":    "fooversion",
						"cryptoType": "foocryptotype",
						"tm":         "234567",
						"encrypted":  "true",
					},
					Entries: []wallet.Entry{
						{
							Address: addrs[2],
							Public:  pubkeys[2],
							Secret:  seckeys[2],
						},
						{
							Address: addrs[3],
							Public:  pubkeys[3],
							Secret:  seckeys[3],
						},
					},
				},
			},
			httpResponse: []*WalletResponse{
				{
					Meta: readable.WalletMeta{
						Coin:       "foocoin",
						Filename:   "foofilename2",
						Label:      "foolabel2",
						Type:       "footype",
						Version:    "fooversion",
						CryptoType: "foocryptotype",
						Timestamp:  123456,
						Encrypted:  false,
					},
					Entries: []readable.WalletEntry{
						{
							Address: addrs[1].String(),
							Public:  pubkeys[1].Hex(),
						},
					},
				},
				{
					Meta: readable.WalletMeta{
						Coin:       "foocoin",
						Filename:   "foofilename3",
						Label:      "foolabel3",
						Type:       "footype",
						Version:    "fooversion",
						CryptoType: "foocryptotype",
						Timestamp:  234567,
						Encrypted:  true,
					},
					Entries: []readable.WalletEntry{
						{
							Address: addrs[2].String(),
							Public:  pubkeys[2].Hex(),
						},
						{
							Address: addrs[3].String(),
							Public:  pubkeys[3].Hex(),
						},
					},
				},
				{
					Meta: readable.WalletMeta{
						Coin:       "foocoin",
						Filename:   "foofilename",
						Label:      "foolabel",
						Type:       "footype",
						Version:    "fooversion",
						CryptoType: "foocryptotype",
						Timestamp:  345678,
						Encrypted:  true,
					},
					Entries: []readable.WalletEntry{
						{
							Address: addrs[0].String(),
							Public:  pubkeys[0].Hex(),
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		gateway := &MockGatewayer{}
		gateway.On("GetWallets").Return(tc.getWalletsResponse, tc.getWalletsErr)

		endpoint := "/api/v1/wallets"

		req, err := http.NewRequest(tc.method, endpoint, nil)
		require.NoError(t, err)

		csrfStore := &CSRFStore{
			Enabled: true,
		}
		setCSRFParameters(csrfStore, tokenValid, req)

		rr := httptest.NewRecorder()
		handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

		handler.ServeHTTP(rr, req)

		status := rr.Code
		require.Equal(t, tc.status, status, "got `%v` want `%v`",
			tc.name, status, tc.status)

		if status != http.StatusOK {
			require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
				strings.TrimSpace(rr.Body.String()), status, tc.err)
		} else {
			var msg []*WalletResponse
			err := json.Unmarshal(rr.Body.Bytes(), &msg)
			require.NoError(t, err)
			require.NotNil(t, msg)
			require.Equal(t, tc.httpResponse, msg, tc.name)
		}
	}
}

func TestWalletUnloadHandler(t *testing.T) {
	tt := []struct {
		name            string
		method          string
		status          int
		err             string
		walletID        string
		unloadWalletErr error
		csrfDisabled    bool
	}{
		{
			name:     "405",
			method:   http.MethodGet,
			status:   http.StatusMethodNotAllowed,
			err:      "405 Method Not Allowed",
			walletID: "wallet.wlt",
		},
		{
			name:   "400 - missing wallet id",
			method: http.MethodPost,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - missing wallet id",
		},
		{
			name:            "403 - Forbidden - wallet API disabled",
			method:          http.MethodPost,
			status:          http.StatusForbidden,
			err:             "403 Forbidden",
			walletID:        "wallet.wlt",
			unloadWalletErr: wallet.ErrWalletAPIDisabled,
		},
		{
			name:     "200 - ok",
			method:   http.MethodPost,
			status:   http.StatusOK,
			walletID: "wallet.wlt",
		},
		{
			name:         "200 - ok, csrf disabled",
			method:       http.MethodPost,
			status:       http.StatusOK,
			walletID:     "wallet.wlt",
			csrfDisabled: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("UnloadWallet", tc.walletID).Return(tc.unloadWalletErr)

			endpoint := "/api/v1/wallet/unload"
			v := url.Values{}
			v.Add("id", tc.walletID)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			}
		})
	}
}

func TestEncryptWallet(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type gatewayReturnPair struct {
		w   *wallet.Wallet
		err error
	}
	tt := []struct {
		name          string
		method        string
		wltID         string
		password      string
		gatewayReturn gatewayReturnPair
		status        int
		expectWallet  WalletResponse
		expectErr     string
	}{
		{
			name:     "200 - OK",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				w: &wallet.Wallet{
					Meta: map[string]string{
						"filename":  "wallet.wlt",
						"seed":      "seed",
						"lastSeed":  "lastSeed",
						"secrets":   "secrets",
						"encrypted": "true",
					},
					Entries: cloneEntries(entries),
				},
			},
			status: http.StatusOK,
			expectWallet: WalletResponse{
				Meta: readable.WalletMeta{
					Filename:  "wallet.wlt",
					Encrypted: true,
				},
				Entries: responseEntries,
			},
		},
		{
			name:     "403 Forbidden",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletAPIDisabled,
			},
			status:    http.StatusForbidden,
			expectErr: "403 Forbidden",
		},
		{
			name:      "405 Method Not Allowed",
			method:    http.MethodGet,
			wltID:     "wallet.wlt",
			password:  "pwd",
			status:    http.StatusMethodNotAllowed,
			expectErr: "405 Method Not Allowed",
		},
		{
			name:     "400 - Missing Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrMissingPassword,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing password",
		},
		{
			name:      "400 - Missing Wallet Id",
			method:    http.MethodPost,
			wltID:     "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing wallet id",
		},
		{
			name:     "401 Unauthorized - Invalid Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrInvalidPassword,
			},
			status:    http.StatusUnauthorized,
			expectErr: "401 Unauthorized - invalid password",
		},
		{
			name:     "404 - Wallet Not Found",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotExist,
			},
			status:    http.StatusNotFound,
			expectErr: "404 Not Found",
		},
		{
			name:     "400 - Wallet Is Encrypted",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletEncrypted,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - wallet is encrypted",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("EncryptWallet", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturn.w, tc.gatewayReturn.err)

			endpoint := "/api/v1/wallet/encrypt"
			v := url.Values{}
			v.Add("id", tc.wltID)
			v.Add("password", tc.password)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: true,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`, body: %v", status, tc.status, rr.Body.String())

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
				if status == http.StatusUnauthorized {
					require.Equal(t, HTTP401AuthHeader, rr.Header().Get("WWW-Authenticate"))
				}
				return
			}

			var rlt WalletResponse
			err = json.NewDecoder(rr.Body).Decode(&rlt)
			require.NoError(t, err)
			require.Equal(t, tc.expectWallet, rlt)
		})
	}
}

func TestDecryptWallet(t *testing.T) {
	entries, responseEntries := makeEntries([]byte("seed"), 5)
	type gatewayReturnPair struct {
		w   *wallet.Wallet
		err error
	}

	tt := []struct {
		name          string
		method        string
		wltID         string
		password      string
		gatewayReturn gatewayReturnPair
		status        int
		expectWallet  WalletResponse
		expectErr     string
		csrfDisabled  bool
	}{
		{
			name:     "200 OK",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				w: &wallet.Wallet{
					Meta: map[string]string{
						"filename":  "wallet",
						"seed":      "seed",
						"lastSeed":  "lastSeed",
						"secrets":   "",
						"encrypted": "false",
					},
					Entries: cloneEntries(entries),
				},
			},
			status: http.StatusOK,
			expectWallet: WalletResponse{
				Meta: readable.WalletMeta{
					Filename:  "wallet",
					Encrypted: false,
				},
				Entries: responseEntries,
			},
		},
		{
			name:     "200 OK CSRF disabled",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				w: &wallet.Wallet{
					Meta: map[string]string{
						"filename":  "wallet",
						"seed":      "seed",
						"lastSeed":  "lastSeed",
						"secrets":   "",
						"encrypted": "false",
					},
					Entries: cloneEntries(entries),
				},
			},
			status: http.StatusOK,
			expectWallet: WalletResponse{
				Meta: readable.WalletMeta{
					Filename:  "wallet",
					Encrypted: false,
				},
				Entries: responseEntries,
			},
			csrfDisabled: true,
		},

		{
			name:     "403 Forbidden",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletAPIDisabled,
			},
			status:    http.StatusForbidden,
			expectErr: "403 Forbidden",
		},
		{
			name:      "405 Method Not Allowed",
			method:    http.MethodGet,
			status:    http.StatusMethodNotAllowed,
			expectErr: "405 Method Not Allowed",
		},
		{
			name:      "400 - Missing Wallet ID",
			method:    http.MethodPost,
			wltID:     "",
			password:  "",
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing wallet id",
		},
		{
			name:     "400 - Missing Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrMissingPassword,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - missing password",
		},
		{
			name:     "400 - Wallet IS Not Encrypted",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotEncrypted,
			},
			status:    http.StatusBadRequest,
			expectErr: "400 Bad Request - wallet is not encrypted",
		},
		{
			name:     "401 Unauthorized - Invalid Password",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrInvalidPassword,
			},
			status:    http.StatusUnauthorized,
			expectErr: "401 Unauthorized - invalid password",
		},
		{
			name:     "404 - Wallet Does Not Exist",
			method:   http.MethodPost,
			wltID:    "wallet.wlt",
			password: "pwd",
			gatewayReturn: gatewayReturnPair{
				err: wallet.ErrWalletNotExist,
			},
			status:    http.StatusNotFound,
			expectErr: "404 Not Found",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &MockGatewayer{}
			gateway.On("DecryptWallet", tc.wltID, []byte(tc.password)).Return(tc.gatewayReturn.w, tc.gatewayReturn.err)

			endpoint := "/api/v1/wallet/decrypt"
			v := url.Values{}
			v.Add("id", tc.wltID)
			v.Add("password", tc.password)

			req, err := http.NewRequest(tc.method, endpoint, strings.NewReader(v.Encode()))
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			setCSRFParameters(csrfStore, tokenValid, req)

			rr := httptest.NewRecorder()
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "wrong status code: got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.expectErr, strings.TrimSpace(rr.Body.String()))
				if status == http.StatusUnauthorized {
					require.Equal(t, HTTP401AuthHeader, rr.Header().Get("WWW-Authenticate"))
				}
				return
			}

			var r WalletResponse
			err = json.NewDecoder(rr.Body).Decode(&r)
			require.NoError(t, err)
			require.Equal(t, tc.expectWallet, r)
		})
	}
}

// makeEntries derives N wallet address entries from given seed
// Returns set of wallet.Entry and wallet.ReadableEntry, the readable
// entries' secrets are removed.
func makeEntries(seed []byte, n int) ([]wallet.Entry, []readable.WalletEntry) { // nolint: unparam
	seckeys := cipher.GenerateDeterministicKeyPairs(seed, n)
	var entries []wallet.Entry
	var responseEntries []readable.WalletEntry
	for i, seckey := range seckeys {
		pubkey := cipher.PubKeyFromSecKey(seckey)
		entries = append(entries, wallet.Entry{
			Address: cipher.AddressFromPubKey(pubkey),
			Public:  pubkey,
			Secret:  seckey,
		})
		responseEntries = append(responseEntries, readable.WalletEntry{
			Address: entries[i].Address.String(),
			Public:  entries[i].Public.Hex(),
		})
	}
	return entries, responseEntries
}

func cloneEntries(es []wallet.Entry) []wallet.Entry {
	var entries []wallet.Entry
	entries = append(entries, es...)
	return entries
}
