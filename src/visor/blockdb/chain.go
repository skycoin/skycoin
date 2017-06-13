package blockdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	blockMetaKey = []byte("blockchain_meta")
	blockHeadKey = []byte("blockchain_head")
)

var (
	// InitialHead represents empty hash head
	InitialHead = cipher.SHA256{}
)

// ChainMeta contains meta info of blockchain
type ChainMeta struct {
	v *bucket.Bucket // blockchain meta info bucket
}

// NewChainMeta creates blockchain instance
func NewChainMeta(db *bolt.DB) (*ChainMeta, error) {
	// create blockchain meta bucket if not exist
	meta, err := bucket.New(blockMetaKey, db)
	if err != nil {
		return nil, fmt.Errorf("create blockchain meta info bucket failed: %v", err)
	}

	return &ChainMeta{
		v: meta,
	}, nil
}

// Head returns the block head hash
func (cm *ChainMeta) Head() cipher.SHA256 {
	v := cm.v.Get(blockHeadKey)
	if v == nil {
		return cipher.SHA256{}
	}

	return cipher.MustSHA256FromHex(string(v))
}

// UpdateHead updates the head of the blockchain
func (cm *ChainMeta) UpdateHead(head cipher.SHA256) error {
	return cm.v.Put(blockHeadKey, []byte(head.Hex()))
}
