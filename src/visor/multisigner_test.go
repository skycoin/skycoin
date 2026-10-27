package visor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// newMultiSignerTestVisor makes a visor with the genesis block and the
// multi-signer rule. A non-nil seckey makes it a signer.
func newMultiSignerTestVisor(t *testing.T, rule coin.MultiSignerRule, seckey cipher.SecKey) *Visor {
	db, shutdown := prepareDB(t)
	t.Cleanup(shutdown)

	cfg := NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress
	cfg.MultiSigner = rule
	if !seckey.Null() {
		cfg.IsBlockPublisher = true
		cfg.Arbitrating = true
		cfg.BlockchainSeckey = seckey
	}

	bc, err := NewBlockchain(db, BlockchainConfig{Pubkey: genPublic, MultiSigner: rule, Arbitrating: cfg.Arbitrating})
	require.NoError(t, err)
	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	v := &Visor{
		Config:      cfg,
		unconfirmed: unconfirmed,
		blockchain:  bc,
		db:          db,
		history:     historydb.New(),
	}
	addGenesisBlockToVisor(t, v)
	return v
}

func TestVisorMultiSigner(t *testing.T) {
	pks := make([]cipher.PubKey, 4)
	sks := make([]cipher.SecKey, 4)
	for i := range pks {
		pks[i], sks[i] = cipher.GenerateKeyPair()
	}
	rule := coin.MultiSignerRule{Signers: pks, Threshold: 3, ForkSeq: 1}

	signers := make([]*Visor, 4)
	for i := range signers {
		signers[i] = newMultiSignerTestVisor(t, rule, sks[i])
		require.NoError(t, signers[i].Config.Verify())
		require.Equal(t, i, signers[i].SignerIndex())
	}
	follower := newMultiSignerTestVisor(t, rule, cipher.SecKey{})
	require.Equal(t, -1, follower.SignerIndex())

	execute := func(v *Visor, b coin.SignedBlock) error {
		return v.db.Update("", func(tx *dbutil.Tx) error {
			return v.executeSignedBlock(tx, b)
		})
	}

	// Signer 1 proposes block 1 from a transaction in its pool
	var gb *coin.SignedBlock
	var err error
	require.NoError(t, follower.db.View("", func(tx *dbutil.Tx) error {
		gb, err = follower.blockchain.GetGenesisBlock(tx)
		return err
	}))
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	txn := makeUnspentsTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, 10, params.UserVerifyTxn.MaxDropletPrecision)
	prop := signers[1]
	require.NoError(t, prop.db.Update("", func(tx *dbutil.Tx) error {
		_, softErr, err := prop.unconfirmed.InjectTransaction(tx, prop.blockchain, txn, params.MainNetDistribution, prop.Config.UnconfirmedVerifyTxn)
		require.Nil(t, softErr)
		return err
	}))
	prop.Config.MaxBlockTransactionsSize, err = txn.Size()
	require.NoError(t, err)

	sb, err := prop.CreateBlockProposal(time.Unix(int64(genTime)+100, 0)) //nolint:gosec
	require.NoError(t, err)
	require.Equal(t, uint64(1), sb.Seq())
	require.Equal(t, 1, rule.ProposerIndex(sb))
	head, err := prop.GetHeadBlock()
	require.NoError(t, err)
	require.Equal(t, uint64(0), head.Seq(), "a proposal is not executed")

	// The proposal alone, or with one cosignature, is not final
	require.Equal(t, coin.ErrNotEnoughBlockSigners{Have: 1, Need: 3}, execute(follower, sb))
	for _, v := range signers {
		require.NoError(t, v.VerifyBlock(sb))
	}
	c2, err := signers[2].SignBlockCosig(sb.Block)
	require.NoError(t, err)
	c3, err := signers[3].SignBlockCosig(sb.Block)
	require.NoError(t, err)

	two := sb
	two.Cosigs = []coin.BlockCosig{c2}
	require.Equal(t, coin.ErrNotEnoughBlockSigners{Have: 2, Need: 3}, execute(follower, two))

	// A block signed by the old single publisher key is not accepted after the fork
	single := sb
	single.Sig = cipher.MustSignHash(sb.HashHeader(), genSecret)
	single.Cosigs = []coin.BlockCosig{c2, c3}
	require.Equal(t, coin.ErrBlockProposerNotSigner, execute(follower, single))

	// k signatures make it final. Junk cosignatures are dropped before storing.
	final := sb
	_, outSK := cipher.GenerateKeyPair()
	final.Cosigs = []coin.BlockCosig{c2, c2, c3, {Signer: 0, Sig: cipher.MustSignHash(sb.HashHeader(), outSK)}}
	require.NoError(t, execute(follower, final))
	stored, err := follower.GetSignedBlockBySeq(1)
	require.NoError(t, err)
	require.Equal(t, []coin.BlockCosig{c2, c3}, stored.Cosigs)

	// The startup check verifies the stored signatures under the rule
	require.NoError(t, CheckDatabase(follower.db, genPublic, coin.PQBlockRule{}, rule, nil))
	hash := sb.HashHeader()
	require.NoError(t, follower.db.Update("", func(tx *dbutil.Tx) error {
		return dbutil.Delete(tx, blockdb.BlockCosigsBkt, hash[:])
	}))
	require.Equal(t, coin.ErrNotEnoughBlockSigners{Have: 1, Need: 3}, CheckDatabase(follower.db, genPublic, coin.PQBlockRule{}, rule, nil))
}

func TestVisorBlockVotes(t *testing.T) {
	pk, sk := cipher.GenerateKeyPair()
	v := newMultiSignerTestVisor(t, coin.MultiSignerRule{Signers: []cipher.PubKey{pk}, Threshold: 1, ForkSeq: 1}, sk)

	h1 := cipher.SumSHA256([]byte("one"))
	h2 := cipher.SumSHA256([]byte("two"))

	_, ok, err := v.BlockVote(5)
	require.NoError(t, err)
	require.False(t, ok)

	require.NoError(t, v.RecordBlockVote(5, h1))
	require.NoError(t, v.RecordBlockVote(5, h1), "the same block again is fine")
	require.Equal(t, ErrAlreadyVoted{Seq: 5, Hash: h1}, v.RecordBlockVote(5, h2))
	require.NoError(t, v.RecordBlockVote(6, h2))

	got, ok, err := v.BlockVote(5)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, h1, got)
}

func TestConfigVerifyMultiSigner(t *testing.T) {
	pks := make([]cipher.PubKey, 4)
	sks := make([]cipher.SecKey, 4)
	for i := range pks {
		pks[i], sks[i] = cipher.GenerateKeyPair()
	}

	cfg := NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = genPublic
	cfg.MultiSigner = coin.MultiSignerRule{Signers: pks, Threshold: 3, ForkSeq: 1}
	require.NoError(t, cfg.Verify(), "a plain node")

	cfg.IsBlockPublisher = true
	cfg.BlockchainSeckey = genSecret
	require.ErrorContains(t, cfg.Verify(), "not one of block_signers")

	cfg.BlockchainSeckey = sks[2]
	require.NoError(t, cfg.Verify(), "a signer whose key is not the genesis key")

	cfg.MultiSigner.Threshold = 2
	require.ErrorContains(t, cfg.Verify(), "more than 2/3")
	cfg.MultiSigner.Threshold = 3

	cfg.PQBlocks = coin.PQBlockRule{Pubkey: []byte{1}, ForkSeq: 1}
	cfg.IsBlockPublisher = false
	require.Error(t, cfg.Verify())
}

// A signer whose key is not the genesis key uses the configured genesis
// signature instead of signing the genesis block itself.
func TestMultiSignerGenesis(t *testing.T) {
	pk, sk := cipher.GenerateKeyPair()
	rule := coin.MultiSignerRule{Signers: []cipher.PubKey{pk}, Threshold: 1, ForkSeq: 1}

	db, shutdown := prepareDB(t)
	defer shutdown()
	bc, err := NewBlockchain(db, BlockchainConfig{Pubkey: genPublic, MultiSigner: rule})
	require.NoError(t, err)

	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)

	cfg := NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = genPublic
	cfg.BlockchainSeckey = sk
	cfg.IsBlockPublisher = true
	cfg.MultiSigner = rule
	cfg.GenesisAddress = genAddress
	cfg.GenesisCoinVolume = genCoins
	cfg.GenesisTimestamp = genTime
	cfg.GenesisSignature = cipher.MustSignHash(gb.HashHeader(), genSecret)

	v := &Visor{Config: cfg, blockchain: bc, db: db, history: historydb.New()}
	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)
	v.unconfirmed = unconfirmed

	require.NoError(t, db.Update("", v.maybeCreateGenesisBlock))
	got, err := v.GetSignedBlockBySeq(0)
	require.NoError(t, err)
	require.Equal(t, cfg.GenesisSignature, got.Sig)
}
