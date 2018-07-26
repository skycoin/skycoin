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
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

// PrepareDB creates and opens a temporary test DB and returns it with a cleanup callback
func PrepareDB(t *testing.T) (*dbutil.DB, func()) {
	f, err := ioutil.TempFile("", "testdb")
	require.NoError(t, err)

	db, err := bolt.Open(f.Name(), 0700, nil)
	require.NoError(t, err)

	return dbutil.WrapDB(db), func() {
		db.Close()
		f.Close()
		os.Remove(f.Name())
	}
}

// RequireError requires that an error is not nil and that its message matches
func RequireError(t *testing.T, err error, msg string) {
	t.Helper()
	require.Error(t, err)
	require.Equal(t, msg, err.Error())
}

// MakeAddress creates a cipher.Address
func MakeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

// RandBytes returns n random bytes
func RandBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return b
}

// RandSHA256 returns a random SHA256 hash
func RandSHA256(t *testing.T) cipher.SHA256 {
	return cipher.SumSHA256(RandBytes(t, 128))
}

// SHA256FromHex converts an SHA256 hex string to a cipher.SHA256
func SHA256FromHex(t *testing.T, hex string) cipher.SHA256 {
	sha, err := cipher.SHA256FromHex(hex)
	require.NoError(t, err)
	return sha
}

// RandSig returns a random cipher.Sig
func RandSig(t *testing.T) cipher.Sig {
	return cipher.NewSig(RandBytes(t, 65))
}
