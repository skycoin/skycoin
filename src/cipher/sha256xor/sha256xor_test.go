package sha256xor

import (
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	secp256k1 "github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/stretchr/testify/require"
)

func TestEncrypt(t *testing.T) {
	tt := []struct {
		name     string
		data     []byte
		password []byte
		err      error
	}{
		{
			"data length=1 password is empty=true",
			randBytes(t, 1),
			nil,
			errors.New("missing password"),
		},
		{
			"data length=1  password is empty=false",
			randBytes(t, 1),
			[]byte("key"),
			nil,
		},
		{
			"data length<32 password is empty=false",
			randBytes(t, 2),
			[]byte("pwd"),
			nil,
		},
		{
			"data length=32 password is empty=false",
			randBytes(t, 32),
			[]byte("pwd"),
			nil,
		},
		{
			"data length=2*32 password is empty=false",
			randBytes(t, 64),
			[]byte("9JMkCPphe73NQvGhmab"),
			nil,
		},
		{
			"data length>2*32 password is empty=false",
			randBytes(t, 65),
			[]byte("9JMkCPphe73NQvGhmab"),
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			edata, err := Encrypt(tc.data, tc.password)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			n := (lengthSize + len(tc.data)) / blockSize
			m := (lengthSize + len(tc.data)) % blockSize
			if m > 0 {
				n++
			}

			totalEncryptedDataLen := checksumSize + nonceSize + 32 + n*blockSize // 32 is the hash data length
			require.Equal(t, totalEncryptedDataLen, len(edata))
			var checksum SHA256
			copy(checksum[:], edata[:checksumSize])
			require.Equal(t, checksum, SumSHA256(edata[checksumSize:]))
		})
	}

	// test data of length range in 32 to 64, for testing the 32 bytes padding.
	pwd := []byte("pwd")
	for i := 33; i <= 64; i++ {
		name := fmt.Sprintf("data length=%d password is empty=false", i)
		t.Run(name, func(t *testing.T) {
			data := randBytes(t, i)
			edata, err := Encrypt(data, pwd)
			require.NoError(t, err)

			n := (lengthSize + len(data)) / blockSize
			m := (lengthSize + len(data)) % blockSize
			if m > 0 {
				n++
			}

			totalEncryptedDataLen := checksumSize + nonceSize + 32 + n*blockSize // 32 is the hash data length
			require.Equal(t, totalEncryptedDataLen, len(edata))
			var checksum SHA256
			copy(checksum[:], edata[:checksumSize])
			require.Equal(t, checksum, SumSHA256(edata[checksumSize:]))
		})
	}
}

func TestDecrypt(t *testing.T) {
	data := randBytes(t, 32)
	tt := []struct {
		name          string
		encryptedData func() []byte // encrypted data
		password      []byte
		err           error
	}{
		{
			"invalid data length",
			func() []byte {
				return makeEncryptedData(data, 65, []byte("pwd"))
			},
			[]byte("pwd"),
			errors.New("invalid data length"),
		},
		{
			"invalid checksum",
			func() []byte {
				edata := makeEncryptedData(data, 32, []byte("pwd"))
				// Changes the encrypted data, so that the checksum could not match
				edata[len(edata)-1]++
				return edata
			},
			[]byte("pwd"),
			errors.New("invalid data, checksum is not matched"),
		},
		{
			"empty password",
			func() []byte {
				return makeEncryptedData(data, 32, []byte("pwd"))
			},
			[]byte(""),
			errors.New("missing password"),
		},
		{
			"nil password",
			func() []byte {
				return makeEncryptedData(data, 32, []byte("pwd"))
			},
			nil,
			errors.New("missing password"),
		},
		{
			"invalid password",
			func() []byte {
				return makeEncryptedData(data, 32, []byte("pwd"))
			},
			[]byte("wrong password"),
			errors.New("invalid password"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			edata := tc.encryptedData()
			d, err := Decrypt(edata, tc.password)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			require.Equal(t, d, data)
		})
	}

	// test data of length in range of 0 to 64
	for i := 0; i <= 64; i++ {
		name := fmt.Sprintf("data length=%d", i)
		t.Run(name, func(t *testing.T) {
			data := randBytes(t, i)
			edata := makeEncryptedData(data, uint32(len(data)), []byte("pwd"))
			d, err := Decrypt(edata, []byte("pwd"))
			require.NoError(t, err)
			require.Equal(t, data, d)
		})
	}
}

// encrypts data, manually set the data length, so we could test invalid data length cases.
func makeEncryptedData(data []byte, dataLength uint32, password []byte) []byte {
	dataLenBytes := make([]byte, lengthSize)
	binary.LittleEndian.PutUint32(dataLenBytes, dataLength)

	ldata := append(dataLenBytes, data...)

	// Pads length + data with null to 32 bytes
	l := len(ldata) // hash + length + data
	n := l / blockSize
	m := l % blockSize
	if m > 0 {
		paddingNull := make([]byte, blockSize-m)
		ldata = append(ldata, paddingNull...)
		n++
	}

	// Hash(length+data+padding)
	dataHash := SumSHA256(ldata)

	// Initialize blocks with data hash
	blocks := []SHA256{dataHash}
	for i := 0; i < n; i++ {
		var b SHA256
		copy(b[:], ldata[i*blockSize:(i+1)*blockSize])
		blocks = append(blocks, b)
	}

	// Generates a nonce
	nonce := RandByte(nonceSize)
	// Hash the nonce
	hashNonce := SumSHA256(nonce)
	// Hash the password
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
	checkSum := SumSHA256(nonceAndDataBytes)

	return append(checkSum[:], nonceAndDataBytes...)
}
