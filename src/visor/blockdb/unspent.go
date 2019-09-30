package blockdb

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

var (
	xorhashKey         = []byte("xorhash")
	addrIndexHeightKey = []byte("addr_index_height")

	// UnspentPoolBkt holds unspent outputs, indexed by unspent output hash
	UnspentPoolBkt = []byte("unspent_pool")
	// UnspentPoolAddrIndexBkt maps addresses to their unspent outputs
	UnspentPoolAddrIndexBkt = []byte("unspent_pool_addr_index")
	// UnspentMetaBkt holds unspent output metadata
	UnspentMetaBkt = []byte("unspent_meta")
)

// ErrUnspentNotExist is returned if an unspent is not found in the pool
type ErrUnspentNotExist struct {
	UxID string
}

// NewErrUnspentNotExist creates ErrUnspentNotExist from a UxID
func NewErrUnspentNotExist(uxID string) error {
	return ErrUnspentNotExist{
		UxID: uxID,
	}
}

func (e ErrUnspentNotExist) Error() string {
	return fmt.Sprintf("unspent output of %s does not exist", e.UxID)
}

// AddressHashes maps addresses to a set of hashes
type AddressHashes map[cipher.Address][]cipher.SHA256

// Flatten flattens all hash sets from AddressHashes to one slice
func (a AddressHashes) Flatten() []cipher.SHA256 {
	total := 0
	for _, h := range a {
		total += len(h)
	}

	hashes := make([]cipher.SHA256, total)

	i := 0
	for _, h := range a {
		copy(hashes[i:], h)
		i += len(h)
	}

	return hashes
}

type unspentMeta struct{}

func (m unspentMeta) getXorHash(tx *dbutil.Tx) (cipher.SHA256, error) {
	v, err := dbutil.GetBucketValue(tx, UnspentMetaBkt, xorhashKey)
	if err != nil {
		return cipher.SHA256{}, err
	} else if v == nil {
		return cipher.SHA256{}, nil
	}

	return cipher.SHA256FromBytes(v)
}

func (m *unspentMeta) setXorHash(tx *dbutil.Tx, hash cipher.SHA256) error {
	return dbutil.PutBucketValue(tx, UnspentMetaBkt, xorhashKey, hash[:])
}

func (m *unspentMeta) getAddrIndexHeight(tx *dbutil.Tx) (uint64, bool, error) {
	v, err := dbutil.GetBucketValue(tx, UnspentMetaBkt, addrIndexHeightKey)
	if err != nil {
		return 0, false, err
	} else if v == nil {
		return 0, false, nil
	}

	return dbutil.Btoi(v), true, nil
}

func (m *unspentMeta) setAddrIndexHeight(tx *dbutil.Tx, height uint64) error {
	return dbutil.PutBucketValue(tx, UnspentMetaBkt, addrIndexHeightKey, dbutil.Itob(height))
}

type pool struct{}

func (pl pool) get(tx *dbutil.Tx, hash cipher.SHA256) (*coin.UxOut, error) {
	var out coin.UxOut

	v, err := dbutil.GetBucketValueNoCopy(tx, UnspentPoolBkt, hash[:])
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	if err := decodeUxOutExact(v, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (pl pool) getAll(tx *dbutil.Tx) (coin.UxArray, error) {
	var uxa coin.UxArray

	if err := dbutil.ForEach(tx, UnspentPoolBkt, func(_, v []byte) error {
		var ux coin.UxOut
		if err := decodeUxOutExact(v, &ux); err != nil {
			return err
		}

		uxa = append(uxa, ux)
		return nil
	}); err != nil {
		return nil, err
	}

	return uxa, nil
}

func (pl pool) put(tx *dbutil.Tx, hash cipher.SHA256, ux coin.UxOut) error {
	buf, err := encodeUxOut(&ux)
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, UnspentPoolBkt, hash[:], buf)
}

func (pl *pool) delete(tx *dbutil.Tx, hash cipher.SHA256) error {
	return dbutil.Delete(tx, UnspentPoolBkt, hash[:])
}

type poolAddrIndex struct{}

func (p poolAddrIndex) get(tx *dbutil.Tx, addr cipher.Address) ([]cipher.SHA256, error) {
	var hashes hashesWrapper

	v, err := dbutil.GetBucketValueNoCopy(tx, UnspentPoolAddrIndexBkt, addr.Bytes())
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	if err := decodeHashesWrapperExact(v, &hashes); err != nil {
		return nil, err
	}

	return hashes.Hashes, nil
}

func (p poolAddrIndex) put(tx *dbutil.Tx, addr cipher.Address, hashes []cipher.SHA256) error {
	if len(hashes) == 0 {
		return errors.New("poolAddrIndex.put cannot put empty hash array")
	}

	hashesMap := make(map[cipher.SHA256]struct{}, len(hashes))
	for _, h := range hashes {
		if _, ok := hashesMap[h]; ok {
			return errors.New("poolAddrIndex.put: hashes array contains duplicate")
		}

		hashesMap[h] = struct{}{}
	}

	buf, err := encodeHashesWrapper(&hashesWrapper{
		Hashes: hashes,
	})
	if err != nil {
		return err
	}

	return dbutil.PutBucketValue(tx, UnspentPoolAddrIndexBkt, addr.Bytes(), buf)
}

// adjust adds and removes hashes from an address -> hashes index
// TODO -- if necessary, this can be optimized further to accept multiple addresses at once,
// so that all get queries can be performed before the set
func (p poolAddrIndex) adjust(tx *dbutil.Tx, addr cipher.Address, addHashes, rmHashes []cipher.SHA256) error {
	if len(addHashes) == 0 && len(rmHashes) == 0 {
		return nil
	}

	existingHashes, err := p.get(tx, addr)
	if err != nil {
		return err
	}

	rmHashesMap := make(map[cipher.SHA256]struct{}, len(rmHashes))
	for _, h := range rmHashes {
		rmHashesMap[h] = struct{}{}
	}

	if len(rmHashesMap) != len(rmHashes) {
		return errors.New("poolAddrIndex.adjust: rmHashes contains duplicates")
	}

	newHashesSize := len(existingHashes) - len(rmHashes)
	if newHashesSize < 0 {
		return errors.New("poolAddrIndex.adjust: rmHashes is longer than existingHashes")
	}

	newHashes := make([]cipher.SHA256, 0, newHashesSize)
	newHashesMap := make(map[cipher.SHA256]struct{}, newHashesSize)

	rmHashesCount := 0
	for _, h := range existingHashes {
		if _, ok := rmHashesMap[h]; ok {
			rmHashesCount++
		} else {
			newHashes = append(newHashes, h)
			newHashesMap[h] = struct{}{}
		}
	}

	if rmHashesCount != len(rmHashes) {
		return fmt.Errorf("poolAddrIndex.adjust: rmHashes contains %d hashes not indexed for address %s", len(rmHashes)-rmHashesCount, addr.String())
	}

	for _, h := range addHashes {
		if _, ok := rmHashesMap[h]; ok {
			return errors.New("poolAddrIndex.adjust: hash appears in both addHashes and rmHashes")
		}

		if _, ok := newHashesMap[h]; !ok {
			newHashes = append(newHashes, h)
			newHashesMap[h] = struct{}{}
		} else {
			return fmt.Errorf("poolAddrIndex.adjust: uxout hash %s is already indexed for address %s", h.Hex(), addr.String())
		}
	}

	// Delete the row if hashes is empty, so that the length of the bucket can
	// be used to determine the number of addresses with unspents
	if len(newHashes) == 0 {
		return dbutil.Delete(tx, UnspentPoolAddrIndexBkt, addr.Bytes())
	}

	return p.put(tx, addr, newHashes)
}

// Unspents unspent outputs pool
type Unspents struct {
	pool          *pool
	poolAddrIndex *poolAddrIndex
	meta          *unspentMeta
}

// NewUnspentPool creates new unspent pool instance
func NewUnspentPool() *Unspents {
	return &Unspents{
		pool:          &pool{},
		poolAddrIndex: &poolAddrIndex{},
		meta:          &unspentMeta{},
	}
}

// MaybeBuildIndexes builds indexes if necessary
func (up *Unspents) MaybeBuildIndexes(tx *dbutil.Tx, headSeq uint64) error {
	logger.Info("Unspents.MaybeBuildIndexes")

	// Compare the addrIndexHeight to the head block,
	// if not equal, rebuild the address index
	addrIndexHeight, ok, err := up.meta.getAddrIndexHeight(tx)
	if err != nil {
		return err
	}

	if ok && addrIndexHeight == headSeq {
		return nil
	}

	if addrIndexHeight > headSeq {
		logger.Critical().Warningf("addrIndexHeight > headSeq (%d > %d)", addrIndexHeight, headSeq)
	}

	logger.Infof("Rebuilding unspent_pool_addr_index (addrHeightIndexExists=%v, addrIndexHeight=%d, headSeq=%d)", ok, addrIndexHeight, headSeq)

	return up.buildAddrIndex(tx)
}

func (up *Unspents) buildAddrIndex(tx *dbutil.Tx) error {
	logger.Info("Building unspent address index")

	if err := dbutil.Reset(tx, UnspentPoolAddrIndexBkt); err != nil {
		return err
	}

	addrHashes := make(map[cipher.Address][]cipher.SHA256)

	var maxBlockSeq uint64
	if err := dbutil.ForEach(tx, UnspentPoolBkt, func(k, v []byte) error {
		var ux coin.UxOut
		if err := decodeUxOutExact(v, &ux); err != nil {
			return err
		}

		if ux.Head.BkSeq > maxBlockSeq {
			maxBlockSeq = ux.Head.BkSeq
		}

		h := ux.Hash()

		if !bytes.Equal(k[:], h[:]) {
			return errors.New("Unspent pool uxout.Hash() does not match its key")
		}

		addrHashes[ux.Body.Address] = append(addrHashes[ux.Body.Address], h)

		return nil
	}); err != nil {
		return err
	}

	if len(addrHashes) == 0 {
		logger.Infof("No unspents to index")
		return nil
	}

	for addr, hashes := range addrHashes {
		if err := up.poolAddrIndex.put(tx, addr, hashes); err != nil {
			return err
		}
	}

	if err := up.meta.setAddrIndexHeight(tx, maxBlockSeq); err != nil {
		return err
	}

	logger.Infof("Indexed unspents for %d addresses", len(addrHashes))

	return nil
}

// ProcessBlock adds unspents from a block to the unspent pool
func (up *Unspents) ProcessBlock(tx *dbutil.Tx, b *coin.SignedBlock) error {
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
	rmAddrHashes := make(map[cipher.Address][]cipher.SHA256)
	for _, ux := range uxs {
		xorHash = xorHash.Xor(ux.SnapshotHash())

		h := ux.Hash()

		if err := up.pool.delete(tx, h); err != nil {
			return err
		}

		rmAddrHashes[ux.Body.Address] = append(rmAddrHashes[ux.Body.Address], h)
	}

	// Create new outputs
	txnUxHashes := make([]cipher.SHA256, len(txnUxs))
	addAddrHashes := make(map[cipher.Address][]cipher.SHA256)
	for i, ux := range txnUxs {
		h := ux.Hash()
		txnUxHashes[i] = h
		addAddrHashes[ux.Body.Address] = append(addAddrHashes[ux.Body.Address], h)
	}

	// Check that the uxout exists in the pool already, otherwise xorHash will be calculated wrong
	for _, h := range txnUxHashes {
		if hasKey, err := up.Contains(tx, h); err != nil {
			return err
		} else if hasKey {
			return fmt.Errorf("attempted to insert uxout:%v twice into the unspent pool", h.Hex())
		}
	}

	for i, ux := range txnUxs {
		// Add new outputs
		if err := up.pool.put(tx, txnUxHashes[i], ux); err != nil {
			return err
		}

		// Recalculate xorHash
		xorHash = xorHash.Xor(ux.SnapshotHash())
	}

	// Set xorHash
	if err := up.meta.setXorHash(tx, xorHash); err != nil {
		return err
	}

	// Update indexes
	for addr, rmHashes := range rmAddrHashes {
		addHashes := addAddrHashes[addr]

		if err := up.poolAddrIndex.adjust(tx, addr, addHashes, rmHashes); err != nil {
			return err
		}

		delete(addAddrHashes, addr)
	}

	for addr, addHashes := range addAddrHashes {
		if err := up.poolAddrIndex.adjust(tx, addr, addHashes, nil); err != nil {
			return err
		}
	}

	// Check that the addrIndexHeight is incremental
	addrIndexHeight, ok, err := up.meta.getAddrIndexHeight(tx)
	if err != nil {
		return err
	}

	if b.Block.Head.BkSeq == 0 {
		if ok {
			err := errors.New("addrIndexHeight is set but no block has been indexed yet")
			logger.Critical().Error(err.Error())
			return err
		}
	} else if b.Block.Head.BkSeq != addrIndexHeight+1 {
		err := errors.New("unspent pool processing blocks out of order")
		logger.Critical().Error(err.Error())
		return err
	}

	// Update the addrIndexHeight
	return up.meta.setAddrIndexHeight(tx, b.Block.Head.BkSeq)
}

// GetArray returns UxOut for a set of hashes, will return error if any of the hashes do not exist in the pool.
func (up *Unspents) GetArray(tx *dbutil.Tx, hashes []cipher.SHA256) (coin.UxArray, error) {
	var uxa coin.UxArray

	for _, h := range hashes {
		ux, err := up.pool.get(tx, h)
		if err != nil {
			return nil, err
		} else if ux == nil {
			return nil, NewErrUnspentNotExist(h.Hex())
		}

		uxa = append(uxa, *ux)
	}

	return uxa, nil
}

// Get returns the uxout value of given hash
func (up *Unspents) Get(tx *dbutil.Tx, h cipher.SHA256) (*coin.UxOut, error) {
	return up.pool.get(tx, h)
}

// GetAll returns Pool as an array. Note: they are not in any particular order.
func (up *Unspents) GetAll(tx *dbutil.Tx) (coin.UxArray, error) {
	return up.pool.getAll(tx)
}

// Len returns the unspent outputs num
func (up *Unspents) Len(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, UnspentPoolBkt)
}

// Contains check if the hash of uxout does exist in the pool
func (up *Unspents) Contains(tx *dbutil.Tx, h cipher.SHA256) (bool, error) {
	return dbutil.BucketHasKey(tx, UnspentPoolBkt, h[:])
}

// GetUnspentHashesOfAddrs returns a map of addresses to their unspent output hashes
func (up *Unspents) GetUnspentHashesOfAddrs(tx *dbutil.Tx, addrs []cipher.Address) (AddressHashes, error) {
	addrHashes := make(AddressHashes, len(addrs))

	for _, addr := range addrs {
		hashes, err := up.poolAddrIndex.get(tx, addr)
		if err != nil {
			return nil, err
		}

		addrHashes[addr] = hashes
	}

	return addrHashes, nil
}

// GetUnspentsOfAddrs returns a map of addresses to their unspent outputs
func (up *Unspents) GetUnspentsOfAddrs(tx *dbutil.Tx, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrUxs := make(coin.AddressUxOuts, len(addrs))

	for _, addr := range addrs {
		hashes, err := up.poolAddrIndex.get(tx, addr)
		if err != nil {
			return nil, err
		}

		uxa, err := up.GetArray(tx, hashes)
		if err != nil {
			switch e := err.(type) {
			case ErrUnspentNotExist:
				logger.Critical().Errorf("Unspent hash %s indexed under address %s does not exist in unspent pool", e.UxID, addr.String())
			}
			return nil, err
		}

		addrUxs[addr] = uxa
	}

	return addrUxs, nil
}

// GetUxHash returns unspent output checksum for the Block.
// Must be called after Block is fully initialized,
// and before its outputs are added to the unspent pool
func (up *Unspents) GetUxHash(tx *dbutil.Tx) (cipher.SHA256, error) {
	return up.meta.getXorHash(tx)
}

// AddressCount returns the total number of addresses with unspents
func (up *Unspents) AddressCount(tx *dbutil.Tx) (uint64, error) {
	return dbutil.Len(tx, UnspentPoolAddrIndexBkt)
}
