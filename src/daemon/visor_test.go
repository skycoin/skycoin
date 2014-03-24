package daemon

import (
    "crypto/rand"
    "github.com/skycoin/encoder"
    "github.com/skycoin/gnet"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/skycoin/skycoin/src/visor"
    "github.com/skycoin/skycoin/src/wallet"
    "github.com/stretchr/testify/assert"
    "os"
    "path/filepath"
    "sort"
    "testing"
    "time"
)

const (
    testMasterKeysFile = "testmaster.keys"
    testWalletFile     = "testwallet.wlt"
    testWalletDir      = "./"
    testBlocksigsFile  = "testblockchain.sigs"
    testBlockchainFile = "testblockchain.bin"
)

var (
    fullWalletFile = filepath.Join(testWalletDir, testWalletFile)
)

func randBytes(t *testing.T, n int) []byte {
    b := make([]byte, n)
    x, err := rand.Read(b)
    assert.Equal(t, n, x)
    assert.Nil(t, err)
    return b
}

func randSHA256(t *testing.T) coin.SHA256 {
    return coin.SumSHA256(randBytes(t, 32))
}

func createGenesisSignature(master wallet.WalletEntry) coin.Sig {
    c := visor.NewVisorConfig()
    bc := coin.NewBlockchain()
    gb := bc.CreateGenesisBlock(master.Address, c.GenesisTimestamp,
        c.GenesisCoinVolume)
    return coin.SignHash(gb.HashHeader(), master.Secret)
}

// Returns an appropriate VisorConfig and a master visor
func setupVisor() (VisorConfig, *visor.Visor) {
    coin.SetAddressVersion("test")

    // Make a new master visor + blockchain
    // Get the signed genesis block,
    mw := wallet.NewWalletEntry()
    mvc := visor.NewVisorConfig()
    mvc.IsMaster = true
    mvc.MasterKeys = mw
    mvc.CoinHourBurnFactor = 0
    mvc.GenesisSignature = createGenesisSignature(mw)
    mv := visor.NewVisor(mvc)

    // Use the master values for a client configuration
    c := NewVisorConfig()
    c.Config.WalletDirectory = testWalletDir
    c.Config.IsMaster = false
    c.Config.MasterKeys = mw
    c.Config.MasterKeys.Secret = coin.SecKey{}
    c.Config.GenesisSignature = mvc.GenesisSignature
    c.Config.GenesisTimestamp = mvc.GenesisTimestamp
    return c, mv
}

func setupMasterVisor() VisorConfig {
    cleanupVisor()
    coin.SetAddressVersion("test")
    c := NewVisorConfig()
    c.Config.IsMaster = true
    mw := wallet.NewWalletEntry()
    c.Config.MasterKeys = mw
    c.Config.GenesisSignature = createGenesisSignature(mw)
    return c
}

func cleanupVisor() {
    os.Remove(fullWalletFile)
    os.Remove(testMasterKeysFile)
    os.Remove(testBlockchainFile)
    os.Remove(testBlocksigsFile)
    os.Remove(fullWalletFile + ".tmp")
    os.Remove(testMasterKeysFile + ".tmp")
    os.Remove(testBlockchainFile + ".tmp")
    os.Remove(testBlocksigsFile + ".tmp")
    os.Remove(fullWalletFile + ".bak")
    os.Remove(testMasterKeysFile + ".bak")
    os.Remove(testBlockchainFile + ".bak")
    os.Remove(testBlocksigsFile + ".bak")
    wallets, err := filepath.Glob("*." + wallet.WalletExt)
    if err != nil {
        logger.Critical("Failed to glob wallet files: %v", err)
    } else {
        for _, w := range wallets {
            os.Remove(w)
            os.Remove(w + ".bak")
            os.Remove(w + ".tmp")
        }
    }
}

// Returns a daemon with the visor enabled, but networking disabled
func newVisorDaemon(vc VisorConfig) (*Daemon, chan int) {
    quit := make(chan int)
    c := NewConfig()
    c.Daemon.DisableNetworking = true
    c.Visor = vc
    c.Visor.Config.CoinHourBurnFactor = 0
    d := NewDaemon(c)
    return d, quit
}

// Writes a wallet entry to disk at filename
func writeMasterKeysFile() (wallet.WalletEntry, error) {
    we := wallet.NewWalletEntry()
    rwe := wallet.NewReadableWalletEntry(&we)
    err := rwe.Save(testMasterKeysFile)
    return we, err
}

func assertFileExists(t *testing.T, filename string) {
    stat, err := os.Stat(filename)
    assert.Nil(t, err)
    assert.NotNil(t, stat)
    if stat != nil {
        assert.True(t, stat.Mode().IsRegular())
    }
}

func assertFileNotExists(t *testing.T, filename string) {
    _, err := os.Stat(filename)
    assert.NotNil(t, err)
    assert.True(t, os.IsNotExist(err))
}

func createUnconfirmedTxn() visor.UnconfirmedTxn {
    now := util.Now()
    return visor.UnconfirmedTxn{
        Txn: coin.Transaction{
            Head: coin.TransactionHeader{
                Hash: coin.SumSHA256([]byte("cascas")),
            },
        },
        Received:  now,
        Checked:   now,
        Announced: util.ZeroTime(),
    }
}

func addUnconfirmedTxn(v *Visor) visor.UnconfirmedTxn {
    ut := createUnconfirmedTxn()
    v.Visor.Unconfirmed.Txns[ut.Txn.Hash()] = ut
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
    we := wallet.NewWalletEntry()
    return mv.Spend(mv.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0}, 0,
        we.Address)
}

func makeValidTxnNoError(t *testing.T, mv *visor.Visor) coin.Transaction {
    we := wallet.NewWalletEntry()
    tx, err := mv.Spend(mv.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0}, 0,
        we.Address)
    assert.Nil(t, err)
    return tx
}

func transferCoins(mv *visor.Visor, v *visor.Visor) error {
    // Give the nonmaster some money to spend
    addr := v.Wallets[0].GetAddresses()[0]
    tx, err := mv.Spend(mv.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0}, 0,
        addr)
    if err != nil {
        return err
    }
    mv.RecordTxn(tx)
    sb, err := mv.CreateAndExecuteBlock()
    if err != nil {
        return err
    }
    return v.ExecuteSignedBlock(sb)
}

func makeMoreBlocks(mv *visor.Visor, n int,
    now uint64) ([]visor.SignedBlock, error) {
    dest := wallet.NewWalletEntry()
    blocks := make([]visor.SignedBlock, n)
    for i := 0; i < n; i++ {
        tx, err := mv.Spend(mv.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0},
            0, dest.Address)
        if err != nil {
            return nil, err
        }
        mv.RecordTxn(tx)
        sb, err := mv.CreateBlock(now + uint64(i) + 1)
        if err != nil {
            return nil, err
        }
        err = mv.ExecuteSignedBlock(sb)
        if err != nil {
            return nil, err
        }
        blocks[i] = sb
    }
    return blocks, nil
}

func makeBlocks(mv *visor.Visor, n int) ([]visor.SignedBlock, error) {
    return makeMoreBlocks(mv, n, uint64(util.UnixNow()))
}

/* Tests for daemon's loop related to visor */

func testBlockCreationTicker(t *testing.T, vcfg VisorConfig, master bool,
    mv *visor.Visor) {
    vcfg.Config.BlockCreationInterval = 1
    defer gnet.EraseMessages()
    c := NewConfig()
    c.Visor = vcfg
    c.Daemon.DisableNetworking = false
    d := NewDaemon(c)
    if !master {
        err := transferCoins(mv, d.Visor.Visor)
        assert.Nil(t, err)
    }
    quit := make(chan int)
    defer closeDaemon(d, quit)
    gc := setupExistingPool(d.Pool)
    go d.Pool.Pool.ConnectionWriteLoop(gc)
    assert.True(t, gc.LastSent.IsZero())
    assert.Equal(t, len(d.Pool.Pool.Pool), 1)
    assert.Equal(t, d.Pool.Pool.Pool[gc.Id], gc)
    assert.Equal(t, len(d.Pool.Pool.Addresses), 1)
    start := 0
    if !master {
        start = 1
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(start))
    go d.Start(quit)
    time.Sleep(time.Second + (time.Millisecond * 50))
    // Creation should not have occured, because no transaction
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(start))
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)

    // Creation should occur with a transaction, if not a master
    // Make a transaction
    assert.False(t, d.Visor.Config.Disabled)
    assert.True(t, gc.LastSent.IsZero())
    dest := wallet.NewWalletEntry()
    tx, err := d.Visor.Spend(d.Visor.Visor.Wallets[0].GetID(),
        visor.Balance{10 * 1e6, 0}, 0, dest.Address, d.Pool)
    wait()
    assert.Nil(t, err)
    assert.Equal(t, d.Pool.Pool.Pool[gc.Id], gc)
    assert.Equal(t, len(d.Pool.Pool.DisconnectQueue), 0)
    assert.Equal(t, len(d.Pool.Pool.Pool), 1)
    assert.Equal(t, len(d.Pool.Pool.Addresses), 1)
    // Since the daemon loop is running, it will have processed the SendResult
    // Instead we can check if the txn was announced
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    ut := d.Visor.Visor.Unconfirmed.Txns[tx.Hash()]
    assert.False(t, ut.Announced.IsZero())
    assert.False(t, gc.LastSent.IsZero())
    ls := gc.LastSent

    // Now, block should be created
    time.Sleep(time.Second + (time.Millisecond * 50))
    final := start
    if master {
        final += 1
    }
    assert.Equal(t, len(d.Pool.Pool.Pool), 1)
    // Again, we can't check SendResults since the daemon loop is running.
    // We can only check that LastSent was updated, if its the master and it
    // created the block.
    if master {
        assert.True(t, gc.LastSent.After(ls))
    } else {
        assert.Equal(t, gc.LastSent, ls)
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(final))
    assert.False(t, gc.LastSent.IsZero())
}

func TestBlockCreationTicker(t *testing.T) {
    defer cleanupVisor()
    vcfg, mv := setupVisor()
    // No blocks should get created if we are not master
    testBlockCreationTicker(t, vcfg, false, mv)
}

func TestBlockCreationTickerMaster(t *testing.T) {
    defer cleanupVisor()
    vcfg := setupMasterVisor()
    // Master should make a block
    testBlockCreationTicker(t, vcfg, true, nil)
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
    assert.Equal(t, len(d.Visor.Visor.Unconfirmed.Txns), 0)
}

func TestBlocksRequestTicker(t *testing.T) {
    defer cleanupVisor()
    vc, _ := setupVisor()
    vc.BlocksRequestRate = time.Millisecond * 10
    d, quit := newVisorDaemon(vc)
    d.Config.DisableNetworking = false
    gc := gnetConnection(addr)
    go d.Pool.Pool.ConnectionWriteLoop(gc)
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
    d.Config.DisableNetworking = false
    gc := gnetConnection(addr)
    go d.Pool.Pool.ConnectionWriteLoop(gc)
    d.Pool.Pool.Pool[gc.Id] = gc
    d.Pool.Pool.Addresses[gc.Addr()] = gc
    assert.True(t, gc.LastSent.IsZero())
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
    c.Config.WalletDirectory = testWalletDir
    v := NewVisor(c)
    assert.NotPanics(t, v.Shutdown)
    // Should not save anything
    assertFileNotExists(t, testBlockchainFile)
    assertFileNotExists(t, testBlocksigsFile)
    assertFileNotExists(t, fullWalletFile)
    wallets, err := filepath.Glob(testWalletDir + "*.wlt")
    assert.Nil(t, err)
    assert.Equal(t, len(wallets), 0)
    cleanupVisor()

    c.Disabled = false
    v = NewVisor(c)
    v.Visor.Wallets[0].SetFilename(testWalletFile)
    assert.NotPanics(t, v.Shutdown)
    assertFileExists(t, testBlockchainFile)
    assertFileExists(t, testBlocksigsFile)
    assertFileExists(t, fullWalletFile)
    cleanupVisor()

    // If master, no wallet should be saved
    c = setupMasterVisor()
    c.Config.BlockchainFile = testBlockchainFile
    c.Config.BlockSigsFile = testBlocksigsFile
    c.Config.WalletDirectory = testWalletDir
    v = NewVisor(c)
    v.Visor.Wallets[0].SetFilename(testWalletFile)
    assert.NotPanics(t, v.Shutdown)
    assertFileExists(t, testBlockchainFile)
    assertFileExists(t, testBlocksigsFile)
    assertFileNotExists(t, fullWalletFile)
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
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 1)
    wait()

    v.Config.Disabled = true
    v.RefreshUnconfirmed()
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 1)

    v.Config.Disabled = false
    v.RefreshUnconfirmed()
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 0)
}

func TestVisorRequestBlocks(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    go p.Pool.ConnectionWriteLoop(gc)

    // Disabled
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Valid
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() { v.RequestBlocks(p) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GetBlocksMessage)
    assert.True(t, ok)
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorAnnounceBlocks(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    go p.Pool.ConnectionWriteLoop(gc)

    // Disabled
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Valid send
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.False(t, v.Config.Disabled)
    assert.NotPanics(t, func() { v.AnnounceBlocks(p) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*AnnounceBlocksMessage)
    assert.True(t, ok)
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorRequestBlocksFromAddr(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    go p.Pool.ConnectionWriteLoop(gc)

    // Disabled
    vc.Disabled = true
    v := NewVisor(vc)
    assert.NotPanics(t, func() {
        err := v.RequestBlocksFromAddr(p, addr)
        assert.NotNil(t, err)
        assert.Equal(t, err.Error(), "Visor disabled")
    })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    vc.Disabled = false
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        assert.Nil(t, v.RequestBlocksFromAddr(p, addr))
    })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GetBlocksMessage)
    assert.True(t, ok)
    assert.False(t, gc.LastSent.IsZero())

    // No connection found for addr
    gc.LastSent = util.ZeroTime()
    gc.Conn = NewDummyConn(addr)
    delete(p.Pool.Pool, gc.Id)
    delete(p.Pool.Addresses, gc.Addr())
    assert.NotPanics(t, func() {
        assert.NotNil(t, v.RequestBlocksFromAddr(p, addr))
    })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())
}

func TestVisorBroadcastBlock(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    // Should not send anything if disabled
    assert.NotPanics(t, func() {
        v.broadcastBlock(visor.SignedBlock{}, p)
    })
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Sending
    gc.Conn = NewDummyConn(addr)
    vc.Disabled = false
    v = NewVisor(vc)
    sb := v.Visor.GetGenesisBlock()
    assert.NotPanics(t, func() {
        v.broadcastBlock(sb, p)
    })
    go p.Pool.ConnectionWriteLoop(gc)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GiveBlocksMessage)
    assert.True(t, ok)
    assert.Nil(t, sr.Error)
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorBroadcastTransaction(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    go p.Pool.ConnectionWriteLoop(gc)
    vc, _ := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    ut := createUnconfirmedTxn()
    assert.NotPanics(t, func() {
        v.broadcastTransaction(ut.Txn, p)
    })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Sending
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.NotPanics(t, func() {
        v.broadcastTransaction(ut.Txn, p)
    })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GiveTxnsMessage)
    assert.True(t, ok)
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorSpend(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    go p.Pool.ConnectionWriteLoop(gc)
    vc, mv := setupVisor()
    vc.Disabled = true
    v := NewVisor(vc)
    // Spending while disabled
    _, err := v.Spend("xxx", visor.Balance{10e6, 0}, 0,
        mv.Wallets[0].GetAddresses()[0], p)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Visor disabled")
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Spending but spend fails (no money)
    vc.Disabled = false
    v = NewVisor(vc)
    _, err = v.Spend(v.Visor.Wallets[0].GetID(),
        visor.Balance{1000 * 10e6, 0}, 0, mv.Wallets[0].GetAddresses()[0],
        p)
    wait()
    assert.NotNil(t, err)
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Spending succeeds, and announced
    vc, mv = setupVisor()
    vc.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v = NewVisor(vc)
    assert.Nil(t, transferCoins(mv, v.Visor))
    _, err = v.Spend(v.Visor.Wallets[0].GetID(), visor.Balance{10e6, 0},
        0, mv.Wallets[0].GetAddresses()[0], p)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Equal(t, sr.Connection, gc)
    assert.Nil(t, sr.Error)
    _, ok := sr.Message.(*GiveTxnsMessage)
    assert.True(t, ok)
    assert.Nil(t, err)
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 1)
    assert.False(t, gc.LastSent.IsZero())
}

func TestVisorResendTransaction(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    go p.Pool.ConnectionWriteLoop(gc)
    vc, mv := setupVisor()
    v := NewVisor(vc)
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 0)

    // Nothing should happen if txn unknown
    v.ResendTransaction(coin.SumSHA256([]byte("garbage")), p)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // give the visor some coins, and make a spend to add a txn
    assert.Nil(t, transferCoins(mv, v.Visor))
    tx, err := v.Spend(v.Visor.Wallets[0].GetID(), visor.Balance{10e6, 0}, 0,
        mv.Wallets[0].GetAddresses()[0], p)
    assert.Nil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    <-p.Pool.SendResults
    assert.Equal(t, len(v.Visor.Unconfirmed.Txns), 1)
    h := tx.Hash()
    ut := v.Visor.Unconfirmed.Txns[h]
    ut.Announced = util.ZeroTime()
    v.Visor.Unconfirmed.Txns[h] = ut
    assert.True(t, v.Visor.Unconfirmed.Txns[h].Announced.IsZero())
    // Reset the sent timer since we made a successful spend
    gc.LastSent = util.ZeroTime()

    // Nothing should send if disabled
    v.Config.Disabled = true
    v.ResendTransaction(h, p)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    ann := v.Visor.Unconfirmed.Txns[h].Announced
    assert.True(t, ann.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // Should have resent
    v.Config.Disabled = false
    gc.Conn = NewDummyConn(addr)
    v.ResendTransaction(h, p)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GiveTxnsMessage)
    assert.True(t, ok)
    ann = v.Visor.Unconfirmed.Txns[h].Announced
    // Announced state should not be updated until we process it
    assert.True(t, ann.IsZero())
    assert.False(t, gc.LastSent.IsZero())
}

func TestCreateAndPublishBlock(t *testing.T) {
    defer cleanupVisor()
    defer gnet.EraseMessages()
    p, gc := setupPool()
    vc, mv := setupVisor()
    dest := wallet.NewWalletEntry()
    go p.Pool.ConnectionWriteLoop(gc)

    // Disabled
    vc.Disabled = true
    vc.Config = mv.Config
    v := NewVisor(vc)
    v.Visor = mv
    err := v.CreateAndPublishBlock(p)
    assert.NotNil(t, err)
    wait()
    assert.Equal(t, err.Error(), "Visor disabled")
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(0))

    // Created and sent
    vc.Disabled = false
    vc.Config.IsMaster = true
    vc.Config = mv.Config
    v = NewVisor(vc)
    gc.Conn = NewDummyConn(addr)
    _, err = v.Spend(v.Visor.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0}, 0,
        dest.Address, p)
    assert.Nil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    <-p.Pool.SendResults
    err = v.CreateAndPublishBlock(p)
    wait()
    assert.Nil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*GiveBlocksMessage)
    assert.True(t, ok)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))

    // Can't create, don't have coins
    // First, spend all of our coins
    // vc2, _ := setupVisor()
    // vc2.Config.GenesisSignature = vc.Config.GenesisSignature
    // vc2.Config.MasterKeys = vc.Config.MasterKeys
    // vc2.Config.IsMaster = true
    // vc2.Disabled = false
    vc.Config.IsMaster = true
    vc.Disabled = false
    v = NewVisor(vc)
    tx, err := v.Spend(v.Visor.Wallets[0].GetID(),
        visor.Balance{vc.Config.GenesisCoinVolume, 0},
        vc.Config.GenesisCoinVolume, dest.Address, p)
    mv.RecordTxn(tx)
    wait()
    assert.Nil(t, err)
    assert.Equal(t, len(p.Pool.SendResults), 1)
    for len(p.Pool.SendResults) > 0 {
        <-p.Pool.SendResults
    }
    err = v.CreateAndPublishBlock(p)
    assert.Nil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    for len(p.Pool.SendResults) > 0 {
        <-p.Pool.SendResults
    }
    // No coins to spend, fail
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))
    _, err = v.Spend(v.Visor.Wallets[0].GetID(), visor.Balance{10 * 1e6, 0}, 0,
        dest.Address, p)
    assert.NotNil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    err = v.CreateAndPublishBlock(p)
    assert.NotNil(t, err)
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.Equal(t, v.Visor.MostRecentBkSeq(), uint64(1))
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

    // Test serialization
    m = NewGetBlocksMessage(uint64(106))
    b := encoder.Serialize(m)
    m2 := GetBlocksMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestGetBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)
    p := d.Pool.Pool
    go p.ConnectionWriteLoop(gc)
    assert.Nil(t, transferCoins(mv, d.Visor.Visor))
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(1))
    m := NewGetBlocksMessage(uint64(7))
    m.c = messageContext(addr)
    go p.ConnectionWriteLoop(m.c.Conn)
    defer m.c.Conn.Close()

    // Disabled
    d.Visor.Config.Disabled = true
    m.Process(d)
    wait()
    assert.Equal(t, len(p.SendResults), 0)
    // Disabled should not record a blockchain length
    assert.Equal(t, len(d.Visor.blockchainLengths), 0)

    // Enabled handler, should record bc length not not send anything since
    // we have nothing new enough
    d.Visor.Config.Disabled = false
    m.Process(d)
    wait()
    assert.Equal(t, len(p.SendResults), 0)
    assert.Equal(t, len(d.Visor.blockchainLengths), 1)
    assert.Equal(t, d.Visor.blockchainLengths[addr], uint64(7))
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // Working send
    m.LastBlock = uint64(0)
    m.c.Conn.Conn = NewDummyConn(addr)
    m.Process(d)
    wait()
    assert.Equal(t, len(p.SendResults), 1)
    if len(p.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, m.c.Conn)
    _, ok := sr.Message.(*GiveBlocksMessage)
    assert.True(t, ok)
    assert.Equal(t, len(d.Visor.blockchainLengths), 1)
    assert.Equal(t, d.Visor.blockchainLengths[addr], uint64(0))
    assert.False(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())
}

func TestGiveBlocksMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    _, mv := setupVisor()
    blocks := []visor.SignedBlock{mv.GetGenesisBlock()}
    m := NewGiveBlocksMessage(blocks)
    assert.Equal(t, m.Blocks, blocks)
    testSimpleMessageHandler(t, d, m)

    // Test serialization
    bks, err := makeBlocks(mv, 4)
    assert.Nil(t, err)
    blocks = append(blocks, bks...)
    m = NewGiveBlocksMessage(blocks)
    b := encoder.Serialize(m)
    m2 := GiveBlocksMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestGiveBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)
    go d.Pool.Pool.ConnectionWriteLoop(gc)

    blocks, err := makeBlocks(mv, 2)
    assert.Nil(t, err)
    assert.Equal(t, len(blocks), 2)
    m := NewGiveBlocksMessage(blocks)
    m.c = messageContext(addr)

    // Disabled should have nothing happen
    d.Visor.Config.Disabled = true
    m.Process(d)
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(0))
    assert.True(t, gc.LastSent.IsZero())

    // Not disabled and blocks were reannounced
    d.Visor.Config.Disabled = false
    gc.Conn = NewDummyConn(addr)
    assert.Equal(t, len(blocks), 2)
    m = NewGiveBlocksMessage(blocks)
    assert.Equal(t, len(m.Blocks), 2)
    m.c = messageContext(addr)
    m.Process(d)
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
    if len(d.Pool.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-d.Pool.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*AnnounceBlocksMessage)
    assert.True(t, ok)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(2))
    assert.False(t, gc.LastSent.IsZero())

    // Send blocks we have and some we dont, as long as they are in order
    // we can use the ones at the end
    gc.LastSent = util.ZeroTime()
    moreBlocks, err := makeMoreBlocks(mv, 2,
        blocks[len(blocks)-1].Block.Head.Time)
    assert.Nil(t, err)
    blocks = append(blocks, moreBlocks...)
    m = NewGiveBlocksMessage(blocks)
    m.c = messageContext(addr)
    m.Process(d)
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
    if len(d.Pool.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr = <-d.Pool.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, gc)
    _, ok = sr.Message.(*AnnounceBlocksMessage)
    assert.True(t, ok)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(4))
    assert.False(t, gc.LastSent.IsZero())

    // Send invalid blocks
    gc.LastSent = util.ZeroTime()
    bb := visor.SignedBlock{
        Block: coin.Block{
            Head: coin.BlockHeader{
                BkSeq: uint64(7),
            }}}
    m = NewGiveBlocksMessage([]visor.SignedBlock{bb})
    m.c = messageContext(addr)
    m.Process(d)
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(4))
    assert.True(t, gc.LastSent.IsZero())
}

func TestAnnounceBlocksMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    m := NewAnnounceBlocksMessage(uint64(7))
    assert.Equal(t, m.MaxBkSeq, uint64(7))
    testSimpleMessageHandler(t, d, m)

    // Test serialization
    m = NewAnnounceBlocksMessage(uint64(101))
    b := encoder.Serialize(m)
    m2 := AnnounceBlocksMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestAnnounceBlocksMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    p := d.Pool
    gc := setupExistingPool(p)
    go p.Pool.ConnectionWriteLoop(gc)
    defer gc.Close()
    assert.Nil(t, transferCoins(mv, d.Visor.Visor))
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(1))

    // Disabled, nothing should happen
    d.Visor.Config.Disabled = true
    m := NewAnnounceBlocksMessage(uint64(2))
    m.c = messageContext(addr)
    defer m.c.Conn.Close()
    go p.Pool.ConnectionWriteLoop(m.c.Conn)
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // We know all the blocks
    d.Visor.Config.Disabled = false
    m.MaxBkSeq = uint64(1)
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 0)
    assert.True(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // We send a GetBlocksMessage in response to a higher MaxBkSeq
    m.MaxBkSeq = uint64(7)
    assert.False(t, d.Visor.Visor.MostRecentBkSeq() >= m.MaxBkSeq)
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Nil(t, sr.Error)
    assert.Equal(t, sr.Connection, m.c.Conn)
    _, ok := sr.Message.(*GetBlocksMessage)
    assert.True(t, ok)
    assert.False(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())
}

func TestAnnounceTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewAnnounceTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)

    // Test serialization
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn.Hash())
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn.Hash())
    m = NewAnnounceTxnsMessage(txns)
    assert.Equal(t, len(m.Txns), 3)
    b := encoder.Serialize(m)
    m2 := AnnounceTxnsMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestAnnounceTxnsMessageProcess(t *testing.T) {
    v, _ := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)
    go d.Pool.Pool.ConnectionWriteLoop(gc)

    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewAnnounceTxnsMessage(txns)
    m.c = messageContext(addr)
    go d.Pool.Pool.ConnectionWriteLoop(m.c.Conn)
    defer m.c.Conn.Close()

    // Disabled, nothing should happen
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.True(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())

    // We don't know some, request them
    d.Visor.Config.Disabled = false
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
    if len(d.Pool.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-d.Pool.Pool.SendResults
    assert.Equal(t, sr.Connection, m.c.Conn)
    assert.Nil(t, sr.Error)
    _, ok := sr.Message.(*GetTxnsMessage)
    assert.True(t, ok)
    assert.False(t, m.c.Conn.LastSent.IsZero())
    // Should not have been broadcast
    assert.True(t, gc.LastSent.IsZero())

    // We know all the reported txns, nothing should be sent
    d.Visor.Visor.Unconfirmed.Txns[tx.Txn.Hash()] = tx
    m.c.Conn.Conn = NewDummyConn(addr)
    m.c.Conn.LastSent = util.ZeroTime()
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.True(t, m.c.Conn.LastSent.IsZero())
    assert.True(t, gc.LastSent.IsZero())
}

func TestGetTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewGetTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)

    // Test serialization
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn.Hash())
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn.Hash())
    m = NewGetTxnsMessage(txns)
    assert.Equal(t, len(m.Txns), 3)
    b := encoder.Serialize(m)
    m2 := GetTxnsMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestGetTxnsMessageProcess(t *testing.T) {
    v, _ := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)
    p := d.Pool
    go p.Pool.ConnectionWriteLoop(gc)
    tx := createUnconfirmedTxn()
    tx.Txn.Head.Hash = coin.SumSHA256([]byte("asdadwadwada"))
    txns := []coin.SHA256{tx.Txn.Hash()}
    m := NewGetTxnsMessage(txns)
    m.c = messageContext(addr)
    go p.Pool.ConnectionWriteLoop(m.c.Conn)
    defer m.c.Conn.Close()

    // We don't have any to reply with
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // Disabled, nothing should happen
    d.Visor.Visor.Unconfirmed.Txns[tx.Txn.Hash()] = tx
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    assert.True(t, m.c.Conn.LastSent.IsZero())

    // We have some to reply with
    d.Visor.Config.Disabled = false
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(p.Pool.SendResults), 1)
    if len(p.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-p.Pool.SendResults
    assert.Equal(t, sr.Connection, m.c.Conn)
    assert.Nil(t, sr.Error)
    _, ok := sr.Message.(*GiveTxnsMessage)
    assert.True(t, ok)
    assert.False(t, m.c.Conn.LastSent.IsZero())
    // Should not be broadcast to others
    assert.True(t, gc.LastSent.IsZero())
}

func TestGiveTxnsMessageHandle(t *testing.T) {
    d := newDefaultDaemon()
    defer shutdown(d)
    tx := createUnconfirmedTxn()
    txns := coin.Transactions{tx.Txn}
    m := NewGiveTxnsMessage(txns)
    assert.Equal(t, m.Txns, txns)
    testSimpleMessageHandler(t, d, m)

    // Test serialization
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn)
    tx = createUnconfirmedTxn()
    txns = append(txns, tx.Txn)
    m = NewGiveTxnsMessage(txns)
    assert.Equal(t, len(m.Txns), 3)
    b := encoder.Serialize(m)
    m2 := GiveTxnsMessage{}
    assert.Nil(t, encoder.DeserializeRaw(b, &m2))
    assert.Equal(t, *m, m2)
}

func TestGiveTxnsMessageProcess(t *testing.T) {
    v, mv := setupVisor()
    d, _ := newVisorDaemon(v)
    defer shutdown(d)
    gc := setupExistingPool(d.Pool)
    go d.Pool.Pool.ConnectionWriteLoop(gc)

    utx := createUnconfirmedTxn()
    txns := coin.Transactions{utx.Txn}
    m := NewGiveTxnsMessage(txns)
    m.c = messageContext(addr)

    // No valid txns, nothing should be sent
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(mv.Unconfirmed.Txns), 0)
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.True(t, gc.LastSent.IsZero())

    // Disabled, nothing should happen
    tx, err := makeValidTxn(mv)
    assert.Nil(t, err)
    m.Txns = coin.Transactions{tx}
    d.Visor.Config.Disabled = true
    assert.NotPanics(t, func() { m.Process(d) })
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 0)
    assert.Equal(t, len(mv.Unconfirmed.Txns), 0)
    assert.True(t, gc.LastSent.IsZero())

    // A valid txn, we should broadcast. Txn's announce state should be updated
    d.Visor.Config.Disabled = false
    assert.True(t, gc.LastSent.IsZero())
    assert.NotPanics(t, func() { m.Process(d) })
    assert.Equal(t, len(d.Visor.Visor.Unconfirmed.Txns), 1)
    wait()
    assert.Equal(t, len(d.Pool.Pool.SendResults), 1)
    if len(d.Pool.Pool.SendResults) == 0 {
        t.Fatal("SendResults empty, would block")
    }
    sr := <-d.Pool.Pool.SendResults
    assert.Equal(t, sr.Connection, gc)
    _, ok := sr.Message.(*AnnounceTxnsMessage)
    assert.True(t, ok)
    assert.Nil(t, err)
    assert.False(t, gc.LastSent.IsZero())
    _, ok = d.Visor.Visor.Unconfirmed.Txns[tx.Hash()]
    assert.True(t, ok)
}

/* Misc */

func assertSendingTxnsMessageInterface(t *testing.T, i interface{},
    hashes []coin.SHA256, isSending bool) {
    m, ok := i.(SendingTxnsMessage)
    assert.Equal(t, ok, isSending)
    if isSending {
        assert.Equal(t, m.GetTxns(), hashes)
    }
}

func TestSendingTxnsMessageInterface(t *testing.T) {
    hashes := []coin.SHA256{randSHA256(t), randSHA256(t)}

    // GetTxnsMessage should not be a SendingTxnsMessage, it is a request for
    // them
    getx := NewGetTxnsMessage(hashes)
    assertSendingTxnsMessageInterface(t, getx, hashes, false)

    // AnnounceTxnsMessage is a SendingTxnsMessage
    annx := NewAnnounceTxnsMessage(hashes)
    assertSendingTxnsMessageInterface(t, annx, hashes, true)

    // GiveTxnsMessage is a SendingTxnsMessage
    defer cleanupVisor()
    _, v := setupVisor()
    txns := coin.Transactions{
        makeValidTxnNoError(t, v),
        makeValidTxnNoError(t, v),
    }
    givx := NewGiveTxnsMessage(txns)
    assertSendingTxnsMessageInterface(t, givx, txns.Hashes(), true)
}

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
