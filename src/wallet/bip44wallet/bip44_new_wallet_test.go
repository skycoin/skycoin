package bip44wallet

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
	"github.com/stretchr/testify/require"
)

func TestBip44WalletNew(t *testing.T) {
	bip44SkycoinType := bip44.CoinTypeSkycoin
	newBip44Type := bip44.CoinType(1000)

	tt := []struct {
		name                string
		filename            string
		label               string
		seed                string
		seedPassphrase      string
		coinType            meta.CoinType
		bip44CoinType       *bip44.CoinType
		cryptoType          crypto.CryptoType
		expectBip44CoinType bip44.CoinType
		err                 error
	}{
		{
			name:                "skycoin default crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinType("skycoin"),
			cryptoType:          crypto.DefaultCryptoType,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin default crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeBitcoin,
			cryptoType:          crypto.DefaultCryptoType,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin crypto type sha256xor",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeSkycoin,
			cryptoType:          crypto.CryptoTypeSha256Xor,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin crypto type sha256xor",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeBitcoin,
			cryptoType:          crypto.CryptoTypeSha256Xor,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin no crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeSkycoin,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin no crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeBitcoin,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin explicit bip44 coin type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeSkycoin,
			bip44CoinType:       &bip44SkycoinType,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "skycoin new bip44 coin type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            meta.CoinTypeSkycoin,
			bip44CoinType:       &newBip44Type,
			expectBip44CoinType: newBip44Type,
		},
		{
			name:           "no filename",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeSkycoin,
			err:            errors.New("Filename not set"),
		},
		{
			name:           "no coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			err:            errors.New("Missing coin type"),
		},
		{
			name:           "skycoin empty seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           "",
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeSkycoin,
			cryptoType:     crypto.DefaultCryptoType,
			err:            errors.New("Seed missing in unencrypted bip44 wallet"),
		},
		{
			name:           "skycoin invalid seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           invalidBip44Seed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeSkycoin,
			cryptoType:     crypto.DefaultCryptoType,
			err:            errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
		{
			name:           "new coin type, no bi44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinType("unknown"),
			err:            errors.New("Missing bip44 coin type"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
				Filename:       tc.filename,
				Label:          tc.label,
				Seed:           tc.seed,
				SeedPassphrase: tc.seedPassphrase,
				CoinType:       tc.coinType,
				CryptoType:     tc.cryptoType,
				Bip44CoinType:  tc.bip44CoinType,
			})

			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			require.Equal(t, tc.filename, w.Meta.Filename())
			require.Equal(t, tc.label, w.Meta.Label())
			require.Equal(t, tc.seed, w.Meta.Seed())
			require.Equal(t, tc.seedPassphrase, w.Meta.SeedPassphrase())
			require.Equal(t, tc.coinType, w.Meta.Coin())
			require.Equal(t, walletType, w.Meta.Type())
			require.False(t, w.Meta.IsEncrypted())
			require.NotEmpty(t, w.Meta.Timestamp())
			require.NotNil(t, w.decoder)
			bip44Coin := w.Bip44Coin()
			require.Equal(t, tc.expectBip44CoinType, *bip44Coin)
			require.Empty(t, w.Meta.Secrets())

			if tc.cryptoType != "" {
				require.Equal(t, tc.cryptoType, w.Meta.CryptoType())
			} else {
				require.Equal(t, crypto.DefaultCryptoType, w.Meta.CryptoType())
			}
		})
	}
}

func TestWalletCreateAccount(t *testing.T) {
	w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       meta.CoinTypeSkycoin,
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	ai, err = w.NewAccount("account2")
	require.Equal(t, uint32(1), ai)

	require.Equal(t, uint32(2), w.accounts.len())
}

func TestWalletAccountCreateAddresses(t *testing.T) {
	w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       meta.CoinTypeSkycoin,
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)
	require.Equal(t, uint32(0), ai)

	addrs, err := w.NewExternalAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr := make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinExternalAddresses[:2], addrsStr)

	addrs, err = w.NewChangeAddresses(ai, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr = make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinChangeAddresses[:2], addrsStr)
}

func TestBip44WalletLock(t *testing.T) {
	w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       meta.CoinTypeSkycoin,
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.NewExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.NewChangeAddresses(ai, 2)
	require.NoError(t, err)

	err = w.Lock([]byte("123456"))
	require.NoError(t, err)

	require.Empty(t, w.Seed())
	require.Empty(t, w.SeedPassphrase())
	require.NotEmpty(t, w.Secrets())
	require.True(t, w.IsEncrypted())

	// confirms that no secrets exist in the accounts
	ss := make(secrets.Secrets)
	w.accounts.packSecrets(ss)
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
	w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       meta.CoinTypeSkycoin,
		CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.NewExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.NewChangeAddresses(ai, 2)
	require.NoError(t, err)

	cw := w.Clone()

	err = cw.Lock([]byte("123456"))
	require.NoError(t, err)

	// generates addresses after locking
	_, err = cw.NewExternalAddresses(ai, 2)
	require.NoError(t, err)
	_, err = cw.NewChangeAddresses(ai, 3)
	require.NoError(t, err)

	// unlock with wrong password
	_, err = cw.Unlock([]byte("12345"))
	require.Equal(t, errors.New("Invalid password"), err)

	// unlock with correct password
	wlt, err := cw.Unlock([]byte("123456"))
	require.NoError(t, err)

	el, err := wlt.accounts.entriesLen(ai, bip44.ExternalChainIndex)
	require.NoError(t, err)
	require.Equal(t, uint32(4), el)

	cl, err := wlt.accounts.entriesLen(ai, bip44.ChangeChainIndex)
	require.NoError(t, err)
	require.Equal(t, uint32(5), cl)

	// confirms that unlocking wallet won't lose data
	require.Empty(t, wlt.Secrets())
	require.False(t, wlt.IsEncrypted())
	require.Equal(t, w.Seed(), wlt.Seed())
	require.Equal(t, w.SeedPassphrase(), wlt.SeedPassphrase())

	// pack the origin wallet's secrets
	originSS := make(secrets.Secrets)
	w.accounts.packSecrets(originSS)

	// pack the unlocked wallet's secrets
	ss := make(secrets.Secrets)
	wlt.accounts.packSecrets(ss)

	// compare these two secrets, they should have the same keys and values
	require.Equal(t, len(originSS)+5, len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}

func TestBip44WalletNewSerializeDeserialize(t *testing.T) {
	w, err := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       meta.CoinTypeSkycoin,
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.NewExternalAddresses(ai, 2)
	require.NoError(t, err)

	_, err = w.NewChangeAddresses(ai, 2)
	require.NoError(t, err)

	b, err := w.Serialize()
	require.NoError(t, err)
	t.Log(string(b))

	wlt := Bip44WalletNew{}
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
	require.Equal(t, w.accounts.len(), wlt.accounts.len())
	originSS := make(secrets.Secrets)
	ss := make(secrets.Secrets)
	w.accounts.packSecrets(originSS)
	wlt.accounts.packSecrets(ss)

	require.Equal(t, len(originSS), len(ss))
	for k, v := range originSS {
		vv, ok := ss[k]
		require.True(t, ok)
		require.Equal(t, v, vv)
	}
}

func TestBip44WalletDiffNoneSecrets(t *testing.T) {
	tt := []struct {
		name                       string
		options                    Bip44WalletCreateOptions
		changeWalletFunc           func(t *testing.T, w *Bip44WalletNew)
		password                   []byte
		err                        error
		expectedMetaChange         meta.Meta
		expectedNewExternalAddrNum int
		expectedNewChangeAddrNum   int
	}{
		{
			name: "new external addresses",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
		},
		{
			name: "new change addresses",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewChangeAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewChangeAddrNum: 3,
		},
		{
			name: "new external and change addresses",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
				_, err = w.NewChangeAddresses(0, 1)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
			expectedNewChangeAddrNum:   1,
		},
		{
			name: "change label",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				w.Meta[meta.MetaLabel] = "label_changed"
			},
			expectedMetaChange: meta.Meta{meta.MetaLabel: "label_changed"},
		},
		{
			name: "change nothing",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
			},
		},

		{
			name: "new external addresses, lock",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
		},
		{
			name: "new change addresses, lock",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewChangeAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewChangeAddrNum: 3,
		},
		{
			name: "new external and change addresses, lock",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
				_, err = w.NewChangeAddresses(0, 1)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
			expectedNewChangeAddrNum:   1,
		},
		{
			name: "change label, lock",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
				w.Meta[meta.MetaLabel] = "label_changed"
			},
			expectedMetaChange: meta.Meta{meta.MetaLabel: "label_changed"},
		},
		{
			name: "change nothing, lock",
			options: Bip44WalletCreateOptions{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       meta.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Bip44WalletNew) {
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewBip44WalletNew(tc.options)
			require.NoError(t, err)

			ai, err := w.NewAccount("account1")
			require.NoError(t, err)

			_, err = w.NewExternalAddresses(ai, 2)
			require.NoError(t, err)

			_, err = w.NewChangeAddresses(ai, 2)
			require.NoError(t, err)

			if len(tc.password) > 0 {
				require.NoError(t, w.Lock(tc.password))
			}

			// Apply changes
			w2 := w.Clone()
			tc.changeWalletFunc(t, &w2)

			diff, err := w.DiffNoneSecrets(&w2)

			require.Empty(t, diff.Meta[meta.MetaSecrets])
			require.Empty(t, diff.Meta[meta.MetaSeed])
			require.Empty(t, diff.Meta[meta.MetaSeedPassphrase])
			require.Empty(t, diff.Meta[meta.MetaEncrypted])
			require.Empty(t, diff.Meta[meta.MetaAccountsHash])

			require.Equal(t, len(tc.expectedMetaChange), len(diff.Meta))
			for k, v := range diff.Meta {
				require.Equal(t, tc.expectedMetaChange[k], v)
			}

			require.Equal(t, tc.expectedNewExternalAddrNum, diff.Accounts[0].NewExternalAddressNum)
			require.Equal(t, tc.expectedNewChangeAddrNum, diff.Accounts[0].NewChangeAddressNum)
		})
	}
}
