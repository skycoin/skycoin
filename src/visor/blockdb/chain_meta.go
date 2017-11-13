package blockdb

import (
	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	// blockchain meta info bucket
	blockchainMetaBkt = []byte("blockchain_meta")
	// blockchain head sequence number
	headSeqKey = []byte("head_seq")
)

type chainMeta struct{}

func newChainMeta(db *dbutil.DB) (*chainMeta, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			blockchainMetaBkt,
		})
	}); err != nil {
		return nil, err
	}

	return &chainMeta{}, nil
}

func (m chainMeta) setHeadSeq(tx *bolt.Tx, seq uint64) error {
	return dbutil.PutBucketValue(tx, blockchainMetaBkt, headSeqKey, dbutil.Itob(seq))
}

func (m chainMeta) getHeadSeq(tx *bolt.Tx) (uint64, error) {
	v, err := dbutil.GetBucketValue(tx, blockchainMetaBkt, headSeqKey)
	if err != nil {
		return 0, err
	} else if v == nil {
		return 0, nil
	}

	return dbutil.Btoi(v), nil
}
