package blockdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/bucket"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024

	genTime      uint64 = 1000
	incTime      uint64 = 3600 * 1000
	genCoins     uint64 = 1000e6
	genCoinHours uint64 = 1000 * 1000

	failedWhenSave bool
)

func _feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func cleanState() {
	failedWhenSave = false
}

type fakeBlockTree struct {
	blocks     map[string]*coin.Block
	saveFailed bool
}

func newFakeBlockTree() *fakeBlockTree {
	return &fakeBlockTree{
		blocks: make(map[string]*coin.Block),
	}
}

func (bt fakeBlockTree) AddBlockWithTx(tx *bolt.Tx, b *coin.Block) error {
	if bt.saveFailed {
		failedWhenSave = true
		return errors.New("intentional failed")
	}
	bt.blocks[b.HashHeader().Hex()] = b
	return nil
}

func (bt fakeBlockTree) GetBlock(hash cipher.SHA256) *coin.Block {
	if failedWhenSave {
		return nil
	}
	return bt.blocks[hash.Hex()]
}

func (bt fakeBlockTree) GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block {
	return nil
}

type fakeSignatureStore struct {
	db         *bolt.DB
	sigs       map[string]cipher.Sig
	saveFailed bool
	getSigErr  error
}

func newFakeSigStore() *fakeSignatureStore {
	return &fakeSignatureStore{
		sigs: make(map[string]cipher.Sig),
	}
}

func (ss fakeSignatureStore) AddWithTx(tx *bolt.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	if ss.saveFailed {
		failedWhenSave = true
		return errors.New("intentional failed")
	}

	ss.sigs[hash.Hex()] = sig
	return nil
}

func (ss fakeSignatureStore) Get(hash cipher.SHA256) (cipher.Sig, bool, error) {
	if failedWhenSave {
		return cipher.Sig{}, false, nil
	}

	if ss.getSigErr != nil {
		return cipher.Sig{}, false, ss.getSigErr
	}

	sig, ok := ss.sigs[hash.Hex()]
	return sig, ok, nil
}

type fakeUnspentPool struct {
	outs       map[cipher.SHA256]coin.UxOut
	uxHash     cipher.SHA256
	saveFailed bool
}

func newFakeUnspentsPool() *fakeUnspentPool {
	return &fakeUnspentPool{
		outs: make(map[cipher.SHA256]coin.UxOut),
	}
}

func (fup fakeUnspentPool) Len() uint64 {
	return uint64(len(fup.outs))
}

func (fup fakeUnspentPool) Get(h cipher.SHA256) (coin.UxOut, bool) {
	out, ok := fup.outs[h]
	return out, ok
}

func (fup fakeUnspentPool) GetAll() (coin.UxArray, error) {
	outs := make(coin.UxArray, 0, len(fup.outs))
	for _, out := range fup.outs {
		outs = append(outs, out)
	}

	return outs, nil
}

func (fup fakeUnspentPool) GetArray(hashes []cipher.SHA256) (coin.UxArray, error) {
	outs := make(coin.UxArray, 0, len(hashes))
	for _, h := range hashes {
		ux, ok := fup.outs[h]
		if !ok {
			return nil, fmt.Errorf("unspent output of %s does not exist", h.Hex())
		}

		outs = append(outs, ux)
	}
	return outs, nil
}

func (fup fakeUnspentPool) GetUxHash() cipher.SHA256 {
	return fup.uxHash
}

func (fup fakeUnspentPool) GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts {
	addrOutMap := map[cipher.Address]coin.UxArray{}
	for _, out := range fup.outs {
		addr := out.Body.Address
		addrOutMap[addr] = append(addrOutMap[addr], out)
	}

	return addrOutMap
}

func (fup fakeUnspentPool) ProcessBlock(b *coin.SignedBlock) bucket.TxHandler {
	return func(tx *bolt.Tx) (bucket.Rollback, error) {
		if fup.saveFailed {
			failedWhenSave = true
			return func() {}, errors.New("intentional failed")
		}
		return func() {}, nil
	}
}

func (fup fakeUnspentPool) Contains(h cipher.SHA256) bool {
	_, ok := fup.outs[h]
	return ok
}

func TestNewBlockchain(t *testing.T) {
	// walker := func(hps []coin.HashPair) cipher.SHA256 {
	// 	return hps[0].Hash
	// }

	// tt := []struct {
	// 	name string
	// 	wlk  func(hps []coin.HashPair) cipher.SHA256
	// 	tree blockTree
	// 	sigs signatureStore
	// 	err  error
	// }{
	// 	{
	// 		"ok",
	// 		false,
	// 		walker,
	// 		&fakeBlockTree{},
	// 		&fakeSignatureStore{},
	// 		nil,
	// 	},
	// }

	// for _, tc := range tt {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		db, err := testutil.PrepareDB(t)
	// 		require.NoError(t, err)
	// 		bc, err := NewBlockchain(db, walker)
	// 		require.Equal(t, , actual interface{}, msgAndArgs ...interface{})
	// 	})
	// }

	// bc, err := NewBlockchain(db, func(hps []coin.HashPair) cipher.SHA256 {
	// 	return hps[0].Hash
	// })

	// assert.Nil(t, err)
	// assert.NotNil(t, bc.db)
	// assert.NotNil(t, bc.UnspentPool())
	// assert.NotNil(t, bc.meta)

	// // check the existence of buckets
	// db.View(func(tx *bolt.Tx) error {
	// 	assert.NotNil(t, tx.Bucket([]byte("unspent_pool")))
	// 	assert.NotNil(t, tx.Bucket([]byte("unspent_meta")))
	// 	assert.NotNil(t, tx.Bucket([]byte("blockchain_meta")))
	// 	return nil
	// })
}

func DefaultWalker(hps []coin.HashPair) cipher.SHA256 {
	return hps[0].Hash
}

func makeGenesisBlock(t *testing.T) coin.SignedBlock {
	gb, err := coin.NewGenesisBlock(genAddress, genCoinHours, genTime)
	require.NoError(t, err)

	sig := cipher.SignHash(gb.HashHeader(), genSecret)
	return coin.SignedBlock{
		Block: *gb,
		Sig:   sig,
	}
}

func TestBlockchainAddBlockWithTx(t *testing.T) {
	type expect struct {
		err           error
		sigSaved      bool
		blockSaved    bool
		genesisCached bool
		headSeq       uint64
	}

	tt := []struct {
		name     string
		tree     BlockTree
		sigs     BlockSigs
		unspents UnspentPool
		expect   expect
	}{
		{
			"ok",
			newFakeBlockTree(),
			newFakeSigStore(),
			newFakeUnspentsPool(),
			expect{
				nil,
				true,
				true,
				true,
				uint64(0),
			},
		},
		{
			"save sig failed",
			newFakeBlockTree(),
			fakeSignatureStore{saveFailed: true},
			newFakeUnspentsPool(),
			expect{
				errors.New("save signature failed: intentional failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
		{
			"save block failed",
			fakeBlockTree{saveFailed: true},
			newFakeSigStore(),
			newFakeUnspentsPool(),
			expect{
				errors.New("save block failed: intentional failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
		{
			"unspent process block failed",
			newFakeBlockTree(),
			newFakeSigStore(),
			fakeUnspentPool{saveFailed: true},
			expect{
				errors.New("intentional failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cleanState()
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()
			bc, err := createBlockchain(db,
				DefaultWalker,
				tc.tree,
				tc.sigs,
				tc.unspents)
			require.NoError(t, err)

			gb := makeGenesisBlock(t)

			err = db.Update(func(tx *bolt.Tx) error {
				return bc.AddBlockWithTx(tx, &gb)
			})

			require.Equal(t, tc.expect.err, err)

			// check sig
			_, ok, err := tc.sigs.Get(gb.HashHeader())
			require.NoError(t, err)
			require.Equal(t, tc.expect.sigSaved, ok)

			// check block in tree
			b := tc.tree.GetBlock(gb.HashHeader())
			require.Equal(t, tc.expect.blockSaved, b != nil)

			// check cache of head seq
			require.Equal(t, tc.expect.headSeq, bc.cache.headSeq)

			require.Equal(t, tc.expect.genesisCached, bc.cache.genesisBlock != nil)
			if tc.expect.genesisCached {
				require.Equal(t, gb.HashHeader().Hex(), bc.cache.genesisBlock.HashHeader().Hex())
			}
		})
	}

}

func TestBlockchainHead(t *testing.T) {
	cleanState()
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	_, err = bc.Head()
	require.EqualError(t, err, "found no head block: 0")

	gb := makeGenesisBlock(t)
	db.Update(func(tx *bolt.Tx) error {
		err := bc.AddBlockWithTx(tx, &gb)
		require.NoError(t, err)
		return nil
	})

	b, err := bc.Head()
	require.NoError(t, err)
	require.Equal(t, gb.HashHeader().Hex(), b.HashHeader().Hex())
}

func TestBlockchainLen(t *testing.T) {
	bc := Blockchain{}
	require.Equal(t, uint64(0), bc.Len())

	gb := makeGenesisBlock(t)
	bc.cache.genesisBlock = &gb
	require.Equal(t, uint64(1), bc.Len())

	bc.cache.headSeq = 1
	require.Equal(t, uint64(2), bc.Len())
}

func TestBlockchainGetBlockByHash(t *testing.T) {
	gb := makeGenesisBlock(t)

	type expect struct {
		err error
		b   *coin.SignedBlock
	}

	tt := []struct {
		name   string
		tree   BlockTree
		sigs   BlockSigs
		hash   cipher.SHA256
		expect expect
	}{
		{
			"ok",
			&fakeBlockTree{
				blocks: map[string]*coin.Block{
					gb.HashHeader().Hex(): &gb.Block,
				},
			},
			&fakeSignatureStore{
				sigs: map[string]cipher.Sig{
					gb.HashHeader().Hex(): gb.Sig,
				},
			},
			gb.HashHeader(),
			expect{
				nil,
				&gb,
			},
		},
		{
			"block not exist",
			&fakeBlockTree{
				blocks: map[string]*coin.Block{},
			},
			&fakeSignatureStore{
				sigs: map[string]cipher.Sig{},
			},
			gb.HashHeader(),
			expect{
				nil,
				nil,
			},
		},
		{
			"signature not exist",
			&fakeBlockTree{
				blocks: map[string]*coin.Block{
					gb.HashHeader().Hex(): &gb.Block,
				},
			},
			&fakeSignatureStore{
				sigs: map[string]cipher.Sig{},
			},
			gb.HashHeader(),
			expect{
				fmt.Errorf("find no signature of block: %v", gb.HashHeader().Hex()),
				nil,
			},
		},
		{
			"get signature error",
			&fakeBlockTree{
				blocks: map[string]*coin.Block{
					gb.HashHeader().Hex(): &gb.Block,
				},
			},
			&fakeSignatureStore{
				getSigErr: errors.New("intentional error"),
				sigs:      map[string]cipher.Sig{},
			},
			gb.HashHeader(),
			expect{
				fmt.Errorf("find signature of block: %v failed: intentional error", gb.HashHeader().Hex()),
				nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cleanState()

			bc := Blockchain{
				tree: tc.tree,
				sigs: tc.sigs,
			}

			b, err := bc.GetBlockByHash(tc.hash)
			require.Equal(t, tc.expect.err, err)
			require.Equal(t, tc.expect.b, b)
		})
	}
}
