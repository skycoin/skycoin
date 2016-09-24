package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// Output augment output struct,
type UxOut struct {
	Out           coin.UxOut
	SpentTxID     cipher.SHA256 // id of tx which spent this output.
	SpentBlockSeq uint64        // block seq that spent the output.
}

func (o UxOut) Hash() cipher.SHA256 {
	return o.Out.Hash()
}

// Outputs bucket stores outputs, outID as key and Output as value.
type Outputs struct {
	bkt *bucket.Bucket
}

func newOutputsBkt(db *bolt.DB) (*Outputs, error) {
	bkt, err := bucket.New([]byte("uxouts"), db)
	if err != nil {
		return nil, err
	}
	return &Outputs{bkt}, nil
}

func (op *Outputs) Set(out UxOut) error {
	key := out.Hash()
	bin := encoder.Serialize(out)
	return op.bkt.Put(key[:], bin)
}

func (op *Outputs) Get(uxID cipher.SHA256) (*UxOut, error) {
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
