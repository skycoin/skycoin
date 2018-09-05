package api

// APIs for blockchain related information

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http" // http,json helpers
	"github.com/skycoin/skycoin/src/visor"
)

// blockchainMetadataHandler returns the blockchain metadata
// Method: GET
// URI: /api/v1/blockchain/metadata
func blockchainMetadataHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		visorMetadata, err := gateway.GetBlockchainMetadata()
		if err != nil {
			err = fmt.Errorf("gateway.GetBlockchainMetadata failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		// This can happen if the node is shut down at the right moment, guard against a panic
		if visorMetadata == nil {
			err = errors.New("gateway.GetBlockchainMetadata metadata is nil")
			wh.Error500(w, err.Error())
			return
		}

		metadata := readable.NewBlockchainMetadata(*visorMetadata)

		wh.SendJSONOr500(logger, w, metadata)
	}
}

// blockchainProgressHandler returns the blockchain sync progress
// Method: GET
// URI: /api/v1/blockchain/progress
func blockchainProgressHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		progress, err := gateway.GetBlockchainProgress()
		if err != nil {
			err = fmt.Errorf("gateway.GetBlockchainProgress failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		// This can happen if the node is shut down at the right moment, guard against a panic
		if progress == nil {
			err = errors.New("gateway.GetBlockchainProgress progress is nil")
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, readable.NewBlockchainProgress(progress))
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
			var b *coin.SignedBlock
			var inputs [][]visor.TransactionInput

			switch {
			case hash != "":
				b, inputs, err = gateway.GetSignedBlockByHashVerbose(h)
			case seq != "":
				b, inputs, err = gateway.GetSignedBlockBySeqVerbose(uSeq)
			}

			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			if b == nil {
				wh.Error404(w, "")
				return
			}

			rb, err := readable.NewBlockVerbose(&b.Block, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
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

		rb, err := readable.NewBlock(&b.Block)
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
			blocks, inputs, err := gateway.GetBlocksInRangeVerbose(start, end)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rb, err := readable.NewBlocksVerbose(blocks, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
		} else {
			blocks, err := gateway.GetBlocksInRange(start, end)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rb, err := readable.NewBlocks(blocks)
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
			blocks, inputs, err := gateway.GetLastBlocksVerbose(n)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rb, err := readable.NewBlocksVerbose(blocks, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
			return
		}

		blocks, err := gateway.GetLastBlocks(n)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		rb, err := readable.NewBlocks(blocks)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rb)
	}
}
