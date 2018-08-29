package api

// APIs for blockchain related information

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/visor" //http,json helpers
)

// blockchainProgressHandler returns the blockchain metadata
// Method: GET
// URI: /api/v1/blockchain/metadata
func blockchainHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadata, err := gateway.GetBlockchainMetadata()
		if err != nil {
			err = fmt.Errorf("gateway.GetBlockchainMetadata failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, metadata)
	}
}

// blockchainProgressHandler returns the blockchain sync progress
// Method: GET
// URI: /api/v1/blockchain/progress
func blockchainProgressHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		progress, err := gateway.GetBlockchainProgress()
		if err != nil {
			err = fmt.Errorf("gateway.GetBlockchainProgress failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, progress)
	}
}

func parseBoolFlag(v string) (bool, error) {
	if v == "" {
		return false, nil
	}

	return strconv.ParseBool(v)
}

// blockHandler returns a block by hash or seq
// Method: GET
// URI: /api/v1/block
// Args:
// 	hash [transaction hash string]
//  seq [int]
// 	Note: only one of hash or seq is allowed
func blockHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		hash := r.FormValue("hash")
		seq := r.FormValue("seq")

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		switch {
		case hash == "" && seq == "":
			wh.Error400(w, "should specify one filter, hash or seq")
			return
		case hash != "" && seq != "":
			wh.Error400(w, "should only specify one filter, hash or seq")
			return
		}

		var h cipher.SHA256
		if hash != "" {
			var err error
			h, err = cipher.SHA256FromHex(hash)
			if err != nil {
				wh.Error400(w, err.Error())
				return
			}
		}

		var uSeq uint64
		if seq != "" {
			var err error
			uSeq, err = strconv.ParseUint(seq, 10, 64)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("Invalid seq value %q", seq))
				return
			}
		}

		if verbose {
			var b *visor.ReadableBlockVerbose

			switch {
			case hash != "":
				b, err = gateway.GetBlockByHashVerbose(h)
			case seq != "":
				b, err = gateway.GetBlockBySeqVerbose(uSeq)
			}

			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			if b == nil {
				wh.Error404(w, "")
				return
			}

			wh.SendJSONOr500(logger, w, b)
			return
		}

		var b *coin.SignedBlock
		switch {
		case hash != "":
			b, err = gateway.GetSignedBlockByHash(h)
		case seq != "":
			b, err = gateway.GetSignedBlockBySeq(uSeq)
		}

		if err != nil {
			wh.Error500(w, err.Error())
			return
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

		wh.SendJSONOr500(logger, w, rb)
	}
}

// blocksHandler returns blocks between a start and end point.
// The block sequences include both the start and end point.
// Method: GET
// URI: /api/v1/blocks
// Args:
//	start [int]
//	end [int]
//  verbose [bool]
func blocksHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		sstart := r.FormValue("start")
		start, err := strconv.ParseUint(sstart, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid start value %q", sstart))
			return
		}

		send := r.FormValue("end")
		end, err := strconv.ParseUint(send, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid end value %q", send))
			return
		}

		if verbose {
			rb, err := gateway.GetBlocksVerbose(start, end)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
		} else {
			rb, err := gateway.GetBlocks(start, end)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
		}
	}
}

// lastBlocksHandler returns the most recent N blocks on the blockchain
// Method: GET
// URI: /api/v1/last_blocks
// Args:
//	num [int]
//  verbose [bool]
func lastBlocksHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		num := r.FormValue("num")
		n, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("Invalid num value %q", num))
			return
		}

		if verbose {
			rb, err := gateway.GetLastBlocksVerbose(n)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			wh.SendJSONOr500(logger, w, rb)
		} else {
			rb, err := gateway.GetLastBlocks(n)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			wh.SendJSONOr500(logger, w, rb)
		}
	}
}
