// Package gui Api-related information for the GUI
package gui

import (
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/daemon"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/wallet"
)

// Generating secret key, address, public key by given
// GET/POST
// 	bc - bool - is bitcoin type (optional) - default: true
//	n - int - Generation count (optional) - default: 1
//	s - bool - is hide secret key (optional) - default: false
//	seed - string - seed hash
func apiCreateAddressHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seed := r.FormValue("seed")
		if seed == "" {
			wh.Error400(w, "Empty seed")
			return
		}

		isBitcoin := true
		isBitcoinOpt := r.FormValue("bc")
		if isBitcoinOpt != "" {
			var err error
			if isBitcoin, err = strconv.ParseBool(isBitcoinOpt); err != nil {
				wh.Error400(w, `Invalid bool for "bc"`)
				return
			}
		}

		genCount := 1
		genCountOpt := r.FormValue("n")
		if genCountOpt != "" {
			var err error
			if genCount, err = strconv.Atoi(genCountOpt); err != nil {
				wh.Error400(w, `Invalid int for "n"`)
				return
			}
		}

		if genCount < 1 {
			wh.Error400(w, `"n" must be > 0`)
			return
		}

		// TODO -- hideSecKey should probably default to true
		hideSecKey := false
		hideSecKeyOpt := r.FormValue("s")
		if hideSecKeyOpt != "" {
			var err error
			if hideSecKey, err = strconv.ParseBool(hideSecKeyOpt); err != nil {
				wh.Error400(w, `Invalid bool for "s"`)
				return
			}
		}

		var coinType wallet.CoinType
		if isBitcoin {
			coinType = wallet.CoinTypeBitcoin
		} else {
			coinType = wallet.CoinTypeSkycoin
		}

		wallet, err := wallet.CreateAddresses(coinType, seed, genCount, hideSecKey)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendOr404(w, wallet)
	}
}

// RegisterAPIHandlers registers api handlers
func RegisterAPIHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	//  Generates wallet bitcoin/skycoin addresses and seckey,pubkey
	// GET/POST
	// 	bc - bool - is bitcoin type (optional) - default: true
	//	n - int - Generation count (optional) - default: 1
	//	s - bool - is hide secret key (optional) - default: false
	//	seed - string - seed hash
	mux.HandleFunc("/api/create-address", apiCreateAddressHandler(gateway))
}
