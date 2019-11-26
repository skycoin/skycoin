package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
)

func getLockedMap(distributionAddresses [4]cipher.Address) map[cipher.Address]struct{} {
	// distributionAddresses := []cipher.Address{"a1", "b1", "c1", "d1"}
	dmap := map[cipher.Address]struct{}{}
	for _, addr := range distributionAddresses {
		dmap[addr] = struct{}{}
	}
	return dmap
}

func getAllAccounts(distributionAddresses, otherAddresses [4]cipher.Address) map[cipher.Address]uint64 {
	accMap := map[cipher.Address]uint64{}
	accMap[distributionAddresses[0]] = 1000000
	accMap[distributionAddresses[1]] = 1000000
	accMap[distributionAddresses[2]] = 2123456
	accMap[distributionAddresses[3]] = 1000000
	accMap[otherAddresses[0]] = 3010000
	accMap[otherAddresses[1]] = 2010000
	accMap[otherAddresses[2]] = 4010000
	accMap[otherAddresses[3]] = 1000000

	return accMap
}

func TestRichlist(t *testing.T) {
	otherAddresses := [4]cipher.Address{
		cipher.MustDecodeBase58Address("2cmpPv9PJfKFStekrKZXBnAfLKE6cB7qMrS"),
		cipher.MustDecodeBase58Address("jhLw4EXNn2E7zVjrmi8fGsATZfRnAXfqRj"),
		cipher.MustDecodeBase58Address("R7zjFhmW3KqGz6r92VFpJTpRWCzaXSokYb"),
		cipher.MustDecodeBase58Address("JFQRvKXBoTt6D8aiFVrGquemzbrQDGKTAR"),
	}

	distributionAddresses := [4]cipher.Address{
		cipher.MustDecodeBase58Address("DniB7KqDRNx8CjM6vruaKwbQPgWj1GSj5t"),
		cipher.MustDecodeBase58Address("FbJuRez3RKpYsTSYTVyAQt146vzcFNkqpU"),
		cipher.MustDecodeBase58Address("2mxNdCnUd7vF1uSpRhSDMEhZHKAyL9r1Uys"),
		cipher.MustDecodeBase58Address("uBcaMg2vGpy45K7NVsRGmuNXQdaB8kgHfM"),
	}

	expectedRichlist := Richlist{
		RichlistBalance{Address: otherAddresses[2], Coins: 4010000, Locked: false},
		RichlistBalance{Address: otherAddresses[0], Coins: 3010000, Locked: false},
		RichlistBalance{Address: distributionAddresses[2], Coins: 2123456, Locked: true},
		RichlistBalance{Address: otherAddresses[1], Coins: 2010000, Locked: false},
		RichlistBalance{Address: distributionAddresses[0], Coins: 1000000, Locked: true},
		RichlistBalance{Address: distributionAddresses[1], Coins: 1000000, Locked: true},
		RichlistBalance{Address: distributionAddresses[3], Coins: 1000000, Locked: true},
		RichlistBalance{Address: otherAddresses[3], Coins: 1000000, Locked: false},
	}

	accMap := getAllAccounts(distributionAddresses, otherAddresses)
	distributionMap := getLockedMap(distributionAddresses)

	richlist, err := NewRichlist(map[cipher.Address]uint64{}, map[cipher.Address]struct{}{})
	require.NoError(t, err)
	require.Equal(t, Richlist{}, richlist)

	richlist, err = NewRichlist(accMap, distributionMap)
	require.NoError(t, err)
	require.Equal(t, expectedRichlist, richlist)

	cases := []struct {
		name        string
		filterMap   map[cipher.Address]struct{}
		richlistLen int
		result      Richlist
	}{
		{
			name: "filterRichlist",
			filterMap: map[cipher.Address]struct{}{
				otherAddresses[0]: struct{}{},
				otherAddresses[1]: struct{}{},
				otherAddresses[3]: struct{}{},
			},
			richlistLen: 5,
			result: Richlist{
				RichlistBalance{Address: otherAddresses[2], Locked: false, Coins: 4010000},
				RichlistBalance{Address: distributionAddresses[2], Locked: true, Coins: 2123456},
				RichlistBalance{Address: distributionAddresses[0], Locked: true, Coins: 1000000},
				RichlistBalance{Address: distributionAddresses[1], Locked: true, Coins: 1000000},
				RichlistBalance{Address: distributionAddresses[3], Locked: true, Coins: 1000000},
			},
		},

		{
			name:        "allRichlist",
			filterMap:   map[cipher.Address]struct{}{},
			richlistLen: 8,
			result:      expectedRichlist,
		},
		{
			name: "lockedRichlist",
			filterMap: map[cipher.Address]struct{}{
				otherAddresses[0]: struct{}{},
				otherAddresses[1]: struct{}{},
				otherAddresses[2]: struct{}{},
				otherAddresses[3]: struct{}{},
			},
			result: Richlist{
				RichlistBalance{Address: distributionAddresses[2], Locked: true, Coins: 2123456},
				RichlistBalance{Address: distributionAddresses[0], Locked: true, Coins: 1000000},
				RichlistBalance{Address: distributionAddresses[1], Locked: true, Coins: 1000000},
				RichlistBalance{Address: distributionAddresses[3], Locked: true, Coins: 1000000},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := richlist.FilterAddresses(tc.filterMap)
			require.Equal(t, len(tc.result), len(result), "%d != %d", len(tc.result), len(result))

			require.Equal(t, tc.result, result)
		})
	}
}
