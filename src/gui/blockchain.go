// Blockchain related information for the GUI
package gui

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

func RegisterBlockchainHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	mux.HandleFunc("/blockchain/metadata", blockchainHandler(gateway))
	mux.HandleFunc("/blockchain/progress", blockchainProgressHandler(gateway))

	// get block by hash or seq
	mux.HandleFunc("/block", getBlock(gateway))
	// get block by seq
	// mux.HandleFunc("/block/seq", getBlockBySeq(gateway))
	// get blocks in specific range
	mux.HandleFunc("/blocks", getBlocks(gateway))
	// get last 10 blocks
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
		if r.Method != "GET" {
			wh.Error405(w, "")
			return
		}
		hash := r.FormValue("hash")
		seq := r.FormValue("seq")
		var b *coin.Block
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

			b = gate.V.GetBlockByHash(h)
		case seq != "":
			uSeq, err := strconv.ParseUint(seq, 10, 64)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			b = gate.V.GetBlockBySeq(uSeq)
		}

		if b == nil {
			wh.SendOr404(w, nil)
			return
		}
		wh.SendOr404(w, visor.NewReadableBlock(b))
	}
}

func getBlocks(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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
		wh.SendOr404(w, gateway.GetBlocks(start, end))
	}
}

// get last N blocks
func getLastBlocks(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		wh.SendOr404(w, gateway.GetLastBlocks(n))
	}
}
