package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/wallet"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestSerializedBlockchain(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()
    bc := &coin.Blockchain{}
    bc.Blocks = []coin.Block{}
    for i := uint64(0); i < 10; i++ {
        bc.Blocks = append(bc.Blocks, coin.Block{})
        bc.Blocks[i].Head.BkSeq = i
    }
    assert.Equal(t, len(bc.Blocks), 10)
    bc.Unspent = coin.NewUnspentPool()
    for i := uint64(0); i < 10; i++ {
        bc.Unspent.Add(makeUxOut(t))
    }
    assert.Equal(t, len(bc.Unspent.Pool), 10)

    sbc := NewSerializedBlockchain(bc)
    assert.Equal(t, sbc.Blocks, bc.Blocks)
    assert.Equal(t, len(sbc.Unspents), len(bc.Unspent.Pool))

    // Back to blockchain works
    assert.Equal(t, bc, sbc.ToBlockchain())

    // Saving and reloading works
    assert.Nil(t, sbc.Save(testBlockchainFile))
    assertFileExists(t, testBlockchainFile)

    sbc2, err := LoadSerializedBlockchain(testBlockchainFile)
    assert.Nil(t, err)
    assert.Equal(t, sbc, sbc2)
    assert.Equal(t, bc, sbc2.ToBlockchain())

    bc2, err := LoadBlockchain(testBlockchainFile)
    assert.Nil(t, err)
    assert.Equal(t, bc, bc2)
}

func TestLoadBlockchain(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()

    // Loading a non-existent blockchain should return error
    bc, err := LoadBlockchain(testBlockchainFile)
    assert.NotNil(t, err)
    assert.True(t, os.IsNotExist(err))

    // Loading a real blockchain should be fine
    vc := newMasterVisorConfig(t)
    v := NewVisor(vc)
    v.Config.BlockchainFile = testBlockchainFile
    assert.Nil(t, transferCoinsToSelf(v, v.Config.MasterKeys.Address))
    assert.Equal(t, len(v.blockchain.Blocks), 2)
    v.SaveBlockchain()
    assertFileExists(t, testBlockchainFile)
    bc, err = LoadBlockchain(testBlockchainFile)
    assert.Nil(t, err)
    assert.Equal(t, v.blockchain, bc)

    // Loading a corrupted blockchain should return error
    corruptFile(t, testBlockchainFile)
    _, err = LoadBlockchain(testBlockchainFile)
    assert.NotNil(t, err)
}

func TestLoadBlockchainPrivate(t *testing.T) {
    defer cleanupVisor()
    cleanupVisor()

    we := wallet.NewWalletEntry()

    // No filename should return fresh blockchain
    bc := loadBlockchain("", we.Address)
    assert.Equal(t, len(bc.Blocks), 0)

    // Filename with no file should return fresh blockchain
    assertFileNotExists(t, testBlockchainFile)
    bc = loadBlockchain(testBlockchainFile, we.Address)
    assert.Equal(t, len(bc.Blocks), 0)

    // Loading an empty blockchain should panic
    assert.Nil(t, SaveBlockchain(bc, testBlockchainFile))
    assertFileExists(t, testBlockchainFile)
    assert.Panics(t, func() {
        loadBlockchain(testBlockchainFile, we.Address)
    })

    // Loading a blockchain with a different genesis address should panic
    vc := newMasterVisorConfig(t)
    bc.CreateGenesisBlock(vc.MasterKeys.Address, 0, 100e6)
    assert.Equal(t, len(bc.Blocks), 1)
    assert.Nil(t, SaveBlockchain(bc, testBlockchainFile))
    assertFileExists(t, testBlockchainFile)
    assert.Panics(t, func() {
        loadBlockchain(testBlockchainFile, coin.Address{})
    })

    // Loading a corrupt blockchain should panic
    corruptFile(t, testBlockchainFile)
    assert.Panics(t, func() {
        loadBlockchain(testBlockchainFile, we.Address)
    })
    cleanupVisor()

    // Loading a valid blockchain should be safe
    vc = newMasterVisorConfig(t)
    vc.BlockchainFile = testBlockchainFile
    v := NewVisor(vc)
    assert.Nil(t, transferCoinsToSelf(v, v.Config.MasterKeys.Address))
    assert.Equal(t, len(v.blockchain.Blocks), 2)
    assert.Nil(t, v.SaveBlockchain())
    assertFileExists(t, testBlockchainFile)
    bc = loadBlockchain(testBlockchainFile, v.Config.MasterKeys.Address)
    assert.Equal(t, v.blockchain, bc)
}
