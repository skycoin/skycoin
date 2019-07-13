package wallet

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip39"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/cipher/encrypt"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	log = logging.MustGetLogger("wallet_test")
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

var u = flag.Bool("u", false, "update test wallet file in ./testdata")

func init() {
	flag.Parse()

	// Change the scrypt N value in cryptoTable to make test faster, otherwise
	// it would take more than 200 seconds to finish.
	cryptoTable[CryptoTypeScryptChacha20poly1305] = encrypt.ScryptChacha20poly1305{
		N:      1 << 15,
		R:      encrypt.ScryptR,
		P:      encrypt.ScryptP,
		KeyLen: encrypt.ScryptKeyLen,
	}

	// When -u flag is specified, update the following wallet files:
	//     - ./testdata/scrypt-chacha20poly1305-encrypted.wlt
	//     - ./testdata/sha256xor-encrypted.wlt
	if *u {
		// Update ./testdata/scrypt-chacha20poly1305-encrypted.wlt
		//     - Create an unencrypted wallet
		//     - Generate an address
		//     - Lock the wallet with scrypt-chacha20poly1305 crypto type and password of "pwd".
		w, err := NewWallet("scrypt-chacha20poly1305-encrypted.wlt", Options{
			Seed:  "seed-scrypt-chacha20poly1305",
			Label: "scrypt-chacha20poly1305",
		})
		if err != nil {
			log.Panic(err)
		}

		if _, err := w.GenerateAddresses(1); err != nil {
			log.Panic(err)
		}

		if err := Lock(w, []byte("pwd"), CryptoTypeScryptChacha20poly1305); err != nil {
			log.Panic(err)
		}

		if err := Save(w, "./testdata"); err != nil {
			log.Panic(err)
		}

		// Update ./testdata/sha256xor-encrypted.wlt
		//     - Create an sha256xor encrypted wallet with password: "pwd".
		w1, err := NewWallet("sha256xor-encrypted.wlt", Options{
			Seed:       "seed-sha256xor",
			Label:      "sha256xor",
			Encrypt:    true,
			Password:   []byte("pwd"),
			CryptoType: CryptoTypeSha256Xor,
		})
		if err != nil {
			log.Panic(err)
		}

		if err := Save(w1, "./testdata"); err != nil {
			log.Panic(err)
		}
	}
}

func TestNewWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name    string
		wltName string
		opts    Options
		expect  expect
	}{
		{
			name:    "ok, empty collection wallet",
			wltName: "test-collection.wlt",
			opts: Options{
				Type: WalletTypeCollection,
			},
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test-collection.wlt",
					"coin":     string(CoinTypeSkycoin),
					"type":     WalletTypeDeterministic,
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok, xpub wallet",
			wltName: "test-xpub.wlt",
			opts: Options{
				Type: WalletTypeXPub,
				XPub: "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
			},
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test-collection.wlt",
					"coin":     string(CoinTypeSkycoin),
					"type":     WalletTypeDeterministic,
					"version":  Version,
					"xpub":     "xpub6CkxdS1d4vNqqcnf9xPgqR5e2jE2PZKmKSw93QQMjHE1hRk22nU4zns85EDRgmLWYXYtu62XexwqaET33XA28c26NbXCAUJh1xmqq6B3S2v",
				},
				err: nil,
			},
		},
		{
			name:    "ok all defaults",
			wltName: "test.wlt",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeBip44,
			},
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     string(CoinTypeSkycoin),
					"type":     WalletTypeDeterministic,
					"seed":     "testseed123",
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok with seed set, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Seed: "testseed123",
				Type: WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     string(CoinTypeSkycoin),
					"type":     WalletTypeDeterministic,
					"seed":     "testseed123",
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok with label and seed set, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label: "wallet1",
				Seed:  "testseed123",
				Type:  WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     string(CoinTypeSkycoin),
					"type":     WalletTypeDeterministic,
					"seed":     "testseed123",
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok with label, seed and coin set, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label: "wallet1",
				Coin:  CoinTypeBitcoin,
				Seed:  "testseed123",
				Type:  WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     string(CoinTypeBitcoin),
					"type":     WalletTypeDeterministic,
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
		{
			name:    "ok default crypto type, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Seed:     "testseed123",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: nil,
			},
		},
		{
			name:    "encrypt without password, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label:   "wallet1",
				Coin:    CoinTypeSkycoin,
				Seed:    "testseed123",
				Encrypt: true,
				Type:    WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: ErrMissingPassword,
			},
		},
		{
			name:    "create with no seed, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			expect: expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      WalletTypeDeterministic,
					"encrypted": "true",
				},
				err: ErrMissingSeed,
			},
		},
		{
			name:    "password=pwd encrypt=false, deterministic",
			wltName: "test.wlt",
			opts: Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Encrypt:  false,
				Seed:     "seed",
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			expect: expect{
				err: ErrMissingEncrypt,
			},
		},
		{
			name:    "ok bip44",
			wltName: "bip44.wlt",
			opts: Options{
				Label: "bip44wallet1",
				Type:  WalletTypeBip44,
				Seed:  "voyage say extend find sheriff surge priority merit ignore maple cash argue",
			},
			expect: expect{
				meta: map[string]string{
					"label":     "bip44wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      string(WalletTypeBip44),
					"version":   Version,
					"bip44Coin": "8000",
				},
			},
		},
		{
			name:    "invalid xpub wallet",
			wltName: "test-xpub.wlt",
			opts: Options{
				Type: WalletTypeXPub,
				XPub: "xpubbad",
			},
			expect: expect{
				err: NewError(errors.New("invalid xpub key: Serialized keys should be exactly 82 bytes")),
			},
		},
	}

	for _, tc := range tt {
		// test all supported crypto types
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			if tc.opts.Encrypt {
				tc.opts.CryptoType = ct
			}
			t.Run(name, func(t *testing.T) {
				w, err := NewWallet(tc.wltName, tc.opts)

				if tc.expect.err == nil {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
					require.Equal(t, tc.expect.err, err, "%s != %s", tc.expect.err.Error(), err.Error())
					return
				}

				require.Equal(t, tc.opts.Encrypt, w.IsEncrypted())

				if w.IsEncrypted() {
					// Confirms the seeds and entry secrets are all empty
					require.Equal(t, "", w.Seed())
					require.Equal(t, "", w.LastSeed())

					for _, e := range w.GetEntries() {
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
		opts    Options
		lockPwd []byte
		err     error
	}{
		{
			name: "ok deterministic",
			opts: Options{
				Seed: "seed",
				Type: WalletTypeDeterministic,
			},
			lockPwd: []byte("pwd"),
		},
		{
			name: "ok bip44",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeBip44,
			},
			lockPwd: []byte("pwd"),
		},
		{
			name: "password is nil",
			opts: Options{
				Seed: "seed",
				Type: WalletTypeDeterministic,
			},
			lockPwd: nil,
			err:     ErrMissingPassword,
		},
		{
			name: "wallet already encrypted",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			lockPwd: []byte("pwd"),
			err:     ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			if tc.opts.Encrypt {
				tc.opts.CryptoType = ct
			}
			t.Run(name, func(t *testing.T) {
				wltName := NewWalletFilename()
				w, err := NewWallet(wltName, tc.opts)
				require.NoError(t, err)

				if !w.IsEncrypted() {
					// Generates 2 addresses
					_, err = w.GenerateAddresses(2)
					require.NoError(t, err)
				}

				err = Lock(w, tc.lockPwd, ct)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.True(t, w.IsEncrypted())

				// Checks if the seeds are wiped
				require.Empty(t, w.Seed())
				require.Empty(t, w.LastSeed())

				// Checks if the entries are encrypted
				for _, e := range w.GetEntries() {
					require.Equal(t, cipher.SecKey{}, e.Secret)
				}
			})

		}
	}
}

func TestWalletUnlock(t *testing.T) {
	tt := []struct {
		name      string
		opts      Options
		unlockPwd []byte
		err       error
	}{
		{
			name: "ok deterministic",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			unlockPwd: []byte("pwd"),
		},
		{
			name: "ok bip44",
			opts: Options{
				Seed:     bip39.MustNewDefaultMnemonic(),
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeBip44,
			},
			unlockPwd: []byte("pwd"),
		},
		{
			name: "unlock with nil password",
			opts: Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
				Type:     WalletTypeDeterministic,
			},
			unlockPwd: nil,
			err:       ErrMissingPassword,
		},
		{
			name: "unlock undecrypted wallet",
			opts: Options{
				Seed:    "seed",
				Encrypt: false,
				Type:    WalletTypeDeterministic,
			},
			unlockPwd: []byte("pwd"),
			err:       ErrWalletNotEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			if tc.opts.Encrypt {
				tc.opts.CryptoType = ct
			}
			t.Run(name, func(t *testing.T) {
				w := makeWallet(t, tc.opts, 1)
				// Tests the unlock method
				wlt, err := Unlock(w, tc.unlockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.False(t, wlt.IsEncrypted())

				// Checks the seeds
				require.Equal(t, tc.opts.Seed, wlt.Seed())

				// Checks the generated addresses
				require.Equal(t, 1, wlt.EntriesLen())

				switch tc.opts.Type {
				case WalletTypeBip44:
					require.Empty(t, wlt.LastSeed())
					keys := generateBip44Chain(t, wlt.Seed(), wlt.SeedPassphrase(), bip44.ExternalChainIndex, 1)
					for i, e := range wlt.GetEntries() {
						sk := cipher.MustNewSecKey(keys[i].Key)
						addr := cipher.MustAddressFromSecKey(sk)
						require.Equal(t, addr, e.Address)
					}

				case WalletTypeDeterministic:
					sd, sks := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(wlt.Seed()), 1)
					require.Equal(t, hex.EncodeToString(sd), wlt.LastSeed())

					for i, e := range wlt.GetEntries() {
						addr := cipher.MustAddressFromSecKey(sks[i])
						require.Equal(t, addr, e.Address)
					}
				default:
					t.Fatalf("unhandled wallet type %q", tc.opts.Type)
				}

				// Checks the original seeds
				require.NotEqual(t, tc.opts.Seed, w.Seed())

				// Checks if the seckeys in entries of original wallet are empty
				for _, e := range w.GetEntries() {
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
	for _, walletType := range []string{
		WalletTypeBip44,
		WalletTypeDeterministic,
	} {
		for ct := range cryptoTable {
			t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
				w, err := NewWallet("wallet", Options{
					Label: "wallet",
					Seed:  bip39.MustNewDefaultMnemonic(),
					Type:  walletType,
				})
				require.NoError(t, err)
				_, err = w.GenerateAddresses(9)
				require.NoError(t, err)
				require.Equal(t, 10, w.EntriesLen())

				if walletType == WalletTypeBip44 {
					// Generate change entries, to verify that their secret keys
					// are protected and revealed when locked and unlocked
					for i := 0; i < 5; i++ {
						e, err := w.(*Bip44Wallet).GenerateChangeEntry()
						require.NoError(t, err)
						require.Equal(t, bip44.ChangeChainIndex, e.Change)
					}

					require.Equal(t, 15, w.EntriesLen())
					nExternal := 0
					nChange := 0
					for _, e := range w.GetEntries() {
						switch e.Change {
						case bip44.ExternalChainIndex:
							nExternal++
						case bip44.ChangeChainIndex:
							nChange++
						default:
							t.Fatalf("invalid change chain index: %d", e.Change)
						}
					}
					require.Equal(t, 10, nExternal)
					require.Equal(t, 5, nChange)
				}

				// clone the wallet
				cw := w.Clone()
				require.Equal(t, w, cw)

				// lock the cloned wallet
				err = Lock(cw, []byte("pwd"), ct)
				require.NoError(t, err)

				checkNoSensitiveData(t, cw)

				// unlock the cloned wallet
				ucw, err := Unlock(cw, []byte("pwd"))
				require.NoError(t, err)

				require.Equal(t, w, ucw)
			})
		}
	}
}

func makeWallet(t *testing.T, opts Options, addrNum uint64) Wallet { //nolint:unparam
	// Create an unlocked wallet, then generate addresses, lock if the options.Encrypt is true.
	preOpts := opts
	opts.Encrypt = false
	opts.Password = nil
	w, err := NewWallet("t.wlt", opts)
	require.NoError(t, err)

	if addrNum > 1 {
		_, err = w.GenerateAddresses(addrNum - 1)
		require.NoError(t, err)
	}
	if preOpts.Encrypt {
		err = Lock(w, preOpts.Password, preOpts.CryptoType)
		require.NoError(t, err)
	}
	return w
}

func TestLoadWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name     string
		filename string
		expect   expect
	}{
		{
			name:     "ok",
			filename: "./testdata/test1.wlt",
			expect: expect{
				meta: map[string]string{
					"coin":     string(CoinTypeSkycoin),
					"filename": "test1.wlt",
					"label":    "test3",
					"lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
					"seed":     "buddy fossil side modify turtle door label grunt baby worth brush master",
					"tm":       "1503458909",
					"type":     WalletTypeDeterministic,
					"version":  "0.1",
				},
				err: nil,
			},
		},
		{
			name:     "wallet file doesn't exist",
			filename: "not_exist_file.wlt",
			expect: expect{
				meta: map[string]string{},
				err:  fmt.Errorf("wallet \"not_exist_file.wlt\" doesn't exist"),
			},
		},
		{
			name:     "invalid wallet: no type",
			filename: "./testdata/invalid_wallets/no_type.wlt",
			expect: expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet \"./testdata/invalid_wallets/no_type.wlt\": invalid wallet type"),
			},
		},
		{
			name:     "invalid wallet: invalid type",
			filename: "./testdata/invalid_wallets/err_type.wlt",
			expect: expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet \"./testdata/invalid_wallets/err_type.wlt\": invalid wallet type"),
			},
		},
		{
			name:     "invalid wallet: no coin",
			filename: "./testdata/invalid_wallets/no_coin.wlt",
			expect: expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet \"./testdata/invalid_wallets/no_coin.wlt\": invalid coin type"),
			},
		},
		{
			name:     "invalid wallet: no seed",
			filename: "./testdata/invalid_wallets/no_seed.wlt",
			expect: expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet \"no_seed.wlt\": seed missing in unencrypted deterministic wallet"),
			},
		},
		{
			name:     "version=0.2 encrypted=true crypto=scrypt-chacha20poly1305",
			filename: "./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			expect: expect{
				meta: map[string]string{
					"coin":       string(CoinTypeSkycoin),
					"cryptoType": "scrypt-chacha20poly1305",
					"encrypted":  "true",
					"filename":   "scrypt-chacha20poly1305-encrypted.wlt",
					"label":      "scrypt-chacha20poly1305",
					"lastSeed":   "",
					"seed":       "",
					"type":       WalletTypeDeterministic,
					"version":    "0.2",
				},
				err: nil,
			},
		},
		{
			name:     "version=0.2 encrypted=true crypto=sha256xor",
			filename: "./testdata/sha256xor-encrypted.wlt",
			expect: expect{
				meta: map[string]string{
					"coin":       string(CoinTypeSkycoin),
					"cryptoType": "sha256-xor",
					"encrypted":  "true",
					"filename":   "sha256xor-encrypted.wlt",
					"label":      "sha256xor",
					"lastSeed":   "",
					"seed":       "",
					"type":       WalletTypeDeterministic,
					"version":    "0.2",
				},
				err: nil,
			},
		},
		{
			name:     "version=0.2 encrypted=false",
			filename: "./testdata/v2_no_encrypt.wlt",
			expect: expect{
				meta: map[string]string{
					"coin":       string(CoinTypeSkycoin),
					"cryptoType": "scrypt-chacha20poly1305",
					"encrypted":  "false",
					"filename":   "v2_no_encrypt.wlt",
					"label":      "v2_no_encrypt",
					"lastSeed":   "c79454cf362b3f55e5effce09f664311650a44b9c189b3c8eed1ae9bd696cd9e",
					"secrets":    "",
					"seed":       "seed",
					"type":       WalletTypeDeterministic,
					"version":    "0.2",
				},
				err: nil,
			},
		},
		{
			name:     "version=0.3 encrypted=false type=bip44",
			filename: "./testdata/test5-bip44.wlt",
			expect: expect{
				meta: map[string]string{
					"bip44Coin":      fmt.Sprint(bip44.CoinTypeSkycoin),
					"coin":           string(CoinTypeSkycoin),
					"cryptoType":     "",
					"encrypted":      "false",
					"filename":       "test5-bip44.wlt",
					"label":          "test5-bip44",
					"lastSeed":       "",
					"secrets":        "",
					"seed":           "voyage say extend find sheriff surge priority merit ignore maple cash argue",
					"seedPassphrase": "",
					"type":           WalletTypeBip44,
					"version":        "0.3",
				},
				err: nil,
			},
		},
		{
			name:     "version=0.3 encrypted=false type=bip44 seed-passphrase=true",
			filename: "./testdata/test6-passphrase-bip44.wlt",
			expect: expect{
				meta: map[string]string{
					"bip44Coin":      fmt.Sprint(bip44.CoinTypeSkycoin),
					"coin":           string(CoinTypeSkycoin),
					"cryptoType":     "",
					"encrypted":      "false",
					"filename":       "test6-passphrase-bip44.wlt",
					"label":          "test6-passphrase-bip44",
					"lastSeed":       "",
					"secrets":        "",
					"seed":           "voyage say extend find sheriff surge priority merit ignore maple cash argue",
					"seedPassphrase": "foobar",
					"type":           WalletTypeBip44,
					"version":        "0.3",
				},
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := Load(tc.filename)
			if err != nil {
				require.Equal(t, tc.expect.err.Error(), err.Error())
			}
			if err != nil {
				return
			}

			for k, v := range tc.expect.meta {
				vv := w.Find(k)
				require.Equal(t, v, vv)
			}

			if w.IsEncrypted() {
				require.NotEmpty(t, w.Secrets())
			}
		})
	}
}

func TestWalletGenerateAddress(t *testing.T) {
	tt := []struct {
		name               string
		opts               Options
		num                uint64
		oneAddressEachTime bool
		err                error
	}{
		{
			name: "ok with one address deterministic",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeDeterministic,
			},
			num: 1,
		},
		{
			name: "ok with two address deterministic",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeDeterministic,
			},
			num: 2,
		},
		{
			name: "ok with three address and generate one address each time deterministic",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeDeterministic,
			},
			num:                2,
			oneAddressEachTime: true,
		},
		{
			name: "wallet is encrypted deterministic",
			opts: Options{
				Seed:     bip39.MustNewDefaultMnemonic(),
				Type:     WalletTypeDeterministic,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			num:                2,
			oneAddressEachTime: true,
			err:                ErrWalletEncrypted,
		},

		{
			name: "ok with one address bip44",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeBip44,
			},
			num: 1,
		},
		{
			name: "ok with two address bip44",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeBip44,
			},
			num: 2,
		},
		{
			name: "ok with three address and generate one address each time bip44",
			opts: Options{
				Seed: bip39.MustNewDefaultMnemonic(),
				Type: WalletTypeBip44,
			},
			num:                2,
			oneAddressEachTime: true,
		},
		{
			name: "wallet is encrypted bip44",
			opts: Options{
				Seed:     bip39.MustNewDefaultMnemonic(),
				Type:     WalletTypeBip44,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			num:                2,
			oneAddressEachTime: true,
			err:                ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("crypto=%v %v", ct, tc.name)
			if tc.opts.Encrypt {
				tc.opts.CryptoType = ct
			}

			t.Run(name, func(t *testing.T) {
				// create wallet
				w, err := NewWallet("test.wlt", tc.opts)
				require.NoError(t, err)

				// generate addresses
				if tc.oneAddressEachTime {
					_, err = w.GenerateAddresses(tc.num - 1)
					require.Equal(t, tc.err, err)
					if err != nil {
						return
					}
				} else {
					for i := uint64(0); i < tc.num-1; i++ {
						_, err := w.GenerateAddresses(1)
						require.Equal(t, tc.err, err)
						if err != nil {
							return
						}
					}
				}

				// check the entry number
				require.Equal(t, w.EntriesLen(), int(tc.num))

				addrs := w.GetAddresses()

				switch tc.opts.Type {
				case WalletTypeDeterministic:
					_, keys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(tc.opts.Seed), int(tc.num))
					for i, k := range keys {
						a := cipher.MustAddressFromSecKey(k)
						require.Equal(t, a.String(), addrs[i].String())
					}
				case WalletTypeBip44:
					keys := generateBip44Chain(t, tc.opts.Seed, tc.opts.SeedPassphrase, bip44.ExternalChainIndex, int(tc.num))
					for i, k := range keys {
						sk := cipher.MustNewSecKey(k.Key)
						a := cipher.MustAddressFromSecKey(sk)
						require.Equal(t, a.String(), addrs[i].String())
					}
				default:
					t.Fatalf("unhandled wallet type %q", tc.opts.Type)
				}
			})
		}
	}
}

// generateBip44Chain generates N keys for the leaf change chain
func generateBip44Chain(t *testing.T, seed, seedPassphrase string, change uint32, num int) []*bip32.PrivateKey {
	ss, err := bip39.NewSeed(seed, seedPassphrase)
	require.NoError(t, err)

	cc, err := bip44.NewCoin(ss, bip44.CoinTypeSkycoin)
	require.NoError(t, err)

	acct, err := cc.Account(0)
	require.NoError(t, err)

	chain, err := acct.NewPrivateChildKey(change)
	require.NoError(t, err)

	keys := make([]*bip32.PrivateKey, num)
	for i := 0; i < num; i++ {
		k, err := chain.NewPrivateChildKey(uint32(i))
		require.NoError(t, err)
		keys[i] = k
	}

	return keys
}

func TestWalletGetEntry(t *testing.T) {
	tt := []struct {
		name    string
		wltFile string
		address string
		find    bool
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			"JUdRuTiqD1mGcw358twMg3VPpXpzbkdRvJ",
			true,
		},
		{
			"entry not exist",
			"./testdata/test1.wlt",
			"2ULfxDUuenUY5V4Pr8whmoAwFdUseXNyjXC",
			false,
		},
		{
			"scrypt-chacha20poly1305 encrypted wallet",
			"./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			"LxcitUpWQZbPjgEPs6R1i3G4Xa31nPMoSG",
			true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := Load(tc.wltFile)
			require.NoError(t, err)
			a, err := cipher.DecodeBase58Address(tc.address)
			require.NoError(t, err)
			e, ok := w.GetEntry(a)
			require.Equal(t, tc.find, ok)
			if ok {
				require.Equal(t, tc.address, e.Address.String())
			}
		})
	}
}

func TestWalletCollectionAddEntry(t *testing.T) {
	test1SecKey, err := cipher.SecKeyFromHex("1fc5396e91e60b9fc613d004ea5bd2ccea17053a12127301b3857ead76fdb93e")
	require.NoError(t, err)

	_, s := cipher.GenerateKeyPair()

	cases := []struct {
		name    string
		wltFile string
		secKey  cipher.SecKey
		err     error
	}{
		{
			"ok",
			"./testdata/test4-collection.wlt",
			s,
			nil,
		},
		{
			"dup entry",
			"./testdata/test4-collection.wlt",
			test1SecKey,
			errors.New("wallet already contains entry with this address"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w, err := Load(tc.wltFile)
			require.NoError(t, err)

			a := cipher.MustAddressFromSecKey(tc.secKey)
			p := cipher.MustPubKeyFromSecKey(tc.secKey)
			require.Equal(t, tc.err, w.(*CollectionWallet).AddEntry(Entry{
				Address: a,
				Public:  p,
				Secret:  tc.secKey,
			}))
		})
	}
}

func TestWalletGuard(t *testing.T) {
	cases := []struct {
		name       string
		walletType string
	}{
		{
			name:       "deterministic",
			walletType: WalletTypeDeterministic,
		},
		{
			name:       "bip44",
			walletType: WalletTypeBip44,
		},
	}

	for ct := range cryptoTable {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
				validate := func(w Wallet) {
					require.Equal(t, "", w.Seed())
					require.Equal(t, "", w.LastSeed())
					require.Equal(t, "", w.SeedPassphrase())
					for _, e := range w.GetEntries() {
						require.True(t, e.Secret.Null())
					}
				}

				seed := bip39.MustNewDefaultMnemonic()
				w, err := NewWallet("t.wlt", Options{
					Seed:       seed,
					Encrypt:    true,
					Password:   []byte("pwd"),
					CryptoType: ct,
					Type:       tc.walletType,
				})
				require.NoError(t, err)

				err = GuardUpdate(w, []byte("pwd"), func(w Wallet) error {
					require.Equal(t, seed, w.Seed())
					w.SetLabel("label")
					return nil
				})
				require.NoError(t, err)
				require.Equal(t, "label", w.Label())
				validate(w)

				err = GuardView(w, []byte("pwd"), func(w Wallet) error {
					require.Equal(t, "label", w.Label())
					w.SetLabel("new label")
					return nil
				})
				require.NoError(t, err)

				require.Equal(t, "label", w.Label())
				validate(w)
			})
		}
	}
}

func TestRemoveBackupFiles(t *testing.T) {
	type wltInfo struct {
		wltName string
		version string
	}

	tt := []struct {
		name                   string
		initFiles              []wltInfo
		expectedRemainingFiles map[string]struct{}
	}{
		{
			name:                   "no file",
			initFiles:              []wltInfo{},
			expectedRemainingFiles: map[string]struct{}{},
		},
		{
			name: "wlt v0.1=1 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=2 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=2 delete 2 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=3 delete 3 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t3.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			name: "wlt v0.1=3 bak v0.1=1 no delete",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t2.wlt",
					"0.1",
				},
				{
					"t3.wlt",
					"0.1",
				},
				{
					"t4.wlt.bak",
					"0.1",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t4.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.2=3 bak v0.2=1 no delete",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.2",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.1=1 bak v0.1=1 wlt v0.2=2 bak v0.2=2 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t2.wlt.bak",
					"0.2",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t2.wlt.bak": struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			name: "wlt v0.1=1 bak v0.1=2 wlt v0.2=2 bak v0.2=1 delete 1 bak",
			initFiles: []wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
				{
					"t2.wlt",
					"0.2",
				},
				{
					"t2.wlt.bak",
					"0.1",
				},
				{
					"t3.wlt",
					"0.2",
				},
				{
					"t3.wlt.bak",
					"0.2",
				},
			},
			expectedRemainingFiles: map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t2.wlt.bak": struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			dir := prepareWltDir()
			// Initialize files
			for _, f := range tc.initFiles {
				w, err := NewWallet(f.wltName, Options{
					Seed: "s1",
					Type: WalletTypeDeterministic,
				})
				require.NoError(t, err)
				w.SetVersion(f.version)

				require.NoError(t, Save(w, dir))
			}

			require.NoError(t, removeBackupFiles(dir))

			// Get all remaining files
			fs, err := ioutil.ReadDir(dir)
			require.NoError(t, err)
			require.Len(t, fs, len(tc.expectedRemainingFiles))
			for _, f := range fs {
				_, ok := tc.expectedRemainingFiles[f.Name()]
				require.True(t, ok)
			}
		})
	}
}

func TestWalletValidate(t *testing.T) {
	goodMetaUnencrypted := map[string]string{
		"filename":  "foo.wlt",
		"type":      WalletTypeDeterministic,
		"coin":      string(CoinTypeSkycoin),
		"encrypted": "false",
		"seed":      "fooseed",
		"lastSeed":  "foolastseed",
	}

	goodMetaEncrypted := map[string]string{
		"filename":   "foo.wlt",
		"type":       WalletTypeDeterministic,
		"coin":       string(CoinTypeSkycoin),
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
			meta: delField(goodMetaUnencrypted, metaFilename),
			err:  errors.New("filename not set"),
		},
		{
			name: "wallet type missing",
			meta: delField(goodMetaUnencrypted, metaType),
			err:  errors.New("type field not set"),
		},
		{
			name: "invalid wallet type",
			meta: setField(goodMetaUnencrypted, metaType, "footype"),
			err:  ErrInvalidWalletType,
		},
		{
			name: "coin field missing",
			meta: delField(goodMetaUnencrypted, metaCoin),
			err:  errors.New("coin field not set"),
		},
		{
			name: "encrypted field invalid",
			meta: setField(goodMetaUnencrypted, metaEncrypted, "foo"),
			err:  errors.New("encrypted field is not a valid bool"),
		},
		{
			name: "unencrypted missing seed",
			meta: delField(goodMetaUnencrypted, metaSeed),
			err:  errors.New("seed missing in unencrypted deterministic wallet"),
		},
		{
			name: "unencrypted missing last seed",
			meta: delField(goodMetaUnencrypted, metaLastSeed),
			err:  errors.New("lastSeed missing in unencrypted deterministic wallet"),
		},
		{
			name: "crypto type missing",
			meta: delField(goodMetaEncrypted, metaCryptoType),
			err:  errors.New("crypto type field not set"),
		},
		{
			name: "crypto type invalid",
			meta: setField(goodMetaEncrypted, metaCryptoType, "foocryptotype"),
			err:  errors.New("unknown crypto type"),
		},
		{
			name: "secrets missing",
			meta: delField(goodMetaEncrypted, metaSecrets),
			err:  errors.New("wallet is encrypted, but secrets field not set"),
		},
		{
			name: "secrets empty",
			meta: setField(goodMetaEncrypted, metaSecrets, ""),
			err:  errors.New("wallet is encrypted, but secrets field not set"),
		},
		{
			name: "valid unencrypted",
			meta: goodMetaUnencrypted,
		},
		{
			name: "valid encrypted",
			meta: goodMetaEncrypted,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := &DeterministicWallet{
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
