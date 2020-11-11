package visor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/testutil"
	_require "github.com/skycoin/skycoin/src/testutil/require"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/timeutil"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

const (
	blockchainPubkeyStr = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
)

func prepareDB(t *testing.T) (*dbutil.DB, func()) {
	db, shutdown := testutil.PrepareDB(t)

	err := CreateBuckets(db)
	if err != nil {
		shutdown()
		t.Fatalf("CreateBuckets failed: %v", err)
	}

	return db, shutdown
}

func readAll(t *testing.T, f string) []byte {
	fi, err := os.Open(f)
	require.NoError(t, err)
	defer fi.Close()

	b, err := ioutil.ReadAll(fi)
	require.NoError(t, err)

	return b
}

func mustParsePubkey(t *testing.T) cipher.PubKey {
	// Parse the blockchain pubkey associated with this corrupted test db
	t.Helper()
	pubkey, err := cipher.PubKeyFromHex(blockchainPubkeyStr)
	require.NoError(t, err)
	return pubkey
}

func writeDBFile(t *testing.T, badDBFile string, badDBData []byte) {
	t.Logf("Writing the original bad db file back to %s", badDBFile)
	fi, err := os.OpenFile(badDBFile, os.O_WRONLY, 0600)
	require.NoError(t, err)
	defer fi.Close()

	_, err = io.Copy(fi, bytes.NewBuffer(badDBData))
	require.NoError(t, err)
}

func findCorruptDBFiles(t *testing.T, badDBFile string) []string {
	corruptFiles, err := filepath.Glob(badDBFile + ".corrupt.*")
	require.NoError(t, err)
	return corruptFiles
}

func removeCorruptDBFiles(t *testing.T, badDBFile string) {
	corruptFiles := findCorruptDBFiles(t, badDBFile)
	for _, m := range corruptFiles {
		err := os.Remove(m)
		require.NoError(t, err)
	}
}

func addGenesisBlockToVisor(t *testing.T, vs *Visor) *coin.SignedBlock {
	// create genesis block
	gb, err := coin.NewGenesisBlock(genAddress, genCoins, genTime)
	require.NoError(t, err)
	gbSig := cipher.MustSignHash(gb.HashHeader(), genSecret)
	vs.Config.GenesisSignature = gbSig

	sb := coin.SignedBlock{
		Block: *gb,
		Sig:   gbSig,
	}

	// add genesis block to blockchain
	err = vs.db.Update("", func(tx *dbutil.Tx) error {
		return vs.executeSignedBlock(tx, sb)
	})
	require.NoError(t, err)

	return &sb
}

func TestErrMissingSignatureRecreateDB(t *testing.T) {
	badDBFile := "./testdata/data.db.nosig" // about 8MB size
	badDBData := readAll(t, badDBFile)

	pubkey := mustParsePubkey(t)

	// Remove any existing corrupt db files from testdata
	removeCorruptDBFiles(t, badDBFile)
	corruptFiles := findCorruptDBFiles(t, badDBFile)
	require.Len(t, corruptFiles, 0)

	// Cleanup
	defer func() {
		// Write the bad db data back to badDBFile
		writeDBFile(t, badDBFile, badDBData)
		// Remove leftover corrupt db copies
		removeCorruptDBFiles(t, badDBFile)
	}()

	// Make sure that the database file causes ErrMissingSignature error
	t.Logf("Checking that %s is a corrupted database", badDBFile)
	func() {
		db, err := OpenDB(badDBFile, false)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		bc, err := NewBlockchain(db, BlockchainConfig{
			Pubkey:      pubkey,
			Arbitrating: false,
		})
		require.NoError(t, err)

		// err = db.View("", func(tx *dbutil.Tx) error {
		f := func(tx *dbutil.Tx, b *coin.SignedBlock) error {
			return bc.VerifySignature(b)
		}

		err = bc.WalkChain(BlockchainVerifyTheadNum, f, nil)

		require.Error(t, err)
		require.IsType(t, blockdb.ErrMissingSignature{}, err)
	}()

	// Loading this invalid db should cause ResetCorruptDB() to recreate the db
	t.Logf("Loading the corrupted db from %s", badDBFile)
	badDB, err := OpenDB(badDBFile, false)
	require.NoError(t, err)
	require.NotNil(t, badDB)
	require.NotEmpty(t, badDB.Path())
	t.Logf("badDB.Path() == %s", badDB.Path())

	db, err := ResetCorruptDB(badDB, pubkey, nil)
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)

	require.NotNil(t, db)

	// A corrupted database file should exist
	corruptFiles = findCorruptDBFiles(t, badDBFile)
	require.Len(t, corruptFiles, 1)

	// A new db should be written in place of the old bad db, and not be corrupted
	t.Logf("Checking that the new db file is valid")
	func() {
		db, err := OpenDB(badDBFile, false)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		// The new db is not corrupted and loads without error
		bc, err := NewBlockchain(db, BlockchainConfig{
			Pubkey:      pubkey,
			Arbitrating: false,
		})
		require.NoError(t, err)
		require.NotNil(t, bc)
	}()
}

func TestHistorydbVerifier(t *testing.T) {
	tt := []struct {
		name      string
		dbPath    string
		expectErr error
	}{
		{
			name:   "db is ok",
			dbPath: "./testdata/data.db.ok",
		},
		{
			name:      "missing transaction",
			dbPath:    "./testdata/data.db.notxn",
			expectErr: historydb.NewErrHistoryDBCorrupted(errors.New("HistoryDB.Verify: transaction 98db7eb30e13853d3dd93d5d8b4061596d5d288b6f8b92c4d43c46c6599f67fb does not exist in historydb")),
		},
		{
			name:      "missing uxout",
			dbPath:    "./testdata/data.db.nouxout",
			expectErr: historydb.NewErrHistoryDBCorrupted(errors.New("HistoryDB.Verify: transaction (input|output) 2f87d77c2a7d00b547db1af50e0ba04bafc5b05711e4939e9ec2640a21127dc0 does not exist in historydb")),
		},
		{
			name:      "missing addr transaction index",
			dbPath:    "./testdata/data.db.no-addr-txn-index",
			expectErr: historydb.NewErrHistoryDBCorrupted(errors.New(`HistoryDB.Verify: index of address transaction \[2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF:98db7eb30e13853d3dd93d5d8b4061596d5d288b6f8b92c4d43c46c6599f67fb\] does not exist in historydb`)),
		},
		{
			name:      "missing addr uxout index",
			dbPath:    "./testdata/data.db.no-addr-uxout-index",
			expectErr: historydb.NewErrHistoryDBCorrupted(errors.New(`HistoryDB.Verify: index of address uxout \[2fGC7kwAM9yZyEF1QqBqp8uo9RUsF6ENGJF:2f87d77c2a7d00b547db1af50e0ba04bafc5b05711e4939e9ec2640a21127dc0\] does not exist in historydb`)),
		},
	}

	pubKeyStr := "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
	pubkey := cipher.MustPubKeyFromHex(pubKeyStr)
	history := historydb.New()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db, err := OpenDB(tc.dbPath, true)
			require.NoError(t, err)
			bc, err := NewBlockchain(db, BlockchainConfig{
				Pubkey: pubkey,
			})
			require.NoError(t, err)

			indexesMap := historydb.NewIndexesMap()
			f := func(tx *dbutil.Tx, b *coin.SignedBlock) error {
				return history.Verify(tx, b, indexesMap)
			}

			err = bc.WalkChain(2, f, nil)
			if tc.expectErr == nil {
				require.Nil(t, err)
				return
			}

			// Confirms that the error type is matched
			require.IsType(t, tc.expectErr, err)
			// Confirms the error message is matched
			require.Regexp(t, tc.expectErr.Error(), err.Error())
		})
	}

}

func TestVisorCreateBlock(t *testing.T) {
	when := uint64(time.Now().UTC().Unix())

	db, shutdown := prepareDB(t)
	defer shutdown()

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey: genPublic,
	})

	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	his := historydb.New()

	cfg := NewConfig()
	cfg.IsBlockPublisher = false
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		unconfirmed: unconfirmed,
		blockchain:  bc,
		db:          db,
		history:     his,
	}

	// CreateBlock panics if called when not a block publisher
	_require.PanicsWithLogMessage(t, "Only a block publisher node can create blocks", func() {
		err := db.Update("", func(tx *dbutil.Tx) error {
			_, err := v.createBlock(tx, when)
			return err
		})
		require.NoError(t, err)
	})

	v.Config.IsBlockPublisher = true
	v.Config.BlockchainSeckey = genSecret

	addGenesisBlockToVisor(t, v)
	var gb *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		gb, err = v.blockchain.GetGenesisBlock(tx)
		return err
	})
	require.NoError(t, err)
	require.NotNil(t, gb)

	// If no transactions in the unconfirmed pool, return an error
	err = db.Update("", func(tx *dbutil.Tx) error {
		_, err = v.createBlock(tx, when)
		testutil.RequireError(t, err, "No transactions")
		return nil
	})
	require.NoError(t, err)

	// Create enough unspent outputs to create all of these transactions
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	nUnspents := 100
	txn := makeUnspentsTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, nUnspents, params.UserVerifyTxn.MaxDropletPrecision)

	var known bool
	var softErr *ErrTxnViolatesSoftConstraint
	err = db.Update("", func(tx *dbutil.Tx) error {
		var err error
		known, softErr, err = unconfirmed.InjectTransaction(tx, bc, txn, params.MainNetDistribution, v.Config.UnconfirmedVerifyTxn)
		return err
	})
	require.NoError(t, err)
	require.False(t, known)
	require.Nil(t, softErr)

	v.Config.MaxBlockTransactionsSize, err = txn.Size()
	require.NoError(t, err)
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))

	var length uint64
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		length, err = unconfirmed.Len(tx)
		return err
	})
	require.NoError(t, err)

	require.Equal(t, uint64(0), length)
	v.Config.MaxBlockTransactionsSize = 1024 * 4

	// Create various transactions and add them to unconfirmed pool
	uxs = coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])
	var coins uint64 = 9e6
	var f uint64 = 10
	toAddr := testutil.MakeAddress()

	// Add more transactions than is allowed in a block, to verify truncation
	var txns coin.Transactions
	var i int
	truncatedTxns, err := txns.TruncateBytesTo(v.Config.MaxBlockTransactionsSize)
	require.NoError(t, err)
	for len(txns) == len(truncatedTxns) {
		tx := makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins, f)
		txns = append(txns, tx)
		i++
		truncatedTxns, err = txns.TruncateBytesTo(v.Config.MaxBlockTransactionsSize)
		require.NoError(t, err)
	}
	require.NotEqual(t, 0, len(txns))

	// Use different f sizes to verify f ordering
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins, f*5))
	i++
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins, f*10))
	i++

	// Use invalid decimal places to verify decimal place filtering.
	// The fs are set higher to ensure that they are not filtered due to truncating with a low f
	// Spending 9.1 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1e5, f*20))
	i++
	// Spending 9.01 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1e4, f*30))
	i++
	// Spending 9.0001 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1e3, f*40))
	i++
	// Spending 9.0001 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1e2, f*50))
	i++
	// Spending 9.00001 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1e1, f*60))
	i++
	// Spending 9.000001 SKY
	txns = append(txns, makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins+1, f*70))
	// i++

	// Confirm that at least one transaction has an invalid decimal output
	foundInvalidCoins := false
	for _, txn := range txns {
		for _, o := range txn.Out {
			if err := params.DropletPrecisionCheck(v.Config.UnconfirmedVerifyTxn.MaxDropletPrecision, o.Coins); err != nil {
				foundInvalidCoins = true
				break
			}
		}
	}
	require.True(t, foundInvalidCoins)

	// Inject transactions into the unconfirmed pool
	for i, txn := range txns {
		var known bool
		var softErr *ErrTxnViolatesSoftConstraint
		err = db.Update("", func(tx *dbutil.Tx) error {
			var err error
			known, softErr, err = unconfirmed.InjectTransaction(tx, bc, txn, params.MainNetDistribution, v.Config.UnconfirmedVerifyTxn)
			return err
		})
		require.False(t, known)
		require.NoError(t, err)

		// The last 3 transactions will have a soft constraint violation for too many decimal places,
		// but would still be injected into the pool
		if i < len(txns)-3 {
			require.Nil(t, softErr)
		} else {
			testutil.RequireError(t, softErr, "Transaction violates soft constraint: invalid amount, too many decimal places")
		}
	}

	// Make sure all transactions were injected
	var allInjectedTxns []coin.Transaction
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		allInjectedTxns, err = unconfirmed.AllRawTransactions(tx)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, len(txns), len(allInjectedTxns))

	err = db.Update("", func(tx *dbutil.Tx) error {
		var err error
		sb, err = v.createBlock(tx, when+100)
		return err
	})
	require.NoError(t, err)
	require.Equal(t, when+100, sb.Block.Head.Time)

	blockTxns := sb.Block.Body.Transactions
	require.NotEqual(t, len(txns), len(blockTxns), "Transactions should be truncated")
	require.Equal(t, 18, len(blockTxns))

	// Check fee ordering
	err = db.View("", func(tx *dbutil.Tx) error {
		inUxs, err := v.blockchain.Unspent().GetArray(tx, blockTxns[0].In)
		require.NoError(t, err)
		prevFee, err := fee.TransactionFee(&blockTxns[0], sb.Head.Time, inUxs)
		require.NoError(t, err)

		for i := 1; i < len(blockTxns); i++ {
			inUxs, err := v.blockchain.Unspent().GetArray(tx, blockTxns[i].In)
			require.NoError(t, err)
			f, err := fee.TransactionFee(&blockTxns[i], sb.Head.Time, inUxs)
			require.NoError(t, err)
			require.True(t, f <= prevFee)
			prevFee = f
		}

		return nil
	})

	require.NoError(t, err)

	// Check that decimal rules are enforced
	for i, txn := range blockTxns {
		for j, o := range txn.Out {
			err := params.DropletPrecisionCheck(v.Config.CreateBlockVerifyTxn.MaxDropletPrecision, o.Coins)
			require.NoError(t, err, "txout %d.%d coins=%d", i, j, o.Coins)
		}
	}
}

func TestVisorInjectTransaction(t *testing.T) {
	when := uint64(time.Now().UTC().Unix())

	db, shutdown := prepareDB(t)
	defer shutdown()

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey: genPublic,
	})
	require.NoError(t, err)

	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	his := historydb.New()

	cfg := NewConfig()
	cfg.IsBlockPublisher = false
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		unconfirmed: unconfirmed,
		blockchain:  bc,
		db:          db,
		history:     his,
	}

	// CreateBlock panics if called when not a block publisher
	_require.PanicsWithLogMessage(t, "Only a block publisher node can create blocks", func() {
		err := db.Update("", func(tx *dbutil.Tx) error {
			_, err := v.createBlock(tx, when)
			return err
		})
		require.NoError(t, err)
	})

	v.Config.IsBlockPublisher = true
	v.Config.BlockchainSeckey = genSecret

	addGenesisBlockToVisor(t, v)

	var gb *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		gb, err = v.blockchain.GetGenesisBlock(tx)
		return err
	})
	require.NoError(t, err)
	require.NotNil(t, gb)

	// If no transactions in the unconfirmed pool, return an error
	err = db.Update("", func(tx *dbutil.Tx) error {
		_, err := v.createBlock(tx, when)
		return err
	})
	testutil.RequireError(t, err, "No transactions")

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	toAddr := testutil.MakeAddress()
	var coins uint64 = 10e6

	// Create a transaction with valid decimal places
	txn := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectForeignTransaction(txn)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)

	// Execute a block to clear this transaction from the pool
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 2, len(sb.Body.Transactions[0].Out))

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)

		length, err = bc.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), length)

		return nil
	})
	require.NoError(t, err)

	uxs = coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])

	// Check transactions with overflowing output coins fail
	txn = makeOverflowCoinsSpendTxn(t, coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr)
	_, softErr, err = v.InjectForeignTransaction(txn)
	require.IsType(t, ErrTxnViolatesHardConstraint{}, err)
	testutil.RequireError(t, err.(ErrTxnViolatesHardConstraint).Err, "Output coins overflow")
	require.Nil(t, softErr)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)
		return nil
	})
	require.NoError(t, err)

	// Check transactions with overflowing output hours fail
	// It should not be injected; when injecting a txn, the overflowing output hours is treated
	// as a hard constraint. It is only a soft constraint when the txn is included in a signed block.
	txn = makeOverflowHoursSpendTxn(t, coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr)
	_, softErr, err = v.InjectForeignTransaction(txn)
	require.Nil(t, softErr)
	require.IsType(t, ErrTxnViolatesHardConstraint{}, err)
	testutil.RequireError(t, err.(ErrTxnViolatesHardConstraint).Err, "Transaction output hours overflow")

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)
		return nil
	})
	require.NoError(t, err)

	// Create a transaction with invalid decimal places
	// It's still injected, because this is considered a soft error
	invalidCoins := coins + (params.UserVerifyTxn.MaxDropletDivisor() / 10)
	txn = makeSpendTxn(t, uxs, []cipher.SecKey{genSecret, genSecret}, toAddr, invalidCoins)
	_, softErr, err = v.InjectForeignTransaction(txn)
	require.NoError(t, err)
	testutil.RequireError(t, softErr.Err, params.ErrInvalidDecimals.Error())

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), length)
		return nil
	})
	require.NoError(t, err)

	// Create a transaction with null address output
	uxs = coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	txn = makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	txn.Out[0].Address = cipher.Address{}
	known, _, _, err = v.InjectUserTransaction(txn)
	require.False(t, known)
	require.IsType(t, ErrTxnViolatesUserConstraint{}, err)
	testutil.RequireError(t, err, "Transaction violates user constraint: Transaction output is sent to the null address")
}

func makeOverflowCoinsSpendTxn(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address) coin.Transaction {
	spendTxn := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		err := spendTxn.PushInput(ux.Hash())
		require.NoError(t, err)
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	hours := totalHours / 12

	// These two outputs' coins added up will overflow
	err := spendTxn.PushOutput(toAddr, 18446744073709551000, hours)
	require.NoError(t, err)
	err = spendTxn.PushOutput(toAddr, totalCoins, hours)
	require.NoError(t, err)

	spendTxn.SignInputs(keys)
	err = spendTxn.UpdateHeader()
	require.NoError(t, err)
	return spendTxn
}

func makeOverflowHoursSpendTxn(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address) coin.Transaction {
	spendTxn := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		err := spendTxn.PushInput(ux.Hash())
		require.NoError(t, err)
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	hours := totalHours / 12

	// These two outputs' hours added up will overflow
	err := spendTxn.PushOutput(toAddr, totalCoins/2, 18446744073709551615)
	require.NoError(t, err)
	err = spendTxn.PushOutput(toAddr, totalCoins-totalCoins/2, hours)
	require.NoError(t, err)

	spendTxn.SignInputs(keys)
	err = spendTxn.UpdateHeader()
	require.NoError(t, err)
	return spendTxn
}

func makeTestData(t *testing.T, n int) ([]historydb.Transaction, []coin.SignedBlock, []UnconfirmedTransaction, uint64) { //nolint:unparam
	var txns []historydb.Transaction
	var blocks []coin.SignedBlock
	var uncfmTxns []UnconfirmedTransaction
	for i := uint64(0); i < uint64(n); i++ {
		tm := time.Now().UTC().Unix() + int64(i)*int64(time.Second)
		txns = append(txns, historydb.Transaction{
			BlockSeq: i,
			Txn: coin.Transaction{
				InnerHash: testutil.RandSHA256(t),
			},
		})

		blocks = append(blocks, coin.SignedBlock{
			Block: coin.Block{
				Head: coin.BlockHeader{
					BkSeq: i,
					Time:  uint64(tm),
				},
			},
		})

		uncfmTxns = append(uncfmTxns, UnconfirmedTransaction{
			Transaction: coin.Transaction{
				InnerHash: testutil.RandSHA256(t),
			},
			Received: time.Now().UTC().Unix() + int64(n)*int64(time.Second),
		})
	}

	return txns, blocks, uncfmTxns, uint64(n)
}

func makeUncfmUxs(txns []UnconfirmedTransaction) coin.UxArray {
	var uxs coin.UxArray
	for i := range txns {
		uxs = append(uxs, coin.UxOut{
			Head: coin.UxHead{
				Time: uint64(txns[i].Received),
			},
			Body: coin.UxBody{
				SrcTransaction: txns[i].Transaction.Hash(),
			},
		})
	}
	return uxs
}

type txnsAndUncfmTxns struct {
	Txns      []historydb.Transaction
	UncfmTxns []UnconfirmedTransaction
}
type expectTxnResult struct {
	txns      []Transaction
	uncfmTxns []Transaction
	err       error
}

func TestGetTransactions(t *testing.T) {
	// Generates test data
	txns, blocks, uncfmTxns, headSeq := makeTestData(t, 10)
	// Generates []Transaction
	var lTxns []Transaction
	for i := range txns {
		height := headSeq - txns[i].BlockSeq + 1
		lTxns = append(lTxns, Transaction{
			Transaction: txns[i].Txn,
			Status:      NewConfirmedTransactionStatus(height, txns[i].BlockSeq),
			Time:        blocks[i].Time(),
		})
	}

	// Generate unconfirmed []Transaction
	var luncfmTxns []Transaction
	for i, txn := range uncfmTxns {
		luncfmTxns = append(luncfmTxns, Transaction{
			Transaction: uncfmTxns[i].Transaction,
			Status:      NewUnconfirmedTransactionStatus(),
			Time:        uint64(timeutil.NanoToTime(txn.Received).Unix()),
		})
	}

	// Generates addresses
	var addrs []cipher.Address
	for i := 0; i < 10; i++ {
		addrs = append(addrs, testutil.MakeAddress())
	}

	tt := []struct {
		name      string
		addrTxns  map[cipher.Address]txnsAndUncfmTxns
		blocks    []coin.SignedBlock
		bcHeadSeq uint64
		filters   []TxFilter
		expect    expectTxnResult
	}{
		{
			"addrFilter=1 addr=1 txns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=0 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=0 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=2 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=1 txns=2 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=0 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=0 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=1 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=2 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: nil,
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=2 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=2 unconfirmedTxns=3",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: luncfmTxns[:3],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=3 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: nil,
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:3],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=3 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:3],
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=3 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:3],
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 addr=2 txns=3 unconfirmedTxns=3",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: uncfmTxns[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
			},
			expectTxnResult{
				txns:      lTxns[:3],
				uncfmTxns: luncfmTxns[:3],
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 addrsFilter=1 confirmed=false txns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxns=2 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmedTxFilter=1 confirmed=false addr=1 txns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txns=1 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:1],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txns=2 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txns=2 unconfirmedTxns=3",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:3],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2/1 txns=2 unconfirmedTxns=3",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[:2],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2/2 txns=2 unconfirmedTxns=3",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: uncfmTxns[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[1:2]),
				NewConfirmedTxFilter(false),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: luncfmTxns[2:3],
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txns=0 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      nil,
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      nil,
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txns=1 unconfirmedTxns=0",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txns=1 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:1],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:1],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txns=2 unconfirmedTxns=1",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txns=2 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2/1 txns=3 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: uncfmTxns[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:1]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:2],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2/2 txns=3 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: uncfmTxns[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[1:2]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[2:3],
				uncfmTxns: nil,
				err:       nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2 txns=3 unconfirmedTxns=2",
			map[cipher.Address]txnsAndUncfmTxns{
				addrs[0]: txnsAndUncfmTxns{
					Txns:      txns[:2],
					UncfmTxns: uncfmTxns[:1],
				},
				addrs[1]: txnsAndUncfmTxns{
					Txns:      txns[2:3],
					UncfmTxns: uncfmTxns[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				NewAddrsFilter(addrs[:2]),
				NewConfirmedTxFilter(true),
			},
			expectTxnResult{
				txns:      lTxns[:3],
				uncfmTxns: nil,
				err:       nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			matchDBTx := mock.MatchedBy(func(tx *dbutil.Tx) bool {
				return true
			})

			his := newHistoryerMock2()
			uncfmTxnPool := NewUnconfirmedTransactionPoolerMock2()
			forEachFunc := mock.MatchedBy(func(f func(hash cipher.SHA256, txn UnconfirmedTransaction) error) bool {
				for i, txn := range uncfmTxns {
					if err := f(txn.Transaction.Hash(), uncfmTxns[i]); err != nil {
						return false
					}
				}
				return true
			})

			uncfmTxnPool.On("ForEach", matchDBTx, forEachFunc).Return(nil)
			for i, txn := range uncfmTxns {
				uncfmTxnPool.On("Get", matchDBTx, txn.Transaction.Hash()).Return(&uncfmTxns[i], nil)
			}

			for addr, txns := range tc.addrTxns {
				for i, txn := range txns.Txns {
					his.On("GetTransaction", matchDBTx, txn.Hash()).Return(&txns.Txns[i], nil)
				}

				his.txns = append(his.txns, txns.Txns...)
				uncfmTxnPool.On("GetUnspentsOfAddr", matchDBTx, addr).Return(makeUncfmUxs(txns.UncfmTxns), nil)
				for i, uncfmTxn := range txns.UncfmTxns {
					uncfmTxnPool.On("Get", matchDBTx, uncfmTxn.Transaction.Hash()).Return(&txns.UncfmTxns[i], nil)
				}
				uncfmTxnPool.txns = append(uncfmTxnPool.txns, txns.UncfmTxns...)
			}

			var hisHashes []cipher.SHA256
			var hisAddrs []cipher.Address
			for _, flt := range tc.filters {
				switch f := flt.(type) {
				case AddrsFilter:
					hisAddrs = f.Addrs
					for _, a := range f.Addrs {
						txns, ok := tc.addrTxns[a]
						require.True(t, ok)
						for _, txn := range txns.Txns {
							hisHashes = append(hisHashes, txn.Hash())
						}
					}
				}
			}

			his.On("GetTransactionHashesForAddresses", matchDBTx, hisAddrs).Return(hisHashes, nil)

			bc := &MockBlockchainer{}
			for i, b := range tc.blocks {
				bc.On("GetSignedBlockBySeq", matchDBTx, b.Seq()).Return(&tc.blocks[i], nil)
			}

			bc.On("HeadSeq", matchDBTx).Return(tc.bcHeadSeq, true, nil)

			db, shutdown := prepareDB(t)
			defer shutdown()

			txnModel := transactionModel{
				history:     his,
				unconfirmed: uncfmTxnPool,
				blockchain:  bc,
			}

			v := &Visor{
				db:          db,
				history:     his,
				unconfirmed: uncfmTxnPool,
				blockchain:  bc,
				txns:        &txnModel,
			}

			retTxns, _, err := v.GetTransactions(tc.filters, AscOrder, nil)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Len(t, retTxns, len(tc.expect.txns)+len(tc.expect.uncfmTxns))

			// Splits confirmed and unconfirmed txns in returned transactions
			uncfmTxnMap := make(map[cipher.SHA256]Transaction)
			txnMap := make(map[cipher.SHA256]Transaction)
			for i, txn := range retTxns {
				h := txn.Transaction.Hash()
				if retTxns[i].Status.Confirmed {
					txnMap[h] = retTxns[i]
				} else {
					uncfmTxnMap[h] = retTxns[i]
				}
			}

			// Confirms that all expected confirmed transactions must be in the txnMap
			for _, txn := range tc.expect.txns {
				h := txn.Transaction.Hash()
				retTxn, ok := txnMap[h]
				require.True(t, ok)
				require.Equal(t, txn, retTxn)
			}

			// Confirms that all expected unconfirmed transactions must be in the uncfmTxnMap
			for _, txn := range tc.expect.uncfmTxns {
				h := txn.Transaction.Hash()
				retTxn, ok := uncfmTxnMap[h]
				require.True(t, ok)
				require.Equal(t, txn, retTxn)
			}
		})
	}
}

func TestRefreshUnconfirmed(t *testing.T) {
	db, shutdown := prepareDB(t)
	defer shutdown()

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey: genPublic,
	})
	require.NoError(t, err)

	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	his := historydb.New()

	cfg := NewConfig()
	cfg.IsBlockPublisher = true
	cfg.BlockchainSeckey = genSecret
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		unconfirmed: unconfirmed,
		blockchain:  bc,
		db:          db,
		history:     his,
	}

	addGenesisBlockToVisor(t, v)
	var gb *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		gb, err = v.blockchain.GetGenesisBlock(tx)
		return err
	})
	require.NoError(t, err)
	require.NotNil(t, gb)

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	toAddr := testutil.MakeAddress()
	var coins uint64 = 10e6

	// Create a valid transaction that will remain valid
	validTxn := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectForeignTransaction(validTxn)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), length)
		return nil
	})
	require.NoError(t, err)

	// Create a transaction with invalid decimal places
	// It's still injected, because this is considered a soft error
	// This transaction will stay invalid on refresh
	invalidCoins := coins + (params.UserVerifyTxn.MaxDropletDivisor() / 10)
	alwaysInvalidTxn := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, toAddr, invalidCoins)
	_, softErr, err = v.InjectForeignTransaction(alwaysInvalidTxn)
	require.NoError(t, err)
	testutil.RequireError(t, softErr.Err, params.ErrInvalidDecimals.Error())

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), length)
		return nil
	})
	require.NoError(t, err)

	// Create a transaction that exceeds UnconfirmedVerifyTxn.MaxTransactionSize
	// It's still injected, because this is considered a soft error
	// This transaction will become valid on refresh (by increasing UnconfirmedVerifyTxn.MaxTransactionSize)
	originalMaxUnconfirmedTxnSize := v.Config.UnconfirmedVerifyTxn.MaxTransactionSize
	v.Config.UnconfirmedVerifyTxn.MaxTransactionSize = 1
	sometimesInvalidTxn := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	_, softErr, err = v.InjectForeignTransaction(sometimesInvalidTxn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	testutil.RequireError(t, softErr.Err, ErrTxnExceedsMaxBlockSize.Error())

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(3), length)
		return nil
	})
	require.NoError(t, err)

	// The first txn remains valid,
	// the second txn remains invalid,
	// the third txn becomes valid
	v.Config.UnconfirmedVerifyTxn.MaxTransactionSize = originalMaxUnconfirmedTxnSize
	hashes, err := v.RefreshUnconfirmed()
	require.NoError(t, err)
	require.Equal(t, []cipher.SHA256{sometimesInvalidTxn.Hash()}, hashes)

	// Reduce the max block size to affirm that the valid transaction becomes invalid
	// The first txn becomes invalid,
	// the second txn remains invalid,
	// the third txn becomes invalid again
	v.Config.UnconfirmedVerifyTxn.MaxTransactionSize = 1
	hashes, err = v.RefreshUnconfirmed()
	require.NoError(t, err)
	require.Nil(t, hashes)

	// Restore the max block size to affirm the expected transaction validity shifts
	// The first txn was valid, became invalid, and is now valid again
	// The second txn was always invalid
	// The third txn was invalid, became valid, became invalid, and is now valid again
	v.Config.UnconfirmedVerifyTxn.MaxTransactionSize = originalMaxUnconfirmedTxnSize
	hashes, err = v.RefreshUnconfirmed()
	require.NoError(t, err)

	// Sort hashes for deterministic comparison
	expectedHashes := []cipher.SHA256{validTxn.Hash(), sometimesInvalidTxn.Hash()}
	sort.Slice(expectedHashes, func(i, j int) bool {
		return bytes.Compare(expectedHashes[i][:], expectedHashes[j][:]) < 0
	})
	sort.Slice(hashes, func(i, j int) bool {
		return bytes.Compare(hashes[i][:], hashes[j][:]) < 0
	})
	require.Equal(t, expectedHashes, hashes)
}

func TestRemoveInvalidUnconfirmedDoubleSpendArbitrating(t *testing.T) {
	db, shutdown := prepareDB(t)
	defer shutdown()

	bc, err := NewBlockchain(db, BlockchainConfig{
		Pubkey:      genPublic,
		Arbitrating: true,
	})
	require.NoError(t, err)

	unconfirmed, err := NewUnconfirmedTransactionPool(db)
	require.NoError(t, err)

	his := historydb.New()

	cfg := NewConfig()
	cfg.IsBlockPublisher = true
	cfg.Arbitrating = true
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress
	cfg.BlockchainSeckey = genSecret

	v := &Visor{
		Config:      cfg,
		unconfirmed: unconfirmed,
		blockchain:  bc,
		db:          db,
		history:     his,
	}

	addGenesisBlockToVisor(t, v)
	var gb *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		gb, err = v.blockchain.GetGenesisBlock(tx)
		return err
	})
	require.NoError(t, err)
	require.NotNil(t, gb)

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	// Create two valid transactions, both spending the same inputs, one with a higher fee
	// Then, create a block from these transactions.
	// The one with the higher fee should be included in the block, and the other should be ignored.
	// A call to RemoveInvalidUnconfirmed will remove the other txn, because it would now be a double spend.

	var coins uint64 = 10e6
	txn1 := makeSpendTxn(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectForeignTransaction(txn1)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), length)
		return nil
	})
	require.NoError(t, err)

	var fee uint64 = 1
	txn2 := makeSpendTxWithFee(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins, fee)
	known, softErr, err = v.InjectForeignTransaction(txn2)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), length)
		return nil
	})
	require.NoError(t, err)

	// Execute a block, txn2 should be included because it has a higher fee
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 2, len(sb.Body.Transactions[0].Out))
	require.Equal(t, txn2.Hash().Hex(), sb.Body.Transactions[0].Hash().Hex())

	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), length)

		length, err = bc.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(2), length)

		return nil
	})
	require.NoError(t, err)

	// Call RemoveInvalidUnconfirmed, the first txn will be removed because it is now a double-spend txn
	removed, err := v.RemoveInvalidUnconfirmed()
	require.NoError(t, err)
	require.Equal(t, []cipher.SHA256{txn1.Hash()}, removed)
	err = db.View("", func(tx *dbutil.Tx) error {
		length, err := unconfirmed.Len(tx)
		require.NoError(t, err)
		require.Equal(t, uint64(0), length)
		return nil
	})
	require.NoError(t, err)
}

func makeTxn(t *testing.T, headTime uint64, in, out []coin.UxOut, keys []cipher.SecKey) (coin.Transaction, []TransactionInput) {
	inputs := make([]cipher.SHA256, len(in))
	for i, input := range in {
		inputs[i] = input.Hash()
	}

	outputs := make([]coin.TransactionOutput, len(out))
	for i, output := range out {
		outputs[i] = coin.TransactionOutput{
			Address: output.Body.Address,
			Coins:   output.Body.Coins,
			Hours:   output.Body.Hours,
		}
	}

	txn := coin.Transaction{
		In:  inputs,
		Out: outputs,
	}

	txn.SignInputs(keys)
	err := txn.UpdateHeader()
	require.NoError(t, err)

	txnInputs, err := NewTransactionInputs(in, headTime)
	require.NoError(t, err)
	return txn, txnInputs
}

func TestVerifyTxnVerbose(t *testing.T) {
	head := coin.SignedBlock{
		Block: coin.Block{
			Head: coin.BlockHeader{
				Time: uint64(time.Now().UTC().Unix()),
			},
		},
	}

	hashes := make([]cipher.SHA256, 20)
	for i := 0; i < 20; i++ {
		hashes[i] = testutil.RandSHA256(t)
	}

	keys := make([]cipher.SecKey, 5)
	for i := 0; i < 5; i++ {
		_, keys[i] = cipher.GenerateKeyPair()
	}

	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = cipher.MustAddressFromSecKey(keys[i])
	}

	srcTxnHashes := make([]cipher.SHA256, 5)
	inputs := make([]coin.UxOut, 5)
	historyOutputs := make([]historydb.UxOut, 5)

	for i := 0; i < 5; i++ {
		srcTxnHashes[i] = testutil.RandSHA256(t)
		inputs[i] = coin.UxOut{
			Head: coin.UxHead{
				Time: head.Time(),
			},
			Body: coin.UxBody{
				SrcTransaction: srcTxnHashes[i],
				Address:        addrs[i],
				Coins:          10e6,
				Hours:          1000,
			},
		}

		historyOutputs[i] = historydb.UxOut{
			Out: inputs[i],
		}
	}

	outputs := make([]coin.UxOut, 5)
	for i := 0; i < 5; i++ {
		outputs[i] = coin.UxOut{
			Head: coin.UxHead{
				Time: head.Time(),
			},
			Body: coin.UxBody{
				Address: testutil.MakeAddress(),
				Coins:   10e6,
				Hours:   400 + uint64(i)*200,
			},
		}
	}

	// add uxout with math.MaxUint64 hours
	outputs = append(outputs, coin.UxOut{
		Head: coin.UxHead{
			Time: head.Time(),
		},
		Body: coin.UxBody{
			Address: testutil.MakeAddress(),
			Coins:   10e6,
			Hours:   math.MaxUint64,
		},
	})

	// add output which has 11e6 coins
	outputs = append(outputs, coin.UxOut{
		Head: coin.UxHead{
			Time: head.Time(),
		},
		Body: coin.UxBody{
			Address: testutil.MakeAddress(),
			Coins:   11e6,
			Hours:   500,
		},
	})

	// create a transaction
	txn, spentInputs := makeTxn(t, head.Time(), inputs[:1], outputs[:1], keys[:1])

	// create a transaction which sends coin to null address
	toNullAddrTxn, toNullAddrSpentUxBalances := makeTxn(t, head.Time(), inputs[:1], outputs[:1], keys[:1])
	toNullAddrTxn.Out[0].Address = cipher.Address{}

	// create a transaction with insufficient coin hours
	inSufficientCoinHoursTxn, _ := makeTxn(t, head.Time(), inputs[:1], outputs[4:5], keys[:1])

	// create a transaction with zero fee
	zeroFeeTxn, _ := makeTxn(t, head.Time(), inputs[:1], outputs[3:4], keys[:1])

	// create a transaction with output coin hours overflow
	coinHourOverflowTxn, _ := makeTxn(t, head.Time(), inputs[:1], outputs[4:], keys[:1])

	// create a transaction with insufficient fee
	insufficientFeeOuts := []coin.UxOut{
		coin.UxOut{
			Head: coin.UxHead{
				Time: head.Time(),
			},
			Body: coin.UxBody{
				Address: testutil.MakeAddress(),
				Coins:   1e6,
				Hours:   950,
			},
		},
	}
	insufficientFeeTxn, _ := makeTxn(t, head.Time(), inputs[:1], insufficientFeeOuts[:], keys[:1])

	// create a transaction with insufficient coins
	insufficientCoinsTxn, _ := makeTxn(t, head.Time(), inputs[:1], outputs[6:], keys[:1])

	// create a transaction with valid signature signed by key that doesn't own the input
	badSigTxn, badSigTxnSpentInputs := makeTxn(t, head.Time(), inputs[:2], outputs[:2], keys[1:3])

	// create a transaction with an empty signature
	fullyUnsignedTxn, fullyUnsignedSpentInputs := makeTxn(t, head.Time(), inputs[:2], outputs[:2], keys[:2])
	fullyUnsignedTxn.Sigs[0] = cipher.Sig{}
	fullyUnsignedTxn.Sigs[1] = cipher.Sig{}
	err := fullyUnsignedTxn.UpdateHeader()
	require.NoError(t, err)

	partiallyUnsignedTxn, partiallyUnsignedSpentInputs := makeTxn(t, head.Time(), inputs[:2], outputs[:2], keys[:2])
	partiallyUnsignedTxn.Sigs[0] = cipher.Sig{}
	err = partiallyUnsignedTxn.UpdateHeader()
	require.NoError(t, err)

	type testCase struct {
		name        string
		txn         coin.Transaction
		isConfirmed bool
		isUnsigned  bool
		signed      TxnSignedFlag
		inputs      []TransactionInput
		err         error

		maxUserTransactionSize uint32

		getArrayRet coin.UxArray
		getArrayErr error

		getHistoryTxnRet *historydb.Transaction
		getHistoryTxnErr error

		getHistoryUxOutsRet []historydb.UxOut
		getHistoryUxOutsErr error

		getSignedBlocksBySeqRet *coin.SignedBlock
		getSignedBlocksBySeqErr error
	}

	baseCases := []testCase{
		{
			name:        "transaction has been spent",
			txn:         txn,
			isConfirmed: true,
			signed:      TxnSigned,
			inputs:      spentInputs[:],

			getArrayErr: blockdb.ErrUnspentNotExist{UxID: inputs[0].Hash().Hex()},
			getHistoryTxnRet: &historydb.Transaction{
				Txn:      txn,
				BlockSeq: 10,
			},
			getHistoryUxOutsRet: historyOutputs[:1],
			getSignedBlocksBySeqRet: &coin.SignedBlock{
				Block: coin.Block{
					Head: coin.BlockHeader{
						Time: 10000000,
					},
				},
			},
		},
		{
			name:        "transaction has been spent, get previous block error",
			txn:         txn,
			isConfirmed: true,
			signed:      TxnSigned,
			inputs:      nil,
			err:         errors.New("GetSignedBlockBySeq failed"),

			getArrayErr: blockdb.ErrUnspentNotExist{UxID: inputs[0].Hash().Hex()},
			getHistoryTxnRet: &historydb.Transaction{
				Txn:      txn,
				BlockSeq: 10,
			},
			getHistoryUxOutsRet:     historyOutputs[:1],
			getSignedBlocksBySeqErr: errors.New("GetSignedBlockBySeq failed"),
		},
		{
			name:        "transaction has been spent, previous block not found",
			txn:         txn,
			isConfirmed: true,
			signed:      TxnSigned,
			inputs:      nil,
			err:         fmt.Errorf("VerifyTxnVerbose: previous block seq=%d not found", 9),

			getArrayErr: blockdb.ErrUnspentNotExist{UxID: inputs[0].Hash().Hex()},
			getHistoryTxnRet: &historydb.Transaction{
				Txn:      txn,
				BlockSeq: 10,
			},
			getHistoryUxOutsRet:     historyOutputs[:1],
			getSignedBlocksBySeqRet: nil,
			getSignedBlocksBySeqErr: nil,
		},
		{
			name:        "transaction does not exist in either unspents or historydb",
			txn:         txn,
			isConfirmed: false,
			signed:      TxnSigned,
			err:         ErrTxnViolatesHardConstraint{fmt.Errorf("transaction input of %s does not exist in either unspent pool or historydb", inputs[0].Hash().Hex())},

			getArrayErr: blockdb.ErrUnspentNotExist{UxID: inputs[0].Hash().Hex()},
		},
		{
			name:        "transaction violate user constratins, send to null address",
			txn:         toNullAddrTxn,
			isConfirmed: false,
			signed:      TxnSigned,
			err:         ErrTxnViolatesUserConstraint{errors.New("Transaction output is sent to the null address")},
			inputs:      toNullAddrSpentUxBalances[:],

			getArrayRet: inputs[:1],
		},
		{
			name:                   "transaction violate soft constraints, transaction size bigger than max block size",
			signed:                 TxnSigned,
			maxUserTransactionSize: 1,
			txn:                    txn,
			inputs:                 spentInputs[:],
			err:                    ErrTxnViolatesSoftConstraint{errors.New("Transaction size bigger than max block size")},

			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate soft constraints, Insufficient coinhours for transaction outputs",
			signed:      TxnSigned,
			txn:         inSufficientCoinHoursTxn,
			err:         ErrTxnViolatesSoftConstraint{fee.ErrTxnInsufficientCoinHours},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate soft constraints, zero fee",
			signed:      TxnSigned,
			txn:         zeroFeeTxn,
			err:         ErrTxnViolatesSoftConstraint{fee.ErrTxnNoFee},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate soft constraints, coin hour overflow",
			signed:      TxnSigned,
			txn:         coinHourOverflowTxn,
			err:         ErrTxnViolatesSoftConstraint{errors.New("Transaction output hours overflow")},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate soft constraints, insufficient fee",
			signed:      TxnSigned,
			txn:         insufficientFeeTxn,
			err:         ErrTxnViolatesSoftConstraint{fee.ErrTxnInsufficientFee},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate hard constraints, insufficient coins",
			signed:      TxnSigned,
			txn:         insufficientCoinsTxn,
			err:         ErrTxnViolatesHardConstraint{errors.New("Insufficient coins")},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
		{
			name:        "transaction violate hard constraints, bad signature",
			signed:      TxnSigned,
			txn:         badSigTxn,
			err:         ErrTxnViolatesHardConstraint{errors.New("Signature not valid for output being spent")},
			getArrayRet: inputs[:2],
			inputs:      badSigTxnSpentInputs[:],
		},
	}

	signedOnlyCases := []testCase{
		{
			name:        "TxnSignedFlag=TxnSigned transaction violate hard constraints, fully unsigned",
			signed:      TxnSigned,
			isUnsigned:  true,
			txn:         fullyUnsignedTxn,
			err:         ErrTxnViolatesHardConstraint{errors.New("Unsigned input in transaction")},
			inputs:      fullyUnsignedSpentInputs[:],
			getArrayRet: inputs[:2],
		},
		{
			name:        "TxnSignedFlag=TxnSigned transaction violate hard constraints, partially unsigned",
			signed:      TxnSigned,
			isUnsigned:  true,
			txn:         partiallyUnsignedTxn,
			err:         ErrTxnViolatesHardConstraint{errors.New("Unsigned input in transaction")},
			inputs:      partiallyUnsignedSpentInputs[:],
			getArrayRet: inputs[:2],
		},
		{
			name:        "TxnSignedFlag=TxnSigned ok fully signed",
			signed:      TxnSigned,
			txn:         txn,
			inputs:      spentInputs,
			getArrayRet: inputs[:1],
		},
	}

	unsignedOnlyCases := []testCase{
		{
			name:        "TxnSignedFlag=TxnUnsigned ok, fully unsigned",
			signed:      TxnUnsigned,
			isUnsigned:  true,
			txn:         fullyUnsignedTxn,
			inputs:      fullyUnsignedSpentInputs,
			getArrayRet: inputs[:2],
		},
		{
			name:        "TxnSignedFlag=TxnUnsigned ok, partially unsigned",
			signed:      TxnUnsigned,
			isUnsigned:  true,
			txn:         partiallyUnsignedTxn,
			inputs:      partiallyUnsignedSpentInputs,
			getArrayRet: inputs[:2],
		},
		{
			name:        "TxnSignedFlag=TxnUnsigned transaction violate hard constraints: signed",
			signed:      TxnUnsigned,
			txn:         txn,
			err:         ErrTxnViolatesHardConstraint{errors.New("Unsigned transaction must contain a null signature")},
			inputs:      spentInputs[:],
			getArrayRet: inputs[:1],
		},
	}

	cases := []testCase{}

	for _, tc := range baseCases {
		c := tc
		c.name = fmt.Sprintf("TxnSignedFlag=TxnSigned %s", tc.name)
		c.signed = TxnSigned
		cases = append(cases, c)

		c = tc
		c.name = fmt.Sprintf("TxnSignedFlag=TxnUnsigned %s", tc.name)
		c.signed = TxnUnsigned
		c.txn.Sigs = make([]cipher.Sig, len(c.txn.Sigs))
		copy(c.txn.Sigs, tc.txn.Sigs)
		c.txn.Sigs[0] = cipher.Sig{}
		c.isUnsigned = true
		cases = append(cases, c)
	}

	cases = append(cases, signedOnlyCases...)
	cases = append(cases, unsignedOnlyCases...)

	matchDBTx := mock.MatchedBy(func(tx *dbutil.Tx) bool {
		return true
	})

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db, shutdown := testutil.PrepareDB(t)
			defer shutdown()

			history := &MockHistoryer{}
			bc := &MockBlockchainer{}
			unspent := &MockUnspentPooler{}

			bc.On("Unspent").Return(unspent)
			bc.On("Head", matchDBTx).Return(&head, nil)
			if tc.getHistoryTxnRet != nil {
				bc.On("GetSignedBlockBySeq", matchDBTx, tc.getHistoryTxnRet.BlockSeq-1).Return(tc.getSignedBlocksBySeqRet, tc.getSignedBlocksBySeqErr)
			}

			unspent.On("GetArray", matchDBTx, tc.txn.In).Return(tc.getArrayRet, tc.getArrayErr)

			history.On("GetTransaction", matchDBTx, tc.txn.Hash()).Return(tc.getHistoryTxnRet, tc.getHistoryTxnErr)
			history.On("GetUxOuts", matchDBTx, tc.txn.In).Return(tc.getHistoryUxOutsRet, tc.getHistoryUxOutsErr)

			v := &Visor{
				blockchain: bc,
				db:         db,
				history:    history,
				Config:     Config{},
			}

			originalMaxUnconfirmedTxnSize := params.UserVerifyTxn.MaxTransactionSize
			defer func() {
				params.UserVerifyTxn.MaxTransactionSize = originalMaxUnconfirmedTxnSize
			}()

			if tc.maxUserTransactionSize != 0 {
				params.UserVerifyTxn.MaxTransactionSize = tc.maxUserTransactionSize
			}

			var isConfirmed bool
			var inputs []TransactionInput
			err := v.db.View("VerifyTxnVerbose", func(tx *dbutil.Tx) error {
				var err error
				inputs, isConfirmed, err = v.VerifyTxnVerbose(&tc.txn, tc.signed)
				return err
			})

			require.Equal(t, tc.isConfirmed, isConfirmed)
			require.Equal(t, tc.inputs, inputs)

			require.Equal(t, tc.err, err, "%v != %v", tc.err, err)

			if tc.isUnsigned && tc.signed == TxnSigned {
				require.Error(t, err)
			}

			if !tc.isUnsigned && tc.signed == TxnUnsigned {
				require.Error(t, err)
			}
		})
	}
}

// historyerMock2 embeds historyerMock, and rewrite the ForEach method
type historyerMock2 struct {
	MockHistoryer
	txns []historydb.Transaction
}

func newHistoryerMock2() *historyerMock2 {
	return &historyerMock2{}
}

func (h *historyerMock2) ForEachTxn(tx *dbutil.Tx, f func(cipher.SHA256, *historydb.Transaction) error) error {
	for i := range h.txns {
		if err := f(h.txns[i].Hash(), &h.txns[i]); err != nil {
			return err
		}
	}
	return nil
}

// MockUnconfirmedTransactionPooler2 embeds UnconfirmedTxnPoolerMock, and rewrite the GetFiltered method
type MockUnconfirmedTransactionPooler2 struct {
	MockUnconfirmedTransactionPooler
	txns []UnconfirmedTransaction
}

func NewUnconfirmedTransactionPoolerMock2() *MockUnconfirmedTransactionPooler2 {
	return &MockUnconfirmedTransactionPooler2{}
}

func (m *MockUnconfirmedTransactionPooler2) GetFiltered(tx *dbutil.Tx, f func(tx UnconfirmedTransaction) bool) ([]UnconfirmedTransaction, error) {
	var txns []UnconfirmedTransaction
	for i := range m.txns {
		if f(m.txns[i]) {
			txns = append(txns, m.txns[i])
		}
	}
	return txns, nil
}

func TestFbyAddresses(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	tests := []struct {
		name    string
		addrs   []cipher.Address
		outputs []coin.UxOut
		want    []coin.UxOut
	}{
		{
			"filter with one address",
			[]cipher.Address{addrs[0]},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple addresses",
			[]cipher.Address{addrs[0], addrs[1]},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		outs := FbyAddresses(tt.addrs)(tt.outputs)
		require.Equal(t, outs, coin.UxArray(tt.want))
	}
}

func TestFbyHashes(t *testing.T) {
	uxs := make(coin.UxArray, 5)
	addrs := make([]cipher.Address, 5)
	for i := 0; i < 5; i++ {
		addrs[i] = testutil.MakeAddress()
		uxs[i] = coin.UxOut{
			Body: coin.UxBody{
				Address: addrs[i],
			},
		}
	}

	tests := []struct {
		name    string
		hashes  []cipher.SHA256
		outputs coin.UxArray
		want    coin.UxArray
	}{
		{
			"filter with one hash",
			[]cipher.SHA256{uxs[0].Hash()},
			uxs[:2],
			uxs[:1],
		},
		{
			"filter with multiple hash",
			[]cipher.SHA256{uxs[0].Hash(), uxs[1].Hash()},
			uxs[:3],
			uxs[:2],
		},
	}
	for _, tt := range tests {
		outs := FbyHashes(tt.hashes)(tt.outputs)
		require.Equal(t, outs, tt.want)
	}
}
