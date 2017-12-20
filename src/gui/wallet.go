package gui

// Wallet-related information for the GUI
import (
	"errors"
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

// SpendResult represents the result of spending
type SpendResult struct {
	Balance     *wallet.BalancePair        `json:"balance,omitempty"`
	Transaction *visor.ReadableTransaction `json:"txn,omitempty"`
	Error       string                     `json:"error,omitempty"`
}

// Loads wallet from seed, will scan ahead N address and
// load addresses till the last one that have coins.
// URI: /wallet/create
// Method: POST
// Args:
//     seed: Wallet seed [required]
//     label: Wallet label [required]
//     scan: The number of addresses to scan ahead for balances [optional, must be > 0]
// 	   encrypt: Whether encrypt the wallet [optional, must be 0 or 1]
// 	   password: Password for encrypting wallet [optional, required if the encrypt value is 1]
func walletCreateHandler(gateway *daemon.Gateway) http.HandlerFunc {
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

		// Get wallet encrypt options
		encrypt, password, err := getWalletEncryptOptions(r)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wlt, err := gateway.CreateWallet("", wallet.Options{
			Seed:     seed,
			Label:    label,
			Encrypt:  encrypt,
			Password: password,
		})
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		wlt, err = gateway.ScanAheadWalletAddresses(wlt.Filename(), password, scanN-1)
		if err != nil {
			logger.Error("gateway.ScanAheadWalletAddresses failed: %v", err)
			wh.Error500(w)
			return
		}

		wh.SendOr500(w, wallet.NewReadableWallet(wlt))
	}
}

// Get wallet encrypt options, check the encrypt and password paramenter value.
// The encrypt value must be 0 or 1; password must be provided if the encrypt value is 1.
func getWalletEncryptOptions(r *http.Request) (bool, []byte, error) {
	encrypt := r.FormValue("encrypt")
	password := r.FormValue("password")
	switch encrypt {
	case "1":
		if password == "" {
			return false, nil, errors.New("password is required for wallet with encryption")
		}
		return true, []byte(password), nil
	case "0":
		return false, nil, nil
	default:
		return false, nil, fmt.Errorf("invalid [encrypt] value: %v, must be 0 or 1 ", encrypt)
	}
}

// URI: /wallet/newAddress
// Method: POST
// Args:
//     id: Wallet ID [required]
//     num: Number of address need to create [optional, if not set the default value is 1]
//     password: Wallet password [optional, required if the wallet is encrypted]
func walletNewAddressesHandler(gateway *daemon.Gateway) http.HandlerFunc {
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

		// The number of address that need to create, default is 1
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

		// Get wallet password
		password := r.FormValue("password")

		addrs, err := gateway.NewAddresses(wltID, []byte(password), n)
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

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
// URI: /wallet/balance
// Method: GET
// Args:
//     id: Wallet ID [required]
func walletBalanceHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
//    id: Wallet ID [required]
//	  dst: Recipient address [required]
// 	  coins: The number of droplet you will send [required]
//	  password: Wallet password [optional, required if the wallet is encrypted]
func walletSpendHandler(gateway *daemon.Gateway) http.HandlerFunc {
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

		password := r.FormValue("password")

		ret := spend(gateway, wltID, []byte(password), coins, dst)
		if ret.Error != "" {
			logger.Error(ret.Error)
		}

		wh.SendOr404(w, ret)
	}
}

// Spends coins from specific wallet.
// Set password as nil if spend coins from unencrypted wallet,
// otherwise the password must be provided.
func spend(gateway *daemon.Gateway, walletID string, password []byte, coins uint64, dest cipher.Address) *SpendResult {
	var tx *coin.Transaction
	var b wallet.BalancePair
	var err error
	for {
		tx, err = gateway.Spend(walletID, password, coins, dest)
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

// Updates wallet label
// Method: GET
// Args:
//     id: Wallet ID [required]
//     label: Wallet label [required]
func walletUpdateLabelHandler(gateway *daemon.Gateway) http.HandlerFunc {
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
// Method: GET
// URI: /wallet
// Args:
//     id: Wallet ID [required]
func walletGetHandler(gateway *daemon.Gateway) http.HandlerFunc {
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

// Returns wallet related unconfirmed transactions
// [Deprecated]:
//     No one is using this api and the URI is easy to misunderstanding,
//     this api should returns all transactions that are related to a wallet, while
//     the implementation actually returns the unconfirmed pending transactions of a wallet.
// Method: GET
// Args:
//     id: Wallet ID [required]
// Returns all pending transanction for all addresses by selected Wallet
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
// URI: /wallets
func walletsHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wlts := gateway.GetWallets().ToReadable()
		wh.SendOr404(w, wlts)
	}
}

// Loads/unloads wallets from the wallet directory
// URI: /wallets/reload
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

// Returns the path of wallet directory
// URI: /wallet/folderName
func getWalletFolder(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ret := WalletFolder{
			Address: gateway.GetWalletDir(),
		}
		wh.SendOr404(w, ret)
	}
}

// Generates and returns a new wallet seed
// URI: /wallet/newSeed
func newWalletSeedHandler(gateway *daemon.Gateway) http.HandlerFunc {
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

// Encrypts wallet
// Method: GET
// URI: /wallet/encrypt
// Args:
//     id: Wallet ID [required]
//     password: Password for encrypting the wallet [required]
func walletEncryptHandler(gateway *daemon.Gateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		id := r.FormValue("id")
		if id == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		password := r.FormValue("password")
		if password == "" {
			wh.Error400(w, "missing password")
			return
		}

		wlt, err := gateway.EncryptWallet(id, []byte(password))
		if err != nil {
			switch err.(type) {
			case wallet.ErrWalletNotExist:
				wh.Error400(w, err.Error())
				return
			default:
				logger.Error("Encrypt wallet failed: %v", err)
				wh.Error500(w)
				return
			}
		}

		wh.SendOr404(w, wallet.NewReadableWallet(wlt))
	}
}

// RegisterWalletHandlers registers wallet handlers
func RegisterWalletHandlers(mux *http.ServeMux, gateway *daemon.Gateway) {
	//  Gets a wallet
	//	Method: GET
	//	Args:
	//     id: Wallet ID [required]
	mux.HandleFunc("/wallet", walletGetHandler(gateway))

	// Loads wallet from seed, will scan ahead N address and
	// load addresses till the last one that have coins.
	// Method: POST
	// Args:
	//     seed: Wallet seed [required]
	//     label: Wallet label [required]
	//     scan: The number of addresses to scan ahead for balances [optional, must be > 0]
	// 	   encrypt: Whether encrypt the wallet [optional, must be 0 or 1]
	// 	   password: Password for encrypting wallet [optional, required if the encrypt value is 1]
	mux.HandleFunc("/wallet/create", walletCreateHandler(gateway))

	// Generates new addresses in specific wallet
	// Method: POST
	// Args:
	//   id: Wallet ID [required]
	//     num: The number of address want to generate [required]
	//     password: Wallet password [optional, required if the wallet is encrypted]
	mux.HandleFunc("/wallet/newAddress", walletNewAddressesHandler(gateway))

	// Returns the confirmed and predicted balance for a specific wallet.
	// The predicted balance is the confirmed balance minus any pending
	// spent amount.
	// Method: GET
	// Args:
	//     id: Wallet ID [required]
	mux.HandleFunc("/wallet/balance", walletBalanceHandler(gateway))

	// Creates and broadcasts a transaction sending money from one of our wallets
	// to destination address.
	// Method: POST
	// Args:
	//    id: Wallet ID [required]
	//	  dst: Recipient address [required]
	// 	  coins: The number of droplet you will send [required]
	//	  password: Wallet password [optional, required if the wallet is encrypted]
	mux.HandleFunc("/wallet/spend", walletSpendHandler(gateway))

	// Returns wallet related unconfirmed transactions
	// [Deprecated]:
	//     No one is using this api and the URI is easy to misunderstanding,
	//     this api should returns all transactions that are related to a wallet, while
	//     the implementation actually returns the unconfirmed pending transactions of a wallet.
	// Method: GET
	// Args:
	//     id: Wallet ID [required]
	// Returns all pending transanction for all addresses by selected Wallet
	mux.HandleFunc("/wallet/transactions", walletTransactionsHandler(gateway))

	// Updates wallet label
	// Method: GET
	// Args:
	//     id: Wallet ID [required]
	//     label: Wallet label [required]
	mux.HandleFunc("/wallet/update", walletUpdateLabelHandler(gateway))

	// Returns all loaded wallets
	mux.HandleFunc("/wallets", walletsHandler(gateway))

	// Rescans the wallet directory and loads/unloads wallets based on which
	// files are present. Returns "success" if it works. Otherwise returns
	// 500 status with error message.
	mux.HandleFunc("/wallets/reload", walletsReloadHandler(gateway))

	// Returns the path of wallet directory
	mux.HandleFunc("/wallets/folderName", getWalletFolder(gateway))

	// Generates a new wallet seed
	mux.Handle("/wallet/newSeed", newWalletSeedHandler(gateway))

	// Encrypts wallet
	// Method: GET
	// Args:
	//     id: Wallet ID [required]
	//     password: Password for encrypting wallet [required]
	mux.Handle("/wallet/encrypt", walletEncryptHandler(gateway))
}
