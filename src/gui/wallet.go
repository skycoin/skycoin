package gui

// Wallet-related information for the GUI
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util"
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

// WalletRPC wallet rpc
type WalletRPC struct {
	Wallets         wallet.Wallets
	WalletDirectory string
	Options         []wallet.Option
}

// NotesRPC note rpc
type NotesRPC struct {
	Notes           wallet.Notes
	WalletDirectory string
}

// Wg use a global for now
var Wg *WalletRPC

// Ng global note
var Ng *NotesRPC

// InitWalletRPC init wallet rpc
func InitWalletRPC(walletDir string, options ...wallet.Option) {
	Wg = NewWalletRPC(walletDir, options...)
	Ng = NewNotesRPC(walletDir)
}

// NewNotesRPC new notes rpc
func NewNotesRPC(walletDir string) *NotesRPC {
	rpc := &NotesRPC{}
	if err := os.MkdirAll(walletDir, os.FileMode(0700)); err != nil {
		log.Panicf("Failed to create notes directory %s: %v", walletDir, err)
	}
	rpc.WalletDirectory = walletDir
	w, err := wallet.LoadNotes(rpc.WalletDirectory)
	if err != nil {
		log.Panicf("Failed to load all notes: %v", err)
	}
	wallet.CreateNoteFileIfNotExist(walletDir)
	rpc.Notes = w
	return rpc
}

// NewWalletRPC new wallet rpc
func NewWalletRPC(walletDir string, options ...wallet.Option) *WalletRPC {
	rpc := &WalletRPC{}
	if err := os.MkdirAll(walletDir, os.FileMode(0700)); err != nil {
		log.Panicf("Failed to create wallet directory %s: %v", walletDir, err)
	}

	rpc.WalletDirectory = walletDir
	for i := range options {
		rpc.Options = append(rpc.Options, options[i])
	}

	w, err := wallet.LoadWallets(rpc.WalletDirectory)
	if err != nil {
		log.Panicf("Failed to load all wallets: %v", err)
	}
	rpc.Wallets = w

	if len(rpc.Wallets) == 0 {
		wltName := wallet.NewWalletFilename()
		rpc.CreateWallet(wltName)

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

// ReloadWallets reload wallets
func (wlt *WalletRPC) ReloadWallets() error {
	wallets, err := wallet.LoadWallets(wlt.WalletDirectory)
	if err != nil {
		return err
	}
	wlt.Wallets = wallets
	return nil
}

// SaveWallet saves a wallet
func (wlt *WalletRPC) SaveWallet(walletID string) error {
	if w, ok := wlt.Wallets.Get(walletID); ok {
		return w.Save(wlt.WalletDirectory)
	}
	return fmt.Errorf("Unknown wallet %s", walletID)
}

// SaveWallets saves wallets
func (wlt *WalletRPC) SaveWallets() map[string]error {
	return wlt.Wallets.Save(wlt.WalletDirectory)
}

// CreateWallet creates wallet
func (wlt *WalletRPC) CreateWallet(wltName string, options ...wallet.Option) (wallet.Wallet, error) {
	ops := make([]wallet.Option, 0, len(wlt.Options)+len(options))
	ops = append(ops, wlt.Options...)
	ops = append(ops, options...)
	w := wallet.NewWallet(wltName, ops...)
	// generate a default address
	w.GenerateAddresses(1)

	if err := wlt.Wallets.Add(w); err != nil {
		return wallet.Wallet{}, err
	}

	return w, nil
}

// NewAddresses generate address entries in specific wallet,
// return nil if wallet does not exist.
func (wlt *WalletRPC) NewAddresses(wltID string, num int) ([]cipher.Address, error) {
	return wlt.Wallets.NewAddresses(wltID, num)
}

// GetWalletReadable returns a readable wallet
func (wlt *WalletRPC) GetWalletReadable(walletID string) *wallet.ReadableWallet {
	if w, ok := wlt.Wallets.Get(walletID); ok {
		return wallet.NewReadableWallet(w)
	}
	return nil
}

// GetWalletsReadable returns readable wallets
func (wlt *WalletRPC) GetWalletsReadable() []*wallet.ReadableWallet {
	return wlt.Wallets.ToReadable()
}

// GetNotesReadable returns readable notes
func (nt *NotesRPC) GetNotesReadable() wallet.ReadableNotes {
	return nt.Notes.ToReadable()
}

// GetWallet returns wallet of give id
func (wlt *WalletRPC) GetWallet(walletID string) *wallet.Wallet {
	if w, ok := wlt.Wallets.Get(walletID); ok {
		return &w
	}
	return nil
}

// GetWalletBalance modify to return error
// NOT WORKING
// actually uses visor
func (wlt *WalletRPC) GetWalletBalance(gateway *daemon.Gateway,
	walletID string) (wallet.BalancePair, error) {

	w, ok := wlt.Wallets.Get(walletID)
	if !ok {
		log.Printf("GetWalletBalance: ID NOT FOUND: id= '%s'", walletID)
		return wallet.BalancePair{}, errors.New("Id not found")
	}

	return gateway.WalletBalance(w), nil
}

// HasUnconfirmedTransactions checks if the wallet has pending, unconfirmed transactions
// - do not allow any transactions if there are pending
//Check if any of the outputs are spent
func (wlt *WalletRPC) HasUnconfirmedTransactions(v *visor.Visor,
	wallet *wallet.Wallet) bool {

	if wallet == nil {
		log.Panic("Wallet does not exist")
	}

	// auxs := v.Blockchain.GetUnspent().AllForAddresses(wallet.GetAddresses())
	unspent := v.Blockchain.GetUnspent()
	puxs := v.Unconfirmed.SpendsForAddresses(unspent, wallet.GetAddressSet())

	//no transactions
	if len(puxs) == 0 {
		return true
	}

	return false
}

// SpendResult represents the result of spending
type SpendResult struct {
	Balance     wallet.BalancePair        `json:"balance"`
	Transaction visor.ReadableTransaction `json:"txn"`
	Error       string                    `json:"error"`
}

// Spend TODO
// - split send into
// -- get addresses
// -- get unspent outputs
// -- construct transaction
// -- sign transaction
// -- inject transaction
func Spend(gateway *daemon.Gateway,
	wrpc *WalletRPC,
	walletID string,
	amt wallet.Balance,
	fee uint64,
	dest cipher.Address) *SpendResult {
	var txn coin.Transaction
	var b wallet.BalancePair
	var err error
	for {
		txn, err = Spend2(gateway, wrpc, walletID, amt, fee, dest)
		if err != nil {
			logger.Error("Transaction creation failed: %v", err)
			break
		}
		log.Printf("Spend: \ntx= \n %s \n", visor.TransactionToJSON(txn))

		b, err = wrpc.GetWalletBalance(gateway, walletID)
		if err != nil {
			logger.Error("Get wallet balance failed: %v", err)
			break
		}

		txn, err = gateway.InjectTransaction(txn)
		if err != nil {
			logger.Error("Inject transaction failed: %v", err)
			break
		}
		break
	}

	if err != nil {
		return &SpendResult{
			Error: err.Error(),
		}
	}

	return &SpendResult{
		Balance:     b,
		Transaction: visor.NewReadableTransaction(&visor.Transaction{Txn: txn}),
	}
}

// Spend2 Creates a transaction spending amt with additional fee.  Fee is in addition
// to the base required fee given amt.Hours.
// TODO
// - pull in outputs from blockchain from wallet
// - create transaction here
// - sign transction and return
func Spend2(gateway *daemon.Gateway, wrpc *WalletRPC, walletID string, amt wallet.Balance,
	fee uint64, dest cipher.Address) (coin.Transaction, error) {

	wallet, ok := wrpc.Wallets.Get(walletID)
	if !ok {
		return coin.Transaction{}, fmt.Errorf("Unknown wallet %v", walletID)
	}

	return gateway.CreateSpendingTransaction(wallet, amt, dest)
}

/*
REFACTOR
*/

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
func walletBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		r.ParseForm()
		b, err := Wg.GetWalletBalance(gateway, id)

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
			for i, addr := range addrsStr {
				addrs[i] = cipher.MustDecodeBase58Address(addr)
			}

			bal := gateway.AddressesBalance(addrs)

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

		walletID := r.FormValue("id")
		if walletID == "" {
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

		var hours uint64
		var fee uint64 //doesnt work/do anything right now

		//MOVE THIS INTO HERE
		ret := Spend(gateway, Wg, walletID, wallet.NewBalance(coins, hours), fee, dst)

		if ret.Error != "" {
			wh.Error400(w, fmt.Sprintf("Spend Failed: %s", ret.Error))
			return
		}
		wh.SendOr404(w, ret)
	}
}

// Create a wallet Name is set by creation date
func notesCreate(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("API request made to create a note")
		note := r.FormValue("note")
		transactionID := r.FormValue("transaction_id")
		newNote := wallet.Note{
			TxID:  transactionID,
			Value: note,
		}
		Ng.Notes.SaveNote(Ng.WalletDirectory, newNote)
		rlt := Ng.GetNotesReadable()
		wh.SendOr500(w, rlt)
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
			wlt, err = Wg.CreateWallet(wltName, wallet.OptSeed(seed), wallet.OptLabel(label))
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

// Returns a wallet by ID if GET.  Creates or updates a wallet if POST.
func notesHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//ret := wallet.Wallets.ToPublicReadable()
		ret := Ng.GetNotesReadable()
		wh.SendOr404(w, ret)
	}
}

// Returns JSON of unconfirmed transactions for user's wallet
func walletTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			wallet := Wg.GetWallet(r.FormValue("id"))
			addresses := wallet.GetAddresses()
			ret := gateway.GetUnconfirmedTxns(addresses)

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

// WalletFolder struct
type WalletFolder struct {
	Address string `json:"address"`
}

// Loads/unloads wallets from the wallet directory
func getWalletFolder(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ret := WalletFolder{
			Address: util.UserHome() + "/.skycoin/wallets",
		}
		wh.SendOr404(w, ret)
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
			var addrs []string
			var hashes []string

			trimSpace := func(vs []string) []string {
				for i := range vs {
					vs[i] = strings.TrimSpace(vs[i])
				}
				return vs
			}

			addrStr := r.FormValue("addrs")
			if addrStr != "" {
				addrs = trimSpace(strings.Split(addrStr, ","))
			}

			hashStr := r.FormValue("hashes")
			if hashStr != "" {
				hashes = trimSpace(strings.Split(hashStr, ","))
			}

			filters := []daemon.OutputsFilter{}
			if len(addrs) > 0 {
				filters = append(filters, daemon.FbyAddresses(addrs))
			}

			if len(hashes) > 0 {
				filters = append(filters, daemon.FbyHashes(hashes))
			}

			outs := gateway.GetUnspentOutputs(filters...)

			wh.SendOr404(w, outs)
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

// RegisterWalletHandlers registers wallet handlers
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

	mux.HandleFunc("/wallets/folderName", getWalletFolder(gateway))

	//get set of unspent outputs
	mux.HandleFunc("/outputs", getOutputsHandler(gateway))

	// get balance of addresses
	mux.HandleFunc("/balance", getBalanceHandler(gateway))

	// generate wallet seed
	mux.Handle("/wallet/newSeed", newWalletSeed(gateway))

	// generate wallet seed
	mux.Handle("/notes", notesHandler(gateway))

	mux.Handle("/notes/create", notesCreate(gateway))
}
