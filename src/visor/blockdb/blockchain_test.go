package blockdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024

	genTime      uint64 = 1000
	incTime      uint64 = 3600 * 1000
	genCoins     uint64 = 1000e6
	genCoinHours uint64 = 1000 * 1000
)

func feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

type fakeStorage struct {
	tree    *fakeBlockTree
	sigs    *fakeSignatureStore
	unspent *fakeUnspentPool
}

func newFakeStorage() *fakeStorage {
	var failedWhenSaved bool
	// var failedWhenSaved2 bool
	// var failedWhenSaved3 bool
	return &fakeStorage{
		tree:    newFakeBlockTree(&failedWhenSaved),
		sigs:    newFakeSigStore(&failedWhenSaved),
		unspent: newFakeUnspentPool(&failedWhenSaved),
	}
}

type fakeBlockTree struct {
	blocks     map[string]*coin.Block
	saveFailed bool

	// state tracking: do not configure directly
	// set to true if saveFailed was true and certain operations were performed
	failedWhenSaved *bool
}

func newFakeBlockTree(failedWhenSaved *bool) *fakeBlockTree {
	return &fakeBlockTree{
		blocks:          make(map[string]*coin.Block),
		failedWhenSaved: failedWhenSaved,
	}
}

func (bt *fakeBlockTree) AddBlock(tx *bolt.Tx, b *coin.Block) error {
	if bt.saveFailed {
		if bt.failedWhenSaved != nil {
			*bt.failedWhenSaved = true
		}
		return errors.New("intentionally failed")
	}
	bt.blocks[b.HashHeader().Hex()] = b
	return nil
}

func (bt *fakeBlockTree) GetBlock(tx *bolt.Tx, hash cipher.SHA256) (*coin.Block, error) {
	if bt.failedWhenSaved != nil && *bt.failedWhenSaved {
		return nil, nil
	}
	return bt.blocks[hash.Hex()], nil
}

func (bt *fakeBlockTree) GetBlockInDepth(tx *bolt.Tx, dep uint64, filter Walker) (*coin.Block, error) {
	return nil, nil
}

func (bt *fakeBlockTree) ForEachBlock(tx *bolt.Tx, f func(*coin.Block) error) error {
	return nil
}

type fakeSignatureStore struct {
	db         *bolt.DB
	sigs       map[string]cipher.Sig
	saveFailed bool
	getSigErr  error

	failedWhenSaved *bool
}

func newFakeSigStore(failedWhenSaved *bool) *fakeSignatureStore {
	return &fakeSignatureStore{
		sigs:            make(map[string]cipher.Sig),
		failedWhenSaved: failedWhenSaved,
	}
}

func (ss *fakeSignatureStore) Add(tx *bolt.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	if ss.saveFailed {
		if ss.failedWhenSaved != nil {
			*ss.failedWhenSaved = true
		}
		return errors.New("intentionally failed")
	}

	ss.sigs[hash.Hex()] = sig
	return nil
}

func (ss *fakeSignatureStore) Get(tx *bolt.Tx, hash cipher.SHA256) (cipher.Sig, bool, error) {
	if ss.failedWhenSaved != nil && *ss.failedWhenSaved {
		return cipher.Sig{}, false, nil
	}

	if ss.getSigErr != nil {
		return cipher.Sig{}, false, ss.getSigErr
	}

	sig, ok := ss.sigs[hash.Hex()]
	return sig, ok, nil
}

func (ss *fakeSignatureStore) ForEach(tx *bolt.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
	return nil
}

type fakeUnspentPool struct {
	outs       map[cipher.SHA256]coin.UxOut
	uxHash     cipher.SHA256
	saveFailed bool

	failedWhenSaved *bool
}

func newFakeUnspentPool(failedWhenSaved *bool) *fakeUnspentPool {
	return &fakeUnspentPool{
		outs:            make(map[cipher.SHA256]coin.UxOut),
		failedWhenSaved: failedWhenSaved,
	}
}

func (fup *fakeUnspentPool) Len() uint64 {
	return uint64(len(fup.outs))
}

func (fup *fakeUnspentPool) Get(h cipher.SHA256) (coin.UxOut, bool) {
	out, ok := fup.outs[h]
	return out, ok
}

func (fup *fakeUnspentPool) GetAll() (coin.UxArray, error) {
	outs := make(coin.UxArray, 0, len(fup.outs))
	for _, out := range fup.outs {
		outs = append(outs, out)
	}

	return outs, nil
}

func (fup *fakeUnspentPool) GetArray(hashes []cipher.SHA256) (coin.UxArray, error) {
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

func (fup *fakeUnspentPool) GetUxHash() cipher.SHA256 {
	return fup.uxHash
}

func (fup *fakeUnspentPool) GetUnspentsOfAddrs(addrs []cipher.Address) coin.AddressUxOuts {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
	}

	addrOutMap := make(coin.AddressUxOuts)
	for _, out := range fup.outs {
		addr := out.Body.Address
		addrOutMap[addr] = append(addrOutMap[addr], out)
	}

	return addrOutMap
}

func (fup *fakeUnspentPool) ProcessBlock(b *coin.SignedBlock) dbutil.TxHandler {
	return func(tx *bolt.Tx) (dbutil.Rollback, error) {
		if fup.saveFailed {
			if fup.failedWhenSaved != nil {
				*fup.failedWhenSaved = true
			}
			return func() {}, errors.New("intentionally failed")
		}
		return func() {}, nil
	}
}

func (fup *fakeUnspentPool) Contains(h cipher.SHA256) bool {
	_, ok := fup.outs[h]
	return ok
}

func DefaultWalker(tx *bolt.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
	return hps[0].Hash, true
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

	type failedSaves struct {
		tree    bool
		sigs    bool
		unspent bool
	}

	tt := []struct {
		name        string
		fakeStorage *fakeStorage
		failedSaves failedSaves
		expect      expect
	}{
		{
			"ok",
			newFakeStorage(),
			failedSaves{},
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
			newFakeStorage(),
			failedSaves{
				sigs: true,
			},
			expect{
				errors.New("save signature failed: intentionally failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
		{
			"save block failed",
			newFakeStorage(),
			failedSaves{
				tree: true,
			},
			expect{
				errors.New("save block failed: intentionally failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
		{
			"unspent process block failed",
			newFakeStorage(),
			failedSaves{
				unspent: true,
			},
			expect{
				errors.New("intentionally failed"),
				false,
				false,
				false,
				uint64(0),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			tc.fakeStorage.tree.saveFailed = tc.failedSaves.tree
			tc.fakeStorage.sigs.saveFailed = tc.failedSaves.sigs
			tc.fakeStorage.unspent.saveFailed = tc.failedSaves.unspent

			bc, err := createBlockchain(db, DefaultWalker, tc.fakeStorage.tree, tc.fakeStorage.sigs, tc.fakeStorage.unspent)
			require.NoError(t, err)

			gb := makeGenesisBlock(t)

			err = db.Update(func(tx *bolt.Tx) error {
				return bc.AddBlock(tx, &gb)
			})

			require.Equal(t, tc.expect.err, err)

			// check sig
			err = db.View(func(tx *bolt.Tx) error {
				_, ok, err := tc.fakeStorage.sigs.Get(tx, gb.HashHeader())
				require.NoError(t, err)
				require.Equal(t, tc.expect.sigSaved, ok)

				// check block in tree
				b, err := tc.fakeStorage.tree.GetBlock(tx, gb.HashHeader())
				require.NoError(t, err)
				require.Equal(t, tc.expect.blockSaved, b != nil)

				return nil
			})
			require.NoError(t, err)

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
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = bc.Head(tx)
		require.Equal(t, err, ErrNoHeadBlock)

		gb := makeGenesisBlock(t)

		err := bc.AddBlock(tx, &gb)
		require.NoError(t, err)

		b, err := bc.Head(tx)
		require.NoError(t, err)
		require.Equal(t, gb.HashHeader().Hex(), b.HashHeader().Hex())

		return nil
	})
	require.NoError(t, err)
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
				NewErrSignatureLost(&gb.Block),
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			bc, err := NewBlockchain(db, DefaultWalker)
			require.NoError(t, err)

			bc.tree = tc.tree
			bc.sigs = tc.sigs

			err = db.View(func(tx *bolt.Tx) error {
				b, err := bc.GetSignedBlockByHash(tx, tc.hash)
				require.Equal(t, tc.expect.err, err)
				require.Equal(t, tc.expect.b, b)
				return nil
			})
			require.NoError(t, err)
		})
	}
}
