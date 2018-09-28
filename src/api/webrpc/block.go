package webrpc

import (
	"github.com/skycoin/skycoin/src/readable"
)

// request params: [seq1, seq2, seq3...]
func getBlocksBySeqHandler(req Request, gateway Gatewayer) Response {
	var seqs []uint64
	if err := req.DecodeParams(&seqs); err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(seqs) == 0 {
		return MakeErrorResponse(ErrCodeInvalidParams, "empty params")
	}
	blocks, err := gateway.GetBlocks(seqs)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	rBlocks, err := readable.NewBlocks(blocks)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	return makeSuccessResponse(req.ID, rBlocks)
}

// request params: [number]
func getLastBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the req params
	var num []uint64
	if err := req.DecodeParams(&num); err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(num) != 1 {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	blocks, err := gateway.GetLastBlocks(num[0])
	if err != nil {
		logger.Errorf("gateway.GetLastBlocks failed: %v", err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	rBlocks, err := readable.NewBlocks(blocks)
	if err != nil {
		logger.Errorf("readable.NewBlocks failed: %v", err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	return makeSuccessResponse(req.ID, rBlocks)
}

func getBlocksHandler(req Request, gateway Gatewayer) Response {
	var params []uint64
	if err := req.DecodeParams(&params); err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(params) != 2 {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	blocks, err := gateway.GetBlocksInRange(params[0], params[1])
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	rBlocks, err := readable.NewBlocks(blocks)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	return makeSuccessResponse(req.ID, rBlocks)
}
