package transaction

import (
	"bytes"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
)

func TestSortSpendsCoinsLowToHigh(t *testing.T) {
	// UxBalances are sorted with Coins lowest, then following other order rules
	orderedUxb := []UxBalance{
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 5,
			Coins: 1,
			Hours: 0,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 3,
			Coins: 10,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 1,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 100,
			Hours: 100,
		},
	}

	shuffleWorked := false
	nShuffle := 20
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		// Sanity check that shuffling produces a new result
		if !reflect.DeepEqual(uxb, orderedUxb) {
			shuffleWorked = true
		}

		sortSpendsCoinsLowToHigh(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsLowToHigh(t, uxb)
	}

	require.True(t, shuffleWorked)

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsLowToHigh(uxb)
		verifySortedCoinsLowToHigh(t, uxb)
	}
}

func TestSortSpendsCoinsHighToLow(t *testing.T) {
	// UxBalances are sorted with Coins highest, then following other order rules
	orderedUxb := []UxBalance{
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 10000,
			Hours: 0,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 10,
			Coins: 1000,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 100,
			Hours: 100,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 3,
			Coins: 10,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 1,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
	}

	shuffleWorked := false
	nShuffle := 20
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		rand.Shuffle(len(uxb), func(i, j int) {
			uxb[i], uxb[j] = uxb[j], uxb[i]
		})

		if !reflect.DeepEqual(uxb, orderedUxb) {
			shuffleWorked = true
		}

		sortSpendsCoinsHighToLow(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsHighToLow(t, uxb)
	}

	require.True(t, shuffleWorked)

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsHighToLow(uxb)
		verifySortedCoinsHighToLow(t, uxb)
	}
}

func TestSortSpendsHoursLowToHigh(t *testing.T) {
	// UxBalances are sorted with Coins lowest, then following other order rules
	orderedUxb := []UxBalance{
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 5,
			Coins: 1,
			Hours: 0,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 3,
			Coins: 10,
			Hours: 1,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 1,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("bddf0aaf80f96c144f33ac8a27764a868d37e1c11e568063ebeb1367de859566"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  cipher.MustSHA256FromHex("f569461182b0efe9a5c666e9a35c6602b351021c1803cc740aca548cf6db4cb2"),
			BkSeq: 2,
			Coins: 10,
			Hours: 10,
		},
		{
			Hash:  testutil.RandSHA256(t),
			BkSeq: 4,
			Coins: 100,
			Hours: 100,
		},
	}

	shuffleWorked := false
	nShuffle := 20
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		// Sanity check that shuffling produces a new result
		if !reflect.DeepEqual(uxb, orderedUxb) {
			shuffleWorked = true
		}

		sortSpendsHoursLowToHigh(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedHoursLowToHigh(t, uxb)
	}

	require.True(t, shuffleWorked)

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsHoursLowToHigh(uxb)
		verifySortedHoursLowToHigh(t, uxb)
	}
}

func TestChooseSpendsMaximizeUxOuts(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins <= b.Coins
		})
	}

	// 0 coins (error)
	uxb := makeRandomUxBalances(t)
	verifyChosenCoins(t, uxb, 0, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
		return a.Coins <= b.Coins
	})

	// 0 coins in a UxBalance (panic)
	uxb[1].Coins = 0
	require.Panics(t, func() {
		verifyChosenCoins(t, uxb, 10, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins <= b.Coins
		})
	})

	// MaxUint64 coins (error)
	uxb = makeRandomUxBalances(t)
	verifyChosenCoins(t, uxb, math.MaxUint64, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
		return a.Coins <= b.Coins
	})
}

func TestChooseSpendsMinimizeUxOutsRandom(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins >= b.Coins
		})
	}

	// 0 coins (error)
	uxb := makeRandomUxBalances(t)
	verifyChosenCoins(t, uxb, 0, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
		return a.Coins >= b.Coins
	})

	// 0 coins in a UxBalance (panic)
	uxb[1].Coins = 0
	require.Panics(t, func() {
		verifyChosenCoins(t, uxb, 10, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins >= b.Coins
		})
	})

	// MaxUint64 coins (error)
	uxb = makeRandomUxBalances(t)
	verifyChosenCoins(t, uxb, math.MaxUint64, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
		return a.Coins >= b.Coins
	})
}

func makeRandomUxBalances(t *testing.T) []UxBalance {
	// Generate random 0-100 UxBalances
	// Coins 1-10 (must be >0)
	// Hours 0-10
	// BkSeq 0-10
	// Hash random
	// Small ranges are used for Coins, Hours, BkSeq to increase likelihood
	// that they collide and test deeper sorting comparisons

	n := rand.Intn(101)
	uxb := make([]UxBalance, n)

	// Use a random max range for the hours' rand range to ensure enough
	// balances have zero hours
	hasZeroHoursRange := rand.Intn(3) + 1

	for i := 0; i < n; i++ {
		ux := UxBalance{
			Coins: uint64(rand.Intn(10) + 1), // 1-10
			Hours: uint64(rand.Intn(hasZeroHoursRange)),
			BkSeq: uint64(rand.Intn(11)), // 0-10
			Hash:  testutil.RandSHA256(t),
		}

		uxb[i] = ux
	}

	return uxb
}

func verifyChosenCoins(t *testing.T, uxb []UxBalance, coins uint64, chooseSpends func([]UxBalance, uint64, uint64) ([]UxBalance, error), cmpCoins func(i, j UxBalance) bool) {
	var haveZero, haveNonzero int
	for _, ux := range uxb {
		if ux.Hours == 0 {
			haveZero++
		} else {
			haveNonzero++
		}
	}

	var totalCoins, totalHours uint64
	for _, ux := range uxb {
		totalCoins += ux.Coins
		totalHours += ux.Hours
	}

	chosen, err := chooseSpends(uxb, coins, 0)

	if coins == 0 {
		testutil.RequireError(t, err, ErrZeroSpend.Error())
		return
	}

	if len(uxb) == 0 {
		testutil.RequireError(t, err, ErrNoUnspents.Error())
		return
	}

	if totalHours == 0 {
		testutil.RequireError(t, err, fee.ErrTxnNoFee.Error())
		return
	}

	if coins > totalCoins {
		testutil.RequireError(t, err, ErrInsufficientBalance.Error())
		return
	}

	require.NoError(t, err)
	require.NotEqual(t, 0, len(chosen))

	// Check that there are no duplicated spends chosen
	uxMap := make(map[UxBalance]struct{}, len(chosen))
	for _, ux := range chosen {
		_, ok := uxMap[ux]
		require.False(t, ok)
		uxMap[ux] = struct{}{}
	}

	// The first chosen spend should have non-zero coin hours
	require.NotEqual(t, uint64(0), chosen[0].Hours)

	// Outputs with zero hours should come before any outputs with non-zero hours,
	// except for the first output
	for i := range chosen {
		if i <= 1 {
			continue
		}

		a := chosen[i-1]
		b := chosen[i]

		if b.Hours == 0 {
			require.Equal(t, uint64(0), a.Hours)
		}
	}

	// The initial UxBalance with hours should have more or equal coins than any other UxBalance with hours
	// If it has equal coins, it should have less hours
	for _, ux := range chosen[1:] {
		if ux.Hours != 0 {
			require.True(t, chosen[0].Coins >= ux.Coins)

			if chosen[0].Coins == ux.Coins {
				require.True(t, chosen[0].Hours <= ux.Hours)
			}
		}
	}

	var zeroBalances, nonzeroBalances []UxBalance
	for _, ux := range chosen[1:] {
		if ux.Hours == 0 {
			zeroBalances = append(zeroBalances, ux)
		} else {
			nonzeroBalances = append(nonzeroBalances, ux)
		}
	}

	// Amongst the UxBalances with zero hours, they should be sorted as specified
	verifySortedCoins(t, zeroBalances, cmpCoins)

	// Amongst the UxBalances with non-zero hours, they should be sorted as specified
	verifySortedCoins(t, nonzeroBalances, cmpCoins)

	// If there are any extra UxBalances with non-zero hours, all of the zeros should have been chosen
	if len(nonzeroBalances) > 0 {
		require.Equal(t, haveZero, len(zeroBalances))
	}

	// Excessive UxBalances to satisfy the amount requested should not be included
	var haveCoins uint64
	for i, ux := range chosen {
		haveCoins += ux.Coins
		if haveCoins >= coins {
			require.Equal(t, len(chosen)-1, i)
		}
	}
}

func verifySortedCoins(t *testing.T, uxb []UxBalance, cmpCoins func(a, b UxBalance) bool) {
	if len(uxb) <= 1 {
		return
	}

	for i := range uxb {
		if i == 0 {
			continue
		}

		a := uxb[i-1]
		b := uxb[i]

		require.True(t, cmpCoins(a, b))

		if a.Coins == b.Coins {
			require.True(t, a.Hours <= b.Hours)

			if a.Hours == b.Hours {
				require.True(t, a.BkSeq <= b.BkSeq)

				if a.BkSeq == b.BkSeq {
					cmp := bytes.Compare(a.Hash[:], b.Hash[:])
					require.True(t, cmp < 0)
				}
			}
		}
	}
}

func verifySortedCoinsLowToHigh(t *testing.T, uxb []UxBalance) {
	verifySortedCoins(t, uxb, func(a, b UxBalance) bool {
		return a.Coins <= b.Coins
	})
}

func verifySortedCoinsHighToLow(t *testing.T, uxb []UxBalance) {
	verifySortedCoins(t, uxb, func(a, b UxBalance) bool {
		return a.Coins >= b.Coins
	})
}

func verifySortedHours(t *testing.T, uxb []UxBalance, cmpHours func(a, b UxBalance) bool) {
	if len(uxb) <= 1 {
		return
	}

	for i := range uxb {
		if i == 0 {
			continue
		}

		a := uxb[i-1]
		b := uxb[i]

		require.True(t, cmpHours(a, b))

		if a.Hours == b.Hours {
			require.True(t, a.Coins <= b.Coins)

			if a.Coins == b.Coins {
				require.True(t, a.BkSeq <= b.BkSeq)

				if a.BkSeq == b.BkSeq {
					cmp := bytes.Compare(a.Hash[:], b.Hash[:])
					require.True(t, cmp < 0)
				}
			}
		}
	}
}

func verifySortedHoursLowToHigh(t *testing.T, uxb []UxBalance) {
	verifySortedHours(t, uxb, func(a, b UxBalance) bool {
		return a.Hours <= b.Hours
	})
}
