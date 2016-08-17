package visor

import (
	"bytes"
	"log"
	"sort"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/stretchr/testify/assert"
)

func tNow() uint64 {
	return uint64(time.Now().UTC().Unix())
}

func assertError(t *testing.T, err error, msg string) {
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), msg)
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeUxBalances(b []wallet.Balance, headTime uint64) coin.UxArray {
	uxs := make(coin.UxArray, len(b))
	for i, _ := range b {
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
	for i, _ := range uxs {
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
	for i, _ := range uxa {
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
	for i, _ := range uxs {
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
	for i, _ := range uxs[:len(uxs)-1] {
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

func TestCalculateBurnAndChange(t *testing.T) {
	// Not enough for fee
	burn, change, err := calculateBurnAndChange(100, 50, 200, 1)
	assertError(t, err, "Insufficient total")

	// Not enough for spend
	burn, change, err = calculateBurnAndChange(100, 200, 10, 1)
	assertError(t, err, "Insufficient total")

	// 0 factor is no burn
	burn, change, err = calculateBurnAndChange(100, 50, 10, 0)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(0))
	assert.Equal(t, change, uint64(40))

	// 1 factor is 100% burn
	// 110 - 10 / 2 = 50
	// 110 - 50 - 10 - 50 = 0
	// 50 + 0 / 1 = 50
	burn, change, err = calculateBurnAndChange(110, 50, 10, 1)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(50))
	assert.Equal(t, change, uint64(0))

	// 50% burn
	// (100 - 10) / 3 = 30
	// 100 - 10 - 50 - 30 = 10
	// 50 + 10 / 2 = 30
	burn, change, err = calculateBurnAndChange(100, 50, 10, 2)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(30))
	assert.Equal(t, change, uint64(10))

	// 33% burn
	// 100 / 4 = 25
	// 100 - 0 - 60 - 25 = 15
	// 60 + 15 / 3 = 25
	burn, change, err = calculateBurnAndChange(100, 60, 0, 3)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(25))
	assert.Equal(t, change, uint64(15))

	// 25% burn
	// 100 / 5 = 20
	// 100 - 0 - 70 - 20 = 10
	// 70 + 10 / 4 = 20
	burn, change, err = calculateBurnAndChange(100, 70, 0, 4)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(20))
	assert.Equal(t, change, uint64(10))

	// Leftover coins from division remainder go to change, not burn
	for i := uint64(1); i < uint64(5); i++ {
		burn, change, err = calculateBurnAndChange(100+i, 70, 0, 4)
		assert.Nil(t, err)
		assert.Equal(t, burn, uint64(20))
		assert.Equal(t, change, uint64(10+i))
	}

	// 25% burn
	// 105 / 5 = 21
	// 105 - 0 - 70 - 21 = 14
	// 70 + 14 / 4 = 24
	burn, change, err = calculateBurnAndChange(105, 70, 0, 4)
	assert.Nil(t, err)
	assert.Equal(t, burn, uint64(21))
	assert.Equal(t, change, uint64(14))
}

func TestCreateSpendsNotEnoughCoins(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{10e6, 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{1e6, 100},
		wallet.Balance{8e6, 0},
	}, now)
	_, err := createSpends(now, uxs, amt, 0, 0)
	assertError(t, err, "Not enough coins")
}

func TestCreateSpendsNotEnoughHours(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{10e6, 110}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{2e6, 100},
		wallet.Balance{8e6, 0},
	}, now)
	_, err := createSpends(now, uxs, amt, 0, 0)
	assertError(t, err, "Not enough hours")
}

func TestIgnoreBadCoins(t *testing.T) {
	// We would satisfy this spend if the bad coins were not skipped
	now := tNow()
	amt := wallet.Balance{10e6, 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{2e6, 50},
		wallet.Balance{8e6, 0},
		wallet.Balance{0, 100},
		wallet.Balance{1e6 + 1, 100},
	}, now)
	_, err := createSpends(now, uxs, amt, 0, 0)
	assertError(t, err, "Not enough hours")
}

func TestBadSpending(t *testing.T) {
	_, err := createSpends(tNow(), coin.UxArray{},
		wallet.Balance{1e6 + 1, 1000}, 0, 1)
	assertError(t, err, "Coins must be multiple of 1e6")
	_, err = createSpends(tNow(), coin.UxArray{},
		wallet.Balance{0, 100}, 0, 1)
	assertError(t, err, "Zero spend amount")
}

func TestCreateSpendsExact(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{10e6, 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{1e6, 50},
		wallet.Balance{8e6, 40},
		wallet.Balance{2e6, 60},
	}, now)
	// Force them to get sorted
	uxs[2].Head.BkSeq = uint64(0)
	uxs[1].Head.BkSeq = uint64(1)
	uxs[0].Head.BkSeq = uint64(2)
	cuxs := append(coin.UxArray{}, uxs...)
	spends, err := createSpends(now, uxs, amt, 0, 0)
	assert.Nil(t, err)
	assert.Equal(t, len(spends), 2)
	assert.Equal(t, spends, coin.UxArray{cuxs[2], cuxs[1]})
}

func TestCreateSpendsWithBurn(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{10e6, 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{1e6, 50},
		wallet.Balance{8e6, 40},
		wallet.Balance{2e6, 60},
	}, now)
	// Force them to get sorted
	uxs[2].Head.BkSeq = uint64(0)
	uxs[1].Head.BkSeq = uint64(1)
	uxs[0].Head.BkSeq = uint64(2)
	cuxs := append(coin.UxArray{}, uxs...)
	// Should spend 8e6,2e6 for the exact amount, but have to add 1e6 to
	// obtain +50 for a 50% fee
	spends, err := createSpends(now, uxs, amt, 0, 2)
	assert.Nil(t, err)
	assert.Equal(t, len(spends), 3)
	assert.Equal(t, spends, coin.UxArray{cuxs[2], cuxs[1], cuxs[0]})

	have := wallet.Balance{0, 0}
	for _, ux := range spends {
		have = have.Add(wallet.NewBalanceFromUxOut(now, &ux))
	}
	burn, change, err := calculateBurnAndChange(have.Hours, amt.Hours, 0, 2)
	assert.Equal(t, burn, uint64(50))
	assert.Equal(t, change, uint64(0))
	assert.Nil(t, err)
}

func TestCreateSpendsNotEnoughBurn(t *testing.T) {
	now := tNow()
	amt := wallet.Balance{10e6, 100}
	uxs := makeUxBalances([]wallet.Balance{
		wallet.Balance{1e6, 40},
		wallet.Balance{8e6, 40},
		wallet.Balance{2e6, 60},
	}, now)
	// Force them to get sorted
	uxs[2].Head.BkSeq = uint64(0)
	uxs[1].Head.BkSeq = uint64(1)
	uxs[0].Head.BkSeq = uint64(2)
	_, err := createSpends(now, uxs, amt, 0, 2)
	assertError(t, err, "Not enough hours to burn")
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
	spends, err := createSpends(now, uxs, amt, 0, 0)
	assert.True(t, sort.IsSorted(OldestUxOut(uxs)))
	assert.Nil(t, err)
	assert.Equal(t, spends, cuxs[:len(spends)])
	assert.Equal(t, len(spends), 4)
	assert.Equal(t, spends, coin.UxArray{ouxs[4], ouxs[2], ouxs[1], ouxs[3]})

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
	assert.Equal(t, len(ouxs), len(spends))
	assert.Equal(t, ouxs, spends)
}

func TestCreateSpendingTransaction(t *testing.T) {
	// Setup
	w := wallet.NewSimpleWallet()
	for i := 0; i < 4; i++ {
		w.CreateEntry()
	}
	uncf := NewUnconfirmedTxnPool()
	now := tNow()
	a := makeAddress()

	// Failing createSpends
	amt := wallet.Balance{0, 0}
	unsp := coin.NewUnspentPool()
	_, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, 0, 0, a)
	assert.NotNil(t, err)

	// Valid txn, fee, no change
	uxs := makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	tx, err := CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
	assert.Nil(t, err)
	assert.Equal(t, len(tx.Out), 1)
	assert.Equal(t, tx.Out[0], coin.TransactionOutput{
		Coins:   25e6,
		Hours:   200,
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
	tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
	assert.Nil(t, err)
	assert.Equal(t, len(tx.Out), 2)
	assert.Equal(t, tx.Out[0], coin.TransactionOutput{
		Coins:   1e6,
		Hours:   (150 + 200 + 125) - (200 + 100),
		Address: w.GetAddresses()[0],
	})
	assert.Equal(t, tx.Out[1], coin.TransactionOutput{
		Coins:   25e6,
		Hours:   200,
		Address: a,
	})
	assert.Equal(t, len(tx.In), 3)
	assert.Equal(t, tx.In, []cipher.SHA256{
		uxs[0].Hash(), uxs[1].Hash(), uxs[2].Hash(),
	})
	assert.Nil(t, tx.Verify())

	// Valid txn, but wastes coin hours
	uxs = makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 200},
	}, now, w.GetAddresses()[:2])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	_, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
	assertError(t, err, "Have enough coins, but not enough to send coin "+
		"hours change back. Would spend 50 more hours than requested.")

	// Would be valid, but unconfirmed subtraction causes it to not be
	// First, make a txn to subtract
	uxs = makeUxBalancesForAddresses([]wallet.Balance{
		wallet.Balance{10e6, 150},
		wallet.Balance{15e6, 150},
	}, now, w.GetAddresses()[:2])
	unsp = coin.NewUnspentPool()
	addUxArrayToUnspentPool(&unsp, uxs)
	amt = wallet.Balance{25e6, 200}
	tx, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
	assert.Nil(t, err)
	// Add it to the unconfirmed pool (bypass InjectTxn to avoid blockchain)
	uncf.Txns[tx.Hash()] = uncf.createUnconfirmedTxn(&unsp, tx,
		w.GetAddressSet())
	// Make a spend that must not reuse previous addresses
	_, err = CreateSpendingTransaction(w, uncf, &unsp, now, amt, 100, 0, a)
	assertError(t, err, "Not enough coins")
}
