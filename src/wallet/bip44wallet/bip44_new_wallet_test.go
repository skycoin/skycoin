package bip44wallet

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
	"github.com/stretchr/testify/require"
)

func TestBip44WalletNew(t *testing.T) {
	tt := []struct {
		name           string
		filename       string
		label          string
		seed           string
		seedPassphrase string
		coinType       meta.CoinType
		cryptoType     crypto.CryptoType
		err            error
	}{
		{
			name:           "skycoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinType("skycoin"),
			cryptoType:     crypto.DefaultCryptoType,
		},
		{
			name:           "bitcoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeBitcoin,
			cryptoType:     crypto.DefaultCryptoType,
		},
		{
			name:           "skycoin crypto type sha256xor",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeSkycoin,
			cryptoType:     crypto.CryptoTypeSha256Xor,
		},
		{
			name:           "bitcoin crypto type sha256xor",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeBitcoin,
			cryptoType:     crypto.CryptoTypeSha256Xor,
		},
		{
			name:           "skycoin no crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeSkycoin,
		},
		{
			name:           "bitcoin no crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       meta.CoinTypeBitcoin,
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
			err:            errors.New("Coin field not set"),
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
			})

			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			require.Equal(t, wallet.Version, w.Meta.Version())
			require.Equal(t, tc.filename, w.Meta.Filename())
			require.Equal(t, tc.label, w.Meta.Label())
			require.Equal(t, tc.seed, w.Meta.Seed())
			require.Equal(t, tc.seedPassphrase, w.Meta.SeedPassphrase())
			require.Equal(t, tc.coinType, w.Meta.Coin())
			require.Equal(t, wallet.WalletTypeBip44, w.Meta.Type())
			require.False(t, w.Meta.IsEncrypted())
			require.NotEmpty(t, w.Meta.Timestamp())
			require.NotNil(t, w.decoder)
			bip44Coin, ok := w.Bip44Coin()
			require.True(t, ok)
			require.Equal(t, resolveCoinAdapter(tc.coinType).Bip44CoinType(), bip44Coin)
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

	addrs, err := w.NewAddresses(ai, bip44.ExternalChainIndex, 2)
	require.NoError(t, err)
	require.Equal(t, 2, len(addrs))
	addrsStr := make([]string, 2)
	for i, a := range addrs {
		addrsStr[i] = a.String()
	}
	require.Equal(t, testSkycoinExternalAddresses[:2], addrsStr)

	addrs, err = w.NewAddresses(ai, bip44.ChangeChainIndex, 2)
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

	_, err = w.NewAddresses(ai, bip44.ExternalChainIndex, 2)
	require.NoError(t, err)

	_, err = w.NewAddresses(ai, bip44.ChangeChainIndex, 2)
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
	})
	require.NoError(t, err)

	ai, err := w.NewAccount("account1")
	require.NoError(t, err)

	_, err = w.NewAddresses(ai, bip44.ExternalChainIndex, 2)
	require.NoError(t, err)

	_, err = w.NewAddresses(ai, bip44.ChangeChainIndex, 2)
	require.NoError(t, err)

	cw := w.clone()

	err = cw.Lock([]byte("123456"))
	require.NoError(t, err)

	// unlock with wrong password
	_, err = cw.Unlock([]byte("12345"))
	require.Equal(t, wallet.ErrInvalidPassword, err)

	// unlock with the correct password
	wlt, err := cw.Unlock([]byte("123456"))
	require.NoError(t, err)

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
	require.Equal(t, len(originSS), len(ss))
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

	_, err = w.NewAddresses(ai, bip44.ExternalChainIndex, 2)
	require.NoError(t, err)

	_, err = w.NewAddresses(ai, bip44.ChangeChainIndex, 2)
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
