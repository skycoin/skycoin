package historydb

import (
	"math/rand"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func TestGetLastTxs(t *testing.T) {
	testData := []uint64{0, 3, lastTxNum, lastTxNum + 10}
	for i := range testData {
		func(i uint64) {
			db, teardown := testutil.PrepareDB(t)
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

func TestTransactionGet(t *testing.T) {
	txs := make([]Transaction, 0, 3)
	for i := 0; i < 3; i++ {
		txs = append(txs, makeTransaction())
	}

	testCases := []struct {
		name   string
		hash   cipher.SHA256
		expect *Transaction
	}{
		{
			"get first",
			txs[0].Hash(),
			&txs[0],
		},
		{
			"get second",
			txs[1].Hash(),
			&txs[1],
		},
		{
			"not exist",
			txs[2].Hash(),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := testutil.PrepareDB(t)
			defer td()
			txsBkt, err := newTransactionsBkt(db)
			require.Nil(t, err)

			// init the bkt
			for _, tx := range txs[:2] {
				require.Nil(t, txsBkt.Add(&tx))
			}

			// get slice
			ts, err := txsBkt.Get(tc.hash)
			require.Nil(t, err)
			require.Equal(t, tc.expect, ts)
		})
	}
}

func TestTransactionGetSlice(t *testing.T) {
	txs := make([]Transaction, 0, 4)
	for i := 0; i < 4; i++ {
		txs = append(txs, makeTransaction())
	}

	testCases := []struct {
		name   string
		hashes []cipher.SHA256
		expect []Transaction
	}{
		{
			"get one",
			[]cipher.SHA256{
				txs[0].Hash(),
			},
			txs[:1],
		},
		{
			"get two",
			[]cipher.SHA256{
				txs[0].Hash(),
				txs[1].Hash(),
			},
			txs[:2],
		},
		{
			"get all",
			[]cipher.SHA256{
				txs[0].Hash(),
				txs[1].Hash(),
				txs[2].Hash(),
			},
			txs[:3],
		},
		{
			"not exist",
			[]cipher.SHA256{
				txs[3].Hash(),
			},
			[]Transaction{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := testutil.PrepareDB(t)
			defer td()
			txsBkt, err := newTransactionsBkt(db)
			require.Nil(t, err)

			// init the bkt
			for _, tx := range txs[:3] {
				require.Nil(t, txsBkt.Add(&tx))
			}

			// get slice
			ts, err := txsBkt.GetSlice(tc.hashes)
			require.Nil(t, err)
			require.Equal(t, tc.expect, ts)
		})
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
