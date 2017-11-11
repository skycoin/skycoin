package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var addressTxnsBkt = []byte("address_txns")

// addressTxn buckets for storing address related transactions
// address as key, transaction id slice as value
type addressTxns struct{}

func newAddressTxns(db *bolt.DB) (*addressTxns, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(addressTxnsBkt)
		return err
	}); err != nil {
		return nil, err
	}

	return &addressTxns{}, nil
}

// Get returns the transaction hashes of given address
func (atx *addressTxns) Get(tx *bolt.Tx, address cipher.Address) ([]cipher.SHA256, error) {
	var txHashes []cipher.SHA256
	if err := dbutil.GetBucketObjectDecoded(tx, addressTxnsBkt, address.Bytes(), &txHashes); err != nil {
		switch err.(type) {
		case dbutil.ObjectNotExistErr:
			return nil, nil
		default:
			return nil, err
		}
	}

	return txHashes, nil
}

// Add adds a hash to an address's hash list
func (atx *addressTxns) Add(tx *bolt.Tx, addr cipher.Address, hash cipher.SHA256) error {
	hashes, err := atx.Get(tx, addr)
	if err != nil {
		return err
	}

	// check for duplicates
	for _, u := range hashes {
		if u == hash {
			return nil
		}
	}

	hashes = append(hashes, hash)
	return dbutil.PutBucketValue(tx, addressTxnsBkt, addr.Bytes(), encoder.Serialize(hashes))
}

// IsEmpty checks if address transactions bucket is empty
func (atx *addressTxns) IsEmpty(tx *bolt.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, addressTxnsBkt)
}

// Reset resets the bucket
func (atx *addressTxns) Reset(tx *bolt.Tx) error {
	return dbutil.Reset(tx, addressTxnsBkt)
}
