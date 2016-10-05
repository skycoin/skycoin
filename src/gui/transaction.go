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

func RegisterTxHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get set of pending transactions
	mux.HandleFunc("/pendingTxs", getPendingTxs(gateway))
	// get latest confirmed transactions
	mux.HandleFunc("/lastTxs", getLastTxs(gateway))
	// get txn by txid
	mux.HandleFunc("/transaction", getTransactionByID(gateway))
	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
	// get raw tx by txid.
	mux.HandleFunc("/rawtx", getRawTx(gateway))
}

// Returns pending transactions
func getPendingTxs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
			return
		}

		V := gateway.V
		ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(V.Unconfirmed.Txns))
		for _, unconfirmedTxn := range V.Unconfirmed.Txns {
			readable := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
			ret = append(ret, &readable)
		}

		wh.SendOr404(w, &ret)
	}
}

// getLastTxs get the last confirmed txs.
func getLastTxs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
			return
		}
		txs, err := gateway.V.GetLastTxs()
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}

		resTxs := make([]visor.TransactionResult, len(txs))
		for i, tx := range txs {
			head := gateway.V.GetHeadBlock()
			height := head.Seq() - tx.BlockSeq + 1
			resTxs[i] = visor.TransactionResult{
				Transaction: visor.NewReadableTransaction(&tx.Tx),
				Status:      visor.NewConfirmedTransactionStatus(height),
			}
		}

		wh.SendOr404(w, &resTxs)
	}
}

func getTransactionByID(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		tx, err := gate.V.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		resTx := visor.TransactionResult{
			Transaction: visor.NewReadableTransaction(&tx.Txn),
			Status:      tx.Status,
		}
		wh.SendOr404(w, &resTx)
	}
}

//Implement
func injectTransaction(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			wh.Error405(w, "")
			return
		}
		// get the rawtransaction
		v := struct {
			Rawtx []byte `json:"rawtx"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			logger.Error("bad request: %v", err)
			wh.Error400(w, err.Error())
			return
		}

		txn := coin.TransactionDeserialize(v.Rawtx)
		if err := visor.VerifyTransactionFee(gateway.D.Visor.Visor.Blockchain, &txn); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		t, err := gateway.D.Visor.InjectTransaction(txn, gateway.D.Pool)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("inject tx failed:%v", err))
			return
		}

		wh.SendOr404(w, t.Hash().Hex())
	}
}

func getRawTx(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
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

		tx, err := gate.V.GetTransaction(h)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		d := tx.Txn.Serialize()
		wh.SendOr404(w, hex.EncodeToString(d))
		return
	}
}
