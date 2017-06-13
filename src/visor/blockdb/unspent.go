package blockdb

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	xorhashKey = []byte("xorhash")
)

// UnspentPool unspent outputs pool
type UnspentPool struct {
	db   *bolt.DB
	pool *bucket.Bucket
	meta *bucket.Bucket // mainly for store xorhash and parsed height
}

// NewUnspentPool creates new unspent pool instance
func NewUnspentPool(db *bolt.DB) (*UnspentPool, error) {
	up := &UnspentPool{
		db: db,
	}

	pool, err := bucket.New([]byte("unspent_pool"), db)
	if err != nil {
		return nil, err
	}
	up.pool = pool

	meta, err := bucket.New([]byte("unspent_pool_meta"), db)
	if err != nil {
		return nil, err
	}

	// initialize xorhash with empty SHA256
	meta.Put(xorhashKey, []byte(cipher.SHA256{}.Hex()))

	up.meta = meta

	return up, nil
}

// Add adds a UxOut to pool
func (up *UnspentPool) Add(ux coin.UxOut) error {
	h := ux.Hash()
	if up.Contains(h) {
		return fmt.Errorf("attemps to insert uxout:%v twice into the unspent pool", h.Hex())
	}
	v := encoder.Serialize(ux)
	xorhash, err := up.getXorHash()
	if err != nil {
		return err
	}

	xorhash = xorhash.Xor(ux.SnapshotHash())
	if err := up.setXorHash(xorhash); err != nil {
		return err
	}

	return up.pool.Put([]byte(h.Hex()), v)
}

// GetArray returns UxOut by given hash array, will return error when
// if any of the hashes is not exist.
func (up *UnspentPool) GetArray(hashes []cipher.SHA256) (coin.UxArray, error) {
	uxs := make(coin.UxArray, 0, len(hashes))
	for i := range hashes {
		ux, ok, err := up.Get(hashes[i])
		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, fmt.Errorf("unspent output of %s does not exist", hashes[i].Hex())
		}

		uxs = append(uxs, ux)
	}
	return uxs, nil
}

// Get returns the uxout value of give hash
func (up *UnspentPool) Get(h cipher.SHA256) (coin.UxOut, bool, error) {
	if v := up.pool.Get([]byte(h.Hex())); v != nil {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return coin.UxOut{}, false, fmt.Errorf("get uxout from pool failed: %v", err)
		}
		return ux, true, nil
	}

	return coin.UxOut{}, false, nil
}

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *UnspentPool) GetAll() (ua coin.UxArray, err error) {
	arr := make(coin.UxArray, 0, up.pool.Len())
	if err := up.pool.ForEach(func(k, v []byte) error {
		var uxout coin.UxOut
		if err := encoder.DeserializeRaw(v, &uxout); err != nil {
			return err
		}

		arr = append(arr, uxout)
		return nil
	}); err != nil {
		return coin.UxArray{}, err
	}

	return arr, nil
}

// delete removes an unsepnt from the pool by hash
func (up *UnspentPool) delete(h cipher.SHA256) error {
	if v := up.pool.Get([]byte(h.Hex())); v != nil {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		hash, err := up.getXorHash()
		if err != nil {
			return err
		}

		hash = hash.Xor(ux.SnapshotHash())
		if err := up.setXorHash(hash); err != nil {
			return err
		}

		// delete unspent hash from pool
		if err := up.pool.Delete([]byte(h.Hex())); err != nil {
			return err
		}
	}

	return nil
}

// Delete delete unspent of given hashes
func (up *UnspentPool) Delete(hashes []cipher.SHA256) error {
	for _, hash := range hashes {
		if err := up.delete(hash); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the pool size
func (up *UnspentPool) Len() uint64 {
	return uint64(up.pool.Len())
}

// Collides checks for hash collisions with existing hashes
func (up *UnspentPool) Collides(hashes []cipher.SHA256) bool {
	for i := range hashes {
		if up.Contains(hashes[i]) {
			return true
		}
	}
	return false
}

// Contains check if the hash of uxout does exist in the pool
func (up *UnspentPool) Contains(h cipher.SHA256) bool {
	v := up.pool.Get([]byte(h.Hex()))
	return v != nil
}

// GetUnspentsOfAddr returns all unspent outputs of given address
func (up *UnspentPool) GetUnspentsOfAddr(addr cipher.Address) (coin.UxArray, error) {
	uxs := coin.UxArray{}
	if err := up.pool.ForEach(func(k, v []byte) error {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		if ux.Body.Address == addr {
			uxs = append(uxs, ux)
		}
		return nil
	}); err != nil {
		return coin.UxArray{}, fmt.Errorf("get unsepnts of address %v failed: %v",
			addr.String(), err)
	}

	return uxs, nil
}

// GetUnspentsOfAddrs returns unspent outputs map of given addresses,
// the address as return map key, unspent outputs as value.
func (up *UnspentPool) GetUnspentsOfAddrs(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]byte, len(addrs))
	for _, a := range addrs {
		addrm[a] = byte(1)
	}

	addrUxs := coin.AddressUxOuts{}
	if err := up.pool.ForEach(func(k, v []byte) error {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		if _, ok := addrm[ux.Body.Address]; ok {
			addrUxs[ux.Body.Address] = append(addrUxs[ux.Body.Address], ux)
		}
		return nil
	}); err != nil {
		return coin.AddressUxOuts{},
			fmt.Errorf("get unsepnts of address array failed: %v", err)
	}

	return addrUxs, nil
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *UnspentPool) GetUxHash() (cipher.SHA256, error) {
	return up.getXorHash()
}

// getXorHash returns XorHash
func (up *UnspentPool) getXorHash() (cipher.SHA256, error) {
	v := up.meta.Get(xorhashKey)
	hash, err := cipher.SHA256FromHex(string(v))
	if err != nil {
		return cipher.SHA256{}, fmt.Errorf("parse xorhash failed: %v", err)
	}

	return hash, nil
}

// setXorHash updates hash
func (up *UnspentPool) setXorHash(hash cipher.SHA256) error {
	if err := up.meta.Put(xorhashKey, []byte(hash.Hex())); err != nil {
		return fmt.Errorf("set xorhash failed: %v", err)
	}
	return nil
}
