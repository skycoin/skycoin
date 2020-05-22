package bip44wallet

import (
	"errors"
	"fmt"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/stretchr/testify/require"
)

var (
	skycoinExternalAddrs = skycoinAddressStringsToAddress(testSkycoinExternalAddresses)
	bitcoinExternalAddrs = bitcoinAddressStringsToAddress(testBitcoinExternalAddresses)
)

type mockTxnsFinder map[cipher.Addresser]bool

func (mb mockTxnsFinder) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	if len(addrs) == 0 {
		return nil, nil
	}
	active := make([]bool, len(addrs))
	for i, addr := range addrs {
		active[i] = mb[addr]
	}
	return active, nil
}

func TestBip44NewWallet(t *testing.T) {
	bip44SkycoinType := bip44.CoinTypeSkycoin
	newBip44Type := bip44.CoinType(1000)

	type expect struct {
		coinType      wallet.CoinType
		bip44CoinType bip44.CoinType
		entriesLen    int
		accountName   string
	}

	tt := []struct {
		name           string
		filename       string
		label          string
		seed           string
		opts           []wallet.Option
		seedPassphrase string
		expect         expect
		err            error
	}{
		{
			name:           "skycoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "bitcoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeBitcoin),
			},
			expect: expect{
				wallet.CoinTypeBitcoin,
				bip44.CoinTypeBitcoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin explicit bip44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
				wallet.OptionBip44Coin(&bip44SkycoinType),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin geneateN=5",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionGenerateN(5),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				5,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin geneateN < scanN entriesLen=generateN",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionGenerateN(3),
				wallet.OptionScanN(4),
				wallet.OptionTransactionsFinder(mockTxnsFinder{
					skycoinExternalAddrs[0]: true,
				}),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				3,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin geneateN < scanN entriesLen>generateN",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionGenerateN(3),
				wallet.OptionScanN(4),
				wallet.OptionTransactionsFinder(mockTxnsFinder{
					skycoinExternalAddrs[3]: true,
				}),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				4,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin scanN=2 entriesLen=2",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionScanN(2),
				wallet.OptionTransactionsFinder(mockTxnsFinder{
					skycoinExternalAddrs[1]: true,
				}),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				2,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin scanN=1 entriesLen=1 has txn",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionScanN(1),
				wallet.OptionTransactionsFinder(mockTxnsFinder{
					skycoinExternalAddrs[0]: true,
				}),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin scanN=1 entriesLen=1 no txn",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionScanN(1),
				wallet.OptionTransactionsFinder(mockTxnsFinder{}),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin set default account name",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionDefaultBip44AccountName("marketing"),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				"marketing",
			},
		},
		{
			name:           "skycoin new bip44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
				wallet.OptionBip44Coin(&newBip44Type),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				newBip44Type,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "no filename",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: errors.New("filename not set"),
		},
		{
			name:           "no coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			expect: expect{
				wallet.CoinTypeSkycoin,
				bip44.CoinTypeSkycoin,
				1,
				DefaultAccountName,
			},
		},
		{
			name:           "skycoin empty seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           "",
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: wallet.ErrMissingSeed,
		},
		{
			name:           "skycoin invalid seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           invalidBip44Seed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
			},
			err: errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
		{
			name:           "new coin type, no bi44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType("unknown"),
			},
			err: errors.New("bip44 coin type not set"),
		},
	}

	for _, tc := range tt {
		for _, encrypt := range []bool{false, true} {
			for _, ct := range crypto.TypesInsecure() {
				name := fmt.Sprintf("%s crypto=%v encrypt-%v", tc.name, ct, encrypt)
				opts := tc.opts
				if encrypt {
					opts = append(opts, wallet.OptionEncrypt(true))
					opts = append(opts, wallet.OptionPassword([]byte("pwd")))
					opts = append(opts, wallet.OptionCryptoType(ct))
				}

				t.Run(name, func(t *testing.T) {
					w, err := NewWallet(tc.filename, tc.label, tc.seed, tc.seedPassphrase, opts...)
					require.Equal(t, tc.err, err, fmt.Sprintf("want: %v got: %v", tc.err, err))
					if err != nil {
						return
					}
					require.Equal(t, tc.filename, w.Meta.Filename())
					require.Equal(t, tc.label, w.Meta.Label())
					require.Equal(t, WalletType, w.Meta.Type())
					require.NotEmpty(t, w.Meta.Timestamp())
					require.NotNil(t, w.decoder)
					bip44Coin := w.Bip44Coin()
					require.Equal(t, tc.expect.bip44CoinType, *bip44Coin)
					require.Equal(t, tc.expect.coinType, w.Meta.Coin())

					if encrypt {
						require.Equal(t, ct, w.Meta.CryptoType())
						// confirms that seeds and entry secrets are all empty
						checkNoSensitiveData(t, w)
						return
					}

					require.Equal(t, tc.seed, w.Meta.Seed())
					require.Equal(t, tc.seedPassphrase, w.Meta.SeedPassphrase())
					require.False(t, w.Meta.IsEncrypted())
					require.Empty(t, w.Meta.Secrets())
					el, err := w.EntriesLen()
					require.NoError(t, err)
					require.Equal(t, tc.expect.entriesLen, el)
					a, err := w.account(0)
					require.NoError(t, err)
					require.Equal(t, tc.expect.accountName, a.Name)
				})
			}
		}
	}
}

func TestBip44NewWalletDefaultCrypto(t *testing.T) {
	for _, coinType := range []wallet.CoinType{wallet.CoinTypeSkycoin} {
		//for _, coinType := range []wallet.CoinType{wallet.CoinTypeSkycoin, wallet.CoinTypeBitcoin} {
		for _, encrypt := range []bool{true} {
			name := fmt.Sprintf("coinType %v encrypt %v", coinType, encrypt)
			t.Run(name, func(t *testing.T) {
				opts := []wallet.Option{
					wallet.OptionCoinType(coinType),
				}
				if encrypt {
					opts = append(opts, wallet.OptionEncrypt(true))
					opts = append(opts, wallet.OptionPassword([]byte("pwd")))
				}
				w, err := NewWallet(
					"test.wlt",
					"test",
					testSeed,
					testSeedPassphrase,
					opts...,
				)
				require.NoError(t, err)
				require.Equal(t, "test.wlt", w.Meta.Filename())
				require.Equal(t, "test", w.Meta.Label())
				require.Equal(t, WalletType, w.Meta.Type())
				require.NotEmpty(t, w.Meta.Timestamp())
				require.NotNil(t, w.decoder)
				bip44Coin := w.Bip44Coin()
				switch coinType {
				case wallet.CoinTypeSkycoin:
					require.Equal(t, bip44.CoinTypeSkycoin, *bip44Coin)
				case wallet.CoinTypeBitcoin:
					require.Equal(t, bip44.CoinTypeBitcoin, *bip44Coin)
				}
				require.Equal(t, coinType, w.Meta.Coin())

				if encrypt {
					require.Equal(t, crypto.DefaultCryptoType, w.Meta.CryptoType())
					// confirms that seeds and entry secrets are all empty
					checkNoSensitiveData(t, w)
					return
				}

				require.Equal(t, testSeed, w.Meta.Seed())
				require.Equal(t, testSeedPassphrase, w.Meta.SeedPassphrase())
				require.False(t, w.Meta.IsEncrypted())
				require.Empty(t, w.Meta.Secrets())
			})
		}
	}
}

func checkNoSensitiveData(t *testing.T, w *Wallet) {
	// confirms that seeds and entry secrets are all empty
	require.True(t, w.IsEncrypted())
	require.Equal(t, w.Seed(), "")
	require.Equal(t, w.SeedPassphrase(), "")
	for _, a := range w.Accounts() {
		for _, c := range []uint32{bip44.ExternalChainIndex, bip44.ChangeChainIndex} {
			entries, err := w.entries(a.Index, c)
			require.NoError(t, err)

			// confirms account secrets are empty
			require.Empty(t, w.accountManager.(*bip44Accounts).accounts[a.Index].PrivateKey)

			// confirms no secrets in the entries
			for _, e := range entries {
				require.Equal(t, cipher.SecKey{}, e.Secret)
			}
		}
	}

	require.NotEmpty(t, w.Meta.Secrets())
	return
}

func TestWalletLock(t *testing.T) {
	tt := []struct {
		name    string
		wltName string
		opts    []wallet.Option
		lockPwd []byte
		err     error
	}{
		{
			name:    "ok",
			lockPwd: []byte("pwd"),
		},
		{
			name:    "password is nil",
			lockPwd: nil,
			err:     wallet.ErrMissingPassword,
		},
		{
			name: "wallet already encrypted",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			lockPwd: []byte("pwd"),
			err:     wallet.ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			opts := tc.opts
			opts = append(opts, wallet.OptionCryptoType(ct))

			t.Run(name, func(t *testing.T) {
				w, err := NewWallet("test.wlt", "test", testSeed, testSeedPassphrase, opts...)
				require.NoError(t, err)

				// create a default account
				_, err = w.NewAccount("account1")
				require.NoError(t, err)

				if !w.IsEncrypted() {
					// Generates 2 addresses
					_, err = w.GenerateAddresses(2)
					require.NoError(t, err)
				}

				err = w.Lock(tc.lockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				checkNoSensitiveData(t, w)
				//require.True(t, w.IsEncrypted())
				//
				//// Checks if the seeds are wiped
				//require.Empty(t, w.Seed())
				//require.Empty(t, w.LastSeed())
				//
				//// Checks if the entries are encrypted
				//entries, err := w.GetEntries()
				//require.NoError(t, err)
				//
				//for _, e := range entries {
				//	require.Equal(t, cipher.SecKey{}, e.Secret)
				//}
			})

		}
	}
}

func TestWalletCreateAccount(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	ai, err = w.NewAccount("account2")
	require.Equal(t, uint32(1), ai)

	require.Equal(t, uint32(2), w.accountManager.len())
}

func TestWalletAccountCreateAddresses(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	addrs, err := w.newExternalAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr := make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinExternalAddresses[:2], addrsStr)

	addrs, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr = make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinChangeAddresses[:2], addrsStr)
}

func TestBip44WalletLock(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	err = w.Lock([]byte("123456"))
	require.NoError(t, err)

	require.Empty(t, w.Seed())
	require.Empty(t, w.SeedPassphrase())
	require.NotEmpty(t, w.Secrets())
	require.True(t, w.IsEncrypted())

	// confirms that no secrets exist in the accounts
	ss := make(wallet.Secrets)
	w.accountManager.packSecrets(ss)
	require.Equal(t, 4, len(ss))
	for k, v := range ss {
		if k == secretBip44AccountPrivateKey {
			require.Empty(t, v)
		} else {
			require.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", v)
		}
	}
}

// - Test wallet unlock
func TestBip44WalletUnlock(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin),
		wallet.OptionCryptoType(crypto.CryptoTypeScryptChacha20poly1305Insecure))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	cw := w.Clone().(*Wallet)

	err = cw.Lock([]byte("123456"))
	require.NoError(t, err)

	// generates addresses after locking
	_, err = cw.newExternalAddresses(ai, 2)
	require.NoError(t, err)
	_, err = cw.newChangeAddresses(ai, 3)
	require.NoError(t, err)

	// unlock with wrong password
	_, err = cw.Unlock([]byte("12345"))
	require.Equal(t, errors.New("invalid password"), err)

	// unlock with correct password
	wlt, err := cw.Unlock([]byte("123456"))
	require.NoError(t, err)

	el, err := wlt.EntriesLen(wallet.OptionAccount(ai), wallet.OptionChange(false))
	require.NoError(t, err)
	require.Equal(t, 4, el)

	cl, err := wlt.EntriesLen(wallet.OptionAccount(ai), wallet.OptionChange(true))
	require.NoError(t, err)
	require.Equal(t, 5, cl)

	// confirms that unlocking wallet won't lose data
	require.Empty(t, wlt.Secrets())
	require.False(t, wlt.IsEncrypted())
	require.Equal(t, w.Seed(), wlt.Seed())
	require.Equal(t, w.SeedPassphrase(), wlt.SeedPassphrase())

	// pack the origin wallet's secrets
	originSS := make(wallet.Secrets)
	w.accountManager.packSecrets(originSS)

	// pack the unlocked wallet's secrets
	ss := make(wallet.Secrets)
	wlt.(*Wallet).accountManager.packSecrets(ss)

	// compare these two secrets, they should have the same keys and values
	require.Equal(t, len(originSS)+5, len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}

func TestBip44WalletNewSerializeDeserialize(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.newExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.newChangeAddresses(ai, 2)
	require.NoError(t, err)

	b, err := w.Serialize()
	require.NoError(t, err)
	t.Log(string(b))

	wlt := Wallet{}
	err = wlt.Deserialize(b)
	require.NoError(t, err)

	// Confirms that serialize/deserialize do not lose meta data
	require.Equal(t, len(w.Meta), len(wlt.Meta))
	for k, v := range wlt.Meta {
		vv, ok := w.Meta[k]
		require.Truef(t, ok, "key:%s", k)
		require.Equal(t, v, vv)
	}

	// confirms that serialize/deserialize do not lose accounts data
	require.Equal(t, w.accountManager.len(), wlt.accountManager.len())
	originSS := make(wallet.Secrets)
	ss := make(wallet.Secrets)
	w.accountManager.packSecrets(originSS)
	wlt.accountManager.packSecrets(ss)

	require.Equal(t, len(originSS), len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}

func skycoinAddressStringsToAddress(addrsStr []string) []cipher.Addresser {
	var addrs []cipher.Addresser
	for _, addr := range addrsStr {
		a := cipher.MustDecodeBase58Address(addr)
		addrs = append(addrs, a)
	}

	return addrs
}

func bitcoinAddressStringsToAddress(addrsStr []string) []cipher.Addresser {
	var addrs []cipher.Addresser
	for _, addr := range addrsStr {
		a := cipher.MustDecodeBase58BitcoinAddress(addr)
		addrs = append(addrs, a)
	}
	return addrs
}
