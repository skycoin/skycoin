package deterministic

import (
	"fmt"
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

	//testDefaultMnemonicSeed := bip39.MustNewDefaultMnemonic()

	tt := []struct {
		name    string
		wltName string
		seed    string
		label   string
		options []wallet.Option
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
			options: []wallet.Option{
				CoinType(wallet.CoinTypeBitcoin),
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
			options: []wallet.Option{
				Encrypt(true),
				Password([]byte("pwd")),
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
			options: []wallet.Option{
				Encrypt(true),
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
			options: []wallet.Option{
				Encrypt(true),
				Password([]byte("pwd")),
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
			options: []wallet.Option{
				Encrypt(false),
				Password([]byte("pwd")),
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

			// apply the options to an temporary wallet
			opts := tc.options
			mOpts := &moreOptions{}
			for _, opt := range tc.options {
				opt(mOpts)
			}

			if mOpts.Encrypt {
				// append the insecure crypto type
				opts = append(opts, CryptoType(ct))
			}

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
					//entries, err := w.GetEntries()
					entries := w.GetEntries()
					//require.NoError(t, err)

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
		opts    []wallet.Option
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
		for _, ct := range crypto.TypesInsecure() {
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
