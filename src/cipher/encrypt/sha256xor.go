package encrypt

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/SkycoinProject/skycoin/src/cipher"
	secp256k1 "github.com/SkycoinProject/skycoin/src/cipher/secp256k1-go"
)

const (
	// Data size of each block
	sha256XorBlockSize = 32 // 32 bytes
	// Nonce data size
	sha256XorNonceSize = 32 // 32 bytes
	// Checksum data size
	sha256XorChecksumSize = 32 // 32 bytes
	// Data length size
	sha256XorDataLengthSize = 4 // 4 bytes

)

// Error definition
var (
	ErrMissingPassword       = errors.New("missing password")
	ErrDataTooLarge          = errors.New("data length overflowed, it must <= math.MaxUint32(4294967295)")
	ErrInvalidChecksumLength = errors.New("invalid checksum length")
	ErrInvalidChecksum       = errors.New("invalid data, checksum is not matched")
	ErrInvalidNonceLength    = errors.New("invalid nonce length")
	ErrInvalidBlockSize      = errors.New("invalid block size, must be multiple of 32 bytes")
	ErrReadDataHashFailed    = errors.New("read data hash failed: read length != 32")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrReadDataLengthFailed  = errors.New("read data length failed")
	ErrInvalidDataLength     = errors.New("invalid data length")
)

// DefaultSha256Xor default sha256xor encryptor
var DefaultSha256Xor = Sha256Xor{}

// Sha256Xor provides methods to do encryption and decryption
type Sha256Xor struct{}

// Encrypt encrypts the data with password
//
// 1> Add 32 bits length prefix to indicate the length of data. <length(4 bytes)><data>
// 2> Pad the length + data to 32 bytes with nulls at end
// 2> SHA256(<length(4 bytes)><data><padding>) and prefix the hash. <hash(32 bytes)><length(4 bytes)><data><padding>
// 3> Split the whole data(hash+length+data+padding) into 256 bits(32 bytes) blocks
// 4> Each block is encrypted by XORing the unencrypted block with SHA256(SHA256(password), SHA256(index, SHA256(nonce))
// 	  - index is 0 for the first block of 32 bytes, 1 for the second block of 32 bytes, 2 for third block
// 5> Prefix nonce and SHA256 the nonce with blocks to get checksum, and prefix the checksum
// 6> Finally, the data format is: base64(<checksum(32 bytes)><nonce(32 bytes)><block0.Hex(), block1.Hex()...>)
func (s Sha256Xor) Encrypt(data []byte, password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	if uint(len(data)) > math.MaxUint32 {
		return nil, ErrDataTooLarge
	}

	// Sets data length prefix
	dataLenBytes := make([]byte, sha256XorDataLengthSize)
	binary.LittleEndian.PutUint32(dataLenBytes, uint32(len(data)))

	// Prefixes data with length
	ldata := append(dataLenBytes, data...)

	// Pads length + data with null to 32 bytes
	l := len(ldata) // hash + length + data
	n := l / sha256XorBlockSize
	m := l % sha256XorBlockSize
	if m > 0 {
		paddingNull := make([]byte, sha256XorBlockSize-m)
		ldata = append(ldata, paddingNull...)
		n++
	}

	// Hash(length+data+padding)
	dataHash := cipher.SumSHA256(ldata)

	// Initialize blocks with data hash
	blocks := []cipher.SHA256{dataHash}
	for i := 0; i < n; i++ {
		var b cipher.SHA256
		copy(b[:], ldata[i*sha256XorBlockSize:(i+1)*sha256XorBlockSize])
		blocks = append(blocks, b)
	}

	// Generates a nonce
	nonce := cipher.RandByte(sha256XorNonceSize)
	// Hash the nonce
	hashNonce := cipher.SumSHA256(nonce)
	// Derives key by secp256k1 hashing password
	key := secp256k1.Secp256k1Hash(password)

	var encryptedData []byte
	// Encodes the blocks
	for i := range blocks {
		// Hash(password, hash(index, hash(nonce)))
		h := hashKeyIndexNonce(key, int64(i), hashNonce)
		encryptedHash := blocks[i].Xor(h)
		encryptedData = append(encryptedData, encryptedHash[:]...)
	}

	// Prefix the nonce
	nonceAndDataBytes := append(nonce, encryptedData...)
	// Calculates the checksum
	checkSum := cipher.SumSHA256(nonceAndDataBytes)

	finalData := append(checkSum[:], nonceAndDataBytes...)

	// Base64 encodes the data
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(finalData)))
	enc.Encode(buf, finalData)
	return buf, nil
}

// Decrypt decrypts the data
func (s Sha256Xor) Decrypt(data []byte, password []byte) ([]byte, error) {
	if len(password) == 0 {
		return nil, ErrMissingPassword
	}

	// Base64 decodes data
	enc := base64.StdEncoding
	encData := make([]byte, enc.DecodedLen(len(data)))
	n, err := enc.Decode(encData, data)
	if err != nil {
		return nil, err
	}

	encData = encData[:n]

	// Derives key by secp256k1 hashing password
	key := secp256k1.Secp256k1Hash(password)

	buf := bytes.NewBuffer(encData)

	// Gets checksum
	var checkSum cipher.SHA256
	n, err = buf.Read(checkSum[:])
	if err != nil {
		return nil, err
	}

	if n != sha256XorChecksumSize {
		return nil, ErrInvalidChecksumLength
	}

	// Checks the checksum
	csh := cipher.SumSHA256(buf.Bytes())
	if csh != checkSum {
		return nil, ErrInvalidChecksum
	}

	// Gets the nonce
	nonce := make([]byte, sha256XorNonceSize)
	n, err = buf.Read(nonce)
	if err != nil {
		return nil, err
	}

	if n != sha256XorNonceSize {
		return nil, ErrInvalidNonceLength
	}

	var decodeData []byte
	hashNonce := cipher.SumSHA256(nonce)
	i := 0
	for {
		var block cipher.SHA256
		n, err := buf.Read(block[:])
		if err == io.EOF {
			break
		}

		if n != sha256XorBlockSize {
			return nil, ErrInvalidBlockSize
		}

		// Decodes the block
		dataHash := block.Xor(hashKeyIndexNonce(key, int64(i), hashNonce))
		decodeData = append(decodeData, dataHash[:]...)
		i++
	}

	buf = bytes.NewBuffer(decodeData)

	// Gets the hash
	var dataHash cipher.SHA256
	n, err = buf.Read(dataHash[:])
	if err != nil {
		return nil, fmt.Errorf("read data hash failed: %v", err)
	}

	if n != 32 {
		return nil, ErrReadDataHashFailed
	}

	// Checks the hash
	if dataHash != cipher.SumSHA256(buf.Bytes()) {
		return nil, ErrInvalidPassword
	}

	// Reads out the data length
	dataLenBytes := make([]byte, sha256XorDataLengthSize)
	n, err = buf.Read(dataLenBytes)
	if err != nil {
		return nil, err
	}

	if n != sha256XorDataLengthSize {
		return nil, ErrReadDataLengthFailed
	}

	l := binary.LittleEndian.Uint32(dataLenBytes)

	if uint64(buf.Len()) > math.MaxUint32 {
		return nil, ErrDataTooLarge
	}

	if l > uint32(buf.Len()) {
		return nil, ErrInvalidDataLength
	}

	// Reads out the raw data
	rawData := make([]byte, l)
	n, err = buf.Read(rawData)
	if err != nil {
		return nil, err
	}

	if uint32(n) != l {
		return nil, fmt.Errorf("read data failed, expect %d bytes, but get %d bytes", l, n)
	}

	return rawData, nil
}

// hash(password, hash(index, hash(nonce)))
func hashKeyIndexNonce(key []byte, index int64, nonceHash cipher.SHA256) cipher.SHA256 {
	// convert index to 256bit number
	indexBytes := make([]byte, 32)
	binary.PutVarint(indexBytes, index)

	// hash(index, nonceHash)
	indexNonceHash := cipher.SumSHA256(append(indexBytes, nonceHash[:]...))

	// hash(hash(password), indexNonceHash)
	var keyHash cipher.SHA256
	copy(keyHash[:], key[:])
	return cipher.AddSHA256(keyHash, indexNonceHash)
}
