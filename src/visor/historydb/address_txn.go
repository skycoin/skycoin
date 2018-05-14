package historydb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// AddressTxnsBkt maps addresses to transaction hashes
var AddressTxnsBkt = []byte("address_txns")

// addressTxn buckets for storing address related transactions
// address as key, transaction id slice as value
type addressTxns struct{}

// Get returns the transaction hashes of given address
func (atx *addressTxns) Get(tx *dbutil.Tx, address cipher.Address) ([]cipher.SHA256, error) {
	var txHashes []cipher.SHA256
	if ok, err := dbutil.GetBucketObjectDecoded(tx, AddressTxnsBkt, address.Bytes(), &txHashes); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return txHashes, nil
}

// Add adds a hash to an address's hash list
func (atx *addressTxns) Add(tx *dbutil.Tx, addr cipher.Address, hash cipher.SHA256) error {
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
	return dbutil.PutBucketValue(tx, AddressTxnsBkt, addr.Bytes(), encoder.Serialize(hashes))
}

// IsEmpty checks if address transactions bucket is empty
func (atx *addressTxns) IsEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, AddressTxnsBkt)
}

// Reset resets the bucket
func (atx *addressTxns) Reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, AddressTxnsBkt)
}
