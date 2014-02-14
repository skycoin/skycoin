package visor

import (
    "errors"
    "github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    "io/ioutil"

    "hash/fnv"
    "hash"
)


var (
    fnvs    hash.Hash = fnvs.New64a()
)

func FnvsHash(data []byte) uint64 {
    fnvs.Reset()
    fnvs.Write(data)
    return fnvs.Sum64(nil)
}

//Todo:
// - store signature and block in same
// - prefix block with size
// - blockchain should be append only

type SignedBlock struct {
    BkSeq uint64
    Sig   coin.Sig
    Block coin.Block
}

// Used to serialize the BlockSigs.Sigs map
type BlockSerialized struct {
    Length uint64
    Data []byte
}

// Used to serialize the BlockSigs.Sigs map
type BlockchainFile struct {
    BA []BlockSerialized //
}

func (self *BlockchainFile) Load(filename string) error {

}

func SaveBlockchain(filename string) (BlockSigs, error) {

}

func LoadBlockchain(filename string) (BlockSigs, error) {
    

}

func LoadBlockchain(filename string) (BlockSigs, error) {
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
    sigs := make([]BlockSigSerialized, 0, len(self.Sigs))
    for k, v := range self.Sigs {
        sigs = append(sigs, BlockSigSerialized{
            BkSeq: k,
            Sig:   v,
        })
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
    seq := sb.Block.Header.BkSeq
    self.Sigs[seq] = sb.Sig
    if seq > self.MaxSeq {
        self.MaxSeq = seq
    }
}
