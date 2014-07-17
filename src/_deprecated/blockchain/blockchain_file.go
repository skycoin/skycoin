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
    "io/ioutil"
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
    //BkSeq uint64
    Sig   cipher.Sig
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

    var dln uint64 = uint64(len(b)+24) //length include prefix
    var seq uint64 = uint64(sb.Block.Head.BkSeq)
    var chk uint64 = FnvsHash(b)

    le_PutUint64(prefix[0:8], dln)
    le_PutUint64(prefix[8:16], seq)
    le_PutUint64(prefix[16:24], chk)

    b = append(prefix[:], b...)

    return b
}

//Decode Length, return length of next data element
func Dec_len(b []byte) uint64 {
    var dln uint64 = le_Uint64(b[0:8])
    return dln
}

//decode block to bytes
func Dec_block(b []byte) (SignedBlock, error) {
    if len(b) < 16 {
        log.Panic()
    }

    var dln uint64 = le_Uint64(b[0:8])
    var seq uint64 = le_Uint64(b[8:16])
    var chk uint64 = le_Uint64(b[16:24])

    if dln != uint64(len(b)) {
        log.Panic("Dec_block, length check failed")
    }

    b = b[24:] //cleave off header

    if chk != FnvsHash(b) {
        log.Panic("Dec_block, checksum failed")
    }

    var sb SignedBlock
    err := encoder.DeserializeRaw(b, &sb)
    if err != nil {
        log.Panic("Dec_block, deserialization failed")
    }

    if seq != sb.Block.Head.BkSeq {
        log.Panic("Dec_block, seq mismatch")
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

//can steam and buffer into 50 meg buffer
//does not need to read in all blocks at once
func (self *BlockchainFile) Load(filename string) ([]SignedBlock, error) {
    var sb []SignedBlock

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    for len(data) > 4 {
        var dln uint64 = Dec_len(data)
        bd := data[:dln] //block data
        data = data[dln:]
        b, err := Dec_block(bd)
        if err != nil {
            log.Panic("BlockchainFile, Load, Decoding Block failed")
        }
        sb = append(sb, b)
    }
    return sb, nil
}