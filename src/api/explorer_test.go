package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
)

func makeSuccessCoinSupplyResult(t *testing.T, allUnspents readable.UnspentOutputsSummary) *CoinSupply {
	unlockedAddrs := visor.GetUnlockedDistributionAddresses()
	var unlockedSupply uint64
	// check confirmed unspents only
	// Search map of unlocked addresses
	// used to filter unspents
	unlockedAddrSet := newStringSet(unlockedAddrs)
	for _, u := range allUnspents.HeadOutputs {
		// check if address is an unlocked distribution address
		if _, ok := unlockedAddrSet[u.Address]; ok {
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
	lockedAddrSet := newStringSet(lockedAddrs)

	// get total coins hours which excludes locked distribution addresses
	var totalCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		if _, ok := lockedAddrSet[out.Address]; !ok {
			totalCoinHours += out.Hours
		}
	}

	// get current coin hours which excludes all distribution addresses
	var currentCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		// check if address not in locked distribution addresses
		if _, ok := lockedAddrSet[out.Address]; !ok {
			// check if address not in unlocked distribution addresses
			if _, ok := unlockedAddrSet[out.Address]; !ok {
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
	successAddressRaw, err := cipher.DecodeBase58Address(successAddress)
	require.NoError(t, err)

	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	validHashRaw, err := cipher.SHA256FromHex(validHash)
	require.NoError(t, err)

	type verboseResult struct {
		Transactions []visor.Transaction
		Inputs       [][]visor.TransactionInput
	}

	tt := []struct {
		name                                   string
		method                                 string
		status                                 int
		err                                    string
		addressParam                           string
		gatewayGetTransactionsForAddressErr    error
		gatewayGetTransactionsForAddressResult verboseResult
		result                                 []readable.TransactionVerbose
		csrfDisabled                           bool
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
			name:                                "500 - gw GetVerboseTransactionsForAddress error",
			method:                              http.MethodGet,
			status:                              http.StatusInternalServerError,
			err:                                 "500 Internal Server Error - gateway.GetVerboseTransactionsForAddress failed: gatewayGetTransactionsForAddressErr",
			addressParam:                        address.String(),
			gatewayGetTransactionsForAddressErr: errors.New("gatewayGetTransactionsForAddressErr"),
		},
		{
			name:         "200",
			method:       http.MethodGet,
			status:       http.StatusOK,
			addressParam: address.String(),
			gatewayGetTransactionsForAddressResult: verboseResult{
				Transactions: []visor.Transaction{
					{
						Transaction: coin.Transaction{
							In: []cipher.SHA256{
								validHashRaw,
							},
						},
					},
				},
				Inputs: [][]visor.TransactionInput{
					[]visor.TransactionInput{
						{
							UxOut: coin.UxOut{
								Body: coin.UxBody{
									Address: successAddressRaw,
									Coins:   99000000,
									Hours:   100,
								},
							},
							CalculatedHours: 101,
						},
					},
				},
			},
			result: []readable.TransactionVerbose{
				{
					Status: &readable.TransactionStatus{
						Unconfirmed: true,
					},
					BlockTransactionVerbose: readable.BlockTransactionVerbose{
						Hash:      "4fa025f043d1e5e8895ca4dc6602dac8d5c315544c166044d80c98a09e950c71",
						InnerHash: "0000000000000000000000000000000000000000000000000000000000000000",
						Fee:       101,
						Sigs:      []string{},
						In: []readable.TransactionInput{
							{
								Hash:            "e8ca653d9953b548f0098dd303f8166e636856a5c40e478e3756e440c01e9cb9",
								Address:         successAddress,
								Coins:           "99.000000",
								Hours:           100,
								CalculatedHours: 101,
							},
						},
						Out: []readable.TransactionOutput{},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/explorer/address"
			gateway := &MockGatewayer{}
			gateway.On("GetVerboseTransactionsForAddress", address).Return(tc.gatewayGetTransactionsForAddressResult.Transactions,
				tc.gatewayGetTransactionsForAddressResult.Inputs, tc.gatewayGetTransactionsForAddressErr)

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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []readable.TransactionVerbose
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}

func TestCoinSupply(t *testing.T) {
	unlockedAddrs := visor.GetUnlockedDistributionAddresses()
	successGatewayGetUnspentOutputsResult := readable.UnspentOutputsSummary{
		HeadOutputs: readable.UnspentOutputs{
			readable.UnspentOutput{
				Coins: "0",
			},
			readable.UnspentOutput{
				Coins: "0",
			},
		},
	}

	unlockedAddrsRaw := make([]cipher.Address, len(unlockedAddrs))
	for i, addr := range unlockedAddrs {
		a, err := cipher.DecodeBase58Address(addr)
		require.NoError(t, err)
		unlockedAddrsRaw[i] = a
	}

	var filterInUnlocked []visor.OutputsFilter
	filterInUnlocked = append(filterInUnlocked, visor.FbyAddresses(unlockedAddrs))
	tt := []struct {
		name                           string
		method                         string
		status                         int
		err                            string
		gatewayGetUnspentOutputsArg    []visor.OutputsFilter
		gatewayGetUnspentOutputsResult *visor.UnspentOutputsSummary
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
			err:    "500 Internal Server Error - gateway.GetUnspentOutputsSummary failed: gatewayGetUnspentOutputsErr",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:   "500 - gatewayGetUnspentOutputsErr",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - gateway.GetUnspentOutputsSummary failed: gatewayGetUnspentOutputsErr",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:   "500 - too large HeadOutputs item",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			err:    "500 Internal Server Error - Failed to convert coins to string: Droplet string conversion failed: Value is too large",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsResult: &visor.UnspentOutputsSummary{
				Confirmed: []visor.UnspentOutput{
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   9223372036854775807,
								Address: unlockedAddrsRaw[0],
							},
						},
					},
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   1000000,
								Address: unlockedAddrsRaw[0],
							},
						},
					},
				},
			},
		},
		{
			name:   "200",
			method: http.MethodGet,
			status: http.StatusOK,

			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsResult: &visor.UnspentOutputsSummary{
				Confirmed: []visor.UnspentOutput{
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins: 0,
							},
						},
					},
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins: 0,
							},
						},
					},
				},
			},
			result: makeSuccessCoinSupplyResult(t, successGatewayGetUnspentOutputsResult),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			endpoint := "/api/v1/coinSupply"
			gateway := &MockGatewayer{}
			gateway.On("GetUnspentOutputsSummary", mock.Anything).Return(tc.gatewayGetUnspentOutputsResult, tc.gatewayGetUnspentOutputsErr)

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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`", strings.TrimSpace(rr.Body.String()), status, tc.err)
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
			err:    "500 Internal Server Error - gatewayGetRichlistErr",
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
				Richlist: []readable.RichlistBalance{
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
				Richlist: []readable.RichlistBalance{
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
			endpoint := "/api/v1/richlist"
			gateway := &MockGatewayer{}
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
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
			err:    "500 Internal Server Error - gatewayGetAddressCountErr",
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
			endpoint := "/api/v1/addresscount"
			gateway := &MockGatewayer{}
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
			handler := newServerMux(defaultMuxConfig(), gateway, csrfStore, nil)
			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "got `%v` want `%v`", status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "got `%v`| %d, want `%v`",
					strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg Result
				err = json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
