package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// Output augment output struct,
type Output struct {
	coin.TransactionOutput
	CreateTxID      cipher.SHA256 // id of tx which spent this output.
	CreatedBlockSeq uint64        // block seq that created the output.
	SpentTxID       cipher.SHA256
	SpentBlockSeq   uint64 // block seq that spent the output.
}

// Outputs bucket stores outputs, outID as key and Output as value.
type Outputs struct {
	bkt *bucket.Bucket
}

func newOutputs(db *bolt.DB) (*Outputs, error) {
	bkt, err := bucket.New([]byte("outputs"), db)
	if err != nil {
		return nil, err
	}
	return &Outputs{bkt}, nil
}

func (op *Outputs) Set(out Output) error {
	key := out.UxId(out.CreateTxID)
	bin := encoder.Serialize(out)
	return op.bkt.Put(key[:], bin)
}

func (op *Outputs) Get(uxID cipher.SHA256) (*Output, error) {
	bin := op.bkt.Get(uxID[:])
	if bin == nil {
		return nil, nil
	}

	out := Output{}
	if err := encoder.DeserializeRaw(bin, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
