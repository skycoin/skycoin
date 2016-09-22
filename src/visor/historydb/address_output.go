package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// address provides apis for getting all spent outputs from specific address, and
// getting all inputs of specific address.

// addressInput bucket stores outputs that send to the addresses,
// address as key, out id as value.
type addressIn struct {
	bkt *bucket.Bucket
}

// addressOut bucket stores output that address spent.
type addressOut struct {
	bkt *bucket.Bucket
}

func newAddressIn(db *bolt.DB) (*addressIn, error) {
	bkt, err := bucket.New([]byte("address_in"), db)
	if err != nil {
		return nil, err
	}

	return &addressIn{bkt}, nil
}

func (ai *addressIn) Add(address cipher.Address, outID cipher.SHA256) error {
	return ai.bkt.Put(address.Bytes(), outID[:])
}

func newAddressOut(db *bolt.DB) (*addressOut, error) {
	bkt, err := bucket.New([]byte("address_out"), db)
	if err != nil {
		return nil, err
	}
	return &addressOut{bkt}, nil
}

func (ao *addressOut) Add(address cipher.Address, outID cipher.SHA256) error {
	return ao.bkt.Put(address.Bytes(), outID[:])
}
