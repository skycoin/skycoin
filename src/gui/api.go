// Package gui Api-related information for the GUI
package gui

import (
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	// "github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	// "github.com/skycoin/skycoin/src/visor"
	// "github.com/skycoin/skycoin/src/wallet"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// Wallet @todo remove duplicate struct of src/wallet/deterministic.go Wallet struct
// when that will be adopted to have many entries
type Wallet struct {
	Meta    map[string]string `json:"meta"`
	Entries []KeyEntry        `json:"entries"`
}

// KeyEntry wallet entry
type KeyEntry struct {
	Address string `json:"address"`
	Public  string `json:"public_key"`
	Secret  string `json:"secret_key"`
}

var (
	isBitcoin  bool
	hideSecKey bool
)

// Generating secret key, address, public key by given
// GET/POST
// 	bc - bool - is bitcoin type (optional) - default: true
//	n - int - Generation count (optional) - default: 1
//	s - bool - is hide secret key (optional) - default: false
//	seed - string - seed hash
func apiCreateAddressHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var seed = r.FormValue("seed")
		var err error

		if seed == "" {
			wh.Error400(w, "Empty seed")
			return
		}

		isBitcoin, err = strconv.ParseBool(r.FormValue("bc"))
		if err != nil {
			isBitcoin = true
		}

		genCount, err := strconv.Atoi(r.FormValue("n"))
		if err != nil {
			genCount = 1
		}

		hideSecKey, err = strconv.ParseBool(r.FormValue("s"))
		if err != nil {
			hideSecKey = false
		}

		wallet := Wallet{
			Meta:    make(map[string]string), //map[string]string
			Entries: make([]KeyEntry, genCount),
		}

		if isBitcoin == false {
			wallet.Meta = map[string]string{"coin": "skycoin"}
		} else {
			wallet.Meta = map[string]string{"coin": "bitcoin"}
		}

		wallet.Meta["seed"] = seed

		seckeys := cipher.GenerateDeterministicKeyPairs([]byte(seed), genCount)

		for i, sec := range seckeys {
			pub := cipher.PubKeyFromSecKey(sec)
			wallet.Entries[i] = getKeyEntry(pub, sec)
		}

		ret := wallet

		wh.SendOr404(w, ret)
	}
}

func getKeyEntry(pub cipher.PubKey, sec cipher.SecKey) KeyEntry {

	var e KeyEntry

	//skycoin address
	if isBitcoin == false {
		e = KeyEntry{
			Address: cipher.AddressFromPubKey(pub).String(),
			Public:  pub.Hex(),
			Secret:  sec.Hex(),
		}
	}

	//bitcoin address
	if isBitcoin == true {
		e = KeyEntry{
			Address: cipher.BitcoinAddressFromPubkey(pub),
			Public:  pub.Hex(),
			Secret:  cipher.BitcoinWalletImportFormatFromSeckey(sec),
		}
	}

	//hide the secret key
	if hideSecKey == true {
		e.Secret = ""
	}

	return e
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
