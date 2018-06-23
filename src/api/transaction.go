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
	"github.com/skycoin/skycoin/src/wallet"

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

// VerifyTxnRequest represents the data struct of the request for /api/v2/transaction/verify
type VerifyTxnRequest struct {
	EncodedTransaction string `json:"encoded_transaction"`
}

// VerifyTxnResponse the response data struct for /api/v2/transaction/verify
type VerifyTxnResponse struct {
	Confirmed   bool               `json:"confirmed"`
	Transaction CreatedTransaction `json:"transaction"`
}

func writeHTTPResponse(w http.ResponseWriter, resp HTTPResponse) {
	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		wh.Error500(w, "json.MarshalIndent failed")
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if resp.Error == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		if resp.Error.Code < 400 || resp.Error.Code >= 600 {
			logger.Critical().Errorf("writeHTTPResponse invalid error status code: %d", resp.Error.Code)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(resp.Error.Code)
		}
	}

	if _, err := w.Write(out); err != nil {
		logger.WithError(err).Error("http Write failed")
	}
}

// Decode and verify an encoded transaction
// Method: POST
// URI: /api/v2/transaction/verify
func verifyTxnHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req VerifyTxnRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := decodeTxn(req.EncodedTransaction)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("decode transaction failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		var resp HTTPResponse
		inputs, isTxnConfirmed, err := gateway.VerifyTxnVerbose(txn)
		if err != nil {
			switch err.(type) {
			case visor.ErrTxnViolatesSoftConstraint,
				visor.ErrTxnViolatesHardConstraint,
				visor.ErrTxnViolatesUserConstraint:
				resp.Error = &HTTPError{
					Code:    http.StatusUnprocessableEntity,
					Message: err.Error(),
				}
			default:
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
				writeHTTPResponse(w, resp)
				return
			}
		}

		verifyTxnResp := VerifyTxnResponse{
			Confirmed: isTxnConfirmed,
		}

		if len(inputs) != len(txn.In) {
			inputs = nil
		}
		verboseTxn, err := newCreatedTransactionFuzzy(txn, inputs)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		verifyTxnResp.Transaction = *verboseTxn

		resp.Data = verifyTxnResp

		if isTxnConfirmed && resp.Error == nil {
			resp.Error = &HTTPError{
				Code:    http.StatusUnprocessableEntity,
				Message: "transaction has been spent",
			}
		}

		writeHTTPResponse(w, resp)
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

	return &txn, nil
}

// newCreatedTransactionFuzzy creates a CreatedTransaction but accomodates possibly invalid txn input
func newCreatedTransactionFuzzy(txn *coin.Transaction, inputs []wallet.UxBalance) (*CreatedTransaction, error) {
	if len(txn.In) != len(inputs) && len(inputs) != 0 {
		return nil, errors.New("len(txn.In) != len(inputs)")
	}

	var outputHours uint64
	var feeInvalid bool
	for _, o := range txn.Out {
		var err error
		outputHours, err = coin.AddUint64(outputHours, o.Hours)
		if err != nil {
			feeInvalid = true
		}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = coin.AddUint64(inputHours, i.Hours)
		if err != nil {
			feeInvalid = true
		}
	}

	if inputHours < outputHours {
		feeInvalid = true
	}

	var fee uint64
	if !feeInvalid {
		fee = inputHours - outputHours
	}

	sigs := make([]string, len(txn.Sigs))
	for i, s := range txn.Sigs {
		sigs[i] = s.Hex()
	}

	txid := txn.Hash()
	out := make([]CreatedTransactionOutput, len(txn.Out))
	for i, o := range txn.Out {
		co, err := NewCreatedTransactionOutput(o, txid)
		if err != nil {
			logger.WithError(err).Error("NewCreatedTransactionOutput failed")
			continue
		}
		out[i] = *co
	}

	in := make([]CreatedTransactionInput, len(txn.In))
	if len(inputs) == 0 {
		for i, h := range txn.In {
			in[i] = CreatedTransactionInput{
				UxID: h.Hex(),
			}
		}
	} else {
		for i, o := range inputs {
			ci, err := NewCreatedTransactionInput(o)
			if err != nil {
				logger.WithError(err).Error("NewCreatedTransactionInput failed")
				continue
			}
			in[i] = *ci
		}
	}

	return &CreatedTransaction{
		Length:    txn.Length,
		Type:      txn.Type,
		TxID:      txid.Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fmt.Sprint(fee),

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}
