package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
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
			wh.Error405(w)
			return
		}

		txns, err := gateway.GetAllUnconfirmedTxns()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		ret := make([]*visor.ReadableUnconfirmedTxnV2, 0, len(txns))
		for _, unconfirmedTxn := range txns {
			readable, err := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			readableV2, err := NewReadableUnconfirmedTxnV2(gateway, readable)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			ret = append(ret, readableV2)
		}

		wh.SendJSONOr500(logger, w, &ret)
	}
}

func getTransactionByIDV2(gateway Gatewayer) http.HandlerFunc {
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

		rbTxn, err := visor.NewReadableTransaction(txn)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		rbTxnV2, err := NewReadableTransactionV2(gateway, rbTxn)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		resTxn := TransactionResultV2{
			Transaction: *rbTxnV2,
			Status:      txn.Status,
			Time:        txn.Time,
		}
		wh.SendJSONOr500(logger, w, &resTxn)
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

		txnRltsV2 := make([]TransactionResultV2, 0, len(txnRlts.Txns))

		for _, txn := range txnRlts.Txns {
			rbTxnV2, err := NewReadableTransactionV2(gateway, &txn.Transaction)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
			trV2 := TransactionResultV2{
				Transaction: *rbTxnV2,
				Status:      txn.Status,
				Time:        txn.Time,
			}
			txnRltsV2 = append(txnRltsV2, trV2)
		}

		wh.SendJSONOr500(logger, w, txnRltsV2)
	}
}
