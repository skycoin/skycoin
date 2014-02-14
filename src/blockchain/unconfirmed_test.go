package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/stretchr/testify/assert"
    "sort"
    "testing"
    "time"
)

func makeValidTxn(mv *Visor) (coin.Transaction, error) {
    we := NewWalletEntry()
    return mv.Spend(Balance{10 * 1e6, 0}, 0, we.Address)
}

func makeValidTxnNoChange(mv *Visor) (coin.Transaction, error) {
    we := NewWalletEntry()
    b := mv.Balance(mv.Config.MasterKeys.Address)
    return mv.Spend(b, 0, we.Address)
}

func makeInvalidTxn(mv *Visor) (coin.Transaction, error) {
    we := NewWalletEntry()
    txn, err := mv.Spend(Balance{10 * 1e6, 0}, 0, we.Address)
    if err != nil {
        return txn, err
    }
    txn.Out[0].DestinationAddress = coin.Address{}
    return txn, nil
}

func assertValidUnspent(t *testing.T, bc *coin.Blockchain,
    unspent *coin.UnspentPool, tx coin.Transaction) {
    expect := bc.TxUxOut(tx, coin.BlockHeader{})
    assert.NotEqual(t, len(expect), 0)
    assert.Equal(t, len(expect), len(unspent.Arr))
    for _, ux := range expect {
        assert.True(t, unspent.Has(ux.Hash()))
    }
}

func assertValidUnconfirmed(t *testing.T, txns map[coin.SHA256]UnconfirmedTxn,
    txn coin.Transaction, didAnnounce, isOurReceive, isOurSpend bool) {
    ut, ok := txns[txn.Hash()]
    assert.True(t, ok)
    assert.Equal(t, ut.Txn, txn)
    assert.Equal(t, ut.IsOurReceive, isOurReceive)
    assert.Equal(t, ut.IsOurSpend, isOurSpend)
    assert.Equal(t, ut.Announced.IsZero(), !didAnnounce)
    assert.False(t, ut.Received.IsZero())
    assert.False(t, ut.Checked.IsZero())
}

func TestUnconfirmedTxnHash(t *testing.T) {
    utx := createUnconfirmedTxn()
    assert.Equal(t, utx.Hash(), utx.Txn.Hash())
    assert.NotEqual(t, utx.Hash(), utx.Txn.Header.Hash)
}

func TestNewUnconfirmedTxnPool(t *testing.T) {
    ut := NewUnconfirmedTxnPool()
    assert.NotNil(t, ut.Txns)
    assert.Equal(t, len(ut.Txns), 0)
}

func TestSetAnnounced(t *testing.T) {
    ut := NewUnconfirmedTxnPool()
    assert.Equal(t, len(ut.Txns), 0)
    // Unknown should be safe and a noop
    assert.NotPanics(t, func() {
        ut.SetAnnounced(coin.SHA256{}, util.Now())
    })
    assert.Equal(t, len(ut.Txns), 0)
    utx := createUnconfirmedTxn()
    assert.True(t, utx.Announced.IsZero())
    ut.Txns[utx.Hash()] = utx
    now := util.Now()
    ut.SetAnnounced(utx.Hash(), now)
    assert.Equal(t, ut.Txns[utx.Hash()].Announced, now)
}

func TestRecordTxn(t *testing.T) {
    defer cleanupVisor()
    // Test with invalid txn
    mv := setupMasterVisor()
    ut := NewUnconfirmedTxnPool()
    txn, err := makeInvalidTxn(mv)
    assert.Nil(t, err)
    assert.NotNil(t, ut.RecordTxn(mv.blockchain, txn, nil, false))
    assert.Equal(t, len(ut.Txns), 0)

    // Test didAnnounce=false
    mv = setupMasterVisor()
    ut = NewUnconfirmedTxnPool()
    txn, err = makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, nil, false))
    assertValidUnspent(t, mv.blockchain, &ut.Unspent, txn)
    assertValidUnconfirmed(t, ut.Txns, txn, false, false, false)

    // Test didAnnounce=true
    mv = setupMasterVisor()
    ut = NewUnconfirmedTxnPool()
    txn, err = makeValidTxn(mv)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, nil, true))
    assertValidUnspent(t, mv.blockchain, &ut.Unspent, txn)
    assertValidUnconfirmed(t, ut.Txns, txn, true, false, false)

    // Test where we are receiver of ux outputs
    mv = setupMasterVisor()
    assert.Equal(t, len(mv.blockchain.Unspent.Arr), 1)
    ut = NewUnconfirmedTxnPool()
    txn, err = makeValidTxn(mv)
    assert.Nil(t, err)
    addrs := make(map[coin.Address]byte, 1)
    addrs[txn.Out[1].DestinationAddress] = byte(1)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, addrs, false))
    assertValidUnspent(t, mv.blockchain, &ut.Unspent, txn)
    assertValidUnconfirmed(t, ut.Txns, txn, false, true, false)

    // Test where we are spender of ux outputs
    mv = setupMasterVisor()
    assert.Equal(t, len(mv.blockchain.Unspent.Arr), 1)
    ut = NewUnconfirmedTxnPool()
    txn, err = makeValidTxnNoChange(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 1)
    ux, ok := mv.blockchain.Unspent.Get(txn.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, addrs, false))
    assertValidUnspent(t, mv.blockchain, &ut.Unspent, txn)
    assertValidUnconfirmed(t, ut.Txns, txn, false, false, true)

    // Test where we are both spender and receiver of ux outputs
    mv = setupMasterVisor()
    ut = NewUnconfirmedTxnPool()
    txn, err = makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 2)
    addrs[txn.Out[0].DestinationAddress] = byte(1)
    ux, ok = mv.blockchain.Unspent.Get(txn.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, addrs, false))
    assertValidUnspent(t, mv.blockchain, &ut.Unspent, txn)
    assertValidUnconfirmed(t, ut.Txns, txn, false, true, true)
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)

    // Test duplicate Record, should be no-op besides state change
    assert.Nil(t, ut.RecordTxn(mv.blockchain, txn, addrs, true))
    assertValidUnconfirmed(t, ut.Txns, txn, true, true, true)
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)
}

func TestRawTxns(t *testing.T) {
    ut := NewUnconfirmedTxnPool()
    utxs := make(coin.Transactions, 4)
    for i := 0; i < len(utxs); i++ {
        utx := addUnconfirmedTxnToPool(ut)
        utxs[i] = utx.Txn
    }
    sort.Sort(utxs)
    txns := ut.RawTxns()
    sort.Sort(txns)
    for i, tx := range txns {
        assert.Equal(t, utxs[i], tx)
    }
}

func TestRemoveTxn(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    ut := NewUnconfirmedTxnPool()

    utx, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, utx, nil, false))
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)

    // Unknown txn is no-op
    badh := randSHA256()
    assert.NotEqual(t, badh, utx.Hash())
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)
    ut.removeTxn(mv.blockchain, badh)
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)

    // Known txn updates Txns, predicted Unspents
    utx2, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, ut.RecordTxn(mv.blockchain, utx2, nil, false))
    assert.Equal(t, len(ut.Txns), 2)
    assert.Equal(t, len(ut.Unspent.Arr), 4)
    ut.removeTxn(mv.blockchain, utx.Hash())
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)
    ut.removeTxn(mv.blockchain, utx.Hash())
    assert.Equal(t, len(ut.Txns), 1)
    assert.Equal(t, len(ut.Unspent.Arr), 2)
    ut.removeTxn(mv.blockchain, utx2.Hash())
    assert.Equal(t, len(ut.Txns), 0)
    assert.Equal(t, len(ut.Unspent.Arr), 0)
}

func TestRemoveTxns(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    up := NewUnconfirmedTxnPool()

    // Include an unknown hash, and omit a known hash. The other two should
    // be removed.
    hashes := make([]coin.SHA256, 0, 3)
    hashes = append(hashes, randSHA256()) // unknown hash
    ut, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut, nil, false))
    hashes = append(hashes, ut.Hash())
    ut2, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut2, nil, false))
    hashes = append(hashes, ut2.Hash())
    ut3, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut3, nil, false))

    assert.Equal(t, len(up.Unspent.Arr), 3*2)
    assert.Equal(t, len(up.Txns), 3)
    up.removeTxns(mv.blockchain, hashes)
    assert.Equal(t, len(up.Unspent.Arr), 1*2)
    assert.Equal(t, len(up.Txns), 1)
    _, ok := up.Txns[ut3.Hash()]
    assert.True(t, ok)
}

func TestRemoveTransactions(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    up := NewUnconfirmedTxnPool()

    // Include an unknown txn, and omit a known hash. The other two should
    // be removed.
    unkUt, err := makeValidTxn(mv)
    assert.Nil(t, err)
    txns := make(coin.Transactions, 0, 3)
    txns = append(txns, unkUt) // unknown txn
    ut, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut, nil, false))
    txns = append(txns, ut)
    ut2, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut2, nil, false))
    txns = append(txns, ut2)
    ut3, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ut3, nil, false))

    assert.Equal(t, len(up.Unspent.Arr), 3*2)
    assert.Equal(t, len(up.Txns), 3)
    up.RemoveTransactions(mv.blockchain, txns)
    assert.Equal(t, len(up.Unspent.Arr), 1*2)
    assert.Equal(t, len(up.Txns), 1)
    _, ok := up.Txns[ut3.Hash()]
    assert.True(t, ok)
}

func testRefresh(t *testing.T, mv *Visor,
    refresh func(checkPeriod, maxAge time.Duration)) {
    up := mv.UnconfirmedTxns
    // Add a transaction that is invalid, but will not be checked yet
    // Add a transaction that is invalid, and will be checked and removed
    invalidTxUnchecked, err := makeValidTxn(mv)
    assert.Nil(t, err)
    invalidTxChecked, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, invalidTxUnchecked.Verify())
    assert.Nil(t, invalidTxChecked.Verify())
    // Invalidate it by spending the output that this txn references
    invalidator, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, invalidator, nil, false))
    assert.Equal(t, len(up.Txns), 1)
    _, err = mv.CreateAndExecuteBlock()
    assert.Nil(t, err)
    assert.Equal(t, len(up.Txns), 0)
    assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxUnchecked))
    assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxChecked))

    invalidUtxUnchecked := UnconfirmedTxn{
        Txn:       invalidTxUnchecked,
        Received:  util.Now(),
        Checked:   util.Now(),
        Announced: util.ZeroTime(),
    }
    invalidUtxChecked := invalidUtxUnchecked
    invalidUtxChecked.Txn = invalidTxChecked
    invalidUtxUnchecked.Checked = util.Now().Add(time.Hour)
    invalidUtxChecked.Checked = util.Now().Add(-time.Hour)
    up.Txns[invalidUtxUnchecked.Hash()] = invalidUtxUnchecked
    up.Txns[invalidUtxChecked.Hash()] = invalidUtxChecked
    assert.Equal(t, len(up.Txns), 2)
    bh := coin.BlockHeader{}
    for _, ux := range mv.blockchain.TxUxOut(invalidTxUnchecked, bh) {
        up.Unspent.Add(ux)
    }
    for _, ux := range mv.blockchain.TxUxOut(invalidTxChecked, bh) {
        up.Unspent.Add(ux)
    }
    // Add a transaction that is valid, and will not be checked yet
    validTxUnchecked, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, validTxUnchecked, nil, false))
    assert.Equal(t, len(up.Txns), 3)
    validUtxUnchecked := up.Txns[validTxUnchecked.Hash()]
    validUtxUnchecked.Checked = util.Now().Add(time.Hour)
    up.Txns[validUtxUnchecked.Hash()] = validUtxUnchecked
    // Add a transaction that is valid, and will be checked
    validTxChecked, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, validTxChecked, nil, false))
    assert.Equal(t, len(up.Txns), 4)
    validUtxChecked := up.Txns[validTxChecked.Hash()]
    validUtxChecked.Checked = util.Now().Add(-time.Hour)
    up.Txns[validUtxChecked.Hash()] = validUtxChecked
    // Add a transaction that is expired
    validTxExpired, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, validTxExpired, nil, false))
    assert.Equal(t, len(up.Txns), 5)
    validUtxExpired := up.Txns[validTxExpired.Hash()]
    validUtxExpired.Received = util.Now().Add(-time.Hour)
    up.Txns[validTxExpired.Hash()] = validUtxExpired

    // Pre-sanity check
    assert.Equal(t, len(up.Unspent.Arr), 2*5)
    assert.Equal(t, len(up.Txns), 5)

    // Refresh
    checkPeriod := time.Second * 2
    maxAge := time.Second * 4
    refresh(checkPeriod, maxAge)

    // All utxns that are unchecked should be exactly the same
    assert.Equal(t, up.Txns[validUtxUnchecked.Hash()], validUtxUnchecked)
    assert.Equal(t, up.Txns[invalidUtxUnchecked.Hash()], invalidUtxUnchecked)
    // The valid one that is checked should have its checked status updated
    validUtxCheckedUpdated := up.Txns[validUtxChecked.Hash()]
    assert.True(t, validUtxCheckedUpdated.Checked.After(validUtxChecked.Checked))
    validUtxChecked.Checked = validUtxCheckedUpdated.Checked
    assert.Equal(t, validUtxChecked, validUtxCheckedUpdated)
    // The invalid checked one and the expired one should be removed
    _, ok := up.Txns[invalidUtxChecked.Hash()]
    assert.False(t, ok)
    _, ok = up.Txns[validUtxExpired.Hash()]
    assert.False(t, ok)
    // Also, the unspents should have 2 * nRemaining
    assert.Equal(t, len(up.Unspent.Arr), 2*3)
    assert.Equal(t, len(up.Txns), 3)
}

func TestRefresh(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    testRefresh(t, mv, func(checkPeriod, maxAge time.Duration) {
        mv.UnconfirmedTxns.Refresh(mv.blockchain, checkPeriod, maxAge)
    })
}

func TestGetOldOwnedTransactions(t *testing.T) {
    mv := setupMasterVisor()
    up := mv.UnconfirmedTxns

    // Add a transaction that is not ours, both new and old
    notOursNew, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, notOursNew, nil, true))
    notOursOld, err := makeValidTxn(mv)
    assert.Nil(t, err)
    assert.Nil(t, up.RecordTxn(mv.blockchain, notOursOld, nil, false))
    // Add a transaction that is our spend, both new and old
    ourSpendNew, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs := make(map[coin.Address]byte, 1)
    ux, ok := mv.blockchain.Unspent.Get(ourSpendNew.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourSpendNew, addrs, true))
    ourSpendOld, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 1)
    ux, ok = mv.blockchain.Unspent.Get(ourSpendNew.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourSpendOld, addrs, false))
    // Add a transaction that is our receive, both new and old
    ourReceiveNew, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 1)
    addrs[ourReceiveNew.Out[1].DestinationAddress] = byte(1)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourReceiveNew, addrs, true))
    ourReceiveOld, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 1)
    addrs[ourReceiveOld.Out[1].DestinationAddress] = byte(1)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourReceiveOld, addrs, false))
    // Add a transaction that is both our spend and receive, both new and old
    ourBothNew, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 2)
    ux, ok = mv.blockchain.Unspent.Get(ourBothNew.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    addrs[ourBothNew.Out[1].DestinationAddress] = byte(1)
    assert.Equal(t, len(addrs), 2)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourBothNew, addrs, true))
    ourBothOld, err := makeValidTxn(mv)
    assert.Nil(t, err)
    addrs = make(map[coin.Address]byte, 1)
    ux, ok = mv.blockchain.Unspent.Get(ourBothOld.In[0].UxOut)
    assert.True(t, ok)
    addrs[ux.Body.Address] = byte(1)
    addrs[ourBothOld.Out[1].DestinationAddress] = byte(1)
    assert.Equal(t, len(addrs), 2)
    assert.Nil(t, up.RecordTxn(mv.blockchain, ourBothOld, addrs, false))

    // Get the old owned txns
    utxns := up.GetOldOwnedTransactions(time.Hour)

    // Check that the 3 txns are ones we are interested in and old enough
    assert.Equal(t, len(utxns), 3)
    mapTxns := make(map[coin.SHA256]bool)
    txns := make(coin.Transactions, len(utxns))
    for i, utx := range utxns {
        txns[i] = utx.Txn
        assert.True(t, utx.IsOurSpend || utx.IsOurReceive)
        assert.True(t, utx.Announced.IsZero())
        mapTxns[utx.Hash()] = true
    }
    assert.Equal(t, len(mapTxns), 3)
    sort.Sort(txns)
    expectTxns := coin.Transactions{ourSpendOld, ourReceiveOld, ourBothOld}
    sort.Sort(expectTxns)
    assert.Equal(t, txns, expectTxns)
}

func TestFilterKnown(t *testing.T) {
    up := NewUnconfirmedTxnPool()

    uts := make([]UnconfirmedTxn, 4)
    for i := 0; i < len(uts); i++ {
        ut := createUnconfirmedTxn()
        uts[i] = ut
        up.Txns[ut.Hash()] = ut
    }
    assert.Equal(t, len(up.Txns), 4)

    hashes := []coin.SHA256{
        uts[0].Hash(),
        uts[1].Hash(),
        randSHA256(),
        randSHA256(),
    }

    known := up.FilterKnown(hashes)
    assert.Equal(t, len(known), 2)
    for i, h := range known {
        assert.Equal(t, h, hashes[i+2])
    }
    _, ok := up.Txns[known[0]]
    assert.False(t, ok)
    _, ok = up.Txns[known[1]]
    assert.False(t, ok)
}

func TestGetKnown(t *testing.T) {
    up := NewUnconfirmedTxnPool()

    uts := make([]UnconfirmedTxn, 4)
    for i := 0; i < len(uts); i++ {
        ut := createUnconfirmedTxn()
        uts[i] = ut
        up.Txns[ut.Hash()] = ut
    }
    assert.Equal(t, len(up.Txns), 4)

    hashes := []coin.SHA256{
        uts[0].Hash(),
        uts[1].Hash(),
        randSHA256(),
        randSHA256(),
    }

    known := up.GetKnown(hashes)
    assert.Equal(t, len(known), 2)
    for i, tx := range known {
        assert.Equal(t, tx.Hash(), hashes[i])
        assert.Equal(t, tx, uts[i].Txn)
    }
    _, ok := up.Txns[known[0].Hash()]
    assert.True(t, ok)
    _, ok = up.Txns[known[1].Hash()]
    assert.True(t, ok)
}
