package webrpc

// StatusResult result struct of get_status
type StatusResult struct {
	Running bool `json:"running"`
}

func getStatusHandler(req Request, _ Gatewayer) Response {
	if len(req.Params) > 0 {
		return makeErrorResponse(&RPCError{
			Code:    errCodeInvalidParams,
			Message: errMsgInvalidParams,
		})
	}
	return makeSuccessResponse(ptrString(req.ID), StatusResult{true})
}
