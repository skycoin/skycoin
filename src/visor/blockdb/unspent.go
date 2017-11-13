package blockdb

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	xorhashKey = []byte("xorhash")

	// bucket for unspent pool
	unspentPoolBkt = []byte("unspent_pool")
	// bucket for unspent meta info
	unspentMetaBkt = []byte("unspent_meta")
)

// UnspentGetter provides unspend pool related
// querying methods
type UnspentGetter interface {
	// GetUnspentsOfAddrs returns all unspent outputs of given addresses
	GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts
	Get(cipher.SHA256) (coin.UxOut, bool)
}

type unspentMeta struct{}

func (m unspentMeta) getXorHash(tx *bolt.Tx) (cipher.SHA256, error) {
	v, err := dbutil.GetBucketValue(tx, unspentMetaBkt, xorhashKey)
	if err != nil {
		return cipher.SHA256{}, err
	} else if v == nil {
		return cipher.SHA256{}, nil
	}

	var hash cipher.SHA256
	copy(hash[:], v[:])
	return hash, nil
}

func (m *unspentMeta) setXorHash(tx *bolt.Tx, hash cipher.SHA256) error {
	return dbutil.PutBucketValue(tx, unspentMetaBkt, xorhashKey, hash[:])
}

type pool struct{}

func (pl pool) get(tx *bolt.Tx, hash cipher.SHA256) (*coin.UxOut, bool, error) {
	var out coin.UxOut

	if ok, err := dbutil.GetBucketObjectDecoded(tx, unspentPoolBkt, hash[:], &out); err != nil {
		return nil, false, err
	} else if !ok {
		return nil, false, nil
	}

	return &out, true, nil
}

func (pl pool) set(tx *bolt.Tx, hash cipher.SHA256, ux coin.UxOut) error {
	return dbutil.PutBucketValue(tx, unspentPoolBkt, hash[:], encoder.Serialize(ux))
}

func (pl *pool) delete(tx *bolt.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, unspentPoolBkt, hash[:])
}

// Unspents unspent outputs pool
type Unspents struct {
	db    *dbutil.DB
	pool  *pool
	meta  *unspentMeta
	cache struct {
		pool   map[string]coin.UxOut
		uxhash cipher.SHA256
	}
	sync.RWMutex
}

// NewUnspentPool creates new unspent pool instance
func NewUnspentPool(db *dbutil.DB) (*Unspents, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		return dbutil.CreateBuckets(tx, [][]byte{
			unspentPoolBkt,
			unspentMetaBkt,
		})
	}); err != nil {
		return nil, err
	}

	up := &Unspents{
		db:   db,
		pool: &pool{},
		meta: &unspentMeta{},
	}
	up.cache.pool = make(map[string]coin.UxOut)

	// Load from db
	if err := db.View(func(tx *bolt.Tx) error {
		return up.syncCache(tx)
	}); err != nil {
		return nil, err
	}

	return up, nil
}

func (up *Unspents) syncCache(tx *bolt.Tx) error {
	// Load unspent outputs
	if err := dbutil.ForEach(tx, unspentPoolBkt, func(k, v []byte) error {
		var hash cipher.SHA256
		copy(hash[:], k[:])

		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return fmt.Errorf("load unspent outputs from db failed: %v", err)
		}

		up.cache.pool[hash.Hex()] = ux
		return nil
	}); err != nil {
		return err
	}

	// Load uxhash
	uxhash, err := up.meta.getXorHash(tx)
	if err != nil {
		return err
	}

	up.cache.uxhash = uxhash
	return nil
}

// ProcessBlock adds unspents from a block to the unspent pool
func (up *Unspents) ProcessBlock(b *coin.SignedBlock) dbutil.TxHandler {
	return func(tx *bolt.Tx) (dbutil.Rollback, error) {
		var (
			delUxs    []coin.UxOut
			addUxs    []coin.UxOut
			uxHash    cipher.SHA256
			oldUxHash = up.cache.uxhash
		)

		for _, txn := range b.Body.Transactions {
			// Get uxouts that need to be deleted
			uxs, err := up.getArray(txn.In)
			if err != nil {
				return func() {}, err
			}

			delUxs = append(delUxs, uxs...)

			// Remove spent outputs
			if err := up.delete(tx, txn.In); err != nil {
				return func() {}, err
			}

			// Create new outputs
			txUxs := coin.CreateUnspents(b.Head, txn)
			addUxs = append(addUxs, txUxs...)
			for i := range txUxs {
				uxHash, err = up.add(tx, txUxs[i])
				if err != nil {
					return func() {}, err
				}
			}
		}

		// update caches
		up.Lock()
		defer up.Unlock()
		up.deleteUxFromCache(delUxs)
		up.addUxToCache(addUxs)
		up.updateUxHashInCache(uxHash)

		return func() {
			up.Lock()
			defer up.Unlock()
			// reverse the cache
			up.deleteUxFromCache(addUxs)
			up.addUxToCache(delUxs)
			up.updateUxHashInCache(oldUxHash)
		}, nil
	}
}

func (up *Unspents) add(tx *bolt.Tx, ux coin.UxOut) (uxhash cipher.SHA256, err error) {
	// will rollback all updates if return is not nil
	// in case of unexpected panic, we must catch it and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unspent pool add uxout failed: %v", err)
		}
	}()

	// check if the uxout does exist in the pool
	h := ux.Hash()
	if up.Contains(h) {
		return cipher.SHA256{}, fmt.Errorf("attemps to insert uxout:%v twice into the unspent pool", h.Hex())
	}

	xorhash, err := up.meta.getXorHash(tx)
	if err != nil {
		return cipher.SHA256{}, err
	}

	xorhash = xorhash.Xor(ux.SnapshotHash())
	if err := up.meta.setXorHash(tx, xorhash); err != nil {
		return cipher.SHA256{}, err
	}

	err = up.pool.set(tx, h, ux)
	if err != nil {
		return cipher.SHA256{}, err
	}

	return xorhash, nil
}

func (up *Unspents) deleteUxFromCache(uxs []coin.UxOut) {
	for _, ux := range uxs {
		delete(up.cache.pool, ux.Hash().Hex())
	}
}

func (up *Unspents) addUxToCache(uxs []coin.UxOut) {
	for i, ux := range uxs {
		up.cache.pool[ux.Hash().Hex()] = uxs[i]
	}
}

func (up *Unspents) updateUxHashInCache(hash cipher.SHA256) {
	up.cache.uxhash = hash
}

// GetArray returns UxOut for a set of hashes, will return error if any of the hashes do not exist in the pool.
func (up *Unspents) GetArray(hashes []cipher.SHA256) (coin.UxArray, error) {
	up.RLock()
	defer up.RUnlock()
	return up.getArray(hashes)
}

func (up *Unspents) getArray(hashes []cipher.SHA256) (coin.UxArray, error) {
	uxs := make(coin.UxArray, 0, len(hashes))
	for i := range hashes {
		ux, ok := up.cache.pool[hashes[i].Hex()]
		if !ok {
			return nil, fmt.Errorf("unspent output of %s does not exist", hashes[i].Hex())
		}

		uxs = append(uxs, ux)
	}
	return uxs, nil
}

// Get returns the uxout value of given hash
func (up *Unspents) Get(h cipher.SHA256) (coin.UxOut, bool) {
	up.RLock()
	defer up.RUnlock()

	ux, ok := up.cache.pool[h.Hex()]
	return ux, ok
}

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *Unspents) GetAll() (coin.UxArray, error) {
	up.RLock()
	defer up.RUnlock()

	arr := make(coin.UxArray, 0, len(up.cache.pool))
	for _, ux := range up.cache.pool {
		arr = append(arr, ux)
	}

	return arr, nil
}

// delete delete unspent of given hashes
func (up *Unspents) delete(tx *bolt.Tx, hashes []cipher.SHA256) error {
	var uxHash cipher.SHA256
	for _, hash := range hashes {
		ux, ok, err := up.pool.get(tx, hash)
		if err != nil {
			return err
		}

		if !ok {
			continue
		}

		uxHash, err = up.meta.getXorHash(tx)
		if err != nil {
			return err
		}

		uxHash = uxHash.Xor(ux.SnapshotHash())

		// update uxhash
		if err = up.meta.setXorHash(tx, uxHash); err != nil {
			return err
		}

		if err := up.pool.delete(tx, hash); err != nil {
			return err
		}
	}

	return nil
}

// Len returns the unspent outputs num
func (up *Unspents) Len() uint64 {
	up.RLock()
	defer up.RUnlock()
	return uint64(len(up.cache.pool))
}

// Contains check if the hash of uxout does exist in the pool
func (up *Unspents) Contains(h cipher.SHA256) bool {
	up.RLock()
	defer up.RUnlock()
	_, ok := up.cache.pool[h.Hex()]
	return ok
}

// GetUnspentsOfAddrs returns unspent outputs map of given addresses,
// the address as return map key, unspent outputs as value.
func (up *Unspents) GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts {
	up.RLock()
	defer up.RUnlock()

	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
	}

	addrUxs := coin.AddressUxOuts{}
	for _, ux := range up.cache.pool {
		if _, ok := addrm[ux.Body.Address]; ok {
			addrUxs[ux.Body.Address] = append(addrUxs[ux.Body.Address], ux)
		}
	}
	return addrUxs
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *Unspents) GetUxHash() cipher.SHA256 {
	up.RLock()
	defer up.RUnlock()
	return up.cache.uxhash
}
