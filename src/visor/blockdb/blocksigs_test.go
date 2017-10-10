package blockdb

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/testutil"
)

func TestNewBlockSigs(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	sigs, err := NewBlockSigs(db)
	require.NoError(t, err)
	require.NotNil(t, sigs)

	// check the bucket
	require.NotNil(t, sigs.Sigs)

	db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockSigsBkt)
		require.NotNil(t, bkt)
		return nil
	})
}

func TestBlockSigsGet(t *testing.T) {
	type hashSig struct {
		hash cipher.SHA256
		sig  cipher.Sig
	}

	type expect struct {
		exist bool
		sig   cipher.Sig
		err   error
	}

	hashSigs := []hashSig{}
	for i := 0; i < 5; i++ {
		_, s := cipher.GenerateKeyPair()
		h := randSHA256(t)

		sig := cipher.SignHash(h, s)
		hashSigs = append(hashSigs, hashSig{
			hash: h,
			sig:  sig,
		})
	}

	tt := []struct {
		name   string
		init   []hashSig
		hash   cipher.SHA256
		expect expect
	}{
		{
			"ok",
			hashSigs[:],
			hashSigs[0].hash,
			expect{
				true,
				hashSigs[0].sig,
				nil,
			},
		},
		{
			"not exist",
			hashSigs[1:],
			hashSigs[0].hash,
			expect{
				false,
				cipher.Sig{},
				nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			// init db
			db.Update(func(tx *bolt.Tx) error {
				bkt, err := tx.CreateBucketIfNotExists(blockSigsBkt)
				require.NoError(t, err)
				for _, hs := range tc.init {
					err = bkt.Put(hs.hash[:], encoder.Serialize(hs.sig))
					require.NoError(t, err)
				}
				return nil
			})

			sigs, err := NewBlockSigs(db)
			require.NoError(t, err)
			sg, ok, err := sigs.Get(tc.hash)
			require.Equal(t, tc.expect.err, err)
			require.Equal(t, tc.expect.exist, ok)
			if ok {
				require.Equal(t, tc.expect.sig, sg)
			}
		})
	}
}

func TestBlockSigsAddWithTx(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	_, s := cipher.GenerateKeyPair()
	h := randSHA256(t)
	sig := cipher.SignHash(h, s)

	sigs, err := NewBlockSigs(db)
	require.NoError(t, err)

	db.Update(func(tx *bolt.Tx) error {
		return sigs.AddWithTx(tx, h, sig)
	})

	// check the db
	db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockSigsBkt)
		v := bkt.Get(h[:])
		require.NotNil(t, v)
		var s cipher.Sig
		err := encoder.DeserializeRaw(v, &s)
		require.NoError(t, err)
		require.Equal(t, sig, s)
		return nil
	})
}
