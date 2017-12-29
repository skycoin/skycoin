package webrpc

import (
	"github.com/skycoin/skycoin/src/util/uxotutil"
)

// OutputsTopn the output json format
type OutputsTopn struct {
	Outputs []uxotutil.AccountJSON `json:"richlist"`
}

//TopnParas the argument for topn outputs
type TopnParas struct {
	Topn                int  `json:"topn"`
	IncludeDistribution bool `json:"include_distribution"`
}

func getRichlistHandler(req Request, gateway Gatewayer) Response {
	topn := TopnParas{}
	if err := req.DecodeParams(&topn); err != nil {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}
	topnAcc, err := gateway.GetRichlist(topn.Topn, topn.IncludeDistribution)
	if err != nil {
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}
	return makeSuccessResponse(req.ID, OutputsTopn{topnAcc})
}
