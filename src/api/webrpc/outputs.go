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

	outs, err := gateway.GetUnspentOutputs(daemon.FbyAddresses(addrs))
	if err != nil {
		logger.Error("get unspent outputs failed: %v", err)
		return makeErrorResponse(errCodeInternalError, fmt.Sprintf("gateway.GetUnspentOutputs failed: %v", err))
	}

	return makeSuccessResponse(req.ID, OutputsResult{outs})
}

func getOutputsWithFiltersHandler(req Request, gateway Gatewayer) Response {
	var filters map[string][]string
	if err := req.DecodeParams(&filters); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	addrs, okaddr := filters["addrs"]
	hashes, okhash := filters["hashes"]
	addrLen := len(addrs)
	hashLen := len(hashes)

	// Check that either addrs or hashes are provided
	if (!okaddr && !okhash) || (addrLen == 0 && hashLen == 0) {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	// Trim whitespace from addrs and hashes
	for i, a := range addrs {
		addrs[i] = strings.Trim(a, " ")
	}
	for i, h := range hashes {
		hashes[i] = strings.Trim(h, " ")
	}

	// validate those addresses
	for _, a := range addrs {
		if _, err := cipher.DecodeBase58Address(a); err != nil {
			return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid address: %v", a))
		}
	}

	// Filter unspent outputs
	outputFilters := []daemon.OutputsFilter{}
	if addrLen > 0 && hashLen > 0 {
		// Filter by addr and hash
		outputFilters = append(outputFilters, daemon.FbyAddresses(addrs))
		outputFilters = append(outputFilters, daemon.FbyHashes(hashes))
	} else if addrLen > 0 {
		// Filter by addr
		outputFilters = append(outputFilters, daemon.FbyAddresses(addrs))
	} else {
		// Filter by hash
		outputFilters = append(outputFilters, daemon.FbyHashes(hashes))
	}

	outs, err := gateway.GetUnspentOutputs(outputFilters...)
	if err != nil {
		logger.Error("get unspent outputs with filters failed: %v", err)
		return makeErrorResponse(errCodeInternalError)
	}

	return makeSuccessResponse(req.ID, OutputsResult{outs})
}
