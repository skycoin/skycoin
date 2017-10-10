package historydb

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/bucket"
	"github.com/stretchr/testify/assert"
)

func TestNewHistoryMeta(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	hm, err := newHistoryMeta(db)
	assert.Nil(t, err)
	db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte("history_meta"))
		assert.NotNil(t, bkt)
		return nil
	})

	v := hm.v.Get(parsedHeightKey)
	assert.Nil(t, v)
}

func TestHistoryMetaGetParsedHeight(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	hm, err := newHistoryMeta(db)
	assert.Nil(t, err)

	assert.Equal(t, int64(-1), hm.ParsedHeight())

	assert.Nil(t, hm.v.Put(parsedHeightKey, bucket.Itob(10)))
	assert.Equal(t, int64(10), hm.ParsedHeight())
}

func TestHistoryMetaSetParsedHeight(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	hm, err := newHistoryMeta(db)
	assert.Nil(t, err)
	assert.Nil(t, hm.SetParsedHeight(0))
	assert.Equal(t, uint64(0), bucket.Btoi(hm.v.Get(parsedHeightKey)))

	// set 10
	hm.SetParsedHeight(10)
	assert.Equal(t, uint64(10), bucket.Btoi(hm.v.Get(parsedHeightKey)))
}
