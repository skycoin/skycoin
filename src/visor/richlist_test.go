package visor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getDistributionMap() map[string]struct{} {
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
	accMap["d2"] = 1010000

	return accMap
}

func TestAccountSort(t *testing.T) {
	accMap := getAllAccounts()
	distributionMap := getDistributionMap()
	accMgr := NewAccountMgr(accMap, distributionMap)
	accMgr.Sort()
	cases := []struct {
		topn                int
		includeDistribution bool
		topnCount           int
		result              []AccountJSON
		err                 error
	}{
		{
			topn:                5,
			includeDistribution: true,
			topnCount:           5,
			result: []AccountJSON{AccountJSON{Addr: "c1", Coins: "2.123456", Locked: true},
				AccountJSON{Addr: "d1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "b1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "a1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "c2", Coins: "4.010000", Locked: false},
			},
			err: nil,
		},

		{
			topn:                4,
			includeDistribution: false,
			topnCount:           4,
			result: []AccountJSON{AccountJSON{Addr: "c2", Coins: "4.010000", Locked: false},
				AccountJSON{Addr: "a2", Coins: "3.010000", Locked: false},
				AccountJSON{Addr: "b2", Coins: "2.010000", Locked: false},
				AccountJSON{Addr: "d2", Coins: "1.010000", Locked: false},
			},
			err: nil,
		},
		{
			topn:                -1,
			includeDistribution: false,
			topnCount:           4,
			result: []AccountJSON{AccountJSON{Addr: "c2", Coins: "4.010000", Locked: false},
				AccountJSON{Addr: "a2", Coins: "3.010000", Locked: false},
				AccountJSON{Addr: "b2", Coins: "2.010000", Locked: false},
				AccountJSON{Addr: "d2", Coins: "1.010000", Locked: false},
			},
			err: nil,
		},
		{
			topn:                -1,
			includeDistribution: true,
			topnCount:           8,
			result: []AccountJSON{AccountJSON{Addr: "c1", Coins: "2.123456", Locked: true},
				AccountJSON{Addr: "d1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "b1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "a1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "c2", Coins: "4.010000", Locked: false},
				AccountJSON{Addr: "a2", Coins: "3.010000", Locked: false},
				AccountJSON{Addr: "b2", Coins: "2.010000", Locked: false},
				AccountJSON{Addr: "d2", Coins: "1.010000", Locked: false},
			},
			err: nil,
		},
		{
			topn:                0,
			includeDistribution: true,
			topnCount:           8,
			result: []AccountJSON{AccountJSON{Addr: "c1", Coins: "2.123456", Locked: true},
				AccountJSON{Addr: "d1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "b1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "a1", Coins: "1.000000", Locked: true},
				AccountJSON{Addr: "c2", Coins: "4.010000", Locked: false},
				AccountJSON{Addr: "a2", Coins: "3.010000", Locked: false},
				AccountJSON{Addr: "b2", Coins: "2.010000", Locked: false},
				AccountJSON{Addr: "d2", Coins: "1.010000", Locked: false},
			},
			err: nil,
		},
	}
	for _, tc := range cases {
		name := fmt.Sprintf("topn=%d include=%v", tc.topn, tc.includeDistribution)
		t.Run(name, func(t *testing.T) {
			result, err := accMgr.GetTopn(tc.topn, tc.includeDistribution)
			if tc.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.topnCount, len(result), "%d != %d", tc.topnCount, len(result))
				assert.Equal(t, tc.result, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err)
			}
		})
	}
}
