package sweep_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/crypto"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/wallet"
	_ "github.com/skycoin/skycoin/src/wallet/bip44wallet"   // registers the bip44 wallet type
	_ "github.com/skycoin/skycoin/src/wallet/deterministic" // registers the deterministic wallet type
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

const genesisTime uint64 = 1e9

type node struct {
	t       *testing.T
	v       *visor.Visor
	ws      *wallet.Service
	sk      cipher.SecKey
	genesis cipher.Address
	when    uint64
}

func newNode(t *testing.T) *node {
	pk, sk := cipher.GenerateKeyPair()
	db, err := visor.OpenDB(filepath.Join(t.TempDir(), "data.db"), false)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() }) //nolint:errcheck,gosec

	c := visor.NewConfig()
	c.BlockPublisher = true
	c.Arbitrating = true
	c.BlockchainPubkey, c.BlockchainSeckey = pk, sk
	c.GenesisAddress = cipher.AddressFromPubKey(pk)
	c.GenesisCoinVolume = 1000e6
	c.GenesisTimestamp = genesisTime
	c.Distribution = params.MainNetDistribution

	ws, err := wallet.NewService(wallet.Config{
		WalletDir:       t.TempDir(),
		CryptoType:      crypto.CryptoTypeScryptChacha20poly1305Insecure,
		EnableWalletAPI: true,
	})
	require.NoError(t, err)
	v, err := visor.New(c, db, ws)
	require.NoError(t, err)
	require.NoError(t, v.Init())
	return &node{t: t, v: v, ws: ws, sk: sk, genesis: c.GenesisAddress, when: genesisTime}
}

// InjectBroadcastTransaction stands in for the daemon, which has no peers here.
func (n *node) InjectBroadcastTransaction(txn coin.Transaction) error {
	_, _, _, err := n.v.InjectUserTransaction(txn)
	return err
}

// block publishes the unconfirmed pool as the next block, an hour after the last.
func (n *node) block() coin.Block {
	utxns, err := n.v.GetAllUnconfirmedTransactions()
	require.NoError(n.t, err)
	txns := make(coin.Transactions, len(utxns))
	for i, u := range utxns {
		txns[i] = u.Transaction
	}
	n.when += 3600
	b, err := n.v.CreateBlockFromTxns(txns, n.when)
	require.NoError(n.t, err)
	sb := coin.SignedBlock{Block: b, Sig: cipher.MustSignHash(b.HashHeader(), n.sk)}
	require.NoError(n.t, n.v.ExecuteSignedBlock(sb))
	return b
}

func (n *node) unspentOf(a cipher.Address) coin.UxArray {
	st, err := n.v.SweepState([]cipher.Address{a}, nil)
	require.NoError(n.t, err)
	var out coin.UxArray
	for _, u := range st.Unspent {
		out = append(out, u.UxOut)
	}
	return out
}

// pay sends from the genesis address and confirms it in a block.
func (n *node) pay(to []coin.TransactionOutput) {
	uxs := n.unspentOf(n.genesis)
	require.Len(n.t, uxs, 1)
	head, err := n.v.GetHeadBlock()
	require.NoError(n.t, err)
	inHours, err := uxs[0].CoinHours(head.Time())
	require.NoError(n.t, err)
	var txn coin.Transaction
	require.NoError(n.t, txn.PushInput(uxs[0].Hash()))
	var coins, hours uint64
	for _, o := range to {
		require.NoError(n.t, txn.PushOutput(o.Address, o.Coins, o.Hours))
		coins += o.Coins
		hours += o.Hours
	}
	require.NoError(n.t, txn.PushOutput(n.genesis, uxs[0].Body.Coins-coins, (inHours-hours)/2))
	require.NoError(n.t, txn.SignInputs([]cipher.SecKey{n.sk}))
	require.NoError(n.t, txn.UpdateHeader())
	_, _, _, err = n.v.InjectUserTransaction(txn)
	require.NoError(n.t, err)
	n.block()
}

func (n *node) poolHashes() map[cipher.SHA256]bool {
	utxns, err := n.v.GetAllUnconfirmedTransactions()
	require.NoError(n.t, err)
	out := make(map[cipher.SHA256]bool)
	for _, u := range utxns {
		out[u.Transaction.Hash()] = true
	}
	return out
}

// steal is what an attacker with the private key of ux sends: one output to
// thief, keeping hours on it, so it burns less than everything.
func steal(t *testing.T, ux coin.UxOut, sk cipher.SecKey, thief cipher.Address, keepHours uint64) coin.Transaction {
	var txn coin.Transaction
	require.NoError(t, txn.PushInput(ux.Hash()))
	require.NoError(t, txn.PushOutput(thief, ux.Body.Coins, keepHours))
	require.NoError(t, txn.SignInputs([]cipher.SecKey{sk}))
	require.NoError(t, txn.UpdateHeader())
	return txn
}

func entries(t *testing.T, w wallet.Wallet, acct uint32) wallet.Entries {
	es, err := w.GetEntries(wallet.OptionAccount(acct), wallet.OptionExternal())
	require.NoError(t, err)
	return es
}

func TestSweepOnRealChain(t *testing.T) {
	n := newNode(t)
	const id = "reward.wlt"
	_, err := n.ws.CreateWallet(id, wallet.Options{
		Label: "reward", Seed: "voyage say extend find sheriff surge priority merit ignore maple cash argue",
		Type: wallet.WalletTypeBip44, GenerateN: 5,
	})
	require.NoError(t, err)
	_, err = n.ws.NewAccount(id, nil, "savings")
	require.NoError(t, err)
	_, err = n.ws.NewAddresses(id, nil, wallet.OptionAccount(1), wallet.OptionExternal(), wallet.OptionGenerateN(2))
	require.NoError(t, err)
	w, err := n.ws.GetWallet(id)
	require.NoError(t, err)
	shared, savings := entries(t, w, 0), entries(t, w, 1)
	addr := func(e wallet.Entry) cipher.Address { return e.Address.(cipher.Address) }

	// Day payouts land on never-sent children of the shared account. Another
	// payout lands on a never-sent single address of account 1, which must stay.
	n.pay([]coin.TransactionOutput{
		{Address: addr(shared[1]), Coins: 10e6, Hours: 400},
		{Address: addr(shared[3]), Coins: 5e6, Hours: 300},
		{Address: addr(savings[1]), Coins: 2e6, Hours: 100},
	})
	payouts := append(n.unspentOf(addr(shared[1])), n.unspentOf(addr(shared[3]))...)
	require.Len(t, payouts, 2)

	// A prepared attacker is first to spend one payout, but burns less than all of its hours.
	thief := testutil.MakeAddress()
	attack := steal(t, n.unspentOf(addr(shared[1]))[0], shared[1].Secret, thief, 100)
	require.NoError(t, n.InjectBroadcastTransaction(attack))

	s, err := sweep.NewService(sweep.Config{Dir: t.TempDir()}, n.ws, n.v, n)
	require.NoError(t, err)
	_, err = s.SetRule(sweep.Rule{Wallet: id, SharedAccounts: []uint32{0}, Enabled: true})
	require.NoError(t, err)

	// Coins already sitting on exposed addresses are swept as soon as the rule is enabled.
	s.Pass()
	st := s.Status(id)
	require.Empty(t, st.LastError)
	require.NotEmpty(t, st.Pending)
	require.Equal(t, 1, st.Contested)
	require.Equal(t, uint64(15e6), st.WaitingCoins)
	dest := cipher.MustDecodeBase58Address(st.Destination)
	require.Equal(t, addr(savings[2]), dest, "the next unused address of account 1, past savings[1] which received")
	pool := n.poolHashes()
	require.Len(t, pool, 1+len(st.Pending))

	// No double submission while the sweep is unconfirmed.
	s.Pass()
	require.Equal(t, pool, n.poolHashes())

	b := n.block()
	var included []cipher.SHA256
	for _, txn := range b.Body.Transactions {
		included = append(included, txn.Hash())
	}
	require.NotContains(t, included, attack.Hash(), "the attacker paid less per byte and lost")
	for _, h := range st.Pending {
		require.Contains(t, included, cipher.MustSHA256FromHex(h), "the sweep landed in the next block")
	}

	require.Empty(t, n.unspentOf(addr(shared[1])))
	require.Empty(t, n.unspentOf(addr(shared[3])))
	require.Empty(t, n.unspentOf(thief))
	got := n.unspentOf(dest)
	coins, err := got.Coins()
	require.NoError(t, err)
	require.Equal(t, uint64(15e6), coins)
	for _, u := range got {
		require.Zero(t, u.Body.Hours, "every coin hour was burned")
	}
	require.Len(t, n.unspentOf(addr(savings[1])), 1, "a never-sent single address is not swept")

	s.Pass()
	st = s.Status(id)
	require.Empty(t, st.Pending)
	require.Equal(t, 1, st.Confirmed)
	require.Zero(t, st.WaitingCoins)

	// Equal fee: the attacker builds the same shape as the sweep, burning every
	// hour. The two tie on fee per kB and the lower hash wins.
	n.pay([]coin.TransactionOutput{{Address: addr(shared[2]), Coins: 3e6, Hours: 500}})
	payout := n.unspentOf(addr(shared[2]))[0]
	s.Pass()
	st = s.Status(id)
	require.Len(t, st.Pending, 1)
	sweepHash := cipher.MustSHA256FromHex(st.Pending[0])
	tie := steal(t, payout, shared[2].Secret, thief, 0)
	require.NoError(t, n.InjectBroadcastTransaction(tie))

	b = n.block()
	var hashes []cipher.SHA256
	for _, txn := range b.Body.Transactions {
		if txn.In[0] == payout.Hash() {
			hashes = append(hashes, txn.Hash())
		}
	}
	require.Len(t, hashes, 1, "the payout is spent once")
	winner := sweepHash
	if th := tie.Hash(); bytes.Compare(th[:], sweepHash[:]) < 0 {
		winner = tie.Hash()
	}
	require.Equal(t, winner, hashes[0])
	t.Logf("tie: sweep %s, attacker %s, winner %s", sweepHash.Hex(), tie.Hash().Hex(), winner.Hex())

	s.Pass()
	st = s.Status(id)
	require.Empty(t, st.Pending)
	if winner == sweepHash {
		require.Equal(t, 2, st.Confirmed)
	} else {
		require.Equal(t, []string{sweepHash.Hex()}, st.Displaced)
	}
}

func TestSweepSpentSingleAddress(t *testing.T) {
	n := newNode(t)
	const id = "d.wlt"
	_, err := n.ws.CreateWallet(id, wallet.Options{Label: "d", Seed: "single seed", Type: wallet.WalletTypeDeterministic, GenerateN: 2})
	require.NoError(t, err)
	w, err := n.ws.GetWallet(id)
	require.NoError(t, err)
	es, err := w.GetEntries()
	require.NoError(t, err)
	a, b := es[0].Address.(cipher.Address), es[1].Address.(cipher.Address)

	n.pay([]coin.TransactionOutput{{Address: a, Coins: 4e6, Hours: 200}, {Address: b, Coins: 6e6, Hours: 200}})
	s, err := sweep.NewService(sweep.Config{Dir: t.TempDir()}, n.ws, n.v, n)
	require.NoError(t, err)
	_, err = s.SetRule(sweep.Rule{Wallet: id, Enabled: true})
	require.NoError(t, err)
	s.Pass()
	require.Empty(t, s.Status(id).Pending, "neither address has sent, so neither is swept")

	// a sends once, revealing its public key, and keeps change on itself.
	out := n.unspentOf(a)[0]
	var txn coin.Transaction
	require.NoError(t, txn.PushInput(out.Hash()))
	require.NoError(t, txn.PushOutput(testutil.MakeAddress(), 1e6, 10))
	require.NoError(t, txn.PushOutput(a, 3e6, 10))
	require.NoError(t, txn.SignInputs([]cipher.SecKey{es[0].Secret}))
	require.NoError(t, txn.UpdateHeader())
	require.NoError(t, n.InjectBroadcastTransaction(txn))
	n.block()

	s.Pass()
	st := s.Status(id)
	require.Empty(t, st.LastError)
	require.Equal(t, 1, st.Exposed)
	require.Len(t, st.Pending, 1)
	n.block()
	require.Empty(t, n.unspentOf(a))
	require.Len(t, n.unspentOf(b), 1)
	dest := cipher.MustDecodeBase58Address(st.Destination)
	require.NotEqual(t, a, dest)
	require.NotEqual(t, b, dest, "b appears on chain, so the next unused address is a new one")
	coins, err := n.unspentOf(dest).Coins()
	require.NoError(t, err)
	require.Equal(t, uint64(3e6), coins)
}
