package testutil

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()
// PrepareDB prepares test instance of the DB
func PrepareDB(t *testing.T) (*bolt.DB, func()) {
	f, err := ioutil.TempFile("", "testdb")
	require.NoError(t, err)

	db, err := bolt.Open(f.Name(), 0700, nil)
	require.NoError(t, err)

	return db, func() {
		db.Close()
		os.Remove(f.Name())
	}
}
// RequireError checks that the error exists as required
func RequireError(t *testing.T, err error, msg string) {
	t.Helper()
	require.Error(t, err)
	require.Equal(t, msg, err.Error())
}

// MakeAddress generates new cipher.Address
func MakeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

// RandBytes random bytes with n length
func RandBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return b
}

// RandSHA256 generates random SHA256
func RandSHA256(t *testing.T) cipher.SHA256 {
	return cipher.SumSHA256(RandBytes(t, 128))
}
