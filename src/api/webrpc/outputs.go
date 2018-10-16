package webrpc

import (
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/visor"
)

// OutputsResult the output json format
type OutputsResult struct {
	Outputs readable.UnspentOutputsSummary `json:"outputs"`
}

func getOutputsHandler(req Request, gateway Gatewayer) Response {
	var addrs []string
	if err := req.DecodeParams(&addrs); err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(addrs) == 0 {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	for i, a := range addrs {
		addrs[i] = strings.Trim(a, " ")
	}

	// validate those addresses
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return MakeErrorResponse(ErrCodeInvalidParams, fmt.Sprintf("invalid address: %v", a))
		}
	}

	summary, err := gateway.GetUnspentOutputsSummary([]visor.OutputsFilter{visor.FbyAddresses(addrs)})
	if err != nil {
		logger.Errorf("get unspent outputs failed: %v", err)
		return MakeErrorResponse(ErrCodeInternalError, fmt.Sprintf("gateway.GetUnspentOutputsSummary failed: %v", err))
	}

	rSummary, err := readable.NewUnspentOutputsSummary(summary)
	if err != nil {
		logger.Error(err.Error())
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	return makeSuccessResponse(req.ID, OutputsResult{*rSummary})
}
