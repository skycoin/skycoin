package webrpc

import (
	"encoding/json"
	"strconv"
)

// get last blocks
func getLastBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the req params
	n, err := strconv.ParseUint(req.Params["num"], 10, 64)
	if err != nil {
		return makeErrorResponse("", &RPCError{
			Code:    errCodeInvalidParams,
			Message: errMsgInvalidParams,
		})
	}

	blocks := gateway.GetLastBlocks(n)
	d, err := json.Marshal(blocks)
	if err != nil {
		logger.Errorf("%v", err)
		return makeErrorResponse("", &RPCError{
			Code:    errCodeInternalError,
			Message: errMsgInternalErr,
		})
	}

	return makeSuccessResponse(req.ID, string(d))
}

func getBlocksHandler(req Request, gateway Gatewayer) Response {
	// validate the params
	start, end := req.Params["start"], req.Params["end"]
	s, err := strconv.ParseUint(start, 10, 64)
	if err != nil {
		return makeErrorResponse("", &RPCError{
			Code:    errCodeInvalidParams,
			Message: errMsgInvalidParams,
		})
	}

	e, err := strconv.ParseUint(end, 10, 64)
	if err != nil {
		return makeErrorResponse("", &RPCError{
			Code:    errCodeInvalidParams,
			Message: errMsgInvalidParams,
		})
	}

	blocks := gateway.GetBlocks(s, e)
	d, err := json.Marshal(blocks)
	if err != nil {
		logger.Errorf("%v", err)
		return makeErrorResponse("", &RPCError{
			Code:    errCodeInternalError,
			Message: errMsgInternalErr,
		})
	}

	return makeSuccessResponse(req.ID, string(d))
}
