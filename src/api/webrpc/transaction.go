package webrpc

import (
	"encoding/hex"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
)

type TxnResult struct {
	Transaction *visor.TransactionResult `json:"transaction"`
}

type InjectResult struct {
	Txid string `json:"txid"`
}

func getTransactionHandler(req Request, gateway Gatewayer) Response {
	var txid []string
	if err := req.DecodeParams(&txid); err != nil {
		logger.Critical("decode params failed:%v", err)
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(txid) != 1 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	t, err := cipher.SHA256FromHex(txid[0])
	if err != nil {
		logger.Critical("decode txid err: %v", err)
		return makeErrorResponse(errCodeInvalidParams, "invalid transaction hash")
	}
	txn, err := gateway.GetTransaction(t)
	if err != nil {
		logger.Debugf("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}

	return makeSuccessResponse(req.ID, TxnResult{txn})
}

func injectTransactionHandler(req Request, gateway Gatewayer) Response {
	var rawtx []string
	if err := req.DecodeParams(&rawtx); err != nil {
		logger.Critical("decode params failed:%v", err)
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(rawtx) != 1 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	b, err := hex.DecodeString(rawtx[0])
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid raw transaction:%v", err))
	}

	txn := coin.TransactionDeserialize(b)
	t, err := gateway.InjectTransaction(txn)
	if err != nil {
		return makeErrorResponse(errCodeInternalError, fmt.Sprintf("inject transaction failed:%v", err))
	}

	return makeSuccessResponse(req.ID, InjectResult{t.Hash().Hex()})
}
