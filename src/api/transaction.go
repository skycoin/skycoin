package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// Returns pending transactions
func getPendingTxns(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txns, err := gateway.GetAllUnconfirmedTxns()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(txns))
		for _, unconfirmedTxn := range txns {
			readable, err := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			ret = append(ret, readable)
		}

		wh.SendJSONOr500(logger, w, &ret)
	}
}

func getTransactionByID(gate Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		txid := r.FormValue("txid")
		if txid == "" {
			wh.Error400(w, "txid is empty")
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		txn, err := gate.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}
		if txn == nil {
			wh.Error404(w, "")
			return
		}

		rbTxn, err := visor.NewReadableTransaction(txn)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		resTxn := daemon.TransactionResult{
			Transaction: *rbTxn,
			Status:      txn.Status,
			Time:        txn.Time,
		}
		wh.SendJSONOr500(logger, w, &resTxn)
	}
}

// Returns transactions that match the filters.
// Method: GET
// URI: /api/v1/transactions
// Args:
//     addrs: Comma seperated addresses [optional, returns all transactions if no address provided]
//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
func getTransactions(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		// Gets 'addrs' parameter value
		addrs, err := parseAddressesFromStr(r.FormValue("addrs"))
		if err != nil {
			wh.Error400(w, fmt.Sprintf("parse parameter: 'addrs' failed: %v", err))
			return
		}

		// Initialize transaction filters
		flts := []visor.TxFilter{visor.AddrsFilter(addrs)}

		// Gets the 'confirmed' parameter value
		confirmedStr := r.FormValue("confirmed")
		if confirmedStr != "" {
			confirmed, err := strconv.ParseBool(confirmedStr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("invalid 'confirmed' value: %v", err))
				return
			}

			flts = append(flts, visor.ConfirmedTxFilter(confirmed))
		}

		// Gets transactions
		txns, err := gateway.GetTransactions(flts...)
		if err != nil {
			err = fmt.Errorf("gateway.GetTransactions failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		// Converts visor.Transaction to daemon.TransactionResult
		txnRlts, err := daemon.NewTransactionResults(txns)
		if err != nil {
			err = fmt.Errorf("daemon.NewTransactionResults failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txnRlts.Txns)
	}
}

// parseAddressesFromStr parses comma seperated addresses string into []cipher.Address
func parseAddressesFromStr(s string) ([]cipher.Address, error) {
	addrsStr := splitCommaString(s)

	var addrs []cipher.Address
	for _, s := range addrsStr {
		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, a)
	}

	return addrs, nil
}

// URI: /api/v1/injectTransaction
// Method: POST
// Content-Type: application/json
// Body: {"rawtx": "<encoded transaction>"}
// Response:
//      400 - bad transaction
//      503 - network unavailable for broadcasting transaction
//      200 - ok, returns the transaction hash in hex as string
func injectTransaction(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}
		// get the rawtransaction
		v := struct {
			Rawtx string `json:"rawtx"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		b, err := hex.DecodeString(v.Rawtx)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		txn, err := coin.TransactionDeserialize(b)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		// TODO -- move this to a more general verification layer, see https://github.com/skycoin/skycoin/issues/1342
		// Check that the transaction does not send to an empty address,
		// if this is happening, assume there is a bug in the code that generated the transaction
		for _, o := range txn.Out {
			if o.Address.Null() {
				wh.Error400(w, "Transaction.Out contains an output sending to an empty address")
				return
			}
		}

		if err := gateway.InjectBroadcastTransaction(txn); err != nil {
			err = fmt.Errorf("inject tx failed: %v", err)
			wh.Error503(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txn.Hash().Hex())
	}
}

func resendUnconfirmedTxns(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		rlt, err := gateway.ResendUnconfirmedTxns()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rlt)
		return
	}
}

func getRawTxn(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txid := r.FormValue("txid")
		if txid == "" {
			wh.Error400(w, "txid is empty")
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		txn, err := gateway.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if txn == nil {
			wh.Error404(w, "")
			return
		}

		d := txn.Txn.Serialize()
		wh.SendJSONOr500(logger, w, hex.EncodeToString(d))
		return
	}
}

// VerifyTxnRequest represents the data struct of the request for /transaction/verify
type VerifyTxnRequest struct {
	EncodedTransaction string `json:"encoded_transaction"`
}

// VerifyTxnResponse the response data struct for /transaction/verify api
type VerifyTxnResponse struct {
	Transaction CreatedTransaction `json:"transaction"`
}

// Decode and verify an encoded transaction
// Method: POST
// URI: /api/v1/transaction/verify
func verifyTxnHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			wh.Error415(w)
			return
		}

		var req VerifyTxnRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		txn, err := decodeTxn(req.EncodedTransaction)
		if err != nil {
			wh.Error422(w, err.Error())
			return
		}

		inputs, err := gateway.VerifyTxnVerbose(txn)
		if err != nil {
			wh.Error422(w, err.Error())
			return
		}

		txnRsp, err := NewCreatedTransaction(txn, inputs)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		rsp := VerifyTxnResponse{Transaction: *txnRsp}

		wh.SendJSONOr500(logger, w, rsp)
	}
}

func decodeTxn(encodedTxn string) (*coin.Transaction, error) {
	var txn coin.Transaction
	b, err := hex.DecodeString(encodedTxn)
	if err != nil {
		return nil, err
	}

	txn, err = coin.TransactionDeserialize(b)
	if err != nil {
		return nil, err
	}

	for _, o := range txn.Out {
		if o.Address.Null() {
			return nil, errors.New("Transaction.Out contains an output sending to an empty address")
		}
	}

	return &txn, nil
}
