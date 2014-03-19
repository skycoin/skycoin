package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
    "time"
)

/* Helper functions */

func setupVisorWriting(vc VisorConfig) *Visor {
    vc.WalletDirectory = testWalletFile
    vc.BlockSigsFile = testBlocksigsFile
    vc.BlockchainFile = testBlockchainFile
    return NewVisor(vc)
}

func writeVisorFilesDirect(t *testing.T, v *Visor) {
    assert.Nil(t, v.SaveWallet())
    assert.Nil(t, v.SaveBlockSigs())
    assert.Nil(t, v.SaveBlockchain())
    assertFileExists(t, v.Config.WalletDirectory)
    assertFileExists(t, v.Config.BlockSigsFile)
    assertFileExists(t, v.Config.BlockchainFile)
}

func writeVisorFiles(t *testing.T, vc VisorConfig) *Visor {
    cleanupVisor()
    v := setupVisorWriting(vc)
    writeVisorFilesDirect(t, v)
    return v
}

func newWalletEntry(t *testing.T) WalletEntry {
    we := NewWalletEntry()
    assert.Nil(t, we.Verify())
    return we
}

func setupGenesis(t *testing.T) (WalletEntry, coin.Sig, uint64) {
    we := newWalletEntry(t)
    vc := NewVisorConfig()
    vc.IsMaster = true
    vc.MasterKeys = we
    vc.GenesisSignature = createGenesisSignature(we)
    v := NewVisor(vc)
    we.Secret = coin.SecKey{}
    return we, v.blockSigs.Sigs[0], v.blockchain.Blocks[0].Head.Time
}

func newGenesisConfig(t *testing.T) VisorConfig {
    refvc := NewVisorConfig()
    we, sig, ts := setupGenesis(t)
    refvc.MasterKeys = we
    refvc.GenesisSignature = sig
    refvc.GenesisTimestamp = ts
    refvc.IsMaster = false
    return refvc
}

func corruptFile(t *testing.T, filename string) {
    f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0600)
    assert.Nil(t, err)
    _, err = f.Write([]byte("xxxxx"))
    assert.Nil(t, err)
    f.Close()
}

func setupChildVisorConfig(refvc VisorConfig, master bool) VisorConfig {
    vc := NewVisorConfig()
    vc.IsMaster = master
    vc.MasterKeys = refvc.MasterKeys
    vc.WalletDirectory = testWalletFile
    vc.BlockchainFile = testBlockchainFile
    vc.BlockSigsFile = testBlocksigsFile
    return vc
}

func newMasterVisorConfig(t *testing.T) VisorConfig {
    vc := NewVisorConfig()
    vc.CoinHourBurnFactor = 0
    mw := newWalletEntry(t)
    vc.MasterKeys = mw
    vc.GenesisSignature = createGenesisSignature(mw)
    vc.IsMaster = true
    return vc
}

func addValidTxns(t *testing.T, v *Visor, n int) coin.Transactions {
    txns := make(coin.Transactions, n)
    for i := 0; i < len(txns); i++ {
        txn, err := makeValidTxn(v)
        assert.Nil(t, err)
        txns[i] = txn
    }
    for _, txn := range txns {
        err, known := v.RecordTxn(txn)
        assert.Nil(t, err)
        assert.False(t, known)
    }
    txns = coin.SortTransactions(txns, getFee)
    assert.Equal(t, len(v.Unconfirmed.Txns), n)
    return txns
}

func addSignedBlockAt(t *testing.T, v *Visor, when uint64) SignedBlock {
    we := NewWalletEntry()
    tx, err := v.Spend(Balance{1e6, 0}, 0, we.Address)
    assert.Nil(t, err)
    err, known := v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    sb, err := v.CreateBlock(when)
    assert.Nil(t, err)
    if err != nil {
        return sb
    }
    err = v.ExecuteSignedBlock(sb)
    assert.Nil(t, err)
    return sb
}

func addSignedBlock(t *testing.T, v *Visor) SignedBlock {
    return addSignedBlockAt(t, v, uint64(util.UnixNow()))
}

func addSignedBlocks(t *testing.T, v *Visor, n int) []SignedBlock {
    sbs := make([]SignedBlock, n)
    now := uint64(util.UnixNow())
    for i := 0; i < n; i++ {
        sbs[i] = addSignedBlockAt(t, v, now+1+uint64(i))
    }
    return sbs
}

func assertSignedBlocks(t *testing.T, v *Visor, sbs []SignedBlock,
    start, ct uint64) {
    have := v.MostRecentBkSeq()
    if have <= start {
        assert.Equal(t, len(sbs), 0)
    } else if have-start < ct {
        assert.Equal(t, len(sbs), int(have-start))
    } else {
        assert.Equal(t, len(sbs), int(ct))
    }
    for i, sb := range sbs {
        assert.Nil(t, v.verifySignedBlock(&sb))
        assert.Equal(t, sb.Sig, v.blockSigs.Sigs[sb.Block.Head.BkSeq])
        assert.Equal(t, sb.Block.Head.BkSeq, start+uint64(i)+1)
        assert.Equal(t, v.blockchain.Blocks[start+uint64(i)+1], sb.Block)
    }
}

func assertReadableBlocks(t *testing.T, v *Visor, rbs []ReadableBlock,
    sbs []SignedBlock) {
    assert.Equal(t, len(rbs), len(sbs))
    for i, rb := range rbs {
        assertReadableBlock(t, rb, sbs[i].Block)
    }
}

func assertBlocks(t *testing.T, v *Visor, bs []coin.Block, sbs []SignedBlock) {
    assert.Equal(t, len(bs), len(sbs))
    for i, b := range bs {
        assert.Equal(t, b, sbs[i].Block)
    }
}

/* Actual tests */

func TestNewVisorConfig(t *testing.T) {
    vc := NewVisorConfig()
    assert.False(t, vc.IsMaster)
    assert.True(t, vc.CanSpend)
    assert.Equal(t, vc.WalletDirectory, "")
    assert.Equal(t, vc.BlockchainFile, "")
    assert.Equal(t, vc.BlockSigsFile, "")
    assert.Panics(t, func() { vc.MasterKeys.Verify() })
    assert.NotNil(t, vc.MasterKeys.VerifyPublic())
    assert.Equal(t, vc.GenesisSignature, coin.Sig{})
    assert.Equal(t, vc.WalletSizeMin, 1)
}

func TestNewVisor(t *testing.T) {
    defer cleanupVisor()

    // Not master, Invalid master keys
    cleanupVisor()
    we := NewWalletEntry()
    we.Public = coin.PubKey{}
    vc := NewVisorConfig()
    vc.IsMaster = false
    assert.Panics(t, func() { NewVisor(vc) })
    vc.MasterKeys = we
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, invalid master keys
    cleanupVisor()
    vc.IsMaster = true
    vc.MasterKeys = WalletEntry{}
    assert.Panics(t, func() { NewVisor(vc) })
    vc.MasterKeys = we
    assert.Panics(t, func() { NewVisor(vc) })

    // Not master, no wallet, blockchain, blocksigs file
    cleanupVisor()
    vc = NewVisorConfig()
    vc.IsMaster = false
    we, sig, ts := setupGenesis(t)
    vc.MasterKeys = we
    vc.GenesisSignature = sig
    vc.GenesisTimestamp = ts
    vc.WalletSizeMin = 10
    v := NewVisor(vc)
    assert.Equal(t, v.Wallet.NumEntries(), 10)
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)

    // Master, no wallet, blockchain, blocksigs file
    cleanupVisor()
    vc = NewVisorConfig()
    we = newWalletEntry(t)
    vc.MasterKeys = we
    vc.GenesisSignature = createGenesisSignature(we)
    vc.WalletSizeMin = 10
    vc.IsMaster = true
    v = NewVisor(vc)
    // Wallet should only have 1 entry if master
    assert.Equal(t, v.Wallet.NumEntries(), 1)
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)

    // Not master, has all files
    cleanupVisor()
    refvc := newGenesisConfig(t)
    refv := writeVisorFiles(t, refvc)
    vc = setupChildVisorConfig(refvc, false)
    v = NewVisor(vc)
    assert.Equal(t, v.Wallet, refv.Wallet)
    assert.Equal(t, v.blockchain, refv.blockchain)
    assert.Equal(t, v.blockSigs, refv.blockSigs)

    // Master, has all files
    cleanupVisor()
    refvc = newMasterVisorConfig(t)
    refv = writeVisorFiles(t, refvc)
    vc = setupChildVisorConfig(refvc, true)
    v = NewVisor(vc)
    assert.Equal(t, v.Wallet, refv.Wallet)
    assert.Equal(t, v.blockchain, refv.blockchain)

    // Not master, wallet is corrupt
    cleanupVisor()
    refvc = newGenesisConfig(t)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testWalletFile)
    vc = setupChildVisorConfig(refvc, false)
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, wallet is corrupt.  Nothing happens because master ignores
    // wallet
    cleanupVisor()
    refvc = newMasterVisorConfig(t)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testWalletFile)
    vc = setupChildVisorConfig(refvc, true)
    assert.NotPanics(t, func() { NewVisor(vc) })

    // Not master, blocksigs is corrupt
    cleanupVisor()
    refvc = newGenesisConfig(t)
    assertFileNotExists(t, testWalletFile)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testBlocksigsFile)
    vc = setupChildVisorConfig(refvc, false)
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, blocksigs is corrupt
    cleanupVisor()
    refvc = newMasterVisorConfig(t)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testBlocksigsFile)
    assertFileExists(t, testBlocksigsFile)
    vc = setupChildVisorConfig(refvc, true)
    assert.Panics(t, func() { NewVisor(vc) })

    // Not master, blockchain is corrupt
    cleanupVisor()
    refvc = newGenesisConfig(t)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testBlockchainFile)
    vc = setupChildVisorConfig(refvc, false)
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, blockchain is corrupt
    cleanupVisor()
    refvc = newMasterVisorConfig(t)
    refv = writeVisorFiles(t, refvc)
    corruptFile(t, testBlockchainFile)
    vc = setupChildVisorConfig(refvc, true)
    assert.Panics(t, func() { NewVisor(vc) })

    // Not master, blocksigs is not valid for blockchain
    cleanupVisor()
    refvc = newGenesisConfig(t)
    refv = setupVisorWriting(refvc)
    // Corrupt the signature
    refv.blockSigs.Sigs[uint64(0)] = coin.Sig{}
    writeVisorFilesDirect(t, refv)
    vc = setupChildVisorConfig(refvc, false)
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, blocksigs is not valid for blockchain
    cleanupVisor()
    refvc = newMasterVisorConfig(t)
    refv = setupVisorWriting(refvc)
    // Corrupt the signature
    refv.blockSigs.Sigs[uint64(0)] = coin.Sig{}
    writeVisorFilesDirect(t, refv)
    vc = setupChildVisorConfig(refvc, true)
    assert.Panics(t, func() { NewVisor(vc) })
}

func TestNewMinimalVisor(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    vc.WalletSizeMin = 10000
    v := NewMinimalVisor(vc)
    assert.Equal(t, v.Config, vc)
    assert.NotNil(t, v.Unconfirmed)
    assert.Nil(t, v.Wallet)
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
}

func TestCreateGenesisBlock(t *testing.T) {
    defer cleanupVisor()
    // Test as master, successful
    vc := newMasterVisorConfig(t)
    v := NewMinimalVisor(vc)
    assert.True(t, v.Config.IsMaster)
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
    sb := v.CreateGenesisBlock()
    assert.NotEqual(t, sb.Block, coin.Block{})
    assert.NotEqual(t, sb.Sig, coin.Sig{})
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Nil(t, v.blockSigs.Verify(vc.MasterKeys.Public, v.blockchain))

    // Test as not master, successful
    vc = newGenesisConfig(t)
    v = NewMinimalVisor(vc)
    assert.False(t, v.Config.IsMaster)
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
    sb = v.CreateGenesisBlock()
    assert.NotEqual(t, sb.Block, coin.Block{})
    assert.NotEqual(t, sb.Sig, coin.Sig{})
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Nil(t, v.blockSigs.Verify(vc.MasterKeys.Public, v.blockchain))
    assert.Equal(t, v.Config.GenesisSignature, sb.Sig)
    assert.Equal(t, v.blockchain.Blocks[0].Head.Time, v.Config.GenesisTimestamp)

    // Test as master, blockSigs invalid for pubkey
    vc = newMasterVisorConfig(t)
    vc.MasterKeys.Public = coin.PubKey{}
    v = NewMinimalVisor(vc)
    assert.True(t, v.Config.IsMaster)
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
    assert.Panics(t, func() { v.CreateGenesisBlock() })

    // Test as not master, blockSigs invalid for pubkey
    vc = newGenesisConfig(t)
    vc.MasterKeys.Public = coin.PubKey{}
    v = NewMinimalVisor(vc)
    assert.False(t, v.Config.IsMaster)
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
    assert.Panics(t, func() { v.CreateGenesisBlock() })

    // Test as master, signing failed
    vc = newMasterVisorConfig(t)
    vc.MasterKeys.Secret = coin.SecKey{}
    vc.GenesisSignature = coin.Sig{}
    assert.Equal(t, vc.MasterKeys.Secret, coin.SecKey{})
    v = NewMinimalVisor(vc)
    assert.True(t, v.Config.IsMaster)
    assert.Equal(t, v.Config, vc)
    assert.Equal(t, v.Config.MasterKeys.Secret, coin.SecKey{})
    assert.Equal(t, len(v.blockchain.Blocks), 0)
    assert.Equal(t, len(v.blockSigs.Sigs), 0)
    assert.Panics(t, func() { v.CreateGenesisBlock() })
}

func TestVisorRefreshUnconfirmed(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    testRefresh(t, mv, func(checkPeriod, maxAge time.Duration) {
        mv.Config.UnconfirmedCheckInterval = checkPeriod
        mv.Config.UnconfirmedMaxAge = maxAge
        mv.RefreshUnconfirmed()
    })
}

func TestVisorSaveBlockchain(t *testing.T) {
    cleanupVisor()
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    vc.BlockchainFile = ""

    // Test with no blockchain file set
    v := NewVisor(vc)
    assertFileNotExists(t, testBlockchainFile)
    err := v.SaveBlockchain()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No BlockchainFile location set")
    assertFileNotExists(t, testBlockchainFile)

    // Test with blockchain file set
    vc.BlockchainFile = testBlockchainFile
    v = NewVisor(vc)
    assert.Nil(t, v.SaveBlockchain())
    assertFileExists(t, testBlockchainFile)
    assert.NotPanics(t, func() {
        loadBlockchain(testBlockchainFile, vc.MasterKeys.Address)
    })
    bc := loadBlockchain(testBlockchainFile, vc.MasterKeys.Address)
    assert.Equal(t, v.blockchain, bc)

}

func TestVisorSaveWallet(t *testing.T) {
    cleanupVisor()
    defer cleanupVisor()
    vc := newGenesisConfig(t)
    vc.WalletDirectory = ""
    vc.WalletSizeMin = 10
    assert.False(t, vc.IsMaster)

    // Test with no wallet file set
    v := NewVisor(vc)
    assertFileNotExists(t, testWalletFile)
    err := v.SaveWallet()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No WalletDirectory location set")
    assertFileNotExists(t, testWalletFile)

    // Test with wallet file set
    vc.WalletDirectory = testWalletFile
    v = NewVisor(vc)
    assert.Nil(t, v.SaveWallet())
    assertFileExists(t, testWalletFile)
    w, err := LoadSimpleWallet(testWalletFile)
    assert.Nil(t, err)
    assert.Equal(t, v.Wallet, w)
}

func TestVisorSaveBlockSigs(t *testing.T) {
    cleanupVisor()
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    vc.BlockSigsFile = ""

    // Test with no blocksigs file set
    v := NewVisor(vc)
    assertFileNotExists(t, testBlocksigsFile)
    err := v.SaveBlockSigs()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No BlockSigsFile location set")
    assertFileNotExists(t, testBlocksigsFile)

    vc.BlockSigsFile = testBlocksigsFile
    v = NewVisor(vc)
    assert.Nil(t, v.SaveBlockSigs())
    assertFileExists(t, testBlocksigsFile)

    bs, err := LoadBlockSigs(testBlocksigsFile)
    assert.Nil(t, err)
    assert.Equal(t, v.blockSigs, bs)
}

func TestCreateAddressAndSave(t *testing.T) {
    cleanupVisor()
    defer cleanupVisor()
    vc := newGenesisConfig(t)
    vc.WalletDirectory = ""
    vc.WalletSizeMin = 10

    // Test with no wallet file set
    v := NewVisor(vc)
    assert.Equal(t, v.Wallet.NumEntries(), 10)
    we, err := v.CreateAddressAndSave()
    assert.NotNil(t, err)
    assertFileNotExists(t, testWalletFile)
    assert.Equal(t, v.Wallet.NumEntries(), 11)
    assert.Nil(t, we.Verify())

    // Test with wallet file set
    v.Config.WalletDirectory = testWalletFile
    we, err = v.CreateAddressAndSave()
    assert.Nil(t, err)
    assertFileExists(t, testWalletFile)
    assert.Equal(t, v.Wallet.NumEntries(), 12)
    assert.Nil(t, we.Verify())

    w, err := LoadSimpleWallet(testWalletFile)
    assert.Nil(t, err)
    assert.Equal(t, v.Wallet, w)
}

func TestCreateAndExecuteBlock(t *testing.T) {
    defer cleanupVisor()

    // Test as not master, should fail
    vc := newGenesisConfig(t)
    v := NewVisor(vc)
    assert.Panics(t, func() { v.CreateAndExecuteBlock() })

    // Test as master, no txns
    vc = newMasterVisorConfig(t)
    v = NewVisor(vc)
    _, err := v.CreateAndExecuteBlock()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "No transactions")

    // Test as master, more txns than allowed
    vc.BlockCreationInterval = uint64(101)
    v = NewVisor(vc)
    txns := addValidTxns(t, v, 3)
    txns = coin.SortTransactions(txns, v.blockchain.TransactionFee)
    v.Config.MaxBlockSize = txns[0].Size()
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    sb, err := v.CreateAndExecuteBlock()
    assert.Nil(t, err)

    assert.Equal(t, len(sb.Block.Body.Transactions), 1)
    assert.Equal(t, sb.Block.Body.Transactions[0], txns[0])
    assert.Equal(t, len(v.blockchain.Blocks), 2)
    assert.Equal(t, len(v.blockSigs.Sigs), 2)
    assert.Equal(t, v.blockchain.Blocks[1], sb.Block)
    assert.Equal(t, v.blockSigs.Sigs[1], sb.Sig)
    assert.Equal(t, len(v.Unconfirmed.Txns), 2)
    assert.True(t, sb.Block.Head.Time > v.blockchain.Blocks[0].Head.Time)
    rawTxns := v.Unconfirmed.RawTxns()
    assert.Equal(t, len(rawTxns), 2)
    for _, tx := range sb.Block.Body.Transactions {
        assert.NotEqual(t, tx.Hash(), rawTxns[0].Hash())
        assert.NotEqual(t, tx.Hash(), rawTxns[1].Hash())
    }
    if txns[1].Hash() == rawTxns[0].Hash() {
        assert.Equal(t, txns[2].Hash(), rawTxns[1].Hash())
    } else {
        assert.Equal(t, txns[2].Hash(), rawTxns[0].Hash())
    }
    assert.Nil(t, v.blockSigs.Verify(v.Config.MasterKeys.Public, v.blockchain))

    // No txns, forcing NewBlockFromTransactions to fail
    v = NewVisor(vc)
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    txns = addValidTxns(t, v, 3)
    v.Config.MaxBlockSize = 0
    sb, err = v.CreateAndExecuteBlock()
    assert.NotNil(t, err)
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Equal(t, len(v.Unconfirmed.Txns), 3)
}

func TestVisorSpend(t *testing.T) {
    defer cleanupVisor()
    we := NewWalletEntry()
    addr := we.Address
    vc := newMasterVisorConfig(t)
    assert.Equal(t, vc.CoinHourBurnFactor, uint64(0))
    v := NewVisor(vc)
    ogb := v.TotalBalance()

    // Test can't spend
    v = NewVisor(vc)
    b := Balance{10e6, 0}
    v.Config.CanSpend = false
    _, err := v.Spend(b, 0, addr)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Spending disabled")

    // Test spend 0 amount
    v = NewVisor(vc)
    b = Balance{0, 0}
    _, err = v.Spend(b, 0, addr)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Zero spend amount")

    // Test lacking funds
    v = NewVisor(vc)
    b = Balance{10e16, 10e16}
    _, err = v.Spend(b, 10e16, addr)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Not enough coins")

    // Test created txn too large
    v = NewVisor(vc)
    v.Config.MaxBlockSize = 0
    b = Balance{10e6, 10}
    assert.Panics(t, func() { v.Spend(b, 0, addr) })

    // Test simple spend (we have only 1 address to spend from, no fee)
    v = NewVisor(vc)
    assert.Equal(t, v.Config.CoinHourBurnFactor, uint64(0))
    b = Balance{10e6, 10}
    tx, err := v.Spend(b, 0, addr)
    assert.Nil(t, err)
    assert.Equal(t, len(tx.In), 1)
    assert.Equal(t, len(tx.Out), 2)
    // Hash should be updated
    assert.NotEqual(t, tx.Head.Hash, coin.SHA256{})
    // Should be 1 signature for the single input
    assert.Equal(t, len(tx.Head.Sigs), 1)
    // Spent amount should be correct
    assert.Equal(t, tx.Out[1].Address, addr)
    assert.Equal(t, tx.Out[1].Coins, b.Coins)
    assert.Equal(t, tx.Out[1].Hours, b.Hours)
    // Change amount should be correct
    ourAddr := v.Wallet.GetAddresses()[0]
    assert.Equal(t, tx.Out[0].Address, ourAddr)
    assert.Equal(t, tx.Out[0].Coins, ogb.Coins-b.Coins)
    assert.Equal(t, tx.Out[0].Hours, ogb.Hours-b.Hours)
    assert.Nil(t, tx.Verify())
}

func TestExecuteSignedBlock(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()
    we := NewWalletEntry()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    tx, err := v.Spend(Balance{1e6, 0}, 0, we.Address)
    assert.Nil(t, err)
    err, known := v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    now := uint64(util.UnixNow())

    // Invalid signed block
    sb, err := v.CreateBlock(now)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Nil(t, err)
    sb.Sig = coin.Sig{}
    err = v.ExecuteSignedBlock(sb)
    assert.NotNil(t, err)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)

    // Invalid block
    sb, err = v.CreateBlock(now)
    assert.Nil(t, err)
    // TODO -- empty BodyHash is being accepted, fix blockchain verification
    sb.Block.Head.BodyHash = coin.SHA256{}
    sb.Block.Body.Transactions = make(coin.Transactions, 0)
    sb = v.SignBlock(sb.Block)
    err = v.ExecuteSignedBlock(sb)
    assert.NotNil(t, err)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)

    // Valid block
    sb, err = v.CreateBlock(now)
    assert.Nil(t, err)
    err = v.ExecuteSignedBlock(sb)
    assert.Nil(t, err)
    assert.Equal(t, len(v.blockSigs.Sigs), 2)
    assert.Equal(t, v.blockSigs.Sigs[uint64(1)], sb.Sig)
    assert.Equal(t, v.blockchain.Blocks[1], sb.Block)
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)

    // Test a valid block created by a master but executing in non master
    vc2, mv := setupVisorConfig()
    v2 := NewVisor(vc2)
    addr := v2.Wallet.GetAddresses()[0]
    tx, err = mv.Spend(Balance{1e6, 0}, 0, addr)
    assert.Nil(t, err)
    err, known = mv.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    sb, err = mv.CreateAndExecuteBlock()
    assert.Nil(t, err)
    err = v2.ExecuteSignedBlock(sb)
    assert.Nil(t, err)
    assert.Equal(t, len(v2.blockSigs.Sigs), 2)
    assert.Equal(t, v2.blockSigs.Sigs[uint64(1)], sb.Sig)
    assert.Equal(t, v2.blockchain.Blocks[1], sb.Block)
    assert.Equal(t, len(v2.Unconfirmed.Txns), 0)
}

func TestGetSignedBlocksSince(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    // No blocks
    sbs := v.GetSignedBlocksSince(0, 10)
    assert.Equal(t, len(sbs), 0)

    // All available blocks
    addSignedBlocks(t, v, 10)
    sbs = v.GetSignedBlocksSince(2, 4)
    assertSignedBlocks(t, v, sbs, 2, 4)

    // No available blocks
    sbs = v.GetSignedBlocksSince(100, 20)
    assert.Equal(t, len(sbs), 0)

    // Some, but not all
    sbs = v.GetSignedBlocksSince(7, 5)
    assertSignedBlocks(t, v, sbs, 7, 5)
}

func TestGetGenesisBlock(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)

    // Panics with no signed genesis block
    v := NewMinimalVisor(vc)
    assert.Panics(t, func() { v.GetGenesisBlock() })

    // Panics with no blocks
    v = NewMinimalVisor(vc)
    v.blockSigs.Sigs[0] = coin.Sig{}
    assert.Panics(t, func() { v.GetGenesisBlock() })

    // Correct result
    v = NewVisor(vc)
    gb := v.GetGenesisBlock()
    assert.Equal(t, v.blockSigs.Sigs[0], gb.Sig)
    assert.Equal(t, v.blockchain.Blocks[0], gb.Block)
}

func TestMostRecentBkSeq(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)
    assert.Equal(t, v.MostRecentBkSeq(), uint64(0))
    addSignedBlocks(t, v, 10)
    assert.Equal(t, v.MostRecentBkSeq(), uint64(10))
    v = NewMinimalVisor(vc)
    assert.Panics(t, func() { v.MostRecentBkSeq() })
}

func TestGetBlockchainMetadata(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)
    addSignedBlocks(t, v, 8)
    addUnconfirmedTxn(v)
    addUnconfirmedTxn(v)
    bcm := v.GetBlockchainMetadata()
    assert.Equal(t, bcm.Unspents, uint64(9))
    assert.Equal(t, bcm.Unconfirmed, uint64(2))
    assertReadableBlockHeader(t, bcm.Head, v.blockchain.Head().Head)
}

func TestGetReadableBlock(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    rb, err := v.GetReadableBlock(1)
    assert.NotNil(t, err)
    sb := addSignedBlock(t, v)
    rb, err = v.GetReadableBlock(1)
    assert.Nil(t, err)
    assertReadableBlock(t, rb, sb.Block)
}

func TestGetReadableBlocks(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    rbs := v.GetReadableBlocks(1, 10)
    assert.Equal(t, len(rbs), 0)
    rbs = v.GetReadableBlocks(0, 10)
    sbs := []SignedBlock{SignedBlock{
        Sig:   v.blockSigs.Sigs[0],
        Block: v.blockchain.Blocks[0],
    }}
    assertReadableBlocks(t, v, rbs, sbs)
    sbs = append(sbs, addSignedBlocks(t, v, 5)...)
    rbs = v.GetReadableBlocks(0, 10)
    assertReadableBlocks(t, v, rbs, sbs)
    rbs = v.GetReadableBlocks(2, 4)
    sbs = sbs[2:4]
    assertReadableBlocks(t, v, rbs, sbs)
}

func TestGetBlock(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    b, err := v.GetBlock(1)
    assert.NotNil(t, err)
    sb := addSignedBlock(t, v)
    b, err = v.GetBlock(1)
    assert.Nil(t, err)
    assert.Equal(t, b, sb.Block)
}

func TestGetBlocks(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    bs := v.GetBlocks(1, 10)
    assert.Equal(t, len(bs), 0)
    bs = v.GetBlocks(0, 10)
    sbs := []SignedBlock{SignedBlock{
        Sig:   v.blockSigs.Sigs[0],
        Block: v.blockchain.Blocks[0],
    }}
    assertBlocks(t, v, bs, sbs)
    sbs = append(sbs, addSignedBlocks(t, v, 5)...)
    bs = v.GetBlocks(0, 10)
    assertBlocks(t, v, bs, sbs)
    bs = v.GetBlocks(2, 4)
    sbs = sbs[2:4]
    assertBlocks(t, v, bs, sbs)
}

func TestVisorSetAnnounced(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    now := util.Now()
    utx := addUnconfirmedTxn(v)
    assert.True(t, utx.Announced.IsZero())
    assert.True(t, v.Unconfirmed.Txns[utx.Hash()].Announced.IsZero())
    v.SetAnnounced(utx.Hash(), now)
    assert.False(t, v.Unconfirmed.Txns[utx.Hash()].Announced.IsZero())
    assert.Equal(t, v.Unconfirmed.Txns[utx.Hash()].Announced, now)
}

func TestVisorRecordTxn(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    // Setup txns
    tx, err := makeValidTxn(v)
    assert.Nil(t, err)
    we := v.Wallet.CreateEntry()
    tx2, err := v.Spend(Balance{1e6, 0}, 0, we.Address)
    assert.Nil(t, err)

    // Valid record, did not announce
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    err, known := v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())

    // Invalid txn
    tx.Out = make([]coin.TransactionOutput, 0)
    err, known = v.RecordTxn(tx)
    assert.NotNil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())

    // Make sure isOurSpend and isOurReceive is correct
    tx = tx2
    err, known = v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 2)
    assert.True(t, v.Unconfirmed.Txns[tx.Hash()].Announced.IsZero())
    assert.True(t, v.Unconfirmed.Txns[tx.Hash()].IsOurReceive)
    assert.True(t, v.Unconfirmed.Txns[tx.Hash()].IsOurSpend)
}

func TestGetAddressTransactions(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    // An confirmed txn
    we := v.Wallet.CreateEntry()
    tx, err := v.Spend(Balance{1e6, 0}, 0, we.Address)
    assert.Nil(t, err)
    err, known := v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    _, err = v.CreateAndExecuteBlock()
    assert.Nil(t, err)
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    txns := v.GetAddressTransactions(we.Address)
    assert.Equal(t, len(txns), 1)
    assert.Equal(t, txns[0].Txn, tx)
    assert.True(t, txns[0].Status.Confirmed)
    assert.Equal(t, txns[0].Status.Height, uint64(1))

    // An unconfirmed txn
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    assert.Equal(t, len(v.Unconfirmed.Unspent), 0)
    we = v.Wallet.CreateEntry()
    tx, err = v.Spend(Balance{2e6, 0}, 0, we.Address)
    err, known = v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    assert.Equal(t, len(v.Unconfirmed.Txns), 1)
    assert.Equal(t, len(v.Unconfirmed.Unspent), 1)
    assert.Equal(t, len(v.Unconfirmed.Unspent[tx.Hash()]), 2)
    found := false
    for _, uxs := range v.Unconfirmed.Unspent {
        if found {
            break
        }
        for _, ux := range uxs {
            if ux.Body.Address == we.Address {
                found = true
                break
            }
        }
    }
    auxs := v.Unconfirmed.Unspent.AllForAddress(we.Address)
    assert.Equal(t, len(auxs), 1)
    assert.True(t, found)
    txns = v.GetAddressTransactions(we.Address)
    assert.Equal(t, len(txns), 1)
    assert.Equal(t, txns[0].Txn, tx)
    assert.True(t, txns[0].Status.Unconfirmed)

    // An unconfirmed txn, but pool is corrupted
    assert.True(t, len(v.Unconfirmed.Unspent) > 0)
    ux := coin.UxOut{}
    found = false
    for _, uxs := range v.Unconfirmed.Unspent {
        if len(uxs) > 0 {
            ux = uxs[0]
            found = true
            break
        }
    }
    assert.True(t, found)
    srcTxn := ux.Body.SrcTransaction
    delete(v.Unconfirmed.Txns, srcTxn)
    txns = v.GetAddressTransactions(we.Address)
    assert.Equal(t, len(txns), 0)
}

func TestGetTransaction(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    // Unknown
    tx, err := makeValidTxn(v)
    assert.Nil(t, err)
    tx2 := v.GetTransaction(tx.Hash())
    assert.True(t, tx2.Status.Unknown)

    // Unconfirmed
    err, known := v.RecordTxn(tx)
    assert.Nil(t, err)
    assert.False(t, known)
    tx2 = v.GetTransaction(tx.Hash())
    assert.True(t, tx2.Status.Unconfirmed)
    assert.Equal(t, tx, tx2.Txn)

    // Confirmed
    _, err = v.CreateAndExecuteBlock()
    assert.Nil(t, err)
    assert.Equal(t, len(v.Unconfirmed.Txns), 0)
    tx2 = v.GetTransaction(tx.Hash())
    assert.True(t, tx2.Status.Confirmed)
    assert.Equal(t, tx2.Status.Height, uint64(1))
    assert.Equal(t, tx, tx2.Txn)
}

func TestBalances(t *testing.T) {
    defer cleanupVisor()
    v, mv := setupVisor()
    we := v.Wallet.CreateEntry()
    we2 := v.Wallet.CreateEntry()
    startCoins := mv.Config.GenesisCoinVolume

    // Without predicted outputs
    assert.Nil(t, transferCoinsAdvanced(mv, v, Balance{10e6, 10}, 0, we.Address))
    assert.Nil(t, transferCoinsAdvanced(mv, v, Balance{10e6, 10}, 0, we.Address))
    assert.Nil(t, transferCoinsAdvanced(mv, v, Balance{5e6, 5}, 0, we2.Address))
    assert.Equal(t, v.TotalBalance(), Balance{25e6, 25})
    // assert.Equal(t, v.TotalBalancePredicted(), Balance{25e6, 25})
    mvBalance := Balance{startCoins - 25e6, startCoins - 25}
    assert.Equal(t, mv.TotalBalance(), mvBalance)
    // assert.Equal(t, mv.TotalBalancePredicted(), mvBalance)
    assert.Equal(t, v.Balance(we.Address), Balance{20e6, 20})
    assert.Equal(t, v.Balance(we2.Address), Balance{5e6, 5})
    // assert.Equal(t, v.BalancePredicted(we.Address), Balance{20e6, 20})
    // assert.Equal(t, v.BalancePredicted(we2.Address), Balance{5e6, 5})

    // // With predicted outputs
    // tx, err := mv.Spend(Balance{1e6, 10}, 0, we.Address)
    // assert.Nil(t, err)
    // v.RecordTxn(tx)
    // assert.Nil(t, err)
    // assert.False(t, known)
    // assert.Nil(t, mv.RecordTxn(tx))
    // // Regular balance should not have changed
    // assert.Equal(t, v.TotalBalance(), Balance{25e6, 25})
    // assert.Equal(t, v.TotalBalancePredicted(), Balance{35e6, 35})
    // assert.Equal(t, mv.TotalBalance(), mvBalance)
    // mvBalancePredicted := mvBalance.Sub(Balance{1e6, 10})
    // assert.Equal(t, mv.TotalBalancePredicted(), mvBalancePredicted)
    // assert.Equal(t, v.Balance(we.Address), Balance{20e6, 20})
    // assert.Equal(t, v.Balance(we2.Address), Balance{5e6, 5})
    // assert.Equal(t, v.BalancePredicted(we.Address), Balance{30e6, 30})
    // assert.Equal(t, v.BalancePredicted(we2.Address), Balance{5e6, 5})

    // Can't check CoinHours() is actually used because we'd have to wait
    // an hour
}

func TestVisorVerifySignedBlock(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)
    we := v.Wallet.CreateEntry()

    // Master should verify its own blocks correctly
    txn, err := v.Spend(Balance{1e6, 0}, 0, we.Address)
    assert.Nil(t, err)
    err, known := v.RecordTxn(txn)
    assert.Nil(t, err)
    assert.False(t, known)
    b, err := v.CreateBlock(uint64(util.UnixNow()))
    assert.Nil(t, err)
    assert.Nil(t, v.verifySignedBlock(&b))
    badb := b
    badb.Sig = coin.Sig{}
    assert.NotNil(t, v.verifySignedBlock(&badb))

    // Non master should verify signed blocks generated by master
    mv := v
    v = setupVisorFromMaster(mv)
    assert.Nil(t, v.verifySignedBlock(&b))
    assert.NotNil(t, v.verifySignedBlock(&badb))
}

func TestVisorSignBlock(t *testing.T) {
    defer cleanupVisor()
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)

    // Non master should panic
    b := v.blockchain.Blocks[0]
    v.Config.IsMaster = false
    assert.Panics(t, func() { v.SignBlock(b) })

    // Master should generate valid signed block
    v.Config.IsMaster = true
    sb := v.SignBlock(b)
    assert.Nil(t, v.verifySignedBlock(&sb))
}

func TestLoadWallet(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()

    // Test with no filename (not saving or loading)
    w := loadSimpleWallet("", 20)
    assert.Equal(t, w.NumEntries(), 20)
    assertFileNotExists(t, testWalletFile)

    // Test with filename, file does not exist for loading
    w = loadSimpleWallet(testWalletFile, 20)
    assert.Equal(t, w.NumEntries(), 20)
    assertFileExists(t, testWalletFile)

    // Test with filename, file exists for loading
    w2 := loadSimpleWallet(testWalletFile, 20)
    assert.Equal(t, w2.NumEntries(), 20)
    assert.Equal(t, w, w2)
    assertFileExists(t, testWalletFile)

    // Test with filename, file exists for loading, and we Populate more
    w2 = loadSimpleWallet(testWalletFile, 30)
    assert.Equal(t, w2.NumEntries(), 30)
    for a, we := range w.Entries {
        we2, ok := w2.Entries[a]
        assert.True(t, ok)
        assert.Equal(t, we, we2)
    }
    assertFileExists(t, testWalletFile)

    // Test with filename, file exists for loading, and we Populate less
    w = loadSimpleWallet(testWalletFile, 10)
    assert.Equal(t, w.NumEntries(), 30)
    assert.Equal(t, w, w2)
    assertFileExists(t, testWalletFile)

    // Test with corrupted wallet file for loading
    corruptFile(t, testWalletFile)
    assert.Panics(t, func() { loadSimpleWallet(testWalletFile, 10) })

    // Can't test saving failure since can't force save failure
}

func TestCreateMasterWallet(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()
    we := NewWalletEntry()
    w := CreateMasterWallet(we)
    assert.Equal(t, w.NumEntries(), 1)
    assert.Equal(t, w.GetAddresses()[0], we.Address)

    // Having a wallet file present should not affect loading master wallet
    w.Save(testWalletFile)
    we = NewWalletEntry()
    w = CreateMasterWallet(we)
    assert.Equal(t, w.NumEntries(), 1)
    assert.Equal(t, w.GetAddresses()[0], we.Address)

    // Creating with an invalid wallet entry should panic
    we = NewWalletEntry()
    we.Secret = coin.SecKey{}
    assert.Panics(t, func() { CreateMasterWallet(we) })
    we = NewWalletEntry()
    we.Public = coin.PubKey{}
    assert.Panics(t, func() { CreateMasterWallet(we) })
}
