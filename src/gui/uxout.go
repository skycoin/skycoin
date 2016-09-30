package gui

import (
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

type UxOutJson struct {
	Time          uint64 `json:"time"`
	SrcBkSeq      uint64 `json:"src_block_seq"`
	SrcTx         string `json:"src_tx"`
	OwnerAddress  string `json:"owner_address"`
	Coins         uint64 `json:"coins"`
	Hours         uint64 `json:"hours"`
	SpentBlockSeq uint64 `json:"spent_block_seq"` // block seq that spent the output.
	SpentTxID     string `json:"spent_tx"`        // id of tx which spent this output.
}

func RegisterUxOutHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get uxout by id.
	mux.HandleFunc("/uxout", getUxOutByID(gateway))
}

func getUxOutByID(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rlt struct {
			Success bool                 `json:"success"`
			Reason  string               `json:"reason, omitempty"`
			Uxout   visor.ReadableOutput `json:"output, omitempty"`
		}

		for {
			uxid := r.FormValue("uxid")
			if uxid == "" {
				rlt.Reason = "uxid is empty"
				break
			}
			id, err := cipher.SHA256FromHex(uxid)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}
			uxout, err := gateway.V.GetUxOutByID(id)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}
			rlt.Success = true
			rlt.Uxout = visor.NewReadableOutput(uxout.Out)
			break
		}
		wh.SendOr404(w, &rlt)
	}
}
