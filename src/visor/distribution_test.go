package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
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

		isLocked := TransactionIsLocked(params.MainNetDistribution, uxArray)
		require.Equal(t, expectedIsLocked, isLocked)
	}

	for _, a := range params.MainNetDistribution.LockedAddresses() {
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
