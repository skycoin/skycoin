package webrpc

import (
	"fmt"
	"time"
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
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	blocks, err := gw.GetLastBlocks(1)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}
	if len(blocks) == 0 {
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	b := blocks[0]
	return makeSuccessResponse(req.ID, StatusResult{
		Running:            true,
		BlockNum:           b.Head.BkSeq + 1,
		LastBlockHash:      b.Head.Hash().Hex(),
		TimeSinceLastBlock: fmt.Sprintf("%vs", uint64(time.Now().UTC().Unix())-b.Head.Time),
	})
}
