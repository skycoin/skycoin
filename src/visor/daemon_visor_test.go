package visor

/*
These are tests that used to be in daemon/visor_test.go,
but belong in package visor instead.

They have been moved here without checking if they duplicate any
existing test in visor_test.go.

It is assumed that these tests may provide coverage not present in visor_test.go

They could be merged into visor_test.go, but for simplicity they were only moved here
*/

import (
	"errors"
	"testing"

	"github.com/skycoin/skycoin/src/transaction"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

func setupSimpleVisor(t *testing.T, db *dbutil.DB, bc *Blockchain) *Visor {
	cfg := NewConfig()

	pool, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	return &Visor{
		Config:      cfg,
		unconfirmed: pool,
		blockchain:  bc,
		db:          db,
	}
}

func TestVerifyTransactionInvalidFee(t *testing.T) {
	// Test that a soft constraint is enforced
	// Full verification tests are in visor/blockchain_verify_test.go
	db, close := prepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var f uint64
	addr := testutil.MakeAddress()

	txn := CreateGenesisSpendTransaction(t, db, bc, addr, coins, hours, f)

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	_, softErr, err := v.InjectForeignTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, transaction.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)
}

func TestVerifyTransactionInvalidSignature(t *testing.T) {
	// Test that a hard constraint is enforced
	// Full verification tests are in visor/blockchain_verify_test.go
	db, close := prepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours uint64
	var fee uint64
	addr := testutil.MakeAddress()

	txn := CreateGenesisSpendTransaction(t, db, bc, addr, coins, hours, fee)

	// Invalidate signatures
	txn.Sigs = nil

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	_, softErr, err := v.InjectForeignTransaction(txn)
	require.Nil(t, softErr)
	testutil.RequireError(t, err, transaction.NewErrTxnViolatesHardConstraint(errors.New("Invalid number of signatures")).Error())
}

func TestInjectValidTransaction(t *testing.T) {
	db, close := prepareDB(t)
	defer close()

	_, s := cipher.GenerateKeyPair()
	// Setup blockchain
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours uint64
	var fee uint64
	addr := testutil.MakeAddress()

	txn := CreateGenesisSpendTransaction(t, db, bc, addr, coins, hours, fee)

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	// The unconfirmed pool should be empty
	txns, err := v.GetAllUnconfirmedTransactions()
	require.NoError(t, err)
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectForeignTransaction(txn)
	require.Nil(t, softErr)
	require.NoError(t, err)

	// The transaction should appear in the unconfirmed pool
	txns, err = v.GetAllUnconfirmedTransactions()
	require.NoError(t, err)
	require.Len(t, txns, 1)
	require.Equal(t, txns[0].Transaction, txn)
}

func TestInjectTransactionSoftViolationNoFee(t *testing.T) {
	db, close := prepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var f uint64
	addr := testutil.MakeAddress()

	txn := CreateGenesisSpendTransaction(t, db, bc, addr, coins, hours, f)

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	// The unconfirmed pool should be empty
	txns, err := v.GetAllUnconfirmedTransactions()
	require.NoError(t, err)
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectForeignTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, transaction.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)

	// The transaction should appear in the unconfirmed pool
	txns, err = v.GetAllUnconfirmedTransactions()
	require.NoError(t, err)
	require.Len(t, txns, 1)
}
