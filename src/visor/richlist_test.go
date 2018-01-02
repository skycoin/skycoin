package visor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getLockedMap() map[string]struct{} {
	distributionAddresses := []string{"a1", "b1", "c1", "d1"}
	dmap := map[string]struct{}{}
	for _, addr := range distributionAddresses {
		dmap[addr] = struct{}{}
	}
	return dmap
}

func getAllAccounts() map[string]uint64 {
	accMap := map[string]uint64{}
	accMap["a1"] = 1000000
	accMap["b1"] = 1000000
	accMap["c1"] = 2123456
	accMap["d1"] = 1000000
	accMap["a2"] = 3010000
	accMap["b2"] = 2010000
	accMap["c2"] = 4010000
	accMap["d2"] = 1000000

	return accMap
}

func TestRichlist(t *testing.T) {
	expectedRichlist := Richlist{
		RichlistBalance{Address: "c2", Coins: "4.010000", Locked: false, coins: 4010000},
		RichlistBalance{Address: "a2", Coins: "3.010000", Locked: false, coins: 3010000},
		RichlistBalance{Address: "c1", Coins: "2.123456", Locked: true, coins: 2123456},
		RichlistBalance{Address: "b2", Coins: "2.010000", Locked: false, coins: 2010000},
		RichlistBalance{Address: "a1", Coins: "1.000000", Locked: true, coins: 1000000},
		RichlistBalance{Address: "b1", Coins: "1.000000", Locked: true, coins: 1000000},
		RichlistBalance{Address: "d1", Coins: "1.000000", Locked: true, coins: 1000000},
		RichlistBalance{Address: "d2", Coins: "1.000000", Locked: false, coins: 1000000},
	}

	accMap := getAllAccounts()
	distributionMap := getLockedMap()

	richlist, err := NewRichlist(map[string]uint64{}, map[string]struct{}{})
	assert.NoError(t, err)
	assert.Equal(t, Richlist{}, richlist)

	richlist, err = NewRichlist(accMap, distributionMap)
	assert.NoError(t, err)
	assert.Equal(t, expectedRichlist, richlist)
	cases := []struct {
		name        string
		filterMap   map[string]struct{}
		richerCount int
		result      Richlist
	}{
		{
			name:        "filterRichlist",
			filterMap:   map[string]struct{}{"a2": struct{}{}, "b2": struct{}{}, "d2": struct{}{}},
			richerCount: 5,
			result: Richlist{
				RichlistBalance{Address: "c2", Coins: "4.010000", Locked: false, coins: 4010000},
				RichlistBalance{Address: "c1", Coins: "2.123456", Locked: true, coins: 2123456},
				RichlistBalance{Address: "a1", Coins: "1.000000", Locked: true, coins: 1000000},
				RichlistBalance{Address: "b1", Coins: "1.000000", Locked: true, coins: 1000000},
				RichlistBalance{Address: "d1", Coins: "1.000000", Locked: true, coins: 1000000},
			},
		},

		{
			name:        "allRichlist",
			filterMap:   map[string]struct{}{},
			richerCount: 8,
			result:      expectedRichlist,
		},
		{
			name:        "lockedRichlist",
			filterMap:   map[string]struct{}{"c2": struct{}{}, "a2": struct{}{}, "b2": struct{}{}, "d2": struct{}{}},
			richerCount: 4,
			result: Richlist{
				RichlistBalance{Address: "c1", Coins: "2.123456", Locked: true, coins: 2123456},
				RichlistBalance{Address: "a1", Coins: "1.000000", Locked: true, coins: 1000000},
				RichlistBalance{Address: "b1", Coins: "1.000000", Locked: true, coins: 1000000},
				RichlistBalance{Address: "d1", Coins: "1.000000", Locked: true, coins: 1000000},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := richlist.FilterAddresses(tc.filterMap)
			assert.Equal(t, tc.richerCount, len(result), "%d != %d", tc.richerCount, len(result))
			assert.Equal(t, tc.result, result)
		})
	}
}
