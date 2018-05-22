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
func getPendingTxs(gateway Gatewayer) http.HandlerFunc {
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

		tx, err := gate.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}
		if tx == nil {
			wh.Error404(w, "")
			return
		}

		rbTx, err := visor.NewReadableTransaction(tx)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		resTx := daemon.TransactionResult{
			Transaction: *rbTx,
			Status:      tx.Status,
			Time:        tx.Time,
		}
		wh.SendJSONOr500(logger, w, &resTx)
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
		txRlts, err := daemon.NewTransactionResults(txns)
		if err != nil {
			err = fmt.Errorf("daemon.NewTransactionResults failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txRlts.Txns)
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

func getRawTx(gateway Gatewayer) http.HandlerFunc {
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

		tx, err := gateway.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if tx == nil {
			wh.Error404(w, "")
			return
		}

		d := tx.Txn.Serialize()
		wh.SendJSONOr500(logger, w, hex.EncodeToString(d))
		return
	}
}

func verifyTxHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			wh.Error415(w)
			return
		}

		var v struct {
			EncodedTransaction string `json:"encoded_transaction"`
		}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		tx, err := decodeAndVerifyTx(gateway, v.EncodedTransaction)
		if err != nil {
			wh.Error422(w, err.Error())
			return
		}

		inputs, err := gateway.GetUxBalances(tx.In)
		if err != nil {
			wh.Error503(w, err.Error())
			return
		}

		txRsp, err := NewCreatedTransaction(tx, inputs)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		rsp := struct {
			Transaction interface{} `json:"transaction"`
		}{
			Transaction: txRsp,
		}

		wh.SendJSONOr500(logger, w, rsp)
	}
}

func decodeAndVerifyTx(gateway Gatewayer, encodedTx string) (*coin.Transaction, error) {
	var tx coin.Transaction
	b, err := hex.DecodeString(encodedTx)
	if err != nil {
		return nil, err
	}

	tx, err = coin.TransactionDeserialize(b)
	if err != nil {
		return nil, err
	}

	for _, o := range tx.Out {
		if o.Address.Null() {
			return nil, errors.New("Transaction.Out contains an output sending to an empty address")
		}
	}

	if err := gateway.VerifySingleTxnAllConstraints(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}
