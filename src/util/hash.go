// Hash func related utilities
package util

import (
    "crypto/sha256"
    "hash"
    "io"
    "io/ioutil"
    "os"
)

var sha256_hash hash.Hash = sha256.New()

// Hashes a file
func HashFile(fp string) ([]byte, error) {
    out := make([]byte, 32)
    f, err := os.Open(fp)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    h := sha256.New()
    _, err = io.Copy(h, f)
    if err != nil {
        return nil, err
    }
    sum := h.Sum(nil)
    copy(out[:32], sum[:32])
    return out, err
}

// Hashes a file, reading file into memory before hashing
func HashFileBuffered(fp string) ([]byte, error) {
    b, err := ioutil.ReadFile(fp)
    if err != nil {
        return nil, err
    }
    out := make([]byte, 32)
    HashBuffer(b, out)
    for i, _ := range b {
        b[i] = 0
    }   // prevent memory leak?
    return out, nil
}

func SHA256(data []byte) []byte {
    sha256_hash.Reset()
    sha256_hash.Write(data)
    sum := sha256_hash.Sum(nil)
    out := make([]byte, 32)
    copy(out, sum[:32])
    return out
}

func HashBuffer(data []byte, out []byte) {
    h := SHA256(data)
    copy(out[:32], h[:32])
}

func HashFileSize(fp string) (hash []byte, size uint64, err error) {
    hash = make([]byte, 32)
    f, err := os.Open(fp)
    if err != nil {
        return hash, 0, err
    }
    defer f.Close()
    h := sha256.New()
    _, err = io.Copy(h, f)
    if err != nil {
        return hash, 0, err
    }
    sum := h.Sum(nil)
    copy(hash[:32], sum[:32])
    stat, err := f.Stat()
    if err != nil {
        return hash, 0, err
    }
    size = uint64(stat.Size()) //may for for +2 GB files on 32 bit
    return hash, size, nil
}
