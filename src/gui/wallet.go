package gui

// Wallet-related information for the GUI
import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"

	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"

	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
)

// SpendResult represents the result of spending
type SpendResult struct {
	Balance     *wallet.BalancePair        `json:"balance,omitempty"`
	Transaction *visor.ReadableTransaction `json:"txn,omitempty"`
	Error       string                     `json:"error,omitempty"`
}

type UnconfirmedTxnsResponse struct {
	Transactions []visor.ReadableUnconfirmedTxn `json:"transactions"`
}

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
// URI: /wallet/balance
// Method: GET
// Args:
//     id: wallet id [required]
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
			switch err {
			case wallet.ErrWalletNotExist:
				wh.Error404(w)
				break
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
				break
			default:
				wh.Error500Msg(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, b)
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
// URI: /wallet/spend
// Method: POST
// Args:
//     id: wallet id
//     dst: recipient address
//     coins: the number of droplet you will send
// Response:
//     balance: new balance of the wallet
//     txn: spent transaction
//     error: an error that may have occured after broadcast the transaction to the network
//         if this field is not empty, the spend succeeded, but the response data could not be prepared
func walletSpendHandler(gateway Gatewayer) http.HandlerFunc {
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

		tx, err := gateway.Spend(wltID, coins, dst)
		switch err {
		case nil:
		case fee.ErrTxnNoFee, wallet.ErrSpendingUnconfirmed, wallet.ErrInsufficientBalance:
			wh.Error400(w, err.Error())
			return
		case wallet.ErrWalletNotExist:
			wh.Error404(w)
			return
		case wallet.ErrWalletApiDisabled:
			wh.Error403(w)
			return
		default:
			wh.Error500Msg(w, err.Error())
			return
		}

		txStr, err := visor.TransactionToJSON(*tx)
		if err != nil {
			logger.Error(err.Error())
			wh.SendJSONOr500(logger, w, SpendResult{
				Error: err.Error(),
			})
			return
		}

		logger.Info("Spend: \ntx= \n %s \n", txStr)

		var ret SpendResult

		ret.Transaction, err = visor.NewReadableTransaction(&visor.Transaction{Txn: *tx})
		if err != nil {
			err = fmt.Errorf("Creation of new readable transaction failed: %v", err)
			logger.Error(err.Error())
			ret.Error = err.Error()
			wh.SendJSONOr500(logger, w, ret)
			return
		}

		// Get the new wallet balance
		b, err := gateway.GetWalletBalance(wltID)
		if err != nil {
			err = fmt.Errorf("Get wallet balance failed: %v", err)
			logger.Error(err.Error())
			ret.Error = err.Error()
			wh.SendJSONOr500(logger, w, ret)
			return
		}
		ret.Balance = &b

		wh.SendJSONOr500(logger, w, ret)
	}
}

// Loads wallet from seed, will scan ahead N address and
// load addresses till the last one that have coins.
// Method: POST
// Args:
//     seed: wallet seed [required]
//     label: wallet label [required]
//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
func walletCreate(gateway Gatewayer) http.HandlerFunc {
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
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
				return
			default:
				wh.Error400(w, err.Error())
				return
			}
		}

		wlt, err = gateway.ScanAheadWalletAddresses(wlt.GetFilename(), scanN-1)
		if err != nil {
			logger.Error("gateway.ScanAheadWalletAddresses failed: %v", err)
			wh.Error500(w)
			return
		}

		rlt := wallet.NewReadableWallet(wlt)
		wh.SendJSONOr500(logger, w, rlt)
	}
}

// Genreates new addresses
// URI: /wallet/newAddress
// Method: POST
// Args:
//     id: wallet id [required]
//     num: number of address need to create [optional, if not set the default value is 1]
func walletNewAddresses(gateway Gatewayer) http.HandlerFunc {
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
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
				return
			default:
				wh.Error400(w, err.Error())
				return
			}
		}

		var rlt = struct {
			Addresses []string `json:"addresses"`
		}{}

		for _, a := range addrs {
			rlt.Addresses = append(rlt.Addresses, a.String())
		}

		wh.SendJSONOr500(logger, w, rlt)
		return
	}
}

// Update wallet label
// URI: /wallet/update
// Method: POST
// Args:
//     id: wallet id [required]
//     label: the label the wallet will be updated to [required]
func walletUpdateHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

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
			logger.Errorf("update wallet label failed: %v", err)

			switch err {
			case wallet.ErrWalletNotExist:
				wh.Error404(w)
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error500Msg(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}

// Returns a wallet by id
// URI: /wallet
// Method: GET
// Args:
//     id: wallet id [required]
func walletGet(gateway Gatewayer) http.HandlerFunc {
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
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error400(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, wlt)
	}
}

// Returns JSON of unconfirmed transactions for user's wallet
// URI: /wallet/transactions
// Method: GET
// Args:
//     id: wallet id [required]
func walletTransactionsHandler(gateway Gatewayer) http.HandlerFunc {
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
			logger.Error("get wallet unconfirmed transactions failed: %v", err)
			switch err {
			case wallet.ErrWalletNotExist:
				wh.Error404(w)
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error500Msg(w, err.Error())
			}
			return
		}

		unconfirmedTxns, err := visor.NewReadableUnconfirmedTxns(txns)
		if err != nil {
			wh.Error500Msg(w, err.Error())
			return
		}

		unconfirmedTxnResp := UnconfirmedTxnsResponse{
			Transactions: unconfirmedTxns,
		}
		wh.SendJSONOr500(logger, w, unconfirmedTxnResp)
	}
}

// Returns all loaded wallets
// RUI: /wallets
// Method: GET
func walletsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wlts, err := gateway.GetWallets()
		if err != nil {
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error500(w)
			}
			return
		}
		wh.SendJSONOr500(logger, w, wlts.ToReadable())
	}
}

// WalletFolder struct
type WalletFolder struct {
	Address string `json:"address"`
}

// Returns the wallet directory path
// URI: /wallets/folderName
// Method: GET
func getWalletFolder(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr, err := gateway.GetWalletDir()
		if err != nil {
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error500(w)
			}
			return
		}
		ret := WalletFolder{
			Address: addr,
		}
		wh.SendJSONOr500(logger, w, ret)
	}
}

// Generates wallet seed
// URI: /wallet/newSeed
// Method: GET
// Args:
//     entropy: entropy bitsize [optional, default value of 128 will be used if not set]
func newWalletSeed(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		if gateway.IsWalletAPIDisabled() {
			wh.Error403(w)
			return
		}

		entropyValue := r.FormValue("entropy")
		if entropyValue == "" {
			entropyValue = "128"
		}

		entropyBits, err := strconv.Atoi(entropyValue)
		if err != nil {
			wh.Error400(w, "invalid entropy")
			return
		}

		// Entropy bit size can either be 128 or 256
		if entropyBits != 128 && entropyBits != 256 {
			wh.Error400(w, "entropy length must be 128 or 256")
			return
		}

		entropy, err := bip39.NewEntropy(entropyBits)
		if err != nil {
			logger.Error("bip39.NewEntropy failed: %v", err)
			wh.Error500(w)
			return
		}

		mnemonic, err := bip39.NewMnemonic(entropy)
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
		wh.SendJSONOr500(logger, w, rlt)
	}
}

func walletUnloadHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		if err := gateway.UnloadWallet(id); err != nil {
			switch err {
			case wallet.ErrWalletApiDisabled:
				wh.Error403(w)
			default:
				wh.Error500(w)
			}
		}
	}
}
