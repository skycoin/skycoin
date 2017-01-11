package webrpc

import (
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

type OutputsResult struct {
	Outputs visor.ReadableOutputSet `json:"outputs"`
}

func getOutputsHandler(req Request, gateway Gatewayer) Response {
	var addrs []string
	if err := req.DecodeParams(&addrs); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(addrs) == 0 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	for i, a := range addrs {
		addrs[i] = strings.Trim(a, " ")
	}

	// validate those addresses
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid address: %v", a))
		}
	}

	outs := gateway.GetUnspentOutputs(daemon.FbyAddresses(addrs))
	return makeSuccessResponse(req.ID, OutputsResult{outs})
}
