package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor"
)

// Returns pending transactions api/v2
func getPendingTxnsV2(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		txns, err := gateway.GetPendingTxnsV2()
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("GetPendingTxnsV2 failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}
		var resp HTTPResponse
		resp.Data = txns
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
			resp := NewHTTPErrorResponse(http.StatusBadRequest, fmt.Sprintf("SHA256FromHex failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		txn, err := gateway.GetTransactionV2(h)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("GetTransactionV2 failed : %v", err))
			writeHTTPResponse(w, resp)
			return
		}
		if txn == nil {
			resp := NewHTTPErrorResponse(http.StatusNotFound, fmt.Sprintf("Transaction not found : %v", txid))
			writeHTTPResponse(w, resp)
			return
		}
		var resp HTTPResponse
		resp.Data = txn
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
		txnRltsV2, err := gateway.GetTransactionsV2(flts...)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, fmt.Sprintf("gateway.GetTransactions failed: %v", err))
			writeHTTPResponse(w, resp)
			return
		}

		var resp HTTPResponse
		resp.Data = txnRltsV2
		writeHTTPResponse(w, resp)
	}
}
