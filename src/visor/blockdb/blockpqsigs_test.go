package blockdb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// A database from before the post-quantum code has no block_pq_sigs bucket, and a
// node that opens it read-only cannot create one. Reading its blocks must still work.
func TestBlockchainWithoutPQSigsBucket(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	gb := makeGenesisBlock(t)
	err = db.Update("", func(tx *dbutil.Tx) error {
		if err := bc.AddBlock(tx, &gb); err != nil {
			return err
		}
		return tx.DeleteBucket(BlockPQSigsBkt)
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		require.False(t, dbutil.Exists(tx, BlockPQSigsBkt))

		head, err := bc.Head(tx)
		require.NoError(t, err)
		require.Nil(t, head.PQSig)

		b, err := bc.GetSignedBlockBySeq(tx, 0)
		require.NoError(t, err)
		require.Equal(t, gb.HashHeader(), b.HashHeader())
		require.Nil(t, b.PQSig)

		sig, err := bc.GetBlockPQ(tx, &b.Block)
		require.NoError(t, err)
		require.Nil(t, sig)
		return nil
	})
	require.NoError(t, err)
}
