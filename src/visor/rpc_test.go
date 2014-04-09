package visor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
)

func TestGetWalletBalance(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	v := setupMasterVisor()
	b := wallet.Balance{v.Config.GenesisCoinVolume, v.Config.GenesisCoinVolume}
	spend := wallet.Balance{1e6, 1e6}
	id := v.Wallets[0].GetID()

	// Make a pending transfer of coins, so we can see predicted is correct
	v2, _ := setupVisor()
	addr := v2.Wallets[0].GetAddresses()[0]
	tx, err := v.Spend(id, spend, 0, addr)
	assert.Nil(t, err)
	v.RecordTxn(tx)
	assert.Equal(t, len(v.Unconfirmed.Txns), 1)

	assert.Equal(t, rpc.GetWalletBalance(v, id).Confirmed, b)
	// The predicted balance is 0, because we spent everything contained
	// in the genesis block
	assert.Equal(t, rpc.GetWalletBalance(v, id).Predicted, wallet.Balance{})

	assert.Nil(t, rpc.GetWalletBalance(nil, id))
}

func TestReloadWallets(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	// Create 2 wallets
	v, _ := setupVisor()
	v.CreateWallet()
	assert.Equal(t, len(v.Wallets), 2)
	v.SaveWallets()
	found, err := filepath.Glob(testWalletDir + "*.wlt")
	assert.Nil(t, err)
	assert.Equal(t, len(found), 2)
	wallets := v.Wallets

	// Create a new visor, which will load those 2 wallets
	v, _ = setupVisor()
	assert.Equal(t, len(v.Wallets), 2)
	for _, w := range v.Wallets {
		assert.True(t, w.GetID() == wallets[0].GetID() ||
			w.GetID() == wallets[1].GetID())
	}

	// Create another wallet through a separate channel
	v2, _ := setupVisor()
	addw := v2.CreateWallet()
	v2.SaveWallets()

	// Delete one of the original wallets
	os.Remove(wallets[0].GetFilename())

	// Check that the changed set is in use
	assert.Nil(t, rpc.ReloadWallets(v))
	assert.Equal(t, len(v.Wallets), 2)
	for _, w := range v.Wallets {
		assert.True(t, w.GetID() == wallets[1].GetID() ||
			w.GetID() == addw.GetID())
	}

	assert.Nil(t, rpc.ReloadWallets(nil))
}

func TestSaveWallets(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	v, _ := setupVisor()
	assert.Equal(t, len(v.Wallets), 1)
	v.CreateWallet()
	assert.Equal(t, len(v.Wallets), 2)
	errs := rpc.SaveWallets(v)
	assert.Nil(t, errs)
	assert.Equal(t, len(errs), 0)
	wallets, err := filepath.Glob(testWalletDir + "*.wlt")
	assert.Nil(t, err)
	assert.Equal(t, len(wallets), 2)
	w, _ := setupVisor()
	for i, _ := range v.Wallets {
		vw := v.Wallets[i]
		found := false
		for _, ww := range w.Wallets {
			if vw.GetID() != ww.GetID() {
				continue
			}
			assert.Equal(t, vw.GetID(), ww.GetID())
			assert.Equal(t, vw.GetEntries(), ww.GetEntries())
			assert.Equal(t, vw.GetFilename(), ww.GetFilename())
			assert.Equal(t, vw.GetName(), ww.GetName())
			found = true
		}
		assert.True(t, found)
	}

	assert.Nil(t, rpc.CreateWallet(nil))
}

func TestCreateWallet(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	v, _ := setupVisor()
	assert.Equal(t, len(v.Wallets), 1)
	w := rpc.CreateWallet(v)
	assert.NotNil(t, w)
	assert.Equal(t, len(v.Wallets), 2)

	assert.Nil(t, rpc.CreateWallet(nil))
}

func TestGetWallet(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	v, _ := setupVisor()
	w := v.Wallets[0]
	assert.Equal(t, wallet.NewReadableWallet(w), rpc.GetWallet(v, w.GetID()))
	w2, err := rpc.GetWallet(v, w.GetID()).ToWallet()
	assert.Nil(t, err)
	assert.Equal(t, w, w2)

	assert.Nil(t, rpc.GetWallet(nil, w.GetID()))
}

func TestGetWallets(t *testing.T) {
	defer cleanupVisor()
	rpc := RPC{}
	v, _ := setupVisor()
	v.CreateWallet()
	v.CreateWallet()
	assert.Equal(t, len(v.Wallets), 3)
	wallets := rpc.GetWallets(v)
	assert.Equal(t, len(wallets), 3)
	expect := make([]*wallet.ReadableWallet, 3)
	for i, _ := range v.Wallets {
		expect[i] = wallet.NewReadableWallet(v.Wallets[i])
	}
	for i, w := range expect {
		w2 := wallets[i]
		assert.Equal(t, w2.ID, w.ID)
		assert.Equal(t, w2.Type, w.Type)
		assert.Equal(t, w2.Name, w.Name)
		assert.Equal(t, w2.Filename, w.Filename)
		assert.Equal(t, w2.Extra, w.Extra)
		assert.Equal(t, len(w2.Entries), len(w.Entries))
		for j, e := range w.Entries {
			e2 := w2.Entries[j]
			assert.Equal(t, e2.Public, e.Public)
			assert.Equal(t, e2.Address, e.Address)
			// Secret key should not be exposed by GetWallets
			assert.Equal(t, e2.Secret, "")
			assert.NotEqual(t, e.Secret, "")
		}
	}

	assert.Nil(t, rpc.GetWallets(nil))
}
