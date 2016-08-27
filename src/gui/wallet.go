// Wallet-related information for the GUI
package gui

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
)

//var Wallets wallet.Wallets

/*
REFACTOR
*/

/*
This section is redundant
- after moving the wallets out of visor and daemon, wallet should be in the wallet module
- there is no need for multiple wallets in the same application
*/
//type WalletRPC struct{}

type WalletRPC struct {
	Wallets         wallet.Wallets
	WalletDirectory string
}

//use a global for now
var Wg *WalletRPC

func InitWalletRPC(walletDir string) {
	Wg = NewWalletRPC(walletDir)
}

func NewWalletRPC(walletDir string) *WalletRPC {
	rpc := &WalletRPC{}

	if err := os.MkdirAll(walletDir, os.FileMode(0700)); err != nil {
		log.Panicf("Failed to create wallet directory %s: %v", walletDir, err)
	}

	rpc.WalletDirectory = walletDir

	w, err := wallet.LoadWallets(rpc.WalletDirectory)
	if err != nil {
		log.Panicf("Failed to load all wallets: %v", err)
	}
	rpc.Wallets = w

	if len(rpc.Wallets) == 0 {
		rpc.Wallets.Add(wallet.NewWallet("")) //deterministic
		errs := rpc.Wallets.Save(rpc.WalletDirectory)
		if len(errs) != 0 {
			log.Panicf("Failed to save wallets to %s: %v", rpc.WalletDirectory, errs)
		}
	}

	return rpc
}

func (self *WalletRPC) ReloadWallets() error {
	wallets, err := wallet.LoadWallets(self.WalletDirectory)
	if err != nil {
		return err
	}
	self.Wallets = wallets
	return nil
}

func (self *WalletRPC) SaveWallet(walletID wallet.WalletID) error {
	w := self.Wallets.Get(walletID)
	if w == nil {
		return fmt.Errorf("Unknown wallet %s", walletID)
	}
	return w.Save(self.WalletDirectory)
}

func (self *WalletRPC) SaveWallets() map[wallet.WalletID]error {
	return self.Wallets.Save(self.WalletDirectory)
}

func (self *WalletRPC) CreateWallet(seed string) wallet.Wallet {
	w := wallet.NewWallet(seed)
	self.Wallets.Add(w)
	return w
}

func (self *WalletRPC) GetWalletsReadable() []*wallet.ReadableWallet {
	return self.Wallets.ToReadable()
}

func (self *WalletRPC) GetWalletReadable(walletID wallet.WalletID) *wallet.ReadableWallet {
	w := self.Wallets.Get(walletID)
	if w == nil {
		return nil
	} else {
		return wallet.NewReadableWallet(*w)
	}
}

func (self *WalletRPC) GetWallet(walletID wallet.WalletID) *wallet.Wallet {
	w := self.Wallets.Get(walletID)
	if w == nil {
		return nil
	} else {
		return w
	}
}

//modify to return error
// NOT WORKING
// actually uses visor
func (self *WalletRPC) GetWalletBalance(v *visor.Visor,
	walletID wallet.WalletID) (wallet.BalancePair, error) {

	wlt := self.Wallets.Get(walletID)
	if wlt == nil {
		log.Printf("GetWalletBalance: ID NOT FOUND: id= '%s'", walletID)
		return wallet.BalancePair{}, errors.New("Id not found")
	}
	auxs := v.Blockchain.Unspent.AllForAddresses(wlt.GetAddresses())
	puxs := v.Unconfirmed.SpendsForAddresses(&v.Blockchain.Unspent,
		wlt.GetAddressSet())

	coins1, hours1 := v.AddressBalance(auxs)
	coins2, hours2 := v.AddressBalance(auxs.Sub(puxs))

	confirmed := wallet.Balance{coins1, hours1}
	predicted := wallet.Balance{coins2, hours2}

	return wallet.BalancePair{confirmed, predicted}, nil
}

/*
Checks if the wallet has pending, unconfirmed transactions
- do not allow any transactions if there are pending
*/
//Check if any of the outputs are spent
func (self *WalletRPC) HasUnconfirmedTransactions(v *visor.Visor,
	wallet *wallet.Wallet) bool {

	if wallet == nil {
		log.Panic("Wallet does not exist")
	}

	auxs := v.Blockchain.Unspent.AllForAddresses(wallet.GetAddresses())
	puxs := v.Unconfirmed.SpendsForAddresses(&v.Blockchain.Unspent,
		wallet.GetAddressSet())

	_ = auxs
	_ = puxs

	//no transactions
	if len(puxs) == 0 {
		return true
	}

	return false

}

type SpendResult struct {
	Balance     wallet.BalancePair        `json:"balance"`
	Transaction visor.ReadableTransaction `json:"txn"`
	Error       string                    `json:"error"`
}

// TODO
// - split send into
// -- get addresses
// -- get unspent outputs
// -- construct transaction
// -- sign transaction
// -- inject transaction
func Spend(d *daemon.Daemon, v *daemon.Visor, wrpc *WalletRPC,
	walletID wallet.WalletID, amt wallet.Balance, fee uint64,
	dest cipher.Address) *SpendResult {

	txn, err := Spend2(v.Visor, wrpc, walletID, amt, fee, dest)
	errString := ""
	if err != nil {
		errString = err.Error()
		logger.Error("Failed to make a spend: %v", err)
	}

	b, _ := wrpc.GetWalletBalance(v.Visor, walletID)

	if err != nil {
		log.Printf("transaction creation failed: %v", err)
	} else {
		log.Printf("Spend: \ntx= \n %s \n", visor.TransactionToJSON(txn))
	}

	v.Visor.InjectTxn(txn)
	//could call daemon, inject transaction
	//func (self *Visor) InjectTransaction(txn coin.Transaction, pool *Pool) (coin.Transaction, error) {
	//func (self *Visor) BroadcastTransaction(t coin.Transaction, pool *Pool)

	v.BroadcastTransaction(txn, d.Pool)

	return &SpendResult{
		Balance:     b,
		Transaction: visor.NewReadableTransaction(&txn),
		Error:       errString,
	}
}

// Creates a transaction spending amt with additional fee.  Fee is in addition
// to the base required fee given amt.Hours.
// TODO
// - pull in outputs from blockchain from wallet
// - create transaction here
// - sign transction and return
func Spend2(self *visor.Visor, wrpc *WalletRPC, walletID wallet.WalletID, amt wallet.Balance,
	fee uint64, dest cipher.Address) (coin.Transaction, error) {

	wallet := wrpc.Wallets.Get(walletID)
	if wallet == nil {
		return coin.Transaction{}, fmt.Errorf("Unknown wallet %v", walletID)
	}
	//pull in outputs and do this here
	//FIX
	tx, err := visor.CreateSpendingTransaction(*wallet, self.Unconfirmed,
		&self.Blockchain.Unspent, self.Blockchain.Time(), amt, dest)
	if err != nil {
		return tx, err
	}

	if err := tx.Verify(); err != nil {
		log.Panicf("Invalid transaction, %v", err)
	}

	if err := visor.VerifyTransactionFee(self.Blockchain, &tx); err != nil {
		log.Panicf("Created invalid spending txn: visor fail, %v", err)
	}
	if err := self.Blockchain.VerifyTransaction(tx); err != nil {
		log.Panicf("Created invalid spending txn: blockchain fail, %v", err)
	}
	return tx, err
}

/*
REFACTOR
*/

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
func walletBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		//addr := r.FormValue("addr")

		//r.ParseForm()
		//r.ParseMultipartForm()
		//log.Println(r.Form)

		//r.URL.String()
		r.ParseForm()

		b, err := Wg.GetWalletBalance(gateway.D.Visor.Visor, wallet.WalletID(id))

		if err != nil {
			_ = err
		}
		//log.Printf("%v, %v, %v \n", r.URL.String(), r.RequestURI, r.Form)
		SendOr404(w, b)
	}
}

func getBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			addrsParam := r.URL.Query().Get("addrs")
			addrsStr := strings.Split(addrsParam, ",")
			addrs := make([]cipher.Address, len(addrsStr))
			addrSet := make(map[cipher.Address]byte)
			for i, addr := range addrsStr {
				addrs[i] = cipher.MustDecodeBase58Address(addr)
				addrSet[addrs[i]] = byte(1)
			}

			v := gateway.D.Visor.Visor
			auxs := v.Blockchain.Unspent.AllForAddresses(addrs)
			puxs := v.Unconfirmed.SpendsForAddresses(&v.Blockchain.Unspent, addrSet)

			coins1, hours1 := v.AddressBalance(auxs)
			coins2, hours2 := v.AddressBalance(auxs.Sub(puxs))

			confirmed := wallet.Balance{coins1, hours1}
			predicted := wallet.Balance{coins2, hours2}
			bal := struct {
				Confirmed wallet.Balance `json:"confirmed"`
				Predicted wallet.Balance `json:"predicted"`
			}{
				Confirmed: confirmed,
				Predicted: predicted,
			}

			SendOr404(w, bal)
		}
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
func walletSpendHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//log.Printf("Spend1")

		if r.FormValue("id") == "" {
			Error400(w, "Missing wallet_id")
			return
		}

		walletId := wallet.WalletID(r.FormValue("id"))
		if walletId == "" {
			Error400(w, "Invalid Wallet Id")
			return
		}
		sdst := r.FormValue("dst")
		if sdst == "" {
			Error400(w, "Missing destination address \"dst\"")
			return
		}
		dst, err := cipher.DecodeBase58Address(sdst)
		if err != nil {
			//Error400(w, "Invalid destination address: %v", err)
			Error400(w, "Invalid destination address: %v", err.Error())
			return
		}

		//set fee automatically for now
		/*
			sfee := r.FormValue("fee")
			fee, err := strconv.ParseUint(sfee, 10, 64)
			if err != nil {
				Error400(w, "Invalid \"fee\" value")
				return
			}
		*/
		//var fee uint64 = 0

		scoins := r.FormValue("coins")
		//shours := r.FormValue("hours")
		coins, err := strconv.ParseUint(scoins, 10, 64)
		if err != nil {
			Error400(w, "Invalid \"coins\" value")
			return
		}

		var hours uint64 = 0
		var fee uint64 = 0 //doesnt work/do anything right now

		//MOVE THIS INTO HERE
		ret := Spend(gateway.D, gateway.D.Visor, Wg, walletId, wallet.NewBalance(coins, hours), fee, dst)

		if ret.Error != "" {
			Error400(w, "Spend Failed: %s", ret.Error)
		}
		SendOr404(w, ret)
	}
}

// Create a wallet if no ID provided.  Otherwise update an existing wallet.
// Name is set by creation date
func walletCreate(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("API request made to create a wallet")
		seed := r.FormValue("seed")
		w1 := Wg.CreateWallet(seed) //use seed!
		iw := wallet.NewReadableWallet(w1)
		if iw != nil {
			if err := Wg.SaveWallet(w1.GetID()); err != nil {
				m := "Failed to save wallet after renaming: %v"
				logger.Critical(m, err)
			}
		}
		SendOr500(w, iw)
	}
}

//all this does is update the wallet name
// Does Nothing
func walletUpdate(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update wallet
		id := wallet.WalletID(r.FormValue("id"))
		//name := r.FormValue("name")
		w1 := Wg.GetWallet(id)
		if w1 != nil {
			if err := Wg.SaveWallet(w1.GetID()); err != nil {
				m := "Failed to save wallet after renaming: %v"
				logger.Critical(m, err)
			}
		}
		iw := wallet.NewReadableWallet(*w1)
		SendOr404(w, iw)
	}
}

// Returns a wallet by ID if GET.  Creates or updates a wallet if POST.
func walletGet(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ret := Wg.GetWallet(wallet.WalletID(r.FormValue("id")))
			SendOr404(w, ret)
		}
	}
}

// Returns JSON of pending transactions for user's wallet
func walletTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			wallet := Wg.GetWallet(wallet.WalletID(r.FormValue("id")))
			addresses := wallet.GetAddresses()
			ret := gateway.Visor.GetWalletTransactions(gateway.V, addresses)

			SendOr404(w, ret)
		}
	}
}

// Returns all loaded wallets
func walletsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ret := wallet.Wallets.ToPublicReadable()
		ret := Wg.GetWalletsReadable()
		SendOr404(w, ret)
	}
}

// Saves all loaded wallets
func walletsSaveHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errs := Wg.SaveWallets() // (map[wallet.WalletID]error)
		if len(errs) != 0 {
			err := ""
			for id, e := range errs {
				err += string(id) + ": " + e.Error()
			}
			Error500(w, err)
		}
	}
}

// Loads/unloads wallets from the wallet directory
func walletsReloadHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Wg.ReloadWallets()
		if err != nil {
			Error500(w, err.(error).Error())
		}
	}
}

// Returns the outputs for all addresses
func getOutputsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			uxouts := gateway.Visor.GetUnspentOutputReadables(gateway.V)
			rawaddrs := r.URL.Query().Get("addrs")
			if rawaddrs == "" {
				SendOr404(w, uxouts)
				return
			}

			// split the addrs,
			addrs := strings.Split(rawaddrs, ",")
			addrMap := make(map[string]bool)
			for _, addr := range addrs {
				addrMap[addr] = true
			}

			ret := []visor.ReadableOutput{}
			for _, u := range uxouts {
				if _, ok := addrMap[u.Address]; ok {
					ret = append(ret, u)
				}
			}
			SendOr404(w, ret)
		}
	}
}

// Returns pending transactions
// TODO: FIX!!! Iterates all blocks since begining
// Gets list of transactions
// TODO: this will slow down exponentially as blockchain size increases
// TODO: split function for determining if transaction is confirmed, into another function
// TODO: only iterate, last 50 blocks, to determine if transaction is confirmed
// TODO: use transaction ID hash, not readable, to confirm transaction
func getTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
			blks := V.Blockchain.Blocks

			//only look at last 50 blocks, for checking if transaction is confirmed
			const max_history = 50
			x1 := len(blks)               // start
			x2 := len(blks) - max_history //stop
			if x2 < 0 {
				x2 = 0
			}
			blks = blks[x2:x1] //only look at last 50 blocks

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

func getTransactionHandler(gate *daemon.Gateway) http.HandlerFunc {
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
		SendOr404(w, visor.NewReadableTransaction(&tx.Txn))
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

func RegisterWalletHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	// Returns wallet info
	// GET Arguments:
	//      id - Wallet ID.

	//  Gets a wallet .  Will be assigned name if present.
	mux.HandleFunc("/wallet", walletGet(gateway))

	// POST/GET Arguments:
	//		seed [optional]
	//create new wallet
	mux.HandleFunc("/wallet/create", walletCreate(gateway))

	//update an existing wallet
	//does nothing
	mux.HandleFunc("/wallet/update", walletUpdate(gateway))

	// Returns the confirmed and predicted balance for a specific wallet.
	// The predicted balance is the confirmed balance minus any pending
	// spent amount.
	// GET arguments:
	//      id: Wallet ID
	mux.HandleFunc("/wallet/balance", walletBalanceHandler(gateway))

	// Sends coins&hours to another address.
	// POST arguments:
	//  id: Wallet ID
	//  coins: Number of coins to spend
	//  hours: Number of hours to spends
	//  fee: Number of hours to use as fee, on top of the default fee.
	//  Returns total amount spent if successful, otherwise error describing
	//  failure status.
	mux.HandleFunc("/wallet/spend", walletSpendHandler(gateway))

	// GET Arguments:
	//		id: Wallet ID
	// Returns all pending transanction for all addresses by selected Wallet
	mux.HandleFunc("/wallet/transactions", walletTransactionsHandler(gateway))

	// Returns all loaded wallets
	mux.HandleFunc("/wallets", walletsHandler(gateway))
	// Saves all wallets to disk. Returns nothing if it works. Otherwise returns
	// 500 status with error message.

	mux.HandleFunc("/wallets/save", walletsSaveHandler(gateway))
	// Rescans the wallet directory and loads/unloads wallets based on which
	// files are present. Returns nothing if it works. Otherwise returns
	// 500 status with error message.
	mux.HandleFunc("/wallets/reload", walletsReloadHandler(gateway))

	//get set of unspent outputs
	mux.HandleFunc("/outputs", getOutputsHandler(gateway))

	// get balance of addresses
	mux.HandleFunc("/balance", getBalanceHandler(gateway))

	//get set of pending transaction
	mux.HandleFunc("/transactions", getTransactionsHandler(gateway))

	// get txn of address.
	mux.HandleFunc("/transaction", getTransactionHandler(gateway))

	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
}
