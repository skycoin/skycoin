package visor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024
)

var genTime uint64 = 1000
var incTime uint64 = 3600 * 1000
var genCoins uint64 = 1000e6
var genCoinHours uint64 = 1000 * 1000

var failedWhenSave bool

func tNow() uint64 {
	return uint64(utc.UnixNow())
}

func feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func addGenesisBlock(t *testing.T, bc *Blockchain) *coin.SignedBlock {
	// create genesis block
	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	gbSig := cipher.SignHash(gb.HashHeader(), genSecret)

	// add genesis block to blockchain
	bc.db.Update(func(tx *bolt.Tx) error {
		return bc.store.AddBlock(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   gbSig,
		})
	})
	return &coin.SignedBlock{
		Block: *gb,
		Sig:   gbSig,
	}
}

func makeSpendTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction {
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	hours := totalHours / 4

	spendTx.PushOutput(toAddr, coins, hours)
	spendTx.PushOutput(uxs[0].Body.Address, totalCoins-coins, totalHours/4)
	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()
	return spendTx
}

func makeLostCoinTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction { // nolint: unparam
	tx := coin.Transaction{}
	var (
		totalCoins uint64
		totalHours uint64
	)

	for _, ux := range uxs {
		tx.PushInput(ux.Hash())
		totalCoins += ux.Body.Coins
		totalHours += ux.Body.Hours
	}

	tx.PushOutput(toAddr, coins, totalHours/4)
	changeCoins := totalCoins - coins
	if changeCoins > 0 {
		tx.PushOutput(uxs[0].Body.Address, changeCoins-1, totalHours/4)
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return tx
}

/* Helpers */
type fakeChainStore struct {
	len    uint64
	blocks []coin.SignedBlock
	up     blockdb.UnspentPool
}

func (fcs *fakeChainStore) Head(tx *bolt.Tx) (*coin.SignedBlock, error) {
	l := len(fcs.blocks)
	if l == 0 {
		return nil, blockdb.ErrNoHeadBlock
	}

	return &fcs.blocks[l-1], nil
}

func (fcs *fakeChainStore) HeadSeq() uint64 {
	h, err := fcs.Head(nil)
	if err != nil {
		return 0
	}
	return h.Seq()
}

func (fcs *fakeChainStore) Len() uint64 {
	return uint64(len(fcs.blocks))
}

func (fcs *fakeChainStore) AddBlock(tx *bolt.Tx, b *coin.SignedBlock) error {
	return nil
}

func (fcs *fakeChainStore) GetBlockSignature(tx *bolt.Tx, b *coin.Block) (cipher.Sig, bool, error) {
	return cipher.Sig{}, false, nil
}

func (fcs *fakeChainStore) GetBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.Block, error) {
	return nil, nil
}

func (fcs *fakeChainStore) GetSignedBlockByHash(tx *bolt.Tx, hash cipher.SHA256) (*coin.SignedBlock, error) {
	return nil, nil
}

func (fcs *fakeChainStore) GetSignedBlockBySeq(tx *bolt.Tx, seq uint64) (*coin.SignedBlock, error) {
	l := len(fcs.blocks)
	if seq >= uint64(l) {
		return nil, nil
	}

	return &fcs.blocks[seq], nil
}

func (fcs *fakeChainStore) UnspentPool() blockdb.UnspentPool {
	return nil
}

func (fcs *fakeChainStore) GetGenesisBlock() *coin.SignedBlock {
	if len(fcs.blocks) > 0 {
		return &fcs.blocks[0]
	}
	return nil
}

func (fcs *fakeChainStore) ForEachBlock(tx *bolt.Tx, f func(*coin.Block) error) error {
	return nil
}

func (fcs *fakeChainStore) ForEachSignature(tx *bolt.Tx, f func(cipher.SHA256, cipher.Sig) error) error {
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View(func(tx *bolt.Tx) error {
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

			require.Equal(t, tc.isGenesis, bc.isGenesisBlock(*tc.b))
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			bc := &Blockchain{
				db:    db,
				store: tc.store,
			}

			err := db.View(func(tx *bolt.Tx) error {
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			bs, err := bc.GetBlocks(tc.req.st, tc.req.ed)
			require.NoError(t, err)
			require.Equal(t, len(tc.expect), len(bs))
			require.Equal(t, tc.expect, bs)
		})
	}
}

func TestGetLastBlocks(t *testing.T) {
	blocks := makeBlocks(t, 5)
	tt := []struct {
		name   string
		store  chainStore
		n      uint64
		expcet []coin.SignedBlock
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			bc := Blockchain{
				db:    db,
				store: tc.store,
			}

			bs, err := bc.GetLastBlocks(tc.n)
			require.NoError(t, err)
			require.Equal(t, tc.expcet, bs)
		})
	}

}

// newBlock calls bc.NewBlock in a bolt.Tx
func newBlock(t *testing.T, bc *Blockchain, txn coin.Transaction, timestamp uint64) *coin.Block {
	var b *coin.Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		var err error
		b, err = bc.NewBlock(tx, coin.Transactions{txn}, timestamp)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	return b
}

// blockchainHead calls bc.Head in a bolt.Tx
func blockchainHead(t *testing.T, bc *Blockchain) *coin.SignedBlock {
	var head *coin.SignedBlock
	err := bc.db.View(func(tx *bolt.Tx) error {
		var err error
		head, err = bc.Head(tx)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	return head
}

func TestVerifyTransaction(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlock(t, bc)

	toAddr := testutil.MakeAddress()
	coins := uint64(10e6)

	head := blockchainHead(t, bc)

	// create normal spending tx
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	tx := makeSpendTx(uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	err = bc.VerifyTransaction(head, tx)
	require.NoError(t, err)

	originInnerHash := tx.InnerHash
	// test invalid header hash
	tx.InnerHash = cipher.SHA256{}
	err = bc.VerifyTransaction(head, tx)
	require.Equal(t, errors.New("Invalid header hash"), err)

	// set back the originInnerHash
	tx.InnerHash = originInnerHash

	// create new block to spend the coins
	b := newBlock(t, bc, tx, genTime+100)

	// add the block to blockchain
	err = bc.db.Update(func(tx *bolt.Tx) error {
		return bc.store.AddBlock(tx, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
	})
	require.NoError(t, err)

	// none exist ux, the ux already spent
	err = bc.VerifyTransaction(head, tx)
	er := fmt.Errorf("unspent output of %s does not exist", tx.In[0].Hex())
	require.Equal(t, er, err)

	// check invalid sig
	uxs = coin.CreateUnspents(b.Head, tx)
	_, key := cipher.GenerateKeyPair()
	toAddr2 := testutil.MakeAddress()
	tx2 := makeSpendTx(uxs, []cipher.SecKey{key, key}, toAddr2, 5e6)
	err = bc.VerifyTransaction(head, tx2)
	require.Equal(t, errors.New("Signature not valid for output being spent"), err)

	// create lost coin transaction
	uxs2 := coin.CreateUnspents(b.Head, tx)
	toAddr3 := testutil.MakeAddress()
	lostCoinTx := makeLostCoinTx(coin.UxArray{uxs2[1]}, []cipher.SecKey{genSecret}, toAddr3, 10e5)
	err = bc.VerifyTransaction(head, lostCoinTx)
	require.Equal(t, errors.New("Transactions may not create or destroy coins"), err)
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
			errors.New("Signature not valid for output being spent"),
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
			"arbitratint no transactions",
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
			db, closeDB := testutil.PrepareDB(t)
			defer closeDB()

			// create chain store
			store, err := blockdb.NewBlockchain(db, DefaultWalker)
			require.NoError(t, err)

			// create Blockchain
			bc := &Blockchain{
				arbitrating: tc.arbitrating,
				db:          db,
				store:       store,
			}

			// init chain
			head := addGenesisBlock(t, bc)
			tm := head.Time()
			for i, spend := range tc.initChain {
				uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[spend.TxIndex])
				tx := makeSpendTx(coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)

				b := newBlock(t, bc, tx, tm+uint64(i*100))

				sb := &coin.SignedBlock{
					Block: *b,
					Sig:   cipher.SignHash(b.HashHeader(), genSecret),
				}
				db.Update(func(tx *bolt.Tx) error {
					return bc.store.AddBlock(tx, sb)
				})
				head = sb
			}

			// create spending transactions
			txs := make([]coin.Transaction, len(tc.spends))
			for i, spend := range tc.spends {
				uxs := coin.CreateUnspents(head.Head, head.Body.Transactions[spend.TxIndex])
				tx := makeSpendTx(coin.UxArray{uxs[spend.UxIndex]}, spend.Keys, spend.ToAddr, spend.Coins)
				txs[i] = tx
			}

			head = blockchainHead(t, bc)
			_, err = bc.processTransactions(head, txs)
			require.EqualValues(t, tc.err, err)
		})
	}

}

func TestVerifyUxHash(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlock(t, bc)
	uxHash := bc.Unspent().GetUxHash()
	tx := coin.Transaction{}
	b, err := coin.NewBlock(gb.Block, genTime+100, uxHash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)

	err = bc.verifyUxHash(*b)
	require.NoError(t, err)

	b2, err := coin.NewBlock(gb.Block, genTime+10, testutil.RandSHA256(t), coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)

	err = bc.verifyUxHash(*b2)
	require.Equal(t, errors.New("UxHash does not match"), err)
}

func TestProcessBlock(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

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
		Sig:   cipher.SignHash(gb.HashHeader(), genSecret),
	}

	// Test with empty blockchain
	db.Update(func(tx *bolt.Tx) error {
		_, err := bc.processBlock(tx, nil, sb)
		require.NoError(t, err)
		return nil
	})

	// Add genesis block to chain store
	db.Update(func(tx *bolt.Tx) error {
		err := bc.store.AddBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})

	// Create new block
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	toAddr := testutil.MakeAddress()
	tx := makeSpendTx(uxs, []cipher.SecKey{genSecret}, toAddr, 10e6)
	uxhash := bc.Unspent().GetUxHash()
	b, err := coin.NewBlock(*gb, genTime+100, uxhash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)

	head := blockchainHead(t, bc)

	db.Update(func(tx *bolt.Tx) error {
		_, err := bc.processBlock(tx, head, coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
		require.NoError(t, err)
		return nil
	})
}

func TestExecuteBlock(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

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
		Sig:   cipher.SignHash(gb.HashHeader(), genSecret),
	}

	// test with empty chain
	db.Update(func(tx *bolt.Tx) error {
		err := bc.ExecuteBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})

	// new block
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	toAddr := testutil.MakeAddress()
	tx := makeSpendTx(uxs, []cipher.SecKey{genSecret}, toAddr, 10e6)
	uxhash := bc.Unspent().GetUxHash()

	b, err := coin.NewBlock(*gb, genTime+100, uxhash, coin.Transactions{tx}, feeCalc)
	require.NoError(t, err)
	db.Update(func(tx *bolt.Tx) error {
		err := bc.ExecuteBlock(tx, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
		require.NoError(t, err)
		return nil
	})
}
