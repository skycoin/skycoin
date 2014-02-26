package visor

import (
    "bytes"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/stretchr/testify/assert"
    "log"
    "sort"
    "testing"
)

func assertError(t *testing.T, err error, msg string) {
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), msg)
}

func makeAddress() coin.Address {
    p, _ := coin.GenerateKeyPair()
    return coin.AddressFromPubKey(p)
}

func makeUxBalances(b []Balance, headTime uint64) coin.UxArray {
    uxs := make(coin.UxArray, len(b))
    for i, _ := range b {
        uxs[i] = coin.UxOut{
            Head: coin.UxHead{
                Time: headTime,
            },
            Body: coin.UxBody{
                SrcTransaction: randSHA256(),
                Address:        makeAddress(),
                Coins:          b[i].Coins,
                Hours:          b[i].Hours,
            },
        }
    }
    return uxs
}

func makeUxBalancesForAddresses(b []Balance, headTime uint64,
    addrs []coin.Address) coin.UxArray {
    if len(b) != len(addrs) {
        log.Panic("Need as many addresses and balances")
    }
    uxs := makeUxBalances(b, headTime)
    for i, _ := range uxs {
        uxs[i].Head.BkSeq = uint64(i)
        uxs[i].Body.Address = addrs[i]
    }
    return uxs
}

func makeUxOut(t *testing.T) coin.UxOut {
    return coin.UxOut{
        Head: coin.UxHead{
            BkSeq: 1,
            Time:  coin.Now(),
        },
        Body: coin.UxBody{
            SrcTransaction: randSHA256(),
            Address:        makeAddress(),
            Coins:          1e6,
            Hours:          1024,
        },
    }
}

func makeUxArray(t *testing.T, n int) coin.UxArray {
    uxa := make(coin.UxArray, n)
    for i, _ := range uxa {
        uxa[i] = makeUxOut(t)
    }
    return uxa
}

func addUxArrayToUnspentPool(u *coin.UnspentPool, uxs coin.UxArray) {
    for _, ux := range uxs {
        u.Add(ux)
    }
}

func TestOldestUxOut(t *testing.T) {
    uxs := OldestUxOut(makeUxArray(t, 4))
    for i, _ := range uxs {
        uxs[i].Head.BkSeq = uint64(i)
    }
    assert.True(t, sort.IsSorted(uxs))
    assert.Equal(t, uxs.Len(), 4)

    uxs.Swap(0, 1)
    assert.False(t, sort.IsSorted(uxs))
    assert.Equal(t, uxs[0].Head.BkSeq, uint64(1))
    assert.Equal(t, uxs[1].Head.BkSeq, uint64(0))
    uxs.Swap(0, 1)
    assert.True(t, sort.IsSorted(uxs))
    assert.Equal(t, uxs[0].Head.BkSeq, uint64(0))
    assert.Equal(t, uxs[1].Head.BkSeq, uint64(1))

    // Test hash sorting
    uxs[1].Head.BkSeq = uint64(0)
    h0 := uxs[0].Hash()
    h1 := uxs[1].Hash()
    firstLower := bytes.Compare(h0[:], h1[:]) < 0
    if firstLower {
        uxs.Swap(0, 1)
    }
    assert.False(t, sort.IsSorted(uxs))
    sort.Sort(uxs)

    cmpHash := false
    cmpSeq := false
    for i, _ := range uxs[:len(uxs)-1] {
        j := i + 1
        if uxs[i].Head.BkSeq == uxs[j].Head.BkSeq {
            ih := uxs[i].Hash()
            jh := uxs[j].Hash()
            assert.True(t, bytes.Compare(ih[:], jh[:]) < 0)
            cmpHash = true
        } else {
            assert.True(t, uxs[i].Head.BkSeq < uxs[j].Head.BkSeq)
            cmpSeq = true
        }
    }
    assert.True(t, cmpHash)
    assert.True(t, cmpSeq)

    // Duplicate output panics
    uxs = append(uxs, uxs[0])
    assert.Panics(t, func() { sort.Sort(uxs) })
}

func TestCreateSpendsNotEnoughCoins(t *testing.T) {
    now := coin.Now()
    amt := Balance{10e6, 100}
    uxs := makeUxBalances([]Balance{
        Balance{1e6, 100},
        Balance{8e6, 0},
    }, now)
    _, err := createSpends(now, uxs, amt, 0, 0)
    assertError(t, err, "Not enough coins")
}

func TestCreateSpendsNotEnoughHours(t *testing.T) {
    now := coin.Now()
    amt := Balance{10e6, 110}
    uxs := makeUxBalances([]Balance{
        Balance{2e6, 100},
        Balance{8e6, 0},
    }, now)
    _, err := createSpends(now, uxs, amt, 0, 0)
    assertError(t, err, "Not enough hours")
}

func TestIgnoreBadCoins(t *testing.T) {
    // We would satisfy this spend if the bad coins were not skipped
    now := coin.Now()
    amt := Balance{10e6, 100}
    uxs := makeUxBalances([]Balance{
        Balance{2e6, 50},
        Balance{8e6, 0},
        Balance{0, 100},
        Balance{1e6 + 1, 100},
    }, now)
    _, err := createSpends(now, uxs, amt, 0, 0)
    assertError(t, err, "Not enough hours")
}

func TestBadSpending(t *testing.T) {
    _, err := createSpends(coin.Now(), coin.UxArray{},
        Balance{1e6 + 1, 1000}, 0, 1)
    assertError(t, err, "Coins must be multiple of 1e6")
    _, err = createSpends(coin.Now(), coin.UxArray{},
        Balance{0, 100}, 0, 1)
    assertError(t, err, "Zero spend amount")
}

func TestCreateSpendsExact(t *testing.T) {
    now := coin.Now()
    amt := Balance{10e6, 100}
    uxs := makeUxBalances([]Balance{
        Balance{1e6, 50},
        Balance{8e6, 40},
        Balance{2e6, 60},
    }, now)
    // Force them to get sorted
    uxs[2].Head.BkSeq = uint64(0)
    uxs[1].Head.BkSeq = uint64(1)
    uxs[0].Head.BkSeq = uint64(2)
    cuxs := append(coin.UxArray{}, uxs...)
    spends, err := createSpends(now, uxs, amt, 0, 0)
    assert.Nil(t, err)
    assert.Equal(t, len(spends), 2)
    assert.Equal(t, spends, coin.UxArray{cuxs[2], cuxs[1]})
}

func TestCreateSpends(t *testing.T) {
    now := coin.Now()
    amt := Balance{12e6, 125}
    uxs := makeUxBalances([]Balance{
        Balance{1e6, 50},
        Balance{8e6, 10}, // 3
        Balance{2e6, 80}, // 2
        Balance{5e6, 15}, // 4
        Balance{7e6, 20}, // 1
    }, now)
    uxs[4].Head.BkSeq = uint64(1)
    uxs[3].Head.BkSeq = uint64(4)
    uxs[2].Head.BkSeq = uint64(2)
    uxs[1].Head.BkSeq = uint64(3)
    uxs[0].Head.BkSeq = uint64(5)
    if sort.IsSorted(OldestUxOut(uxs)) {
        uxs[0], uxs[1] = uxs[1], uxs[0]
    }
    assert.False(t, sort.IsSorted(OldestUxOut(uxs)))
    expectedSorting := coin.UxArray{uxs[4], uxs[2], uxs[1], uxs[3], uxs[0]}
    cuxs := append(coin.UxArray{}, uxs...)
    sort.Sort(OldestUxOut(cuxs))
    assert.Equal(t, expectedSorting, cuxs)
    assert.True(t, sort.IsSorted(OldestUxOut(cuxs)))
    assert.False(t, sort.IsSorted(OldestUxOut(uxs)))

    ouxs := append(coin.UxArray{}, uxs...)
    spends, err := createSpends(now, uxs, amt, 0, 0)
    assert.True(t, sort.IsSorted(OldestUxOut(uxs)))
    assert.Nil(t, err)
    assert.Equal(t, spends, cuxs[:len(spends)])
    assert.Equal(t, len(spends), 4)
    assert.Equal(t, spends, coin.UxArray{ouxs[4], ouxs[2], ouxs[1], ouxs[3]})

    // Recalculate what it should be
    b := Balance{0, 0}
    ouxs = make(coin.UxArray, 0, len(spends))
    for _, ux := range cuxs {
        if b.Coins > amt.Coins && b.Hours >= amt.Hours {
            break
        }
        b = b.Add(NewBalanceFromUxOut(now, &ux))
        ouxs = append(ouxs, ux)
    }
    assert.Equal(t, len(ouxs), len(spends))
    assert.Equal(t, ouxs, spends)
}

func TestCreateSpendingTransaction(t *testing.T) {
    // Setup
    w := NewSimpleWallet()
    w.Populate(4)
    uncf := NewUnconfirmedTxnPool()
    now := coin.Now()
    a := makeAddress()

    // Failing createSpends
    amt := Balance{0, 0}
    unsp := coin.NewUnspentPool()
    _, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, 0, 0, a)
    assert.NotNil(t, err)

    // Valid txn, fee, no change
    uxs := makeUxBalancesForAddresses([]Balance{
        Balance{10e6, 150},
        Balance{15e6, 150},
    }, now, w.GetAddresses()[:2])
    unsp = coin.NewUnspentPool()
    addUxArrayToUnspentPool(&unsp, uxs)
    amt = Balance{25e6, 200}
    tx, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
    assert.Nil(t, err)
    assert.Equal(t, len(tx.Out), 1)
    assert.Equal(t, tx.Out[0], coin.TransactionOutput{
        Coins:   25e6,
        Hours:   200,
        Address: a,
    })
    assert.Equal(t, len(tx.In), 2)
    assert.Equal(t, tx.In, []coin.SHA256{uxs[0].Hash(), uxs[1].Hash()})
    assert.Nil(t, tx.Verify())

    // Valid txn, change
    uxs = makeUxBalancesForAddresses([]Balance{
        Balance{10e6, 150},
        Balance{15e6, 200},
        Balance{1e6, 125},
    }, now, w.GetAddresses()[:3])
    unsp = coin.NewUnspentPool()
    addUxArrayToUnspentPool(&unsp, uxs)
    amt = Balance{25e6, 200}
    tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
    assert.Nil(t, err)
    assert.Equal(t, len(tx.Out), 2)
    assert.Equal(t, tx.Out[0], coin.TransactionOutput{
        Coins:   1e6,
        Hours:   (150 + 200 + 125) - (200 + 100),
        Address: w.GetAddresses()[0],
    })
    assert.Equal(t, tx.Out[1], coin.TransactionOutput{
        Coins:   25e6,
        Hours:   200,
        Address: a,
    })
    assert.Equal(t, len(tx.In), 3)
    assert.Equal(t, tx.In, []coin.SHA256{
        uxs[0].Hash(), uxs[1].Hash(), uxs[2].Hash(),
    })
    assert.Nil(t, tx.Verify())

    // Valid txn, but wastes coin hours
    uxs = makeUxBalancesForAddresses([]Balance{
        Balance{10e6, 150},
        Balance{15e6, 200},
    }, now, w.GetAddresses()[:2])
    unsp = coin.NewUnspentPool()
    addUxArrayToUnspentPool(&unsp, uxs)
    amt = Balance{25e6, 200}
    _, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
    assertError(t, err, "Have enough coins, but not enough to send coin "+
        "hours change back. Would spend 50 more hours than requested.")

    // Would be valid, but unconfirmed subtraction causes it to not be
    // First, make a txn to subtract
    uxs = makeUxBalancesForAddresses([]Balance{
        Balance{10e6, 150},
        Balance{15e6, 150},
    }, now, w.GetAddresses()[:2])
    unsp = coin.NewUnspentPool()
    addUxArrayToUnspentPool(&unsp, uxs)
    amt = Balance{25e6, 200}
    tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
    assert.Nil(t, err)
    // Add it to the unconfirmed pool (bypass RecordTxn to avoid blockchain)
    uncf.Txns[tx.Hash()] = uncf.createUnconfirmedTxn(&unsp, tx,
        w.GetAddressSet())
    // Make a spend that must not reuse previous addresses
    _, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
    assertError(t, err, "Not enough coins")
}
