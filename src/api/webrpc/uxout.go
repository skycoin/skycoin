package webrpc

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// AddrUxoutResult the address uxout json format
type AddrUxoutResult struct {
	Address string                 `json:"address"`
	UxOuts  []*historydb.UxOutJSON `json:"uxouts"`
}

func getAddrUxOutsHandler(req Request, gateway Gatewayer) Response {
	var addrs []string
	if err := req.DecodeParams(&addrs); err != nil {
		logger.Critical().Errorf("decode params failed: %v", err)
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(addrs) == 0 {
		logger.Error("empty request params")
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	results := make([]AddrUxoutResult, len(addrs))

	for i, addr := range addrs {
		// decode address
		a, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			logger.Error(err)
			return MakeErrorResponse(ErrCodeInvalidParams, fmt.Sprintf("%v", err))
		}
		results[i].Address = addr
		uxouts, err := gateway.GetAddrUxOuts([]cipher.Address{a})
		if err != nil {
			logger.Error(err)
			return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
		}

		//Convert slice UxOut to slice of UxOutJson
		uxs := make([]*historydb.UxOutJSON, len(uxouts))
		for i, ux := range uxouts {
			uxs[i] = historydb.NewUxOutJSON(ux)
		}

		results[i].UxOuts = append(results[i].UxOuts, uxs...)
	}

	return makeSuccessResponse(req.ID, &results)
}
