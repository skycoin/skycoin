package coin

import (
	"bytes"
	"math"
	"sort"
	"testing"

	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

func makeTransactionWithSecret(t *testing.T) (Transaction, cipher.SecKey) {
	tx := Transaction{}
	ux, s := makeUxOutWithSecret(t)
	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{s})
	tx.PushOutput(makeAddress(), 1e6, 50)
	tx.PushOutput(makeAddress(), 5e6, 50)
	tx.UpdateHeader()
	return tx, s
}

func makeTransaction(t *testing.T) Transaction {
	tx, _ := makeTransactionWithSecret(t)
	return tx
}

func makeTransactions(t *testing.T, n int) Transactions {
	txns := make(Transactions, n)
	for i, _ := range txns {
		txns[i] = makeTransaction(t)
	}
	return txns
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func manualTransactionsIsSorted(t *testing.T, txns Transactions,
	getFee FeeCalculator) bool {
	isSorted := true
	for i := 0; i < len(txns)-1; i++ {
		ifee, err := getFee(&txns[i])
		assert.Nil(t, err)
		jfee, err := getFee(&txns[i+1])
		assert.Nil(t, err)
		if ifee == jfee {
			hi := txns[i].Hash()
			hj := txns[i+1].Hash()
			if bytes.Compare(hi[:], hj[:]) > 0 {
				isSorted = false
				break
			}
		} else {
			if ifee < jfee {
				isSorted = false
				break
			}
		}
	}
	return isSorted
}

func copyTransaction(tx Transaction) Transaction {
	txo := Transaction{}
	txo.Length = tx.Length
	txo.Type = tx.Type
	txo.InnerHash = tx.InnerHash
	txo.Sigs = make([]cipher.Sig, len(tx.Sigs))
	copy(txo.Sigs, tx.Sigs)
	txo.In = make([]cipher.SHA256, len(tx.In))
	copy(txo.In, tx.In)
	txo.Out = make([]TransactionOutput, len(tx.Out))
	copy(txo.Out, tx.Out)
	return txo
}

func TestTransactionVerify(t *testing.T) {
	// Mismatch header hash
	tx := makeTransaction(t)
	tx.InnerHash = cipher.SHA256{}
	assertError(t, tx.Verify(), "Invalid header hash")

	// No inputs
	tx = makeTransaction(t)
	tx.In = make([]cipher.SHA256, 0)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "No inputs")

	// No outputs
	tx = makeTransaction(t)
	tx.Out = make([]TransactionOutput, 0)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "No outputs")

	// Invalid number of sigs
	tx = makeTransaction(t)
	tx.Sigs = make([]cipher.Sig, 0)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Invalid number of signatures")
	tx.Sigs = make([]cipher.Sig, 20)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Invalid number of signatures")

	// Too many sigs & inputs
	tx = makeTransaction(t)
	tx.Sigs = make([]cipher.Sig, math.MaxUint16)
	tx.In = make([]cipher.SHA256, math.MaxUint16)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Too many signatures and inputs")

	// Duplicate inputs
	tx, s := makeTransactionWithSecret(t)
	tx.PushInput(tx.In[0])
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{s, s})
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Duplicate spend")

	// Duplicate outputs
	tx = makeTransaction(t)
	to := tx.Out[0]
	tx.PushOutput(to.Address, to.Coins, to.Hours)
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Duplicate output in transaction")

	// Invalid signature, empty
	tx = makeTransaction(t)
	tx.Sigs[0] = cipher.Sig{}
	assertError(t, tx.Verify(), "Failed to recover public key")
	// We can't check here for other invalid signatures:
	//      - Signatures signed by someone else, spending coins they don't own
	//      - Signature is for wrong hash
	// This must be done by blockchain tests, because we need the address
	// from the unspent being spent

	// Output coins are not multiples of 1e6
	tx = makeTransaction(t)
	tx.Out[0].Coins += 10
	tx.UpdateHeader()
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{genSecret})
	assert.NotEqual(t, tx.Out[0].Coins%1e6, uint64(0))
	assertError(t, tx.Verify(), "Transaction outputs must be multiple of "+
		"1e6 base units")

	// Output coins are 0
	tx = makeTransaction(t)
	tx.Out[0].Coins = 0
	tx.UpdateHeader()
	assertError(t, tx.Verify(), "Zero coin output")

	// Valid
	tx = makeTransaction(t)
	tx.Out[0].Coins = 10e6
	tx.Out[1].Coins = 1e6
	tx.UpdateHeader()
	assert.Nil(t, tx.Verify())
}

func TestTransactionPushInput(t *testing.T) {
	tx := &Transaction{}
	ux := makeUxOut(t)
	assert.Equal(t, tx.PushInput(ux.Hash()), uint16(0))
	assert.Equal(t, len(tx.In), 1)
	assert.Equal(t, tx.In[0], ux.Hash())
	tx.In = append(tx.In, make([]cipher.SHA256, math.MaxUint16)...)
	ux = makeUxOut(t)
	assert.Panics(t, func() { tx.PushInput(ux.Hash()) })
}

func TestTransactionPushOutput(t *testing.T) {
	tx := &Transaction{}
	a := makeAddress()
	tx.PushOutput(a, 100, 150)
	assert.Equal(t, len(tx.Out), 1)
	assert.Equal(t, tx.Out[0], TransactionOutput{
		Address: a,
		Coins:   100,
		Hours:   150,
	})
	for i := 1; i < 20; i++ {
		a := makeAddress()
		tx.PushOutput(a, uint64(i*100), uint64(i*50))
		assert.Equal(t, len(tx.Out), i+1)
		assert.Equal(t, tx.Out[i], TransactionOutput{
			Address: a,
			Coins:   uint64(i * 100),
			Hours:   uint64(i * 50),
		})
	}
}

func TestTransactionSignInputs(t *testing.T) {
	tx := &Transaction{}
	// Panics if txns already signed
	tx.Sigs = append(tx.Sigs, cipher.Sig{})
	assert.Panics(t, func() { tx.SignInputs([]cipher.SecKey{}) })
	// Panics if not enough keys
	tx = &Transaction{}
	ux, s := makeUxOutWithSecret(t)
	tx.PushInput(ux.Hash())
	ux2, s2 := makeUxOutWithSecret(t)
	tx.PushInput(ux2.Hash())
	tx.PushOutput(makeAddress(), 40, 80)
	assert.Equal(t, len(tx.Sigs), 0)
	assert.Panics(t, func() { tx.SignInputs([]cipher.SecKey{s}) })
	assert.Equal(t, len(tx.Sigs), 0)
	// Valid signing
	h := tx.HashInner()
	assert.NotPanics(t, func() { tx.SignInputs([]cipher.SecKey{s, s2}) })
	assert.Equal(t, len(tx.Sigs), 2)
	assert.Equal(t, tx.HashInner(), h)
	p := cipher.PubKeyFromSecKey(s)
	a := cipher.AddressFromPubKey(p)
	p = cipher.PubKeyFromSecKey(s2)
	a2 := cipher.AddressFromPubKey(p)
	assert.Nil(t, cipher.ChkSig(a, cipher.AddSHA256(h, tx.In[0]), tx.Sigs[0]))
	assert.Nil(t, cipher.ChkSig(a2, cipher.AddSHA256(h, tx.In[1]), tx.Sigs[1]))
	assert.NotNil(t, cipher.ChkSig(a, h, tx.Sigs[1]))
	assert.NotNil(t, cipher.ChkSig(a2, h, tx.Sigs[0]))
}

func TestTransactionHash(t *testing.T) {
	tx := makeTransaction(t)
	assert.NotEqual(t, tx.Hash(), cipher.SHA256{})
	assert.NotEqual(t, tx.HashInner(), tx.Hash())
}

func TestTransactionUpdateHeader(t *testing.T) {
	tx := makeTransaction(t)
	h := tx.InnerHash
	tx.InnerHash = cipher.SHA256{}
	tx.UpdateHeader()
	assert.NotEqual(t, tx.InnerHash, cipher.SHA256{})
	assert.Equal(t, tx.InnerHash, h)
	assert.Equal(t, tx.InnerHash, tx.HashInner())
}

func TestTransactionHashInner(t *testing.T) {
	tx := makeTransaction(t)

	h := tx.HashInner()
	assert.NotEqual(t, h, cipher.SHA256{})

	// If tx.In is changed, hash should change
	tx2 := copyTransaction(tx)
	ux := makeUxOut(t)
	tx2.In[0] = ux.Hash()
	assert.NotEqual(t, tx, tx2)
	assert.Equal(t, tx2.In[0], ux.Hash())
	assert.NotEqual(t, tx.HashInner(), tx2.HashInner())

	// If tx.Out is changed, hash should change
	tx2 = copyTransaction(tx)
	a := makeAddress()
	tx2.Out[0].Address = a
	assert.NotEqual(t, tx, tx2)
	assert.Equal(t, tx2.Out[0].Address, a)
	assert.NotEqual(t, tx.HashInner(), tx2.HashInner())

	// If tx.Head is changed, hash should not change
	tx2 = copyTransaction(tx)
	tx.Sigs = append(tx.Sigs, cipher.Sig{})
	assert.Equal(t, tx.HashInner(), tx2.HashInner())
}

func TestTransactionSerialization(t *testing.T) {
	tx := makeTransaction(t)
	b := tx.Serialize()
	tx2 := TransactionDeserialize(b)
	assert.Equal(t, tx, tx2)
	// Invalid deserialization
	assert.Panics(t, func() { TransactionDeserialize([]byte{0x04}) })
}

func TestTransactionOutputHours(t *testing.T) {
	tx := Transaction{}
	tx.PushOutput(makeAddress(), 1e6, 100)
	tx.PushOutput(makeAddress(), 1e6, 200)
	tx.PushOutput(makeAddress(), 1e6, 500)
	tx.PushOutput(makeAddress(), 1e6, 0)
	assert.Equal(t, tx.OutputHours(), uint64(800))
}

func TestTransactionFees(t *testing.T) {
	bc := NewBlockchain(&FakeTree{}, nil)
	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	assert.Equal(t, bc.Len(), uint64(1))
	_, ux := addBlockToBlockchain(t, bc)
	assert.Equal(t, bc.Len(), uint64(3))

	// Valid txn, 100 hours fee
	tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	fee, err := Transactions{tx}.Fees(bc.TransactionFee)
	assert.Nil(t, err)
	assert.Equal(t, fee, uint64(100))

	// Multiple txns, 100 hours fee each
	tx2, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	fee, err = Transactions{tx, tx2}.Fees(bc.TransactionFee)
	assert.Nil(t, err)
	assert.Equal(t, fee, uint64(200))

	// Txn spending unknown output
	tx = Transaction{}
	unknownUx := makeUxOut(t)
	tx.PushInput(unknownUx.Hash())
	_, err = Transactions{tx}.Fees(bc.TransactionFee)
	assertError(t, err, "Unspent output does not exist")

	// Txn spending more hours than avail
	tx, _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret, 100,
		100)
	tx.PushOutput(makeAddress(), 1e6, 10000)
	_, err = Transactions{tx}.Fees(bc.TransactionFee)
	assertError(t, err, "Insufficient coinhours for transaction outputs")
}

func TestNewSortableTransactions(t *testing.T) {
	bc := NewBlockchain(&FakeTree{}, nil)

	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	_, ux := addBlockToBlockchain(t, bc)
	txns := make(Transactions, 4)
	for i := range txns {
		tx, _ := makeTransactionForChainWithHoursFee(t, bc, ux, genSecret,
			100, uint64(i*100))
		txns[i] = tx
	}
	sTxns := newSortableTransactions(txns, bc.TransactionFee)
	assert.Equal(t, len(sTxns.Txns), len(txns))
	assert.Equal(t, len(sTxns.Fees), len(txns))
	assert.Equal(t, len(sTxns.Hashes), len(txns))
	for i, tx := range txns {
		assert.Equal(t, sTxns.Txns[i], tx)
		assert.Equal(t, sTxns.Hashes[i], tx.Hash())
		fee, err := bc.TransactionFee(&tx)
		assert.Nil(t, err)
		assert.Equal(t, sTxns.Fees[i], (fee*1024)/uint64(tx.Size()))
	}
}

func TestTransactionsSize(t *testing.T) {
	txns := makeTransactions(t, 10)
	size := 0
	for _, tx := range txns {
		size += len(encoder.Serialize(&tx))
	}
	assert.NotEqual(t, size, 0)
	assert.Equal(t, txns.Size(), size)
}

func TestTransactionSorting(t *testing.T) {
	bc := NewBlockchain(&FakeTree{}, nil)
	bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	_, ux := addBlockToBlockchain(t, bc)
	txns := make(Transactions, 4)
	for i := 0; i < len(txns); i++ {
		fee := uint64(0)
		if i == 0 || i == 2 {
			fee = uint64(1000)
		} else {
			fee = uint64(i * 100)
		}
		txns[i], _ = makeTransactionForChainWithHoursFee(t, bc, ux, genSecret,
			100, fee)
	}

	// TODO -- check that things are actually sorted, and test with something
	// that has matching fee

	// Sort(), IsSorted(), Less()
	isSorted := manualTransactionsIsSorted(t, txns, bc.TransactionFee)
	sTxns := newSortableTransactions(txns, bc.TransactionFee)
	for i := range txns {
		assert.Equal(t, sTxns.Txns[i], txns[i])
		assert.Equal(t, sTxns.Hashes[i], txns[i].Hash())
		fee, err := bc.TransactionFee(&txns[i])
		assert.Nil(t, err)
		assert.Equal(t, sTxns.Fees[i], (fee*1024)/uint64(txns[i].Size()))
	}

	assert.Equal(t, sort.IsSorted(sTxns), isSorted)
	assert.Equal(t, sTxns.IsSorted(), isSorted)
	if isSorted {
		txns[0], txns[1] = txns[1], txns[0]
		sTxns = newSortableTransactions(txns, bc.TransactionFee)
		assert.False(t, sTxns.Less(0, 1))
		assert.True(t, sTxns.Less(1, 0))
	}
	sTxns = newSortableTransactions(txns, bc.TransactionFee)
	assert.False(t, manualTransactionsIsSorted(t, txns, bc.TransactionFee))
	assert.False(t, sort.IsSorted(sTxns))
	assert.False(t, sTxns.IsSorted())
	txns2 := SortTransactions(txns, bc.TransactionFee)
	assert.True(t, manualTransactionsIsSorted(t, txns2, bc.TransactionFee))
	sTxns = newSortableTransactions(txns2, bc.TransactionFee)
	assert.True(t, sort.IsSorted(sTxns))
	assert.True(t, sTxns.IsSorted())
	for i := 0; i < len(txns2)-1; i++ {
		assert.True(t, sTxns.Less(i, i+1))
		assert.False(t, sTxns.Less(i+1, i))
	}

	// Check that sorting works
	sTxns = newSortableTransactions(txns, bc.TransactionFee)
	sTxns.Sort()
	hashChecked := false
	for i := range txns[:len(txns)-1] {
		j := i + 1
		assert.True(t, sTxns.Fees[i] >= sTxns.Fees[j])
		if sTxns.Fees[i] == sTxns.Fees[j] {
			hashChecked = true
			cmp := bytes.Compare(sTxns.Hashes[i][:], sTxns.Hashes[j][:])
			assert.True(t, cmp < 0)
		}
	}
	assert.True(t, hashChecked)

	// Len()
	assert.Equal(t, len(txns), sTxns.Len())
	assert.Equal(t, len(sTxns.Txns), sTxns.Len())
	assert.Equal(t, len(sTxns.Fees), sTxns.Len())
	assert.Equal(t, len(sTxns.Hashes), sTxns.Len())
	assert.Equal(t, 4, sTxns.Len())

	// Swap()
	tx1 := sTxns.Txns[0]
	tx2 := sTxns.Txns[1]
	fee1 := sTxns.Fees[0]
	fee2 := sTxns.Fees[1]
	hash1 := sTxns.Hashes[0]
	hash2 := sTxns.Hashes[1]
	sTxns.Swap(0, 1)
	assert.Equal(t, sTxns.Txns[0], tx2)
	assert.Equal(t, sTxns.Txns[1], tx1)
	assert.Equal(t, sTxns.Fees[0], fee2)
	assert.Equal(t, sTxns.Fees[1], fee1)
	assert.Equal(t, sTxns.Hashes[0], hash2)
	assert.Equal(t, sTxns.Hashes[1], hash1)
	sTxns.Swap(0, 1)
	assert.Equal(t, sTxns.Txns[0], tx1)
	assert.Equal(t, sTxns.Txns[1], tx2)
	assert.Equal(t, sTxns.Fees[0], fee1)
	assert.Equal(t, sTxns.Fees[1], fee2)
	assert.Equal(t, sTxns.Hashes[0], hash1)
	assert.Equal(t, sTxns.Hashes[1], hash2)
	sTxns.Swap(1, 0)
	assert.Equal(t, sTxns.Txns[0], tx2)
	assert.Equal(t, sTxns.Txns[1], tx1)
	assert.Equal(t, sTxns.Fees[0], fee2)
	assert.Equal(t, sTxns.Fees[1], fee1)
	assert.Equal(t, sTxns.Hashes[0], hash2)
	assert.Equal(t, sTxns.Hashes[1], hash1)
	sTxns.Swap(1, 0)
	assert.Equal(t, sTxns.Txns[0], tx1)
	assert.Equal(t, sTxns.Txns[1], tx2)
	assert.Equal(t, sTxns.Fees[0], fee1)
	assert.Equal(t, sTxns.Fees[1], fee2)
	assert.Equal(t, sTxns.Hashes[0], hash1)
	assert.Equal(t, sTxns.Hashes[1], hash2)

	// SortTransaction()
	sTxns.Sort()
	assert.True(t, sTxns.IsSorted())
	assert.NotEqual(t, txns, sTxns.Txns)
	txns2 = SortTransactions(txns, bc.TransactionFee)
	assert.Equal(t, sTxns.Txns, txns2)
	sTxns2 := newSortableTransactions(txns2, bc.TransactionFee)
	assert.True(t, sTxns2.IsSorted())
	assert.Equal(t, sTxns, sTxns2)
}

func TestTransactionsHashes(t *testing.T) {
	txns := make(Transactions, 4)
	for i := 0; i < len(txns); i++ {
		txns[i] = makeTransaction(t)
	}
	hashes := txns.Hashes()
	assert.Equal(t, len(hashes), 4)
	for i, h := range hashes {
		assert.Equal(t, h, txns[i].Hash())
	}
}

func TestTransactionsTruncateBytesTo(t *testing.T) {
	txns := makeTransactions(t, 10)
	trunc := 0
	for i := 0; i < len(txns)/2; i++ {
		trunc += txns[i].Size()
	}
	// Truncating halfway
	txns2 := txns.TruncateBytesTo(trunc)
	assert.Equal(t, len(txns2), len(txns)/2)
	assert.Equal(t, txns2.Size(), trunc)

	// Stepping into next boundary has same cutoff, must exceed
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	assert.Equal(t, len(txns2), len(txns)/2)
	assert.Equal(t, txns2.Size(), trunc-1)

	// Moving to 1 before next level
	trunc += txns[5].Size() - 2
	txns2 = txns.TruncateBytesTo(trunc)
	assert.Equal(t, len(txns2), len(txns)/2)
	assert.Equal(t, txns2.Size(), trunc-txns[5].Size()+1)

	// Moving to next level
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	assert.Equal(t, len(txns2), len(txns)/2+1)
	assert.Equal(t, txns2.Size(), trunc)

	// Truncating to full available amt
	trunc = txns.Size()
	txns2 = txns.TruncateBytesTo(trunc)
	assert.Equal(t, txns, txns2)
	assert.Equal(t, txns2.Size(), trunc)

	// Truncating over amount
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	assert.Equal(t, txns, txns2)
	assert.Equal(t, txns2.Size(), trunc-1)

	// Truncating to 0
	trunc = 0
	txns2 = txns.TruncateBytesTo(0)
	assert.Equal(t, len(txns2), 0)
	assert.Equal(t, txns2.Size(), trunc)
}

func TestFullTransaction(t *testing.T) {
	p1, s1 := cipher.GenerateKeyPair()
	a1 := cipher.AddressFromPubKey(p1)
	bc := NewBlockchain(&FakeTree{}, nil)
	bc.CreateGenesisBlock(a1, _genCoins, _genTime)
	tx := Transaction{}
	ux := bc.unspent.Array()[0]
	tx.PushInput(ux.Hash())
	p2, s2 := cipher.GenerateKeyPair()
	a2 := cipher.AddressFromPubKey(p2)
	tx.PushOutput(a1, ux.Body.Coins-6e6, 100)
	tx.PushOutput(a2, 1e6, 100)
	tx.PushOutput(a2, 5e6, 100)
	tx.SignInputs([]cipher.SecKey{s1})
	tx.UpdateHeader()
	assert.Nil(t, tx.Verify())
	assert.Nil(t, bc.VerifyTransaction(tx))
	b, err := bc.NewBlockFromTransactions(Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)
	_, err = bc.ExecuteBlock(&b)
	assert.Nil(t, err)

	txo := CreateUnspents(bc.Head().Head, tx)
	tx = Transaction{}
	assert.Equal(t, txo[0].Body.Address, a1)
	assert.Equal(t, txo[1].Body.Address, a2)
	assert.Equal(t, txo[2].Body.Address, a2)
	ux0, ok := bc.unspent.Get(txo[0].Hash())
	assert.True(t, ok)
	ux1, ok := bc.unspent.Get(txo[1].Hash())
	assert.True(t, ok)
	ux2, ok := bc.unspent.Get(txo[2].Hash())
	assert.True(t, ok)
	tx.PushInput(ux0.Hash())
	tx.PushInput(ux1.Hash())
	tx.PushInput(ux2.Hash())
	tx.PushOutput(a2, 10e6, 200)
	tx.PushOutput(a1, ux.Body.Coins-10e6, 100)
	tx.SignInputs([]cipher.SecKey{s1, s2, s2})
	tx.UpdateHeader()
	assert.Nil(t, tx.Verify())
	assert.Nil(t, bc.VerifyTransaction(tx))
	b, err = bc.NewBlockFromTransactions(Transactions{tx}, bc.Time()+_incTime)
	assert.Nil(t, err)
	_, err = bc.ExecuteBlock(&b)
	assert.Nil(t, err)
}
