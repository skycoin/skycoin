package wallet

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
)

const blockSize = 32

// encrypt encrypts the data with password
// 1> Split data into 256 bits(32 bytes) blocks (pad to 32 bytes with nulls at end)
// 2> Each block is encrypted by XORing the unencrypted block with SHA256(SHA256(password), SHA256(index, SHA256(nonce))
// 	  - index is 0 for the first block of 32 bytes, 1 for the second block of 32 bytes, 2 for third block
// 3> SHA256 the nonce with comma seperated, hex encoded blocks of 32 bytes(256 bits)
// 4> Encode <checksum(32 bytes)><nonce(32 bytes)><block0.Hex(), block1.Hex()...> with base64
func encrypt(data []byte, password []byte) (string, error) {
	// split the data into 256 bit blocks(pad with null to 32 bytes)
	l := len(data)
	n := l / blockSize
	var blocks [][blockSize]byte
	for i := 0; i < n; i++ {
		var b [blockSize]byte
		copy(b[:], data[i*blockSize:(i+1)*blockSize])
		blocks = append(blocks, b)
	}

	// append last block if exist
	if l%blockSize > 0 {
		b := [blockSize]byte{}
		copy(b[:], data[n*blockSize:])
		blocks = append(blocks, b)
	}

	nonce := cipher.RandByte(blockSize)
	var encryptedBlocks []string
	// encode the blocks
	for i := range blocks {
		h := hashPwdIndexNonce(password, int64(i), nonce)
		bh := cipher.SHA256(blocks[i])
		encryptedBlocks = append(encryptedBlocks, bh.Xor(h).Hex())
	}

	encryptedData := strings.Join(encryptedBlocks, ",")
	nonceAndDataBytes := append(nonce, []byte(encryptedData)...)
	checkSum := cipher.SumSHA256(nonceAndDataBytes)
	var buf bytes.Buffer
	_, err := buf.Write(checkSum[:])
	if err != nil {
		return "", err
	}

	_, err = buf.Write(nonceAndDataBytes)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func decrypt(data string, password []byte) ([]byte, error) {
	// decode with base64
	dataBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(dataBytes)

	checkSumBytes := make([]byte, blockSize)
	n, err := buf.Read(checkSumBytes)
	if err != nil {
		return nil, err
	}

	if n != blockSize {
		return nil, errors.New("decode checksum failed")
	}

	var checkSum cipher.SHA256
	copy(checkSum[:], checkSumBytes)

	// verify the checksum
	csh := cipher.SumSHA256(buf.Bytes())
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
		bh, err := cipher.SHA256FromHex(encryptedBlocks[i])
		if err != nil {
			return nil, err
		}

		dataHash := bh.Xor(hashPwdIndexNonce(password, int64(i), nonce))
		decodeData = append(decodeData, dataHash[:]...)
	}

	// remove the padding null
	return bytes.TrimRight(decodeData, "\x00"), nil
}

// hash(password, hash(index, hash(nonce)))
func hashPwdIndexNonce(password []byte, index int64, nonce []byte) cipher.SHA256 {
	// convert index to 256bit number
	indexBytes := make([]byte, 32)
	binary.PutVarint(indexBytes, index)

	// hash(index, hash(nonce))
	nonceHash := cipher.SumSHA256(nonce)
	indexNonceHash := cipher.SumSHA256(append(indexBytes, nonceHash[:]...))

	// hash(hash(password), indexNonceHash)
	return cipher.AddSHA256(cipher.SumSHA256(password), indexNonceHash)
}
