package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encrypt"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

var u = flag.Bool("u", false, "update test wallet file in ./testdata")

func init() {
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
			Seed:  "seed",
			Label: "scrypt-chacha20poly1305",
		})
		if err != nil {
			log.Panic(err)
		}

		if _, err := w.GenerateAddresses(1); err != nil {
			log.Panic(err)
		}

		if err := w.lock([]byte("pwd"), CryptoTypeScryptChacha20poly1305); err != nil {
			log.Panic(err)
		}

		if err := w.Save("./testdata"); err != nil {
			log.Panic(err)
		}

		// Update ./testdata/sha256xor-encrypted.wlt
		//     - Create an sha256xor encrypted wallet with password: "pwd".
		w1, err := NewWallet("sha256xor-encrypted.wlt", Options{
			Seed:       "seed",
			Label:      "sha256xor",
			Encrypt:    true,
			Password:   []byte("pwd"),
			CryptoType: CryptoTypeSha256Xor,
		})
		if err != nil {
			log.Panic(err)
		}

		if err := w1.Save("./testdata"); err != nil {
			log.Panic(err)
		}
	}
}

type mockBalanceGetter map[cipher.Address]BalancePair

func (mb mockBalanceGetter) GetBalanceOfAddrs(addrs []cipher.Address) ([]BalancePair, error) {
	var bals []BalancePair
	for _, addr := range addrs {
		bal := mb[addr]
		bals = append(bals, bal)
	}
	return bals, nil
}

// 10 addresses of seed1
var addrsOfSeed1 = []string{
	"2GBifzJEehbDX7Mkk63Prfa4MQQQyRzBLfe",
	"q2kU13X8XsAg8cS8BuSeSVzjPF9AT9ghAa",
	"2WXvTagXtrc1Qq71yjNXw86TC6SRgfVRH1B",
	"2NUNw748b9mT2FHRxgJL5KjBHasLfdP32Sh",
	"2V1CnVzWoXDaCX6wHU4tLJkWaFmLcQBb2q4",
	"wBkMr936thcr57wxyrH6ffvA99JN2Q1MN1",
	"2f92Wht7VQefAyoJUz3SEnfwT6wTdeAcq3L",
	"27UM5jPFYVuve3ceEHAYGaJSmkynQYmwPcH",
	"xjWbVN7ihReasVFwXJSSYYWF7rgQa22auC",
	"2LyanokLYFeBfBsNkRYHp2qtN8naGFJqeUw",
}

var childSeedsOfSeed1 = []string{
	"22b826c586039f8078433be26618ca1024e883d97de2267313bb78068f634c5a",
	"68efbbdf8aa06368cfc55e252d1e782bbd7651e590ee59e94ab579d2e44c20ad",
	"8894c818732375680284be4509d153272726f42296b85ecac1fb66b9dc7484b9",
	"6603375ee19c1e9fffe369e3f62e9deaa6931c1183d7da7f24ecbbd591061502",
	"91a63f939149d423ea39701d8ed16cfb16a3554c184d214d2289018ddb9e73de",
	"f0f4f008aa3e7cd32ee953507856fb46e37b734fd289dc01449133d7e37a1f07",
	"6b194da58a5ba5660cf2b00076cf6a2962fe8fe0523abca5647c87df3352866a",
	"b47a2678f7e797d3ada86e7e36855f572a18ab78dcbe54ed0613bba69fd76f8d",
	"fe064533108dadbef13be3a95f547ba03423aa6a701c40aaaed775cb783b12b3",
	"d554da211321a437e4d08f2a57e3ef255cffa89dd182e0fd52a4fd5bdfcab1ae",
}

func fromAddrString(t *testing.T, addrStrs []string) []cipher.Address {
	addrs := make([]cipher.Address, 0, len(addrStrs))
	for _, addr := range addrStrs {
		a, err := cipher.DecodeBase58Address(addr)
		require.NoError(t, err)
		addrs = append(addrs, a)
	}
	return addrs
}

func TestNewWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name    string
		wltName string
		ops     Options
		expect  expect
	}{
		{
			"ok with seed set",
			"test.wlt",
			Options{
				Seed: "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     "skycoin",
					"type":     "deterministic",
					"seed":     "testseed123",
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			"ok with label and seed set",
			"test.wlt",
			Options{
				Label: "wallet1",
				Seed:  "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     "skycoin",
					"type":     "deterministic",
					"seed":     "testseed123",
					"version":  Version,
				},
				err: nil,
			},
		},
		{
			"ok with label, seed and coin set",
			"test.wlt",
			Options{
				Label: "wallet1",
				Coin:  CoinTypeBitcoin,
				Seed:  "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     string(CoinTypeBitcoin),
					"type":     "deterministic",
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
		{
			"ok default crypto type",
			"test.wlt",
			Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Seed:     "testseed123",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      "deterministic",
					"encrypted": "true",
				},
				err: nil,
			},
		},
		{
			"encrypt without password",
			"test.wlt",
			Options{
				Label:   "wallet1",
				Coin:    CoinTypeSkycoin,
				Seed:    "testseed123",
				Encrypt: true,
			},
			expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      "deterministic",
					"encrypted": "true",
				},
				err: ErrMissingPassword,
			},
		},
		{
			"create with no seed",
			"test.wlt",
			Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(CoinTypeSkycoin),
					"type":      "deterministic",
					"encrypted": "true",
				},
				err: ErrMissingSeed,
			},
		},
		{
			"password=pwd encrypt=false",
			"test.wlt",
			Options{
				Label:    "wallet1",
				Coin:     CoinTypeSkycoin,
				Encrypt:  false,
				Seed:     "seed",
				Password: []byte("pwd"),
			},
			expect{
				err: ErrMissingEncrypt,
			},
		},
	}

	for _, tc := range tt {
		// test all supported crypto types
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			if tc.ops.Encrypt {
				tc.ops.CryptoType = ct
			}
			t.Run(name, func(t *testing.T) {
				w, err := NewWallet(tc.wltName, tc.ops)
				require.Equal(t, tc.expect.err, err)
				if err != nil {
					return
				}

				require.Equal(t, tc.ops.Encrypt, w.IsEncrypted())

				if w.IsEncrypted() {
					// Confirms the seeds and entry secrets are all empty
					require.Equal(t, "", w.seed())
					require.Equal(t, "", w.lastSeed())

					for _, e := range w.Entries {
						require.Empty(t, e.Secret)
					}

					// Confirms that secrets field is not empty
					require.NotEmpty(t, w.secrets())
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
			"ok",
			Options{
				Seed: "seed",
			},
			[]byte("pwd"),
			nil,
		},
		{
			"password is nil",
			Options{
				Seed: "seed",
			},
			nil,
			ErrMissingPassword,
		},
		{
			"wallet already encrypted",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			ErrWalletEncrypted,
		},
	}

	for _, tc := range tt {
		for ct := range cryptoTable {
			name := fmt.Sprintf("%v crypto=%v", tc.name, ct)
			if tc.opts.Encrypt {
				tc.opts.CryptoType = ct
			}
			t.Run(name, func(t *testing.T) {
				wltName := newWalletFilename()
				w, err := NewWallet(wltName, tc.opts)
				require.NoError(t, err)

				if !w.IsEncrypted() {
					// Generates 2 addresses
					_, err = w.GenerateAddresses(2)
					require.NoError(t, err)
				}

				err = w.lock(tc.lockPwd, ct)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.True(t, w.IsEncrypted())

				// Checks if the seeds are wiped
				require.Empty(t, w.seed())
				require.Empty(t, w.lastSeed())

				// Checks if the entries are encrypted
				for i := range w.Entries {
					require.Equal(t, cipher.SecKey{}, w.Entries[i].Secret)
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
			"ok",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			[]byte("pwd"),
			nil,
		},
		{
			"unlock with nil password",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			nil,
			ErrMissingPassword,
		},
		{
			"unlock undecrypted wallet",
			Options{
				Seed:    "seed",
				Encrypt: false,
			},
			[]byte("pwd"),
			ErrWalletNotEncrypted,
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
				wlt, err := w.unlock(tc.unlockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.False(t, wlt.IsEncrypted())

				// Checks the seeds
				require.Equal(t, tc.opts.Seed, wlt.seed())

				// Checks the generated addresses
				sd, sks := cipher.GenerateDeterministicKeyPairsSeed([]byte(wlt.seed()), 1)
				require.Equal(t, uint64(1), uint64(len(wlt.Entries)))

				// Checks the last seed
				require.Equal(t, hex.EncodeToString(sd), wlt.lastSeed())

				for i := range wlt.Entries {
					addr := cipher.AddressFromSecKey(sks[i])
					require.Equal(t, addr, wlt.Entries[i].Address)
				}

				// Checks the original seeds
				require.NotEqual(t, tc.opts.Seed, w.seed())

				// Checks if the seckeys in entries of original wallet are empty
				for i := range w.Entries {
					require.Equal(t, cipher.SecKey{}, w.Entries[i].Secret)
				}

				// Checks if the seed and lastSeed in original wallet are sitll empty
				require.Empty(t, w.seed())
				require.Empty(t, w.lastSeed())
			})
		}
	}
}

func TestLockAndUnLock(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			w, err := NewWallet("wallet", Options{
				Label: "wallet",
				Seed:  "seed",
			})
			require.NoError(t, err)
			_, err = w.GenerateAddresses(10)
			require.NoError(t, err)
			require.Len(t, w.Entries, 10)

			// clone the wallet
			cw := w.clone()
			require.Equal(t, w, cw)

			// lock the cloned wallet
			err = cw.lock([]byte("pwd"), ct)
			require.NoError(t, err)

			// unlock the cloned wallet
			ucw, err := cw.unlock([]byte("pwd"))
			require.NoError(t, err)

			require.Equal(t, w, ucw)
		})
	}
}

func makeWallet(t *testing.T, opts Options, addrNum uint64) *Wallet {
	// Create an unlocked wallet, then generate addresses, lock if the options.Encrypt is true.
	preOpts := opts
	opts.Encrypt = false
	opts.Password = nil
	w, err := NewWallet("t.wlt", opts)
	require.NoError(t, err)

	_, err = w.GenerateAddresses(addrNum)
	require.NoError(t, err)
	if preOpts.Encrypt {
		err = w.lock(preOpts.Password, preOpts.CryptoType)
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
		name   string
		file   string
		expect expect
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			expect{
				meta: map[string]string{
					"coin":     string(CoinTypeSkycoin),
					"filename": "test1.wlt",
					"label":    "test3",
					"lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
					"seed":     "buddy fossil side modify turtle door label grunt baby worth brush master",
					"tm":       "1503458909",
					"type":     "deterministic",
					"version":  "0.1",
				},
				err: nil,
			},
		},
		{
			"wallet file doesn't exist",
			"not_exist_file.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("load wallet file failed, wallet not_exist_file.wlt doesn't exist"),
			},
		},
		{
			"invalid wallet: no type",
			"./testdata/invalid_wallets/no_type.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_type.wlt: type field not set"),
			},
		},
		{
			"invalid wallet: invalid type",
			"./testdata/invalid_wallets/err_type.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet err_type.wlt: wallet type invalid"),
			},
		},
		{
			"invalid wallet: no coin",
			"./testdata/invalid_wallets/no_coin.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_coin.wlt: coin field not set"),
			},
		},
		{
			"invalid wallet: no seed",
			"./testdata/invalid_wallets/no_seed.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_seed.wlt: seed field not set"),
			},
		},
		{
			"version=0.2 encrypted=true crypto=scrypt-chacha20poly1305",
			"./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			expect{
				meta: map[string]string{
					"coin":       "skycoin",
					"cryptoType": "scrypt-chacha20poly1305",
					"encrypted":  "true",
					"filename":   "scrypt-chacha20poly1305-encrypted.wlt",
					"label":      "scrypt-chacha20poly1305",
					"lastSeed":   "",
					"seed":       "",
					"type":       "deterministic",
					"version":    "0.2",
				},
				err: nil,
			},
		},
		{
			"version=0.2 encrypted=true crypto=sha256xor",
			"./testdata/sha256xor-encrypted.wlt",
			expect{
				meta: map[string]string{
					"coin":       "skycoin",
					"cryptoType": "sha256-xor",
					"encrypted":  "true",
					"filename":   "sha256xor-encrypted.wlt",
					"label":      "sha256xor",
					"lastSeed":   "",
					"seed":       "",
					"type":       "deterministic",
					"version":    "0.2",
				},
				err: nil,
			},
		},
		{
			"version=0.2 encrypted=flase",
			"./testdata/v2_no_encrypt.wlt",
			expect{
				meta: map[string]string{
					"coin":       "skycoin",
					"cryptoType": "scrypt-chacha20poly1305",
					"encrypted":  "false",
					"filename":   "v2_no_encrypt.wlt",
					"label":      "v2_no_encrypt",
					"lastSeed":   "c79454cf362b3f55e5effce09f664311650a44b9c189b3c8eed1ae9bd696cd9e",
					"secrets":    "",
					"seed":       "seed",
					"type":       "deterministic",
					"version":    "0.2",
				},
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := Load(tc.file)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			for k, v := range tc.expect.meta {
				vv := w.Meta[k]
				require.Equal(t, v, vv)
			}

			if w.IsEncrypted() {
				require.NotEmpty(t, w.Meta[metaSecrets])
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
			"ok with one address",
			Options{
				Seed: "seed",
			},
			1,
			false,
			nil,
		},
		{
			"ok with two address",
			Options{
				Seed: "seed",
			},
			2,
			false,
			nil,
		},
		{
			"ok with three address and generate one address each time",
			Options{
				Seed: "seed",
			},
			2,
			true,
			nil,
		},
		{
			"wallet is encrypted",
			Options{
				Seed:     "seed",
				Encrypt:  true,
				Password: []byte("pwd"),
			},
			2,
			true,
			ErrWalletEncrypted,
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
				require.Equal(t, int(tc.num), len(w.Entries))

				addrs := w.GetAddresses()

				_, keys := cipher.GenerateDeterministicKeyPairsSeed([]byte(tc.opts.Seed), int(tc.num))
				for i, k := range keys {
					a := cipher.AddressFromSecKey(k)
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
			"scrypt-chacha20poly1305 encrytped wallet",
			"./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			"2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM",
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

func TestWalletAddEntry(t *testing.T) {
	test1SecKey, err := cipher.SecKeyFromHex("1fc5396e91e60b9fc613d004ea5bd2ccea17053a12127301b3857ead76fdb93e")
	require.NoError(t, err)

	_, s := cipher.GenerateKeyPair()
	seckeys := []cipher.SecKey{
		test1SecKey,
		s,
	}

	tt := []struct {
		name    string
		wltFile string
		secKey  cipher.SecKey
		err     error
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			seckeys[1],
			nil,
		},
		{
			"dup entry",
			"./testdata/test1.wlt",
			seckeys[0],
			errors.New("duplicate address entry"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := Load(tc.wltFile)
			require.NoError(t, err)
			a := cipher.AddressFromSecKey(tc.secKey)
			p := cipher.PubKeyFromSecKey(tc.secKey)
			require.Equal(t, tc.err, w.AddEntry(Entry{
				Address: a,
				Public:  p,
				Secret:  s,
			}))
		})
	}
}

func TestWalletGuard(t *testing.T) {
	for ct := range cryptoTable {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			validate := func(w *Wallet) {
				require.Equal(t, "", w.seed())
				require.Equal(t, "", w.lastSeed())
				for _, e := range w.Entries {
					require.Equal(t, cipher.SecKey{}, e.Secret)
				}
			}

			w, err := NewWallet("t.wlt", Options{
				Seed:       "seed",
				Encrypt:    true,
				Password:   []byte("pwd"),
				CryptoType: ct,
			})
			require.NoError(t, err)

			require.NoError(t, w.guardUpdate([]byte("pwd"), func(w *Wallet) error {
				require.Equal(t, "seed", w.seed())
				w.setLabel("label")
				return nil
			}))
			require.Equal(t, "label", w.Label())
			validate(w)

			w.guardView([]byte("pwd"), func(w *Wallet) error {
				require.Equal(t, "label", w.Label())
				w.setLabel("new label")
				return nil
			})

			require.Equal(t, "label", w.Label())
			validate(w)

		})
	}
}

type distributeSpendHoursTestCase struct {
	name              string
	inputHours        uint64
	nAddrs            uint64
	haveChange        bool
	expectChangeHours uint64
	expectAddrHours   []uint64
}

var burnFactor2TestCases = []distributeSpendHoursTestCase{
	{
		name:            "no input hours, one addr, no change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "no input hours, two addrs, no change",
		inputHours:      0,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{0, 0},
	},
	{
		name:            "no input hours, one addr, change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      true,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "one input hour, one addr, no change",
		inputHours:      1,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "two input hours, one addr, no change",
		inputHours:      2,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:              "two input hours, one addr, change",
		inputHours:        2,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:              "three input hours, one addr, change",
		inputHours:        3,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:            "three input hours, one addr, no change",
		inputHours:      3,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:            "three input hours, two addrs, no change",
		inputHours:      3,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{1, 0},
	},
	{
		name:            "four input hours, one addr, no change",
		inputHours:      4,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:              "four input hours, one addr, change",
		inputHours:        4,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "four input hours, two addr, change",
		inputHours:        4,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1, 0},
	},
	{
		name:              "30 (divided by 2, odd number) input hours, two addr, change",
		inputHours:        30,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{4, 3},
	},
	{
		name:              "33 (odd number) input hours, two addr, change",
		inputHours:        33,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{4, 4},
	},
	{
		name:              "33 (odd number) input hours, three addr, change",
		inputHours:        33,
		nAddrs:            3,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{3, 3, 2},
	},
}

var burnFactor3TestCases = []distributeSpendHoursTestCase{
	{
		name:            "no input hours, one addr, no change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "no input hours, two addrs, no change",
		inputHours:      0,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{0, 0},
	},
	{
		name:            "no input hours, one addr, change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      true,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "one input hour, one addr, no change",
		inputHours:      1,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "two input hours, one addr, no change",
		inputHours:      2,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:            "three input hours, one addr, no change",
		inputHours:      3,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:              "two input hours, one addr, change",
		inputHours:        2,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:              "three input hours, one addr, change",
		inputHours:        3,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "four input hours, one addr, change",
		inputHours:        4,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:            "four input hours, one addr, no change",
		inputHours:      4,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:            "four input hours, two addrs, no change",
		inputHours:      4,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{1, 1},
	},
	{
		name:            "five input hours, one addr, no change",
		inputHours:      5,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{3},
	},
	{
		name:              "five input hours, one addr, change",
		inputHours:        5,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "five input hours, two addr, change",
		inputHours:        5,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1, 0},
	},
	{
		name:              "32 input hours, two addr, change",
		inputHours:        32,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 11,
		expectAddrHours:   []uint64{5, 5},
	},
	{
		name:              "35 input hours, two addr, change",
		inputHours:        35,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 12,
		expectAddrHours:   []uint64{6, 5},
	},
	{
		name:              "32 input hours, three addr, change",
		inputHours:        32,
		nAddrs:            3,
		haveChange:        true,
		expectChangeHours: 11,
		expectAddrHours:   []uint64{4, 3, 3},
	},
}

func TestWalletDistributeSpendHours(t *testing.T) {
	var cases []distributeSpendHoursTestCase
	switch fee.BurnFactor {
	case 2:
		cases = burnFactor2TestCases
	case 3:
		cases = burnFactor3TestCases
	default:
		t.Fatalf("No test cases defined for fee.BurnFactor=%d", fee.BurnFactor)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			changeHours, addrHours, totalHours := DistributeSpendHours(tc.inputHours, tc.nAddrs, tc.haveChange)
			require.Equal(t, tc.expectChangeHours, changeHours)
			require.Equal(t, tc.expectAddrHours, addrHours)
			require.Equal(t, tc.nAddrs, uint64(len(addrHours)))

			outputHours := changeHours
			for _, h := range addrHours {
				outputHours += h
			}
			require.True(t, tc.inputHours >= outputHours)
			require.Equal(t, outputHours, totalHours)

			if tc.inputHours != 0 {
				err := fee.VerifyTransactionFeeForHours(outputHours, tc.inputHours-outputHours)
				require.NoError(t, err)
			}
		})
	}

	// Tests over range of values
	for inputHours := uint64(0); inputHours <= 1e3; inputHours++ {
		for nAddrs := uint64(1); nAddrs < 16; nAddrs++ {
			for _, haveChange := range []bool{true, false} {
				name := fmt.Sprintf("inputHours=%d nAddrs=%d haveChange=%v", inputHours, nAddrs, haveChange)
				t.Run(name, func(t *testing.T) {
					changeHours, addrHours, totalHours := DistributeSpendHours(inputHours, nAddrs, haveChange)
					require.Equal(t, nAddrs, uint64(len(addrHours)))

					var sumAddrHours uint64
					for _, h := range addrHours {
						sumAddrHours += h
					}

					if haveChange {
						remainingHours := (inputHours - fee.RequiredFee(inputHours))
						splitRemainingHours := remainingHours / 2
						require.True(t, changeHours == splitRemainingHours || changeHours == splitRemainingHours+1)
						require.Equal(t, splitRemainingHours, sumAddrHours)
					} else {
						require.Equal(t, uint64(0), changeHours)
						require.Equal(t, inputHours-fee.RequiredFee(inputHours), sumAddrHours)
					}

					outputHours := sumAddrHours + changeHours
					require.True(t, inputHours >= outputHours)
					require.Equal(t, outputHours, totalHours)

					if inputHours != 0 {
						err := fee.VerifyTransactionFeeForHours(outputHours, inputHours-outputHours)
						require.NoError(t, err)
					}

					// addrHours at the beginning and end of the array should not differ by more than one
					max := addrHours[0]
					min := addrHours[len(addrHours)-1]
					require.True(t, max-min <= 1)
				})
			}
		}
	}
}

func uxBalancesEqual(a, b []UxBalance) bool {
	if len(a) != len(b) {
		return false
	}

	for i, x := range a {
		if x != b[i] {
			return false
		}
	}

	return true
}

func TestWalletSortSpendsLowToHigh(t *testing.T) {
	// UxBalances are sorted with Coins lowest, then following other order rules
	orderedUxb := []UxBalance{
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 5,
			Coins: 1,
			Hours: 0,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 3,
			Coins: 10,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 1,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 100,
			Hours: 100,
		},
	}

	shuffleWorked := false
	nShuffle := 20
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		// Sanity check that shuffling produces a new result
		if !uxBalancesEqual(uxb, orderedUxb) {
			shuffleWorked = true
		}

		sortSpendsCoinsLowToHigh(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsLowToHigh(t, uxb)
	}

	require.True(t, shuffleWorked)

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsHighToLow(uxb)
		verifySortedCoinsHighToLow(t, uxb)
	}
}

func TestWalletSortSpendsHighToLow(t *testing.T) {
	// UxBalances are sorted with Coins highest, then following other order rules
	orderedUxb := []UxBalance{
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 10000,
			Hours: 0,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 10,
			Coins: 1000,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 100,
			Hours: 100,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 3,
			Coins: 10,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 1,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
	}

	shuffleWorked := false
	nShuffle := 20
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		if !uxBalancesEqual(uxb, orderedUxb) {
			shuffleWorked = true
		}

		sortSpendsCoinsHighToLow(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsHighToLow(t, uxb)
	}

	require.True(t, shuffleWorked)

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsHighToLow(uxb)
		verifySortedCoinsHighToLow(t, uxb)
	}
}

func TestWalletChooseSpendsMaximizeUxOuts(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins <= b.Coins
		})
	}
}

func TestWalletChooseSpendsMinimizeUxOuts(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins >= b.Coins
		})
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
			"no file",
			[]wltInfo{},
			map[string]struct{}{},
		},
		{
			"wlt v0.1=1 bak v0.1=1 delete 1 bak",
			[]wltInfo{
				{
					"t1.wlt",
					"0.1",
				},
				{
					"t1.wlt.bak",
					"0.1",
				},
			},
			map[string]struct{}{
				"t1.wlt": struct{}{},
			},
		},
		{
			"wlt v0.1=2 bak v0.1=1 delete 1 bak",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
			},
		},
		{
			"wlt v0.1=3 bak v0.1=1 delete 1 bak",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			"wlt v0.1=3 bak v0.1=2 delete 2 bak",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			"wlt v0.1=3 bak v0.1=3 delete 3 bak",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt": struct{}{},
				"t2.wlt": struct{}{},
				"t3.wlt": struct{}{},
			},
		},
		{
			"wlt v0.1=3 bak v0.1=1 no delete",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t4.wlt.bak": struct{}{},
			},
		},
		{
			"wlt v0.2=3 bak v0.2=1 no delete",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			"wlt v0.1=1 bak v0.1=1 wlt v0.2=2 bak v0.2=2 delete 1 bak",
			[]wltInfo{
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
			map[string]struct{}{
				"t1.wlt":     struct{}{},
				"t2.wlt":     struct{}{},
				"t2.wlt.bak": struct{}{},
				"t3.wlt":     struct{}{},
				"t3.wlt.bak": struct{}{},
			},
		},
		{
			"wlt v0.1=1 bak v0.1=2 wlt v0.2=2 bak v0.2=1 delete 1 bak",
			[]wltInfo{
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
			map[string]struct{}{
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
				})
				require.NoError(t, err)
				w.setVersion(f.version)

				require.NoError(t, w.Save(dir))
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

func makeRandomUxBalances(t *testing.T) []UxBalance {
	// Generate random 0-100 UxBalances
	// Coins 1-10 (must be >0)
	// Hours 0-10
	// BkSeq 0-10
	// Hash random
	// Small ranges are used for Coins, Hours, BkSeq to increase likelihood
	// that they collide and test deeper sorting comparisons

	n := rand.Intn(101)
	uxb := make([]UxBalance, n)

	// Use a random max range for the hours' rand range to ensure enough
	// balances have zero hours
	hasZeroHoursRange := rand.Intn(3) + 1

	for i := 0; i < n; i++ {
		ux := UxBalance{
			Coins: uint64(rand.Intn(10) + 1), // 1-10
			Hours: uint64(rand.Intn(hasZeroHoursRange)),
			BkSeq: uint64(rand.Intn(11)), // 0-10
			Hash:  testutil.RandSHA256(t),
		}

		uxb[i] = ux
	}

	return uxb
}

func verifyChosenCoins(t *testing.T, uxb []UxBalance, coins uint64, chooseSpends func([]UxBalance, uint64) ([]UxBalance, error), cmpCoins func(i, j UxBalance) bool) {
	var haveZero, haveNonzero int
	for _, ux := range uxb {
		if ux.Hours == 0 {
			haveZero++
		} else {
			haveNonzero++
		}
	}

	var totalCoins, totalHours uint64
	for _, ux := range uxb {
		totalCoins += ux.Coins
		totalHours += ux.Hours
	}

	chosen, err := chooseSpends(uxb, coins)

	if coins == 0 {
		testutil.RequireError(t, err, "zero spend amount")
		return
	}

	if len(uxb) == 0 {
		testutil.RequireError(t, err, "no unspents to spend")
		return
	}

	if totalHours == 0 {
		testutil.RequireError(t, err, fee.ErrTxnNoFee.Error())
		return
	}

	if coins > totalCoins {
		testutil.RequireError(t, err, ErrInsufficientBalance.Error())
		return
	}

	require.NoError(t, err)
	require.NotEqual(t, 0, len(chosen))

	// Check that there are no duplicated spends chosen
	uxMap := make(map[UxBalance]struct{}, len(chosen))
	for _, ux := range chosen {
		_, ok := uxMap[ux]
		require.False(t, ok)
		uxMap[ux] = struct{}{}
	}

	// The first chosen spend should have non-zero coin hours
	require.NotEqual(t, uint64(0), chosen[0].Hours)

	// Outputs with zero hours should come before any outputs with non-zero hours,
	// except for the first output
	for i := range chosen {
		if i <= 1 {
			continue
		}

		a := chosen[i-1]
		b := chosen[i]

		if b.Hours == 0 {
			require.Equal(t, uint64(0), a.Hours)
		}
	}

	// The initial UxBalance with hours should have more or equal coins than any other UxBalance with hours
	// If it has equal coins, it should have less hours
	for _, ux := range chosen[1:] {
		if ux.Hours != 0 {
			require.True(t, chosen[0].Coins >= ux.Coins)

			if chosen[0].Coins == ux.Coins {
				require.True(t, chosen[0].Hours <= ux.Hours)
			}
		}
	}

	var zeroBalances, nonzeroBalances []UxBalance
	for _, ux := range chosen[1:] {
		if ux.Hours == 0 {
			zeroBalances = append(zeroBalances, ux)
		} else {
			nonzeroBalances = append(nonzeroBalances, ux)
		}
	}

	// Amongst the UxBalances with zero hours, they should be sorted as specified
	verifySortedCoins(t, zeroBalances, cmpCoins)

	// Amongst the UxBalances with non-zero hours, they should be sorted as specified
	verifySortedCoins(t, nonzeroBalances, cmpCoins)

	// If there are any extra UxBalances with non-zero hours, all of the zeros should have been chosen
	if len(nonzeroBalances) > 0 {
		require.Equal(t, haveZero, len(zeroBalances))
	}

	// Excessive UxBalances to satisfy the amount requested should not be included
	var haveCoins uint64
	for i, ux := range chosen {
		haveCoins += ux.Coins
		if haveCoins >= coins {
			require.Equal(t, len(chosen)-1, i)
		}
	}
}

func verifySortedCoins(t *testing.T, uxb []UxBalance, cmpCoins func(a, b UxBalance) bool) {
	if len(uxb) <= 1 {
		return
	}

	for i := range uxb {
		if i == 0 {
			continue
		}

		a := uxb[i-1]
		b := uxb[i]

		require.True(t, cmpCoins(a, b))

		if a.Coins == b.Coins {
			require.True(t, a.Hours <= b.Hours)

			if a.Hours == b.Hours {
				require.True(t, a.BkSeq <= b.BkSeq)

				if a.BkSeq == b.BkSeq {
					cmp := bytes.Compare(a.Hash[:], b.Hash[:])
					require.True(t, cmp < 0)
				}
			}
		}
	}
}

func verifySortedCoinsLowToHigh(t *testing.T, uxb []UxBalance) {
	verifySortedCoins(t, uxb, func(a, b UxBalance) bool {
		return a.Coins <= b.Coins
	})
}

func verifySortedCoinsHighToLow(t *testing.T, uxb []UxBalance) {
	verifySortedCoins(t, uxb, func(a, b UxBalance) bool {
		return a.Coins >= b.Coins
	})
}
