package gui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// Returns pending transactions
func getPendingTxs(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txns := gateway.GetAllUnconfirmedTxns()
		ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(txns))
		for _, unconfirmedTxn := range txns {
			txInputsData, err := getTransactionInputsData(&unconfirmedTxn.Txn, gateway)
			if err != nil {
				logger.Error("%v", err)
				wh.Error500(w)
				return
			}

			readable, err := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn, txInputsData)
			if err != nil {
				logger.Error("%v", err)
				wh.Error500(w)
				return
			}
			ret = append(ret, readable)
		}

		wh.SendOr404(w, &ret)
	}
}

// DEPRECATED: last txs can't recover from db when restart
// , and it's not used actually
func getLastTxs(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}
		txs, err := gateway.GetLastTxs()
		if err != nil {
			logger.Error("gateway.GetLastTxs failed: %v", err)
			wh.Error500(w)
			return
		}

		resTxs := make([]visor.TransactionResult, len(txs))
		for i, tx := range txs {
			txInputsData, err := getTransactionInputsData(&tx.Txn, gateway)
			if err != nil {
				logger.Error("%v", err)
				wh.Error500(w)
				return
			}

			rbTx, err := visor.NewReadableTransaction(tx, txInputsData)
			if err != nil {
				logger.Error("%v", err)
				wh.Error500(w)
				return
			}

			resTxs[i] = visor.TransactionResult{
				Transaction: *rbTx,
				Status:      tx.Status,
			}
		}

		wh.SendOr404(w, &resTxs)
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
			wh.Error404(w)
			return
		}

		txInputsData, err := getTransactionInputsData(&tx.Txn, gate)
		if err != nil {
			logger.Error("%v", err)
			wh.Error500(w)
			return
		}

		rbTx, err := visor.NewReadableTransaction(tx, txInputsData)
		if err != nil {
			logger.Error("%v", err)
			wh.Error500(w)
			return
		}

		resTx := visor.TransactionResult{
			Transaction: *rbTx,
			Status:      tx.Status,
		}
		wh.SendOr404(w, &resTx)
	}
}

// Returns transactions that match the filters.
// Method: GET
// URI: /transactions
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
			wh.Error400(w, fmt.Sprintf("parse parament: 'addrs' failed: %v", err))
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
			logger.Error("get transactions failed: %v", err)
			wh.Error500(w)
			return
		}

		// Converts visor.Transaction to visor.TransactionResult
		txRlts, err := NewTransactionResults(gateway, txns)
		if err != nil {
			logger.Error("Converts []visor.Transaction to visor.TransactionResults failed: %v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr404(w, txRlts.Txns)
	}
}

// Parses comma seperated addresses string into []cipher.Address,
func parseAddressesFromStr(addrStr string) ([]cipher.Address, error) {
	if addrStr == "" {
		return nil, nil
	}

	var addrs []cipher.Address
	for _, as := range strings.Split(addrStr, ",") {
		s := strings.TrimSpace(as)
		if s == "" {
			continue
		}

		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return nil, err
		}

		addrs = append(addrs, a)
	}

	return addrs, nil
}

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
			logger.Error("bad request: %v", err)
			wh.Error400(w, err.Error())
			return
		}

		b, err := hex.DecodeString(v.Rawtx)
		if err != nil {
			logger.Error("%v", err)
			wh.Error400(w, err.Error())
			return
		}

		txn, err := coin.TransactionDeserialize(b)
		if err != nil {
			logger.Error("%v", err)
			wh.Error400(w, err.Error())
			return
		}

		if err := gateway.InjectBroadcastTransaction(txn); err != nil {
			logger.Error("%v", err)
			wh.Error400(w, fmt.Sprintf("inject tx failed: %v", err))
			return
		}

		wh.SendOr404(w, txn.Hash().Hex())
	}
}

func resendUnconfirmedTxns(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		rlt := gateway.ResendUnconfirmedTxns()
		wh.SendOr404(w, rlt)
		return
	}
}

func getRawTx(gate Gatewayer) http.HandlerFunc {
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
			wh.Error404(w)
			return
		}

		d := tx.Txn.Serialize()
		wh.SendOr404(w, hex.EncodeToString(d))
		return
	}
}

// NewTransactionResults converts []Transaction to []TransactionResults
func NewTransactionResults(gateway Gatewayer, txs []visor.Transaction) (*visor.TransactionResults, error) {
	txRlts := make([]visor.TransactionResult, 0, len(txs))
	for _, tx := range txs {
		txInputsData, err := getTransactionInputsData(&tx.Txn, gateway)
		if err != nil {
			return nil, err
		}

		rbTx, err := visor.NewReadableTransaction(&tx, txInputsData)
		if err != nil {
			return nil, err
		}

		txRlts = append(txRlts, visor.TransactionResult{
			Transaction: *rbTx,
			Status:      tx.Status,
			Time:        tx.Time,
		})
	}

	return &visor.TransactionResults{
		Txns: txRlts,
	}, nil
}

// getTransactionInputsData returns the inputs data of a transaction
func getTransactionInputsData(tx *coin.Transaction, gw Gatewayer) ([]*historydb.UxOut, error) {
	txInputsData := make([]*historydb.UxOut, 0, len(tx.In))

	for _, in := range tx.In {

		uxout, err := getInputData(in, gw)
		if err != nil {
			return nil, err
		}

		txInputsData = append(txInputsData, uxout)
	}

	return txInputsData, nil
}
