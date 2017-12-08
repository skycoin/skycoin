package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func TestNewWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name    string
		wltName string
		ops     Options
		expect  expect
	}{
		{
			"ok with seed set",
			"test.wlt",
			Options{
				Seed: "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "",
					"filename": "test.wlt",
					"coin":     "skycoin",
					"type":     "deterministic",
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
		{
			"ok with label and seed set",
			"test.wlt",
			Options{
				Label: "wallet1",
				Seed:  "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     "skycoin",
					"type":     "deterministic",
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
		{
			"ok with label, seed and coin set",
			"test.wlt",
			Options{
				Label: "wallet1",
				Coin:  CoinTypeBitcoin,
				Seed:  "testseed123",
			},
			expect{
				meta: map[string]string{
					"label":    "wallet1",
					"filename": "test.wlt",
					"coin":     string(CoinTypeBitcoin),
					"type":     "deterministic",
					"seed":     "testseed123",
				},
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w, err := NewWallet(tc.wltName, tc.ops)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}
			require.NoError(t, w.Validate())
			for k, v := range tc.expect.meta {
				vv, ok := w.Meta[k]
				require.True(t, ok)
				require.Equal(t, v, vv)
			}
		})
	}
}

func TestLoadWallet(t *testing.T) {
	type expect struct {
		meta map[string]string
		err  error
	}

	tt := []struct {
		name   string
		file   string
		expect expect
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			expect{
				meta: map[string]string{
					"coin":     string(CoinTypeSkycoin),
					"filename": "test1.wlt",
					"label":    "test3",
					"lastSeed": "9182b02c0004217ba9a55593f8cf0abecc30d041e094b266dbb5103e1919adaf",
					"seed":     "buddy fossil side modify turtle door label grunt baby worth brush master",
					"tm":       "1503458909",
					"type":     "deterministic",
					"version":  "0.1",
				},
				err: nil,
			},
		},
		{
			"wallet file doesn't exist",
			"not_exist_file.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("load wallet file failed, wallet not_exist_file.wlt doesn't exist"),
			},
		},
		{
			"invalid wallet: no type",
			"./testdata/invalid_wallets/no_type.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_type.wlt: type field not set"),
			},
		},
		{
			"invalid wallet: invalid type",
			"./testdata/invalid_wallets/err_type.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet err_type.wlt: wallet type invalid"),
			},
		},
		{
			"invalid wallet: no coin",
			"./testdata/invalid_wallets/no_coin.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_coin.wlt: coin field not set"),
			},
		},
		{
			"invalid wallet: no seed",
			"./testdata/invalid_wallets/no_seed.wlt",
			expect{
				meta: map[string]string{},
				err:  fmt.Errorf("invalid wallet no_seed.wlt: seed field not set"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := Wallet{}
			err := w.Load(tc.file)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			for k, v := range tc.expect.meta {
				vv := w.Meta[k]
				require.Equal(t, v, vv)
			}
		})
	}
}

func TestWalletGetEntry(t *testing.T) {
	tt := []struct {
		name    string
		wltFile string
		address string
		find    bool
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			"JUdRuTiqD1mGcw358twMg3VPpXpzbkdRvJ",
			true,
		},
		{
			"entry not exist",
			"./testdata/test1.wlt",
			"2ULfxDUuenUY5V4Pr8whmoAwFdUseXNyjXC",
			false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := Wallet{}
			require.NoError(t, w.Load(tc.wltFile))
			a, err := cipher.DecodeBase58Address(tc.address)
			require.NoError(t, err)
			e, ok := w.GetEntry(a)
			require.Equal(t, tc.find, ok)
			if ok {
				require.Equal(t, tc.address, e.Address.String())
			}
		})
	}
}

func TestWalletAddEntry(t *testing.T) {
	test1SecKey, err := cipher.SecKeyFromHex("1fc5396e91e60b9fc613d004ea5bd2ccea17053a12127301b3857ead76fdb93e")
	require.NoError(t, err)

	_, s := cipher.GenerateKeyPair()
	seckeys := []cipher.SecKey{
		test1SecKey,
		s,
	}

	tt := []struct {
		name    string
		wltFile string
		secKey  cipher.SecKey
		err     error
	}{
		{
			"ok",
			"./testdata/test1.wlt",
			seckeys[1],
			nil,
		},
		{
			"dup entry",
			"./testdata/test1.wlt",
			seckeys[0],
			errors.New("duplicate address entry"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := Wallet{}
			require.NoError(t, w.Load(tc.wltFile))
			a := cipher.AddressFromSecKey(tc.secKey)
			p := cipher.PubKeyFromSecKey(tc.secKey)
			require.Equal(t, tc.err, w.AddEntry(Entry{
				Address: a,
				Public:  p,
				Secret:  s,
			}))
		})
	}
}

type distributeSpendHoursTestCase struct {
	name              string
	inputHours        uint64
	nAddrs            uint64
	haveChange        bool
	expectChangeHours uint64
	expectAddrHours   []uint64
}

var burnFactor2TestCases = []distributeSpendHoursTestCase{
	{
		name:            "no input hours, one addr, no change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "no input hours, two addrs, no change",
		inputHours:      0,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{0, 0},
	},
	{
		name:            "no input hours, one addr, change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      true,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "one input hour, one addr, no change",
		inputHours:      1,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "two input hours, one addr, no change",
		inputHours:      2,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:              "two input hours, one addr, change",
		inputHours:        2,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:              "three input hours, one addr, change",
		inputHours:        3,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:            "three input hours, one addr, no change",
		inputHours:      3,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:            "three input hours, two addrs, no change",
		inputHours:      3,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{1, 0},
	},
	{
		name:            "four input hours, one addr, no change",
		inputHours:      4,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:              "four input hours, one addr, change",
		inputHours:        4,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "four input hours, two addr, change",
		inputHours:        4,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1, 0},
	},
	{
		name:              "30 (divided by 2, odd number) input hours, two addr, change",
		inputHours:        30,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{4, 3},
	},
	{
		name:              "33 (odd number) input hours, two addr, change",
		inputHours:        33,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{4, 4},
	},
	{
		name:              "33 (odd number) input hours, three addr, change",
		inputHours:        33,
		nAddrs:            3,
		haveChange:        true,
		expectChangeHours: 8,
		expectAddrHours:   []uint64{3, 3, 2},
	},
}

var burnFactor3TestCases = []distributeSpendHoursTestCase{
	{
		name:            "no input hours, one addr, no change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "no input hours, two addrs, no change",
		inputHours:      0,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{0, 0},
	},
	{
		name:            "no input hours, one addr, change",
		inputHours:      0,
		nAddrs:          1,
		haveChange:      true,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "one input hour, one addr, no change",
		inputHours:      1,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{0},
	},
	{
		name:            "two input hours, one addr, no change",
		inputHours:      2,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{1},
	},
	{
		name:            "three input hours, one addr, no change",
		inputHours:      3,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:              "two input hours, one addr, change",
		inputHours:        2,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{0},
	},
	{
		name:              "three input hours, one addr, change",
		inputHours:        3,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "four input hours, one addr, change",
		inputHours:        4,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 1,
		expectAddrHours:   []uint64{1},
	},
	{
		name:            "four input hours, one addr, no change",
		inputHours:      4,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{2},
	},
	{
		name:            "four input hours, two addrs, no change",
		inputHours:      4,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{1, 1},
	},
	{
		name:            "five input hours, one addr, no change",
		inputHours:      5,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{3},
	},
	{
		name:              "five input hours, one addr, change",
		inputHours:        5,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1},
	},
	{
		name:              "five input hours, two addr, change",
		inputHours:        5,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1, 0},
	},
	{
		name:              "32 input hours, two addr, change",
		inputHours:        32,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 11,
		expectAddrHours:   []uint64{5, 5},
	},
	{
		name:              "35 input hours, two addr, change",
		inputHours:        35,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 12,
		expectAddrHours:   []uint64{6, 5},
	},
	{
		name:              "32 input hours, three addr, change",
		inputHours:        32,
		nAddrs:            3,
		haveChange:        true,
		expectChangeHours: 11,
		expectAddrHours:   []uint64{4, 3, 3},
	},
}

func TestWalletDistributeSpendHours(t *testing.T) {
	var cases []distributeSpendHoursTestCase
	switch fee.BurnFactor {
	case 2:
		cases = burnFactor2TestCases
	case 3:
		cases = burnFactor3TestCases
	default:
		t.Fatalf("No test cases defined for fee.BurnFactor=%d", fee.BurnFactor)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			changeHours, addrHours, totalHours := DistributeSpendHours(tc.inputHours, tc.nAddrs, tc.haveChange)
			require.Equal(t, tc.expectChangeHours, changeHours)
			require.Equal(t, tc.expectAddrHours, addrHours)
			require.Equal(t, tc.nAddrs, uint64(len(addrHours)))

			outputHours := changeHours
			for _, h := range addrHours {
				outputHours += h
			}
			require.True(t, tc.inputHours >= outputHours)
			require.Equal(t, outputHours, totalHours)

			if tc.inputHours != 0 {
				err := fee.VerifyTransactionFeeForHours(outputHours, tc.inputHours-outputHours)
				require.NoError(t, err)
			}
		})
	}

	// Tests over range of values
	for inputHours := uint64(0); inputHours <= 1e3; inputHours++ {
		for nAddrs := uint64(1); nAddrs < 16; nAddrs++ {
			for _, haveChange := range []bool{true, false} {
				name := fmt.Sprintf("inputHours=%d nAddrs=%d haveChange=%v", inputHours, nAddrs, haveChange)
				t.Run(name, func(t *testing.T) {
					changeHours, addrHours, totalHours := DistributeSpendHours(inputHours, nAddrs, haveChange)
					require.Equal(t, nAddrs, uint64(len(addrHours)))

					var sumAddrHours uint64
					for _, h := range addrHours {
						sumAddrHours += h
					}

					if haveChange {
						remainingHours := (inputHours - fee.RequiredFee(inputHours))
						splitRemainingHours := remainingHours / 2
						require.True(t, changeHours == splitRemainingHours || changeHours == splitRemainingHours+1)
						require.Equal(t, splitRemainingHours, sumAddrHours)
					} else {
						require.Equal(t, uint64(0), changeHours)
						require.Equal(t, inputHours-fee.RequiredFee(inputHours), sumAddrHours)
					}

					outputHours := sumAddrHours + changeHours
					require.True(t, inputHours >= outputHours)
					require.Equal(t, outputHours, totalHours)

					if inputHours != 0 {
						err := fee.VerifyTransactionFeeForHours(outputHours, inputHours-outputHours)
						require.NoError(t, err)
					}

					// addrHours at the beginning and end of the array should not differ by more than one
					max := addrHours[0]
					min := addrHours[len(addrHours)-1]
					require.True(t, max-min <= 1)
				})
			}
		}
	}
}

func TestWalletSortSpendsLowToHigh(t *testing.T) {
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

	for i := 0; i < 10; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		require.NotEqual(t, uxb, orderedUxb)

		sortSpendsCoinsLowToHigh(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsLowToHigh(t, uxb)
	}

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsHighToLow(uxb)
		verifySortedCoinsHighToLow(t, uxb)
	}
}

func TestWalletSortSpendsHighToLow(t *testing.T) {
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

	nShuffle := 10
	for i := 0; i < nShuffle; i++ {
		// Shuffle the list
		uxb := make([]UxBalance, len(orderedUxb))
		copy(uxb, orderedUxb)

		for i := range uxb {
			j := rand.Intn(i + 1)
			uxb[i], uxb[j] = uxb[j], uxb[i]
		}

		require.NotEqual(t, uxb, orderedUxb)

		sortSpendsCoinsHighToLow(uxb)

		for i, ux := range uxb {
			require.Equal(t, orderedUxb[i], ux, "index %d", i)
		}

		verifySortedCoinsHighToLow(t, uxb)
	}

	nRand := 1000
	for i := 0; i < nRand; i++ {
		uxb := makeRandomUxBalances(t)

		sortSpendsCoinsHighToLow(uxb)
		verifySortedCoinsHighToLow(t, uxb)
	}
}

func TestWalletChooseSpendsMaximizeUxOuts(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMaximizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins <= b.Coins
		})
	}
}

func TestWalletChooseSpendsMinimizeUxOuts(t *testing.T) {
	nRand := 10000
	for i := 0; i < nRand; i++ {
		coins := uint64((rand.Intn(3)+1)*10 + rand.Intn(3)) // 10,20,30 + 0,1,2
		uxb := makeRandomUxBalances(t)

		verifyChosenCoins(t, uxb, coins, ChooseSpendsMinimizeUxOuts, func(a, b UxBalance) bool {
			return a.Coins >= b.Coins
		})
	}
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

func verifyChosenCoins(t *testing.T, uxb []UxBalance, coins uint64, chooseSpends func([]UxBalance, uint64) ([]UxBalance, error), cmpCoins func(i, j UxBalance) bool) {
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

	chosen, err := chooseSpends(uxb, coins)

	if coins == 0 {
		testutil.RequireError(t, err, "zero spend amount")
		return
	}

	if len(uxb) == 0 {
		testutil.RequireError(t, err, "no unspents to spend")
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
