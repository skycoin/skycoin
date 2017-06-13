package blockdb

import "testing"
import "github.com/stretchr/testify/assert"
import "github.com/boltdb/bolt"

func TestNewChainMeta(t *testing.T) {
	db, td, err := setup(t)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer td()

	_, err = NewChainMeta(db)
	assert.Nil(t, err)
	db.View(func(tx *bolt.Tx) error {
		chainMetaBkt := tx.Bucket(blockMetaKey)
		assert.NotNil(t, chainMetaBkt)
		return nil
	})
}

func TestGetHead(t *testing.T) {
	db, td, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer td()

	bc, err := NewChainMeta(db)
	assert.Nil(t, err)

	assert.Equal(t, int64(-1), bc.Head())
}
