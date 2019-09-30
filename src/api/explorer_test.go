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

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/readable"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/SkycoinProject/skycoin/src/visor"
)

func makeSuccessCoinSupplyResult(t *testing.T, allUnspents readable.UnspentOutputsSummary) *CoinSupply {
	unlockedAddrs := params.MainNetDistribution.UnlockedAddressesDecoded()
	var unlockedSupply uint64
	// check confirmed unspents only
	// Search map of unlocked addresses
	// used to filter unspents
	unlockedAddrSet := newAddrSet(unlockedAddrs)
	for _, u := range allUnspents.HeadOutputs {
		// check if address is an unlocked distribution address
		if _, ok := unlockedAddrSet[cipher.MustDecodeBase58Address(u.Address)]; ok {
			coins, err := droplet.FromString(u.Coins)
			require.NoError(t, err)
			unlockedSupply += coins
		}
	}
	// "total supply" is the number of coins unlocked.
	// Each distribution address was allocated params.MainNetDistribution.AddressInitialBalance coins.
	totalSupply := uint64(len(unlockedAddrs)) * params.MainNetDistribution.AddressInitialBalance()
	totalSupply *= droplet.Multiplier

	// "current supply" is the number of coins distribution from the unlocked pool
	currentSupply := totalSupply - unlockedSupply

	currentSupplyStr, err := droplet.ToString(currentSupply)
	require.NoError(t, err)

	totalSupplyStr, err := droplet.ToString(totalSupply)
	require.NoError(t, err)

	maxSupplyStr, err := droplet.ToString(params.MainNetDistribution.MaxCoinSupply * droplet.Multiplier)
	require.NoError(t, err)

	// locked distribution addresses
	lockedAddrs := params.MainNetDistribution.LockedAddressesDecoded()
	lockedAddrSet := newAddrSet(lockedAddrs)

	// get total coins hours which excludes locked distribution addresses
	var totalCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		if _, ok := lockedAddrSet[cipher.MustDecodeBase58Address(out.Address)]; !ok {
			totalCoinHours += out.Hours
		}
	}

	// get current coin hours which excludes all distribution addresses
	var currentCoinHours uint64
	for _, out := range allUnspents.HeadOutputs {
		// check if address not in locked distribution addresses
		if _, ok := lockedAddrSet[cipher.MustDecodeBase58Address(out.Address)]; !ok {
			// check if address not in unlocked distribution addresses
			if _, ok := unlockedAddrSet[cipher.MustDecodeBase58Address(out.Address)]; !ok {
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
		UnlockedAddresses:     params.MainNetDistribution.UnlockedAddresses(),
		LockedAddresses:       params.MainNetDistribution.LockedAddresses(),
	}
	return &cs
}

func TestCoinSupply(t *testing.T) {
	addrs := []cipher.Address{
		testutil.MakeAddress(),
		testutil.MakeAddress(),
	}

	unlockedAddrs := params.MainNetDistribution.UnlockedAddressesDecoded()
	successGatewayGetUnspentOutputsResult := readable.UnspentOutputsSummary{
		HeadOutputs: readable.UnspentOutputs{
			readable.UnspentOutput{
				Address: addrs[0].String(),
				Coins:   "0",
			},
			readable.UnspentOutput{
				Address: addrs[1].String(),
				Coins:   "0",
			},
		},
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
			name:                        "500 - gatewayGetUnspentOutputsErr",
			method:                      http.MethodGet,
			status:                      http.StatusInternalServerError,
			err:                         "500 Internal Server Error - gateway.GetUnspentOutputsSummary failed: gatewayGetUnspentOutputsErr",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:                        "500 - gatewayGetUnspentOutputsErr",
			method:                      http.MethodGet,
			status:                      http.StatusInternalServerError,
			err:                         "500 Internal Server Error - gateway.GetUnspentOutputsSummary failed: gatewayGetUnspentOutputsErr",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsErr: errors.New("gatewayGetUnspentOutputsErr"),
		},
		{
			name:                        "500 - too large HeadOutputs item",
			method:                      http.MethodGet,
			status:                      http.StatusInternalServerError,
			err:                         "500 Internal Server Error - Failed to convert coins to string: Droplet string conversion failed: Value is too large",
			gatewayGetUnspentOutputsArg: filterInUnlocked,
			gatewayGetUnspentOutputsResult: &visor.UnspentOutputsSummary{
				Confirmed: []visor.UnspentOutput{
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   9223372036854775807,
								Address: unlockedAddrs[0],
							},
						},
					},
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   1000000,
								Address: unlockedAddrs[0],
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
								Coins:   0,
								Address: addrs[0],
							},
						},
					},
					visor.UnspentOutput{
						UxOut: coin.UxOut{
							Body: coin.UxBody{
								Coins:   0,
								Address: addrs[1],
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
			gateway.On("VisorConfig").Return(visor.Config{
				Distribution: params.MainNetDistribution,
			})

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			cfg := defaultMuxConfig()
			cfg.disableCSRF = tc.csrfDisabled

			handler := newServerMux(cfg, gateway)
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
					Address: cipher.MustDecodeBase58Address("2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF"),
					Coins:   1000000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1"),
					Coins:   500000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW"),
					Coins:   500000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("2TmvdBWJgxMwGs84R4drS9p5fYkva4dGdfs"),
					Coins:   244458e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("24gvUHXHtSg5drKiFsMw7iMgoN2PbLub53C"),
					Coins:   195503e6,
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
					Address: cipher.MustDecodeBase58Address("2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF"),
					Coins:   1000000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("27jg25DZX21MXMypVbKJMmgCJ5SPuEunMF1"),
					Coins:   500000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("2fGi2jhvp6ppHg3DecguZgzqvpJj2Gd4KHW"),
					Coins:   500000e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("2TmvdBWJgxMwGs84R4drS9p5fYkva4dGdfs"),
					Coins:   244458e6,
					Locked:  false,
				},
				{
					Address: cipher.MustDecodeBase58Address("24gvUHXHtSg5drKiFsMw7iMgoN2PbLub53C"),
					Coins:   195503e6,
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

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}
			handler := newServerMux(muxConfig{
				host:           configuredHost,
				appLoc:         ".",
				disableCSRF:    tc.csrfDisabled,
				disableCSP:     true,
				enabledAPISets: allAPISetsEnabled,
			}, gateway)
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
			name:                      "500 - gw GetAddressCount error",
			method:                    http.MethodGet,
			status:                    http.StatusInternalServerError,
			err:                       "500 Internal Server Error - gatewayGetAddressCountErr",
			gatewayGetAddressCountErr: errors.New("gatewayGetAddressCountErr"),
		},
		{
			name:                         "200",
			method:                       http.MethodGet,
			status:                       http.StatusOK,
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
			gateway.On("AddressCount").Return(tc.gatewayGetAddressCountResult, tc.gatewayGetAddressCountErr)

			req, err := http.NewRequest(tc.method, endpoint, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()

			if tc.csrfDisabled {
				setCSRFParameters(t, tokenInvalid, req)
			} else {
				setCSRFParameters(t, tokenValid, req)
			}

			handler := newServerMux(muxConfig{
				host:           configuredHost,
				appLoc:         ".",
				disableCSRF:    tc.csrfDisabled,
				disableCSP:     true,
				enabledAPISets: allAPISetsEnabled,
			}, gateway)
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
