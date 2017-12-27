package webrpc

import (
	"fmt"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util/uxotutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/spaco/spo/src/util/droplet"
)

// OutputsResult the output json format
type OutputsResult struct {
	Outputs visor.ReadableOutputSet `json:"outputs"`
}

// OutputsTopn the output json format
type OutputsTopn struct {
	Outputs []uxotutil.AccountJSON `json:"richlist"`
}

//TopnParas the argument for topn outputs
type TopnParas struct {
	Topn                int  `json:"topn"`
	IncludeDistribution bool `json:"include_distribution"`
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
		return makeErrorResponse(errCodeInternalError)
	}
	return makeSuccessResponse(req.ID, OutputsResult{outs})
}

func getTopnUxoutHandler(req Request, gateway Gatewayer) Response {
	topn := TopnParas{}
	if err := req.DecodeParams(&topn); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}
	outsall, err := gateway.GetUnspentOutputs(daemon.FbyAddressesNotIncluded([]string{}))
	if err != nil {
		logger.Error("get topn unspent outputs failed: %v", err)
		return makeErrorResponse(errCodeInternalError)
	}

	allAccounts := map[string]uint64{}
	for _, out := range outsall.HeadOutputs {
		amt, err := droplet.FromString(out.Coins)
		if err != nil {
			logger.Error("get topn unspent outputs failed: %v", err)
			return makeErrorResponse(errCodeInternalError)
		}
		if _, ok := allAccounts[out.Address]; ok {
			allAccounts[out.Address] += amt
		} else {
			allAccounts[out.Address] = amt
		}
	}
	distributionMap := getDistributiomAddressMap()
	amgr := uxotutil.NewAccountMgr(allAccounts, distributionMap)
	amgr.Sort()
	topnAcc, err := amgr.GetTopn(topn.Topn, topn.IncludeDistribution)
	if err != nil {
		logger.Error("get topn unspent outputs failed: %v", err)
		return makeErrorResponse(errCodeInternalError)
	}

	return makeSuccessResponse(req.ID, OutputsTopn{topnAcc})
}

func getDistributiomAddressMap() map[string]struct{} {
	distributionMap := map[string]struct{}{}
	addresses := visor.GetLockedDistributionAddresses()
	for _, address := range addresses {
		distributionMap[address] = struct{}{}
	}
	return distributionMap
}
