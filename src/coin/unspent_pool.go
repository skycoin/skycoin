package coin

import (
	"errors"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
)

// UnspentPool manages UxOuts
type UnspentPool struct {
	// Maps from UxOut.Hash() to UxOut
	Pool map[cipher.SHA256]UxOut
	// Total running hash
	XorHash cipher.SHA256
}

// NewUnspentPool creates unspent pool
func NewUnspentPool() UnspentPool {
	return UnspentPool{
		Pool:    make(map[cipher.SHA256]UxOut),
		XorHash: cipher.SHA256{},
	}
}

// Rebuild reconstructs the indices from the underlying Array
func (up *UnspentPool) Rebuild(uxs UxArray) {
	up.Pool = make(map[cipher.SHA256]UxOut, len(uxs))
	xh := cipher.SHA256{}
	for i := range uxs {
		h := uxs[i].Hash()
		up.Pool[h] = uxs[i]
		xh = xh.Xor(uxs[i].SnapshotHash())
	}
	up.XorHash = xh
	if len(up.Pool) != len(uxs) {
		log.Panic("Corrupt UnspentPool array: contains duplicate UxOut")
	}
}

// Array returns Pool as an array. Note: they are not in any particular order.
func (up *UnspentPool) Array() UxArray {
	arr := make(UxArray, len(up.Pool))
	i := 0
	for _, v := range up.Pool {
		arr[i] = v
		i++
	}
	return arr
}

// Add adds a UxOut to the UnspentPool
func (up *UnspentPool) Add(ux UxOut) {
	h := ux.Hash()
	if up.Has(h) {
		log.Panic("Attempt to insert UxOut twice")
	}
	up.Pool[h] = ux
	up.XorHash = up.XorHash.Xor(ux.SnapshotHash())
}

// Get returns a UxOut by hash, and whether it actually exists (if it does not
// exist, the map would return an empty UxOut)
func (up *UnspentPool) Get(h cipher.SHA256) (UxOut, bool) {
	ux, ok := up.Pool[h]
	return ux, ok
}

// GetMultiple returns a UxArray for hashes, or error if any not found
func (up *UnspentPool) GetMultiple(hashes []cipher.SHA256) (UxArray, error) {
	uxia := make(UxArray, len(hashes))
	for i := range hashes {
		uxi, exists := up.Get(hashes[i])
		if !exists {
			return nil, errors.New("Unspent output does not exist")
		}
		uxia[i] = uxi
	}
	return uxia, nil
}

// Collides checks for hash collisions with existing hashes
func (up *UnspentPool) Collides(hashes []cipher.SHA256) bool {
	for i := range hashes {
		if _, ok := up.Pool[hashes[i]]; ok {
			return true
		}
	}
	return false
}

// Has returns true if an unspent exists for this hash
func (up *UnspentPool) Has(h cipher.SHA256) bool {
	_, ok := up.Pool[h]
	return ok
}

// Del removes an unspent from the pool, by hash
func (up *UnspentPool) Del(h cipher.SHA256) {
	if ux, ok := up.Pool[h]; ok {
		up.XorHash = up.XorHash.Xor(ux.SnapshotHash())
		delete(up.Pool, h)
	}
}

// DelMultiple delete multiple hashes in a batch
func (up *UnspentPool) DelMultiple(hashes []cipher.SHA256) {
	for i := range hashes {
		up.Del(hashes[i])
	}
}

// AllForAddress returns all Unspents for a single address
// Warning: Not threadsafe.
// Deprecate: User application should not be querying this
// ^^^ Can't do a Spend without this info
func (up UnspentPool) AllForAddress(a cipher.Address) UxArray {
	uxo := make(UxArray, 0)
	for _, ux := range up.Pool {
		if ux.Body.Address == a {
			uxo = append(uxo, ux)
		}
	}
	return uxo
}

// AllForAddresses Returns Unspents for multiple addresses
// Warning: Not threadsafe.
// Deprecate: User application should not be querying this
// ^^^ Can't do a Spend without this info
func (up UnspentPool) AllForAddresses(addrs []cipher.Address) AddressUxOuts {
	m := make(map[cipher.Address]byte, len(addrs))
	for _, a := range addrs {
		m[a] = byte(1)
	}
	uxo := make(AddressUxOuts)
	for _, ux := range up.Pool {
		if _, exists := m[ux.Body.Address]; exists {
			uxo[ux.Body.Address] = append(uxo[ux.Body.Address], ux)
		}
	}
	return uxo
}

// GetUxHash returns unspent output checksum for the Block. Must be called after Block
// is fully initialized, and before its outputs are added to the unspent pool
func (up UnspentPool) GetUxHash() cipher.SHA256 {
	return up.XorHash
}

// Clone returns the copy of self
func (up *UnspentPool) Clone() UnspentPool {
	upClone := UnspentPool{
		Pool:    make(map[cipher.SHA256]UxOut),
		XorHash: up.XorHash,
	}

	for k, v := range up.Pool {
		upClone.Pool[k] = v
	}
	return upClone
}
