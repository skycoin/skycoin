// Wallet-related information for the GUI
package gui

import (
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
	bip39 "github.com/tyler-smith/go-bip39"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
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
		wltName := wallet.NewWalletFilename()
		rpc.CreateWallet("", wltName, "")

		if err := rpc.SaveWallet(wltName); err != nil {
			log.Panicf("Failed to save wallets to %s: %v", rpc.WalletDirectory, err)
		}

		// newWlt := wallet.NewWallet("", wltName, wltName) //deterministic
		// newWlt.GenerateAddresses(1)
		// rpc.Wallets.Add(newWlt)
		// errs := rpc.Wallets.Save(rpc.WalletDirectory)
		// if len(errs) != 0 {
		// 	log.Panicf("Failed to save wallets to %s: %v", rpc.WalletDirectory, errs)
		// }
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

func (self *WalletRPC) SaveWallet(walletID string) error {
	if w, ok := self.Wallets.Get(walletID); ok {
		return w.Save(self.WalletDirectory)
	}
	return fmt.Errorf("Unknown wallet %s", walletID)
}

func (self *WalletRPC) SaveWallets() map[string]error {
	return self.Wallets.Save(self.WalletDirectory)
}

func (self *WalletRPC) CreateWallet(seed, wltName, label string) (wallet.Wallet, error) {
	w := wallet.NewWallet(seed, wltName, label)
	// generate a default address
	w.GenerateAddresses(1)

	if err := self.Wallets.Add(w); err != nil {
		return wallet.Wallet{}, err
	}

	return w, nil
}

// NewAddresses generate address entries in specific wallet,
// return nil if wallet does not exist.
func (rpc *WalletRPC) NewAddresses(wltID string, num int) ([]cipher.Address, error) {
	return rpc.Wallets.NewAddresses(wltID, num)
}

func (self *WalletRPC) GetWalletsReadable() []*wallet.ReadableWallet {
	return self.Wallets.ToReadable()
}

func (self *WalletRPC) GetWalletReadable(walletID string) *wallet.ReadableWallet {
	if w, ok := self.Wallets.Get(walletID); ok {
		return wallet.NewReadableWallet(w)
	}
	return nil
}

func (self *WalletRPC) GetWallet(walletID string) *wallet.Wallet {
	if w, ok := self.Wallets.Get(walletID); ok {
		return &w
	}
	return nil
}

//modify to return error
// NOT WORKING
// actually uses visor
func (self *WalletRPC) GetWalletBalance(v *visor.Visor,
	walletID string) (wallet.BalancePair, error) {

	wlt, ok := self.Wallets.Get(walletID)
	if !ok {
		log.Printf("GetWalletBalance: ID NOT FOUND: id= '%s'", walletID)
		return wallet.BalancePair{}, errors.New("Id not found")
	}
	auxs := v.Blockchain.GetUnspent().AllForAddresses(wlt.GetAddresses())
	unspent := v.Blockchain.GetUnspent()
	puxs := v.Unconfirmed.SpendsForAddresses(unspent, wlt.GetAddressSet())

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

	auxs := v.Blockchain.GetUnspent().AllForAddresses(wallet.GetAddresses())
	unspent := v.Blockchain.GetUnspent()
	puxs := v.Unconfirmed.SpendsForAddresses(unspent, wallet.GetAddressSet())

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
	walletID string, amt wallet.Balance, fee uint64,
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
		Transaction: visor.NewReadableTransaction(&visor.Transaction{Txn: txn}),
		Error:       errString,
	}
}

// Creates a transaction spending amt with additional fee.  Fee is in addition
// to the base required fee given amt.Hours.
// TODO
// - pull in outputs from blockchain from wallet
// - create transaction here
// - sign transction and return
func Spend2(self *visor.Visor, wrpc *WalletRPC, walletID string, amt wallet.Balance,
	fee uint64, dest cipher.Address) (coin.Transaction, error) {

	wallet, ok := wrpc.Wallets.Get(walletID)
	if !ok {
		return coin.Transaction{}, fmt.Errorf("Unknown wallet %v", walletID)
	}
	//pull in outputs and do this here
	//FIX
	unspent := self.Blockchain.GetUnspent()
	tx, err := visor.CreateSpendingTransaction(wallet, self.Unconfirmed,
		unspent, self.Blockchain.Time(), amt, dest)
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

		b, err := Wg.GetWalletBalance(gateway.D.Visor.Visor, id)

		if err != nil {
			_ = err
		}
		//log.Printf("%v, %v, %v \n", r.URL.String(), r.RequestURI, r.Form)
		wh.SendOr404(w, b)
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
			auxs := v.Blockchain.GetUnspent().AllForAddresses(addrs)
			unspent := v.Blockchain.GetUnspent()
			puxs := v.Unconfirmed.SpendsForAddresses(unspent, addrSet)

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

			wh.SendOr404(w, bal)
		}
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
func walletSpendHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//log.Printf("Spend1")

		if r.FormValue("id") == "" {
			wh.Error400(w, "Missing wallet_id")
			return
		}

		walletId := r.FormValue("id")
		if walletId == "" {
			wh.Error400(w, "Invalid Wallet Id")
			return
		}
		sdst := r.FormValue("dst")
		if sdst == "" {
			wh.Error400(w, "Missing destination address \"dst\"")
			return
		}
		dst, err := cipher.DecodeBase58Address(sdst)
		if err != nil {
			//Error400(w, "Invalid destination address: %v", err)
			wh.Error400(w, "Invalid destination address: %v", err.Error())
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
			wh.Error400(w, "Invalid \"coins\" value")
			return
		}

		var hours uint64 = 0
		var fee uint64 = 0 //doesnt work/do anything right now

		//MOVE THIS INTO HERE
		ret := Spend(gateway.D, gateway.D.Visor, Wg, walletId, wallet.NewBalance(coins, hours), fee, dst)

		if ret.Error != "" {
			wh.Error400(w, "Spend Failed: %s", ret.Error)
		}
		wh.SendOr404(w, ret)
	}
}

// Create a wallet Name is set by creation date
func walletCreate(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("API request made to create a wallet")
		seed := r.FormValue("seed")
		label := r.FormValue("label")
		wltName := wallet.NewWalletFilename()
		var wlt wallet.Wallet
		var err error
		// the wallet name may dup, rename it till no conflict.
		for {
			wlt, err = Wg.CreateWallet(seed, wltName, label)
			if err != nil && strings.Contains(err.Error(), "renaming") {
				wltName = wallet.NewWalletFilename()
				continue
			}
			break
		}

		if err := Wg.SaveWallet(wlt.GetID()); err != nil {
			wh.Error400(w, err.Error())
			return
		}

		rlt := wallet.NewReadableWallet(wlt)
		wh.SendOr500(w, rlt)
	}
}

func walletNewAddresses(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			wh.Error405(w, "")
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "wallet id not set")
			return
		}

		addrs, err := Wg.NewAddresses(wltID, 1)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if err := Wg.SaveWallet(wltID); err != nil {
			wh.Error500(w, "")
			return
		}

		var rlt = struct {
			Address string `json:"address"`
		}{
			addrs[0].String(),
		}
		wh.SendOr404(w, rlt)
		return
	}
}

// Update wallet label
func walletUpdateHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update wallet
		id := r.FormValue("id")
		if id == "" {
			wh.Error400(w, "wallet id is empty")
			return
		}

		label := r.FormValue("label")
		if label == "" {
			wh.Error400(w, "label is empty")
			return
		}

		wlt := Wg.GetWallet(id)
		if wlt == nil {
			wh.Error404(w, fmt.Sprintf("wallet of id: %v does not exist", id))
			return
		}

		wlt.SetLabel(label)
		if err := Wg.SaveWallet(wlt.GetID()); err != nil {
			m := "Failed to save wallet: %v"
			logger.Critical(m, "Failed to update label of wallet %v", id)
			wh.Error500(w, "Update wallet failed")
			return
		}

		wh.SendOr404(w, "success")
	}
}

// Returns a wallet by ID if GET.  Creates or updates a wallet if POST.
func walletGet(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			ret := Wg.GetWallet(r.FormValue("id"))
			wh.SendOr404(w, ret)
		}
	}
}

// Returns JSON of pending transactions for user's wallet
func walletTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			wallet := Wg.GetWallet(r.FormValue("id"))
			addresses := wallet.GetAddresses()
			ret := gateway.Visor.GetWalletTransactions(gateway.V, addresses)

			wh.SendOr404(w, ret)
		}
	}
}

// Returns all loaded wallets
func walletsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ret := wallet.Wallets.ToPublicReadable()
		ret := Wg.GetWalletsReadable()
		wh.SendOr404(w, ret)
	}
}

// Saves all loaded wallets
func walletsSaveHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		errs := Wg.SaveWallets() // (map[string]error)
		if len(errs) != 0 {
			err := ""
			for id, e := range errs {
				err += id + ": " + e.Error()
			}
			wh.Error500(w, err)
		}
	}
}

// Loads/unloads wallets from the wallet directory
func walletsReloadHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := Wg.ReloadWallets()
		if err != nil {
			wh.Error500(w, err.(error).Error())
		}
	}
}

// getOutputsHandler get utxos base on the filters in url params.
// mode: GET
// url: /outputs?addrs=[:addrs]&hashes=[:hashes]
// if addrs and hashes are not specificed, return all unspent outputs.
// if both addrs and hashes are specificed, then both those filters are need to be matched.
// if only specify one filter, then return outputs match the filter.
func getOutputsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			uxouts := gateway.Visor.GetUnspentOutputReadables(gateway.V)
			rawaddrs := r.FormValue("addrs")
			hashes := r.FormValue("hashes")

			if rawaddrs == "" && hashes == "" {
				wh.SendOr404(w, uxouts)
				return
			}

			addrMatch := []visor.ReadableOutput{}
			if rawaddrs != "" {
				addrs := strings.Split(rawaddrs, ",")
				addrMap := make(map[string]bool)
				for _, addr := range addrs {
					addrMap[addr] = true
				}

				for _, u := range uxouts {
					if _, ok := addrMap[u.Address]; ok {
						addrMatch = append(addrMatch, u)
					}
				}
			}

			hsMatch := []visor.ReadableOutput{}
			hsMatchMap := map[string]bool{}
			if hashes != "" {
				hs := strings.Split(hashes, ",")
				hsMap := make(map[string]bool)
				for _, h := range hs {
					hsMap[h] = true
				}

				for _, u := range uxouts {
					if _, ok := hsMap[u.Hash]; ok {
						hsMatch = append(hsMatch, u)
						hsMatchMap[u.Hash] = true
					}
				}
			}

			ret := []visor.ReadableOutput{}
			if rawaddrs != "" && hashes != "" {
				for _, u := range addrMatch {
					if _, ok := hsMatchMap[u.Hash]; ok {
						ret = append(ret, u)
					}
				}
				wh.SendOr404(w, ret)
				return
			}

			wh.SendOr404(w, append(addrMatch, hsMatch...))
		}
	}
}

func newWalletSeed(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			wh.Error500(w)
			return
		}

		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			wh.Error500(w)
			return
		}

		var rlt = struct {
			Seed string `json:"seed"`
		}{
			mnemonic,
		}

		wh.SendOr404(w, rlt)
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

	mux.HandleFunc("/wallet/newAddress", walletNewAddresses(gateway))

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

	// Update wallet label
	// 		GET Arguments:
	// 			id: wallet id
	// 			label: wallet label
	mux.HandleFunc("/wallet/update", walletUpdateHandler(gateway))

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

	// generate wallet seed
	mux.Handle("/wallet/seed", newWalletSeed(gateway))
}
