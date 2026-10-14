package daemon

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
)

func pqTestTxnWithWitness(t *testing.T) coin.Transaction {
	sk, err := pqsig.NewSecKey()
	require.NoError(t, err)
	var txn coin.Transaction
	require.NoError(t, txn.PushInput(testutil.RandSHA256(t)))
	require.NoError(t, txn.PushOutput(testutil.MakeAddress(), 1e6, 1))
	txn.InnerHash = txn.HashInner()
	require.NoError(t, txn.SignPQInput(sk, 0))
	require.NoError(t, txn.UpdateHeader())
	return txn
}

func TestGivePQWitnessesMessageEncoding(t *testing.T) {
	txn := pqTestTxnWithWitness(t)
	msgs := newGivePQWitnessesMessages([]coin.Transaction{{}, txn})
	require.Len(t, msgs, 1)
	require.Len(t, msgs[0].Witnesses, 1, "standard transactions have no witness")

	m := msgs[0]
	buf := make([]byte, m.EncodeSize())
	require.NoError(t, m.Encode(buf))
	var got GivePQWitnessesMessage
	n, err := got.Decode(buf)
	require.NoError(t, err)
	require.Equal(t, uint64(len(buf)), n)
	require.Equal(t, m.Witnesses, got.Witnesses)

	_, err = got.Decode(buf[:len(buf)-1])
	require.Equal(t, encoder.ErrBufferUnderflow, err)
}

func TestNewGivePQWitnessesMessagesSplitsBySize(t *testing.T) {
	txns := make([]coin.Transaction, 80)
	for i := range txns {
		txns[i] = pqTestTxnWithWitness(t)
	}
	msgs := newGivePQWitnessesMessages(txns)
	require.Greater(t, len(msgs), 1)
	total := 0
	for _, m := range msgs {
		require.LessOrEqual(t, m.EncodeSize()+4, uint64(256*1024))
		total += len(m.Witnesses)
	}
	require.Equal(t, len(txns), total)
}

func TestAttachPQWitness(t *testing.T) {
	txn := pqTestTxnWithWitness(t)
	dm := &Daemon{pqWitnesses: newPQSigCache()}
	dm.addPQWitnesses(newGivePQWitnessesMessages([]coin.Transaction{txn})[0].Witnesses)

	stripped := txn
	stripped.PQ = nil
	dm.attachPQWitness(&stripped)
	require.Equal(t, txn.PQ, stripped.PQ)
	require.NoError(t, stripped.Verify())

	other := pqTestTxnWithWitness(t)
	other.PQ = nil
	dm.attachPQWitness(&other)
	require.Empty(t, other.PQ, "no witness was received for it")
}
