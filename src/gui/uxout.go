package gui

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor/historydb"
)

func RegisterUxOutHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get uxout by id.
	mux.HandleFunc("/uxout", getUxOutByID(gateway))
	// get address in uxouts
	mux.HandleFunc("/address_in_uxouts", getRecvUxOutOfAddr(gateway))
	// get address out uxouts
	mux.HandleFunc("/address_out_uxouts", getSpentOutUxOutOfAddr(gateway))
}

func getUxOutByID(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		uxout, err := gateway.V.GetUxOutByID(id)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if uxout == nil {
			wh.Error404(w, "not found")
			return
		}

		wh.SendOr404(w, historydb.NewUxOutJSON(uxout))
	}
}

func getRecvUxOutOfAddr(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		uxs, err := gateway.GetRecvUxOutOfAddr(cipherAddr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxOuts := make([]*historydb.UxOutJSON, len(uxs))
		for i, ux := range uxs {
			uxOuts[i] = historydb.NewUxOutJSON(ux)
		}
		wh.SendOr404(w, &uxOuts)
	}
}

func getSpentOutUxOutOfAddr(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		uxs, err := gateway.GetSpentUxOutOfAddr(cipherAddr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxOuts := make([]*historydb.UxOutJSON, len(uxs))
		for i, ux := range uxs {
			uxOuts[i] = newUxOutJson(ux)
		}
		wh.SendOr404(w, &uxOuts)
	}
}
