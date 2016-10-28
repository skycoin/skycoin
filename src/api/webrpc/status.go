package webrpc

func getStatus(req Request) Response {
	return makeSuccessResponse(req.ID, "running")
}
