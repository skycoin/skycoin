package params

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDistributionAddressArrays(t *testing.T) {
	require.Len(t, MainNetDistribution.Addresses, 100)

	// At the time of this writing, there should be 25 addresses in the
	// unlocked pool and 75 in the locked pool.
	require.Len(t, MainNetDistribution.UnlockedAddresses(), 25)
	require.Len(t, MainNetDistribution.LockedAddresses(), 75)

	addrsMap := make(map[string]struct{})
	for _, a := range MainNetDistribution.Addresses {
		// Check no duplicate address in distribution addresses
		_, ok := addrsMap[a]
		require.False(t, ok)
		addrsMap[a] = struct{}{}
	}

	unlocked := MainNetDistribution.UnlockedAddresses()
	unlockedMap := make(map[string]struct{})
	for _, a := range unlocked {
		// Check no duplicate address in unlocked addresses
		_, ok := unlockedMap[a]
		require.False(t, ok)

		// Check unlocked address in set of all addresses
		_, ok = addrsMap[a]
		require.True(t, ok)

		unlockedMap[a] = struct{}{}
	}

	locked := MainNetDistribution.LockedAddresses()
	lockedMap := make(map[string]struct{})
	for _, a := range locked {
		// Check no duplicate address in locked addresses
		_, ok := lockedMap[a]
		require.False(t, ok)

		// Check locked address in set of all addresses
		_, ok = addrsMap[a]
		require.True(t, ok)

		// Check locked address not in unlocked addresses
		_, ok = unlockedMap[a]
		require.False(t, ok)

		lockedMap[a] = struct{}{}
	}
}
