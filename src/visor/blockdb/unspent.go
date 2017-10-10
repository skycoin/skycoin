package blockdb

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/bucket"
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

// Unspents unspent outputs pool
type Unspents struct {
	db    *bolt.DB
	pool  *pool
	meta  *unspentMeta
	cache struct {
		pool   map[string]coin.UxOut
		uxhash cipher.SHA256
	}
	sync.Mutex
}

type unspentMeta struct {
	bucket.Bucket
}

func newUnspentMeta(db *bolt.DB) (*unspentMeta, error) {
	bkt, err := bucket.New(unspentMetaBkt, db)
	if err != nil {
		return nil, fmt.Errorf("failed to create unspent_meta bucket: %v", err)
	}

	return &unspentMeta{
		Bucket: *bkt,
	}, nil
}

func (m unspentMeta) getXorHashWithTx(tx *bolt.Tx) (cipher.SHA256, error) {
	if v := m.GetWithTx(tx, xorhashKey); v != nil {
		var hash cipher.SHA256
		copy(hash[:], v[:])
		return hash, nil
	}

	return cipher.SHA256{}, nil
}

func (m *unspentMeta) setXorHashWithTx(tx *bolt.Tx, hash cipher.SHA256) error {
	return m.PutWithTx(tx, xorhashKey, hash[:])
}

type pool struct {
	bucket.Bucket
}

func newPool(db *bolt.DB) (*pool, error) {
	bkt, err := bucket.New(unspentPoolBkt, db)
	if err != nil {
		return nil, err
	}

	return &pool{
		Bucket: *bkt,
	}, nil
}

func (pl pool) getWithTx(tx *bolt.Tx, hash cipher.SHA256) (*coin.UxOut, bool, error) {
	if v := pl.GetWithTx(tx, hash[:]); v != nil {
		var out coin.UxOut
		if err := encoder.DeserializeRaw(v, &out); err != nil {
			return nil, false, err
		}
		return &out, true, nil
	}
	return nil, false, nil
}

func (pl pool) setWithTx(tx *bolt.Tx, hash cipher.SHA256, ux coin.UxOut) error {
	v := encoder.Serialize(ux)
	return pl.PutWithTx(tx, hash[:], v)
}

func (pl *pool) deleteWithTx(tx *bolt.Tx, hash cipher.SHA256) error {
	return pl.DeleteWithTx(tx, hash[:])
}

// NewUnspentPool creates new unspent pool instance
func NewUnspentPool(db *bolt.DB) (*Unspents, error) {
	up := &Unspents{db: db}
	up.cache.pool = make(map[string]coin.UxOut)

	pool, err := newPool(db)
	if err != nil {
		return nil, err
	}
	up.pool = pool

	meta, err := newUnspentMeta(db)
	if err != nil {
		return nil, err
	}
	up.meta = meta

	// load from db
	if err := up.syncCache(); err != nil {
		return nil, err
	}

	return up, nil
}

func (up *Unspents) syncCache() error {
	// load unspent outputs
	if err := up.pool.ForEach(func(k, v []byte) error {
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

	// load uxhash
	uxhash, err := up.getUxHashFromDB()
	if err != nil {
		return err
	}

	up.cache.uxhash = uxhash
	return nil
}

func (up *Unspents) ProcessBlock(b *coin.SignedBlock) bucket.TxHandler {
	return func(tx *bolt.Tx) (bucket.Rollback, error) {
		var (
			delUxs    []coin.UxOut
			addUxs    []coin.UxOut
			uxHash    cipher.SHA256
			oldUxHash = up.cache.uxhash
		)

		for _, txn := range b.Body.Transactions {
			// get uxouts that need to be deleted
			uxs, err := up.getArray(txn.In)
			if err != nil {
				return func() {}, err
			}

			delUxs = append(delUxs, uxs...)

			// Remove spent outputs
			if _, err = up.deleteWithTx(tx, txn.In); err != nil {
				return func() {}, err
			}

			// Create new outputs
			txUxs := coin.CreateUnspents(b.Head, txn)
			addUxs = append(addUxs, txUxs...)
			for i := range txUxs {
				uxHash, err = up.addWithTx(tx, txUxs[i])
				if err != nil {
					return func() {}, err
				}
			}
		}

		// update caches
		up.Lock()
		up.deleteUxFromCache(delUxs)
		up.addUxToCache(addUxs)
		up.updateUxHashInCache(uxHash)
		up.Unlock()

		return func() {
			up.Lock()
			// reverse the cache
			up.deleteUxFromCache(addUxs)
			up.addUxToCache(delUxs)
			up.updateUxHashInCache(oldUxHash)
			up.Unlock()
		}, nil
	}
}

func (up *Unspents) addWithTx(tx *bolt.Tx, ux coin.UxOut) (uxhash cipher.SHA256, err error) {
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

	xorhash, err := up.meta.getXorHashWithTx(tx)
	if err != nil {
		return cipher.SHA256{}, err
	}

	xorhash = xorhash.Xor(ux.SnapshotHash())
	if err := up.meta.setXorHashWithTx(tx, xorhash); err != nil {
		return cipher.SHA256{}, err
	}

	err = up.pool.setWithTx(tx, h, ux)
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

// GetArray returns UxOut by given hash array, will return error when
// if any of the hashes is not exist.
func (up *Unspents) GetArray(hashes []cipher.SHA256) (coin.UxArray, error) {
	up.Lock()
	defer up.Unlock()
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
	up.Lock()
	ux, ok := up.cache.pool[h.Hex()]
	up.Unlock()

	return ux, ok
}

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *Unspents) GetAll() (coin.UxArray, error) {
	up.Lock()
	arr := make(coin.UxArray, 0, len(up.cache.pool))
	for _, ux := range up.cache.pool {
		arr = append(arr, ux)
	}
	up.Unlock()

	return arr, nil
}

// delete delete unspent of given hashes
func (up *Unspents) deleteWithTx(tx *bolt.Tx, hashes []cipher.SHA256) (cipher.SHA256, error) {
	var uxHash cipher.SHA256
	for _, hash := range hashes {
		ux, ok, err := up.pool.getWithTx(tx, hash)
		if err != nil {
			return cipher.SHA256{}, err
		}

		if !ok {
			continue
		}

		uxHash, err = up.meta.getXorHashWithTx(tx)
		if err != nil {
			return cipher.SHA256{}, err
		}

		uxHash = uxHash.Xor(ux.SnapshotHash())

		// update uxhash
		if err = up.meta.setXorHashWithTx(tx, uxHash); err != nil {
			return cipher.SHA256{}, err
		}

		if err := up.pool.deleteWithTx(tx, hash); err != nil {
			return cipher.SHA256{}, err
		}
	}

	return uxHash, nil
}

// Len returns the unspent outputs num
func (up *Unspents) Len() uint64 {
	up.Lock()
	defer up.Unlock()
	return uint64(len(up.cache.pool))
}

// Contains check if the hash of uxout does exist in the pool
func (up *Unspents) Contains(h cipher.SHA256) bool {
	up.Lock()
	_, ok := up.cache.pool[h.Hex()]
	up.Unlock()
	return ok
}

// GetUnspentsOfAddrs returns unspent outputs map of given addresses,
// the address as return map key, unspent outputs as value.
func (up *Unspents) GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts {
	up.Lock()
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
	up.Unlock()
	return addrUxs
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *Unspents) GetUxHash() cipher.SHA256 {
	up.Lock()
	defer up.Unlock()
	return up.cache.uxhash
}

func (up *Unspents) getUxHashFromDB() (cipher.SHA256, error) {
	if v := up.meta.Get(xorhashKey); v != nil {
		var hash cipher.SHA256
		copy(hash[:], v[:])
		return hash, nil
	}
	return cipher.SHA256{}, nil
}
