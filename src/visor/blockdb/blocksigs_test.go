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

	_, err := newBlockSigs(db)
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockSigsBkt)
		require.NotNil(t, bkt)
		return nil
	})

	require.NoError(t, err)
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
			err := db.Update(func(tx *bolt.Tx) error {
				bkt, err := tx.CreateBucketIfNotExists(blockSigsBkt)
				require.NoError(t, err)
				for _, hs := range tc.init {
					err = bkt.Put(hs.hash[:], encoder.Serialize(hs.sig))
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			sigs, err := newBlockSigs(db)
			require.NoError(t, err)

			err = db.View(func(tx *bolt.Tx) error {
				sg, ok, err := sigs.Get(tx, tc.hash)
				require.Equal(t, tc.expect.err, err)
				require.Equal(t, tc.expect.exist, ok)
				if ok {
					require.Equal(t, tc.expect.sig, sg)
				}

				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestBlockSigsAddWithTx(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	_, s := cipher.GenerateKeyPair()
	h := randSHA256(t)
	sig := cipher.SignHash(h, s)

	sigs, err := newBlockSigs(db)
	require.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		return sigs.Add(tx, h, sig)
	})
	require.NoError(t, err)

	// check the db
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blockSigsBkt)
		v := bkt.Get(h[:])
		require.NotNil(t, v)
		var s cipher.Sig
		err := encoder.DeserializeRaw(v, &s)
		require.NoError(t, err)
		require.Equal(t, sig, s)
		return nil
	})
	require.NoError(t, err)
}
