package webrpc

import (
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

// OutputsResult the output json format
type OutputsResult struct {
	Outputs visor.ReadableOutputSet `json:"outputs"`
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

	outs, err := gateway.GetUnspentOutputs(daemon.FbyAddresses(addrs))
	if err != nil {
		logger.Errorf("get unspent outputs failed: %v", err)
		return MakeErrorResponse(ErrCodeInternalError, fmt.Sprintf("gateway.GetUnspentOutputs failed: %v", err))
	}

	return makeSuccessResponse(req.ID, OutputsResult{*outs})
}
