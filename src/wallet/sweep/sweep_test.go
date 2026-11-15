package sweep

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/wallet"
	_ "github.com/skycoin/skycoin/src/wallet/bip44wallet"   // registers the bip44 wallet type
	_ "github.com/skycoin/skycoin/src/wallet/collection"    // registers the collection wallet type
	_ "github.com/skycoin/skycoin/src/wallet/deterministic" // registers the deterministic wallet type
)

const testSeed = "voyage say extend find sheriff surge priority merit ignore maple cash argue"

func newWallets(t *testing.T) *wallet.Service {
	ws, err := wallet.NewService(wallet.Config{
		WalletDir:       t.TempDir(),
		CryptoType:      crypto.CryptoTypeScryptChacha20poly1305Insecure,
		EnableWalletAPI: true,
	})
	require.NoError(t, err)
	return ws
}

// newBip44 makes a bip44 wallet with two accounts of five external addresses each.
func newBip44(t *testing.T, ws *wallet.Service, id string, password []byte) wallet.Wallet {
	_, err := ws.CreateWallet(id, wallet.Options{
		Label: id, Seed: testSeed, Type: wallet.WalletTypeBip44, GenerateN: 5,
		Encrypt: len(password) > 0, Password: password,
	})
	require.NoError(t, err)
	_, err = ws.NewAccount(id, password, "second")
	require.NoError(t, err)
	_, err = ws.NewAddresses(id, nil, wallet.OptionAccount(1), wallet.OptionExternal(), wallet.OptionGenerateN(4))
	require.NoError(t, err)
	w, err := ws.GetWallet(id)
	require.NoError(t, err)
	return w
}

func accountAddrs(t *testing.T, w wallet.Wallet, acct uint32) []cipher.Address {
	es, err := w.GetEntries(wallet.OptionAccount(acct), wallet.OptionExternal())
	require.NoError(t, err)
	out := make([]cipher.Address, len(es))
	for i, e := range es {
		out[i] = e.Address.(cipher.Address)
	}
	return out
}

func TestExposed(t *testing.T) {
	ws := newWallets(t)
	w := newBip44(t, ws, "b.wlt", nil)
	addrs, err := walletAddresses(w)
	require.NoError(t, err)
	a0, a1 := accountAddrs(t, w, 0), accountAddrs(t, w, 1)

	t.Run("every address of a shared account, including those that never sent", func(t *testing.T) {
		exp := exposed(Rule{SharedAccounts: []uint32{0}}, addrs, nil)
		change, err := w.GetEntries(wallet.OptionAccount(0), wallet.OptionChange())
		require.NoError(t, err)
		require.NotEmpty(t, change)
		require.True(t, exp[change[0].Address.(cipher.Address)], "the change chain is derivable from the account xpub too")
		require.Len(t, exp, len(a0)+len(change))
		for _, a := range a0 {
			require.True(t, exp[a])
		}
		for _, a := range a1 {
			require.False(t, exp[a])
		}
	})

	t.Run("an address that sent is exposed, one that never sent is not", func(t *testing.T) {
		exp := exposed(Rule{}, addrs, map[cipher.Address]bool{a1[2]: true})
		require.Equal(t, map[cipher.Address]bool{a1[2]: true}, exp)
	})

	t.Run("a hand added address must have sent or be shared", func(t *testing.T) {
		r := Rule{Wallet: "b.wlt", Addresses: []string{a1[0].String()}}
		require.ErrorContains(t, checkManual(r, addrs, nil), "never sent")
		require.NoError(t, checkManual(r, addrs, map[cipher.Address]bool{a1[0]: true}))
		r = Rule{Wallet: "b.wlt", SharedAccounts: []uint32{1}, Addresses: []string{a1[0].String()}}
		require.NoError(t, checkManual(r, addrs, nil))
		r = Rule{Wallet: "b.wlt", Addresses: []string{testutil.MakeAddress().String()}}
		require.ErrorContains(t, checkManual(r, addrs, nil), "not in wallet")
	})

	t.Run("deterministic wallet addresses are exposed only once they sent", func(t *testing.T) {
		_, err := ws.CreateWallet("d.wlt", wallet.Options{Label: "d", Seed: "d seed", Type: wallet.WalletTypeDeterministic, GenerateN: 3})
		require.NoError(t, err)
		d, err := ws.GetWallet("d.wlt")
		require.NoError(t, err)
		das, err := walletAddresses(d)
		require.NoError(t, err)
		require.Len(t, das, 3)
		require.Empty(t, exposed(Rule{}, das, nil))
		exp := exposed(Rule{}, das, map[cipher.Address]bool{das[1].addr: true})
		require.Equal(t, map[cipher.Address]bool{das[1].addr: true}, exp)
	})
}

func TestGapNeeds(t *testing.T) {
	ws := newWallets(t)
	w := newBip44(t, ws, "b.wlt", nil)
	addrs, err := walletAddresses(w)
	require.NoError(t, err)
	a0 := accountAddrs(t, w, 0)
	r := Rule{SharedAccounts: []uint32{0}, Gap: 3}

	require.Empty(t, gapNeeds(r, addrs, nil), "five unused addresses cover a gap of three")
	require.Equal(t, map[uint32]uint64{0: 2}, gapNeeds(r, addrs, map[cipher.Address]bool{a0[3]: true}))
	require.Equal(t, map[uint32]uint64{0: 20}, gapNeeds(Rule{SharedAccounts: []uint32{0}}, addrs, map[cipher.Address]bool{a0[4]: true}))
}

func TestRuleCheck(t *testing.T) {
	one, zero := uint32(1), uint32(0)
	ok := Rule{Wallet: "w", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAccount, Account: &one}}
	require.NoError(t, ok.check())

	bad := ok
	bad.Destination.Account = &zero
	require.ErrorContains(t, bad.check(), "is shared")

	bad = Rule{Wallet: "w", SharedAccounts: []uint32{0, 2}, Destination: Destination{Kind: DestAccount, Account: &one}}
	require.ErrorContains(t, bad.check(), "must come after shared account 2")

	require.ErrorContains(t, Rule{Wallet: "w", Destination: Destination{Kind: "x"}}.check(), "unknown destination")
	require.ErrorContains(t, Rule{Wallet: "w", Destination: Destination{Kind: DestAddress, Address: "nope"}}.check(), "invalid destination")
}

func TestCheckDestination(t *testing.T) {
	a, b := testutil.MakeAddress(), testutil.MakeAddress()
	require.NoError(t, checkDestination(a, map[cipher.Address]bool{b: true}, nil))
	require.ErrorIs(t, checkDestination(a, map[cipher.Address]bool{a: true}, nil), ErrDestinationExposed)
	require.ErrorIs(t, checkDestination(a, nil, map[cipher.Address]bool{a: true}), ErrDestinationExposed)
}

func makeUx(t *testing.T, addr cipher.Address, coins, hours uint64, seed byte) coin.UxOut {
	t.Helper()
	return coin.UxOut{
		Head: coin.UxHead{Time: 1000, BkSeq: 1},
		Body: coin.UxBody{SrcTransaction: cipher.SumSHA256([]byte{seed}), Address: addr, Coins: coins, Hours: hours},
	}
}

func TestPlanBurnsEveryHour(t *testing.T) {
	to := testutil.MakeAddress()
	from := testutil.MakeAddress()
	uxs := coin.UxArray{makeUx(t, from, 2e6, 50, 1), makeUx(t, from, 3e6, 70, 2), makeUx(t, from, 1e6, 0, 3)}

	plan, err := Plan(uxs, 1000, to, params.UserVerifyTxn.MaxTransactionSize)
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
	inHours, err := sw.Inputs.CoinHours(1000)
	require.NoError(t, err)
	txnFee, err := fee.TransactionFee(sw.Txn, 1000, sw.Inputs)
	require.NoError(t, err)
	require.Equal(t, inHours, txnFee)
	require.NoError(t, sw.Txn.VerifyUnsigned(), "the unsigned shape is a valid transaction body")
}

func TestPlanNoHours(t *testing.T) {
	from := testutil.MakeAddress()
	_, err := Plan(coin.UxArray{makeUx(t, from, 1e6, 0, 1)}, 1000, testutil.MakeAddress(), params.UserVerifyTxn.MaxTransactionSize)
	require.ErrorIs(t, err, ErrNoHours)
}

// TestPlanNoSubsetPaysMore checks that no subset of a transaction's inputs pays
// more per byte than the transaction, so an attacker who spends some of the same
// outputs from their own hours cannot outbid the sweep.
func TestPlanNoSubsetPaysMore(t *testing.T) {
	from, to := testutil.MakeAddress(), testutil.MakeAddress()
	hours := []uint64{5000, 4900, 30, 20, 4800, 1, 2, 600}
	var uxs coin.UxArray
	for i, h := range hours {
		uxs = append(uxs, makeUx(t, from, uint64(i+1)*1e6, h, byte(i))) //nolint:gosec // test value
	}
	plan, err := Plan(uxs, 1000, to, params.UserVerifyTxn.MaxTransactionSize)
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
					c, err := sw.Inputs[i].CoinHours(1000)
					require.NoError(t, err)
					h += c
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

func TestPlanSplitsBySize(t *testing.T) {
	from, to := testutil.MakeAddress(), testutil.MakeAddress()
	var uxs coin.UxArray
	for i := range 700 {
		uxs = append(uxs, coin.UxOut{
			Head: coin.UxHead{Time: 1000, BkSeq: 1},
			Body: coin.UxBody{SrcTransaction: cipher.SumSHA256([]byte{byte(i), byte(i >> 8)}), Address: from, Coins: 1e6, Hours: 100},
		})
	}
	plan, err := Plan(uxs, 1000, to, params.UserVerifyTxn.MaxTransactionSize)
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

// fakeChain holds outputs and history for unit tests.
type fakeChain struct {
	seq     uint64
	time    uint64
	unspent []coin.UxOut
	sent    map[cipher.Address]bool
	pool    map[cipher.SHA256]bool
}

func (c *fakeChain) HeadBkSeq() (uint64, bool, error) { return c.seq, true, nil }

func (c *fakeChain) SweepState(addrs []cipher.Address, txids []cipher.SHA256) (ChainState, error) {
	want := make(map[cipher.Address]bool)
	for _, a := range addrs {
		want[a] = true
	}
	st := ChainState{HeadSeq: c.seq, HeadTime: c.time, Sent: map[cipher.Address]bool{}, Seen: map[cipher.Address]bool{}, Txns: map[cipher.SHA256]TxnState{}}
	for _, u := range c.unspent {
		if want[u.Body.Address] {
			st.Unspent = append(st.Unspent, UxOut{UxOut: u})
			st.Seen[u.Body.Address] = true
		}
	}
	for a := range c.sent {
		if want[a] {
			st.Sent[a], st.Seen[a] = true, true
		}
	}
	for _, h := range txids {
		if c.pool[h] {
			st.Txns[h] = TxnPending
		}
	}
	return st, nil
}

type captureInjector struct {
	chain *fakeChain
	txns  []coin.Transaction
}

func (c *captureInjector) InjectBroadcastTransaction(txn coin.Transaction) error {
	c.txns = append(c.txns, txn)
	if c.chain != nil {
		if c.chain.pool == nil {
			c.chain.pool = map[cipher.SHA256]bool{}
		}
		c.chain.pool[txn.Hash()] = true
	}
	return nil
}

func TestLockedWalletPauses(t *testing.T) {
	ws := newWallets(t)
	pw := []byte("pwd")
	w := newBip44(t, ws, "e.wlt", pw)
	a0 := accountAddrs(t, w, 0)
	chain := &fakeChain{seq: 2, time: 1000 + 3600*10, unspent: []coin.UxOut{makeUx(t, a0[3], 7e6, 10, 1)}}
	inj := &captureInjector{chain: chain}
	dir := t.TempDir()
	s, err := NewService(Config{Dir: dir}, ws, chain, inj)
	require.NoError(t, err)

	_, err = s.SetRule(Rule{Wallet: "e.wlt", SharedAccounts: []uint32{0}, Enabled: true})
	require.NoError(t, err)
	s.Pass()
	st := s.Status("e.wlt")
	require.Equal(t, "paused: wallet locked, 7.000000 coins waiting on exposed addresses", st.Paused)
	require.Equal(t, uint64(7e6), st.WaitingCoins)
	require.False(t, st.Unlocked)
	require.Empty(t, inj.txns)

	require.Error(t, s.Unlock("e.wlt", []byte("wrong"), time.Minute))
	require.NoError(t, s.Unlock("e.wlt", pw, time.Minute))
	require.True(t, s.Status("e.wlt").Unlocked)
	s.Pass()
	st = s.Status("e.wlt")
	require.Empty(t, st.Paused)
	require.Empty(t, st.LastError)
	require.Len(t, inj.txns, 1)
	require.Len(t, st.Pending, 1)

	// The sweep went to account 1, the first account after the shared one.
	w, err = ws.GetWallet("e.wlt")
	require.NoError(t, err)
	require.Equal(t, accountAddrs(t, w, 1)[0], inj.txns[0].Out[0].Address)
	require.True(t, inj.txns[0].FullySigned())
	require.NoError(t, inj.txns[0].VerifyInputSignatures(coin.UxArray{chain.unspent[0]}))

	// A second pass while the sweep is pending sends nothing new.
	s.Pass()
	require.Len(t, inj.txns, 1)

	// The password expires.
	s.now = func() time.Time { return time.Now().Add(2 * time.Minute) }
	require.False(t, s.Status("e.wlt").Unlocked)
	_, ok := s.password("e.wlt")
	require.False(t, ok)

	// Nothing about the password reaches the rules file.
	rules, err := loadRules(dir)
	require.NoError(t, err)
	require.Len(t, rules, 1)
}

func TestSetRuleRefusesExposedDestination(t *testing.T) {
	ws := newWallets(t)
	w := newBip44(t, ws, "b.wlt", nil)
	a0, a1 := accountAddrs(t, w, 0), accountAddrs(t, w, 1)
	spent := testutil.MakeAddress()
	chain := &fakeChain{seq: 1, time: 1000, sent: map[cipher.Address]bool{spent: true}}
	s, err := NewService(Config{Dir: t.TempDir()}, ws, chain, &captureInjector{})
	require.NoError(t, err)

	_, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAddress, Address: a0[2].String()}})
	require.ErrorIs(t, err, ErrDestinationExposed, "an address of the shared account")
	_, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAddress, Address: spent.String()}})
	require.ErrorIs(t, err, ErrDestinationExposed, "an address that has sent")
	zero := uint32(0)
	_, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAccount, Account: &zero}})
	require.Error(t, err, "the shared account itself")
	_, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestWallet}})
	require.NoError(t, err, "accepted, but a wallet destination of the same bip44 wallet is account 0")
	r, _ := s.Rule("b.wlt")
	_, err = s.destination(r, listing{w: w}, nil)
	require.ErrorIs(t, err, ErrDestinationExposed)

	r, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAddress, Address: a1[1].String()}})
	require.NoError(t, err)
	require.Equal(t, DestAddress, r.Destination.Kind)

	r, err = s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0}})
	require.NoError(t, err)
	require.Equal(t, DestAccount, r.Destination.Kind, "a bip44 wallet defaults to another account")
}

func TestAccountDestinationCreatesAccount(t *testing.T) {
	ws := newWallets(t)
	newBip44(t, ws, "b.wlt", nil)
	s, err := NewService(Config{Dir: t.TempDir()}, ws, &fakeChain{seq: 1, time: 1000}, &captureInjector{})
	require.NoError(t, err)
	r, err := s.SetRule(Rule{Wallet: "b.wlt", SharedAccounts: []uint32{0, 1}, Enabled: true})
	require.NoError(t, err)

	l, err := s.list(r, nil)
	require.NoError(t, err)
	dest, err := s.destination(r, l, nil)
	require.NoError(t, err)
	w, err := ws.GetWallet("b.wlt")
	require.NoError(t, err)
	require.Len(t, w.Accounts(), 3, "account 2 was created")
	require.Equal(t, accountAddrs(t, w, 2)[0], dest)

	// A used destination address is followed by a new one.
	l.st.Seen = map[cipher.Address]bool{dest: true}
	l.w = w
	next, err := s.destination(r, l, nil)
	require.NoError(t, err)
	require.NotEqual(t, dest, next)
	w, err = ws.GetWallet("b.wlt")
	require.NoError(t, err)
	require.Equal(t, accountAddrs(t, w, 2)[1], next)
}

func TestRulesFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	one := uint32(1)
	in := map[string]Rule{
		"a.wlt": {Wallet: "a.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAccount, Account: &one}, Enabled: true},
		"b.wlt": {Wallet: "b.wlt", Destination: Destination{Kind: DestWallet}},
	}
	require.NoError(t, saveRules(dir, in))
	out, err := loadRules(dir)
	require.NoError(t, err)
	require.Equal(t, in, out)

	empty, err := loadRules(t.TempDir())
	require.NoError(t, err)
	require.Empty(t, empty)
	require.False(t, errors.Is(err, ErrNoHours))
}
