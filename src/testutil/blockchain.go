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

func RequireError(t *testing.T, err error, msg string) {
	require.Error(t, err)
	require.Equal(t, msg, err.Error())
}

func MakeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func RandBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return b
}
