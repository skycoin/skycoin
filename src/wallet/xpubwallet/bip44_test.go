package xpubwallet

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/bip44wallet"
)

// bip39 test vector from https://github.com/trezor/python-mnemonic/blob/master/vectors.json
const testMnemonic = "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

// bip44Fixture returns the account and external chain xpubs of a bip44 wallet
// built from testMnemonic, and the wallet's first 6 external addresses.
func bip44Fixture(t *testing.T) (accountXPub, chainXPub string, want []cipher.Addresser) {
	w, err := bip44wallet.NewWallet("bip44.wlt", "bip44", testMnemonic, "",
		wallet.OptionCoinType(wallet.CoinTypeSkycoin), wallet.OptionGenerateN(6))
	require.NoError(t, err)

	entries, err := w.GetEntries(wallet.OptionExternal())
	require.NoError(t, err)
	want = entries.GetAddresses()
	require.Len(t, want, 6)

	accountXPub, err = w.GetXPubKey("0")
	require.NoError(t, err)
	chainXPub, err = w.GetXPubKey("0/0")
	require.NoError(t, err)
	return accountXPub, chainXPub, want
}

func TestAccountXPubUsesExternalChain(t *testing.T) {
	accountXPub, chainXPub, want := bip44Fixture(t)

	for name, xp := range map[string]string{"account": accountXPub, "chain": chainXPub} {
		t.Run(name, func(t *testing.T) {
			w, err := NewWallet("xpub.wlt", "xpub", xp, wallet.OptionGenerateN(5))
			require.NoError(t, err)
			require.Equal(t, xp, w.XPub())

			got, err := w.GetAddresses()
			require.NoError(t, err)
			require.Equal(t, want[:5], got)

			// A save and load keeps the chain, and new addresses continue on it.
			b, err := w.Serialize()
			require.NoError(t, err)
			w2, err := Loader{}.Load(b)
			require.NoError(t, err)
			_, err = w2.GenerateAddresses(wallet.OptionGenerateN(1))
			require.NoError(t, err)

			got, err = w2.GetAddresses()
			require.NoError(t, err)
			require.Equal(t, want, got)
		})
	}

	w, err := NewWallet("xpub.wlt", "xpub", accountXPub)
	require.NoError(t, err)
	require.Equal(t, "0", w.Meta[wallet.MetaXPubChain])

	w, err = NewWallet("xpub.wlt", "xpub", chainXPub)
	require.NoError(t, err)
	_, ok := w.Meta[wallet.MetaXPubChain]
	require.False(t, ok)
}

func TestLegacyAccountXPubWalletKeepsItsChain(t *testing.T) {
	accountXPub, _, _ := bip44Fixture(t)

	key, err := bip32.DeserializeEncodedPublicKey(accountXPub)
	require.NoError(t, err)
	legacy := make(wallet.Entries, 3)
	for i := range legacy {
		child, err := key.NewPublicChildKey(uint32(i))
		require.NoError(t, err)
		pk := cipher.MustNewPubKey(child.Key)
		legacy[i] = wallet.Entry{Address: cipher.AddressFromPubKey(pk), Public: pk, ChildNumber: uint32(i)}
	}

	// A wallet saved before the xpubChain field existed.
	old := &Wallet{
		Meta: wallet.Meta{
			wallet.MetaFilename: "old.wlt",
			wallet.MetaType:     WalletType,
			wallet.MetaVersion:  wallet.Version,
			wallet.MetaCoin:     string(wallet.CoinTypeSkycoin),
			wallet.MetaXPub:     accountXPub,
		},
		entries: legacy[:2],
		xpub:    key,
	}
	b, err := old.Serialize()
	require.NoError(t, err)

	w, err := Loader{}.Load(b)
	require.NoError(t, err)
	_, err = w.GenerateAddresses(wallet.OptionGenerateN(1))
	require.NoError(t, err)
	got, err := w.GetAddresses()
	require.NoError(t, err)
	require.Equal(t, legacy.GetAddresses(), got)

	// Marking the same file as external chain must not silently mix chains.
	old.Meta[wallet.MetaXPubChain] = "0"
	b, err = old.Serialize()
	require.NoError(t, err)
	_, err = Loader{}.Load(b)
	require.EqualError(t, err, "stored entries were not derived from the wallet's xpub chain")

	old.Meta[wallet.MetaXPubChain] = "1"
	b, err = old.Serialize()
	require.NoError(t, err)
	_, err = Loader{}.Load(b)
	require.EqualError(t, err, `invalid xpubChain value "1"`)
}
