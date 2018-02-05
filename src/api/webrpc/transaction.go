package webrpc

import (
	"encoding/hex"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// TxnResult wraps the visor.TransactionResult
type TxnResult struct {
	Transaction *visor.TransactionResult `json:"transaction"`
}

// TxIDJson wraps txid with json tags
type TxIDJson struct {
	Txid string `json:"txid"`
}

func getTransactionHandler(req Request, gateway Gatewayer) Response {
	var txid []string
	if err := req.DecodeParams(&txid); err != nil {
		logger.Critical("decode params failed: %v", err)
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

	if txn == nil {
		return makeErrorResponse(errCodeInvalidRequest, "transaction doesn't exist")
	}

	txInputsData, err := GetTransactionInputsData(&txn.Txn, gateway)
	if err != nil {
		return makeErrorResponse(errCodeInternalError, "invalid transaction")
	}

	tx, err := visor.NewTransactionResult(txn, txInputsData)
	if err != nil {
		logger.Error("%v", err)
		return makeErrorResponse(errCodeInternalError, errMsgInternalError)
	}

	return makeSuccessResponse(req.ID, TxnResult{tx})
}

func injectTransactionHandler(req Request, gateway Gatewayer) Response {
	var rawtx []string
	if err := req.DecodeParams(&rawtx); err != nil {
		logger.Critical("decode params failed: %v", err)
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	if len(rawtx) != 1 {
		return makeErrorResponse(errCodeInvalidParams, errMsgInvalidParams)
	}

	b, err := hex.DecodeString(rawtx[0])
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("invalid raw transaction: %v", err))
	}

	txn, err := coin.TransactionDeserialize(b)
	if err != nil {
		return makeErrorResponse(errCodeInvalidParams, fmt.Sprintf("%v", err))
	}

	if err := gateway.InjectBroadcastTransaction(txn); err != nil {
		return makeErrorResponse(errCodeInternalError, fmt.Sprintf("inject transaction failed: %v", err))
	}

	return makeSuccessResponse(req.ID, TxIDJson{txn.Hash().Hex()})
}

// GetTransactionInputsData returns the inputs data of a transaction
func GetTransactionInputsData(tx *coin.Transaction, gateway Gatewayer) ([]*historydb.UxOut, error) {
	txInputsData := make([]*historydb.UxOut, 0, len(tx.In))

	for _, in := range tx.In {

		uxout, err := gateway.GetUxOutByID(in)
		if err != nil {
			logger.Error("%v", err)
			return nil, err
		}
		if uxout == nil {
			logger.Error("uxout of %d does not exist in history db", in)
			return nil, fmt.Errorf("uxout of %d does not exist in history db", in)
		}

		txInputsData = append(txInputsData, uxout)
	}

	return txInputsData, nil
}
