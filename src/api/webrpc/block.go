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
	// validate the params
	var param struct {
		Start *uint64
		End   *uint64
	}

	if err := req.DecodeParams(&param); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if param.Start == nil || param.End == nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks := gateway.GetBlocks(*param.Start, *param.End)
	return makeSuccessResponse(req.ID, blocks)
}
