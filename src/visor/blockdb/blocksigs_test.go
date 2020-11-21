package blockdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

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
		h := testutil.RandSHA256(t)

		sig := cipher.MustSignHash(h, s)
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
			db, closeDB := prepareDB(t)
			defer closeDB()

			// init db
			err := db.Update("", func(tx *dbutil.Tx) error {
				bkt, err := tx.CreateBucketIfNotExists(BlockSigsBkt)
				require.NoError(t, err)
				for _, hs := range tc.init {
					err = bkt.Put(hs.hash[:], encoder.Serialize(hs.sig))
					require.NoError(t, err)
				}
				return nil
			})
			require.NoError(t, err)

			sigs := &blockSigs{}

			err = db.View("", func(tx *dbutil.Tx) error {
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

func TestBlockSigsAdd(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	_, s := cipher.GenerateKeyPair()
	h := testutil.RandSHA256(t)
	sig := cipher.MustSignHash(h, s)

	sigs := &blockSigs{}

	err := db.Update("", func(tx *dbutil.Tx) error {
		return sigs.Add(tx, h, sig)
	})
	require.NoError(t, err)

	// check the db
	err = db.View("", func(tx *dbutil.Tx) error {
		bkt := tx.Bucket(BlockSigsBkt)
		v := bkt.Get(h[:])
		require.NotNil(t, v)
		var s cipher.Sig
		err := encoder.DeserializeRawExact(v, &s)
		require.NoError(t, err)
		require.Equal(t, sig, s)
		return nil
	})
	require.NoError(t, err)
}
