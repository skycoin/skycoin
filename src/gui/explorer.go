package gui


import (
	"net/http"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/historydb"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)


func RegisterExploerHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// get set of pending transactions
	mux.HandleFunc("/explorer/address", getTransactionsForAddress(gateway))
}


func getTransactionsForAddress(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			wh.Error405(w, "")
			return
		}
		addr := r.FormValue("address")
		if addr == "" {
			wh.Error400(w, "address is empty")
			return
		}

		cipherAddr, err := cipher.DecodeBase58Address(addr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		uxs, err := gateway.GetAddressUxOuts(cipherAddr)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		resTxs := make([]visor.ReadableAddressTransaction, len(uxs))

		for i, ux := range uxs {
			sourceTxnNumber,err := cipher.SHA256FromHex(ux.Out.Body.SrcTransaction.Hex())
			if err!=nil{
				wh.Error400(w, "Transaction id is not good")
				return
			}
			sourceTransaction, err := gateway.V.GetTransaction(sourceTxnNumber)
			in := make([]visor.ReadableTransactionInput, len(sourceTransaction.Txn.In))
			for i, _ := range sourceTransaction.Txn.In {
				var uxout *historydb.UxOut
				var err error
				c := make(chan struct{})
				id, err := cipher.SHA256FromHex(sourceTransaction.Txn.In[i].Hex())
				if err != nil {
					wh.Error400(w, err.Error())
					return
				}
				gateway.Requests <- func() {
					uxout, err = gateway.V.GetUxOutByID(id)
					c <- struct{}{}
				}
				<-c
				in[i] = visor.NewReadableTransactionInput(sourceTransaction.Txn.In[i].Hex(), uxout.Out.Body.Address.String())
			}

			resTxs[i] = visor.NewReadableAddressTransaction(sourceTransaction, in);
		}
		wh.SendOr404(w, &resTxs)
	}
}

