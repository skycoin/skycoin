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

// Height returns history parsed height
func (hm *historyMeta) ParsedHeight(tx *dbutil.Tx) (uint64, bool, error) {
	v, err := dbutil.GetBucketValue(tx, HistoryMetaBkt, parsedHeightKey)
	if err != nil {
		return 0, false, err
	} else if v == nil {
		return 0, false, nil
	}

	return dbutil.Btoi(v), true, nil
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
