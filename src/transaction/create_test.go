package transaction

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/util/fee"
)

func TestCreate(t *testing.T) {
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

	validParams := Params{
		HoursSelection: HoursSelection{
			Type: HoursSelectionTypeManual,
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

	newShareFactor := func(a string) *decimal.Decimal {
		d, err := decimal.NewFromString(a)
		require.NoError(t, err)
		return &d
	}

	firstAddress := func(uxa coin.UxArray) cipher.Address {
		require.NotEmpty(t, uxa)

		addresses := make([]cipher.Address, len(uxa))
		for i, a := range uxa {
			addresses[i] = a.Body.Address
		}

		sort.Slice(addresses, func(i, j int) bool {
			x := addresses[i].Bytes()
			y := addresses[j].Bytes()
			return bytes.Compare(x, y) < 0
		})

		return addresses[0]
	}

	cases := []struct {
		name            string
		err             error
		params          Params
		unspents        []coin.UxOut
		addressUnspents coin.AddressUxOuts
		chosenUnspents  []coin.UxOut
		headTime        uint64
		changeOutput    *coin.TransactionOutput
		toExpectedHours []uint64
	}{
		{
			name:   "params invalid",
			params: Params{},
			err:    NewError(errors.New("To is required")),
		},

		{
			name: "overflowing coin hours in params",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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
			err: NewError(errors.New("total output hours error: uint64 addition overflow")),
		},

		{
			name: "overflowing coins in params",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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
			err: NewError(errors.New("total output coins error: uint64 addition overflow")),
		},

		{
			name: "no unspents",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   10,
						Coins:   1e6,
					},
				},
			},
			err: ErrNoUnspents,
		},

		{
			name: "insufficient coins",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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
			err:      ErrInsufficientBalance,
		},

		{
			name: "insufficient hours",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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
			err:      ErrInsufficientHours,
		},

		{
			name: "manual, 1 output, no change",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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

		// TODO -- belongs in visor_wallet_test.go
		// {
		// 	name: "manual, 1 output, no change, unknown address in auxs",
		// 	params: Params{
		// 		ChangeAddress: &changeAddress,
		// 		HoursSelection: HoursSelection{
		// 			Type: HoursSelectionTypeManual,
		// 		},
		// 		Wallet: CreateTransactionWalletParams{},
		// 		To: []coin.TransactionOutput{
		// 			{
		// 				Address: addrs[0],
		// 				Hours:   50,
		// 				Coins:   2e6,
		// 			},
		// 		},
		// 	},
		// 	addressUnspents: coin.AddressUxOuts{
		// 		testutil.MakeAddress(): []coin.UxOut{extraUxouts[0][0]},
		// 	},
		// 	err: ErrUnknownAddress,
		// },

		{
			name: "manual, 1 output, change",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   130,
						Coins:   2e6 + 1,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   50,
				Coins:   2e6 - 1,
			},
		},

		{
			name: "manual, 1 output, change, unspecified change address",
			params: Params{
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   2e6 + 1,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput: &coin.TransactionOutput{
				Address: firstAddress([]coin.UxOut{originalUxouts[0], originalUxouts[1]}),
				Hours:   130,
				Coins:   2e6 - 1,
			},
		},

		{
			// there are leftover coin hours and an additional input is added
			// to force change to save the leftover coin hours
			name: "manual, 1 output, forced change",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   0,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   272,
				Coins:   2e6,
			},
		},

		{
			// there are leftover coin hours and no coins change,
			// but there are no more unspents to use to force a change output
			name: "manual, 1 output, forced change rejected no more unspents",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   80,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       originalUxouts[:2],
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput:   nil,
		},

		{
			// there are leftover coin hours and no coins change,
			// but the hours cost of saving them with an additional input is less than is leftover
			name: "manual, 1 output, forced change rejected",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   175,
						Coins:   2e6 * 2,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1]},
			changeOutput:   nil,
		},

		{
			name: "manual, multiple outputs",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6 + 1,
					},
					{
						Address: addrs[1],
						Hours:   70,
						Coins:   2e6,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   102,
				Coins:   2e6 - 1,
			},
		},

		{
			name: "manual, multiple outputs, varied addressUnspents",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Hours:   50,
						Coins:   1e6 + 1,
					},
					{
						Address: addrs[1],
						Hours:   70,
						Coins:   2e6,
					},
				},
			},
			addressUnspents: coin.AddressUxOuts{
				extraWalletAddrs[0]: []coin.UxOut{extraUxouts[0][0]},
				extraWalletAddrs[3]: []coin.UxOut{extraUxouts[3][1], extraUxouts[3][2]},
				extraWalletAddrs[5]: []coin.UxOut{extraUxouts[5][6]},
			},
			chosenUnspents: []coin.UxOut{extraUxouts[0][0], extraUxouts[3][1], extraUxouts[3][2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   102,
				Coins:   2e6 - 1,
			},
		},

		{
			name: "auto, multiple outputs, share factor 0.5",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0.5"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   136,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{27, 54, 54, 1},
		},

		{
			name: "auto, multiple outputs, share factor 0.5, switch to 1.0 because no change could be made",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0.5"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e6 - 1e3,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:        []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			chosenUnspents:  []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			toExpectedHours: []uint64{46, 90, 90, 45, 1},
		},

		{
			name: "auto, multiple outputs, share factor 0",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("0"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   272,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{0, 0, 0, 0},
		},

		{
			name: "auto, multiple outputs, share factor 1",
			params: Params{
				ChangeAddress: &changeAddress,
				HoursSelection: HoursSelection{
					Type:        HoursSelectionTypeAuto,
					Mode:        HoursSelectionModeShare,
					ShareFactor: newShareFactor("1"),
				},
				To: []coin.TransactionOutput{
					{
						Address: addrs[0],
						Coins:   1e6,
					},
					{
						Address: addrs[0],
						Coins:   2e6,
					},
					{
						Address: addrs[1],
						Coins:   2e6,
					},
					{
						Address: addrs[4],
						Coins:   1e3,
					},
				},
			},
			unspents:       uxouts,
			chosenUnspents: []coin.UxOut{originalUxouts[0], originalUxouts[1], originalUxouts[2]},
			changeOutput: &coin.TransactionOutput{
				Address: changeAddress,
				Hours:   0,
				Coins:   2e6 - (1e6 + 1e3),
			},
			toExpectedHours: []uint64{55, 108, 108, 1},
		},

		{
			name:     "no coin hours in inputs",
			unspents: uxoutsNoHours[:],
			params: Params{
				HoursSelection: HoursSelection{
					Type: HoursSelectionTypeManual,
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

		{
			name:     "duplicate unspent output",
			unspents: append(uxouts, uxouts[:2]...),
			params:   validParams,
			err:      errors.New("Duplicate UxBalance in array"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
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

			t.Log("len of addrUxOuts:", len(addrUxOuts.Flatten()))
			txn, inputs, err := Create(tc.params, addrUxOuts, tc.headTime)
			require.Equal(t, tc.err, err, "%v != %v", tc.err, err)
			if tc.err != nil {
				return
			}

			err = txn.VerifyUnsigned()
			require.NoError(t, err)

			t.Log("len of txn.In:", len(txn.In))
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
				t.Log(u.Hash())
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

			t.Log(len(chosenUnspentHashes))
			t.Log(len(sortedTxnIn))
			t.Log(len(txn.In))
			require.Equal(t, chosenUnspentHashes, sortedTxnIn)

			sort.Slice(inputs, func(i, j int) bool {
				h1 := inputs[i].Hash
				h2 := inputs[j].Hash
				return bytes.Compare(h1[:], h2[:]) < 0
			})

			chosenUnspentsUxBalances := make([]UxBalance, len(chosenUnspents))
			for i, o := range chosenUnspents {
				b, err := NewUxBalance(tc.headTime, o)
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
		})
	}
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

func TestDecimalPtrEqual(t *testing.T) {
	// Sanity check for decimal comparison after converting to a pointer;
	// if this test fails then the Create method can go into an infinite loop
	oneDecimal := decimal.New(1, 0)
	oneDecimalPtr := &oneDecimal
	oneDecimal2 := decimal.New(1, 0)
	require.True(t, oneDecimalPtr.Equal(oneDecimal2))
}
