package api

// APIs for wallet-related information

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/go-bip39"
	"github.com/skycoin/skycoin/src/daemon"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/util/fee"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/wallet"
)

// SpendResult represents the result of spending
type SpendResult struct {
	Balance     *readable.BalancePair `json:"balance,omitempty"`
	Transaction *readable.Transaction `json:"txn,omitempty"`
	Error       string                `json:"error,omitempty"`
}

// UnconfirmedTxnsResponse contains unconfirmed transaction data
type UnconfirmedTxnsResponse struct {
	Transactions []readable.UnconfirmedTransactions `json:"transactions"`
}

// UnconfirmedTxnsVerboseResponse contains verbose unconfirmed transaction data
type UnconfirmedTxnsVerboseResponse struct {
	Transactions []readable.UnconfirmedTransactionVerbose `json:"transactions"`
}

// BalanceResponse address balance summary struct
type BalanceResponse struct {
	readable.BalancePair
	Addresses readable.AddressBalances `json:"addresses"`
}

// WalletResponse wallet response struct for http apis
type WalletResponse struct {
	Meta    readable.WalletMeta    `json:"meta"`
	Entries []readable.WalletEntry `json:"entries"`
}

// NewWalletResponse creates WalletResponse struct from *wallet.Wallet
func NewWalletResponse(w *wallet.Wallet) (*WalletResponse, error) {
	var wr WalletResponse

	wr.Meta.Coin = w.Meta["coin"]
	wr.Meta.Filename = w.Meta["filename"]
	wr.Meta.Label = w.Meta["label"]
	wr.Meta.Type = w.Meta["type"]
	wr.Meta.Version = w.Meta["version"]
	wr.Meta.CryptoType = w.Meta["cryptoType"]

	// Converts "encrypted" string to boolean if any
	if encryptedStr, ok := w.Meta["encrypted"]; ok {
		encrypted, err := strconv.ParseBool(encryptedStr)
		if err != nil {
			return nil, err
		}
		wr.Meta.Encrypted = encrypted
	}

	if tmStr, ok := w.Meta["tm"]; ok {
		// Converts "tm" string to integer timestamp.
		tm, err := strconv.ParseInt(tmStr, 10, 64)
		if err != nil {
			return nil, err
		}
		wr.Meta.Timestamp = tm
	}

	for _, e := range w.Entries {
		wr.Entries = append(wr.Entries, readable.WalletEntry{
			Address: e.Address.String(),
			Public:  e.Public.Hex(),
		})
	}

	return &wr, nil
}

// Returns the wallet's balance, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
// URI: /api/v1/wallet/balance
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

		walletBalance, addressBalances, err := gateway.GetWalletBalance(wltID)
		if err != nil {
			logger.Errorf("Get wallet balance failed: %v", err)
			switch err {
			case wallet.ErrWalletNotExist:
				wh.Error404(w, "")
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, BalanceResponse{
			BalancePair: readable.NewBalancePair(walletBalance),
			Addresses:   readable.NewAddressBalances(addressBalances),
		})
	}
}

// Returns the balance of one or more addresses, both confirmed and predicted.  The predicted
// balance is the confirmed balance minus the pending spends.
// URI: /api/v1s/balance
// Method: GET, POST
// Args:
//     addrs: command separated list of addresses [required]
func balanceHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		addrsParam := r.FormValue("addrs")
		addrsStr := splitCommaString(addrsParam)

		addrs := make([]cipher.Address, 0, len(addrsStr))
		for _, addr := range addrsStr {
			a, err := cipher.DecodeBase58Address(addr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("address %s is invalid: %v", addr, err))
				return
			}
			addrs = append(addrs, a)
		}

		if len(addrs) == 0 {
			wh.Error400(w, "addrs is required")
			return
		}

		bals, err := gateway.GetBalanceOfAddrs(addrs)
		if err != nil {
			err = fmt.Errorf("gateway.GetBalanceOfAddrs failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		// create map of address to balance
		addressBalances := make(readable.AddressBalances, len(addrs))
		for idx, addr := range addrs {
			addressBalances[addr.String()] = readable.NewBalancePair(bals[idx])
		}

		var balance wallet.BalancePair
		for _, bal := range bals {
			var err error
			balance.Confirmed, err = balance.Confirmed.Add(bal.Confirmed)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			balance.Predicted, err = balance.Predicted.Add(bal.Predicted)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}
		}

		wh.SendJSONOr500(logger, w, BalanceResponse{
			BalancePair: readable.NewBalancePair(balance),
			Addresses:   addressBalances,
		})
	}
}

// Creates and broadcasts a transaction sending money from one of our wallets
// to destination address.
// URI: /api/v1/wallet/spend
// Method: POST
// Args:
//     id: wallet id
//     dst: recipient address
//     coins: the number of droplet you will send
//     password: wallet password
// Response:
//     balance: new balance of the wallet
//     txn: spent transaction
//     error: an error that may have occurred after broadcast the transaction to the network
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

		tx, err := gateway.Spend(wltID, []byte(r.FormValue("password")), coins, dst)
		switch err {
		case nil:
		case fee.ErrTxnNoFee,
			wallet.ErrSpendingUnconfirmed,
			wallet.ErrInsufficientBalance,
			wallet.ErrWalletNotEncrypted,
			wallet.ErrMissingPassword,
			wallet.ErrWalletEncrypted,
			wallet.ErrInvalidPassword:
			wh.Error400(w, err.Error())
			return
		case wallet.ErrWalletAPIDisabled,
			daemon.ErrSpendMethodDisabled:
			wh.Error403(w, "")
			return
		case wallet.ErrWalletNotExist:
			wh.Error404(w, "")
			return
		default:
			wh.Error500(w, err.Error())
			return
		}

		var ret SpendResult

		ret.Transaction, err = readable.NewTransaction(*tx, false)
		if err != nil {
			err = fmt.Errorf("readable.NewTransaction failed: %v", err)
			logger.Error(err)
			ret.Error = err.Error()
			wh.SendJSONOr500(logger, w, ret)
			return
		}

		// Get the new wallet balance
		walletBalance, _, err := gateway.GetWalletBalance(wltID)
		if err != nil {
			err = fmt.Errorf("gateway.GetWalletBalance failed: %v", err)
			logger.Error(err)
			ret.Error = err.Error()
			wh.SendJSONOr500(logger, w, ret)
			return
		}
		b := readable.NewBalancePair(walletBalance)
		ret.Balance = &b

		wh.SendJSONOr500(logger, w, ret)
	}
}

// Loads wallet from seed, will scan ahead N address and
// load addresses till the last one that have coins.
// URI: /api/v1/wallet/create
// Method: POST
// Args:
//     seed: wallet seed [required]
//     label: wallet label [required]
//     scan: the number of addresses to scan ahead for balances [optional, must be > 0]
//     encrypt: bool value, whether encrypt the wallet [optional]
//     password: password for encrypting wallet [optional, must be provided if "encrypt" is set]
func walletCreateHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			wh.Error405(w)
			return
		}

		seed := r.FormValue("seed")
		if seed == "" {
			wh.Error400(w, "missing seed")
			return
		}

		label := r.FormValue("label")
		if label == "" {
			wh.Error400(w, "missing label")
			return
		}

		password := r.FormValue("password")
		defer func() {
			password = ""
		}()

		var encrypt bool
		encryptStr := r.FormValue("encrypt")
		if encryptStr != "" {
			var err error
			encrypt, err = strconv.ParseBool(encryptStr)
			if err != nil {
				wh.Error400(w, fmt.Sprintf("invalid encrypt value: %v", err))
				return
			}
		}

		if encrypt && len(password) == 0 {
			wh.Error400(w, "missing password")
			return
		}

		if !encrypt && len(password) > 0 {
			wh.Error400(w, "encrypt must be true as password is provided")
			return
		}

		scanNStr := r.FormValue("scan")
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
			Seed:     seed,
			Label:    label,
			Encrypt:  encrypt,
			Password: []byte(password),
			ScanN:    scanN,
		})
		if err != nil {
			switch err.(type) {
			case wallet.Error:
				switch err {
				case wallet.ErrWalletAPIDisabled:
					wh.Error403(w, "")
					return
				default:
					wh.Error400(w, err.Error())
					return
				}
			default:
				wh.Error500(w, err.Error())
				return
			}

		}

		rlt, err := NewWalletResponse(wlt)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		wh.SendJSONOr500(logger, w, rlt)
	}
}

// Genreates new addresses
// URI: /api/v1/wallet/newAddress
// Method: POST
// Args:
//     id: wallet id [required]
//     num: number of address need to create [optional, if not set the default value is 1]
//     password: wallet password [optional, must be provided if the wallet is encrypted]
func walletNewAddressesHandler(gateway Gatewayer) http.HandlerFunc {
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

		password := r.FormValue("password")
		defer func() {
			password = ""
		}()

		addrs, err := gateway.NewAddresses(wltID, []byte(password), n)
		if err != nil {
			switch err {
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error400(w, err.Error())
			}
			return
		}

		var rlt = struct {
			Addresses []string `json:"addresses"`
		}{}

		for _, a := range addrs {
			rlt.Addresses = append(rlt.Addresses, a.String())
		}

		wh.SendJSONOr500(logger, w, rlt)
	}
}

// Update wallet label
// URI: /api/v1/wallet/update
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
				wh.Error404(w, "")
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		wh.SendJSONOr500(logger, w, "success")
	}
}

// Returns a wallet by id
// URI: /api/v1/wallet
// Method: GET
// Args:
//     id: wallet id [required]
func walletHandler(gateway Gatewayer) http.HandlerFunc {
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
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error400(w, err.Error())
			}
			return
		}
		rlt, err := NewWalletResponse(wlt)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		wh.SendJSONOr500(logger, w, rlt)
	}
}

// walletTransactionsHandler returns all unconfirmed transactions for all addresses in a given wallet
// URI: /api/v1/wallet/transactions
// Method: GET
// Args:
//	id: wallet id [required]
//	verbose: [bool] include verbose transaction input data
func walletTransactionsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		verbose, err := parseBoolFlag(r.FormValue("verbose"))
		if err != nil {
			wh.Error400(w, "Invalid value for verbose")
			return
		}

		wltID := r.FormValue("id")
		if wltID == "" {
			wh.Error400(w, "missing wallet id")
			return
		}

		handleWalletError := func(err error) {
			switch err {
			case nil:
			case wallet.ErrWalletNotExist:
				wh.Error404(w, "")
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
			}
		}

		if verbose {
			txns, inputs, err := gateway.GetWalletUnconfirmedTransactionsVerbose(wltID)
			if err != nil {
				logger.Errorf("get wallet unconfirmed transactions verbose failed: %v", err)
				handleWalletError(err)
				return
			}

			vb := make([]readable.UnconfirmedTransactionVerbose, len(txns))
			for i, txn := range txns {
				v, err := readable.NewUnconfirmedTransactionVerbose(&txn, inputs[i])
				if err != nil {
					wh.Error500(w, err.Error())
					return
				}
				vb[i] = *v
			}

			wh.SendJSONOr500(logger, w, UnconfirmedTxnsVerboseResponse{
				Transactions: vb,
			})
		} else {
			txns, err := gateway.GetWalletUnconfirmedTransactions(wltID)
			if err != nil {
				logger.Errorf("get wallet unconfirmed transactions failed: %v", err)
				handleWalletError(err)
				return
			}

			unconfirmedTxns, err := readable.NewUnconfirmedTransactions(txns)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wh.SendJSONOr500(logger, w, UnconfirmedTxnsResponse{
				Transactions: unconfirmedTxns,
			})
		}
	}
}

// Returns all loaded wallets
// URI: /api/v1/wallets
// Method: GET
func walletsHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wlts, err := gateway.GetWallets()
		if err != nil {
			switch err {
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		wrs := make([]*WalletResponse, 0, len(wlts))
		for _, wlt := range wlts {
			wr, err := NewWalletResponse(wlt)
			if err != nil {
				wh.Error500(w, err.Error())
				return
			}

			wrs = append(wrs, wr)
		}

		sort.Slice(wrs, func(i, j int) bool {
			return wrs[i].Meta.Timestamp < wrs[j].Meta.Timestamp
		})

		wh.SendJSONOr500(logger, w, wrs)
	}
}

// WalletFolder struct
type WalletFolder struct {
	Address string `json:"address"`
}

// Returns the wallet directory path
// URI: /api/v1/wallets/folderName
// Method: GET
func walletFolderHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr, err := gateway.GetWalletDir()
		if err != nil {
			switch err {
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
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
// URI: /api/v1/wallet/newSeed
// Method: GET
// Args:
//     entropy: entropy bitsize [optional, default value of 128 will be used if not set]
func newSeedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
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
			err = fmt.Errorf("bip39.NewEntropy failed: %v", err)
			wh.Error500(w, err.Error())
			return
		}

		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			err = fmt.Errorf("bip39.NewDefaultMnemonic failed: %v", err)
			wh.Error500(w, err.Error())
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

// Returns seed of wallet of given id
// URI: /api/v1/wallet/seed
// Method: POST
// Args:
//     id: wallet id
//     password: wallet password
func walletSeedHandler(gateway Gatewayer) http.HandlerFunc {
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

		password := r.FormValue("password")
		defer func() {
			password = ""
		}()

		seed, err := gateway.GetWalletSeed(id, []byte(password))
		if err != nil {
			switch err {
			case wallet.ErrMissingPassword,
				wallet.ErrWalletNotEncrypted,
				wallet.ErrInvalidPassword:
				wh.Error400(w, err.Error())
			case wallet.ErrWalletAPIDisabled, wallet.ErrSeedAPIDisabled:
				wh.Error403(w, "")
			case wallet.ErrWalletNotExist:
				wh.Error404(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		v := struct {
			Seed string `json:"seed"`
		}{
			Seed: seed,
		}

		wh.SendJSONOr500(logger, w, v)
	}
}

// Unloads wallet from the wallet service
// URI: /api/v1/wallet/unload
// Method: POST
// Args:
//     id: wallet id
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
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			default:
				wh.Error500(w, err.Error())
			}
		}
	}
}

// Encrypts wallet
// URI: /api/v1/wallet/encrypt
// Method: POST
// Args:
//     id: wallet id
//     password: wallet password
func walletEncryptHandler(gateway Gatewayer) http.HandlerFunc {
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

		password := r.FormValue("password")
		defer func() {
			password = ""
		}()

		wlt, err := gateway.EncryptWallet(id, []byte(password))
		if err != nil {
			switch err {
			case wallet.ErrWalletEncrypted,
				wallet.ErrMissingPassword,
				wallet.ErrInvalidPassword:
				wh.Error400(w, err.Error())
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			case wallet.ErrWalletNotExist:
				wh.Error404(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		// Make sure the sensitive data are wiped
		rlt, err := NewWalletResponse(wlt)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		wh.SendJSONOr500(logger, w, rlt)
	}
}

// Decrypts wallet
// URI: /api/v1/wallet/decrypt
// Method: POST
// Args:
//     id: wallet id
//     password: wallet password
func walletDecryptHandler(gateway Gatewayer) http.HandlerFunc {
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

		password := r.FormValue("password")
		defer func() {
			password = ""
		}()

		wlt, err := gateway.DecryptWallet(id, []byte(password))
		if err != nil {
			switch err {
			case wallet.ErrMissingPassword,
				wallet.ErrWalletNotEncrypted,
				wallet.ErrInvalidPassword:
				wh.Error400(w, err.Error())
			case wallet.ErrWalletAPIDisabled:
				wh.Error403(w, "")
			case wallet.ErrWalletNotExist:
				wh.Error404(w, "")
			default:
				wh.Error500(w, err.Error())
			}
			return
		}

		rlt, err := NewWalletResponse(wlt)
		if err != nil {
			wh.Error500(w, err.Error())
			return
		}
		wh.SendJSONOr500(logger, w, rlt)
	}
}

// WalletRecoverRequest is the request data for POST /api/v2/wallet/recover
type WalletRecoverRequest struct {
	ID       string `json:"id"`
	Seed     string `json:"seed"`
	Password string `json:"password"`
}

// URI: /api/v2/wallet/recover
// Method: POST
// Args:
//	id: wallet id
//  seed: wallet seed
//  password: [optional] new password
// Recovers an encrypted wallet by providing the seed.
// The first address will be generated from seed and compared to the first address
// of the specified wallet. If they match, the wallet will be regenerated
// with an optional password.
// If the wallet is not encrypted, an error is returned.
func walletRecoverHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
			writeHTTPResponse(w, resp)
			return
		}

		if r.Header.Get("Content-Type") != ContentTypeJSON {
			resp := NewHTTPErrorResponse(http.StatusUnsupportedMediaType, "")
			writeHTTPResponse(w, resp)
			return
		}

		var req WalletRecoverRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		if req.ID == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "id is required")
			writeHTTPResponse(w, resp)
			return
		}

		if req.Seed == "" {
			resp := NewHTTPErrorResponse(http.StatusBadRequest, "seed is required")
			writeHTTPResponse(w, resp)
			return
		}

		var password []byte
		if req.Password != "" {
			password = []byte(req.Password)
		}

		defer func() {
			req.Seed = ""
			req.Password = ""
			password = nil
		}()

		wlt, err := gateway.RecoverWallet(req.ID, req.Seed, password)
		if err != nil {
			var resp HTTPResponse
			switch err {
			case wallet.ErrWalletNotEncrypted, wallet.ErrWalletRecoverSeedWrong:
				resp = NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
			case wallet.ErrWalletNotExist:
				resp = NewHTTPErrorResponse(http.StatusNotFound, "")
			case wallet.ErrWalletAPIDisabled:
				resp = NewHTTPErrorResponse(http.StatusForbidden, "")
			default:
				resp = NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			}
			writeHTTPResponse(w, resp)
			return
		}

		rlt, err := NewWalletResponse(wlt)
		if err != nil {
			resp := NewHTTPErrorResponse(http.StatusInternalServerError, err.Error())
			writeHTTPResponse(w, resp)
			return
		}

		writeHTTPResponse(w, HTTPResponse{
			Data: rlt,
		})
	}
}
