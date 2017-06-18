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
	xorhashKey      = []byte("xorhash")
	parsedHeightKey = []byte("parsed_height")
)

// UnspentPool unspent outputs pool
type UnspentPool struct {
	db   *bolt.DB
	pool *bucket.Bucket
	meta *bucket.Bucket
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

// Reset resets the unspent pool buckets
func (up *UnspentPool) Reset() error {
	return up.db.Update(func(tx *bolt.Tx) error {
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
	})
}

// ProcessBlock processes block
func (up *UnspentPool) ProcessBlock(b *coin.Block) error {
	txns := b.Body.Transactions
	return up.db.Update(func(tx *bolt.Tx) error {
		for _, txn := range txns {
			// Remove spent outputs
			if err := up.delete(tx, txn.In); err != nil {
				return err
			}

			// Create new outputs
			txUxs := coin.CreateUnspents(b.Head, txn)
			for i := range txUxs {
				if err := up.add(tx, txUxs[i]); err != nil {
					return err
				}
			}

			// uxs = append(uxs, txUxs...)
		}
		return nil
	})
}

// Add adds a UxOut to pool
func (up *UnspentPool) Add(ux coin.UxOut) error {
	return up.db.Update(func(tx *bolt.Tx) error {
		return up.add(tx, ux)
	})
}

func (up *UnspentPool) add(tx *bolt.Tx, ux coin.UxOut) (err error) {
	// all updates will rollback if return is not nil
	// in case of unexpected panic, we must catch it and return error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unspent pool add uxout failed: %v", err)
		}
	}()

	h := ux.Hash()
	if up.Contains(h) {
		return fmt.Errorf("attemps to insert uxout:%v twice into the unspent pool", h.Hex())
	}

	mb := tx.Bucket(up.meta.Name)
	xorhash, err := getXorHash(mb)
	if err != nil {
		return err
	}

	xorhash = xorhash.Xor(ux.SnapshotHash())

	pb := tx.Bucket(up.pool.Name)
	if err := setUxOut(pb, h, ux); err != nil {
		return err
	}

	if err = setXorHash(mb, xorhash); err != nil {
		return err
	}

	return setHeadSeq(mb, ux.Head.BkSeq)
}

func setUxOut(bkt *bolt.Bucket, hash cipher.SHA256, ux coin.UxOut) error {
	v := encoder.Serialize(ux)
	return bkt.Put([]byte(hash.Hex()), v)
}

func getXorHash(bkt *bolt.Bucket) (cipher.SHA256, error) {
	if v := bkt.Get(xorhashKey); v != nil {
		return cipher.SHA256FromHex(string(v))
	}

	return cipher.SHA256{}, nil
}

func setXorHash(bkt *bolt.Bucket, hash cipher.SHA256) error {
	return bkt.Put(xorhashKey, []byte(hash.Hex()))
}

func setHeadSeq(bkt *bolt.Bucket, seq uint64) error {
	return bkt.Put(parsedHeightKey, bucket.Itob(seq))
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
func (up *UnspentPool) del(tx *bolt.Tx, h cipher.SHA256) error {
	if v := up.pool.Get([]byte(h.Hex())); v != nil {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		mb := tx.Bucket(up.meta.Name)
		hash, err := getXorHash(mb)
		if err != nil {
			return err
		}

		hash = hash.Xor(ux.SnapshotHash())

		// update uxhash
		if err = setXorHash(mb, hash); err != nil {
			return err
		}

		pb := tx.Bucket(up.pool.Name)
		return deleteUxOut(pb, h)
	}

	return nil
}

func deleteUxOut(bkt *bolt.Bucket, hash cipher.SHA256) error {
	return bkt.Delete([]byte(hash.Hex()))
}

// delete delete unspent of given hashes
func (up *UnspentPool) delete(tx *bolt.Tx, hashes []cipher.SHA256) error {
	for _, hash := range hashes {
		if err := up.del(tx, hash); err != nil {
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
	if v := up.meta.Get(xorhashKey); v != nil {
		hash, err := cipher.SHA256FromHex(string(v))
		if err != nil {
			return cipher.SHA256{}, fmt.Errorf("parse xorhash failed: %v", err)
		}
		return hash, nil
	}

	return cipher.SHA256{}, nil
}

// HeadSeq returns head block seq
func (up *UnspentPool) HeadSeq() int64 {
	if v := up.meta.Get(parsedHeightKey); v != nil {
		return int64(bucket.Btoi(v))
	}

	return -1
}
