package historydb

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

func TestHistoryMetaGetSetParsedHeight(t *testing.T) {
	db, td := prepareDB(t)
	defer td()

	hm := &historyMeta{}

	err := db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.parsedBlockSeq(tx)
		require.NoError(t, err)
		require.False(t, ok)
		require.Equal(t, uint64(0), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
		err := hm.setParsedBlockSeq(tx, 10)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.parsedBlockSeq(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, uint64(10), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
		err := hm.setParsedBlockSeq(tx, 0)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		height, ok, err := hm.parsedBlockSeq(tx)
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, uint64(0), height)
		return err
	})
	require.NoError(t, err)

}
