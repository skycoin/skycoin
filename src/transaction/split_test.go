package transaction

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
)

func makeSplitUxOuts(t *testing.T, addr cipher.Address, headTime uint64, coins func(i int) uint64, n int) coin.AddressUxOuts {
	_, secKeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte("seed"), 1)
	uxouts := make([]coin.UxOut, n)
	for i := range uxouts {
		uxouts[i] = makeUxOut(t, secKeys[0], coins(i), uint64(10+i%7)) //nolint:gosec
		uxouts[i].Body.Address = addr
		uxouts[i].Body.SrcTransaction = cipher.SumSHA256([]byte(strconv.Itoa(i)))
		uxouts[i].Head.BkSeq = uint64(i + 1) //nolint:gosec
		uxouts[i].Head.Time = headTime
	}
	return coin.AddressUxOuts{addr: uxouts}
}

func oneCoin(int) uint64 { return 1e6 }

// verifySplit checks the parts of a split send against p and the outputs it had.
func verifySplit(t *testing.T, p Params, auxs coin.AddressUxOuts, headTime uint64, txns []*coin.Transaction, inputs [][]UxBalance) {
	t.Helper()
	require.Len(t, inputs, len(txns))

	uxb, err := NewUxBalances(auxs.Flatten(), headTime)
	require.NoError(t, err)
	known := make(map[cipher.SHA256]UxBalance, len(uxb))
	for _, u := range uxb {
		known[u.Hash] = u
	}

	type paid struct{ coins, hours uint64 }
	got := make(map[cipher.Address]paid)
	seen := make(map[cipher.SHA256]bool)
	for i, txn := range txns {
		require.NoError(t, CheckSize(txn, inputs[i]))
		size, err := txn.Size()
		require.NoError(t, err)
		require.LessOrEqual(t, size, params.UserVerifyTxn.MaxTransactionSize)
		require.NoError(t, VerifySingleTxnUserConstraints(*txn))

		require.Len(t, inputs[i], len(txn.In))
		var inCoins, inHours uint64
		for j, h := range txn.In {
			require.Equal(t, h, inputs[i][j].Hash)
			u, ok := known[h]
			require.True(t, ok, "input is not one of the outputs given")
			require.False(t, seen[h], "input is spent by two parts")
			seen[h] = true
			inCoins += u.Coins
			inHours += u.Hours
		}

		var outCoins, outHours uint64
		var nChange int
		for _, o := range txn.Out {
			outCoins += o.Coins
			outHours += o.Hours
			if o.Address == *p.ChangeAddress {
				nChange++
				continue
			}
			g := got[o.Address]
			g.coins += o.Coins
			g.hours += o.Hours
			got[o.Address] = g
		}
		require.LessOrEqual(t, nChange, 1)
		require.Equal(t, inCoins, outCoins, "part %d does not balance", i)
		require.NoError(t, fee.VerifyTransactionFeeForHours(outHours, inHours-outHours, params.UserVerifyTxn.BurnFactor), "part %d", i)
	}

	want := make(map[cipher.Address]paid)
	for _, o := range p.To {
		w := want[o.Address]
		w.coins += o.Coins
		w.hours += o.Hours
		want[o.Address] = w
	}
	for a, w := range want {
		require.Equal(t, w.coins, got[a].coins, "coins sent to %s", a)
		if p.HoursSelection.Type == HoursSelectionTypeManual {
			require.Equal(t, w.hours, got[a].hours, "hours sent to %s", a)
		}
	}
	require.Len(t, got, len(want))
}

func TestCreateSplit(t *testing.T) {
	headTime := uint64(time.Now().UTC().Unix()) //nolint:gosec
	from := testutil.MakeAddress()
	change := testutil.MakeAddress()
	share := decimal.NewFromFloat(0.5)
	auto := HoursSelection{Type: HoursSelectionTypeAuto, Mode: HoursSelectionModeShare, ShareFactor: &share}
	manual := HoursSelection{Type: HoursSelectionTypeManual}
	to := []cipher.Address{testutil.MakeAddress(), testutil.MakeAddress(), testutil.MakeAddress()}
	rnd := rand.New(rand.NewSource(1)) //nolint:gosec
	randomCoins := make([]uint64, 1500)
	var randomTotal uint64
	for i := range randomCoins {
		randomCoins[i] = uint64(1+rnd.Intn(50)) * 1e6 //nolint:gosec
		randomTotal += randomCoins[i]
	}

	cases := []struct {
		name   string
		nOuts  int
		coins  func(i int) uint64
		hs     HoursSelection
		to     []coin.TransactionOutput
		nParts int
	}{
		{
			name:   "fits in one",
			nOuts:  300,
			coins:  oneCoin,
			hs:     auto,
			to:     []coin.TransactionOutput{{Address: to[0], Coins: 300e6}},
			nParts: 1,
		},
		{
			name:   "one recipient",
			nOuts:  1000,
			coins:  oneCoin,
			hs:     auto,
			to:     []coin.TransactionOutput{{Address: to[0], Coins: 900e6}},
			nParts: 3,
		},
		{
			name:  "several recipients cut across parts",
			nOuts: 1000,
			coins: oneCoin,
			hs:    auto,
			to: []coin.TransactionOutput{
				{Address: to[0], Coins: 250e6},
				{Address: to[1], Coins: 333e6},
				{Address: to[2], Coins: 317e6},
			},
			nParts: 3,
		},
		{
			name:   "manual hours",
			nOuts:  1000,
			coins:  oneCoin,
			hs:     manual,
			to:     []coin.TransactionOutput{{Address: to[0], Coins: 700e6, Hours: 2001}},
			nParts: 3,
		},
		{
			name:   "outputs of different sizes",
			nOuts:  len(randomCoins),
			coins:  func(i int) uint64 { return randomCoins[i] },
			hs:     auto,
			to:     []coin.TransactionOutput{{Address: to[0], Coins: randomTotal - 100e6}},
			nParts: 5,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			auxs := makeSplitUxOuts(t, from, headTime, tc.coins, tc.nOuts)
			p := Params{HoursSelection: tc.hs, ChangeAddress: &change, To: tc.to}

			txns, inputs, err := CreateSplit(p, auxs, headTime)
			require.NoError(t, err)
			require.Len(t, txns, tc.nParts)
			verifySplit(t, p, auxs, headTime, txns, inputs)

			if tc.nParts == 1 {
				txn, in, err := Create(p, auxs, headTime)
				require.NoError(t, err)
				require.Equal(t, txn, txns[0])
				require.Equal(t, in, inputs[0])
				return
			}

			// Every part but the last pays from a full chunk and keeps change,
			// so the unsent hours are not all given to the recipients.
			for _, txn := range txns[:len(txns)-1] {
				require.Equal(t, change, txn.Out[len(txn.Out)-1].Address)
			}
		})
	}
}

func TestCreateSplitInsufficientBalance(t *testing.T) {
	headTime := uint64(time.Now().UTC().Unix()) //nolint:gosec
	auxs := makeSplitUxOuts(t, testutil.MakeAddress(), headTime, oneCoin, 500)
	share := decimal.NewFromFloat(0.5)
	p := Params{
		HoursSelection: HoursSelection{Type: HoursSelectionTypeAuto, Mode: HoursSelectionModeShare, ShareFactor: &share},
		To:             []coin.TransactionOutput{{Address: testutil.MakeAddress(), Coins: 501e6}},
	}
	_, _, err := CreateSplit(p, auxs, headTime)
	require.Equal(t, ErrInsufficientBalance, err)
}

func TestCreateSplitPQ(t *testing.T) {
	headTime := uint64(time.Now().UTC().Unix()) //nolint:gosec
	from := cipher.AddressFromPQPubKey([]byte("pq key"))
	change := testutil.MakeAddress()
	auxs := makeSplitUxOuts(t, from, headTime, oneCoin, 40)
	share := decimal.NewFromFloat(0.5)
	p := Params{
		HoursSelection: HoursSelection{Type: HoursSelectionTypeAuto, Mode: HoursSelectionModeShare, ShareFactor: &share},
		ChangeAddress:  &change,
		To:             []coin.TransactionOutput{{Address: testutil.MakeAddress(), Coins: 30e6}},
	}

	// Each post-quantum input adds about 3.7 KB once signed, so 30 of them
	// need several parts even though the unsigned transaction is small.
	_, _, err := Create(p, auxs, headTime)
	require.ErrorAs(t, err, new(TooLargeError))

	txns, inputs, err := CreateSplit(p, auxs, headTime)
	require.NoError(t, err)
	require.Greater(t, len(txns), 3)
	verifySplit(t, p, auxs, headTime, txns, inputs)
}

func TestTakeCoins(t *testing.T) {
	a, b := testutil.MakeAddress(), testutil.MakeAddress()
	to := []coin.TransactionOutput{
		{Address: a, Coins: 10e6, Hours: 100},
		{Address: b, Coins: 30e6, Hours: 7},
	}

	part, rest := takeCoins(to, 40e6)
	require.Equal(t, to, part)
	require.Empty(t, rest)

	part, rest = takeCoins(to, 25e6)
	require.Equal(t, []coin.TransactionOutput{to[0], {Address: b, Coins: 15e6, Hours: 3}}, part)
	require.Equal(t, []coin.TransactionOutput{{Address: b, Coins: 15e6, Hours: 4}}, rest)

	part, rest = takeCoins(to, 10e6)
	require.Equal(t, to[:1], part)
	require.Equal(t, to[1:], rest)

	part, rest = takeCoins(to, 5e6)
	require.Equal(t, []coin.TransactionOutput{{Address: a, Coins: 5e6, Hours: 50}}, part)
	require.Equal(t, []coin.TransactionOutput{{Address: a, Coins: 5e6, Hours: 50}, to[1]}, rest)
}
