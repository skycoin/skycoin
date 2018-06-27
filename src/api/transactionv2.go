package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

// TransactionResultV2 represents transaction result api/v2
type TransactionResultV2 struct {
	Status      visor.TransactionStatus     `json:"status"`
	Time        uint64                      `json:"time"`
	Transaction visor.ReadableTransactionV2 `json:"txn"`
}

// TransactionResultsV2 array of transaction results api/v2
type TransactionResultsV2 struct {
	Txns []TransactionResultV2 `json:"txns"`
}

// Returns pending transactions api/v2
func getPendingTxnsV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		txns, err := gateway.GetAllUnconfirmedTxns()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("GetAllUnconfirmedTxns failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		ret := make([]*visor.ReadableUnconfirmedTxnV2, 0, len(txns))
		for _, unconfirmedTxn := range txns {
			readable, err := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableUnconfirmedTxn failed: %v", err))
				writeHTTPResponse(w, resp)
				return
			}
			readableV2, err := NewReadableUnconfirmedTxnV2(gateway, readable)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableUnconfirmedTxnV2 failed: %v", err))
				writeHTTPResponse(w, resp)
				return
			}
			ret = append(ret, readableV2)
		}
		var resp HTTPResponse
		resp.Data = ret
		writeHTTPResponse(w, resp)
	}
}

func getTransactionByIDV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
		}
		txid := r.FormValue("txid")
		if txid == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "txid is empty")
			writeHTTPResponse(w, resp)
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("SHA256FromHex failed : %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := gateway.GetTransaction(h)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("GetTransaction failed : %v", err))
			writeHTTPResponse(w, resp)
			return
		}
		if txn == nil {
			resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("Transaction not found : %v", txid))
			writeHTTPResponse(w, resp)
			return
		}

		rbTxn, err := visor.NewReadableTransaction(txn)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableTransaction failed : %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		rbTxnV2, err := NewReadableTransactionV2(gateway, rbTxn)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableTransactionV2 failed : %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		resTxn := TransactionResultV2{
			Transaction: *rbTxnV2,
			Status:      txn.Status,
			Time:        txn.Time,
		}
		var resp HTTPResponse
		resp.Data = resTxn
		writeHTTPResponse(w, resp)
	}
}

// Returns transactions that match the filters.
// Method: GET
// URI: /api/v2/transactions
// Args:
//     addrs: Comma seperated addresses [optional, returns all transactions if no address provided]
//     confirmed: Whether the transactions should be confirmed [optional, must be 0 or 1; if not provided, returns all]
func getTransactionsV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		// Gets 'addrs' parameter value
		addrs, err := parseAddressesFromStr(r.FormValue("addrs"))
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("parse parameter: 'addrs' failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		// Initialize transaction filters
		flts := []visor.TxFilter{visor.AddrsFilter(addrs)}

		// Gets the 'confirmed' parameter value
		confirmedStr := r.FormValue("confirmed")
		if confirmedStr != "" {
			confirmed, err := strconv.ParseBool(confirmedStr)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("invalid 'confirmed' value: %v", err))
				writeHTTPResponse(w, resp)
				return
			}

			flts = append(flts, visor.ConfirmedTxFilter(confirmed))
		}

		// Gets transactions
		txns, err := gateway.GetTransactions(flts...)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("gateway.GetTransactions failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		// Converts visor.Transaction to daemon.TransactionResult
		txnRlts, err := daemon.NewTransactionResults(txns)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("gateway.NewTransactionResults failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		txnRltsV2 := make([]TransactionResultV2, 0, len(txnRlts.Txns))

		for _, txn := range txnRlts.Txns {
			rbTxnV2, err := NewReadableTransactionV2(gateway, &txn.Transaction)
			if err != nil {
				resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("NewReadableTransactionV2 failed: %v", err))
				writeHTTPResponse(w, resp)
				return
			}
			trV2 := TransactionResultV2{
				Transaction: *rbTxnV2,
				Status:      txn.Status,
				Time:        txn.Time,
			}
			txnRltsV2 = append(txnRltsV2, trV2)
		}

		var resp HTTPResponse
		resp.Data = txnRltsV2
		writeHTTPResponse(w, resp)
	}
}
