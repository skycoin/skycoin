package gui

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
)

func RegisterTxHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	//get set of pending transaction
	mux.HandleFunc("/transactions", getTransactions(gateway))
	// get txn by txid.
	mux.HandleFunc("/transaction", getTransactionByID(gateway))
	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
	// get raw tx by txid.
	mux.HandleFunc("/rawtx", getRawTx(gateway))
}

// Returns pending transactions
// TODO: FIX!!! Iterates all blocks since begining
// Gets list of transactions
// TODO: this will slow down exponentially as blockchain size increases
// TODO: split function for determining if transaction is confirmed, into another function
// TODO: only iterate, last 50 blocks, to determine if transaction is confirmed
// TODO: use transaction ID hash, not readable, to confirm transaction
func getTransactions(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		V := gateway.V
		isConfirmed := r.URL.Query().Get("confirm")

		//default case
		if isConfirmed != "1" {
			ret := make([]*visor.ReadableUnconfirmedTxn, 0, len(V.Unconfirmed.Txns))
			for _, unconfirmedTxn := range V.Unconfirmed.Txns {
				readable := visor.NewReadableUnconfirmedTxn(&unconfirmedTxn)
				ret = append(ret, &readable)
			}
			SendOr404(w, ret)
		}

		//WARNING: TODO: This iterates all blocks and all transactions
		//TODO: need way to determine if transaction is "confirmed", without iterating all blocks
		if isConfirmed == "1" {
			// blks := V.Blockchain.Blocks

			//only look at last 50 blocks, for checking if transaction is confirmed
			// const max_history = 50
			// x1 := len(blks)               // start
			// x2 := len(blks) - max_history //stop
			// if x2 < 0 {
			// 	x2 = 0
			// }
			blks := V.Blockchain.GetLatestBlocks(50)
			// blks = blks[x2:x1] //only look at last 50 blocks

			totalTxns := []coin.Transaction{}
			//WARNING: Iterates all blocks, since start
			//TODO: use transaction hash, not input/output
			for _, b := range blks {
				totalTxns = append(totalTxns, b.Body.Transactions...)
			}

			rdTxns := make([]visor.ReadableTransaction, len(totalTxns))
			for i, txn := range totalTxns {
				rdTxns[i] = visor.NewReadableTransaction(&txn)
			}

			rltTxns := []visor.ReadableTransaction{}
			// addr := r.URL.Query().Get("addr")
			input := r.URL.Query().Get("input")
			output := r.URL.Query().Get("output")

			if input != "" {
				uxids := getUxidsOfAddr(input, rdTxns)
				for _, uxid := range uxids {
					for _, txn := range rdTxns {
						for _, in := range txn.In {
							if in == uxid {
								rltTxns = append(rltTxns, txn)
								break
							}
						}
					}
				}
			}

			if output != "" {
				outTxns := []visor.ReadableTransaction{}
				if input != "" {
					outTxns = rltTxns
				} else {
					outTxns = rdTxns
				}

				txs := []visor.ReadableTransaction{}
				for _, txn := range outTxns {
					for _, out := range txn.Out {
						if out.Address == output {
							txs = append(txs, txn)
							break
						}
					}
				}
				rltTxns = txs
			}
			SendOr404(w, rltTxns)
		}

	}
}

func getUxidsOfAddr(addr string, rdTxns []visor.ReadableTransaction) []string {
	uxids := []string{}
	for _, txn := range rdTxns {
		for _, out := range txn.Out {
			if out.Address == addr {
				uxids = append(uxids, out.Hash)
			}
		}
	}
	return uxids
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
			SendOr404(w, rlt)
			return
		}

		txn := coin.TransactionDeserialize(v.Rawtx)

		if err := visor.VerifyTransactionFee(gateway.D.Visor.Visor.Blockchain, &txn); err != nil {
			rlt.Reason = err.Error()
			SendOr404(w, rlt)
			return
		}

		t, err := gateway.D.Visor.InjectTransaction(txn, gateway.D.Pool)
		if err != nil {
			logger.Error("inject tx failed:%v", err)
			rlt.Reason = "inject tx failed"
			SendOr404(w, rlt)
			return
		}

		rlt.Success = true
		rlt.Txid = t.Hash().Hex()

		//ret := gateway.Visor.GetUnspentOutputReadables(gateway.V)
		SendOr404(w, rlt)
	}
}

func getTransactionByID(gate *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txid := r.URL.Query().Get("txid")
		if txid == "" {
			return
		}

		h, err := cipher.SHA256FromHex(txid)
		if err != nil {
			return
		}
		tx := gate.V.GetTransaction(h)
		rlt := visor.TransactionResult{
			Transaction: visor.NewReadableTransaction(&tx.Txn),
			Status:      tx.Status,
		}
		SendOr404(w, &rlt)
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

		SendOr404(w, res)
	}
}
