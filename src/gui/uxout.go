package gui

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"github.com/skycoin/skycoin/src/visor/historydb"
)

type uxOutJSON struct {
	Time          uint64 `json:"time"`
	SrcBkSeq      uint64 `json:"src_block_seq"`
	SrcTx         string `json:"src_tx"`
	OwnerAddress  string `json:"owner_address"`
	Coins         uint64 `json:"coins"`
	Hours         uint64 `json:"hours"`
	SpentBlockSeq uint64 `json:"spent_block_seq"` // block seq that spent the output.
	SpentTxID     string `json:"spent_tx"`        // id of tx which spent this output.
}

func newUxOutJson(out *historydb.UxOut) *uxOutJSON {
	return &uxOutJSON{
		Time:          out.Out.Head.Time,
		SrcBkSeq:      out.Out.Head.BkSeq,
		SrcTx:         out.Out.Body.SrcTransaction.Hex(),
		OwnerAddress:  out.Out.Body.Address.String(),
		Coins:         out.Out.Body.Coins,
		Hours:         out.Out.Body.Hours,
		SpentBlockSeq: out.SpentBlockSeq,
		SpentTxID:     out.SpentTxID.Hex(),
	}
}

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

		wh.SendOr404(w, newUxOutJson(uxout))
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

		uxs, err := gateway.V.GetRecvUxOutOfAddr(cipherAddr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxOuts := make([]*uxOutJSON, len(uxs))
		for i, ux := range uxs {
			uxOuts[i] = newUxOutJson(ux)
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

		uxs, err := gateway.V.GetSpentUxOutOfAddr(cipherAddr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxOuts := make([]*uxOutJSON, len(uxs))
		for i, ux := range uxs {
			uxOuts[i] = newUxOutJson(ux)
		}
		wh.SendOr404(w, &uxOuts)
	}
}
