package historydb

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/stretchr/testify/assert"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (*bolt.DB, func(), error) {
	dbName := fmt.Sprintf("%ddb", rand.Int31n(10000))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), dbName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return nil, teardown, err
	}

	util.DataDir = tmpDir
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}

	teardown = func() {
		db.Close()
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	return db, teardown, nil
}

func TestGetLastTxs(t *testing.T) {
	testData := []uint64{0, 3, lastTxNum, lastTxNum + 10}
	for i := range testData {
		func(i uint64) {
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()
			txIns, err := newTransactionsBkt(db)
			if err != nil {
				t.Fatal(err)
			}

			var txs []cipher.SHA256
			for j := uint64(0); j < testData[i]; j++ {
				tx := makeTransaction()
				txs = append(txs, tx.Hash())
				if err := txIns.Add(&tx); err != nil {
					t.Fatal(err)
				}
			}
			if testData[i] > lastTxNum {
				txs = txs[len(txs)-lastTxNum:]
			}
			lastTxHash := txIns.GetLastTxs()
			assert.Equal(t, txs, lastTxHash)
		}(uint64(i))
	}
}

func makeTransaction() Transaction {
	tx := Transaction{}
	ux, s := makeUxOutWithSecret()
	tx.Tx.PushInput(ux.Hash())
	tx.Tx.SignInputs([]cipher.SecKey{s})
	tx.Tx.PushOutput(makeAddress(), 1e6, 50)
	tx.Tx.PushOutput(makeAddress(), 5e6, 50)
	tx.Tx.UpdateHeader()
	return tx
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeUxBodyWithSecret() (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(randBytes(128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret() (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret()
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	rand.Read(b)
	return b
}
