package coin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/encoder"
	"github.com/SkycoinProject/skycoin/src/testutil"
	_require "github.com/SkycoinProject/skycoin/src/testutil/require"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
)

func makeTransactionFromUxOuts(t *testing.T, uxs []UxOut, secs []cipher.SecKey) Transaction {
	require.Equal(t, len(uxs), len(secs))

	txn := Transaction{}

	err := txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 5e6, 50)
	require.NoError(t, err)

	for _, ux := range uxs {
		err = txn.PushInput(ux.Hash())
		require.NoError(t, err)
	}

	txn.SignInputs(secs)

	err = txn.UpdateHeader()
	require.NoError(t, err)
	return txn
}

func makeTransactionFromUxOut(t *testing.T, ux UxOut, s cipher.SecKey) Transaction {
	return makeTransactionFromUxOuts(t, []UxOut{ux}, []cipher.SecKey{s})
}

func makeTransaction(t *testing.T) Transaction {
	ux, s := makeUxOutWithSecret(t)
	return makeTransactionFromUxOut(t, ux, s)
}

func makeTransactionMultipleInputs(t *testing.T, n int) (Transaction, []cipher.SecKey) {
	uxs := make([]UxOut, n)
	secs := make([]cipher.SecKey, n)
	for i := 0; i < n; i++ {
		ux, s := makeUxOutWithSecret(t)
		uxs[i] = ux
		secs[i] = s
	}
	return makeTransactionFromUxOuts(t, uxs, secs), secs
}

func makeTransactions(t *testing.T, n int) Transactions { //nolint:unparam
	txns := make(Transactions, n)
	for i := range txns {
		txns[i] = makeTransaction(t)
	}
	return txns
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func copyTransaction(txn Transaction) Transaction {
	txo := Transaction{}
	txo.Length = txn.Length
	txo.Type = txn.Type
	txo.InnerHash = txn.InnerHash
	txo.Sigs = make([]cipher.Sig, len(txn.Sigs))
	copy(txo.Sigs, txn.Sigs)
	txo.In = make([]cipher.SHA256, len(txn.In))
	copy(txo.In, txn.In)
	txo.Out = make([]TransactionOutput, len(txn.Out))
	copy(txo.Out, txn.Out)
	return txo
}

func TestTransactionVerify(t *testing.T) {
	// Mismatch header hash
	txn := makeTransaction(t)
	txn.InnerHash = cipher.SHA256{}
	testutil.RequireError(t, txn.Verify(), "InnerHash does not match computed hash")

	// No inputs
	txn = makeTransaction(t)
	txn.In = make([]cipher.SHA256, 0)
	err := txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "No inputs")

	// No outputs
	txn = makeTransaction(t)
	txn.Out = make([]TransactionOutput, 0)
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "No outputs")

	// Invalid number of sigs
	txn = makeTransaction(t)
	txn.Sigs = make([]cipher.Sig, 0)
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Invalid number of signatures")
	txn.Sigs = make([]cipher.Sig, 20)
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Invalid number of signatures")

	// Too many sigs & inputs
	txn = makeTransaction(t)
	txn.Sigs = make([]cipher.Sig, math.MaxUint16+1)
	txn.In = make([]cipher.SHA256, math.MaxUint16+1)
	testutil.RequireError(t, txn.Verify(), "Too many signatures and inputs")

	// Duplicate inputs
	ux, s := makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	err = txn.PushInput(txn.In[0])
	require.NoError(t, err)
	txn.Sigs = nil
	txn.SignInputs([]cipher.SecKey{s, s})
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Duplicate spend")

	// Duplicate outputs
	txn = makeTransaction(t)
	to := txn.Out[0]
	err = txn.PushOutput(to.Address, to.Coins, to.Hours)
	require.NoError(t, err)
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Duplicate output in transaction")

	// Invalid signature, empty
	txn = makeTransaction(t)
	txn.Sigs[0] = cipher.Sig{}
	testutil.RequireError(t, txn.Verify(), "Unsigned input in transaction")

	// Invalid signature, not empty
	// A stable invalid signature must be used because random signatures could appear valid
	// Note: Transaction.Verify() only checks that the signature is a minimally valid signature
	badSig := "9a0f86874a4d9541f58a1de4db1c1b58765a868dc6f027445d0a2a8a7bddd1c45ea559fcd7bef45e1b76ccdaf8e50bbebd952acbbea87d1cb3f7a964bc89bf1ed5"
	txn = makeTransaction(t)
	txn.Sigs[0] = cipher.MustSigFromHex(badSig)
	testutil.RequireError(t, txn.Verify(), "Failed to recover pubkey from signature")

	// We can't check here for other invalid signatures:
	//      - Signatures signed by someone else, spending coins they don't own
	//      - Signatures signing a different message
	// This must be done by blockchain tests, because we need the address
	// from the unspent being spent
	// The verification here only checks that the signature is valid at all

	// Output coins are 0
	txn = makeTransaction(t)
	txn.Out[0].Coins = 0
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Zero coin output")

	// Output coin overflow
	txn = makeTransaction(t)
	txn.Out[0].Coins = math.MaxUint64 - 3e6
	err = txn.UpdateHeader()
	require.NoError(t, err)
	testutil.RequireError(t, txn.Verify(), "Output coins overflow")

	// Output coins are not multiples of 1e6 (valid, decimal restriction is not enforced here)
	txn = makeTransaction(t)
	txn.Out[0].Coins += 10
	err = txn.UpdateHeader()
	require.NoError(t, err)
	txn.Sigs = nil
	txn.SignInputs([]cipher.SecKey{genSecret})
	require.NotEqual(t, txn.Out[0].Coins%1e6, uint64(0))
	require.NoError(t, txn.Verify())

	// Valid
	txn = makeTransaction(t)
	txn.Out[0].Coins = 10e6
	txn.Out[1].Coins = 1e6
	err = txn.UpdateHeader()
	require.NoError(t, err)
	require.NoError(t, txn.Verify())
}

func TestTransactionVerifyUnsigned(t *testing.T) {
	txn, _ := makeTransactionMultipleInputs(t, 2)
	err := txn.VerifyUnsigned()
	testutil.RequireError(t, err, "Unsigned transaction must contain a null signature")

	// Invalid signature, not empty
	// A stable invalid signature must be used because random signatures could appear valid
	// Note: Transaction.Verify() only checks that the signature is a minimally valid signature
	badSig := "9a0f86874a4d9541f58a1de4db1c1b58765a868dc6f027445d0a2a8a7bddd1c45ea559fcd7bef45e1b76ccdaf8e50bbebd952acbbea87d1cb3f7a964bc89bf1ed5"
	txn, _ = makeTransactionMultipleInputs(t, 2)
	txn.Sigs[0] = cipher.Sig{}
	txn.Sigs[1] = cipher.MustSigFromHex(badSig)
	testutil.RequireError(t, txn.VerifyUnsigned(), "Failed to recover pubkey from signature")

	txn.Sigs = nil
	err = txn.VerifyUnsigned()
	testutil.RequireError(t, err, "Invalid number of signatures")

	// Transaction is unsigned if at least 1 signature is null
	txn, _ = makeTransactionMultipleInputs(t, 3)
	require.True(t, len(txn.Sigs) > 1)
	txn.Sigs[0] = cipher.Sig{}
	err = txn.VerifyUnsigned()
	require.NoError(t, err)

	// Transaction is unsigned if all signatures are null
	for i := range txn.Sigs {
		txn.Sigs[i] = cipher.Sig{}
	}
	err = txn.VerifyUnsigned()
	require.NoError(t, err)
}

func TestTransactionVerifyInput(t *testing.T) {
	// Invalid uxIn args
	txn := makeTransaction(t)
	_require.PanicsWithLogMessage(t, "txn.In != uxIn", func() {
		_ = txn.VerifyInputSignatures(nil) //nolint:errcheck
	})
	_require.PanicsWithLogMessage(t, "txn.In != uxIn", func() {
		_ = txn.VerifyInputSignatures(UxArray{}) //nolint:errcheck
	})
	_require.PanicsWithLogMessage(t, "txn.In != uxIn", func() {
		_ = txn.VerifyInputSignatures(make(UxArray, 3)) //nolint:errcheck
	})

	// txn.In != txn.Sigs
	ux, s := makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	txn.Sigs = []cipher.Sig{}
	_require.PanicsWithLogMessage(t, "txn.In != txn.Sigs", func() {
		_ = txn.VerifyInputSignatures(UxArray{ux}) //nolint:errcheck
	})

	ux, s = makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	txn.Sigs = append(txn.Sigs, cipher.Sig{})
	_require.PanicsWithLogMessage(t, "txn.In != txn.Sigs", func() {
		_ = txn.VerifyInputSignatures(UxArray{ux}) //nolint:errcheck
	})

	// txn.InnerHash != txn.HashInner()
	ux, s = makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	txn.InnerHash = cipher.SHA256{}
	_require.PanicsWithLogMessage(t, "Invalid Tx Inner Hash", func() {
		_ = txn.VerifyInputSignatures(UxArray{ux}) //nolint:errcheck
	})

	// txn.In does not match uxIn hashes
	ux, s = makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	_require.PanicsWithLogMessage(t, "Ux hash mismatch", func() {
		_ = txn.VerifyInputSignatures(UxArray{UxOut{}}) //nolint:errcheck
	})

	// Unsigned txn
	ux, s = makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	txn.Sigs[0] = cipher.Sig{}
	err := txn.VerifyInputSignatures(UxArray{ux})
	testutil.RequireError(t, err, "Unsigned input in transaction")

	// Signature signed by someone else
	ux, _ = makeUxOutWithSecret(t)
	_, s2 := makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s2)
	err = txn.VerifyInputSignatures(UxArray{ux})
	testutil.RequireError(t, err, "Signature not valid for output being spent")

	// Valid
	ux, s = makeUxOutWithSecret(t)
	txn = makeTransactionFromUxOut(t, ux, s)
	err = txn.VerifyInputSignatures(UxArray{ux})
	require.NoError(t, err)
}

func TestTransactionPushInput(t *testing.T) {
	txn := &Transaction{}
	ux := makeUxOut(t)
	require.NoError(t, txn.PushInput(ux.Hash()))
	require.Equal(t, len(txn.In), 1)
	require.Equal(t, txn.In[0], ux.Hash())
	txn.In = append(txn.In, make([]cipher.SHA256, math.MaxUint16)...)
	ux = makeUxOut(t)
	err := txn.PushInput(ux.Hash())
	testutil.RequireError(t, err, "Max transaction inputs reached")
}

func TestTransactionPushOutput(t *testing.T) {
	txn := &Transaction{}
	a := makeAddress()
	err := txn.PushOutput(a, 100, 150)
	require.NoError(t, err)
	require.Equal(t, len(txn.Out), 1)
	require.Equal(t, txn.Out[0], TransactionOutput{
		Address: a,
		Coins:   100,
		Hours:   150,
	})
	for i := 1; i < 20; i++ {
		a := makeAddress()
		err := txn.PushOutput(a, uint64(i*100), uint64(i*50))
		require.NoError(t, err)
		require.Equal(t, len(txn.Out), i+1)
		require.Equal(t, txn.Out[i], TransactionOutput{
			Address: a,
			Coins:   uint64(i * 100),
			Hours:   uint64(i * 50),
		})
	}

	txn.Out = append(txn.Out, make([]TransactionOutput, math.MaxUint16-len(txn.Out))...)
	err = txn.PushOutput(a, 999, 999)
	testutil.RequireError(t, err, "Max transaction outputs reached")
}

func TestTransactionSignInput(t *testing.T) {
	txn, seckeys := makeTransactionMultipleInputs(t, 3)
	require.True(t, txn.IsFullySigned())

	// Input is already signed
	err := txn.SignInput(seckeys[0], 0)
	testutil.RequireError(t, err, "Input already signed")
	require.True(t, txn.IsFullySigned())

	// Input is not signed
	txn.Sigs[1] = cipher.Sig{}
	require.False(t, txn.IsFullySigned())
	err = txn.SignInput(seckeys[1], 1)
	require.NoError(t, err)
	require.True(t, txn.IsFullySigned())
	err = txn.SignInput(seckeys[1], 1)
	testutil.RequireError(t, err, "Input already signed")

	// Transaction has no sigs; sigs array is initialized
	txn.Sigs = nil
	require.False(t, txn.IsFullySigned())
	err = txn.SignInput(seckeys[2], 2)
	require.NoError(t, err)
	require.False(t, txn.IsFullySigned())
	require.Len(t, txn.Sigs, 3)
	require.True(t, txn.Sigs[0].Null())
	require.True(t, txn.Sigs[1].Null())
	require.False(t, txn.Sigs[2].Null())

	// SignInputs on a partially signed transaction fails
	require.Panics(t, func() {
		txn.SignInputs(seckeys)
	})

	// Signing the rest of the inputs individually works
	err = txn.SignInput(seckeys[1], 1)
	require.NoError(t, err)
	require.False(t, txn.IsFullySigned())
	err = txn.SignInput(seckeys[0], 0)
	require.NoError(t, err)
	require.True(t, txn.IsFullySigned())

	// Can use SignInputs on allocated array of empty sigs
	txn.Sigs = make([]cipher.Sig, 3)
	txn.SignInputs(seckeys)
	require.True(t, txn.IsFullySigned())
}

func TestTransactionSignInputs(t *testing.T) {
	txn := &Transaction{}
	// Panics if txns already signed
	txn.Sigs = append(txn.Sigs, cipher.Sig{})
	require.Panics(t, func() { txn.SignInputs([]cipher.SecKey{}) })
	// Panics if not enough keys
	txn = &Transaction{}
	ux, s := makeUxOutWithSecret(t)
	err := txn.PushInput(ux.Hash())
	require.NoError(t, err)
	ux2, s2 := makeUxOutWithSecret(t)
	err = txn.PushInput(ux2.Hash())
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 40, 80)
	require.NoError(t, err)
	require.Equal(t, len(txn.Sigs), 0)
	require.Panics(t, func() { txn.SignInputs([]cipher.SecKey{s}) })
	require.Equal(t, len(txn.Sigs), 0)
	// Valid signing
	h := txn.HashInner()
	require.NotPanics(t, func() { txn.SignInputs([]cipher.SecKey{s, s2}) })
	require.Equal(t, len(txn.Sigs), 2)
	h2 := txn.HashInner()
	require.Equal(t, h2, h)
	p := cipher.MustPubKeyFromSecKey(s)
	a := cipher.AddressFromPubKey(p)
	p = cipher.MustPubKeyFromSecKey(s2)
	a2 := cipher.AddressFromPubKey(p)
	require.NoError(t, cipher.VerifyAddressSignedHash(a, txn.Sigs[0], cipher.AddSHA256(h, txn.In[0])))
	require.NoError(t, cipher.VerifyAddressSignedHash(a2, txn.Sigs[1], cipher.AddSHA256(h, txn.In[1])))
	require.Error(t, cipher.VerifyAddressSignedHash(a, txn.Sigs[1], h))
	require.Error(t, cipher.VerifyAddressSignedHash(a2, txn.Sigs[0], h))
}

func TestTransactionHash(t *testing.T) {
	txn := makeTransaction(t)
	h := txn.Hash()
	h2 := txn.HashInner()
	require.NotEqual(t, h, cipher.SHA256{})
	require.NotEqual(t, h2, h)
}

func TestTransactionUpdateHeader(t *testing.T) {
	txn := makeTransaction(t)
	h := txn.InnerHash
	txn.InnerHash = cipher.SHA256{}
	err := txn.UpdateHeader()
	require.NoError(t, err)
	require.NotEqual(t, txn.InnerHash, cipher.SHA256{})
	require.Equal(t, txn.InnerHash, h)
	require.Equal(t, txn.InnerHash, txn.HashInner())
}

func TestTransactionHashInner(t *testing.T) {
	txn := makeTransaction(t)

	require.NotEqual(t, cipher.SHA256{}, txn.HashInner())

	// If txn.In is changed, inner hash should change
	txn2 := copyTransaction(txn)
	ux := makeUxOut(t)
	txn2.In[0] = ux.Hash()
	require.NotEqual(t, txn, txn2)
	require.Equal(t, txn2.In[0], ux.Hash())
	require.NotEqual(t, txn.HashInner(), txn2.HashInner())

	// If txn.Out is changed, inner hash should change
	txn2 = copyTransaction(txn)
	a := makeAddress()
	txn2.Out[0].Address = a
	require.NotEqual(t, txn, txn2)
	require.Equal(t, txn2.Out[0].Address, a)
	require.NotEqual(t, txn.HashInner(), txn2.HashInner())

	// If txn.Head is changed, inner hash should not change
	txn2 = copyTransaction(txn)
	txn.Sigs = append(txn.Sigs, cipher.Sig{})
	require.Equal(t, txn.HashInner(), txn2.HashInner())
}

func TestTransactionSerialization(t *testing.T) {
	txn := makeTransaction(t)
	b, err := txn.Serialize()
	require.NoError(t, err)
	txn2, err := DeserializeTransaction(b)
	require.NoError(t, err)
	require.Equal(t, txn, txn2)

	// Check reserializing deserialized txn
	b2, err := txn2.Serialize()
	require.NoError(t, err)
	txn3, err := DeserializeTransaction(b2)
	require.NoError(t, err)
	require.Equal(t, txn2, txn3)

	// Check hex encode/decode followed by deserialize
	s := hex.EncodeToString(b)
	sb, err := hex.DecodeString(s)
	require.NoError(t, err)
	txn4, err := DeserializeTransaction(sb)
	require.NoError(t, err)
	require.Equal(t, txn2, txn4)

	// Invalid deserialization
	require.Panics(t, func() {
		MustDeserializeTransaction([]byte{0x04})
	})

	// SerializeHex
	x, err := txn.SerializeHex()
	require.NoError(t, err)
	txn5, err := DeserializeTransactionHex(x)
	require.NoError(t, err)
	require.Equal(t, txn, txn5)

	// Invalid hex deserialization
	require.Panics(t, func() {
		MustDeserializeTransactionHex("foo")
	})

	ss, err := txn.Serialize()
	require.NoError(t, err)
	require.Equal(t, ss, txn.MustSerialize())
	sshh, err := txn.SerializeHex()
	require.NoError(t, err)
	require.Equal(t, sshh, txn.MustSerializeHex())
}

func TestTransactionOutputHours(t *testing.T) {
	txn := Transaction{}
	err := txn.PushOutput(makeAddress(), 1e6, 100)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 200)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 500)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 1e6, 0)
	require.NoError(t, err)
	hours, err := txn.OutputHours()
	require.NoError(t, err)
	require.Equal(t, hours, uint64(800))

	err = txn.PushOutput(makeAddress(), 1e6, math.MaxUint64-700)
	require.NoError(t, err)
	_, err = txn.OutputHours()
	testutil.RequireError(t, err, "Transaction output hours overflow")
}

func TestTransactionsSize(t *testing.T) {
	txns := makeTransactions(t, 10)
	var size uint32
	for _, txn := range txns {
		encodedLen, err := mathutil.IntToUint32(len(encoder.Serialize(&txn)))
		require.NoError(t, err)
		size, err = mathutil.AddUint32(size, encodedLen)
		require.NoError(t, err)
	}

	require.NotEqual(t, size, 0)
	s, err := txns.Size()
	require.NoError(t, err)
	require.Equal(t, s, size)
}

func TestTransactionsHashes(t *testing.T) {
	txns := make(Transactions, 4)
	for i := 0; i < len(txns); i++ {
		txns[i] = makeTransaction(t)
	}
	hashes := txns.Hashes()
	require.Equal(t, len(hashes), 4)
	for i, h := range hashes {
		require.Equal(t, h, txns[i].Hash())
	}
}

func TestTransactionsTruncateBytesTo(t *testing.T) {
	txns := makeTransactions(t, 10)
	var trunc uint32
	for i := 0; i < len(txns)/2; i++ {
		size, err := txns[i].Size()
		require.NoError(t, err)
		trunc, err = mathutil.AddUint32(trunc, size)
		require.NoError(t, err)
	}

	// Truncating halfway
	txns2, err := txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, len(txns2), len(txns)/2)
	totalSize, err := txns2.Size()
	require.NoError(t, err)
	require.Equal(t, totalSize, trunc)

	// Stepping into next boundary has same cutoff, must exceed
	trunc++
	txns2, err = txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, len(txns2), len(txns)/2)
	totalSize, err = txns2.Size()
	require.NoError(t, err)
	require.Equal(t, totalSize, trunc-1)

	// Moving to 1 before next level
	size5, err := txns[5].Size()
	require.NoError(t, err)
	require.True(t, size5 >= 2)
	trunc, err = mathutil.AddUint32(trunc, size5-2)
	require.NoError(t, err)
	txns2, err = txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, len(txns2), len(txns)/2)

	totalSize, err = txns2.Size()
	require.NoError(t, err)
	size5, err = txns[5].Size()
	require.NoError(t, err)
	require.Equal(t, totalSize, trunc-size5+1)

	// Moving to next level
	trunc++
	txns2, err = txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, len(txns2), len(txns)/2+1)
	size, err := txns2.Size()
	require.NoError(t, err)
	require.Equal(t, size, trunc)

	// Truncating to full available amt
	trunc, err = txns.Size()
	require.NoError(t, err)
	txns2, err = txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, txns, txns2)
	size, err = txns2.Size()
	require.NoError(t, err)
	require.Equal(t, size, trunc)

	// Truncating over amount
	trunc++
	txns2, err = txns.TruncateBytesTo(trunc)
	require.NoError(t, err)
	require.Equal(t, txns, txns2)
	size, err = txns2.Size()
	require.NoError(t, err)
	require.Equal(t, size, trunc-1)

	// Truncating to 0
	trunc = 0
	txns2, err = txns.TruncateBytesTo(0)
	require.NoError(t, err)
	require.Equal(t, len(txns2), 0)
	size, err = txns2.Size()
	require.NoError(t, err)
	require.Equal(t, size, trunc)
}

func TestVerifyTransactionCoinsSpending(t *testing.T) {
	// Input coins overflow
	// Insufficient coins
	// Destroy coins

	type ux struct {
		coins uint64
		hours uint64
	}

	cases := []struct {
		name   string
		inUxs  []ux
		outUxs []ux
		err    error
	}{
		{
			name: "Input coins overflow",
			inUxs: []ux{
				{
					coins: math.MaxUint64 - 1e6 + 1,
					hours: 10,
				},
				{
					coins: 1e6,
					hours: 0,
				},
			},
			err: errors.New("Transaction input coins overflow"),
		},

		{
			name: "Output coins overflow",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: math.MaxUint64 - 10e6 + 1,
					hours: 0,
				},
				{
					coins: 20e6,
					hours: 1,
				},
			},
			err: errors.New("Transaction output coins overflow"),
		},

		{
			name: "Insufficient coins",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 20e6,
					hours: 1,
				},
				{
					coins: 10e6,
					hours: 1,
				},
			},
			err: errors.New("Insufficient coins"),
		},

		{
			name: "Destroyed coins",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 5e6,
					hours: 1,
				},
				{
					coins: 10e6,
					hours: 1,
				},
			},
			err: errors.New("Transactions may not destroy coins"),
		},

		{
			name: "valid",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 10e6,
					hours: 11,
				},
				{
					coins: 10e6,
					hours: 1,
				},
				{
					coins: 5e6,
					hours: 0,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var uxIn, uxOut UxArray

			for _, ch := range tc.inUxs {
				uxIn = append(uxIn, UxOut{
					Body: UxBody{
						Coins: ch.coins,
						Hours: ch.hours,
					},
				})
			}

			for _, ch := range tc.outUxs {
				uxOut = append(uxOut, UxOut{
					Body: UxBody{
						Coins: ch.coins,
						Hours: ch.hours,
					},
				})
			}

			err := VerifyTransactionCoinsSpending(uxIn, uxOut)
			require.Equal(t, tc.err, err)
		})
	}
}

func TestVerifyTransactionHoursSpending(t *testing.T) {
	// Input hours overflow
	// Insufficient hours
	// NOTE: does not check for hours overflow, that had to be moved to soft constraints
	// NOTE: if uxIn.CoinHours() fails during the addition of earned hours to base hours,
	// the error is ignored and treated as 0 hours

	type ux struct {
		coins uint64
		hours uint64
	}

	cases := []struct {
		name     string
		inUxs    []ux
		outUxs   []ux
		headTime uint64
		err      string
	}{
		{
			name: "Input hours overflow",
			inUxs: []ux{
				{
					coins: 3e6,
					hours: math.MaxUint64 - 1e6 + 1,
				},
				{
					coins: 1e6,
					hours: 1e6,
				},
			},
			err: "Transaction input hours overflow",
		},

		{
			name: "Insufficient coin hours",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 15e6,
					hours: 10,
				},
				{
					coins: 10e6,
					hours: 11,
				},
			},
			err: "Insufficient coin hours",
		},

		{
			name: "coin hours time calculation overflow",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 10e6,
					hours: 11,
				},
				{
					coins: 10e6,
					hours: 1,
				},
				{
					coins: 5e6,
					hours: 0,
				},
			},
			headTime: math.MaxUint64,
			err:      "UxOut.CoinHours: Calculating whole coin seconds overflows uint64 seconds=18446744073709551615 coins=10 uxid=",
		},

		{
			name:     "Invalid (coin hours overflow when adding earned hours, which is treated as 0, and now enough coin hours)",
			headTime: 1e6,
			inUxs: []ux{
				{
					coins: 10e6,
					hours: math.MaxUint64,
				},
			},
			outUxs: []ux{
				{
					coins: 10e6,
					hours: 1,
				},
			},
			err: "Insufficient coin hours",
		},

		{
			name:     "Valid (coin hours overflow when adding earned hours, which is treated as 0, but not sending any hours)",
			headTime: 1e6,
			inUxs: []ux{
				{
					coins: 10e6,
					hours: math.MaxUint64,
				},
			},
			outUxs: []ux{
				{
					coins: 10e6,
					hours: 0,
				},
			},
		},

		{
			name: "Valid (base inputs have insufficient coin hours, but have sufficient after adjusting coinhours by headTime)",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 15e6,
					hours: 10,
				},
				{
					coins: 10e6,
					hours: 11,
				},
			},
			headTime: 1492707255,
		},

		{
			name: "valid",
			inUxs: []ux{
				{
					coins: 10e6,
					hours: 10,
				},
				{
					coins: 15e6,
					hours: 10,
				},
			},
			outUxs: []ux{
				{
					coins: 10e6,
					hours: 11,
				},
				{
					coins: 10e6,
					hours: 1,
				},
				{
					coins: 5e6,
					hours: 0,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var uxIn, uxOut UxArray

			for _, ch := range tc.inUxs {
				uxIn = append(uxIn, UxOut{
					Body: UxBody{
						Coins: ch.coins,
						Hours: ch.hours,
					},
				})
			}

			for _, ch := range tc.outUxs {
				uxOut = append(uxOut, UxOut{
					Body: UxBody{
						Coins: ch.coins,
						Hours: ch.hours,
					},
				})
			}

			err := VerifyTransactionHoursSpending(tc.headTime, uxIn, uxOut)
			if tc.err == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.True(t, strings.HasPrefix(err.Error(), tc.err))
			}
		})
	}
}

func TestTransactionsFees(t *testing.T) {
	calc := func(txn *Transaction) (uint64, error) {
		return 1, nil
	}

	var txns Transactions

	// Nil txns
	fee, err := txns.Fees(calc)
	require.NoError(t, err)
	require.Equal(t, uint64(0), fee)

	txns = append(txns, Transaction{})
	txns = append(txns, Transaction{})

	// 2 transactions, calc() always returns 1
	fee, err = txns.Fees(calc)
	require.NoError(t, err)
	require.Equal(t, uint64(2), fee)

	// calc error
	failingCalc := func(txn *Transaction) (uint64, error) {
		return 0, errors.New("bad calc")
	}
	_, err = txns.Fees(failingCalc)
	testutil.RequireError(t, err, "bad calc")

	// summing of calculated fees overflows
	overflowCalc := func(txn *Transaction) (uint64, error) {
		return math.MaxUint64, nil
	}

	_, err = txns.Fees(overflowCalc)
	testutil.RequireError(t, err, "Transactions fee totals overflow")
}

func TestSortTransactions(t *testing.T) {
	n := 6
	var txns Transactions
	for i := 0; i < n; i++ {
		txn := Transaction{}
		err := txn.PushOutput(makeAddress(), 1e6, uint64(i*1e3))
		require.NoError(t, err)
		err = txn.UpdateHeader()
		require.NoError(t, err)
		txns = append(txns, txn)
	}

	hashSortedTxns := append(Transactions{}, txns...)

	sort.Slice(hashSortedTxns, func(i, j int) bool {
		ihash := hashSortedTxns[i].Hash()
		jhash := hashSortedTxns[j].Hash()
		return bytes.Compare(ihash[:], jhash[:]) < 0
	})

	cases := []struct {
		name       string
		feeCalc    FeeCalculator
		txns       Transactions
		sortedTxns Transactions
	}{
		{
			name:       "already sorted",
			txns:       Transactions{txns[0], txns[1]},
			sortedTxns: Transactions{txns[0], txns[1]},
			feeCalc: func(txn *Transaction) (uint64, error) {
				return 1e8 - txn.Out[0].Hours, nil
			},
		},

		{
			name:       "reverse sorted",
			txns:       Transactions{txns[1], txns[0]},
			sortedTxns: Transactions{txns[0], txns[1]},
			feeCalc: func(txn *Transaction) (uint64, error) {
				return 1e8 - txn.Out[0].Hours, nil
			},
		},

		{
			name:       "hash tiebreaker",
			txns:       Transactions{hashSortedTxns[1], hashSortedTxns[0]},
			sortedTxns: Transactions{hashSortedTxns[0], hashSortedTxns[1]},
			feeCalc: func(txn *Transaction) (uint64, error) {
				return 1e8, nil
			},
		},

		{
			name:       "invalid fee multiplication is capped",
			txns:       Transactions{txns[1], txns[2], txns[0]},
			sortedTxns: Transactions{txns[2], txns[0], txns[1]},
			feeCalc: func(txn *Transaction) (uint64, error) {
				if txn.Hash() == txns[2].Hash() {
					return math.MaxUint64 / 2, nil
				}
				return 1e8 - txn.Out[0].Hours, nil
			},
		},

		{
			name:       "failed fee calc is filtered",
			txns:       Transactions{txns[1], txns[2], txns[0]},
			sortedTxns: Transactions{txns[0], txns[1]},
			feeCalc: func(txn *Transaction) (uint64, error) {
				if txn.Hash() == txns[2].Hash() {
					return 0, errors.New("fee calc failed")
				}
				return 1e8 - txn.Out[0].Hours, nil
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txns, err := SortTransactions(tc.txns, tc.feeCalc)
			require.NoError(t, err)
			require.Equal(t, tc.sortedTxns, txns)
		})
	}
}

func TestTransactionSignedUnsigned(t *testing.T) {
	txn, _ := makeTransactionMultipleInputs(t, 2)
	require.True(t, txn.IsFullySigned())
	require.True(t, txn.hasNonNullSignature())
	require.False(t, txn.IsFullyUnsigned())
	require.False(t, txn.hasNullSignature())

	txn.Sigs[1] = cipher.Sig{}
	require.False(t, txn.IsFullySigned())
	require.True(t, txn.hasNonNullSignature())
	require.False(t, txn.IsFullyUnsigned())
	require.True(t, txn.hasNullSignature())

	txn.Sigs[0] = cipher.Sig{}
	require.False(t, txn.IsFullySigned())
	require.False(t, txn.hasNonNullSignature())
	require.True(t, txn.IsFullyUnsigned())
	require.True(t, txn.hasNullSignature())

	txn.Sigs = nil
	require.False(t, txn.IsFullySigned())
	require.False(t, txn.hasNonNullSignature())
	require.True(t, txn.IsFullyUnsigned())
	require.False(t, txn.hasNullSignature())
}
