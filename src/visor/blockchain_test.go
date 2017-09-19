// +build ignore
// These tests need to be rewritten to conform with blockdb changes

package visor

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/stretchr/testify/assert"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024
)

var _genTime uint64 = 1000
var _incTime uint64 = 3600 * 1000
var _genCoins uint64 = 1000e6
var _genCoinHours uint64 = 1000 * 1000

func tNow() uint64 {
	return uint64(utc.UnixNow())
}

func _feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func _makeFeeCalc(fee uint64) coin.FeeCalculator {
	return func(t *coin.Transaction) (uint64, error) {
		return fee, nil
	}
}

/* Helpers */

type FakeTree struct {
	blocks []*coin.Block
}

func (ft *FakeTree) AddBlock(b *coin.Block) error {
	ft.blocks = append(ft.blocks, b)
	return nil
}

func (ft *FakeTree) RemoveBlock(b *coin.Block) error {
	return nil
}

func (ft *FakeTree) GetBlock(hash cipher.SHA256) *coin.Block {
	for _, b := range ft.blocks {
		if b.HashHeader() == hash {
			return b
		}
	}
	return nil
}

func (ft *FakeTree) GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block {
	if dep >= uint64(len(ft.blocks)) {
		return nil
	}
	return ft.blocks[int(dep)]
}

func assertError(t *testing.T, err error, msg string) {
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), msg)
}

// func makeNewBlock(uxHash cipher.SHA256) (*coin.Block, error) {
// 	body := coin.BlockBody{
// 		Transactions: coin.Transactions{coin.Transaction{}},
// 	}

// 	prev := coin.Block{
// 		Body: body,
// 		Head: coin.BlockHeader{
// 			Version:  0x02,
// 			Time:     100,
// 			BkSeq:    0,
// 			Fee:      10,
// 			PrevHash: cipher.SHA256{},
// 			BodyHash: body.Hash(),
// 		}}
// 	return coin.NewBlock(prev, 100+20, uxHash, coin.Transactions{coin.Transaction{}}, _feeCalc)
// }

func makeTransactionForChainWithHoursFee(t *testing.T, bc *Blockchain,
	ux coin.UxOut, sec cipher.SecKey, hours, fee uint64) (coin.Transaction, cipher.SecKey) {
	chrs := ux.CoinHours(bc.Time())
	if chrs < hours+fee {
		log.Panicf("CoinHours underflow. Have %d, need at least %d", chrs,
			hours+fee)
	}
	assert.Equal(t, cipher.AddressFromPubKey(cipher.PubKeyFromSecKey(sec)), ux.Body.Address)
	knownUx, exists := bc.Unspent().Get(ux.Hash())
	assert.True(t, exists)
	assert.Equal(t, knownUx, ux)
	tx := coin.Transaction{}
	tx.PushInput(ux.Hash())
	p, newSec := cipher.GenerateKeyPair()
	addr := cipher.AddressFromPubKey(p)
	tx.PushOutput(addr, 1e6, hours)
	coinsOut := ux.Body.Coins - 1e6
	if coinsOut > 0 {
		tx.PushOutput(genAddress, coinsOut, chrs-hours-fee)
	}
	tx.SignInputs([]cipher.SecKey{sec})
	assert.Equal(t, len(tx.Sigs), 1)
	assert.Nil(t, cipher.ChkSig(ux.Body.Address, cipher.AddSHA256(tx.HashInner(), tx.In[0]), tx.Sigs[0]))
	tx.UpdateHeader()
	assert.Nil(t, tx.Verify())
	err := bc.VerifyTransaction(tx)
	assert.Nil(t, err)
	return tx, newSec
}

func makeTransactionForChainWithFee(t *testing.T, bc *Blockchain,
	fee uint64) coin.Transaction {
	ux := coin.UxOut{}
	hrs := uint64(100)
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	for _, u := range uxs {
		if u.CoinHours(bc.Time()) > hrs {
			ux = u
			break
		}
	}
	assert.Equal(t, ux.Body.Address, genAddress)
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, hrs,
		fee)
	return tx
}

func makeTransactionForChain(t *testing.T, bc *Blockchain) coin.Transaction {
	return makeTransactionForChainWithFee(t, bc, 500000001)
}

// addBlockToBlockchain test helper function
// Adds 2 blocks to the blockchain and return an unspent that has >0 coin hours
func addBlockToBlockchain(t *testing.T, bc *Blockchain) (coin.Block, coin.UxOut) {
	// Split the genesis block into two transactions
	unspents, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	assert.Equal(t, len(unspents), 1)
	ux := unspents[0]
	assert.Equal(t, ux.Body.Address, genAddress)
	pub := cipher.PubKeyFromSecKey(genSecret)
	assert.Equal(t, genAddress, cipher.AddressFromPubKey(pub))
	sig := cipher.SignHash(ux.Hash(), genSecret)
	assert.Nil(t, cipher.ChkSig(ux.Body.Address, ux.Hash(), sig))

	tx, sec := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 0, 0)
	b, err := bc.NewBlockFromTransactions(coin.Transactions{tx}, _incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, *b, tx)
	unspents, err = bc.Unspent().GetAll()
	require.NoError(t, err)
	assert.Equal(t, len(unspents), 2)

	// Spend one of them
	// The other will have hours now
	ux = coin.UxOut{}
	for _, u := range unspents {
		if u.Body.Address != genAddress {
			ux = u
			break
		}
	}
	assert.NotEqual(t, ux.Body.Address, cipher.Address{})
	assert.NotEqual(t, ux.Body.Address, genAddress)
	pub = cipher.PubKeyFromSecKey(sec)
	addr := cipher.AddressFromPubKey(pub)
	assert.Equal(t, ux.Body.Address, addr)
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, sec, 0, 0)
	b, err = bc.NewBlockFromTransactions(coin.Transactions{tx},
		bc.Time()+_incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, *b, tx)

	unspents, err = bc.Unspent().GetAll()
	require.NoError(t, err)
	assert.Equal(t, unspents, 2)

	// Check that the output in the 2nd block is owned by genesis,
	// and has coin hours
	for _, u := range unspents {
		if u.Body.Address == genAddress {
			ux = u
			break
		}
	}
	assert.Equal(t, ux.Body.Address, genAddress)
	assert.Equal(t, ux.Head.BkSeq, uint64(1))
	assert.True(t, ux.CoinHours(bc.Time()) > 0)

	return *b, ux
}

func splitUnspent(t *testing.T, bc *Blockchain, ux coin.UxOut) coin.UxArray {
	tx := coin.Transaction{}
	hrs := ux.CoinHours(bc.Time())
	if hrs < 2 {
		log.Panic("Not enough hours, would generate duplicate output")
	}
	assert.Equal(t, ux.Body.Address, genAddress)
	tx.PushInput(ux.Hash())
	coinsA := ux.Body.Coins / 2
	coinsB := coinsA
	if (ux.Body.Coins/1e6)%2 == 1 {
		coinsA = (ux.Body.Coins - 1e6) / 2
		coinsB = coinsA + 1e6
	}
	tx.PushOutput(genAddress, coinsA, hrs/4)
	tx.PushOutput(genAddress, coinsB, hrs/2)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	b, err := bc.NewBlockFromTransactions(coin.Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)

	var uxs coin.UxArray
	sig := cipher.SignHash(b.HashHeader(), genSecret)
	sb := coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}
	err = executeBlock(bc, &sb)
	assert.Nil(t, err)

	uxs, err = bc.Unspent().GetAll()
	require.NoError(t, err)

	assert.Equal(t, len(uxs), 2)
	return uxs
}

func executeBlock(bc *Blockchain, sb *coin.SignedBlock) error {
	return bc.db.Update(func(tx *bolt.Tx) error {
		return bc.ExecuteBlockWithTx(tx, sb)
	})
}

func makeMultipleOutputs(t *testing.T, bc *Blockchain) {
	txn := coin.Transaction{}
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	ux := uxs[0]
	txn.PushInput(ux.Hash())
	txn.PushOutput(genAddress, 1e6, 100)
	txn.PushOutput(genAddress, 2e6, 100)
	txn.PushOutput(genAddress, _genCoins-3e6, 100)
	txn.SignInputs([]cipher.SecKey{genSecret})
	txn.UpdateHeader()
	assert.Nil(t, txn.Verify())
	assert.Nil(t, bc.VerifyTransaction(txn))
	b, err := bc.NewBlockFromTransactions(coin.Transactions{txn},
		bc.Time()+_incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, *b, txn)
}

func assertExecuteBlock(t *testing.T, bc *Blockchain, b coin.Block,
	tx coin.Transaction) {
	seq := bc.HeadSeq()
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	nUxs := len(uxs)
	sb := coin.SignedBlock{
		Block: b,
		Sig:   cipher.SignHash(b.HashHeader(), genSecret),
	}
	err = executeBlock(bc, &sb)
	assert.Nil(t, err)

	assert.Equal(t, bc.HeadSeq(), seq+1)
	assert.Equal(t, len(uxs), len(tx.Out))
	assert.False(t, uxs.HasDupes())
	uxs, err = bc.Unspent().GetAll()
	require.NoError(t, err)
	assert.Equal(t, len(uxs), nUxs+len(tx.Out)-len(tx.In))
	for _, ux := range uxs {
		ux2, exist := bc.Unspent().Get(ux.Hash())
		assert.True(t, exist)
		assert.Equal(t, ux, ux2)
	}

	head, err := bc.Head()
	require.NoError(t, err)
	uxs2 := coin.CreateUnspents(head.Head, tx)
	assert.Equal(t, len(uxs2), len(uxs))
	for i, u := range uxs2 {
		assert.Equal(t, u.Body, uxs[i].Body)
	}
}

func assertValidUnspents(t *testing.T, bh coin.BlockHeader, tx coin.Transaction,
	uxo coin.UxArray) {
	assert.Equal(t, len(tx.Out), len(uxo))
	for i, ux := range uxo {
		assert.Equal(t, bh.Time, ux.Head.Time)
		assert.Equal(t, bh.BkSeq, ux.Head.BkSeq)
		assert.Equal(t, tx.Hash(), ux.Body.SrcTransaction)
		assert.Equal(t, tx.Out[i].Address, ux.Body.Address)
		assert.Equal(t, tx.Out[i].Coins, ux.Body.Coins)
		assert.Equal(t, tx.Out[i].Hours, ux.Body.Hours)
	}
}

// func TestNewBlockchain(t *testing.T) {
// 	db, cleanDB := testutil.PrepareDB(t)
// 	defer cleanDB()

// 	p, _ := cipher.GenerateKeyPair()

// 	bc, err := NewBlockchain(db, p, Arbitrating(true))
// 	require.NoError(t, err)

// 	// assert.Equal(t, len(b.GetUnspent().Pool), 0)
// }

// func TestCreateMasterGenesisBlock(t *testing.T) {
// 	db, cleanDB := testutil.PrepareDB(t)
// 	defer cleanDB()

// 	p, _ := cipher.GenerateKeyPair()
// 	bc, err := NewBlockchain(db, p, Arbitrating(true))
// 	require.NoError(t, err)

// 	genAddress := makeAddress()
// 	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
// 	require.NotNil(t, err)

// 	// assert.Equal(t, len(b.Blocks), 1)
// 	// assert.Equal(t, b.Blocks[0], gb)
// 	uxs, err := bc.Unspent().GetAll()
// 	require.NoError(t, err)
// 	assert.Equal(t, len(uxs), 1)
// 	assert.Equal(t, uxs[0].Body.Address, genAddress)
// 	assert.NotEqual(t, gb.Head.Time, uint64(0))
// 	assert.Equal(t, gb.Head.BkSeq, uint64(0))
// 	// Panicing
// 	assert.Panics(t, func() { bc.CreateGenesisBlock(genAddress, _genCoins, _genTime) })
// }

func TestCreateGenesisBlock(t *testing.T) {
	now := tNow()

	db, cleanDB := testutil.PrepareDB(t)
	defer cleanDB()

	p, _ := cipher.GenerateKeyPair()
	bc, err := NewBlockchain(db, p, Arbitrating(true))
	require.NoError(t, err)

	bc.CreateGenesisBlock(genAddress, _genCoins, now)
	assert.Equal(t, bc.Time(), now)
	assert.Equal(t, bc.Head().Head.BkSeq, uint64(0))
	assert.Equal(t, len(bc.Head().Body.Transactions), 1)
	assert.Equal(t, len(bc.Head().Body.Transactions[0].Out), 1)
	assert.Equal(t, len(bc.Head().Body.Transactions[0].In), 0)
	txn := bc.Head().Body.Transactions[0]
	txo := txn.Out[0]
	assert.Equal(t, txo.Address, genAddress)
	assert.Equal(t, txo.Coins, _genCoins)
	assert.Equal(t, txo.Hours, _genCoins)
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	assert.Equal(t, len(uxs), 1)
	ux := uxs[0]
	assert.Equal(t, ux.Head.BkSeq, uint64(0))
	assert.Equal(t, ux.Head.Time, now)
	assert.Equal(t, ux.Head.Time, bc.Head().Head.Time)
	assert.Equal(t, ux.Body.SrcTransaction, txn.InnerHash)
	assert.Equal(t, ux.Body.Address, genAddress)
	assert.Equal(t, ux.Body.Coins, _genCoins)
	assert.Equal(t, txo.Coins, ux.Body.Coins)
	assert.Equal(t, txo.Hours, ux.Body.Hours)
	// 1 hour per coin, at init
	assert.Equal(t, ux.Body.Hours, _genCoins)
	h := cipher.Merkle([]cipher.SHA256{bc.Head().Body.Transactions[0].Hash()})
	assert.Equal(t, bc.Head().Head.BodyHash, h)
	assert.Equal(t, bc.Head().Head.PrevHash, cipher.SHA256{})
	// TODO -- check valid snapshot
	assert.NotEqual(t, bc.Head().Head.UxHash, [4]byte{})
	expect := coin.CreateUnspents(bc.Head().Head, txn)
	expect.Sort()
	uxs, err = bc.Unspent().GetAll()
	require.NoError(t, err)
	uxs.Sort()
	assert.Equal(t, expect, uxs)
	// Panicing
	_, err = bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.Error(t, err)
}

func TestBlockchainHead(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	b, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := b.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	assert.Equal(t, b.Head(), &gb)
	nb, _ := addBlockToBlockchain(t, b)
	assert.Equal(t, b.Head(), &nb)
}

func TestBlockchainTime(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	b, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := b.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	assert.Equal(t, b.Time(), gb.Head.Time)
	nb, _ := addBlockToBlockchain(t, b)
	assert.Equal(t, b.Time(), nb.Head.Time)
}

func TestNewBlockFromTransactions(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	// gb.Head.Version = 0x0F
	// bc.Blocks[0] = gb
	// assert.Equal(t, bc.GetGenesisBlock().Head.Version, uint32(0x0F))

	assert.Equal(t, bc.Len(), uint64(1))
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, bc.Len(), uint64(3))

	// No transactions
	_, err = bc.NewBlockFromTransactions(coin.Transactions{},
		bc.Time()+_incTime)
	assertError(t, err, "No transactions")
	assert.Equal(t, bc.Len(), uint64(3))

	// Bad currentTime, must be greater than head time
	fee := uint64(100)
	txn, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		fee)
	txns := coin.Transactions{txn}
	assert.Panics(t, func() {
		bc.NewBlockFromTransactions(txns, bc.Time())
	})

	// Valid transaction
	hrs := ux.CoinHours(bc.Time())
	seq := bc.Head().Head.BkSeq
	b, err := bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 1)
	assert.Equal(t, b.Body.Transactions[0], txn)
	assert.Equal(t, b.Head.BkSeq, seq+1)
	assert.Equal(t, b.Head.Time, bc.Time()+_incTime)
	assert.Equal(t, b.Head.Version, gb.Head.Version)
	assert.Equal(t, b.Head.Fee, fee)
	assert.Equal(t, b.Head.Fee, hrs-txn.OutputHours())
	assert.NotEqual(t, b.Head.Fee, uint64(0))

	// Invalid transaction
	txn.InnerHash = cipher.SHA256{}
	txns = coin.Transactions{txn}
	_, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assertError(t, err, "Invalid header hash")

	// Multiple transactions, sorted
	// First, split our genesis block into two, so we can make 2 valid txns
	uxs := splitUnspent(t, bc, ux)

	// tNow, make two valid txns
	txn = coin.Transaction{}
	txn.PushInput(uxs[0].Hash())
	txn.PushOutput(genAddress, uxs[0].Body.Coins, uxs[0].Body.Hours)
	txn.SignInputs([]cipher.SecKey{genSecret})
	txn.UpdateHeader()
	txn2 := coin.Transaction{}
	txn2.PushInput(uxs[1].Hash())
	txn2.PushOutput(genAddress, uxs[1].Body.Coins, uxs[1].Body.Hours)
	txn2.SignInputs([]cipher.SecKey{genSecret})
	txn2.UpdateHeader()

	// Combine them and sort
	txns = coin.Transactions{txn, txn2}
	txns = coin.SortTransactions(txns, bc.TransactionFee)
	b, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 2)
	assert.Equal(t, b.Body.Transactions, txns)

	// Order should be preserved
	txns2 := coin.Transactions{txn, txn2}
	sTxns := coin.NewSortableTransactions(txns2, bc.TransactionFee)
	if sTxns.IsSorted() {
		txns2[0], txns2[1] = txns2[1], txns2[0]
	}
	b, err = bc.NewBlockFromTransactions(txns2, bc.Time()+_incTime)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 2)
	assert.Equal(t, b.Body.Transactions, txns2)
}

func TestCreateUnspents(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	// 1 out
	tx := coin.Transaction{}
	tx.PushOutput(genAddress, 11e6, 255)
	bh := coin.BlockHeader{
		Time:  tNow(),
		BkSeq: uint64(1),
	}
	uxout := coin.CreateUnspents(bh, tx)
	assert.Equal(t, len(uxout), 1)
	assertValidUnspents(t, bh, tx, uxout)

	// Multiple outs.  Should work regardless of validity
	tx = coin.Transaction{}
	ux := testutil.MakeUxOut(t)
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, 100, 150)
	tx.PushOutput(genAddress, 200, 77)
	bh.BkSeq++
	uxout = coin.CreateUnspents(bh, tx)
	assert.Equal(t, len(uxout), 2)
	assertValidUnspents(t, bh, tx, uxout)

	// No outs
	tx = coin.Transaction{}
	uxout = coin.CreateUnspents(bh, tx)
	assertValidUnspents(t, bh, tx, uxout)
}

func TestVerifyTransactionSpending(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)

	// Overspending hours

	tx := coin.Transaction{}
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	tx.PushInput(uxs[0].Hash())
	tx.PushOutput(genAddress, 1e6, uxs[0].Body.Hours)
	tx.PushOutput(genAddress, uxs[0].Body.Coins-1e6, 1)
	uxIn, err := bc.Unspent().GetArray(tx.In)
	assert.Nil(t, err)
	uxOut := coin.CreateUnspents(bc.Head().Head, tx)
	assertError(t, coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut),
		"Insufficient coin hours")

	// add block to blockchain.
	_, ux := addBlockToBlockchain(t, bc)
	// addBlockToBlockchain(t, bc)

	// Valid
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	uxIn, err = bc.Unspent().GetArray(tx.In)
	assert.Nil(t, err)
	uxOut = coin.CreateUnspents(bc.Head().Head, tx)
	assert.Nil(t, coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut))

	// Destroying coins
	tx = coin.Transaction{}
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, 1e6, 100)
	tx.PushOutput(genAddress, 10e6, 100)
	uxIn, err = bc.Unspent().GetArray(tx.In)
	assert.Nil(t, err)
	uxOut = coin.CreateUnspents(bc.Head().Head, tx)
	err = coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(),
		"Transactions may not create or destroy coins")
	assertError(t, coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut),
		"Transactions may not create or destroy coins")

	// Insufficient coins
	tx = coin.Transaction{}
	tx.PushInput(ux.Hash())
	p, s := cipher.GenerateKeyPair()
	a := cipher.AddressFromPubKey(p)
	coins := ux.Body.Coins
	assert.True(t, coins > 1e6)
	tx.PushOutput(a, 1e6, 100)
	tx.PushOutput(genAddress, coins-1e6, 100)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	b, err := bc.NewBlockFromTransactions(coin.Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)
	sig := cipher.SignHash(b.HashHeader(), genSecret)
	sb := coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}
	err = executeBlock(bc, &sb)
	assert.Nil(t, err)
	tx = coin.Transaction{}
	tx.PushInput(uxs[0].Hash())
	tx.PushOutput(a, 10e6, 1)
	tx.SignInputs([]cipher.SecKey{s})
	tx.UpdateHeader()
	uxIn, err = bc.Unspent().GetArray(tx.In)
	assert.Nil(t, err)
	uxOut = coin.CreateUnspents(bc.Head().Head, tx)
	assertError(t, coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut),
		"Insufficient coins")
}

func TestVerifyTransaction(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	// Genesis block is not valid by normal standards
	assert.NotNil(t, bc.VerifyTransaction(gb.Body.Transactions[0]))
	assert.Equal(t, bc.Len(), uint64(1))
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, bc.Len(), uint64(3))

	// Valid txn
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	assert.Nil(t, bc.VerifyTransaction(tx))
	assert.Equal(t, bc.Len(), uint64(3))

	// Failure, spending unknown output
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.In[0] = cipher.SHA256{}
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Unspent output does not exist")
	assert.Equal(t, bc.Len(), uint64(3))

	// Failure, duplicate input
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.In = append(tx.In, tx.In[0])
	tx.SignInputs([]cipher.SecKey{genSecret, genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Duplicate spend")
	assert.Equal(t, bc.Len(), uint64(3))

	// Failure, zero coin output
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.PushOutput(genAddress, 0, 100)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Zero coin output")

	// // Failure, hash collision with unspents
	// tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	// uxOut := coin.CreateUnspents(bc.Head().Head, tx)
	// bc.Unspent().Add(uxOut[0])
	// assertError(t, bc.VerifyTransaction(tx),
	// 	"New unspent collides with existing unspent")

	// Failure, not spending enough coins
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.PushOutput(genAddress, 10e6, 100)
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Insufficient coins")

	// Failure, spending outputs we don't own
	_, s := cipher.GenerateKeyPair()
	tx = coin.Transaction{}
	uxs, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	for _, u := range uxs {
		if u.Body.Address != genAddress {
			ux = u
			break
		}
	}
	assert.NotEqual(t, ux.Body.Address, genAddress)
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, ux.Body.Coins, ux.Body.Hours)
	tx.SignInputs([]cipher.SecKey{s})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx),
		"Signature not valid for output being spent")

	// Failure, wrong signature for txn hash
	tx = coin.Transaction{}
	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.PushOutput(genAddress, ux.Body.Coins, ux.Body.Hours)
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx),
		"Signature not valid for output being spent")
}

func TestBlockchainProcessBlock(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	// Genesis block not valid after the fact
	// assert.NotNil(t, bc.verifyBlock(gb))
	assert.Equal(t, bc.Len(), uint64(1))
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, bc.Len(), uint64(3))

	// Valid block
	tx := coin.Transaction{}
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, ux.Body.Coins, ux.CoinHours(bc.Time()))
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	b, err := bc.NewBlockFromTransactions(coin.Transactions{tx}, bc.Time()+_incTime)
	assert.Equal(t, len(b.Body.Transactions), 1)
	assert.Equal(t, len(b.Body.Transactions[0].Out), 1)
	assert.Nil(t, err)
	// assert.Nil(t, bc.verifyBlock(b))

	// Invalid block header
	b.Head.BkSeq = gb.Head.BkSeq
	assert.Equal(t, len(b.Body.Transactions), 1)
	assert.Equal(t, len(b.Body.Transactions[0].Out), 1)
	_, err = bc.processBlock(*b)
	testutil.RequireError(t, err, "BkSeq invalid")

	// Invalid transactions, makes duplicate outputs
	b.Head.BkSeq = bc.Head().Head.BkSeq + 1
	b.Body.Transactions = append(b.Body.Transactions, b.Body.Transactions[0])
	b.Head.BodyHash = b.HashBody()
	_, err = bc.processBlock(*b)
	testutil.RequireError(t, err, "Duplicate unspent output across transactions")
}

func TestVerifyUxHash(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	_, err = bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	require.NoError(t, err)

	b := coin.Block{Body: coin.BlockBody{}, Head: coin.BlockHeader{}}
	b.Body.Transactions = append(b.Body.Transactions, testutil.MakeTransaction(t))
	uxHash := bc.Unspent().GetUxHash()
	copy(b.Head.UxHash[:], uxHash[:])
	assert.Nil(t, bc.verifyUxHash(b))
	b.Head.UxHash = cipher.SHA256{}
	testutil.RequireError(t, bc.verifyUxHash(b), "UxHash does not match")
}

func TestVerifyBlockHeader(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	b := coin.Block{Body: coin.BlockBody{}}
	b.Body.Transactions = append(b.Body.Transactions, testutil.MakeTransaction(t))
	h := coin.BlockHeader{}

	h.BkSeq = 1
	h.Time = gb.Head.Time + 1
	h.PrevHash = gb.HashHeader()
	h.BodyHash = b.HashBody()

	// Valid header
	b.Head = h
	assert.Nil(t, bc.verifyBlockHeader(b))

	// Invalid bkSeq
	i := h
	i.BkSeq++
	b.Head = i
	assertError(t, bc.verifyBlockHeader(b), "BkSeq invalid")

	// Invalid time
	i = h
	i.Time = gb.Head.Time
	b.Head = i
	assertError(t, bc.verifyBlockHeader(b),
		"Block time must be > head time")
	b.Head.Time--
	assertError(t, bc.verifyBlockHeader(b),
		"Block time must be > head time")

	// Invalid prevHash
	i = h
	i.PrevHash = cipher.SHA256{}
	b.Head = i
	assertError(t, bc.verifyBlockHeader(b),
		"PrevHash does not match current head")

	// Invalid bodyHash
	i = h
	i.BodyHash = cipher.SHA256{}
	b.Head = i
	assertError(t, bc.verifyBlockHeader(b),
		"Computed body hash does not match")
}

func TestTransactionFee(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()
	bc, err := NewBlockchain(db, genPublic)
	require.NoError(t, err)
	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	assert.Equal(t, bc.Len(), uint64(1))
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, bc.Len(), uint64(3))

	// Valid txn, 100 hours fee
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	fee, err := bc.TransactionFee(&tx)
	assert.Nil(t, err)
	assert.Equal(t, fee, uint64(100))

	// Txn spending unknown output
	tx = coin.Transaction{}
	unknownUx := testutil.MakeUxOut(t)
	tx.PushInput(unknownUx.Hash())
	_, err = bc.TransactionFee(&tx)
	assertError(t, err, "Unspent output does not exist")

	// Txn spending more hours than avail
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 100)
	tx.PushOutput(testutil.MakeAddress(), 1e6, 10000)
	_, err = bc.TransactionFee(&tx)
	assertError(t, err, "Insufficient coinhours for transaction outputs")
}

// func TestProcessTransactions(t *testing.T) {
// 	db, closeDB := testutil.PrepareDB(t)
// 	defer closeDB()
// 	bc, err := NewBlockchain(db, genPublic)
// 	require.NoError(t, err)
// 	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
// 	fmt.Println("genesis time:", bc.GetGenesisBlock().Time())
// 	assert.Equal(t, bc.Len(), uint64(1))
// 	_, ux := addBlockToBlockchain(t, bc)
// 	assert.Equal(t, bc.Len(), uint64(3))

// 	// Invalid, no transactions in block
// 	// arbitrating=false
// 	bc.arbitrating = false
// 	txns, err := bc.processTransactions(coin.Transactions{})
// 	assert.Nil(t, txns)
// 	assertError(t, err, "No transactions")
// 	// arbitrating=true
// 	bc.arbitrating = true
// 	txns, err = bc.processTransactions(coin.Transactions{})
// 	assert.Equal(t, len(txns), 0)
// 	assert.Nil(t, err)

// 	// Invalid, txn.Verify() fails
// 	// TODO -- combine all txn.Verify() failures into one test
// 	// method, and call it from here, from ExecuteBlock(), from
// 	// Verify(), from VerifyTransaction()
// 	txns = coin.Transactions{}
// 	txn := coin.Transaction{}
// 	txn.PushInput(ux.Hash())
// 	txn.PushOutput(genAddress, 777, 100)
// 	txn.SignInputs([]cipher.SecKey{genSecret})
// 	txn.UpdateHeader()
// 	txns = append(txns, txn)
// 	// arbitrating=false
// 	bc.arbitrating = false
// 	txns2, err := bc.processTransactions(txns)
// 	assert.Nil(t, txns2)
// 	assertError(t, err,
// 		"Transaction outputs must be multiple of 1e6 base units")
// 	// arbitrating=true
// 	bc.arbitrating = true
// 	txns2, err = bc.processTransactions(txns)
// 	assert.NotNil(t, txns2)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(txns2), 0)

// 	// Invalid, duplicate unspent will be created by these txns
// 	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
// 		100)
// 	txns = coin.Transactions{txn, txn}
// 	// arbitrating=false
// 	bc.arbitrating = false
// 	txns2, err = bc.processTransactions(txns)
// 	assertError(t, err, "Duplicate unspent output across transactions")
// 	assert.Nil(t, txns2)
// 	// arbitrating=true.  One of the offending transactions should be removed
// 	bc.arbitrating = true
// 	txns2, err = bc.processTransactions(txns)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(txns2), 1)
// 	assert.Equal(t, txns2[0], txn)

// 	// Check that a new output will not collide with the existing pool
// 	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
// 		100)
// 	txns = coin.Transactions{txn}
// 	uxb := coin.UxBody{
// 		SrcTransaction: txn.Hash(),
// 		Coins:          txn.Out[0].Coins,
// 		Hours:          txn.Out[0].Hours,
// 		Address:        txn.Out[0].Address,
// 	}
// 	bc.GetUnspent().Add(coin.UxOut{Body: uxb})
// 	// arbitrating=false
// 	txns2, err = bc.processTransactions(txns, false)
// 	assertError(t, err, "New unspent collides with existing unspent")
// 	assert.Nil(t, txns2)
// 	// arbitrating=true
// 	txns2, err = bc.processTransactions(txns, true)
// 	assert.Equal(t, len(txns2), 0)
// 	assert.NotNil(t, txns2)
// 	assert.Nil(t, err)

// 	// Spending of duplicate inputs being spent across txns
// 	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
// 		100)
// 	txn2, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
// 		100)
// 	txn2.Out = nil
// 	txn2.PushOutput(makeAddress(), 1e6, 100)
// 	txn2.PushOutput(makeAddress(), ux.Body.Coins-1e6, 100)
// 	txn2.Sigs = nil
// 	txn2.SignInputs([]cipher.SecKey{genSecret})
// 	txn2.UpdateHeader()
// 	txns = coin.SortTransactions(coin.Transactions{txn, txn2}, bc.TransactionFee)
// 	// arbitrating=false
// 	txns2, err = bc.processTransactions(txns, false)
// 	assertError(t, err, "Cannot spend output twice in the same block")
// 	assert.Nil(t, txns2)
// 	// arbitrating=true
// 	txns2, err = bc.processTransactions(txns, true)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(txns2), 1)
// 	assert.Equal(t, txns2[0], txns[0])
// }

// func TestExecuteBlock(t *testing.T) {
// 	ft := FakeTree{}
// 	bc := NewBlockchain(&ft, nil)
// 	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
// 	assert.Equal(t, bc.Len(), uint64(1))
// 	_, ux := addBlockToBlockchain(t, bc)
// 	assert.Equal(t, bc.Len(), uint64(3))

// 	// Invalid block returns error
// 	b := coin.Block{}
// 	uxs, err := bc.ExecuteBlock(&b)
// 	assert.NotNil(t, err)
// 	assert.Nil(t, uxs)

// 	// Valid block, spends are removed from the unspent pool, new ones are
// 	// added.  Blocks is updated, and new unspents are returns
// 	assert.Equal(t, bc.Len(), uint64(3))
// 	assert.Equal(t, len(bc.GetUnspent().Pool), 2)
// 	spuxs := splitUnspent(t, bc, ux)
// 	tx := coin.Transaction{}
// 	tx.PushInput(spuxs[0].Hash())
// 	coins := spuxs[0].Body.Coins
// 	extra := coins % 4e6
// 	coins = (coins - extra) / 4
// 	tx.PushOutput(genAddress, coins+extra, spuxs[0].Body.Hours/5)
// 	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/6)
// 	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/7)
// 	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/8)
// 	tx.SignInputs([]cipher.SecKey{genSecret})
// 	tx.UpdateHeader()
// 	tx2 := coin.Transaction{}
// 	tx2.PushInput(spuxs[1].Hash())
// 	tx2.PushOutput(genAddress, spuxs[1].Body.Coins, spuxs[1].Body.Hours/10)
// 	tx2.SignInputs([]cipher.SecKey{genSecret})
// 	tx2.UpdateHeader()
// 	txns := coin.Transactions{tx, tx2}
// 	sTxns := coin.NewSortableTransactions(txns, bc.TransactionFee)
// 	unswapped := sTxns.IsSorted()
// 	txns = coin.SortTransactions(txns, bc.TransactionFee)
// 	txns, err = bc.verifyTransactions(txns)
// 	assert.Nil(t, err)
// 	seq := bc.Head().Head.BkSeq
// 	b, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
// 	assert.Equal(t, b.Head.BkSeq, seq+1)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(b.Body.Transactions), 2)
// 	assert.Equal(t, b.Body.Transactions, txns)
// 	uxs, err = bc.ExecuteBlock(&b)
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(uxs), 5)
// 	// Check that all unspents look correct and are in the unspent pool
// 	txOuts := []coin.TransactionOutput{}
// 	if unswapped {
// 		txOuts = append(txOuts, tx.Out...)
// 		txOuts = append(txOuts, tx2.Out...)
// 	} else {
// 		txOuts = append(txOuts, tx2.Out...)
// 		txOuts = append(txOuts, tx.Out...)
// 	}
// 	for i, ux := range uxs {
// 		if unswapped {
// 			if i < len(tx.Out) {
// 				assert.Equal(t, ux.Body.SrcTransaction, tx.Hash())
// 			} else {
// 				assert.Equal(t, ux.Body.SrcTransaction, tx2.Hash())
// 			}
// 		} else {
// 			if i < len(tx2.Out) {
// 				assert.Equal(t, ux.Body.SrcTransaction, tx2.Hash())
// 			} else {
// 				assert.Equal(t, ux.Body.SrcTransaction, tx.Hash())
// 			}
// 		}
// 		assert.Equal(t, ux.Body.Address, txOuts[i].Address)
// 		assert.Equal(t, ux.Body.Coins, txOuts[i].Coins)
// 		assert.Equal(t, ux.Body.Hours, txOuts[i].Hours)
// 		assert.Equal(t, ux.Head.BkSeq, b.Head.BkSeq)
// 		assert.Equal(t, ux.Head.Time, b.Head.Time)
// 		assert.True(t, bc.GetUnspent().Has(ux.Hash()))
// 	}
// 	// Check that all spends are no longer in the pool
// 	txIns := []cipher.SHA256{}
// 	txIns = append(txIns, tx.In...)
// 	txIns = append(txIns, tx2.In...)
// 	for _, ux := range txIns {
// 		assert.False(t, bc.GetUnspent().Has(ux))
// 	}
// }

func TestMakeTx(t *testing.T) {
	db, closedb := testutil.PrepareDB(t)
	defer closedb()

	bc, err := NewBlockchain(db, genPublic, Arbitrating(true))
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)

	toAddr := testutil.MakeAddress()
	tx, err := testutil.MakeTx(bc.Unspent().GetUnspentsOfAddr, []cipher.SecKey{genSecret}, toAddr, 10e6)
	require.NoError(t, err)
	tt := Transaction{
		Txn: *tx,
	}
	rtx := NewReadableTransaction(&tt)
	v, err := json.MarshalIndent(rtx, "", "    ")
	require.NoError(t, err)
	fmt.Println(string(v))
	uxHash := bc.Unspent().GetUxHash()
	fmt.Println(bc.Time())
	b, err := coin.NewBlock(gb, bc.Time()+_incTime, uxHash, coin.Transactions{*tx}, _feeCalc)
	require.NoError(t, err)
	rb := NewReadableBlock(b)
	bv, err := json.MarshalIndent(rb, "", "    ")
	require.NoError(t, err)
	fmt.Println(string(bv))

	sb := &coin.SignedBlock{
		Block: *b,
		Sig:   cipher.SignHash(b.HashHeader(), genSecret),
	}
	err = executeBlock(bc, sb)
	require.NoError(t, err)
}

type Spending struct {
	FromSeckeys []cipher.SecKey
	ToAddr      cipher.Address
	Coins       uint64
}

func MakeSpendingChain(t *testing.T, db *bolt.DB, spds []Spending) (*Blockchain, error) {
	bc, err := NewBlockchain(db, genPublic, Arbitrating(true))
	require.NoError(t, err)
	gb, err := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	require.NoError(t, err)
	preBlock := gb
	currentTime := bc.Time() + _incTime
	for _, spending := range spds {
		tx, err := testutil.MakeTx(bc.Unspent().GetUnspentsOfAddr, spending.FromSeckeys, spending.ToAddr, spending.Coins)
		require.NoError(t, err)
		uxHash := bc.Unspent().GetUxHash()
		b, err := coin.NewBlock(preBlock, currentTime, uxHash, coin.Transactions{*tx}, _feeCalc)

		require.NoError(t, err)
		err = executeBlock(bc, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
		preBlock = *b
		currentTime += _incTime
		require.NoError(t, err)
	}

	return bc, nil
}

func TestSpendingChain(t *testing.T) {
	db, closedb := testutil.PrepareDB(t)
	defer closedb()

	toAddrs := []cipher.Address{}
	keys := []cipher.SecKey{}
	for i := 0; i < 3; i++ {
		_, s := cipher.GenerateKeyPair()
		keys = append(keys, s)
		toAddrs = append(toAddrs, cipher.AddressFromSecKey(s))
	}

	spendChan := []Spending{
		Spending{
			FromSeckeys: []cipher.SecKey{genSecret},
			ToAddr:      toAddrs[0],
			Coins:       10e6,
		},
		Spending{
			FromSeckeys: []cipher.SecKey{keys[0]},
			ToAddr:      toAddrs[1],
			Coins:       1e6,
		},
	}

	bc, err := MakeSpendingChain(t, db, spendChan)
	require.NoError(t, err)
	require.Equal(t, uint64(3), bc.Len())
}
