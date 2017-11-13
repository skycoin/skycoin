package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	historyMetaBkt  = []byte("history_meta")
	parsedHeightKey = []byte("parsed_height")
)

// historyMeta bucket for storing block history meta info
type historyMeta struct{}

func newHistoryMeta(db *dbutil.DB) (*historyMeta, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			historyMetaBkt,
		})
	}); err != nil {
		return nil, err
	}

	return &historyMeta{}, nil
}

// Height returns history parsed height, if no block was parsed, return -1.
func (hm *historyMeta) ParsedHeight(tx *bolt.Tx) (int64, error) {
	v, err := dbutil.GetBucketValue(tx, historyMetaBkt, parsedHeightKey)
	if err != nil {
		return 0, err
	} else if v == nil {
		return -1, nil
	}

	return int64(dbutil.Btoi(v)), nil
}

// SetParsedHeight updates history parsed height
func (hm *historyMeta) SetParsedHeight(tx *bolt.Tx, h uint64) error {
	return dbutil.PutBucketValue(tx, historyMetaBkt, parsedHeightKey, dbutil.Itob(h))
}

// IsEmpty checks if history meta bucket is empty
func (hm *historyMeta) IsEmpty(tx *bolt.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, historyMetaBkt)
}

// Reset resets the bucket
func (hm *historyMeta) Reset(tx *bolt.Tx) error {
	return dbutil.Reset(tx, historyMetaBkt)
}
