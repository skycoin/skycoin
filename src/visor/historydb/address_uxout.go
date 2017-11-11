package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var addressUxBkt = []byte("address_in")

// bucket for storing address with UxOut, key as address, value as UxOut.
type addressUx struct{}

// create address affected UxOuts bucket.
func newAddressUx(db *bolt.DB) (*addressUx, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(addressUxBkt)
		return err
	}); err != nil {
		return nil, err
	}

	return &addressUx{}, nil
}

// Get return nil on not found.
func (au *addressUx) Get(tx *bolt.Tx, address cipher.Address) ([]cipher.SHA256, error) {
	var uxHashes []cipher.SHA256

	if err := dbutil.GetBucketObjectDecoded(tx, addressUxBkt, address.Bytes(), &uxHashes); err != nil {
		switch err.(type) {
		case dbutil.ObjectNotExistErr:
			return nil, nil
		default:
			return nil, err
		}
	}

	return uxHashes, nil
}

// Add adds a hash to an address's hash list
func (au *addressUx) Add(tx *bolt.Tx, address cipher.Address, uxHash cipher.SHA256) error {
	hashes, err := au.Get(tx, address)
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
	return dbutil.PutBucketValue(tx, addressUxBkt, address.Bytes(), encoder.Serialize(hashes))
}

// IsEmpty checks if the addressUx bucket is empty
func (au *addressUx) IsEmpty(tx *bolt.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, addressUxBkt)
}

// Reset resets the bucket
func (au *addressUx) Reset(tx *bolt.Tx) error {
	return dbutil.Reset(tx, addressUxBkt)
}
