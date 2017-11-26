package cipher

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

const (
	// the data size of each block
	blockSize = 32 // 32 bytes
	// the data length prefix size
	lenPrefixSize = 4 // 4 bytes
)

// Encrypt encrypts the data with password
// 1> Add 32 bits length prefix to indicate the length of data
// 2> Split data into 256 bits(32 bytes) blocks (pad to 32 bytes with nulls at end)
// 3> Each block is encrypted by XORing the unencrypted block with SHA256(SHA256(password), SHA256(index, SHA256(nonce))
// 	  - index is 0 for the first block of 32 bytes, 1 for the second block of 32 bytes, 2 for third block
// 4> SHA256 the nonce with comma seperated, hex encoded blocks of 32 bytes(256 bits)
// 5> Encode <checksum(32 bytes)><nonce(32 bytes)><block0.Hex(), block1.Hex()...> with base64
func Encrypt(data []byte, password []byte) ([]byte, error) {
	// set data length prefix
	prefix := make([]byte, lenPrefixSize)
	binary.PutUvarint(prefix, uint64(len(data)))
	pdata := append(prefix, data...)

	// split the data into 256 bit blocks(pad with null to 32 bytes)
	l := len(pdata)
	n := l / blockSize
	var blocks [][blockSize]byte
	for i := 0; i < n; i++ {
		var b [blockSize]byte
		copy(b[:], pdata[i*blockSize:(i+1)*blockSize])
		blocks = append(blocks, b)
	}

	// append last block if exist
	if l%blockSize > 0 {
		b := [blockSize]byte{}
		copy(b[:], pdata[n*blockSize:])
		blocks = append(blocks, b)
	}

	nonce := RandByte(blockSize)
	var encryptedBlocks []string
	// encode the blocks
	for i := range blocks {
		h := hashPwdIndexNonce(password, int64(i), nonce)
		bh := SHA256(blocks[i])
		encryptedBlocks = append(encryptedBlocks, bh.Xor(h).Hex())
	}

	encryptedData := strings.Join(encryptedBlocks, ",")
	nonceAndDataBytes := append(nonce, []byte(encryptedData)...)
	checkSum := SumSHA256(nonceAndDataBytes)
	var buf bytes.Buffer
	_, err := buf.Write(checkSum[:])
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(nonceAndDataBytes)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decrypt(data []byte, password []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)

	checkSumBytes := make([]byte, blockSize)
	n, err := buf.Read(checkSumBytes)
	if err != nil {
		return nil, err
	}

	if n != blockSize {
		return nil, errors.New("decode checksum failed")
	}

	var checkSum SHA256
	copy(checkSum[:], checkSumBytes)

	// verify the checksum
	csh := SumSHA256(buf.Bytes())
	if csh.Hex() != checkSum.Hex() {
		return nil, errors.New("invalid checksum")
	}

	nonce := make([]byte, blockSize)
	n, err = buf.Read(nonce)
	if err != nil {
		return nil, err
	}

	if n != blockSize {
		return nil, errors.New("decode nonce failed")
	}

	encryptedBlocks := strings.Split(buf.String(), ",")
	var decodeData []byte
	for i := range encryptedBlocks {
		bh, err := SHA256FromHex(encryptedBlocks[i])
		if err != nil {
			return nil, err
		}

		dataHash := bh.Xor(hashPwdIndexNonce(password, int64(i), nonce))
		decodeData = append(decodeData, dataHash[:]...)
	}

	buf = bytes.NewBuffer(decodeData)
	l, err := binary.ReadUvarint(bytes.NewReader(decodeData[:lenPrefixSize]))
	if err != nil {
		return nil, fmt.Errorf("read prefix length failed: %v", err)
	}

	if l > uint64(len(decodeData[4:])) {
		return nil, fmt.Errorf("prefix length > data length")
	}

	return decodeData[lenPrefixSize : l+lenPrefixSize], nil
}

// hash(password, hash(index, hash(nonce)))
func hashPwdIndexNonce(password []byte, index int64, nonce []byte) SHA256 {
	// convert index to 256bit number
	indexBytes := make([]byte, 32)
	binary.PutVarint(indexBytes, index)

	// hash(index, hash(nonce))
	nonceHash := SumSHA256(nonce)
	indexNonceHash := SumSHA256(append(indexBytes, nonceHash[:]...))

	// hash(hash(password), indexNonceHash)
	return AddSHA256(SumSHA256(password), indexNonceHash)
}
