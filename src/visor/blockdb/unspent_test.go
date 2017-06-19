package blockdb

import (
	"crypto/rand"
	"testing"

	"strings"

	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
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
	db, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	up, err := NewUnspentPool(db)
	assert.Nil(t, err)

	assert.Equal(t, 0, up.pool.Len())
	v := up.meta.Get(xorhashKey)
	assert.Nil(t, v)
}

func TestUnspentPoolAdd(t *testing.T) {
	db, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	up, err := NewUnspentPool(db)
	assert.Nil(t, err)

	ux := makeUxOut(t)
	assert.Nil(t, up.Add(ux))
	assert.Equal(t, uint64(1), up.Len())

	v := up.pool.Get([]byte(ux.Hash().Hex()))
	assert.NotNil(t, v)
	var uxc coin.UxOut
	err = encoder.DeserializeRaw(v, &uxc)
	assert.Nil(t, err)
	assert.Equal(t, ux, uxc)

	xorhash, err := up.GetUxHash()
	assert.Nil(t, err)

	assert.NotEqual(t, cipher.SHA256{}, xorhash)
	assert.Equal(t, ux.SnapshotHash(), xorhash)

	// Duplicate add, return err
	err = up.Add(ux)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "twice into the unspent pool"))

	// Add second, must be ok
	ux2 := makeUxOut(t)
	assert.Nil(t, up.Add(ux2))
	assert.Equal(t, uint64(2), up.Len())
	var ux2c coin.UxOut
	v2 := up.pool.Get([]byte(ux2.Hash().Hex()))
	assert.NotNil(t, v2)
	assert.Nil(t, encoder.DeserializeRaw(v2, &ux2c))
	assert.Equal(t, ux2, ux2c)

	h := ux.SnapshotHash()
	h = h.Xor(ux2.SnapshotHash())
	uph, err := up.GetUxHash()
	assert.Nil(t, err)
	assert.Equal(t, h, uph)
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
			}

			ux, ok, err := up.Get(tc.hash)
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
	db, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	up, err := NewUnspentPool(db)
	assert.Nil(t, err)

	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		err = up.Add(ux)
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
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

func TestUnspentPoolDelete(t *testing.T) {
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
			}

			err = up.db.Update(func(tx *bolt.Tx) error {
				return up.delete(tx, tc.deleteHashes)
			})
			assert.Equal(t, tc.error, err)

			if err != nil {
				return
			}

			xh, err := up.GetUxHash()
			assert.Nil(t, err)
			assert.Equal(t, tc.xorhash, xh)

			for _, hash := range tc.deleteHashes {
				_, ok, err := up.Get(hash)
				assert.Nil(t, err)
				assert.False(t, ok)
			}
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
			}

			unspents, err := up.GetUnspentsOfAddr(tc.addr)
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
			}

			unspents, err := up.GetUnspentsOfAddrs(tc.addrs)
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
			db, teardown, err := setup(t)
			if err != nil {
				t.Fatal(err)
			}
			defer teardown()

			up, err := NewUnspentPool(db)
			assert.Nil(t, err)
			for _, ux := range tc.unspents {
				assert.Nil(t, up.Add(ux))
			}

			assert.Equal(t, tc.isCollides, up.Collides(tc.hashes))
		})
	}
}
