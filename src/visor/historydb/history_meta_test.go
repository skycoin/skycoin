package historydb

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
)

func TestHistoryMetaGetSetParsedHeight(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	hm, err := newHistoryMeta(db)
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		height, err := hm.ParsedHeight(tx)
		require.NoError(t, err)
		require.Equal(t, int64(-1), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		err := hm.SetParsedHeight(tx, 10)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		height, err := hm.ParsedHeight(tx)
		require.NoError(t, err)
		require.Equal(t, int64(10), height)
		return err
	})
	require.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		err := hm.SetParsedHeight(tx, 0)
		require.NoError(t, err)
		return err
	})
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		height, err := hm.ParsedHeight(tx)
		require.NoError(t, err)
		require.Equal(t, int64(0), height)
		return err
	})
	require.NoError(t, err)

}
