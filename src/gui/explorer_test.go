package gui

import (
	"net/http"
	"testing"

	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/pkg/errors"
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

// GetAddressTxns returns a *visor.TransactionResults
func (gw *FakeGateway) GetAddressTxns(a cipher.Address) (*visor.TransactionResults, error) {
	args := gw.Called(a)
	return args.Get(0).(*visor.TransactionResults), args.Error(1)
}

// GetUxOutByID gets UxOut by hash id.
func (gw *FakeGateway) GetUxOutByID(id cipher.SHA256) (*historydb.UxOut, error) {
	args := gw.Called(id)
	return args.Get(0).(*historydb.UxOut), args.Error(1)

}

// GetUnspentOutputs gets unspent outputs and returns the filtered results,
// Note: all filters will be executed as the pending sequence in 'AND' mode.
func (gw *FakeGateway) GetUnspentOutputs(filters ...daemon.OutputsFilter) (visor.ReadableOutputSet, error) {
	args := gw.Called(filters)
	return args.Get(0).(visor.ReadableOutputSet), args.Error(1)
}

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
	invalidHash := "caicb"
	validHash := "79216473e8f2c17095c6887cc9edca6c023afedfac2e0c5460e8b6f359684f8b"
	tt := []struct {
		name                        string
		method                      string
		url                         string
		status                      int
		err                         string
		addressParam                string
		gatewayGetAddressTxnsResult *visor.TransactionResults
		gatewayGetAddressTxnsErr    error
		gatewayGetUxOutByIDArg      cipher.SHA256
		gatewayGetUxOutByIDResult   *historydb.UxOut
		gatewayGetUxOutByIDErr      error
		result                      []ReadableTransaction
	}{
		{
			"405",
			http.MethodPost,
			"/explorer/address",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			"0",
			&visor.TransactionResults{},
			nil,
			cipher.SHA256{},
			nil,
			nil,
			nil,
		},
		{
			"400 - address is empty",
			http.MethodGet,
			"/explorer/address",
			http.StatusBadRequest,
			"400 Bad Request - address is empty",
			"",
			&visor.TransactionResults{},
			nil,
			cipher.SHA256{},
			nil,
			nil,
			nil,
		},
		{
			"400 - invalid address",
			http.MethodGet,
			"/explorer/address",
			http.StatusBadRequest,
			"400 Bad Request - invalid address",
			"badAddress",
			&visor.TransactionResults{},
			nil,
			cipher.SHA256{},
			nil,
			nil,
			nil,
		},
		{
			"500 - gw GetAddressTxns error",
			http.MethodGet,
			"/explorer/address",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			address.String(),
			&visor.TransactionResults{},
			errors.New("gatewayGetAddressTxnsErr"),
			cipher.SHA256{},
			nil,
			nil,
			nil,
		},
		{
			"500 - cipher.SHA256FromHex(tx.Transaction.In) error",
			http.MethodGet,
			"/explorer/address",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			address.String(),
			&visor.TransactionResults{
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
			nil,
			cipher.SHA256{},
			nil,
			nil,
			nil,
		},
		{
			"500 - GetUxOutByID error",
			http.MethodGet,
			"/explorer/address",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			address.String(),
			&visor.TransactionResults{
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
			nil,
			testutil.SHA256FromHex(t, validHash),
			nil,
			errors.New("gatewayGetUxOutByIDErr"),
			nil,
		},
		{
			"500 - GetUxOutByID nil result",
			http.MethodGet,
			"/explorer/address",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			address.String(),
			&visor.TransactionResults{
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
			nil,
			testutil.SHA256FromHex(t, validHash),
			nil,
			nil,
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/explorer/address",
			http.StatusOK,
			"",
			address.String(),
			&visor.TransactionResults{
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
			nil,
			testutil.SHA256FromHex(t, validHash),
			&historydb.UxOut{},
			nil,
			[]ReadableTransaction{
				{
					In: []visor.ReadableTransactionInput{
						{
							Hash:    validHash,
							Address: successAddress,
						},
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetAddressTxns", address).Return(tc.gatewayGetAddressTxnsResult, tc.gatewayGetAddressTxnsErr)
			gateway.On("GetUxOutByID", tc.gatewayGetUxOutByIDArg).Return(tc.gatewayGetUxOutByIDResult, tc.gatewayGetUxOutByIDErr)

			v := url.Values{}
			var urlFull = tc.url
			if tc.addressParam != "" {
				v.Add("address", tc.addressParam)
			}

			if len(v) > 0 {
				urlFull += "?" + v.Encode()
			}

			req, err := http.NewRequest(tc.method, urlFull, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getTransactionsForAddress(gateway))

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg []ReadableTransaction
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
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
	var filterInUnlocked []daemon.OutputsFilter = nil
	filterInUnlocked = append(filterInUnlocked, daemon.FbyAddresses(unlockedAddrs))
	tt := []struct {
		name                           string
		method                         string
		url                            string
		status                         int
		err                            string
		gatewayGetUnspentOutputsArg    []daemon.OutputsFilter
		gatewayGetUnspentOutputsResult visor.ReadableOutputSet
		gatewayGetUnspentOutputsErr    error
		result                         *CoinSupply
	}{
		{
			"405",
			http.MethodPost,
			"/coinSupply",
			http.StatusMethodNotAllowed,
			"405 Method Not Allowed",
			[]daemon.OutputsFilter{},
			visor.ReadableOutputSet{},
			nil,
			nil,
		},
		{
			"500 - gatewayGetUnspentOutputsErr",
			http.MethodGet,
			"/coinSupply",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			filterInUnlocked,
			visor.ReadableOutputSet{},
			errors.New("gatewayGetUnspentOutputsErr"),
			nil,
		},
		{
			"500 - gatewayGetUnspentOutputsErr",
			http.MethodGet,
			"/coinSupply",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			filterInUnlocked,
			visor.ReadableOutputSet{},
			errors.New("gatewayGetUnspentOutputsErr"),
			nil,
		},
		{
			"500 - too large HeadOutputs item",
			http.MethodGet,
			"/coinSupply",
			http.StatusInternalServerError,
			"500 Internal Server Error",
			filterInUnlocked,
			visor.ReadableOutputSet{
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
			nil,
			nil,
		},
		{
			"200",
			http.MethodGet,
			"/coinSupply",
			http.StatusOK,
			"",
			filterInUnlocked,
			visor.ReadableOutputSet{
				HeadOutputs: visor.ReadableOutputs{
					visor.ReadableOutput{
						Coins: "0",
						//Coins: "9223372036854",
					},
					visor.ReadableOutput{
						Coins: "0",
					},
				},
			},
			nil,
			makeSuccessCoinSupplyResult(t, successGatewayGetUnspentOutputsResult),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &FakeGateway{
				t: t,
			}
			gateway.On("GetUnspentOutputs", mock.Anything).Return(tc.gatewayGetUnspentOutputsResult, tc.gatewayGetUnspentOutputsErr)

			req, err := http.NewRequest(tc.method, tc.url, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getCoinSupply(gateway))

			handler.ServeHTTP(rr, req)

			status := rr.Code
			require.Equal(t, tc.status, status, "case: %s, handler returned wrong status code: got `%v` want `%v`", tc.name, status, tc.status)

			if status != http.StatusOK {
				require.Equal(t, tc.err, strings.TrimSpace(rr.Body.String()), "case: %s, handler returned wrong error message: got `%v`| %s, want `%v`",
					tc.name, strings.TrimSpace(rr.Body.String()), status, tc.err)
			} else {
				var msg *CoinSupply
				err := json.Unmarshal(rr.Body.Bytes(), &msg)
				require.NoError(t, err)
				require.Equal(t, tc.result, msg)
			}
		})
	}
}
