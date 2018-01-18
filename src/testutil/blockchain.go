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
	"github.com/skycoin/skycoin/src/coin"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

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

func RequireError(t *testing.T, err error, msg string) {
	t.Helper()
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

func RandSHA256(t *testing.T) cipher.SHA256 {
	return cipher.SumSHA256(RandBytes(t, 128))
}

func SHA256FromHex(t *testing.T, hex string) cipher.SHA256 {
	sha, err := cipher.SHA256FromHex(hex)
	require.NoError(t, err)
	return sha
}
func MakeTransaction(t *testing.T) coin.Transaction {
	tx, _ := makeTransactionWithSecret(t)
	return tx
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeTransactionWithSecret(t *testing.T) (coin.Transaction, cipher.SecKey) {
	tx := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)

	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{s})
	tx.PushOutput(makeAddress(), 1e6, 50)
	tx.PushOutput(makeAddress(), 5e6, 50)
	tx.UpdateHeader()
	return tx, s
}

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}
