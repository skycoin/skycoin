package sb

import (
	"crypto/sha256"
	"hash"
	"io"
	"io/ioutil"
	"os"
)

var sha256_hash hash.Hash = sha256.New()

func _HashFile(fp string, out []byte) error {
	f, err := os.Open(fp)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return err
	}
	sum := h.Sum(nil)
	copy(out[0:32], sum[0:32])
	return nil
}

func HashFile(fp string) ([32]byte, error) {
	var out [32]byte
	err := _HashFile(fp, out[:])
	return out, err
}

//reads file into memory before hashing
func HashFile2(fp string) ([32]byte, error) {
	var out [32]byte
	b, err := ioutil.ReadFile(fp) //read file
	if err != nil {
		return out, err
	}
	HashBuffer(b, out[:])
	for i, _ := range b {
		b[i] = 0
	} // prevent memory leak?
	return out, nil
}

func HashFileSize(fp string) (hash [32]byte, size uint64, err error) {
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
	copy(hash[0:32], sum[0:32])
	stat, err := f.Stat()
	if err != nil {
		return hash, 0, err
	}
	size = uint64(stat.Size()) //may for for +2 GB files on 32 bit
	return hash, size, nil
}
