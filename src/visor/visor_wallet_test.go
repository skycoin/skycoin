package visor

import (
	"bytes"
	"errors"
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

func TestWalletCreateTransaction(t *testing.T) {
	// TODO
}

func TestCreateTransactionParamsValidate(t *testing.T) {
	var nullAddress cipher.Address
	addr := testutil.MakeAddress()
	hash := testutil.RandSHA256(t)

	cases := []struct {
		name string
		p    CreateTransactionParams
		err  error
	}{
		{
			name: "both addrs and uxouts specified",
			p: CreateTransactionParams{
				Addresses: []cipher.Address{addr},
				UxOuts:    []cipher.SHA256{hash},
			},
			err: ErrCreateTransactionParamsConflict,
		},

		{
			name: "null address in addrs",
			p: CreateTransactionParams{
				Addresses: []cipher.Address{nullAddress},
			},
			err: ErrIncludesNullAddress,
		},

		{
			name: "duplicate address in addrs",
			p: CreateTransactionParams{
				Addresses: []cipher.Address{addr, addr},
			},
			err: ErrDuplicateAddresses,
		},

		{
			name: "duplicate hash in uxouts",
			p: CreateTransactionParams{
				UxOuts: []cipher.SHA256{hash, hash},
			},
			err: ErrDuplicateUxOuts,
		},

		{
			name: "ok, addrs specified",
			p: CreateTransactionParams{
				Addresses: []cipher.Address{addr},
			},
		},

		{
			name: "ok, uxouts specified",
			p: CreateTransactionParams{
				UxOuts: []cipher.SHA256{hash},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.p.Validate()
			require.Equal(t, tc.err, err, "%v != %v", tc.err, err)
		})
	}
}

func TestWalletCreateTransactionValidation(t *testing.T) {
	// This only tests that WalletCreateTransaction and WalletCreateTransactionSigned fails on invalid inputs;
	// success tests are performed by live integration tests

	validParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   10,
				Hours:   10,
			},
		},
	}

	cases := []struct {
		name string
		p    transaction.Params
		wp   CreateTransactionParams
		err  error
	}{
		{
			name: "bad transaction.Params",
			p:    transaction.Params{},
			err:  transaction.ErrMissingReceivers,
		},
		{
			name: "bad CreateTransactionParams",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: []cipher.Address{testutil.MakeAddress()},
				UxOuts:    []cipher.SHA256{testutil.RandSHA256(t)},
			},
			err: ErrCreateTransactionParamsConflict,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// setup visor
			v := &Visor{}

			_, _, err := v.WalletCreateTransaction("foo.wlt", tc.p, tc.wp)
			require.Equal(t, tc.err, err)

			_, _, err = v.WalletCreateTransactionSigned("foo.wlt", nil, tc.p, tc.wp)
			require.Equal(t, tc.err, err)

			if tc.err != nil {
				return
			}

			// Valid WalletCreateTransaction and WalletCreateTransactionSigned calls are tested in live integration tests
		})
	}
}

func TestCreateTransactionValidation(t *testing.T) {
	// This only tests that CreateTransaction fails on invalid inputs;
	// success tests are performed by live integration tests

	validParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   10,
				Hours:   10,
			},
		},
	}

	cases := []struct {
		name string
		p    transaction.Params
		wp   CreateTransactionParams
		err  error
	}{
		{
			name: "bad transaction.Params",
			p:    transaction.Params{},
			err:  transaction.ErrMissingReceivers,
		},
		{
			name: "bad CreateTransactionParams",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: []cipher.Address{testutil.MakeAddress()},
				UxOuts:    []cipher.SHA256{testutil.RandSHA256(t)},
			},
			err: ErrCreateTransactionParamsConflict,
		},
		{
			name: "Addresses and UxOuts both empty",
			p:    validParams,
			err:  errors.New("UxOuts or Addresses must not be empty"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// setup visor
			v := &Visor{}

			_, _, err := v.CreateTransaction(tc.p, tc.wp)
			require.Equal(t, tc.err, err)

			if tc.err != nil {
				return
			}

			// Valid CreateTransaction calls are tested in live integration tests
		})
	}
}

func TestGetCreateTransactionAuxsUxOut(t *testing.T) {
	allAddrs := make([]cipher.Address, 10)
	for i := range allAddrs {
		allAddrs[i] = testutil.MakeAddress()
	}

	hashes := make([]cipher.SHA256, 20)
	for i := range hashes {
		hashes[i] = testutil.RandSHA256(t)
	}

	srcTxns := make([]cipher.SHA256, 20)
	for i := range srcTxns {
		srcTxns[i] = testutil.RandSHA256(t)
	}

	cases := []struct {
		name              string
		ignoreUnconfirmed bool
		uxOuts            []cipher.SHA256
		expectedAuxs      coin.AddressUxOuts
		err               error

		forEachErr      error
		unconfirmedTxns coin.Transactions
		getArrayInputs  []cipher.SHA256
		getArrayRet     coin.UxArray
		getArrayErr     error
	}{
		{
			name:   "uxouts specified, ok",
			uxOuts: hashes[5:10],
			unconfirmedTxns: coin.Transactions{
				coin.Transaction{
					In: hashes[0:2],
				},
				coin.Transaction{
					In: hashes[2:4],
				},
			},
			getArrayInputs: hashes[5:10],
			getArrayRet: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[6],
						Address:        allAddrs[3],
					},
				},
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[6],
							Address:        allAddrs[3],
						},
					},
				},
			},
		},

		{
			name:       "uxouts specified, unconfirmed spend",
			uxOuts:     hashes[0:4],
			err:        ErrSpendingUnconfirmed,
			forEachErr: ErrSpendingUnconfirmed,
			unconfirmedTxns: coin.Transactions{
				coin.Transaction{
					In: hashes[6:10],
				},
				coin.Transaction{
					In: hashes[3:6],
				},
			},
		},

		{
			name:              "uxouts specified, unconfirmed spend ignored",
			ignoreUnconfirmed: true,
			uxOuts:            hashes[5:10],
			unconfirmedTxns: coin.Transactions{
				coin.Transaction{
					In: hashes[0:2],
				},
				coin.Transaction{
					In: hashes[2:4],
				},
				coin.Transaction{
					In: hashes[8:10],
				},
			},
			getArrayInputs: hashes[5:8], // the 8th & 9th hash are filtered because it is an unconfirmed spend
			getArrayRet: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
			},
		},

		{
			name:   "uxouts specified, unknown uxout",
			uxOuts: hashes[5:10],
			err: blockdb.ErrUnspentNotExist{
				UxID: "foo",
			},
			getArrayErr: blockdb.ErrUnspentNotExist{
				UxID: "foo",
			},
			unconfirmedTxns: coin.Transactions{
				coin.Transaction{
					In: hashes[0:2],
				},
				coin.Transaction{
					In: hashes[2:4],
				},
			},
			getArrayInputs: hashes[5:10],
			getArrayRet: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[4],
						Address:        testutil.MakeAddress(),
					},
				},
			},
		},
	}

	matchDBTx := mock.MatchedBy(func(tx *dbutil.Tx) bool {
		return true
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := testutil.PrepareDB(t)
			defer shutdown()

			unconfirmed := &MockUnconfirmedTransactionPooler{}
			bc := &MockBlockchainer{}
			unspent := &MockUnspentPooler{}
			require.Implements(t, (*blockdb.UnspentPooler)(nil), unspent)

			v := &Visor{
				Unconfirmed: unconfirmed,
				Blockchain:  bc,
				DB:          db,
			}

			unconfirmed.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(func(args mock.Arguments) {
				// Simulate the ForEach callback method, unless ForEach was configured to return an error
				fn := args.Get(1).(func(cipher.SHA256, UnconfirmedTransaction) error)
				for _, u := range tc.unconfirmedTxns {
					err := fn(u.Hash(), UnconfirmedTransaction{
						Transaction: u,
					})

					// If any of the input hashes are in an unconfirmed transaction,
					// the callback handler should have returned ErrSpendingUnconfirmed
					// unless IgnoreUnconfirmed is true
					hasUnconfirmedHash := hashesIntersect(u.In, tc.uxOuts)

					if hasUnconfirmedHash {
						if tc.ignoreUnconfirmed {
							require.NoError(t, err)
						} else {
							require.Equal(t, ErrSpendingUnconfirmed, err)
						}
					} else {
						require.NoError(t, err)
					}
				}
			})

			unspent.On("GetArray", matchDBTx, mock.MatchedBy(func(args []cipher.SHA256) bool {
				// Compares two []coin.UxOuts for equality, ignoring the order of elements in the slice
				if len(args) != len(tc.getArrayInputs) {
					return false
				}

				x := make([]cipher.SHA256, len(tc.getArrayInputs))
				copy(x, tc.getArrayInputs)
				y := make([]cipher.SHA256, len(args))
				copy(y, args)

				sort.Slice(x, func(a, b int) bool {
					return bytes.Compare(x[a][:], x[b][:]) < 0
				})
				sort.Slice(y, func(a, b int) bool {
					return bytes.Compare(y[a][:], y[b][:]) < 0
				})

				return reflect.DeepEqual(x, y)
			})).Return(tc.getArrayRet, tc.getArrayErr)

			bc.On("Unspent").Return(unspent)

			var auxs coin.AddressUxOuts
			err := v.DB.View("", func(tx *dbutil.Tx) error {
				var err error
				auxs, err = v.getCreateTransactionAuxsUxOut(tx, tc.uxOuts, tc.ignoreUnconfirmed)
				return err
			})

			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.expectedAuxs, auxs)
		})
	}
}

func TestGetCreateTransactionAuxsAddress(t *testing.T) {
	allAddrs := make([]cipher.Address, 10)
	for i := range allAddrs {
		allAddrs[i] = testutil.MakeAddress()
	}

	hashes := make([]cipher.SHA256, 20)
	for i := range hashes {
		hashes[i] = testutil.RandSHA256(t)
	}

	srcTxns := make([]cipher.SHA256, 20)
	for i := range srcTxns {
		srcTxns[i] = testutil.RandSHA256(t)
	}

	cases := []struct {
		name              string
		ignoreUnconfirmed bool
		addrs             []cipher.Address
		expectedAuxs      coin.AddressUxOuts
		err               error

		forEachErr                 error
		unconfirmedTxns            coin.Transactions
		getArrayInputs             []cipher.SHA256
		getArrayRet                coin.UxArray
		getArrayErr                error
		getUnspentHashesOfAddrsRet blockdb.AddressHashes
	}{
		{
			name:           "ok",
			addrs:          allAddrs,
			getArrayInputs: hashes[0:4],
			getArrayRet: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[6],
						Address:        allAddrs[3],
					},
				},
			},
			getUnspentHashesOfAddrsRet: blockdb.AddressHashes{
				allAddrs[1]: hashes[0:2],
				allAddrs[3]: hashes[2:4],
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[6],
							Address:        allAddrs[3],
						},
					},
				},
			},
		},

		{
			name:       "err, unconfirmed spends",
			addrs:      allAddrs,
			err:        ErrSpendingUnconfirmed,
			forEachErr: ErrSpendingUnconfirmed,
			getUnspentHashesOfAddrsRet: blockdb.AddressHashes{
				allAddrs[1]: hashes[0:2],
				allAddrs[3]: hashes[2:4],
			},
		},

		{
			name:              "ignore unconfirmed",
			ignoreUnconfirmed: true,
			addrs:             allAddrs,
			unconfirmedTxns: coin.Transactions{
				{
					In: []cipher.SHA256{hashes[1]},
				},
				{
					In: []cipher.SHA256{hashes[2]},
				},
			},
			getArrayInputs: []cipher.SHA256{hashes[0], hashes[3]},
			getArrayRet: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[6],
						Address:        allAddrs[3],
					},
				},
			},
			getUnspentHashesOfAddrsRet: blockdb.AddressHashes{
				allAddrs[1]: hashes[0:2],
				allAddrs[3]: hashes[2:4],
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					coin.UxOut{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[6],
							Address:        allAddrs[3],
						},
					},
				},
			},
		},
	}

	matchDBTx := mock.MatchedBy(func(tx *dbutil.Tx) bool {
		return true
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := testutil.PrepareDB(t)
			defer shutdown()

			unconfirmed := &MockUnconfirmedTransactionPooler{}
			bc := &MockBlockchainer{}
			unspent := &MockUnspentPooler{}
			require.Implements(t, (*blockdb.UnspentPooler)(nil), unspent)

			v := &Visor{
				Unconfirmed: unconfirmed,
				Blockchain:  bc,
				DB:          db,
			}

			unconfirmed.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(func(args mock.Arguments) {
				// Simulate the ForEach callback method, unless ForEach was configured to return an error
				if tc.forEachErr != nil {
					return
				}
				fn := args.Get(1).(func(cipher.SHA256, UnconfirmedTransaction) error)
				for _, u := range tc.unconfirmedTxns {
					err := fn(u.Hash(), UnconfirmedTransaction{
						Transaction: u,
					})

					// If any of the input hashes are in an unconfirmed transaction,
					// the callback handler should have returned ErrSpendingUnconfirmed
					// unless IgnoreUnconfirmed is true
					hasUnconfirmedHash := hashesIntersect(u.In, tc.getUnspentHashesOfAddrsRet.Flatten())

					if hasUnconfirmedHash {
						if tc.ignoreUnconfirmed {
							require.NoError(t, err)
						} else {
							require.Equal(t, ErrSpendingUnconfirmed, err)
						}
					} else {
						require.NoError(t, err)
					}
				}
			})

			unspent.On("GetArray", matchDBTx, mock.MatchedBy(func(args []cipher.SHA256) bool {
				// Compares two []coin.UxOuts for equality, ignoring the order of elements in the slice
				if len(args) != len(tc.getArrayInputs) {
					return false
				}

				x := make([]cipher.SHA256, len(tc.getArrayInputs))
				copy(x, tc.getArrayInputs)
				y := make([]cipher.SHA256, len(args))
				copy(y, args)

				sort.Slice(x, func(a, b int) bool {
					return bytes.Compare(x[a][:], x[b][:]) < 0
				})
				sort.Slice(y, func(a, b int) bool {
					return bytes.Compare(y[a][:], y[b][:]) < 0
				})

				return reflect.DeepEqual(x, y)
			})).Return(tc.getArrayRet, tc.getArrayErr)

			unspent.On("GetUnspentHashesOfAddrs", matchDBTx, tc.addrs).Return(tc.getUnspentHashesOfAddrsRet, nil)

			bc.On("Unspent").Return(unspent)

			var auxs coin.AddressUxOuts
			err := v.DB.View("", func(tx *dbutil.Tx) error {
				var err error
				auxs, err = v.getCreateTransactionAuxsAddress(tx, tc.addrs, tc.ignoreUnconfirmed)
				return err
			})

			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.expectedAuxs, auxs)
		})
	}
}

// hashesIntersect returns true if there are any hashes common to x and y
func hashesIntersect(x, y []cipher.SHA256) bool {
	for _, a := range x {
		for _, b := range y {
			if a == b {
				return true
			}
		}
	}
	return false
}
