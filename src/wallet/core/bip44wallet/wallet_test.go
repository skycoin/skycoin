package bip44wallet

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
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
		coinType            wallet.CoinType
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
			coinType:            wallet.CoinType("skycoin"),
			cryptoType:          crypto.DefaultCryptoType,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin default crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeBitcoin,
			cryptoType:          crypto.DefaultCryptoType,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin crypto type sha256xor",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeSkycoin,
			cryptoType:          crypto.CryptoTypeSha256Xor,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin crypto type sha256xor",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeBitcoin,
			cryptoType:          crypto.CryptoTypeSha256Xor,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin no crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeSkycoin,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "bitcoin no crypto type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeBitcoin,
			expectBip44CoinType: bip44.CoinTypeBitcoin,
		},
		{
			name:                "skycoin explicit bip44 coin type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeSkycoin,
			bip44CoinType:       &bip44SkycoinType,
			expectBip44CoinType: bip44.CoinTypeSkycoin,
		},
		{
			name:                "skycoin new bip44 coin type",
			filename:            "test.wlt",
			label:               "test",
			seed:                testSeed,
			seedPassphrase:      testSeedPassphrase,
			coinType:            wallet.CoinTypeSkycoin,
			bip44CoinType:       &newBip44Type,
			expectBip44CoinType: newBip44Type,
		},
		{
			name:           "no filename",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeSkycoin,
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
			coinType:       wallet.CoinTypeSkycoin,
			cryptoType:     crypto.DefaultCryptoType,
			err:            errors.New("Seed missing in unencrypted bip44 wallet"),
		},
		{
			name:           "skycoin invalid seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           invalidBip44Seed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinTypeSkycoin,
			cryptoType:     crypto.DefaultCryptoType,
			err:            errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
		{
			name:           "new coin type, no bi44 coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       wallet.CoinType("unknown"),
			err:            errors.New("Missing bip44 coin type"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewWallet(Options{
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
	w, err := NewWallet(Options{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       wallet.CoinTypeSkycoin,
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
	w, err := NewWallet(Options{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       wallet.CoinTypeSkycoin,
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
	w, err := NewWallet(Options{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       wallet.CoinTypeSkycoin,
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
	ss := make(wallet.Secrets)
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
	w, err := NewWallet(Options{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       wallet.CoinTypeSkycoin,
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
	originSS := make(wallet.Secrets)
	w.accounts.packSecrets(originSS)

	// pack the unlocked wallet's secrets
	ss := make(wallet.Secrets)
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
	w, err := NewWallet(Options{
		Filename:       "test.wlt",
		Label:          "test",
		Seed:           testSeed,
		SeedPassphrase: testSeedPassphrase,
		CoinType:       wallet.CoinTypeSkycoin,
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
	require.Equal(t, w.accounts.len(), wlt.accounts.len())
	originSS := make(wallet.Secrets)
	ss := make(wallet.Secrets)
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
		options                    Options
		changeWalletFunc           func(t *testing.T, w *Wallet)
		password                   []byte
		err                        error
		expectedMetaChange         wallet.Meta
		expectedNewExternalAddrNum int
		expectedNewChangeAddrNum   int
	}{
		{
			name: "new external addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
		},
		{
			name: "new change addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				_, err := w.NewChangeAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewChangeAddrNum: 3,
		},
		{
			name: "new external and change addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Wallet) {
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
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				w.Meta[wallet.MetaLabel] = "label_changed"
			},
			expectedMetaChange: wallet.Meta{wallet.MetaLabel: "label_changed"},
		},
		{
			name: "change nothing",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			changeWalletFunc: func(t *testing.T, w *Wallet) {
			},
		},

		{
			name: "new external addresses, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				_, err := w.NewExternalAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewExternalAddrNum: 3,
		},
		{
			name: "new change addresses, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				_, err := w.NewChangeAddresses(0, 3)
				require.NoError(t, err)
			},
			expectedNewChangeAddrNum: 3,
		},
		{
			name: "new external and change addresses, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
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
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				w.Meta[wallet.MetaLabel] = "label_changed"
			},
			expectedMetaChange: wallet.Meta{wallet.MetaLabel: "label_changed"},
		},
		{
			name: "change nothing, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
			},
		},
		{
			name: "change secrets, should not be collected by diff",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				w.Meta[wallet.MetaSecrets] = "changed secrets"
				w.Meta[wallet.MetaSeed] = "new seed"
				w.Meta[wallet.MetaSeedPassphrase] = "new seed passphrase"
				w.Meta[wallet.MetaEncrypted] = "true"
				w.Meta[wallet.MetaAccountsHash] = "new accounts hash"
			},
		},
		{
			name: "change secrets, should not be collected by diff, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				w.Meta[wallet.MetaSecrets] = "changed secrets"
				w.Meta[wallet.MetaSeed] = "new seed"
				w.Meta[wallet.MetaSeedPassphrase] = "new seed passphrase"
				w.Meta[wallet.MetaEncrypted] = "false"
				w.Meta[wallet.MetaAccountsHash] = "new accounts hash"
			},
		},
		{
			name: "change immutable meta, should not be collected by diff, lock",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			password: []byte("12345"),
			changeWalletFunc: func(t *testing.T, w *Wallet) {
				im := immutableMeta()
				for k := range im {
					w.Meta[k] = w.Meta[k] + "-changed"
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewWallet(tc.options)
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
			require.NoError(t, err)

			require.Empty(t, diff.Meta[wallet.MetaSecrets])
			require.Empty(t, diff.Meta[wallet.MetaSeed])
			require.Empty(t, diff.Meta[wallet.MetaSeedPassphrase])
			require.Empty(t, diff.Meta[wallet.MetaEncrypted])
			require.Empty(t, diff.Meta[wallet.MetaAccountsHash])

			require.Equal(t, len(tc.expectedMetaChange), len(diff.Meta))
			for k, v := range diff.Meta {
				require.Equal(t, tc.expectedMetaChange[k], v)
			}

			require.Equal(t, tc.expectedNewExternalAddrNum, diff.Accounts[0].NewExternalAddressNum)
			require.Equal(t, tc.expectedNewChangeAddrNum, diff.Accounts[0].NewChangeAddressNum)
		})
	}
}

func TestBip44WalletCommitDiffs(t *testing.T) {
	tt := []struct {
		name                       string
		options                    Options
		password                   []byte
		diffs                      *WalletDiff
		err                        error
		expectedMeta               wallet.Meta
		expectedNewExternalAddrNum int
		expectedNewChangeAddrNum   int
	}{
		{
			name: "new external addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Accounts: []AccountDiff{{NewExternalAddressNum: 1}},
			},
			expectedNewExternalAddrNum: 1,
		},
		{
			name: "new 5 external addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Accounts: []AccountDiff{{NewExternalAddressNum: 5}},
			},
			expectedNewExternalAddrNum: 5,
		},
		{
			name: "new change addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Accounts: []AccountDiff{{NewChangeAddressNum: 1}},
			},
			expectedNewChangeAddrNum: 1,
		},
		{
			name: "new 5 change addresses",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Accounts: []AccountDiff{{NewChangeAddressNum: 5}},
			},
			expectedNewChangeAddrNum: 5,
		},
		{
			name: "change label",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Meta: wallet.Meta{wallet.MetaLabel: "test_changed"},
			},
			expectedMeta: wallet.Meta{wallet.MetaLabel: "test_changed"},
		},
		{
			name: "change secrets, seed, seedPassphrase",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Meta: wallet.Meta{
					wallet.MetaSecrets:        "secrets_changed",
					wallet.MetaSeed:           "seed_changed",
					wallet.MetaSeedPassphrase: "seed_passphrase_changed",
				},
			},
			expectedMeta: wallet.Meta{
				wallet.MetaSecrets:        "secrets_changed",
				wallet.MetaSeed:           "seed_changed",
				wallet.MetaSeedPassphrase: "seed_passphrase_changed",
			},
		},
		{
			name: "change immutable filename, coin type, crypto type, wallet type, no commit",
			options: Options{
				Filename:       "test.wlt",
				Label:          "test",
				Seed:           testSeed,
				SeedPassphrase: testSeedPassphrase,
				CoinType:       wallet.CoinTypeSkycoin,
				CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
			},
			diffs: &WalletDiff{
				Meta: wallet.Meta{
					wallet.MetaFilename:   "test_changed.wlt",
					wallet.MetaCoin:       "coin_changed",
					wallet.MetaBip44Coin:  "bip44_coin_type_changed",
					wallet.MetaCryptoType: "crypto_changed",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewWallet(tc.options)
			require.NoError(t, err)

			aid, err := w.NewAccount("default")
			require.NoError(t, err)

			// Lock wallet if password is provided
			if len(tc.password) > 0 {
				err = w.Lock(tc.password)
				require.NoError(t, err)
			}

			err = w.CommitDiffs(tc.diffs)
			require.NoError(t, err)

			// Confirms that the meta data is applied
			for k, v := range tc.expectedMeta {
				require.Equal(t, v, w.Meta[k])
			}

			// Confirms that external chain addresses length is matched
			el, err := w.ExternalEntriesLen(aid)
			require.Equal(t, uint32(tc.expectedNewExternalAddrNum), el)

			// Confirms that change chain addresses length is matched
			cl, err := w.ChangeEntriesLen(aid)
			require.NoError(t, err)
			require.Equal(t, uint32(tc.expectedNewChangeAddrNum), cl)
		})
	}
}
