package visor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

// newPQTestVisor makes a visor with the genesis block and the post-quantum rule.
func newPQTestVisor(t *testing.T, rule coin.PQBlockRule, pqSeckey pqsig.SecKey) *Visor {
	db, shutdown := prepareDB(t)
	t.Cleanup(shutdown)

	bc, err := NewBlockchain(db, BlockchainConfig{Pubkey: genPublic, PQBlocks: rule})
	require.NoError(t, err)
	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	cfg := NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress
	cfg.PQBlocks = rule
	if pqSeckey != nil {
		cfg.IsBlockPublisher = true
		cfg.BlockchainSeckey = genSecret
		cfg.BlockchainPQSeckey = pqSeckey
	}

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

func TestVisorPQBlocks(t *testing.T) {
	pqSeckey, err := pqsig.NewSecKey()
	require.NoError(t, err)
	pqPubkey, err := pqSeckey.PubKey()
	require.NoError(t, err)
	rule := coin.PQBlockRule{Pubkey: pqPubkey, ForkSeq: 1}

	// The publisher signs block 1 with both keys and stores both signatures.
	pub := newPQTestVisor(t, rule, pqSeckey)
	require.NoError(t, pub.Config.Verify())

	var gb *coin.SignedBlock
	require.NoError(t, pub.db.View("", func(tx *dbutil.Tx) error {
		gb, err = pub.blockchain.GetGenesisBlock(tx)
		return err
	}))
	require.Empty(t, gb.PQSig, "the genesis block is exempt")

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	txn := makeUnspentsTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, 10, params.UserVerifyTxn.MaxDropletPrecision)
	require.NoError(t, pub.db.Update("", func(tx *dbutil.Tx) error {
		_, softErr, err := pub.unconfirmed.InjectTransaction(tx, pub.blockchain, txn, params.MainNetDistribution, pub.Config.UnconfirmedVerifyTxn)
		require.Nil(t, softErr)
		return err
	}))
	pub.Config.MaxBlockTransactionsSize, err = txn.Size()
	require.NoError(t, err)

	sb, err := pub.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Len(t, sb.PQSig, pqsig.SigSize)
	require.NoError(t, sb.VerifyPQSignature(pqPubkey))

	var stored *coin.SignedBlock
	require.NoError(t, pub.db.View("", func(tx *dbutil.Tx) error {
		stored, err = pub.blockchain.GetSignedBlockBySeq(tx, 1)
		return err
	}))
	require.Equal(t, sb.PQSig, stored.PQSig)

	execute := func(v *Visor, b coin.SignedBlock) error {
		return v.db.Update("", func(tx *dbutil.Tx) error {
			return v.executeSignedBlock(tx, b)
		})
	}

	// A follower under the rule rejects block 1 without its signature, or with a bad one.
	follower := newPQTestVisor(t, rule, nil)
	missing := sb
	missing.PQSig = nil
	require.Equal(t, pqsig.ErrMissingSig, execute(follower, missing))

	bad := sb
	bad.PQSig = append([]byte(nil), sb.PQSig...)
	bad.PQSig[0] ^= 1
	require.Equal(t, pqsig.ErrInvalidSig, execute(follower, bad))

	require.NoError(t, execute(follower, sb))

	// A follower whose fork height is later takes the block without one.
	later := newPQTestVisor(t, coin.PQBlockRule{Pubkey: pqPubkey, ForkSeq: 2}, nil)
	require.NoError(t, execute(later, missing))
}

func TestConfigVerifyPQ(t *testing.T) {
	pqSeckey, err := pqsig.NewSecKey()
	require.NoError(t, err)
	pqPubkey, err := pqSeckey.PubKey()
	require.NoError(t, err)
	otherSeckey, err := pqsig.NewSecKey()
	require.NoError(t, err)

	cfg := NewConfig()
	cfg.Distribution = params.MainNetDistribution
	cfg.BlockchainPubkey = genPublic
	cfg.BlockchainSeckey = genSecret
	cfg.IsBlockPublisher = true
	cfg.PQBlocks = coin.PQBlockRule{Pubkey: pqPubkey, ForkSeq: 10}
	require.ErrorContains(t, cfg.Verify(), "post-quantum secret key is not set")

	cfg.BlockchainPQSeckey = otherSeckey
	require.ErrorContains(t, cfg.Verify(), "invalid post-quantum seckey for pubkey")

	cfg.BlockchainPQSeckey = pqSeckey
	require.NoError(t, cfg.Verify())

	cfg.IsBlockPublisher = false
	cfg.PQBlocks.Pubkey = pqPubkey[:100]
	require.ErrorContains(t, cfg.Verify(), "invalid post-quantum blockchain pubkey")
}

func TestVerifyPQAddresses(t *testing.T) {
	pqAddr := cipher.AddressFromPQPubKey(make([]byte, pqsig.PubKeySize))
	payPQ := coin.Transaction{Out: []coin.TransactionOutput{{Address: pqAddr, Coins: 1e6}}}
	spendPQ := coin.Transaction{Type: coin.TxnTypePQ}
	plain := coin.Transaction{Out: []coin.TransactionOutput{{Address: genAddress, Coins: 1e6}}}

	never := Blockchain{cfg: BlockchainConfig{}}
	require.NoError(t, never.verifyPQAddresses(plain, 100))
	require.Error(t, never.verifyPQAddresses(payPQ, 100))
	require.Error(t, never.verifyPQAddresses(spendPQ, 100))

	fork := Blockchain{cfg: BlockchainConfig{PQAddressForkSeq: 10}}
	require.Error(t, fork.verifyPQAddresses(payPQ, 9))
	require.Error(t, fork.verifyPQAddresses(spendPQ, 9))
	require.NoError(t, fork.verifyPQAddresses(payPQ, 10))
	require.NoError(t, fork.verifyPQAddresses(spendPQ, 10))
}
