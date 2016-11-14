package coin

import (
	"errors"
	"log"

	"github.com/skycoin/skycoin/src/cipher"
)

// Manages UxOuts
type UnspentPool struct {
	// Maps from UxOut.Hash() to UxOut
	Pool map[cipher.SHA256]UxOut
	// Total running hash
	XorHash cipher.SHA256
}

func NewUnspentPool() UnspentPool {
	return UnspentPool{
		Pool:    make(map[cipher.SHA256]UxOut),
		XorHash: cipher.SHA256{},
	}
}

// Reconstructs the indices from the underlying Array
func (self *UnspentPool) Rebuild(uxs UxArray) {
	self.Pool = make(map[cipher.SHA256]UxOut, len(uxs))
	xh := cipher.SHA256{}
	for i, _ := range uxs {
		h := uxs[i].Hash()
		self.Pool[h] = uxs[i]
		xh = xh.Xor(uxs[i].SnapshotHash())
	}
	self.XorHash = xh
	if len(self.Pool) != len(uxs) {
		log.Panic("Corrupt UnspentPool array: contains duplicate UxOut")
	}
}

// Returns Pool as an array. Note: they are not in any particular order.
func (self *UnspentPool) Array() UxArray {
	arr := make(UxArray, len(self.Pool))
	i := 0
	for _, v := range self.Pool {
		arr[i] = v
		i++
	}
	return arr
}

// Adds a UxOut to the UnspentPool
func (self *UnspentPool) Add(ux UxOut) {
	h := ux.Hash()
	if self.Has(h) {
		log.Panic("Attempt to insert UxOut twice")
	}
	self.Pool[h] = ux
	self.XorHash = self.XorHash.Xor(ux.SnapshotHash())
}

// Returns a UxOut by hash, and whether it actually exists (if it does not
// exist, the map would return an empty UxOut)
func (self *UnspentPool) Get(h cipher.SHA256) (UxOut, bool) {
	ux, ok := self.Pool[h]
	return ux, ok
}

// Returns a UxArray for hashes, or error if any not found
func (self *UnspentPool) GetMultiple(hashes []cipher.SHA256) (UxArray, error) {
	uxia := make(UxArray, len(hashes))
	for i, _ := range hashes {
		uxi, exists := self.Get(hashes[i])
		if !exists {
			return nil, errors.New("Unspent output does not exist")
		}
		uxia[i] = uxi
	}
	return uxia, nil
}

// Checks for hash collisions with existing hashes
func (self *UnspentPool) Collides(hashes []cipher.SHA256) bool {
	for i, _ := range hashes {
		if _, ok := self.Pool[hashes[i]]; ok {
			return true
		}
	}
	return false
}

// Returns true if an unspent exists for this hash
func (self *UnspentPool) Has(h cipher.SHA256) bool {
	_, ok := self.Pool[h]
	return ok
}

// Removes an unspent from the pool, by hash
func (self *UnspentPool) Del(h cipher.SHA256) {
	if ux, ok := self.Pool[h]; ok {
		self.XorHash = self.XorHash.Xor(ux.SnapshotHash())
		delete(self.Pool, h)
	}
}

// Delete multiple hashes in a batch
func (self *UnspentPool) DelMultiple(hashes []cipher.SHA256) {
	for i := range hashes {
		self.Del(hashes[i])
	}
}

// Returns all Unspents for a single address
// Warning: Not threadsafe.
// Deprecate: User application should not be querying this
// ^^^ Can't do a Spend without this info
func (self UnspentPool) AllForAddress(a cipher.Address) UxArray {
	uxo := make(UxArray, 0)
	for _, ux := range self.Pool {
		if ux.Body.Address == a {
			uxo = append(uxo, ux)
		}
	}
	return uxo
}

// Returns Unspents for multiple addresses
// Warning: Not threadsafe.
// Deprecate: User application should not be querying this
// ^^^ Can't do a Spend without this info
func (self UnspentPool) AllForAddresses(addrs []cipher.Address) AddressUxOuts {
	m := make(map[cipher.Address]byte, len(addrs))
	for _, a := range addrs {
		m[a] = byte(1)
	}
	uxo := make(AddressUxOuts)
	for _, ux := range self.Pool {
		if _, exists := m[ux.Body.Address]; exists {
			uxo[ux.Body.Address] = append(uxo[ux.Body.Address], ux)
		}
	}
	return uxo
}

// GetUxHash returns unspent output checksum for the Block. Must be called after Block
// is fully initialized, and before its outputs are added to the unspent pool
func (self UnspentPool) GetUxHash() cipher.SHA256 {
	return self.XorHash
}
