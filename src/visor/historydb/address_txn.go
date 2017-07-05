package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var addressTxnsBktName = []byte("address_txns")

// addressTxn buckets for storing address related transactions
// address as key, transaction id slice as value
type addressTxns struct {
	bkt *bucket.Bucket
}

func newAddressTxnsBkt(db *bolt.DB) (*addressTxns, error) {
	bkt, err := bucket.New(addressTxnsBktName, db)
	if err != nil {
		return nil, err
	}

	return &addressTxns{bkt}, nil
}

// Get returns the transaction hashes of given address
func (atx *addressTxns) Get(address cipher.Address) ([]cipher.SHA256, error) {
	var txHashes []cipher.SHA256
	v := atx.bkt.Get(address.Bytes())
	if v == nil {
		return []cipher.SHA256{}, nil
	}

	if err := encoder.DeserializeRaw(v, &txHashes); err != nil {
		return []cipher.SHA256{}, err
	}

	return txHashes, nil
}

// IsEmpty checks if address transactions bucket is empty
func (atx *addressTxns) IsEmpty() bool {
	return atx.bkt.IsEmpty()
}

// Reset resets the bucket
func (atx *addressTxns) Reset() error {
	return atx.bkt.Reset()
}

func setAddressTxns(bkt *bolt.Bucket, addr cipher.Address, hash cipher.SHA256) error {
	// get hashes
	addrBytes := addr.Bytes()
	v := bkt.Get(addrBytes)
	if v == nil {
		bin := encoder.Serialize([]cipher.SHA256{hash})
		return bkt.Put(addrBytes, bin)
	}

	var hashes []cipher.SHA256
	if err := encoder.DeserializeRaw(v, &hashes); err != nil {
		return err
	}

	// check dup
	for _, u := range hashes {
		if u == hash {
			return nil
		}
	}

	hashes = append(hashes, hash)
	bin := encoder.Serialize(hashes)
	return bkt.Put(addrBytes, bin)
}
