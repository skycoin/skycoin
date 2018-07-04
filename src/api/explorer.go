package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor"
)

func getCoinSupply(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		supply, err := gateway.GetCoinSupply()
		if err != nil {
			err = fmt.Errorf("gateway.GetCoinSupply failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}
		if supply != nil {
			wh.SendJSONOr500(logger, w, supply)
		}
	}
}

// method: GET
// url: /explorer/address?address=${address}
func getTransactionsForAddress(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		txns, err := gateway.GetTransactionsForAddress(cipherAddr)
		if err != nil {
			err = fmt.Errorf("gateway.GetTransactionsForAddress failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txns)
	}
}

// Richlist is the API response for /richlist, contains top address balances
type Richlist struct {
	Richlist visor.Richlist `json:"richlist"`
}

// method: GET
// url: /richlist?n=${number}&include-distribution=${bool}
func getRichlist(gateway Gatewayer) http.HandlerFunc {
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
			Richlist: richlist,
		})
	}
}

// method: GET
// url: /addresscount
func getAddressCount(gateway Gatewayer) http.HandlerFunc {
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
