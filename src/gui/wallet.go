// Wallet-related information for the GUI
package gui

import (
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
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
var Wg *WalletRPC = NewWalletRPC()

func NewWalletRPC() *WalletRPC {
	rpc := WalletRPC{}

	//wallet directory
	//cleanup, pass as parameter during init

	DataDirectory := util.InitDataDir("")
	rpc.WalletDirectory = filepath.Join(DataDirectory, "wallets/")
	logger.Debug("Wallet Directory= %v", rpc.WalletDirectory)
	util.InitDataDir(rpc.WalletDirectory)

	rpc.Wallets = wallet.Wallets{}

	//util.InitDataDir(".skycoin")
	//util.InitDataDir(".skycoin/wallets")

	//if rpc.WalletDirectory != "" {
	w, err := wallet.LoadWallets(rpc.WalletDirectory)
	if err != nil {
		log.Panicf("Failed to load all wallets: %v", err)
	}
	rpc.Wallets = w
	//}
	if len(rpc.Wallets) == 0 {
		rpc.Wallets.Add(wallet.NewWallet("")) //deterministic
		if rpc.WalletDirectory != "" {
			errs := rpc.Wallets.Save(rpc.WalletDirectory)
			if len(errs) != 0 {
				log.Panicf("Failed to save wallets: %v", errs)
			}
		}
	}

	return &rpc
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
	walletID wallet.WalletID) wallet.BalancePair {

	wlt := self.Wallets.Get(walletID)
	if wlt == nil {
		log.Printf("GetWalletBalance: ID NOT FOUND")
		return wallet.BalancePair{}
	}
	auxs := v.Blockchain.Unspent.AllForAddresses(wlt.GetAddresses())
	puxs := v.Unconfirmed.SpendsForAddresses(&v.Blockchain.Unspent,
		wlt.GetAddressSet())

	coins1, hours1 := v.AddressBalance(auxs)
	coins2, hours2 := v.AddressBalance(auxs.Sub(puxs))

	confirmed := wallet.Balance{coins1, hours1}
	predicted := wallet.Balance{coins2, hours2}

	return wallet.BalancePair{confirmed, predicted}
}

/*
REFACTOR
*/

// Returns the confirmed & predicted balance for a single address
// IMPLEMENT?
/*
func AddressBalance(self *visor.Visor, addr cipher.Address) wallet.BalancePair {
	auxs := self.Blockchain.Unspent.AllForAddress(addr)
	puxs := self.Unconfirmed.SpendsForAddress(&self.Blockchain.Unspent, addr)
	confirmed := self.balance(auxs)
	predicted := self.balance(auxs.Sub(puxs))
	return wallet.BalancePair{confirmed, predicted}
}
*/

// Returns the confirmed & predicted balance for a Wallet
/*
func (self *visor.Visor) WalletBalance(walletID wallet.WalletID) wallet.BalancePair {
	wlt := self.Wallets.Get(walletID)
	if wlt == nil {
		return wallet.BalancePair{}
	}
	auxs := self.Blockchain.Unspent.AllForAddresses(wlt.GetAddresses())
	puxs := self.Unconfirmed.SpendsForAddresses(&self.Blockchain.Unspent,
		wlt.GetAddressSet())
	confirmed := self.totalBalance(auxs)
	predicted := self.totalBalance(auxs.Sub(puxs))
	return wallet.BalancePair{confirmed, predicted}
}
*/

/*
// Return the total balance of all loaded wallets
func (self *visor.Visor) TotalBalance() wallet.BalancePair {
	b := wallet.BalancePair{}
	for _, w := range self.Wallets {
		c := self.WalletBalance(w.GetID())
		b.Confirmed = b.Confirmed.Add(c.Confirmed)
		b.Predicted = b.Confirmed.Add(c.Predicted)
	}
	return b
}

// Computes the total balance for a cipher.Address's coin.UxOuts
func (self *visor.Visor) balance(uxs coin.UxArray) wallet.Balance {
	prevTime := self.Blockchain.Time()
	b := wallet.NewBalance(0, 0)
	for _, ux := range uxs {
		b = b.Add(wallet.NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
	}
	return b
}

// Computes the total balance for cipher.Addresses and their coin.UxOuts
func (self *visor.Visor) totalBalance(auxs coin.AddressUxOuts) wallet.Balance {
	prevTime := self.Blockchain.Time()
	b := wallet.NewBalance(0, 0)
	for _, uxs := range auxs {
		for _, ux := range uxs {
			b = b.Add(wallet.NewBalance(ux.Body.Coins, ux.CoinHours(prevTime)))
		}
	}
	return b
}
*/

/*
REFACTOR
*/

/*
func Spend(self *daemon.Gateway, wrpc WalletRPC, walletID wallet.WalletID, amt wallet.Balance,
	fee uint64, dest cipher.Address) interface{} {
	self.Requests <- func() interface{} {
		return Spend2(self.D.Visor, wrpc,
			walletID, amt, fee, dest)
	}
	r := <-self.Responses
	return r
}
*/

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
func Spend(v *daemon.Visor, wrpc *WalletRPC,
	walletID wallet.WalletID, amt wallet.Balance, fee uint64,
	dest cipher.Address) *SpendResult {

	txn, err := Spend2(v.Visor, wrpc, walletID, amt, fee, dest)
	errString := ""
	if err != nil {
		errString = err.Error()
		logger.Error("Failed to make a spend: %v", err)
	}
	b := wrpc.GetWalletBalance(v.Visor, walletID)

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
	tx, err := visor.CreateSpendingTransaction(*wallet, self.Unconfirmed,
		&self.Blockchain.Unspent, self.Blockchain.Time(), amt, fee,
		dest)
	if err != nil {
		return tx, err
	}

	if err := tx.Verify(); err != nil {
		log.Panicf("Invalid transaction, %v", err)
	}

	if err := visor.VerifyTransaction(self.Blockchain, &tx); err != nil {
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
		SendOr404(w, Wg.GetWalletBalance(gateway.D.Visor.Visor, wallet.WalletID(id)))
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
func walletSpendHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletId := wallet.WalletID(r.FormValue("id"))
		if walletId == "" {
			Error400(w, "Missing wallet_id")
			return
		}
		sdst := r.FormValue("dst")
		if sdst == "" {
			Error400(w, "Missing destination address \"dst\"")
			return
		}
		dst, err := cipher.DecodeBase58Address(sdst)
		if err != nil {
			Error400(w, "Invalid destination address")
			return
		}
		sfee := r.FormValue("fee")
		fee, err := strconv.ParseUint(sfee, 10, 64)
		if err != nil {
			Error400(w, "Invalid \"fee\" value")
			return
		}
		scoins := r.FormValue("coins")
		shours := r.FormValue("hours")
		coins, err := strconv.ParseUint(scoins, 10, 64)
		if err != nil {
			Error400(w, "Invalid \"coins\" value")
			return
		}
		hours, err := strconv.ParseUint(shours, 10, 64)
		if err != nil {
			Error400(w, "Invalid \"hours\" value")
			return
		}
		SendOr404(w, Spend(gateway.D.Visor, Wg, walletId, wallet.NewBalance(coins, hours),
			fee, dst))
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

// Returns the outputs for a wallet
func getOutputsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ret := gateway.Visor.GetUnspentOutputReadables(gateway.V)
		SendOr404(w, ret)
	}
}

func injectTransaction(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ret := gateway.Visor.GetUnspentOutputReadables(gateway.V)
		//SendOr404(w, ret)
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

	//inject a transaction into network
	mux.HandleFunc("/injectTransaction", injectTransaction(gateway))
}
