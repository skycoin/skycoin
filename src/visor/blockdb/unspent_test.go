package blockdb

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

type spending struct {
	ToAddr cipher.Address
	Coins  uint64
}

func randBytes(t *testing.T, n int) []byte { // nolint: unparam
	b := make([]byte, n)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return b
}

func randSHA256(t *testing.T) cipher.SHA256 {
	return cipher.SumSHA256(randBytes(t, 128))
}

func makeUxBody(t *testing.T) coin.UxBody {
	body, _ := makeUxBodyWithSecret(t)
	return body
}

func makeUxOut(t *testing.T) coin.UxOut {
	ux, _ := makeUxOutWithSecret(t)
	return ux
}

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(randBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) { // nolint: unparam
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func TestNewUnspentPool(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	up, err := NewUnspentPool(db)
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		length, err := dbutil.Len(tx, unspentPoolBkt)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)

		h, err := up.meta.getXorHash(tx)
		require.NoError(t, err)
		require.Equal(t, cipher.SHA256{}, h)
		return nil

	})
	require.NoError(t, err)
}

func addUxOut(db *dbutil.DB, up *Unspents, ux coin.UxOut) error {
	return db.Update(func(tx *bolt.Tx) error {
		return up.pool.set(tx, ux.Hash(), ux)
	})
}

func TestUnspentPoolGet(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	testCases := []struct {
		name     string
		unspents coin.UxArray
		hash     cipher.SHA256
		ux       *coin.UxOut
	}{
		{
			"not exist",
			uxs[:2],
			uxs[2].Hash(),
			nil,
		},
		{
			"find one",
			uxs[:2],
			uxs[1].Hash(),
			&uxs[1],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			require.NoError(t, err)
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			err = db.View(func(tx *bolt.Tx) error {
				ux, err := up.Get(tx, tc.hash)
				require.NoError(t, err)
				require.Equal(t, tc.ux, ux)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestUnspentPoolLen(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	db, closedb := testutil.PrepareDB(t)
	defer closedb()

	up, err := NewUnspentPool(db)
	require.NoError(t, err)

	for _, ux := range uxs {
		err := addUxOut(db, up, ux)
		require.NoError(t, err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		length, err := up.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(5), length)
		return nil
	})
	require.NoError(t, err)
}

func TestUnspentPoolGetUxHash(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	db, closedb := testutil.PrepareDB(t)
	defer closedb()

	up, err := NewUnspentPool(db)
	require.NoError(t, err)

	for _, ux := range uxs {
		err := addUxOut(db, up, ux)
		require.NoError(t, err)
		err = db.Update(func(tx *bolt.Tx) error {
			uxHash, err := up.GetUxHash(tx)
			require.NoError(t, err)

			xorHash, err := up.meta.getXorHash(tx)
			require.NoError(t, err)
			require.Equal(t, xorHash.Hex(), uxHash.Hex())
			return nil
		})
		require.NoError(t, err)
	}
}

func TestUnspentPoolGetArray(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	up, err := NewUnspentPool(db)
	require.NoError(t, err)

	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		err = addUxOut(db, up, ux)
		require.NoError(t, err)
		uxs = append(uxs, ux)
	}

	outsideUx := makeUxOut(t)

	testCases := []struct {
		name     string
		hashes   []cipher.SHA256
		err      error
		unspents coin.UxArray
	}{
		{
			"get first",
			[]cipher.SHA256{uxs[0].Hash()},
			nil,
			uxs[:1],
		},
		{
			"get second",
			[]cipher.SHA256{uxs[1].Hash()},
			nil,
			uxs[1:2],
		},
		{
			"get two",
			[]cipher.SHA256{uxs[0].Hash(), uxs[1].Hash()},
			nil,
			uxs[0:2],
		},
		{
			"get not exist",
			[]cipher.SHA256{outsideUx.Hash()},
			fmt.Errorf("unspent output does not exist: %s", outsideUx.Hash().Hex()),
			coin.UxArray{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.View(func(tx *bolt.Tx) error {
				uxs, err := up.GetArray(tx, tc.hashes)
				require.Equal(t, tc.err, err)
				if err == nil {
					require.Equal(t, tc.unspents, uxs)
				}
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestUnspentPoolGetAll(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	testCases := []struct {
		name     string
		unspents coin.UxArray
		expect   coin.UxArray
	}{
		{
			"empty",
			coin.UxArray{},
			coin.UxArray{},
		},
		{
			"one unspent",
			uxs[:1],
			uxs[:1],
		},
		{
			"two unspent",
			uxs[:2],
			uxs[:2],
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			require.NoError(t, err)
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			err = db.View(func(tx *bolt.Tx) error {
				unspents, err := up.GetAll(tx)
				require.NoError(t, err)

				uxm := make(map[cipher.SHA256]struct{})
				for _, ux := range unspents {
					uxm[ux.Hash()] = struct{}{}
				}

				for _, ux := range tc.expect {
					_, ok := uxm[ux.Hash()]
					require.True(t, ok)
				}

				return nil
			})
			require.NoError(t, err)
		})
	}
}

func BenchmarkUnspentPoolGetAll(b *testing.B) {
	var t testing.T
	db, teardown := testutil.PrepareDB(&t)
	defer teardown()

	up, err := NewUnspentPool(db)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		ux := makeUxOut(&t)
		if err := addUxOut(db, up, ux); err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	for i := 0; i < b.N; i++ {
		err := db.View(func(tx *bolt.Tx) error {
			_, err = up.GetAll(tx)
			return err
		})
		if err != nil {
			b.Fatal(err)
		}
	}
	fmt.Println(time.Since(start))
}

func TestGetUnspentOfAddrs(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	uxs[4].Body.Address = uxs[0].Body.Address

	testCases := []struct {
		name     string
		unspents coin.UxArray
		addrs    []cipher.Address
		expect   coin.UxArray
	}{
		{
			"one one addr one unspent",
			uxs[:],
			[]cipher.Address{uxs[1].Body.Address},
			uxs[1:2],
		},
		{
			"one addr two unspents",
			uxs[:],
			[]cipher.Address{uxs[0].Body.Address},
			[]coin.UxOut{
				uxs[0],
				uxs[4],
			},
		},
		{
			"two addrs three unspents",
			uxs[:],
			[]cipher.Address{uxs[0].Body.Address, uxs[1].Body.Address},
			[]coin.UxOut{
				uxs[0],
				uxs[1],
				uxs[4],
			},
		},
		{
			"two addrs two unspents",
			uxs[:],
			[]cipher.Address{uxs[2].Body.Address, uxs[1].Body.Address},
			[]coin.UxOut{
				uxs[1],
				uxs[2],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			require.NoError(t, err)
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			var unspents coin.AddressUxOuts
			err = db.View(func(tx *bolt.Tx) error {
				var err error
				unspents, err = up.GetUnspentsOfAddrs(tx, tc.addrs)
				require.NoError(t, err)
				return nil
			})
			require.NoError(t, err)

			uxm := make(map[cipher.SHA256]struct{}, len(unspents))
			for _, uxs := range unspents {
				for _, ux := range uxs {
					uxm[ux.Hash()] = struct{}{}
				}
			}

			require.Equal(t, len(uxm), len(tc.expect))

			for _, ux := range tc.expect {
				_, ok := uxm[ux.Hash()]
				require.True(t, ok)
			}
		})
	}
}

func TestUnspentProcessBlock(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	tt := []struct {
		name   string
		init   coin.UxArray
		inputs coin.UxArray
	}{
		{
			"ok",
			uxs,
			uxs[:1],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closedb := testutil.PrepareDB(t)
			defer closedb()

			up, err := NewUnspentPool(db)
			require.NoError(t, err)

			for _, ux := range tc.init {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			txn := coin.Transaction{}
			for _, in := range tc.inputs {
				txn.PushInput(in.Hash())
			}

			a := testutil.MakeAddress()
			txn.PushOutput(a, 1e6, uxs[0].Body.Hours/2)

			var block *coin.Block
			var oldUxHash cipher.SHA256

			err = db.Update(func(tx *bolt.Tx) error {
				uxHash, err := up.GetUxHash(tx)
				require.NoError(t, err)

				block, err = coin.NewBlock(coin.Block{}, uint64(time.Now().Unix()), uxHash, coin.Transactions{txn}, feeCalc)
				require.NoError(t, err)

				oldUxHash, err = up.GetUxHash(tx)
				require.NoError(t, err)

				err = up.ProcessBlock(tx, &coin.SignedBlock{Block: *block})
				require.NoError(t, err)

				return nil
			})
			require.NoError(t, err)

			txOuts := coin.CreateUnspents(block.Head, txn)

			err = db.View(func(tx *bolt.Tx) error {
				// check that the inputs should already been deleted from unspent pool
				for _, in := range tc.inputs {
					v, err := up.Get(tx, in.Hash())
					require.NoError(t, err)
					require.Nil(t, v)
				}

				// check the new generate unspent
				hasKey, err := up.Contains(tx, txOuts[0].Hash())
				require.NoError(t, err)
				require.True(t, hasKey)

				// check uxHash
				for _, in := range tc.inputs {
					oldUxHash = oldUxHash.Xor(in.SnapshotHash())
				}

				uxHash := oldUxHash.Xor(txOuts[0].SnapshotHash())
				newUxHash, err := up.GetUxHash(tx)
				require.NoError(t, err)
				require.Equal(t, uxHash.Hex(), newUxHash.Hex())

				return nil
			})
			require.NoError(t, err)
		})
	}

}
