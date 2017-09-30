package gui

// Blockchain related information for the GUI

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor" //http,json helpers

	"github.com/skycoin/skycoin/src/daemon"
)

const lastBlockNum = 10

// RegisterBlockchainHandlers registers blockchain handlers
func RegisterBlockchainHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	mux.HandleFunc("/blockchain/metadata", blockchainHandler(gateway))
	mux.HandleFunc("/blockchain/progress", blockchainProgressHandler(gateway))

	// get block by hash or seq
	mux.HandleFunc("/block", getBlock(gateway))
	// get block by seq
	// mux.HandleFunc("/block/seq", getBlockBySeq(gateway))
	// get blocks in specific range
	mux.HandleFunc("/blocks", getBlocks(gateway))
	// get last N blocks
	mux.HandleFunc("/last_blocks", getLastBlocks(gateway))
}

func blockchainHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetBlockchainMetadata())
	}
}

func blockchainProgressHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wh.SendOr404(w, gateway.GetBlockchainProgress())
	}
}

// get block by hash or seq
// method: GET
// url: /block?hash=[:hash]  or /block?seq[:seq]
// params: hash or seq, should only specify one filter.
func getBlock(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		hash := r.FormValue("hash")
		seq := r.FormValue("seq")
		var b coin.SignedBlock
		var exist bool
		switch {
		case hash == "" && seq == "":
			wh.Error400(w, "should specify one filter, hash or seq")
			return
		case hash != "" && seq != "":
			wh.Error400(w, "should only specify one filter, hash or seq")
			return
		case hash != "":
			h, err := cipher.SHA256FromHex(hash)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			b, exist = gate.GetBlockByHash(h)
		case seq != "":
			uSeq, err := strconv.ParseUint(seq, 10, 64)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			b, exist = gate.GetBlockBySeq(uSeq)
		}

		if !exist {
			wh.SendOr404(w, nil)
			return
		}

		rb, err := visor.NewReadableBlock(&b.Block)
		if err != nil {
			logger.Error("%v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr404(w, rb)
	}
}

func getBlocks(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
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
		rb, err := gateway.GetBlocks(start, end)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Get blocks failed: %v", err))
			return
		}
		wh.SendOr404(w, rb)
	}
}

// get last N blocks
func getLastBlocks(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		num := r.FormValue("num")
		if num == "" {
			wh.Error400(w, "Param: num is empty")
			return
		}

		n, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		rb, err := gateway.GetLastBlocks(n)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Get last %v blocks failed: %v", n, err))
			return
		}

		wh.SendOr404(w, rb)
	}
}
