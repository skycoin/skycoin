package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
)

func TestTransactionIsLocked(t *testing.T) {
	test := func(addrStr string, expectedIsLocked bool) {
		addr := cipher.MustDecodeBase58Address(addrStr)

		uxOut := coin.UxOut{
			Body: coin.UxBody{
				Address: addr,
			},
		}
		uxArray := coin.UxArray{uxOut}

		locked, err := TransactionIsLocked(params.MainNetDistribution, uxArray)
		require.NoError(t, err)
		require.Equal(t, expectedIsLocked, locked)
	}

	locked, err := params.MainNetDistribution.LockedAddresses()
	require.NoError(t, err)
	for _, a := range locked {
		test(a, true)
	}

	for _, a := range params.MainNetDistribution.UnlockedAddresses() {
		test(a, false)
	}

	// A random address should not be locked
	pubKey, _ := cipher.GenerateKeyPair()
	addr := cipher.AddressFromPubKey(pubKey)
	test(addr.String(), false)
}
