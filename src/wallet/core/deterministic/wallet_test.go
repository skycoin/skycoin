package deterministic

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name    string
		wltName string
		seed    string
		label   string
		opts    []wallet.Option
		expect  expect
	}{
		{
			name:    "ok all defaults",
			wltName: "test.wlt",
			label:   "",
			seed:    "testseed123",
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     string(wallet.CoinTypeSkycoin),
					"type":     wallet.WalletTypeDeterministic,
					"seed":     "testseed123",
					"version":  wallet.Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok with label, seed and coin set, deterministic",
			wltName: "test.wlt",
			label:   "test",
			seed:    "testseed123",
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeBitcoin),
			},
			expect: expect{
				meta: map[string]string{
					"label":    "test",
					"filename": "test.wlt",
					"coin":     string(wallet.CoinTypeBitcoin),
					"type":     wallet.WalletTypeDeterministic,
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
		{
			name:    "ok default crypto type, deterministic",
			wltName: "test.wlt",
			label:   "test",
			seed:    "testseed123",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			expect: expect{
				meta: map[string]string{
					"label":     "test",
					"coin":      string(wallet.CoinTypeSkycoin),
					"type":      wallet.WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: nil,
			},
		},
		{
			name:    "encrypt without password, deterministic",
			wltName: "test.wlt",
			label:   "wallet1",
			seed:    "testseed123",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
			},
			expect: expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(wallet.CoinTypeSkycoin),
					"type":      wallet.WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: wallet.ErrMissingPassword,
			},
		},
		{
			name:    "create with no seed, deterministic",
			wltName: "test.wlt",
			label:   "test",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			expect: expect{
				meta: map[string]string{
					"label":     "test",
					"coin":      string(wallet.CoinTypeSkycoin),
					"type":      wallet.WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: wallet.ErrMissingSeed,
			},
		},
		{
			name:    "password=pwd encrypt=false, deterministic",
			wltName: "test.wlt",
			label:   "test",
			seed:    "seed",
			opts: []wallet.Option{
				wallet.OptionEncrypt(false),
				wallet.OptionPassword([]byte("pwd")),
			},
			expect: expect{
				err: wallet.ErrMissingEncrypt,
			},
		},
	}

	for _, tc := range tt {
		// test all supported crypto types
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			opts := useInsecureCrypto(tc.opts, ct)
			t.Run(name, func(t *testing.T) {
				w, err := NewWallet(tc.wltName, tc.label, tc.seed, opts...)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				//require.Equal(t, tc.opts.Encrypt, w.IsEncrypted())
				// confirms the meta data
				for k, v := range tc.expect.meta {
					require.Equal(t, v, w.Meta[k])
				}

				if w.IsEncrypted() {
					// Confirms the seeds and entry secrets are all empty
					require.Equal(t, "", w.Seed())
					require.Equal(t, "", w.LastSeed())
					entries, err := w.GetEntries()
					require.NoError(t, err)

					for _, e := range entries {
						require.True(t, e.Secret.Null())
					}

					// Confirms that secrets field is not empty
					require.NotEmpty(t, w.Secrets())
				}
			})
		}
	}
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
			name:    "ok deterministic",
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
			opts := useInsecureCrypto(tc.opts, ct)

			t.Run(name, func(t *testing.T) {
				wltName := wallet.NewWalletFilename()
				w, err := NewWallet(wltName, "test", "testseed123", opts...)
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

				require.True(t, w.IsEncrypted())

				// Checks if the seeds are wiped
				require.Empty(t, w.Seed())
				require.Empty(t, w.LastSeed())

				// Checks if the entries are encrypted
				entries, err := w.GetEntries()
				require.NoError(t, err)

				for _, e := range entries {
					require.Equal(t, cipher.SecKey{}, e.Secret)
				}
			})

		}
	}
}

func TestWalletUnlock(t *testing.T) {
	tt := []struct {
		name      string
		opts      []wallet.Option
		unlockPwd []byte
		err       error
	}{
		{
			name: "ok deterministic",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
				wallet.OptionGenerateN(1),
			},
			unlockPwd: []byte("pwd"),
		},
		{
			name: "unlock with nil password",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			unlockPwd: nil,
			err:       wallet.ErrMissingPassword,
		},
		{
			name: "unlock with wrong password",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			unlockPwd: []byte("wrong_pwd"),
			err:       wallet.ErrInvalidPassword,
		},
		{
			name:      "unlock undecrypted wallet",
			unlockPwd: []byte("pwd"),
			err:       wallet.ErrWalletNotEncrypted,
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)

			opts := useInsecureCrypto(tc.opts, ct)

			t.Run(name, func(t *testing.T) {
				w, err := NewWallet("test.wlt", "test", "testseed123", opts...)
				require.NoError(t, err)
				// Tests the unlock method
				wlt, err := w.Unlock(tc.unlockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.False(t, wlt.IsEncrypted())

				// Checks the seeds
				require.Equal(t, "testseed123", wlt.Seed())

				// Checks the generated addresses
				el, err := wlt.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, 1, el)

				sd, sks := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(wlt.Seed()), 1)
				require.Equal(t, hex.EncodeToString(sd), wlt.LastSeed())
				entries, err := wlt.GetEntries()
				require.NoError(t, err)
				for i, e := range entries {
					addr := cipher.MustAddressFromSecKey(sks[i])
					require.Equal(t, addr, e.Address)
				}

				// Checks the original seeds
				require.NotEqual(t, "testseed123", w.Seed())

				// Checks if the seckeys in entries of original wallet are empty
				entries, err = w.GetEntries()
				require.NoError(t, err)
				for _, e := range entries {
					require.True(t, e.Secret.Null())
				}

				// Checks if the seed and lastSeed in original wallet are still empty
				require.Empty(t, w.Seed())
				require.Empty(t, w.LastSeed())
				require.Empty(t, w.SeedPassphrase())
			})
		}
	}
}

func TestLockAndUnLock(t *testing.T) {
	for _, ct := range crypto.TypesInsecure() {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			w, err := NewWallet("wallet.wlt", "test", bip39.MustNewDefaultMnemonic(), wallet.OptionCryptoType(ct))
			require.NoError(t, err)
			_, err = w.GenerateAddresses(9)
			require.NoError(t, err)
			el, err := w.EntriesLen()
			require.NoError(t, err)
			require.Equal(t, 9, el)

			// clone the wallet
			cw := w.Clone()
			require.Equal(t, w, cw)

			// lock the cloned wallet
			err = cw.Lock([]byte("pwd"))
			require.NoError(t, err)

			checkNoSensitiveData(t, cw)

			// unlock the cloned wallet
			ucw, err := cw.Unlock([]byte("pwd"))
			require.NoError(t, err)

			require.Equal(t, w, ucw)
		})
	}
}

func TestWalletGenerateAddress(t *testing.T) {
	seed := bip39.MustNewDefaultMnemonic()

	tt := []struct {
		name               string
		seed               string
		opts               []wallet.Option
		num                uint64
		oneAddressEachTime bool
		err                error
	}{
		{
			name: "ok with none address deterministic",
			seed: seed,
			num:  0,
		},
		{
			name: "ok with one address deterministic",
			seed: seed,
			num:  1,
		},
		{
			name: "ok with two address deterministic",
			seed: seed,
			num:  2,
		},
		{
			name:               "ok with three address and generate one address each time deterministic",
			seed:               seed,
			num:                2,
			oneAddressEachTime: true,
		},
		{
			name: "wallet is encrypted deterministic",
			seed: seed,
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
				wallet.OptionPassword([]byte("pwd")),
			},
			num:                2,
			oneAddressEachTime: true,
			err:                wallet.ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for _, ct := range crypto.TypesInsecure() {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			opts := useInsecureCrypto(tc.opts, ct)
			t.Run(name, func(t *testing.T) {
				// create wallet
				w, err := NewWallet("test.wlt", "test", tc.seed, opts...)
				require.NoError(t, err)

				// generate addresses
				if !tc.oneAddressEachTime {
					_, err = w.GenerateAddresses(tc.num)
					require.Equal(t, tc.err, err)
					if err != nil {
						return
					}
				} else {
					for i := uint64(0); i < tc.num; i++ {
						_, err := w.GenerateAddresses(1)
						require.Equal(t, tc.err, err)
						if err != nil {
							return
						}
					}
				}

				// check the entry number
				l, err := w.EntriesLen()
				require.NoError(t, err)
				require.Equal(t, int(tc.num), l)

				addrs, err := w.GetAddresses()
				require.NoError(t, err)

				_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(tc.seed), int(tc.num))
				for i, k := range keys {
					a := cipher.MustAddressFromSecKey(k)
					require.Equal(t, a.String(), addrs[i].String())
				}
			})
		}
	}
}

func TestWalletGetEntry(t *testing.T) {
	tt := []struct {
		name    string
		wltFile string
		address string
		err     error
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			"JUdRuTiqD1mGcw358twMg3VPpXpzbkdRvJ",
			nil,
		},
		{
			"entry not exist",
			"./testdata/test1.wlt",
			"2ULfxDUuenUY5V4Pr8whmoAwFdUseXNyjXC",
			wallet.ErrEntryNotFound,
		},
		{
			"scrypt-chacha20poly1305 encrypted wallet",
			"./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			"LxcitUpWQZbPjgEPs6R1i3G4Xa31nPMoSG",
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data, err := ioutil.ReadFile(tc.wltFile)
			require.NoError(t, err)
			ld := Loader{}
			w, err := ld.Load(data)
			require.NoError(t, err)

			a, err := cipher.DecodeBase58Address(tc.address)
			require.NoError(t, err)
			e, err := w.GetEntry(a)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			require.Equal(t, tc.address, e.Address.String())
		})
	}
}

func useInsecureCrypto(opts []wallet.Option, ct crypto.CryptoType) []wallet.Option {
	mOpts := &wallet.AdvancedOptions{}
	for _, opt := range opts {
		opt(mOpts)
	}

	if mOpts.Encrypt {
		// append the insecure crypto type
		opts = append(opts, wallet.OptionCryptoType(ct))
	}
	return opts
}

func checkNoSensitiveData(t *testing.T, w wallet.Wallet) {
	require.Empty(t, w.Seed())
	require.Empty(t, w.LastSeed())
	entries, err := w.GetEntries()
	require.NoError(t, err)
	for _, e := range entries {
		require.True(t, e.Secret.Null())
	}
}
