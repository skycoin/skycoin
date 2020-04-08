package wallet

import (
	"errors"
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/stretchr/testify/require"
)

func TestBip44WalletNew(t *testing.T) {
	tt := []struct {
		name           string
		filename       string
		label          string
		seed           string
		seedPassphrase string
		coinType       CoinType
		cryptoType     CryptoType
		err            error
	}{
		{
			name:           "skycoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
			cryptoType:     DefaultCryptoType,
		},
		{
			name:           "bitcoin default crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeBitcoin,
			cryptoType:     DefaultCryptoType,
		},
		{
			name:           "skycoin crypto type sha256xor",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
			cryptoType:     CryptoTypeSha256Xor,
		},
		{
			name:           "bitcoin crypto type sha256xor",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeBitcoin,
			cryptoType:     CryptoTypeSha256Xor,
		},
		{
			name:           "skycoin no crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
		},
		{
			name:           "bitcoin no crypto type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeBitcoin,
		},
		{
			name:           "no filename",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
			err:            errors.New("filename not set"),
		},
		{
			name:           "no coin type",
			filename:       "test.wlt",
			label:          "test",
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			err:            errors.New("coin field not set"),
		},
		{
			name:           "skycoin empty seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           "",
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
			cryptoType:     DefaultCryptoType,
			err:            errors.New("seed missing in unencrypted bip44 wallet"),
		},
		{
			name:           "skycoin invalid seed",
			filename:       "test.wlt",
			label:          "test",
			seed:           invalidBip44Seed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
			cryptoType:     DefaultCryptoType,
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
			require.Equal(t, Version, w.Meta.Version())
			require.Equal(t, tc.filename, w.Meta.Filename())
			require.Equal(t, tc.label, w.Meta.Label())
			require.Equal(t, tc.seed, w.Meta.Seed())
			require.Equal(t, tc.seedPassphrase, w.Meta.SeedPassphrase())
			require.Equal(t, tc.coinType, w.Meta.Coin())
			require.Equal(t, WalletTypeBip44, w.Meta.Type())
			require.False(t, w.Meta.IsEncrypted())
			require.NotEmpty(t, w.Meta.Timestamp())
			require.NotNil(t, w.decoder)
			require.Equal(t, resolveCoinAdapter(tc.coinType).Bip44CoinType(), w.Meta.Bip44Coin())
			require.Empty(t, w.Meta.Secrets())

			if tc.cryptoType != "" {
				require.Equal(t, tc.cryptoType, w.Meta.CryptoType())
			} else {
				require.Equal(t, DefaultCryptoType, w.Meta.CryptoType())
			}
		})
	}

	v, err := w.MarshalToJSON()
	require.NoError(t, err)
	t.Log(string(v))
}

func makeAccount(t *testing.T, seed string, password string) *bip44.Account {
	s, err := bip39.NewSeed(seed, password)
	require.NoError(t, err)

	c, err := bip44.NewCoin(s, bip44.CoinTypeSkycoin)
	require.NoError(t, err)
	a, err := c.Account(0)
	require.NoError(t, err)
	return a
}

func TestMakeChainKeys(t *testing.T) {
	a := makeAccount(t, testSeed, testSeedPassphrase)
	p1PubKey, err := a.NewPublicChildKey(0)
	require.NoError(t, err)
	p2PubKey, err := a.NewPublicChildKey(1)
	require.NoError(t, err)

	// confirms that multiple times of calling of makeChainPubKeys will not affect the result
	for i := 0; i < 10; i++ {
		pexternalKey, pchangeKey, err := makeChainPubKeys(a)
		require.NoError(t, err)

		require.Equal(t, p1PubKey, pexternalKey)
		require.Equal(t, p2PubKey, pchangeKey)
	}
}

func TestNewBip44Account(t *testing.T) {
	tt := []struct {
		name        string
		accountName string
		index       uint32
		coinType    bip44.CoinType
		err         error
	}{
		{
			name:        "index 0, test, skycoin",
			accountName: "test",
			index:       0,
			coinType:    bip44.CoinTypeSkycoin,
		},
		{
			name:        "index 1, test1, skycoin",
			accountName: "test1",
			index:       1,
			coinType:    bip44.CoinTypeSkycoin,
		},
		{
			name:        "index 2, test2, bitcoin",
			accountName: "test2",
			index:       2,
			coinType:    bip44.CoinTypeBitcoin,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			a := makeAccount(t, testSeed, testSeedPassphrase)
			ba, err := newBip44Account(a, tc.index, tc.accountName, tc.coinType)
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.accountName, ba.Name)
			require.Equal(t, tc.index, ba.Index)
			require.Equal(t, tc.coinType, ba.CoinType)
			require.Equal(t, 2, len(ba.Chains))

			externalKey, changeKey, err := makeChainPubKeys(a)
			require.NoError(t, err)
			require.Equal(t, *externalKey, ba.Chains[0].PubKey)
			require.Equal(t, *changeKey, ba.Chains[1].PubKey)
			require.Equal(t, 0, len(ba.Chains[0].Entries))
			require.Equal(t, 0, len(ba.Chains[1].Entries))
		})
	}
}

func TestBip44ChainNewAddress(t *testing.T) {

}

func TestBip44AccountNewAddress(t *testing.T) {

}
