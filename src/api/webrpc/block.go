package webrpc

// request params: [seq1, seq2, seq3...]
func getBlocksBySeqHandler(req Request, gateway Gatewayer) Response {
	var seqs []uint64
	if err := req.DecodeParams(&seqs); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(seqs) == 0 {
		return makeErrorResponse(errCodeInvalidParams, "empty params")
	}
	blocks, err := gateway.GetBlocksInDepth(seqs)
	if err != nil {
		logger.Error("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}
	return makeSuccessResponse(req.ID, blocks)
}

// request params: [number]
func getLastBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the req params
	var num []uint64
	if err := req.DecodeParams(&num); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(num) != 1 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks, err := gateway.GetLastBlocks(num[0])
	if err != nil {
		logger.Error("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}
	return makeSuccessResponse(req.ID, blocks)
}

func getBlocksHandler(req Request, gateway Gatewayer) Response {
	var params []uint64
	if err := req.DecodeParams(&params); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(params) != 2 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks, err := gateway.GetBlocks(params[0], params[1])
	if err != nil {
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}
	return makeSuccessResponse(req.ID, blocks)
}
