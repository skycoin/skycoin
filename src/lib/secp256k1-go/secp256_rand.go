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
