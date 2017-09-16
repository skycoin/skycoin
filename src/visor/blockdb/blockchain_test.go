package blockdb

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/assert"
)

func _feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func TestNewBlockchain(t *testing.T) {
	db, td := testutil.PrepareDB(t)
	defer td()

	bc, err := NewBlockchain(db)
	assert.Nil(t, err)

	assert.NotNil(t, bc.db)
	assert.NotNil(t, bc.UnspentPool())
	assert.NotNil(t, bc.meta)

	// check the existence of buckets
	db.View(func(tx *bolt.Tx) error {
		assert.NotNil(t, tx.Bucket([]byte("unspent_pool")))
		assert.NotNil(t, tx.Bucket([]byte("unspent_meta")))
		assert.NotNil(t, tx.Bucket([]byte("blockchain_meta")))
		return nil
	})
}
