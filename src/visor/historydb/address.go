package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// address provides apis for getting all spent outputs from specific address, and
// getting all inputs of specific address.

// addressInput bucket stores outputs that send to the addresses,
// address as key, out id as value.
// type addressIn struct {
// 	bkt *bucket.Bucket
// }

// addressOut bucket stores output that address spent.
// type addressOut struct {
// 	bkt *bucket.Bucket
// }

type addressUx struct {
	bkt *bucket.Bucket
}

func newAddressUx(db *bolt.DB, name []byte) (*addressUx, error) {
	bkt, err := bucket.New(name, db)
	if err != nil {
		return nil, err
	}

	return &addressUx{bkt}, nil
}

func newAddressInBkt(db *bolt.DB) (*addressUx, error) {
	return newAddressUx(db, []byte("address_in"))
}

func newAddressOutBkt(db *bolt.DB) (*addressUx, error) {
	return newAddressUx(db, []byte("address_out"))
}

// func newAddressIn(db *bolt.DB) (*addressIn, error) {
// 	bkt, err := bucket.New([]byte("address_in"), db)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &addressIn{bkt}, nil
// }

// Get return nil on not found.
func (au *addressUx) Get(address cipher.Address) ([]cipher.SHA256, error) {
	uxHashes := []cipher.SHA256{}
	bin := au.bkt.Get(address.Bytes())
	if bin == nil {
		return nil, nil
	}
	if err := encoder.DeserializeRaw(bin, &uxHashes); err != nil {
		return nil, err
	}
	return uxHashes, nil
}

func (au *addressUx) Add(address cipher.Address, uxHash cipher.SHA256) error {
	hashes, err := au.Get(address)
	if err != nil {
		return err
	}

	if hashes == nil {
		bin := encoder.Serialize([]cipher.SHA256{uxHash})
		return au.bkt.Put(address.Bytes(), bin)
	}

	// check dup
	for _, u := range hashes {
		if u == uxHash {
			return nil
		}
	}

	hashes = append(hashes, uxHash)
	bin := encoder.Serialize(hashes)
	return au.bkt.Put(address.Bytes(), bin)
}

// func newAddressOut(db *bolt.DB) (*addressOut, error) {
// 	bkt, err := bucket.New([]byte("address_out"), db)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &addressOut{bkt}, nil
// }

// func (ao *addressOut) Add(address cipher.Address, uxHash cipher.SHA256) error {
// 	return ao.bkt.Put(address.Bytes(), uxHash[:])
// }
