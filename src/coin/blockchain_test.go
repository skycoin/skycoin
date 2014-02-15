package coin

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

var (
    genPublic, genSecret = GenerateKeyPair()
    genAddress           = AddressFromPubKey(genPublic)
)

/* Helpers */

func makeNewBlock() Block {
    body := BlockBody{
        Transactions: nil,
    }
    prev := Block{
        Body: body,
        Header: BlockHeader{
            Version:  0x02,
            Time:     100,
            BkSeq:    0,
            Fee:      10,
            PrevHash: SHA256{},
            BodyHash: body.Hash(),
        }}
    return newBlock(&prev, 20)
}

func makeTransactionForChain(t *testing.T, bc *Blockchain) Transaction {
    tx := Transaction{}
    ux := bc.Unspent.Arr[0]
    assert.Equal(t, ux.Body.Address, genAddress)
    tx.PushInput(bc.Unspent.Arr[0].Hash())
    tx.PushOutput(makeAddress(), 1e6, 100)
    tx.PushOutput(genAddress, ux.Body.Coins-1e6, ux.Body.Hours-100)
    tx.SignInputs([]SecKey{genSecret})
    tx.UpdateHeader()
    assert.Nil(t, tx.Verify())
    assert.Nil(t, bc.VerifyTransaction(tx))
    return tx

}

func addTransactionToBlock(t *testing.T, b *Block) Transaction {
    tx := makeTransaction(t)
    b.Body.Transactions = append(b.Body.Transactions, tx)
    return tx
}

func addBlockToBlockchain(t *testing.T, bc *Blockchain) Block {
    tx := makeTransactionForChain(t, bc)
    b, err := bc.NewBlockFromTransactions(Transactions{tx}, 10)
    assert.Nil(t, err)
    err = bc.ExecuteBlock(b)
    assert.Nil(t, err)
    return b
}

func makeMultipleOutputs(t *testing.T, bc *Blockchain) {
    txn := Transaction{}
    ux := bc.Unspent.Arr[0]
    txn.PushInput(ux.Hash())
    txn.PushOutput(genAddress, 1e6, 100)
    txn.PushOutput(genAddress, 2e6, 100)
    txn.PushOutput(genAddress, genesisCoinVolume-3e6, 100)
    txn.SignInputs([]SecKey{genSecret})
    txn.UpdateHeader()
    assert.Nil(t, txn.Verify())
    assert.Nil(t, bc.VerifyTransaction(txn))
    b, err := bc.NewBlockFromTransactions(Transactions{txn}, 100)
    assert.Nil(t, err)
    assert.Nil(t, bc.ExecuteBlock(b))
}

/* Tests */

func TestConstantsStayConstant(t *testing.T) {
    assert.Equal(t, genesisCoinVolume, uint64(100*1e6*1e6))
    assert.Equal(t, genesisCoinHours, uint64(1024*1024))
}

func TestNewBlock(t *testing.T) {
    prev := Block{Header: BlockHeader{Version: 0x02, Time: 100, BkSeq: 0}}
    b := newBlock(&prev, 33)
    assert.Equal(t, b.Body, BlockBody{})
    assert.Equal(t, b.Header.PrevHash, prev.HashHeader())
    assert.Equal(t, b.Header.Time, prev.Header.Time+33)
    assert.Equal(t, b.Header.BkSeq, uint64(1))
}

func TestBlockHashHeader(t *testing.T) {
    b := makeNewBlock()
    assert.Equal(t, b.HashHeader(), b.Header.Hash())
    assert.NotEqual(t, b.HashHeader(), SHA256{})
}

func TestBlockHashBody(t *testing.T) {
    b := makeNewBlock()
    assert.Equal(t, b.HashBody(), SHA256{})
    assert.Equal(t, b.HashBody(), b.Body.Hash())
    tx := addTransactionToBlock(t, &b)
    assert.NotEqual(t, b.HashBody(), SHA256{})
    assert.Equal(t, b.HashBody(), Merkle([]SHA256{tx.Hash()}))
    assert.Equal(t, b.HashBody(), b.Body.Hash())
}

func TestBlockUpdateHeader(t *testing.T) {
    b := makeNewBlock()
    tx := addTransactionToBlock(t, &b)
    b.Header.BodyHash = SHA256{}
    b.UpdateHeader()
    assert.NotEqual(t, b.Header.BodyHash, SHA256{})
    assert.Equal(t, b.Header.BodyHash, Merkle([]SHA256{tx.Hash()}))
    // Changing txns should change hash
    h := b.Header.BodyHash
    addTransactionToBlock(t, &b)
    b.UpdateHeader()
    assert.NotEqual(t, b.Header.BodyHash, h)
    assert.NotEqual(t, b.Header.BodyHash, SHA256{})
}

func TestBlockString(t *testing.T) {
    b := makeNewBlock()
    assert.Equal(t, b.String(), b.Header.String())
}

func TestBlockGetTransaction(t *testing.T) {
    b := makeNewBlock()
    _, ok := b.GetTransaction(SHA256{})
    assert.False(t, ok)
    tx := addTransactionToBlock(t, &b)
    tx2, ok := b.GetTransaction(tx.Hash())
    assert.True(t, ok)
    assert.Equal(t, tx, tx2)
    tx3 := addTransactionToBlock(t, &b)
    tx4, ok := b.GetTransaction(tx3.Hash())
    assert.True(t, ok)
    assert.Equal(t, tx3, tx4)
    _, ok = b.GetTransaction(SHA256{})
    assert.False(t, ok)
}

func TestNewBlockHeader(t *testing.T) {
    b := makeNewBlock()
    prev := &b.Header
    bh := newBlockHeader(prev, 22)
    assert.Equal(t, bh.PrevHash, prev.Hash())
    assert.NotEqual(t, bh.PrevHash, prev.PrevHash)
    assert.Equal(t, bh.Time, prev.Time+22)
    assert.Equal(t, bh.BkSeq, prev.BkSeq+1)
}

func TestBlockHeaderHash(t *testing.T) {
    b := makeNewBlock()
    h := b.Header.Hash()
    assert.Equal(t, h, b.Header.Hash())
    assert.NotEqual(t, b.Header.Hash(), SHA256{})
    // Change header should change hash
    b.Header.BkSeq = uint64(5)
    assert.NotEqual(t, h, b.Header.Hash())
}

func TestBlockHeaderBytes(t *testing.T) {
    b := makeNewBlock()
    by := b.Header.Bytes()
    assert.NotNil(t, by)
    assert.NotEqual(t, len(by), 0)
    b.Header.BkSeq += 1
    by2 := b.Header.Bytes()
    assert.False(t, bytes.Equal(by, by2))
}

func TestBlockHeaderString(t *testing.T) {
    b := makeNewBlock()
    assert.NotEqual(t, b.Header.String(), "")
}

func TestBlockBodyHash(t *testing.T) {
    b := makeNewBlock()
    assert.Equal(t, b.Body.Hash(), Merkle([]SHA256{}))
    tx1 := addTransactionToBlock(t, &b)
    assert.Equal(t, b.Body.Hash(), Merkle([]SHA256{tx1.Hash()}))
    tx2 := addTransactionToBlock(t, &b)
    assert.Equal(t, b.Body.Hash(), Merkle([]SHA256{tx1.Hash(), tx2.Hash()}))
}

func TestBlockBodyBytes(t *testing.T) {
    b := makeNewBlock()
    assert.NotPanics(t, func() { b.Body.Bytes() })
    by := b.Body.Bytes()
    assert.NotNil(t, by)
    addTransactionToBlock(t, &b)
    by2 := b.Body.Bytes()
    assert.True(t, len(by2) > len(by))
}

func TestNewBlockchain(t *testing.T) {
    b := NewBlockchain()
    assert.NotNil(t, b.Blocks)
    assert.Equal(t, len(b.Blocks), 0)
    assert.NotNil(t, b.Unspent)
    assert.Equal(t, len(b.Unspent.Arr), 0)
}

func TestCreateMasterGenesisBlock(t *testing.T) {
    b := NewBlockchain()
    a := makeAddress()
    gb := b.CreateMasterGenesisBlock(a)
    assert.Equal(t, len(b.Blocks), 1)
    assert.Equal(t, b.Blocks[0], gb)
    assert.Equal(t, len(b.Unspent.Arr), 1)
    assert.Equal(t, b.Unspent.Arr[0].Body.Address, a)
    assert.NotEqual(t, gb.Header.Time, uint64(0))
    assert.Equal(t, gb.Header.BkSeq, uint64(0))
    // Panicing
    assert.Panics(t, func() { b.CreateMasterGenesisBlock(a) })
}

func TestCreateGenesisBlock(t *testing.T) {
    b := NewBlockchain()
    now := Now()
    gb := b.CreateGenesisBlock(genAddress, now)
    assert.Equal(t, gb.Header.Time, now)
    assert.Equal(t, gb.Header.BkSeq, uint64(0))
    ux := b.Unspent.Arr[0]
    assert.Equal(t, ux.Head.BkSeq, uint64(0))
    assert.Equal(t, ux.Head.Time, now)
    assert.Equal(t, ux.Body.SrcTransaction, SHA256{})
    assert.Equal(t, ux.Body.Address, genAddress)
    assert.Equal(t, ux.Body.Coins, genesisCoinVolume)
    assert.Equal(t, ux.Body.Hours, genesisCoinHours)
    assert.Equal(t, gb.Header.BodyHash, SHA256{})
    assert.Equal(t, gb.Header.PrevHash, SHA256{})
    // Panicing
    assert.Panics(t, func() { b.CreateGenesisBlock(genAddress, now) })
}

func TestBlockchainHead(t *testing.T) {
    b := NewBlockchain()
    gb := b.CreateMasterGenesisBlock(genAddress)
    assert.Equal(t, *(b.Head()), gb)
    nb := addBlockToBlockchain(t, b)
    assert.Equal(t, *(b.Head()), nb)
}

func TestBlockchainTime(t *testing.T) {
    b := NewBlockchain()
    gb := b.CreateMasterGenesisBlock(genAddress)
    assert.Equal(t, b.Time(), gb.Header.Time)
    nb := addBlockToBlockchain(t, b)
    assert.Equal(t, b.Time(), nb.Header.Time)
}

func TestNewBlockFromTransactions(t *testing.T) {
    bc := NewBlockchain()
    gb := bc.CreateMasterGenesisBlock(genAddress)
    _, err := bc.NewBlockFromTransactions(Transactions{}, 22)
    assert.Error(t, err, "No valid transactions")
    txn := makeTransactionForChain(t, bc)
    txns := Transactions{txn}
    assert.Panics(t, func() { bc.NewBlockFromTransactions(txns, 0) })
    b, err := bc.NewBlockFromTransactions(txns, 22)
    assert.Nil(t, err)
    assert.Equal(t, len(b.Body.Transactions), 1)
    assert.Equal(t, b.Body.Transactions[0], txn)
    assert.Equal(t, b.Header.BkSeq, uint64(1))
    assert.Equal(t, b.Header.Time, gb.Header.Time+22)
    txn.Header.Hash = SHA256{}
    txns = Transactions{txn}
    _, err = bc.NewBlockFromTransactions(txns, 22)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No valid transactions")
}

func TestTxUxIn(t *testing.T) {
    bc := NewBlockchain()
    bc.CreateMasterGenesisBlock(genAddress)
    txn := makeTransactionForChain(t, bc)
    txin, err := bc.txUxIn(txn)
    assert.Nil(t, err)
    assert.Equal(t, len(txin), 1)
    assert.Equal(t, len(txin), len(txn.In))
    assert.Equal(t, txin[0], bc.Unspent.Arr[0])

    // Empty txn
    txn = Transaction{}
    txin, err = bc.txUxIn(txn)
    assert.Nil(t, err)
    assert.Equal(t, len(txin), 0)

    // Spending unknown output
    txn = makeTransactionForChain(t, bc)
    txn.In[0].UxOut = SHA256{}
    _, err = bc.txUxIn(txn)
    assert.Error(t, err, "Unspent output does not exist")

    // Multiple inputs
    makeMultipleOutputs(t, bc)
    txn = Transaction{}
    ux0 := bc.Unspent.Arr[0]
    ux1 := bc.Unspent.Arr[1]
    txn.PushInput(ux0.Hash())
    txn.PushInput(ux1.Hash())
    txn.PushOutput(genAddress, ux0.Body.Coins+ux1.Body.Coins, ux0.Body.Hours)
    txn.SignInputs([]SecKey{genSecret, genSecret})
    txn.UpdateHeader()
    assert.Nil(t, txn.Verify())
    assert.Nil(t, bc.VerifyTransaction(txn))
    txin, err = bc.txUxIn(txn)
    assert.Nil(t, err)
    assert.Equal(t, len(txin), 2)
    assert.Equal(t, txin[0], ux0)
    assert.Equal(t, txin[1], ux1)
}

func TestNow(t *testing.T) {
    now := Now()
    now2 := uint64(time.Now().UTC().Unix())
    assert.True(t, now == now2 || now2-1 == now)
}
