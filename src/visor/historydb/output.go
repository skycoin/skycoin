package historydb

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// UxOutsBkt holds unspent outputs
var UxOutsBkt = []byte("uxouts")

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
type UxOuts struct{}

// Set sets out value
func (ux *UxOuts) Set(tx *dbutil.Tx, out UxOut) error {
	hash := out.Hash()
	return dbutil.PutBucketValue(tx, UxOutsBkt, hash[:], encoder.Serialize(out))
}

// Get gets UxOut of given id
func (ux *UxOuts) Get(tx *dbutil.Tx, uxID cipher.SHA256) (*UxOut, error) {
	var out UxOut

	if ok, err := dbutil.GetBucketObjectDecoded(tx, UxOutsBkt, uxID[:], &out); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &out, nil
}

// GetArray returns UxOuts for a set of uxids, will return error if any of the uxids do not exist
func (ux *UxOuts) GetArray(tx *dbutil.Tx, uxIDs []cipher.SHA256) ([]*UxOut, error) {
	var outs []*UxOut
	for _, uxID := range uxIDs {
		out, err := ux.Get(tx, uxID)
		if err != nil {
			return nil, err
		} else if out == nil {
			return nil, NewErrUxOutNotExist(uxID.Hex())
		}

		outs = append(outs, out)
	}

	return outs, nil
}

// IsEmpty checks if the uxout bucekt is empty
func (ux *UxOuts) IsEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, UxOutsBkt)
}

// Reset resets the bucket
func (ux *UxOuts) Reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, UxOutsBkt)
}

// ErrUxOutNotExist is returned if an uxout is not found in historydb
type ErrUxOutNotExist struct {
	UxID string
}

// NewErrUxOutNotExist creates ErrUxOutNotExist from a UxID
func NewErrUxOutNotExist(uxID string) error {
	return ErrUxOutNotExist{UxID: uxID}
}

func (e ErrUxOutNotExist) Error() string {
	return fmt.Sprintf("uxout of %s does not exist", e.UxID)
}
