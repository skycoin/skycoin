package blockchain

import (
    //"errors"
    "github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/coin"
    "github.com/skycoin/skycoin/src/util"
    //"io/ioutil"

    "hash/fnv"
    "hash"
    "log"
)


var (
    fnvs hash.Hash64 = fnv.New64a()
)

func FnvsHash(data []byte) uint64 {
    fnvs.Reset()
    fnvs.Write(data)
    return fnvs.Sum64()
}

func le_Uint64(b []byte) uint64 {
    return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
        uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func le_PutUint64(b []byte, v uint64) {
    b[0] = byte(v)
    b[1] = byte(v >> 8)
    b[2] = byte(v >> 16)
    b[3] = byte(v >> 24)
    b[4] = byte(v >> 32)
    b[5] = byte(v >> 40)
    b[6] = byte(v >> 48)
    b[7] = byte(v >> 56)
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


//type BlockSerialized struct {
//    Length uint64
//    Chksum uint64
//    Data []byte
//}

// Used to serialize the BlockSigs.Sigs map
type BlockchainFile struct {
    Blocks []SignedBlock //
}

//encode block to bytes
func Enc_block(sb SignedBlock) []byte {
    var prefix []byte = make([]byte, 16)
    b := encoder.Serialize(sb)

    var dln uint64 = uint64(len(b))
    var chk uint64 = FnvsHash(b)

    le_PutUint64(prefix[0:8], dln)
    le_PutUint64(prefix[8:16], chk)

    b = append(prefix[:], b...)

    return b
}

//decode block to bytes
func Dec_block(b []byte) (SignedBlock, error) {
    if len(b) < 16 {
        log.Panic()
    }

    var dln uint64 = le_Uint64(b[0:8])
    var chk uint64 = le_Uint64(b[8:16])

    b = b[16:]

    if dln != uint64(len(b)) {
        log.Panic("Dec_block, length check failed")
    }

    if chk != FnvsHash(b) {
        log.Panic("Dec_block, checksum failed")
    }

    var sb SignedBlock

    err := encoder.DeserializeRaw(b, &sb)

    if err != nil {
        log.Panic("Dec_block, deserialization failed")
    }

    return sb, nil
}

//TODO: write individual blocks. append as they come in
func (self *BlockchainFile) Save(filename string) error {
    var data []byte
    for _,b := range self.Blocks {
        data = append(data, Enc_block(b)...)
    }
    util.SaveBinary(filename, data, 0644)
    return nil
}

func (self *BlockchainFile) Load(filename string) error {
    var data []byte
    for _,b := range self.Blocks {
        data = append(data, Enc_block(b)...)
    }
    util.SaveBinary(filename, data, 0644)
    return nil
}


/*
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
*/

// Checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.

/*
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
*/