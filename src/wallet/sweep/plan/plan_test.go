package plan

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
)

func u32(v uint32) *uint32 { return &v }
func off() *bool           { b := false; return &b }

// wallet is two bip44 accounts of four external and one change address each.
type wallet struct {
	addrs []Addr
	ext   [2][]cipher.Address
	chg   [2]cipher.Address
}

func newWallet() wallet {
	var w wallet
	for acct := range uint32(2) {
		for child := range uint32(4) {
			a := testutil.MakeAddress()
			w.ext[acct] = append(w.ext[acct], a)
			w.addrs = append(w.addrs, Addr{Address: a, Account: acct, Child: child, Bip44: true})
		}
		c := testutil.MakeAddress()
		w.chg[acct] = c
		w.addrs = append(w.addrs, Addr{Address: c, Account: acct, Change: true, Bip44: true})
	}
	return w
}

func set(as ...cipher.Address) map[cipher.Address]bool {
	m := make(map[cipher.Address]bool, len(as))
	for _, a := range as {
		m[a] = true
	}
	return m
}

func TestExposed(t *testing.T) {
	w := newWallet()

	t.Run("every address of a shared account, including those that never sent", func(t *testing.T) {
		exp := Rule{SharedAccounts: []uint32{0}}.Exposed(w.addrs, nil)
		require.Equal(t, set(append(w.ext[0], w.chg[0])...), exp, "the change chain is derivable from the account xpub too")
	})

	t.Run("an address that sent is exposed, one that never sent is not", func(t *testing.T) {
		exp := Rule{}.Exposed(w.addrs, set(w.ext[1][2]))
		require.Equal(t, set(w.ext[1][2]), exp)
	})

	t.Run("single addresses are exposed only once they sent", func(t *testing.T) {
		a, b := testutil.MakeAddress(), testutil.MakeAddress()
		das := []Addr{{Address: a}, {Address: b}}
		require.Empty(t, Rule{Addresses: []string{a.String()}}.Exposed(das, nil), "listing a single address that never sent does not expose it")
		require.Equal(t, set(b), Rule{}.Exposed(das, set(b)))
	})

	t.Run("a post-quantum address is never exposed", func(t *testing.T) {
		pq := testutil.MakeAddress()
		pq.Version = cipher.AddressVersionPQ
		require.Empty(t, Rule{}.Exposed([]Addr{{Address: pq}}, set(pq)))
	})
}

func TestIncludeSpentOff(t *testing.T) {
	w := newWallet()
	sent := set(w.ext[0][1], w.ext[0][2], w.ext[1][1], w.ext[1][3])

	require.Equal(t, sent, Rule{}.Exposed(w.addrs, sent), "by default every address that sent is swept")

	r := Rule{IncludeSpent: off()}
	require.Empty(t, r.Exposed(w.addrs, sent), "with nothing listed, nothing is swept")

	r.Addresses = []string{w.ext[0][1].String()}
	require.Equal(t, set(w.ext[0][1]), r.Exposed(w.addrs, sent))

	r.Mappings = []Mapping{
		{From: Source{Address: w.ext[0][2].String()}, To: Destination{Kind: DestWallet}},
		{From: Source{Account: u32(1)}, To: Destination{Kind: DestWallet}},
	}
	require.Equal(t, set(w.ext[0][1], w.ext[0][2], w.ext[1][1], w.ext[1][3]), r.Exposed(w.addrs, sent),
		"an account source covers the addresses of the account that sent, and only those")

	r = Rule{IncludeSpent: off(), SharedAccounts: []uint32{1}}
	require.Equal(t, set(append(w.ext[1], w.chg[1])...), r.Exposed(w.addrs, sent), "a shared account is listed")
}

func TestRoute(t *testing.T) {
	w := newWallet()
	sent := set(w.ext[0][0], w.ext[0][1], w.ext[1][0], w.ext[1][1])
	r := Rule{
		Mappings: []Mapping{
			{From: Source{Account: u32(1)}, To: Destination{Kind: DestWallet, Wallet: "cold"}},
			{From: Source{Address: w.ext[1][1].String()}, To: Destination{Kind: DestAddress, Address: testutil.MakeAddress().String()}},
		},
	}

	t.Run("level 1 off leaves unmapped addresses alone", func(t *testing.T) {
		routes, exposed, uncovered := r.Route(w.addrs, sent)
		require.Equal(t, 4, exposed)
		require.Equal(t, 2, uncovered)
		require.Equal(t, map[cipher.Address]int{w.ext[1][0]: 0, w.ext[1][1]: 1}, routes, "the address mapping wins over its account's")
	})

	t.Run("level 1 on takes what no mapping covers", func(t *testing.T) {
		r := r
		r.Enabled = true
		routes, _, uncovered := r.Route(w.addrs, sent)
		require.Zero(t, uncovered)
		require.Equal(t, map[cipher.Address]int{w.ext[0][0]: Level1, w.ext[0][1]: Level1, w.ext[1][0]: 0, w.ext[1][1]: 1}, routes)
		require.Equal(t, "cold", r.Target(0).Wallet)
		require.Equal(t, r.Destination, r.Target(Level1))
	})

	t.Run("include spent off with level 1 sweeps only what is listed", func(t *testing.T) {
		r := Rule{Enabled: true, IncludeSpent: off(), Mappings: []Mapping{{From: Source{Account: u32(1)}, To: Destination{Kind: DestWallet}}}}
		routes, exposed, uncovered := r.Route(w.addrs, sent)
		require.Equal(t, 2, exposed)
		require.Zero(t, uncovered)
		require.Equal(t, map[cipher.Address]int{w.ext[1][0]: 0, w.ext[1][1]: 0}, routes)
	})
}

func TestFind(t *testing.T) {
	w := newWallet()
	sent := set(w.ext[0][0], w.ext[1][0])
	r := Rule{Enabled: true, Mappings: []Mapping{{From: Source{Account: u32(1)}, To: Destination{Kind: DestWallet}}}}
	out := func(a cipher.Address, coins, hours uint64, seed byte, contested bool) Output {
		return Output{Hash: cipher.SumSHA256([]byte{seed}), Address: a, Coins: coins, Hours: hours, Contested: contested}
	}
	unspent := []Output{
		out(w.ext[0][0], 1e6, 10, 1, false),
		out(w.ext[0][0], 2e6, 20, 2, true),
		out(w.ext[1][0], 3e6, 30, 3, false),
		out(w.ext[0][3], 4e6, 40, 4, false), // never sent
		out(w.ext[1][0], 5e6, 50, 5, false),
	}
	p := r.Find(w.addrs, sent, unspent, map[cipher.SHA256]bool{unspent[4].Hash: true})
	require.Len(t, p.Groups, 2)
	require.Equal(t, Level1, p.Groups[0].Target)
	require.Equal(t, uint64(3e6), p.Groups[0].Coins)
	require.Equal(t, 0, p.Groups[1].Target)
	require.Equal(t, []Output{unspent[2]}, p.Groups[1].Outputs, "an output a pending sweep spends is left out")
	require.Equal(t, uint64(6e6), p.Coins)
	require.Equal(t, uint64(60), p.Hours)
	require.Equal(t, 1, p.Contested)
	require.Equal(t, 2, p.Exposed)
}

func TestCheckListed(t *testing.T) {
	w := newWallet()
	r := Rule{Wallet: "b.wlt", Addresses: []string{w.ext[1][0].String()}}
	require.ErrorContains(t, r.CheckListed(w.addrs, nil, 2), "never sent")
	require.NoError(t, r.CheckListed(w.addrs, set(w.ext[1][0]), 2))
	r.SharedAccounts = []uint32{1}
	require.NoError(t, r.CheckListed(w.addrs, nil, 2))

	r = Rule{Wallet: "b.wlt", Addresses: []string{testutil.MakeAddress().String()}}
	require.ErrorContains(t, r.CheckListed(w.addrs, nil, 2), "not in wallet")

	r = Rule{Wallet: "b.wlt", Mappings: []Mapping{{From: Source{Address: w.ext[0][3].String()}, To: Destination{Kind: DestWallet}}}}
	require.ErrorContains(t, r.CheckListed(w.addrs, nil, 2), "never sent", "a mapping source must be exposed too")

	r = Rule{Wallet: "b.wlt", Mappings: []Mapping{{From: Source{Account: u32(2)}, To: Destination{Kind: DestWallet}}}}
	require.ErrorContains(t, r.CheckListed(w.addrs, nil, 2), "no account 2")
	require.ErrorContains(t, Rule{Wallet: "b.wlt", SharedAccounts: []uint32{5}}.CheckListed(w.addrs, nil, 2), "no account 5")
}

func TestRuleCheck(t *testing.T) {
	ok := Rule{Wallet: "w", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAccount, Account: u32(1)}}
	require.NoError(t, ok.Check())

	bad := ok
	bad.Destination.Account = u32(0)
	require.ErrorContains(t, bad.Check(), "is shared")

	bad = Rule{Wallet: "w", SharedAccounts: []uint32{0, 2}, Destination: Destination{Kind: DestAccount, Account: u32(1)}}
	require.ErrorContains(t, bad.Check(), "must come after shared account 2")

	require.ErrorContains(t, Rule{Wallet: "w", Destination: Destination{Kind: "x"}}.Check(), "unknown destination")
	require.ErrorContains(t, Rule{Wallet: "w", Destination: Destination{Kind: DestAddress, Address: "nope"}}.Check(), "invalid destination")

	addr := testutil.MakeAddress().String()
	wallet := Destination{Kind: DestWallet}
	for name, tc := range map[string]struct {
		m   []Mapping
		err string
	}{
		"address source":         {m: []Mapping{{From: Source{Address: addr}, To: wallet}}},
		"account source":         {m: []Mapping{{From: Source{Account: u32(2)}, To: Destination{Kind: DestAccount, Account: u32(3)}}}},
		"empty source":           {m: []Mapping{{To: wallet}}, err: "source is empty"},
		"both":                   {m: []Mapping{{From: Source{Address: addr, Account: u32(1)}, To: wallet}}, err: "not both"},
		"bad source":             {m: []Mapping{{From: Source{Address: "nope"}, To: wallet}}, err: "invalid source"},
		"address twice":          {m: []Mapping{{From: Source{Address: addr}, To: wallet}, {From: Source{Address: addr}, To: wallet}}, err: "mapped twice"},
		"account twice":          {m: []Mapping{{From: Source{Account: u32(1)}, To: wallet}, {From: Source{Account: u32(1)}, To: wallet}}, err: "mapped twice"},
		"to a source account":    {m: []Mapping{{From: Source{Account: u32(1)}, To: Destination{Kind: DestAccount, Account: u32(1)}}}, err: "is a source"},
		"to a shared account":    {m: []Mapping{{From: Source{Address: addr}, To: Destination{Kind: DestAccount, Account: u32(0)}}}, err: "mapping 1: destination account 0 is shared"},
		"no destination kind":    {m: []Mapping{{From: Source{Address: addr}}}, err: "unknown destination kind"},
		"bad destination":        {m: []Mapping{{From: Source{Address: addr}, To: Destination{Kind: DestAddress, Address: "x"}}}, err: "invalid destination"},
		"wallet with an account": {m: []Mapping{{From: Source{Address: addr}, To: Destination{Kind: DestWallet, Account: u32(4)}}}, err: "takes only a wallet"},
	} {
		t.Run(name, func(t *testing.T) {
			r := Rule{Wallet: "w", SharedAccounts: []uint32{0}, Destination: wallet, Mappings: tc.m}
			if tc.err == "" {
				require.NoError(t, r.Check())
			} else {
				require.ErrorContains(t, r.Check(), tc.err)
			}
		})
	}
}

func TestDestinationAccount(t *testing.T) {
	next := Destination{Kind: DestAccount}
	require.Zero(t, Rule{}.DestinationAccount(next))
	require.Equal(t, uint32(3), Rule{SharedAccounts: []uint32{2, 0}}.DestinationAccount(next))
	r := Rule{SharedAccounts: []uint32{0}, Mappings: []Mapping{{From: Source{Account: u32(4)}, To: Destination{Kind: DestWallet}}}}
	require.Equal(t, uint32(5), r.DestinationAccount(next), "the next account comes after the source accounts too")
	require.Equal(t, uint32(7), r.DestinationAccount(Destination{Kind: DestAccount, Account: u32(7)}))
}

func TestActive(t *testing.T) {
	require.False(t, Rule{}.Active())
	require.True(t, Rule{Enabled: true}.Active())
	require.True(t, Rule{Mappings: []Mapping{{}}}.Active(), "mappings sweep with level 1 off")
	require.True(t, Rule{}.SweepsSpent())
	require.False(t, Rule{IncludeSpent: off()}.SweepsSpent())
}

func TestGapNeeds(t *testing.T) {
	w := newWallet()
	r := Rule{SharedAccounts: []uint32{0}, Gap: 3}
	require.Empty(t, r.GapNeeds(w.addrs, nil), "four unused addresses cover a gap of three")
	require.Equal(t, map[uint32]uint64{0: 2}, r.GapNeeds(w.addrs, set(w.ext[0][2])))
	require.Equal(t, map[uint32]uint64{0: 20}, Rule{SharedAccounts: []uint32{0}}.GapNeeds(w.addrs, set(w.ext[0][3])))
}

func TestCheckDestination(t *testing.T) {
	a, b := testutil.MakeAddress(), testutil.MakeAddress()
	require.NoError(t, CheckDestination(a, set(b), nil))
	require.ErrorIs(t, CheckDestination(a, set(a), nil), ErrDestinationExposed)
	require.ErrorIs(t, CheckDestination(a, nil, set(a)), ErrDestinationExposed)
	require.Error(t, CheckDestination(cipher.Address{}, nil, nil))
}

func makeUx(addr cipher.Address, coins, hours uint64, seed byte) coin.UxOut {
	return coin.UxOut{
		Head: coin.UxHead{Time: 1000, BkSeq: 1},
		Body: coin.UxBody{SrcTransaction: cipher.SumSHA256([]byte{seed}), Address: addr, Coins: coins, Hours: hours},
	}
}

func outputs(t *testing.T, uxs coin.UxArray) []Output {
	t.Helper()
	out := make([]Output, len(uxs))
	for i, u := range uxs {
		h, err := u.CoinHours(1000)
		require.NoError(t, err)
		out[i] = Output{Hash: u.Hash(), Address: u.Body.Address, Coins: u.Body.Coins, Hours: h}
	}
	return out
}

func TestBuildBurnsEveryHour(t *testing.T) {
	to, from := testutil.MakeAddress(), testutil.MakeAddress()
	uxs := coin.UxArray{makeUx(from, 2e6, 50, 1), makeUx(from, 3e6, 70, 2), makeUx(from, 1e6, 0, 3)}

	plan, err := Build(outputs(t, uxs), to, params.UserVerifyTxn.MaxTransactionSize)
	require.NoError(t, err)
	require.Len(t, plan, 1)
	sw := plan[0]
	require.Len(t, sw.Txn.In, 3, "the output without hours joins the last transaction")
	require.Len(t, sw.Txn.Out, 1)
	require.Equal(t, uint64(6e6), sw.Txn.Out[0].Coins)
	require.Zero(t, sw.Txn.Out[0].Hours)
	require.Equal(t, uint64(120), sw.Fee)

	// The verifier accepts burning every hour, at any burn factor.
	for _, bf := range []uint32{params.MinBurnFactor, params.UserVerifyTxn.BurnFactor, 10, 1000} {
		require.NoError(t, fee.VerifyTransactionFee(sw.Txn, sw.Fee, bf))
	}
	byHash := make(map[cipher.SHA256]coin.UxOut)
	for _, u := range uxs {
		byHash[u.Hash()] = u
	}
	in := make(coin.UxArray, len(sw.Txn.In))
	for i, h := range sw.Txn.In {
		in[i] = byHash[h]
	}
	txnFee, err := fee.TransactionFee(sw.Txn, 1000, in)
	require.NoError(t, err)
	require.Equal(t, uint64(120), txnFee)
	require.NoError(t, sw.Txn.VerifyUnsigned(), "the unsigned shape is a valid transaction body")
}

func TestBuildNoHours(t *testing.T) {
	from := testutil.MakeAddress()
	_, err := Build(outputs(t, coin.UxArray{makeUx(from, 1e6, 0, 1)}), testutil.MakeAddress(), params.UserVerifyTxn.MaxTransactionSize)
	require.ErrorIs(t, err, ErrNoHours)
}

// TestBuildNoSubsetPaysMore checks that no subset of a transaction's inputs pays
// more per byte than the transaction, so an attacker who spends some of the same
// outputs from their own hours cannot outbid the sweep.
func TestBuildNoSubsetPaysMore(t *testing.T) {
	from, to := testutil.MakeAddress(), testutil.MakeAddress()
	hours := []uint64{5000, 4900, 30, 20, 4800, 1, 2, 600}
	var uxs coin.UxArray
	for i, h := range hours {
		uxs = append(uxs, makeUx(from, uint64(i+1)*1e6, h, byte(i))) //nolint:gosec // test value
	}
	plan, err := Build(outputs(t, uxs), to, params.UserVerifyTxn.MaxTransactionSize)
	require.NoError(t, err)
	require.Greater(t, len(plan), 1, "the low hour outputs go in their own transaction")

	sizes := sizeTable{to: to}
	var swept int
	for _, sw := range plan {
		swept += len(sw.Inputs)
		size, err := sw.Txn.Size()
		require.NoError(t, err)
		n := len(sw.Inputs)
		for mask := 1; mask < 1<<n; mask++ {
			var h uint64
			var k int
			for i := range n {
				if mask&(1<<i) != 0 {
					h += sw.Inputs[i].Hours
					k++
				}
			}
			ks, err := sizes.size(k)
			require.NoError(t, err)
			require.False(t, ratioLess(sw.Fee, size, h, ks), "subset %b pays more per byte than the sweep", mask)
		}
	}
	require.Equal(t, len(uxs), swept)
}

func TestBuildSplitsBySize(t *testing.T) {
	from, to := testutil.MakeAddress(), testutil.MakeAddress()
	var outs []Output
	for i := range 700 {
		outs = append(outs, Output{Hash: cipher.SumSHA256([]byte{byte(i), byte(i >> 8)}), Address: from, Coins: 1e6, Hours: 100})
	}
	plan, err := Build(outs, to, params.UserVerifyTxn.MaxTransactionSize)
	require.NoError(t, err)
	require.Greater(t, len(plan), 1)
	seen := make(map[cipher.SHA256]bool)
	for _, sw := range plan {
		size, err := sw.Txn.Size()
		require.NoError(t, err)
		require.LessOrEqual(t, size, params.UserVerifyTxn.MaxTransactionSize)
		for _, h := range sw.Txn.In {
			require.False(t, seen[h])
			seen[h] = true
		}
	}
	require.Len(t, seen, 700)
}
