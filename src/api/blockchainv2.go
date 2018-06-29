package api

// APIs for blockchain related information (v2)

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor" //http,json helpers
)

// get block by hash or seq api/v2
// method: GET
// url: /block?hash=[:hash]  or /block?seq[:seq]
// params: hash or seq, should only specify one filter.
func getBlockV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		hash := r.FormValue("hash")
		seq := r.FormValue("seq")
		var b *visor.ReadableBlockV2
		switch {
		case hash == "" && seq == "":
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "should specify one filter, hash or seq")
			writeHTTPResponse(w, resp)
			return
		case hash != "" && seq != "":
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "should only specify one filter, hash or seq")
			writeHTTPResponse(w, resp)
			return
		case hash != "":
			h, err := cipher.SHA256FromHex(hash)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid hash value: %v. %v", hash, err))
				writeHTTPResponse(w, resp)
				return
			}

			b, err = gateway.GetBlockByHashV2(h)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("GetSignedBlockByHash failed: %v, %v", hash, err))
				writeHTTPResponse(w, resp)
				return
			}
		case seq != "":
			uSeq, err := strconv.ParseUint(seq, 10, 64)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid seq value: %v. %v", seq, err))
				writeHTTPResponse(w, resp)
				return
			}

			b, err = gateway.GetBlockBySeqV2(uSeq)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("GetSignedBlockBySeq failed: %v", err))
				writeHTTPResponse(w, resp)
				return
			}
		}

		if b == nil {
			resp := NewHTTPErrorResponse(http.StatusNotFound, "Block not found")
			writeHTTPResponse(w, resp)
			return
		}

		var resp HTTPResponse
		resp.Data = b
		writeHTTPResponse(w, resp)
	}
}

func getBlocksV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}
		sstart := r.FormValue("start")
		start, err := strconv.ParseUint(sstart, 10, 64)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid start value \"%s\"", sstart))
			writeHTTPResponse(w, resp)
			return
		}

		send := r.FormValue("end")
		end, err := strconv.ParseUint(send, 10, 64)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid end value \"%s\"", send))
			writeHTTPResponse(w, resp)
			return
		}
		rbs, err := gateway.GetBlocksV2(start, end)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Get blocks failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}
		var resp HTTPResponse
		resp.Data = rbs
		writeHTTPResponse(w, resp)
	}
}

// get last N blocks
func getLastBlocksV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		num := r.FormValue("num")
		if num == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "Param: num is empty")
			writeHTTPResponse(w, resp)
			return
		}

		n, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("Invalid num value \"%s\"", num))
			writeHTTPResponse(w, resp)
			return
		}

		rb, err := gateway.GetLastBlocks(n)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("Get last %v blocks failed: %v", n, err))
			writeHTTPResponse(w, resp)
			return
		}

		rbv2, err := NewReadableBlocksV2(gateway, rb)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableBlocksV2 failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}
		var resp HTTPResponse
		resp.Data = rbv2
		writeHTTPResponse(w, resp)
	}
}
