package encrypt

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	secp256k1 "github.com/SkycoinProject/skycoin/src/cipher/secp256k1-go"
	"github.com/SkycoinProject/skycoin/src/testutil"
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
			testutil.RandBytes(t, 1),
			nil,
			errors.New("missing password"),
		},
		{
			"data length=1  password is empty=false",
			testutil.RandBytes(t, 1),
			[]byte("key"),
			nil,
		},
		{
			"data length<32 password is empty=false",
			testutil.RandBytes(t, 2),
			[]byte("pwd"),
			nil,
		},
		{
			"data length=32 password is empty=false",
			testutil.RandBytes(t, 32),
			[]byte("pwd"),
			nil,
		},
		{
			"data length=2*32 password is empty=false",
			testutil.RandBytes(t, 64),
			[]byte("9JMkCPphe73NQvGhmab"),
			nil,
		},
		{
			"data length>2*32 password is empty=false",
			testutil.RandBytes(t, 65),
			[]byte("9JMkCPphe73NQvGhmab"),
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			edata, err := Sha256Xor{}.Encrypt(tc.data, tc.password)
			require.Equal(t, tc.err, err)
			if err != nil {
				return
			}

			n := (sha256XorDataLengthSize + len(tc.data)) / sha256XorBlockSize
			m := (sha256XorDataLengthSize + len(tc.data)) % sha256XorBlockSize
			if m > 0 {
				n++
			}

			rdata, err := base64.StdEncoding.DecodeString(string(edata))
			require.NoError(t, err)

			totalEncryptedDataLen := sha256XorChecksumSize + sha256XorNonceSize + 32 + n*sha256XorBlockSize // 32 is the hash data length
			require.Equal(t, totalEncryptedDataLen, len(rdata))
			var checksum cipher.SHA256
			copy(checksum[:], rdata[:sha256XorChecksumSize])
			require.Equal(t, checksum, cipher.SumSHA256(rdata[sha256XorChecksumSize:]))
		})
	}

	// test data of length range in 32 to 64, for testing the 32 bytes padding.
	pwd := []byte("pwd")
	for i := 33; i <= 64; i++ {
		name := fmt.Sprintf("data length=%d password is empty=false", i)
		t.Run(name, func(t *testing.T) {
			data := testutil.RandBytes(t, i)
			edata, err := Sha256Xor{}.Encrypt(data, pwd)
			require.NoError(t, err)

			n := (sha256XorDataLengthSize + len(data)) / sha256XorBlockSize
			m := (sha256XorDataLengthSize + len(data)) % sha256XorBlockSize
			if m > 0 {
				n++
			}

			rdata, err := base64.StdEncoding.DecodeString(string(edata))
			require.NoError(t, err)

			totalEncryptedDataLen := sha256XorChecksumSize + sha256XorNonceSize + 32 + n*sha256XorBlockSize // 32 is the hash data length
			require.Equal(t, totalEncryptedDataLen, len(rdata))
			var checksum cipher.SHA256
			copy(checksum[:], rdata[:sha256XorChecksumSize])
			require.Equal(t, checksum, cipher.SumSHA256(rdata[sha256XorChecksumSize:]))
		})
	}
}

func TestDecrypt(t *testing.T) {
	data := testutil.RandBytes(t, 32)
	tt := []struct {
		name          string
		encryptedData func() []byte // encrypted data
		password      []byte
		err           error
	}{
		{
			"invalid data length",
			func() []byte {
				return makeEncryptedData(t, data, 65, []byte("pwd"))
			},
			[]byte("pwd"),
			errors.New("invalid data length"),
		},
		{
			"invalid checksum",
			func() []byte {
				edata := makeEncryptedData(t, data, 32, []byte("pwd"))
				// Changes the encrypted data, so that the checksum could not match
				rd, err := base64.StdEncoding.DecodeString(string(edata))
				require.NoError(t, err)
				rd[len(rd)-1]++
				return []byte(base64.StdEncoding.EncodeToString(rd))
			},
			[]byte("pwd"),
			errors.New("invalid data, checksum is not matched"),
		},
		{
			"empty password",
			func() []byte {
				return makeEncryptedData(t, data, 32, []byte("pwd"))
			},
			[]byte(""),
			errors.New("missing password"),
		},
		{
			"nil password",
			func() []byte {
				return makeEncryptedData(t, data, 32, []byte("pwd"))
			},
			nil,
			errors.New("missing password"),
		},
		{
			"invalid password",
			func() []byte {
				return makeEncryptedData(t, data, 32, []byte("pwd"))
			},
			[]byte("wrong password"),
			errors.New("invalid password"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			edata := tc.encryptedData()
			d, err := Sha256Xor{}.Decrypt(edata, tc.password)
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
			data := testutil.RandBytes(t, i)
			edata := makeEncryptedData(t, data, uint32(len(data)), []byte("pwd"))
			d, err := Sha256Xor{}.Decrypt(edata, []byte("pwd"))
			require.NoError(t, err)
			require.Equal(t, data, d)
		})
	}
}

// encrypts data, manually set the data length, so we could test invalid data length cases.
func makeEncryptedData(t *testing.T, data []byte, dataLength uint32, password []byte) []byte {
	dataLenBytes := make([]byte, sha256XorDataLengthSize)
	binary.LittleEndian.PutUint32(dataLenBytes, dataLength)

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
	nonce := testutil.RandBytes(t, int(sha256XorNonceSize))
	// Hash the nonce
	hashNonce := cipher.SumSHA256(nonce)
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
	checkSum := cipher.SumSHA256(nonceAndDataBytes)
	rd := append(checkSum[:], nonceAndDataBytes...)
	enc := base64.StdEncoding
	buf := make([]byte, enc.EncodedLen(len(rd)))
	enc.Encode(buf, rd)
	return buf
}
