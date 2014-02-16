package coin

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "math"
    "sort"
    "testing"
)

func makeTransactionWithSecret(t *testing.T) (Transaction, SecKey) {
    tx := Transaction{}
    ux, s := makeUxOutWithSecret(t)
    tx.PushInput(ux.Hash())
    tx.SignInputs([]SecKey{s})
    tx.PushOutput(makeAddress(), 1e6, 50)
    tx.PushOutput(makeAddress(), 5e6, 50)
    tx.UpdateHeader()
    return tx, s
}

func makeTransaction(t *testing.T) Transaction {
    tx, _ := makeTransactionWithSecret(t)
    return tx
}

func makeAddress() Address {
    p, _ := GenerateKeyPair()
    return AddressFromPubKey(p)
}

func manualTransactionsIsSorted(txns Transactions) bool {
    isSorted := true
    for i := 0; i < len(txns)-1; i++ {
        hi := txns[i].Hash()
        hj := txns[i+1].Hash()
        if bytes.Compare(hi[:], hj[:]) > 0 {
            isSorted = false
            break
        }
    }
    return isSorted
}

func copyTransaction(tx Transaction) Transaction {
    txo := Transaction{}
    txo.Head = tx.Head
    txo.Head.Sigs = make([]Sig, len(tx.Head.Sigs))
    copy(txo.Head.Sigs, tx.Head.Sigs)
    txo.In = make([]SHA256, len(tx.In))
    copy(txo.In, tx.In)
    txo.Out = make([]TransactionOutput, len(tx.Out))
    copy(txo.Out, tx.Out)
    return txo
}

func TestTransactionVerify(t *testing.T) {
    // Mismatch header hash
    tx := makeTransaction(t)
    tx.Head.Hash = SHA256{}
    assertError(t, tx.Verify(), "Invalid header hash")

    // No inputs
    tx = makeTransaction(t)
    tx.In = make([]SHA256, 0)
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "No inputs")

    // No outputs
    tx = makeTransaction(t)
    tx.Out = make([]TransactionOutput, 0)
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "No outputs")

    // Invalid number of sigs
    tx = makeTransaction(t)
    tx.Head.Sigs = make([]Sig, 0)
    assertError(t, tx.Verify(), "Invalid number of signatures")
    tx.Head.Sigs = make([]Sig, 20)
    assertError(t, tx.Verify(), "Invalid number of signatures")

    // Too many sigs & inputs
    tx = makeTransaction(t)
    tx.Head.Sigs = make([]Sig, math.MaxUint16)
    tx.In = make([]SHA256, math.MaxUint16)
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "Too many signatures and inputs")

    // Duplicate inputs
    tx, s := makeTransactionWithSecret(t)
    tx.PushInput(tx.In[0])
    tx.Head.Sigs = nil
    tx.SignInputs([]SecKey{s, s})
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "Duplicate spend")

    // Duplicate outputs
    tx = makeTransaction(t)
    to := tx.Out[0]
    tx.PushOutput(to.Address, to.Coins, to.Hours)
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "Duplicate output in transaction")

    // Invalid signature, empty
    tx = makeTransaction(t)
    tx.Head.Sigs[0] = Sig{}
    assertError(t, tx.Verify(), "Failed to recover public key")
    // Invalid signature, haKify(), "Invalid transaction signature")
    // We can't check here for other invalid signatures:
    //      - Signatures signed by someone else, spending coins they don't own
    //      - Signature is for wrong hash
    // This must be done by blockchain tests, because we need the address
    // from the unspent being spent

    // Output coins are not multiples of 1e6
    tx = makeTransaction(t)
    tx.Out[0].Coins += 10
    tx.UpdateHeader()
    tx.Head.Sigs = nil
    tx.SignInputs([]SecKey{genSecret})
    assert.NotEqual(t, tx.Out[0].Coins%1e6, uint64(0))
    assertError(t, tx.Verify(), "Transaction outputs must be multiple of "+
        "1e6 base units")

    // Output coins are 0
    tx = makeTransaction(t)
    tx.Out[0].Coins = 0
    tx.UpdateHeader()
    assertError(t, tx.Verify(), "Zero coin output")

    // Valid
    tx = makeTransaction(t)
    tx.Out[0].Coins = 10e6
    tx.Out[1].Coins = 1e6
    tx.UpdateHeader()
    assert.Nil(t, tx.Verify())
}

func TestTransactionPushInput(t *testing.T) {
    tx := &Transaction{}
    ux := makeUxOut(t)
    assert.Equal(t, tx.PushInput(ux.Hash()), uint16(0))
    assert.Equal(t, len(tx.In), 1)
    assert.Equal(t, tx.In[0], ux.Hash())
    tx.In = append(tx.In, make([]SHA256, math.MaxUint16)...)
    ux = makeUxOut(t)
    assert.Panics(t, func() { tx.PushInput(ux.Hash()) })
}

func TestTransactionPushOutput(t *testing.T) {
    tx := &Transaction{}
    a := makeAddress()
    tx.PushOutput(a, 100, 150)
    assert.Equal(t, len(tx.Out), 1)
    assert.Equal(t, tx.Out[0], TransactionOutput{
        Address: a,
        Coins:   100,
        Hours:   150,
    })
    for i := 1; i < 20; i++ {
        a := makeAddress()
        tx.PushOutput(a, uint64(i*100), uint64(i*50))
        assert.Equal(t, len(tx.Out), i+1)
        assert.Equal(t, tx.Out[i], TransactionOutput{
            Address: a,
            Coins:   uint64(i * 100),
            Hours:   uint64(i * 50),
        })
    }
}

func TestTransactionSignInput(t *testing.T) {
    tx := &Transaction{}
    // Panics if too many inputs exist
    tx.In = append(tx.In, make([]SHA256, math.MaxUint16+2)...)
    _, s := GenerateKeyPair()
    assert.Panics(t, func() { tx.signInput(0, s, SHA256{}) })

    // Panics if idx too large for number of inputs
    tx = &Transaction{}
    ux, s := makeUxOutWithSecret(t)
    tx.PushInput(ux.Hash())
    assert.Panics(t, func() { tx.signInput(1, s, SHA256{}) })

    // Sigs should be extended if needed
    assert.Equal(t, len(tx.Head.Sigs), 0)
    ux2, s2 := makeUxOutWithSecret(t)
    tx.PushInput(ux2.Hash())
    tx.signInput(1, s2, tx.hashInner())
    assert.Equal(t, len(tx.Head.Sigs), 2)
    assert.Equal(t, tx.Head.Sigs[0], Sig{})
    assert.NotEqual(t, tx.Head.Sigs[1], Sig{})
    // Signing the earlier sig should be ok
    tx.signInput(0, s, tx.hashInner())
    assert.Equal(t, len(tx.Head.Sigs), 2)
    assert.NotEqual(t, tx.Head.Sigs[0], Sig{})
    assert.NotEqual(t, tx.Head.Sigs[1], Sig{})
}

func TestTransactionSignInputs(t *testing.T) {
    tx := &Transaction{}
    // Panics if txns already signed
    tx.Head.Sigs = append(tx.Head.Sigs, Sig{})
    assert.Panics(t, func() { tx.SignInputs([]SecKey{}) })
    // Panics if not enough keys
    tx = &Transaction{}
    ux, s := makeUxOutWithSecret(t)
    tx.PushInput(ux.Hash())
    ux2, s2 := makeUxOutWithSecret(t)
    tx.PushInput(ux2.Hash())
    tx.PushOutput(makeAddress(), 40, 80)
    assert.Equal(t, len(tx.Head.Sigs), 0)
    assert.Panics(t, func() { tx.SignInputs([]SecKey{s}) })
    assert.Equal(t, len(tx.Head.Sigs), 0)
    // Valid signing
    h := tx.hashInner()
    assert.NotPanics(t, func() { tx.SignInputs([]SecKey{s, s2}) })
    assert.Equal(t, len(tx.Head.Sigs), 2)
    assert.Equal(t, tx.hashInner(), h)
    p := PubKeyFromSecKey(s)
    a := AddressFromPubKey(p)
    p = PubKeyFromSecKey(s2)
    a2 := AddressFromPubKey(p)
    assert.Nil(t, ChkSig(a, h, tx.Head.Sigs[0]))
    assert.Nil(t, ChkSig(a2, h, tx.Head.Sigs[1]))
    assert.NotNil(t, ChkSig(a, h, tx.Head.Sigs[1]))
    assert.NotNil(t, ChkSig(a2, h, tx.Head.Sigs[0]))
}

func TestTransactionHash(t *testing.T) {
    tx := makeTransaction(t)
    assert.NotEqual(t, tx.Hash(), SHA256{})
    assert.NotEqual(t, tx.hashInner(), tx.Hash())
}

func TestTransactionUpdateHeader(t *testing.T) {
    tx := makeTransaction(t)
    h := tx.Head.Hash
    tx.Head.Hash = SHA256{}
    tx.UpdateHeader()
    assert.NotEqual(t, tx.Head.Hash, SHA256{})
    assert.Equal(t, tx.Head.Hash, h)
    assert.Equal(t, tx.Head.Hash, tx.hashInner())
}

func TestTransactionHashInner(t *testing.T) {
    tx := makeTransaction(t)

    h := tx.hashInner()
    assert.NotEqual(t, h, SHA256{})

    // If tx.In is changed, hash should change
    tx2 := copyTransaction(tx)
    ux := makeUxOut(t)
    tx2.In[0] = ux.Hash()
    assert.NotEqual(t, tx, tx2)
    assert.Equal(t, tx2.In[0], ux.Hash())
    assert.NotEqual(t, tx.hashInner(), tx2.hashInner())

    // If tx.Out is changed, hash should change
    tx2 = copyTransaction(tx)
    a := makeAddress()
    tx2.Out[0].Address = a
    assert.NotEqual(t, tx, tx2)
    assert.Equal(t, tx2.Out[0].Address, a)
    assert.NotEqual(t, tx.hashInner(), tx2.hashInner())

    // If tx.Head is changed, hash should not change
    tx2 = copyTransaction(tx)
    tx.Head.Sigs = append(tx.Head.Sigs, Sig{})
    assert.Equal(t, tx.hashInner(), tx2.hashInner())
}

func TestTransactionSerialization(t *testing.T) {
    tx := makeTransaction(t)
    b := tx.Serialize()
    tx2 := TransactionDeserialize(b)
    assert.Equal(t, tx, tx2)
    // Invalid deserialization
    assert.Panics(t, func() { TransactionDeserialize([]byte{0x04}) })
}

func TestTransactionSorting(t *testing.T) {
    txns := make(Transactions, 4)
    for i := 0; i < len(txns); i++ {
        txns[i] = makeTransaction(t)
    }

    // Sort(), IsSorted(), Less()
    isSorted := manualTransactionsIsSorted(txns)
    assert.Equal(t, sort.IsSorted(txns), isSorted)
    assert.Equal(t, txns.IsSorted(), isSorted)
    if isSorted {
        txns[0], txns[1] = txns[1], txns[0]
        assert.False(t, txns.Less(0, 1))
        assert.True(t, txns.Less(1, 0))
    }
    assert.False(t, manualTransactionsIsSorted(txns))
    assert.False(t, sort.IsSorted(txns))
    assert.False(t, txns.IsSorted())
    txns.Sort()
    assert.True(t, manualTransactionsIsSorted(txns))
    assert.True(t, sort.IsSorted(txns))
    assert.True(t, txns.IsSorted())
    for i := 0; i < len(txns)-1; i++ {
        assert.True(t, txns.Less(i, i+1))
        assert.False(t, txns.Less(i+1, i))
    }

    // Len()
    assert.Equal(t, len(txns), txns.Len())
    assert.Equal(t, 4, txns.Len())

    // Swap()
    tx1 := txns[0]
    tx2 := txns[1]
    txns.Swap(0, 1)
    assert.Equal(t, txns[0], tx2)
    assert.Equal(t, txns[1], tx1)
    txns.Swap(0, 1)
    assert.Equal(t, txns[0], tx1)
    assert.Equal(t, txns[1], tx2)
    txns.Swap(1, 0)
    assert.Equal(t, txns[0], tx2)
    assert.Equal(t, txns[1], tx1)
    txns.Swap(1, 0)
    assert.Equal(t, txns[0], tx1)
    assert.Equal(t, txns[1], tx2)
}

func TestTransactionsHashes(t *testing.T) {
    txns := make(Transactions, 4)
    for i := 0; i < len(txns); i++ {
        txns[i] = makeTransaction(t)
    }
    hashes := txns.Hashes()
    assert.Equal(t, len(hashes), 4)
    for i, h := range hashes {
        assert.Equal(t, h, txns[i].Hash())
    }
}

func TestFullTransaction(t *testing.T) {
    p1, s1 := GenerateKeyPair()
    a1 := AddressFromPubKey(p1)
    bc := NewBlockchain()
    bc.CreateMasterGenesisBlock(a1)
    tx := Transaction{}
    ux := bc.Unspent.Arr[0]
    tx.PushInput(ux.Hash())
    p2, s2 := GenerateKeyPair()
    a2 := AddressFromPubKey(p2)
    tx.PushOutput(a1, ux.Body.Coins-6e6, 100)
    tx.PushOutput(a2, 1e6, 100)
    tx.PushOutput(a2, 5e6, 100)
    tx.SignInputs([]SecKey{s1})
    tx.UpdateHeader()
    assert.Nil(t, tx.Verify())
    assert.Nil(t, bc.VerifyTransaction(tx))
    b, err := bc.NewBlockFromTransactions(Transactions{tx}, 10)
    assert.Nil(t, err)
    _, err = bc.ExecuteBlock(b)
    assert.Nil(t, err)

    txo := CreateExpectedUnspents(tx)
    tx = Transaction{}
    assert.Equal(t, txo[0].Body.Address, a1)
    assert.Equal(t, txo[1].Body.Address, a2)
    assert.Equal(t, txo[2].Body.Address, a2)
    ux0, ok := bc.Unspent.Get(txo[0].Hash())
    assert.True(t, ok)
    ux1, ok := bc.Unspent.Get(txo[1].Hash())
    assert.True(t, ok)
    ux2, ok := bc.Unspent.Get(txo[2].Hash())
    assert.True(t, ok)
    tx.PushInput(ux0.Hash())
    tx.PushInput(ux1.Hash())
    tx.PushInput(ux2.Hash())
    tx.PushOutput(a2, 10e6, 200)
    tx.PushOutput(a1, ux.Body.Coins-10e6, 100)
    tx.SignInputs([]SecKey{s1, s2, s2})
    tx.UpdateHeader()
    assert.Nil(t, tx.Verify())
    assert.Nil(t, bc.VerifyTransaction(tx))
    b, err = bc.NewBlockFromTransactions(Transactions{tx}, 10)
    assert.Nil(t, err)
    _, err = bc.ExecuteBlock(b)
    assert.Nil(t, err)
}
