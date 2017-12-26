package gui

// Wallet-related information for the GUI
import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	bip39 "github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

type Gatewayer interface {
	Spend(wltID string, coins uint64, dest cipher.Address) (*coin.Transaction, error)
	GetWalletBalance(wltID string) (wallet.BalancePair, error)
}

// SpendResult represents the result of spending
type SpendResult struct {
	Balance     *wallet.BalancePair        `json:"balance,omitempty"`
	Transaction *visor.ReadableTransaction `json:"txn,omitempty"`
	Error       string                     `json:"error,omitempty"`
}

// Spend spend coins from specific wallet
func Spend(gateway Gatewayer, walletID string, coins uint64, dest cipher.Address) *SpendResult {
	var tx *coin.Transaction
	var b wallet.BalancePair
	var err error
	for {
		tx, err = gateway.Spend(walletID, coins, dest)
		if err != nil {
			break
		}

		var txStr string
		txStr, err = visor.TransactionToJSON(*tx)
		if err != nil {
			break
		}

		logger.Info("Spend: \ntx= \n %s \n", txStr)

		b, err = gateway.GetWalletBalance(walletID)
		if err != nil {
			err = fmt.Errorf("Get wallet balance failed: %v", err)
			break
		}

		break
	}

	if err != nil {
		return &SpendResult{
			Error: err.Error(),
		}
	}

	rbTx, err := visor.NewReadableTransaction(&visor.Transaction{Txn: *tx})
	if err != nil {
		logger.Error("%v", err)
		return &SpendResult{}
	}

	return &SpendResult{
		Balance:     &b,
		Transaction: rbTx,
	}
}

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
func walletBalanceHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		b, err := gateway.GetWalletBalance(wltID)
		if err != nil {
			logger.Error("Get wallet balance failed: %v", err)
			return
		}
		wh.SendOr404(w, b)
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
// URI: /wallet/spend
// Method: POST
// Args:
//  id: wallet id
//	dst: recipient address
// 	coins: the number of droplet you will send
func WalletSpendHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		sdst := r.FormValue("dst")
		if sdst == "" {
			wh.Error400(w, "missing destination address \"dst\"")
			return
		}
		dst, err := cipher.DecodeBase58Address(sdst)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("invalid destination address: %v", err))
			return
		}

		scoins := r.FormValue("coins")
		coins, err := strconv.ParseUint(scoins, 10, 64)
		if err != nil {
			wh.Error400(w, `invalid "coins" value`)
			return
		}

		if coins <= 0 {
			wh.Error400(w, `invalid "coins" value, must > 0`)
			return
		}

		ret := Spend(gateway, wltID, coins, dst)
		if ret.Error != "" {
			logger.Error(ret.Error)
		}

		wh.SendOr404(w, ret)
	}
}

// Loads wallet from seed, will scan ahead N address and
// load addresses till the last one that have coins.
// Method: POST
// Args:
//     seed: wallet seed [required]
//     label: wallet label [required]
//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
func walletCreate(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		seed := r.FormValue("seed")
		label := r.FormValue("label")
		scanNStr := r.FormValue("scan")

		if seed == "" {
			wh.Error400(w, "missing seed")
			return
		}

		if label == "" {
			wh.Error400(w, "missing label")
			return
		}

		var scanN uint64 = 1
		if scanNStr != "" {
			var err error
			scanN, err = strconv.ParseUint(scanNStr, 10, 64)
			if err != nil {
				wh.Error400(w, "invalid scan value")
				return
			}
		}

		if scanN == 0 {
			wh.Error400(w, "scan must be > 0")
			return
		}

		wlt, err := gateway.CreateWallet("", wallet.Options{
			Seed:  seed,
			Label: label,
		})
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wlt, err = gateway.ScanAheadWalletAddresses(wlt.GetFilename(), scanN-1)
		if err != nil {
			logger.Error("gateway.ScanAheadWalletAddresses failed: %v", err)
			wh.Error500(w)
			return
		}

		rlt := wallet.NewReadableWallet(wlt)
		wh.SendOr500(w, rlt)
	}
}

// method: POST
// url: /wallet/newAddress
// params:
// 		id: wallet id
// 	   num: number of address need to create, if not set the default value is 1
func walletNewAddresses(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		// the number of address that need to create, default is 1
		var n uint64 = 1
		var err error
		num := r.FormValue("num")
		if num != "" {
			n, err = strconv.ParseUint(num, 10, 64)
			if err != nil {
				wh.Error400(w, "invalid num value")
				return
			}
		}

		addrs, err := gateway.NewAddresses(wltID, n)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		var rlt = struct {
			Address []string `json:"addresses"`
		}{}

		for _, a := range addrs {
			rlt.Address = append(rlt.Address, a.String())
		}

		wh.SendOr404(w, rlt)
		return
	}
}

// Update wallet label
func walletUpdateHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Update wallet
		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		label := r.FormValue("label")
		if label == "" {
			wh.Error400(w, "missing label")
			return
		}

		if err := gateway.UpdateWalletLabel(wltID, label); err != nil {
			wh.Error400(w, fmt.Sprintf("update wallet label failed: %v", err))
			return
		}

		wh.SendOr404(w, "success")
	}
}

// Returns a wallet by id
func walletGet(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, fmt.Sprintf("missing wallet id"))
			return
		}

		wlt, err := gateway.GetWallet(wltID)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wh.SendOr404(w, wlt)
	}
}

// Returns JSON of unconfirmed transactions for user's wallet
func walletTransactionsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		txns, err := gateway.GetWalletUnconfirmedTxns(wltID)
		if err != nil {
			wh.Error400(w, fmt.Sprintf("get wallet unconfirmed transactions failed: %v", err))
			return
		}

		wh.SendOr404(w, txns)
	}
}

// Returns all loaded wallets
func walletsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wlts := gateway.GetWallets().ToReadable()
		wh.SendOr404(w, wlts)
	}
}

// Loads/unloads wallets from the wallet directory
func walletsReloadHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := gateway.ReloadWallets(); err != nil {
			logger.Error("reload wallet failed: %v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr404(w, "success")
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
			Address: gateway.GetWalletDir(),
		}
		wh.SendOr404(w, ret)
	}
}

func newWalletSeed(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mnemonic, err := bip39.NewDefaultMnemomic()
		if err != nil {
			logger.Error("bip39.NewDefaultMnemomic failed: %v", err)
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

	// Loads wallet from seed, will scan ahead N address and
	// load addresses till the last one that have coins.
	// Method: POST
	// Args:
	//     seed: wallet seed [required]
	//     label: wallet label [required]
	//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
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
	mux.HandleFunc("/wallet/spend", WalletSpendHandler(gateway))

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

	// Rescans the wallet directory and loads/unloads wallets based on which
	// files are present. Returns nothing if it works. Otherwise returns
	// 500 status with error message.
	mux.HandleFunc("/wallets/reload", walletsReloadHandler(gateway))

	mux.HandleFunc("/wallets/folderName", getWalletFolder(gateway))

	// generate wallet seed
	mux.Handle("/wallet/newSeed", newWalletSeed(gateway))

	// generate wallet seed
}
