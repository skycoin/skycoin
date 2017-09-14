package blockdb

import (
	"crypto/rand"
	"testing"

	"fmt"

	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/assert"
)

func randBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	assert.Equal(t, n, x) //end unit testing.
	assert.Nil(t, err)
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

func TestNewUnspentPool(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	up, err := NewUnspentPool(db)
	assert.Nil(t, err)

	assert.Equal(t, 0, up.pool.Len())
	v := up.meta.Get(xorhashKey)
	assert.Nil(t, v)
}

func addUxOut(up *UnspentPool, ux coin.UxOut) error {
	if err := up.db.Update(func(tx *bolt.Tx) error {
		_, err := up.addWithTx(tx, ux)
		return err
	}); err != nil {
		return err
	}
	up.addUxToCache([]coin.UxOut{ux})
	return nil
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
		ux       coin.UxOut
		exist    bool
	}{
		{
			"not exist",
			uxs[:2],
			uxs[2].Hash(),
			coin.UxOut{},
			false,
		},
		{
			"find one",
			uxs[:2],
			uxs[1].Hash(),
			uxs[1],
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			ux, ok := up.Get(tc.hash)
			assert.Nil(t, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.ux, ux)
			assert.Equal(t, tc.exist, ok)
		})
	}
}

func TestUnspentPoolGetArray(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	up, err := NewUnspentPool(db)
	assert.Nil(t, err)

	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		err = addUxOut(up, ux)
		assert.Nil(t, err)
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
			fmt.Errorf("unspent output of %s does not exist", outsideUx.Hash().Hex()),
			coin.UxArray{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uxs, err := up.GetArray(tc.hashes)
			assert.Equal(t, tc.err, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.unspents, uxs)
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
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			unspents, err := up.GetAll()
			assert.Nil(t, err)
			uxm := make(map[cipher.SHA256]byte)
			for _, ux := range unspents {
				uxm[ux.Hash()] = byte(1)
			}

			for _, ux := range tc.expect {
				_, ok := uxm[ux.Hash()]
				assert.True(t, ok)
			}
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
		if err := addUxOut(up, ux); err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	for i := 0; i < b.N; i++ {
		_, err = up.GetAll()
		if err != nil {
			b.Fatal(err)
		}
	}
	fmt.Println(time.Since(start))
}

func TestUnspentPoolDeleteWithTx(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	testCases := []struct {
		name         string
		unspents     coin.UxArray
		deleteHashes []cipher.SHA256
		error
		xorhash cipher.SHA256
	}{
		{
			"delete one ok",
			uxs[:2],
			[]cipher.SHA256{uxs[0].Hash()},
			nil,
			uxs[1].SnapshotHash(),
		},
		{
			"delete multilpe ok",
			uxs[:3],
			[]cipher.SHA256{uxs[0].Hash(), uxs[1].Hash()},
			nil,
			uxs[2].SnapshotHash(),
		},
		{
			"delete all ok",
			uxs[:3],
			[]cipher.SHA256{uxs[0].Hash(), uxs[1].Hash(), uxs[2].Hash()},
			nil,
			cipher.SHA256{},
		},
		{
			"delete middle one",
			uxs[:3],
			[]cipher.SHA256{uxs[1].Hash()},
			nil,
			func() cipher.SHA256 {
				h := uxs[0].SnapshotHash()
				return h.Xor(uxs[2].SnapshotHash())
			}(),
		},
		{
			"delete unknow hash",
			uxs[:2],
			[]cipher.SHA256{uxs[2].Hash()},
			nil,
			func() cipher.SHA256 {
				h := uxs[0].SnapshotHash()
				return h.Xor(uxs[1].SnapshotHash())
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			err = up.db.Update(func(tx *bolt.Tx) error {
				if _, err := up.deleteWithTx(tx, tc.deleteHashes); err != nil {
					return err
				}

				meta := unspentMeta{tx.Bucket(up.meta.Name)}
				xorhash, err := meta.getXorHash()
				assert.Nil(t, err)

				assert.Equal(t, tc.xorhash, xorhash)

				uxOut := uxOuts{tx.Bucket(up.pool.Name)}

				for _, hash := range tc.deleteHashes {
					_, ok, err := uxOut.get(hash)
					assert.Nil(t, err)
					assert.False(t, ok)
				}
				return nil
			})
			assert.Equal(t, tc.error, err)
		})
	}
}

func TestGetUnspentOfAddr(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	uxs[4].Body.Address = uxs[0].Body.Address

	testCases := []struct {
		name     string
		unspents coin.UxArray
		addr     cipher.Address
		expect   coin.UxArray
	}{
		{
			"one unspent",
			uxs[:],
			uxs[1].Body.Address,
			uxs[1:2],
		},
		{
			"two unspents",
			uxs[:],
			uxs[0].Body.Address,
			[]coin.UxOut{
				uxs[0],
				uxs[4],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			unspents := up.GetUnspentsOfAddr(tc.addr)
			assert.Nil(t, err)
			uxm := make(map[cipher.SHA256]byte, len(unspents))
			for _, ux := range unspents {
				uxm[ux.Hash()] = byte(1)
			}

			assert.Equal(t, len(uxm), len(tc.expect))

			for _, ux := range tc.expect {
				_, ok := uxm[ux.Hash()]
				assert.True(t, ok)
			}
		})
	}
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
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			unspents := up.GetUnspentsOfAddrs(tc.addrs)
			assert.Nil(t, err)
			uxm := make(map[cipher.SHA256]byte, len(unspents))
			for _, uxs := range unspents {
				for _, ux := range uxs {
					uxm[ux.Hash()] = byte(1)
				}
			}

			assert.Equal(t, len(uxm), len(tc.expect))

			for _, ux := range tc.expect {
				_, ok := uxm[ux.Hash()]
				assert.True(t, ok)
			}
		})
	}
}

func TestCollides(t *testing.T) {
	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		uxs = append(uxs, ux)
	}

	testCases := []struct {
		name       string
		unspents   coin.UxArray
		hashes     []cipher.SHA256
		isCollides bool
	}{
		{
			"no collides",
			uxs[:2],
			[]cipher.SHA256{uxs[2].Hash()},
			false,
		},
		{
			"one collides",
			uxs[:2],
			[]cipher.SHA256{uxs[1].Hash()},
			true,
		},
		{
			"multiple collides",
			uxs[:3],
			[]cipher.SHA256{uxs[1].Hash(), uxs[0].Hash()},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := testutil.PrepareDB(t)
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, addUxOut(up, ux))
			}

			assert.Equal(t, tc.isCollides, up.Collides(tc.hashes))
		})
	}
}
