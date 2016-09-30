package gui

import (
	"encoding/hex"
	"encoding/json"
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
	// get txn by txid.
	mux.HandleFunc("/transaction", getTransactionByID(gateway))
	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
	// get raw tx by txid.
	mux.HandleFunc("/rawtx", getRawTx(gateway))
}

// Returns pending transactions
func getPendingTxs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		V := gateway.V
		ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(V.Unconfirmed.Txns))
		for _, unconfirmedTxn := range V.Unconfirmed.Txns {
			readable := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
			ret = append(ret, &readable)
		}

		var rlt struct {
			Success bool                            `json:"success"`
			Txns    []*visor.ReadableUnconfirmedTxn `json:"unconfirm_txs"`
		} {
			Success: true,
			Txns: ret,
		}

		wh.SendOr404(w, &rlt)
	}
}

// getLastTxs get the last confirmed txs.
func getLastTxs(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rlt struct {
			Success bool                      `json:"success"`
			Reason  string                    `json:"reason, omitempty"`
			Txs     []visor.TransactionResult `json:"transactions, omitempty"`
		}
		for {
			txs, err := gateway.V.GetLastTxs()
			if err != nil {
				rlt.Reason = err.Error()
				break
			}
			rlt.Success = true
			rlt.Txs = make([]visor.TransactionResult, len(txs))
			for i, tx := range txs {
				head := gateway.V.GetHeadBlock()
				height := head.Seq() - tx.BlockSeq + 1
				rlt.Txs[i] = visor.TransactionResult{
					Transaction: visor.NewReadableTransaction(&tx.Tx),
					Status:      visor.NewConfirmedTransactionStatus(height),
				}
			}
			break
		}
		wh.SendOr404(w, &rlt)
	}
}

func getTransactionByID(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rlt struct {
			Success bool                    `json:"success"`
			Reason  string                  `json:"reason,omitempty"`
			Tx      visor.TransactionResult `json:"transaction, omitempty"`
		}
		for {
			txid := r.FormValue("txid")
			if txid == "" {
				rlt.Reason = "txid is empty"
				break
			}

			h, err := cipher.SHA256FromHex(txid)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}

			tx, err := gate.V.GetTransaction(h)
			if err != nil {
				rlt.Reason = err.Error()
				break
			}

			rlt.Success = true
			rlt.Tx = visor.TransactionResult{
				Transaction: visor.NewReadableTransaction(&tx.Txn),
				Status:      tx.Status,
			}
			break
		}
		wh.SendOr404(w, &rlt)
	}
}


//Implement
func injectTransaction(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the rawtransaction
		v := struct {
			Rawtx []byte `json:"rawtx"`
		}{}

		rlt := struct {
			Success bool   `json:"success"`
			Reason  string `json:"reason"`
			Txid    string `json:"txid"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			logger.Error("bad request: %v", err)
			rlt.Reason = "bad request"
			wh.SendOr404(w, rlt)
			return
		}

		txn := coin.TransactionDeserialize(v.Rawtx)

		if err := visor.VerifyTransactionFee(gateway.D.Visor.Visor.Blockchain, &txn); err != nil {
			rlt.Reason = err.Error()
			wh.SendOr404(w, rlt)
			return
		}

		t, err := gateway.D.Visor.InjectTransaction(txn, gateway.D.Pool)
		if err != nil {
			logger.Error("inject tx failed:%v", err)
			rlt.Reason = "inject tx failed"
			wh.SendOr404(w, rlt)
			return
		}

		rlt.Success = true
		rlt.Txid = t.Hash().Hex()

		//ret := gateway.Visor.GetUnspentOutputReadables(gateway.V)
		wh.SendOr404(w, rlt)
	}
}

func getRawTx(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txid := r.FormValue("txid")
		if txid == "" {
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			return
		}
		tx := gate.V.GetTransaction(h)
		d := tx.Txn.Serialize()
		res := struct {
			Rawtx string `json:"rawtx"`
		}{
			hex.EncodeToString(d),
		}

		wh.SendOr404(w, res)
	}
}
