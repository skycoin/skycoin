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
	meta *bucket.Bucket
}

type unspentMeta struct {
	*bolt.Bucket
}

func (m unspentMeta) getXorHash() (cipher.SHA256, error) {
	if v := m.Get(xorhashKey); v != nil {
		var hash cipher.SHA256
		copy(hash[:], v[:])
		return hash, nil
	}

	return cipher.SHA256{}, nil
}

func (m *unspentMeta) setXorHash(hash cipher.SHA256) error {
	return m.Put(xorhashKey, hash[:])
}

type uxOut struct {
	*bolt.Bucket
}

func (uo uxOut) get(hash cipher.SHA256) (*coin.UxOut, bool, error) {
	if v := uo.Get(hash[:]); v != nil {
		var out coin.UxOut
		if err := encoder.DeserializeRaw(v, &out); err != nil {
			return nil, false, err
		}
		return &out, true, nil
	}
	return nil, false, nil
}

func (uo uxOut) set(hash cipher.SHA256, ux coin.UxOut) error {
	v := encoder.Serialize(ux)
	return uo.Put(hash[:], v)
}

func (uo *uxOut) delete(hash cipher.SHA256) error {
	return uo.Delete(hash[:])
}

// NewUnspentPool creates new unspent pool instance
func NewUnspentPool(db *bolt.DB) (*UnspentPool, error) {
	up := &UnspentPool{db: db}

	pool, err := bucket.New([]byte("unspent_pool"), db)
	if err != nil {
		return nil, err
	}
	up.pool = pool

	meta, err := bucket.New([]byte("unspent_meta"), db)
	if err != nil {
		return nil, err
	}
	up.meta = meta

	return up, nil
}

// resetWithTx resets the unspent pool buckets
func (up *UnspentPool) resetWithTx(tx *bolt.Tx) error {
	if err := tx.DeleteBucket(up.pool.Name); err != nil {
		return err
	}

	_, err := tx.CreateBucket(up.pool.Name)
	if err != nil {
		return err
	}

	if err := tx.DeleteBucket(up.meta.Name); err != nil {
		return err
	}

	_, err = tx.CreateBucket(up.meta.Name)
	return err
}

func (up *UnspentPool) addWithTx(tx *bolt.Tx, ux coin.UxOut) (err error) {
	// all updates will rollback if return is not nil
	// in case of unexpected panic, we must catch it and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unspent pool add uxout failed: %v", err)
		}
	}()

	// check if the uxout does exist in the pool
	h := ux.Hash()
	if up.Contains(h) {
		return fmt.Errorf("attemps to insert uxout:%v twice into the unspent pool", h.Hex())
	}

	meta := unspentMeta{tx.Bucket(up.meta.Name)}
	xorhash, err := meta.getXorHash()
	if err != nil {
		return err
	}

	xorhash = xorhash.Xor(ux.SnapshotHash())
	if err := meta.setXorHash(xorhash); err != nil {
		return err
	}

	return uxOut{tx.Bucket(up.pool.Name)}.set(h, ux)
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
	if v := up.pool.Get(h[:]); v != nil {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return coin.UxOut{}, false, fmt.Errorf("get uxout from pool failed: %v", err)
		}
		return ux, true, nil
	}

	return coin.UxOut{}, false, nil
}

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *UnspentPool) GetAll() (coin.UxArray, error) {
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

// delete delete unspent of given hashes
func (up *UnspentPool) deleteWithTx(tx *bolt.Tx, hashes []cipher.SHA256) error {
	uxouts := uxOut{tx.Bucket(up.pool.Name)}
	meta := unspentMeta{tx.Bucket(up.meta.Name)}

	for _, hash := range hashes {
		ux, ok, err := uxouts.get(hash)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		uxHash, err := meta.getXorHash()
		if err != nil {
			return err
		}

		uxHash = uxHash.Xor(ux.SnapshotHash())

		// update uxhash
		if err = meta.setXorHash(uxHash); err != nil {
			return err
		}

		if err := uxouts.delete(hash); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the unspent outputs num
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
	v := up.pool.Get(h[:])
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
		return coin.UxArray{}, fmt.Errorf("get unspents of address %v failed: %v",
			addr.String(), err)
	}

	return uxs, nil
}

// GetUnspentsOfAddrs returns unspent outputs map of given addresses,
// the address as return map key, unspent outputs as value.
func (up *UnspentPool) GetUnspentsOfAddrs(addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
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
			fmt.Errorf("get unspents of address array failed: %v", err)
	}

	return addrUxs, nil
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *UnspentPool) GetUxHash() (cipher.SHA256, error) {
	if v := up.meta.Get(xorhashKey); v != nil {
		var hash cipher.SHA256
		copy(hash[:], v[:])
		return hash, nil
	}

	return cipher.SHA256{}, nil
}
