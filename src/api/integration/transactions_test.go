package integration_test

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/wallet"
)

func TestStableInjectTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	cases := []struct {
		name string
		txn  coin.Transaction
		code int
		err  string
	}{
		{
			name: "database is read only",
			txn:  coin.Transaction{},
			code: http.StatusInternalServerError,
			err:  "500 Internal Server Error - database is in read-only mode",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := c.InjectTransaction(&tc.txn)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)

			// Result should be a valid txid
			require.NotEmpty(t, result)
			h, err := cipher.SHA256FromHex(result)
			require.NoError(t, err)
			require.NotEqual(t, cipher.SHA256{}, h)
		})
	}
}

func TestLiveInjectTransactionDisableNetworking(t *testing.T) {
	if !doLive(t) {
		return
	}

	if !liveDisableNetworking(t) {
		t.Skip("Networking must be disabled for this test")
		return
	}

	requireWalletEnv(t)

	c := newClient()

	w, totalCoins, totalHours, password := prepareAndCheckWallet(t, c, 2e6, 20)

	defaultChangeAddress := w.GetEntryAt(0).Address.String()

	type testCase struct {
		name         string
		createTxnReq api.WalletCreateTransactionRequest
		err          string
		code         int
	}

	cases := []testCase{
		{
			name: "valid request, networking disabled",
			err:  "503 Service Unavailable - Outgoing connections are disabled",
			code: http.StatusServiceUnavailable,
			createTxnReq: api.WalletCreateTransactionRequest{
				WalletID: w.Filename(),
				Password: password,
				CreateTransactionRequest: api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(1).Address.String(),
							Coins:   toDropletString(t, totalCoins),
							Hours:   fmt.Sprint(totalHours / 2),
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			txnResp, err := c.WalletCreateTransaction(tc.createTxnReq)
			require.NoError(t, err)

			txid, err := c.InjectEncodedTransaction(txnResp.EncodedTransaction)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)

				// A second injection will fail with the same error,
				// since the transaction should not be saved to the DB
				_, err = c.InjectEncodedTransaction(txnResp.EncodedTransaction)
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NotEmpty(t, txid)
			require.Equal(t, txnResp.Transaction.TxID, txid)

			h, err := cipher.SHA256FromHex(txid)
			require.NoError(t, err)
			require.NotEqual(t, cipher.SHA256{}, h)
		})
	}
}

func TestLiveInjectTransactionEnableNetworking(t *testing.T) {
	if !doLive(t) {
		return
	}

	if liveDisableNetworking(t) {
		t.Skip("This tests requires networking enabled")
		return
	}

	requireWalletEnv(t)

	c := newClient()
	w, totalCoins, _, password := prepareAndCheckWallet(t, c, 2e6, 2)

	defaultChangeAddress := w.GetEntryAt(0).Address.String()

	// prepareTxnFunc prepares a valid transaction
	prepareTxnFunc := func(t *testing.T, toAddr string, coins uint64, shareFactor string) (coin.Transaction, *api.CreateTransactionResponse) {
		createTxnReq := api.WalletCreateTransactionRequest{
			WalletID: w.Filename(),
			Password: password,
			CreateTransactionRequest: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: shareFactor,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						// Address: w.GetEntryAt(1).Address.String(),
						Address: toAddr,
						Coins:   toDropletString(t, coins),
					},
				},
			},
		}

		txnResp, err := c.WalletCreateTransaction(createTxnReq)
		require.NoError(t, err)

		txn, err := coin.DeserializeTransactionHex(txnResp.EncodedTransaction)
		require.NoError(t, err)
		return txn, txnResp
	}

	reSignTxnFunc := func(t *testing.T, txn coin.Transaction, txnRsp *api.CreateTransactionResponse, wlt wallet.Wallet) coin.Transaction {
		walletPassword := os.Getenv("WALLET_PASSWORD")
		err := wallet.GuardView(wlt, []byte(walletPassword), func(unlockWlt wallet.Wallet) error {
			keyMap := make(map[string]cipher.SecKey, unlockWlt.EntriesLen())
			for _, e := range unlockWlt.GetEntries() {
				addr := cipher.MustAddressFromSecKey(e.Secret)
				keyMap[addr.String()] = e.Secret
			}

			// Get seckeys in wallet of input addresses
			keys := make([]cipher.SecKey, len(txnRsp.Transaction.In))
			for i, in := range txnRsp.Transaction.In {
				k, ok := keyMap[in.Address]
				if !ok {
					t.Fatal("seckey does not exist")
					return errors.New("seckey does not exist")
				}
				keys[i] = k
			}
			// clear the old signatures
			txn.Sigs = []cipher.Sig{}
			txn.SignInputs(keys)
			return nil
		})
		require.NoError(t, err)
		return txn
	}

	tt := []struct {
		name      string
		createTxn func(t *testing.T) *coin.Transaction
		code      int
		err       string
		checkTxn  func(t *testing.T, tx *readable.TransactionWithStatus)
	}{
		{
			name: "send all coins to the first address",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				return &txn
			},
			checkTxn: func(t *testing.T, tx *readable.TransactionWithStatus) {
				// Confirms the total output coins are equal to the totalCoins
				var coins uint64
				for _, o := range tx.Transaction.Out {
					c, err := droplet.FromString(o.Coins)
					require.NoError(t, err)
					coins, err = mathutil.AddUint64(coins, c)
					require.NoError(t, err)
				}

				// Confirms the address balance are equal to the totalCoins
				coins, _ = getAddressBalance(t, c, w.GetEntryAt(0).Address.String())
				require.Equal(t, totalCoins, coins)
			},
			code: http.StatusOK,
		},
		{
			// send 0.003 coin to the second address,
			// this amount is chosen to not interfere with TestLiveWalletCreateTransaction
			name: "send 0.003 coin to second address",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(1).Address.String(), 3e3, "0.5")
				return &txn
			},
			checkTxn: func(t *testing.T, tx *readable.TransactionWithStatus) {
				// Confirms there're two outputs, one to the second address, one as change output to the first address.
				require.Len(t, tx.Transaction.Out, 2)

				// Gets the output of the second address in the transaction
				getAddrOutputInTxn := func(t *testing.T, tx *readable.TransactionWithStatus, addr string) *readable.TransactionOutput {
					for _, output := range tx.Transaction.Out {
						if output.Address == addr {
							return &output
						}
					}
					t.Fatalf("transaction doesn't have output to address: %v", addr)
					return nil
				}

				out := getAddrOutputInTxn(t, tx, w.GetEntryAt(1).Address.String())

				// Confirms the second address has 0.003 coin
				require.Equal(t, out.Coins, "0.003000")
				require.Equal(t, out.Address, w.GetEntryAt(1).Address.String())

				coin, err := droplet.FromString(out.Coins)
				require.NoError(t, err)

				// Gets the expected change coins
				expectChangeCoins := totalCoins - coin

				// Gets the real change coins
				changeOut := getAddrOutputInTxn(t, tx, w.GetEntryAt(0).Address.String())
				changeCoins, err := droplet.FromString(changeOut.Coins)
				require.NoError(t, err)
				// Confirms the change coins are matched.
				require.Equal(t, expectChangeCoins, changeCoins)
			},
			code: http.StatusOK,
		},
		{
			name: "send to null address",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")

				// set the transaction output address as null
				txn.Out[0].Address = cipher.Address{}
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates user constraint: Transaction output is sent to the null address",
		},
		{
			// Use an input from block 1024: 2f842b0fbf5ef2dd59c8b5127795f1e88bfa6b510a41c62eac28fc2006d279e3
			name: "double spend",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")

				hash, err := cipher.SHA256FromHex("2f842b0fbf5ef2dd59c8b5127795f1e88bfa6b510a41c62eac28fc2006d279e3")
				require.NoError(t, err)
				txn.In[0] = hash
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: unspent output of 2f842b0fbf5ef2dd59c8b5127795f1e88bfa6b510a41c62eac28fc2006d279e3 does not exist",
		},
		{
			name: "output hours overflow",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(1).Address.String(), 1e6, "1")

				// set one output hours as math.MaxUint64
				txn.Out[0].Hours = math.MaxUint64 - 1
				txn.Out[1].Hours = 100
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Transaction output hours overflow",
		},
		{
			name: "no inputs",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")

				txn.In = []cipher.SHA256{}
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: No inputs",
		},
		{
			name: "no outputs",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")

				txn.Out = []coin.TransactionOutput{}
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: No outputs",
		},
		{
			name: "invalid number of signatures",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Sigs = []cipher.Sig{}
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Invalid number of signatures",
		},
		{
			name: "duplicate spend",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				// Make duplicate inputs
				txn.In = append(txn.In, txn.In[0])
				// Make duplicate sigs
				txn.Sigs = append(txn.Sigs, txn.Sigs[0])
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Duplicate spend",
		},
		{
			name: "transaction type invalid",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Type = 1
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: transaction type invalid",
		},
		{
			name: "zero coin output",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Out[0].Coins = 0
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Zero coin output",
		},
		{
			name: "output coins overflow",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), 1e6, "1")
				txn.Out[0].Coins = math.MaxUint64 - 1
				txn.Out[1].Coins = 2
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Output coins overflow",
		},
		{
			name: "incorrect transaction length",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Length = 1
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Incorrect transaction length",
		},
		{
			name: "duplicate output",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Out = append(txn.Out, txn.Out[0])
				txn.InnerHash = txn.HashInner()
				size, _, err := txn.SizeHash()
				require.NoError(t, err)
				txn.Length = size
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Duplicate output in transaction",
		},
		{
			name: "inner hash does not match",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.InnerHash = testutil.RandSHA256(t)
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: InnerHash does not match computed hash",
		},
		{
			name: "unsigned input",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Sigs[0] = cipher.Sig{}
				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Unsigned input in transaction",
		},
		{
			name: "invalid sig",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				txn.Sigs[0] = testutil.RandSig(t)

				txn.InnerHash = txn.HashInner()
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Failed to recover pubkey from signature",
		},
		{
			name: "signature not valid for output being spent",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, _ := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				// Use a wrong private key to sign txn.In[0] and change txn.Sigs[0]
				_, seckey := cipher.GenerateKeyPair()
				h := cipher.AddSHA256(txn.InnerHash, txn.In[0])
				txn.Sigs[0] = cipher.MustSignHash(h, seckey)

				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Signature not valid for output being spent",
		},
		{
			name: "insufficient coins",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, txnRsp := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				// Make output coins > input coins
				txn.Out[0].Coins = txn.Out[0].Coins + 1
				txn.InnerHash = txn.HashInner()
				// Sign txn again as the inner hash is changed
				txn = reSignTxnFunc(t, txn, txnRsp, w)

				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Insufficient coins",
		},
		{
			name: "transaction may not destry coins",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, txnRsp := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				// Make output coins < input coins
				txn.Out[0].Coins = txn.Out[0].Coins - 1
				txn.InnerHash = txn.HashInner()
				// Sign the txn again as the inner hash is changed
				txn = reSignTxnFunc(t, txn, txnRsp, w)
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Transactions may not destroy coins",
		},
		{
			name: "insufficient coin hours",
			createTxn: func(t *testing.T) *coin.Transaction {
				txn, txnRsp := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), totalCoins, "1")
				// Make up more output coin hours
				txn.Out[0].Hours = txn.Out[0].Hours + 1e6
				// Recalculate inner hash
				txn.InnerHash = txn.HashInner()
				// Sign the txn again
				txn = reSignTxnFunc(t, txn, txnRsp, w)
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates hard constraint: Insufficient coin hours",
		},
		{
			name: "invalid amount, too many decimal places",
			createTxn: func(t *testing.T) *coin.Transaction {
				// Make a txn with txn.Out[0].Coins equal 1e3, as we have at least 2e6 coins
				// so there will have at least two outputs.
				txn, txnRsp := prepareTxnFunc(t, w.GetEntryAt(0).Address.String(), 1e3, "1")
				// Make txn.Out[0].Coins too many decimal places
				txn.Out[0].Coins = 5e2
				// Move the remaining 5e2 from the first output to the second output, so that
				// we won't lose coins
				txn.Out[1].Coins = txn.Out[1].Coins + 5e2
				txn.InnerHash = txn.HashInner()

				txn = reSignTxnFunc(t, txn, txnRsp, w)
				return &txn
			},
			code: http.StatusBadRequest,
			err:  "400 Bad Request - Transaction violates soft constraint: invalid amount, too many decimal places",
		},
	}
	// TODO:
	// The following test cases can not be added here, as we cannot use the big size
	// transaction is not allowed in encoder/decoder, which will fail the test
	// in building transaction step.
	//
	// 1. Make up a txn which exceeds max block size to violate the soft constraint
	// 2. Make up a transaction that can exceed max block size
	// 		expected err: "400 Bad Request - Transaction violates hard constraint: Transaction size bigger than max block size"
	// 3. Make up a transaction that has inputs/outputs exceed max
	//		expected err:
	// 		- "400 Bad Request - Transaction violates hard constraint: Too many signatures and inputs"
	// 		- "400 Bad Request - Transaction violates hard constraint: Too many outputs"
	// TODO:
	// 1. Add test case to inject transaction who has inputs locked.

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			txn := tc.createTxn(t)
			txid, err := c.InjectTransaction(txn)
			if tc.code != http.StatusOK {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, txn.Hash().Hex(), txid)

			tk := time.NewTicker(time.Second)
			var txnStatus *readable.TransactionWithStatus
		loop:
			for {
				select {
				case <-time.After(30 * time.Second):
					t.Fatal("Waiting for transaction to be confirmed timeout")
				case <-tk.C:
					txnStatus = getTransaction(t, c, txn.Hash().Hex())
					if txnStatus.Status.Confirmed {
						break loop
					}
				}
			}
			tc.checkTxn(t, txnStatus)
		})
	}

	// Test to inject invalid rawtx
	_, err := c.InjectEncodedTransaction("invalidrawtx")
	assertResponseError(t, err, 400, "400 Bad Request - Transaction violates user constraint: Transaction output is sent to the null address")

}

func TestLiveWalletSignTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()

	w, _, _, password := prepareAndCheckWallet(t, c, 2e6, 20)

	// Fetch outputs held by the wallet
	addrs := make([]string, w.EntriesLen())
	for i, e := range w.GetEntries() {
		addrs[i] = e.SkycoinAddress().String()
	}

	summary, err := c.OutputsForAddresses(addrs)
	require.NoError(t, err)
	// Abort if the transaction is spending summary
	require.Empty(t, summary.OutgoingOutputs)
	// Need at least 2 summary for the created transaction
	require.True(t, len(summary.HeadOutputs) > 1)

	// Use the first two outputs for a transaction
	headOutputs := summary.HeadOutputs[:2]
	outputs, err := headOutputs.ToUxArray()
	require.NoError(t, err)
	totalCoins, err := outputs.Coins()
	require.NoError(t, err)
	totalCoinsStr, err := droplet.ToString(totalCoins)
	require.NoError(t, err)

	uxOutHashes := make([]string, len(outputs))
	for i, o := range outputs {
		uxOutHashes[i] = o.Hash().Hex()
	}

	// Create an unsigned transaction using two inputs
	// Ensure at least 2 inputs
	// Specify outputs in the request to create txn
	// Specify unsigned in the request to create txn
	txnResp, err := c.WalletCreateTransaction(api.WalletCreateTransactionRequest{
		Unsigned: true,
		WalletID: w.Filename(),
		Password: password,
		CreateTransactionRequest: api.CreateTransactionRequest{
			UxOuts: uxOutHashes,
			HoursSelection: api.HoursSelection{
				Type:        transaction.HoursSelectionTypeAuto,
				Mode:        transaction.HoursSelectionModeShare,
				ShareFactor: "0.5",
			},
			To: []api.Receiver{
				{
					Address: w.GetEntryAt(0).SkycoinAddress().String(),
					Coins:   totalCoinsStr,
				},
			},
		},
	})
	require.NoError(t, err)

	// Create an invalid txn with an extra null sig
	invalidTxn := coin.MustDeserializeTransactionHex(txnResp.EncodedTransaction)
	invalidTxn.Sigs = append(invalidTxn.Sigs, cipher.Sig{})
	require.NotEqual(t, len(invalidTxn.In), len(invalidTxn.Sigs))

	type testCase struct {
		name        string
		req         api.WalletSignTransactionRequest
		fullySigned bool
		err         string
		code        int
	}

	cases := []testCase{
		{
			name: "sign one input",
			req: api.WalletSignTransactionRequest{
				WalletID:           w.Filename(),
				Password:           password,
				SignIndexes:        []int{1},
				EncodedTransaction: txnResp.EncodedTransaction,
			},
			fullySigned: false,
		},

		{
			name: "sign all inputs",
			req: api.WalletSignTransactionRequest{
				WalletID:           w.Filename(),
				Password:           password,
				SignIndexes:        nil,
				EncodedTransaction: txnResp.EncodedTransaction,
			},
			fullySigned: true,
		},

		{
			name: "sign invalid txn",
			req: api.WalletSignTransactionRequest{
				WalletID:           w.Filename(),
				Password:           password,
				SignIndexes:        nil,
				EncodedTransaction: invalidTxn.MustSerializeHex(),
			},
			code: http.StatusBadRequest,
			err:  "Transaction violates hard constraint: Invalid number of signatures",
		},
	}

	doTest := func(tc testCase) {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := c.WalletSignTransaction(tc.req)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)

			txn, err := coin.DeserializeTransactionHex(tc.req.EncodedTransaction)
			require.NoError(t, err)

			// TxID should have changed
			require.NotEqual(t, txn.Hash(), resp.Transaction.TxID)
			// Length, InnerHash should not have changed
			require.Equal(t, txn.Length, resp.Transaction.Length)
			require.Equal(t, txn.InnerHash.Hex(), resp.Transaction.InnerHash)

			_, err = c.VerifyTransaction(api.VerifyTransactionRequest{
				EncodedTransaction: resp.EncodedTransaction,
				Unsigned:           false,
			})
			if tc.fullySigned {
				require.NoError(t, err)
			} else {
				testutil.RequireError(t, err, "Transaction violates hard constraint: Unsigned input in transaction")
			}

			_, err = c.VerifyTransaction(api.VerifyTransactionRequest{
				EncodedTransaction: resp.EncodedTransaction,
				Unsigned:           true,
			})
			if tc.fullySigned {
				testutil.RequireError(t, err, "Transaction violates hard constraint: Unsigned transaction must contain a null signature")
			} else {
				require.NoError(t, err)
			}
		})
	}

	for _, tc := range cases {
		doTest(tc)
	}

	// Create a partially signed transaction then sign the remainder of it
	resp, err := c.WalletSignTransaction(api.WalletSignTransactionRequest{
		WalletID:           w.Filename(),
		Password:           password,
		SignIndexes:        []int{1},
		EncodedTransaction: txnResp.EncodedTransaction,
	})
	require.NoError(t, err)

	doTest(testCase{
		name: "sign partially signed transaction",
		req: api.WalletSignTransactionRequest{
			WalletID:           w.Filename(),
			Password:           password,
			EncodedTransaction: resp.EncodedTransaction,
		},
		fullySigned: true,
	})
}

func toDropletString(t *testing.T, i uint64) string {
	x, err := droplet.ToString(i)
	require.NoError(t, err)
	return x
}

func TestStableCreateTransaction(t *testing.T) {
	if !doStable(t) {
		return
	}

	c := newClient()

	unknownHash := testutil.RandSHA256(t).Hex()

	cases := []struct {
		name string
		req  api.CreateTransactionRequest
		err  string
		code int
	}{
		{
			name: "invalid no uxouts for addresses",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses: []string{testutil.MakeAddress().String()},
				To: []api.Receiver{
					{
						Address: testutil.MakeAddress().String(),
						Coins:   "1.000000",
						Hours:   "100",
					},
				},
			},
			code: http.StatusBadRequest,
			err:  "no unspents to spend",
		},

		{
			name: "invalid uxouts do not exist",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				UxOuts: []string{unknownHash},
				To: []api.Receiver{
					{
						Address: testutil.MakeAddress().String(),
						Coins:   "1.000000",
						Hours:   "100",
					},
				},
			},
			code: http.StatusBadRequest,
			err:  fmt.Sprintf("unspent output of %s does not exist", unknownHash),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotEmpty(t, tc.err)
			_, err := c.CreateTransaction(tc.req)
			assertResponseError(t, err, tc.code, tc.err)
		})
	}
}

type liveCreateTxnTestCase struct {
	name                 string
	req                  api.CreateTransactionRequest
	outputs              []coin.TransactionOutput
	outputsSubset        []coin.TransactionOutput
	err                  string
	code                 int
	ignoreHours          bool
	additionalRespVerify func(t *testing.T, r *api.CreateTransactionResponse)
}

func makeLiveCreateTxnTestCases(t *testing.T, w wallet.Wallet, totalCoins, totalHours uint64) []liveCreateTxnTestCase {
	remainingHours := fee.RemainingHours(totalHours, params.UserVerifyTxn.BurnFactor)
	require.True(t, remainingHours > 1)
	unknownOutput := testutil.RandSHA256(t)
	defaultChangeAddress := w.GetEntryAt(0).Address.String()

	// Get all outputs
	c := newClient()
	outputs, err := c.Outputs()
	require.NoError(t, err)

	// Split outputs into those held by the wallet and those not
	addresses := make([]string, w.EntriesLen())
	addressMap := make(map[string]struct{}, w.EntriesLen())
	for i, e := range w.GetEntries() {
		addresses[i] = e.Address.String()
		addressMap[e.Address.String()] = struct{}{}
	}

	var walletOutputHashes []string
	var walletOutputs readable.UnspentOutputs
	walletAuxs := make(map[string][]string)
	for _, o := range outputs.HeadOutputs {
		if _, ok := addressMap[o.Address]; ok {
			walletOutputs = append(walletOutputs, o)
			walletOutputHashes = append(walletOutputHashes, o.Hash)
			walletAuxs[o.Address] = append(walletAuxs[o.Address], o.Hash)
		}
	}

	require.NotEmpty(t, walletOutputs)

	return []liveCreateTxnTestCase{
		{
			name: "invalid decimals",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   "0.0001",
						Hours:   "1",
					},
				},
			},
			err:  "to[0].coins has too many decimal places",
			code: http.StatusBadRequest,
		},

		{
			name: "overflowing hours",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   "0.001",
						Hours:   "1",
					},
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   "0.001",
						Hours:   fmt.Sprint(uint64(math.MaxUint64)),
					},
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   "0.001",
						Hours:   fmt.Sprint(uint64(math.MaxUint64) - 1),
					},
				},
			},
			err:  "total output hours error: uint64 addition overflow",
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient coins",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   fmt.Sprint(totalCoins + 1),
						Hours:   "1",
					},
				},
			},
			err:  "balance is not sufficient",
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient hours",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(0).Address.String(),
						Coins:   toDropletString(t, totalCoins),
						Hours:   fmt.Sprint(totalHours + 1),
					},
				},
			},
			err:  "hours are not sufficient",
			code: http.StatusBadRequest,
		},

		{
			// NOTE: this test will fail if "totalCoins - 1e3" does not require
			// all of the outputs to be spent, e.g. if there is an output with
			// "totalCoins - 1e3" coins in it.
			// TODO -- Check that the wallet does not have an output of 0.001,
			// because then this test cannot be performed, since there is no
			// way to use all outputs and produce change in that case.
			name: "valid request, manual one output with change, spend all",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins-1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.GetEntryAt(0).SkycoinAddress(),
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
		},

		{
			// NOTE: this test will fail if "totalCoins - 1e3" does not require
			// all of the outputs to be spent, e.g. if there is an output with
			// "totalCoins - 1e3" coins in it.
			// TODO -- Check that the wallet does not have an output of 0.001,
			// because then this test cannot be performed, since there is no
			// way to use all outputs and produce change in that case.
			name: "valid request, manual one output with change, spend all, unspecified change address",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses: addresses,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins-1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					// Address omitted -- will be checked later in the test body
					Coins: 1e3,
					Hours: remainingHours - 1,
				},
			},
		},

		{
			name: "valid request, manual one output with change, don't spend all",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, 1e3),
						Hours:   "1",
					},
				},
			},
			outputsSubset: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   1e3,
					Hours:   1,
				},
				// NOTE: change omitted,
				// change is too difficult to predict in this case, we are
				// just checking that not all uxouts get spent in the transaction
			},
		},

		{
			name: "valid request, manual one output no change",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins,
					Hours:   1,
				},
			},
		},

		{
			// NOTE: no reliable way to test the ignore unconfirmed behavior,
			// this test only checks that if IgnoreUnconfirmed is specified,
			// the API doesn't throw up some parsing error
			name: "valid request, manual one output no change, ignore unconfirmed",
			req: api.CreateTransactionRequest{
				IgnoreUnconfirmed: true,
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins,
					Hours:   1,
				},
			},
		},

		{
			name: "valid request, auto one output no change, share factor recalculates to 1.0",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: "0.5",
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins),
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins,
					Hours:   remainingHours,
				},
			},
		},

		{
			name: "valid request, auto two outputs with change",
			req: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: "0.5",
				},
				Addresses:     addresses,
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, 1e3),
					},
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins-2e3),
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   1e3,
				},
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins - 2e3,
				},
				{
					Address: w.GetEntryAt(0).SkycoinAddress(),
					Coins:   1e3,
				},
			},
			ignoreHours: true, // the hours are too unpredictable
		},

		{
			name: "uxout does not exist",
			req: api.CreateTransactionRequest{
				UxOuts: []string{unknownOutput.Hex()},
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins),
						Hours:   "1",
					},
				},
			},
			err:  fmt.Sprintf("unspent output of %s does not exist", unknownOutput.Hex()),
			code: http.StatusBadRequest,
		},

		{
			name: "insufficient balance with uxouts",
			req: api.CreateTransactionRequest{
				UxOuts: []string{walletOutputs[0].Hash},
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins+1e3),
						Hours:   "1",
					},
				},
			},
			err:  "balance is not sufficient",
			code: http.StatusBadRequest,
		},

		{
			// NOTE: expects wallet to have multiple outputs with non-zero coins
			name: "insufficient hours with uxouts",
			req: api.CreateTransactionRequest{
				UxOuts: []string{walletOutputs[0].Hash},
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, 1e3),
						Hours:   fmt.Sprint(totalHours + 1),
					},
				},
			},
			err:  "hours are not sufficient",
			code: http.StatusBadRequest,
		},

		{
			name: "valid request, uxouts specified",
			req: api.CreateTransactionRequest{
				// NOTE: all uxouts are provided, which has the same behavior as
				// not providing any uxouts or addresses.
				// Using a subset of uxouts makes the wallet setup very
				// difficult, especially to make deterministic, in the live test
				// More complex cases should be covered by unit tests
				UxOuts: walletOutputHashes,
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins-1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.GetEntryAt(0).SkycoinAddress(),
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
			additionalRespVerify: func(t *testing.T, r *api.CreateTransactionResponse) {
				require.Equal(t, len(walletOutputHashes), len(r.Transaction.In))
			},
		},

		{
			name: "valid request, addresses specified",
			req: api.CreateTransactionRequest{
				// NOTE: all addresses are provided, which has the same behavior as
				// not providing any addresses.
				// Using a subset of addresses makes the wallet setup very
				// difficult, especially to make deterministic, in the live test
				// More complex cases should be covered by unit tests
				Addresses: addresses,
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &defaultChangeAddress,
				To: []api.Receiver{
					{
						Address: w.GetEntryAt(1).Address.String(),
						Coins:   toDropletString(t, totalCoins-1e3),
						Hours:   "1",
					},
				},
			},
			outputs: []coin.TransactionOutput{
				{
					Address: w.GetEntryAt(1).SkycoinAddress(),
					Coins:   totalCoins - 1e3,
					Hours:   1,
				},
				{
					Address: w.GetEntryAt(0).SkycoinAddress(),
					Coins:   1e3,
					Hours:   remainingHours - 1,
				},
			},
		},
	}
}

func TestLiveCreateTransaction(t *testing.T) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()

	w, totalCoins, totalHours, _ := prepareAndCheckWallet(t, c, 2e6, 20)

	cases := makeLiveCreateTxnTestCases(t, w, totalCoins, totalHours)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.False(t, len(tc.outputs) != 0 && len(tc.outputsSubset) != 0, "outputs and outputsSubset can't both be set")

			result, err := c.CreateTransaction(tc.req)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)

			assertCreateTransactionResult(t, c, tc, result, true, nil)
		})
	}
}

func TestLiveWalletCreateTransactionSpecificUnsigned(t *testing.T) {
	testLiveWalletCreateTransactionSpecific(t, true)
}

func TestLiveWalletCreateTransactionSpecificSigned(t *testing.T) {
	testLiveWalletCreateTransactionSpecific(t, false)
}

func testLiveWalletCreateTransactionSpecific(t *testing.T, unsigned bool) {
	if !doLive(t) {
		return
	}

	requireWalletEnv(t)

	c := newClient()

	w, totalCoins, totalHours, password := prepareAndCheckWallet(t, c, 2e6, 20)

	remainingHours := fee.RemainingHours(totalHours, params.UserVerifyTxn.BurnFactor)
	require.True(t, remainingHours > 1)

	// Split outputs into those held by the wallet and those not
	addresses := make([]string, w.EntriesLen())
	addressMap := make(map[string]struct{}, w.EntriesLen())
	for i, e := range w.GetEntries() {
		addresses[i] = e.Address.String()
		addressMap[e.Address.String()] = struct{}{}
	}

	outputs, err := c.Outputs()
	require.NoError(t, err)

	var walletOutputs readable.UnspentOutputs
	walletAuxs := make(map[string][]string)
	var nonWalletOutputs readable.UnspentOutputs
	for _, o := range outputs.HeadOutputs {
		if _, ok := addressMap[o.Address]; ok {
			walletOutputs = append(walletOutputs, o)
			walletAuxs[o.Address] = append(walletAuxs[o.Address], o.Hash)
		} else {
			nonWalletOutputs = append(nonWalletOutputs, o)
		}
	}

	require.NotEmpty(t, walletOutputs)
	require.NotEmpty(t, nonWalletOutputs)

	defaultChangeAddress := w.GetEntryAt(0).Address.String()

	baseCases := makeLiveCreateTxnTestCases(t, w, totalCoins, totalHours)

	type liveWalletCreateTxnTestCase struct {
		liveCreateTxnTestCase
		password string
		walletID string
	}

	cases := make([]liveWalletCreateTxnTestCase, len(baseCases))
	for i, tc := range baseCases {
		cases[i] = liveWalletCreateTxnTestCase{
			liveCreateTxnTestCase: tc,
			walletID:              w.Filename(),
			password:              password,
		}
	}

	cases = append(cases, []liveWalletCreateTxnTestCase{
		{
			walletID: w.Filename(),
			password: password,
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "uxout not held by the wallet",
				req: api.CreateTransactionRequest{
					UxOuts: []string{nonWalletOutputs[0].Hash},
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(1).Address.String(),
							Coins:   nonWalletOutputs[0].Coins,
							Hours:   "1",
						},
					},
				},
				err:  "uxout is not owned by any address in the wallet",
				code: http.StatusBadRequest,
			},
		},

		{
			walletID: w.Filename(),
			password: password,
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "specified addresses not in wallet",
				req: api.CreateTransactionRequest{
					Addresses: []string{testutil.MakeAddress().String()},
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(1).Address.String(),
							Coins:   toDropletString(t, totalCoins),
							Hours:   "1",
						},
					},
				},
				err:  "address not found in wallet",
				code: http.StatusBadRequest,
			},
		},

		{
			walletID: w.Filename(),
			password: password,
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "valid request, addresses and uxouts not specified",
				req: api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(1).Address.String(),
							Coins:   toDropletString(t, totalCoins-1e3),
							Hours:   "1",
						},
					},
				},
				outputs: []coin.TransactionOutput{
					{
						Address: w.GetEntryAt(1).SkycoinAddress(),
						Coins:   totalCoins - 1e3,
						Hours:   1,
					},
					{
						Address: w.GetEntryAt(0).SkycoinAddress(),
						Coins:   1e3,
						Hours:   remainingHours - 1,
					},
				},
			},
		},
	}...)

	if w.IsEncrypted() {
		cases = append(cases, liveWalletCreateTxnTestCase{
			walletID: w.Filename(),
			password: password + "foo",
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "invalid password",
				req: api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(0).Address.String(),
							Coins:   "1000",
							Hours:   "1",
						},
					},
				},
				err:  "invalid password",
				code: http.StatusBadRequest,
			},
		})

		cases = append(cases, liveWalletCreateTxnTestCase{
			walletID: w.Filename(),
			password: "",
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "password not provided",
				req: api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(0).Address.String(),
							Coins:   "1000",
							Hours:   "1",
						},
					},
				},
				err:  "missing password",
				code: http.StatusBadRequest,
			},
		})

	} else {
		err := "wallet is not encrypted"
		if unsigned {
			err = "password must not be used for unsigned transactions"
		}

		cases = append(cases, liveWalletCreateTxnTestCase{
			walletID: w.Filename(),
			password: password + "foo",
			liveCreateTxnTestCase: liveCreateTxnTestCase{
				name: "password provided for unencrypted wallet",
				req: api.CreateTransactionRequest{
					HoursSelection: api.HoursSelection{
						Type: transaction.HoursSelectionTypeManual,
					},
					ChangeAddress: &defaultChangeAddress,
					To: []api.Receiver{
						{
							Address: w.GetEntryAt(0).Address.String(),
							Coins:   "1000",
							Hours:   "1",
						},
					},
				},
				err:  err,
				code: http.StatusBadRequest,
			},
		})
	}

	for _, tc := range cases {
		name := fmt.Sprintf("unsigned=%v %s", unsigned, tc.name)
		t.Run(name, func(t *testing.T) {
			require.False(t, len(tc.outputs) != 0 && len(tc.outputsSubset) != 0, "outputs and outputsSubset can't both be set")

			// Fetch a copy of the wallet to look for modifications to the wallet
			// after the transaction is created
			w, err := c.Wallet(tc.walletID)
			require.NoError(t, err)

			req := api.WalletCreateTransactionRequest{
				CreateTransactionRequest: tc.req,
				WalletID:                 tc.walletID,
				Password:                 tc.password,
				Unsigned:                 unsigned,
			}

			if tc.err != "" {
				tc.err = fmt.Sprintf("%d %s - %s", tc.code, http.StatusText(tc.code), tc.err)
			}

			result, err := c.WalletCreateTransaction(req)
			if tc.err != "" {
				assertResponseError(t, err, tc.code, tc.err)
				return
			}

			require.NoError(t, err)

			assertCreateTransactionResult(t, c, tc.liveCreateTxnTestCase, result, unsigned, w)
		})
	}
}

func getLastChangeEntry(t *testing.T, w *api.WalletResponse) *readable.WalletEntry {
	require.Equal(t, wallet.WalletTypeBip44, w.Meta.Type)

	// Find the last "change" entry
	require.NotEmpty(t, w.Entries)
	sort.Slice(w.Entries, func(i, j int) bool {
		if *w.Entries[i].Change == *w.Entries[j].Change {
			return *w.Entries[i].ChildNumber > *w.Entries[j].ChildNumber
		}
		return *w.Entries[i].Change > *w.Entries[j].Change
	})

	lastChangeEntry := w.Entries[0]
	if *lastChangeEntry.Change != bip44.ChangeChainIndex {
		// no change entry
		return nil
	}

	return &lastChangeEntry
}

func isNullAddress(a string) bool {
	if a == "" {
		return true
	}

	addr := cipher.MustDecodeBase58Address(a)
	return addr.Null()
}

func assertCreateTransactionResult(t *testing.T, c *api.Client, tc liveCreateTxnTestCase, result *api.CreateTransactionResponse, unsigned bool, w *api.WalletResponse) {
	if len(tc.outputsSubset) == 0 {
		require.Equal(t, len(tc.outputs), len(result.Transaction.Out))
	}

	for i, o := range tc.outputs {
		coins, err := droplet.FromString(result.Transaction.Out[i].Coins)
		require.NoError(t, err)
		require.Equal(t, o.Coins, coins, "[%d] %d != %d", i, o.Coins, coins)

		if !tc.ignoreHours {
			hours, err := strconv.ParseUint(result.Transaction.Out[i].Hours, 10, 64)
			require.NoError(t, err)
			require.Equal(t, o.Hours, hours, "[%d] %d != %d", i, o.Hours, hours)
		}

		if o.Address.Null() {
			// The final change output may not have the address specified,
			// if the ChangeAddress was not specified in the wallet params.
			require.Equal(t, i, len(tc.outputs)-1)
			require.Nil(t, tc.req.ChangeAddress)
			changeAddr := result.Transaction.Out[i].Address
			require.False(t, isNullAddress(changeAddr))

			if w != nil && w.Meta.Type == wallet.WalletTypeBip44 {
				// Check that the change address was a new address generated
				// from the wallet's change path

				// Get the update wallet from the API.
				// Look for the last change address.
				// It should match the change address that was used.
				// Compare it to the previous wallet
				w2, err := c.Wallet(w.Meta.Filename)
				require.NoError(t, err)
				lastChangeEntry := getLastChangeEntry(t, w2)

				// Compare it to the initial wallet state.
				// It should be a new address with an incremented child number
				prevLastChangeEntry := getLastChangeEntry(t, w)
				require.NotEqual(t, prevLastChangeEntry, lastChangeEntry)
				if prevLastChangeEntry == nil {
					require.Equal(t, uint32(0), *lastChangeEntry.ChildNumber)
				} else {
					require.Equal(t, *prevLastChangeEntry.ChildNumber+1, *lastChangeEntry.ChildNumber)
				}

				// Make sure that the last change address in the wallet was used
				require.False(t, isNullAddress(lastChangeEntry.Address))
				require.Equal(t, changeAddr, lastChangeEntry.Address)
			} else {
				// Check that the automatically-selected change address was one
				// of the addresses for the UTXOs spent by the transaction
				changeAddrFound := false
				for _, x := range result.Transaction.In {
					require.False(t, isNullAddress(x.Address))
					if changeAddr == x.Address {
						changeAddrFound = true
						break
					}
				}

				require.True(t, changeAddrFound)
			}
		} else {
			require.Equal(t, o.Address.String(), result.Transaction.Out[i].Address)
		}
	}

	// The wallet should be unmodified if the wallet type is not bip44
	if w != nil && w.Meta.Type != wallet.WalletTypeBip44 {
		w2, err := c.Wallet(w.Meta.Filename)
		require.NoError(t, err)
		require.Equal(t, w, w2)
	}

	assertEncodeTxnMatchesTxn(t, result)
	assertRequestedCoins(t, tc.req.To, result.Transaction.Out)
	assertCreatedTransactionValid(t, result.Transaction, unsigned)

	if tc.req.HoursSelection.Type == transaction.HoursSelectionTypeManual {
		assertRequestedHours(t, tc.req.To, result.Transaction.Out)
	}

	if tc.additionalRespVerify != nil {
		tc.additionalRespVerify(t, result)
	}

	assertVerifyTransaction(t, result.EncodedTransaction, unsigned)
}

func TestLiveWalletCreateTransactionRandomUnsigned(t *testing.T) {
	testLiveWalletCreateTransactionRandom(t, true)
}

func TestLiveWalletCreateTransactionRandomSigned(t *testing.T) {
	testLiveWalletCreateTransactionRandom(t, false)
}

func testLiveWalletCreateTransactionRandom(t *testing.T, unsigned bool) {
	if !doLive(t) {
		return
	}

	debug := false
	tLog := func(t *testing.T, args ...interface{}) {
		if debug {
			t.Log(args...)
		}
	}
	tLogf := func(t *testing.T, msg string, args ...interface{}) {
		if debug {
			t.Logf(msg, args...)
		}
	}

	requireWalletEnv(t)

	c := newClient()

	w, totalCoins, totalHours, password := prepareAndCheckWallet(t, c, 2e6, 20)

	if w.IsEncrypted() {
		t.Skip("Skipping TestLiveWalletCreateTransactionRandom tests with encrypted wallet")
		return
	}

	remainingHours := fee.RemainingHours(totalHours, params.UserVerifyTxn.BurnFactor)
	require.True(t, remainingHours > 1)

	assertTxnOutputCount := func(t *testing.T, changeAddress string, nOutputs int, result *api.CreateTransactionResponse) {
		nResultOutputs := len(result.Transaction.Out)
		require.True(t, nResultOutputs == nOutputs || nResultOutputs == nOutputs+1)
		hasChange := nResultOutputs == nOutputs+1
		changeOutput := result.Transaction.Out[nResultOutputs-1]
		if hasChange {
			require.Equal(t, changeOutput.Address, changeAddress)
		}

		tLog(t, "hasChange", hasChange)
		if hasChange {
			tLog(t, "changeCoins", changeOutput.Coins)
			tLog(t, "changeHours", changeOutput.Hours)
		}
	}

	iterations := 250
	maxOutputs := 10
	destAddrs := make([]cipher.Address, maxOutputs)
	for i := range destAddrs {
		destAddrs[i] = testutil.MakeAddress()
	}

	for i := 0; i < iterations; i++ {
		tLog(t, "iteration", i)
		tLog(t, "totalCoins", totalCoins)
		tLog(t, "totalHours", totalHours)

		spendableHours := fee.RemainingHours(totalHours, params.UserVerifyTxn.BurnFactor)
		tLog(t, "spendableHours", spendableHours)

		coins := rand.Intn(int(totalCoins)) + 1
		coins -= coins % int(params.UserVerifyTxn.MaxDropletDivisor())
		if coins == 0 {
			coins = int(params.UserVerifyTxn.MaxDropletDivisor())
		}
		hours := rand.Intn(int(spendableHours + 1))
		nOutputs := rand.Intn(maxOutputs) + 1

		tLog(t, "sendCoins", coins)
		tLog(t, "sendHours", hours)

		changeAddress := w.GetEntryAt(0).Address.String()

		shareFactor := strconv.FormatFloat(rand.Float64(), 'f', 8, 64)

		tLog(t, "shareFactor", shareFactor)

		to := make([]api.Receiver, 0, nOutputs)
		remainingHours := hours
		remainingCoins := coins
		for i := 0; i < nOutputs; i++ {
			if remainingCoins == 0 {
				break
			}

			receiver := api.Receiver{}
			receiver.Address = destAddrs[rand.Intn(len(destAddrs))].String()

			if i == nOutputs-1 {
				var err error
				receiver.Coins, err = droplet.ToString(uint64(remainingCoins))
				require.NoError(t, err)
				receiver.Hours = fmt.Sprint(remainingHours)

				remainingCoins = 0
				remainingHours = 0
			} else {
				receiverCoins := rand.Intn(remainingCoins) + 1
				receiverCoins -= receiverCoins % int(params.UserVerifyTxn.MaxDropletDivisor())
				if receiverCoins == 0 {
					receiverCoins = int(params.UserVerifyTxn.MaxDropletDivisor())
				}

				var err error
				receiver.Coins, err = droplet.ToString(uint64(receiverCoins))
				require.NoError(t, err)
				remainingCoins -= receiverCoins

				receiverHours := rand.Intn(remainingHours + 1)
				receiver.Hours = fmt.Sprint(receiverHours)
				remainingHours -= receiverHours
			}

			to = append(to, receiver)
		}

		// Remove duplicate outputs
		dup := make(map[api.Receiver]struct{}, len(to))
		newTo := make([]api.Receiver, 0, len(dup))
		for _, o := range to {
			if _, ok := dup[o]; !ok {
				dup[o] = struct{}{}
				newTo = append(newTo, o)
			}
		}
		to = newTo

		nOutputs = len(to)
		tLog(t, "nOutputs", nOutputs)

		rand.Shuffle(len(to), func(i, j int) {
			to[i], to[j] = to[j], to[i]
		})

		for i, o := range to {
			tLogf(t, "to[%d].Hours %s\n", i, o.Hours)
		}

		autoTo := make([]api.Receiver, len(to))
		for i, o := range to {
			autoTo[i] = api.Receiver{
				Address: o.Address,
				Coins:   o.Coins,
				Hours:   "",
			}
		}

		// Remove duplicate outputs
		dup = make(map[api.Receiver]struct{}, len(autoTo))
		newAutoTo := make([]api.Receiver, 0, len(dup))
		for _, o := range autoTo {
			if _, ok := dup[o]; !ok {
				dup[o] = struct{}{}
				newAutoTo = append(newAutoTo, o)
			}
		}
		autoTo = newAutoTo

		nAutoOutputs := len(autoTo)
		tLog(t, "nAutoOutputs", nAutoOutputs)

		for i, o := range autoTo {
			tLogf(t, "autoTo[%d].Coins %s\n", i, o.Coins)
		}

		// Auto, random share factor

		result, err := c.WalletCreateTransaction(api.WalletCreateTransactionRequest{
			WalletID: w.Filename(),
			Password: password,
			Unsigned: unsigned,
			CreateTransactionRequest: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: shareFactor,
				},
				ChangeAddress: &changeAddress,
				To:            autoTo,
			},
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction, unsigned)
		assertVerifyTransaction(t, result.EncodedTransaction, unsigned)

		// Auto, share factor 0

		result, err = c.WalletCreateTransaction(api.WalletCreateTransactionRequest{
			WalletID: w.Filename(),
			Password: password,
			Unsigned: unsigned,
			CreateTransactionRequest: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: "0",
				},
				ChangeAddress: &changeAddress,
				To:            autoTo,
			},
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction, unsigned)
		assertVerifyTransaction(t, result.EncodedTransaction, unsigned)

		// Check that the non-change outputs have 0 hours
		for _, o := range result.Transaction.Out[:nAutoOutputs] {
			require.Equal(t, "0", o.Hours)
		}

		// Auto, share factor 1

		result, err = c.WalletCreateTransaction(api.WalletCreateTransactionRequest{
			Unsigned: unsigned,
			WalletID: w.Filename(),
			Password: password,
			CreateTransactionRequest: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type:        transaction.HoursSelectionTypeAuto,
					Mode:        transaction.HoursSelectionModeShare,
					ShareFactor: "1",
				},
				ChangeAddress: &changeAddress,
				To:            autoTo,
			},
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nAutoOutputs, result)
		assertRequestedCoins(t, autoTo, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction, unsigned)
		assertVerifyTransaction(t, result.EncodedTransaction, unsigned)

		// Check that the change output has 0 hours
		if len(result.Transaction.Out) > nAutoOutputs {
			require.Equal(t, "0", result.Transaction.Out[nAutoOutputs].Hours)
		}

		// Manual

		result, err = c.WalletCreateTransaction(api.WalletCreateTransactionRequest{
			Unsigned: unsigned,
			WalletID: w.Filename(),
			Password: password,
			CreateTransactionRequest: api.CreateTransactionRequest{
				HoursSelection: api.HoursSelection{
					Type: transaction.HoursSelectionTypeManual,
				},
				ChangeAddress: &changeAddress,
				To:            to,
			},
		})
		require.NoError(t, err)

		assertEncodeTxnMatchesTxn(t, result)
		assertTxnOutputCount(t, changeAddress, nOutputs, result)
		assertRequestedCoins(t, to, result.Transaction.Out)
		assertRequestedHours(t, to, result.Transaction.Out)
		assertCreatedTransactionValid(t, result.Transaction, unsigned)
		assertVerifyTransaction(t, result.EncodedTransaction, unsigned)
	}
}

func assertEncodeTxnMatchesTxn(t *testing.T, result *api.CreateTransactionResponse) {
	require.NotEmpty(t, result.EncodedTransaction)
	emptyTxn := &coin.Transaction{}
	require.NotEqual(t, emptyTxn.MustSerializeHex(), result.EncodedTransaction)

	txn, err := result.Transaction.ToTransaction()
	require.NoError(t, err)

	require.Equal(t, txn.MustSerializeHex(), result.EncodedTransaction)
	require.Equal(t, int(txn.Length), len(txn.MustSerialize()))
}

func assertRequestedCoins(t *testing.T, to []api.Receiver, out []api.CreatedTransactionOutput) {
	var requestedCoins uint64
	for _, o := range to {
		c, err := droplet.FromString(o.Coins)
		require.NoError(t, err)
		requestedCoins += c
	}

	var sentCoins uint64
	for _, o := range out[:len(to)] { // exclude change output
		c, err := droplet.FromString(o.Coins)
		require.NoError(t, err)
		sentCoins += c
	}

	require.Equal(t, requestedCoins, sentCoins)
}

func assertRequestedHours(t *testing.T, to []api.Receiver, out []api.CreatedTransactionOutput) {
	for i, o := range out[:len(to)] { // exclude change output
		toHours, err := strconv.ParseUint(to[i].Hours, 10, 64)
		require.NoError(t, err)

		outHours, err := strconv.ParseUint(o.Hours, 10, 64)
		require.NoError(t, err)
		require.Equal(t, toHours, outHours)
	}
}

func assertVerifyTransaction(t *testing.T, encodedTransaction string, unsigned bool) {
	c := newClient()
	_, err := c.VerifyTransaction(api.VerifyTransactionRequest{
		EncodedTransaction: encodedTransaction,
		Unsigned:           false,
	})
	if unsigned {
		assertResponseError(t, err, http.StatusUnprocessableEntity, "Transaction violates hard constraint: Unsigned input in transaction")
	} else {
		require.NoError(t, err)
	}

	_, err = c.VerifyTransaction(api.VerifyTransactionRequest{
		EncodedTransaction: encodedTransaction,
		Unsigned:           true,
	})
	if unsigned {
		require.NoError(t, err)
	} else {
		assertResponseError(t, err, http.StatusUnprocessableEntity, "Transaction violates hard constraint: Unsigned transaction must contain a null signature")
	}
}

func assertCreatedTransactionValid(t *testing.T, r api.CreatedTransaction, unsigned bool) {
	require.NotEmpty(t, r.In)
	require.NotEmpty(t, r.Out)

	require.Equal(t, len(r.In), len(r.Sigs))
	if unsigned {
		for _, s := range r.Sigs {
			ss := cipher.MustSigFromHex(s)
			require.True(t, ss.Null())
		}
	}

	fee, err := strconv.ParseUint(r.Fee, 10, 64)
	require.NoError(t, err)

	require.NotEqual(t, uint64(0), fee)

	var inputHours uint64
	var inputCoins uint64
	for _, in := range r.In {
		require.NotNil(t, in.CalculatedHours)
		calculatedHours, err := strconv.ParseUint(in.CalculatedHours, 10, 64)
		require.NoError(t, err)
		inputHours, err = mathutil.AddUint64(inputHours, calculatedHours)
		require.NoError(t, err)

		require.NotNil(t, in.Hours)
		hours, err := strconv.ParseUint(in.Hours, 10, 64)
		require.NoError(t, err)

		require.True(t, hours <= calculatedHours)

		require.NotNil(t, in.Coins)
		coins, err := droplet.FromString(in.Coins)
		require.NoError(t, err)
		inputCoins, err = mathutil.AddUint64(inputCoins, coins)
		require.NoError(t, err)
	}

	var outputHours uint64
	var outputCoins uint64
	for _, out := range r.Out {
		hours, err := strconv.ParseUint(out.Hours, 10, 64)
		require.NoError(t, err)
		outputHours, err = mathutil.AddUint64(outputHours, hours)
		require.NoError(t, err)

		coins, err := droplet.FromString(out.Coins)
		require.NoError(t, err)
		outputCoins, err = mathutil.AddUint64(outputCoins, coins)
		require.NoError(t, err)
	}

	require.True(t, inputHours > outputHours)
	require.Equal(t, inputHours-outputHours, fee)

	require.Equal(t, inputCoins, outputCoins)

	require.Equal(t, uint8(0), r.Type)
	require.NotEmpty(t, r.Length)
}

func getTransaction(t *testing.T, c *api.Client, txid string) *readable.TransactionWithStatus {
	tx, err := c.Transaction(txid)
	if err != nil {
		t.Fatalf("%v", err)
	}

	return tx
}

// getAddressBalance gets balance of given address.
// Returns coins and coin hours.
func getAddressBalance(t *testing.T, c *api.Client, addr string) (uint64, uint64) { //nolint:unparam
	bp, err := c.Balance([]string{addr})
	if err != nil {
		t.Fatalf("%v", err)
	}
	return bp.Confirmed.Coins, bp.Confirmed.Hours
}
