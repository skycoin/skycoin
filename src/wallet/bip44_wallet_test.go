package wallet

import (
	"testing"

	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/stretchr/testify/require"
)

func TestBip44WalletAssign(t *testing.T) {
	w, err := NewBip44Wallet("test.wlt", Options{
		Seed:           "enact seek among recall one save armed parrot license ask giant fog",
		Coin:           meta.CoinTypeSkycoin,
		SeedPassphrase: "12345",
		CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
	}, nil)

	require.NoError(t, err)
	bw := w.(*Bip44Wallet)
	_, err = bw.NewExternalAddresses(defaultAccount, 4)
	require.NoError(t, err)

	require.Equal(t, 5, w.EntriesLen())

	_, err = bw.NewChangeAddresses(defaultAccount, 2)
	require.NoError(t, err)

	require.Equal(t, 7, bw.EntriesLen())

	w1, err := NewBip44Wallet("test1.wlt", Options{
		Seed:           "keep analyst jeans trip erosion race fantasy point spray dinner finger palm",
		Coin:           meta.CoinTypeSkycoin,
		SeedPassphrase: "54321",
		CryptoType:     crypto.CryptoTypeScryptChacha20poly1305Insecure,
	}, nil)

	require.NoError(t, err)

	bw1 := w1.(*Bip44Wallet)

	// Confirms there is one default address
	require.Equal(t, 1, bw1.EntriesLen())

	// Do assignment
	*bw1 = *bw

	// Confirms the entries length is correct
	require.Equal(t, 7, bw.EntriesLen())

	es, err := bw1.ExternalEntries(defaultAccount)
	require.NoError(t, err)
	require.Equal(t, 5, len(es))

	// Confirms that the seed is the same
	require.Equal(t, "enact seek among recall one save armed parrot license ask giant fog", bw1.Seed())
	// Confirms  that the seed passphrase is the same
	require.Equal(t, "12345", bw1.SeedPassphrase())
}

// TODO: generate a change address if there is no change entry
