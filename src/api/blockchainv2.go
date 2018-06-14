package api

// APIs for blockchain related information (v2)

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor" //http,json helpers
)

// get block by hash or seq api/v2
// method: GET
// url: /block?hash=[:hash]  or /block?seq[:seq]
// params: hash or seq, should only specify one filter.
func getBlockV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		hash := r.FormValue("hash")
		seq := r.FormValue("seq")
		var b *coin.SignedBlock
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

			b, err = gateway.GetSignedBlockByHash(h)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
		case seq != "":
			uSeq, err := strconv.ParseUint(seq, 10, 64)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}

			b, err = gateway.GetSignedBlockBySeq(uSeq)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
		}

		if b == nil {
			wh.Error404(w, "")
			return
		}

		rb, err := visor.NewReadableBlock(&b.Block)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		rbv2, err := NewReadableBlockV2(gateway, rb)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rbv2)
	}
}

func getBlocksV2(gateway Gatewayer) http.HandlerFunc {
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
		rbv2, err := NewReadableBlocksV2(gateway, rb)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("NewReadableBlocksV2: %v", err))
			return
		}
		wh.SendJSONOr500(logger, w, rbv2)
	}
}

// get last N blocks
func getLastBlocksV2(gateway Gatewayer) http.HandlerFunc {
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

		rbv2, err := NewReadableBlocksV2(gateway, rb)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("NewReadableBlocksV2: %v", err))
			return
		}

		wh.SendJSONOr500(logger, w, rbv2)
	}
}
