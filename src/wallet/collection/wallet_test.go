package collection

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/crypto"
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

func TestNewWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name    string
		wltName string
		label   string
		opts    []wallet.Option
		expect  expect
	}{
		{
			name:    "ok all defaults",
			wltName: "test.wlt",
			label:   "",
			expect: expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     string(wallet.CoinTypeSkycoin),
					"type":     wallet.WalletTypeCollection,
					"version":  wallet.Version,
				},
				err: nil,
			},
		},
		{
			name:    "ok with label,and coin set, collection",
			wltName: "test.wlt",
			label:   "test",
			opts: []wallet.Option{
				wallet.OptionCoinType(wallet.CoinTypeBitcoin),
			},
			expect: expect{
				meta: map[string]string{
					"label":    "test",
					"filename": "test.wlt",
					"coin":     string(wallet.CoinTypeBitcoin),
					"type":     wallet.WalletTypeCollection,
				},
				err: nil,
			},
		},
		{
			name:    "ok default crypto type, collection",
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
					"type":      wallet.WalletTypeCollection,
					"encrypted": "true",
				},
				err: nil,
			},
		},
		{
			name:    "encrypt without password, collection",
			wltName: "test.wlt",
			label:   "wallet1",
			opts: []wallet.Option{
				wallet.OptionEncrypt(true),
			},
			expect: expect{
				meta: map[string]string{
					"label":     "wallet1",
					"coin":      string(wallet.CoinTypeSkycoin),
					"type":      wallet.WalletTypeCollection,
					"encrypted": "true",
				},
				err: wallet.ErrMissingPassword,
			},
		},
		{
			name:    "password=pwd encrypt=false",
			wltName: "test.wlt",
			label:   "test",
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
			opts := tc.opts
			opts = append(opts, wallet.OptionCryptoType(ct))

			t.Run(name, func(t *testing.T) {
				w, err := NewWallet(tc.wltName, tc.label, opts...)
				require.Equal(t, tc.expect.err, err, fmt.Sprintf("want:%v get:%v", tc.expect.err, err))
				if err != nil {
					return
				}

				// require.Equal(t, tc.opts.Encrypt, w.IsEncrypted())
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
				wltName := wallet.NewWalletFilename()
				w, err := NewWallet(wltName, "test", opts...)
				require.NoError(t, err)

				// add entries
				for _, e := range testSkycoinEntries {
					w.AddEntry(e)
				}

				err = w.Lock(tc.lockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.True(t, w.IsEncrypted())

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

func TestWalletAddEntry(t *testing.T) {
	w, err := NewWallet(
		"collection",
		"collection",
		wallet.OptionCryptoType(crypto.CryptoTypeScryptChacha20poly1305Insecure))
	require.NoError(t, err)

	keys, err := cipher.GenerateDeterministicKeyPairs([]byte("testseed123"), 5)
	require.NoError(t, err)

	entries := [5]wallet.Entry{}
	for i := range keys {
		pubkey, err := cipher.PubKeyFromSecKey(keys[i])
		require.NoError(t, err)
		addr, err := cipher.AddressFromSecKey(keys[i])
		require.NoError(t, err)
		entry := wallet.Entry{
			Address: addr,
			Public:  pubkey,
			Secret:  keys[i],
		}
		entries[i] = entry
		err = w.AddEntry(entry)
		require.NoError(t, err)
	}

	el, err := w.EntriesLen()
	require.NoError(t, err)
	require.Equal(t, 5, el)

	es, err := w.GetEntries()
	require.NoError(t, err)
	for i := range es {
		require.Equal(t, entries[i], es[i])
	}

	// try to add dup entry
	err = w.AddEntry(entries[0])
	require.EqualError(t, err, "wallet already contains entry with this address")

	// try to add entry with invalid seckey
	invalidKey := keys[0]
	invalidKey[len(invalidKey)-1] = 0
	err = w.AddEntry(wallet.Entry{Secret: invalidKey})
	require.EqualError(t, err, "invalid public key for secret key")

	// mismatch public key
	entry := entries[0]
	entry.Public = entries[1].Public
	err = w.AddEntry(entry)
	require.EqualError(t, err, "invalid public key for secret key")

	// lock the wallet and try to add an entry
	err = w.Lock([]byte("password"))
	require.NoError(t, err)

	err = w.AddEntry(wallet.Entry{})
	require.Equal(t, wallet.ErrWalletEncrypted, err)
}

func TestWalletUnlock(t *testing.T) {
	tt := []struct {
		name      string
		isLock    bool
		opts      []wallet.Option
		unlockPwd []byte
		err       error
	}{
		{
			name:      "ok",
			isLock:    true,
			unlockPwd: []byte("pwd"),
		},
		{
			name:      "unlock with nil password",
			isLock:    true,
			unlockPwd: nil,
			err:       wallet.ErrMissingPassword,
		},
		{
			name:      "unlock with wrong password",
			isLock:    true,
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

			keys, err := cipher.GenerateDeterministicKeyPairs([]byte("testseed123"), 5)
			require.NoError(t, err)

			entries := [5]wallet.Entry{}
			for i := range keys {
				pubkey, err := cipher.PubKeyFromSecKey(keys[i])
				require.NoError(t, err)
				addr, err := cipher.AddressFromSecKey(keys[i])
				require.NoError(t, err)
				entry := wallet.Entry{
					Address: addr,
					Public:  pubkey,
					Secret:  keys[i],
				}
				entries[i] = entry
			}

			opts := tc.opts
			opts = append(opts, wallet.OptionCryptoType(ct))

			t.Run(name, func(t *testing.T) {
				w, err := NewWallet("test.wlt", "test", opts...)
				require.NoError(t, err)

				// Add entries
				for _, e := range entries {
					err = w.AddEntry(e)
					require.NoError(t, err)
				}

				if tc.isLock {
					err = w.Lock([]byte("pwd"))
					require.NoError(t, err)
				}

				// Tests the unlock method
				wlt, err := w.Unlock(tc.unlockPwd)
				require.Equal(t, tc.err, err)
				if err != nil {
					return
				}

				require.False(t, wlt.IsEncrypted())

				es, err := wlt.GetEntries()
				require.NoError(t, err)
				for i, e := range es {
					require.Equal(t, entries[i], e)
				}
			})
		}
	}
}

func TestLockAndUnLock(t *testing.T) {
	keys, err := cipher.GenerateDeterministicKeyPairs([]byte("testseed123"), 5)
	require.NoError(t, err)

	entries := [5]wallet.Entry{}
	for i := range keys {
		pubkey, err := cipher.PubKeyFromSecKey(keys[i])
		require.NoError(t, err)
		addr, err := cipher.AddressFromSecKey(keys[i])
		require.NoError(t, err)
		entry := wallet.Entry{
			Address: addr,
			Public:  pubkey,
			Secret:  keys[i],
		}
		entries[i] = entry
	}

	for _, ct := range crypto.TypesInsecure() {
		t.Run(fmt.Sprintf("crypto=%v", ct), func(t *testing.T) {
			w, err := NewWallet("wallet.wlt", "test", wallet.OptionCryptoType(ct))
			require.NoError(t, err)
			for _, e := range entries {
				err = w.AddEntry(e)
				require.NoError(t, err)
			}

			// clone the wallet
			cw := w.Clone()
			require.Equal(t, w, cw)

			// lock the cloned wallet
			err = cw.Lock([]byte("pwd"))
			require.NoError(t, err)

			require.True(t, cw.IsEncrypted())

			// Checks if the entries are encrypted
			es, err := cw.GetEntries()
			require.NoError(t, err)

			for _, e := range es {
				require.Equal(t, cipher.SecKey{}, e.Secret)
			}

			// unlock the cloned wallet
			ucw, err := cw.Unlock([]byte("pwd"))
			require.NoError(t, err)

			require.Equal(t, w, ucw)
		})
	}
}

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

func TestWalletGetEntry(t *testing.T) {
	tt := []struct {
		name    string
		wltFile string
		address string
		err     error
	}{
		{
			"ok",
			"./testdata/test-collection.wlt",
			"JUdRuTiqD1mGcw358twMg3VPpXpzbkdRvJ",
			nil,
		},
		{
			"entry not exist",
			"./testdata/test-collection.wlt",
			"2ULfxDUuenUY5V4Pr8whmoAwFdUseXNyjXC",
			wallet.ErrEntryNotFound,
		},
		{
			"scrypt-chacha20poly1305 encrypted wallet",
			"./testdata/scrypt-chacha20poly1305-encrypted.wlt",
			"JUdRuTiqD1mGcw358twMg3VPpXpzbkdRvJ",
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

func TestWalletSerialize(t *testing.T) {
	w, err := NewWallet("test.wlt", "test")
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		require.NoError(t, w.AddEntry(testSkycoinEntries[i]))
	}

	w.SetTimestamp(0)
	b, err := w.Serialize()
	require.NoError(t, err)

	// load wallet file and compare
	fb, err := ioutil.ReadFile("./testdata/wallet_serialize.wlt")
	require.NoError(t, err)
	fb = bytes.TrimRight(fb, "\n")
	require.Equal(t, fb, b)

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
