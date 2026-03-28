package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
)

// walletResponse mirrors api.WalletResponse for the thin client
type walletResponse struct {
	Meta     readable.WalletMeta      `json:"meta"`
	Entries  []readable.WalletEntry   `json:"entries"`
	Accounts []readable.WalletAccount `json:"accounts,omitempty"`
}

func newWalletResponse(w wallet.Wallet) (*walletResponse, error) {
	var wr walletResponse

	wr.Meta.Coin = w.Coin()
	wr.Meta.Filename = w.Filename()
	wr.Meta.Label = w.Label()
	wr.Meta.Type = w.Type()
	wr.Meta.Version = w.Version()
	wr.Meta.CryptoType = w.CryptoType()
	wr.Meta.Encrypted = w.IsEncrypted()
	wr.Meta.Timestamp = w.Timestamp()
	wr.Meta.Temp = w.IsTemp()

	var options []wallet.Option
	switch w.Type() {
	case wallet.WalletTypeBip44:
		bip44Coin := w.Bip44Coin()
		if bip44Coin == nil {
			return nil, fmt.Errorf("wallet has no Bip44Coin meta data")
		}
		wr.Meta.Bip44Coin = bip44Coin
		options = append(options, wallet.OptionExternal(), wallet.OptionChange())

		// Populate per-account structure
		accounts := w.Accounts()
		wr.Accounts = make([]readable.WalletAccount, len(accounts))
		for ai, acct := range accounts {
			wa := readable.WalletAccount{
				Name:  acct.Name,
				Index: acct.Index,
			}

			extEntries, err := w.GetEntries(wallet.OptionAccount(acct.Index), wallet.OptionExternal())
			if err != nil {
				return nil, fmt.Errorf("failed to get external entries for account %d: %v", acct.Index, err)
			}
			wa.ExternalEntries = walletEntriesToReadable(extEntries)

			chgEntries, err := w.GetEntries(wallet.OptionAccount(acct.Index), wallet.OptionChange())
			if err != nil {
				return nil, fmt.Errorf("failed to get change entries for account %d: %v", acct.Index, err)
			}
			wa.ChangeEntries = walletEntriesToReadable(chgEntries)

			wr.Accounts[ai] = wa
		}
	case wallet.WalletTypeXPub:
		wr.Meta.XPub = w.XPub()
	}

	entries, err := w.GetEntries(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet entries: %v", err)
	}
	wr.Entries = make([]readable.WalletEntry, len(entries))

	for i, e := range entries {
		wr.Entries[i] = readable.WalletEntry{
			Address: e.Address.String(),
			Public:  e.Public.Hex(),
		}

		switch w.Type() {
		case wallet.WalletTypeBip44:
			childNumber := e.ChildNumber
			wr.Entries[i].ChildNumber = &childNumber
			change := e.Change
			wr.Entries[i].Change = &change
		case wallet.WalletTypeXPub:
			childNumber := e.ChildNumber
			wr.Entries[i].ChildNumber = &childNumber
		}
	}

	return &wr, nil
}

// walletEntriesToReadable converts wallet entries to readable format with child number info
func walletEntriesToReadable(entries wallet.Entries) []readable.WalletEntry {
	result := make([]readable.WalletEntry, len(entries))
	for i, e := range entries {
		childNumber := e.ChildNumber
		change := e.Change
		result[i] = readable.WalletEntry{
			Address:     e.Address.String(),
			Public:      e.Public.Hex(),
			ChildNumber: &childNumber,
			Change:      &change,
		}
	}
	return result
}

// handleWalletAPI handles wallet-related API requests locally when --wallet-dir is set.
// Returns true if the request was handled, false if it should be proxied to the remote node.
func handleWalletAPI(c *gin.Context, apiPath string, wltService *wallet.Service, nodeURL string) bool {
	path := strings.TrimSuffix(apiPath, "/")
	method := c.Request.Method

	switch {
	// Pure wallet management endpoints (served locally)
	case path == "/v1/wallet" && method == http.MethodGet:
		handleGetWallet(c, wltService)
	case path == "/v1/wallet/create" && method == http.MethodPost:
		handleCreateWallet(c, wltService)
	case path == "/v1/wallet/createTemp" && method == http.MethodPost:
		handleCreateTempWallet(c, wltService)
	case path == "/v1/wallet/newAddress" && method == http.MethodPost:
		handleNewAddresses(c, wltService)
	case path == "/v1/wallet/update" && method == http.MethodPost:
		handleUpdateWallet(c, wltService)
	case path == "/v1/wallets" && method == http.MethodGet:
		handleGetWallets(c, wltService)
	case path == "/v1/wallets/folderName" && method == http.MethodGet:
		handleWalletFolder(c, wltService)
	case path == "/v1/wallet/unload" && method == http.MethodPost:
		handleUnloadWallet(c, wltService)
	case path == "/v1/wallet/encrypt" && method == http.MethodPost:
		handleEncryptWallet(c, wltService)
	case path == "/v1/wallet/decrypt" && method == http.MethodPost:
		handleDecryptWallet(c, wltService)
	case path == "/v1/wallet/seed" && method == http.MethodPost:
		handleWalletSeed(c, wltService)
	case path == "/v1/wallet/newSeed" && method == http.MethodGet:
		handleNewSeed(c)
	case path == "/v2/wallet/seed/verify" && method == http.MethodPost:
		handleVerifySeed(c)
	case path == "/v2/wallet/recover" && method == http.MethodPost:
		handleRecoverWallet(c, wltService)
	case path == "/v1/wallet/xpub" && method == http.MethodGet:
		handleWalletXPub(c, wltService)
	case path == "/v1/wallet/newAccount" && method == http.MethodPost:
		handleNewAccount(c, wltService)

	// Transaction signing: sign locally using wallet's stored keys
	case path == "/v2/wallet/transaction/sign" && method == http.MethodPost:
		handleSignTransaction(c, wltService)

	// Hybrid endpoints: resolve wallet addresses locally, query remote node
	case path == "/v1/wallet/balance" && method == http.MethodGet:
		handleWalletBalance(c, wltService, nodeURL)

	default:
		return false
	}
	return true
}

// error response helpers matching the daemon API format
func errBadRequest(c *gin.Context, msg string) {
	c.String(http.StatusBadRequest, "400 Bad Request - %s", msg)
}

func errForbidden(c *gin.Context) {
	c.String(http.StatusForbidden, "403 Forbidden")
}

func errNotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "404 Not Found")
}

func errInternal(c *gin.Context, msg string) {
	c.String(http.StatusInternalServerError, "500 Internal Server Error - %s", msg)
}

func handleWalletError(c *gin.Context, err error) {
	switch err {
	case wallet.ErrWalletNotExist:
		errNotFound(c)
	case wallet.ErrWalletAPIDisabled:
		errForbidden(c)
	default:
		errBadRequest(c, err.Error())
	}
}

// GET /api/v1/wallet
func handleGetWallet(c *gin.Context, s *wallet.Service) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	wlt, err := s.GetWallet(wltID)
	if err != nil {
		handleWalletError(c, err)
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/create
func handleCreateWallet(c *gin.Context, s *wallet.Service) {
	walletType := c.Request.FormValue("type")
	if walletType == "" {
		walletType = wallet.WalletTypeDeterministic
	}

	seed := c.Request.FormValue("seed")
	label := c.Request.FormValue("label")
	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	var encrypt bool
	encryptStr := c.Request.FormValue("encrypt")
	if encryptStr != "" {
		var err error
		encrypt, err = strconv.ParseBool(encryptStr)
		if err != nil {
			errBadRequest(c, fmt.Sprintf("invalid encrypt value: %v", err))
			return
		}
	}

	if encrypt && len(password) == 0 {
		errBadRequest(c, "missing password")
		return
	}

	if !encrypt && len(password) > 0 {
		errBadRequest(c, "encrypt must be true as password is provided")
		return
	}

	var bip44Coin *bip44.CoinType
	bip44CoinStr := c.Request.FormValue("bip44-coin")
	if bip44CoinStr != "" {
		if walletType != wallet.WalletTypeBip44 {
			errBadRequest(c, "bip44-coin is only valid for bip44 type wallets")
			return
		}
		bip44CoinInt, err := strconv.ParseUint(bip44CoinStr, 10, 32)
		if err != nil {
			errBadRequest(c, "invalid bip44-coin value")
			return
		}
		ct := bip44.CoinType(bip44CoinInt)
		bip44Coin = &ct
	}

	secKeys, err := wallet.ParsePrivateKeys(c.Request.FormValue("private-keys"))
	if err != nil {
		errBadRequest(c, "invalid collection private keys")
		return
	}
	defer func() {
		for i := range secKeys {
			secKeys[i] = cipher.SecKey{}
		}
	}()

	// Determine coin type based on segwit preference.
	// The wallet service defaults new BTC BIP44 wallets to segwit.
	// If the user explicitly unchecks segwit, override to legacy.
	var coinType wallet.CoinType
	segwitStr := c.Request.FormValue("segwit")
	if segwitStr == "false" && walletType == wallet.WalletTypeBip44 {
		coinType = wallet.CoinTypeBitcoin
	}

	wlt, err := s.CreateWallet("", wallet.Options{
		Seed:                  seed,
		Label:                 label,
		Encrypt:               encrypt,
		Password:              []byte(password),
		Type:                  walletType,
		Coin:                  coinType,
		SeedPassphrase:        c.Request.FormValue("seed-passphrase"),
		Bip44Coin:             bip44Coin,
		XPub:                  c.Request.FormValue("xpub"),
		CollectionPrivateKeys: secKeys,
	})
	if err != nil {
		switch err.(type) {
		case wallet.Error:
			handleWalletError(c, err)
		default:
			errInternal(c, err.Error())
		}
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/createTemp
func handleCreateTempWallet(c *gin.Context, s *wallet.Service) {
	walletType := c.Request.FormValue("type")
	if walletType == "" {
		walletType = wallet.WalletTypeDeterministic
	}

	seed := c.Request.FormValue("seed")
	label := c.Request.FormValue("label")

	var bip44Coin *bip44.CoinType
	bip44CoinStr := c.Request.FormValue("bip44-coin")
	if bip44CoinStr != "" {
		if walletType != wallet.WalletTypeBip44 {
			errBadRequest(c, "bip44-coin is only valid for bip44 type wallets")
			return
		}
		bip44CoinInt, err := strconv.ParseUint(bip44CoinStr, 10, 32)
		if err != nil {
			errBadRequest(c, "invalid bip44-coin value")
			return
		}
		ct := bip44.CoinType(bip44CoinInt)
		bip44Coin = &ct
	}

	secKeys, err := wallet.ParsePrivateKeys(c.Request.FormValue("private-keys"))
	if err != nil {
		errBadRequest(c, "invalid collection private keys")
		return
	}
	defer func() {
		for i := range secKeys {
			secKeys[i] = cipher.SecKey{}
		}
	}()

	wlt, err := s.CreateWallet("", wallet.Options{
		Temp:                  true,
		Seed:                  seed,
		Label:                 label,
		Type:                  walletType,
		Bip44Coin:             bip44Coin,
		XPub:                  c.Request.FormValue("xpub"),
		CollectionPrivateKeys: secKeys,
	})
	if err != nil {
		switch err.(type) {
		case wallet.Error:
			handleWalletError(c, err)
		default:
			errInternal(c, err.Error())
		}
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/newAddress
func handleNewAddresses(c *gin.Context, s *wallet.Service) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	var opts []wallet.Option
	num := c.Request.FormValue("num")
	if num != "" {
		n, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			errBadRequest(c, "invalid num value")
			return
		}
		opts = append(opts, wallet.OptionGenerateN(n))
	}

	account := c.Request.FormValue("account")
	if account != "" {
		a, err := strconv.ParseUint(account, 10, 32)
		if err != nil {
			errBadRequest(c, "invalid account value")
			return
		}
		opts = append(opts, wallet.OptionAccount(uint32(a))) //nolint:gosec
	}

	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	privateKeys, err := wallet.ParsePrivateKeys(c.Request.FormValue("private-keys"))
	if err != nil {
		errBadRequest(c, "invalid private keys")
		return
	}
	if len(privateKeys) > 0 {
		opts = append(opts, wallet.OptionCollectionPrivateKeys(privateKeys))
	}

	addrs, err := s.NewAddresses(wltID, []byte(password), opts...)
	if err != nil {
		handleWalletError(c, err)
		return
	}

	var rlt = struct {
		Addresses []string `json:"addresses"`
	}{}
	for _, a := range addrs {
		rlt.Addresses = append(rlt.Addresses, a.String())
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/update
func handleUpdateWallet(c *gin.Context, s *wallet.Service) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	label := c.Request.FormValue("label")
	if label == "" {
		errBadRequest(c, "missing label")
		return
	}

	if err := s.UpdateWalletLabel(wltID, label); err != nil {
		handleWalletError(c, err)
		return
	}

	c.JSON(http.StatusOK, "success")
}

// GET /api/v1/wallets
func handleGetWallets(c *gin.Context, s *wallet.Service) {
	wlts, err := s.GetWallets()
	if err != nil {
		handleWalletError(c, err)
		return
	}

	wrs := make([]*walletResponse, 0, len(wlts))
	for _, wlt := range wlts {
		wr, err := newWalletResponse(wlt)
		if err != nil {
			errInternal(c, err.Error())
			return
		}
		wrs = append(wrs, wr)
	}

	sort.Slice(wrs, func(i, j int) bool {
		return wrs[i].Meta.Timestamp < wrs[j].Meta.Timestamp
	})

	c.JSON(http.StatusOK, wrs)
}

// GET /api/v1/wallets/folderName
func handleWalletFolder(c *gin.Context, s *wallet.Service) {
	addr, err := s.WalletDir()
	if err != nil {
		handleWalletError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"address": addr})
}

// POST /api/v1/wallet/unload
func handleUnloadWallet(c *gin.Context, s *wallet.Service) {
	id := c.Request.FormValue("id")
	if id == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	if err := s.UnloadWallet(id); err != nil {
		handleWalletError(c, err)
	}
}

// POST /api/v1/wallet/encrypt
func handleEncryptWallet(c *gin.Context, s *wallet.Service) {
	id := c.Request.FormValue("id")
	if id == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	wlt, err := s.EncryptWallet(id, []byte(password))
	if err != nil {
		switch err {
		case wallet.ErrWalletEncrypted,
			wallet.ErrMissingPassword,
			wallet.ErrEncryptTempWallet,
			wallet.ErrInvalidPassword:
			errBadRequest(c, err.Error())
		case wallet.ErrWalletAPIDisabled:
			errForbidden(c)
		case wallet.ErrWalletNotExist:
			errNotFound(c)
		default:
			errInternal(c, err.Error())
		}
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/decrypt
func handleDecryptWallet(c *gin.Context, s *wallet.Service) {
	id := c.Request.FormValue("id")
	if id == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	wlt, err := s.DecryptWallet(id, []byte(password))
	if err != nil {
		switch err {
		case wallet.ErrMissingPassword,
			wallet.ErrWalletNotEncrypted,
			wallet.ErrInvalidPassword:
			errBadRequest(c, err.Error())
		case wallet.ErrWalletAPIDisabled:
			errForbidden(c)
		case wallet.ErrWalletNotExist:
			errNotFound(c)
		default:
			errInternal(c, err.Error())
		}
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// POST /api/v1/wallet/seed
func handleWalletSeed(c *gin.Context, s *wallet.Service) {
	id := c.Request.FormValue("id")
	if id == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	password := c.Request.FormValue("password")
	defer func() { password = "" }()

	seed, seedPassphrase, err := s.GetWalletSeed(id, []byte(password))
	if err != nil {
		switch err {
		case wallet.ErrMissingPassword,
			wallet.ErrWalletNotEncrypted,
			wallet.ErrInvalidPassword:
			errBadRequest(c, err.Error())
		case wallet.ErrWalletAPIDisabled, wallet.ErrSeedAPIDisabled:
			errForbidden(c)
		case wallet.ErrWalletNotExist:
			errNotFound(c)
		default:
			errInternal(c, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"seed":            seed,
		"seed_passphrase": seedPassphrase,
	})
}

// GET /api/v1/wallet/newSeed
func handleNewSeed(c *gin.Context) {
	entropyValue := c.Request.FormValue("entropy")
	if entropyValue == "" {
		entropyValue = "128"
	}

	entropyBits, err := strconv.Atoi(entropyValue)
	if err != nil {
		errBadRequest(c, "invalid entropy")
		return
	}

	if entropyBits != 128 && entropyBits != 256 {
		errBadRequest(c, "entropy length must be 128 or 256")
		return
	}

	entropy, err := bip39.NewEntropy(entropyBits)
	if err != nil {
		errInternal(c, fmt.Sprintf("bip39.NewEntropy failed: %v", err))
		return
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		errInternal(c, fmt.Sprintf("bip39.NewMnemonic failed: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"seed": mnemonic})
}

// POST /api/v2/wallet/seed/verify
func handleVerifySeed(c *gin.Context) {
	var req struct {
		Seed string `json:"seed"`
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusBadRequest}})
		return
	}

	if req.Seed == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "seed is required", "code": http.StatusBadRequest}})
		return
	}

	if err := bip39.ValidateMnemonic(req.Seed); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusUnprocessableEntity}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": struct{}{}})
}

// POST /api/v2/wallet/recover
func handleRecoverWallet(c *gin.Context, s *wallet.Service) {
	var req struct {
		ID             string `json:"id"`
		Seed           string `json:"seed"`
		SeedPassphrase string `json:"seed_passphrase"`
		Password       string `json:"password"`
	}
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusBadRequest}})
		return
	}

	if req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "id is required", "code": http.StatusBadRequest}})
		return
	}

	if req.Seed == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": "seed is required", "code": http.StatusBadRequest}})
		return
	}

	var password []byte
	if req.Password != "" {
		password = []byte(req.Password)
	}

	defer func() {
		req.Seed = ""
		req.SeedPassphrase = ""
		req.Password = ""
		password = nil
	}()

	wlt, err := s.RecoverWallet(req.ID, req.Seed, req.SeedPassphrase, password)
	if err != nil {
		switch err.(type) {
		case wallet.Error:
			switch err {
			case wallet.ErrWalletNotExist:
				c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"message": "wallet not found", "code": http.StatusNotFound}})
			case wallet.ErrWalletAPIDisabled:
				c.JSON(http.StatusForbidden, gin.H{"error": gin.H{"message": "wallet API disabled", "code": http.StatusForbidden}})
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusBadRequest}})
			}
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusInternalServerError}})
		}
		return
	}

	rlt, err := newWalletResponse(wlt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": err.Error(), "code": http.StatusInternalServerError}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rlt})
}

// GET /api/v1/wallet/xpub
func handleWalletXPub(c *gin.Context, s *wallet.Service) {
	id := c.Request.FormValue("id")
	if id == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	path := c.Request.FormValue("path")
	if path == "" {
		errBadRequest(c, "missing xpub path")
		return
	}

	if _, err := bip44wallet.ParseXPubPath(path); err != nil {
		errBadRequest(c, err.Error())
		return
	}

	wlt, err := s.GetWallet(id)
	if err != nil {
		handleWalletError(c, err)
		return
	}

	if wlt.Type() != wallet.WalletTypeBip44 {
		errBadRequest(c, fmt.Sprintf("%v type wallet does not support xpub key", wlt.Type()))
		return
	}

	bip44Wlt, ok := wlt.(*bip44wallet.Wallet)
	if !ok {
		errBadRequest(c, "invalid bip44 wallet type assertion")
		return
	}

	xpubKey, err := bip44Wlt.GetXPubKey(path)
	if err != nil {
		errBadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"xpub_key": xpubKey})
}

// GET /api/v1/wallet/balance (hybrid: local wallet addresses -> remote balance query)
func handleWalletBalance(c *gin.Context, s *wallet.Service, nodeURL string) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	wlt, err := s.GetWallet(wltID)
	if err != nil {
		handleWalletError(c, err)
		return
	}

	// Get all addresses from the wallet
	var opts []wallet.Option
	if wlt.Type() == wallet.WalletTypeBip44 {
		opts = append(opts, wallet.OptionExternal(), wallet.OptionChange())
	}

	addrs, err := wlt.GetAddresses(opts...)
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to get wallet addresses: %v", err))
		return
	}

	if len(addrs) == 0 {
		// No addresses, return zero balance
		c.JSON(http.StatusOK, gin.H{
			"confirmed": gin.H{"coins": 0, "hours": 0},
			"predicted": gin.H{"coins": 0, "hours": 0},
			"addresses": gin.H{},
		})
		return
	}

	// Build comma-separated address string
	addrStrs := make([]string, len(addrs))
	for i, a := range addrs {
		addrStrs[i] = a.String()
	}
	addrParam := strings.Join(addrStrs, ",")

	// Query remote node for balance
	balanceURL := fmt.Sprintf("%s/api/v1/balance?addrs=%s", nodeURL, addrParam)
	log.Printf("[WALLET] Querying balance for %d addresses from %s", len(addrs), nodeURL)

	resp, err := http.Get(balanceURL) //nolint:gosec
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to query node balance: %v", err))
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing balance response body: %v", cerr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to read balance response: %v", err))
		return
	}

	// Pass through the response from the remote node
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// POST /api/v1/wallet/newAccount
func handleNewAccount(c *gin.Context, s *wallet.Service) {
	wltID := c.Request.FormValue("id")
	if wltID == "" {
		errBadRequest(c, "missing wallet id")
		return
	}

	name := c.Request.FormValue("name")

	var resultWlt wallet.Wallet
	err := s.Update(wltID, func(w wallet.Wallet) error {
		bip44Wlt, ok := w.(*bip44wallet.Wallet)
		if !ok {
			return fmt.Errorf("wallet %s is not a bip44 wallet", wltID)
		}

		acctIndex, err := bip44Wlt.NewAccount(name)
		if err != nil {
			return fmt.Errorf("failed to create account: %v", err)
		}

		// Generate a default address on the new account's external chain
		_, err = bip44Wlt.GenerateAddresses(wallet.OptionAccount(acctIndex), wallet.OptionExternal())
		if err != nil {
			return fmt.Errorf("failed to generate address for new account: %v", err)
		}

		resultWlt = w
		return nil
	})
	if err != nil {
		errInternal(c, err.Error())
		return
	}

	rlt, err := newWalletResponse(resultWlt)
	if err != nil {
		errInternal(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, rlt)
}

// handleSignTransaction signs an unsigned transaction using the local wallet's keys.
// POST /api/v2/wallet/transaction/sign
// JSON body: { "wallet_id": "...", "encoded_transaction": "...", "input_addresses": ["addr1", ...], "password": "..." }
// input_addresses must be in the same order as the transaction inputs.
func handleSignTransaction(c *gin.Context, wltService *wallet.Service) {
	var req struct {
		WalletID           string   `json:"wallet_id"`
		EncodedTransaction string   `json:"encoded_transaction"`
		InputAddresses     []string `json:"input_addresses"`
		Password           string   `json:"password"`
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		errBadRequest(c, "failed to read request body")
		return
	}
	if err := json.Unmarshal(body, &req); err != nil {
		errBadRequest(c, "invalid JSON")
		return
	}

	if req.WalletID == "" {
		errBadRequest(c, "wallet_id is required")
		return
	}
	if req.EncodedTransaction == "" {
		errBadRequest(c, "encoded_transaction is required")
		return
	}

	// Deserialize the unsigned transaction
	txn, err := coin.DeserializeTransactionHex(req.EncodedTransaction)
	if err != nil {
		errBadRequest(c, fmt.Sprintf("failed to deserialize transaction: %v", err))
		return
	}

	if len(req.InputAddresses) != len(txn.In) {
		errBadRequest(c, fmt.Sprintf("input_addresses length (%d) must match transaction inputs length (%d)", len(req.InputAddresses), len(txn.In)))
		return
	}

	// Get the wallet
	wlt, err := wltService.GetWallet(req.WalletID)
	if err != nil {
		if err == wallet.ErrWalletNotExist {
			c.String(http.StatusNotFound, "404 Not Found - wallet not found")
		} else {
			errInternal(c, err.Error())
		}
		return
	}

	// Sign the transaction using wallet entries matched by address
	signFunc := func(w wallet.Wallet) error {
		seckeys := make([]cipher.SecKey, len(txn.In))
		for i, addrStr := range req.InputAddresses {
			addr, err := cipher.DecodeBase58Address(addrStr)
			if err != nil {
				return fmt.Errorf("invalid input address %q: %v", addrStr, err)
			}

			entry, err := w.GetEntry(addr)
			if err != nil {
				return fmt.Errorf("address %s not found in wallet: %v", addrStr, err)
			}

			seckeys[i] = entry.Secret
		}

		txn.SignInputs(seckeys)
		return txn.UpdateHeader()
	}

	if wlt.IsEncrypted() {
		if err := wallet.GuardUpdate(wlt, []byte(req.Password), func(w wallet.Wallet) error {
			return signFunc(w)
		}); err != nil {
			errBadRequest(c, err.Error())
			return
		}
	} else {
		if err := signFunc(wlt); err != nil {
			errBadRequest(c, err.Error())
			return
		}
	}

	signedHex, err := txn.SerializeHex()
	if err != nil {
		errInternal(c, fmt.Sprintf("failed to serialize signed transaction: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"encoded_transaction": signedHex,
		},
	})
}
