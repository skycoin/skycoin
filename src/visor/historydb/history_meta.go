package historydb

import (
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// HistoryMetaBkt holds history metadata
	HistoryMetaBkt  = []byte("history_meta")
	parsedHeightKey = []byte("parsed_height")
)

// historyMeta bucket for storing block history meta info
type historyMeta struct{}

// Height returns history parsed height, if no block was parsed, return -1.
func (hm *historyMeta) ParsedHeight(tx *dbutil.Tx) (int64, error) {
	v, err := dbutil.GetBucketValue(tx, HistoryMetaBkt, parsedHeightKey)
	if err != nil {
		return 0, err
	} else if v == nil {
		return -1, nil
	}

	return int64(dbutil.Btoi(v)), nil
}

// SetParsedHeight updates history parsed height
func (hm *historyMeta) SetParsedHeight(tx *dbutil.Tx, h uint64) error {
	return dbutil.PutBucketValue(tx, HistoryMetaBkt, parsedHeightKey, dbutil.Itob(h))
}

// IsEmpty checks if history meta bucket is empty
func (hm *historyMeta) IsEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, HistoryMetaBkt)
}

// Reset resets the bucket
func (hm *historyMeta) Reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, HistoryMetaBkt)
}
