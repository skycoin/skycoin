package webrpc

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

	blocks := gateway.GetLastBlocks(num[0])
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

	blocks := gateway.GetBlocks(params[0], params[1])
	return makeSuccessResponse(req.ID, blocks)
}
