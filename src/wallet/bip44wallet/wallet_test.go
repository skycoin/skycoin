package bip44wallet

import (
	"errors"
	"fmt"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/require"
)

var (
	skycoinExternalAddrs = skycoinAddressStringsToAddress(testSkycoinExternalAddresses)
	skycoinChangeAddrs   = skycoinAddressStringsToAddress(testSkycoinChangeAddresses)
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
		isTemp        bool
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
				2, // 1 external, 1 change
				DefaultAccountName,
				false,
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
				2, // 1 external, 1 change,
				DefaultAccountName,
				false,
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
				2, // one external, 1 change,
				DefaultAccountName,
				false,
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
				6, // 5 external, 1 change
				DefaultAccountName,
				false,
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
				4, // 3 external, 1 change
				DefaultAccountName,
				false,
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
				5, // 4 external, 1 change
				DefaultAccountName,
				false,
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
				3, // 2 external, one change
				DefaultAccountName,
				false,
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
				2, // 1 external, 1 change
				DefaultAccountName,
				false,
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
				2, // 1 external, 1 change
				DefaultAccountName,
				false,
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
				2, // 1 external, 1 change
				"marketing",
				false,
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
				2, // 1 external, 1 change
				DefaultAccountName,
				false,
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
				2, // 1 external, 1 change
				DefaultAccountName,
				false,
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
		{
			name:           "temp wallet",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeSkycoin),
				wallet.OptionBip44Coin(&newBip44Type),
				wallet.OptionTemp(true),
			},
			expect: expect{
				wallet.CoinTypeSkycoin,
				newBip44Type,
				2, // 1 external, 1 change
				DefaultAccountName,
				true,
			},
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

				if !w.IsEncrypted() {
					// Generates 2 addresses
					_, err = w.GenerateAddresses(wallet.OptionGenerateN(2))
					require.NoError(t, err)

					_, err = w.GenerateAddresses(wallet.OptionGenerateN(2), wallet.OptionChange())
					require.NoError(t, err)
				}

				err = w.Lock(tc.lockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				checkNoSensitiveData(t, w)
			})
		}
	}
}

// - Test wallet unlock
func TestWalletUnlock(t *testing.T) {
	tt := []struct {
		name string
		opts []wallet.Option
		pwd  []byte
		err  error
	}{
		{
			name: "ok bip44 wallet",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
				wallet.OptionGenerateN(5),
			},
			pwd: []byte("pwd"),
		},
		{
			name: "password is nil",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			pwd: nil,
			err: wallet.ErrMissingPassword,
		},
		{
			name: "wrong password",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			pwd: []byte("wrong_pwd"),
			err: wallet.ErrInvalidPassword,
		},
		{
			name: "wallet not encrypted",
			pwd:  []byte("pwd"),
			err:  wallet.ErrWalletNotEncrypted,
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

				cw := w.Clone()

				// Unlock the wallet
				unlockWlt, err := cw.Unlock(tc.pwd)
				require.Equal(t, tc.err, err, fmt.Sprintf("want: %v get: %v", tc.err, err))
				if err != nil {
					return
				}

				require.False(t, unlockWlt.IsEncrypted())
				require.Empty(t, unlockWlt.Secrets())

				// Checks the seeds and seed passphrase
				require.Equal(t, testSeed, unlockWlt.Seed())
				require.Equal(t, testSeedPassphrase, unlockWlt.SeedPassphrase())

				// Checks the generated external addresses
				el, err := unlockWlt.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 6, el)

				// pack the origin wallet's secrets
				originSS := make(wallet.Secrets)
				w.accountManager.packSecrets(originSS)

				// pack the unlocked wallet's secrets
				ss := make(wallet.Secrets)
				unlockWlt.(*Wallet).accountManager.packSecrets(ss)

				// compare these two secrets, they should have the same keys and values
				// len(ss) - 1, to remove the private account key
				require.Equal(t, len(originSS), len(ss)-1)
				for k := range originSS {
					vv, ok := ss[k]
					require.True(t, ok)

					addr, err := cipher.DecodeBase58Address(k)
					require.NoError(t, err)
					sk, err := cipher.SecKeyFromHex(vv)
					require.NoError(t, err)
					genAddr, err := cipher.AddressFromSecKey(sk)
					require.NoError(t, err)
					require.Equal(t, addr, genAddr)
				}
			})
		}
	}
}

func TestLockAndUnLock(t *testing.T) {
	for _, ct := range crypto.TypesInsecure() {
		w, err := NewWallet("wallet.wlt", "test", testSeed, testSeedPassphrase, wallet.OptionCryptoType(ct))
		require.NoError(t, err)
		_, err = w.GenerateAddresses(wallet.OptionGenerateN(9))
		require.NoError(t, err)
		el, err := w.EntriesLen()
		require.NoError(t, err)
		// 1 default address + 9
		require.Equal(t, 11, el)

		// clone the wallet
		cw := w.Clone()

		// lock the cloned wallet
		err = cw.Lock([]byte("pwd"))
		require.NoError(t, err)

		checkNoSensitiveData(t, cw.(*Wallet))

		// unlock the cloned wallet
		ucw, err := cw.Unlock([]byte("pwd"))
		require.NoError(t, err)

		// set the account and decoder to nil
		w.accountManager = nil
		ucw.(*Wallet).accountManager = nil
		w.decoder = nil
		ucw.(*Wallet).decoder = nil
		require.Equal(t, w, ucw)
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
	require.Equal(t, uint32(1), ai)

	ai, err = w.NewAccount("account2")
	require.Equal(t, uint32(2), ai)

	require.Equal(t, uint32(3), w.accountManager.len())
}

func TestWalletAccountCreateAddresses(t *testing.T) {
	w, err := NewWallet(
		"test.wlt",
		"test",
		testSeed,
		testSeedPassphrase,
		wallet.OptionCoinType(wallet.CoinTypeSkycoin))
	require.NoError(t, err)

	addrs, err := w.newExternalAddresses(0, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr := make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinExternalAddresses[1:3], addrsStr)

	addrs, err = w.newChangeAddresses(0, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr = make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinChangeAddresses[1:3], addrsStr)
}

func TestWalletGenerateAddress(t *testing.T) {
	tt := []struct {
		name               string
		opts               []wallet.Option
		num                uint64
		expectNum          uint64
		oneAddressEachTime bool
		err                error
	}{
		{
			name:      "ok with none address",
			num:       0,
			expectNum: 2, // 1 external, 1 change
		},
		{
			name:      "ok with one address",
			num:       1,
			expectNum: 3, // 2 external, 1 change
		},
		{
			name:      "ok with two address",
			num:       2,
			expectNum: 4, // 3 external, 1 change
		},
		{
			name:               "ok with three address and generate one address each time deterministic",
			num:                2,
			oneAddressEachTime: true,
			expectNum:          4, // 3 external, 1 change
		},
		{
			name: "encrypt wallet",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			num:                2,
			oneAddressEachTime: true,
			expectNum:          4, // 3 external, 1 change
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("%v crypto=%v num=%d", tc.name, ct, tc.num)
			opts := tc.opts
			opts = append(opts, wallet.OptionCryptoType(ct))

			t.Run(name, func(t *testing.T) {
				// create wallet
				w, err := NewWallet("test.wlt", "test", testSeed, testSeedPassphrase, opts...)
				require.NoError(t, err)

				// generate address
				if !tc.oneAddressEachTime {
					_, err := w.GenerateAddresses(wallet.OptionGenerateN(tc.num))
					require.NoError(t, err)
					if err != nil {
						return
					}
				} else {
					for i := uint64(0); i < tc.num; i++ {
						_, err := w.GenerateAddresses(wallet.OptionGenerateN(1))
						require.Equal(t, tc.err, err)
						if err != nil {
							return
						}
					}
				}

				// check the entry number
				l, err := w.EntriesLen()
				require.NoError(t, err)
				// 1 default address + tc.num = wallet.EntriesLen()
				require.Equal(t, int(tc.expectNum), l)

				addrs, err := w.GetAddresses()
				require.NoError(t, err)

				expectAddrs := make([]cipher.Addresser, tc.num+2)
				copy(expectAddrs, skycoinExternalAddrs[:tc.num+1])
				copy(expectAddrs[tc.num+1:], skycoinChangeAddrs[:1])
				require.Equal(t, expectAddrs, addrs)
			})
		}
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

func TestPeekChangeAddress(t *testing.T) {
	w, err := NewWallet("test.wlt", "test", testSeed, testSeedPassphrase)
	require.NoError(t, err)

	addr, err := w.PeekChangeAddress(mockTxnsFinder{})
	require.NoError(t, err)

	require.Equal(t, skycoinChangeAddrs[0], addr)

	addr, err = w.PeekChangeAddress(mockTxnsFinder{skycoinChangeAddrs[0]: true})
	require.NoError(t, err)
	require.Equal(t, skycoinChangeAddrs[1], addr)

	addr, err = w.PeekChangeAddress(mockTxnsFinder{skycoinChangeAddrs[1]: true})
	require.NoError(t, err)
	require.Equal(t, skycoinChangeAddrs[2], addr)
}

func TestScanAddresses(t *testing.T) {
	eAddrs := skycoinExternalAddrs
	cAddrs := skycoinChangeAddrs

	tt := []struct {
		name                 string
		scanN                uint32
		txnFinder            wallet.TransactionsFinder
		expectAddrs          []cipher.Addresser
		expectAllChangeAddrs []cipher.Addresser
		err                  error
	}{
		{
			name:                 "no txns",
			scanN:                10,
			txnFinder:            mockTxnsFinder{},
			expectAllChangeAddrs: cAddrs[:1],
		},
		{
			name:                 "external addr with txn",
			scanN:                10,
			txnFinder:            mockTxnsFinder{eAddrs[1]: true},
			expectAddrs:          eAddrs[1:2],
			expectAllChangeAddrs: cAddrs[:1],
		},
		{
			name:      "change addr with txn",
			scanN:     10,
			txnFinder: mockTxnsFinder{cAddrs[1]: true},
			// The default change address already exist, thus no more new change addresses will be created
			expectAllChangeAddrs: cAddrs[0:2],
		},
		{
			name:                 "external and change addrs with txns",
			scanN:                10,
			txnFinder:            mockTxnsFinder{eAddrs[1]: true, cAddrs[1]: true},
			expectAddrs:          []cipher.Addresser{eAddrs[1]},
			expectAllChangeAddrs: cAddrs[:2],
		},
		{
			name: "external and change addrs with txns 2", scanN: 10,
			txnFinder:            mockTxnsFinder{eAddrs[2]: true, cAddrs[1]: true},
			expectAddrs:          eAddrs[1:3],
			expectAllChangeAddrs: cAddrs[:2],
		},
		{
			name:  "external and change addrs with txns 3",
			scanN: 10,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs:          eAddrs[1:5],
			expectAllChangeAddrs: cAddrs[:5],
		},
		{
			name:  "not enough addresses scanned",
			scanN: 3,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAllChangeAddrs: cAddrs[:1],
		},
		{
			name:  "just enough addresses scanned",
			scanN: 4,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs:          eAddrs[1:5],
			expectAllChangeAddrs: cAddrs[:5],
		},
		{
			name:  "more addresses scanned",
			scanN: 6,
			txnFinder: mockTxnsFinder{
				eAddrs[4]: true,
				cAddrs[4]: true,
			},
			expectAddrs:          eAddrs[1:5],
			expectAllChangeAddrs: cAddrs[:5],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewWallet("test.wlt", "test", testSeed, testSeedPassphrase, wallet.OptionTransactionsFinder(tc.txnFinder))
			require.NoError(t, err)

			addrs, err := w.ScanAddresses(uint64(tc.scanN), tc.txnFinder)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expectAddrs, addrs)

			// get the change address, as the ScanAddresses function won't return the change addresses
			changeAddrs, err := w.GetAddresses(wallet.OptionChange())
			require.NoError(t, err)
			require.Equal(t, tc.expectAllChangeAddrs, changeAddrs)
		})
	}
}

func getExternalAddrs(t *testing.T) []cipher.Addresser {
	return skycoinAddressStringsToAddress(testSkycoinExternalAddresses)
}

func getChangeAddrs(t *testing.T) []cipher.Addresser {
	return skycoinAddressStringsToAddress(testSkycoinChangeAddresses)
}
