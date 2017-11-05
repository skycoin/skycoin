package webrpc

import (
	"fmt"
)

// StatusResult result struct of get_status
type StatusResult struct {
	Running            bool   `json:"running"`
	BlockNum           uint64 `json:"num_of_blocks"`
	LastBlockHash      string `json:"hash_of_last_block"`
	TimeSinceLastBlock string `json:"time_since_last_block"`
}

func getStatusHandler(req Request, gw Gatewayer) Response {
	if len(req.Params) > 0 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	blocks, err := gw.GetLastBlocks(1)
	if err != nil {
		logger.Error("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}
	if len(blocks.Blocks) == 0 {
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}

	b := blocks.Blocks[0]
	return makeSuccessResponse(req.ID, StatusResult{
		Running:            true,
		BlockNum:           b.Head.BkSeq + 1,
		LastBlockHash:      b.Head.BlockHash,
		TimeSinceLastBlock: fmt.Sprintf("%vs", gw.GetTimeNow()-b.Head.Time),
	})
}
