package wallet

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	invalidBip44Seed = "invalid bip44 seed"
)

func TestNewBip44Account(t *testing.T) {
	tt := []struct {
		name           string
		accountName    string
		index          int
		seed           string
		seedPassphrase string
		coinType       CoinType
		err            error
	}{
		{
			name:           "index 0, skycoin, with seed passphrase",
			accountName:    "test",
			index:          0,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
		},
		{
			name:        "index 0, skycoin, without seed passphrase",
			accountName: "test",
			index:       0,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:           "index 1, skycoin, with seed passphrase",
			accountName:    "test",
			index:          1,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeSkycoin,
		},
		{
			name:        "index 1, skycoin, without seed passphrase",
			accountName: "test",
			index:       1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:           "index 2, bitcoin, with seed passphrase",
			accountName:    "test",
			index:          2,
			seed:           testSeed,
			seedPassphrase: testSeedPassphrase,
			coinType:       CoinTypeBitcoin,
		},
		{
			name:        "index 2, bitcoin, without seed passphrase",
			accountName: "test",
			index:       2,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
		},
		{
			name:        "index 0x80000000 -1, skycoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000 - 1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:        "index 0x80000000 -1, bitcoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000 - 1,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
		},
		{
			name:        "index 0x80000000, skycoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "index 0x80000000, bitcoin, without seed passphrase",
			accountName: "test",
			index:       0x80000000,
			seed:        testSeed,
			coinType:    CoinTypeBitcoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "index uint32(-1), skycoin, without seed passphrase",
			accountName: "test",
			index:       -1,
			seed:        testSeed,
			coinType:    CoinTypeSkycoin,
			err:         errors.New("bip44 account index must be < 0x80000000"),
		},
		{
			name:        "skycoin, invalid bip44 seed",
			accountName: "test",
			index:       0,
			seed:        "abc def ghe a",
			coinType:    CoinTypeSkycoin,
			err:         errors.New("Mnemonic must have 12, 15, 18, 21 or 24 words"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ba, err := newBip44Account(bip44AccountCreateOptions{
				name:           tc.accountName,
				index:          uint32(tc.index),
				seed:           tc.seed,
				seedPassphrase: testSeedPassphrase,
				coinType:       tc.coinType,
			})
			if err != nil {
				require.Equal(t, tc.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.accountName, ba.Name)
			require.Equal(t, uint32(tc.index), ba.Index)
			require.Equal(t, tc.coinType, ba.CoinType)
			require.Equal(t, 2, len(ba.Chains))

			externalKey, changeKey, err := makeChainPubKeys(&ba.Account)
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
