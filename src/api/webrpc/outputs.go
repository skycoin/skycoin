package webrpc

import "strings"
import "github.com/skycoin/skycoin/src/cipher"
import "fmt"
import "github.com/skycoin/skycoin/src/visor"

type OutputsResult struct {
	Outputs []visor.ReadableOutput `json:"outputs"`
}

func getOutputsHandler(req Request, gateway Gatewayer) Response {
	addrStr := req.Params["addresses"]
	if addrStr == "" {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	addrs := strings.Split(addrStr, ",")
	for i, a := range addrs {
		addrs[i] = strings.Trim(a, " ")
	}

	// validate those addresses
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid address: %v", a))
		}
	}

	outs := gateway.GetUnspentByAddrs(addrs)
	return makeSuccessResponse(req.ID, OutputsResult{outs})
}
