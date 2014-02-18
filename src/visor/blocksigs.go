package visor

import (
    "errors"
    "github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "io/ioutil"
)

type SignedBlock struct {
    Block coin.Block
    Sig   coin.Sig
}

// Used to serialize the BlockSigs.Sigs map
type BlockSigSerialized struct {
    BkSeq uint64
    Sig   coin.Sig
}

// Used to serialize the BlockSigs.Sigs map
type BlockSigsSerialized struct {
    Sigs []BlockSigSerialized
}

// Manages known BlockSigs as received.
// TODO -- support out of order blocks.  This requires a change to the
// message protocol to support ranges similar to bitcoin's locator hashes.
// We also need to keep track of whether a block has been executed so that
// as continuity is established we can execute chains of blocks.
// TODO -- Since we will need to hold blocks that cannot be verified
// immediately against the blockchain, we need to be able to hold multiple
// BlockSigs per BkSeq, or use hashes as keys.  For now, this is not a
// problem assuming the signed blocks created from master are valid blocks,
// because we can check the signature independently of the blockchain.
type BlockSigs struct {
    Sigs   map[uint64]coin.Sig
    MaxSeq uint64
}

func NewBlockSigs() BlockSigs {
    bs := BlockSigs{
        Sigs:   make(map[uint64]coin.Sig),
        MaxSeq: 0,
    }
    return bs
}

func LoadBlockSigs(filename string) (BlockSigs, error) {
    bs := NewBlockSigs()
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return bs, err
    }
    sigs := BlockSigsSerialized{make([]BlockSigSerialized, 0)}
    err = encoder.DeserializeRaw(data, &sigs)
    if err != nil {
        return bs, err
    }
    bs.Sigs = make(map[uint64]coin.Sig, len(sigs.Sigs))
    for _, s := range sigs.Sigs {
        bs.Sigs[s.BkSeq] = s.Sig
        if s.BkSeq > bs.MaxSeq {
            bs.MaxSeq = s.BkSeq
        }
    }
    return bs, nil
}

func (self *BlockSigs) Save(filename string) error {
    // Convert the Sigs map to an array of element
    sigs := make([]BlockSigSerialized, len(self.Sigs))
    i := 0
    for k, v := range self.Sigs {
        sigs[i] = BlockSigSerialized{
            BkSeq: k,
            Sig:   v,
        }
        i++
    }
    bss := BlockSigsSerialized{sigs}
    data := encoder.Serialize(bss)
    return util.SaveBinary(filename, data, 0644)
}

// Checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (self *BlockSigs) Verify(masterPublic coin.PubKey, bc *coin.Blockchain) error {
    blocks := uint64(len(bc.Blocks))
    if blocks != uint64(len(self.Sigs)) {
        return errors.New("Missing signatures for blocks or vice versa")
    }

    // For now, block sigs must all be sequential and continuous
    if self.MaxSeq+1 != blocks {
        return errors.New("MaxSeq does not match blockchain size")
    }
    for i := uint64(0); i < self.MaxSeq; i++ {
        if _, ok := self.Sigs[i]; !ok {
            return errors.New("Blocksigs missing signature")
        }
    }

    for k, v := range self.Sigs {
        err := coin.VerifySignature(masterPublic, v, bc.Blocks[k].HashHeader())
        if err != nil {
            return err
        }
    }
    return nil
}

// Adds a SignedBlock
func (self *BlockSigs) record(sb *SignedBlock) {
    seq := sb.Block.Head.BkSeq
    self.Sigs[seq] = sb.Sig
    if seq > self.MaxSeq {
        self.MaxSeq = seq
    }
}
