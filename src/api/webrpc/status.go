package webrpc

func getStatusHandler(req Request, _ Gatewayer) Response {
	return makeSuccessResponse(req.ID, `{"running": true}`)
}
