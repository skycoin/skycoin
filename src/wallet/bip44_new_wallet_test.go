package wallet

import (
	"testing"

	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/stretchr/testify/require"
)

const (
	testSeed           = "attitude coach wet rely typical habit alien security deny imitate spike slab"
	testSeedPassphrase = "pwd"
)

func TestBip44WalletNew(t *testing.T) {
	w := NewBip44WalletNew(Bip44WalletCreateOptions{
		Filename: "test.wlt",
		Label:    "test",
		Seed:     testSeed,
		Coin:     CoinTypeSkycoin,
	})

	ai, err := w.NewAccount("test")
	require.NoError(t, err)

	as, err := w.NewAddresses(ai, bip44.ExternalChainIndex, 2)
	require.NoError(t, err)
	t.Log(as[0])

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
