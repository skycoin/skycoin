package visor

import (
    "crypto/rand"
    "encoding/json"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/stretchr/testify/assert"
    "os"
    "reflect"
    "testing"
    "time"
)

const (
    testMasterKeysFile = "testmaster.keys"
    testWalletFile     = "testwallet.json"
    testBlocksigsFile  = "testblockchain.sigs"
    testBlockchainFile = "testblockchain.bin"
)

// Returns an appropriate VisorConfig and a master visor
func setupVisorConfig() (VisorConfig, *Visor) {
    coin.SetAddressVersion("test")

    // Make a new master visor + blockchain
    // Get the signed genesis block,
    mw := NewWalletEntry()
    mvc := NewVisorConfig()
    mvc.IsMaster = true
    mvc.MasterKeys = mw
    mv := NewVisor(mvc)
    sb := mv.GetGenesisBlock()

    // Use the master values for a client configuration
    c := NewVisorConfig()
    c.IsMaster = false
    c.MasterKeys = mw
    c.MasterKeys.Secret = coin.SecKey{}
    c.GenesisTimestamp = sb.Block.Header.Time
    c.GenesisSignature = sb.Sig
    return c, mv
}

func setupVisor() (v *Visor, mv *Visor) {
    vc, mv := setupVisorConfig()
    v = NewVisor(vc)
    return
}

func setupMasterVisorConfig() VisorConfig {
    // Create testmaster.keys file
    coin.SetAddressVersion("test")
    c := NewVisorConfig()
    c.IsMaster = true
    c.MasterKeys = NewWalletEntry()
    return c
}

func setupMasterVisor() *Visor {
    return NewVisor(setupMasterVisorConfig())
}

func cleanupVisor() {
    os.Remove(testMasterKeysFile)
    os.Remove(testBlockchainFile)
    os.Remove(testBlocksigsFile)
    os.Remove(testWalletFile)
}

func createUnconfirmedTxn() UnconfirmedTxn {
    ut := UnconfirmedTxn{}
    ut.Txn = coin.Transaction{}
    b := make([]byte, 32)
    rand.Read(b)
    ut.Txn.Header.Hash = coin.SumSHA256(b)
    ut.Received = time.Now().UTC()
    ut.Checked = ut.Received
    ut.Announced = ut.Received
    ut.IsOurSpend = true
    ut.IsOurReceive = true
    return ut
}

func addUnconfirmedTxn(v *Visor) UnconfirmedTxn {
    ut := createUnconfirmedTxn()
    v.UnconfirmedTxns.Txns[ut.Txn.Header.Hash] = ut
    return ut
}

func transferCoins(mv *Visor, v *Visor) error {
    // Give the nonmaster some money to spend
    addr := v.Wallet.Entries[0].Address
    tx, err := mv.Spend(Balance{10 * 1e6, 0}, 0, addr)
    if err != nil {
        return err
    }
    mv.RecordTxn(tx, false)
    sb, err := mv.CreateBlock()
    if err != nil {
        return err
    }
    return v.ExecuteSignedBlock(sb)
}

func assertJSONSerializability(t *testing.T, thing interface{}) {
    b, err := json.Marshal(thing)
    assert.Nil(t, err)
    rt := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(thing)).Interface())
    newThing := reflect.New(rt).Interface()
    err = json.Unmarshal(b, newThing)
    assert.Nil(t, err)
    assert.True(t, reflect.DeepEqual(thing, newThing))
}

func TestNewBlockchainMetadata(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    addUnconfirmedTxn(v)
    addUnconfirmedTxn(v)

    bcm := NewBlockchainMetadata(v)
    assert.Equal(t, bcm.Unspents, uint64(2))
    assert.Equal(t, bcm.Unconfirmed, uint64(2))
    assertReadableBlockHeader(t, bcm.Head, v.blockchain.Head().Header)
    assertJSONSerializability(t, &bcm)
}

func TestNewTransactionStatus(t *testing.T) {
    ts := NewUnconfirmedTransactionStatus()
    assert.True(t, ts.Unconfirmed)
    assert.False(t, ts.Unknown)
    assert.False(t, ts.Confirmed)
    assert.Equal(t, ts.Height, uint64(0))
    assertJSONSerializability(t, &ts)

    ts = NewUnknownTransactionStatus()
    assert.False(t, ts.Unconfirmed)
    assert.True(t, ts.Unknown)
    assert.False(t, ts.Confirmed)
    assert.Equal(t, ts.Height, uint64(0))
    assertJSONSerializability(t, &ts)

    ts = NewConfirmedTransactionStatus(uint64(7))
    assert.False(t, ts.Unconfirmed)
    assert.False(t, ts.Unknown)
    assert.True(t, ts.Confirmed)
    assert.Equal(t, ts.Height, uint64(7))
    assertJSONSerializability(t, &ts)

    assert.Panics(t, func() { NewConfirmedTransactionStatus(uint64(0)) })
}

func assertReadableTransactionHeader(t *testing.T,
    rth ReadableTransactionHeader, th coin.TransactionHeader) {
    assert.Equal(t, len(rth.Sigs), len(th.Sigs))
    assert.NotPanics(t, func() {
        for i, s := range rth.Sigs {
            assert.Equal(t, coin.MustSigFromHex(s), th.Sigs[i])
        }
        assert.Equal(t, coin.MustSHA256FromHex(rth.Hash), th.Hash)
    })
    assertJSONSerializability(t, &rth)
}

func TestReadableTransactionHeader(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    b := mv.blockchain.Head()
    th := b.Body.Transactions[0].Header
    rth := NewReadableTransactionHeader(&th)
    assertReadableTransactionHeader(t, rth, th)
}

func assertReadableTransactionOutput(t *testing.T,
    rto ReadableTransactionOutput, to coin.TransactionOutput) {
    assert.NotPanics(t, func() {
        assert.Equal(t, coin.MustDecodeBase58Address(rto.DestinationAddress),
            to.DestinationAddress)
    })
    assert.Equal(t, rto.Coins, to.Coins)
    assert.Equal(t, rto.Hours, to.Hours)
    assertJSONSerializability(t, &rto)
}

func TestReadableTransactionOutput(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    b := mv.blockchain.Head()
    to := b.Body.Transactions[0].Out[0]

    rto := NewReadableTransactionOutput(&to)
    assertReadableTransactionOutput(t, rto, to)
}

func assertReadableTransactionInput(t *testing.T,
    rti ReadableTransactionInput, ti coin.TransactionInput) {
    assert.NotPanics(t, func() {
        assert.Equal(t, coin.MustSHA256FromHex(rti.UxOut), ti.UxOut)
    })
    assertJSONSerializability(t, &rti)
}

func TestReadableTransactionInput(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    b := mv.blockchain.Head()
    ti := b.Body.Transactions[0].In[0]
    rti := NewReadableTransactionInput(&ti)
    assertReadableTransactionInput(t, rti, ti)
}

func assertReadableTransaction(t *testing.T, rtx ReadableTransaction,
    tx coin.Transaction) {
    assert.Equal(t, len(tx.In), len(rtx.In))
    assert.Equal(t, len(tx.Out), len(rtx.Out))
    assertReadableTransactionHeader(t, rtx.Header, tx.Header)
    for i, ti := range rtx.In {
        assertReadableTransactionInput(t, ti, tx.In[i])
    }
    for i, to := range rtx.Out {
        assertReadableTransactionOutput(t, to, tx.Out[i])
    }
    assertJSONSerializability(t, &rtx)
}

func TestReadableTransaction(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    b := mv.blockchain.Head()
    tx := b.Body.Transactions[0]

    rtx := NewReadableTransaction(&tx)
    assertReadableTransaction(t, rtx, tx)
}

func assertReadableBlockHeader(t *testing.T, rb ReadableBlockHeader,
    bh coin.BlockHeader) {
    assert.Equal(t, rb.Version, bh.Version)
    assert.Equal(t, rb.Time, bh.Time)
    assert.Equal(t, rb.BkSeq, bh.BkSeq)
    assert.Equal(t, rb.Fee, bh.Fee)
    assert.NotPanics(t, func() {
        assert.Equal(t, coin.MustSHA256FromHex(rb.PrevHash), bh.PrevHash)
        assert.Equal(t, coin.MustSHA256FromHex(rb.BodyHash), bh.BodyHash)
    })
    assertJSONSerializability(t, &rb)
}

func TestNewReadableBlockHeader(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    bh := mv.blockchain.Head().Header
    assert.Equal(t, bh.BkSeq, uint64(1))
    rb := NewReadableBlockHeader(&bh)
    assertReadableBlockHeader(t, rb, bh)
}

func assertReadableBlockBody(t *testing.T, rbb ReadableBlockBody,
    bb coin.BlockBody) {
    assert.Equal(t, len(rbb.Transactions), len(bb.Transactions))
    for i, rt := range rbb.Transactions {
        assertReadableTransaction(t, rt, bb.Transactions[i])
    }
    assertJSONSerializability(t, &rbb)
}

func assertReadableBlock(t *testing.T, rb ReadableBlock, b coin.Block) {
    assertReadableBlockHeader(t, rb.Header, b.Header)
    assertReadableBlockBody(t, rb.Body, b.Body)
    assertJSONSerializability(t, &rb)
}

func TestNewReadableBlock(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    assert.Nil(t, transferCoins(mv, v))
    b := *(mv.blockchain.Head())
    assert.Equal(t, b.Header.BkSeq, uint64(1))
    rb := NewReadableBlock(&b)
    assertReadableBlock(t, rb, b)
}
