package api

// APIs for wallet-related information

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/readable"
	wh "github.com/skycoin/skycoin/src/util/http"
	"github.com/skycoin/skycoin/src/wallet"
)

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

// NewWalletResponse creates WalletResponse struct from wallet.Wallet
func NewWalletResponse(w wallet.Wallet) (*WalletResponse, error) {
	var wr WalletResponse

	wr.Meta.Coin = w.Coin()
	wr.Meta.Filename = w.Filename()
	wr.Meta.Label = w.Label()
	wr.Meta.Type = w.Type()
	wr.Meta.Version = w.Version()
	wr.Meta.CryptoType = w.CryptoType()
	wr.Meta.Encrypted = w.IsEncrypted()
	wr.Meta.Timestamp = w.Timestamp()

	for _, e := range w.GetEntries() {
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
		addrs, err := parseAddressesFromStr(addrsParam)
		if err != nil {
			wh.Error400(w, err.Error())
			return
		}

		if len(addrs) == 0 {
			wh.Error400(w, "addrs is required")
			return
		}

		bals, err := gateway.GetBalanceOfAddresses(addrs)
		if err != nil {
			err = fmt.Errorf("gateway.GetBalanceOfAddresses failed: %v", err)
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
		}, gateway)
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
func walletFolderHandler(s Walleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		addr, err := s.WalletDir()
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

// VerifySeedRequest is the request data for POST /api/v2/wallet/seed/verify
type VerifySeedRequest struct {
	Seed string `json:"seed"`
}

// walletVerifySeedHandler verifies a wallet seed
// Method: POST
// URI: /api/v2/wallet/seed/verify
func walletVerifySeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		resp := NewHTTPErrorResponse(http.StatusMethodNotAllowed, "")
		writeHTTPResponse(w, resp)
		return
	}

	var req VerifySeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	if req.Seed == "" {
		resp := NewHTTPErrorResponse(http.StatusBadRequest, "seed is required")
		writeHTTPResponse(w, resp)
		return
	}

	if err := bip39.ValidateMnemonic(req.Seed); err != nil {
		resp := NewHTTPErrorResponse(http.StatusUnprocessableEntity, err.Error())
		writeHTTPResponse(w, resp)
		return
	}

	writeHTTPResponse(w, HTTPResponse{Data: struct{}{}})
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
