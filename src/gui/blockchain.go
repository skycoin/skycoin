// Blockchain related information for the GUI
package gui

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor" //http,json helpers

	"github.com/skycoin/skycoin/src/daemon"
)

func blockchainHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetBlockchainMetadata())
	}
}

func blockchainBlockHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sseq := r.FormValue("seq")
		seq, err := strconv.ParseUint(sseq, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid seq value \"%s\"", sseq))
			return
		}
		wh.SendOr404(w, gateway.GetBlock(seq))
	}
}

func blockchainBlocksHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sstart := r.FormValue("start")
		start, err := strconv.ParseUint(sstart, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid start value \"%s\"", sstart))
			return
		}
		send := r.FormValue("end")
		end, err := strconv.ParseUint(send, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid end value \"%s\"", send))
			return
		}
		wh.SendOr404(w, gateway.GetBlocks(start, end))
	}
}

func blockchainProgressHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetBlockchainProgress())
	}
}

func RegisterBlockchainHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	mux.HandleFunc("/blockchain", blockchainHandler(gateway))
	// get block by hash
	mux.HandleFunc("/blockchain/block/hash", getBlockByHash(gateway))
	// get block by seq
	mux.HandleFunc("/blockchain/block/seq", getBlockBySeq(gateway))
	mux.HandleFunc("/blockchain/blocks", blockchainBlocksHandler(gateway))
	mux.HandleFunc("/blockchain/progress", blockchainProgressHandler(gateway))
}

func getBlockByHash(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.FormValue("hash")
		var rlt struct {
			Success bool                `json:"success"`
			Reason  string              `json:"reason,omitempty"`
			Block   visor.ReadableBlock `json:"block,omitempty"`
		}
		for {
			if hash == "" {
				rlt.Reason = "block hash is empty"
				break
			}

			h, err := cipher.SHA256FromHex(hash)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}

			b := gate.V.GetBlockByHash(h)
			if b == nil {
				rlt.Reason = fmt.Sprintf("block of hash:%s does not exist", hash)
				break
			}

			rlt.Success = true
			rlt.Block = visor.NewReadableBlock(b)
			break
		}
		wh.SendOr404(w, &rlt)
	}
}

func getBlockBySeq(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seq := r.FormValue("seq")
		var rlt struct {
			Success bool                `json:"success"`
			Reason  string              `json:"reason,omitempty"`
			Block   visor.ReadableBlock `json:"block,omitempty"`
		}
		for {
			if seq == "" {
				rlt.Reason = "block seq is empty"
				break
			}

			uSeq, err := strconv.ParseUint(seq, 10, 64)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}

			b := gate.V.GetBlockBySeq(uSeq)
			if b == nil {
				rlt.Reason = fmt.Sprintf("block in seq:%s does not exist", seq)
				break
			}

			rlt.Success = true
			rlt.Block = visor.NewReadableBlock(b)
			break
		}
		wh.SendOr404(w, &rlt)
	}
}
