package coin

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/testutil"
	_require "github.com/skycoin/skycoin/src/testutil/require"
)

func makeTransactionFromUxOut(ux UxOut, s cipher.SecKey) Transaction {
	tx := Transaction{}
	tx.PushInput(ux.Hash())
	tx.PushOutput(makeAddress(), 1e6, 50)
	tx.PushOutput(makeAddress(), 5e6, 50)
	tx.SignInputs([]cipher.SecKey{s})
	tx.UpdateHeader()
	return tx
}

func makeTransaction(t *testing.T) Transaction {
	ux, s := makeUxOutWithSecret(t)
	return makeTransactionFromUxOut(ux, s)
}

func makeTransactions(t *testing.T, n int) Transactions { // nolint: unparam
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

func copyTransaction(tx Transaction) Transaction {
	txo := Transaction{}
	txo.Length = tx.Length
	txo.Type = tx.Type
	txo.InnerHash = tx.InnerHash
	txo.Sigs = make([]cipher.Sig, len(tx.Sigs))
	copy(txo.Sigs, tx.Sigs)
	txo.In = make([]cipher.SHA256, len(tx.In))
	copy(txo.In, tx.In)
	txo.Out = make([]TransactionOutput, len(tx.Out))
	copy(txo.Out, tx.Out)
	return txo
}

func TestTransactionVerify(t *testing.T) {
	// Mismatch header hash
	tx := makeTransaction(t)
	tx.InnerHash = cipher.SHA256{}
	testutil.RequireError(t, tx.Verify(), "Invalid header hash")

	// No inputs
	tx = makeTransaction(t)
	tx.In = make([]cipher.SHA256, 0)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "No inputs")

	// No outputs
	tx = makeTransaction(t)
	tx.Out = make([]TransactionOutput, 0)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "No outputs")

	// Invalid number of sigs
	tx = makeTransaction(t)
	tx.Sigs = make([]cipher.Sig, 0)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Invalid number of signatures")
	tx.Sigs = make([]cipher.Sig, 20)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Invalid number of signatures")

	// Too many sigs & inputs
	tx = makeTransaction(t)
	tx.Sigs = make([]cipher.Sig, math.MaxUint16)
	tx.In = make([]cipher.SHA256, math.MaxUint16)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Too many signatures and inputs")

	// Duplicate inputs
	ux, s := makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	tx.PushInput(tx.In[0])
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{s, s})
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Duplicate spend")

	// Duplicate outputs
	tx = makeTransaction(t)
	to := tx.Out[0]
	tx.PushOutput(to.Address, to.Coins, to.Hours)
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Duplicate output in transaction")

	// Invalid signature, empty
	tx = makeTransaction(t)
	tx.Sigs[0] = cipher.Sig{}
	testutil.RequireError(t, tx.Verify(), "Failed to recover public key")
	// We can't check here for other invalid signatures:
	//      - Signatures signed by someone else, spending coins they don't own
	//      - Signature is for wrong hash
	// This must be done by blockchain tests, because we need the address
	// from the unspent being spent

	// Output coins are 0
	tx = makeTransaction(t)
	tx.Out[0].Coins = 0
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Zero coin output")

	// Output coin overflow
	tx = makeTransaction(t)
	tx.Out[0].Coins = math.MaxUint64 - 3e6
	tx.UpdateHeader()
	testutil.RequireError(t, tx.Verify(), "Output coins overflow")

	// Output coins are not multiples of 1e6 (valid, decimal restriction is not enforced here)
	tx = makeTransaction(t)
	tx.Out[0].Coins += 10
	tx.UpdateHeader()
	tx.Sigs = nil
	tx.SignInputs([]cipher.SecKey{genSecret})
	require.NotEqual(t, tx.Out[0].Coins%1e6, uint64(0))
	require.NoError(t, tx.Verify())

	// Valid
	tx = makeTransaction(t)
	tx.Out[0].Coins = 10e6
	tx.Out[1].Coins = 1e6
	tx.UpdateHeader()
	require.Nil(t, tx.Verify())
}

func TestTransactionVerifyInput(t *testing.T) {
	// Invalid uxIn args
	tx := makeTransaction(t)
	_require.PanicsWithLogMessage(t, "tx.In != uxIn", func() {
		tx.VerifyInput(nil)
	})
	_require.PanicsWithLogMessage(t, "tx.In != uxIn", func() {
		tx.VerifyInput(UxArray{})
	})
	_require.PanicsWithLogMessage(t, "tx.In != uxIn", func() {
		tx.VerifyInput(make(UxArray, 3))
	})

	// tx.In != tx.Sigs
	ux, s := makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	tx.Sigs = []cipher.Sig{}
	_require.PanicsWithLogMessage(t, "tx.In != tx.Sigs", func() {
		tx.VerifyInput(UxArray{ux})
	})

	ux, s = makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	tx.Sigs = append(tx.Sigs, cipher.Sig{})
	_require.PanicsWithLogMessage(t, "tx.In != tx.Sigs", func() {
		tx.VerifyInput(UxArray{ux})
	})

	// tx.InnerHash != tx.HashInner()
	ux, s = makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	tx.InnerHash = cipher.SHA256{}
	_require.PanicsWithLogMessage(t, "Invalid Tx Inner Hash", func() {
		tx.VerifyInput(UxArray{ux})
	})

	// tx.In does not match uxIn hashes
	ux, s = makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	_require.PanicsWithLogMessage(t, "Ux hash mismatch", func() {
		tx.VerifyInput(UxArray{UxOut{}})
	})

	// Invalid signature
	ux, s = makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	tx.Sigs[0] = cipher.Sig{}
	err := tx.VerifyInput(UxArray{ux})
	testutil.RequireError(t, err, "Signature not valid for output being spent")

	// Valid
	ux, s = makeUxOutWithSecret(t)
	tx = makeTransactionFromUxOut(ux, s)
	err = tx.VerifyInput(UxArray{ux})
	require.NoError(t, err)
}

func TestTransactionPushInput(t *testing.T) {
	tx := &Transaction{}
	ux := makeUxOut(t)
	require.Equal(t, tx.PushInput(ux.Hash()), uint16(0))
	require.Equal(t, len(tx.In), 1)
	require.Equal(t, tx.In[0], ux.Hash())
	tx.In = append(tx.In, make([]cipher.SHA256, math.MaxUint16)...)
	ux = makeUxOut(t)
	require.Panics(t, func() { tx.PushInput(ux.Hash()) })
}

func TestTransactionPushOutput(t *testing.T) {
	tx := &Transaction{}
	a := makeAddress()
	tx.PushOutput(a, 100, 150)
	require.Equal(t, len(tx.Out), 1)
	require.Equal(t, tx.Out[0], TransactionOutput{
		Address: a,
		Coins:   100,
		Hours:   150,
	})
	for i := 1; i < 20; i++ {
		a := makeAddress()
		tx.PushOutput(a, uint64(i*100), uint64(i*50))
		require.Equal(t, len(tx.Out), i+1)
		require.Equal(t, tx.Out[i], TransactionOutput{
			Address: a,
			Coins:   uint64(i * 100),
			Hours:   uint64(i * 50),
		})
	}
}

func TestTransactionSignInputs(t *testing.T) {
	tx := &Transaction{}
	// Panics if txns already signed
	tx.Sigs = append(tx.Sigs, cipher.Sig{})
	require.Panics(t, func() { tx.SignInputs([]cipher.SecKey{}) })
	// Panics if not enough keys
	tx = &Transaction{}
	ux, s := makeUxOutWithSecret(t)
	tx.PushInput(ux.Hash())
	ux2, s2 := makeUxOutWithSecret(t)
	tx.PushInput(ux2.Hash())
	tx.PushOutput(makeAddress(), 40, 80)
	require.Equal(t, len(tx.Sigs), 0)
	require.Panics(t, func() { tx.SignInputs([]cipher.SecKey{s}) })
	require.Equal(t, len(tx.Sigs), 0)
	// Valid signing
	h := tx.HashInner()
	require.NotPanics(t, func() { tx.SignInputs([]cipher.SecKey{s, s2}) })
	require.Equal(t, len(tx.Sigs), 2)
	require.Equal(t, tx.HashInner(), h)
	p := cipher.PubKeyFromSecKey(s)
	a := cipher.AddressFromPubKey(p)
	p = cipher.PubKeyFromSecKey(s2)
	a2 := cipher.AddressFromPubKey(p)
	require.Nil(t, cipher.ChkSig(a, cipher.AddSHA256(h, tx.In[0]), tx.Sigs[0]))
	require.Nil(t, cipher.ChkSig(a2, cipher.AddSHA256(h, tx.In[1]), tx.Sigs[1]))
	require.NotNil(t, cipher.ChkSig(a, h, tx.Sigs[1]))
	require.NotNil(t, cipher.ChkSig(a2, h, tx.Sigs[0]))
}

func TestTransactionHash(t *testing.T) {
	tx := makeTransaction(t)
	require.NotEqual(t, tx.Hash(), cipher.SHA256{})
	require.NotEqual(t, tx.HashInner(), tx.Hash())
}

func TestTransactionUpdateHeader(t *testing.T) {
	tx := makeTransaction(t)
	h := tx.InnerHash
	tx.InnerHash = cipher.SHA256{}
	tx.UpdateHeader()
	require.NotEqual(t, tx.InnerHash, cipher.SHA256{})
	require.Equal(t, tx.InnerHash, h)
	require.Equal(t, tx.InnerHash, tx.HashInner())
}

func TestTransactionHashInner(t *testing.T) {
	tx := makeTransaction(t)

	h := tx.HashInner()
	require.NotEqual(t, h, cipher.SHA256{})

	// If tx.In is changed, hash should change
	tx2 := copyTransaction(tx)
	ux := makeUxOut(t)
	tx2.In[0] = ux.Hash()
	require.NotEqual(t, tx, tx2)
	require.Equal(t, tx2.In[0], ux.Hash())
	require.NotEqual(t, tx.HashInner(), tx2.HashInner())

	// If tx.Out is changed, hash should change
	tx2 = copyTransaction(tx)
	a := makeAddress()
	tx2.Out[0].Address = a
	require.NotEqual(t, tx, tx2)
	require.Equal(t, tx2.Out[0].Address, a)
	require.NotEqual(t, tx.HashInner(), tx2.HashInner())

	// If tx.Head is changed, hash should not change
	tx2 = copyTransaction(tx)
	tx.Sigs = append(tx.Sigs, cipher.Sig{})
	require.Equal(t, tx.HashInner(), tx2.HashInner())
}

func TestTransactionSerialization(t *testing.T) {
	tx := makeTransaction(t)
	b := tx.Serialize()
	tx2, err := TransactionDeserialize(b)
	require.NoError(t, err)
	require.Equal(t, tx, tx2)
	// Invalid deserialization
	require.Panics(t, func() { MustTransactionDeserialize([]byte{0x04}) })
}

func TestTransactionOutputHours(t *testing.T) {
	tx := Transaction{}
	tx.PushOutput(makeAddress(), 1e6, 100)
	tx.PushOutput(makeAddress(), 1e6, 200)
	tx.PushOutput(makeAddress(), 1e6, 500)
	tx.PushOutput(makeAddress(), 1e6, 0)
	hours, err := tx.OutputHours()
	require.NoError(t, err)
	require.Equal(t, hours, uint64(800))

	tx.PushOutput(makeAddress(), 1e6, math.MaxUint64-700)
	_, err = tx.OutputHours()
	testutil.RequireError(t, err, "Transaction output hours overflow")
}

type outAddr struct {
	Addr  cipher.Address
	Coins uint64
	Hours uint64
}

func makeTx(s cipher.SecKey, ux *UxOut, outs []outAddr, tm uint64, seq uint64) (*Transaction, UxArray, error) {
	if ux == nil {
		// genesis block tx.
		tx := Transaction{}
		tx.PushOutput(outs[0].Addr, outs[0].Coins, outs[0].Hours)
		_, s = cipher.GenerateKeyPair()
		ux := UxOut{
			Head: UxHead{
				Time:  100,
				BkSeq: 0,
			},
			Body: UxBody{
				SrcTransaction: tx.InnerHash,
				Address:        outs[0].Addr,
				Coins:          outs[0].Coins,
				Hours:          outs[0].Hours,
			},
		}
		return &tx, []UxOut{ux}, nil
	}

	tx := Transaction{}
	tx.PushInput(ux.Hash())
	tx.SignInputs([]cipher.SecKey{s})
	for _, o := range outs {
		tx.PushOutput(o.Addr, o.Coins, o.Hours)
	}
	tx.UpdateHeader()

	uxo := make(UxArray, len(tx.Out))
	for i := range tx.Out {
		uxo[i] = UxOut{
			Head: UxHead{
				Time:  tm,
				BkSeq: seq,
			},
			Body: UxBody{
				SrcTransaction: tx.Hash(),
				Address:        tx.Out[i].Address,
				Coins:          tx.Out[i].Coins,
				Hours:          tx.Out[i].Hours,
			},
		}
	}
	return &tx, uxo, nil
}

func TestTransactionsSize(t *testing.T) {
	txns := makeTransactions(t, 10)
	size := 0
	for _, tx := range txns {
		size += len(encoder.Serialize(&tx))
	}
	require.NotEqual(t, size, 0)
	require.Equal(t, txns.Size(), size)
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
	trunc := 0
	for i := 0; i < len(txns)/2; i++ {
		trunc += txns[i].Size()
	}
	// Truncating halfway
	txns2 := txns.TruncateBytesTo(trunc)
	require.Equal(t, len(txns2), len(txns)/2)
	require.Equal(t, txns2.Size(), trunc)

	// Stepping into next boundary has same cutoff, must exceed
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	require.Equal(t, len(txns2), len(txns)/2)
	require.Equal(t, txns2.Size(), trunc-1)

	// Moving to 1 before next level
	trunc += txns[5].Size() - 2
	txns2 = txns.TruncateBytesTo(trunc)
	require.Equal(t, len(txns2), len(txns)/2)
	require.Equal(t, txns2.Size(), trunc-txns[5].Size()+1)

	// Moving to next level
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	require.Equal(t, len(txns2), len(txns)/2+1)
	require.Equal(t, txns2.Size(), trunc)

	// Truncating to full available amt
	trunc = txns.Size()
	txns2 = txns.TruncateBytesTo(trunc)
	require.Equal(t, txns, txns2)
	require.Equal(t, txns2.Size(), trunc)

	// Truncating over amount
	trunc++
	txns2 = txns.TruncateBytesTo(trunc)
	require.Equal(t, txns, txns2)
	require.Equal(t, txns2.Size(), trunc-1)

	// Truncating to 0
	trunc = 0
	txns2 = txns.TruncateBytesTo(0)
	require.Equal(t, len(txns2), 0)
	require.Equal(t, txns2.Size(), trunc)
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
	calc := func(tx *Transaction) (uint64, error) {
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
	failingCalc := func(tx *Transaction) (uint64, error) {
		return 0, errors.New("bad calc")
	}
	_, err = txns.Fees(failingCalc)
	testutil.RequireError(t, err, "bad calc")

	// summing of calculated fees overflows
	overflowCalc := func(tx *Transaction) (uint64, error) {
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
		txn.PushOutput(makeAddress(), 1e6, uint64(i*1e3))
		txn.UpdateHeader()
		txns = append(txns, txn)
	}

	var hashSortedTxns Transactions
	for _, txn := range txns {
		hashSortedTxns = append(hashSortedTxns, txn)
	}

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
			txns := SortTransactions(tc.txns, tc.feeCalc)
			require.Equal(t, tc.sortedTxns, txns)
		})
	}
}
