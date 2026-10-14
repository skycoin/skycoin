package coin

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/pqsig"
	"github.com/skycoin/skycoin/src/testutil"
)

// pqTestTxn spends one output at a secp256k1 address and one at a post-quantum address.
func pqTestTxn(t *testing.T) (Transaction, UxArray, cipher.SecKey, pqsig.SecKey) {
	_, secpSk := cipher.GenerateKeyPair()
	secpAddr := cipher.MustAddressFromSecKey(secpSk)
	pqSk, err := pqsig.NewSecKey()
	require.NoError(t, err)
	pqPk, err := pqSk.PubKey()
	require.NoError(t, err)
	pqAddr := cipher.AddressFromPQPubKey(pqPk)
	require.Equal(t, cipher.AddressVersionPQ, pqAddr.Version)

	uxIn := UxArray{
		{Body: UxBody{SrcTransaction: testutil.RandSHA256(t), Address: secpAddr, Coins: 5e6, Hours: 10}},
		{Body: UxBody{SrcTransaction: testutil.RandSHA256(t), Address: pqAddr, Coins: 5e6, Hours: 10}},
	}
	var txn Transaction
	for _, ux := range uxIn {
		require.NoError(t, txn.PushInput(ux.Hash()))
	}
	require.NoError(t, txn.PushOutput(testutil.MakeAddress(), 10e6, 5))
	txn.InnerHash = txn.HashInner()
	return txn, uxIn, secpSk, pqSk
}

func signPQTestTxn(t *testing.T, txn *Transaction, secpSk cipher.SecKey, pqSk pqsig.SecKey) {
	require.NoError(t, txn.SignInput(secpSk, 0))
	require.NoError(t, txn.SignPQInput(pqSk, 1))
	require.NoError(t, txn.UpdateHeader())
}

func TestPQTransaction(t *testing.T) {
	txn, uxIn, secpSk, pqSk := pqTestTxn(t)
	require.True(t, txn.IsFullyUnsigned())

	require.NoError(t, txn.SignInput(secpSk, 0))
	require.NoError(t, txn.UpdateHeader())
	require.False(t, txn.IsFullySigned())
	require.NoError(t, txn.VerifyUnsigned())
	require.NoError(t, txn.VerifyPartialInputSignatures(uxIn))

	require.NoError(t, txn.SignPQInput(pqSk, 1))
	require.Error(t, txn.SignPQInput(pqSk, 1), "already signed")
	require.NoError(t, txn.UpdateHeader())
	require.Equal(t, TxnTypePQ, txn.Type)
	require.True(t, txn.IsFullySigned())
	require.NoError(t, txn.Verify())
	require.NoError(t, txn.VerifyInputSignatures(uxIn))

	// The witness is not part of the transaction ID.
	id := txn.Hash()
	stripped := txn
	stripped.PQ = nil
	require.Equal(t, id, stripped.Hash())

	// The raw form carries the witness.
	raw, err := txn.SerializeHexWithPQ()
	require.NoError(t, err)
	back, err := DeserializeTransactionHexWithPQ(raw)
	require.NoError(t, err)
	require.Equal(t, txn.PQ, back.PQ)
	require.Equal(t, id, back.Hash())
	require.NoError(t, back.VerifyInputSignatures(uxIn))
}

func TestPQTransactionRejects(t *testing.T) {
	txn, uxIn, secpSk, pqSk := pqTestTxn(t)
	signPQTestTxn(t, &txn, secpSk, pqSk)

	t.Run("missing witness", func(t *testing.T) {
		bad := txn
		bad.PQ = nil
		require.Error(t, bad.Verify())
	})

	t.Run("bad signature", func(t *testing.T) {
		bad := txn
		bad.PQ = append([]PQInputSig(nil), txn.PQ...)
		bad.PQ[1].Sig = append([]byte(nil), txn.PQ[1].Sig...)
		bad.PQ[1].Sig[0] ^= 1
		require.Equal(t, pqsig.ErrInvalidSig, bad.Verify())
	})

	t.Run("key that does not own the address", func(t *testing.T) {
		other, err := pqsig.NewSecKey()
		require.NoError(t, err)
		bad, _, _, _ := pqTestTxn(t)
		bad.In, bad.Out, bad.InnerHash = txn.In, txn.Out, txn.InnerHash
		signPQTestTxn(t, &bad, secpSk, other)
		require.NoError(t, bad.Verify(), "the signature matches its own key")
		require.Error(t, bad.VerifyInputSignatures(uxIn), "but that key does not own the input")
	})

	t.Run("secp256k1 signature on a post-quantum input", func(t *testing.T) {
		bad, _, _, _ := pqTestTxn(t)
		bad.In, bad.Out, bad.InnerHash = txn.In, txn.Out, txn.InnerHash
		require.NoError(t, bad.SignInput(secpSk, 0))
		require.NoError(t, bad.SignInput(secpSk, 1))
		require.NoError(t, bad.UpdateHeader())
		require.Error(t, bad.VerifyInputSignatures(uxIn))
	})

	t.Run("post-quantum signature on a secp256k1 input", func(t *testing.T) {
		bad, _, _, _ := pqTestTxn(t)
		bad.In, bad.Out, bad.InnerHash = txn.In, txn.Out, txn.InnerHash
		require.NoError(t, bad.SignPQInput(pqSk, 0))
		require.NoError(t, bad.SignPQInput(pqSk, 1))
		require.NoError(t, bad.UpdateHeader())
		require.Error(t, bad.VerifyInputSignatures(uxIn))
	})

	t.Run("standard type with a witness", func(t *testing.T) {
		bad := txn
		bad.Type = TxnTypeStandard
		require.Error(t, bad.Verify())
	})
}

func TestPQWitnessDecodeLimits(t *testing.T) {
	_, _, err := DecodePQWitness([]byte{1, 0, 0, 0, 0xff, 0xff, 0, 0})
	require.Equal(t, ErrPQWitnessMalformed, err)
	_, _, err = DecodePQWitness([]byte{1, 0})
	require.Equal(t, ErrPQWitnessMalformed, err)
}
