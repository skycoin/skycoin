package api

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/visor"
)

// pendingTxnsHandler returns pending (unconfirmed) transactions
// Method: GET
// URI: /api/v1/pendingTxs
// Args:
//	verbose: [bool] include verbose transaction input data
func pendingTxnsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		if verbose {
			txns, inputs, err := gateway.GetAllUnconfirmedTransactionsVerbose()
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			vb, err := readable.NewUnconfirmedTransactionsVerbose(txns, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, vb)
		} else {
			txns, err := gateway.GetAllUnconfirmedTransactions()
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			ret, err := readable.NewUnconfirmedTransactions(txns)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, ret)
		}
	}
}

// TransactionEncodedResponse represents the data struct of the response to /api/v1/transaction?encoded=1
type TransactionEncodedResponse struct {
	Status             readable.TransactionStatus `json:"status"`
	Time               uint64                     `json:"time"`
	EncodedTransaction string                     `json:"encoded_transaction"`
}

// transactionHandler returns a transaction identified by its txid hash
// Method: GET
// URI: /api/v1/transaction
// Args:
//	txid: transaction hash
//	verbose: [bool] include verbose transaction input data
//  encoded: [bool] return as a raw encoded transaction
func transactionHandler(gateway Gatewayer) http.HandlerFunc {
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

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		encoded, err := parseBoolFlag(r.FormValue("encoded"))
		if err != nil {
			wh.Error400(w, "Invalid value for encoded")
			return
		}

		if verbose && encoded {
			wh.Error400(w, "verbose and encoded cannot be combined")
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if verbose {
			txn, inputs, err := gateway.GetTransactionWithInputs(h)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			if txn == nil {
				wh.Error404(w, "")
				return
			}

			rTxn, err := readable.NewTransactionWithStatusVerbose(txn, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, rTxn)
			return
		}

		txn, err := gateway.GetTransaction(h)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		if txn == nil {
			wh.Error404(w, "")
			return
		}

		if encoded {
			txnHex, err := txn.Transaction.SerializeHex()
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, TransactionEncodedResponse{
				EncodedTransaction: txnHex,
				Status:             readable.NewTransactionStatus(txn.Status),
				Time:               txn.Time,
			})
			return
		}

		rTxn, err := readable.NewTransactionWithStatus(txn)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, rTxn)
	}
}

// TransactionsWithStatus array of transaction results
type TransactionsWithStatus struct {
	Transactions []readable.TransactionWithStatus `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatus) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatus converts []Transaction to TransactionsWithStatus
func NewTransactionsWithStatus(txns []visor.Transaction) (*TransactionsWithStatus, error) {
	txnRlts := make([]readable.TransactionWithStatus, 0, len(txns))
	for _, txn := range txns {
		rTxn, err := readable.NewTransactionWithStatus(&txn)
		if err != nil {
			return nil, err
		}
		txnRlts = append(txnRlts, *rTxn)
	}

	return &TransactionsWithStatus{
		Transactions: txnRlts,
	}, nil
}

// TransactionsWithStatusVerbose array of transaction results
type TransactionsWithStatusVerbose struct {
	Transactions []readable.TransactionWithStatusVerbose `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatusVerbose) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatusVerbose converts []Transaction to []TransactionsWithStatusVerbose
func NewTransactionsWithStatusVerbose(txns []visor.Transaction, inputs [][]visor.TransactionInput) (*TransactionsWithStatusVerbose, error) {
	if len(txns) != len(inputs) {
		return nil, errors.New("NewTransactionsWithStatusVerbose: len(txns) != len(inputs)")
	}

	txnRlts := make([]readable.TransactionWithStatusVerbose, len(txns))
	for i, txn := range txns {
		rTxn, err := readable.NewTransactionWithStatusVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}
		txnRlts[i] = *rTxn
	}

	return &TransactionsWithStatusVerbose{
		Transactions: txnRlts,
	}, nil
}

// Returns transactions that match the filters.
// Method: GET, POST
// URI: /api/v1/transactions
// Args:
//     addrs: Comma separated addresses [optional, returns all transactions if no address provided]
//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
//	   verbose: [bool] include verbose transaction input data
func transactionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		// Gets 'addrs' parameter value
		addrs, err := parseAddressesFromStr(r.FormValue("addrs"))
		if err != nil {
			wh.Error400(w, fmt.Sprintf("parse parameter: 'addrs' failed: %v", err))
			return
		}

		// Initialize transaction filters
		flts := []visor.TxFilter{visor.NewAddrsFilter(addrs)}

		// Gets the 'confirmed' parameter value
		confirmedStr := r.FormValue("confirmed")
		if confirmedStr != "" {
			confirmed, err := strconv.ParseBool(confirmedStr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("invalid 'confirmed' value: %v", err))
				return
			}

			flts = append(flts, visor.NewConfirmedTxFilter(confirmed))
		}

		if verbose {
			txns, inputs, err := gateway.GetTransactionsWithInputs(flts)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns, err := NewTransactionsWithStatusVerbose(txns, inputs)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns.Sort()

			wh.SendJSONOr500(logger, w, rTxns.Transactions)
		} else {
			txns, err := gateway.GetTransactions(flts)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns, err := NewTransactionsWithStatus(txns)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			rTxns.Sort()

			wh.SendJSONOr500(logger, w, rTxns.Transactions)
		}
	}
}

// InjectTransactionRequest is sent to POST /api/v1/injectTransaction
type InjectTransactionRequest struct {
	RawTxn      string `json:"rawtx"`
	NoBroadcast bool   `json:"no_broadcast,omitempty"`
}

// URI: /api/v1/injectTransaction
// Method: POST
// Content-Type: application/json
// Body: {"rawtx": "<hex encoded transaction>"}
// Response:
//      200 - ok, returns the transaction hash in hex as string
//      400 - bad transaction
//		500 - other error
//      503 - network unavailable for broadcasting transaction
func injectTransactionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		var v InjectTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if v.RawTxn == "" {
			wh.Error400(w, "rawtx is required")
			return
		}

		txn, err := coin.DeserializeTransactionHex(v.RawTxn)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if v.NoBroadcast {
			if err := gateway.InjectTransaction(txn); err != nil {
				switch err.(type) {
				case visor.ErrTxnViolatesUserConstraint,
					visor.ErrTxnViolatesHardConstraint,
					visor.ErrTxnViolatesSoftConstraint:
					wh.Error400(w, err.Error())
				default:
					wh.Error500(w, err.Error())
				}
				return
			}
		} else {
			if err := gateway.InjectBroadcastTransaction(txn); err != nil {
				switch err.(type) {
				case visor.ErrTxnViolatesUserConstraint,
					visor.ErrTxnViolatesHardConstraint,
					visor.ErrTxnViolatesSoftConstraint:
					wh.Error400(w, err.Error())
				default:
					if daemon.IsBroadcastFailure(err) {
						wh.Error503(w, err.Error())
					} else {
						wh.Error500(w, err.Error())
					}
				}
				return
			}
		}

		wh.SendJSONOr500(logger, w, txn.Hash().Hex())
	}
}

// ResendResult the result of rebroadcasting transaction
type ResendResult struct {
	Txids []string `json:"txids"`
}

// NewResendResult creates a ResendResult from a list of transaction ID hashes
func NewResendResult(hashes []cipher.SHA256) ResendResult {
	txids := make([]string, len(hashes))
	for i, h := range hashes {
		txids[i] = h.Hex()
	}
	return ResendResult{
		Txids: txids,
	}
}

// URI: /api/v1/resendUnconfirmedTxns
// Method: POST
// Broadcasts all unconfirmed transactions from the unconfirmed transaction pool
// Response:
//      200 - ok, returns the transaction hashes that were resent
//      405 - method not POST
//		500 - other error
//      503 - network unavailable for broadcasting transaction
func resendUnconfirmedTxnsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		hashes, err := gateway.ResendUnconfirmedTxns()
		if err != nil {
			if daemon.IsBroadcastFailure(err) {
				wh.Error503(w, err.Error())
			} else {
				wh.Error500(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, NewResendResult(hashes))
	}
}

// URI: /api/v1/rawtx
// Method: GET
// Args:
//	txid: transaction ID hash
// Returns the hex-encoded byte serialization of a transaction.
// The transaction may be confirmed or unconfirmed.
func rawTxnHandler(gateway Gatewayer) http.HandlerFunc {
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

		txnHex, err := txn.Transaction.SerializeHex()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		wh.SendJSONOr500(logger, w, txnHex)
	}
}

// VerifyTransactionRequest represents the data struct of the request for /api/v2/transaction/verify
type VerifyTransactionRequest struct {
	Unsigned           bool   `json:"unsigned"`
	EncodedTransaction string `json:"encoded_transaction"`
}

// VerifyTransactionResponse the response data struct for /api/v2/transaction/verify
type VerifyTransactionResponse struct {
	Unsigned    bool               `json:"unsigned"`
	Confirmed   bool               `json:"confirmed"`
	Transaction CreatedTransaction `json:"transaction"`
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

		var req VerifyTransactionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if req.EncodedTransaction == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "encoded_transaction is required")
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := decodeTxn(req.EncodedTransaction)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("decode transaction failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		signed := visor.TxnSigned
		if req.Unsigned {
			signed = visor.TxnUnsigned
		}

		var resp HTTPResponse
		inputs, isTxnConfirmed, err := gateway.VerifyTxnVerbose(txn, signed)
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

		verifyTxnResp := VerifyTransactionResponse{
			Confirmed: isTxnConfirmed,
			Unsigned:  !txn.IsFullySigned(),
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

	txn, err = coin.DeserializeTransaction(b)
	if err != nil {
		return nil, err
	}

	return &txn, nil
}

// newCreatedTransactionFuzzy creates a CreatedTransaction but accommodates possibly invalid txn input
func newCreatedTransactionFuzzy(txn *coin.Transaction, inputs []visor.TransactionInput) (*CreatedTransaction, error) {
	if len(txn.In) != len(inputs) && len(inputs) != 0 {
		return nil, errors.New("len(txn.In) != len(inputs)")
	}

	var outputHours uint64
	var feeInvalid bool
	for _, o := range txn.Out {
		var err error
		outputHours, err = mathutil.AddUint64(outputHours, o.Hours)
		if err != nil {
			feeInvalid = true
		}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = mathutil.AddUint64(inputHours, i.CalculatedHours)
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

	txID := txn.Hash()
	out := make([]CreatedTransactionOutput, len(txn.Out))
	for i, o := range txn.Out {
		co, err := NewCreatedTransactionOutput(o, txID)
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
		TxID:      txID.Hex(),
		InnerHash: txn.InnerHash.Hex(),
		Fee:       fmt.Sprint(fee),

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}
