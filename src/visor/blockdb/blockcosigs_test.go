package blockdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

func TestBlockCosigsBucket(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	gb := makeGenesisBlock(t)
	_, sk1 := cipher.GenerateKeyPair()
	_, sk2 := cipher.GenerateKeyPair()
	h := gb.HashHeader()
	c1, err := coin.SignBlockCosig(h, 1, sk1)
	require.NoError(t, err)
	c2, err := coin.SignBlockCosig(h, 3, sk2)
	require.NoError(t, err)
	gb.Cosigs = []coin.BlockCosig{c1, c2}

	other := cipher.SumSHA256([]byte("no such block"))

	err = db.Update("", func(tx *dbutil.Tx) error {
		return bc.AddBlock(tx, &gb)
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		cs, err := bc.GetBlockCosigs(tx, h)
		require.NoError(t, err)
		require.Equal(t, gb.Cosigs, cs)

		cs, err = bc.GetBlockCosigs(tx, other)
		require.NoError(t, err)
		require.Nil(t, cs)

		b, err := bc.GetSignedBlockBySeq(tx, 0)
		require.NoError(t, err)
		require.Equal(t, gb.Cosigs, b.Cosigs)

		b, err = bc.GetSignedBlockByHash(tx, h)
		require.NoError(t, err)
		require.Equal(t, gb.Cosigs, b.Cosigs)
		return nil
	})
	require.NoError(t, err)
}

// A block without cosignatures stores nothing in the bucket, and a database
// from before the bucket existed, opened read-only, still reads.
func TestBlockchainWithoutCosigsBucket(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	gb := makeGenesisBlock(t)
	err = db.Update("", func(tx *dbutil.Tx) error {
		if err := bc.AddBlock(tx, &gb); err != nil {
			return err
		}
		n := 0
		if err := dbutil.ForEach(tx, BlockCosigsBkt, func(_, _ []byte) error {
			n++
			return nil
		}); err != nil {
			return err
		}
		require.Zero(t, n)
		return tx.DeleteBucket(BlockCosigsBkt)
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		b, err := bc.GetSignedBlockBySeq(tx, 0)
		require.NoError(t, err)
		require.Nil(t, b.Cosigs)
		return nil
	})
	require.NoError(t, err)
}
