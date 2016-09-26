package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// UxOut augment coin.UxOut struct
type UxOut struct {
	Out           coin.UxOut
	SpentTxID     cipher.SHA256 // id of tx which spent this output.
	SpentBlockSeq uint64        // block seq that spent the output.
}

func (o UxOut) Hash() cipher.SHA256 {
	return o.Out.Hash()
}

// UxOuts bucket stores outputs, UxOut hash as key and Output as value.
type UxOuts struct {
	bkt *bucket.Bucket
}

func newOutputsBkt(db *bolt.DB) (*UxOuts, error) {
	bkt, err := bucket.New([]byte("uxouts"), db)
	if err != nil {
		return nil, err
	}
	return &UxOuts{bkt}, nil
}

func (op *UxOuts) Set(out UxOut) error {
	key := out.Hash()
	bin := encoder.Serialize(out)
	return op.bkt.Put(key[:], bin)
}

func (op *UxOuts) Get(uxID cipher.SHA256) (*UxOut, error) {
	bin := op.bkt.Get(uxID[:])
	if bin == nil {
		return nil, nil
	}

	out := UxOut{}
	if err := encoder.DeserializeRaw(bin, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
