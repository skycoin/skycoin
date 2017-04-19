package visor

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"crypto/rand"
// 	"testing"
// 	"time"

// 	"github.com/skycoin/skycoin/src/cipher"
// 	"github.com/skycoin/skycoin/src/coin"
// 	"github.com/skycoin/skycoin/src/util"
// 	"github.com/skycoin/skycoin/src/wallet"
// 	"github.com/stretchr/testify/assert"
// )

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

func makeValidTxn() (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	uncf := NewUnconfirmedTxnPool()
	now := tNow()
	a := makeAddress()
	uxs := makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp := coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt := wallet.Balance{10 * 1e6, 0}
	return CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
}

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

func makeValidTxnNoChange() (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	uncf := NewUnconfirmedTxnPool()
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

func makeInvalidTxn() (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	uncf := NewUnconfirmedTxnPool()
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
	unspent TxnUnspents, tx coin.Transaction) {
	expect := coin.CreateUnspents(bc.Head().Head, tx)
	assert.NotEqual(t, len(expect), 0)
	sum := 0
	for _, uxs := range unspent {
		sum += len(uxs)
	}
	assert.Equal(t, len(expect), sum)
	uxs := unspent[tx.Hash()]
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

func createUnconfirmedTxns(t *testing.T, up *UnconfirmedTxnPool, n int) []UnconfirmedTxn {
	uts := make([]UnconfirmedTxn, 4)
	usp := coin.NewUnspentPool()
	for i := 0; i < len(uts); i++ {
		tx, _ := makeValidTxn()
		ut := up.createUnconfirmedTxn(&usp, tx)
		uts[i] = ut
		up.Txns[ut.Hash()] = ut
	}
	assert.Equal(t, len(up.Txns), 4)
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

// func TestUnconfirmedTxnHash(t *testing.T) {
// 	up := NewUnconfirmedTxnPool()
// 	uts := createUnconfirmedTxns(t, up, 1)
// 	utx := uts[0]
// 	assert.Equal(t, utx.Hash(), utx.Txn.Hash())
// 	assert.NotEqual(t, utx.Hash(), utx.Txn.Head.Hash)
// }

func TestNewUnconfirmedTxnPool(t *testing.T) {
	ut := NewUnconfirmedTxnPool()
	assert.NotNil(t, ut.Txns)
	assert.Equal(t, len(ut.Txns), 0)
}

func TestSetAnnounced(t *testing.T) {
	ut := NewUnconfirmedTxnPool()
	assert.Equal(t, len(ut.Txns), 0)
	// Unknown should be safe and a noop
	assert.NotPanics(t, func() {
		ut.SetAnnounced(cipher.SHA256{}, util.Now())
	})
	assert.Equal(t, len(ut.Txns), 0)
	utx := createUnconfirmedTxns(t, ut, 1)[0]
	assert.True(t, nanoToTime(utx.Announced).IsZero())
	ut.Txns[utx.Hash()] = utx
	now := util.Now()
	ut.SetAnnounced(utx.Hash(), now)
	assert.Equal(t, ut.Txns[utx.Hash()].Announced, now)
}

func TestInjectTxn(t *testing.T) {
	defer cleanupVisor()
	// Test with invalid txn
	mv := setupMasterVisor()
	ut := NewUnconfirmedTxnPool()
	txn, err := makeInvalidTxn(mv)
	assert.Nil(t, err)
	err, known := ut.InjectTxn(mv.blockchain, txn, nil, testBlockSize, 0)
	assert.NotNil(t, err)
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 0)

	// Test didAnnounce=false
	mv = setupMasterVisor()
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = ut.InjectTxn(mv.blockchain, txn, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test didAnnounce=true
	mv = setupMasterVisor()
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = ut.InjectTxn(mv.blockchain, txn, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test txn too large
	mv = setupMasterVisor()
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = ut.InjectTxn(mv.blockchain, txn, nil, txn.Size()-1, 1)
	assertError(t, err, "Transaction too large")
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 0)

	// Test where we are receiver of ux outputs
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxn(mv)
	assert.Nil(t, err)
	addrs := make(map[cipher.Address]byte, 1)
	addrs[txn.Out[1].Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test where we are spender of ux outputs
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxnNoChange(mv)
	assert.Nil(t, err)
	addrs = make(map[cipher.Address]byte, 1)
	ux, ok := mv.blockchain.Unspent.Get(txn.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test where we are both spender and receiver of ux outputs
	mv = setupMasterVisor()
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxn(mv)
	assert.Nil(t, err)
	addrs = make(map[cipher.Address]byte, 2)
	addrs[txn.Out[0].Address] = byte(1)
	ux, ok = mv.blockchain.Unspent.Get(txn.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}

	// Test duplicate Record, should be no-op besides state change
	utx := ut.Txns[txn.Hash()]
	// Set a placeholder value on the utx to check if we overwrote it
	utx.Announced = util.ZeroTime().Add(time.Minute)
	ut.Txns[txn.Hash()] = utx
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.True(t, known)
	utx2 := ut.Txns[txn.Hash()]
	assert.Equal(t, utx2.Announced, util.ZeroTime().Add(time.Minute))
	utx2.Announced = util.ZeroTime()
	ut.Txns[utx2.Hash()] = utx2
	// Received & checked should be updated
	assert.True(t, utx2.Received.After(utx.Received))
	assert.True(t, utx2.Checked.After(utx.Checked))
	assertValidUnconfirmed(t, ut.Txns, txn)
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}

	// Test with valid fee, exact
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxnWithFeeFactor(mv, 4, 0)
	assert.Nil(t, err)
	addrs = make(map[cipher.Address]byte, 1)
	addrs[txn.Out[1].Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test with valid fee, surplus
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxnWithFeeFactor(mv, 4, 100)
	assert.Nil(t, err)
	addrs = make(map[cipher.Address]byte, 1)
	addrs[txn.Out[1].Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	assert.Nil(t, err)
	assert.False(t, known)
	assertValidUnspent(t, mv.blockchain, ut.Unspent, txn)
	assertValidUnconfirmed(t, ut.Txns, txn)

	// Test with invalid fee
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxnWithFeeFactor(mv, 5, 0)
	assert.Nil(t, err)
	_, err = mv.blockchain.TransactionFee(&txn)
	assert.Nil(t, err)
	addrs = make(map[cipher.Address]byte, 1)
	addrs[txn.Out[1].Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	assertError(t, err, "Transaction fee minimum not met")
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 0)

	// Test with bc.TransactionFee failing
	mv = setupMasterVisor()
	assert.Equal(t, len(mv.blockchain.Unspent.Pool), 1)
	ut = NewUnconfirmedTxnPool()
	txn, err = makeValidTxnWithFeeFactor(mv, 4, 100)
	assert.Nil(t, err)
	txn.Out[1].Hours = 1e16
	txn.Head.Sigs = make([]cipher.Sig, 0)
	txn.SignInputs([]cipher.SecKey{mv.Config.MasterKeys.Secret})
	txn.UpdateHeader()
	addrs = make(map[cipher.Address]byte, 1)
	addrs[txn.Out[1].Address] = byte(1)
	err, known = ut.InjectTxn(mv.blockchain, txn, addrs, testBlockSize, 4)
	assertError(t, err, "Insufficient coinhours for transaction outputs")
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 0)
}

func TestRawTxns(t *testing.T) {
	ut := NewUnconfirmedTxnPool()
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
	defer cleanupVisor()
	mv := setupMasterVisor()
	ut := NewUnconfirmedTxnPool()

	utx, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known := ut.InjectTxn(mv.blockchain, utx, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}

	// Unknown txn is no-op
	badh := randSHA256()
	assert.NotEqual(t, badh, utx.Hash())
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	ut.removeTxn(mv.blockchain, badh)
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}

	// Known txn updates Txns, predicted Unspents
	utx2, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = ut.InjectTxn(mv.blockchain, utx2, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assert.Equal(t, len(ut.Txns), 2)
	assert.Equal(t, len(ut.Unspent), 2)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	ut.removeTxn(mv.blockchain, utx.Hash())
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	ut.removeTxn(mv.blockchain, utx.Hash())
	assert.Equal(t, len(ut.Txns), 1)
	assert.Equal(t, len(ut.Unspent), 1)
	for _, uxs := range ut.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	ut.removeTxn(mv.blockchain, utx2.Hash())
	assert.Equal(t, len(ut.Txns), 0)
	assert.Equal(t, len(ut.Unspent), 0)
}

func TestRemoveTxns(t *testing.T) {
	defer cleanupVisor()
	mv := setupMasterVisor()
	up := NewUnconfirmedTxnPool()

	// Include an unknown hash, and omit a known hash. The other two should
	// be removed.
	hashes := make([]cipher.SHA256, 0, 3)
	hashes = append(hashes, randSHA256()) // unknown hash
	ut, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known := up.InjectTxn(mv.blockchain, ut, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	hashes = append(hashes, ut.Hash())
	ut2, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = up.InjectTxn(mv.blockchain, ut2, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	hashes = append(hashes, ut2.Hash())
	ut3, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = up.InjectTxn(mv.blockchain, ut3, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)

	assert.Equal(t, len(up.Unspent), 3)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 3)
	up.removeTxns(mv.blockchain, hashes)
	assert.Equal(t, len(up.Unspent), 1)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 1)
	_, ok := up.Txns[ut3.Hash()]
	assert.True(t, ok)
}

func TestRemoveTransactions(t *testing.T) {
	defer cleanupVisor()
	mv := setupMasterVisor()
	up := NewUnconfirmedTxnPool()

	// Include an unknown txn, and omit a known hash. The other two should
	// be removed.
	unkUt, err := makeValidTxn(mv)
	assert.Nil(t, err)
	txns := make(coin.Transactions, 0, 3)
	txns = append(txns, unkUt) // unknown txn
	ut, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known := up.InjectTxn(mv.blockchain, ut, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	txns = append(txns, ut)
	ut2, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = up.InjectTxn(mv.blockchain, ut2, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	txns = append(txns, ut2)
	ut3, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known = up.InjectTxn(mv.blockchain, ut3, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)

	assert.Equal(t, len(up.Unspent), 3)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 3)
	up.RemoveTransactions(mv.blockchain, txns)
	assert.Equal(t, len(up.Unspent), 1)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 1)
	_, ok := up.Txns[ut3.Hash()]
	assert.True(t, ok)
}

func testRefresh(t *testing.T, mv *Visor,
	refresh func(checkPeriod, maxAge time.Duration)) {
	up := mv.Unconfirmed
	// Add a transaction that is invalid, but will not be checked yet
	// Add a transaction that is invalid, and will be checked and removed
	invalidTxUnchecked, err := makeValidTxn(mv)
	assert.Nil(t, err)
	invalidTxChecked, err := makeValidTxn(mv)
	assert.Nil(t, err)
	assert.Nil(t, invalidTxUnchecked.Verify())
	assert.Nil(t, invalidTxChecked.Verify())
	// Invalidate it by spending the output that this txn references
	invalidator, err := makeValidTxn(mv)
	assert.Nil(t, err)
	err, known := up.InjectTxn(mv.blockchain, invalidator, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	assert.Equal(t, len(up.Txns), 1)
	_, err = mv.CreateAndExecuteBlock()
	assert.Nil(t, err)
	assert.Equal(t, len(up.Txns), 0)
	assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxUnchecked))
	assert.NotNil(t, mv.blockchain.VerifyTransaction(invalidTxChecked))

	invalidUtxUnchecked := UnconfirmedTxn{
		Txn:       invalidTxUnchecked,
		Received:  util.Now(),
		Checked:   util.Now(),
		Announced: util.ZeroTime(),
	}
	invalidUtxChecked := invalidUtxUnchecked
	invalidUtxChecked.Txn = invalidTxChecked
	invalidUtxUnchecked.Checked = util.Now().Add(time.Hour)
	invalidUtxChecked.Checked = util.Now().Add(-time.Hour)
	up.Txns[invalidUtxUnchecked.Hash()] = invalidUtxUnchecked
	up.Txns[invalidUtxChecked.Hash()] = invalidUtxChecked
	assert.Equal(t, len(up.Txns), 2)
	uncheckedHash := invalidTxUnchecked.Hash()
	checkedHash := invalidTxChecked.Hash()
	up.Unspent[uncheckedHash] = coin.CreateUnspents(coin.BlockHeader{},
		invalidTxUnchecked)
	up.Unspent[checkedHash] = coin.CreateUnspents(coin.BlockHeader{},
		invalidTxChecked)

	// Create a transaction that is valid, and will not be checked yet
	validTxUnchecked, err := makeValidTxn(mv)
	assert.Nil(t, err)
	// Create a transaction that is valid, and will be checked
	validTxChecked, err := makeValidTxn(mv)
	assert.Nil(t, err)
	// Create a transaction that is expired
	validTxExpired, err := makeValidTxn(mv)
	assert.Nil(t, err)

	// Add the transaction that is valid, and will not be checked yet
	err, known = up.InjectTxn(mv.blockchain, validTxUnchecked, nil,
		testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	validUtxUnchecked := up.Txns[validTxUnchecked.Hash()]
	validUtxUnchecked.Checked = util.Now().Add(time.Hour)
	up.Txns[validUtxUnchecked.Hash()] = validUtxUnchecked
	assert.Equal(t, len(up.Txns), 3)

	// Add the transaction that is valid, and will be checked
	err, known = up.InjectTxn(mv.blockchain, validTxChecked, nil,
		testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	validUtxChecked := up.Txns[validTxChecked.Hash()]
	validUtxChecked.Checked = util.Now().Add(-time.Hour)
	up.Txns[validUtxChecked.Hash()] = validUtxChecked
	assert.Equal(t, len(up.Txns), 4)

	// Add the transaction that is expired
	err, known = up.InjectTxn(mv.blockchain, validTxExpired, nil,
		testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	validUtxExpired := up.Txns[validTxExpired.Hash()]
	validUtxExpired.Received = util.Now().Add(-time.Hour)
	up.Txns[validTxExpired.Hash()] = validUtxExpired
	assert.Equal(t, len(up.Txns), 5)

	// Pre-sanity check
	assert.Equal(t, len(up.Unspent), 5)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 5)

	// Refresh
	checkPeriod := time.Second * 2
	maxAge := time.Second * 4
	refresh(checkPeriod, maxAge)

	// All utxns that are unchecked should be exactly the same
	assert.Equal(t, up.Txns[validUtxUnchecked.Hash()], validUtxUnchecked)
	assert.Equal(t, up.Txns[invalidUtxUnchecked.Hash()], invalidUtxUnchecked)
	// The valid one that is checked should have its checked status updated
	validUtxCheckedUpdated := up.Txns[validUtxChecked.Hash()]
	assert.True(t,
		validUtxCheckedUpdated.Checked.After(validUtxChecked.Checked))
	validUtxChecked.Checked = validUtxCheckedUpdated.Checked
	assert.Equal(t, validUtxChecked, validUtxCheckedUpdated)
	// The invalid checked one and the expired one should be removed
	_, ok := up.Txns[invalidUtxChecked.Hash()]
	assert.False(t, ok)
	_, ok = up.Txns[validUtxExpired.Hash()]
	assert.False(t, ok)
	// Also, the unspents should have 2 * nRemaining
	assert.Equal(t, len(up.Unspent), 3)
	for _, uxs := range up.Unspent {
		assert.Equal(t, len(uxs), 2)
	}
	assert.Equal(t, len(up.Txns), 3)
}

func TestRefresh(t *testing.T) {
	defer cleanupVisor()
	mv := setupMasterVisor()
	testRefresh(t, mv, func(checkPeriod, maxAge time.Duration) {
		mv.Unconfirmed.Refresh(mv.blockchain, checkPeriod, maxAge)
	})
}

func TestGetOldOwnedTransactions(t *testing.T) {
	mv := setupMasterVisor()
	up := mv.Unconfirmed

	// Setup txns
	notOursNew, err := makeValidTxn(mv)
	assert.Nil(t, err)
	notOursOld, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourSpendNew, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourSpendOld, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourReceiveNew, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourReceiveOld, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourBothNew, err := makeValidTxn(mv)
	assert.Nil(t, err)
	ourBothOld, err := makeValidTxn(mv)
	assert.Nil(t, err)

	// Add a transaction that is not ours, both new and old
	err, known := up.InjectTxn(mv.blockchain, notOursNew, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	up.SetAnnounced(notOursNew.Hash(), util.Now())
	err, known = up.InjectTxn(mv.blockchain, notOursOld, nil, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)

	// Add a transaction that is our spend, both new and old
	addrs := make(map[cipher.Address]byte, 1)
	ux, ok := mv.blockchain.Unspent.Get(ourSpendNew.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	err, known = up.InjectTxn(mv.blockchain, ourSpendNew, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	up.SetAnnounced(ourSpendNew.Hash(), util.Now())
	addrs = make(map[cipher.Address]byte, 1)
	ux, ok = mv.blockchain.Unspent.Get(ourSpendNew.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	err, known = up.InjectTxn(mv.blockchain, ourSpendOld, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)

	// Add a transaction that is our receive, both new and old
	addrs = make(map[cipher.Address]byte, 1)
	addrs[ourReceiveNew.Out[1].Address] = byte(1)
	err, known = up.InjectTxn(mv.blockchain, ourReceiveNew, addrs,
		testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	up.SetAnnounced(ourReceiveNew.Hash(), util.Now())
	addrs = make(map[cipher.Address]byte, 1)
	addrs[ourReceiveOld.Out[1].Address] = byte(1)
	err, known = up.InjectTxn(mv.blockchain, ourReceiveOld, addrs,
		testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	// Add a transaction that is both our spend and receive, both new and old
	addrs = make(map[cipher.Address]byte, 2)
	ux, ok = mv.blockchain.Unspent.Get(ourBothNew.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	addrs[ourBothNew.Out[1].Address] = byte(1)
	assert.Equal(t, len(addrs), 2)
	err, known = up.InjectTxn(mv.blockchain, ourBothNew, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
	up.SetAnnounced(ourBothNew.Hash(), util.Now())
	addrs = make(map[cipher.Address]byte, 1)
	ux, ok = mv.blockchain.Unspent.Get(ourBothOld.In[0])
	assert.True(t, ok)
	addrs[ux.Body.Address] = byte(1)
	addrs[ourBothOld.Out[1].Address] = byte(1)
	assert.Equal(t, len(addrs), 2)
	err, known = up.InjectTxn(mv.blockchain, ourBothOld, addrs, testBlockSize, 0)
	assert.Nil(t, err)
	assert.False(t, known)
}

func TestFilterKnown(t *testing.T) {
	up := NewUnconfirmedTxnPool()
	uts := createUnconfirmedTxns(t, up, 4)
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
	_, ok := up.Txns[known[0]]
	assert.False(t, ok)
	_, ok = up.Txns[known[1]]
	assert.False(t, ok)
}

func TestGetKnown(t *testing.T) {
	up := NewUnconfirmedTxnPool()
	uts := createUnconfirmedTxns(t, up, 4)
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
	_, ok := up.Txns[known[0].Hash()]
	assert.True(t, ok)
	_, ok = up.Txns[known[1].Hash()]
	assert.True(t, ok)
}

func TestSpendsForAddresses(t *testing.T) {
	up := NewUnconfirmedTxnPool()
	unspent := coin.NewUnspentPool()
	addrs := make(map[cipher.Address]byte, 0)
	n := 4
	useAddrs := make([]cipher.Address, n)
	for i, _ := range useAddrs {
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
		uxa := coin.CreateUnspents(coin.BlockHeader{}, txn)
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
		up.Txns[ut.Hash()] = ut
	}

	// Now look them up
	assert.Equal(t, len(addrs), 3)
	assert.Equal(t, len(up.Txns), 4)
	auxs := up.SpendsForAddresses(&unspent, addrs)
	assert.Equal(t, len(auxs), 3)
	assert.Equal(t, len(auxs[useAddrs[0]]), 2)
	assert.Equal(t, len(auxs[useAddrs[2]]), 1)
	assert.Equal(t, len(auxs[useAddrs[3]]), 1)
	assert.Equal(t, auxs[useAddrs[0]], coin.UxArray{uxs[0], uxs[1]})
	assert.Equal(t, auxs[useAddrs[2]], coin.UxArray{uxs[2]})
	assert.Equal(t, auxs[useAddrs[3]], coin.UxArray{uxs[3]})
}
