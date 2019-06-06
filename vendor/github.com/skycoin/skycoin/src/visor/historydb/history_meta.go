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

// parsedBlockSeq returns history parsed block seq
func (hm *historyMeta) parsedBlockSeq(tx *dbutil.Tx) (uint64, bool, error) {
	v, err := dbutil.GetBucketValue(tx, HistoryMetaBkt, parsedHeightKey)
	if err != nil {
		return 0, false, err
	} else if v == nil {
		return 0, false, nil
	}

	return dbutil.Btoi(v), true, nil
}

// setParsedBlockSeq updates history parsed block seq
func (hm *historyMeta) setParsedBlockSeq(tx *dbutil.Tx, h uint64) error {
	return dbutil.PutBucketValue(tx, HistoryMetaBkt, parsedHeightKey, dbutil.Itob(h))
}

// reset resets the bucket
func (hm *historyMeta) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, HistoryMetaBkt)
}
