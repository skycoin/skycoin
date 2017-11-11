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
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/wallet"
)

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
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

func getFee(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func makeValidTxnNoChange(t *testing.T) (coin.Transaction, error) {
	w := wallet.NewWallet("test")
	w.GenerateAddresses(2)
	db, close := testutil.PrepareDB(t)
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
	db, close := testutil.PrepareDB(t)
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

func requireValidUnspent(t *testing.T, bc *Blockchain,
	unspent *txUnspents, tx coin.Transaction) {
	expect := coin.CreateUnspents(bc.Head().Head, tx)
	require.NotEqual(t, len(expect), 0)
	sum := 0

	unspent.forEach(func(_ cipher.SHA256, uxo coin.UxArray) {
		sum += len(uxo)
	})

	require.Equal(t, len(expect), sum)
	uxs, err := unspent.get(tx.Hash())
	require.NoError(t, err)
	for _, ux := range expect {
		found := false
		for _, u := range uxs {
			if u.Hash() == ux.Hash() {
				found = true
				break
			}
		}
		require.True(t, found)
	}
}

func requireValidUnconfirmed(t *testing.T, txns map[cipher.SHA256]UnconfirmedTxn,
	txn coin.Transaction) {
	ut, ok := txns[txn.Hash()]
	require.True(t, ok)
	require.Equal(t, ut.Txn, txn)
	require.True(t, nanoToTime(ut.Announced).IsZero())
	require.False(t, nanoToTime(ut.Received).IsZero())
	require.False(t, nanoToTime(ut.Checked).IsZero())
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
	require.Equal(t, up.Txns.len(), 4)
	return uts
}

func TestUnconfirmedTxnHash(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	uts := createUnconfirmedTxns(t, up, bc, 1)
	utx := uts[0]
	require.Equal(t, utx.Hash(), utx.Txn.Hash())
}

func TestNewUnconfirmedTxnPool(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	require.NotNil(t, ut.Txns)
	require.Equal(t, ut.Txns.len(), 0)
}

func TestSetAnnounced(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	require.Equal(t, ut.Txns.len(), 0)
	// Unknown should be safe and a noop
	require.NotPanics(t, func() {
		ut.SetAnnounced(cipher.SHA256{}, utc.Now())
	})
	require.Equal(t, ut.Txns.len(), 0)
	bc := makeBlockchain()
	utx := createUnconfirmedTxns(t, ut, bc, 1)[0]
	require.True(t, nanoToTime(utx.Announced).IsZero())
	ut.Txns.put(&utx)
	now := utc.Now()
	ut.SetAnnounced(utx.Hash(), now)
	v, ok := ut.Txns.get(utx.Hash())
	require.True(t, ok)
	require.Equal(t, v.Announced, now.UnixNano())
}

func makeBlockchain() *Blockchain {
	ft := FakeTree{}
	b := NewBlockchain(&ft, nil)
	b.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	return b
}

func TestInjectTxn(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)
	bc := makeBlockchain()
	txn := makeTransactionForChain(t, bc)
	known, err := ut.InjectTransaction(bc, txn)
	require.NoError(t, err)

	require.False(t, known)
	require.Equal(t, ut.Txns.len(), 1)

	requireValidUnspent(t, bc, ut.Unspent, txn)
	allUncfmTxs, err := ut.Txns.getAll()
	require.NoError(t, err)
	uncfmMap := make(map[cipher.SHA256]UnconfirmedTxn, len(allUncfmTxs))
	for _, txn := range allUncfmTxs {
		uncfmMap[txn.Hash()] = txn
	}
	requireValidUnconfirmed(t, uncfmMap, txn)

	// Test duplicate Record, should be no-op besides state change
	utx, ok := ut.Txns.get(txn.Hash())
	require.True(t, ok)
	// Set a placeholder value on the utx to check if we overwrote it
	utx.Announced = time.Time{}.Add(time.Minute).UnixNano()
	ut.Txns.put(utx)
	known, err = ut.InjectTransaction(bc, txn)
	require.NoError(t, err)
	require.True(t, known)
	utx2, ok := ut.Txns.get(txn.Hash())
	require.True(t, ok)
	require.Equal(t, utx2.Announced, time.Time{}.Add(time.Minute).UnixNano())
	utx2.Announced = time.Time{}.UnixNano()
	ut.Txns.put(utx2)
	// Received & checked should be updated
	require.True(t, nanoToTime(utx2.Received).After(nanoToTime(utx.Received)))
	require.True(t, nanoToTime(utx2.Checked).After(nanoToTime(utx.Checked)))
	all, err := ut.Txns.getAll()
	require.NoError(t, err)
	allMap := make(map[cipher.SHA256]UnconfirmedTxn, len(all))
	for _, tx := range all {
		allMap[tx.Hash()] = tx
	}
	requireValidUnconfirmed(t, allMap, txn)
	require.Equal(t, ut.Txns.len(), 1)
	require.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxa coin.UxArray) {
		require.Equal(t, len(uxa), 2)
	})
}

func TestRawTxns(t *testing.T) {
	db, close := testutil.PrepareDB(t)
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
		require.Equal(t, utxs[i], tx)
	}
}

func TestRemoveTxn(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	ut := NewUnconfirmedTxnPool(db)

	bc := makeBlockchain()
	utx := makeTransactionForChain(t, bc)
	known, err := ut.InjectTransaction(bc, utx)
	require.NoError(t, err)
	require.False(t, known)
	require.Equal(t, ut.Txns.len(), 1)
	require.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxa coin.UxArray) {
		require.Equal(t, len(uxa), 2)
	})

	// Unknown txn is no-op
	badh := randSHA256()
	require.NotEqual(t, badh, utx.Hash())
	require.Equal(t, ut.Txns.len(), 1)
	require.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	ut.removeTxn(bc, badh)
	require.Equal(t, ut.Txns.len(), 1)
	require.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})

	// Known txn updates Txns, predicted Unspents
	utx2 := makeTransactionForChain(t, bc)
	require.NoError(t, err)
	known, err = ut.InjectTransaction(bc, utx2)
	require.NoError(t, err)
	require.False(t, known)
	require.Equal(t, ut.Txns.len(), 2)
	require.Equal(t, ut.Unspent.len(), 2)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	ut.removeTxn(bc, utx.Hash())
	require.Equal(t, ut.Txns.len(), 1)
	require.Equal(t, ut.Unspent.len(), 1)
	ut.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})

	ut.removeTxn(bc, utx2.Hash())
	require.Equal(t, ut.Len(), 0)
	require.Equal(t, ut.Unspent.len(), 0)
}

func TestRemoveTxns(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)

	// Include an unknown hash, and omit a known hash. The other two should
	// be removed.
	hashes := make([]cipher.SHA256, 0, 3)
	hashes = append(hashes, randSHA256()) // unknown hash
	bc := makeBlockchain()
	ut := makeTransactionForChain(t, bc)
	known, err := up.InjectTransaction(bc, ut)
	require.NoError(t, err)
	require.False(t, known)
	hashes = append(hashes, ut.Hash())
	ut2 := makeTransactionForChain(t, bc)
	known, err = up.InjectTransaction(bc, ut2)
	require.False(t, known)
	hashes = append(hashes, ut2.Hash())
	ut3 := makeTransactionForChain(t, bc)
	known, err = up.InjectTransaction(bc, ut3)
	require.False(t, known)

	require.Equal(t, up.Unspent.len(), 3)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	require.Equal(t, up.Len(), 3)
	up.removeTxns(bc, hashes)
	require.Equal(t, up.Unspent.len(), 1)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	require.Equal(t, up.Len(), 1)
	_, ok := up.Txns.get(ut3.Hash())
	require.True(t, ok)
}

func TestRemoveTransactions(t *testing.T) {
	bc := makeBlockchain()
	db, close := testutil.PrepareDB(t)
	defer close()
	up := NewUnconfirmedTxnPool(db)

	// Include an unknown txn, and omit a known hash. The other two should
	// be removed.
	unkUt := makeTransactionForChain(t, bc)
	txns := make(coin.Transactions, 0, 3)
	txns = append(txns, unkUt) // unknown txn
	ut := makeTransactionForChain(t, bc)
	known, err := up.InjectTransaction(bc, ut)
	require.NoError(t, err)
	require.False(t, known)
	txns = append(txns, ut)
	ut2 := makeTransactionForChain(t, bc)
	require.NoError(t, err)
	known, err = up.InjectTransaction(bc, ut2)
	require.NoError(t, err)
	require.False(t, known)
	txns = append(txns, ut2)
	ut3 := makeTransactionForChain(t, bc)
	known, err = up.InjectTransaction(bc, ut3)
	require.NoError(t, err)
	require.False(t, known)

	require.Equal(t, up.Unspent.len(), 3)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	require.Equal(t, up.Txns.len(), 3)
	up.RemoveTransactions(bc, txns)
	require.Equal(t, up.Unspent.len(), 1)
	up.Unspent.forEach(func(_ cipher.SHA256, uxs coin.UxArray) {
		require.Equal(t, len(uxs), 2)
	})
	require.Equal(t, up.Txns.len(), 1)
	_, ok := up.Txns.get(ut3.Hash())
	require.True(t, ok)
}

func TestGetUnknown(t *testing.T) {
	db, close := testutil.PrepareDB(t)
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

	known, err := up.GetUnknown(hashes)
	require.NoError(t, err)
	require.Equal(t, len(known), 2)
	for i, h := range known {
		require.Equal(t, h, hashes[i+2])
	}
	_, ok := up.Txns.get(known[0])
	require.False(t, ok)
	_, ok = up.Txns.get(known[1])
	require.False(t, ok)
}

func TestGetKnown(t *testing.T) {
	db, close := testutil.PrepareDB(t)
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
	require.Equal(t, len(known), 2)
	for i, tx := range known {
		require.Equal(t, tx.Hash(), hashes[i])
		require.Equal(t, tx, uts[i].Txn)
	}
	_, ok := up.Txns.get(known[0].Hash())
	require.True(t, ok)
	_, ok = up.Txns.get(known[1].Hash())
	require.True(t, ok)
}

func TestTxUnspentsGetAllForAddress(t *testing.T) {
	db, close := testutil.PrepareDB(t)
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
		require.Nil(t, up.Unspent.put(txn.Hash(), uxa))
		uxs = append(uxs, uxa...)
	}
	require.Equal(t, len(uxs), 4)

	uxa := up.Unspent.getAllForAddress(useAddrs[0])
	require.Equal(t, 2, len(uxa))
	for _, u := range uxa {
		var has bool
		for _, ux := range uxs[:2] {
			if ux.Hash() == u.Hash() {
				has = true
			}
		}
		require.True(t, has)
	}
}

func TestSpendsForAddresses(t *testing.T) {
	db, close := testutil.PrepareDB(t)
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
	require.Equal(t, len(uxs), 4)

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
	require.Equal(t, len(addrs), 3)
	require.Equal(t, up.Txns.len(), 4)

	auxs := up.SpendsForAddresses(&unspent, addrs)
	require.Equal(t, len(auxs), 3)
	require.Equal(t, len(auxs[useAddrs[0]]), 2)
	require.Equal(t, len(auxs[useAddrs[2]]), 1)
	require.Equal(t, len(auxs[useAddrs[3]]), 1)
	for _, u := range uxs[:2] {
		var has bool
		for _, u1 := range auxs[useAddrs[0]] {
			if u.Hash() == u1.Hash() {
				has = true
				break
			}
		}
		require.True(t, has)
	}
	require.Equal(t, auxs[useAddrs[2]], coin.UxArray{uxs[2]})
	require.Equal(t, auxs[useAddrs[3]], coin.UxArray{uxs[3]})
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
			db, close := testutil.PrepareDB(t)
			defer close()
			up := NewUnconfirmedTxnPool(db)
			for _, u := range tc.init {
				err := up.Txns.put(&u)
				require.NoError(t, err)
			}
			// check
			u, ok := up.Txns.get(tc.get.Hash())
			require.Equal(t, tc.exist, ok)
			if !ok {
				return
			}
			require.Equal(t, tc.get, *u)
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
			db, close := testutil.PrepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				err := bkt.put(&u)
				require.NoError(t, err)
			}

			// update
			err := bkt.update(uctxs[tc.index].Hash(), func(u *UnconfirmedTxn) {
				u.Announced = tc.timestamp
				u.Checked = tc.timestamp
				u.Received = tc.timestamp
			})
			require.Equal(t, tc.err, err)

			uctxs[tc.index].Announced = tc.timestamp
			uctxs[tc.index].Received = tc.timestamp
			uctxs[tc.index].Checked = tc.timestamp

			for _, u := range tc.init {
				ux, ok := bkt.get(u.Hash())
				require.True(t, ok)
				require.Equal(t, u, *ux)
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
	require.NoError(t, err)
	defer os.Remove(f)

	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs {
		err := bkt.put(&u)
		require.NoError(t, err)
	}

	vm, err := bkt.getAll()
	require.NoError(t, err)
	require.Equal(t, uctxs, vm)

	db.Close()

	db, err = bolt.Open(f, 0700, nil)
	require.NoError(t, err)
	defer db.Close()
	bkt = newUncfmTxBkt(db)

	vm, err = bkt.getAll()
	require.NoError(t, err)
	require.Equal(t, uctxs, vm)
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
			db, close := testutil.PrepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				require.Nil(t, bkt.put(&u))
			}

			err := bkt.rangeUpdate(func(key cipher.SHA256, ux *UnconfirmedTxn) {
				if key == uctxs[tc.index].Hash() {
					ux.Announced = tc.time
					ux.Checked = tc.time
					ux.Received = tc.time
				}
			})
			require.NoError(t, err)

			uctxs[tc.index].Announced = tc.time
			uctxs[tc.index].Checked = tc.time
			uctxs[tc.index].Received = tc.time

			for _, u := range uctxs {
				ux, ok := bkt.get(u.Hash())
				require.True(t, ok)
				require.Equal(t, u, *ux)
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
			db, close := testutil.PrepareDB(t)
			defer close()
			bkt := newUncfmTxBkt(db)
			for _, u := range tc.init {
				require.Nil(t, bkt.put(&u))
			}

			key := uctxs[tc.index].Hash()
			require.Nil(t, bkt.delete(key))

			_, ok := bkt.get(key)
			require.Equal(t, tc.exist, ok)
			require.Equal(t, tc.exist, bkt.isExist(key))
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

	db, close := testutil.PrepareDB(t)
	defer close()
	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs {
		um[u.Hash()] = u
		require.Nil(t, bkt.put(&u))
	}

	var count int
	bkt.forEach(func(k cipher.SHA256, ux *UnconfirmedTxn) error {
		require.Equal(t, um[k], *ux)
		count++
		return nil
	})
	require.Equal(t, len(uctxs), count)
}

func TestUnconfirmedTxLen(t *testing.T) {
	uctxs := []UnconfirmedTxn{
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
		createUnconfirmedTxn(),
	}
	db, close := testutil.PrepareDB(t)
	defer close()
	bkt := newUncfmTxBkt(db)
	for _, u := range uctxs[:2] {
		require.Nil(t, bkt.put(&u))
	}
	require.Equal(t, len(uctxs[:2]), bkt.len())

	// add the last one
	require.Nil(t, bkt.put(&uctxs[2]))
	require.Equal(t, bkt.len(), len(uctxs))

	for i := 0; i < len(uctxs); i++ {
		require.Nil(t, bkt.delete(uctxs[i].Hash()))
		require.Equal(t, bkt.len(), len(uctxs)-1-i)
	}
}
