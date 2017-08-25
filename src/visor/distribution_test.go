package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

func TestDistributionAddressArrays(t *testing.T) {
	require.Len(t, GetDistributionAddresses(), 100)

	// At the time of this writing, there should be 25 addresses in the
	// unlocked pool and 75 in the locked pool.
	require.Len(t, GetUnlockedDistributionAddresses(), 25)
	require.Len(t, GetLockedDistributionAddresses(), 75)

	all := GetDistributionAddresses()
	allMap := make(map[string]struct{})
	for _, a := range all {
		// Check no duplicate address in distribution addresses
		_, ok := allMap[a]
		require.False(t, ok)
		allMap[a] = struct{}{}
	}

	unlocked := GetUnlockedDistributionAddresses()
	unlockedMap := make(map[string]struct{})
	for _, a := range unlocked {
		// Check no duplicate address in unlocked addresses
		_, ok := unlockedMap[a]
		require.False(t, ok)

		// Check unlocked address in set of all addresses
		_, ok = allMap[a]
		require.True(t, ok)

		unlockedMap[a] = struct{}{}
	}

	locked := GetLockedDistributionAddresses()
	lockedMap := make(map[string]struct{})
	for _, a := range locked {
		// Check no duplicate address in locked addresses
		_, ok := lockedMap[a]
		require.False(t, ok)

		// Check locked address in set of all addresses
		_, ok = allMap[a]
		require.True(t, ok)

		// Check locked address not in unlocked addresses
		_, ok = unlockedMap[a]
		require.False(t, ok)

		lockedMap[a] = struct{}{}
	}
}

func TestTransactionIsLocked(t *testing.T) {
	test := func(addrStr string, expectedIsLocked bool) {
		addr := cipher.MustDecodeBase58Address(addrStr)

		uxOut := coin.UxOut{
			Body: coin.UxBody{
				Address: addr,
			},
		}
		uxArray := coin.UxArray{uxOut}

		isLocked := TransactionIsLocked(uxArray)
		require.Equal(t, expectedIsLocked, isLocked)
	}

	for _, a := range GetLockedDistributionAddresses() {
		test(a, true)
	}

	for _, a := range GetUnlockedDistributionAddresses() {
		test(a, false)
	}

	// A random address should not be locked
	pubKey, _ := cipher.GenerateKeyPair()
	addr := cipher.AddressFromPubKey(pubKey)
	test(addr.String(), false)
}
