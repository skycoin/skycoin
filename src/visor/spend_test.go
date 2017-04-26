package visor

import (
	"bytes"
	"log"
	"sort"
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
)

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeUxBalances(b []wallet.Balance, headTime uint64) coin.UxArray {
	uxs := make(coin.UxArray, len(b))
	for i := range b {
		uxs[i] = coin.UxOut{
			Head: coin.UxHead{
				Time: headTime,
			},
			Body: coin.UxBody{
				SrcTransaction: randSHA256(),
				Address:        makeAddress(),
				Coins:          b[i].Coins,
				Hours:          b[i].Hours,
			},
		}
	}
	return uxs
}

func makeUxBalancesForAddresses(b []wallet.Balance, headTime uint64,
	addrs []cipher.Address) coin.UxArray {
	if len(b) != len(addrs) {
		log.Panic("Need as many addresses and balances")
	}
	uxs := makeUxBalances(b, headTime)
	for i := range uxs {
		uxs[i].Head.BkSeq = uint64(i)
		uxs[i].Body.Address = addrs[i]
	}
	return uxs
}

func makeUxOut(t *testing.T) coin.UxOut {
	return coin.UxOut{
		Head: coin.UxHead{
			BkSeq: 1,
			Time:  tNow(),
		},
		Body: coin.UxBody{
			SrcTransaction: randSHA256(),
			Address:        makeAddress(),
			Coins:          1e6,
			Hours:          1024,
		},
	}
}

func makeUxArray(t *testing.T, n int) coin.UxArray {
	uxa := make(coin.UxArray, n)
	for i := range uxa {
		uxa[i] = makeUxOut(t)
	}
	return uxa
}

func addUxArrayToUnspentPool(u *coin.UnspentPool, uxs coin.UxArray) {
	for _, ux := range uxs {
		u.Add(ux)
	}
}

func TestOldestUxOut(t *testing.T) {
	uxs := OldestUxOut(makeUxArray(t, 4))
	for i := range uxs {
		uxs[i].Head.BkSeq = uint64(i)
	}
	assert.True(t, sort.IsSorted(uxs))
	assert.Equal(t, uxs.Len(), 4)

	uxs.Swap(0, 1)
	assert.False(t, sort.IsSorted(uxs))
	assert.Equal(t, uxs[0].Head.BkSeq, uint64(1))
	assert.Equal(t, uxs[1].Head.BkSeq, uint64(0))
	uxs.Swap(0, 1)
	assert.True(t, sort.IsSorted(uxs))
	assert.Equal(t, uxs[0].Head.BkSeq, uint64(0))
	assert.Equal(t, uxs[1].Head.BkSeq, uint64(1))

	// Test hash sorting
	uxs[1].Head.BkSeq = uint64(0)
	h0 := uxs[0].Hash()
	h1 := uxs[1].Hash()
	firstLower := bytes.Compare(h0[:], h1[:]) < 0
	if firstLower {
		uxs.Swap(0, 1)
	}
	assert.False(t, sort.IsSorted(uxs))
	sort.Sort(uxs)

	cmpHash := false
	cmpSeq := false
	for i := range uxs[:len(uxs)-1] {
		j := i + 1
		if uxs[i].Head.BkSeq == uxs[j].Head.BkSeq {
			ih := uxs[i].Hash()
			jh := uxs[j].Hash()
			assert.True(t, bytes.Compare(ih[:], jh[:]) < 0)
			cmpHash = true
		} else {
			assert.True(t, uxs[i].Head.BkSeq < uxs[j].Head.BkSeq)
			cmpSeq = true
		}
	}
	assert.True(t, cmpHash)
	assert.True(t, cmpSeq)

	// Duplicate output panics
	uxs = append(uxs, uxs[0])
	assert.Panics(t, func() { sort.Sort(uxs) })
}

func TestCreateSpendsNotEnoughCoins(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{Coins: 10e6, Hours: 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{Coins: 1e6, Hours: 100},
		wallet.Balance{Coins: 8e6, Hours: 0},
	}, now)
	_, err := createSpends(now, uxs, amt)
	assertError(t, err, "Not enough coins")
}

func TestBadSpending(t *testing.T) {
	_, err := createSpends(tNow(), coin.UxArray{},
		wallet.Balance{Coins: 1e6 + 1, Hours: 1000})
	assertError(t, err, "Coins must be multiple of 1e6")
	_, err = createSpends(tNow(), coin.UxArray{},
		wallet.Balance{Coins: 0, Hours: 100})
	assertError(t, err, "Zero spend amount")
}

func TestCreateSpends(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{12e6, 125}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{1e6, 50},
		wallet.Balance{8e6, 10}, // 3
		wallet.Balance{2e6, 80}, // 2
		wallet.Balance{5e6, 15}, // 4
		wallet.Balance{7e6, 20}, // 1
	}, now)
	uxs[4].Head.BkSeq = uint64(1)
	uxs[3].Head.BkSeq = uint64(4)
	uxs[2].Head.BkSeq = uint64(2)
	uxs[1].Head.BkSeq = uint64(3)
	uxs[0].Head.BkSeq = uint64(5)
	if sort.IsSorted(OldestUxOut(uxs)) {
		uxs[0], uxs[1] = uxs[1], uxs[0]
	}
	assert.False(t, sort.IsSorted(OldestUxOut(uxs)))
	expectedSorting := coin.UxArray{uxs[4], uxs[2], uxs[1], uxs[3], uxs[0]}
	cuxs := append(coin.UxArray{}, uxs...)
	sort.Sort(OldestUxOut(cuxs))
	assert.Equal(t, expectedSorting, cuxs)
	assert.True(t, sort.IsSorted(OldestUxOut(cuxs)))
	assert.False(t, sort.IsSorted(OldestUxOut(uxs)))

	ouxs := append(coin.UxArray{}, uxs...)
	spends, err := createSpends(now, uxs, amt)
	assert.True(t, sort.IsSorted(OldestUxOut(uxs)))
	assert.Nil(t, err)
	assert.Equal(t, spends, cuxs[:len(spends)])
	assert.Equal(t, len(spends), 5)
	// assert.Equal(t, spends, coin.UxArray{ouxs[4], ouxs[2], ouxs[1], ouxs[3]})

	// Recalculate what it should be
	b := wallet.Balance{0, 0}
	ouxs = make(coin.UxArray, 0, len(spends))
	for _, ux := range cuxs {
		if b.Coins > amt.Coins && b.Hours >= amt.Hours {
			break
		}
		b = b.Add(wallet.NewBalanceFromUxOut(now, &ux))
		ouxs = append(ouxs, ux)
	}
	// assert.Equal(t, len(ouxs), len(spends))
	// assert.Equal(t, ouxs, spends)
}

func TestCreateSpendingTransaction(t *testing.T) {
	// Setup
	db, close := prepareDB(t)
	defer close()
	w := wallet.NewWallet("fortest.wlt")

	w.GenerateAddresses(4)
	uncf := NewUnconfirmedTxnPool(db)
	now := tNow()
	a := makeAddress()

	// Failing createSpends
	amt := wallet.Balance{0, 0}
	unsp := coin.NewUnspentPool()
	_, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	assert.NotNil(t, err)

	// Valid txn, fee, no change
	uxs := makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	tx, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	assert.Nil(t, err)
	assert.Equal(t, len(tx.Out), 1)
	assert.Equal(t, tx.Out[0], coin.TransactionOutput{
		Coins:   25e6,
		Hours:   37,
		Address: a,
	})
	assert.Equal(t, len(tx.In), 2)
	assert.Equal(t, tx.In, []cipher.SHA256{uxs[0].Hash(), uxs[1].Hash()})
	assert.Nil(t, tx.Verify())

	// Valid txn, change
	uxs = makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 200},
		wallet.Balance{1e6, 125},
	}, now, w.GetAddresses()[:3])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	assert.Nil(t, err)
	assert.Equal(t, len(tx.Out), 2)
	assert.Equal(t, tx.Out[0], coin.TransactionOutput{
		Coins:   1e6,
		Hours:   (150 + 200 + 125) / 8,
		Address: w.GetAddresses()[0],
	})
	assert.Equal(t, tx.Out[1], coin.TransactionOutput{
		Coins:   25e6,
		Hours:   (150 + 200 + 125) / 8,
		Address: a,
	})
	assert.Equal(t, len(tx.In), 3)
	assert.Equal(t, tx.In, []cipher.SHA256{
		uxs[0].Hash(), uxs[1].Hash(), uxs[2].Hash(),
	})
	assert.Nil(t, tx.Verify())

	// Would be valid, but unconfirmed subtraction causes it to not be
	// First, make a txn to subtract
	uxs = makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	assert.Nil(t, err)
	// Add it to the unconfirmed pool (bypass InjectTxn to avoid blockchain)
	ux := uncf.createUnconfirmedTxn(&unsp, tx)
	uncf.Txns.put(&ux)
	// Make a spend that must not reuse previous addresses
	_, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, a)
	assertError(t, err, "Not enough coins")
}
