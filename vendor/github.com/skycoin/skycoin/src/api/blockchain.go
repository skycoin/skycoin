package api

// APIs for blockchain related information

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
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

		headSeq, _, err := gateway.HeadBkSeq()
		if err != nil {
			err = fmt.Errorf("gateway.HeadBkSeq failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		progress := gateway.GetBlockchainProgress(headSeq)

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

			rb, err := readable.NewBlockVerbose(b.Block, inputs)
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

		rb, err := readable.NewBlock(b.Block)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rb)
	}
}

// blocksHandler returns blocks between a start and end point,
// or an explicit list of sequences.
// If using start and end, the block sequences include both the start and end point.
// Explicit sequences cannot be combined with start and end.
// Method: GET, POST
// URI: /api/v1/blocks
// Args:
//	start [int]
//	end [int]
//  seqs [comma separated list of ints]
//  verbose [bool]
func blocksHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		sStart := r.FormValue("start")
		sEnd := r.FormValue("end")
		sSeqs := r.FormValue("seqs")

		if sSeqs != "" && (sStart != "" || sEnd != "") {
			wh.Error400(w, "seqs cannot be used with start or end")
			return
		}

		if sSeqs == "" && sStart == "" && sEnd == "" {
			wh.Error400(w, "At least one of seqs or start or end are required")
			return
		}

		var start uint64
		var end uint64
		var seqs []uint64

		if sStart != "" {
			var err error
			start, err = strconv.ParseUint(sStart, 10, 64)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("Invalid start value %q", sStart))
				return
			}
		}

		if sEnd != "" {
			var err error
			end, err = strconv.ParseUint(sEnd, 10, 64)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("Invalid end value %q", sEnd))
				return
			}
		}

		if sSeqs != "" {
			ssSeqs := strings.Split(sSeqs, ",")
			seqs = make([]uint64, len(ssSeqs))
			seqsMap := make(map[uint64]struct{}, len(ssSeqs))
			for i, s := range ssSeqs {
				x, err := strconv.ParseUint(s, 10, 64)
				if err != nil {
					wh.Error400(w, fmt.Sprintf("Invalid sequence %q at seqs[%d]", s, i))
					return
				}

				if _, ok := seqsMap[x]; ok {
					wh.Error400(w, fmt.Sprintf("Duplicate sequence %d at seqs[%d]", x, i))
					return
				}
				seqsMap[x] = struct{}{}

				seqs[i] = x
			}
		}

		if verbose {
			var blocks []coin.SignedBlock
			var inputs [][][]visor.TransactionInput
			var err error

			if len(seqs) > 0 {
				blocks, inputs, err = gateway.GetBlocksVerbose(seqs)
			} else {
				blocks, inputs, err = gateway.GetBlocksInRangeVerbose(start, end)
			}

			if err != nil {
				switch err.(type) {
				case visor.ErrBlockNotExist:
					wh.Error404(w, err.Error())
				default:
					wh.Error500(w, err.Error())
				}
				return
			}

			rb, err := readable.NewBlocksVerbose(blocks, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rb)
		} else {
			var blocks []coin.SignedBlock
			var err error

			if len(seqs) > 0 {
				blocks, err = gateway.GetBlocks(seqs)
			} else {
				blocks, err = gateway.GetBlocksInRange(start, end)
			}

			if err != nil {
				switch err.(type) {
				case visor.ErrBlockNotExist:
					wh.Error404(w, err.Error())
				default:
					wh.Error500(w, err.Error())
				}
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
