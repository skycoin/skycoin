package sweep

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
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

func TestWalletAddresses(t *testing.T) {
	ws := newWallets(t)
	w := newBip44(t, ws, "b.wlt", nil)
	addrs, err := walletAddresses(w)
	require.NoError(t, err)
	a0, a1 := accountAddrs(t, w, 0), accountAddrs(t, w, 1)
	exp := Rule{SharedAccounts: []uint32{0}}.Exposed(addrs, nil)
	change, err := w.GetEntries(wallet.OptionAccount(0), wallet.OptionChange())
	require.NoError(t, err)
	require.NotEmpty(t, change)
	require.True(t, exp[change[0].Address.(cipher.Address)], "the change chain is listed with its account")
	require.Len(t, exp, len(a0)+len(change))
	for _, a := range a1 {
		require.False(t, exp[a])
	}

	_, err = ws.CreateWallet("d.wlt", wallet.Options{Label: "d", Seed: "d seed", Type: wallet.WalletTypeDeterministic, GenerateN: 3})
	require.NoError(t, err)
	d, err := ws.GetWallet("d.wlt")
	require.NoError(t, err)
	das, err := walletAddresses(d)
	require.NoError(t, err)
	require.Len(t, das, 3)
	require.False(t, das[0].Bip44)
}

func makeUx(t *testing.T, addr cipher.Address, coins uint64, seed byte) coin.UxOut {
	t.Helper()
	return coin.UxOut{
		Head: coin.UxHead{Time: 1000, BkSeq: 1},
		Body: coin.UxBody{SrcTransaction: cipher.SumSHA256([]byte{seed}), Address: addr, Coins: coins, Hours: 10},
	}
}

// fakeChain holds outputs and history for unit tests.
type fakeChain struct {
	mu      sync.Mutex
	seq     uint64
	time    uint64
	unspent []coin.UxOut
	sent    map[cipher.Address]bool
	pool    map[cipher.SHA256]bool
}

func (c *fakeChain) HeadBkSeq() (uint64, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.seq, true, nil
}

func (c *fakeChain) SweepState(addrs []cipher.Address, txids []cipher.SHA256) (ChainState, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
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
		c.chain.mu.Lock()
		defer c.chain.mu.Unlock()
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
	chain := &fakeChain{seq: 2, time: 1000 + 3600*10, unspent: []coin.UxOut{makeUx(t, a0[3], 7e6, 1)}}
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
	_, err = s.destination(r, r.Destination, &listing{w: w}, nil, nil)
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
	dest, err := s.destination(r, r.Destination, &l, nil, nil)
	require.NoError(t, err)
	w, err := ws.GetWallet("b.wlt")
	require.NoError(t, err)
	require.Len(t, w.Accounts(), 3, "account 2 was created")
	require.Equal(t, accountAddrs(t, w, 2)[0], dest)

	// A used destination address is followed by a new one.
	l.st.Seen = map[cipher.Address]bool{dest: true}
	l.w = w
	next, err := s.destination(r, r.Destination, &l, nil, nil)
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
}

func TestRunSweepsAtStartAndOnNewBlocks(t *testing.T) {
	ws := newWallets(t)
	pw := []byte("pwd")
	w := newBip44(t, ws, "e.wlt", pw)
	a0 := accountAddrs(t, w, 0)
	chain := &fakeChain{seq: 0, time: 1000 + 3600*10, unspent: []coin.UxOut{makeUx(t, a0[1], 2e6, 1)}}
	inj := &captureInjector{chain: chain}
	dir := t.TempDir()
	require.NoError(t, saveRules(dir, map[string]Rule{"e.wlt": {Wallet: "e.wlt", SharedAccounts: []uint32{0}, Destination: Destination{Kind: DestAccount}, Enabled: true}}))
	s, err := NewService(Config{Dir: dir, Poll: 5 * time.Millisecond}, ws, chain, inj)
	require.NoError(t, err)
	require.NoError(t, s.Unlock("e.wlt", pw, 0))

	quit := make(chan struct{})
	done := make(chan struct{})
	go func() {
		s.Run(quit)
		close(done)
	}()
	require.Eventually(t, func() bool { return len(s.Status("e.wlt").Pending) == 1 }, 5*time.Second, 5*time.Millisecond, "a rule loaded from disk is swept at start, on a chain still at block 0")

	chain.mu.Lock()
	chain.seq = 1
	chain.unspent = append(chain.unspent, makeUx(t, a0[2], 3e6, 2))
	chain.mu.Unlock()
	require.Eventually(t, func() bool { return len(s.Status("e.wlt").Pending) == 2 }, 5*time.Second, 5*time.Millisecond, "a new block triggers a pass")

	close(quit)
	<-done
	require.False(t, s.Status("e.wlt").Unlocked, "stopping forgets every password")
}

func TestMappingsWithLevel1Off(t *testing.T) {
	ws := newWallets(t)
	w := newBip44(t, ws, "b.wlt", nil)
	_, err := ws.CreateWallet("cold.wlt", wallet.Options{Label: "cold", Seed: "cold seed", Type: wallet.WalletTypeDeterministic, GenerateN: 2})
	require.NoError(t, err)
	cold, err := ws.GetWallet("cold.wlt")
	require.NoError(t, err)
	coldAddrs, err := walletAddresses(cold)
	require.NoError(t, err)
	a0, a1 := accountAddrs(t, w, 0), accountAddrs(t, w, 1)
	fixed := testutil.MakeAddress()

	chain := &fakeChain{seq: 1, time: 1000 + 3600*10, sent: map[cipher.Address]bool{a0[0]: true, a0[1]: true, a1[0]: true}}
	chain.unspent = []coin.UxOut{
		makeUx(t, a0[0], 1e6, 1), // sent, mapped to the fixed address
		makeUx(t, a0[1], 2e6, 2), // sent, not mapped
		makeUx(t, a1[0], 3e6, 3), // sent, covered by the account 1 mapping
		makeUx(t, a1[2], 4e6, 4), // never sent
	}
	inj := &captureInjector{chain: chain}
	s, err := NewService(Config{Dir: t.TempDir()}, ws, chain, inj)
	require.NoError(t, err)

	one := uint32(1)
	_, err = s.SetRule(Rule{Wallet: "b.wlt", Mappings: []Mapping{{From: Source{Address: a1[2].String()}, To: Destination{Kind: DestAddress, Address: fixed.String()}}}})
	require.ErrorContains(t, err, "never sent", "a mapping cannot sweep a single address that never sent")
	_, err = s.SetRule(Rule{Wallet: "b.wlt", Mappings: []Mapping{{From: Source{Address: a0[0].String()}, To: Destination{Kind: DestAddress, Address: a0[1].String()}}}})
	require.ErrorIs(t, err, ErrDestinationExposed, "a mapping destination that has sent")

	r, err := s.SetRule(Rule{Wallet: "b.wlt", Mappings: []Mapping{
		{From: Source{Address: a0[0].String()}, To: Destination{Kind: DestAddress, Address: fixed.String()}},
		{From: Source{Account: &one}, To: Destination{Kind: DestWallet, Wallet: "cold.wlt"}},
	}})
	require.NoError(t, err)
	require.False(t, r.Enabled)
	require.Equal(t, DestAccount, r.Destination.Kind, "level 1 gets its default even while off")

	s.Pass()
	st := s.Status("b.wlt")
	require.Empty(t, st.LastError)
	require.True(t, st.Active)
	require.False(t, st.Enabled)
	require.Equal(t, 3, st.Exposed)
	require.Equal(t, 1, st.Uncovered, "the unmapped address that sent is left alone")
	require.Len(t, inj.txns, 2)
	paid := map[cipher.Address]uint64{}
	for _, txn := range inj.txns {
		require.Len(t, txn.Out, 1)
		require.Zero(t, txn.Out[0].Hours)
		paid[txn.Out[0].Address] += txn.Out[0].Coins
	}
	require.Equal(t, map[cipher.Address]uint64{fixed: 1e6, coldAddrs[0].Address: 3e6}, paid)
	require.ElementsMatch(t, []string{fixed.String(), coldAddrs[0].Address.String()}, st.Destinations)

	// Level 1 takes the rest, to account 2, the first after the source account.
	inj.txns = nil
	r.Enabled = true
	_, err = s.SetRule(r)
	require.NoError(t, err)
	s.Pass()
	require.Len(t, inj.txns, 1)
	w, err = ws.GetWallet("b.wlt")
	require.NoError(t, err)
	require.Len(t, w.Accounts(), 3)
	require.Equal(t, accountAddrs(t, w, 2)[0], inj.txns[0].Out[0].Address)
	require.Equal(t, uint64(2e6), inj.txns[0].Out[0].Coins)
	require.Zero(t, s.Status("b.wlt").Uncovered)

	// Turning off the has-sent inclusion leaves only what is listed.
	f := false
	r.IncludeSpent = &f
	r.Mappings = r.Mappings[:1]
	_, err = s.SetRule(r)
	require.NoError(t, err)
	l, err := s.list(r, nil)
	require.NoError(t, err)
	require.Equal(t, map[cipher.Address]bool{a0[0]: true}, r.Exposed(l.addrs, l.st.Sent))
}
