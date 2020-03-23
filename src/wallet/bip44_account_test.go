package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBip44Account(t *testing.T) {
	tt := []struct {
		name        string
		accountName string
		index       uint32
		coinType    CoinType
		err         error
	}{
		{
			name:        "index 0, test, skycoin",
			accountName: "test",
			index:       0,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:        "index 1, test1, skycoin",
			accountName: "test1",
			index:       1,
			coinType:    CoinTypeSkycoin,
		},
		{
			name:        "index 2, test2, bitcoin",
			accountName: "test2",
			index:       2,
			coinType:    CoinTypeBitcoin,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ba, err := newBip44Account(bip44AccountCreateOptions{
				name:     tc.accountName,
				index:    tc.index,
				seed:     testSeed,
				coinType: tc.coinType,
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
