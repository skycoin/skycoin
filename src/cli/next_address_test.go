package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
)

func TestAddressChainKey(t *testing.T) {
	// bip39 test vector from https://github.com/trezor/python-mnemonic/blob/master/vectors.json
	const mnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	w, err := bip44wallet.NewWallet("bip44.wlt", "bip44", mnemonic, "",
		wallet.OptionCoinType(wallet.CoinTypeSkycoin), wallet.OptionGenerateN(5))
	require.NoError(t, err)
	entries, err := w.GetEntries(wallet.OptionExternal())
	require.NoError(t, err)
	want := entries.GetAddresses()
	require.Len(t, want, 5)

	for _, path := range []string{"0", "0/0"} {
		xpub, err := w.GetXPubKey(path)
		require.NoError(t, err)

		key, err := addressChainKey(xpub)
		require.NoError(t, err, path)

		var got []cipher.Addresser
		for i := uint32(0); i < 5; i++ {
			child, err := key.NewPublicChildKey(i)
			require.NoError(t, err)
			got = append(got, cipher.AddressFromPubKey(cipher.MustNewPubKey(child.Key)))
		}
		require.Equal(t, want, got, path)
	}

	_, err = addressChainKey("not an xpub")
	require.Error(t, err)
}
