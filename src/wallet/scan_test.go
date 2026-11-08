package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

type activeSet map[cipher.Addresser]bool

func (s activeSet) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	active := make([]bool, len(addrs))
	for i, a := range addrs {
		active[i] = s[a]
	}
	return active, nil
}

func TestScanGap(t *testing.T) {
	chain := make([]cipher.Addresser, 60)
	for i := range chain {
		pk, _ := cipher.GenerateKeyPair()
		chain[i] = cipher.AddressFromPubKey(pk)
	}

	tt := []struct {
		name   string
		gap    uint64
		active []int
		keep   int
	}{
		{name: "none used", gap: 20, keep: 0},
		{name: "first used", gap: 20, active: []int{0}, keep: 1},
		{name: "last of the first window", gap: 20, active: []int{19}, keep: 20},
		{name: "past the first window", gap: 20, active: []int{19, 38}, keep: 39},
		{name: "gap of 19 is bridged", gap: 20, active: []int{5, 25}, keep: 26},
		{name: "gap of 20 ends the scan", gap: 20, active: []int{5, 26}, keep: 6},
		{name: "small gap", gap: 2, active: []int{1, 3, 5}, keep: 6},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			set := activeSet{}
			for _, i := range tc.active {
				set[chain[i]] = true
			}
			next := 0
			generate := func(n uint64) ([]cipher.Addresser, error) {
				end := min(next+int(n), len(chain)) //nolint:gosec
				addrs := chain[next:end]
				next = end
				return addrs, nil
			}

			addrs, err := ScanGap(tc.gap, generate, set)
			require.NoError(t, err)
			require.Equal(t, chain[:tc.keep], addrs)
		})
	}
}
