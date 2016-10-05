package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// blocks storage bucket, when parsing blockchain, the parsed blocks will be stored in this bucket.
// block header hash as key, block as value.
type blocks struct {
	bkt *bucket.Bucket
}

func newBlockBkt(db *bolt.DB) (*blocks, error) {
	bkt, err := bucket.New([]byte("blocks"), db)
	if err != nil {
		return nil, err
	}
	return &blocks{bkt}, nil
}

func (bs *blocks) Add(b *coin.Block) error {
	key := b.HashHeader()
	return bs.bkt.Put(key[:], encoder.Serialize(b))
}

// Get gets block of specifich hash header in blocks bucket
func (bs blocks) Get(hash cipher.SHA256) (*coin.Block, error) {
	bin := bs.bkt.Get(hash[:])
	if bin == nil {
		return nil, nil
	}
	var block coin.Block
	if err := encoder.DeserializeRaw(bin, &block); err != nil {
		return nil, err
	}
	return &block, nil
}
