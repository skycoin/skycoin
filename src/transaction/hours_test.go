package transaction

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/mathutil"
)

func TestDistributeCoinHoursProportional(t *testing.T) {
	cases := []struct {
		name   string
		coins  []uint64
		hours  uint64
		output []uint64
		err    error
	}{
		{
			name:  "no coins",
			hours: 1,
			err:   errors.New("DistributeCoinHoursProportional coins array must not be empty"),
		},
		{
			name:  "coins have 0 in them",
			coins: []uint64{1, 2, 0, 3},
			hours: 1,
			err:   errors.New("DistributeCoinHoursProportional coins array has a zero value"),
		},
		{
			name:  "total coins too large while adding",
			coins: []uint64{10, math.MaxUint64 - 9},
			hours: 1,
			err:   mathutil.ErrUint64AddOverflow,
		},
		{
			name:  "total coins too large after adding",
			coins: []uint64{10, math.MaxInt64},
			hours: 1,
			err:   mathutil.ErrUint64OverflowsInt64,
		},
		{
			name:  "single coin too large",
			coins: []uint64{10, math.MaxInt64 + 1},
			hours: 1,
			err:   mathutil.ErrUint64OverflowsInt64,
		},
		{
			name:  "hours too large",
			coins: []uint64{10},
			hours: math.MaxInt64 + 1,
			err:   mathutil.ErrUint64OverflowsInt64,
		},

		{
			name:   "valid, one input",
			coins:  []uint64{1},
			hours:  1,
			output: []uint64{1},
		},

		{
			name:   "zero hours",
			coins:  []uint64{1},
			hours:  0,
			output: []uint64{0},
		},

		{
			name:   "valid, multiple inputs, all equal",
			coins:  []uint64{2, 4, 8, 16},
			hours:  30,
			output: []uint64{2, 4, 8, 16},
		},

		{
			name:   "valid, multiple inputs, rational division in coins and hours",
			coins:  []uint64{2, 4, 8, 16},
			hours:  30,
			output: []uint64{2, 4, 8, 16},
		},

		{
			name:   "valid, multiple inputs, rational division in coins, irrational in hours",
			coins:  []uint64{2, 4, 8, 16},
			hours:  31,
			output: []uint64{3, 4, 8, 16},
		},

		{
			name:   "valid, multiple inputs, irrational division in coins, rational in hours",
			coins:  []uint64{2, 3, 5, 7, 11, 13},
			hours:  41,
			output: []uint64{2, 3, 5, 7, 11, 13},
		},

		{
			name:   "valid, multiple inputs, irrational division in coins and hours",
			coins:  []uint64{2, 3, 5, 7, 11, 13},
			hours:  50,
			output: []uint64{3, 4, 7, 8, 13, 15},
		},

		{
			name:   "valid, multiple inputs that would receive 0 hours but get compensated from remainder as priority",
			coins:  []uint64{16, 8, 4, 2, 1, 1},
			hours:  14,
			output: []uint64{7, 3, 1, 1, 1, 1},
		},

		{
			name:   "not enough hours for everyone",
			coins:  []uint64{1, 1, 1, 1, 1},
			hours:  1,
			output: []uint64{1, 0, 0, 0, 0},
		},

		{
			name:   "not enough hours for everyone 2",
			coins:  []uint64{1, 1, 1, 1, 1},
			hours:  3,
			output: []uint64{1, 1, 1, 0, 0},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hours, err := DistributeCoinHoursProportional(tc.coins, tc.hours)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.output, hours)
			}
		})
	}

	// Randomized tests
	iterations := 10000
	maxCoinsLen := 300
	maxMaxCoins := 100000
	maxHours := 15000000
	coins := make([]uint64, maxCoinsLen)
	for i := 0; i < iterations; i++ {
		coinsLen := rand.Intn(maxCoinsLen) + 1

		maxCoins := rand.Intn(maxMaxCoins) + 1

		var totalCoins uint64
		for i := 0; i < coinsLen; i++ {
			coins[i] = uint64(rand.Intn(maxCoins) + 1)

			var err error
			totalCoins, err = mathutil.AddUint64(totalCoins, coins[i])
			require.NoError(t, err)
		}

		hours := uint64(rand.Intn(maxHours))

		output, err := DistributeCoinHoursProportional(coins[:coinsLen], hours)
		require.NoError(t, err)

		require.Equal(t, coinsLen, len(output))

		var totalHours uint64
		for _, h := range output {
			if hours >= totalCoins {
				require.NotEqual(t, uint64(0), h)
			}

			var err error
			totalHours, err = mathutil.AddUint64(totalHours, h)
			require.NoError(t, err)
		}

		require.Equal(t, hours, totalHours)
	}
}

func TestDistributeSpendHours(t *testing.T) {
	originalBurnFactor := params.UserVerifyTxn.BurnFactor

	cases := []struct {
		burnFactor uint32
		cases      []distributeSpendHoursTestCase
	}{
		{2, burnFactor2TestCases},
		{3, burnFactor3TestCases},
		{10, burnFactor10TestCases},
	}

	tested := false
	for _, tcc := range cases {
		if tcc.burnFactor == params.UserVerifyTxn.BurnFactor {
			tested = true
		}

		for _, tc := range tcc.cases {
			t.Run(tc.name, func(t *testing.T) {
				params.UserVerifyTxn.BurnFactor = tcc.burnFactor
				defer func() {
					params.UserVerifyTxn.BurnFactor = originalBurnFactor
				}()

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
					err := fee.VerifyTransactionFeeForHours(outputHours, tc.inputHours-outputHours, params.UserVerifyTxn.BurnFactor)
					require.NoError(t, err)
				}
			})
		}

		t.Run(fmt.Sprintf("burn-factor-%d-range", tcc.burnFactor), func(t *testing.T) {
			params.UserVerifyTxn.BurnFactor = tcc.burnFactor
			defer func() {
				params.UserVerifyTxn.BurnFactor = originalBurnFactor
			}()

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
								remainingHours := (inputHours - fee.RequiredFee(inputHours, params.UserVerifyTxn.BurnFactor))
								splitRemainingHours := remainingHours / 2
								require.True(t, changeHours == splitRemainingHours || changeHours == splitRemainingHours+1)
								require.Equal(t, splitRemainingHours, sumAddrHours)
							} else {
								require.Equal(t, uint64(0), changeHours)
								require.Equal(t, inputHours-fee.RequiredFee(inputHours, params.UserVerifyTxn.BurnFactor), sumAddrHours)
							}

							outputHours := sumAddrHours + changeHours
							require.True(t, inputHours >= outputHours)
							require.Equal(t, outputHours, totalHours)

							if inputHours != 0 {
								err := fee.VerifyTransactionFeeForHours(outputHours, inputHours-outputHours, params.UserVerifyTxn.BurnFactor)
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
		})
	}

	require.True(t, tested, "configured BurnFactor=%d has not been tested", params.UserVerifyTxn.BurnFactor)
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

var burnFactor10TestCases = []distributeSpendHoursTestCase{
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
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1},
	},
	{
		name:            "four input hours, one addr, no change",
		inputHours:      4,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{3},
	},
	{
		name:            "four input hours, two addrs, no change",
		inputHours:      4,
		nAddrs:          2,
		haveChange:      false,
		expectAddrHours: []uint64{2, 1},
	},
	{
		name:            "five input hours, one addr, no change",
		inputHours:      5,
		nAddrs:          1,
		haveChange:      false,
		expectAddrHours: []uint64{4},
	},
	{
		name:              "five input hours, one addr, change",
		inputHours:        5,
		nAddrs:            1,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{2},
	},
	{
		name:              "five input hours, two addr, change",
		inputHours:        5,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 2,
		expectAddrHours:   []uint64{1, 1},
	},
	{
		name:              "32 input hours, two addr, change",
		inputHours:        32,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 14,
		expectAddrHours:   []uint64{7, 7},
	},
	{
		name:              "35 input hours, two addr, change",
		inputHours:        35,
		nAddrs:            2,
		haveChange:        true,
		expectChangeHours: 16,
		expectAddrHours:   []uint64{8, 7},
	},
	{
		name:              "32 input hours, three addr, change",
		inputHours:        32,
		nAddrs:            3,
		haveChange:        true,
		expectChangeHours: 14,
		expectAddrHours:   []uint64{5, 5, 4},
	},
}
