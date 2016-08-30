// Blockchain related information for the GUI
package gui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/daemon"
)

func blockchainHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendOr404(w, gateway.GetBlockchainMetadata())
	}
}

func blockchainBlockHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseq := r.FormValue("seq")
		seq, err := strconv.ParseUint(sseq, 10, 64)
		if err != nil {
			Error400(w, fmt.Sprintf("Invalid seq value \"%s\"", sseq))
			return
		}
		SendOr404(w, gateway.GetBlock(seq))
	}
}

func blockchainBlocksHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sstart := r.FormValue("start")
		start, err := strconv.ParseUint(sstart, 10, 64)
		if err != nil {
			Error400(w, fmt.Sprintf("Invalid start value \"%s\"", sstart))
			return
		}
		send := r.FormValue("end")
		end, err := strconv.ParseUint(send, 10, 64)
		if err != nil {
			Error400(w, fmt.Sprintf("Invalid end value \"%s\"", send))
			return
		}
		SendOr404(w, gateway.GetBlocks(start, end))
	}
}

func blockchainProgressHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		SendOr404(w, gateway.GetBlockchainProgress())
	}
}

func RegisterBlockchainHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	mux.HandleFunc("/blockchain", blockchainHandler(gateway))
	mux.HandleFunc("/blockchain/block", blockchainBlockHandler(gateway))
	mux.HandleFunc("/blockchain/blocks", blockchainBlocksHandler(gateway))
	mux.HandleFunc("/blockchain/progress", blockchainProgressHandler(gateway))
}
