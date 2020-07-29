package deterministic

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/stretchr/testify/require"
)

var (
	testSeed           = "test123"
	testSkycoinEntries = skycoinEntries([]readableEntry{
		{
			Address: "B4B6Hx1a3WPUHP323Bhqydifeu8TS4Zfan",
			Public:  "0251a81011b0b766242fb3d6777ae1f62e490e73e0c66ed25cbbb45421fa476356",
			Secret:  "91570790c29faa5ecfea981fdcb4bbb81280309f3f17dec4bad6e7697e126410",
		},
		{
			Address: "2FgDYaVqoR3DusaUmuin6xYnSW2FKpCcRrX",
			Public:  "028fcc354cb75dc2041ad2f938f5cd8453f6b799c6550a6f78aef214fa9b13721e",
			Secret:  "df18db782378605e9e5cbff9c845af6139294a114ba91ae530ad1bd20738c9e2",
		},
		{
			Address: "2L8awjtwfe1pMbkHKB9zZdeHYmSNAB1Krug",
			Public:  "021d175a6e13e58d5223d6d0d517eb66e6f0802674b762a4430b84b0349d79cbfb",
			Secret:  "1ad93991934205910301a773cf4f747a9a1fb2f86ec0478d5bb7b8296da9df8e",
		},
		{
			Address: "2VgCCNKj3TXFzddZSrSvbgHREdBB51Em4pN",
			Public:  "024290c4f4c6c2b975af998345564e069d45da9bbd814502600dab23fb38517110",
			Secret:  "cde8f59cf5d052e0c9594493caef8779179ee902870af28ef0e3fcd311774d99",
		},
		{
			Address: "7fqATQGPa6x3Qb7uakFWJM5xqSv3y8wtA1",
			Public:  "0399558a5cfd5b439175776252ec4f01287eaf9a93a791d3f0f12cb098453059a9",
			Secret:  "9238164770d2d7bd9dc1b9f303522b46f08c2b5db175ef7b1f754be9f4439725",
		},
	})
)

func skycoinEntries(es []readableEntry) []wallet.Entry {
	entries := make([]wallet.Entry, len(es))
	for i, e := range es {
		pk, err := cipher.PubKeyFromHex(e.Public)
		if err != nil {
			panic(err)
		}
		sk, err := cipher.SecKeyFromHex(e.Secret)
		if err != nil {
			panic(err)
		}

		entries[i] = wallet.Entry{
			Address: cipher.MustDecodeBase58Address(e.Address),
			Public:  pk,
			Secret:  sk,
		}
	}

	return entries
}

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
				l, err := w.EntriesLen()
				require.NoError(t, err)
				fmt.Println(l)

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
				l, err = w.EntriesLen()
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

func TestWalletValidate(t *testing.T) {
	goodMetaUnencrypted := map[string]string{
		"filename":  "foo.wlt",
		"type":      wallet.WalletTypeDeterministic,
		"coin":      string(wallet.CoinTypeSkycoin),
		"encrypted": "false",
		"seed":      "fooseed",
		"lastSeed":  "foolastseed",
	}

	goodMetaEncrypted := map[string]string{
		"filename":   "foo.wlt",
		"type":       wallet.WalletTypeDeterministic,
		"coin":       string(wallet.CoinTypeSkycoin),
		"encrypted":  "true",
		"cryptoType": "scrypt-chacha20poly1305",
		"seed":       "",
		"lastSeed":   "",
		"secrets":    "xacsdasdasdasd",
	}

	copyMap := func(m map[string]string) map[string]string {
		n := make(map[string]string, len(m))
		for k, v := range m {
			n[k] = v
		}
		return n
	}

	delField := func(m map[string]string, f string) map[string]string {
		n := copyMap(m)
		delete(n, f)
		return n
	}

	setField := func(m map[string]string, f, g string) map[string]string {
		n := copyMap(m)
		n[f] = g
		return n
	}

	cases := []struct {
		name string
		meta map[string]string
		err  error
	}{
		{
			name: "missing filename",
			meta: delField(goodMetaUnencrypted, wallet.MetaFilename),
			err:  errors.New("filename not set"),
		},
		{
			name: "wallet type missing",
			meta: delField(goodMetaUnencrypted, wallet.MetaType),
			err:  errors.New("type field not set"),
		},
		{
			name: "invalid wallet type",
			meta: setField(goodMetaUnencrypted, wallet.MetaType, "footype"),
			err:  wallet.ErrInvalidWalletType,
		},
		{
			name: "coin field missing",
			meta: delField(goodMetaUnencrypted, wallet.MetaCoin),
			err:  errors.New("coin field not set"),
		},
		{
			name: "encrypted field invalid",
			meta: setField(goodMetaUnencrypted, wallet.MetaEncrypted, "foo"),
			err:  errors.New("encrypted field is not a valid bool"),
		},
		{
			name: "unencrypted missing seed",
			meta: delField(goodMetaUnencrypted, wallet.MetaSeed),
			err:  errors.New("seed missing in unencrypted deterministic wallet"),
		},
		{
			name: "unencrypted missing last seed",
			meta: delField(goodMetaUnencrypted, wallet.MetaLastSeed),
			err:  errors.New("lastSeed missing in unencrypted deterministic wallet"),
		},
		{
			name: "crypto type missing",
			meta: delField(goodMetaEncrypted, wallet.MetaCryptoType),
			err:  errors.New("crypto type field not set"),
		},
		{
			name: "crypto type invalid",
			meta: setField(goodMetaEncrypted, wallet.MetaCryptoType, "foocryptotype"),
			err:  errors.New("unknown crypto type"),
		},
		{
			name: "secrets missing",
			meta: delField(goodMetaEncrypted, wallet.MetaSecrets),
			err:  errors.New("wallet is encrypted, but secrets field not set"),
		},
		{
			name: "secrets empty",
			meta: setField(goodMetaEncrypted, wallet.MetaSecrets, ""),
			err:  errors.New("wallet is encrypted, but secrets field not set"),
		},
		{
			name: "valid unencrypted",
			meta: goodMetaUnencrypted,
		}, {
			name: "valid encrypted",
			meta: goodMetaEncrypted,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &Wallet{
				Meta: tc.meta,
			}
			err := w.Validate()

			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.err, err, "%s != %s", tc.err, err)
			}
		})
	}
}

func TestWalletSerialize(t *testing.T) {
	w, err := NewWallet("test.wlt", "test", "test123")
	require.NoError(t, err)

	_, err = w.GenerateAddresses(5)
	require.NoError(t, err)

	w.SetTimestamp(0)
	b, err := w.Serialize()
	require.NoError(t, err)

	// load wallet file and compare
	fb, err := ioutil.ReadFile("./testdata/wallet_serialize.wlt")
	require.NoError(t, err)
	require.Equal(t, bytes.TrimRight(fb, "\n"), b)

	wlt := Wallet{}
	err = wlt.Deserialize(b)
	require.NoError(t, err)
}

func TestWalletDeserialize(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/wallet_serialize.wlt")
	require.NoError(t, err)

	w := Wallet{}
	err = w.Deserialize(b)
	require.NoError(t, err)

	require.Equal(t, w.Filename(), "test.wlt")
	require.Equal(t, w.Label(), "test")
	entries, err := w.GetEntries()
	require.NoError(t, err)
	require.Equal(t, 5, len(entries))
	for i, e := range entries {
		require.Equal(t, testSkycoinEntries[i], e)
	}
}
