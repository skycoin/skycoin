package webrpc

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
)

// AddrUxoutResult the address uxout json format
type AddrUxoutResult struct {
	Address string                 `json:"address"`
	UxOuts  []readable.SpentOutput `json:"uxouts"`
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

	parsedAddrs := make([]cipher.Address, len(addrs))
	for i, addr := range addrs {
		a, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			logger.Error(err)
			return MakeErrorResponse(ErrCodeInvalidParams, fmt.Sprintf("%v", err))
		}

		parsedAddrs[i] = a
	}

	uxOuts, err := gateway.GetSpentOutputsForAddresses(parsedAddrs)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	for i, uxs := range uxOuts {
		rUxs := readable.NewSpentOutputs(uxs)
		results[i].Address = addrs[i]
		results[i].UxOuts = rUxs
	}

	return makeSuccessResponse(req.ID, &results)
}
