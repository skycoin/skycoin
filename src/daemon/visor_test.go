package daemon

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/visor"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
    "time"
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
    os.Remove("testblockchain.bin")
    os.Remove("testblockchain.sigs")
    os.Remove("testwallet.json")
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

/* Tests for daemon's loop related to visor */

func testBlockCreationTicker(t *testing.T, vcfg VisorConfig, master bool,
    mv *visor.Visor) {
    defer cleanupVisor()
    vcfg.Config.BlockCreationInterval = 1
    c := NewConfig()
    c.Visor = vcfg
    c.Daemon.DisableNetworking = true
    d := NewDaemon(c)
    if !master {
        // Give the nonmaster some money to spend
        addr := d.Visor.Visor.Wallet.Entries[0].Address
        tx, err := mv.Spend(visor.Balance{10 * 1e6, 0}, 0, addr)
        assert.Nil(t, err)
        mv.RecordTxn(tx, false)
        sb, err := mv.CreateBlock()
        assert.Nil(t, err)
        err = d.Visor.Visor.ExecuteSignedBlock(sb)
        assert.Nil(t, err)
    }
    quit := make(chan int)
    defer closeDaemon(d, quit)
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
    dest := visor.NewWalletEntry()
    _, err := d.Visor.Spend(visor.Balance{10 * 1e6, 0}, 0, dest.Address,
        d.Pool)
    assert.Nil(t, err)
    time.Sleep(time.Millisecond * 1250)
    final := start
    if master {
        final += 1
    }
    assert.Equal(t, d.Visor.Visor.MostRecentBkSeq(), uint64(final))
}

func TestBlockCreationTicker(t *testing.T) {
    vcfg, mv := setupVisor()
    // No blocks should get created if we are not master
    testBlockCreationTicker(t, vcfg, false, mv)
}

func TestBlockCreationTickerMaster(t *testing.T) {
    vcfg := setupMasterVisor()
    // Master should make a block
    testBlockCreationTicker(t, vcfg, true, nil)
}

func TestUnconfirmedRefreshTicker(t *testing.T) {
    vc, _ := setupVisor()
    vc.Config.UnconfirmedRefreshRate = time.Millisecond * 10
    vc.Config.UnconfirmedCheckInterval = time.Nanosecond
    vc.Config.UnconfirmedMaxAge = time.Nanosecond
    d, quit := newVisorDaemon(vc)
    ut := visor.UnconfirmedTxn{}
    ut.Txn = coin.Transaction{}
    ut.Received = time.Now().UTC()
    ut.Checked = ut.Received
    ut.Announced = ut.Received
    d.Visor.Visor.UnconfirmedTxns.Txns[ut.Txn.Header.Hash] = ut
    time.Sleep(time.Millisecond)
    go d.Start(quit)
    time.Sleep(time.Millisecond * 15)
    closeDaemon(d, quit)
    assert.Equal(t, len(d.Visor.Visor.UnconfirmedTxns.Txns), 0)
}

/* Tests for daemon.Visor */
