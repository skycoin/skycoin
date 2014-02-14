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
    tx.PushOutput(makeAddress(), 100, 50)
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
    txo.Header = tx.Header
    txo.Header.Sigs = make([]Sig, len(tx.Header.Sigs))
    copy(txo.Header.Sigs, tx.Header.Sigs)
    txo.In = make([]TransactionInput, len(tx.In))
    copy(txo.In, tx.In)
    txo.Out = make([]TransactionOutput, len(tx.Out))
    copy(txo.Out, tx.Out)
    return txo
}

func TestTransactionVerify(t *testing.T) {
    // Mismatch header hash
    tx := makeTransaction(t)
    tx.Header.Hash = SHA256{}
    err := tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Invalid header hash")

    // No inputs
    tx = makeTransaction(t)
    tx.In = make([]TransactionInput, 0)
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No inputs")

    // No outputs
    tx = makeTransaction(t)
    tx.Out = make([]TransactionOutput, 0)
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No outputs")

    // Invalid number of sigs
    tx = makeTransaction(t)
    tx.Header.Sigs = make([]Sig, 0)
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Invalid number of signatures")
    tx.Header.Sigs = make([]Sig, 20)
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Invalid number of signatures")

    // Too many sigs & inputs
    tx = makeTransaction(t)
    tx.Header.Sigs = make([]Sig, math.MaxUint16)
    tx.In = make([]TransactionInput, math.MaxUint16)
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Too many signatures and inputs")

    // Duplicate inputs
    tx, s := makeTransactionWithSecret(t)
    tx.PushInput(tx.In[0].UxOut)
    tx.Header.Sigs = nil
    tx.SignInputs([]SecKey{s, s})
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Duplicate spend")

    // Duplicate outputs
    tx = makeTransaction(t)
    to := tx.Out[0]
    tx.PushOutput(to.DestinationAddress, to.Coins, to.Hours)
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Duplicate output in transaction")

    // Valid signatures
    tx = makeTransaction(t)
    tx.Header.Sigs[0] = Sig{}
    err = tx.Verify()
    assert.NotNil(t, err)
    _, s = GenerateKeyPair()
    tx.Header.Sigs[0] = SignHash(tx.hashInner(), s)
    err = tx.Verify()
    assert.NotNil(t, err)
    tx = makeTransaction(t)
    tx.Out[0].Coins += 10
    tx.UpdateHeader()
    err = tx.Verify()
    assert.NotNil(t, err)

    // Output coins are not multiples of 1e6
    tx = makeTransaction(t)
    assert.NotEqual(t, tx.Out[0].Coins%1e6, 0)
    err = tx.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Transaction outputs must be multiple of "+
        "1e6 base units")

    // Valid
    tx = makeTransaction(t)
    tx.Out[0].Coins = 10e6
    tx.UpdateHeader()
    assert.Nil(t, tx.Verify())
}

func TestTransactionPushInput(t *testing.T) {
    tx := &Transaction{}
    ux := makeUxOut(t)
    assert.Equal(t, tx.PushInput(ux.Hash()), uint16(0))
    assert.Equal(t, len(tx.In), 1)
    assert.Equal(t, tx.In[0].UxOut, ux.Hash())
    tx.In = append(tx.In, make([]TransactionInput, math.MaxUint16)...)
    ux = makeUxOut(t)
    assert.Panics(t, func() { tx.PushInput(ux.Hash()) })
}

func TestTransactionPushOutput(t *testing.T) {
    tx := &Transaction{}
    a := makeAddress()
    tx.PushOutput(a, 100, 150)
    assert.Equal(t, len(tx.Out), 1)
    assert.Equal(t, tx.Out[0], TransactionOutput{
        DestinationAddress: a,
        Coins:              100,
        Hours:              150,
    })
    for i := 1; i < 20; i++ {
        a := makeAddress()
        tx.PushOutput(a, uint64(i*100), uint64(i*50))
        assert.Equal(t, len(tx.Out), i+1)
        assert.Equal(t, tx.Out[i], TransactionOutput{
            DestinationAddress: a,
            Coins:              uint64(i * 100),
            Hours:              uint64(i * 50),
        })
    }
}

func TestTransactionSignInput(t *testing.T) {
    tx := &Transaction{}
    // Panics if too many inputs exist
    tx.In = append(tx.In, make([]TransactionInput, math.MaxUint16+2)...)
    _, s := GenerateKeyPair()
    assert.Panics(t, func() { tx.signInput(0, s, SHA256{}) })

    // Panics if idx too large for number of inputs
    tx = &Transaction{}
    ux, s := makeUxOutWithSecret(t)
    tx.PushInput(ux.Hash())
    assert.Panics(t, func() { tx.signInput(1, s, SHA256{}) })

    // Sigs should be extended if needed
    assert.Equal(t, len(tx.Header.Sigs), 0)
    ux2, s2 := makeUxOutWithSecret(t)
    tx.PushInput(ux2.Hash())
    tx.signInput(1, s2, tx.hashInner())
    assert.Equal(t, len(tx.Header.Sigs), 2)
    assert.Equal(t, tx.Header.Sigs[0], Sig{})
    assert.NotEqual(t, tx.Header.Sigs[1], Sig{})
    // Signing the earlier sig should be ok
    tx.signInput(0, s, tx.hashInner())
    assert.Equal(t, len(tx.Header.Sigs), 2)
    assert.NotEqual(t, tx.Header.Sigs[0], Sig{})
    assert.NotEqual(t, tx.Header.Sigs[1], Sig{})
}

func TestTransactionSignInputs(t *testing.T) {
    tx := &Transaction{}
    // Panics if txns already signed
    tx.Header.Sigs = append(tx.Header.Sigs, Sig{})
    assert.Panics(t, func() { tx.SignInputs([]SecKey{}) })
    // Panics if not enough keys
    tx = &Transaction{}
    ux, s := makeUxOutWithSecret(t)
    tx.PushInput(ux.Hash())
    ux2, s2 := makeUxOutWithSecret(t)
    tx.PushInput(ux2.Hash())
    tx.PushOutput(makeAddress(), 40, 80)
    assert.Equal(t, len(tx.Header.Sigs), 0)
    assert.Panics(t, func() { tx.SignInputs([]SecKey{s}) })
    assert.Equal(t, len(tx.Header.Sigs), 0)
    // Valid signing
    h := tx.hashInner()
    assert.NotPanics(t, func() { tx.SignInputs([]SecKey{s, s2}) })
    assert.Equal(t, len(tx.Header.Sigs), 2)
    assert.Equal(t, tx.hashInner(), h)
    p := PubKeyFromSecKey(s)
    a := AddressFromPubKey(p)
    p = PubKeyFromSecKey(s2)
    a2 := AddressFromPubKey(p)
    assert.Nil(t, ChkSig(a, h, tx.Header.Sigs[0]))
    assert.Nil(t, ChkSig(a2, h, tx.Header.Sigs[1]))
    assert.NotNil(t, ChkSig(a, h, tx.Header.Sigs[1]))
    assert.NotNil(t, ChkSig(a2, h, tx.Header.Sigs[0]))
}

func TestTransactionHash(t *testing.T) {
    tx := makeTransaction(t)
    assert.NotEqual(t, tx.Hash(), SHA256{})
    assert.NotEqual(t, tx.hashInner(), tx.Hash())
}

func TestTransactionUpdateHeader(t *testing.T) {
    tx := makeTransaction(t)
    h := tx.Header.Hash
    tx.Header.Hash = SHA256{}
    tx.UpdateHeader()
    assert.NotEqual(t, tx.Header.Hash, SHA256{})
    assert.Equal(t, tx.Header.Hash, h)
    assert.Equal(t, tx.Header.Hash, tx.hashInner())
}

func TestTransactionHashInner(t *testing.T) {
    tx := makeTransaction(t)

    h := tx.hashInner()
    assert.NotEqual(t, h, SHA256{})

    // If tx.In is changed, hash should change
    tx2 := copyTransaction(tx)
    ux := makeUxOut(t)
    tx2.In[0].UxOut = ux.Hash()
    assert.NotEqual(t, tx, tx2)
    assert.Equal(t, tx2.In[0].UxOut, ux.Hash())
    assert.NotEqual(t, tx.hashInner(), tx2.hashInner())

    // If tx.Out is changed, hash should change
    tx2 = copyTransaction(tx)
    a := makeAddress()
    tx2.Out[0].DestinationAddress = a
    assert.NotEqual(t, tx, tx2)
    assert.Equal(t, tx2.Out[0].DestinationAddress, a)
    assert.NotEqual(t, tx.hashInner(), tx2.hashInner())

    // If tx.Header is changed, hash should not change
    tx2 = copyTransaction(tx)
    tx.Header.Sigs = append(tx.Header.Sigs, Sig{})
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
