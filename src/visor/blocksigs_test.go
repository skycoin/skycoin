package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func makeBlocks(mv *Visor, n int) ([]SignedBlock, error) {
    dest := NewWalletEntry()
    blocks := make([]SignedBlock, 0, n)
    for i := 0; i < n; i++ {
        tx, err := mv.Spend(Balance{10 * 1e6, 0}, 0, dest.Address)
        if err != nil {
            return nil, err
        }
        mv.RecordTxn(tx, false)
        sb, err := mv.CreateBlock()
        if err != nil {
            return nil, err
        }
        blocks = append(blocks, sb)
    }
    return blocks, nil
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

func TestNewBlockSigs(t *testing.T) {
    bs := NewBlockSigs()
    assert.NotNil(t, bs.Sigs)
    assert.Equal(t, len(bs.Sigs), 0)
    assert.Equal(t, bs.MaxSeq, uint64(0))
}

func TestSaveLoadBlockSigs(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    sbs, err := makeBlocks(mv, 7)
    assert.Nil(t, err)
    bs := NewBlockSigs()
    for _, sb := range sbs {
        bs.Sigs[sb.Block.Header.BkSeq] = sb.Sig
    }
    // We give it an invalid BkSeq, because the BkSeq should be corrected
    // when loaded
    bs.MaxSeq = uint64(0)

    err = bs.Save(testBlocksigsFile)
    assert.Nil(t, err)
    assertFileExists(t, testBlocksigsFile)

    newBs, err := LoadBlockSigs(testBlocksigsFile)
    assert.Nil(t, err)
    assert.Equal(t, newBs.MaxSeq, uint64(len(mv.blockchain.Blocks)-1))
    assert.Equal(t, len(newBs.Sigs), len(bs.Sigs))
    for k, v := range bs.Sigs {
        w, ok := newBs.Sigs[k]
        assert.True(t, ok)
        assert.Equal(t, v, w)
    }

    // Loading a corrupted file should cause error in deserialization
    f, err := os.OpenFile(testBlocksigsFile, os.O_WRONLY|os.O_TRUNC, 0644)
    assert.Nil(t, err)
    b := make([]byte, 1)
    _, err = f.Write(b)
    assert.Nil(t, err)
    f.Close()

    newBs, err = LoadBlockSigs(testBlocksigsFile)
    assert.NotNil(t, err)
    assert.Equal(t, newBs.MaxSeq, uint64(0))
}

func TestBlockSigsVerify(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    bc := mv.blockchain
    pub := mv.Config.MasterKeys.Public
    sbs, err := makeBlocks(mv, 7)
    assert.Nil(t, err)

    bs := NewBlockSigs()
    bs.Sigs[uint64(0)] = mv.blockSigs.Sigs[0]

    err = bs.Verify(pub, bc)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Missing signatures for blocks or vice versa")

    // MaxSeq incorrect
    for _, sb := range sbs {
        bs.Sigs[sb.Block.Header.BkSeq] = sb.Sig
    }
    err = bs.Verify(pub, bc)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "MaxSeq does not match blockchain size")

    // Block missing from continuous sequence, despite number of sigs correct
    bs.MaxSeq = uint64(len(bc.Blocks) - 1)
    rm := bs.Sigs[uint64(3)]
    delete(bs.Sigs, uint64(3))
    bs.Sigs[uint64(100)] = rm
    err = bs.Verify(pub, bc)
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Blocksigs missing signature")
    delete(bs.Sigs, uint64(100))

    // Invalid signature
    bs.Sigs[uint64(3)] = coin.Sig{}
    err = bs.Verify(pub, bc)
    assert.NotNil(t, err)

    // Valid
    bs.Sigs[uint64(3)] = rm
    err = bs.Verify(pub, bc)
    assert.Nil(t, err)

    // Saving and loading should pass verification
    err = bs.Save(testBlocksigsFile)
    assert.Nil(t, err)
    bs, err = LoadBlockSigs(testBlocksigsFile)
    assert.Nil(t, err)
    err = bs.Verify(pub, bc)
    assert.Nil(t, err)
}

func TestBlockSigsRecord(t *testing.T) {
    defer cleanupVisor()
    mv := setupMasterVisor()
    bs := NewBlockSigs()
    bs.record(&SignedBlock{
        Sig:   mv.blockSigs.Sigs[0],
        Block: mv.blockchain.Blocks[0],
    })
    assert.Equal(t, len(bs.Sigs), 1)
    sbs, err := makeBlocks(mv, 5)
    assert.Nil(t, err)
    for i := 0; i < 5; i++ {
        bs.record(&sbs[i])
        assert.Equal(t, len(bs.Sigs), i+2)
        assert.Equal(t, bs.Sigs[uint64(i+1)], sbs[i].Sig)
        assert.Equal(t, bs.MaxSeq, sbs[i].Block.Header.BkSeq)
    }
}
