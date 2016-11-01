package webrpc

import "strings"
import "github.com/skycoin/skycoin/src/cipher"
import "fmt"

func getOutputsHandler(req Request, gateway Gatewayer) Response {
	addrs := strings.Split(req.Params["addresses"], ",")
	if len(addrs) == 0 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	// validate those addresses
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid address: %v", a))
		}
	}

	outs := gateway.GetUnspentByAddrs(addrs)
	return makeSuccessResponse(req.ID, outs)
}
