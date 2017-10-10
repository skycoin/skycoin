package gui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// RegisterTxHandlers registers transaction handlers
func RegisterTxHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get set of pending transactions
	mux.HandleFunc("/pendingTxs", getPendingTxs(gateway))
	// get latest confirmed transactions
	mux.HandleFunc("/lastTxs", getLastTxs(gateway))
	// get txn by txid
	mux.HandleFunc("/transaction", getTransactionByID(gateway))
	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
	mux.HandleFunc("/resendUnconfirmedTxns", resendUnconfirmedTxns(gateway))
	// get raw tx by txid.
	mux.HandleFunc("/rawtx", getRawTx(gateway))
}

// Returns pending transactions
func getPendingTxs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		txns := gateway.GetAllUnconfirmedTxns()
		ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(txns))
		for _, unconfirmedTxn := range txns {
			readable, err := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
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
func getLastTxs(gateway *daemon.Gateway) http.HandlerFunc {
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
			rbTx, err := visor.NewReadableTransaction(tx)
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

func getTransactionByID(gate *daemon.Gateway) http.HandlerFunc {
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

		rbTx, err := visor.NewReadableTransaction(tx)
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

//Implement
func injectTransaction(gateway *daemon.Gateway) http.HandlerFunc {
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

		txn := coin.TransactionDeserialize(b)
		if err := gateway.InjectTransaction(txn); err != nil {
			wh.Error400(w, fmt.Sprintf("inject tx failed:%v", err))
			return
		}

		wh.SendOr404(w, txn.Hash().Hex())
	}
}

func resendUnconfirmedTxns(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		rlt := gate.ResendUnconfirmedTxns()
		v, _ := json.MarshalIndent(rlt, "", "    ")
		fmt.Println(v)
		wh.SendOr404(w, rlt)
		return
	}
}

func getRawTx(gate *daemon.Gateway) http.HandlerFunc {
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
