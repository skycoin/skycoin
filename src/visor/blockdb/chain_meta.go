package blockdb

import (
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// BlockchainMetaBkt holds blockchain metadata
	BlockchainMetaBkt = []byte("blockchain_meta")
	// blockchain head sequence number
	headSeqKey = []byte("head_seq")
)

type chainMeta struct{}

func (m chainMeta) SetHeadSeq(tx *dbutil.Tx, seq uint64) error {
	return dbutil.PutBucketValue(tx, BlockchainMetaBkt, headSeqKey, dbutil.Itob(seq))
}

func (m chainMeta) GetHeadSeq(tx *dbutil.Tx) (uint64, bool, error) {
	v, err := dbutil.GetBucketValue(tx, BlockchainMetaBkt, headSeqKey)
	if err != nil {
		return 0, false, err
	} else if v == nil {
		return 0, false, nil
	}

	return dbutil.Btoi(v), true, nil
}
