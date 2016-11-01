package webrpc

import "strconv"

func getLastBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the req params
	n, err := strconv.ParseUint(req.Params["num"], 10, 64)
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks := gateway.GetLastBlocks(n)
	return makeSuccessResponse(req.ID, blocks)
}

func getBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the params
	start, end := req.Params["start"], req.Params["end"]
	if start == "" {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if end == "" {
		end = start
	}

	s, err := strconv.ParseUint(start, 10, 64)
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	e, err := strconv.ParseUint(end, 10, 64)
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks := gateway.GetBlocks(s, e)
	return makeSuccessResponse(req.ID, blocks)
}
