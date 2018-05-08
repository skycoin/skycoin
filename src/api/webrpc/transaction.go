package webrpc

import (
	"encoding/hex"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
)

// TxnResult wraps the daemon.TransactionResult
type TxnResult struct {
	Transaction *daemon.TransactionResult `json:"transaction"`
}

// TxIDJson wraps txid with json tags
type TxIDJson struct {
	Txid string `json:"txid"`
}

func getTransactionHandler(req Request, gateway Gatewayer) Response {
	var txid []string
	if err := req.DecodeParams(&txid); err != nil {
		logger.Critical().Errorf("decode params failed: %v", err)
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(txid) != 1 {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	t, err := cipher.SHA256FromHex(txid[0])
	if err != nil {
		logger.Critical().Errorf("decode txid err: %v", err)
		return MakeErrorResponse(ErrCodeInvalidParams, "invalid transaction hash")
	}
	txn, err := gateway.GetTransaction(t)
	if err != nil {
		logger.Debug(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	if txn == nil {
		return MakeErrorResponse(ErrCodeInvalidRequest, "transaction doesn't exist")
	}

	tx, err := daemon.NewTransactionResult(txn)
	if err != nil {
		logger.Error(err)
		return MakeErrorResponse(ErrCodeInternalError, ErrMsgInternalError)
	}

	return makeSuccessResponse(req.ID, TxnResult{tx})
}

func injectTransactionHandler(req Request, gateway Gatewayer) Response {
	var rawtx []string
	if err := req.DecodeParams(&rawtx); err != nil {
		logger.Critical().Errorf("decode params failed: %v", err)
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	if len(rawtx) != 1 {
		return MakeErrorResponse(ErrCodeInvalidParams, ErrMsgInvalidParams)
	}

	b, err := hex.DecodeString(rawtx[0])
	if err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, fmt.Sprintf("invalid raw transaction: %v", err))
	}

	txn, err := coin.TransactionDeserialize(b)
	if err != nil {
		return MakeErrorResponse(ErrCodeInvalidParams, fmt.Sprintf("%v", err))
	}

	if err := gateway.InjectBroadcastTransaction(txn); err != nil {
		return MakeErrorResponse(ErrCodeInternalError, fmt.Sprintf("inject transaction failed: %v", err))
	}

	return makeSuccessResponse(req.ID, TxIDJson{txn.Hash().Hex()})
}
