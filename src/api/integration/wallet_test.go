package integration_test

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestWalletNewSeed(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	cases := []struct {
		name     string
		entropy  int
		numWords int
		errCode  int
		errMsg   string
	}{
		{
			name:     "entropy 128",
			entropy:  128,
			numWords: 12,
		},
		{
			name:     "entropy 256",
			entropy:  256,
			numWords: 24,
		},
		{
			name:    "entropy 100",
			entropy: 100,
			errCode: http.StatusBadRequest,
			errMsg:  "400 Bad Request - entropy length must be 128 or 256",
		},
	}

	c := newClient()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seed, err := c.NewSeed(tc.entropy)
			if tc.errMsg != "" {
				assertResponseError(t, err, tc.errCode, tc.errMsg)
				return
			}

			require.NoError(t, err)
			words := strings.Split(seed, " ")
			require.Len(t, words, tc.numWords)

			// no extra whitespace on the seed
			require.Equal(t, seed, strings.TrimSpace(seed))

			// should generate a different seed each time
			seed2, err := c.NewSeed(tc.entropy)
			require.NoError(t, err)
			require.NotEqual(t, seed, seed2)
		})
	}
}

func TestCreateWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()

	w, seed, clean := createWallet(t, c, false, "", "")
	defer clean()
	require.False(t, w.Meta.Encrypted)

	walletDir := getWalletDir(t, c)

	// Confirms the wallet does exist
	walletPath := filepath.Join(walletDir, w.Meta.Filename)
	_, err := os.Stat(walletPath)
	require.NoError(t, err)

	// Loads the wallet and confirms that the wallet has the same seed
	lw, err := wallet.Load(walletPath)
	require.NoError(t, err)
	require.False(t, lw.IsEncrypted())
	require.Equal(t, seed, lw.Seed())
	require.Equal(t, len(w.Entries), lw.EntriesLen())

	for i := range w.Entries {
		require.Equal(t, w.Entries[i].Address, lw.GetEntryAt(i).Address.String())
		require.Equal(t, w.Entries[i].Public, lw.GetEntryAt(i).Public.Hex())
	}

	// Creates wallet with encryption
	encW, _, encWClean := createWallet(t, c, true, "pwd", "")
	defer encWClean()
	require.True(t, encW.Meta.Encrypted)

	walletPath = filepath.Join(walletDir, encW.Meta.Filename)
	encLW, err := wallet.Load(walletPath)
	require.NoError(t, err)

	// Confirms the loaded wallet is encrypted and has the same address entries
	require.True(t, encLW.IsEncrypted())
	require.Equal(t, len(encW.Entries), encLW.EntriesLen())

	for i := range encW.Entries {
		require.Equal(t, encW.Entries[i].Address, encLW.GetEntryAt(i).Address.String())
		require.Equal(t, encW.Entries[i].Public, encLW.GetEntryAt(i).Public.Hex())
	}
}

func TestGetWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()

	// Create a wallet
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	// Confirms the wallet can be acquired
	w1, err := c.Wallet(w.Meta.Filename)
	require.NoError(t, err)
	require.Equal(t, *w, *w1)
}

func TestGetWallets(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()

	// Creates 2 new wallets
	var ws []api.WalletResponse
	for i := 0; i < 2; i++ {
		w, _, clean := createWallet(t, c, false, "", "")
		defer clean()
		// cleaners = append(cleaners, clean)
		ws = append(ws, *w)
	}

	// Gets wallet from node
	wlts, err := c.Wallets()
	require.NoError(t, err)

	// Create the wallet map
	walletMap := make(map[string]api.WalletResponse)
	for _, w := range wlts {
		walletMap[w.Meta.Filename] = w
	}

	// Confirms the returned wallets contains the wallet we created.
	for _, w := range ws {
		retW, ok := walletMap[w.Meta.Filename]
		require.True(t, ok)
		require.Equal(t, w, retW)
	}
}

// TestWalletNewAddress will generate 30 wallets for testing, and they will
// be removed automatically after testing.
func TestWalletNewAddress(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	// We only test 30 cases, cause the more addresses we generate, the longer
	// it takes, we don't want to spend much time here.
	for i := 1; i <= 30; i++ {
		name := fmt.Sprintf("generate %v addresses", i)
		t.Run(name, func(t *testing.T) {
			c := newClient()
			var encrypt bool
			var password string
			// Test wallet with encryption only when i == 2, so that
			// the tests won't time out.
			if i == 2 {
				encrypt = true
				password = "pwd"
			}

			w, seed, clean := createWallet(t, c, encrypt, password, "")
			defer clean()

			addrs, err := c.NewWalletAddress(w.Meta.Filename, i, password)
			if err != nil {
				t.Fatalf("%v", err)
				return
			}
			require.NoError(t, err)

			seckeys := cipher.MustGenerateDeterministicKeyPairs([]byte(seed), i+1)
			var as []string
			for _, k := range seckeys {
				as = append(as, cipher.MustAddressFromSecKey(k).String())
			}

			// Confirms thoses new generated addresses are the same.
			require.Equal(t, len(addrs), len(as)-1)
			for i := range addrs {
				require.Equal(t, as[i+1], addrs[i])
			}
		})
	}
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	w, _, clean := createWallet(t, c, false, "", "casino away claim road artist where blossom warrior demise royal still palm")
	defer clean()

	bp, err := c.WalletBalance(w.Meta.Filename)
	require.NoError(t, err)

	var expect api.BalanceResponse
	checkGoldenFile(t, "wallet-balance.golden", TestData{*bp, &expect})
}

func TestLiveWalletBalance(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()
	_, walletName, _ := getWalletFromEnv(t, c)
	bp, err := c.WalletBalance(walletName)
	require.NoError(t, err)
	require.NotNil(t, bp)
	require.NotNil(t, bp.Addresses)
}

func TestWalletUpdate(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	err := c.UpdateWallet(w.Meta.Filename, "new wallet")
	require.NoError(t, err)

	// Confirms the wallet has label of "new wallet"
	w1, err := c.Wallet(w.Meta.Filename)
	require.NoError(t, err)
	require.Equal(t, w1.Meta.Label, "new wallet")
}

func TestStableWalletUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	txns, err := c.WalletUnconfirmedTransactions(w.Meta.Filename)
	require.NoError(t, err)

	var expect api.UnconfirmedTxnsResponse
	checkGoldenFile(t, "wallet-transactions.golden", TestData{*txns, &expect})
}

func TestLiveWalletUnconfirmedTransactions(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()
	w, _, _, _ := prepareAndCheckWallet(t, c, 1e6, 1)
	txns, err := c.WalletUnconfirmedTransactions(w.Filename())
	require.NoError(t, err)

	bp, err := c.WalletBalance(w.Filename())
	require.NoError(t, err)
	// There's pending transactions if predicted coins are not the same as confirmed coins
	if bp.Predicted.Coins != bp.Confirmed.Coins {
		require.NotEmpty(t, txns.Transactions)
		return
	}

	require.Empty(t, txns.Transactions)
}

func TestStableWalletUnconfirmedTransactionsVerbose(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	txns, err := c.WalletUnconfirmedTransactionsVerbose(w.Meta.Filename)
	require.NoError(t, err)

	var expect api.UnconfirmedTxnsVerboseResponse
	checkGoldenFile(t, "wallet-transactions-verbose.golden", TestData{*txns, &expect})
}

func TestLiveWalletUnconfirmedTransactionsVerbose(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()
	w, _, _, _ := prepareAndCheckWallet(t, c, 1e6, 1)
	txns, err := c.WalletUnconfirmedTransactionsVerbose(w.Filename())
	require.NoError(t, err)

	bp, err := c.WalletBalance(w.Filename())
	require.NoError(t, err)
	// There's pending transactions if predicted coins are not the same as confirmed coins
	if bp.Predicted.Coins != bp.Confirmed.Coins {
		require.NotEmpty(t, txns.Transactions)
		return
	}

	require.Empty(t, txns.Transactions)
}

func TestWalletFolderName(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()
	folderName, err := c.WalletFolderName()
	require.NoError(t, err)

	require.NotNil(t, folderName)
	require.NotEmpty(t, folderName.Address)
}

func TestEncryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()

	// Create a unencrypted wallet
	w, _, clean := createWallet(t, c, false, "", "")
	defer clean()

	// Encrypts the wallet
	rlt, err := c.EncryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)
	require.NotEmpty(t, rlt.Meta.CryptoType)
	require.True(t, rlt.Meta.Encrypted)

	//  Encrypt the wallet again, should returns error
	_, err = c.EncryptWallet(w.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - wallet is encrypted")

	// Confirms that no sensitive data do exist in wallet file
	wf, err := c.WalletFolderName()
	require.NoError(t, err)
	wltPath := filepath.Join(wf.Address, w.Meta.Filename)
	lw, err := wallet.Load(wltPath)
	require.NoError(t, err)
	require.Empty(t, lw.Seed())
	require.Empty(t, lw.LastSeed())
	require.NotEmpty(t, lw.Secrets())

	// Decrypts the wallet, and confirms that the
	// seed and address entries are the same as it was before being encrypted.
	dw, err := c.DecryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)
	require.Equal(t, w, dw)
}

func TestDecryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	c := newClient()
	w, seed, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	// Decrypt wallet with different password, must fail
	_, err := c.DecryptWallet(w.Meta.Filename, "pwd1")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - invalid password")

	// Decrypt wallet with no password, must fail
	_, err = c.DecryptWallet(w.Meta.Filename, "")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - missing password")

	// Decrypts wallet with correct password
	dw, err := c.DecryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)

	// Confirms that no sensitive data are returned
	require.Empty(t, dw.Meta.CryptoType)
	require.False(t, dw.Meta.Encrypted)

	// Loads wallet from file
	wf, err := c.WalletFolderName()
	require.NoError(t, err)
	wltPath := filepath.Join(wf.Address, w.Meta.Filename)
	lw, err := wallet.Load(wltPath)
	require.NoError(t, err)

	require.Equal(t, lw.Seed(), seed)
	require.Equal(t, 1, lw.EntriesLen())

	// Confirms the last seed is matched
	lseed, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 1)
	require.Equal(t, hex.EncodeToString(lseed), lw.LastSeed())

	// Confirms that the first address is derivied from the private key
	pubkey := cipher.MustPubKeyFromSecKey(seckeys[0])
	require.Equal(t, w.Entries[0].Address, cipher.AddressFromPubKey(pubkey).String())
	require.Equal(t, lw.GetEntryAt(0).Address.String(), w.Entries[0].Address)
}

func TestRecoverWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if doLive(t) && !doLiveWallet(t) {
		return
	}

	// Create an encrypted wallet with some addresses pregenerated,
	// to make sure recover recovers the same number of addresses
	c := newClient()
	wf, err := c.WalletFolderName()
	require.NoError(t, err)

	// Load the wallet from disk to check that it was saved
	checkWalletOnDisk := func(w *api.WalletResponse) {
		wltPath := filepath.Join(wf.Address, w.Meta.Filename)
		lw, err := wallet.Load(wltPath)
		require.NoError(t, err)
		lwr, err := api.NewWalletResponse(lw)
		require.NoError(t, err)
		require.Equal(t, w, lwr)
	}

	w, seed, clean := createWallet(t, c, false, "", "fooseed")
	require.Equal(t, "fooseed", seed)
	defer clean()

	_, err = c.NewWalletAddress(w.Meta.Filename, 10, "")
	require.NoError(t, err)

	w, err = c.Wallet(w.Meta.Filename)
	require.NoError(t, err)

	// Recover fails if the wallet is not encrypted
	_, err = c.RecoverWallet(w.Meta.Filename, "fooseed", "")
	assertResponseError(t, err, http.StatusBadRequest, "wallet is not encrypted")

	_, err = c.EncryptWallet(w.Meta.Filename, "pwd")
	require.NoError(t, err)

	// Recovery fails if the seed doesn't match
	_, err = c.RecoverWallet(w.Meta.Filename, "wrongseed", "")
	assertResponseError(t, err, http.StatusBadRequest, "wallet recovery seed is wrong")

	// Successful recovery with no new password
	w2, err := c.RecoverWallet(w.Meta.Filename, "fooseed", "")
	require.NoError(t, err)
	require.False(t, w2.Meta.Encrypted)
	checkWalletOnDisk(w2)
	require.Equal(t, w, w2)

	_, err = c.EncryptWallet(w.Meta.Filename, "pwd2")
	require.NoError(t, err)

	// Successful recovery with a new password
	w3, err := c.RecoverWallet(w.Meta.Filename, "fooseed", "pwd3")
	require.NoError(t, err)
	require.True(t, w3.Meta.Encrypted)
	require.Equal(t, w3.Meta.CryptoType, "scrypt-chacha20poly1305")
	checkWalletOnDisk(w3)
	w3.Meta.Encrypted = w.Meta.Encrypted
	w3.Meta.CryptoType = w.Meta.CryptoType
	require.Equal(t, w, w3)

	w4, err := c.DecryptWallet(w.Meta.Filename, "pwd3")
	require.NoError(t, err)
	require.False(t, w.Meta.Encrypted)
	require.Equal(t, w, w4)
}

func TestVerifyWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}
	c := newClient()

	// check with correct seed
	isValid, err := c.VerifySeed("nut wife logic sample addict shop before tobacco crisp bleak lawsuit affair")
	require.NoError(t, err)
	require.True(t, isValid)

	// check with incorrect seed
	isValid, err = c.VerifySeed("nut ")
	require.False(t, isValid)
	assertResponseError(t, err, http.StatusUnprocessableEntity, bip39.ErrSurroundingWhitespace.Error())
}

func TestGetWalletSeedDisabledAPI(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if mode(t) == testModeLive && !*testLiveWallet {
		t.Skip("Skipping tests because live mode enabled but wallet tests disabled")
	}

	if mode(t) == testModeEnableSeedAPI {
		t.Skip("Skipping because enable seed API tests is on")
	}

	c := newClient()

	// Create an encrypted wallet
	w, _, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	_, err := c.WalletSeed(w.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusForbidden, "403 Forbidden - Endpoint is disabled")
}

func TestGetWalletSeedEnabledAPI(t *testing.T) {
	if !doEnableSeedAPI(t) {
		return
	}

	c := newClient()

	// Create an encrypted wallet
	w, seed, clean := createWallet(t, c, true, "pwd", "")
	defer clean()

	require.NotEmpty(t, seed)

	sd, err := c.WalletSeed(w.Meta.Filename, "pwd")
	require.NoError(t, err)

	// Confirms the seed are matched
	require.Equal(t, seed, sd)

	// Get seed of wrong wallet id
	_, err = c.WalletSeed("w.wlt", "pwd")
	assertResponseError(t, err, http.StatusNotFound, "404 Not Found")

	// Check with invalid password
	_, err = c.WalletSeed(w.Meta.Filename, "wrong password")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - invalid password")

	// Check with missing password
	_, err = c.WalletSeed(w.Meta.Filename, "")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - missing password")

	// Create unencrypted wallet to check against
	nw, _, nclean := createWallet(t, c, false, "", "")
	defer nclean()
	_, err = c.WalletSeed(nw.Meta.Filename, "pwd")
	assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - wallet is not encrypted")
}

// prepareAndCheckWallet gets wallet from environment, and confirms:
// 1. The minimal coins and coin hours requirements are met.
// 2. The wallet has at least two address entry.
// Returns the loaded wallet, total coins, total coin hours and password of the wallet.
func prepareAndCheckWallet(t *testing.T, c *api.Client, miniCoins, miniCoinHours uint64) (wallet.Wallet, uint64, uint64, string) {
	walletDir, walletName, password := getWalletFromEnv(t, c)
	walletPath := filepath.Join(walletDir, walletName)

	// Checks if the wallet does exist
	if _, err := os.Stat(walletPath); os.IsNotExist(err) {
		t.Fatalf("Wallet %v doesn't exist", walletPath)
	}

	w, err := wallet.Load(walletPath)
	if err != nil {
		t.Fatalf("Load wallet %v failed: %v", walletPath, err)
	}

	if w.IsEncrypted() && password == "" {
		t.Fatalf("Wallet is encrypted, must set WALLET_PASSWORD env var")
	}

	// Generate more addresses if address entries less than 2.
	if w.EntriesLen() < 2 {
		_, err := c.NewWalletAddress(w.Filename(), 2-w.EntriesLen(), password)
		if err != nil {
			t.Fatalf("New wallet address failed: %v", err)
		}

		w, err = wallet.Load(walletPath)
		if err != nil {
			t.Fatalf("Reload wallet %v failed: %v", walletPath, err)
		}
	}

	coins, hours := getWalletBalance(t, c, walletName)
	if coins < miniCoins {
		t.Fatalf("Wallet must have at least %d coins", miniCoins)
	}

	if hours < miniCoinHours {
		t.Fatalf("Wallet must have at least %d coin hours", miniCoinHours)
	}

	if err := wallet.Save(w, walletDir); err != nil {
		t.Fatalf("%v", err)
	}

	return w, coins, hours, password
}

// getWalletFromEnv loads wallet from environment variables.
// Returns wallet dir, wallet name and wallet password is any.
func getWalletFromEnv(t *testing.T, c *api.Client) (string, string, string) {
	walletDir := getWalletDir(t, c)

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("Missing WALLET_NAME environment value")
	}

	walletPassword := os.Getenv("WALLET_PASSWORD")
	return walletDir, walletName, walletPassword
}

func requireWalletEnv(t *testing.T) {
	if !doLiveWallet(t) {
		return
	}

	walletName := os.Getenv("WALLET_NAME")
	if walletName == "" {
		t.Fatal("missing WALLET_NAME environment value")
	}
}

// getWalletBalance gets wallet balance.
// Returns coins and hours
func getWalletBalance(t *testing.T, c *api.Client, walletName string) (uint64, uint64) {
	wp, err := c.WalletBalance(walletName)
	if err != nil {
		t.Fatalf("Get wallet balance of %v failed: %v", walletName, err)
	}

	return wp.Confirmed.Coins, wp.Confirmed.Hours
}

func getWalletDir(t *testing.T, c *api.Client) string {
	wf, err := c.WalletFolderName()
	if err != nil {
		t.Fatalf("%v", err)
	}
	return wf.Address
}

// createWallet creates a wallet with rand seed.
// Returns the generated wallet, seed and clean up function.
func createWallet(t *testing.T, c *api.Client, encrypt bool, password string, seed string) (*api.WalletResponse, string, func()) {
	if seed == "" {
		seed = hex.EncodeToString(cipher.RandByte(32))
	}
	// Use the first 6 letter of the seed as label.
	var w *api.WalletResponse
	var err error
	if encrypt {
		w, err = c.CreateEncryptedWallet(seed, seed[:6], password, 0)
	} else {
		w, err = c.CreateUnencryptedWallet(seed, seed[:6], 0)
	}

	require.NoError(t, err)

	walletDir := getWalletDir(t, c)

	return w, seed, func() {
		// Cleaner function to delete the wallet and bak wallet
		walletPath := filepath.Join(walletDir, w.Meta.Filename)
		err = os.Remove(walletPath)
		require.NoError(t, err)

		bakWalletPath := walletPath + ".bak"
		if _, err := os.Stat(bakWalletPath); !os.IsNotExist(err) {
			// Return directly if no .bak file does exist
			err = os.Remove(bakWalletPath)
			require.NoError(t, err)
		}

		require.NoError(t, err)

		// Removes the wallet from memory
		err = c.UnloadWallet(w.Meta.Filename)
		require.NoError(t, err)
	}
}
