package coin

import (
    "bytes"
    //"encoding/hex"
    "github.com/skycoin/skycoin/src/lib/base58"
    "log"
)

//version is after Key to enable better vanity address generation
type Address struct {
    Key     [20]byte //sha256(sha256(ridmd160(pubkey)))
    Version byte
    ChkSum  [4]byte
}

//address as Base58 encoded string
func (g *Address) String() string {
    return string(base58.Hex2Base58(g.Key[:]))
}

// Returns address as raw bytes, containing version and then key
func (g *Address) Bytes() []byte {
    b := make([]byte, 25)
    copy(b[0:20], g.Key[0:20])
    b[20] = g.Version
    copy(b[21:25], g.ChkSum[0:4])
    return b
}

// Returns address checksum
// 4 byte checksum
func (g *Address) Checksum() []byte {
    r1 := append(g.Key[:],[]byte{g.Version}...)
    r2 := SumSHA256(r1[:])
    return r2[0:4] //4 bytes
}

func (g *Address) setChecksum() {
    copy(g.ChkSum[0:4], g.Checksum())
}

//r3 := SumSHA256(r2[:])
//r4 := HashRipemd160(r3[:])

func (g *Address) ChecksumVerify() int {
    chksum := g.Checksum()
    if len(chksum) != 4 {
        log.Panic("Invalid address checksum")
    }
    if !bytes.Equal(chksum[0:4], g.ChkSum[0:4]) {
        return 0
    }

    return 1
}

// Creates Address from PubKey
// sha256(sha256(ridmd160(pubkey)))
func AddressFromPubKey(pubkey PubKey) Address {
    r1 := SumSHA256(pubkey[:])
    r2 := SumSHA256(r1[:])
    r3 := HashRipemd160(r2[:])
    addr := Address{Version: 0x0f, Key: r3}
    addr.setChecksum()
    return addr
}
