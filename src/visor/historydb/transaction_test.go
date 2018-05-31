package historydb

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func TestTransactionGet(t *testing.T) {
	txns := make([]Transaction, 0, 3)
	for i := 0; i < 3; i++ {
		txns = append(txns, makeTransaction(t))
	}

	testCases := []struct {
		name   string
		hash   cipher.SHA256
		expect *Transaction
	}{
		{
			"get first",
			txns[0].Hash(),
			&txns[0],
		},
		{
			"get second",
			txns[1].Hash(),
			&txns[1],
		},
		{
			"not exist",
			txns[2].Hash(),
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()

			txsBkt, err := newTransactions(db)
			require.NoError(t, err)

			// init the bkt
			err = db.Update("", func(tx *dbutil.Tx) error {
				for _, txn := range txns[:2] {
					err := txsBkt.Add(tx, &txn)
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			// get slice
			err = db.View("", func(tx *dbutil.Tx) error {
				ts, err := txsBkt.Get(tx, tc.hash)
				require.NoError(t, err)
				require.Equal(t, tc.expect, ts)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestTransactionGetSlice(t *testing.T) {
	txns := make([]Transaction, 0, 4)
	for i := 0; i < 4; i++ {
		txns = append(txns, makeTransaction(t))
	}

	testCases := []struct {
		name   string
		hashes []cipher.SHA256
		expect []Transaction
	}{
		{
			"get one",
			[]cipher.SHA256{
				txns[0].Hash(),
			},
			txns[:1],
		},
		{
			"get two",
			[]cipher.SHA256{
				txns[0].Hash(),
				txns[1].Hash(),
			},
			txns[:2],
		},
		{
			"get all",
			[]cipher.SHA256{
				txns[0].Hash(),
				txns[1].Hash(),
				txns[2].Hash(),
			},
			txns[:3],
		},
		{
			"not exist",
			[]cipher.SHA256{
				txns[3].Hash(),
			},
			nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()
			txsBkt, err := newTransactions(db)
			require.NoError(t, err)

			// init the bkt
			err = db.Update("", func(tx *dbutil.Tx) error {
				for _, txn := range txns[:3] {
					err := txsBkt.Add(tx, &txn)
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			// get slice
			err = db.View("", func(tx *dbutil.Tx) error {
				ts, err := txsBkt.GetSlice(tx, tc.hashes)
				require.NoError(t, err)
				require.Equal(t, tc.expect, ts)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func makeTransaction(t *testing.T) Transaction {
	tx := Transaction{}
	ux, s := makeUxOutWithSecret(t)
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

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
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
