package historydb

import (
	"errors"
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

	txnHashes := make([]cipher.SHA256, len(txns))
	for i, x := range txns {
		txnHashes[i] = x.Hash()
	}

	testCases := []struct {
		name   string
		hash   cipher.SHA256
		expect *Transaction
	}{
		{
			name:   "get first",
			hash:   txnHashes[0],
			expect: &txns[0],
		},
		{
			name:   "get second",
			hash:   txnHashes[1],
			expect: &txns[1],
		},
		{
			name:   "not exist",
			hash:   txnHashes[2],
			expect: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()

			txsBkt := &transactions{}

			// init the bkt
			err := db.Update("", func(tx *dbutil.Tx) error {
				for _, txn := range txns[:2] {
					err := txsBkt.put(tx, &txn)
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			// get slice
			err = db.View("", func(tx *dbutil.Tx) error {
				ts, err := txsBkt.get(tx, tc.hash)
				require.NoError(t, err)
				require.Equal(t, tc.expect, ts)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestTransactionGetArray(t *testing.T) {
	txns := make([]Transaction, 0, 4)
	for i := 0; i < 4; i++ {
		txns = append(txns, makeTransaction(t))
	}

	txnHashes := make([]cipher.SHA256, len(txns))
	for i, x := range txns {
		txnHashes[i] = x.Hash()
	}

	testCases := []struct {
		name   string
		hashes []cipher.SHA256
		expect []Transaction
		err    error
	}{
		{
			name: "get one",
			hashes: []cipher.SHA256{
				txnHashes[0],
			},
			expect: txns[:1],
		},

		{
			name: "get two",
			hashes: []cipher.SHA256{
				txnHashes[0],
				txnHashes[1],
			},
			expect: txns[:2],
		},

		{
			name: "get all",
			hashes: []cipher.SHA256{
				txnHashes[0],
				txnHashes[1],
				txnHashes[2],
			},
			expect: txns[:3],
		},

		{
			name: "not exist",
			hashes: []cipher.SHA256{
				txnHashes[3],
			},
			err: errors.New("Transaction not found"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, td := prepareDB(t)
			defer td()
			txsBkt := &transactions{}

			// init the bkt
			err := db.Update("", func(tx *dbutil.Tx) error {
				for _, txn := range txns[:3] {
					err := txsBkt.put(tx, &txn)
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			// get slice
			err = db.View("", func(tx *dbutil.Tx) error {
				ts, err := txsBkt.getArray(tx, tc.hashes)
				if tc.err != nil {
					require.Equal(t, tc.err, err)
					return nil
				}
				require.NoError(t, err)
				require.Equal(t, tc.expect, ts)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func makeTransaction(t *testing.T) Transaction {
	txn := Transaction{}
	ux, s := makeUxOutWithSecret(t)

	err := txn.Txn.PushInput(ux.Hash())
	require.NoError(t, err)
	err = txn.Txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.Txn.PushOutput(makeAddress(), 5e6, 50)
	require.NoError(t, err)
	txn.Txn.SignInputs([]cipher.SecKey{s})
	err = txn.Txn.UpdateHeader()
	require.NoError(t, err)
	return txn
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
