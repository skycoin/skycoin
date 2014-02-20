package coin

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "sort"
    "testing"
)

func randSHA256(t *testing.T) SHA256 {
    return SumSHA256(randBytes(t, 128))
}

func makeUxBody(t *testing.T) UxBody {
    body, _ := makeUxBodyWithSecret(t)
    return body
}

func makeUxOut(t *testing.T) UxOut {
    ux, _ := makeUxOutWithSecret(t)
    return ux
}

func makeUxBodyWithSecret(t *testing.T) (UxBody, SecKey) {
    p, s := GenerateKeyPair()
    return UxBody{
        SrcTransaction: SumSHA256(randBytes(t, 128)),
        Address:        AddressFromPubKey(p),
        Coins:          1e6,
        Hours:          100,
    }, s
}

func makeUxOutWithSecret(t *testing.T) (UxOut, SecKey) {
    body, sec := makeUxBodyWithSecret(t)
    return UxOut{
        Head: UxHead{
            Time:  100,
            BkSeq: 2,
        },
        Body: body,
    }, sec
}

func TestUxBodyHash(t *testing.T) {
    uxb := makeUxBody(t)
    h := uxb.Hash()
    assert.NotEqual(t, h, SHA256{})
}

func TestUxOutHash(t *testing.T) {
    uxb := makeUxBody(t)
    uxo := UxOut{Body: uxb}
    assert.Equal(t, uxb.Hash(), uxo.Hash())
    // Head should not affect hash
    uxo.Head = UxHead{0, 1}
    assert.Equal(t, uxb.Hash(), uxo.Hash())
}

func TestUxOutCoinHours(t *testing.T) {
    uxo := makeUxOut(t)
    // No hours passed
    now := uint64(200)
    assert.Equal(t, uxo.CoinHours(now), uxo.Body.Hours)
    now = uint64(3600) + uxo.Head.Time
    assert.Equal(t, uxo.CoinHours(now), uxo.Body.Hours+(uxo.Body.Coins/1e6))
    now = uint64(3600*6) + uxo.Head.Time
    assert.Equal(t, uxo.CoinHours(now), uxo.Body.Hours+(uxo.Body.Coins/1e6)*6)
    now = uxo.Head.Time / 2
    assert.Equal(t, uxo.CoinHours(now), uxo.Body.Hours)
    uxo.Body.Coins = genesisCoinVolume
    uxo.Body.Hours = genesisCoinHours
    assert.Equal(t, uxo.CoinHours(uxo.Head.Time), uxo.Body.Hours)
    assert.Equal(t, uxo.CoinHours(uxo.Head.Time+3600),
        uxo.Body.Hours+(genesisCoinVolume/1e6))
}

func makeUxArray(t *testing.T, n int) UxArray {
    uxa := make(UxArray, 4)
    for i := 0; i < len(uxa); i++ {
        uxa[i] = makeUxOut(t)
    }
    return uxa
}

func TestUxArrayHashArray(t *testing.T) {
    uxa := makeUxArray(t, 4)
    hashes := uxa.Hashes()
    assert.Equal(t, len(hashes), len(uxa))
    for i, h := range hashes {
        assert.Equal(t, h, uxa[i].Hash())
    }
}

func TestUxArrayHasDupes(t *testing.T) {
    uxa := makeUxArray(t, 4)
    assert.False(t, uxa.HasDupes())
    uxa[0] = uxa[1]
    assert.True(t, uxa.HasDupes())
}

func TestUxArrayRemoveDupes(t *testing.T) {
    uxa := makeUxArray(t, 4)
    assert.False(t, uxa.HasDupes())
    assert.Equal(t, uxa, uxa.removeDupes())
    uxa[0] = uxa[1]
    assert.True(t, uxa.HasDupes())
    uxb := uxa.removeDupes()
    assert.False(t, uxb.HasDupes())
    assert.Equal(t, len(uxb), 3)
    assert.Equal(t, uxb[0], uxa[0])
    assert.Equal(t, uxb[1], uxa[2])
    assert.Equal(t, uxb[2], uxa[3])
}

func TestUxArraySub(t *testing.T) {
    uxa := makeUxArray(t, 4)
    uxb := makeUxArray(t, 4)
    uxc := append(uxa[:1], uxb...)
    uxc = append(uxc, uxa[1:2]...)

    uxd := uxc.Sub(uxa)
    assert.Equal(t, uxd, uxb)

    uxd = uxc.Sub(uxb)
    assert.Equal(t, len(uxd), 2)
    assert.Equal(t, uxd, uxa[:2])

    // No intersection
    uxd = uxa.Sub(uxb)
    assert.Equal(t, uxa, uxd)
    uxd = uxb.Sub(uxa)
    assert.Equal(t, uxd, uxb)
}

func manualUxArrayIsSorted(uxa UxArray) bool {
    isSorted := true
    for i := 0; i < len(uxa)-1; i++ {
        hi := uxa[i].Hash()
        hj := uxa[i+1].Hash()
        if bytes.Compare(hi[:], hj[:]) > 0 {
            isSorted = false
        }
    }
    return isSorted
}

func TestUxArraySorting(t *testing.T) {
    uxa := make(UxArray, 4)
    for i := 0; i < len(uxa); i++ {
        uxa[i] = makeUxOut(t)
    }
    isSorted := manualUxArrayIsSorted(uxa)
    assert.Equal(t, sort.IsSorted(uxa), isSorted)
    assert.Equal(t, uxa.IsSorted(), isSorted)
    // Make sure uxa is not sorted
    if isSorted {
        uxa[0], uxa[1] = uxa[1], uxa[0]
    }
    assert.False(t, manualUxArrayIsSorted(uxa))
    assert.False(t, sort.IsSorted(uxa))
    assert.False(t, uxa.IsSorted())
    uxb := make(UxArray, 4)
    for i, ux := range uxa {
        uxb[i] = ux
    }
    sort.Sort(uxa)
    assert.True(t, sort.IsSorted(uxa))
    assert.True(t, manualUxArrayIsSorted(uxa))
    assert.True(t, uxa.IsSorted())
    assert.False(t, sort.IsSorted(uxb))
    uxb.Sort()
    assert.Equal(t, uxa, uxb)
    assert.True(t, sort.IsSorted(uxb))
    assert.True(t, manualUxArrayIsSorted(uxb))
    assert.True(t, uxb.IsSorted())
}

func TestUxArrayLen(t *testing.T) {
    uxa := make(UxArray, 4)
    assert.Equal(t, len(uxa), uxa.Len())
    assert.Equal(t, 4, uxa.Len())
}

func TestUxArrayLess(t *testing.T) {
    uxa := make(UxArray, 2)
    uxa[0] = makeUxOut(t)
    uxa[1] = makeUxOut(t)
    h := make([]SHA256, 2)
    h[0] = uxa[0].Hash()
    h[1] = uxa[1].Hash()
    assert.Equal(t, uxa.Less(0, 1), bytes.Compare(h[0][:], h[1][:]) < 0)
    assert.Equal(t, uxa.Less(1, 0), bytes.Compare(h[0][:], h[1][:]) > 0)
}

func TestUxArraySwap(t *testing.T) {
    uxa := make(UxArray, 2)
    uxx := makeUxOut(t)
    uxy := makeUxOut(t)
    uxa[0] = uxx
    uxa[1] = uxy
    uxa.Swap(0, 1)
    assert.Equal(t, uxa[0], uxy)
    assert.Equal(t, uxa[1], uxx)
    uxa.Swap(0, 1)
    assert.Equal(t, uxa[0], uxx)
    assert.Equal(t, uxa[1], uxy)
    uxa.Swap(1, 0)
    assert.Equal(t, uxa[1], uxx)
    assert.Equal(t, uxa[0], uxy)
}

func TestNewUnspentPool(t *testing.T) {
    up := NewUnspentPool()
    assert.Equal(t, len(up.Arr), 0)
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, up.XorHash, SHA256{})
}

func TestUnspentPoolRebuild(t *testing.T) {
    up := NewUnspentPool()
    up.Arr = append(up.Arr, makeUxOut(t))
    up.Arr = append(up.Arr, makeUxOut(t))
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, up.XorHash, SHA256{})
    up.Rebuild()
    assert.Equal(t, len(up.hashIndex), 2)
    for _, x := range up.Arr {
        xi, ok := up.hashIndex[x.Hash()]
        assert.True(t, ok)
        assert.Equal(t, x, up.Arr[xi])
    }
    h := SHA256{}
    h = h.Xor(up.Arr[0].Hash())
    h = h.Xor(up.Arr[1].Hash())
    assert.Equal(t, up.XorHash, h)
    assert.NotEqual(t, up.XorHash, SHA256{})

    // Duplicate item in array causes panic
    up.Arr = append(up.Arr, up.Arr[0])
    assert.Panics(t, up.Rebuild)
}

func TestUnspentPoolAdd(t *testing.T) {
    up := NewUnspentPool()
    ux := makeUxOut(t)
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, len(up.Arr), 0)
    up.Add(ux)
    assert.Equal(t, len(up.hashIndex), 1)
    assert.Equal(t, len(up.Arr), 1)
    uxi, ok := up.hashIndex[ux.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux)
    assert.NotEqual(t, up.XorHash, SHA256{})
    assert.Equal(t, up.XorHash, ux.Hash())
    // Duplicate add doesnt change state
    h := up.XorHash
    up.Add(ux)
    assert.Equal(t, len(up.hashIndex), 1)
    assert.Equal(t, len(up.Arr), 1)
    assert.Equal(t, up.XorHash, h)
    // Add a 2nd is ok
    ux2 := makeUxOut(t)
    up.Add(ux2)
    uxi, ok = up.hashIndex[ux2.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux2)
    assert.Equal(t, len(up.hashIndex), 2)
    assert.Equal(t, len(up.Arr), 2)
    h = ux.Hash()
    h = h.Xor(ux2.Hash())
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

func TestUnspentPoolDelFromArray(t *testing.T) {
    up := NewUnspentPool()
    ux := makeUxOut(t)
    up.Add(ux)
    assert.Equal(t, len(up.Arr), 1)
    assert.NotPanics(t, func() { up.delFromArray(0) })
    assert.Equal(t, len(up.Arr), 0)

    up = NewUnspentPool()
    ux2 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    assert.Equal(t, len(up.Arr), 2)
    assert.NotPanics(t, func() { up.delFromArray(1) })
    assert.Equal(t, len(up.Arr), 1)
    assert.Equal(t, up.Arr[0], ux)
    up.delFromArray(0)
    assert.Equal(t, len(up.Arr), 0)

    up = NewUnspentPool()
    ux3 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    up.Add(ux3)
    assert.Equal(t, len(up.Arr), 3)
    assert.NotPanics(t, func() { up.delFromArray(1) })
    assert.Equal(t, len(up.Arr), 2)
    assert.Equal(t, up.Arr[0], ux)
    assert.Equal(t, up.Arr[1], ux3)
    assert.NotPanics(t, func() { up.delFromArray(0) })
    assert.Equal(t, len(up.Arr), 1)
    assert.Equal(t, up.Arr[0], ux3)
}

func TestUnspentPoolDelPrivate(t *testing.T) {
    up := NewUnspentPool()
    ux := makeUxOut(t)
    assert.Equal(t, up.del(ux.Hash()), -1)
    up.Add(ux)
    assert.Equal(t, up.del(ux.Hash()), 0)
    assert.Equal(t, len(up.Arr), 0)
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, up.XorHash, SHA256{})

    ux2 := makeUxOut(t)
    ux3 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    up.Add(ux3)
    assert.Equal(t, len(up.Arr), 3)
    assert.Equal(t, len(up.hashIndex), 3)
    h := ux.Hash()
    h = h.Xor(ux2.Hash())
    h = h.Xor(ux3.Hash())
    assert.Equal(t, up.XorHash, h)
    assert.Equal(t, up.del(ux2.Hash()), 1)
    assert.Equal(t, len(up.Arr), 2)
    assert.Equal(t, len(up.hashIndex), 2)
    uxi, ok := up.hashIndex[ux.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux)
    uxi, ok = up.hashIndex[ux3.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi-1], ux3)
    h = ux.Hash()
    h = h.Xor(ux3.Hash())
    assert.Equal(t, up.XorHash, h)
}

func TestUnspentPoolDelAt(t *testing.T) {
    up := NewUnspentPool()
    ux := makeUxOut(t)
    up.Add(ux)
    up.delAt(0)
    assert.Equal(t, len(up.Arr), 0)
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, up.XorHash, SHA256{})

    ux2 := makeUxOut(t)
    ux3 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    up.Add(ux3)
    assert.Equal(t, len(up.Arr), 3)
    assert.Equal(t, len(up.hashIndex), 3)
    h := ux.Hash()
    h = h.Xor(ux2.Hash())
    h = h.Xor(ux3.Hash())
    assert.Equal(t, up.XorHash, h)
    up.delAt(1)
    assert.Equal(t, len(up.Arr), 2)
    assert.Equal(t, len(up.hashIndex), 2)
    uxi, ok := up.hashIndex[ux.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux)
    uxi, ok = up.hashIndex[ux3.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi-1], ux3)
    h = ux.Hash()
    h = h.Xor(ux3.Hash())
    assert.Equal(t, up.XorHash, h)
}

func TestUnspentPoolUpdateIndices(t *testing.T) {
    up := NewUnspentPool()
    assert.Equal(t, len(up.hashIndex), 0)
    assert.NotPanics(t, func() { up.updateIndices(100) })
    assert.Equal(t, len(up.hashIndex), 0)
    assert.NotPanics(t, func() { up.updateIndices(0) })
    assert.Equal(t, len(up.hashIndex), 0)
    assert.NotPanics(t, func() { up.updateIndices(1) })
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Panics(t, func() { up.updateIndices(-1) })
    assert.Equal(t, len(up.hashIndex), 0)

    ux := makeUxOut(t)
    ux2 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    up.hashIndex[ux.Hash()] = 70
    up.hashIndex[ux2.Hash()] = 71
    up.updateIndices(1)
    assert.Equal(t, up.hashIndex[ux.Hash()], 70)
    assert.Equal(t, up.hashIndex[ux2.Hash()], 1)
    up.hashIndex[ux2.Hash()] = 71
    up.updateIndices(0)
    assert.Equal(t, up.hashIndex[ux.Hash()], 0)
    assert.Equal(t, up.hashIndex[ux2.Hash()], 1)
}

func TestUnspentPoolDel(t *testing.T) {
    up := NewUnspentPool()
    ux := makeUxOut(t)
    ux2 := makeUxOut(t)
    ux3 := makeUxOut(t)
    up.Add(ux)
    up.Add(ux2)
    up.Add(ux3)
    assert.Equal(t, len(up.Arr), 3)
    assert.Equal(t, len(up.hashIndex), 3)
    // Unknown hash
    up.Del(SHA256{})
    assert.Equal(t, len(up.Arr), 3)
    assert.Equal(t, len(up.hashIndex), 3)
    // Delete middle one
    up.Del(ux2.Hash())
    assert.Equal(t, len(up.Arr), 2)
    assert.Equal(t, len(up.hashIndex), 2)
    uxi, ok := up.hashIndex[ux.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux)
    uxi, ok = up.hashIndex[ux3.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux3)
    h := ux.Hash()
    h = h.Xor(ux3.Hash())
    assert.Equal(t, up.XorHash, h)
    // Delete first one
    up.Del(ux.Hash())
    assert.Equal(t, len(up.Arr), 1)
    assert.Equal(t, len(up.hashIndex), 1)
    uxi, ok = up.hashIndex[ux3.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux3)
    assert.Equal(t, up.XorHash, ux3.Hash())
    // Delete remaining one
    up.Del(ux3.Hash())
    assert.Equal(t, len(up.Arr), 0)
    assert.Equal(t, len(up.hashIndex), 0)
    assert.Equal(t, up.XorHash, SHA256{})
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
    assert.Equal(t, len(up.Arr), 4)
    assert.Equal(t, len(up.hashIndex), 4)
    // Delete 1st and 3rd and an unknown
    up.DelMultiple([]SHA256{ux.Hash(), ux3.Hash(), ux5.Hash()})
    assert.Equal(t, len(up.Arr), 2)
    assert.Equal(t, len(up.hashIndex), 2)
    uxi, ok := up.hashIndex[ux2.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux2)
    uxi, ok = up.hashIndex[ux4.Hash()]
    assert.True(t, ok)
    assert.Equal(t, up.Arr[uxi], ux4)
    h := ux2.Hash()
    h = h.Xor(ux4.Hash())
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
    assert.Equal(t, len(unspent.Arr), 2)
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
    txn.In[0] = SHA256{}
    _, err = unspent.GetMultiple(txn.In)
    assertError(t, err, "Unspent output does not exist")

    // Multiple inputs
    unspent.Add(makeUxOut(t))
    unspent.Add(makeUxOut(t))
    txn = Transaction{}
    ux0 = unspent.Arr[0]
    ux1 = unspent.Arr[1]
    txn.PushInput(ux0.Hash())
    txn.PushInput(ux1.Hash())
    txn.PushOutput(genAddress, ux0.Body.Coins+ux1.Body.Coins, ux0.Body.Hours)
    txn.SignInputs([]SecKey{genSecret, genSecret})
    txn.UpdateHeader()
    assert.Nil(t, txn.Verify(testMaxSize))
    txin, err = unspent.GetMultiple(txn.In)
    assert.Nil(t, err)
    assert.Equal(t, len(txin), 2)
    assert.Equal(t, txin[0], ux0)
    assert.Equal(t, txin[1], ux1)

    // Duplicate tx.In
    unspent.Add(makeUxOut(t))
    txn = Transaction{}
    txn.In = append(txn.In, unspent.Arr[0].Hash())
    txn.In = append(txn.In, txn.In[0])
    txn.In = append(txn.In, txn.In[0])
    txin, err = unspent.GetMultiple(txn.In)
    assert.Nil(t, err)
    ux0 = unspent.Arr[0]
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

func TestAddressUxOutsKeys(t *testing.T) {
    unspents := make(AddressUxOuts)
    ux := makeUxOut(t)
    ux2 := makeUxOut(t)
    ux3 := makeUxOut(t)
    unspents[ux.Body.Address] = UxArray{ux}
    unspents[ux2.Body.Address] = UxArray{ux2}
    unspents[ux3.Body.Address] = UxArray{ux3}
    keys := unspents.Keys()
    assert.Equal(t, len(keys), 3)
    dupes := make(map[Address]byte, 3)
    for _, k := range keys {
        dupes[k] = byte(1)
        assert.True(t, k == ux.Body.Address || k == ux2.Body.Address ||
            k == ux3.Body.Address)
    }
    assert.Equal(t, len(keys), len(dupes))
}

func TestAddressUxOutsMerge(t *testing.T) {
    unspents := make(AddressUxOuts)
    unspents2 := make(AddressUxOuts)
    ux := makeUxOut(t)
    ux2 := makeUxOut(t)
    ux3 := makeUxOut(t)
    ux4 := makeUxOut(t)
    ux3.Body.Address = ux.Body.Address

    unspents[ux.Body.Address] = UxArray{ux}
    unspents[ux2.Body.Address] = UxArray{ux2}
    unspents2[ux3.Body.Address] = UxArray{ux3}
    unspents2[ux4.Body.Address] = UxArray{ux4}

    // Valid merge
    keys := []Address{ux.Body.Address, ux2.Body.Address, ux4.Body.Address}
    merged := unspents.Merge(unspents2, keys)
    assert.Equal(t, len(unspents), 2)
    assert.Equal(t, len(unspents2), 2)
    assert.Equal(t, len(merged), 3)
    assert.Equal(t, merged[ux.Body.Address], UxArray{ux, ux3})
    assert.Equal(t, merged[ux2.Body.Address], UxArray{ux2})
    assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})

    // Duplicates should not be merged
    unspents[ux4.Body.Address] = UxArray{ux4}
    unspents[ux.Body.Address] = UxArray{ux, ux3}
    merged = unspents.Merge(unspents2, keys)
    assert.Equal(t, len(merged), 3)
    assert.Equal(t, merged[ux.Body.Address], UxArray{ux, ux3})
    assert.Equal(t, merged[ux2.Body.Address], UxArray{ux2})
    assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})

    // Missing keys should not be merged
    merged = unspents.Merge(unspents2, []Address{})
    assert.Equal(t, len(merged), 0)
    merged = unspents.Merge(unspents2, []Address{ux4.Body.Address})
    assert.Equal(t, len(merged), 1)
    assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})
}

func TestAddressUxOutsSub(t *testing.T) {
    up := make(AddressUxOuts)
    up2 := make(AddressUxOuts)
    uxs := makeUxArray(t, 4)

    uxs[1].Body.Address = uxs[0].Body.Address
    up[uxs[0].Body.Address] = UxArray{uxs[0], uxs[1]}
    up[uxs[2].Body.Address] = UxArray{uxs[2]}
    up[uxs[3].Body.Address] = UxArray{uxs[3]}

    up2[uxs[0].Body.Address] = UxArray{uxs[0]}
    up2[uxs[2].Body.Address] = UxArray{uxs[2]}

    up3 := up.Sub(up2)
    // One address should have been removed, because no elements
    assert.Equal(t, len(up3), 2)
    _, ok := up3[uxs[2].Body.Address]
    assert.False(t, ok)
    // Ux3 should be untouched
    ux3 := up3[uxs[3].Body.Address]
    assert.Equal(t, ux3, UxArray{uxs[3]})
    // Ux0,Ux1 should be missing Ux0
    ux1 := up3[uxs[0].Body.Address]
    assert.Equal(t, ux1, UxArray{uxs[1]})

    // Originals should be unmodified
    assert.Equal(t, len(up), 3)
    assert.Equal(t, len(up[uxs[0].Body.Address]), 2)
    assert.Equal(t, len(up[uxs[2].Body.Address]), 1)
    assert.Equal(t, len(up[uxs[3].Body.Address]), 1)
    assert.Equal(t, len(up2), 2)
    assert.Equal(t, len(up2[uxs[0].Body.Address]), 1)
    assert.Equal(t, len(up2[uxs[2].Body.Address]), 1)
}

func TestAddressUxOutsFlatten(t *testing.T) {
    up := make(AddressUxOuts)
    uxs := makeUxArray(t, 3)
    uxs[2].Body.Address = uxs[1].Body.Address
    emptyAddr := makeAddress()

    // An empty array
    up[emptyAddr] = UxArray{}
    // 1 element array
    up[uxs[0].Body.Address] = UxArray{uxs[0]}
    // 2 element array
    up[uxs[1].Body.Address] = UxArray{uxs[1], uxs[2]}

    flat := up.Flatten()
    assert.Equal(t, len(flat), 3)
    // emptyAddr should not be in the array
    for _, ux := range flat {
        assert.NotEqual(t, ux.Body.Address, emptyAddr)
    }
    if flat[0].Body.Address == uxs[0].Body.Address {
        assert.Equal(t, flat[0], uxs[0])
        assert.Equal(t, flat[0].Body.Address, uxs[0].Body.Address)
        assert.Equal(t, flat[0+1], uxs[1])
        assert.Equal(t, flat[1+1], uxs[2])
        assert.Equal(t, flat[0+1].Body.Address, uxs[1].Body.Address)
        assert.Equal(t, flat[1+1].Body.Address, uxs[2].Body.Address)
    } else {
        assert.Equal(t, flat[0], uxs[1])
        assert.Equal(t, flat[1], uxs[2])
        assert.Equal(t, flat[0].Body.Address, uxs[1].Body.Address)
        assert.Equal(t, flat[1].Body.Address, uxs[2].Body.Address)
        assert.Equal(t, flat[2], uxs[0])
        assert.Equal(t, flat[2].Body.Address, uxs[0].Body.Address)
    }
}
