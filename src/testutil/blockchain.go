package testutil

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

// PrepareDB initializes a temporary bolt.db and provides a cleanup method to defer
func PrepareDB(t *testing.T) (*bolt.DB, func()) {
	f, err := ioutil.TempFile("", "testdb")
	require.Nil(t, err)

	db, err := bolt.Open(f.Name(), 0700, nil)
	require.Nil(t, err)

	return db, func() {
		db.Close()
		os.Remove(f.Name())
	}
}

// RequireError requires that an error occured and compares the error string
func RequireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	require.Equal(t, msg, err.Error())
}

// MakeAddress generates an address
func MakeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

// RandSHA256 returns a random SHA256 hash
func RandSHA256(t *testing.T) cipher.SHA256 {
	return cipher.SumSHA256(RandBytes(t, 128))
}

// RandBytes returns n random bytes
func RandBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	require.Equal(t, n, x)
	require.Nil(t, err)
	return b
}
