package webrpc

func getStatusHandler(req Request, _ Gatewayer) Response {
	return makeSuccessResponse(ptrString(req.ID), ptrString(`{"running": true}`))
}
