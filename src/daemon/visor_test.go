package daemon

import (
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "github.com/stretchr/testify/assert"
    "os"
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
    // Create testmaster.keys file
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
    ut.Received = time.Now().UTC()
    ut.Checked = ut.Received
    ut.Announced = ut.Received
    ut.IsOurSpend = true
    ut.IsOurReceive = true
    return ut
}

func addUnconfirmedTxn(v *Visor) visor.UnconfirmedTxn {
    ut := createUnconfirmedTxn()
    v.Visor.UnconfirmedTxns.Txns[ut.Txn.Header.Hash] = ut
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

func transferCoins(mv *visor.Visor, v *visor.Visor) error {
    // Give the nonmaster some money to spend
    addr := v.Wallet.Entries[0].Address
    tx, err := mv.Spend(visor.Balance{10 * 1e6, 0}, 0, addr)
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    dest := visor.NewWalletEntry()
    _, err := d.Visor.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address,
        d.Pool)
    assert.Nil(t, err)
    if master && !published {
        assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    } else {
        assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
    }
    time.Sleep(time.Millisecond * 1250)
    final := start
    if master {
        final += 1
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(final))
    if !published {
        gc.LastSent = time.Unix(0, 0)
    }
    if published {
        assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
    } else {
        assert.Equal(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
}

func TestBlocksAnnounceTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.BlocksAnnounceRate = time.Millisecond * 10
    d, quit := newVisorDaemon(vc)
    gc := gnetConnection(addr)
    d.Pool.Pool.Pool[gc.Id] = gc
    d.Pool.Pool.Addresses[gc.Addr()] = gc
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    time.Sleep(time.Millisecond * 15)
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    vc.Disabled = false
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
}

func TestVisorAnnounceBlocks(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    vc.Disabled = false
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))

    gc.LastSent = time.Unix(0, 0)
    gc.Conn = NewFailingConn(addr)
    assert.NotPanics(t, func() {
        assert.NotNil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    gc.LastSent = time.Unix(0, 0)
    gc.Conn = NewDummyConn(addr)
    delete(p.Pool.Pool, gc.Id)
    delete(p.Pool.Addresses, gc.Addr())
    assert.NotPanics(t, func() {
        assert.NotNil(t, v.RequestBlocksFromAddr(p, addr))
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // With no transactions, nothing should be sent
    vc.Disabled = false
    vc.TransactionRebroadcastRate = time.Millisecond * 5
    v = NewVisor(vc)
    time.Sleep(time.Millisecond * 20)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // We have a stale owned unconfirmed txn but we failed to send it
    gc.Conn = NewFailingConn(addr)
    v = NewVisor(vc)
    tx := addUnconfirmedTxn(v)
    time.Sleep(time.Millisecond * 20)
    assert.NotPanics(t, func() {
        v.BroadcastOurTransactions(p)
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    newtx := v.Visor.UnconfirmedTxns.Txns[tx.Txn.Header.Hash]
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
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
    newtx = v.Visor.UnconfirmedTxns.Txns[tx.Txn.Header.Hash]
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // Fail to send
    gc.Conn = NewFailingConn(addr)
    vc.Disabled = false
    v = NewVisor(vc)
    sb := v.Visor.GetGenesisBlock()
    assert.NotPanics(t, func() {
        // Returns error if nobody received
        assert.NotNil(t, v.broadcastBlock(sb, p))
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // Succeed in sending
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    sb = v.Visor.GetGenesisBlock()
    assert.NotPanics(t, func() {
        assert.Nil(t, v.broadcastBlock(sb, p))
    })
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
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
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // Fail to send
    gc.Conn = NewFailingConn(addr)
    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        // Returns error if nobody received
        assert.NotNil(t, v.broadcastTransaction(ut.Txn, p))
    })
    assert.Equal(t, gc.LastSent, time.Unix(0, 0))

    // Succeed in sending
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.broadcastTransaction(ut.Txn, p))
    })
    assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
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
            mv.Wallet.Entries[0].Address, p)
        assert.NotNil(t, err)
        assert.Equal(t, err.Error(), "Visor disabled")
        assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    })

    // Spending but spend fails (no money)
    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{1000 * 10e6, 0}, 0,
            mv.Wallet.Entries[0].Address, p)
        assert.NotNil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 0)
        assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    })

    // Spending succeeds, but didn't announce
    vc, mv = setupVisor()
    vc.Disabled = false
    v = NewVisor(vc)
    gc.Conn = NewFailingConn(addr)
    assert.Nil(t, transferCoins(mv, v.Visor))
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{10e6, 0}, 0,
            mv.Wallet.Entries[0].Address, p)
        assert.Nil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
        for _, tx := range v.Visor.UnconfirmedTxns.Txns {
            assert.Equal(t, tx.Announced, time.Unix(0, 0))
        }
        assert.Equal(t, gc.LastSent, time.Unix(0, 0))
    })

    // Spending succeeds, and announced
    vc, mv = setupVisor()
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.Nil(t, transferCoins(mv, v.Visor))
    assert.NotPanics(t, func() {
        _, err := v.Spend(visor.Balance{10e6, 0}, 0,
            mv.Wallet.Entries[0].Address, p)
        assert.Nil(t, err)
        assert.Equal(t, len(v.Visor.UnconfirmedTxns.Txns), 1)
        for _, tx := range v.Visor.UnconfirmedTxns.Txns {
            assert.NotEqual(t, tx.Announced, time.Unix(0, 0))
        }
        assert.NotEqual(t, gc.LastSent, time.Unix(0, 0))
    })
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
