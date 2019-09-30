package blockdb

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

func makeUxBody(t *testing.T) coin.UxBody {
	p, _ := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}
}

func makeUxOut(t *testing.T) coin.UxOut {
	body := makeUxBody(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}
}

func TestNewUnspentPool(t *testing.T) {
	db, teardown := prepareDB(t)
	defer teardown()

	up := NewUnspentPool()

	err := db.View("", func(tx *dbutil.Tx) error {
		length, err := dbutil.Len(tx, UnspentPoolBkt)
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
	return db.Update("", func(tx *dbutil.Tx) error {
		if err := up.pool.put(tx, ux.Hash(), ux); err != nil {
			return err
		}

		return up.poolAddrIndex.adjust(tx, ux.Body.Address, []cipher.SHA256{ux.Hash()}, nil)
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
			db, teardown := prepareDB(t)
			defer teardown()

			up := NewUnspentPool()
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			err := db.View("", func(tx *dbutil.Tx) error {
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

	db, closedb := prepareDB(t)
	defer closedb()

	up := NewUnspentPool()

	for _, ux := range uxs {
		err := addUxOut(db, up, ux)
		require.NoError(t, err)
	}

	err := db.View("", func(tx *dbutil.Tx) error {
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

	db, closedb := prepareDB(t)
	defer closedb()

	up := NewUnspentPool()

	for _, ux := range uxs {
		err := addUxOut(db, up, ux)
		require.NoError(t, err)
		err = db.Update("", func(tx *dbutil.Tx) error {
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
	db, teardown := prepareDB(t)
	defer teardown()

	up := NewUnspentPool()

	var uxs coin.UxArray
	for i := 0; i < 5; i++ {
		ux := makeUxOut(t)
		err := addUxOut(db, up, ux)
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
			NewErrUnspentNotExist(outsideUx.Hash().Hex()),
			coin.UxArray{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.View("", func(tx *dbutil.Tx) error {
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
			db, teardown := prepareDB(t)
			defer teardown()

			up := NewUnspentPool()
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			err := db.View("", func(tx *dbutil.Tx) error {
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
	db, teardown := prepareDB(&t)
	defer teardown()

	up := NewUnspentPool()

	for i := 0; i < 1000; i++ {
		ux := makeUxOut(&t)
		if err := addUxOut(db, up, ux); err != nil {
			b.Fatal(err)
		}
	}

	start := time.Now()
	for i := 0; i < b.N; i++ {
		err := db.View("", func(tx *dbutil.Tx) error {
			_, err := up.GetAll(tx)
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
			name:     "one one addr one unspent",
			unspents: uxs[:],
			addrs:    []cipher.Address{uxs[1].Body.Address},
			expect:   uxs[1:2],
		},
		{
			name:     "one addr two unspents",
			unspents: uxs[:],
			addrs:    []cipher.Address{uxs[0].Body.Address},
			expect: []coin.UxOut{
				uxs[0],
				uxs[4],
			},
		},
		{
			name:     "two addrs three unspents",
			unspents: uxs[:],
			addrs:    []cipher.Address{uxs[0].Body.Address, uxs[1].Body.Address},
			expect: []coin.UxOut{
				uxs[0],
				uxs[1],
				uxs[4],
			},
		},
		{
			name:     "two addrs two unspents",
			unspents: uxs[:],
			addrs:    []cipher.Address{uxs[2].Body.Address, uxs[1].Body.Address},
			expect: []coin.UxOut{
				uxs[1],
				uxs[2],
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := prepareDB(t)
			defer teardown()

			up := NewUnspentPool()
			for _, ux := range tc.unspents {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			var unspents coin.AddressUxOuts
			err := db.View("", func(tx *dbutil.Tx) error {
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

	type testOutputs struct {
		addr  cipher.Address
		coins uint64
		hours uint64
	}

	addr := testutil.MakeAddress()

	tt := []struct {
		name          string
		init          coin.UxArray
		inputs        coin.UxArray
		outputs       []testOutputs
		nIndexedAddrs uint64
	}{
		{
			name:   "spend one create one",
			init:   uxs,
			inputs: uxs[:1],
			outputs: []testOutputs{
				{
					addr:  testutil.MakeAddress(),
					coins: 1e6,
					hours: uxs[0].Body.Hours / 2,
				},
			},
			nIndexedAddrs: 5,
		},

		{
			name:   "spend one create two",
			init:   uxs,
			inputs: uxs[:1],
			outputs: []testOutputs{
				{
					addr:  testutil.MakeAddress(),
					coins: 1e6 / 2,
					hours: uxs[0].Body.Hours / 4,
				},
				{
					addr:  testutil.MakeAddress(),
					coins: 1e6 / 2,
					hours: uxs[0].Body.Hours / 4,
				},
			},
			nIndexedAddrs: 6,
		},

		{
			name:   "spend one create three - two to the same new address and one to the spending address ",
			init:   uxs,
			inputs: uxs[:1],
			outputs: []testOutputs{
				{
					addr:  addr,
					coins: 1e6 / 4,
					hours: uxs[0].Body.Hours / 16,
				},
				{
					addr:  addr,
					coins: 1e6 / 4,
					hours: uxs[0].Body.Hours / 8,
				},
				{
					addr:  uxs[0].Body.Address,
					coins: 1e6 / 4,
					hours: uxs[0].Body.Hours / 8,
				},
			},
			nIndexedAddrs: 6,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closedb := prepareDB(t)
			defer closedb()

			up := NewUnspentPool()

			for _, ux := range tc.init {
				err := addUxOut(db, up, ux)
				require.NoError(t, err)
			}

			txn := coin.Transaction{}
			for _, in := range tc.inputs {
				err := txn.PushInput(in.Hash())
				require.NoError(t, err)
			}

			for _, o := range tc.outputs {
				err := txn.PushOutput(o.addr, o.coins, o.hours)
				require.NoError(t, err)
			}

			var block *coin.Block
			var oldUxHash cipher.SHA256

			err := db.Update("", func(tx *dbutil.Tx) error {
				uxHash, err := up.GetUxHash(tx)
				require.NoError(t, err)

				block, err = coin.NewBlock(coin.Block{}, uint64(time.Now().Unix()), uxHash, coin.Transactions{txn}, feeCalc)
				require.NoError(t, err)

				oldUxHash, err = up.GetUxHash(tx)
				require.NoError(t, err)

				err = up.ProcessBlock(tx, &coin.SignedBlock{
					Block: *block,
				})
				require.NoError(t, err)

				return nil
			})
			require.NoError(t, err)

			txOuts := coin.CreateUnspents(block.Head, txn)

			err = db.View("", func(tx *dbutil.Tx) error {
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
				for _, o := range txOuts[1:] {
					uxHash = uxHash.Xor(o.SnapshotHash())
				}

				newUxHash, err := up.GetUxHash(tx)
				require.NoError(t, err)
				require.Equal(t, uxHash.Hex(), newUxHash.Hex())

				// addr index height should equal the number of blocks added
				addrIndexHeight, ok, err := up.meta.getAddrIndexHeight(tx)
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, uint64(1), addrIndexHeight)

				// addr index should have 5 rows (5 initial addrs, 1 removed as input, 1 added as output)
				addrIndexLength, err := dbutil.Len(tx, UnspentPoolAddrIndexBkt)
				require.NoError(t, err)
				require.Equal(t, tc.nIndexedAddrs, addrIndexLength)

				// new outputs should be added to addr index cache
				expectedAddrHashes := make(map[cipher.Address][]cipher.SHA256)
				for _, o := range txOuts {
					expectedAddrHashes[o.Body.Address] = append(expectedAddrHashes[o.Body.Address], o.Hash())
				}

				for addr, hashes := range expectedAddrHashes {
					addrUxHashes, err := up.poolAddrIndex.get(tx, addr)
					require.NoError(t, err)

					require.Equal(t, len(hashes), len(addrUxHashes))

					sort.Slice(hashes, func(i, j int) bool {
						return bytes.Compare(hashes[i][:], hashes[j][:]) < 1
					})

					sort.Slice(addrUxHashes, func(i, j int) bool {
						return bytes.Compare(addrUxHashes[i][:], addrUxHashes[j][:]) < 1
					})

					require.Equal(t, hashes, addrUxHashes)
				}

				// used up inputs should be removed from addr index cache
				for _, o := range tc.inputs {
					// input addresses that appear in outputs should not be removed
					if _, ok := expectedAddrHashes[o.Body.Address]; ok {
						continue
					}

					addrUxHashes, err := up.poolAddrIndex.get(tx, o.Body.Address)
					require.NoError(t, err)
					require.Nil(t, addrUxHashes)
				}

				// none of the rows in the addr index should have empty arrays of hashes
				err = dbutil.ForEach(tx, UnspentPoolAddrIndexBkt, func(k, v []byte) error {
					_, err := cipher.AddressFromBytes(k)
					require.NoError(t, err)

					var uxHashes []cipher.SHA256
					err = encoder.DeserializeRawExact(v, &uxHashes)
					require.NoError(t, err)
					require.NotEmpty(t, uxHashes)

					return nil
				})
				require.NoError(t, err)

				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestUnspentPoolAddrIndex(t *testing.T) {
	addrs := make([]cipher.Address, 10)
	for i := range addrs {
		addrs[i] = testutil.MakeAddress()
	}

	hashes := make([]cipher.SHA256, 30)
	hashesMap := make(map[cipher.SHA256]struct{})
	for i := range hashes {
		hashes[i] = testutil.RandSHA256(t)
		_, ok := hashesMap[hashes[i]]
		require.False(t, ok)
		hashesMap[hashes[i]] = struct{}{}
	}

	type addrHashMap map[cipher.Address][]cipher.SHA256

	copyHashMap := func(hm addrHashMap) addrHashMap {
		out := make(addrHashMap, len(hm))

		for addr, hashes := range hm {
			copiedHashes := make([]cipher.SHA256, len(hashes))
			copy(copiedHashes[:], hashes[:])
			out[addr] = copiedHashes
		}

		return out
	}

	dup := func(h []cipher.SHA256) []cipher.SHA256 {
		i := make([]cipher.SHA256, len(h))
		copy(i[:], h[:])
		return i
	}

	cases := []struct {
		name      string
		init      addrHashMap
		add       addrHashMap
		remove    addrHashMap
		expect    addrHashMap
		putErr    error
		adjustErr error
	}{
		{
			name: "no initial, add only",
			add: addrHashMap{
				addrs[0]: dup(hashes[0:3]),
				addrs[1]: dup(hashes[3:6]),
			},
			expect: addrHashMap{
				addrs[0]: dup(hashes[0:3]),
				addrs[1]: dup(hashes[3:6]),
			},
		},

		{
			name: "initial, add and remove",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:3]),   // add one to here
				addrs[1]: dup(hashes[3:6]),   // remove one from here
				addrs[2]: dup(hashes[6:9]),   // add and remove one from here
				addrs[3]: dup(hashes[9:12]),  // remove all from here
				addrs[4]: dup(hashes[12:15]), // remove all from here and add one
			},
			add: addrHashMap{
				addrs[0]: dup(hashes[16:17]),
				addrs[2]: dup(hashes[17:18]),
				addrs[4]: dup(hashes[18:19]),
			},
			remove: addrHashMap{
				addrs[1]: dup(hashes[4:5]),
				addrs[2]: dup(hashes[6:7]),
				addrs[3]: dup(hashes[9:12]),
				addrs[4]: dup(hashes[12:15]),
			},
			expect: addrHashMap{
				addrs[0]: append(dup(hashes[0:3]), dup(hashes[16:17])...),
				addrs[1]: append(dup(hashes[3:4]), dup(hashes[5:6])...),
				addrs[2]: append(dup(hashes[7:9]), dup(hashes[17:18])...),
				addrs[4]: dup(hashes[18:19]),
			},
		},

		{
			name: "put error duplicate",
			init: addrHashMap{
				addrs[0]: []cipher.SHA256{hashes[0], hashes[0]},
			},
			putErr: errors.New("poolAddrIndex.put: hashes array contains duplicate"),
		},

		{
			name: "put error empty array",
			init: addrHashMap{
				addrs[0]: []cipher.SHA256{},
			},
			putErr: errors.New("poolAddrIndex.put cannot put empty hash array"),
		},

		{
			name: "adjust error removes have duplicates",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:1]),
			},
			remove: addrHashMap{
				addrs[0]: []cipher.SHA256{hashes[0], hashes[0]},
			},
			adjustErr: errors.New("poolAddrIndex.adjust: rmHashes contains duplicates"),
		},

		{
			name: "adjust error removing more than exists",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:1]),
			},
			remove: addrHashMap{
				addrs[0]: dup(hashes[0:2]),
			},
			adjustErr: errors.New("poolAddrIndex.adjust: rmHashes is longer than existingHashes"),
		},

		{
			name: "adjust error removing hash that does not exist",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:2]),
			},
			remove: addrHashMap{
				addrs[0]: []cipher.SHA256{hashes[0], hashes[11]},
			},
			adjustErr: fmt.Errorf("poolAddrIndex.adjust: rmHashes contains 1 hashes not indexed for address %s", addrs[0].String()),
		},

		{
			name: "adjust error hash in both add and remove",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:10]),
			},
			add: addrHashMap{
				addrs[0]: dup(hashes[4:5]),
			},
			remove: addrHashMap{
				addrs[0]: dup(hashes[1:5]),
			},
			adjustErr: errors.New("poolAddrIndex.adjust: hash appears in both addHashes and rmHashes"),
		},

		{
			name: "adjust error adding hash already indexed",
			init: addrHashMap{
				addrs[0]: dup(hashes[0:10]),
			},
			add: addrHashMap{
				addrs[0]: dup(hashes[4:5]),
			},
			adjustErr: fmt.Errorf("poolAddrIndex.adjust: uxout hash %s is already indexed for address %s", hashes[4].Hex(), addrs[0].String()),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := prepareDB(t)
			defer shutdown()

			tc.init = copyHashMap(tc.init)
			tc.add = copyHashMap(tc.add)
			tc.remove = copyHashMap(tc.remove)
			tc.expect = copyHashMap(tc.expect)

			p := &poolAddrIndex{}
			m := &unspentMeta{}

			// Initialize the data, test that put() works
			err := db.Update("", func(tx *dbutil.Tx) error {
				for addr, hashes := range tc.init {
					if err := p.put(tx, addr, hashes); err != nil {
						return err
					}
				}

				return m.setAddrIndexHeight(tx, uint64(len(tc.init)))
			})

			if tc.putErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.putErr, err)
				return
			}

			// Check the initialized data, test that get() works
			err = db.View("", func(tx *dbutil.Tx) error {
				length, err := dbutil.Len(tx, UnspentPoolAddrIndexBkt)
				require.NoError(t, err)
				require.Equal(t, uint64(len(tc.init)), length)

				height, ok, err := m.getAddrIndexHeight(tx)
				require.NoError(t, err)
				require.True(t, ok)
				require.Equal(t, length, height)

				for addr, expectHashes := range tc.init {
					hashes, err := p.get(tx, addr)
					require.NoError(t, err)
					require.Equal(t, expectHashes, hashes)
				}
				return nil
			})
			require.NoError(t, err)

			// Adjust the data, test that adjust() works
			err = db.Update("", func(tx *dbutil.Tx) error {
				for addr, addHashes := range tc.add {
					rmHashes := tc.remove[addr]
					delete(tc.remove, addr)

					err := p.adjust(tx, addr, addHashes, rmHashes)
					if err != nil {
						return err
					}
				}

				for addr, rmHashes := range tc.remove {
					err := p.adjust(tx, addr, nil, rmHashes)
					if err != nil {
						return err
					}
				}

				return nil
			})

			if tc.adjustErr == nil {
				require.NoError(t, err)
			} else {
				require.Equal(t, tc.adjustErr, err)
				return
			}

			addrHashes := make(addrHashMap)
			err = db.View("", func(tx *dbutil.Tx) error {
				return dbutil.ForEach(tx, UnspentPoolAddrIndexBkt, func(k, v []byte) error {
					addr, err := cipher.AddressFromBytes(k)
					require.NoError(t, err)

					var hashes []cipher.SHA256
					err = encoder.DeserializeRawExact(v, &hashes)
					require.NoError(t, err)

					sort.Slice(hashes, func(i, j int) bool {
						return bytes.Compare(hashes[i][:], hashes[j][:]) < 1
					})

					addrHashes[addr] = hashes

					return nil
				})
			})
			require.NoError(t, err)

			for addr, hashes := range tc.expect {
				sort.Slice(hashes, func(i, j int) bool {
					return bytes.Compare(hashes[i][:], hashes[j][:]) < 1
				})

				tc.expect[addr] = hashes
			}

			require.Equal(t, len(tc.expect), len(addrHashes))
			require.Equal(t, tc.expect, addrHashes)
		})
	}
}

func TestUnspentMaybeBuildIndexesNoIndexNoHead(t *testing.T) {
	// Test with a database that has no unspent addr index
	testUnspentMaybeBuildIndexes(t, 0, nil)
}

func TestUnspentMaybeBuildIndexesNoIndexHaveHead(t *testing.T) {
	// Test with a database that has no unspent addr index
	testUnspentMaybeBuildIndexes(t, 10, nil)
}

func TestUnspentMaybeBuildIndexesPartialIndex(t *testing.T) {
	// Test with a database that has an unspent addr index but the height is wrong
	headHeight := uint64(3)
	testUnspentMaybeBuildIndexes(t, headHeight+1, func(db *dbutil.DB) {
		up := NewUnspentPool()

		// Index the first few blocks
		err := db.Update("", func(tx *dbutil.Tx) error {
			if err := dbutil.CreateBuckets(tx, [][]byte{UnspentPoolAddrIndexBkt}); err != nil {
				return err
			}

			addrHashes := make(map[cipher.Address][]cipher.SHA256)

			if err := dbutil.ForEach(tx, UnspentPoolBkt, func(_, v []byte) error {
				var ux coin.UxOut
				if err := encoder.DeserializeRawExact(v, &ux); err != nil {
					return err
				}

				if ux.Head.BkSeq > headHeight {
					return nil
				}

				h := ux.Hash()
				addrHashes[ux.Body.Address] = append(addrHashes[ux.Body.Address], h)

				return nil
			}); err != nil {
				return err
			}

			for addr, hashes := range addrHashes {
				if err := up.poolAddrIndex.put(tx, addr, hashes); err != nil {
					return err
				}
			}

			return up.meta.setAddrIndexHeight(tx, headHeight)
		})
		require.NoError(t, err)
	})
}

func testUnspentMaybeBuildIndexes(t *testing.T, headIndex uint64, setupDB func(*dbutil.DB)) {
	db, shutdown := setupNoUnspentAddrIndexDB(t)
	defer shutdown()

	if setupDB != nil {
		setupDB(db)
	}

	u := NewUnspentPool()

	// Create the indexes
	err := db.Update("", func(tx *dbutil.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(UnspentPoolAddrIndexBkt); err != nil {
			return err
		}

		return u.MaybeBuildIndexes(tx, headIndex)
	})
	require.NoError(t, err)

	// Check the address->hashes index
	addrHashes := make(map[cipher.Address][]cipher.SHA256)
	err = db.View("", func(tx *dbutil.Tx) error {
		// Perform the unspent lookup the slow way, to confirm it matches the hashed data
		err := dbutil.ForEach(tx, UnspentPoolBkt, func(k, v []byte) error {
			hash, err := cipher.SHA256FromBytes(k)
			require.NoError(t, err)

			var ux coin.UxOut
			err = encoder.DeserializeRawExact(v, &ux)
			require.NoError(t, err)

			require.Equal(t, hash, ux.Hash())

			addrHashes[ux.Body.Address] = append(addrHashes[ux.Body.Address], hash)

			return nil
		})
		require.NoError(t, err)

		length, err := dbutil.Len(tx, UnspentPoolAddrIndexBkt)
		require.NoError(t, err)

		require.Equal(t, uint64(len(addrHashes)), length)

		height, ok, err := u.meta.getAddrIndexHeight(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, uint64(180), height)

		err = dbutil.ForEach(tx, UnspentPoolAddrIndexBkt, func(k, v []byte) error {
			addr, err := cipher.AddressFromBytes(k)
			require.NoError(t, err)

			var hashes []cipher.SHA256
			err = encoder.DeserializeRawExact(v, &hashes)
			require.NoError(t, err)

			expectedHashes, ok := addrHashes[addr]
			require.True(t, ok)

			sort.Slice(expectedHashes, func(i, j int) bool {
				return bytes.Compare(expectedHashes[i][:], expectedHashes[j][:]) < 1
			})

			sort.Slice(hashes, func(i, j int) bool {
				return bytes.Compare(hashes[i][:], hashes[j][:]) < 1
			})

			require.Equal(t, expectedHashes, hashes)

			delete(addrHashes, addr)

			return nil
		})
		require.NoError(t, err)

		require.Empty(t, addrHashes)

		return nil
	})
	require.NoError(t, err)
}

func TestUnspentMaybeBuildIndexesNoRebuild(t *testing.T) {
	// Set addrIndexHeight to head block height, but don't populate the addr index
	// Check that the addr index was not populated, so we know that the index did not get rebuilt if the height matches
	db, shutdown := setupNoUnspentAddrIndexDB(t)
	defer shutdown()

	u := NewUnspentPool()

	// Create the bucket and artificially set the indexed height, without indexing
	headSeq := uint64(180)
	err := db.Update("", func(tx *dbutil.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(UnspentPoolAddrIndexBkt); err != nil {
			return err
		}

		return u.meta.setAddrIndexHeight(tx, headSeq)
	})
	require.NoError(t, err)

	// Attempt to build index based upon the headSeq that we set
	err = db.Update("", func(tx *dbutil.Tx) error {
		return u.MaybeBuildIndexes(tx, headSeq)
	})
	require.NoError(t, err)

	// Check that the addr index is still empty, because the height was the same
	err = db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := u.meta.getAddrIndexHeight(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, headSeq, height)

		length, err := dbutil.Len(tx, UnspentPoolAddrIndexBkt)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)
		return nil
	})
	require.NoError(t, err)
}

func setupNoUnspentAddrIndexDB(t *testing.T) (*dbutil.DB, func()) {
	// Open a test database file that lacks UnspentPoolAddrIndexBkt,
	// copy it to a temp file and open a database around the temp file
	dbFilename := "./testdata/blockchain-180.no-unspent-addr-index.db"
	dbFile, err := os.Open(dbFilename)
	require.NoError(t, err)

	tmpFile, err := ioutil.TempFile("", "testdb")
	require.NoError(t, err)

	_, err = io.Copy(tmpFile, dbFile)
	require.NoError(t, err)

	err = dbFile.Close()
	require.NoError(t, err)

	err = tmpFile.Sync()
	require.NoError(t, err)

	boltDB, err := bolt.Open(tmpFile.Name(), 0700, nil)
	require.NoError(t, err)

	db := dbutil.WrapDB(boltDB)

	return db, func() {
		db.Close()
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
}

func TestAddressHashesFlatten(t *testing.T) {
	addrHashes := make(AddressHashes)

	require.Empty(t, addrHashes.Flatten())

	hashes := make([]cipher.SHA256, 10)
	for i := range hashes {
		hashes[i] = testutil.RandSHA256(t)
	}

	addrHashes = AddressHashes{
		testutil.MakeAddress(): hashes[0:2],
	}
	require.Equal(t, hashes[0:2], addrHashes.Flatten())

	addr1 := testutil.MakeAddress()
	addr2 := testutil.MakeAddress()
	addrHashes = AddressHashes{
		addr1: hashes[0:2],
		addr2: hashes[4:6],
	}

	expectedHashes := append(addrHashes[addr1], addrHashes[addr2]...)
	flattenedHashes := addrHashes.Flatten()

	sort.Slice(expectedHashes, func(a, b int) bool {
		return bytes.Compare(expectedHashes[a][:], expectedHashes[b][:]) < 0
	})
	sort.Slice(flattenedHashes, func(a, b int) bool {
		return bytes.Compare(flattenedHashes[a][:], flattenedHashes[b][:]) < 0
	})

	require.Equal(t, len(expectedHashes), len(flattenedHashes))
	require.Equal(t, expectedHashes, flattenedHashes)
}
