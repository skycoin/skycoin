package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "github.com/stretchr/testify/assert"
    "os"
    "sort"
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
func setupVisor() (VisorConfig, *visor.Visor) {
    coin.SetAddressVersion("test")

    // Make a new master visor + blockchain
    // Get the signed genesis block,
    mw := visor.NewWalletEntry()
    mvc := visor.NewVisorConfig()
    mvc.IsMaster = true
    mvc.MasterKeys = mw
    mv := visor.NewVisor(mvc)
    sb := mv.GetGenesisBlock()

    // Use the master values for a client configuration
    c := NewVisorConfig()
    c.Config.IsMaster = false
    c.Config.MasterKeys = mw
    c.Config.MasterKeys.Secret = coin.SecKey{}
    c.Config.GenesisTimestamp = sb.Block.Header.Time
    c.Config.GenesisSignature = sb.Sig
    return c, mv
}

func setupMasterVisor() VisorConfig {
    cleanupVisor()
    coin.SetAddressVersion("test")
    c := NewVisorConfig()
    c.Config.IsMaster = true
    c.Config.MasterKeys = visor.NewWalletEntry()
    return c
}

func cleanupVisor() {
    os.Remove(testMasterKeysFile)
    os.Remove(testBlockchainFile)
    os.Remove(testBlocksigsFile)
    os.Remove(testWalletFile)
}

// Returns a daemon with the visor enabled, but networking disabled
func newVisorDaemon(vc VisorConfig) (*Daemon, chan int) {
    quit := make(chan int)
    c := NewConfig()
    c.Daemon.DisableNetworking = true
    c.Visor = vc
    d := NewDaemon(c)
    return d, quit
}

// Writes a wallet entry to disk at filename
func writeMasterKeysFile() (visor.WalletEntry, error) {
    we := visor.NewWalletEntry()
    rwe := visor.NewReadableWalletEntry(&we)
    err := rwe.Save(testMasterKeysFile)
    return we, err
}

func assertFileExists(t *testing.T, filename string) {
    stat, err := os.Stat(filename)
    assert.Nil(t, err)
    assert.True(t, stat.Mode().IsRegular())
}

func assertFileNotExists(t *testing.T, filename string) {
    _, err := os.Stat(filename)
    assert.NotNil(t, err)
    assert.True(t, os.IsNotExist(err))
}

func createUnconfirmedTxn() visor.UnconfirmedTxn {
    ut := visor.UnconfirmedTxn{}
    ut.Txn = coin.Transaction{}
    ut.Txn.Header.Hash = coin.SumSHA256([]byte("cascas"))
    ut.Received = util.Now()
    ut.Checked = ut.Received
    ut.Announced = ut.Received
    ut.IsOurSpend = true
    ut.IsOurReceive = true
    return ut
}

func addUnconfirmedTxn(v *Visor) visor.UnconfirmedTxn {
    ut := createUnconfirmedTxn()
    v.Visor.UnconfirmedTxns.Txns[ut.Txn.Hash()] = ut
    return ut
}

func setupExistingPool(p *Pool) *gnet.Connection {
    gc := gnetConnection(addr)
    p.Pool.Pool[gc.Id] = gc
    p.Pool.Addresses[gc.Addr()] = gc
    return gc
}

func setupPool() (*Pool, *gnet.Connection) {
    m := NewMessagesConfig()
    m.Register()
    p := NewPool(NewPoolConfig())
    p.Init(nil)
    return p, setupExistingPool(p)
}

func makeValidTxn(mv *visor.Visor) (coin.Transaction, error) {
    we := visor.NewWalletEntry()
    return mv.Spend(visor.Balance{10 * 1e6, 0}, 0, we.Address)
}

func transferCoins(mv *visor.Visor, v *visor.Visor) error {
    // Give the nonmaster some money to spend
    addr := v.Wallet.GetAddresses()[0]
    tx, err := mv.Spend(visor.Balance{10 * 1e6, 0}, 0, addr)
    if err != nil {
        return err
    }
    mv.RecordTxn(tx, false)
    sb, err := mv.CreateAndExecuteBlock()
    if err != nil {
        return err
    }
    return v.ExecuteSignedBlock(sb)
}

func makeBlocks(mv *visor.Visor, n int) ([]visor.SignedBlock, error) {
    dest := visor.NewWalletEntry()
    blocks := make([]visor.SignedBlock, 0, n)
    for i := 0; i < n; i++ {
        tx, err := mv.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address)
        if err != nil {
            return nil, err
        }
        mv.RecordTxn(tx, false)
        sb, err := mv.CreateAndExecuteBlock()
        if err != nil {
            return nil, err
        }
        blocks = append(blocks, sb)
    }
    return blocks, nil
}

/* Tests for daemon's loop related to visor */

func testBlockCreationTicker(t *testing.T, vcfg VisorConfig, master bool,
    mv *visor.Visor, published bool) {
    vcfg.Config.BlockCreationInterval = 1
    defer gnet.EraseMessages()
    c := NewConfig()
    c.Visor = vcfg
    c.Daemon.DisableNetworking = true
    d := NewDaemon(c)
    if !master {
        err := transferCoins(mv, d.Visor.Visor)
        assert.Nil(t, err)
    }
    quit := make(chan int)
    defer closeDaemon(d, quit)
    gc := setupExistingPool(d.Pool)
    if master && !published {
        gc.Conn = NewFailingConn(addr)
    }
    assert.True(t, gc.LastSent.IsZero())
    assert.Equal(t, len(d.Pool.Pool.Pool), 1)
    assert.Equal(t, len(d.Pool.Pool.Addresses), 1)
    start := 0
    if !master {
        start = 1
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(start))
    go d.Start(quit)
    time.Sleep(time.Millisecond * 1250)
    // Creation should not have occured, because no transaction
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(start))

    // Creation should occur with a transaction, if not a master
    assert.True(t, gc.LastSent.IsZero())
    dest := visor.NewWalletEntry()
    _, err := d.Visor.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address,
        d.Pool)
    assert.Nil(t, err)
    if master && !published {
        assert.True(t, gc.LastSent.IsZero())
    } else {
        assert.False(t, gc.LastSent.IsZero())
    }
    time.Sleep(time.Millisecond * 1250)
    final := start
    if master {
        final += 1
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(final))
    if !published {
        gc.LastSent = util.ZeroTime()
    }
    if published {
        assert.False(t, gc.LastSent.IsZero())
    } else {
        assert.True(t, gc.LastSent.IsZero())
    }
}

func TestBlockCreationTicker(t *testing.T) {
    defer cleanupVisor()
    vcfg, mv := setupVisor()
    // No blocks should get created if we are not master
    testBlockCreationTicker(t, vcfg, false, mv, false)
}

func TestBlockCreationTickerMaster(t *testing.T) {
    defer cleanupVisor()
    vcfg := setupMasterVisor()
    // Master should make a block
    testBlockCreationTicker(t, vcfg, true, nil, true)
}

func TestBlockCreationTickerMasterUnpublished(t *testing.T) {
    defer cleanupVisor()
    vcfg := setupMasterVisor()
    // Master should make a block, but fail to send to anyone
    testBlockCreationTicker(t, vcfg, true, nil, false)
}

func TestUnconfirmedRefreshTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.Config.UnconfirmedRefreshRate = time.Millisecond * 10
    vc.Config.UnconfirmedCheckInterval = time.Nanosecond
    vc.Config.UnconfirmedMaxAge = time.Nanosecond
    d, quit := newVisorDaemon(vc)
    addUnconfirmedTxn(d.Visor)
    time.Sleep(time.Millisecond)
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.Equal(t, len(d.Visor.Visor.UnconfirmedTxns.Txns), 0)
}

func TestBlocksRequestTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.BlocksRequestRate = time.Millisecond * 10
    d, quit := newVisorDaemon(vc)
    gc := gnetConnection(addr)
    d.Pool.Pool.Pool[gc.Id] = gc
    d.Pool.Pool.Addresses[gc.Addr()] = gc
    assert.True(t, gc.LastSent.IsZero())
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.False(t, gc.LastSent.IsZero())
}

func TestBlocksAnnounceTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.BlocksAnnounceRate = time.Millisecond * 10
    d, quit := newVisorDaemon(vc)
    gc := gnetConnection(addr)
    d.Pool.Pool.Pool[gc.Id] = gc
    d.Pool.Pool.Addresses[gc.Addr()] = gc
    assert.True(t, gc.LastSent.IsZero())
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.False(t, gc.LastSent.IsZero())
}

func TestTransactionRebroadcastTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.TransactionRebroadcastRate = time.Millisecond * 10
    d, quit := newVisorDaemon(vc)
    gc := gnetConnection(addr)
    d.Pool.Pool.Pool[gc.Id] = gc
    d.Pool.Pool.Addresses[gc.Addr()] = gc
    addUnconfirmedTxn(d.Visor)
    assert.True(t, gc.LastSent.IsZero())
    time.Sleep(time.Millisecond * 15)
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.False(t, gc.LastSent.IsZero())
}

/* Tests for daemon.Visor */

func TestVisorConfigLoadMasterKeys(t *testing.T) {
    defer cleanupVisor()
    c := NewVisorConfig()
    c.Disabled = true
    mk := c.Config.MasterKeys
    // Shouldn't panic, since its disabled. keys should not be loaded
    assert.NotPanics(t, c.LoadMasterKeys)
    assert.Equal(t, c.Config.MasterKeys, mk)
    c.Disabled = false
    c.MasterKeysFile = testMasterKeysFile
    // Should panic, since keys not found
    assert.Panics(t, c.LoadMasterKeys)

    // Shouldn't panic, and keys should be correct
    we, err := writeMasterKeysFile()
    assert.Nil(t, err)
    assert.NotPanics(t, c.LoadMasterKeys)
    assert.Equal(t, c.Config.MasterKeys, we)
}

func TestNewVisor(t *testing.T) {
    defer cleanupVisor()
    c, _ := setupVisor()
    v := NewVisor(c)
    assert.Equal(t, v.Config, c)
    assert.NotNil(t, v.Visor)
    assert.Equal(t, len(v.blockchainLengths), 0)

    c.Disabled = true
    v = NewVisor(c)
    assert.Equal(t, v.Config, c)
    assert.Nil(t, v.Visor)
    assert.Equal(t, len(v.blockchainLengths), 0)
}

func TestVisorRemoveConnection(t *testing.T) {
    defer cleanupVisor()
    c := NewVisorConfig()
    c.Disabled = true
    v := NewVisor(c)
    assert.NotNil(t, v.blockchainLengths)
    v.blockchainLengths[addr] = 2
    assert.Equal(t, v.blockchainLengths[addr], uint64(2))
    assert.Equal(t, len(v.blockchainLengths), 1)
    v.RemoveConnection(addr)
    assert.Equal(t, v.blockchainLengths[addr], uint64(0))
    assert.Equal(t, len(v.blockchainLengths), 0)
}

func TestVisorShutdown(t *testing.T) {
    defer cleanupVisor()
    c, _ := setupVisor()
    c.Disabled = true
    c.Config.BlockchainFile = testBlockchainFile
    c.Config.BlockSigsFile = testBlocksigsFile
    c.Config.WalletFile = testWalletFile
    v := NewVisor(c)
    assert.NotPanics(t, v.Shutdown)
    // Should not save anything
    assertFileNotExists(t, testBlockchainFile)
    assertFileNotExists(t, testBlocksigsFile)
    assertFileNotExists(t, testWalletFile)
    cleanupVisor()

    c.Disabled = false
    v = NewVisor(c)
    assert.NotPanics(t, v.Shutdown)
    assertFileExists(t, testBlockchainFile)
    assertFileExists(t, testBlocksigsFile)
    assertFileExists(t, testWalletFile)
    cleanupVisor()

    // If master, no wallet should be saved
    c = setupMasterVisor()
    c.Config.BlockchainFile = testBlockchainFile
    c.Config.BlockSigsFile = testBlocksigsFile
    c.Config.WalletFile = testWalletFile
    v = NewVisor(c)
    assert.NotPanics(t, v.Shutdown)
    assertFileExists(t, testBlockchainFile)
    assertFileExists(t, testBlocksigsFile)
    assertFileNotExists(t, testWalletFile)
    cleanupVisor()
}

func TestVisorRefreshUnconfirmed(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.Config.UnconfirmedRefreshRate = time.Millisecond
    vc.Config.UnconfirmedCheckInterval = time.Nanosecond
    vc.Config.UnconfirmedMaxAge = time.Nanosecond
    v := NewVisor(vc)
    addUnconfirmedTxn(v)
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
    wait()
    v.Config.Disabled = true
    v.RefreshUnconfirmed()
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
    v.Config.Disabled = false
    v.RefreshUnconfirmed()
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 0)
}

func TestVisorRequestBlocks(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    assert.True(t, gc.LastSent.IsZero())

    vc.Disabled = false
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    assert.True(t, gc.LastSent.IsZero())

    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorAnnounceBlocks(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.True(t, gc.LastSent.IsZero())

    vc.Disabled = false
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.True(t, gc.LastSent.IsZero())

    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorRequestBlocksFromAddr(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.True(t, gc.LastSent.IsZero())

    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.False(t, gc.LastSent.IsZero())

    gc.LastSent = util.ZeroTime()
    gc.Conn = NewFailingConn(addr)
    assert.NotPanics(t, func() {
        assert.NotNil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.True(t, gc.LastSent.IsZero())

    gc.LastSent = util.ZeroTime()
    gc.Conn = NewDummyConn(addr)
    delete(p.Pool.Pool, gc.Id)
    delete(p.Pool.Addresses, gc.Addr())
    assert.NotPanics(t, func() {
        assert.NotNil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.True(t, gc.LastSent.IsZero())
}

func TestVisorBroadcastOurTransactions(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.True(t, gc.LastSent.IsZero())

    // With no transactions, nothing should be sent
    vc.Disabled = false
    vc.TransactionRebroadcastRate = time.Millisecond * 5
    v = NewVisor(vc)
    time.Sleep(time.Millisecond * 20)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.True(t, gc.LastSent.IsZero())

    // We have a stale owned unconfirmed txn but we failed to send it
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    tx := addUnconfirmedTxn(v)
    time.Sleep(time.Millisecond * 20)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.True(t, gc.LastSent.IsZero())
    newtx := v.Visor.UnconfirmedTxns.Txns[tx.Txn.Hash()]
    assert.Equal(t, newtx.Announced, tx.Announced)

    // We have a stale owned unconfirmed txn, should be sent, and the
    // unconfirmed txn should have announcement updated
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    tx = addUnconfirmedTxn(v)
    time.Sleep(time.Millisecond * 20)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.False(t, gc.LastSent.IsZero())
    newtx = v.Visor.UnconfirmedTxns.Txns[tx.Txn.Hash()]
    assert.True(t, newtx.Announced.After(tx.Announced))
}

func TestVisorBroadcastBlock(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() {
        // Should return error if disabled
        assert.NotNil(t, v.broadcastBlock(visor.SignedBlock{}, p))
    })
    assert.True(t, gc.LastSent.IsZero())

    // Fail to send
    gc.Conn = NewFailingConn(addr)
    vc.Disabled = false
    v = NewVisor(vc)
    sb := v.Visor.GetGenesisBlock()
    assert.NotPanics(t, func() {
        // Returns error if nobody received
        assert.NotNil(t, v.broadcastBlock(sb, p))
    })
    assert.True(t, gc.LastSent.IsZero())

    // Succeed in sending
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    sb = v.Visor.GetGenesisBlock()
    assert.NotPanics(t, func() {
        assert.Nil(t, v.broadcastBlock(sb, p))
    })
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorBroadcastTransaction(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    ut := createUnconfirmedTxn()
    assert.NotPanics(t, func() {
        // Should return error if disabled
        assert.NotNil(t, v.broadcastTransaction(ut.Txn, p))
    })
    assert.True(t, gc.LastSent.IsZero())

    // Fail to send
    gc.Conn = NewFailingConn(addr)
    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        // Returns error if nobody received
        assert.NotNil(t, v.broadcastTransaction(ut.Txn, p))
    })
    assert.True(t, gc.LastSent.IsZero())

    // Succeed in sending
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.broadcastTransaction(ut.Txn, p))
    })
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorSpend(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, mv := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    // Spending while disabled
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{10e6, 0}, 0,
            mv.Wallet.GetAddresses()[0], p)
        assert.NotNil(t, err)
        assert.Equal(t, err.Error(), "Visor disabled")
        assert.True(t, gc.LastSent.IsZero())
    })

    // Spending but spend fails (no money)
    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{1000 * 10e6, 0}, 0,
            mv.Wallet.GetAddresses()[0], p)
        assert.NotNil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 0)
        assert.True(t, gc.LastSent.IsZero())
    })

    // Spending succeeds, but didn't announce
    vc, mv = setupVisor()
    vc.Disabled = false
    v = NewVisor(vc)
    gc.Conn = NewFailingConn(addr)
    assert.Nil(t, transferCoins(mv, v.Visor))
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{10e6, 0}, 0,
            mv.Wallet.GetAddresses()[0], p)
        assert.Nil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
        for _, tx := range v.Visor.UnconfirmedTxns.Txns {
            assert.True(t, tx.Announced.IsZero())
        }
        assert.True(t, gc.LastSent.IsZero())
    })

    // Spending succeeds, and announced
    vc, mv = setupVisor()
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.Nil(t, transferCoins(mv, v.Visor))
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{10e6, 0}, 0,
            mv.Wallet.GetAddresses()[0], p)
        assert.Nil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
        for _, tx := range v.Visor.UnconfirmedTxns.Txns {
            assert.False(t, tx.Announced.IsZero())
        }
        assert.False(t, gc.LastSent.IsZero())
    })
}

func TestVisorResendTransaction(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, mv := setupVisor()
    v := NewVisor(vc)
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 0)

    // Nothing should happen if txn unknown
    assert.False(t, v.ResendTransaction(coin.SumSHA256([]byte("garbage")), p))
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // give the visor some coins, and make a spend to add a txn
    assert.Nil(t, transferCoins(mv, v.Visor))
    tx, err := v.Spend(visor.Balance{10e6, 0}, 0,
        mv.Wallet.GetAddresses()[0], p)
    assert.Nil(t, err)
    assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
    h := tx.Hash()
    ut := v.Visor.UnconfirmedTxns.Txns[h]
    ut.Announced = util.ZeroTime()
    v.Visor.UnconfirmedTxns.Txns[h] = ut
    assert.True(t, v.Visor.UnconfirmedTxns.Txns[h].Announced.IsZero())
    // Reset the sent timer since we made a successful spend
    gc.LastSent = util.ZeroTime()

    // Nothing should send if disabled
    v.Config.Disabled = true
    assert.False(t, v.ResendTransaction(h, p))
    ann := v.Visor.UnconfirmedTxns.Txns[h].Announced
    assert.True(t, ann.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // Nothing should send if failed to send
    gc.Conn = NewFailingConn(addr)
    v.Config.Disabled = false
    assert.False(t, v.ResendTransaction(h, p))
    ann = v.Visor.UnconfirmedTxns.Txns[h].Announced
    assert.True(t, ann.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // Should have resent
    gc.Conn = NewDummyConn(addr)
    assert.True(t, v.ResendTransaction(h, p))
    ann = v.Visor.UnconfirmedTxns.Txns[h].Announced
    assert.False(t, ann.IsZero())
    assert.False(t, gc.LastSent.IsZero())
}

func TestCreateAndPublishBlock(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, mv := setupVisor()
    dest := visor.NewWalletEntry()

    // Disabled
    vc.Disabled = true
    vc.Config = mv.Config
    v := NewVisor(vc)
    v.Visor = mv
    err, published := v.CreateAndPublishBlock(p)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Visor disabled")
    assert.False(t, published)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(0))

    // Created, but failed to send
    vc.Disabled = false
    vc.Config = mv.Config
    v = NewVisor(vc)
    gc.Conn = NewFailingConn(addr)
    _, err = v.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address, p)
    assert.Nil(t, err)
    err, published = v.CreateAndPublishBlock(p)
    assert.Nil(t, err)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))
    assert.False(t, published)

    // Created and sent
    vc.Config.IsMaster = true
    vc.Config = mv.Config
    v = NewVisor(vc)
    gc.Conn = NewDummyConn(addr)
    _, err = v.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address, p)
    assert.Nil(t, err)
    err, published = v.CreateAndPublishBlock(p)
    assert.Nil(t, err)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))
    assert.True(t, published)

    // Can't create, don't have coins
    vc, _ = setupVisor()
    vc.Config.IsMaster = false
    vc.Disabled = false
    v = NewVisor(vc)
    _, err = v.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address, p)
    assert.NotNil(t, err)
    err, published = v.CreateAndPublishBlock(p)
    assert.NotNil(t, err)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(0))
    assert.False(t, published)
}

func TestRecordBlockchainLength(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.recordBlockchainLength(addr, uint64(6)) })
    assert.Equal(t, v.blockchainLengths[addr], uint64(6))
    v.blockchainLengths[addr] = uint64(7)
    assert.NotPanics(t, func() { v.recordBlockchainLength(addr, uint64(5)) })
    assert.Equal(t, v.blockchainLengths[addr], uint64(5))
}

func TestEstimateBlockchainLength(t *testing.T) {
    defer cleanupVisor()
    vc, mv := setupVisor()
    v := NewVisor(vc)
    assert.Nil(t, transferCoins(mv, v.Visor))
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))
    // With no peers reporting, returns our own blockchain length
    assert.Equal(t, v.EstimateBlockchainLength(), uint64(2))

    // With only 1 peer reporting, returns our own blockchain length
    v.recordBlockchainLength(addr, 1)
    assert.Equal(t, v.EstimateBlockchainLength(), uint64(2))

    // If the average reported is lower than our own, return our own
    v.recordBlockchainLength(addr, 1)
    v.recordBlockchainLength(addrb, 1)
    assert.Equal(t, v.EstimateBlockchainLength(), uint64(2))

    // Should return the exact median if odd # of them
    v.recordBlockchainLength(addr, 1)
    v.recordBlockchainLength(addrb, 5)
    v.recordBlockchainLength(addrc, 9)
    assert.Equal(t, v.EstimateBlockchainLength(), uint64(5))

    // Should return the average of the median values, if even # of them
    v.recordBlockchainLength(addr, 1)
    v.recordBlockchainLength(addrb, 5)
    v.recordBlockchainLength(addrc, 9)
    v.recordBlockchainLength(addrd, 100)
    assert.Equal(t, v.EstimateBlockchainLength(), uint64(7))
}

/* Visor Messages */

func TestGetBlocksMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    m := NewGetBlocksMessage(uint64(1))
    assert.Equal(t, m.LastBlock, uint64(1))
    testSimpleMessageHandler(t, d, m)
}

func TestGetBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    assert.Nil(t, transferCoins(mv, d.Visor.Visor))
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(1))
    m := NewGetBlocksMessage(uint64(7))
    m.c = messageContext(addr)
    d.Visor.Config.Disabled = true
    m.Process(d)
    // Disabled should not record a blockchain length
    assert.Equal(t, len(d.Visor.blockchainLengths), 0)

    // Enabled handler, should record bc length not not send anything since
    // we have nothing new enough
    d.Visor.Config.Disabled = false
    m.Process(d)
    assert.Equal(t, len(d.Visor.blockchainLengths), 1)
    assert.Equal(t, d.Visor.blockchainLengths[addr], uint64(7))
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We have something for them, but failed to send
    m.c.Conn.Conn = NewFailingConn(addr)
    m.LastBlock = uint64(0)
    m.Process(d)
    assert.Equal(t, len(d.Visor.blockchainLengths), 1)
    assert.Equal(t, d.Visor.blockchainLengths[addr], uint64(0))
    assert.True(t, m.c.Conn.LastSent.IsZero())

    m.c.Conn.Conn = NewDummyConn(addr)
    m.Process(d)
    assert.Equal(t, len(d.Visor.blockchainLengths), 1)
    assert.Equal(t, d.Visor.blockchainLengths[addr], uint64(0))
    assert.False(t, m.c.Conn.LastSent.IsZero())
}

func TestGiveBlocksMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    _, mv := setupVisor()
    blocks := []visor.SignedBlock{mv.GetGenesisBlock()}
    m := NewGiveBlocksMessage(blocks)
    assert.Equal(t, m.Blocks, blocks)
    testSimpleMessageHandler(t, d, m)
}

func TestGiveBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)

    blocks, err := makeBlocks(mv, 2)
    assert.Nil(t, err)
    assert.Equal(t, len(blocks), 2)
    m := NewGiveBlocksMessage(blocks)
    m.c = messageContext(addr)

    // Disabled should have nothing happen
    d.Visor.Config.Disabled = true
    m.Process(d)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(0))

    // Not disabled, should add blocks, but fail to send
    gc.Conn = NewFailingConn(addr)
    d.Visor.Config.Disabled = false
    m.Process(d)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(2))
    assert.True(t, gc.LastSent.IsZero())

    // Not disabled and blocks were reannounced
    gc.Conn = NewDummyConn(addr)
    blocks, err = makeBlocks(mv, 2)
    assert.Nil(t, err)
    assert.Equal(t, len(blocks), 2)
    assert.Equal(t, len(blocks), 2)
    m = NewGiveBlocksMessage(blocks)
    m.c = messageContext(addr)
    m.Process(d)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(4))
    assert.False(t, gc.LastSent.IsZero())

    // Send blocks we have and some we dont, as long as they are in order
    // we can use the ones at the end
    gc.LastSent = util.ZeroTime()
    moreBlocks, err := makeBlocks(mv, 2)
    assert.Nil(t, err)
    blocks = append(blocks, moreBlocks...)
    m = NewGiveBlocksMessage(blocks)
    m.c = messageContext(addr)
    m.Process(d)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(6))
    assert.False(t, gc.LastSent.IsZero())

    // Send invalid blocks
    gc.LastSent = util.ZeroTime()
    bb := visor.SignedBlock{
        Block: coin.Block{
            Header: coin.BlockHeader{
                BkSeq: uint64(7),
            }}}
    m = NewGiveBlocksMessage([]visor.SignedBlock{bb})
    m.c = messageContext(addr)
    m.Process(d)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(6))
    assert.True(t, gc.LastSent.IsZero())
}

func TestAnnounceBlocksMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    m := NewAnnounceBlocksMessage(uint64(7))
    assert.Equal(t, m.MaxBkSeq, uint64(7))
    testSimpleMessageHandler(t, d, m)
}

func TestAnnounceBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    assert.Nil(t, transferCoins(mv, d.Visor.Visor))
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(1))

    // Disabled, nothing should happen
    d.Visor.Config.Disabled = true
    m := NewAnnounceBlocksMessage(uint64(2))
    m.c = messageContext(addr)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We know all the blocks
    d.Visor.Config.Disabled = false
    m.MaxBkSeq = uint64(1)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We send a GetBlocksMessage in response to a higher MaxBkSeq
    m.MaxBkSeq = uint64(7)
    assert.False(t, d.Visor.Visor.MostRecentBkSeq() >= m.MaxBkSeq)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.False(t, m.c.Conn.LastSent.IsZero())

    // We send a GetBlocksMessage in response to a higher MaxBkSeq,
    // but the send failed
    m.c.Conn.LastSent = util.ZeroTime()
    m.c.Conn.Conn = NewFailingConn(addr)
    m.MaxBkSeq = uint64(7)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())
}

func TestAnnounceTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewAnnounceTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)
}

func TestAnnounceTxnsMessageProcess(t *testing.T) {
    v, _ := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)

    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewAnnounceTxnsMessage(txns)
    m.c = messageContext(addr)

    // Disabled, nothing should happen
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We don't know some, request them
    d.Visor.Config.Disabled = false
    assert.NotPanics(t, func() { m.Process(d) })
    assert.False(t, m.c.Conn.LastSent.IsZero())

    // We don't know some, request them, but fail to request
    m.c.Conn.LastSent = util.ZeroTime()
    m.c.Conn.Conn = NewFailingConn(addr)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We know all the reported txns, nothing should be sent
    d.Visor.Visor.UnconfirmedTxns.Txns[tx.Txn.Hash()] = tx
    m.c.Conn.Conn = NewDummyConn(addr)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())
}

func TestGetTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewGetTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)
}

func TestGetTxnsMessageProcess(t *testing.T) {
    v, _ := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)

    tx := createUnconfirmedTxn()
    tx.Txn.Header.Hash = coin.SumSHA256([]byte("asdadwadwada"))
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewGetTxnsMessage(txns)
    m.c = messageContext(addr)

    // We don't have any to reply with
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // Disabled, nothing should happen
    d.Visor.Visor.UnconfirmedTxns.Txns[tx.Txn.Hash()] = tx
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We have some to reply with
    d.Visor.Config.Disabled = false
    assert.NotPanics(t, func() { m.Process(d) })
    assert.False(t, m.c.Conn.LastSent.IsZero())

    // We have some to reply with, but fail to send
    m.c.Conn.LastSent = util.ZeroTime()
    m.c.Conn.Conn = NewFailingConn(addr)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())
}

func TestGiveTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := []coin.Transaction{tx.Txn}
    m := NewGiveTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)
}

func TestGiveTxnsMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)

    utx := createUnconfirmedTxn()
    txns := []coin.Transaction{utx.Txn}
    m := NewGiveTxnsMessage(txns)
    m.c = messageContext(addr)

    // No valid txns, nothing should be sent
    assert.NotPanics(t, func() { m.Process(d) })
    assert.Equal(t, len(mv.UnconfirmedTxns.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Disabled, nothing should happen
    tx, err := makeValidTxn(mv)
    assert.Nil(t, err)
    m.Txns = []coin.Transaction{tx}
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    assert.Equal(t, len(mv.UnconfirmedTxns.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // A valid txn, we should broadcast. Txn's announce state should be updated
    d.Visor.Config.Disabled = false
    assert.True(t, gc.LastSent.IsZero())
    assert.NotPanics(t, func() { m.Process(d) })
    assert.Equal(t, len(d.Visor.Visor.UnconfirmedTxns.Txns), 1)
    assert.False(t, gc.LastSent.IsZero())
    ut, ok := d.Visor.Visor.UnconfirmedTxns.Txns[tx.Hash()]
    assert.True(t, ok)
    now := util.Now()
    assert.False(t, ut.Announced.IsZero())
    assert.True(t, ut.Announced.Add(time.Second).After(now))

    // A valid txn, but we fail to broadcast.  Txn's announce state should not
    // have been updated
    tx, err = makeValidTxn(mv)
    assert.Nil(t, err)
    m.Txns = []coin.Transaction{tx}
    gc.LastSent = util.ZeroTime()
    gc.Conn = NewFailingConn(addr)
    assert.NotPanics(t, func() { m.Process(d) })
    assert.Equal(t, len(d.Visor.Visor.UnconfirmedTxns.Txns), 2)
    assert.True(t, gc.LastSent.IsZero())
    ut = d.Visor.Visor.UnconfirmedTxns.Txns[tx.Hash()]
    assert.True(t, ut.Announced.IsZero())
}

/* Misc */

func TestBlockchainLengths(t *testing.T) {
    b := make(BlockchainLengths, 10)
    for i := 0; i < 10; i++ {
        b[i] = uint64(9 - i)
    }
    assert.Equal(t, b.Len(), 10)
    assert.True(t, b.Less(4, 3))
    assert.Equal(t, b[4], uint64(5))
    assert.Equal(t, b[3], uint64(6))
    b.Swap(4, 3)
    assert.Equal(t, b[4], uint64(6))
    assert.Equal(t, b[3], uint64(5))
    sort.Sort(b)
    for i := 0; i < 10; i++ {
        assert.Equal(t, b[i], uint64(i))
    }
}
