package webrpc

import (
	"encoding/json"
	"strconv"
)

// get last blocks
func getLastBlocks(req Request, gateway Gatewayer) Response {
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
