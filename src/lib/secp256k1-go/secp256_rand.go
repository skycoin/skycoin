package secp256

import (
    "io"
    "crypto/rand"
)



//use entropy pool etc and cryptographic random number generator
//mix in time
//mix in mix in cpu cycle count
func RandByte(n int) []byte {
    buff := make([]byte, n)
    ret, err := io.ReadFull(rand.Reader, buff)
    if len(buff) != ret || err != nil {
        return nil
    }
    return buff
}

var (
    sha256Hash    hash.Hash = sha256.New()
)

func SumSHA256(b []byte) []byte] {
    sha256Hash.Reset()
    sha256Hash.Write(b)
    sum := sha256Hash.Sum(nil)
    return sum[:]
}