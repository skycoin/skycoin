package webrpc

func getStatus(req Request, _ Gatewayer) Response {
	return makeSuccessResponse(req.ID, `{"running": true}`)
}
