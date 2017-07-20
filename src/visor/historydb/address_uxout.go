package historydb

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

// bucket for storing address with UxOut, key as address, value as UxOut.
type addressUx struct {
	bkt *bucket.Bucket
}

// create address affected UxOuts bucket.
func newAddressUxBkt(db *bolt.DB) (*addressUx, error) {
	bkt, err := bucket.New([]byte("address_in"), db)
	if err != nil {
		return nil, err
	}

	return &addressUx{bkt}, nil
}

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

// IsEmpty checks if the addressUx bucket is empty
func (au *addressUx) IsEmpty() bool {
	return au.bkt.IsEmpty()
}

// Reset resets the bucket
func (au *addressUx) Reset() error {
	return au.bkt.Reset()
}

func setAddressUx(bkt *bolt.Bucket, addr cipher.Address, uxHash cipher.SHA256) error {
	bin := bkt.Get(addr.Bytes())
	if bin == nil {
		return bkt.Put(addr.Bytes(), encoder.Serialize([]cipher.SHA256{uxHash}))
	}

	uxHashes := []cipher.SHA256{}
	if err := encoder.DeserializeRaw(bin, &uxHashes); err != nil {
		return err
	}

	// check dup
	for _, u := range uxHashes {
		if u == uxHash {
			return nil
		}
	}

	uxHashes = append(uxHashes, uxHash)
	return bkt.Put(addr.Bytes(), encoder.Serialize(uxHashes))
}
