package visor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
)

var genTime uint64 = 1000
var genCoins uint64 = 1000e6

func feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func addGenesisBlockToBlockchain(t *testing.T, bc *Blockchain) *coin.SignedBlock {
	// create genesis block
	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	gbSig := cipher.MustSignHash(gb.HashHeader(), genSecret)

	// add genesis block to blockchain
	err = bc.db.Update("", func(tx *dbutil.Tx) error {
		return bc.store.AddBlock(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   gbSig,
		})
	})
	require.NoError(t, err)

	return &coin.SignedBlock{
		Block: *gb,
		Sig:   gbSig,
	}
}

func makeSpendTxn(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction {
	spendTxn := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		err := spendTxn.PushInput(ux.Hash())
		require.NoError(t, err)
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	require.True(t, coins <= totalCoins)

	hours := totalHours / 4

	err := spendTxn.PushOutput(toAddr, coins, hours)
	require.NoError(t, err)
	if totalCoins-coins != 0 {
		err := spendTxn.PushOutput(uxs[0].Body.Address, totalCoins-coins, totalHours/4)
		require.NoError(t, err)
	}
	spendTxn.SignInputs(keys)
	err = spendTxn.UpdateHeader()
	require.NoError(t, err)
	return spendTxn
}

/* Helpers */
type fakeChainStore struct {
	blocks []coin.SignedBlock
}

func (fcs *fakeChainStore) Head(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	l := len(fcs.blocks)
	if l == 0 {
		return nil, blockdb.ErrNoHeadBlock
	}

	return &fcs.blocks[l-1], nil
}

func (fcs *fakeChainStore) HeadSeq(tx *dbutil.Tx) (uint64, bool, error) {
	h, err := fcs.Head(tx)
	if err != nil {
		if err == blockdb.ErrNoHeadBlock {
			return 0, false, nil
		}
		return 0, false, err
	}
	return h.Seq(), true, nil
}

func (fcs *fakeChainStore) Len(tx *dbutil.Tx) (uint64, error) {
	return uint64(len(fcs.blocks)), nil
}

func (fcs *fakeChainStore) AddBlock(tx *dbutil.Tx, b *coin.SignedBlock) error {
	return nil
}

func (fcs *fakeChainStore) GetBlockSignature(tx *dbutil.Tx, b *coin.Block) (cipher.Sig, bool, error) {
	return cipher.Sig{}, false, nil
}

func (fcs *fakeChainStore) GetBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.Block, error) {
	return nil, nil
}

func (fcs *fakeChainStore) GetSignedBlockByHash(tx *dbutil.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
	return nil, nil
}

func (fcs *fakeChainStore) GetSignedBlockBySeq(tx *dbutil.Tx, seq uint64) (*coin.SignedBlock, error) {
	l := len(fcs.blocks)
	if seq >= uint64(l) {
		return nil, nil
	}

	return &fcs.blocks[seq], nil
}

func (fcs *fakeChainStore) UnspentPool() blockdb.UnspentPooler {
	return nil
}

func (fcs *fakeChainStore) GetGenesisBlock(tx *dbutil.Tx) (*coin.SignedBlock, error) {
	if len(fcs.blocks) > 0 {
		return &fcs.blocks[0], nil
	}
	return nil, nil
}

func (fcs *fakeChainStore) ForEachBlock(tx *dbutil.Tx, f func(*coin.Block) error) error {
	return nil
}

func makeBlock(t *testing.T, preBlock coin.Block, tm uint64) *coin.Block {
	uxHash := testutil.RandSHA256(t)
	tx := coin.Transaction{}
	b, err := coin.NewBlock(preBlock, tm, uxHash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)
	return b
}

func makeBlocks(t *testing.T, n int) []coin.SignedBlock {
	var bs []coin.SignedBlock
	preBlock, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	bs = append(bs, coin.SignedBlock{Block: *preBlock})

	now := genTime + 100
	for i := 1; i < n; i++ {
		b := makeBlock(t, *preBlock, now+uint64(i)*100)
		sb := coin.SignedBlock{
			Block: *b,
		}
		bs = append(bs, sb)
		preBlock = b
	}

	return bs
}

func TestBlockchainTime(t *testing.T) {
	bs := makeBlocks(t, 1)
	tt := []struct {
		name  string
		store chainStore
		time  uint64
	}{
		{
			"ok",
			&fakeChainStore{
				blocks: bs[:],
			},
			bs[0].Time(),
		},
		{
			"no head",
			&fakeChainStore{},
			uint64(0),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := prepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View("", func(tx *dbutil.Tx) error {
				tm, err := bc.Time(tx)
				require.NoError(t, err)
				require.Equal(t, tc.time, tm)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestIsGenesisBlock(t *testing.T) {
	bs := makeBlocks(t, 2)
	tt := []struct {
		name      string
		store     chainStore
		b         *coin.Block
		isGenesis bool
	}{
		{
			"genesis block",
			&fakeChainStore{
				blocks: bs[:1],
			},
			&bs[0].Block,
			true,
		},
		{
			"not genesis block",
			&fakeChainStore{
				blocks: bs[:1],
			},
			&bs[1].Block,
			false,
		},
		{
			"empty chain",
			&fakeChainStore{},
			&bs[0].Block,
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			bc := Blockchain{
				store: tc.store,
			}

			isGenesis, err := bc.isGenesisBlock(nil, *tc.b)
			require.NoError(t, err)
			require.Equal(t, tc.isGenesis, isGenesis)
		})
	}
}

func TestVerifyBlockHeader(t *testing.T) {
	bs := makeBlocks(t, 5)
	tt := []struct {
		name  string
		store chainStore
		b     coin.Block
		err   error
	}{
		{
			"ok",
			&fakeChainStore{
				blocks: bs[:1],
			},
			bs[1].Block,
			nil,
		},
		{
			"invalid block seq",
			&fakeChainStore{
				blocks: bs[:1],
			},
			bs[2].Block,
			errors.New("BkSeq invalid"),
		},
		{
			"invalid time",
			&fakeChainStore{
				blocks: bs[:1],
			},
			coin.Block{
				Head: coin.BlockHeader{
					BkSeq: 1,
					Time:  0,
				},
			},

			errors.New("Block time must be > head time"),
		},
		{
			"invalid prehash",
			&fakeChainStore{
				blocks: bs[:1],
			},
			coin.Block{
				Head: coin.BlockHeader{
					BkSeq: 1,
					Time:  bs[1].Time(),
				},
			},

			errors.New("PrevHash does not match current head"),
		},
		{
			"empty blockchain",
			&fakeChainStore{},
			coin.Block{},
			blockdb.ErrNoHeadBlock,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := prepareDB(t)
			defer closeDB()

			bc := &Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View("", func(tx *dbutil.Tx) error {
				err := bc.verifyBlockHeader(tx, tc.b)
				require.Equal(t, tc.err, err)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestGetBlocks(t *testing.T) {
	blocks := makeBlocks(t, 5)
	tt := []struct {
		name  string
		store chainStore
		req   struct {
			st uint64
			ed uint64
		}
		expect []coin.SignedBlock
	}{
		{
			"ok",
			&fakeChainStore{
				blocks: blocks[:],
			},
			struct {
				st uint64
				ed uint64
			}{
				0,
				1,
			},
			blocks[:2],
		},
		{
			"start > end",
			&fakeChainStore{
				blocks: blocks[:],
			},
			struct {
				st uint64
				ed uint64
			}{
				1,
				0,
			},
			nil,
		},
		{
			"start overflow",
			&fakeChainStore{
				blocks: blocks[:],
			},
			struct {
				st uint64
				ed uint64
			}{
				6,
				7,
			},
			nil,
		},
		{
			"start == end",
			&fakeChainStore{
				blocks: blocks[:],
			},
			struct {
				st uint64
				ed uint64
			}{
				0,
				0,
			},
			blocks[:1],
		},
		{
			"end overflow",
			&fakeChainStore{
				blocks: blocks[:],
			},
			struct {
				st uint64
				ed uint64
			}{
				0,
				8,
			},
			blocks[:],
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := prepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View("", func(tx *dbutil.Tx) error {
				bs, err := bc.GetBlocksInRange(tx, tc.req.st, tc.req.ed)
				require.NoError(t, err)
				require.Equal(t, len(tc.expect), len(bs))
				require.Equal(t, tc.expect, bs)
				return nil
			})
			require.NoError(t, err)
		})
	}
}

func TestGetLastBlocks(t *testing.T) {
	blocks := makeBlocks(t, 5)
	tt := []struct {
		name   string
		store  chainStore
		n      uint64
		expect []coin.SignedBlock
	}{
		{
			"get last block",
			&fakeChainStore{
				blocks: blocks[:],
			},
			1,
			blocks[4:5],
		},
		{
			"get last two block",
			&fakeChainStore{
				blocks: blocks[:],
			},
			2,
			blocks[3:5],
		},
		{
			"get all block",
			&fakeChainStore{
				blocks: blocks[:],
			},
			5,
			blocks[0:5],
		},
		{
			"get block from empty chain",
			&fakeChainStore{},
			1,
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, closeDB := prepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View("", func(tx *dbutil.Tx) error {
				bs, err := bc.GetLastBlocks(tx, tc.n)
				require.NoError(t, err)
				require.Equal(t, tc.expect, bs)
				return nil
			})
			require.NoError(t, err)
		})
	}

}

// newBlock calls bc.NewBlock in a dbutil.Tx
func newBlock(t *testing.T, bc *Blockchain, txn coin.Transaction, timestamp uint64) *coin.Block {
	var b *coin.Block
	err := bc.db.View("", func(tx *dbutil.Tx) error {
		var err error
		b, err = bc.NewBlock(tx, coin.Transactions{txn}, timestamp)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	return b
}

type spending struct {
	TxIndex int
	UxIndex int
	Keys    []cipher.SecKey
	ToAddr  cipher.Address
	Coins   uint64
}

func TestProcessTransactions(t *testing.T) {
	toAddrs := make([]cipher.Address, 10)
	keys := make([]cipher.SecKey, 10)
	for i := 0; i < 10; i++ {
		p, s := cipher.GenerateKeyPair()
		toAddrs[i] = cipher.AddressFromPubKey(p)
		keys[i] = s
	}

	tt := []struct {
		name        string
		arbitrating bool
		initChain   []spending
		spends      []spending
		err         error
	}{
		{
			"ok",
			false,
			[]spending{},
			[]spending{
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{genSecret},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
			},
			nil,
		},
		{
			"no transactions",
			false,
			[]spending{},
			[]spending{},
			errors.New("No transactions"),
		},
		{
			"invalid signature",
			false,
			[]spending{},
			[]spending{
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{keys[0]},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
			},
			NewErrTxnViolatesHardConstraint(errors.New("Signature not valid for output being spent")),
		},
		{
			"dup spending",
			false,
			[]spending{},
			[]spending{
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{genSecret},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{genSecret},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
			},
			errors.New("Cannot spend output twice in the same block"),
		},
		{
			"arbitrating no transactions",
			true,
			[]spending{},
			[]spending{},
			nil,
		},
		{
			"invalid signature",
			true,
			[]spending{},
			[]spending{
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{keys[0]},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
			},
			nil,
		},
		{
			"including invalid signature",
			true,
			[]spending{},
			[]spending{
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{genSecret},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
				spending{
					TxIndex: 0,
					UxIndex: 0,
					Keys:    []cipher.SecKey{keys[0]},
					ToAddr:  toAddrs[0],
					Coins:   10e6,
				},
			},
			nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// create test db
			db, closeDB := prepareDB(t)
			defer closeDB()

			err := CreateBuckets(db)
			require.NoError(t, err)

			// create chain store
			store, err := blockdb.NewBlockchain(db, DefaultWalker)
			require.NoError(t, err)

			// create Blockchain
			bc := &Blockchain{
				cfg: BlockchainConfig{
					Arbitrating: tc.arbitrating,
				},
				db:    db,
				store: store,
			}

			// init chain
			head := addGenesisBlockToBlockchain(t, bc)
			tm := head.Time()
			for i, spend := range tc.initChain {
				uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[spend.TxIndex])
				txn := makeSpendTxn(t, coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)

				b := newBlock(t, bc, txn, tm+uint64(i*100))

				sb := &coin.SignedBlock{
					Block: *b,
					Sig:   cipher.MustSignHash(b.HashHeader(), genSecret),
				}
				err = db.Update("", func(tx *dbutil.Tx) error {
					return bc.store.AddBlock(tx, sb)
				})
				require.NoError(t, err)
				head = sb
			}

			// create spending transactions
			txns := make([]coin.Transaction, len(tc.spends))
			for i, spend := range tc.spends {
				uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[spend.TxIndex])
				txn := makeSpendTxn(t, coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)
				txns[i] = txn
			}

			err = db.View("", func(tx *dbutil.Tx) error {
				_, err := bc.processTransactions(tx, txns)
				require.EqualValues(t, tc.err, err)
				return nil
			})
			require.NoError(t, err)
		})
	}

}

func getUxHash(t *testing.T, db *dbutil.DB, bc *Blockchain) cipher.SHA256 {
	var uxHash cipher.SHA256
	err := db.View("", func(tx *dbutil.Tx) error {
		var err error
		uxHash, err = bc.Unspent().GetUxHash(tx)
		return err
	})
	require.NoError(t, err)
	return uxHash
}

func TestVerifyUxHash(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	err := CreateBuckets(db)
	require.NoError(t, err)

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlockToBlockchain(t, bc)
	uxHash := getUxHash(t, db, bc)
	txn := coin.Transaction{}
	b, err := coin.NewBlock(gb.Block, genTime+100, uxHash, coin.Transactions{txn}, feeCalc)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		err = bc.verifyUxHash(tx, *b)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	b2, err := coin.NewBlock(gb.Block, genTime+10, testutil.RandSHA256(t), coin.Transactions{txn}, feeCalc)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		err = bc.verifyUxHash(tx, *b2)
		require.Equal(t, errors.New("UxHash does not match"), err)
		return nil
	})
	require.NoError(t, err)
}

func TestProcessBlock(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	err := CreateBuckets(db)
	require.NoError(t, err)

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)

	sb := coin.SignedBlock{
		Block: *gb,
		Sig:   cipher.MustSignHash(gb.HashHeader(), genSecret),
	}

	// Test with empty blockchain
	err = db.Update("", func(tx *dbutil.Tx) error {
		_, err := bc.processBlock(tx, sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	// Add genesis block to chain store
	err = db.Update("", func(tx *dbutil.Tx) error {
		err := bc.store.AddBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	// Create new block
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	toAddr := testutil.MakeAddress()
	tx := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, toAddr, 10e6)
	uxHash := getUxHash(t, db, bc)
	b, err := coin.NewBlock(*gb, genTime+100, uxHash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)

	err = db.Update("", func(tx *dbutil.Tx) error {
		_, err := bc.processBlock(tx, coin.SignedBlock{
			Block: *b,
			Sig:   cipher.MustSignHash(b.HashHeader(), genSecret),
		})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}

func TestExecuteBlock(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	err := CreateBuckets(db)
	require.NoError(t, err)

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)

	sb := coin.SignedBlock{
		Block: *gb,
		Sig:   cipher.MustSignHash(gb.HashHeader(), genSecret),
	}

	// test with empty chain
	err = db.Update("", func(tx *dbutil.Tx) error {
		err := bc.ExecuteBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	// new block
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	toAddr := testutil.MakeAddress()
	tx := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, toAddr, 10e6)
	uxHash := getUxHash(t, db, bc)
	b, err := coin.NewBlock(*gb, genTime+100, uxHash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)
	err = db.Update("", func(tx *dbutil.Tx) error {
		err := bc.ExecuteBlock(tx, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.MustSignHash(b.HashHeader(), genSecret),
		})
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
}
