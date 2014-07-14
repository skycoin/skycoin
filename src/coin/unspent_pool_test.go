package coin

import (
	"sort"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

func TestNewUnspentPool(t *testing.T) {
	up := NewUnspentPool()
	assert.Equal(t, len(up.Pool), 0)
	assert.Equal(t, up.XorHash, cipher.SHA256{})
}

func TestUnspentPoolRebuild(t *testing.T) {
	up := NewUnspentPool()
	arr := make(UxArray, 0)
	arr = append(arr, makeUxOut(t))
	arr = append(arr, makeUxOut(t))
	assert.Equal(t, len(up.Pool), 0)
	assert.Equal(t, up.XorHash, cipher.SHA256{})
	up.Rebuild(arr)
	assert.Equal(t, len(up.Pool), 2)
	for _, x := range arr {
		ux, ok := up.Pool[x.Hash()]
		assert.True(t, ok)
		assert.Equal(t, x, ux)
	}
	h := cipher.SHA256{}
	h = h.Xor(arr[0].SnapshotHash())
	h = h.Xor(arr[1].SnapshotHash())
	assert.Equal(t, up.XorHash, h)
	assert.NotEqual(t, up.XorHash, cipher.SHA256{})

	// Duplicate item in array causes panic
	arr = append(arr, arr[0])
	assert.Panics(t, func() { up.Rebuild(arr) })
}

func TestUnspentPoolAdd(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	assert.Equal(t, len(up.Pool), 0)
	up.Add(ux)
	assert.Equal(t, len(up.Pool), 1)
	ux, ok := up.Pool[ux.Hash()]
	assert.True(t, ok)
	assert.NotEqual(t, up.XorHash, cipher.SHA256{})
	assert.Equal(t, up.XorHash, ux.SnapshotHash())
	// Duplicate add panics
	h := up.XorHash
	assert.Panics(t, func() { up.Add(ux) })
	assert.Equal(t, len(up.Pool), 1)
	assert.Equal(t, up.XorHash, h)
	// Add a 2nd is ok
	ux2 := makeUxOut(t)
	up.Add(ux2)
	_, ok = up.Pool[ux2.Hash()]
	assert.True(t, ok)
	assert.Equal(t, len(up.Pool), 2)
	h = ux.SnapshotHash()
	h = h.Xor(ux2.SnapshotHash())
	assert.Equal(t, up.XorHash, h)
}

func TestUnspentPoolGet(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	_, ok := up.Get(ux.Hash())
	assert.False(t, ok)
	up.Add(ux)
	ux2, ok := up.Get(ux.Hash())
	assert.True(t, ok)
	assert.Equal(t, ux, ux2)
}

func TestUnspentPoolHas(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	assert.False(t, up.Has(ux.Hash()))
	up.Add(ux)
	assert.True(t, up.Has(ux.Hash()))
}

func TestUnspentPoolDel(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	up.Add(ux)
	up.Add(ux2)
	up.Add(ux3)
	assert.Equal(t, len(up.Pool), 3)
	// Unknown hash
	up.Del(SHA256{})
	assert.Equal(t, len(up.Pool), 3)
	// Delete middle one
	up.Del(ux2.Hash())
	assert.Equal(t, len(up.Pool), 2)
	_, ok := up.Pool[ux.Hash()]
	assert.True(t, ok)
	_, ok = up.Pool[ux3.Hash()]
	assert.True(t, ok)
	h := ux.SnapshotHash()
	h = h.Xor(ux3.SnapshotHash())
	assert.Equal(t, up.XorHash, h)
	// Delete first one
	up.Del(ux.Hash())
	assert.Equal(t, len(up.Pool), 1)
	_, ok = up.Pool[ux3.Hash()]
	assert.True(t, ok)
	assert.Equal(t, up.XorHash, ux3.SnapshotHash())
	// Delete remaining one
	up.Del(ux3.Hash())
	assert.Equal(t, len(up.Pool), 0)
	assert.Equal(t, up.XorHash, cipher.SHA256{})
}

func TestUnspentPoolDelMultiple(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	ux4 := makeUxOut(t)
	ux5 := makeUxOut(t)
	up.Add(ux)
	up.Add(ux2)
	up.Add(ux3)
	up.Add(ux4)
	assert.Equal(t, len(up.Pool), 4)
	// Delete 1st and 3rd and an unknown
	up.DelMultiple([]SHA256{ux.Hash(), ux3.Hash(), ux5.Hash()})
	assert.Equal(t, len(up.Pool), 2)
	_, ok := up.Pool[ux2.Hash()]
	assert.True(t, ok)
	_, ok = up.Pool[ux4.Hash()]
	assert.True(t, ok)
	h := ux2.SnapshotHash()
	h = h.Xor(ux4.SnapshotHash())
	assert.Equal(t, up.XorHash, h)
}

func TestUnspentPoolAllForAddress(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	ux3.Body.Address = ux.Body.Address
	up.Add(ux)
	up.Add(ux2)
	up.Add(ux3)
	// 2 unspents for address
	uxs := up.AllForAddress(ux.Body.Address)
	assert.Equal(t, len(uxs), 2)
	assert.False(t, uxs.HasDupes())
	assert.True(t, uxs[0] == ux || uxs[1] == ux)
	assert.True(t, uxs[0] == ux3 || uxs[1] == ux3)
	// 1 unspent
	uxs = up.AllForAddress(ux2.Body.Address)
	assert.Equal(t, len(uxs), 1)
	assert.Equal(t, uxs[0], ux2)
	// No known addresses
	uxs = up.AllForAddress(Address{})
	assert.Equal(t, len(uxs), 0)
}

func TestUnspentPoolAllForAddresses(t *testing.T) {
	up := NewUnspentPool()
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	ux4 := makeUxOut(t)
	ux3.Body.Address = ux.Body.Address
	up.Add(ux)
	up.Add(ux2)
	up.Add(ux3)
	up.Add(ux4)

	// No addresses
	uxs := up.AllForAddresses([]Address{})
	assert.Equal(t, len(uxs), 0)
	// 1 address
	uxs = up.AllForAddresses([]Address{ux4.Body.Address})
	assert.Equal(t, len(uxs), 1)
	assert.Equal(t, len(uxs[ux4.Body.Address]), 1)
	assert.Equal(t, uxs[ux4.Body.Address][0], ux4)
	// 2 addresses
	uxs = up.AllForAddresses([]Address{ux.Body.Address, ux2.Body.Address})
	assert.Equal(t, len(uxs), 2)
	assert.Equal(t, len(uxs[ux.Body.Address]), 2)
	assert.Equal(t, len(uxs[ux2.Body.Address]), 1)
	got := uxs[ux.Body.Address]
	sort.Sort(got)
	expect := UxArray{ux, ux3}
	sort.Sort(expect)
	assert.Equal(t, got, expect)
	assert.Equal(t, uxs[ux2.Body.Address], UxArray{ux2})
}

func TestUnspentGetMultiple(t *testing.T) {
	unspent := NewUnspentPool()
	// Valid
	txn := Transaction{}
	ux0 := makeUxOut(t)
	ux1 := makeUxOut(t)
	unspent.Add(ux0)
	unspent.Add(ux1)
	assert.Equal(t, len(unspent.Pool), 2)
	txn.PushInput(ux0.Hash())
	txn.PushInput(ux1.Hash())
	txin, err := unspent.GetMultiple(txn.In)
	assert.Nil(t, err)
	assert.Equal(t, len(txin), 2)
	assert.Equal(t, len(txin), len(txn.In))

	// Empty txn
	txn = Transaction{}
	txin, err = unspent.GetMultiple(txn.In)
	assert.Nil(t, err)
	assert.Equal(t, len(txin), 0)

	// Spending unknown output
	txn = makeTransaction(t)
	txn.In[0] = cipher.SHA256{}
	_, err = unspent.GetMultiple(txn.In)
	assertError(t, err, "Unspent output does not exist")

	// Multiple inputs
	ux0 = makeUxOut(t)
	ux1 = makeUxOut(t)
	unspent.Add(ux0)
	unspent.Add(ux1)
	txn = Transaction{}
	txn.PushInput(ux0.Hash())
	txn.PushInput(ux1.Hash())
	txn.PushOutput(genAddress, ux0.Body.Coins+ux1.Body.Coins, ux0.Body.Hours)
	txn.SignInputs([]SecKey{genSecret, genSecret})
	txn.UpdateHeader()
	assert.Nil(t, txn.Verify())
	txin, err = unspent.GetMultiple(txn.In)
	assert.Nil(t, err)
	assert.Equal(t, len(txin), 2)
	assert.Equal(t, txin[0], ux0)
	assert.Equal(t, txin[1], ux1)

	// Duplicate tx.In
	txn = Transaction{}
	txn.In = append(txn.In, ux0.Hash())
	txn.In = append(txn.In, txn.In[0])
	txn.In = append(txn.In, txn.In[0])
	txin, err = unspent.GetMultiple(txn.In)
	assert.Nil(t, err)
	assert.Equal(t, len(txin), 3)
	assert.Equal(t, len(txin), len(txn.In))
	assert.Equal(t, txin[0], ux0)
	assert.Equal(t, txin[1], ux0)
	assert.Equal(t, txin[2], ux0)
	assert.True(t, txin.HasDupes())
}

func TestUnspentCollides(t *testing.T) {
	unspent := NewUnspentPool()
	assert.False(t, unspent.Collides([]SHA256{}))
	assert.False(t, unspent.Collides([]SHA256{randSHA256(t)}))
	ux := makeUxOut(t)
	unspent.Add(ux)
	assert.False(t, unspent.Collides([]SHA256{}))
	assert.False(t, unspent.Collides([]SHA256{randSHA256(t)}))
	assert.True(t, unspent.Collides([]SHA256{ux.Hash()}))
	assert.True(t, unspent.Collides([]SHA256{randSHA256(t), ux.Hash()}))
	assert.True(t, unspent.Collides([]SHA256{ux.Hash(), randSHA256(t)}))
}
