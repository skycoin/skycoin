package cli

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	_ "github.com/skycoin/skycoin/src/wallet/bip44wallet"
	"github.com/skycoin/skycoin/src/wallet/collection"
	_ "github.com/skycoin/skycoin/src/wallet/deterministic"
)

// fakeOutputs serves one confirmed output of 10 coins at addr.
type fakeOutputs struct {
	t    *testing.T
	addr cipher.Address
}

func (f fakeOutputs) OutputsForAddresses(_ []string) (*readable.UnspentOutputsSummary, error) {
	head := coin.BlockHeader{BkSeq: 10, Time: 1700000000}
	ux := coin.UxOut{
		Head: coin.UxHead{BkSeq: 9, Time: head.Time - 3600},
		Body: coin.UxBody{SrcTransaction: testutil.RandSHA256(f.t), Address: f.addr, Coins: 10e6, Hours: 1000},
	}
	out, err := readable.NewUnspentOutput(visor.UnspentOutput{UxOut: ux, CalculatedHours: 1000})
	require.NoError(f.t, err)
	return &readable.UnspentOutputsSummary{
		Head:        readable.NewBlockHeader(head),
		HeadOutputs: readable.UnspentOutputs{out},
	}, nil
}

func TestCreateRawTxnNewChangeAddress(t *testing.T) {
	const seed = "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	to := []SendAmount{{Addr: testutil.MakeAddress().String(), Coins: 1e6}}

	cases := []struct {
		name     string
		typ      string
		password []byte
	}{
		{"deterministic", wallet.WalletTypeDeterministic, nil},
		{"encrypted deterministic", wallet.WalletTypeDeterministic, []byte("pwd")},
		{"bip44", wallet.WalletTypeBip44, nil},
		{"encrypted bip44", wallet.WalletTypeBip44, []byte("pwd")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			w, err := wallet.NewWallet("test.wlt", "label", seed, wallet.Options{
				Type:       tc.typ,
				GenerateN:  1,
				Encrypt:    len(tc.password) != 0,
				Password:   tc.password,
				CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure,
			})
			require.NoError(t, err)
			require.NoError(t, wallet.Save(w, dir))
			first, err := w.GetEntryAt(0)
			require.NoError(t, err)

			file := filepath.Join(dir, "test.wlt")
			txn, err := CreateRawTxnFromWallet(fakeOutputs{t, first.Address.(cipher.Address)}, file, "", to, NewPasswordReader(tc.password), params.MainNetDistribution)
			require.NoError(t, err)
			require.Len(t, txn.Out, 2)
			change := txn.Out[1].Address
			require.NotEqual(t, first.Address.(cipher.Address), change)

			saved, err := wallet.Load(file)
			require.NoError(t, err)
			has, err := saved.HasEntry(change)
			require.NoError(t, err)
			require.True(t, has, "the new change address is saved in the wallet file")
		})
	}

	t.Run("collection keeps an input address", func(t *testing.T) {
		dir := t.TempDir()
		w, err := wallet.NewWallet("test.wlt", "label", "", wallet.Options{Type: wallet.WalletTypeCollection})
		require.NoError(t, err)
		pk, sk := cipher.GenerateKeyPair()
		addr := cipher.AddressFromPubKey(pk)
		require.NoError(t, w.(*collection.Wallet).AddEntry(wallet.Entry{Address: addr, Public: pk, Secret: sk}))
		require.NoError(t, wallet.Save(w, dir))

		txn, err := CreateRawTxnFromWallet(fakeOutputs{t, addr}, filepath.Join(dir, "test.wlt"), "", to, nil, params.MainNetDistribution)
		require.NoError(t, err)
		require.Len(t, txn.Out, 2)
		require.Equal(t, addr, txn.Out[1].Address)
	})
}

func TestSweepWallet(t *testing.T) {
	const seed = "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	dir := t.TempDir()
	w, err := wallet.NewWallet("test.wlt", "label", seed, wallet.Options{Type: wallet.WalletTypeDeterministic, GenerateN: 1})
	require.NoError(t, err)
	require.NoError(t, wallet.Save(w, dir))
	first, err := w.GetEntryAt(0)
	require.NoError(t, err)
	file := filepath.Join(dir, "test.wlt")

	txn, dest, err := sweepWallet(fakeOutputs{t, first.Address.(cipher.Address)}, file, "", nil, params.MainNetDistribution)
	require.NoError(t, err)
	require.Len(t, txn.Out, 1, "a sweep has no change output")
	require.Equal(t, dest, txn.Out[0].Address.String())
	require.Equal(t, uint64(10e6), txn.Out[0].Coins)
	require.NotEqual(t, first.Address.(cipher.Address).String(), dest)

	saved, err := wallet.Load(file)
	require.NoError(t, err)
	has, err := saved.HasEntry(txn.Out[0].Address)
	require.NoError(t, err)
	require.True(t, has)

	other := testutil.MakeAddress()
	txn, dest, err = sweepWallet(fakeOutputs{t, first.Address.(cipher.Address)}, file, other.String(), nil, params.MainNetDistribution)
	require.NoError(t, err)
	require.Equal(t, other.String(), dest)
	require.Equal(t, other, txn.Out[0].Address)
}
