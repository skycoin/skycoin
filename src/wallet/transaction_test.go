package wallet_test

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/transaction"
	"github.com/SkycoinProject/skycoin/src/util/fee"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/core/collection"
)

func TestWalletSignTransaction(t *testing.T) {
	txnSigned, uxs, seckeys := makeTransaction(t, 4)
	require.Equal(t, 4, len(uxs))
	require.Equal(t, 4, len(seckeys))
	txnUnsigned := txnSigned
	txnUnsigned.Sigs = make([]cipher.Sig, len(txnSigned.Sigs))

	w := &collection.Wallet{}
	for _, x := range seckeys {
		p := cipher.MustPubKeyFromSecKey(x)
		a := cipher.AddressFromPubKey(p)
		err := w.AddEntry(wallet.Entry{
			Address: a,
			Public:  p,
			Secret:  x,
		})
		require.NoError(t, err)

		// Add unrelated entries
		err = w.AddEntry(makeEntry())
		require.NoError(t, err)
	}

	badTxnInnerHash, _, _ := makeTransaction(t, 2)
	badTxnInnerHash.InnerHash = testutil.RandSHA256(t)

	badTxnNoInputs, _, _ := makeTransaction(t, 2)
	badTxnNoInputs.Sigs = make([]cipher.Sig, len(badTxnNoInputs.Sigs))
	badTxnNoInputs.In = nil
	err := badTxnNoInputs.UpdateHeader()
	require.NoError(t, err)

	badTxnNoSigs, _, _ := makeTransaction(t, 2)
	badTxnNoSigs.Sigs = nil
	err = badTxnNoSigs.UpdateHeader()
	require.NoError(t, err)

	txnOtherWallet, uxsOtherWallet, secKeysOtherWallet := makeTransaction(t, 4)

	txnPartiallySigned := txnOtherWallet
	err = txnPartiallySigned.Verify()
	require.NoError(t, err)
	err = txnPartiallySigned.VerifyInputSignatures(uxsOtherWallet)
	require.NoError(t, err)

	txnPartiallySigned.Sigs = make([]cipher.Sig, len(txnOtherWallet.Sigs))
	copy(txnPartiallySigned.Sigs, txnOtherWallet.Sigs)
	txnPartiallySigned.Sigs[1] = cipher.Sig{}
	txnPartiallySigned.Sigs[2] = cipher.Sig{}
	err = txnPartiallySigned.UpdateHeader()
	require.NoError(t, err)

	txnPartiallySigned2 := txnPartiallySigned
	txnPartiallySigned2.Sigs = make([]cipher.Sig, len(txnPartiallySigned.Sigs))
	copy(txnPartiallySigned2.Sigs, txnPartiallySigned.Sigs)
	err = txnPartiallySigned2.UpdateHeader()
	require.NoError(t, err)

	txnOtherWallet.Sigs = make([]cipher.Sig, len(txnOtherWallet.Sigs))
	err = txnOtherWallet.UpdateHeader()
	require.NoError(t, err)

	otherWallet := &collection.Wallet{}
	for i := 1; i < 3; i++ {
		p := cipher.MustPubKeyFromSecKey(secKeysOtherWallet[i])
		a := cipher.AddressFromPubKey(p)
		err := otherWallet.AddEntry(wallet.Entry{
			Address: a,
			Public:  p,
			Secret:  secKeysOtherWallet[i],
		})
		require.NoError(t, err)
	}

	cases := []struct {
		name        string
		w           wallet.Wallet
		txn         coin.Transaction
		signIndexes []int
		uxOuts      []coin.UxOut
		err         error
		partial     bool
		complete    bool
	}{
		{
			name:   "signed txn",
			w:      w,
			txn:    txnSigned,
			uxOuts: uxs,
			err:    wallet.NewError(errors.New("Transaction is fully signed")),
		},

		{
			name:        "partially signed txn signing with same index",
			w:           w,
			txn:         txnPartiallySigned,
			uxOuts:      uxs,
			signIndexes: []int{3},
			err:         wallet.NewError(errors.New("Transaction is already signed at index 3")),
		},

		{
			name: "bad txn inner hash",
			w:    w,
			txn:  badTxnInnerHash,
			err:  wallet.NewError(errors.New("Transaction inner hash does not match computed inner hash")),
		},

		{
			name: "txn no inputs",
			w:    w,
			txn:  badTxnNoInputs,
			err:  wallet.NewError(errors.New("No transaction inputs to sign")),
		},

		{
			name: "txn no sigs",
			w:    w,
			txn:  badTxnNoSigs,
			err:  wallet.NewError(errors.New("Transaction signatures array is empty")),
		},

		{
			name:   "len uxouts does not match len inputs",
			w:      w,
			txn:    txnUnsigned,
			uxOuts: uxs[:2],
			err:    errors.New("len(uxOuts) != len(txn.In)"),
		},

		{
			name:        "too many sign indexes",
			w:           w,
			txn:         txnUnsigned,
			uxOuts:      uxs,
			signIndexes: []int{0, 1, 2, 3, 4, 5},
			err:         wallet.NewError(errors.New("Number of signature indexes exceeds number of inputs")),
		},

		{
			name:        "sign index out of range",
			w:           w,
			txn:         txnUnsigned,
			uxOuts:      uxs,
			signIndexes: []int{0, 1, 5, 2},
			err:         wallet.NewError(errors.New("Signature index out of range")),
		},

		{
			name:        "duplicate value in sign indexes",
			w:           w,
			txn:         txnUnsigned,
			uxOuts:      uxs,
			signIndexes: []int{0, 1, 1},
			err:         wallet.NewError(errors.New("Duplicate value in signature indexes")),
		},

		{
			name:   "wallet cannot sign any input",
			w:      w,
			txn:    txnOtherWallet,
			uxOuts: uxsOtherWallet,
			err:    wallet.NewError(errors.New("Wallet cannot sign all requested inputs")),
		},

		{
			name:   "wallet cannot sign some inputs",
			w:      otherWallet,
			txn:    txnOtherWallet,
			uxOuts: uxsOtherWallet,
			err:    wallet.NewError(errors.New("Wallet cannot sign all requested inputs")),
		},

		{
			name:        "wallet cannot sign all specified inputs",
			w:           otherWallet,
			txn:         txnOtherWallet,
			uxOuts:      uxsOtherWallet,
			signIndexes: []int{2, 0},
			err:         wallet.NewError(errors.New("Wallet cannot sign all requested inputs")),
		},

		{
			name:   "valid unsigned txn, all sigs",
			w:      w,
			txn:    txnUnsigned,
			uxOuts: uxs,
			err:    nil,
		},

		{
			name:        "valid unsigned txn, some sigs defined",
			w:           w,
			txn:         txnUnsigned,
			signIndexes: []int{1, 2},
			uxOuts:      uxs,
			err:         nil,
		},

		{
			name:        "valid unsigned txn, all sigs defined",
			w:           w,
			txn:         txnUnsigned,
			signIndexes: []int{0, 1, 2, 3},
			uxOuts:      uxs,
			err:         nil,
		},

		{
			name:        "valid unsigned txn, all sigs defined, unordered",
			w:           w,
			txn:         txnUnsigned,
			signIndexes: []int{2, 1, 3, 0},
			uxOuts:      uxs,
			err:         nil,
		},

		{
			name:        "valid, wallet can sign the specified inputs, but not others",
			w:           otherWallet,
			txn:         txnOtherWallet,
			uxOuts:      uxsOtherWallet,
			signIndexes: []int{2},
			err:         nil,
		},

		{
			name:        "valid, transaction partially signed, unfinished signing",
			w:           otherWallet,
			txn:         txnPartiallySigned,
			uxOuts:      uxsOtherWallet,
			signIndexes: []int{2},
			err:         nil,
			partial:     true,
			complete:    false,
		},

		{
			name:        "valid, transaction partially signed, finished signing",
			w:           otherWallet,
			txn:         txnPartiallySigned2,
			uxOuts:      uxsOtherWallet,
			signIndexes: []int{1, 2},
			err:         nil,
			partial:     true,
			complete:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			signedTxn, err := wallet.SignTransaction(tc.w, &tc.txn, tc.signIndexes, tc.uxOuts)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)

			// The original txn should not be modified
			require.False(t, reflect.DeepEqual(tc.txn, *signedTxn))

			if len(tc.signIndexes) == 0 || len(tc.signIndexes) == len(tc.uxOuts) || tc.complete {
				// Transaction should be fully signed
				require.False(t, signedTxn.IsFullyUnsigned())
				err = signedTxn.Verify()
				require.NoError(t, err)
				err = signedTxn.VerifyInputSignatures(tc.uxOuts)
				require.NoError(t, err)
			} else {
				// index of a valid signature should be found in the signIndexes
				for _, x := range tc.signIndexes {
					require.False(t, signedTxn.Sigs[x].Null())
				}

				if !tc.partial {
					// Number of signatures should equal length of signIndexes
					nSigned := 0
					for _, s := range signedTxn.Sigs {
						if !s.Null() {
							nSigned++
						}
					}

					require.Equal(t, len(tc.signIndexes), nSigned)
				}
			}
		})
	}
}

func TestWalletCreateTransaction(t *testing.T) {
	headTime := uint64(time.Now().UTC().Unix())
	seed := []byte("seed")

	// Generate first keys
	_, secKeys := cipher.MustGenerateDeterministicKeyPairsSeed(seed, 11)
	secKey := secKeys[0]
	addr := cipher.MustAddressFromSecKey(secKey)

	var extraWalletAddrs []cipher.Address
	for _, s := range secKeys[1:] {
		extraWalletAddrs = append(extraWalletAddrs, cipher.MustAddressFromSecKey(s))
	}

	// Create unspent outputs
	var uxouts []coin.UxOut
	var originalUxouts []coin.UxOut
	addrs := []cipher.Address{}
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, uint64(100+i))
		uxout.Head.Time = headTime
		uxouts = append(uxouts, uxout)
		originalUxouts = append(originalUxouts, uxout)

		a := testutil.MakeAddress()
		addrs = append(addrs, a)
	}

	w := &collection.Wallet{}
	for _, x := range secKeys {
		p := cipher.MustPubKeyFromSecKey(x)
		a := cipher.AddressFromPubKey(p)
		err := w.AddEntry(wallet.Entry{
			Address: a,
			Public:  p,
			Secret:  x,
		})
		require.NoError(t, err)

		// Add unrelated entries
		err = w.AddEntry(makeEntry())
		require.NoError(t, err)
	}

	// shuffle the uxouts to test that the uxout sorting during spend selection is working
	rand.Shuffle(len(uxouts), func(i, j int) {
		uxouts[i], uxouts[j] = uxouts[j], uxouts[i]
	})

	// Create extra unspent outputs. These have the same value as uxouts, but are spendable by
	// keys held in extraWalletAddrs
	extraUxouts := make([][]coin.UxOut, len(extraWalletAddrs))
	for j := range extraWalletAddrs {
		s := secKeys[j+1]

		var uxouts []coin.UxOut
		for i := 0; i < 10; i++ {
			uxout := makeUxOut(t, s, 2e6, uint64(100+i))
			uxout.Head.Time = headTime
			uxouts = append(uxouts, uxout)
		}

		extraUxouts[j] = uxouts
	}

	// Create unspent outputs with no hours
	var uxoutsNoHours []coin.UxOut
	for i := 0; i < 10; i++ {
		uxout := makeUxOut(t, secKey, 2e6, 0)
		uxout.Head.Time = headTime
		uxoutsNoHours = append(uxoutsNoHours, uxout)
	}

	// shuffle the uxouts to test that the uxout sorting during spend selection is working
	rand.Shuffle(len(uxoutsNoHours), func(i, j int) {
		uxoutsNoHours[i], uxoutsNoHours[j] = uxoutsNoHours[j], uxoutsNoHours[i]
	})

	changeAddress := testutil.MakeAddress()

	validParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		ChangeAddress: &changeAddress,
		To: []coin.TransactionOutput{
			{
				Address: addrs[0],
				Hours:   10,
				Coins:   1e6,
			},
		},
	}

	unknownAddress := testutil.MakeAddress()

	cases := []struct {
		name            string
		err             error
		params          transaction.Params
		unspents        []coin.UxOut
		addressUnspents coin.AddressUxOuts
		chosenUnspents  []coin.UxOut
		headTime        uint64
		changeOutput    *coin.TransactionOutput
		toExpectedHours []uint64
	}{
		{
			name:   "params invalid",
			params: transaction.Params{},
			err:    transaction.NewError(errors.New("To is required")),
		},

		{
			name:   "unknown address in auxs",
			params: validParams,
			addressUnspents: coin.AddressUxOuts{
				unknownAddress: uxouts,
			},
			err: fmt.Errorf("Address %s from auxs not found in wallet", unknownAddress),
		},

		{
			name: "overflowing coin hours in params",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   math.MaxUint64,
						Coins:   1e6,
					},
					{
						Address: addrs[1],
						Hours:   1,
						Coins:   1e6,
					},
				},
			},
			err: transaction.NewError(errors.New("total output hours error: uint64 addition overflow")),
		},

		{
			name: "overflowing coins in params",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   math.MaxUint64,
					},
					{
						Address: addrs[1],
						Hours:   1,
						Coins:   1,
					},
				},
			},
			err: transaction.NewError(errors.New("total output coins error: uint64 addition overflow")),
		},

		{
			name: "no unspents",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   1e6,
					},
				},
			},
			err: transaction.ErrNoUnspents,
		},

		{
			name: "insufficient coins",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   100e6,
					},
				},
			},
			unspents: uxouts[:1],
			err:      transaction.ErrInsufficientBalance,
		},

		{
			name: "insufficient hours",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   100e6,
						Coins:   1e6,
					},
				},
			},
			unspents: uxouts[:1],
			err:      transaction.ErrInsufficientHours,
		},

		{
			name: "manual, 1 output, no change",
			params: transaction.Params{
				ChangeAddress: &changeAddress,
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   90,
						Coins:   2e6,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0]},
		},

		{
			name:     "no coin hours in inputs",
			unspents: uxoutsNoHours[:],
			params: transaction.Params{
				HoursSelection: transaction.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &changeAddress,
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   1e6,
					},
				},
			},
			err: fee.ErrTxnNoFee,
		},
	}

	type TestResult struct {
		Transaction *coin.Transaction
		Inputs      []transaction.UxBalance
		Unsigned    bool
	}

	results := make(map[string]TestResult, len(cases))

	compareResult := func(t *testing.T, a, b TestResult) {
		require.True(t, reflect.DeepEqual(a.Inputs, b.Inputs))
		require.True(t, reflect.DeepEqual(a.Transaction.Out, b.Transaction.Out))
		require.Equal(t, a.Transaction.InnerHash, b.Transaction.InnerHash)
		require.Equal(t, a.Transaction.Type, b.Transaction.Type)
		require.Equal(t, a.Transaction.Length, b.Transaction.Length)
		require.Equal(t, len(a.Transaction.Sigs), len(b.Transaction.Sigs))

		if a.Unsigned == b.Unsigned {
			require.Equal(t, a.Transaction.Hash(), b.Transaction.Hash())
			// Sigs have a nonce so will vary each run, unset them before comparing the whole transaction
			require.Equal(t, len(a.Transaction.Sigs), len(b.Transaction.Sigs))
			at := *a.Transaction
			bt := *b.Transaction
			at.Sigs = nil
			bt.Sigs = nil
			require.True(t, reflect.DeepEqual(at, bt))
		} else {
			require.NotEqual(t, a.Transaction.Hash(), b.Transaction.Hash())
		}
	}

	bools := []bool{true, false}
	for _, unsigned := range bools {
		for _, tc := range cases {
			name := fmt.Sprintf("unsigned=%v %s", unsigned, tc.name)
			t.Run(name, func(t *testing.T) {
				if tc.headTime == 0 {
					tc.headTime = headTime
				}

				addrUxOuts := coin.AddressUxOuts{
					addr: tc.unspents,
				}

				if tc.addressUnspents != nil {
					addrUxOuts = tc.addressUnspents
				}

				unspents := make(map[cipher.SHA256]coin.UxOut)
				for _, uxs := range addrUxOuts {
					for _, ux := range uxs {
						unspents[ux.Hash()] = ux
					}
				}

				var txn *coin.Transaction
				var inputs []transaction.UxBalance
				var err error
				if unsigned {
					txn, inputs, err = CreateTransaction(w, tc.params, addrUxOuts, tc.headTime)
				} else {
					txn, inputs, err = CreateTransactionSigned(w, tc.params, addrUxOuts, tc.headTime)
				}
				require.Equal(t, tc.err, err, "%v != %v", tc.err, err)
				if tc.err != nil {
					return
				}

				if unsigned {
					err := txn.VerifyUnsigned()
					require.NoError(t, err)
				} else {
					err := txn.Verify()
					require.NoError(t, err)
				}

				require.Equal(t, len(inputs), len(txn.In))

				// Checks duplicate inputs in array
				inputsMap := make(map[cipher.SHA256]struct{})
				for _, i := range inputs {
					_, ok := inputsMap[i.Hash]
					require.False(t, ok)
					inputsMap[i.Hash] = struct{}{}
				}

				for i, inUxid := range txn.In {
					_, ok := unspents[inUxid]
					require.True(t, ok)

					require.Equal(t, inUxid, inputs[i].Hash)
				}

				// Compare the transaction inputs
				chosenUnspents := make([]coin.UxOut, len(tc.chosenUnspents))
				chosenUnspentHashes := make([]cipher.SHA256, len(tc.chosenUnspents))
				for i, u := range tc.chosenUnspents {
					chosenUnspents[i] = u
					chosenUnspentHashes[i] = u.Hash()
				}
				sort.Slice(chosenUnspentHashes, func(i, j int) bool {
					return bytes.Compare(chosenUnspentHashes[i][:], chosenUnspentHashes[j][:]) < 0
				})
				sort.Slice(chosenUnspents, func(i, j int) bool {
					h1 := chosenUnspents[i].Hash()
					h2 := chosenUnspents[j].Hash()
					return bytes.Compare(h1[:], h2[:]) < 0
				})

				sortedTxnIn := make([]cipher.SHA256, len(txn.In))
				copy(sortedTxnIn[:], txn.In[:])

				sort.Slice(sortedTxnIn, func(i, j int) bool {
					return bytes.Compare(sortedTxnIn[i][:], sortedTxnIn[j][:]) < 0
				})

				require.Equal(t, chosenUnspentHashes, sortedTxnIn)

				sort.Slice(inputs, func(i, j int) bool {
					h1 := inputs[i].Hash
					h2 := inputs[j].Hash
					return bytes.Compare(h1[:], h2[:]) < 0
				})

				chosenUnspentsUxBalances := make([]transaction.UxBalance, len(chosenUnspents))
				for i, o := range chosenUnspents {
					b, err := transaction.NewUxBalance(tc.headTime, o)
					require.NoError(t, err)
					chosenUnspentsUxBalances[i] = b
				}

				require.Equal(t, chosenUnspentsUxBalances, inputs)

				// Assign expected hours for comparison
				var to []coin.TransactionOutput
				to = append(to, tc.params.To...)

				if len(tc.toExpectedHours) != 0 {
					require.Equal(t, len(tc.toExpectedHours), len(to))
					for i, h := range tc.toExpectedHours {
						to[i].Hours = h
					}
				}

				// Add the change output if specified
				if tc.changeOutput != nil {
					to = append(to, *tc.changeOutput)
				}

				// Compare transaction outputs
				require.Equal(t, to, txn.Out)

				// compare to previous result for the same test case
				// but either signed or unsigned (both should produce the same transactions, except for signatures)
				result := TestResult{
					Transaction: txn,
					Inputs:      inputs,
					Unsigned:    unsigned,
				}
				if prevResult, ok := results[tc.name]; ok {
					compareResult(t, prevResult, result)
				} else {
					results[tc.name] = result
				}
			})
		}
	}
}

func makeTransaction(t *testing.T, nInputs int) (coin.Transaction, []coin.UxOut, []cipher.SecKey) {
	txn := coin.Transaction{}

	toSign := make([]cipher.SecKey, 0)
	uxs := make([]coin.UxOut, 0)
	for i := 0; i < nInputs; i++ {
		ux, s := makeUxOutWithSecret(t)
		err := txn.PushInput(ux.Hash())
		require.NoError(t, err)
		toSign = append(toSign, s)
		uxs = append(uxs, ux)
	}

	err := txn.PushOutput(makeAddress(), 1e6, 50)
	require.NoError(t, err)
	err = txn.PushOutput(makeAddress(), 5e6, 50)
	require.NoError(t, err)
	txn.SignInputs(toSign)
	err = txn.UpdateHeader()
	require.NoError(t, err)

	return txn, uxs, toSign
}

func makeUxOut(t *testing.T, s cipher.SecKey, coins, hours uint64) coin.UxOut { //nolint:unparam
	body := makeUxBody(t, s, coins, hours)
	tm := rand.Int31n(1000)
	seq := rand.Int31n(100)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  uint64(tm),
			BkSeq: uint64(seq),
		},
		Body: body,
	}
}

func makeUxBody(t *testing.T, s cipher.SecKey, coins, hours uint64) coin.UxBody {
	p := cipher.MustPubKeyFromSecKey(s)
	return coin.UxBody{
		SrcTransaction: cipher.SumSHA256(testutil.RandBytes(t, 128)),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          coins,
		Hours:          hours,
	}
}

func makeUxOutWithSecret(t *testing.T) (coin.UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return coin.UxOut{
		Head: coin.UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func makeUxBodyWithSecret(t *testing.T) (coin.UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return coin.UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeAddress() cipher.Address {
	p, _ := cipher.GenerateKeyPair()
	return cipher.AddressFromPubKey(p)
}

func makeEntry() entry.Entry {
	p, s := cipher.GenerateKeyPair()
	a := cipher.AddressFromPubKey(p)
	return entry.Entry{
		Secret:  s,
		Public:  p,
		Address: a,
	}
}
