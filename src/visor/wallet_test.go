package visor

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/transaction"
	"github.com/SkycoinProject/skycoin/src/visor/blockdb"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
	"github.com/SkycoinProject/skycoin/src/wallet"
	"github.com/SkycoinProject/skycoin/src/wallet/collection"
	"github.com/SkycoinProject/skycoin/src/wallet/crypto"
)

func TestCreateTransaction(t *testing.T) {
	addrs := make([]cipher.Address, 3)
	for i := range addrs {
		addrs[i] = testutil.MakeAddress()
	}

	uxOuts := make([]cipher.SHA256, 3)
	for i := range uxOuts {
		uxOuts[i] = testutil.RandSHA256(t)
	}

	validCreateTxnParams := CreateTransactionParams{
		Addresses: addrs,
	}

	validParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   1e6,
				Hours:   7,
			},
		},
	}

	insufficientBalanceParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   2e6,
				Hours:   7,
			},
		},
	}

	invalidParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: cipher.Address{}, // null address violates user constraints
				Coins:   1e6,
				Hours:   7,
			},
		},
	}

	getArrayRet := coin.UxArray{
		{
			Head: coin.UxHead{
				Time:  uint64(time.Now().Unix()) - 3700,
				BkSeq: 100,
			},
			Body: coin.UxBody{
				SrcTransaction: testutil.RandSHA256(t),
				Address:        addrs[1],
				Coins:          1e6,
				Hours:          100,
			},
		},
	}
	uxOuts[1] = getArrayRet[0].Hash()

	inputs := []TransactionInput{
		{
			UxOut:           getArrayRet[0],
			CalculatedHours: getArrayRet[0].Body.Hours + 1,
		},
	}

	txn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{getArrayRet[0].Hash()},
		Out:    validParams.To,
	}
	err := txn.UpdateHeader()
	require.NoError(t, err)

	invalidParamsTxn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{getArrayRet[0].Hash()},
		Out:    invalidParams.To,
	}
	err = invalidParamsTxn.UpdateHeader()
	require.NoError(t, err)

	headBlock := &coin.SignedBlock{
		Block: coin.Block{
			Head: coin.BlockHeader{
				Time: uint64(time.Now().Unix()),
			},
		},
	}

	cases := []struct {
		name   string
		txn    *coin.Transaction
		inputs []TransactionInput
		p      transaction.Params
		wp     CreateTransactionParams
		err    error

		blockchainHead    *coin.SignedBlock
		blockchainHeadErr error

		unconfirmedTxns []coin.Transaction
		uxOuts          []cipher.SHA256
		forEachErr      error

		getArrayInputs []cipher.SHA256
		getArray       coin.UxArray
		getArrayErr    error

		getUnspentHashesOfAddrs    blockdb.AddressHashes
		getUnspentHashesOfAddrsErr error

		verifyErr error
	}{
		{
			name:              "Blockchain.Head failed",
			p:                 validParams,
			wp:                validCreateTxnParams,
			blockchainHeadErr: errors.New("failure"),
			err:               errors.New("failure"),
		},

		{
			name:                       "GetUnspentHashesOfAddrs failed",
			p:                          validParams,
			wp:                         validCreateTxnParams,
			getUnspentHashesOfAddrsErr: errors.New("failure"),
			err:                        errors.New("failure"),
		},

		{
			name:                    "no unspents found for addresses",
			p:                       validParams,
			wp:                      validCreateTxnParams,
			getUnspentHashesOfAddrs: nil,
			err:                     transaction.ErrNoUnspents,
		},

		{
			name: "Unconfirmed.ForEach failed",
			p:    validParams,
			wp:   validCreateTxnParams,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				addrs[1]: uxOuts,
			},
			forEachErr: errors.New("failure"),
			err:        errors.New("failure"),
		},

		{
			name: "Unspent.GetArray failed",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			getArrayInputs: uxOuts,
			getArrayErr:    errors.New("failure"),
			err:            errors.New("failure"),
		},

		{
			name: "insufficient balance",
			p:    insufficientBalanceParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			err:            transaction.ErrInsufficientBalance,
		},

		{
			name: "invalid params for transaction.Create, uxouts",
			p:    invalidParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            invalidParamsTxn,
			inputs:         inputs,
			err:            transaction.ErrNullAddressReceiver,
		},

		{
			name: "blockchain verify error, uxouts",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
			verifyErr:      NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
			err:            NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
		},

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
			err:  ErrUxOutsOrAddressesRequired,
		},

		{
			name: "ok, addresses",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: addrs,
			},
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				addrs[1]: uxOuts,
			},
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
		},

		{
			name: "ok, uxouts",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b := &MockBlockchainer{}
			ut := &MockUnconfirmedTransactionPooler{}
			up := &MockUnspentPooler{}

			b.On("Head", matchDBTx).Return(tc.blockchainHead, tc.blockchainHeadErr)
			up.On("GetUnspentHashesOfAddrs", matchDBTx, tc.wp.Addresses).Return(tc.getUnspentHashesOfAddrs, tc.getUnspentHashesOfAddrsErr)

			ut.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(unconfirmedForEachMockRun(t, tc.unconfirmedTxns, tc.uxOuts, tc.wp.IgnoreUnconfirmed))

			up.On("GetArray", matchDBTx, mock.MatchedBy(matchUxOutsAnyOrder(tc.getArrayInputs))).Return(tc.getArray, tc.getArrayErr)
			b.On("Unspent").Return(up)

			if tc.txn != nil {
				b.On("VerifySingleTxnSoftHardConstraints", matchDBTx, *tc.txn, params.MainNetDistribution, params.UserVerifyTxn, TxnUnsigned).Return(nil, nil, tc.verifyErr)
			}

			db, shutdown := prepareDB(t)
			defer shutdown()

			v := &Visor{
				db:          db,
				blockchain:  b,
				unconfirmed: ut,
				Config: Config{
					Distribution: params.MainNetDistribution,
				},
			}

			txn, inputs, err := v.CreateTransaction(tc.p, tc.wp)
			require.Equal(t, tc.err, err)
			if tc.err != nil {
				return
			}

			require.Equal(t, tc.txn, txn)
			require.Equal(t, tc.inputs, inputs)
		})
	}
}

func prepareWltDir() string {
	dir, err := ioutil.TempDir("", "wallets")
	if err != nil {
		panic(err)
	}

	return dir
}

func makeEntries(n int) ([]wallet.Entry, []cipher.Address) {
	addrs := make([]cipher.Address, n)
	entries := make([]wallet.Entry, n)
	for i := range addrs {
		p, s := cipher.GenerateKeyPair()
		a := cipher.AddressFromPubKey(p)
		entries[i] = wallet.Entry{
			Address: a,
			Public:  p,
			Secret:  s,
		}
		addrs[i] = a
	}
	return entries, addrs
}

func TestWalletCreateTransaction(t *testing.T) {
	// Create arbitrary entries for a "collection" wallet
	entries, addrs := makeEntries(3)

	// Load the 5th through 8th entries for a known bip44 wallet.
	// These entries will be used for the test data.
	bip44Seed := "voyage say extend find sheriff surge priority merit ignore maple cash argue"
	w, err := wallet.NewWallet(
		"bip44.wlt",
		"label",
		bip44Seed,
		wallet.Options{
			Type:      wallet.WalletTypeBip44,
			GenerateN: 8,
		})
	require.NoError(t, err)
	bip44Entries := func() wallet.Entries {
		entries, err := w.GetEntries()
		require.NoError(t, err)
		return entries[5:8]
	}()
	require.Len(t, bip44Entries, 3)
	bip44Addrs, err := func() ([]cipher.Address, error) {
		addrs, err := w.GetAddresses()
		require.NoError(t, err)
		return wallet.SkycoinAddresses(addrs), nil
	}()
	require.NoError(t, err)
	bip44Addrs = bip44Addrs[5:8]
	require.Len(t, bip44Addrs, 3)

	uxOuts := make([]cipher.SHA256, 3)
	for i := range uxOuts {
		uxOuts[i] = testutil.RandSHA256(t)
	}

	bip44UxOuts := make([]cipher.SHA256, 3)
	for i := range bip44UxOuts {
		bip44UxOuts[i] = testutil.RandSHA256(t)
	}

	validParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   1e6,
				Hours:   7,
			},
		},
	}

	insufficientBalanceParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: testutil.MakeAddress(),
				Coins:   2e6,
				Hours:   7,
			},
		},
	}

	invalidParams := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type: transaction.HoursSelectionTypeManual,
		},
		To: []coin.TransactionOutput{
			{
				Address: cipher.Address{}, // null address violates user constraints
				Coins:   1e6,
				Hours:   7,
			},
		},
	}

	getArrayRet := coin.UxArray{
		{
			Head: coin.UxHead{
				Time:  uint64(time.Now().Unix()) - 3700,
				BkSeq: 100,
			},
			Body: coin.UxBody{
				SrcTransaction: testutil.RandSHA256(t),
				Address:        addrs[1],
				Coins:          1e6,
				Hours:          100,
			},
		},
	}
	uxOuts[1] = getArrayRet[0].Hash()

	unknownUxOutGetArrayRet := append(getArrayRet, coin.UxOut{
		Head: coin.UxHead{
			Time:  uint64(time.Now().Unix()) - 3700,
			BkSeq: 100,
		},
		Body: coin.UxBody{
			SrcTransaction: testutil.RandSHA256(t),
			Address:        testutil.MakeAddress(),
			Coins:          1e6,
			Hours:          100,
		},
	})

	inputs := []TransactionInput{
		{
			UxOut:           getArrayRet[0],
			CalculatedHours: getArrayRet[0].Body.Hours + 1,
		},
	}

	txn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{getArrayRet[0].Hash()},
		Out:    validParams.To,
	}
	err = txn.UpdateHeader()
	require.NoError(t, err)

	invalidParamsTxn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{getArrayRet[0].Hash()},
		Out:    invalidParams.To,
	}
	err = invalidParamsTxn.UpdateHeader()
	require.NoError(t, err)

	bip44GetArrayRet := coin.UxArray{
		{
			Head: coin.UxHead{
				Time:  uint64(time.Now().Unix()) - 3700,
				BkSeq: 100,
			},
			Body: coin.UxBody{
				SrcTransaction: testutil.RandSHA256(t),
				Address:        bip44Addrs[1],
				Coins:          1e6,
				Hours:          100,
			},
		},
	}
	bip44UxOuts[1] = bip44GetArrayRet[0].Hash()

	bip44UnknownUxOutGetArrayRet := append(bip44GetArrayRet, coin.UxOut{
		Head: coin.UxHead{
			Time:  uint64(time.Now().Unix()) - 3700,
			BkSeq: 100,
		},
		Body: coin.UxBody{
			SrcTransaction: testutil.RandSHA256(t),
			Address:        testutil.MakeAddress(),
			Coins:          1e6,
			Hours:          100,
		},
	})

	bip44Inputs := []TransactionInput{
		{
			UxOut:           bip44GetArrayRet[0],
			CalculatedHours: bip44GetArrayRet[0].Body.Hours + 1,
		},
	}

	bip44Txn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{bip44GetArrayRet[0].Hash()},
		Out:    validParams.To,
	}
	err = bip44Txn.UpdateHeader()
	require.NoError(t, err)

	bip44InvalidParamsTxn := &coin.Transaction{
		Length: 183,
		Type:   0,
		Sigs:   make([]cipher.Sig, 1),
		In:     []cipher.SHA256{bip44GetArrayRet[0].Hash()},
		Out:    invalidParams.To,
	}
	err = invalidParamsTxn.UpdateHeader()
	require.NoError(t, err)

	headBlock := &coin.SignedBlock{
		Block: coin.Block{
			Head: coin.BlockHeader{
				Time: uint64(time.Now().Unix()),
			},
		},
	}

	type testCase struct {
		name       string
		txn        *coin.Transaction
		inputs     []TransactionInput
		p          transaction.Params
		wp         CreateTransactionParams
		signed     TxnSignedFlag
		walletID   string
		walletType string
		seed       string
		password   []byte
		err        error

		blockchainHead    *coin.SignedBlock
		blockchainHeadErr error

		unconfirmedTxns []coin.Transaction
		uxOuts          []cipher.SHA256
		forEachErr      error

		getArrayInputs []cipher.SHA256
		getArray       coin.UxArray
		getArrayErr    error

		getUnspentHashesOfAddrs    blockdb.AddressHashes
		getUnspentHashesOfAddrsErr error

		verifyErr error
	}

	baseCases := []testCase{
		{
			name:           "all wallet addresses",
			p:              validParams,
			wp:             CreateTransactionParams{},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				addrs[1]: uxOuts,
			},
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
		},

		{
			name: "specific wallet addresses",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: addrs,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				addrs[1]: uxOuts,
			},
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
		},

		{
			name: "specific uxouts",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
		},

		{
			name: "unknown wallet address",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: append(addrs, testutil.MakeAddress()),
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				addrs[1]: uxOuts,
			},
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			err:            wallet.ErrUnknownAddress,
		},

		{
			name: "unknown wallet uxouts",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       unknownUxOutGetArrayRet,
			err:            wallet.ErrUnknownUxOut,
		},

		{
			name: "insufficient balance",
			p:    insufficientBalanceParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			err:            transaction.ErrInsufficientBalance,
		},

		{
			name: "invalid params for transaction.Create, uxouts",
			p:    invalidParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            invalidParamsTxn,
			inputs:         inputs,
			err:            transaction.ErrNullAddressReceiver,
		},

		{
			name:              "Blockchain.Head failed",
			p:                 validParams,
			walletID:          "foo.wlt",
			walletType:        wallet.WalletTypeCollection,
			blockchainHeadErr: errors.New("failure"),
			err:               errors.New("failure"),
		},

		{
			name: "blockchain verify error, uxouts",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: uxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeCollection,
			blockchainHead: headBlock,
			getArrayInputs: uxOuts,
			getArray:       getArrayRet,
			txn:            txn,
			inputs:         inputs,
			verifyErr:      NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
			err:            NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
		},

		{
			name:           "all wallet addresses bip44 wallet",
			p:              validParams,
			wp:             CreateTransactionParams{},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				bip44Addrs[1]: bip44UxOuts,
			},
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			txn:            bip44Txn,
			inputs:         bip44Inputs,
		},

		{
			name: "specific wallet addresses bip44 wallet",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: bip44Addrs,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				bip44Addrs[1]: bip44UxOuts,
			},
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			txn:            bip44Txn,
			inputs:         bip44Inputs,
		},

		{
			name: "specific uxouts bip44 wallet",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: bip44UxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			txn:            bip44Txn,
			inputs:         bip44Inputs,
		},

		{
			name: "unknown wallet address bip44 wallet",
			p:    validParams,
			wp: CreateTransactionParams{
				Addresses: append(bip44Addrs, testutil.MakeAddress()),
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				bip44Addrs[1]: bip44UxOuts,
			},
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			err:            wallet.ErrUnknownAddress,
		},

		{
			name: "unknown wallet uxouts bip44 wallet",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: bip44UxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getArrayInputs: bip44UxOuts,
			getArray:       bip44UnknownUxOutGetArrayRet,
			err:            wallet.ErrUnknownUxOut,
		},

		{
			name: "insufficient balance bip44 wallet",
			p:    insufficientBalanceParams,
			wp: CreateTransactionParams{
				UxOuts: bip44UxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			err:            transaction.ErrInsufficientBalance,
		},

		{
			name: "invalid params for transaction.Create, uxouts bip44 wallet",
			p:    invalidParams,
			wp: CreateTransactionParams{
				UxOuts: bip44UxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			txn:            bip44InvalidParamsTxn,
			inputs:         bip44Inputs,
			err:            transaction.ErrNullAddressReceiver,
		},

		{
			name:              "Blockchain.Head failed bip44 wallet",
			p:                 validParams,
			walletID:          "foo.wlt",
			walletType:        wallet.WalletTypeBip44,
			seed:              bip44Seed,
			blockchainHeadErr: errors.New("failure"),
			err:               errors.New("failure"),
		},

		{
			name: "blockchain verify error, uxouts bip44 wallet",
			p:    validParams,
			wp: CreateTransactionParams{
				UxOuts: bip44UxOuts,
			},
			walletID:       "foo.wlt",
			walletType:     wallet.WalletTypeBip44,
			seed:           bip44Seed,
			blockchainHead: headBlock,
			getArrayInputs: bip44UxOuts,
			getArray:       bip44GetArrayRet,
			txn:            bip44Txn,
			inputs:         bip44Inputs,
			verifyErr:      NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
			err:            NewErrTxnViolatesSoftConstraint(errors.New("Violates soft constraints")),
		},
	}

	cases := make([]testCase, len(baseCases)*2)
	copy(cases, baseCases)
	copy(cases[len(baseCases):], baseCases)
	for i := range cases[:len(baseCases)] {
		cases[i].signed = TxnUnsigned
		cases[i].password = nil
	}
	for i := range cases[len(baseCases):] {
		cases[i+len(baseCases)].signed = TxnSigned
		cases[i+len(baseCases)].password = []byte("foo")
	}

	for _, tc := range cases {
		name := fmt.Sprintf("signed-flag=%d %s", tc.signed, tc.name)
		t.Run(name, func(t *testing.T) {
			require.NotEmpty(t, tc.walletID)

			ws, err := wallet.NewService(wallet.Config{
				EnableWalletAPI: true,
				CryptoType:      crypto.CryptoTypeScryptChacha20poly1305Insecure,
				WalletDir:       prepareWltDir(),
			})
			require.NoError(t, err)

			generateN := uint64(0)
			if tc.walletType == wallet.WalletTypeBip44 {
				// Generate 8 addresses, the last 3 were used in the test data above
				generateN = 8
			}

			_, err = ws.CreateWallet(tc.walletID, wallet.Options{
				Coin:       wallet.CoinTypeSkycoin,
				Encrypt:    len(tc.password) != 0,
				Password:   tc.password,
				CryptoType: crypto.CryptoTypeScryptChacha20poly1305Insecure,
				Type:       tc.walletType,
				GenerateN:  generateN,
				Seed:       tc.seed,
			})
			require.NoError(t, err)

			err = ws.UpdateSecrets(tc.walletID, tc.password, func(w wallet.Wallet) error {
				switch w.Type() {
				case wallet.WalletTypeCollection:
					// Add 5 unused entries into the wallet, in addition to the 3 above
					uniqueEntries, _ := makeEntries(5)
					for _, e := range append(uniqueEntries, entries...) {
						err := w.(*collection.Wallet).AddEntry(e)
						require.NoError(t, err)
					}
				}
				return nil
			})
			require.NoError(t, err)

			walletAddrs, err := ws.GetAddresses(tc.walletID)
			require.NoError(t, err)

			b := &MockBlockchainer{}
			ut := &MockUnconfirmedTransactionPooler{}
			up := &MockUnspentPooler{}

			addrs := walletAddrs
			if len(tc.wp.Addresses) != 0 {
				addrs = tc.wp.Addresses
			}

			b.On("Head", matchDBTx).Return(tc.blockchainHead, tc.blockchainHeadErr)
			up.On("GetUnspentHashesOfAddrs", matchDBTx, addrs).Return(tc.getUnspentHashesOfAddrs, tc.getUnspentHashesOfAddrsErr)

			ut.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(unconfirmedForEachMockRun(t, tc.unconfirmedTxns, tc.uxOuts, tc.wp.IgnoreUnconfirmed))

			up.On("GetArray", matchDBTx, mock.MatchedBy(matchUxOutsAnyOrder(tc.getArrayInputs))).Return(tc.getArray, tc.getArrayErr)
			b.On("Unspent").Return(up)

			matchTxnIgnoreSigs := mock.MatchedBy(func(txn coin.Transaction) bool {
				switch tc.signed {
				case TxnSigned:
					if !txn.IsFullySigned() {
						return false
					}
					// Unset sigs for comparison to the unsigned txn
					txn.Sigs = make([]cipher.Sig, len(txn.Sigs))
					return reflect.DeepEqual(txn, *tc.txn)
				case TxnUnsigned:
					return reflect.DeepEqual(txn, *tc.txn)
				default:
					return false
				}
			})

			if tc.txn != nil {
				b.On("VerifySingleTxnSoftHardConstraints", matchDBTx, matchTxnIgnoreSigs, params.MainNetDistribution, params.UserVerifyTxn, tc.signed).Return(nil, nil, tc.verifyErr)
			}

			db, shutdown := prepareDB(t)
			defer shutdown()

			v := &Visor{
				db:          db,
				blockchain:  b,
				unconfirmed: ut,
				wallets:     ws,
				Config: Config{
					Distribution: params.MainNetDistribution,
				},
			}

			tf := mockTxnsFinder{}

			var txn *coin.Transaction
			var inputs []TransactionInput
			switch tc.signed {
			case TxnSigned:
				txn, inputs, err = v.WalletCreateTransactionSigned(tc.walletID, tc.password, tc.p, tc.wp, tf)
			case TxnUnsigned:
				txn, inputs, err = v.WalletCreateTransaction(tc.walletID, tc.p, tc.wp, tf)
			default:
				t.Fatal("invalid tc.signed value")
			}
			require.Equal(t, tc.err, err, "%v != %v", tc.err, err)
			if tc.err != nil {
				return
			}

			if tc.signed == TxnSigned {
				require.True(t, txn.IsFullySigned())
				// Unset sigs for comparison to the unsigned transaction
				txn.Sigs = make([]cipher.Sig, len(txn.Sigs))
			}

			require.Equal(t, tc.txn, txn)
			require.Equal(t, tc.inputs, inputs)
		})
	}
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

			_, _, err := v.WalletCreateTransaction("foo.wlt", tc.p, tc.wp, nil)
			require.Equal(t, tc.err, err)

			_, _, err = v.WalletCreateTransactionSigned("foo.wlt", nil, tc.p, tc.wp, nil)
			require.Equal(t, tc.err, err)

			if tc.err != nil {
				return
			}

			// Valid WalletCreateTransaction and WalletCreateTransactionSigned calls are tested in live integration tests
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
		getArray        coin.UxArray
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
			getArray: coin.UxArray{
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
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					{
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
			getArray: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[5],
						Address:        allAddrs[1],
					},
				},
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
			},
		},

		{
			name:              "uxouts specified, all uxouts are unconfirmed",
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
				coin.Transaction{
					In: hashes[5:8],
				},
			},
			err: ErrNoSpendableOutputs,
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
			getArray: coin.UxArray{
				coin.UxOut{
					Body: coin.UxBody{
						SrcTransaction: srcTxns[4],
						Address:        testutil.MakeAddress(),
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := testutil.PrepareDB(t)
			defer shutdown()

			unconfirmed := &MockUnconfirmedTransactionPooler{}
			bc := &MockBlockchainer{}
			unspent := &MockUnspentPooler{}
			require.Implements(t, (*blockdb.UnspentPooler)(nil), unspent)

			v := &Visor{
				unconfirmed: unconfirmed,
				blockchain:  bc,
				db:          db,
			}

			unconfirmed.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(unconfirmedForEachMockRun(t, tc.unconfirmedTxns, tc.uxOuts, tc.ignoreUnconfirmed))

			unspent.On("GetArray", matchDBTx, mock.MatchedBy(matchUxOutsAnyOrder(tc.getArrayInputs))).Return(tc.getArray, tc.getArrayErr)

			bc.On("Unspent").Return(unspent)

			var auxs coin.AddressUxOuts
			err := v.db.View("", func(tx *dbutil.Tx) error {
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

		forEachErr              error
		unconfirmedTxns         coin.Transactions
		getArrayInputs          []cipher.SHA256
		getArray                coin.UxArray
		getArrayErr             error
		getUnspentHashesOfAddrs blockdb.AddressHashes
	}{
		{
			name:           "ok",
			addrs:          allAddrs,
			getArrayInputs: hashes[0:4],
			getArray: coin.UxArray{
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
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				allAddrs[1]: hashes[0:2],
				allAddrs[3]: hashes[2:4],
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					{
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
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
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
			getArray: coin.UxArray{
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
			getUnspentHashesOfAddrs: blockdb.AddressHashes{
				allAddrs[1]: hashes[0:2],
				allAddrs[3]: hashes[2:4],
			},
			expectedAuxs: coin.AddressUxOuts{
				allAddrs[1]: []coin.UxOut{
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[5],
							Address:        allAddrs[1],
						},
					},
				},
				allAddrs[3]: []coin.UxOut{
					{
						Body: coin.UxBody{
							SrcTransaction: srcTxns[6],
							Address:        allAddrs[3],
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := testutil.PrepareDB(t)
			defer shutdown()

			unconfirmed := &MockUnconfirmedTransactionPooler{}
			bc := &MockBlockchainer{}
			unspent := &MockUnspentPooler{}
			require.Implements(t, (*blockdb.UnspentPooler)(nil), unspent)

			v := &Visor{
				unconfirmed: unconfirmed,
				blockchain:  bc,
				db:          db,
			}
			unspent.On("GetUnspentHashesOfAddrs", matchDBTx, tc.addrs).Return(tc.getUnspentHashesOfAddrs, nil)

			unconfirmed.On("ForEach", matchDBTx, mock.MatchedBy(func(f func(cipher.SHA256, UnconfirmedTransaction) error) bool {
				return true
			})).Return(tc.forEachErr).Run(unconfirmedForEachMockRun(t, tc.unconfirmedTxns, tc.getUnspentHashesOfAddrs.Flatten(), tc.ignoreUnconfirmed))

			unspent.On("GetArray", matchDBTx, mock.MatchedBy(matchUxOutsAnyOrder(tc.getArrayInputs))).Return(tc.getArray, tc.getArrayErr)

			bc.On("Unspent").Return(unspent)

			var auxs coin.AddressUxOuts
			err := v.db.View("", func(tx *dbutil.Tx) error {
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

var matchDBTx = mock.MatchedBy(func(tx *dbutil.Tx) bool {
	return true
})

func matchUxOutsAnyOrder(inputs []cipher.SHA256) func(args []cipher.SHA256) bool {
	return func(args []cipher.SHA256) bool {
		// Compares two []coin.UxOuts for equality, ignoring the order of elements in the slice
		if len(args) != len(inputs) {
			return false
		}

		x := make([]cipher.SHA256, len(inputs))
		copy(x, inputs)
		y := make([]cipher.SHA256, len(args))
		copy(y, args)

		sort.Slice(x, func(a, b int) bool {
			return bytes.Compare(x[a][:], x[b][:]) < 0
		})
		sort.Slice(y, func(a, b int) bool {
			return bytes.Compare(y[a][:], y[b][:]) < 0
		})

		return reflect.DeepEqual(x, y)
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

// unconfirmedForEachMockRun simulates the Unconfirmed.ForEach callback method
func unconfirmedForEachMockRun(t *testing.T, unconfirmedTxns []coin.Transaction, uxOuts []cipher.SHA256, ignoreUnconfirmed bool) func(args mock.Arguments) {
	return func(args mock.Arguments) {
		fn := args.Get(1).(func(cipher.SHA256, UnconfirmedTransaction) error)
		for _, u := range unconfirmedTxns {
			err := fn(u.Hash(), UnconfirmedTransaction{
				Transaction: u,
			})

			// If any of the input hashes are in an unconfirmed transaction,
			// the callback handler should have returned ErrSpendingUnconfirmed
			// unless IgnoreUnconfirmed is true
			hasUnconfirmedHash := hashesIntersect(u.In, uxOuts)

			if hasUnconfirmedHash {
				if ignoreUnconfirmed {
					require.NoError(t, err)
				} else {
					require.Equal(t, ErrSpendingUnconfirmed, err)
				}
			} else {
				require.NoError(t, err)
			}

		}
	}
}

type mockTxnsFinder map[cipher.Addresser]bool

func (mb mockTxnsFinder) AddressesActivity(addrs []cipher.Addresser) ([]bool, error) {
	if len(addrs) == 0 {
		return nil, nil
	}
	active := make([]bool, len(addrs))
	for i, addr := range addrs {
		active[i] = mb[addr]
	}
	return active, nil
}
