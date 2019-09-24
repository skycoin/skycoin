package blockdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
)

func prepareDB(t *testing.T) (*dbutil.DB, func()) {
	db, shutdown := testutil.PrepareDB(t)

	err := db.Update("", func(tx *dbutil.Tx) error {
		return CreateBuckets(tx)
	})
	if err != nil {
		shutdown()
		t.Fatalf("CreateBuckets failed: %v", err)
	}

	return db, shutdown
}

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)

	genTime      uint64 = 1000
	genCoinHours uint64 = 1000 * 1000
)

func feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

type fakeStorage struct {
	tree      *fakeBlockTree
	sigs      *fakeSignatureStore
	unspent   *fakeUnspentPool
	chainMeta *fakeChainMeta
}

func newFakeStorage() *fakeStorage {
	var failedWhenSaved bool
	return &fakeStorage{
		tree:      newFakeBlockTree(&failedWhenSaved),
		sigs:      newFakeSigStore(&failedWhenSaved),
		unspent:   newFakeUnspentPool(&failedWhenSaved),
		chainMeta: newFakeChainMeta(),
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

func (bt *fakeBlockTree) AddBlock(tx *dbutil.Tx, b *coin.Block) error {
	if bt.saveFailed {
		if bt.failedWhenSaved != nil {
			*bt.failedWhenSaved = true
		}
		return errors.New("intentionally failed")
	}
	bt.blocks[b.HashHeader().Hex()] = b
	return nil
}

func (bt *fakeBlockTree) GetBlock(tx *dbutil.Tx, hash cipher.SHA256) (*coin.Block, error) {
	if bt.failedWhenSaved != nil && *bt.failedWhenSaved {
		return nil, nil
	}
	return bt.blocks[hash.Hex()], nil
}

func (bt *fakeBlockTree) GetBlockInDepth(tx *dbutil.Tx, depth uint64, filter Walker) (*coin.Block, error) {
	if bt.failedWhenSaved != nil && *bt.failedWhenSaved {
		return nil, nil
	}

	for _, b := range bt.blocks {
		if b.Head.BkSeq == depth {
			return b, nil
		}
	}

	return nil, nil
}

func (bt *fakeBlockTree) ForEachBlock(tx *dbutil.Tx, f func(*coin.Block) error) error {
	return nil
}

type fakeSignatureStore struct {
	sigs       map[string]cipher.Sig
	saveFailed bool
	getSigErr  error

	// state tracking: do not configure directly
	// set to true if saveFailed was true and certain operations were performed
	failedWhenSaved *bool
}

func newFakeSigStore(failedWhenSaved *bool) *fakeSignatureStore {
	return &fakeSignatureStore{
		sigs:            make(map[string]cipher.Sig),
		failedWhenSaved: failedWhenSaved,
	}
}

func (ss *fakeSignatureStore) Add(tx *dbutil.Tx, hash cipher.SHA256, sig cipher.Sig) error {
	if ss.saveFailed {
		if ss.failedWhenSaved != nil {
			*ss.failedWhenSaved = true
		}
		return errors.New("intentionally failed")
	}

	ss.sigs[hash.Hex()] = sig
	return nil
}

func (ss *fakeSignatureStore) Get(tx *dbutil.Tx, hash cipher.SHA256) (cipher.Sig, bool, error) {
	if ss.failedWhenSaved != nil && *ss.failedWhenSaved {
		return cipher.Sig{}, false, nil
	}

	if ss.getSigErr != nil {
		return cipher.Sig{}, false, ss.getSigErr
	}

	sig, ok := ss.sigs[hash.Hex()]
	return sig, ok, nil
}

func (ss *fakeSignatureStore) ForEach(tx *dbutil.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
	return nil
}

type fakeUnspentPool struct {
	outs       map[cipher.SHA256]coin.UxOut
	uxHash     cipher.SHA256
	saveFailed bool

	// state tracking: do not configure directly
	// set to true if saveFailed was true and certain operations were performed
	failedWhenSaved *bool
}

func newFakeUnspentPool(failedWhenSaved *bool) *fakeUnspentPool {
	return &fakeUnspentPool{
		outs:            make(map[cipher.SHA256]coin.UxOut),
		failedWhenSaved: failedWhenSaved,
	}
}

func (fup *fakeUnspentPool) MaybeBuildIndexes(tx *dbutil.Tx, height uint64) error {
	return nil
}

func (fup *fakeUnspentPool) Len(tx *dbutil.Tx) (uint64, error) {
	return uint64(len(fup.outs)), nil
}

func (fup *fakeUnspentPool) Get(tx *dbutil.Tx, h cipher.SHA256) (*coin.UxOut, error) {
	out, ok := fup.outs[h]
	if !ok {
		return nil, nil
	}
	return &out, nil
}

func (fup *fakeUnspentPool) GetAll(tx *dbutil.Tx) (coin.UxArray, error) {
	outs := make(coin.UxArray, 0, len(fup.outs))
	for _, out := range fup.outs {
		outs = append(outs, out)
	}

	return outs, nil
}

func (fup *fakeUnspentPool) GetArray(tx *dbutil.Tx, hashes []cipher.SHA256) (coin.UxArray, error) {
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

func (fup *fakeUnspentPool) GetUxHash(tx *dbutil.Tx) (cipher.SHA256, error) {
	return fup.uxHash, nil
}

func (fup *fakeUnspentPool) GetUnspentHashesOfAddrs(tx *dbutil.Tx, addrs []cipher.Address) (AddressHashes, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
	}

	addrOutMap := make(AddressHashes)
	for _, out := range fup.outs {
		addr := out.Body.Address
		addrOutMap[addr] = append(addrOutMap[addr], out.Hash())
	}

	return addrOutMap, nil
}

func (fup *fakeUnspentPool) GetUnspentsOfAddrs(tx *dbutil.Tx, addrs []cipher.Address) (coin.AddressUxOuts, error) {
	addrm := make(map[cipher.Address]struct{}, len(addrs))
	for _, a := range addrs {
		addrm[a] = struct{}{}
	}

	addrOutMap := make(coin.AddressUxOuts)
	for _, out := range fup.outs {
		addr := out.Body.Address
		addrOutMap[addr] = append(addrOutMap[addr], out)
	}

	return addrOutMap, nil
}

func (fup *fakeUnspentPool) ProcessBlock(tx *dbutil.Tx, b *coin.SignedBlock) error {
	if fup.saveFailed {
		if fup.failedWhenSaved != nil {
			*fup.failedWhenSaved = true
		}
		return errors.New("intentionally failed")
	}
	return nil
}

func (fup *fakeUnspentPool) Contains(tx *dbutil.Tx, h cipher.SHA256) (bool, error) {
	_, ok := fup.outs[h]
	return ok, nil
}

func (fup *fakeUnspentPool) AddressCount(tx *dbutil.Tx) (uint64, error) {
	addrs := make(map[cipher.Address]struct{})
	for _, out := range fup.outs {
		addrs[out.Body.Address] = struct{}{}
	}

	return uint64(len(addrs)), nil
}

type fakeChainMeta struct {
	headSeq   uint64
	didSetSeq bool
}

func newFakeChainMeta() *fakeChainMeta {
	return &fakeChainMeta{}
}

func (fcm *fakeChainMeta) GetHeadSeq(tx *dbutil.Tx) (uint64, bool, error) {
	if !fcm.didSetSeq {
		return 0, false, nil
	}

	return fcm.headSeq, true, nil
}

func (fcm *fakeChainMeta) SetHeadSeq(tx *dbutil.Tx, seq uint64) error {
	fcm.headSeq = seq
	fcm.didSetSeq = true
	return nil
}

func DefaultWalker(tx *dbutil.Tx, hps []coin.HashPair) (cipher.SHA256, bool) {
	return hps[0].Hash, true
}

func makeGenesisBlock(t *testing.T) coin.SignedBlock {
	gb, err := coin.NewGenesisBlock(genAddress, genCoinHours, genTime)
	require.NoError(t, err)

	sig := cipher.MustSignHash(gb.HashHeader(), genSecret)
	return coin.SignedBlock{
		Block: *gb,
		Sig:   sig,
	}
}

func TestBlockchainAddBlock(t *testing.T) {
	type expect struct {
		err        error
		sigSaved   bool
		blockSaved bool
		headSeq    uint64
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
				uint64(0),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := prepareDB(t)
			defer closeDB()

			tc.fakeStorage.tree.saveFailed = tc.failedSaves.tree
			tc.fakeStorage.sigs.saveFailed = tc.failedSaves.sigs
			tc.fakeStorage.unspent.saveFailed = tc.failedSaves.unspent

			bc := &Blockchain{
				db:      db,
				unspent: tc.fakeStorage.unspent,
				meta:    tc.fakeStorage.chainMeta,
				tree:    tc.fakeStorage.tree,
				sigs:    tc.fakeStorage.sigs,
				walker:  DefaultWalker,
			}

			gb := makeGenesisBlock(t)

			err := db.Update("", func(tx *dbutil.Tx) error {
				err := bc.AddBlock(tx, &gb)
				require.Equal(t, tc.expect.err, err)
				return nil
			})
			require.NoError(t, err)

			// check sig
			err = db.View("", func(tx *dbutil.Tx) error {
				_, ok, err := tc.fakeStorage.sigs.Get(tx, gb.HashHeader())
				require.NoError(t, err)
				require.Equal(t, tc.expect.sigSaved, ok)

				// check block in tree
				b, err := tc.fakeStorage.tree.GetBlock(tx, gb.HashHeader())
				require.NoError(t, err)
				require.Equal(t, tc.expect.blockSaved, b != nil)

				// check head seq
				headSeq, ok, err := bc.HeadSeq(tx)
				require.NoError(t, err)

				if tc.expect.err == nil {
					require.True(t, ok)
					require.Equal(t, tc.expect.headSeq, headSeq)
				} else {
					require.False(t, ok)
				}

				// check len
				length, err := bc.Len(tx)
				require.NoError(t, err)

				if tc.expect.err == nil {
					require.Equal(t, uint64(1), length)
				} else {
					require.Equal(t, uint64(0), length)
				}

				// check genesis block
				genesisBlock, err := bc.GetGenesisBlock(tx)
				require.NoError(t, err)

				if tc.expect.err == nil {
					require.NotNil(t, genesisBlock)
					require.Equal(t, gb, *genesisBlock)
				} else {
					require.Nil(t, genesisBlock)
				}

				return nil
			})
			require.NoError(t, err)
		})
	}

}

func TestBlockchainHead(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
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
	db, closeDB := prepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := bc.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)
		return nil
	})
	require.NoError(t, err)

	gb := makeGenesisBlock(t)
	err = db.Update("", func(tx *dbutil.Tx) error {
		err := bc.AddBlock(tx, &gb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := bc.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), length)
		return nil
	})
	require.NoError(t, err)
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
				NewErrMissingSignature(&gb.Block),
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
			db, closeDB := prepareDB(t)
			defer closeDB()

			bc, err := NewBlockchain(db, DefaultWalker)
			require.NoError(t, err)

			bc.tree = tc.tree
			bc.sigs = tc.sigs

			err = db.View("", func(tx *dbutil.Tx) error {
				b, err := bc.GetSignedBlockByHash(tx, tc.hash)
				require.Equal(t, tc.expect.err, err)
				require.Equal(t, tc.expect.b, b)
				return nil
			})
			require.NoError(t, err)
		})
	}
}
