// +build ignore
// These tests need to be rewritten to conform with blockdb changes

package visor

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
)

// const (
// 	testBlockSize = 1024 * 1024
// )

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(randBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          10e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func makeTransactionWithSecret(t *testing.T) (coin.Transaction, cipher.SecKey) {
	tx := coin.Transaction{}
	ux, s := makeUxOutWithSecret(t)
	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{s})
	tx.PushOutput(makeAddress(), 10e6, 100)
	tx.UpdateHeader()
	return tx, s
}

func makeTransaction(t *testing.T) coin.Transaction {
	tx, _ := makeTransactionWithSecret(t)
	return tx
}

func randBytes(t *testing.T, n int) []byte {
	b := make([]byte, n)
	x, err := rand.Read(b)
	assert.Equal(t, n, x)
	assert.Nil(t, err)
	return b
}

func getFee(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func prepareDB(t *testing.T) (*bolt.DB, func()) {
	f := fmt.Sprintf("test%d.db", rand.Intn(1024))
	db, err := bolt.Open(f, 0700, nil)
	assert.Nil(t, err)
	return db, func() {
		db.Close()
		os.Remove(f)
	}
}

// func makeValidTxn(t *testing.T, utp *UnconfirmedTxnPool) (coin.Transaction, error) {
// 	w := wallet.NewWallet("test")
// 	w.GenerateAddresses(2)
// 	now := tNow()
// 	a := makeAddress()

// 	uxs := makeUxBalancesForAddresses([]wallet.Balance{
// 		wallet.Balance{10e6, 150},
// 		wallet.Balance{15e6, 150},
// 	}, now, w.GetAddresses()[:2])
// 	unsp := coin.NewUnspentPool()
// 	addUxArrayToUnspentPool(&unsp, uxs)
// 	amt := wallet.Balance{10 * 1e6, 0}
// 	return CreateSpendingTransaction(w, utp, &unsp, now, amt, a)
// }

// func makeValidTxnWithFeeFactor(mv *Visor,
// 	factor, extra uint64) (coin.Transaction, error) {
// 	we := wallet.NewWalletEntry()
// 	tmp := mv.Config.CoinHourBurnFactor
// 	mv.Config.CoinHourBurnFactor = factor
// 	tx, err := mv.Spend(mv.Wallets[0].GetFilename(), wallet.Balance{10 * 1e6, 1000},
// 		extra, we.Address)
// 	mv.Config.CoinHourBurnFactor = tmp
// 	return tx, err
// }

// func makeValidTxnWithFeeFactorAndExtraChange(mv *Visor,
// 	factor, extra, change uint64) (coin.Transaction, error) {
// 	we := wallet.NewWalletEntry()
// 	tmp := mv.Config.CoinHourBurnFactor
// 	mv.Config.CoinHourBurnFactor = factor
// 	tx, err := mv.Spend(mv.Wallets[0].GetFilename(), wallet.Balance{10 * 1e6, 1002},
// 		extra, we.Address)
// 	mv.Config.CoinHourBurnFactor = tmp
// 	return tx, err
// }

func makeValidTxnNoChange(t *testing.T) (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	db, close := prepareDB(t)
	defer close()
	uncf := NewUnconfirmedTxnPool(db)
	now := tNow()
	a := makeAddress()
	uxs := makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp := coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt := wallet.Balance{25 * 1e6, 0}
	return CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
}

func makeInvalidTxn(t *testing.T) (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	db, close := prepareDB(t)
	defer close()
	uncf := NewUnconfirmedTxnPool(db)
	now := tNow()
	a := makeAddress()
	uxs := makeUxBalancesForAddresses([]wallet.Balance{}, now, w.GetAddresses()[:2])
	unsp := coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt := wallet.Balance{25 * 1e6, 0}
	txn, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	if err != nil {
		return txn, err
	}

	txn.Out[0].Address = cipher.Address{}
	return txn, nil
}

func assertValidUnspent(t *testing.T, bc *Blockchain,
	unspent *txUnspents, tx coin.Transaction) {
	expect := coin.CreateUnspents(bc.Head().Head, tx)
	assert.NotEqual(t, len(expect), 0)
	sum := 0

	unspent.forEach(func(_ cipher.SHA256, uxo coin.UxArray) {
		sum += len(uxo)
	})

	assert.Equal(t, len(expect), sum)
	uxs, err := unspent.get(tx.Hash())
	assert.Nil(t, err)
	for _, ux := range expect {
		found := false
		for _, u := range uxs {
			if u.Hash() == ux.Hash() {
				found = true
				break
			}
		}
		assert.True(t, found)
	}
}

func assertValidUnconfirmed(t *testing.T, txns map[cipher.SHA256]UnconfirmedTxn,
	txn coin.Transaction) {
	ut, ok := txns[txn.Hash()]
	assert.True(t, ok)
	assert.Equal(t, ut.Txn, txn)
	assert.True(t, nanoToTime(ut.Announced).IsZero())
	assert.False(t, nanoToTime(ut.Received).IsZero())
	assert.False(t, nanoToTime(ut.Checked).IsZero())
}

func createUnconfirmedTxns(t *testing.T, up *UnconfirmedTxnPool, bc *Blockchain, n int) []UnconfirmedTxn {
	uts := make([]UnconfirmedTxn, 4)
	// usp := coin.NewUnspentPool()
	for i := 0; i < len(uts); i++ {
		// tx, _ := makeValidTxn(t, up)
		tx := makeTransactionForChain(t, bc)
		ut := up.createUnconfirmedTxn(bc.GetUnspent(), tx)
		uts[i] = ut
		up.Txns.put(&ut)
	}
	assert.Equal(t, up.Txns.len(), 4)
	return uts
}

// func TestVerifyTransaction(t *testing.T) {
// 	mv := setupMasterVisor()
// 	tx, err := makeValidTxn(mv)
// 	assert.Nil(t, err)
// 	err = VerifyTransaction(mv.blockchain, &tx, tx.Size()-1, 0)
// 	assertError(t, err, "Transaction too large")
// 	err = VerifyTransaction(mv.blockchain, &tx, tx.Size(), 2)
// 	assertError(t, err, "Transaction fee minimum not met")
// 	err = VerifyTransaction(mv.blockchain, &tx, tx.Size(), 0)
// 	assert.Nil(t, err)
// 	tx, err = makeValidTxnWithFeeFactor(mv, 4, 10)
// 	assert.Nil(t, err)
// 	err = VerifyTransaction(mv.blockchain, &tx, tx.Size(), 4)
// 	assert.Nil(t, err)

// 	// Make sure that the minimum fee is floor(output/factor)
// 	tx, err = makeValidTxnWithFeeFactorAndExtraChange(mv, 4, 10, 2)
// 	assert.NotEqual(t, tx.OutputHours()%4, 0)
// 	assert.Nil(t, err)
// 	err = VerifyTransaction(mv.blockchain, &tx, tx.Size(), 4)
// 	assert.Nil(t, err)
// }

func TestUnconfirmedTxnHash(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	uts := createUnconfirmedTxns(t, up, bc, 1)
	utx := uts[0]
	assert.Equal(t, utx.Hash(), utx.Txn.Hash())
}

func TestNewUnconfirmedTxnPool(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	assert.NotNil(t, ut.Txns)
	assert.Equal(t, ut.Txns.len(), 0)
}

func TestSetAnnounced(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	assert.Equal(t, ut.Txns.len(), 0)
	// Unknown should be safe and a noop
	assert.NotPanics(t, func() {
		ut.SetAnnounced(cipher.SHA256{}, utc.Now())
	})
	assert.Equal(t, ut.Txns.len(), 0)
	bc := makeBlockchain()
	utx := createUnconfirmedTxns(t, ut, bc, 1)[0]
	assert.True(t, nanoToTime(utx.Announced).IsZero())
	ut.Txns.put(&utx)
	now := utc.Now()
	ut.SetAnnounced(utx.Hash(), now)
	v, ok := ut.Txns.get(utx.Hash())
	assert.True(t, ok)
	assert.Equal(t, v.Announced, now.UnixNano())
}

func makeBlockchain() *Blockchain {
	ft := FakeTree{}
	b := NewBlockchain(&ft, nil)
	b.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	return b
}

func TestInjectTxn(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	txn := makeTransactionForChain(t, bc)
	known, err := ut.InjectTxn(bc, txn)
	assert.Nil(t, err)

	assert.False(t, known)
	assert.Equal(t, ut.Txns.len(), 1)

	assertValidUnspent(t, bc, ut.Unspent, txn)
	allUncfmTxs, err := ut.Txns.getAll()
	assert.Nil(t, err)
	uncfmMap := make(map[cipher.SHA256]UnconfirmedTxn, len(allUncfmTxs))
	for _, txn := range allUncfmTxs {
		uncfmMap[txn.Hash()] = txn
	}
	assertValidUnconfirmed(t, uncfmMap, txn)

	// Test where we are receiver of ux outputs
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxn(mv)
	// assert.Nil(t, err)
	// addrs := make(map[cipher.Address]byte, 1)
	// addrs[txn.Out[1].Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	// assert.Nil(t, err)
	// assert.False(t, known)
	// assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	// assertValidUnconfirmed(t, ut.Txns, txn)

	// // Test where we are spender of ux outputs
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxnNoChange(mv)
	// assert.Nil(t, err)
	// addrs = make(map[cipher.Address]byte, 1)
	// ux, ok := mv.blockchain.Unspent.Get(txn.In[0])
	// assert.True(t, ok)
	// addrs[ux.Body.Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	// assert.Nil(t, err)
	// assert.False(t, known)
	// assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	// assertValidUnconfirmed(t, ut.Txns, txn)

	// // Test where we are both spender and receiver of ux outputs
	// mv = setupMasterVisor()
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxn(mv)
	// assert.Nil(t, err)
	// addrs = make(map[cipher.Address]byte, 2)
	// addrs[txn.Out[0].Address] = byte(1)
	// ux, ok = mv.blockchain.Unspent.Get(txn.In[0])
	// assert.True(t, ok)
	// addrs[ux.Body.Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	// assert.Nil(t, err)
	// assert.False(t, known)
	// assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	// assertValidUnconfirmed(t, ut.Txns, txn)
	// assert.Equal(t, len(ut.Txns), 1)
	// assert.Equal(t, len(ut.Unspent), 1)
	// for _, uxs := range ut.Unspent {
	// 	assert.Equal(t, len(uxs), 2)
	// }

	// Test duplicate Record, should be no-op besides state change
	utx, ok := ut.Txns.get(txn.Hash())
	assert.True(t, ok)
	// Set a placeholder value on the utx to check if we overwrote it
	utx.Announced = time.Time{}.Add(time.Minute).UnixNano()
	ut.Txns.put(utx)
	known, err = ut.InjectTxn(bc, txn)
	assert.Nil(t, err)
	assert.True(t, known)
	utx2, ok := ut.Txns.get(txn.Hash())
	assert.True(t, ok)
	assert.Equal(t, utx2.Announced, time.Time{}.Add(time.Minute).UnixNano())
	utx2.Announced = time.Time{}.UnixNano()
	ut.Txns.put(utx2)
	// Received & checked should be updated
	assert.True(t, nanoToTime(utx2.Received).After(nanoToTime(utx.Received)))
	assert.True(t, nanoToTime(utx2.Checked).After(nanoToTime(utx.Checked)))
	all, err := ut.Txns.getAll()
	assert.Nil(t, err)
	allMap := make(map[cipher.SHA256]UnconfirmedTxn, len(all))
	for _, tx := range all {
		allMap[tx.Hash()] = tx
	}
	assertValidUnconfirmed(t, allMap, txn)
	assert.Equal(t, ut.Txns.len(), 1)
	assert.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxa coin.UxArray) {
		assert.Equal(t, len(uxa), 2)
	})

	// Test with valid fee, exact
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxnWithFeeFactor(mv, 4, 0)
	// assert.Nil(t, err)
	// addrs = make(map[cipher.Address]byte, 1)
	// addrs[txn.Out[1].Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	// assert.Nil(t, err)
	// assert.False(t, known)
	// assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	// assertValidUnconfirmed(t, ut.Txns, txn)

	// // Test with valid fee, surplus
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxnWithFeeFactor(mv, 4, 100)
	// assert.Nil(t, err)
	// addrs = make(map[cipher.Address]byte, 1)
	// addrs[txn.Out[1].Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	// assert.Nil(t, err)
	// assert.False(t, known)
	// assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	// assertValidUnconfirmed(t, ut.Txns, txn)

	// // Test with invalid fee
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxnWithFeeFactor(mv, 5, 0)
	// assert.Nil(t, err)
	// _, err = mv.blockchain.TransactionFee(&txn)
	// assert.Nil(t, err)
	// addrs = make(map[cipher.Address]byte, 1)
	// addrs[txn.Out[1].Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	// assertError(t, err, "Transaction fee minimum not met")
	// assert.False(t, known)
	// assert.Equal(t, len(ut.Txns), 0)

	// // Test with bc.TransactionFee failing
	// mv = setupMasterVisor()
	// assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	// ut = NewUnconfirmedTxnPool()
	// txn, err = makeValidTxnWithFeeFactor(mv, 4, 100)
	// assert.Nil(t, err)
	// txn.Out[1].Hours = 1e16
	// txn.Head.Sigs = make([]cipher.Sig, 0)
	// txn.SignInputs([]cipher.SecKey{mv.Config.MasterKeys.Secret})
	// txn.UpdateHeader()
	// addrs = make(map[cipher.Address]byte, 1)
	// addrs[txn.Out[1].Address] = byte(1)
	// err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	// assertError(t, err, "Insufficient coinhours for transaction outputs")
	// assert.False(t, known)
	// assert.Equal(t, len(ut.Txns), 0)
}

func TestRawTxns(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	utxs := make(coin.Transactions, 4)
	for i := 0; i < len(utxs); i++ {
		utx := addUnconfirmedTxnToPool(ut)
		utxs[i] = utx.Txn
	}
	utxs = coin.SortTransactions(utxs, getFee)
	txns := ut.RawTxns()
	txns = coin.SortTransactions(txns, getFee)
	for i, tx := range txns {
		assert.Equal(t, utxs[i], tx)
	}
}

func TestRemoveTxn(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)

	bc := makeBlockchain()
	utx := makeTransactionForChain(t, bc)
	known, err := ut.InjectTxn(bc, utx)
	assert.Nil(t, err)
	assert.False(t, known)
	assert.Equal(t, ut.Txns.len(), 1)
	assert.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxa coin.UxArray) {
		assert.Equal(t, len(uxa), 2)
	})

	// Unknown txn is no-op
	badh := randSHA256()
	assert.NotEqual(t, badh, utx.Hash())
	assert.Equal(t, ut.Txns.len(), 1)
	assert.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	ut.removeTxn(bc, badh)
	assert.Equal(t, ut.Txns.len(), 1)
	assert.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})

	// Known txn updates Txns, predicted Unspents
	// utx2, err := makeValidTxn(mv)
	utx2 := makeTransactionForChain(t, bc)
	assert.Nil(t, err)
	known, err = ut.InjectTxn(bc, utx2)
	assert.Nil(t, err)
	assert.False(t, known)
	assert.Equal(t, ut.Txns.len(), 2)
	assert.Equal(t, ut.Unspent.len(), 2)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	ut.removeTxn(bc, utx.Hash())
	assert.Equal(t, ut.Txns.len(), 1)
	assert.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	// ut.removeTxn(bc, utx.Hash())
	// assert.Equal(t, ut.Len(), 1)
	// assert.Equal(t, ut.Unspent.len(), 1)
	// ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
	// 	assert.Equal(t, len(uxs), 2)
	// })
	ut.removeTxn(bc, utx2.Hash())
	assert.Equal(t, ut.Len(), 0)
	assert.Equal(t, ut.Unspent.len(), 0)
}

func TestRemoveTxns(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)

	// Include an unknown hash, and omit a known hash. The other two should
	// be removed.
	hashes := make([]cipher.SHA256, 0, 3)
	hashes = append(hashes, randSHA256()) // unknown hash
	bc := makeBlockchain()
	ut := makeTransactionForChain(t, bc)
	known, err := up.InjectTxn(bc, ut)
	assert.Nil(t, err)
	assert.False(t, known)
	hashes = append(hashes, ut.Hash())
	ut2 := makeTransactionForChain(t, bc)
	known, err = up.InjectTxn(bc, ut2)
	assert.False(t, known)
	hashes = append(hashes, ut2.Hash())
	ut3 := makeTransactionForChain(t, bc)
	known, err = up.InjectTxn(bc, ut3)
	assert.False(t, known)

	assert.Equal(t, up.Unspent.len(), 3)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	assert.Equal(t, up.Len(), 3)
	up.removeTxns(bc, hashes)
	assert.Equal(t, up.Unspent.len(), 1)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	assert.Equal(t, up.Len(), 1)
	_, ok := up.Txns.get(ut3.Hash())
	assert.True(t, ok)
}

func TestRemoveTransactions(t *testing.T) {
	bc := makeBlockchain()
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)

	// Include an unknown txn, and omit a known hash. The other two should
	// be removed.
	unkUt := makeTransactionForChain(t, bc)
	txns := make(coin.Transactions, 0, 3)
	txns = append(txns, unkUt) // unknown txn
	ut := makeTransactionForChain(t, bc)
	known, err := up.InjectTxn(bc, ut)
	assert.Nil(t, err)
	assert.False(t, known)
	txns = append(txns, ut)
	ut2 := makeTransactionForChain(t, bc)
	assert.Nil(t, err)
	known, err = up.InjectTxn(bc, ut2)
	assert.Nil(t, err)
	assert.False(t, known)
	txns = append(txns, ut2)
	ut3 := makeTransactionForChain(t, bc)
	known, err = up.InjectTxn(bc, ut3)
	assert.Nil(t, err)
	assert.False(t, known)

	assert.Equal(t, up.Unspent.len(), 3)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	assert.Equal(t, up.Txns.len(), 3)
	up.RemoveTransactions(bc, txns)
	assert.Equal(t, up.Unspent.len(), 1)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		assert.Equal(t, len(uxs), 2)
	})
	assert.Equal(t, up.Txns.len(), 1)
	_, ok := up.Txns.get(ut3.Hash())
	assert.True(t, ok)
}

// func testRefresh(t *testing.T, mv *Visor, refresh func(checkPeriod, maxAge time.Duration)) {
// 	up := mv.Unconfirmed
// 	// Add a transaction that is invalid, but will not be checked yet
// 	// Add a transaction that is invalid, and will be checked and removed
// 	invalidTxUnchecked := makeTransactionForChain(t, bc)
// 	invalidTxChecked := makeTransaction(t, bc)
// 	assert.Nil(t, err)
// 	assert.Nil(t, invalidTxUnchecked.Verify())
// 	assert.Nil(t, invalidTxChecked.Verify())
// 	// Invalidate it by spending the output that this txn references
// 	invalidator, err := makeValidTxn(mv)
// 	assert.Nil(t, err)
// 	err, known := up.InjectTxn(mv.blockchain, invalidator, nil, testBlockSize, 0)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	assert.Equal(t, len(up.Txns), 1)
// 	_, err = mv.CreateAndExecuteBlock()
// 	assert.Nil(t, err)
// 	assert.Equal(t, len(up.Txns), 0)
// 	assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxUnchecked))
// 	assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxChecked))

// 	invalidUtxUnchecked := UnconfirmedTxn{
// 		Txn:       invalidTxUnchecked,
// 		Received:  utc.Now(),
// 		Checked:   utc.Now(),
// 		Announced: time.Time{},
// 	}
// 	invalidUtxChecked := invalidUtxUnchecked
// 	invalidUtxChecked.Txn = invalidTxChecked
// 	invalidUtxUnchecked.Checked = utc.Now().Add(time.Hour)
// 	invalidUtxChecked.Checked = utc.Now().Add(-time.Hour)
// 	up.Txns[invalidUtxUnchecked.Hash()] = invalidUtxUnchecked
// 	up.Txns[invalidUtxChecked.Hash()] = invalidUtxChecked
// 	assert.Equal(t, len(up.Txns), 2)
// 	uncheckedHash := invalidTxUnchecked.Hash()
// 	checkedHash := invalidTxChecked.Hash()
// 	up.Unspent[uncheckedHash] = coin.CreateUnspents(coin.BlockHeader{},
// 		invalidTxUnchecked)
// 	up.Unspent[checkedHash] = coin.CreateUnspents(coin.BlockHeader{},
// 		invalidTxChecked)

// 	// Create a transaction that is valid, and will not be checked yet
// 	validTxUnchecked, err := makeValidTxn(mv)
// 	assert.Nil(t, err)
// 	// Create a transaction that is valid, and will be checked
// 	validTxChecked, err := makeValidTxn(mv)
// 	assert.Nil(t, err)
// 	// Create a transaction that is expired
// 	validTxExpired, err := makeValidTxn(mv)
// 	assert.Nil(t, err)

// 	// Add the transaction that is valid, and will not be checked yet
// 	err, known = up.InjectTxn(mv.blockchain, validTxUnchecked, nil,
// 		testBlockSize, 0)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	validUtxUnchecked := up.Txns[validTxUnchecked.Hash()]
// 	validUtxUnchecked.Checked = utc.Now().Add(time.Hour)
// 	up.Txns[validUtxUnchecked.Hash()] = validUtxUnchecked
// 	assert.Equal(t, len(up.Txns), 3)

// 	// Add the transaction that is valid, and will be checked
// 	err, known = up.InjectTxn(mv.blockchain, validTxChecked, nil,
// 		testBlockSize, 0)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	validUtxChecked := up.Txns[validTxChecked.Hash()]
// 	validUtxChecked.Checked = utc.Now().Add(-time.Hour)
// 	up.Txns[validUtxChecked.Hash()] = validUtxChecked
// 	assert.Equal(t, len(up.Txns), 4)

// 	// Add the transaction that is expired
// 	err, known = up.InjectTxn(mv.blockchain, validTxExpired, nil,
// 		testBlockSize, 0)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	validUtxExpired := up.Txns[validTxExpired.Hash()]
// 	validUtxExpired.Received = utc.Now().Add(-time.Hour)
// 	up.Txns[validTxExpired.Hash()] = validUtxExpired
// 	assert.Equal(t, len(up.Txns), 5)

// 	// Pre-sanity check
// 	assert.Equal(t, len(up.Unspent), 5)
// 	for _, uxs := range up.Unspent {
// 		assert.Equal(t, len(uxs), 2)
// 	}
// 	assert.Equal(t, len(up.Txns), 5)

// 	// Refresh
// 	checkPeriod := time.Second * 2
// 	maxAge := time.Second * 4
// 	refresh(checkPeriod, maxAge)

// 	// All utxns that are unchecked should be exactly the same
// 	assert.Equal(t, up.Txns[validUtxUnchecked.Hash()], validUtxUnchecked)
// 	assert.Equal(t, up.Txns[invalidUtxUnchecked.Hash()], invalidUtxUnchecked)
// 	// The valid one that is checked should have its checked status updated
// 	validUtxCheckedUpdated := up.Txns[validUtxChecked.Hash()]
// 	assert.True(t,
// 		validUtxCheckedUpdated.Checked.After(validUtxChecked.Checked))
// 	validUtxChecked.Checked = validUtxCheckedUpdated.Checked
// 	assert.Equal(t, validUtxChecked, validUtxCheckedUpdated)
// 	// The invalid checked one and the expired one should be removed
// 	_, ok := up.Txns[invalidUtxChecked.Hash()]
// 	assert.False(t, ok)
// 	_, ok = up.Txns[validUtxExpired.Hash()]
// 	assert.False(t, ok)
// 	// Also, the unspents should have 2 * nRemaining
// 	assert.Equal(t, len(up.Unspent), 3)
// 	for _, uxs := range up.Unspent {
// 		assert.Equal(t, len(uxs), 2)
// 	}
// 	assert.Equal(t, len(up.Txns), 3)
// }

// func TestRefresh(t *testing.T) {
// 	bc := makeBlockchain()
// 	testRefresh(t, bc, func(checkPeriod, maxAge time.Duration) {
// 		mv.Unconfirmed.Refresh(mv.blockchain, checkPeriod, maxAge)
// 	})
// }

// func TestGetOldOwnedTransactions(t *testing.T) {
// 	db, close := prepareDB(t)
// 	defer close()
// 	up := NewUnconfirmedTxnPool(db)
// 	bc := makeBlockchain()
// 	// Setup txns
// 	notOursNew := makeTransactionForChain(t, bc)
// 	notOursOld := makeTransactionForChain(t, bc)
// 	ourSpendNew := makeTransactionForChain(t, bc)
// 	ourSpendOld := makeTransactionForChain(t, bc)
// 	ourReceiveNew := makeTransactionForChain(t, bc)
// 	ourReceiveOld := makeTransactionForChain(t, bc)
// 	ourBothNew := makeTransactionForChain(t, bc)
// 	ourBothOld := makeTransactionForChain(t, bc)

// 	// Add a transaction that is not ours, both new and old
// 	err, known := up.InjectTxn(bc, notOursNew)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	up.SetAnnounced(notOursNew.Hash(), utc.Now())
// 	err, known = up.InjectTxn(bc, notOursOld)
// 	assert.Nil(t, err)
// 	assert.False(t, known)

// 	// Add a transaction that is our spend, both new and old
// 	addrs := make(map[cipher.Address]byte, 1)
// 	ux, ok := bc.GetUnspent().Get(ourSpendNew.In[0])
// 	assert.True(t, ok)
// 	addrs[ux.Body.Address] = byte(1)
// 	err, known = up.InjectTxn(bc, ourSpendNew)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	up.SetAnnounced(ourSpendNew.Hash(), utc.Now())
// 	addrs = make(map[cipher.Address]byte, 1)
// 	ux, ok = bc.GetUnspent().Get(ourSpendNew.In[0])
// 	assert.True(t, ok)
// 	addrs[ux.Body.Address] = byte(1)
// 	err, known = up.InjectTxn(bc, ourSpendOld)
// 	assert.Nil(t, err)
// 	assert.False(t, known)

// 	// Add a transaction that is our receive, both new and old
// 	addrs = make(map[cipher.Address]byte, 1)
// 	addrs[ourReceiveNew.Out[1].Address] = byte(1)
// 	err, known = up.InjectTxn(bc, ourReceiveNew)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	up.SetAnnounced(ourReceiveNew.Hash(), utc.Now())
// 	addrs = make(map[cipher.Address]byte, 1)
// 	addrs[ourReceiveOld.Out[1].Address] = byte(1)
// 	err, known = up.InjectTxn(bc, ourReceiveOld)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	// Add a transaction that is both our spend and receive, both new and old
// 	addrs = make(map[cipher.Address]byte, 2)
// 	ux, ok = bc.GetUnspent().Get(ourBothNew.In[0])
// 	assert.True(t, ok)
// 	addrs[ux.Body.Address] = byte(1)
// 	addrs[ourBothNew.Out[1].Address] = byte(1)
// 	assert.Equal(t, len(addrs), 2)
// 	err, known = up.InjectTxn(bc, ourBothNew)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// 	up.SetAnnounced(ourBothNew.Hash(), utc.Now())
// 	addrs = make(map[cipher.Address]byte, 1)
// 	ux, ok = bc.GetUnspent().Get(ourBothOld.In[0])
// 	assert.True(t, ok)
// 	addrs[ux.Body.Address] = byte(1)
// 	addrs[ourBothOld.Out[1].Address] = byte(1)
// 	assert.Equal(t, len(addrs), 2)
// 	err, known = up.InjectTxn(bc, ourBothOld)
// 	assert.Nil(t, err)
// 	assert.False(t, known)
// }

func TestFilterKnown(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	uts := createUnconfirmedTxns(t, up, bc, 4)
	hashes := []cipher.SHA256{
		uts[0].Hash(),
		uts[1].Hash(),
		randSHA256(),
		randSHA256(),
	}

	known := up.FilterKnown(hashes)
	assert.Equal(t, len(known), 2)
	for i, h := range known {
		assert.Equal(t, h, hashes[i+2])
	}
	_, ok := up.Txns.get(known[0])
	assert.False(t, ok)
	_, ok = up.Txns.get(known[1])
	assert.False(t, ok)
}

func TestGetKnown(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	uts := createUnconfirmedTxns(t, up, bc, 4)
	hashes := []cipher.SHA256{
		uts[0].Hash(),
		uts[1].Hash(),
		randSHA256(),
		randSHA256(),
	}

	known := up.GetKnown(hashes)
	assert.Equal(t, len(known), 2)
	for i, tx := range known {
		assert.Equal(t, tx.Hash(), hashes[i])
		assert.Equal(t, tx, uts[i].Txn)
	}
	_, ok := up.Txns.get(known[0].Hash())
	assert.True(t, ok)
	_, ok = up.Txns.get(known[1].Hash())
	assert.True(t, ok)
}

func TestTxUnspentsGetAllForAddress(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)

	addrs := make(map[cipher.Address]byte, 0)
	useAddrs := make([]cipher.Address, 4)
	for i := range useAddrs {
		useAddrs[i] = makeAddress()
	}
	useAddrs[1] = useAddrs[0]
	for _, a := range useAddrs {
		addrs[a] = byte(1)
	}
	// Make confirmed transactions to add to unspent pool
	uxs := make(coin.UxArray, 0)
	for i := 0; i < 4; i++ {
		txn := coin.Transaction{}
		txn.PushInput(randSHA256())
		txn.PushOutput(useAddrs[i], 10e6, 1000)
		uxa := coin.CreateUnspents(coin.BlockHeader{BkSeq: 1}, txn)
		assert.Nil(t, up.Unspent.put(txn.Hash(), uxa))
		uxs = append(uxs, uxa...)
	}
	assert.Equal(t, len(uxs), 4)

	uxa := up.Unspent.getAllForAddress(useAddrs[0])
	assert.Equal(t, 2, len(uxa))
	for _, u := range uxa {
		var has bool
		for _, ux := range uxs[:2] {
			if ux.Hash() == u.Hash() {
				has = true
			}
		}
		assert.True(t, has)
	}
}

func TestSpendsForAddresses(t *testing.T) {
	db, close := prepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)
	unspent := coin.NewUnspentPool()
	addrs := make(map[cipher.Address]byte, 0)
	n := 4
	useAddrs := make([]cipher.Address, n)
	for i := range useAddrs {
		useAddrs[i] = makeAddress()
	}
	useAddrs[1] = useAddrs[0]
	for _, a := range useAddrs {
		addrs[a] = byte(1)
	}
	// Make confirmed transactions to add to unspent pool
	uxs := make(coin.UxArray, 0)
	for i := 0; i < n; i++ {
		txn := coin.Transaction{}
		txn.PushInput(randSHA256())
		txn.PushOutput(useAddrs[i], 10e6, 1000)
		uxa := coin.CreateUnspents(coin.BlockHeader{BkSeq: 1}, txn)
		for _, ux := range uxa {
			unspent.Add(ux)
		}
		uxs = append(uxs, uxa...)
	}
	assert.Equal(t, len(uxs), 4)

	// Make unconfirmed txns that spend those unspents
	for i := 0; i < n; i++ {
		txn := coin.Transaction{}
		txn.PushInput(uxs[i].Hash())
		txn.PushOutput(makeAddress(), 10e6, 1000)
		ut := UnconfirmedTxn{
			Txn: txn,
		}
		up.Txns.put(&ut)
		up.Unspent.put(ut.Hash(), coin.CreateUnspents(coin.BlockHeader{BkSeq: 1}, txn))
	}

	// Now look them up
	assert.Equal(t, len(addrs), 3)
	assert.Equal(t, up.Txns.len(), 4)

	auxs := up.SpendsForAddresses(&unspent, addrs)
	assert.Equal(t, len(auxs), 3)
	assert.Equal(t, len(auxs[useAddrs[0]]), 2)
	assert.Equal(t, len(auxs[useAddrs[2]]), 1)
	assert.Equal(t, len(auxs[useAddrs[3]]), 1)
	for _, u := range uxs[:2] {
		var has bool
		for _, u1 := range auxs[useAddrs[0]] {
			if u.Hash() == u1.Hash() {
				has = true
				break
			}
		}
		assert.True(t, has)
	}
	assert.Equal(t, auxs[useAddrs[2]], coin.UxArray{uxs[2]})
	assert.Equal(t, auxs[useAddrs[3]], coin.UxArray{uxs[3]})
}

func TestUnconfirmTxBktPutAndGet(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}

	now := time.Now()
	for i := range uctxs {
		uctxs[i].Received = now.Add(time.Duration(i) * time.Minute).UnixNano()
		uctxs[i].Checked = uctxs[i].Received
		uctxs[i].Announced = uctxs[i].Received + 100
	}

	testCases := []struct {
		name  string
		init  []UnconfirmedTxn
		get   UnconfirmedTxn
		exist bool
	}{
		{
			"get success",
			uctxs[:2],
			uctxs[1],
			true,
		},
		{
			"get success",
			uctxs[:2],
			uctxs[0],
			true,
		},
		{
			"get failed",
			uctxs[:2],
			uctxs[2],
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, close := prepareDB(t)
			defer close()
			up := NewUnconfirmedTxnPool(db)
			for _, u := range tc.init {
				err := up.Txns.put(&u)
				assert.Nil(t, err)
			}
			// check
			u, ok := up.Txns.get(tc.get.Hash())
			assert.Equal(t, tc.exist, ok)
			if !ok {
				return
			}
			assert.Equal(t, tc.get, *u)
		})
	}
}

func TestUnconfirmTxBktUpdate(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}

	testCases := []struct {
		name      string
		init      []UnconfirmedTxn
		index     int   // update index
		timestamp int64 // update all time to this
		err       error
	}{
		{
			"update success",
			uctxs,
			0,
			time.Now().UnixNano(),
			nil,
		},
		{
			"update not exist",
			uctxs[:2],
			2,
			time.Now().UnixNano(),
			fmt.Errorf("%s does not exist in bucket unconfirmed_txns", uctxs[2].Hash().Hex()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, close := prepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				err := bkt.put(&u)
				assert.Nil(t, err)
			}

			// update
			err := bkt.update(uctxs[tc.index].Hash(), func(u *UnconfirmedTxn) {
				u.Announced = tc.timestamp
				u.Checked = tc.timestamp
				u.Received = tc.timestamp
			})
			assert.Equal(t, tc.err, err)

			uctxs[tc.index].Announced = tc.timestamp
			uctxs[tc.index].Received = tc.timestamp
			uctxs[tc.index].Checked = tc.timestamp

			for _, u := range tc.init {
				ux, ok := bkt.get(u.Hash())
				assert.True(t, ok)
				assert.Equal(t, u, *ux)
			}
		})
	}
}

func TestUnconfirmedBktGetAll(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}

	f := fmt.Sprintf("test%d.db", rand.Intn(1024))
	db, err := bolt.Open(f, 0700, nil)
	assert.Nil(t, err)
	defer os.Remove(f)

	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs {
		err := bkt.put(&u)
		assert.Nil(t, err)
	}

	vm, err := bkt.getAll()
	assert.Nil(t, err)
	assert.Equal(t, uctxs, vm)

	db.Close()

	db, err = bolt.Open(f, 0700, nil)
	assert.Nil(t, err)
	defer db.Close()
	bkt = newUncfmTxBkt(db)

	vm, err = bkt.getAll()
	assert.Nil(t, err)
	assert.Equal(t, uctxs, vm)
}

func TestUnconfirmedTxRangeUpdate(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		// createUnconfirmedTxn(),
		// createUnconfirmedTxn(),
	}

	testCases := []struct {
		name  string
		init  []UnconfirmedTxn
		index int
		time  int64
	}{
		{
			"range update success",
			uctxs,
			0,
			time.Now().UnixNano(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, close := prepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				assert.Nil(t, bkt.put(&u))
			}

			err := bkt.rangeUpdate(func(key cipher.SHA256, ux *UnconfirmedTxn) {
				if key == uctxs[tc.index].Hash() {
					ux.Announced = tc.time
					ux.Checked = tc.time
					ux.Received = tc.time
				}
			})
			assert.Nil(t, err)

			uctxs[tc.index].Announced = tc.time
			uctxs[tc.index].Checked = tc.time
			uctxs[tc.index].Received = tc.time

			for _, u := range uctxs {
				ux, ok := bkt.get(u.Hash())
				assert.True(t, ok)
				assert.Equal(t, u, *ux)
			}
		})
	}
}

func TestUnconfirmedTxDelete(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}

	testCases := []struct {
		name  string
		init  []UnconfirmedTxn
		index int // delete index
		exist bool
	}{
		{
			"delete exist",
			uctxs,
			0,
			false,
		},
		{
			"delete not exist",
			uctxs[:2],
			2,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, close := prepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				assert.Nil(t, bkt.put(&u))
			}

			key := uctxs[tc.index].Hash()
			assert.Nil(t, bkt.delete(key))

			_, ok := bkt.get(key)
			assert.Equal(t, tc.exist, ok)
			assert.Equal(t, tc.exist, bkt.isExist(key))
		})
	}
}

func TestUnconfirmedTxForEach(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}

	um := make(map[cipher.SHA256]UnconfirmedTxn, len(uctxs))

	db, close := prepareDB(t)
	defer close()
	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs {
		um[u.Hash()] = u
		assert.Nil(t, bkt.put(&u))
	}

	var count int
	bkt.forEach(func(k cipher.SHA256, ux *UnconfirmedTxn) error {
		assert.Equal(t, um[k], *ux)
		count++
		return nil
	})
	assert.Equal(t, len(uctxs), count)
}

func TestUnconfirmedTxLen(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}
	db, close := prepareDB(t)
	defer close()
	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs[:2] {
		assert.Nil(t, bkt.put(&u))
	}
	assert.Equal(t, len(uctxs[:2]), bkt.len())

	// add the last one
	assert.Nil(t, bkt.put(&uctxs[2]))
	assert.Equal(t, bkt.len(), len(uctxs))

	for i := 0; i < len(uctxs); i++ {
		assert.Nil(t, bkt.delete(uctxs[i].Hash()))
		assert.Equal(t, bkt.len(), len(uctxs)-1-i)
	}
}
