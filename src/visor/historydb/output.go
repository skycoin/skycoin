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
	SpentTxnID    cipher.SHA256 // id of tx which spent this output.
	SpentBlockSeq uint64        // block seq that spent the output.
}

// Hash returns outhash
func (o UxOut) Hash() cipher.SHA256 {
	return o.Out.Hash()
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

// uxOuts bucket stores outputs, UxOut hash as key and Output as value.
type uxOuts struct{}

// put sets out value
func (ux *uxOuts) put(tx *dbutil.Tx, out UxOut) error {
	hash := out.Hash()
	return dbutil.PutBucketValue(tx, UxOutsBkt, hash[:], encoder.Serialize(out))
}

// get gets UxOut of given id
func (ux *uxOuts) get(tx *dbutil.Tx, uxID cipher.SHA256) (*UxOut, error) {
	var out UxOut

	if ok, err := dbutil.GetBucketObjectDecoded(tx, UxOutsBkt, uxID[:], &out); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &out, nil
}

// getArray returns uxOuts for a set of uxids, will return error if any of the uxids do not exist
func (ux *uxOuts) getArray(tx *dbutil.Tx, uxIDs []cipher.SHA256) ([]UxOut, error) {
	var outs []UxOut
	for _, uxID := range uxIDs {
		out, err := ux.get(tx, uxID)
		if err != nil {
			return nil, err
		} else if out == nil {
			return nil, NewErrUxOutNotExist(uxID.Hex())
		}

		outs = append(outs, *out)
	}

	return outs, nil
}

// isEmpty checks if the uxout bucekt is empty
func (ux *uxOuts) isEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, UxOutsBkt)
}

// reset resets the bucket
func (ux *uxOuts) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, UxOutsBkt)
}
