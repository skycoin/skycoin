package visor

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

const (
	blockchainPubkeyStr = "0328c576d3f420e7682058a981173a4b374c7cc5ff55bf394d3cf57059bbe6456a"
)

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

func TestErrSignatureLostRecreateDB(t *testing.T) {
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
		db, err := OpenDB(badDBFile)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		_, err = NewBlockchain(db, pubkey, Arbitrating(false))
		require.Error(t, err)
		require.IsType(t, blockdb.ErrMissingSignature{}, err)
	}()

	// Loading this invalid db should cause loadBlockchain() to recreate the db
	t.Logf("Loading the corrupted db from %s", badDBFile)
	badDB, err := OpenDB(badDBFile)
	require.NoError(t, err)
	require.NotNil(t, badDB)
	require.NotEmpty(t, badDB.Path())
	t.Logf("badDB.Path() == %s", badDB.Path())

	db, bc, err := loadBlockchain(badDB, pubkey, false)
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)

	require.NotNil(t, db)
	require.NotNil(t, bc)

	// A corrupted database file should exist
	corruptFiles = findCorruptDBFiles(t, badDBFile)
	require.Len(t, corruptFiles, 1)

	// A new db should be written in place of the old bad db, and not be corrupted
	t.Logf("Checking that the new db file is valid")
	func() {
		db, err := OpenDB(badDBFile)
		require.NoError(t, err)
		defer func() {
			err := db.Close()
			assert.NoError(t, err)
		}()

		// The new db is not corrupted and loads without error
		bc, err := NewBlockchain(db, pubkey, Arbitrating(false))
		require.NoError(t, err)
		require.NotNil(t, bc)
	}()
}

func TestVisorCreateBlock(t *testing.T) {
	when := uint64(time.Now().UTC().Unix())

	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	db, bc, err := loadBlockchain(db, genPublic, false)
	require.NoError(t, err)

	unconfirmed := NewUnconfirmedTxnPool(db)

	cfg := NewVisorConfig()
	cfg.DBPath = db.Path()
	cfg.IsMaster = false
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		Unconfirmed: unconfirmed,
		Blockchain:  bc,
		db:          db,
	}

	// CreateBlock panics if called when not master
	require.PanicsWithValue(t, "Only master chain can create blocks", func() {
		v.CreateBlock(when)
	})

	v.Config.IsMaster = true
	v.Config.BlockchainSeckey = genSecret

	addGenesisBlock(t, v.Blockchain)
	gb := v.Blockchain.GetGenesisBlock()
	require.NotNil(t, gb)

	// If no transactions in the unconfirmed pool, return an error
	_, err = v.CreateBlock(when)
	testutil.RequireError(t, err, "No transactions")

	// Create enough unspent outputs to create all of these transactions
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	nUnspents := 100
	txn := makeUnspentsTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, nUnspents, maxDropletDivisor)
	known, err := unconfirmed.InjectTxn(bc, txn)
	require.False(t, known)
	require.NoError(t, err)

	v.Config.MaxBlockSize = txn.Size()
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 0, unconfirmed.Len())
	v.Config.MaxBlockSize = 1024 * 4

	// Create various transactions and add them to unconfirmed pool
	uxs = coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])
	var coins uint64 = 9e6
	var f uint64 = 10
	toAddr := testutil.MakeAddress()

	// Add more transactions than is allowed in a block, to verify truncation
	var txns coin.Transactions
	var i int
	for len(txns) == len(txns.TruncateBytesTo(v.Config.MaxBlockSize)) {
		tx := makeSpendTxWithFee(t, coin.UxArray{uxs[i]}, []cipher.SecKey{genSecret}, toAddr, coins, f)
		txns = append(txns, tx)
		i++
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
	i++

	// Confirm that at least one transaction has an invalid decimal output
	foundInvalidCoins := false
	for _, txn := range txns {
		for _, o := range txn.Out {
			if err := DropletPrecisionCheck(o.Coins); err != nil {
				foundInvalidCoins = true
				break
			}
		}
	}
	require.True(t, foundInvalidCoins)

	// Inject transactions into the unconfirmed pool
	for _, txn := range txns {
		known, err := unconfirmed.InjectTxn(bc, txn)
		require.False(t, known)
		require.NoError(t, err)
	}

	sb, err = v.CreateBlock(when + 1e6)
	require.NoError(t, err)
	require.Equal(t, when+1e6, sb.Block.Head.Time)

	blockTxns := sb.Block.Body.Transactions
	require.NotEqual(t, len(txns), len(blockTxns), "Txns should be truncated")
	require.Equal(t, 18, len(blockTxns))

	// Check f ordering
	inUxs, err := v.Blockchain.Unspent().GetArray(blockTxns[0].In)
	require.NoError(t, err)
	prevFee, err := fee.TransactionFee(&blockTxns[0], sb.Head.Time, inUxs)
	require.NoError(t, err)

	for i := 1; i < len(blockTxns); i++ {
		inUxs, err := v.Blockchain.Unspent().GetArray(blockTxns[i].In)
		require.NoError(t, err)
		f, err := fee.TransactionFee(&blockTxns[i], sb.Head.Time, inUxs)
		require.NoError(t, err)
		require.True(t, f <= prevFee)
		prevFee = f
	}

	// Check that decimal rules are enforced
	for i, txn := range blockTxns {
		for j, o := range txn.Out {
			err := DropletPrecisionCheck(o.Coins)
			require.NoError(t, err, "txout %d.%d coins=%d", i, j, o.Coins)
		}
	}
}

func TestVisorInjectTransaction(t *testing.T) {
	when := uint64(time.Now().UTC().Unix())

	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	db, bc, err := loadBlockchain(db, genPublic, false)
	require.NoError(t, err)

	unconfirmed := NewUnconfirmedTxnPool(db)

	cfg := NewVisorConfig()
	cfg.DBPath = db.Path()
	cfg.IsMaster = false
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		Unconfirmed: unconfirmed,
		Blockchain:  bc,
		db:          db,
	}

	// CreateBlock panics if called when not master
	require.PanicsWithValue(t, "Only master chain can create blocks", func() {
		v.CreateBlock(when)
	})

	v.Config.IsMaster = true
	v.Config.BlockchainSeckey = genSecret

	addGenesisBlock(t, v.Blockchain)
	gb := v.Blockchain.GetGenesisBlock()
	require.NotNil(t, gb)

	// If no transactions in the unconfirmed pool, return an error
	_, err = v.CreateBlock(when)
	testutil.RequireError(t, err, "No transactions")

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	toAddr := testutil.MakeAddress()
	var coins uint64 = 10e6

	// Create an transaction with valid decimal places
	txn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, err := v.InjectTxn(txn)
	require.False(t, known)
	require.NoError(t, err)

	// Execute a block to clear this transaction from the pool
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 2, len(sb.Body.Transactions[0].Out))
	require.Equal(t, 0, unconfirmed.Len())
	require.Equal(t, uint64(2), bc.Len())

	// Create a transaction with invalid decimal places
	uxs = coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])

	invalidCoins := coins + (maxDropletDivisor / 10)
	txn = makeSpendTx(t, uxs, []cipher.SecKey{genSecret, genSecret}, toAddr, invalidCoins)
	_, err = v.InjectTxn(txn)
	testutil.RequireError(t, err, ErrInvalidDecimals.Error())
	require.Equal(t, 0, unconfirmed.Len())
}

func TestVisorCalculatePrecision(t *testing.T) {
	cases := []struct {
		precision uint64
		divisor   uint64
	}{
		{0, 1e6},
		{1, 1e5},
		{2, 1e4},
		{3, 1e3},
		{4, 1e2},
		{5, 1e1},
		{6, 1},
	}

	for _, tc := range cases {
		name := fmt.Sprintf("calculateDivisor(%d)=%d", tc.precision, tc.divisor)
		t.Run(name, func(t *testing.T) {
			divisor := calculateDivisor(tc.precision)
			require.Equal(t, tc.divisor, divisor, "%d != %d", tc.divisor, divisor)
		})
	}

	require.PanicsWithValue(t, "precision must be <= droplet.Exponent", func() {
		calculateDivisor(7)
	})
}
