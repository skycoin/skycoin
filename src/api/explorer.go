package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/droplet"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// CoinSupply records the coin supply info
type CoinSupply struct {
	// Coins distributed beyond the project:
	CurrentSupply string `json:"current_supply"`
	// TotalSupply is CurrentSupply plus coins held by the distribution addresses that are spendable
	TotalSupply string `json:"total_supply"`
	// MaxSupply is the maximum number of coins to be distributed ever
	MaxSupply string `json:"max_supply"`
	// CurrentCoinHourSupply is coins hours in non distribution addresses
	CurrentCoinHourSupply string `json:"current_coinhour_supply"`
	// TotalCoinHourSupply is coin hours in all addresses including unlocked distribution addresses
	TotalCoinHourSupply string `json:"total_coinhour_supply"`
	// Distribution addresses which count towards total supply
	UnlockedAddresses []string `json:"unlocked_distribution_addresses"`
	// Distribution addresses which are locked and do not count towards total supply
	LockedAddresses []string `json:"locked_distribution_addresses"`
}

// newStringSet returns a map-based set for string lookup
func newStringSet(keys []string) map[string]struct{} {
	s := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		s[k] = struct{}{}
	}
	return s
}

// coinSupplyHandler returns coin distribution supply stats
// Method: GET
// URI: /api/v1/coinSupply
func coinSupplyHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		allUnspents, err := gateway.GetUnspentOutputsSummary(nil)
		if err != nil {
			err = fmt.Errorf("gateway.GetUnspentOutputsSummary failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		unlockedAddrs := params.GetUnlockedDistributionAddresses()
		// Search map of unlocked addresses
		// used to filter unspents
		unlockedAddrSet := newStringSet(unlockedAddrs)

		var unlockedSupply uint64
		// check confirmed unspents only
		for _, u := range allUnspents.Confirmed {
			// check if address is an unlocked distribution address
			if _, ok := unlockedAddrSet[u.Body.Address.String()]; ok {
				var err error
				unlockedSupply, err = coin.AddUint64(unlockedSupply, u.Body.Coins)
				if err != nil {
					err = fmt.Errorf("uint64 overflow while adding up unlocked supply coins: %v", err)
					wh.Error500(w, err.Error())
					return
				}
			}
		}

		// "total supply" is the number of coins unlocked.
		// Each distribution address was allocated params.DistributionAddressInitialBalance coins.
		totalSupply := uint64(len(unlockedAddrs)) * params.DistributionAddressInitialBalance
		totalSupply *= droplet.Multiplier

		// "current supply" is the number of coins distributed from the unlocked pool
		currentSupply := totalSupply - unlockedSupply

		currentSupplyStr, err := droplet.ToString(currentSupply)
		if err != nil {
			err = fmt.Errorf("Failed to convert coins to string: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		totalSupplyStr, err := droplet.ToString(totalSupply)
		if err != nil {
			err = fmt.Errorf("Failed to convert coins to string: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		maxSupplyStr, err := droplet.ToString(params.MaxCoinSupply * droplet.Multiplier)
		if err != nil {
			err = fmt.Errorf("Failed to convert coins to string: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		// locked distribution addresses
		lockedAddrs := params.GetLockedDistributionAddresses()
		lockedAddrSet := newStringSet(lockedAddrs)

		// get total coins hours which excludes locked distribution addresses
		var totalCoinHours uint64
		for _, out := range allUnspents.Confirmed {
			if _, ok := lockedAddrSet[out.Body.Address.String()]; !ok {
				var err error
				totalCoinHours, err = coin.AddUint64(totalCoinHours, out.CalculatedHours)
				if err != nil {
					err = fmt.Errorf("uint64 overflow while adding up total coin hours: %v", err)
					wh.Error500(w, err.Error())
					return
				}
			}
		}

		// get current coin hours which excludes all distribution addresses
		var currentCoinHours uint64
		for _, out := range allUnspents.Confirmed {
			// check if address not in locked distribution addresses
			if _, ok := lockedAddrSet[out.Body.Address.String()]; !ok {
				// check if address not in unlocked distribution addresses
				if _, ok := unlockedAddrSet[out.Body.Address.String()]; !ok {
					currentCoinHours += out.CalculatedHours
				}
			}
		}

		if err != nil {
			err = fmt.Errorf("Failed to get total coinhours: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		cs := CoinSupply{
			CurrentSupply:         currentSupplyStr,
			TotalSupply:           totalSupplyStr,
			MaxSupply:             maxSupplyStr,
			CurrentCoinHourSupply: strconv.FormatUint(currentCoinHours, 10),
			TotalCoinHourSupply:   strconv.FormatUint(totalCoinHours, 10),
			UnlockedAddresses:     unlockedAddrs,
			LockedAddresses:       params.GetLockedDistributionAddresses(),
		}

		wh.SendJSONOr500(logger, w, cs)
	}
}

// transactionsForAddressHandler returns all transactions (confirmed and unconfirmed) for an address
// Method: GET
// URI: /explorer/address
// Args:
//	address [string]
func transactionsForAddressHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Critical().Warning("Call to deprecated /api/v1/explorer/address endpoint")

		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr := r.FormValue("address")
		if addr == "" {
			wh.Error400(w, "address is empty")
			return
		}

		cipherAddr, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			wh.Error400(w, "invalid address")
			return
		}

		txns, inputs, err := gateway.GetVerboseTransactionsForAddress(cipherAddr)
		if err != nil {
			err = fmt.Errorf("gateway.GetVerboseTransactionsForAddress failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		vb := make([]readable.TransactionVerbose, len(txns))
		for i, txn := range txns {
			v, err := readable.NewTransactionVerbose(txn, inputs[i])
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			vb[i] = v
		}

		wh.SendJSONOr500(logger, w, vb)
	}
}

// Richlist contains top address balances
type Richlist struct {
	Richlist []readable.RichlistBalance `json:"richlist"`
}

// richlistHandler returns the top skycoin holders
// Method: GET
// URI: /richlist?n=${number}&include-distribution=${bool}
// Args:
//	n [int, number of results to include]
//  include-distribution [bool, include the distribution addresses in the richlist]
func richlistHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		var topn int
		topnStr := r.FormValue("n")
		if topnStr == "" {
			topn = 20
		} else {
			var err error
			topn, err = strconv.Atoi(topnStr)
			if err != nil {
				wh.Error400(w, "invalid n")
				return
			}
		}

		var includeDistribution bool
		includeDistributionStr := r.FormValue("include-distribution")
		if includeDistributionStr == "" {
			includeDistribution = false
		} else {
			var err error
			includeDistribution, err = strconv.ParseBool(includeDistributionStr)
			if err != nil {
				wh.Error400(w, "invalid include-distribution")
				return
			}
		}

		richlist, err := gateway.GetRichlist(includeDistribution)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		if topn > 0 && topn < len(richlist) {
			richlist = richlist[:topn]
		}

		wh.SendJSONOr500(logger, w, Richlist{
			Richlist: readable.NewRichlistBalances(richlist),
		})
	}
}

// addressCountHandler returns the total number of unique address that have coins
// Method: GET
// URI: /addresscount
func addressCountHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addrCount, err := gateway.GetAddressCount()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, &map[string]uint64{"count": addrCount})
	}
}
