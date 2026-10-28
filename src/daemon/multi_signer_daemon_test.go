package daemon

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const msGenCoins = 1000e6

// msChain is a k-of-n multi-signer chain for daemon tests.
type msChain struct {
	pks     []cipher.PubKey
	sks     []cipher.SecKey
	rule    coin.MultiSignerRule
	genPK   cipher.PubKey
	genSK   cipher.SecKey
	genAddr cipher.Address
	genesis coin.SignedBlock
}

// newMSChain makes a chain whose genesis block has time genTime. With
// genesisSigner, signer 0 holds the genesis key.
func newMSChain(t *testing.T, n, k int, forkSeq uint64, genTime time.Time, genesisSigner bool) *msChain {
	c := &msChain{pks: make([]cipher.PubKey, n), sks: make([]cipher.SecKey, n)}
	c.genPK, c.genSK = cipher.GenerateKeyPair()
	c.genAddr = cipher.AddressFromPubKey(c.genPK)
	for i := range c.pks {
		c.pks[i], c.sks[i] = cipher.GenerateKeyPair()
	}
	if genesisSigner {
		c.pks[0], c.sks[0] = c.genPK, c.genSK
	}
	if n > 0 {
		c.rule = coin.MultiSignerRule{Signers: c.pks, Threshold: k, ForkSeq: forkSeq}
	}
	gb, err := coin.NewGenesisBlock(c.genAddr, msGenCoins, uint64(genTime.Unix())) //nolint:gosec
	require.NoError(t, err)
	c.genesis = coin.SignedBlock{Block: *gb, Sig: cipher.MustSignHash(gb.HashHeader(), c.genSK)}
	return c
}

// visorConfig is the config of signer i, or of a plain node for i < 0.
func (c *msChain) visorConfig(i int) visor.Config {
	cfg := visor.NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = c.genPK
	cfg.GenesisAddress = c.genAddr
	cfg.GenesisCoinVolume = msGenCoins
	cfg.GenesisTimestamp = c.genesis.Time()
	cfg.GenesisSignature = c.genesis.Sig
	cfg.MultiSigner = c.rule
	if i >= 0 {
		cfg.IsBlockPublisher = true
		cfg.Arbitrating = true
		cfg.BlockchainSeckey = c.sks[i]
	}
	return cfg
}

func (c *msChain) openVisor(t *testing.T, i int, db *dbutil.DB) *visor.Visor {
	v, err := visor.New(c.visorConfig(i), db, nil)
	require.NoError(t, err)
	_, ok, err := v.HeadBkSeq()
	require.NoError(t, err)
	if !ok {
		require.NoError(t, v.ExecuteSignedBlock(c.genesis))
	}
	return v
}

func (c *msChain) newVisor(t *testing.T, i int) *visor.Visor {
	db, shutdown := testutil.PrepareDB(t)
	t.Cleanup(shutdown)
	return c.openVisor(t, i, db)
}

// nextTxn spends the first output of the head block's first transaction.
func (c *msChain) nextTxn(t *testing.T, v *visor.Visor) coin.Transaction {
	head, err := v.GetHeadBlock()
	require.NoError(t, err)
	ux := coin.CreateUnspents(head.Head, head.Body.Transactions[0])[0]
	var txn coin.Transaction
	require.NoError(t, txn.PushInput(ux.Hash()))
	require.NoError(t, txn.PushOutput(c.genAddr, ux.Body.Coins, ux.Body.Hours/2))
	txn.SignInputs([]cipher.SecKey{c.genSK})
	require.NoError(t, txn.UpdateHeader())
	return txn
}

func inject(t *testing.T, v *visor.Visor, txn coin.Transaction) {
	_, softErr, err := v.InjectForeignTransaction(txn)
	require.NoError(t, err)
	require.Nil(t, softErr)
}

// propose makes signer visor v's proposal for the block after its head.
func (c *msChain) propose(t *testing.T, v *visor.Visor, when time.Time) coin.SignedBlock {
	inject(t, v, c.nextTxn(t, v))
	b, err := v.CreateBlockProposal(when)
	require.NoError(t, err)
	return b
}

func (c *msChain) cosig(t *testing.T, i int, hash cipher.SHA256) BlockCosigEntry {
	cs, err := coin.SignBlockCosig(hash, i, c.sks[i])
	require.NoError(t, err)
	return BlockCosigEntry{Hash: hash, Cosig: cs}
}

func resign(b coin.SignedBlock, sk cipher.SecKey) coin.SignedBlock {
	b.Sig = cipher.MustSignHash(b.HashHeader(), sk)
	return b
}

type sentMsg struct {
	to  string
	msg gnet.Message
}

// sentLog records the messages a daemon sends.
type sentLog struct {
	mu   sync.Mutex
	msgs []sentMsg
}

func (s *sentLog) SendMessage(addr string, m gnet.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msgs = append(s.msgs, sentMsg{addr, m})
	return nil
}

func (s *sentLog) BroadcastMessage(m gnet.Message, addrs []string) ([]uint64, error) {
	if len(addrs) == 0 {
		return nil, gnet.ErrNoAddresses
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range addrs {
		s.msgs = append(s.msgs, sentMsg{a, m})
	}
	return nil, nil
}

func (s *sentLog) reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.msgs = nil
}

// sentOf returns the messages of type T by recipient.
func sentOf[T gnet.Message](s *sentLog) map[string][]T {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string][]T)
	for _, m := range s.msgs {
		if v, ok := m.msg.(T); ok {
			out[m.to] = append(out[m.to], v)
		}
	}
	return out
}

func newMSDaemon(v *visor.Visor) (*Daemon, *sentLog) {
	out := &sentLog{}
	cfg := NewDaemonConfig()
	cfg.BlockRoundTimeout = time.Hour
	return &Daemon{
		config:        cfg,
		visor:         v,
		out:           out,
		announcedTxns: newAnnouncedTxnsCache(),
		pqSigs:        newPQSigCache(),
		pqWitnesses:   newPQSigCache(),
		multiSigner:   newMultiSignerState(),
		connections:   NewConnections(),
		events:        make(chan interface{}, 64),
	}, out
}

func addMSPeer(t *testing.T, dm *Daemon, addr string, id uint64, version int32) {
	_, err := dm.connections.pending(addr)
	require.NoError(t, err)
	_, err = dm.connections.connected(addr, id)
	require.NoError(t, err)
	_, err = dm.connections.introduced(addr, id, &IntroductionMessage{ProtocolVersion: version, Mirror: uint32(id), ListenPort: 6000}) //nolint:gosec
	require.NoError(t, err)
}

func headSeq(t *testing.T, v *visor.Visor) uint64 {
	seq, ok, err := v.HeadBkSeq()
	require.NoError(t, err)
	require.True(t, ok)
	return seq
}

// cosigSigners returns the signers of the cosignatures sent for hash.
func cosigSigners(s *sentLog, hash cipher.SHA256) map[uint32]bool {
	got := make(map[uint32]bool)
	for _, ms := range sentOf[*GiveBlockCosigsMessage](s) {
		for _, m := range ms {
			for _, e := range m.Cosigs {
				if e.Hash == hash {
					got[e.Cosig.Signer] = true
				}
			}
		}
	}
	return got
}

const (
	peerA = "127.0.0.2:6000"
	peerB = "127.0.0.3:6000"
	peerC = "127.0.0.4:6000"
)

func TestMultiSignerCosignsOncePerHeight(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	v := c.openVisor(t, 3, db)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)

	a := c.propose(t, c.newVisor(t, 1), time.Now())
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
	require.Equal(t, map[uint32]bool{3: true}, cosigSigners(out, a.HashHeader()))
	vote, ok, err := v.BlockVote(1)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, a.HashHeader(), vote)

	// Signer 2 is the proposer of round 1 at the same height
	b := c.propose(t, c.newVisor(t, 2), time.Now().Add(time.Second))
	require.NotEqual(t, a.HashHeader(), b.HashHeader())
	require.Equal(t, 2, c.rule.Proposer(1, 1))
	out.reset()
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 1, Block: b})
	require.Empty(t, cosigSigners(out, b.HashHeader()))

	// The vote is read back from the db after a restart
	v2 := c.openVisor(t, 3, db)
	dm2, out2 := newMSDaemon(v2)
	addMSPeer(t, dm2, peerA, 1, multiSignerProtocolVersion)
	dm2.onBlockProposal(peerA, &ProposeBlockMessage{Round: 1, Block: b})
	require.Empty(t, cosigSigners(out2, b.HashHeader()))
	_, held := dm2.multiSigner.proposal(b.HashHeader())
	require.True(t, held, "the proposal is still held, only not signed")

	vote, _, err = v2.BlockVote(1)
	require.NoError(t, err)
	require.Equal(t, a.HashHeader(), vote)
	require.Equal(t, uint64(0), headSeq(t, v2))
}

func TestMultiSignerRejectsBadProposals(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	p1 := c.newVisor(t, 1)
	a := c.propose(t, p1, time.Now())
	_, outsider := cipher.GenerateKeyPair()

	cases := []struct {
		name  string
		round uint64
		block func() coin.SignedBlock
		// held is whether the block is kept for signatures, which a proposer
		// may cause; executable whether k signatures would then make it final
		held, executable bool
	}{
		{"wrong proposer for the round", 0, func() coin.SignedBlock { return resign(a, c.sks[2]) }, false, false},
		{"proposer not listed", 0, func() coin.SignedBlock { return resign(a, outsider) }, false, false},
		{"round too far ahead", 3, func() coin.SignedBlock { return resign(a, c.sks[0]) }, false, false},
		{"not on our head", 0, func() coin.SignedBlock {
			b := a
			b.Head.PrevHash = cipher.SumSHA256([]byte("other"))
			return resign(b, c.sks[1])
		}, true, false},
		{"time not after head", 0, func() coin.SignedBlock {
			b := a
			b.Head.Time = c.genesis.Time()
			return resign(b, c.sks[1])
		}, true, false},
		{"time ahead of our clock", 0, func() coin.SignedBlock {
			b := a
			b.Head.Time = uint64(time.Now().Add(time.Hour).Unix()) //nolint:gosec
			return resign(b, c.sks[1])
		}, true, true},
		{"invalid transaction", 0, func() coin.SignedBlock {
			b := a
			txns := append(coin.Transactions{}, b.Body.Transactions...)
			txn := txns[0]
			txn.Out = append([]coin.TransactionOutput{}, txn.Out...)
			txn.Out[0].Hours++
			txns[0] = txn
			b.Body.Transactions = txns
			b.Head.BodyHash = b.Body.Hash()
			return resign(b, c.sks[1])
		}, true, false},
		{"old height", 0, func() coin.SignedBlock { return c.genesis }, false, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := c.newVisor(t, 3)
			dm, out := newMSDaemon(v)
			addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
			addMSPeer(t, dm, peerB, 2, multiSignerProtocolVersion)

			b := tc.block()
			dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: tc.round, Block: b})
			require.Empty(t, sentOf[*GiveBlockCosigsMessage](out), "not cosigned")
			_, voted, err := v.BlockVote(1)
			require.NoError(t, err)
			require.False(t, voted)
			_, held := dm.multiSigner.proposal(b.HashHeader())
			require.Equal(t, tc.held, held)
			if !tc.held {
				require.Empty(t, sentOf[*ProposeBlockMessage](out), "not relayed")
			}

			// Signers 0 and 2 sign it anyway
			h := b.HashHeader()
			dm.onBlockCosigs(peerA, []BlockCosigEntry{c.cosig(t, 0, h), c.cosig(t, 2, h)})
			want := uint64(0)
			if tc.executable {
				want = 1
			}
			require.Equal(t, want, headSeq(t, v))
		})
	}
}

func TestMultiSignerProposalForLaterHeight(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	b := c.propose(t, c.newVisor(t, 1), time.Now())
	b.Head.BkSeq = 2
	b = resign(b, c.sks[2])

	dm, out := newMSDaemon(c.newVisor(t, 3))
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: b})
	require.Empty(t, sentOf[*GiveBlockCosigsMessage](out))
	require.Empty(t, sentOf[*ProposeBlockMessage](out))
	require.Len(t, sentOf[*GetBlocksMessage](out)[peerA], 1, "asks the sender for the blocks it lacks")
}

func TestMultiSignerThreshold(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	a := c.propose(t, c.newVisor(t, 1), time.Now())
	h := a.HashHeader()

	v := c.newVisor(t, -1)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	addMSPeer(t, dm, peerB, 2, multiSignerProtocolVersion-1)
	addMSPeer(t, dm, peerC, 3, multiSignerProtocolVersion)

	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
	require.Equal(t, uint64(0), headSeq(t, v), "the proposer's signature alone")

	_, outsider := cipher.GenerateKeyPair()
	outsiderCosig, err := coin.SignBlockCosig(h, 0, outsider)
	require.NoError(t, err)
	other := cipher.SumSHA256([]byte("other"))
	two := c.cosig(t, 2, h)
	junk := []BlockCosigEntry{
		two,
		two,
		c.cosig(t, 1, h), // the proposer again
		{Hash: h, Cosig: coin.BlockCosig{Signer: 7, Sig: two.Cosig.Sig}},
		{Hash: h, Cosig: outsiderCosig},
		{Hash: h, Cosig: coin.BlockCosig{Signer: 3, Sig: c.cosig(t, 3, other).Cosig.Sig}},
		{Hash: h, Cosig: coin.BlockCosig{Signer: 0, Sig: cipher.Sig{1, 2, 3}}},
	}
	out.reset()
	dm.onBlockCosigs(peerA, junk)
	require.Equal(t, uint64(0), headSeq(t, v), "k-1 valid signatures")

	relayed := sentOf[*GiveBlockCosigsMessage](out)
	require.NotContains(t, relayed, peerA, "not back to the sender")
	require.NotContains(t, relayed, peerB, "not to an old peer")
	require.Len(t, relayed[peerC], 1)
	require.ElementsMatch(t, []BlockCosigEntry{two, junk[2]}, relayed[peerC][0].Cosigs, "only the valid new ones")

	// The same again is not new and not relayed
	out.reset()
	dm.onBlockCosigs(peerC, junk)
	require.Empty(t, out.msgs)

	dm.onBlockCosigs(peerA, []BlockCosigEntry{c.cosig(t, 3, h)})
	require.Equal(t, uint64(1), headSeq(t, v), "k valid signatures")
	sb, err := v.GetSignedBlockBySeq(1)
	require.NoError(t, err)
	require.Equal(t, h, sb.HashHeader())
	require.Equal(t, []cipher.PubKey{c.pks[1], c.pks[2], c.pks[3]}, sortedSigners(c.rule, *sb))

	gave := sentOf[*GiveBlocksMessage](out)
	require.Len(t, gave[peerA], 1, "the final block is broadcast")
	require.Len(t, gave[peerB], 1)
}

// sortedSigners returns the block signers in signer index order.
func sortedSigners(rule coin.MultiSignerRule, b coin.SignedBlock) []cipher.PubKey {
	in := make(map[cipher.PubKey]bool)
	for _, pk := range rule.BlockSigners(b) {
		in[pk] = true
	}
	var pks []cipher.PubKey
	for _, pk := range rule.Signers {
		if in[pk] {
			pks = append(pks, pk)
		}
	}
	return pks
}

func TestMultiSignerCosigsBeforeProposal(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	a := c.propose(t, c.newVisor(t, 1), time.Now())
	h := a.HashHeader()

	v := c.newVisor(t, -1)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	addMSPeer(t, dm, peerB, 2, multiSignerProtocolVersion)

	dm.onBlockCosigs(peerA, []BlockCosigEntry{c.cosig(t, 2, h), c.cosig(t, 3, h)})
	require.Equal(t, uint64(0), headSeq(t, v))
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
	require.Equal(t, uint64(1), headSeq(t, v))

	// Proposals and signatures for the old height change nothing
	out.reset()
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
	require.Empty(t, out.msgs)
	dm.onBlockCosigs(peerA, []BlockCosigEntry{c.cosig(t, 0, h)})
	require.Equal(t, uint64(1), headSeq(t, v))
	_, held := dm.multiSigner.proposal(h)
	require.False(t, held)
}

func TestMultiSignerHoldsBlockUntilCosigs(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	a := c.propose(t, c.newVisor(t, 1), time.Now())
	h := a.HashHeader()

	v := c.newVisor(t, -1)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)

	m := &GiveBlocksMessage{Blocks: []coin.SignedBlock{a}, c: &gnet.MessageContext{Addr: peerA, ConnID: 1}}
	m.process(dm)
	require.Equal(t, uint64(0), headSeq(t, v))
	_, held := dm.multiSigner.proposal(h)
	require.True(t, held)
	require.Len(t, sentOf[*GetBlocksMessage](out)[peerA], 1, "asks again for the signatures")

	// A second under-signed copy does not ask again so soon
	out.reset()
	m.process(dm)
	require.Empty(t, sentOf[*GetBlocksMessage](out))

	dm.onBlockCosigs(peerA, []BlockCosigEntry{c.cosig(t, 2, h), c.cosig(t, 3, h)})
	require.Equal(t, uint64(1), headSeq(t, v))
}

// newSigners makes a daemon for each signer, each with every other as a peer.
func newSigners(t *testing.T, c *msChain, cfg func(*DaemonConfig)) ([]*Daemon, []*sentLog) {
	dms := make([]*Daemon, len(c.sks))
	outs := make([]*sentLog, len(c.sks))
	for i := range dms {
		v := c.newVisor(t, i)
		inject(t, v, c.nextTxn(t, v))
		dms[i], outs[i] = newMSDaemon(v)
		addMSPeer(t, dms[i], peerA, 1, multiSignerProtocolVersion)
		if cfg != nil {
			cfg(&dms[i].config)
		}
	}
	return dms, outs
}

func proposers(outs []*sentLog) map[int]uint64 {
	got := make(map[int]uint64)
	for i, out := range outs {
		for _, m := range sentOf[*ProposeBlockMessage](out)[peerA] {
			got[i] = m.Round
		}
	}
	return got
}

func TestMultiSignerSignerTick(t *testing.T) {
	t.Run("round 0", func(t *testing.T) {
		c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
		dms, outs := newSigners(t, c, nil)
		for _, dm := range dms {
			dm.signerTick()
		}
		require.Equal(t, map[int]uint64{1: 0}, proposers(outs))
		for _, dm := range dms {
			require.False(t, dm.publishesNextBlock(), "signers publish only together")
		}
		vote, ok, err := dms[1].visor.BlockVote(1)
		require.NoError(t, err)
		require.True(t, ok)
		sent := sentOf[*ProposeBlockMessage](outs[1])[peerA][0]
		require.Equal(t, sent.Block.HashHeader(), vote, "the vote is recorded")
		require.Equal(t, 1, c.rule.ProposerIndex(sent.Block))
	})

	t.Run("round 1 after the timeout", func(t *testing.T) {
		c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
		dms, outs := newSigners(t, c, func(dc *DaemonConfig) { dc.BlockRoundTimeout = 40 * time.Second })
		for _, dm := range dms {
			dm.signerTick()
		}
		require.Equal(t, map[int]uint64{2: 1}, proposers(outs))
	})

	t.Run("before the block interval", func(t *testing.T) {
		c := newMSChain(t, 4, 3, 1, time.Now().Add(-5*time.Second), false)
		dms, outs := newSigners(t, c, nil)
		for _, dm := range dms {
			dm.signerTick()
		}
		require.Empty(t, proposers(outs))
	})

	t.Run("before the fork", func(t *testing.T) {
		c := newMSChain(t, 4, 3, 3, time.Now().Add(-time.Minute), false)
		dms, outs := newSigners(t, c, nil)
		for _, dm := range dms {
			dm.signerTick()
		}
		require.Empty(t, proposers(outs))
	})

	t.Run("not a signer", func(t *testing.T) {
		c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
		v := c.newVisor(t, -1)
		inject(t, v, c.nextTxn(t, v))
		dm, out := newMSDaemon(v)
		addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
		dm.signerTick()
		require.Empty(t, out.msgs)
	})
}

func TestMultiSignerBlockRound(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now(), false)
	dm, _ := newMSDaemon(c.newVisor(t, -1))
	dm.config.BlockRoundTimeout = 0
	dm.config.BlockCreationInterval = 10
	require.Equal(t, 30*time.Second, dm.blockRoundTimeout())
	dm.config.BlockRoundTimeout = time.Second
	require.Equal(t, time.Second, dm.blockRoundTimeout())

	dm.config.BlockRoundTimeout = 0
	head := &coin.SignedBlock{Block: coin.Block{Head: coin.BlockHeader{Time: 1000}}}
	require.Equal(t, uint64(0), dm.blockRound(head, time.Unix(900, 0)))
	require.Equal(t, uint64(0), dm.blockRound(head, time.Unix(1029, 0)))
	require.Equal(t, uint64(1), dm.blockRound(head, time.Unix(1030, 0)))
	require.Equal(t, uint64(2), dm.blockRound(head, time.Unix(1065, 0)))
}

func TestMultiSignerResendVote(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	a := c.propose(t, c.newVisor(t, 1), time.Now().Add(-time.Second))
	h := a.HashHeader()

	t.Run("cosigner", func(t *testing.T) {
		v := c.newVisor(t, 3)
		dm, out := newMSDaemon(v)
		addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
		dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
		first := sentOf[*GiveBlockCosigsMessage](out)[peerA][0].Cosigs[0]

		out.reset()
		dm.signerTick()
		require.Empty(t, out.msgs, "nothing to resend in the same round")

		dm.multiSigner.lastRound = 99
		dm.signerTick()
		props := sentOf[*ProposeBlockMessage](out)[peerA]
		require.Len(t, props, 1)
		require.Equal(t, h, props[0].Block.HashHeader())
		cs := sentOf[*GiveBlockCosigsMessage](out)[peerA]
		require.Len(t, cs, 1)
		require.Equal(t, []BlockCosigEntry{first}, cs[0].Cosigs)
	})

	// Signer 1 is the proposer of the round and has a block to propose, but it
	// has recorded a vote for a, so it resends a and signs nothing new.
	t.Run("recorded vote", func(t *testing.T) {
		v := c.newVisor(t, 1)
		inject(t, v, c.nextTxn(t, v))
		dm, out := newMSDaemon(v)
		addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
		dm.multiSigner.setHeight(1)
		require.True(t, dm.multiSigner.addProposal(a, 0))
		require.NoError(t, v.RecordBlockVote(1, h))

		dm.multiSigner.lastRound = 99
		dm.signerTick()
		props := sentOf[*ProposeBlockMessage](out)[peerA]
		require.Len(t, props, 1)
		require.Equal(t, h, props[0].Block.HashHeader())
		require.Empty(t, sentOf[*GiveBlockCosigsMessage](out))
		vote, _, err := v.BlockVote(1)
		require.NoError(t, err)
		require.Equal(t, h, vote)
	})
}

func TestMultiSignerMessagesHandle(t *testing.T) {
	mc := &gnet.MessageContext{Addr: peerA}
	msgs := []struct {
		m       asyncMessage
		handler string
	}{
		{&ProposeBlockMessage{Round: 1}, "onBlockProposal"},
		{&GiveBlockCosigsMessage{Cosigs: []BlockCosigEntry{{}}}, "onBlockCosigs"},
	}
	for _, tc := range msgs {
		t.Run(tc.handler, func(t *testing.T) {
			d := &mockDaemoner{}
			d.On("recordMessageEvent", tc.m, mc).Return(nil)
			require.NoError(t, tc.m.(gnet.Message).Handle(mc, d))
			d.AssertCalled(t, "recordMessageEvent", tc.m, mc)

			off := &mockDaemoner{}
			off.On("DaemonConfig").Return(DaemonConfig{DisableNetworking: true})
			tc.m.process(off)
			off.AssertNotCalled(t, tc.handler, mock.Anything, mock.Anything)

			on := &mockDaemoner{}
			on.On("DaemonConfig").Return(DaemonConfig{})
			on.On(tc.handler, peerA, mock.Anything).Return()
			tc.m.process(on)
			on.AssertCalled(t, tc.handler, peerA, mock.Anything)
		})
	}
}

func TestMultiSignerMessagesBadEncoding(t *testing.T) {
	var p ProposeBlockMessage
	for _, buf := range [][]byte{nil, {1, 2, 3}, make([]byte, 8), make([]byte, 20)} {
		_, err := p.Decode(buf)
		require.Error(t, err, "%x", buf)
	}

	var g GiveBlockCosigsMessage
	two := make([]byte, 4+blockCosigEntrySize)
	two[0] = 2
	for _, buf := range [][]byte{nil, {1, 0}, two} {
		_, err := g.Decode(buf)
		require.Error(t, err, "%x", buf)
	}
}

func TestMultiSignerPeers(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	dm, out := newMSDaemon(c.newVisor(t, -1))
	require.Empty(t, dm.multiSignerPeers(""))

	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	addMSPeer(t, dm, peerB, 2, multiSignerProtocolVersion-1)
	_, err := dm.connections.pending(peerC)
	require.NoError(t, err)
	_, err = dm.connections.connected(peerC, 3)
	require.NoError(t, err)
	require.Equal(t, []string{peerA}, dm.multiSignerPeers(""), "introduced and new enough")
	require.Empty(t, dm.multiSignerPeers(peerA))

	a := c.propose(t, c.newVisor(t, 1), time.Now())
	dm.relayProposal(&ProposeBlockMessage{Block: a}, "")
	dm.relayBlockCosigs([]BlockCosigEntry{c.cosig(t, 2, a.HashHeader())}, "")
	for _, m := range out.msgs {
		require.Equal(t, peerA, m.to, "%T", m.msg)
	}
	require.Len(t, out.msgs, 2)

	// sendBlockCosigs sends only to a peer that understands them
	out.reset()
	a.Cosigs = []coin.BlockCosig{c.cosig(t, 2, a.HashHeader()).Cosig}
	require.NoError(t, dm.sendBlockCosigs(peerB, []coin.SignedBlock{a}))
	require.Empty(t, out.msgs)
	require.Error(t, dm.sendBlockCosigs("127.0.0.9:6000", []coin.SignedBlock{a}))
	require.NoError(t, dm.sendBlockCosigs(peerA, []coin.SignedBlock{c.genesis}), "no cosignatures")
	require.Empty(t, out.msgs)
	require.NoError(t, dm.sendBlockCosigs(peerA, []coin.SignedBlock{a}))
	require.Len(t, sentOf[*GiveBlockCosigsMessage](out)[peerA], 1)
}

func TestMultiSignerProposalRelayedOnce(t *testing.T) {
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-time.Minute), false)
	a := c.propose(t, c.newVisor(t, 1), time.Now())
	dm, out := newMSDaemon(c.newVisor(t, -1))
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	addMSPeer(t, dm, peerB, 2, multiSignerProtocolVersion)
	addMSPeer(t, dm, peerC, 3, multiSignerProtocolVersion-1)

	dm.onBlockProposal(peerA, &ProposeBlockMessage{Round: 0, Block: a})
	dm.onBlockProposal(peerB, &ProposeBlockMessage{Round: 0, Block: a})
	sent := sentOf[*ProposeBlockMessage](out)
	require.Len(t, sent, 1)
	require.Len(t, sent[peerB], 1)
}

func TestMultiSignerOff(t *testing.T) {
	c := newMSChain(t, 0, 0, 0, time.Now().Add(-time.Minute), false)
	pub := c.visorConfig(-1)
	pub.IsBlockPublisher = true
	pub.BlockchainSeckey = c.genSK
	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()
	pv, err := visor.New(pub, db, nil)
	require.NoError(t, err)
	require.NoError(t, pv.ExecuteSignedBlock(c.genesis))
	require.Equal(t, -1, pv.SignerIndex())
	inject(t, pv, c.nextTxn(t, pv))
	b, err := pv.CreateAndExecuteBlock()
	require.NoError(t, err)

	v := c.newVisor(t, -1)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	require.True(t, dm.publishesNextBlock())

	dm.onBlockProposal(peerA, &ProposeBlockMessage{Block: b})
	dm.onBlockCosigs(peerA, []BlockCosigEntry{{Hash: b.HashHeader()}})
	pdm, pout := newMSDaemon(pv)
	addMSPeer(t, pdm, peerA, 1, multiSignerProtocolVersion)
	pdm.signerTick()
	require.Empty(t, out.msgs)
	require.Empty(t, pout.msgs)
	require.Empty(t, dm.multiSigner.cosigs)

	m := &GiveBlocksMessage{Blocks: []coin.SignedBlock{b}, c: &gnet.MessageContext{Addr: peerA, ConnID: 1}}
	m.process(dm)
	require.Equal(t, uint64(1), headSeq(t, v), "an ordinary block is accepted")
}

func TestMultiSignerBeforeFork(t *testing.T) {
	c := newMSChain(t, 4, 3, 3, time.Now().Add(-time.Minute), true)

	// Signer 0 holds the genesis key and publishes alone before the fork
	pv := c.newVisor(t, 0)
	pdm, pout := newMSDaemon(pv)
	addMSPeer(t, pdm, peerA, 1, multiSignerProtocolVersion)
	require.True(t, pdm.publishesNextBlock())
	inject(t, pv, c.nextTxn(t, pv))
	pdm.signerTick()
	require.Empty(t, pout.msgs, "no proposals before the fork")
	genTime := time.Unix(int64(c.genesis.Time()), 0) //nolint:gosec
	b, err := pv.CreateBlockProposal(genTime.Add(10 * time.Second))
	require.NoError(t, err)
	require.NoError(t, pv.ExecuteSignedBlock(b))

	// Another signer never publishes alone
	odm, _ := newMSDaemon(c.newVisor(t, 1))
	require.False(t, odm.publishesNextBlock())

	v := c.newVisor(t, 2)
	dm, out := newMSDaemon(v)
	addMSPeer(t, dm, peerA, 1, multiSignerProtocolVersion)
	dm.onBlockProposal(peerA, &ProposeBlockMessage{Block: resign(b, c.sks[1])})
	require.Empty(t, out.msgs, "no cosigning before the fork")
	_, voted, err := v.BlockVote(1)
	require.NoError(t, err)
	require.False(t, voted)

	m := &GiveBlocksMessage{Blocks: []coin.SignedBlock{b}, c: &gnet.MessageContext{Addr: peerA, ConnID: 1}}
	m.process(dm)
	require.Equal(t, uint64(1), headSeq(t, v), "a block signed by the genesis key")

	// At the fork the genesis key alone is no longer enough
	b2 := c.propose(t, pv, genTime.Add(20*time.Second))
	require.NoError(t, pv.ExecuteSignedBlock(b2))
	require.False(t, pdm.publishesNextBlock(), "block 3 needs k signatures")
	b3 := c.propose(t, pv, genTime.Add(30*time.Second))
	err = pv.ExecuteSignedBlock(b3)
	var few coin.ErrNotEnoughBlockSigners
	require.ErrorAs(t, err, &few, "the publisher cannot add block 3 alone")
	m = &GiveBlocksMessage{Blocks: []coin.SignedBlock{b2}, c: &gnet.MessageContext{Addr: peerA, ConnID: 1}}
	m.process(dm)
	require.Equal(t, uint64(2), headSeq(t, v))
}

// msRouter carries messages between in-process daemons. Each message is
// encoded and decoded as on the wire.
type msRouter struct {
	t     *testing.T
	nodes map[string]*msNode
	queue []routed
	sent  int
}

type msNode struct {
	dm   *Daemon
	addr string
	id   uint64
	down bool
}

type routed struct {
	from, to string
	msg      gnet.Message
}

type routedSender struct {
	r    *msRouter
	from string
}

func (s routedSender) SendMessage(addr string, m gnet.Message) error {
	s.r.queue = append(s.r.queue, routed{s.from, addr, m})
	return nil
}

func (s routedSender) BroadcastMessage(m gnet.Message, addrs []string) ([]uint64, error) {
	if len(addrs) == 0 {
		return nil, gnet.ErrNoAddresses
	}
	for _, a := range addrs {
		s.r.queue = append(s.r.queue, routed{s.from, a, m})
	}
	return nil, nil
}

func wireCopy(t *testing.T, m gnet.Message) gnet.Message {
	s, ok := m.(gnet.Serializer)
	require.True(t, ok, "%T", m)
	buf := make([]byte, s.EncodeSize())
	require.NoError(t, s.Encode(buf))
	v := reflect.New(reflect.TypeOf(m).Elem()).Interface()
	n, err := v.(gnet.Serializer).Decode(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(len(buf)), n)
	return v.(gnet.Message)
}

// pump delivers queued messages until none are left.
func (r *msRouter) pump() {
	for len(r.queue) > 0 {
		e := r.queue[0]
		r.queue = r.queue[1:]
		r.sent++
		require.Less(r.t, r.sent, 100000, "messages loop")
		to, from := r.nodes[e.to], r.nodes[e.from]
		if to.down || from.down {
			continue
		}
		m := wireCopy(r.t, e.msg)
		require.NoError(r.t, m.Handle(&gnet.MessageContext{Addr: from.addr, ConnID: from.id}, to.dm))
		to.dm.onMessageEvent((<-to.dm.events).(messageEvent))
	}
}

// A 3-of-4 chain makes blocks with one signer offline, and a plain node follows.
func TestMultiSignerNetwork(t *testing.T) {
	if testing.Short() {
		t.Skip("makes blocks in real time")
	}
	c := newMSChain(t, 4, 3, 1, time.Now().Add(-2*time.Second), false)
	r := &msRouter{t: t, nodes: make(map[string]*msNode)}
	var nodes []*msNode
	for i := -1; i < 4; i++ {
		v := c.newVisor(t, i)
		dm, _ := newMSDaemon(v)
		dm.config.BlockCreationInterval = 1
		dm.config.BlockRoundTimeout = 3 * time.Second
		n := &msNode{dm: dm, addr: fmt.Sprintf("127.0.1.%d:6000", i+2), id: uint64(i + 2)} //nolint:gosec
		dm.out = routedSender{r, n.addr}
		r.nodes[n.addr] = n
		nodes = append(nodes, n)
	}
	for _, n := range nodes {
		for _, p := range nodes {
			if p != n {
				addMSPeer(t, n.dm, p.addr, p.id, multiSignerProtocolVersion)
			}
		}
	}
	follower := nodes[0]
	nodes[3].down = true // signer 2

	const target = 4
	injected := uint64(0)
	deadline := time.Now().Add(40 * time.Second)
	for headSeq(t, follower.dm.visor) < target {
		require.True(t, time.Now().Before(deadline), "stalled at %d", headSeq(t, follower.dm.visor))
		h := headSeq(t, follower.dm.visor)
		if injected <= h {
			for _, n := range nodes {
				if !n.down {
					require.Equal(t, h, headSeq(t, n.dm.visor))
					inject(t, n.dm.visor, c.nextTxn(t, n.dm.visor))
				}
			}
			injected = h + 1
		}
		// Delivery is fast next to a round, so each tick's messages arrive before the next
		for _, n := range nodes {
			if !n.down {
				n.dm.signerTick()
				r.pump()
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	want, err := follower.dm.visor.GetSignedBlockBySeq(target)
	require.NoError(t, err)
	for _, n := range nodes {
		if n.down {
			require.Equal(t, uint64(0), headSeq(t, n.dm.visor))
			continue
		}
		got, err := n.dm.visor.GetSignedBlockBySeq(target)
		require.NoError(t, err)
		require.Equal(t, want.HashHeader(), got.HashHeader())
	}
	// Signer 2 has round 0 of block 2, so the turn passed
	require.Equal(t, 2, c.rule.Proposer(2, 0))
	for seq := uint64(1); seq <= target; seq++ {
		b, err := follower.dm.visor.GetSignedBlockBySeq(seq)
		require.NoError(t, err)
		signers := c.rule.BlockSigners(*b)
		require.GreaterOrEqual(t, len(signers), 3)
		require.NotContains(t, signers, c.pks[2])
	}
}
