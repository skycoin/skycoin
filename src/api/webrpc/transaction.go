package webrpc

import (
	"encoding/hex"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

type TxnResult struct {
	Transaction *visor.TransactionResult
}

func getTransactionHandler(req Request, gateway Gatewayer) Response {
	var txid []string
	if err := req.DecodeParams(&txid); err != nil {
		logger.Criticalf("decode params failed:%v", err)
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(txid) != 1 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	t, err := cipher.SHA256FromHex(txid[0])
	if err != nil {
		logger.Criticalf("decode txid err: %v", err)
		return makeErrorResponse(errCodeInvalidParams, "invalid transaction hash")
	}
	txn, err := gateway.GetTransaction(t)
	if err != nil {
		logger.Debugf("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}

	return makeSuccessResponse(req.ID, TxnResult{txn})
}
