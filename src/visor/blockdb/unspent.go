package blockdb

import (
	"fmt"

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

// UnspentGetter provides unspent pool related querying methods
type UnspentGetter interface {
	// GetUnspentsOfAddrs returns all unspent outputs of given addresses
	GetUnspentsOfAddrs(*bolt.Tx, []cipher.Address) coin.AddressUxOuts
	Get(*bolt.Tx, cipher.SHA256) (coin.UxOut, bool)
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

func (pl pool) get(tx *bolt.Tx, hash cipher.SHA256) (*coin.UxOut, error) {
	var out coin.UxOut

	if ok, err := dbutil.GetBucketObjectDecoded(tx, unspentPoolBkt, hash[:], &out); err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}

	return &out, nil
}

func (pl pool) getAll(tx *bolt.Tx) (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := dbutil.ForEach(tx, unspentPoolBkt, func(_, v []byte) error {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		uxa = append(uxa, ux)
		return nil
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

func (pl pool) set(tx *bolt.Tx, hash cipher.SHA256, ux coin.UxOut) error {
	return dbutil.PutBucketValue(tx, unspentPoolBkt, hash[:], encoder.Serialize(ux))
}

func (pl *pool) delete(tx *bolt.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, unspentPoolBkt, hash[:])
}

// Unspents unspent outputs pool
type Unspents struct {
	pool *pool
	meta *unspentMeta
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
		pool: &pool{},
		meta: &unspentMeta{},
	}

	return up, nil
}

// ProcessBlock adds unspents from a block to the unspent pool
func (up *Unspents) ProcessBlock(tx *bolt.Tx, b *coin.SignedBlock) error {
	// Gather all transaction inputs
	var inputs []cipher.SHA256
	var txnUxs coin.UxArray
	for _, txn := range b.Body.Transactions {
		inputs = append(inputs, txn.In...)
		txnUxs = append(txnUxs, coin.CreateUnspents(b.Head, txn)...)
	}

	uxs, err := up.GetArray(tx, inputs)
	if err != nil {
		return err
	}

	xorHash, err := up.meta.getXorHash(tx)
	if err != nil {
		return err
	}

	// Remove spent outputs
	for _, ux := range uxs {
		xorHash = xorHash.Xor(ux.SnapshotHash())

		if err := up.pool.delete(tx, ux.Hash()); err != nil {
			return err
		}
	}

	// Create new outputs
	txnUxHashes := make([]cipher.SHA256, len(txnUxs))
	for i, ux := range txnUxs {
		txnUxHashes[i] = ux.Hash()
	}

	// Check that the uxout does exist in the pool already, otherwise xorHash will be calculated wrong
	for _, h := range txnUxHashes {
		if hasKey, err := up.Contains(tx, h); err != nil {
			return err
		} else if hasKey {
			return fmt.Errorf("attempted to insert uxout:%v twice into the unspent pool", h.Hex())
		}
	}

	// Add new outputs
	for i, ux := range txnUxs {
		if err := up.pool.set(tx, txnUxHashes[i], ux); err != nil {
			return err
		}

		xorHash = xorHash.Xor(ux.SnapshotHash())
	}

	return up.meta.setXorHash(tx, xorHash)
}

// GetArray returns UxOut for a set of hashes, will return error if any of the hashes do not exist in the pool.
func (up *Unspents) GetArray(tx *bolt.Tx, hashes []cipher.SHA256) (coin.UxArray, error) {
	var uxa coin.UxArray

	for _, h := range hashes {
		ux, err := up.pool.get(tx, h)
		if err != nil {
			return nil, err
		} else if ux == nil {
			return nil, fmt.Errorf("unspent output does not exist: %s", h.Hex())
		}

		uxa = append(uxa, *ux)
	}

	return uxa, nil
}

// Get returns the uxout value of given hash
func (up *Unspents) Get(tx *bolt.Tx, h cipher.SHA256) (*coin.UxOut, error) {
	return up.pool.get(tx, h)
}

// // GetForTransactionInputs returns a map of txn hashes to unspent outputs for a transaction's input
// func (up *Unspents) GetForTransactionInputs(tx *bolt.Tx, txns coin.Transactions) (coin.TransactionUnspents, error) {
// 	txnOuts := make(coin.TransactionUnspents, len(txns))

// 	for _, txn := range txns {
// 		outs, err := up.GetArray(tx, txn.In)
// 		if err != nil {
// 			return nil, err
// 		}
// 		txnOuts[txn.Hash()] = outs
// 	}

// 	return txnOuts, nil
// }

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *Unspents) GetAll(tx *bolt.Tx) (coin.UxArray, error) {
	return up.pool.getAll(tx)
}

// Len returns the unspent outputs num
func (up *Unspents) Len(tx *bolt.Tx) (uint64, error) {
	return dbutil.Len(tx, unspentPoolBkt)
}

// Contains check if the hash of uxout does exist in the pool
func (up *Unspents) Contains(tx *bolt.Tx, h cipher.SHA256) (bool, error) {
	return dbutil.BucketHasKey(tx, unspentPoolBkt, h[:])
}

// GetUnspentsOfAddrs returns unspent outputs map of given addresses,
// the address as return map key, unspent outputs as value.
func (up *Unspents) GetUnspentsOfAddrs(tx *bolt.Tx, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
	}

	addrUxs := make(coin.AddressUxOuts, len(addrs))

	if err := dbutil.ForEach(tx, unspentPoolBkt, func(k, v []byte) error {
		var ux coin.UxOut
		if err := encoder.DeserializeRaw(v, &ux); err != nil {
			return err
		}

		addr := ux.Body.Address
		if _, ok := addrm[addr]; ok {
			addrUxs[addr] = append(addrUxs[addr], ux)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return addrUxs, nil
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *Unspents) GetUxHash(tx *bolt.Tx) (cipher.SHA256, error) {
	return up.meta.getXorHash(tx)
}
