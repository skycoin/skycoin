package historydb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// AddressUxBkt maps addresses to unspent outputs
var AddressUxBkt = []byte("address_in")

// bucket for storing address with UxOut, key as address, value as UxOut.
type addressUx struct{}

// get return nil on not found.
func (au *addressUx) get(tx *dbutil.Tx, address cipher.Address) ([]cipher.SHA256, error) {
	var uxHashes []cipher.SHA256

	if ok, err := dbutil.GetBucketObjectDecoded(tx, AddressUxBkt, address.Bytes(), &uxHashes); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return uxHashes, nil
}

// add adds a hash to an address's hash list
func (au *addressUx) add(tx *dbutil.Tx, address cipher.Address, uxHash cipher.SHA256) error {
	hashes, err := au.get(tx, address)
	if err != nil {
		return err
	}

	// check for duplicate hashes
	for _, u := range hashes {
		if u == uxHash {
			return nil
		}
	}

	hashes = append(hashes, uxHash)
	return dbutil.PutBucketValue(tx, AddressUxBkt, address.Bytes(), encoder.Serialize(hashes))
}

// isEmpty checks if the addressUx bucket is empty
func (au *addressUx) isEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, AddressUxBkt)
}

// reset resets the bucket
func (au *addressUx) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, AddressUxBkt)
}
