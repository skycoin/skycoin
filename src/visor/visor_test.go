package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    // "github.com/skycoin/skycoin/src/util"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestNewVisorConfig(t *testing.T) {
    vc := NewVisorConfig()
    assert.False(t, vc.IsMaster)
    assert.True(t, vc.CanSpend)
    assert.Equal(t, vc.WalletFile, "")
    assert.Equal(t, vc.BlockchainFile, "")
    assert.Equal(t, vc.BlockSigsFile, "")
    assert.NotNil(t, vc.MasterKeys.Verify())
    assert.Equal(t, vc.GenesisSignature, coin.Sig{})
    assert.Equal(t, vc.WalletSizeMin, 1)
}

func setupVisorWriting(vc VisorConfig) *Visor {
    vc.WalletFile = testWalletFile
    vc.BlockSigsFile = testBlocksigsFile
    vc.BlockchainFile = testBlockchainFile
    return NewVisor(vc)
}

func writeVisorFilesDirect(t *testing.T, v *Visor) {
    assert.Nil(t, v.SaveWallet())
    assert.Nil(t, v.SaveBlockSigs())
    assert.Nil(t, v.SaveBlockchain())
    assertFileExists(t, v.Config.WalletFile)
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
    v := NewVisor(vc)
    we.Secret = coin.SecKey{}
    return we, v.blockSigs.Sigs[0], v.blockchain.Blocks[0].Header.Time
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
    vc.WalletFile = testWalletFile
    vc.BlockchainFile = testBlockchainFile
    vc.BlockSigsFile = testBlocksigsFile
    return vc
}

func newMasterVisorConfig(t *testing.T) VisorConfig {
    vc := NewVisorConfig()
    vc.MasterKeys = newWalletEntry(t)
    vc.IsMaster = true
    return vc
}

func TestNewVisor(t *testing.T) {
    defer cleanupVisor()

    // Not master, Invalid master keys
    cleanupVisor()
    vc := NewVisorConfig()
    vc.IsMaster = false
    assert.Panics(t, func() { NewVisor(vc) })

    // Master, invalid master keys
    cleanupVisor()
    vc.IsMaster = true
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
    assert.Equal(t, len(v.Wallet.Entries), 10)
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Equal(t, v.masterKeys, vc.MasterKeys)

    // Master, no wallet, blockchain, blocksigs file
    cleanupVisor()
    vc = NewVisorConfig()
    vc.MasterKeys = newWalletEntry(t)
    vc.WalletSizeMin = 10
    vc.IsMaster = true
    v = NewVisor(vc)
    // Wallet should only have 1 entry if master
    assert.Equal(t, len(v.Wallet.Entries), 1)
    assert.Equal(t, len(v.blockchain.Blocks), 1)
    assert.Equal(t, len(v.blockSigs.Sigs), 1)
    assert.Equal(t, v.masterKeys, vc.MasterKeys)

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
