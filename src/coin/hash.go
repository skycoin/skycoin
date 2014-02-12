package coin

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "github.com/skycoin/encoder"
    "github.com/skycoin/skycoin/src/lib/ripemd160"
    "hash"
    "log"
)

var (
    sha256Hash    hash.Hash = sha256.New()
    ripemd160Hash hash.Hash = ripemd160.New()
)

// Ripemd160

type Ripemd160 [20]byte

func (self *Ripemd160) Set(b []byte) {
    if len(b) != 20 {
        log.Panic("Invalid ripemd160 length")
    }
    copy(self[:], b[:])
}

func HashRipemd160(data []byte) Ripemd160 {
    ripemd160Hash.Reset()
    ripemd160Hash.Write(data)
    sum := ripemd160Hash.Sum(nil)
    var h Ripemd160
    copy(h[:], sum[:])
    return h
}

// SHA256

type SHA256 [32]byte

func (g *SHA256) Set(b []byte) {
    if len(b) != 32 {
        log.Panic("Invalid sha256 length")
    }
    copy(g[:], b[:])
}

func (g *SHA256) Hex() string {
    return hex.EncodeToString(g[:])
}

func SumSHA256(b []byte) SHA256 {
    sha256Hash.Reset()
    sha256Hash.Write(b)
    sum := sha256Hash.Sum(nil)
    var h SHA256
    copy(h[:], sum[:])
    return h
}

func SHA256FromHex(hs string) (SHA256, error) {
    h := SHA256{}
    b, err := hex.DecodeString(hs)
    if err != nil {
        return h, err
    }
    if len(b) != len(h) {
        return h, errors.New("Invalid hex length")
    }
    h.Set(b)
    return h, nil
}

func MustSHA256FromHex(hs string) SHA256 {
    h, err := SHA256FromHex(hs)
    if err != nil {
        log.Panic(err)
    }
    return h
}

// Like SumSHA256, but len(b) must equal n, or panic
func MustSumSHA256(b []byte, n int) SHA256 {
    if len(b) != n {
        log.Panic("len(b) != n")
    }
    return SumSHA256(b)
}

// Double SHA256
func SumDoubleSHA256(b []byte) SHA256 {
    h := SumSHA256(b)
    return AddSHA256(h, h)
}

// Returns the SHA256 hash of to two concatenated hashes
func AddSHA256(a1 SHA256, b1 SHA256) SHA256 {
    b := append(a1[:], b1[:]...)
    return MustSumSHA256(b, len(b))
}

func (h1 *SHA256) Xor(h2 SHA256) SHA256 {
    var h3 SHA256
    for i := 0; i < len(h1); i++ {
        h3[i] = h1[i] ^ h2[i]
    }
    return h3
}

//compute root merkle tree hash of hash list
//pad input to power of 16
//group inputs hashes into groups of 16 and hash them down to single hash
//repeat until there is single hash in list
func Merkle(h0 []SHA256) SHA256 {
    //fmt.Printf("Merkle 0: len= %v \n", len(h0))
    if len(h0) == 0 {
        return SHA256{} //zero hash
    }
    np := 0
    for np = 1; np < len(h0); np *= 16 {
    }
    h1 := make([]SHA256, np)

    //var th SHA256 = h0[0]

    var lh0 = len(h0)
    for i := lh0; i < np; i++ { //pad to power of 16
        h1[i] = h0[0] //pad first element till 16
    }

    for len(h1) != 1 {
        //fmt.Printf("Merkle 1: len= %v \n", len(h1))
        h2 := make([]SHA256, len(h1)/16)
        var h3 [16]SHA256
        for i := 0; i < len(h2); i++ {
            for j := 0; j < 16; j++ {
                h3[j] = h1[16*i+j]
            }
            h2[i] = MustSumSHA256(encoder.Serialize(h3), 16*32)
        }
        h1 = h2
    }
    return h1[0]
}

func HashArrayHasDupes(ha []SHA256) bool {
    for i := 0; i < len(ha); i++ {
        for j := i + 1; j < len(ha); j++ {
            if ha[i] == ha[j] {
                return true
            }
        }
    }
    return false
}
