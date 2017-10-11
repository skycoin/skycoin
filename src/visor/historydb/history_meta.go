package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	parsedHeightKey = []byte("parsed_height")
)

// historyMeta bucket for storing block history meta info
type historyMeta struct {
	v *bucket.Bucket
}

func newHistoryMeta(db *bolt.DB) (*historyMeta, error) {
	bkt, err := bucket.New([]byte("history_meta"), db)
	if err != nil {
		return nil, err
	}
	return &historyMeta{v: bkt}, nil
}

// Height returns history parsed height, if no block was parsed, return -1.
func (hm *historyMeta) ParsedHeight() int64 {
	if v := hm.v.Get(parsedHeightKey); v != nil {
		return int64(bucket.Btoi(v))
	}
	return -1
}

// SetParsedHeight updates history parsed height
func (hm *historyMeta) SetParsedHeight(h uint64) error {
	return hm.v.Put(parsedHeightKey, bucket.Itob(h))
}

// IsEmpty checks if history meta bucket is empty
func (hm *historyMeta) IsEmpty() bool {
	return hm.v.IsEmpty()
}

// Reset resets the bucket
func (hm *historyMeta) Reset() error {
	return hm.v.Reset()
}
