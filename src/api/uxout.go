package api

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
)

// URI: /api/v1/uxout
// Method: GET
// Args:
//	uxid: unspent output ID hash
// Returns an unspent output by ID
func uxOutHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		uxid := r.FormValue("uxid")
		if uxid == "" {
			wh.Error400(w, "uxid is empty")
			return
		}

		id, err := cipher.SHA256FromHex(uxid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxout, headTime, err := gateway.GetUxOutByID(id)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if uxout == nil {
			wh.Error404(w, "")
			return
		}

		out, err := readable.NewSpentOutput(uxout, headTime)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, *out)
	}
}

// URI: /api/v1/address_uxouts
// Method: GET
// Args:
//	address
// Returns the historical, spent outputs associated with an address
func addrUxOutsHandler(gateway Gatewayer) http.HandlerFunc {
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
			wh.Error400(w, err.Error())
			return
		}

		uxs, headTime, err := gateway.GetSpentOutputsForAddresses([]cipher.Address{cipherAddr})
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		ret := make([]readable.SpentOutput, 0)
		for _, u := range uxs {
			out, err := readable.NewSpentOutputs(u, headTime)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}
			ret = append(ret, out...)
		}

		wh.SendJSONOr500(logger, w, ret)
	}
}
