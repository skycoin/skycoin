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
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/wallet"
)

func skipWalletIfLive(t *testing.T) bool {
	skip := enabled() && mode(t) == testModeLive && !doLiveWallet(t)
	if skip {
		t.Skip("live wallet tests disabled")
	}
	return skip
}

func TestWalletNewSeed(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if skipWalletIfLive(t) {
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

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()

	mnemonicSeed := bip39.MustNewDefaultMnemonic()

	cases := []struct {
		name           string
		seed           string
		seedPassphrase string
		walletType     string
		encrypted      bool
	}{
		{
			name:       "deterministic encrypted",
			seed:       "fooseed2",
			walletType: wallet.WalletTypeDeterministic,
			encrypted:  true,
		},
		{
			name:       "deterministic unencrypted",
			seed:       "fooseed2",
			walletType: wallet.WalletTypeDeterministic,
		},

		{
			name:           "bip44 with seed passphrase encrypted",
			seed:           mnemonicSeed,
			seedPassphrase: "foobar",
			walletType:     wallet.WalletTypeBip44,
			encrypted:      true,
		},
		{
			name:           "bip44 without seed passphrase encrypted",
			seed:           mnemonicSeed,
			seedPassphrase: "",
			walletType:     wallet.WalletTypeBip44,
			encrypted:      true,
		},
		{
			name:           "bip44 with seed passphrase unencrypted",
			seed:           mnemonicSeed,
			seedPassphrase: "foobar",
			walletType:     wallet.WalletTypeBip44,
		},
		{
			name:           "bip44 without seed passphrase unencrypted",
			seed:           mnemonicSeed,
			seedPassphrase: "",
			walletType:     wallet.WalletTypeBip44,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pwd := ""
			if tc.encrypted {
				pwd = "pwd"
			}

			w, _, clean := createWallet(t, c, tc.encrypted, pwd, tc.seed, tc.seedPassphrase, tc.walletType)
			defer clean()
			require.Equal(t, tc.encrypted, w.Meta.Encrypted)

			walletDir := getWalletDir(t, c)

			// Confirms the wallet does exist
			walletPath := filepath.Join(walletDir, w.Meta.Filename)
			_, err := os.Stat(walletPath)
			require.NoError(t, err)

			// Loads the wallet and confirms that the wallet has the same seed
			lw, err := wallet.Load(walletPath)
			require.NoError(t, err)
			require.Equal(t, len(w.Entries), lw.EntriesLen())
			require.Equal(t, tc.walletType, lw.Type())

			if tc.encrypted {
				require.True(t, lw.IsEncrypted())
				require.Empty(t, lw.Seed())
				require.Empty(t, lw.SeedPassphrase())
			} else {
				require.False(t, lw.IsEncrypted())
				require.Equal(t, tc.seed, lw.Seed())
				require.Equal(t, tc.seedPassphrase, lw.SeedPassphrase())
			}

			for i := range w.Entries {
				require.Equal(t, w.Entries[i].Address, lw.GetEntryAt(i).Address.String())
				require.Equal(t, w.Entries[i].Public, lw.GetEntryAt(i).Public.Hex())

				if tc.encrypted {
					require.True(t, lw.GetEntryAt(i).Secret.Null())
				} else {
					require.False(t, lw.GetEntryAt(i).Secret.Null())
				}

				switch tc.walletType {
				case wallet.WalletTypeBip44:
					require.NotNil(t, w.Entries[i].ChildNumber)
					require.Equal(t, uint32(i), *w.Entries[i].ChildNumber)
					require.NotNil(t, w.Entries[i].Change)
					require.Equal(t, bip44.ExternalChainIndex, *w.Entries[i].Change)
				default:
					require.Nil(t, w.Entries[i].ChildNumber)
					require.Nil(t, w.Entries[i].Change)
				}
			}
		})
	}
}

func TestGetWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()

	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			// Create a wallet
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
			defer clean()

			// Confirms the wallet can be acquired
			w1, err := c.Wallet(w.Meta.Filename)
			require.NoError(t, err)
			require.Equal(t, *w, *w1)
		})
	}
}

func TestGetWallets(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()

	// Creates 2 new wallets of each type
	var ws []api.WalletResponse
	for i := 0; i < 2; i++ {
		for _, walletType := range createWalletTypes {
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
			defer clean()
			// cleaners = append(cleaners, clean)
			ws = append(ws, *w)
		}
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

	if skipWalletIfLive(t) {
		return
	}

	seed := bip39.MustNewDefaultMnemonic()

	cases := []struct {
		name           string
		seed           string
		seedPassphrase string
		walletType     string
	}{
		{
			name:       "deterministic",
			seed:       seed,
			walletType: wallet.WalletTypeDeterministic,
		},
		{
			name:       "bip44 without seed passphrase",
			seed:       seed,
			walletType: wallet.WalletTypeBip44,
		},
		{
			name:           "bip44 with seed passphrase",
			seed:           seed,
			seedPassphrase: "foobar",
			walletType:     wallet.WalletTypeBip44,
		},
	}

	// We only test 30 cases, because the more addresses we generate, the longer
	// it takes, we don't want to spend much time here.
	for _, tc := range cases {
		for i := 1; i <= 30; i++ {
			name := fmt.Sprintf("%s generate %v addresses", tc.name, i)
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

				w, _, clean := createWallet(t, c, encrypt, password, tc.seed, tc.seedPassphrase, tc.walletType)
				defer clean()

				addrs, err := c.NewWalletAddress(w.Meta.Filename, i, password)
				if err != nil {
					t.Fatalf("%v", err)
					return
				}
				require.NoError(t, err)

				switch tc.walletType {
				case wallet.WalletTypeDeterministic:
					seckeys := cipher.MustGenerateDeterministicKeyPairs([]byte(tc.seed), i+1)
					var as []string
					for _, k := range seckeys {
						as = append(as, cipher.MustAddressFromSecKey(k).String())
					}

					// Confirms that the new generated addresses match
					require.Equal(t, len(addrs), len(as)-1)
					for i := range addrs {
						require.Equal(t, as[i+1], addrs[i])
					}
				case wallet.WalletTypeBip44:
					ss, err := bip39.NewSeed(tc.seed, tc.seedPassphrase)
					require.NoError(t, err)

					cc, err := bip44.NewCoin(ss, bip44.CoinTypeSkycoin)
					require.NoError(t, err)

					acct, err := cc.Account(0)
					require.NoError(t, err)

					ext, err := acct.External()
					require.NoError(t, err)

					var as []string
					for j := uint32(0); j < uint32(i+1); j++ {
						k, err := ext.NewPrivateChildKey(j)
						require.NoError(t, err)
						sk := cipher.MustNewSecKey(k.Key)
						as = append(as, cipher.MustAddressFromSecKey(sk).String())
					}

					// Confirms that the new generated addresses match
					require.Equal(t, len(addrs), len(as)-1)
					for i := range addrs {
						require.Equal(t, as[i+1], addrs[i])
					}
				default:
					t.Fatalf("unhandled wallet type %q", tc.walletType)
				}
			})
		}
	}
}

func TestStableWalletBalance(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			seed := "casino away claim road artist where blossom warrior demise royal still palm"
			w, _, clean := createWallet(t, c, false, "", seed, "", walletType)
			defer clean()

			bp, err := c.WalletBalance(w.Meta.Filename)
			require.NoError(t, err)

			var expect api.BalanceResponse
			checkGoldenFile(t, fmt.Sprintf("wallet-balance-%s.golden", walletType), TestData{*bp, &expect})
		})
	}
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

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()
	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
			defer clean()

			err := c.UpdateWallet(w.Meta.Filename, "new wallet")
			require.NoError(t, err)

			// Confirms the wallet has label of "new wallet"
			w1, err := c.Wallet(w.Meta.Filename)
			require.NoError(t, err)
			require.Equal(t, w1.Meta.Label, "new wallet")
		})
	}
}

func TestStableWalletUnconfirmedTransactions(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()
	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
			defer clean()

			txns, err := c.WalletUnconfirmedTransactions(w.Meta.Filename)
			require.NoError(t, err)

			goldenFile := fmt.Sprintf("wallet-%s-transactions.golden", walletType)
			var expect api.UnconfirmedTxnsResponse
			checkGoldenFile(t, goldenFile, TestData{*txns, &expect})
		})
	}
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
	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
			defer clean()

			txns, err := c.WalletUnconfirmedTransactionsVerbose(w.Meta.Filename)
			require.NoError(t, err)

			goldenFile := fmt.Sprintf("wallet-%s-transactions-verbose.golden", walletType)
			var expect api.UnconfirmedTxnsVerboseResponse
			checkGoldenFile(t, goldenFile, TestData{*txns, &expect})
		})
	}
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

	if skipWalletIfLive(t) {
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

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()

	// Create a unencrypted wallet
	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			w, _, clean := createWallet(t, c, false, "", "", "", walletType)
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
		})
	}
}

func TestDecryptWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if skipWalletIfLive(t) {
		return
	}

	c := newClient()

	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			w, seed, clean := createWallet(t, c, true, "pwd", "", "", walletType)
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

			switch walletType {
			case wallet.WalletTypeDeterministic:
				// Confirms the last seed matches
				lseed, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 1)
				require.Equal(t, hex.EncodeToString(lseed), lw.LastSeed())

				// Confirms that the first address is derived from the private key
				pubkey := cipher.MustPubKeyFromSecKey(seckeys[0])
				require.Equal(t, w.Entries[0].Address, cipher.AddressFromPubKey(pubkey).String())
				require.Equal(t, lw.GetEntryAt(0).Address.String(), w.Entries[0].Address)
			case wallet.WalletTypeBip44:
				require.Empty(t, lw.LastSeed())
			default:
				t.Fatalf("unhandled wallet type %q", walletType)
			}
		})
	}
}

func TestRecoverWallet(t *testing.T) {
	if !doLiveOrStable(t) {
		return
	}

	if skipWalletIfLive(t) {
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

	cases := []struct {
		name           string
		seed           string
		seedPassphrase string
		badSeed        string
		walletType     string
	}{
		{
			name:       "deterministic",
			seed:       "fooseed",
			badSeed:    "fooseed2",
			walletType: wallet.WalletTypeDeterministic,
		},
		{
			name:           "bip44 with seed passphrase",
			seed:           "voyage say extend find sheriff surge priority merit ignore maple cash argue",
			seedPassphrase: "foobar",
			badSeed:        "mule seed lady practice desk length roast tongue attract heavy spirit focus",
			walletType:     wallet.WalletTypeBip44,
		},
		{
			name:           "bip44 without seed passphrase",
			seed:           "voyage say extend find sheriff surge priority merit ignore maple cash argue",
			seedPassphrase: "",
			badSeed:        "mule seed lady practice desk length roast tongue attract heavy spirit focus",
			walletType:     wallet.WalletTypeBip44,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, _, clean := createWallet(t, c, false, "", tc.seed, tc.seedPassphrase, tc.walletType)
			defer clean()

			_, err = c.NewWalletAddress(w.Meta.Filename, 10, "")
			require.NoError(t, err)

			w, err = c.Wallet(w.Meta.Filename)
			require.NoError(t, err)

			// Recover fails if the wallet is not encrypted
			_, err = c.RecoverWallet(api.WalletRecoverRequest{
				ID:             w.Meta.Filename,
				Seed:           tc.seed,
				SeedPassphrase: tc.seedPassphrase,
			})
			assertResponseError(t, err, http.StatusBadRequest, "wallet is not encrypted")

			_, err = c.EncryptWallet(w.Meta.Filename, "pwd")
			require.NoError(t, err)

			// Recovery fails if the seed doesn't match
			_, err = c.RecoverWallet(api.WalletRecoverRequest{
				ID:             w.Meta.Filename,
				Seed:           tc.badSeed,
				SeedPassphrase: tc.seedPassphrase,
			})
			assertResponseError(t, err, http.StatusBadRequest, "wallet recovery seed or seed passphrase is wrong")

			// Recovery fails if the seed passphrase doesn't match
			_, err = c.RecoverWallet(api.WalletRecoverRequest{
				ID:             w.Meta.Filename,
				Seed:           tc.seed,
				SeedPassphrase: tc.seedPassphrase + "2",
			})

			switch tc.walletType {
			case wallet.WalletTypeBip44:
				assertResponseError(t, err, http.StatusBadRequest, "wallet recovery seed or seed passphrase is wrong")
			case wallet.WalletTypeDeterministic:
				assertResponseError(t, err, http.StatusBadRequest, "RecoverWallet failed to create temporary wallet for fingerprint comparison: seedPassphrase is only used for \"bip44\" wallets")
			default:
				t.Fatalf("unhandled wallet type %q", tc.walletType)
			}

			// Successful recovery with no new password
			w2, err := c.RecoverWallet(api.WalletRecoverRequest{
				ID:             w.Meta.Filename,
				Seed:           tc.seed,
				SeedPassphrase: tc.seedPassphrase,
			})
			require.NoError(t, err)
			require.False(t, w2.Meta.Encrypted)
			checkWalletOnDisk(w2)
			require.Equal(t, w, w2)

			_, err = c.EncryptWallet(w.Meta.Filename, "pwd2")
			require.NoError(t, err)

			// Successful recovery with a new password
			w3, err := c.RecoverWallet(api.WalletRecoverRequest{
				ID:             w.Meta.Filename,
				Seed:           tc.seed,
				SeedPassphrase: tc.seedPassphrase,
				Password:       "pwd3",
			})
			require.NoError(t, err)
			require.True(t, w3.Meta.Encrypted)
			require.Equal(t, wallet.CryptoTypeScryptChacha20poly1305, w3.Meta.CryptoType)
			checkWalletOnDisk(w3)
			w3.Meta.Encrypted = w.Meta.Encrypted
			w3.Meta.CryptoType = w.Meta.CryptoType
			require.Equal(t, w, w3)

			w4, err := c.DecryptWallet(w.Meta.Filename, "pwd3")
			require.NoError(t, err)
			require.False(t, w.Meta.Encrypted)
			require.Equal(t, w, w4)
		})
	}
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

	for _, walletType := range createWalletTypes {
		t.Run(walletType, func(t *testing.T) {
			// Create an encrypted wallet
			w, _, clean := createWallet(t, c, true, "pwd", "", "", walletType)
			defer clean()

			_, err := c.WalletSeed(w.Meta.Filename, "pwd")
			assertResponseError(t, err, http.StatusForbidden, "403 Forbidden - Endpoint is disabled")
		})
	}
}

func TestGetWalletSeedEnabledAPI(t *testing.T) {
	if !doEnableSeedAPI(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name           string
		walletType     string
		seed1          string
		seed2          string
		seedPassphrase string
	}{
		{
			name:       "deterministic",
			walletType: wallet.WalletTypeDeterministic,
			seed1:      bip39.MustNewDefaultMnemonic(),
			seed2:      bip39.MustNewDefaultMnemonic(),
		},
		{
			name:       "bip44 without seed passphrase",
			walletType: wallet.WalletTypeBip44,
			seed1:      bip39.MustNewDefaultMnemonic(),
			seed2:      bip39.MustNewDefaultMnemonic(),
		},
		{
			name:           "bip44 with seed passphrase",
			walletType:     wallet.WalletTypeBip44,
			seed1:          bip39.MustNewDefaultMnemonic(),
			seed2:          bip39.MustNewDefaultMnemonic(),
			seedPassphrase: "foobar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.seed1)
			require.NotEmpty(t, tc.seed2)

			// Create an encrypted wallet
			w, _, clean := createWallet(t, c, true, "pwd", tc.seed1, tc.seedPassphrase, tc.walletType)
			defer clean()

			resp, err := c.WalletSeed(w.Meta.Filename, "pwd")
			require.NoError(t, err)

			// Confirms the seed are matched
			require.Equal(t, tc.seed1, resp.Seed)
			require.Equal(t, tc.seedPassphrase, resp.SeedPassphrase)

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
			nw, _, nclean := createWallet(t, c, false, "", tc.seed2, tc.seedPassphrase, tc.walletType)
			defer nclean()
			_, err = c.WalletSeed(nw.Meta.Filename, "pwd")
			assertResponseError(t, err, http.StatusBadRequest, "400 Bad Request - wallet is not encrypted")
		})
	}
}

// prepareAndCheckWallet gets wallet from environment, and confirms:
// 1. The minimal coins and coin hours requirements are met.
// 2. The wallet has at least two address entry.
// Returns the loaded wallet, total coins, total coin hours and password of the wallet.
func prepareAndCheckWallet(t *testing.T, c *api.Client, minCoins, minCoinHours uint64) (wallet.Wallet, uint64, uint64, string) {
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
	if coins < minCoins {
		t.Fatalf("Wallet must have at least %d coins", minCoins)
	}

	if hours < minCoinHours {
		t.Fatalf("Wallet must have at least %d coin hours", minCoinHours)
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

var createWalletTypes = []string{wallet.WalletTypeDeterministic, wallet.WalletTypeBip44}

func createWallet(t *testing.T, c *api.Client, encrypt bool, password, seed, seedPassphrase, walletType string) (*api.WalletResponse, string, func()) {
	if seed == "" {
		seed = bip39.MustNewDefaultMnemonic()
	}
	// Use the first 6 letters of the seed as the label
	var w *api.WalletResponse
	var err error
	if encrypt {
		w, err = c.CreateEncryptedWallet(walletType, seed, seedPassphrase, seed[:6], password, 0)
	} else {
		w, err = c.CreateUnencryptedWallet(walletType, seed, seedPassphrase, seed[:6], 0)
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
