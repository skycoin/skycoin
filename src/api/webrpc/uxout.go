package webrpc

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

type addrUxoutResult struct {
	Address     string                 `json:"address"`
	RecvUxOuts  []*historydb.UxOutJSON `json:"recv_uxouts"`
	SpentUxOuts []*historydb.UxOutJSON `json:"spent_uxouts"`
}

func getAddrUxOutsHandler(req Request, gateway Gatewayer) Response {
	var addrs []string
	if err := req.DecodeParams(&addrs); err != nil {
		logger.Critical("decode params failed:%v", err)
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(addrs) == 0 {
		logger.Error("empty request params")
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	results := make([]addrUxoutResult, len(addrs))

	for i, addr := range addrs {
		// decode address
		a, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			logger.Error("%v", err)
			return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("%v", err))
		}
		results[i].Address = addr

		recvUxOuts, err := gateway.GetRecvUxOutOfAddr(a)
		if err != nil {
			logger.Error("%v", err)
			return makeErrorResponse(errCodeInternalError, errMsgInternalError)
		}
		for _, uxout := range recvUxOuts {
			results[i].RecvUxOuts = append(results[i].RecvUxOuts, historydb.NewUxOutJSON(uxout))
		}

		spentUxOuts, err := gateway.GetSpentUxOutOfAddr(a)
		if err != nil {
			logger.Error("%v", err)
			return makeErrorResponse(errCodeInternalError, errMsgInternalError)
		}
		for _, uxout := range spentUxOuts {
			results[i].SpentUxOuts = append(results[i].SpentUxOuts, historydb.NewUxOutJSON(uxout))
		}
	}

	return makeSuccessResponse(req.ID, &results)
}
