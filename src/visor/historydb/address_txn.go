package historydb

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

//go:generate skyencoder -unexported -struct hashesWrapper

// hashesWrapper wraps []cipher.SHA256
type hashesWrapper struct {
	Hashes []cipher.SHA256
}

// AddressTxnsBkt maps addresses to transaction hashes
var AddressTxnsBkt = []byte("address_txns")

// addressTxn buckets for storing address related transactions
// address as key, transaction id slice as value
type addressTxns struct{}

// get returns the transaction hashes of given address
func (atx *addressTxns) get(tx *dbutil.Tx, addr cipher.Address) ([]cipher.SHA256, error) {
	var txnHashes hashesWrapper

	v, err := dbutil.GetBucketValueNoCopy(tx, AddressTxnsBkt, addr.Bytes())
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	if err := decodeHashesWrapperExact(v, &txnHashes); err != nil {
		return nil, err
	}

	return txnHashes.Hashes, nil
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

	buf, err := encodeHashesWrapper(&hashesWrapper{
		Hashes: hashes,
	})
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, AddressTxnsBkt, addr.Bytes(), buf)
}

// contains returns true if an address has transactions
func (atx *addressTxns) contains(tx *dbutil.Tx, addr cipher.Address) (bool, error) {
	return dbutil.BucketHasKey(tx, AddressTxnsBkt, addr.Bytes())
}

// isEmpty checks if address transactions bucket is empty
func (atx *addressTxns) isEmpty(tx *dbutil.Tx) (bool, error) {
	return dbutil.IsEmpty(tx, AddressTxnsBkt)
}

// reset resets the bucket
func (atx *addressTxns) reset(tx *dbutil.Tx) error {
	return dbutil.Reset(tx, AddressTxnsBkt)
}
