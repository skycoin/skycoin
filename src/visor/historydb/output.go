package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// UxOut expend coin.UxOut struct
type UxOut struct {
	Out           coin.UxOut
	SpentTxID     cipher.SHA256 // id of tx which spent this output.
	SpentBlockSeq uint64        // block seq that spent the output.
}

// UxOutJSON UxOut's json format
type UxOutJSON struct {
	Uxid          string `json:"uxid"`
	Time          uint64 `json:"time"`
	SrcBkSeq      uint64 `json:"src_block_seq"`
	SrcTx         string `json:"src_tx"`
	OwnerAddress  string `json:"owner_address"`
	Coins         uint64 `json:"coins"`
	Hours         uint64 `json:"hours"`
	SpentBlockSeq uint64 `json:"spent_block_seq"` // block seq that spent the output.
	SpentTxID     string `json:"spent_tx"`        // id of tx which spent this output.
}

// NewUxOutJSON generates UxOutJSON from UxOut
func NewUxOutJSON(out *UxOut) *UxOutJSON {
	if out == nil {
		return nil
	}

	return &UxOutJSON{
		Uxid:          out.Hash().Hex(),
		Time:          out.Out.Head.Time,
		SrcBkSeq:      out.Out.Head.BkSeq,
		SrcTx:         out.Out.Body.SrcTransaction.Hex(),
		OwnerAddress:  out.Out.Body.Address.String(),
		Coins:         out.Out.Body.Coins,
		Hours:         out.Out.Body.Hours,
		SpentBlockSeq: out.SpentBlockSeq,
		SpentTxID:     out.SpentTxID.Hex(),
	}
}

// Hash returns outhash
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

// Set sets out value
func (ux *UxOuts) Set(out UxOut) error {
	key := out.Hash()
	bin := encoder.Serialize(out)
	return ux.bkt.Put(key[:], bin)
}

// Get gets UxOut of given id
func (ux *UxOuts) Get(uxID cipher.SHA256) (*UxOut, error) {
	bin := ux.bkt.Get(uxID[:])
	if bin == nil {
		return nil, nil
	}

	out := UxOut{}
	if err := encoder.DeserializeRaw(bin, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// IsEmpty checks if the uxout bucekt is empty
func (ux *UxOuts) IsEmpty() bool {
	return ux.bkt.IsEmpty()
}

// Reset resets the bucket
func (ux *UxOuts) Reset() error {
	return ux.bkt.Reset()
}

func getOutput(bkt *bolt.Bucket, hash cipher.SHA256) (*UxOut, error) {
	bin := bkt.Get(hash[:])
	if bin != nil {
		var out UxOut
		if err := encoder.DeserializeRaw(bin, &out); err != nil {
			return nil, err
		}
		return &out, nil
	}

	return nil, nil
}

func setOutput(bkt *bolt.Bucket, ux UxOut) error {
	hash := ux.Hash()
	return bkt.Put(hash[:], encoder.Serialize(ux))
}
