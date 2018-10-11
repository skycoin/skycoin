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

// get returns the transaction hashes of given address
func (atx *addressTxns) get(tx *dbutil.Tx, address cipher.Address) ([]cipher.SHA256, error) {
	var txHashes []cipher.SHA256
	if ok, err := dbutil.GetBucketObjectDecoded(tx, AddressTxnsBkt, address.Bytes(), &txHashes); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return txHashes, nil
}

// add adds a hash to an address's hash list
func (atx *addressTxns) add(tx *dbutil.Tx, addr cipher.Address, hash cipher.SHA256) error {
	hashes, err := atx.get(tx, addr)
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

// isEmpty checks if address transactions bucket is empty
func (atx *addressTxns) isEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, AddressTxnsBkt)
}

// reset resets the bucket
func (atx *addressTxns) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, AddressTxnsBkt)
}
