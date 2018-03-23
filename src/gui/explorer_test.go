package gui

import (
	"errors"
	"net/http"
	"testing"

	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/mock"

	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

func makeSuccessCoinSupplyResult(t *testing.T, allUnspents visor.ReadableOutputSet) *CoinSupply {
	unlockedAddrs := visor.GetUnlockedDistributionAddresses()
	var unlockedSupply uint64
	// check confirmed unspents only
	// Search map of unlocked addresses
	// used to filter unspents
	unlockedAddrMap := daemon.MakeSearchMap(unlockedAddrs)
	for _, u := range allUnspents.HeadOutputs {
		// check if address is an unlocked distribution address
		if _, ok := unlockedAddrMap[u.Address]; ok {
			coins, err := droplet.FromString(u.Coins)
			require.NoError(t, err)
			unlockedSupply += coins
		}
	}
	// "total supply" is the number of coins unlocked.
	// Each distribution address was allocated visor.DistributionAddressInitialBalance coins.
	totalSupply := uint64(len(unlockedAddrs)) * visor.DistributionAddressInitialBalance
	totalSupply *= droplet.Multiplier

	// "current supply" is the number of coins distribution from the unlocked pool
	currentSupply := totalSupply - unlockedSupply

	currentSupplyStr, err := droplet.ToString(currentSupply)
	require.NoError(t, err)

	totalSupplyStr, err := droplet.ToString(totalSupply)
	require.NoError(t, err)

	maxSupplyStr, err := droplet.ToString(visor.MaxCoinSupply * droplet.Multiplier)
	require.NoError(t, err)

	// locked distribution addresses
	lockedAddrs := visor.GetLockedDistributionAddresses()
	lockedAddrMap := daemon.MakeSearchMap(lockedAddrs)

	// get total coins hours which excludes locked distribution addresses
	var totalCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		if _, ok := lockedAddrMap[out.Address]; !ok {
			totalCoinHours += out.Hours
		}
	}

	// get current coin hours which excludes all distribution addresses
	var currentCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		// check if address not in locked distribution addresses
		if _, ok := lockedAddrMap[out.Address]; !ok {
			// check if address not in unlocked distribution addresses
			if _, ok := unlockedAddrMap[out.Address]; !ok {
				currentCoinHours += out.Hours
			}
		}
	}

	cs := CoinSupply{
		CurrentSupply:         currentSupplyStr,
		TotalSupply:           totalSupplyStr,
		MaxSupply:             maxSupplyStr,
		CurrentCoinHourSupply: strconv.FormatUint(currentCoinHours, 10),
		TotalCoinHourSupply:   strconv.FormatUint(totalCoinHours, 10),
		UnlockedAddresses:     unlockedAddrs,
		LockedAddresses:       visor.GetLockedDistributionAddresses(),
	}
	return &cs
}

func TestGetTransactionsForAddress(t *testing.T) {
	address := testutil.MakeAddress()
	successAddress := "111111111111111111111691FSP"
	invalidHash := "cafcb"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	tt := []struct {
		name                        string
		method                      string
		status                      int
		err                         string
		addressParam                string
		gatewayGetAddressTxnsResult *visor.TransactionResults
		gatewayGetAddressTxnsErr    error
		gatewayGetUxOutByIDArg      cipher.SHA256
		gatewayGetUxOutByIDResult   *historydb.UxOut
		gatewayGetUxOutByIDErr      error
		result                      []ReadableTransaction
		csrfDisabled                bool
	}{
		{
			name:         "405",
			method:       http.MethodPost,
			status:       http.StatusMethodNotAllowed,
			err:          "405 Method Not Allowed",
			addressParam: "0",
		},
		{
			name:         "400 - address is empty",
			method:       http.MethodGet,
			status:       http.StatusBadRequest,
			err:          "400 Bad Request - address is empty",
			addressParam: "",
		},
		{
			name:         "400 - invalid address",
			method:       http.MethodGet,
			status:       http.StatusBadRequest,
			err:          "400 Bad Request - invalid address",
			addressParam: "badAddress",
		},
		{
			name:                     "500 - gw GetAddressTxns error",
			method:                   http.MethodGet,
			status:                   http.StatusInternalServerError,
			err:                      "500 Internal Server Error",
			addressParam:             address.String(),
			gatewayGetAddressTxnsErr: errors.New("gatewayGetAddressTxnsErr"),
		},
		{
			name:         "500 - cipher.SHA256FromHex(tx.Transaction.In) error",
			method:       http.MethodGet,
			status:       http.StatusInternalServerError,
			err:          "500 Internal Server Error",
			addressParam: address.String(),
			gatewayGetAddressTxnsResult: &visor.TransactionResults{
				Txns: []visor.TransactionResult{
					{
						Transaction: visor.ReadableTransaction{
							In: []string{
								invalidHash,
							},
						},
					},
				},
			},
		},
		{
			name:         "500 - GetUxOutByID error",
			method:       http.MethodGet,
			status:       http.StatusInternalServerError,
			err:          "500 Internal Server Error",
			addressParam: address.String(),
			gatewayGetAddressTxnsResult: &visor.TransactionResults{
				Txns: []visor.TransactionResult{
					{
						Transaction: visor.ReadableTransaction{
							In: []string{
								validHash,
							},
						},
					},
				},
			},
			gatewayGetUxOutByIDArg: testutil.SHA256FromHex(t, validHash),
			gatewayGetUxOutByIDErr: errors.New("gatewayGetUxOutByIDErr"),
		},
		{
			name:         "500 - GetUxOutByID nil result",
			method:       http.MethodGet,
			status:       http.StatusInternalServerError,
			err:          "500 Internal Server Error",
			addressParam: address.String(),
			gatewayGetAddressTxnsResult: &visor.TransactionResults{
				Txns: []visor.TransactionResult{
					{
						Transaction: visor.ReadableTransaction{
							In: []string{
								validHash,
							},
						},
					},
				},
			},
			gatewayGetUxOutByIDArg: testutil.SHA256FromHex(t, validHash),
		},
		{
			name:         "200",
			method:       http.MethodGet,
			status:       http.StatusOK,
			addressParam: address.String(),
			gatewayGetAddressTxnsResult: &visor.TransactionResults{
				Txns: []visor.TransactionResult{
					{
						Transaction: visor.ReadableTransaction{
							In: []string{
								validHash,
							},
						},
					},
				},
			},
			gatewayGetUxOutByIDArg:    testutil.SHA256FromHex(t, validHash),
			gatewayGetUxOutByIDResult: &historydb.UxOut{},
			result: []ReadableTransaction{
				{
					In: []visor.ReadableTransactionInput{
						{
							Hash:    validHash,
							Address: successAddress,
							Coins:   "0.000000",
							Hours:   0,
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/explorer/address"
			gateway := NewGatewayerMock()
			gateway.On("GetAddressTxns", address).Return(tc.gatewayGetAddressTxnsResult, tc.gatewayGetAddressTxnsErr)
			gateway.On("GetUxOutByID", tc.gatewayGetUxOutByIDArg).Return(tc.gatewayGetUxOutByIDResult, tc.gatewayGetUxOutByIDErr)

			v := url.Values{}
			if tc.addressParam != "" {
				v.Add("address", tc.addressParam)
			}

			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)
			rr := httptest.NewRecorder()
			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []ReadableTransaction
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestCoinSupply(t *testing.T) {
	unlockedAddrs := visor.GetUnlockedDistributionAddresses()
	successGatewayGetUnspentOutputsResult := visor.ReadableOutputSet{
		HeadOutputs: visor.ReadableOutputs{
			visor.ReadableOutput{
				Coins: "0",
			},
			visor.ReadableOutput{
				Coins: "0",
			},
		},
	}
	var filterInUnlocked []daemon.OutputsFilter
	filterInUnlocked = append(filterInUnlocked, daemon.FbyAddresses(unlockedAddrs))
	tt := []struct {
		name                           string
		method                         string
		status                         int
		err                            string
		gatewayGetUnspentOutputsArg    []daemon.OutputsFilter
		gatewayGetUnspentOutputsResult *visor.ReadableOutputSet
		gatewayGetUnspentOutputsErr    error
		result                         *CoinSupply
		csrfDisabled                   bool
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "500 - gatewayGetUnspentOutputsErr",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:   "500 - gatewayGetUnspentOutputsErr",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:   "500 - too large HeadOutputs item",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsResult: &visor.ReadableOutputSet{
				HeadOutputs: visor.ReadableOutputs{
					visor.ReadableOutput{
						Coins:   "9223372036854775807",
						Address: unlockedAddrs[0],
					},
					visor.ReadableOutput{
						Coins: "1",
					},
				},
			},
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,

			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsResult: &visor.ReadableOutputSet{
				HeadOutputs: visor.ReadableOutputs{
					visor.ReadableOutput{
						Coins: "0",
					},
					visor.ReadableOutput{
						Coins: "0",
					},
				},
			},
			result: makeSuccessCoinSupplyResult(t, successGatewayGetUnspentOutputsResult),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/coinSupply"
			gateway := NewGatewayerMock()
			gateway.On("GetUnspentOutputs", mock.Anything).Return(tc.gatewayGetUnspentOutputsResult, tc.gatewayGetUnspentOutputsErr)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *CoinSupply
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetRichlist(t *testing.T) {
	type httpParams struct {
		topn                string
		includeDistribution string
	}
	tt := []struct {
		name                     string
		method                   string
		status                   int
		err                      string
		httpParams               *httpParams
		includeDistribution      bool
		gatewayGetRichlistResult visor.Richlist
		gatewayGetRichlistErr    error
		result                   Richlist
		csrfDisabled             bool
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "400 - bad topn param",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid n",
			httpParams: &httpParams{
				topn: "bad topn",
			},
		},
		{
			name:   "400 - include-distribution",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			err:    "400 Bad Request - invalid include-distribution",
			httpParams: &httpParams{
				topn:                "1",
				includeDistribution: "bad include-distribution",
			},
		},
		{
			name:   "500 - gw GetRichlist error",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			httpParams: &httpParams{
				topn:                "1",
				includeDistribution: "false",
			},
			gatewayGetRichlistErr: errors.New("gatewayGetRichlistErr"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			httpParams: &httpParams{
				topn:                "3",
				includeDistribution: "false",
			},
			gatewayGetRichlistResult: visor.Richlist{
				{
					Address: "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
					Coins:   "1000000.000000",
					Locked:  false,
				},
				{
					Address: "27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1",
					Coins:   "500000.000000",
					Locked:  false,
				},
				{
					Address: "2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW",
					Coins:   "500000.000000",
					Locked:  false,
				},
				{
					Address: "2TmvdBWJgxMwGs84R4drS9p5fYkva4dGdfs",
					Coins:   "244458.000000",
					Locked:  false,
				},
				{
					Address: "24gvUHXHtSg5drKiFsMw7iMgoN2PbLub53C",
					Coins:   "195503.000000",
					Locked:  false,
				},
			},
			result: Richlist{
				Richlist: visor.Richlist{
					{
						Address: "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
						Coins:   "1000000.000000",
						Locked:  false,
					},
					{
						Address: "27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1",
						Coins:   "500000.000000",
						Locked:  false,
					},
					{
						Address: "2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW",
						Coins:   "500000.000000",
						Locked:  false,
					},
				},
			},
		},
		{
			name:   "200 no limit",
			method: http.MethodGet,
			status: http.StatusOK,
			httpParams: &httpParams{
				topn:                "0",
				includeDistribution: "false",
			},
			gatewayGetRichlistResult: visor.Richlist{
				{
					Address: "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
					Coins:   "1000000.000000",
					Locked:  false,
				},
				{
					Address: "27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1",
					Coins:   "500000.000000",
					Locked:  false,
				},
				{
					Address: "2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW",
					Coins:   "500000.000000",
					Locked:  false,
				},
				{
					Address: "2TmvdBWJgxMwGs84R4drS9p5fYkva4dGdfs",
					Coins:   "244458.000000",
					Locked:  false,
				},
				{
					Address: "24gvUHXHtSg5drKiFsMw7iMgoN2PbLub53C",
					Coins:   "195503.000000",
					Locked:  false,
				},
			},
			result: Richlist{
				Richlist: visor.Richlist{
					{
						Address: "2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF",
						Coins:   "1000000.000000",
						Locked:  false,
					},
					{
						Address: "27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1",
						Coins:   "500000.000000",
						Locked:  false,
					},
					{
						Address: "2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW",
						Coins:   "500000.000000",
						Locked:  false,
					},
					{
						Address: "2TmvdBWJgxMwGs84R4drS9p5fYkva4dGdfs",
						Coins:   "244458.000000",
						Locked:  false,
					},
					{
						Address: "24gvUHXHtSg5drKiFsMw7iMgoN2PbLub53C",
						Coins:   "195503.000000",
						Locked:  false,
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/richlist"
			gateway := NewGatewayerMock()
			gateway.On("GetRichlist", tc.includeDistribution).Return(tc.gatewayGetRichlistResult, tc.gatewayGetRichlistErr)

			v := url.Values{}
			if tc.httpParams != nil {
				if tc.httpParams.topn != "" {
					v.Add("n", tc.httpParams.topn)
				}
				if tc.httpParams.includeDistribution != "" {
					v.Add("include-distribution", tc.httpParams.includeDistribution)
				}
			}
			if len(v) > 0 {
				endpoint += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Richlist
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestGetAddressCount(t *testing.T) {
	type Result struct {
		Count uint64
	}
	tt := []struct {
		name                         string
		method                       string
		status                       int
		err                          string
		gatewayGetAddressCountResult uint64
		gatewayGetAddressCountErr    error
		result                       Result
		csrfDisabled                 bool
	}{
		{
			name:   "405",
			method: http.MethodPost,
			status: http.StatusMethodNotAllowed,
			err:    "405 Method Not Allowed",
		},
		{
			name:   "500 - gw GetAddressCount error",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error",
			gatewayGetAddressCountErr: errors.New("gatewayGetAddressCountErr"),
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,
			gatewayGetAddressCountResult: 1,
			result: Result{
				Count: 1,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/addresscount"
			gateway := NewGatewayerMock()
			gateway.On("GetAddressCount").Return(tc.gatewayGetAddressCountResult, tc.gatewayGetAddressCountErr)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			csrfStore := &CSRFStore{
				Enabled: !tc.csrfDisabled,
			}
			if csrfStore.Enabled {
				setCSRFParameters(csrfStore, tokenValid, req)
			} else {
				setCSRFParameters(csrfStore, tokenInvalid, req)
			}
			handler := newServerMux(muxConfig{host: configuredHost, appLoc: "."}, gateway, csrfStore)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Result
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
