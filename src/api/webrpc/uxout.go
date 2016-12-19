package webrpc

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

type addrUxoutResult struct {
	Address string                 `json:"address"`
	UxOuts  []*historydb.UxOutJSON `json:"uxouts"`
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

		uxouts, err := gateway.GetAddrUxOuts(a)
		if err != nil {
			logger.Error("%v", err)
			return makeErrorResponse(errCodeInternalError, errMsgInternalError)
		}
		for _, uxout := range uxouts {
			results[i].UxOuts = append(results[i].UxOuts, historydb.NewUxOutJSON(uxout))
		}
	}

	return makeSuccessResponse(req.ID, &results)
}
