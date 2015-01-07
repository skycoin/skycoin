package coin

import (
	"bytes"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
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

/*



*/

func tNow() uint64 {
	return uint64(time.Now().UTC().Unix())
}

func _feeCalc(t *Transaction) (uint64, error) {
	return 0, nil
}

func _makeFeeCalc(fee uint64) FeeCalculator {
	return func(t *Transaction) (uint64, error) {
		return fee, nil
	}
}

/* Helpers */

func assertError(t *testing.T, err error, msg string) {
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), msg)
}

func makeNewBlock() Block {
	unsp := NewUnspentPool()
	body := BlockBody{
		Transactions: Transactions{Transaction{}},
	}
	prev := Block{
		Body: body,
		Head: BlockHeader{
			Version:  0x02,
			Time:     100,
			BkSeq:    0,
			Fee:      10,
			PrevHash: cipher.SHA256{},
			BodyHash: body.Hash(),
		}}
	return newBlock(prev, 100+20, unsp, Transactions{Transaction{}}, _feeCalc)
}

func makeTransactionForChainWithHoursFee(t *testing.T, bc *Blockchain,
	ux UxOut, sec cipher.SecKey, hours, fee uint64) (Transaction, cipher.SecKey) {
	chrs := ux.CoinHours(bc.Time())
	if chrs < hours+fee {
		log.Panicf("CoinHours underflow. Have %d, need at least %d", chrs,
			hours+fee)
	}
	assert.Equal(t, cipher.AddressFromPubKey(cipher.PubKeyFromSecKey(sec)), ux.Body.Address)
	knownUx, exists := bc.Unspent.Get(ux.Hash())
	assert.True(t, exists)
	assert.Equal(t, knownUx, ux)
	tx := Transaction{}
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
	assert.Nil(t, cipher.ChkSig(ux.Body.Address, cipher.AddSHA256(tx.hashInner(), tx.In[0]), tx.Sigs[0]))
	tx.UpdateHeader()
	assert.Nil(t, tx.Verify())
	err := bc.VerifyTransaction(tx)
	assert.Nil(t, err)
	return tx, newSec
}

func makeTransactionForChainWithFee(t *testing.T, bc *Blockchain,
	fee uint64) Transaction {
	ux := UxOut{}
	hrs := uint64(100)
	for _, u := range bc.Unspent.Array() {
		if ux.CoinHours(bc.Time()) > hrs {
			ux = u
			break
		}
	}
	assert.Equal(t, ux.Body.Address, genAddress)
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, hrs,
		fee)
	return tx
}

func makeTransactionForChain(t *testing.T, bc *Blockchain) Transaction {
	return makeTransactionForChainWithFee(t, bc, 100)
}

func addTransactionToBlock(t *testing.T, b *Block) Transaction {
	tx := makeTransaction(t)
	b.Body.Transactions = append(b.Body.Transactions, tx)
	return tx
}

// Adds 2 blocks to the blockchain and return an unspent that has >0 coin hours
func addBlockToBlockchain(t *testing.T, bc *Blockchain) (Block, UxOut) {
	// Split the genesis block into two transactions
	assert.Equal(t, len(bc.Unspent.Array()), 1)
	ux := bc.Unspent.Array()[0]
	assert.Equal(t, ux.Body.Address, genAddress)
	pub := cipher.PubKeyFromSecKey(genSecret)
	assert.Equal(t, genAddress, cipher.AddressFromPubKey(pub))
	sig := cipher.SignHash(ux.Hash(), genSecret)
	assert.Nil(t, cipher.ChkSig(ux.Body.Address, ux.Hash(), sig))

	tx, sec := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 0, 0)
	b, err := bc.NewBlockFromTransactions(Transactions{tx}, _incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, b, tx)
	assert.Equal(t, len(bc.Unspent.Array()), 2)

	// Spend one of them
	// The other will have hours now
	ux = UxOut{}
	for _, u := range bc.Unspent.Pool {
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
	b, err = bc.NewBlockFromTransactions(Transactions{tx},
		bc.Time()+_incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, b, tx)
	assert.Equal(t, len(bc.Unspent.Array()), 2)

	// Check that the output in the 2nd block is owned by genesis,
	// and has coin hours
	for _, u := range bc.Unspent.Pool {
		if u.Body.Address == genAddress {
			ux = u
			break
		}
	}
	assert.Equal(t, ux.Body.Address, genAddress)
	assert.Equal(t, ux.Head.BkSeq, uint64(1))
	assert.True(t, ux.CoinHours(bc.Time()) > 0)

	return b, ux
}

func splitUnspent(t *testing.T, bc *Blockchain, ux UxOut) UxArray {
	tx := Transaction{}
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
	b, err := bc.NewBlockFromTransactions(Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)
	uxs, err := bc.ExecuteBlock(b)
	assert.Nil(t, err)
	assert.Equal(t, len(uxs), 2)
	return uxs
}

func makeMultipleOutputs(t *testing.T, bc *Blockchain) {
	txn := Transaction{}
	ux := bc.Unspent.Array()[0]
	txn.PushInput(ux.Hash())
	txn.PushOutput(genAddress, 1e6, 100)
	txn.PushOutput(genAddress, 2e6, 100)
	txn.PushOutput(genAddress, _genCoins-3e6, 100)
	txn.SignInputs([]cipher.SecKey{genSecret})
	txn.UpdateHeader()
	assert.Nil(t, txn.Verify())
	assert.Nil(t, bc.VerifyTransaction(txn))
	b, err := bc.NewBlockFromTransactions(Transactions{txn},
		bc.Time()+_incTime)
	assert.Nil(t, err)
	assertExecuteBlock(t, bc, b, txn)
}

func assertExecuteBlock(t *testing.T, bc *Blockchain, b Block,
	tx Transaction) {
	seq := bc.Head().Head.BkSeq
	nUxs := len(bc.Unspent.Pool)
	uxs, err := bc.ExecuteBlock(b)
	assert.Nil(t, err)
	assert.Equal(t, bc.Head().Head.BkSeq, seq+1)
	assert.Equal(t, len(uxs), len(tx.Out))
	assert.False(t, uxs.HasDupes())
	assert.Equal(t, len(bc.Unspent.Pool), nUxs+len(tx.Out)-len(tx.In))
	for _, ux := range uxs {
		assert.True(t, bc.Unspent.Has(ux.Hash()))
		ux2, ok := bc.Unspent.Get(ux.Hash())
		assert.True(t, ok)
		assert.Equal(t, ux, ux2)
	}
	uxs2 := CreateUnspents(bc.Head().Head, tx)
	assert.Equal(t, len(uxs2), len(uxs))
	for i, u := range uxs2 {
		assert.Equal(t, u.Body, uxs[i].Body)
	}
}

func assertValidUnspents(t *testing.T, bh BlockHeader, tx Transaction,
	uxo UxArray) {
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

func _badFeeCalc(t *Transaction) (uint64, error) {
	return 0, errors.New("Bad")
}

/* Tests */

func TestNewBlock(t *testing.T) {
	// TODO -- update this test for newBlock changes
	prev := Block{Head: BlockHeader{Version: 0x02, Time: 100, BkSeq: 98}}
	unsp := NewUnspentPool()
	unsp.XorHash = randSHA256(t)
	txns := Transactions{Transaction{}}
	// invalid txn fees panics
	assert.Panics(t, func() { newBlock(prev, 133, unsp, txns, _badFeeCalc) })
	// no txns panics
	assert.Panics(t, func() {
		newBlock(prev, 133, unsp, nil, _feeCalc)
	})
	assert.Panics(t, func() {
		newBlock(prev, 133, unsp, Transactions{}, _feeCalc)
	})
	// valid block is fine
	fee := uint64(121)
	currentTime := uint64(133)
	b := newBlock(prev, currentTime, unsp, txns, _makeFeeCalc(fee))
	assert.Equal(t, b.Body.Transactions, txns)
	assert.Equal(t, b.Head.Fee, fee*uint64(len(txns)))
	assert.Equal(t, b.Body, BlockBody{txns})
	assert.Equal(t, b.Head.PrevHash, prev.HashHeader())
	assert.Equal(t, b.Head.Time, currentTime)
	assert.Equal(t, b.Head.BkSeq, prev.Head.BkSeq+1)
	assert.Equal(t, b.Head.UxHash,
		getSnapshotHash(unsp, prev.HashHeader()))
}

func TestBlockHashHeader(t *testing.T) {
	b := makeNewBlock()
	assert.Equal(t, b.HashHeader(), b.Head.Hash())
	assert.NotEqual(t, b.HashHeader(), cipher.SHA256{})
}

func TestBlockHashBody(t *testing.T) {
	b := makeNewBlock()
	assert.Equal(t, b.HashBody(), b.Body.Hash())
	hb := b.HashBody()
	hashes := b.Body.Transactions.Hashes()
	tx := addTransactionToBlock(t, &b)
	assert.NotEqual(t, b.HashBody(), hb)
	hashes = append(hashes, tx.Hash())
	assert.Equal(t, b.HashBody(), cipher.Merkle(hashes))
	assert.Equal(t, b.HashBody(), b.Body.Hash())
}

func TestBlockString(t *testing.T) {
	b := makeNewBlock()
	assert.Equal(t, b.String(), b.Head.String())
}

func TestBlockGetTransaction(t *testing.T) {
	b := makeNewBlock()
	_, ok := b.GetTransaction(cipher.SHA256{})
	assert.False(t, ok)
	tx := addTransactionToBlock(t, &b)
	tx2, ok := b.GetTransaction(tx.Hash())
	assert.True(t, ok)
	assert.Equal(t, tx, tx2)
	tx3 := addTransactionToBlock(t, &b)
	tx4, ok := b.GetTransaction(tx3.Hash())
	assert.True(t, ok)
	assert.Equal(t, tx3, tx4)
	_, ok = b.GetTransaction(cipher.SHA256{})
	assert.False(t, ok)
}

func TestNewBlockHeaderPanics(t *testing.T) {
	ct := uint64(100)
	head := BlockHeader{Time: ct}
	assert.Panics(t, func() {
		newBlockHeader(head, NewUnspentPool(), ct, 50, BlockBody{})
	})
	assert.Panics(t, func() {
		newBlockHeader(head, NewUnspentPool(), ct-1, 50, BlockBody{})
	})
	assert.Panics(t, func() {
		newBlockHeader(head, NewUnspentPool(), ct-ct, 50, BlockBody{})
	})
	assert.NotPanics(t, func() {
		newBlockHeader(head, NewUnspentPool(), ct+1, 50, BlockBody{})
	})
}

func TestNewBlockHeader(t *testing.T) {
	// TODO -- update this test for newBlockHeader changes
	b := makeNewBlock()
	prev := b.Head
	unsp := NewUnspentPool()
	unsp.XorHash = randSHA256(t)
	fee := uint64(10)
	bh := newBlockHeader(prev, unsp, prev.Time+22, fee, b.Body)
	assert.Equal(t, bh.PrevHash, prev.Hash())
	assert.NotEqual(t, bh.PrevHash, prev.PrevHash)
	assert.Equal(t, bh.Time, uint64(prev.Time+22))
	assert.Equal(t, bh.BkSeq, uint64(prev.BkSeq+1))
	assert.Equal(t, bh.Fee, fee)
	assert.Equal(t, bh.Version, prev.Version)
	assert.Equal(t, bh.BodyHash, b.Body.Hash())
	assert.Equal(t, bh.UxHash, getSnapshotHash(unsp, prev.Hash()))
}

func TestBlockHeaderHash(t *testing.T) {
	b := makeNewBlock()
	h := b.Head.Hash()
	assert.Equal(t, h, b.Head.Hash())
	assert.NotEqual(t, b.Head.Hash(), cipher.SHA256{})
	// Change header should change hash
	b.Head.BkSeq = uint64(5)
	assert.NotEqual(t, h, b.Head.Hash())
}

func TestBlockHeaderBytes(t *testing.T) {
	b := makeNewBlock()
	by := b.Head.Bytes()
	assert.NotNil(t, by)
	assert.NotEqual(t, len(by), 0)
	b.Head.BkSeq += 1
	by2 := b.Head.Bytes()
	assert.False(t, bytes.Equal(by, by2))
}

func TestBlockHeaderString(t *testing.T) {
	b := makeNewBlock()
	assert.NotEqual(t, b.Head.String(), "")
}

func TestBlockBodyHash(t *testing.T) {
	b := makeNewBlock()
	hashes := b.Body.Transactions.Hashes()
	assert.Equal(t, b.Body.Hash(), cipher.Merkle(hashes))
	tx1 := addTransactionToBlock(t, &b)
	hashes = append(hashes, tx1.Hash())
	assert.Equal(t, b.Body.Hash(), cipher.Merkle(hashes))
	tx2 := addTransactionToBlock(t, &b)
	hashes = append(hashes, tx2.Hash())
	assert.Equal(t, b.Body.Hash(), cipher.Merkle(hashes))
	assert.Equal(t, b.HashBody(), b.Body.Hash())
}

func TestBlockBodyBytes(t *testing.T) {
	b := makeNewBlock()
	assert.NotPanics(t, func() { b.Body.Bytes() })
	by := b.Body.Bytes()
	assert.NotNil(t, by)
	addTransactionToBlock(t, &b)
	by2 := b.Body.Bytes()
	assert.True(t, len(by2) > len(by))
}

func TestBlockBodySize(t *testing.T) {
	b := makeNewBlock()
	addTransactionToBlock(t, &b)
	assert.Equal(t, b.Size(), b.Body.Size())
	assert.Equal(t, b.Body.Size(), b.Body.Transactions.Size())
	size := 0
	for _, x := range b.Body.Transactions {
		size += len(encoder.Serialize(&x))
	}
	assert.NotEqual(t, size, 0)
	assert.Equal(t, b.Size(), size)
}

func TestNewBlockchain(t *testing.T) {
	b := NewBlockchain()
	assert.NotNil(t, b.Blocks)
	assert.Equal(t, len(b.Blocks), 0)
	assert.NotNil(t, b.Unspent)
	assert.Equal(t, len(b.Unspent.Pool), 0)
}

func TestCreateMasterGenesisBlock(t *testing.T) {
	b := NewBlockchain()
	a := makeAddress()
	gb := b.CreateGenesisBlock(a, _genTime, _genCoins)

	assert.Equal(t, len(b.Blocks), 1)
	assert.Equal(t, b.Blocks[0], gb)
	assert.Equal(t, len(b.Unspent.Pool), 1)
	assert.Equal(t, b.Unspent.Array()[0].Body.Address, a)
	assert.NotEqual(t, gb.Head.Time, uint64(0))
	assert.Equal(t, gb.Head.BkSeq, uint64(0))
	// Panicing
	assert.Panics(t, func() { b.CreateGenesisBlock(a, _genTime, _genCoins) })
}

func TestCreateGenesisBlock(t *testing.T) {
	b := NewBlockchain()
	now := tNow()
	gb := b.CreateGenesisBlock(genAddress, now, _genCoins)
	assert.Equal(t, gb.Head.Time, now)
	assert.Equal(t, gb.Head.BkSeq, uint64(0))
	assert.Equal(t, len(gb.Body.Transactions), 1)
	assert.Equal(t, len(gb.Body.Transactions[0].Out), 1)
	assert.Equal(t, len(gb.Body.Transactions[0].In), 0)
	txn := gb.Body.Transactions[0]
	txo := txn.Out[0]
	assert.Equal(t, txo.Address, genAddress)
	assert.Equal(t, txo.Coins, _genCoins)
	assert.Equal(t, txo.Hours, _genCoins)
	assert.Equal(t, len(b.Unspent.Pool), 1)
	ux := b.Unspent.Array()[0]
	assert.Equal(t, ux.Head.BkSeq, uint64(0))
	assert.Equal(t, ux.Head.Time, now)
	assert.Equal(t, ux.Head.Time, gb.Head.Time)
	assert.Equal(t, ux.Body.SrcTransaction, txn.Hash())
	assert.Equal(t, ux.Body.Address, genAddress)
	assert.Equal(t, ux.Body.Coins, _genCoins)
	assert.Equal(t, txo.Coins, ux.Body.Coins)
	assert.Equal(t, txo.Hours, ux.Body.Hours)
	// 1 hour per coin, at init
	assert.Equal(t, ux.Body.Hours, _genCoins)
	h := cipher.Merkle([]cipher.SHA256{gb.Body.Transactions[0].Hash()})
	assert.Equal(t, gb.Head.BodyHash, h)
	assert.Equal(t, gb.Head.PrevHash, cipher.SHA256{})
	// TODO -- check valid snapshot
	assert.NotEqual(t, gb.Head.UxHash, [4]byte{})
	expect := CreateUnspents(gb.Head, txn)
	expect.Sort()
	have := b.Unspent.Array()
	have.Sort()
	assert.Equal(t, expect, have)
	// Panicing
	assert.Panics(t, func() {
		b.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	})
}

func TestBlockchainHead(t *testing.T) {
	b := NewBlockchain()
	gb := b.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	assert.Equal(t, b.Head(), gb)
	nb, _ := addBlockToBlockchain(t, b)
	assert.Equal(t, b.Head(), nb)
}

func TestBlockchainTime(t *testing.T) {
	b := NewBlockchain()
	gb := b.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	assert.Equal(t, b.Time(), gb.Head.Time)
	nb, _ := addBlockToBlockchain(t, b)
	assert.Equal(t, b.Time(), nb.Head.Time)
}

func TestNewBlockFromTransactions(t *testing.T) {
	bc := NewBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	gb.Head.Version = 0x0F
	bc.Blocks[0] = gb
	assert.Equal(t, bc.Blocks[0].Head.Version, uint32(0x0F))
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// No transactions
	_, err := bc.NewBlockFromTransactions(Transactions{},
		bc.Time()+_incTime)
	assertError(t, err, "No transactions")
	assert.Equal(t, len(bc.Blocks), 3)

	// Bad currentTime, must be greater than head time
	fee := uint64(100)
	txn, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		fee)
	txns := Transactions{txn}
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
	txn.Head.Hash = cipher.SHA256{}
	txns = Transactions{txn}
	_, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assertError(t, err, "Invalid header hash")

	// Multiple transactions, sorted
	// First, split our genesis block into two, so we can make 2 valid txns
	uxs := splitUnspent(t, bc, ux)

	// tNow, make two valid txns
	txn = Transaction{}
	txn.PushInput(uxs[0].Hash())
	txn.PushOutput(genAddress, uxs[0].Body.Coins, uxs[0].Body.Hours)
	txn.SignInputs([]cipher.SecKey{genSecret})
	txn.UpdateHeader()
	txn2 := Transaction{}
	txn2.PushInput(uxs[1].Hash())
	txn2.PushOutput(genAddress, uxs[1].Body.Coins, uxs[1].Body.Hours)
	txn2.SignInputs([]cipher.SecKey{genSecret})
	txn2.UpdateHeader()

	// Combine them and sort
	txns = Transactions{txn, txn2}
	txns = SortTransactions(txns, bc.TransactionFee)
	b, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 2)
	assert.Equal(t, b.Body.Transactions, txns)

	// Order should be preserved
	txns2 := Transactions{txn, txn2}
	sTxns := newSortableTransactions(txns2, bc.TransactionFee)
	if sTxns.IsSorted() {
		txns2[0], txns2[1] = txns2[1], txns2[0]
	}
	b, err = bc.NewBlockFromTransactions(txns2, bc.Time()+_incTime)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 2)
	assert.Equal(t, b.Body.Transactions, txns2)
}

func TestVerifyTransactionInputs(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	_, ux := addBlockToBlockchain(t, bc)
	// Valid txn
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	uxIn, err := bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	assert.Nil(t, verifyTransactionInputs(tx, uxIn))
	// Bad sigs
	sig := tx.Sigs[0]
	tx.Sigs[0] = cipher.Sig{}
	assert.NotNil(t, verifyTransactionInputs(tx, uxIn))
	// Too many uxIn
	tx.Sigs[0] = sig
	uxIn, err = bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	assert.Equal(t, len(uxIn), len(tx.In))
	uxIn = append(uxIn, makeUxOut(t))
	assert.True(t, DebugLevel2)
	assert.Panics(t, func() { verifyTransactionInputs(tx, uxIn) })
	// ux hash mismatch
	uxIn, err = bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	tx.In[0] = cipher.SHA256{}
	assert.Panics(t, func() { verifyTransactionInputs(tx, uxIn) })
}

func TestCreateUnspents(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	// 1 out
	tx := Transaction{}
	tx.PushOutput(genAddress, 11e6, 255)
	bh := BlockHeader{
		Time:  tNow(),
		BkSeq: uint64(1),
	}
	uxout := CreateUnspents(bh, tx)
	assert.Equal(t, len(uxout), 1)
	assertValidUnspents(t, bh, tx, uxout)

	// Multiple outs.  Should work regardless of validity
	tx = Transaction{}
	ux := makeUxOut(t)
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, 100, 150)
	tx.PushOutput(genAddress, 200, 77)
	bh.BkSeq += 1
	uxout = CreateUnspents(bh, tx)
	assert.Equal(t, len(uxout), 2)
	assertValidUnspents(t, bh, tx, uxout)

	// No outs
	tx = Transaction{}
	uxout = CreateUnspents(bh, tx)
	assertValidUnspents(t, bh, tx, uxout)
}

func TestVerifyTransactionSpending(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	_, ux := addBlockToBlockchain(t, bc)

	// Valid
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	uxIn, err := bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	uxOut := CreateUnspents(bc.Head().Head, tx)
	assert.Nil(t, verifyTransactionSpending(bc.Time(), tx, uxIn, uxOut))

	// Destroying coins
	tx = Transaction{}
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, 1e6, 0)
	tx.PushOutput(genAddress, 10e6, 0)
	uxIn, err = bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	uxOut = CreateUnspents(bc.Head().Head, tx)
	err = verifyTransactionSpending(bc.Time(), tx, uxIn, uxOut)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(),
		"Transactions may not create or destroy coins")
	assertError(t, verifyTransactionSpending(bc.Time(), tx, uxIn, uxOut),
		"Transactions may not create or destroy coins")

	// Overspending hours
	tx = Transaction{}
	tx.PushInput(bc.Unspent.Array()[0].Hash())
	tx.PushOutput(genAddress, 1e6, bc.Unspent.Array()[0].Body.Hours)
	tx.PushOutput(genAddress, bc.Unspent.Array()[0].Body.Coins-1e6, 1)
	uxIn, err = bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	uxOut = CreateUnspents(bc.Head().Head, tx)
	assertError(t, verifyTransactionSpending(bc.Time(), tx, uxIn, uxOut),
		"Insufficient coin hours")

	// Insufficient coins
	tx = Transaction{}
	tx.PushInput(ux.Hash())
	p, s := cipher.GenerateKeyPair()
	a := cipher.AddressFromPubKey(p)
	coins := ux.Body.Coins
	assert.True(t, coins > 1e6)
	tx.PushOutput(a, 1e6, 100)
	tx.PushOutput(genAddress, coins-1e6, 100)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	b, err := bc.NewBlockFromTransactions(Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)
	uxs, err := bc.ExecuteBlock(b)
	assert.Nil(t, err)
	tx = Transaction{}
	tx.PushInput(uxs[0].Hash())
	tx.PushOutput(a, 10e6, 100)
	tx.SignInputs([]cipher.SecKey{s})
	tx.UpdateHeader()
	uxIn, err = bc.Unspent.GetMultiple(tx.In)
	assert.Nil(t, err)
	uxOut = CreateUnspents(bc.Head().Head, tx)
	assertError(t, verifyTransactionSpending(bc.Time(), tx, uxIn, uxOut),
		"Insufficient coins")
}

func TestVerifyTransaction(t *testing.T) {
	bc := NewBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	// Genesis block is not valid by normal standards
	assert.NotNil(t, bc.VerifyTransaction(gb.Body.Transactions[0]))
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// Valid txn
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	assert.Nil(t, bc.VerifyTransaction(tx))
	assert.Equal(t, len(bc.Blocks), 3)

	// Failure, spending unknown output
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.In[0] = cipher.SHA256{}
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Unspent output does not exist")
	assert.Equal(t, len(bc.Blocks), 3)

	// Failure, duplicate input
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.In = append(tx.In, tx.In[0])
	tx.SignInputs([]cipher.SecKey{genSecret, genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Duplicate spend")
	assert.Equal(t, len(bc.Blocks), 3)

	// Failure, zero coin output
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.Sigs = nil
	tx.PushOutput(genAddress, 0, 100)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Zero coin output")

	// Failure, hash collision with unspents
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	uxOut := CreateUnspents(bc.Head().Head, tx)
	bc.Unspent.Add(uxOut[0])
	assertError(t, bc.VerifyTransaction(tx),
		"New unspent collides with existing unspent")

	// Failure, not spending enough coins
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 50)
	tx.PushOutput(genAddress, 10e6, 100)
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx), "Insufficient coins")

	// Failure, spending outputs we don't own
	_, s := cipher.GenerateKeyPair()
	tx = Transaction{}
	for _, u := range bc.Unspent.Pool {
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
	tx = Transaction{}
	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.PushOutput(genAddress, ux.Body.Coins, ux.Body.Hours)
	tx.UpdateHeader()
	assertError(t, bc.VerifyTransaction(tx),
		"Signature not valid for output being spent")
}

func TestBlockchainVerifyBlock(t *testing.T) {
	bc := NewBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	// Genesis block not valid after the fact
	assert.NotNil(t, bc.VerifyBlock(gb))
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// Valid block
	tx := Transaction{}
	tx.PushInput(ux.Hash())
	tx.PushOutput(genAddress, ux.Body.Coins, ux.CoinHours(bc.Time()))
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	b, err := bc.NewBlockFromTransactions(Transactions{tx}, bc.Time()+_incTime)
	assert.Equal(t, len(b.Body.Transactions), 1)
	assert.Equal(t, len(b.Body.Transactions[0].Out), 1)
	assert.Nil(t, err)
	assert.Nil(t, bc.VerifyBlock(b))

	// Invalid block header
	b.Head.BkSeq = gb.Head.BkSeq
	assert.Equal(t, len(b.Body.Transactions), 1)
	assert.Equal(t, len(b.Body.Transactions[0].Out), 1)
	assertError(t, bc.VerifyBlock(b), "BkSeq invalid")

	// Invalid transactions, makes duplicate outputs
	b.Head.BkSeq = bc.Head().Head.BkSeq + 1
	b.Body.Transactions = append(b.Body.Transactions, b.Body.Transactions[0])
	b.Head.BodyHash = b.HashBody()
	assertError(t, bc.VerifyBlock(b),
		"Duplicate unspent output across transactions")
}

func TestGetUxSnapshot(t *testing.T) {
	unsp := NewUnspentPool()
	xor := randSHA256(t)
	unsp.XorHash = xor
	prev := randSHA256(t)
	sh := getSnapshotHash(unsp, prev)
	expect := cipher.AddSHA256(xor, prev)
	assert.True(t, bytes.Equal(expect[:4], sh[:]))
	assert.NotEqual(t, sh, [4]byte{})
}

func TestVerifyUxSnapshot(t *testing.T) {
	bc := NewBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	b := Block{Body: BlockBody{}, Head: BlockHeader{}}
	b.Body.Transactions = append(b.Body.Transactions, makeTransaction(t))
	bc.Unspent.XorHash = randSHA256(t)
	uxHash := cipher.AddSHA256(bc.Unspent.XorHash, gb.Head.Hash())
	copy(b.Head.UxHash[:], uxHash[:])
	assert.Nil(t, bc.verifyUxSnapshot(b))
	b.Head.UxHash = [4]byte{}
	assertError(t, bc.verifyUxSnapshot(b), "UxSnapshot does not match")
}

func TestVerifyBlockHeader(t *testing.T) {
	bc := NewBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	b := Block{Body: BlockBody{}}
	b.Body.Transactions = append(b.Body.Transactions, makeTransaction(t))
	h := BlockHeader{}

	h.BkSeq = 1
	h.Time = gb.Head.Time + 1
	h.PrevHash = gb.HashHeader()
	h.BodyHash = b.HashBody()

	// Valid header
	b.Head = h
	assert.Nil(t, verifyBlockHeader(bc.Head(), b))

	// Invalid bkSeq
	i := h
	i.BkSeq += 1
	b.Head = i
	assertError(t, verifyBlockHeader(bc.Head(), b), "BkSeq invalid")

	// Invalid time
	i = h
	i.Time = gb.Head.Time
	b.Head = i
	assertError(t, verifyBlockHeader(bc.Head(), b),
		"Block time must be > head time")
	b.Head.Time -= 1
	assertError(t, verifyBlockHeader(bc.Head(), b),
		"Block time must be > head time")

	// Invalid prevHash
	i = h
	i.PrevHash = cipher.SHA256{}
	b.Head = i
	assertError(t, verifyBlockHeader(bc.Head(), b),
		"PrevHash does not match current head")

	// Invalid bodyHash
	i = h
	i.BodyHash = cipher.SHA256{}
	b.Head = i
	assertError(t, verifyBlockHeader(bc.Head(), b),
		"Computed body hash does not match")
}

func TestTransactionFee(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// Valid txn, 100 hours fee
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	fee, err := bc.TransactionFee(&tx)
	assert.Nil(t, err)
	assert.Equal(t, fee, uint64(100))

	// Txn spending unknown output
	tx = Transaction{}
	unknownUx := makeUxOut(t)
	tx.PushInput(unknownUx.Hash())
	_, err = bc.TransactionFee(&tx)
	assertError(t, err, "Unspent output does not exist")

	// Txn spending more hours than avail
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100, 100)
	tx.PushOutput(makeAddress(), 1e6, 10000)
	_, err = bc.TransactionFee(&tx)
	assertError(t, err, "Insufficient coinhours for transaction outputs")
}

func TestProcessTransactions(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// Invalid, no transactions in block
	// arbitrating=false
	txns, err := bc.processTransactions(Transactions{}, false)
	assert.Nil(t, txns)
	assertError(t, err, "No transactions")
	// arbitrating=true
	txns, err = bc.processTransactions(Transactions{}, true)
	assert.Equal(t, len(txns), 0)
	assert.Nil(t, err)

	// Invalid, txn.Verify() fails
	// TODO -- combine all txn.Verify() failures into one test
	// method, and call it from here, from ExecuteBlock(), from
	// Verify(), from VerifyTransaction()
	txns = Transactions{}
	txn := Transaction{}
	txn.PushInput(ux.Hash())
	txn.PushOutput(genAddress, 777, 100)
	txn.SignInputs([]cipher.SecKey{genSecret})
	txn.UpdateHeader()
	txns = append(txns, txn)
	// arbitrating=false
	txns2, err := bc.processTransactions(txns, false)
	assert.Nil(t, txns2)
	assertError(t, err,
		"Transaction outputs must be multiple of 1e6 base units")
	// arbitrating=true
	txns2, err = bc.processTransactions(txns, true)
	assert.NotNil(t, txns2)
	assert.Nil(t, err)
	assert.Equal(t, len(txns2), 0)

	// Invalid, duplicate unspent will be created by these txns
	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	txns = Transactions{txn, txn}
	// arbitrating=false
	txns2, err = bc.processTransactions(txns, false)
	assertError(t, err, "Duplicate unspent output across transactions")
	assert.Nil(t, txns2)
	// arbitrating=true.  One of the offending transactions should be removed
	txns2, err = bc.processTransactions(txns, true)
	assert.Nil(t, err)
	assert.Equal(t, len(txns2), 1)
	assert.Equal(t, txns2[0], txn)

	// Check that a new output will not collide with the existing pool
	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	txns = Transactions{txn}
	uxb := UxBody{
		SrcTransaction: txn.Hash(),
		Coins:          txn.Out[0].Coins,
		Hours:          txn.Out[0].Hours,
		Address:        txn.Out[0].Address,
	}
	bc.Unspent.Add(UxOut{Body: uxb})
	// arbitrating=false
	txns2, err = bc.processTransactions(txns, false)
	assertError(t, err, "New unspent collides with existing unspent")
	assert.Nil(t, txns2)
	// arbitrating=true
	txns2, err = bc.processTransactions(txns, true)
	assert.Equal(t, len(txns2), 0)
	assert.NotNil(t, txns2)
	assert.Nil(t, err)

	// Spending of duplicate inputs being spent across txns
	txn, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	txn2, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	txn2.Out = nil
	txn2.PushOutput(makeAddress(), 1e6, 100)
	txn2.PushOutput(makeAddress(), ux.Body.Coins-1e6, 100)
	txn2.Sigs = nil
	txn2.SignInputs([]cipher.SecKey{genSecret})
	txn2.UpdateHeader()
	txns = SortTransactions(Transactions{txn, txn2}, bc.TransactionFee)
	// arbitrating=false
	txns2, err = bc.processTransactions(txns, false)
	assertError(t, err, "Cannot spend output twice in the same block")
	assert.Nil(t, txns2)
	// arbitrating=true
	txns2, err = bc.processTransactions(txns, true)
	assert.Nil(t, err)
	assert.Equal(t, len(txns2), 1)
	assert.Equal(t, txns2[0], txns[0])
}

func TestExecuteBlock(t *testing.T) {
	bc := NewBlockchain()
	bc.CreateGenesisBlock(genAddress, _genTime, _genCoins)
	assert.Equal(t, len(bc.Blocks), 1)
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, len(bc.Blocks), 3)

	// Invalid block returns error
	b := Block{}
	uxs, err := bc.ExecuteBlock(b)
	assert.NotNil(t, err)
	assert.Nil(t, uxs)

	// Valid block, spends are removed from the unspent pool, new ones are
	// added.  Blocks is updated, and new unspents are returns
	assert.Equal(t, len(bc.Blocks), 3)
	assert.Equal(t, len(bc.Unspent.Pool), 2)
	spuxs := splitUnspent(t, bc, ux)
	tx := Transaction{}
	tx.PushInput(spuxs[0].Hash())
	coins := spuxs[0].Body.Coins
	extra := coins % 4e6
	coins = (coins - extra) / 4
	tx.PushOutput(genAddress, coins+extra, spuxs[0].Body.Hours/5)
	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/6)
	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/7)
	tx.PushOutput(genAddress, coins, spuxs[0].Body.Hours/8)
	tx.SignInputs([]cipher.SecKey{genSecret})
	tx.UpdateHeader()
	tx2 := Transaction{}
	tx2.PushInput(spuxs[1].Hash())
	tx2.PushOutput(genAddress, spuxs[1].Body.Coins, spuxs[1].Body.Hours/10)
	tx2.SignInputs([]cipher.SecKey{genSecret})
	tx2.UpdateHeader()
	txns := Transactions{tx, tx2}
	sTxns := newSortableTransactions(txns, bc.TransactionFee)
	unswapped := sTxns.IsSorted()
	txns = SortTransactions(txns, bc.TransactionFee)
	assert.Nil(t, bc.verifyTransactions(txns))
	seq := bc.Head().Head.BkSeq
	b, err = bc.NewBlockFromTransactions(txns, bc.Time()+_incTime)
	assert.Equal(t, b.Head.BkSeq, seq+1)
	assert.Nil(t, err)
	assert.Equal(t, len(b.Body.Transactions), 2)
	assert.Equal(t, b.Body.Transactions, txns)
	uxs, err = bc.ExecuteBlock(b)
	assert.Nil(t, err)
	assert.Equal(t, len(uxs), 5)
	// Check that all unspents look correct and are in the unspent pool
	txOuts := []TransactionOutput{}
	if unswapped {
		txOuts = append(txOuts, tx.Out...)
		txOuts = append(txOuts, tx2.Out...)
	} else {
		txOuts = append(txOuts, tx2.Out...)
		txOuts = append(txOuts, tx.Out...)
	}
	for i, ux := range uxs {
		if unswapped {
			if i < len(tx.Out) {
				assert.Equal(t, ux.Body.SrcTransaction, tx.Hash())
			} else {
				assert.Equal(t, ux.Body.SrcTransaction, tx2.Hash())
			}
		} else {
			if i < len(tx2.Out) {
				assert.Equal(t, ux.Body.SrcTransaction, tx2.Hash())
			} else {
				assert.Equal(t, ux.Body.SrcTransaction, tx.Hash())
			}
		}
		assert.Equal(t, ux.Body.Address, txOuts[i].Address)
		assert.Equal(t, ux.Body.Coins, txOuts[i].Coins)
		assert.Equal(t, ux.Body.Hours, txOuts[i].Hours)
		assert.Equal(t, ux.Head.BkSeq, b.Head.BkSeq)
		assert.Equal(t, ux.Head.Time, b.Head.Time)
		assert.True(t, bc.Unspent.Has(ux.Hash()))
	}
	// Check that all spends are no longer in the pool
	txIns := []cipher.SHA256{}
	txIns = append(txIns, tx.In...)
	txIns = append(txIns, tx2.In...)
	for _, ux := range txIns {
		assert.False(t, bc.Unspent.Has(ux))
	}
}
