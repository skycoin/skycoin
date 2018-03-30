package visor

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	_require "github.com/skycoin/skycoin/src/testutil/require"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/util/utc"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/historydb"
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
		db, err := OpenDB(badDBFile, false)
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
	badDB, err := OpenDB(badDBFile, false)
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
		db, err := OpenDB(badDBFile, false)
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
	_require.PanicsWithLogMessage(t, "Only master chain can create blocks", func() {
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
	known, softErr, err := unconfirmed.InjectTransaction(bc, txn, v.Config.MaxBlockSize)
	require.False(t, known)
	require.Nil(t, softErr)
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
		known, _, err := unconfirmed.InjectTransaction(bc, txn, v.Config.MaxBlockSize)
		require.False(t, known)
		require.NoError(t, err)
	}

	sb, err = v.CreateBlock(when + 100)
	require.NoError(t, err)
	require.Equal(t, when+100, sb.Block.Head.Time)

	blockTxns := sb.Block.Body.Transactions
	require.NotEqual(t, len(txns), len(blockTxns), "Txns should be truncated")
	require.Equal(t, 18, len(blockTxns))

	// Check fee ordering
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
	_require.PanicsWithLogMessage(t, "Only master chain can create blocks", func() {
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

	// Create a transaction with valid decimal places
	txn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectTransaction(txn)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)

	// Execute a block to clear this transaction from the pool
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 2, len(sb.Body.Transactions[0].Out))
	require.Equal(t, 0, unconfirmed.Len())
	require.Equal(t, uint64(2), bc.Len())

	uxs = coin.CreateUnspents(sb.Head, sb.Body.Transactions[0])

	// Check transactions with overflowing output coins fail
	txn = makeOverflowCoinsSpendTx(t, coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr)
	_, softErr, err = v.InjectTransaction(txn)
	require.IsType(t, ErrTxnViolatesHardConstraint{}, err)
	testutil.RequireError(t, err.(ErrTxnViolatesHardConstraint).Err, "Output coins overflow")
	require.Nil(t, softErr)
	require.Equal(t, 0, unconfirmed.Len())

	// Check transactions with overflowing output hours fail
	// It should not be injected; when injecting a txn, the overflowing output hours is treated
	// as a hard constraint. It is only a soft constraint when the txn is included in a signed block.
	txn = makeOverflowHoursSpendTx(t, coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr)
	_, softErr, err = v.InjectTransaction(txn)
	require.Nil(t, softErr)
	require.IsType(t, ErrTxnViolatesHardConstraint{}, err)
	testutil.RequireError(t, err.(ErrTxnViolatesHardConstraint).Err, "Transaction output hours overflow")
	require.Equal(t, 0, unconfirmed.Len())

	// Create a transaction with invalid decimal places
	// It's still injected, because this is considered a soft error
	invalidCoins := coins + (maxDropletDivisor / 10)
	txn = makeSpendTx(t, uxs, []cipher.SecKey{genSecret, genSecret}, toAddr, invalidCoins)
	_, softErr, err = v.InjectTransaction(txn)
	require.NoError(t, err)
	testutil.RequireError(t, softErr.Err, errInvalidDecimals.Error())
	require.Equal(t, 1, unconfirmed.Len())
}

func makeOverflowCoinsSpendTx(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address) coin.Transaction {
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	hours := totalHours / 12

	// These two outputs' coins added up will overflow
	spendTx.PushOutput(toAddr, 18446744073709551000, hours)
	spendTx.PushOutput(toAddr, totalCoins, hours)

	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()
	return spendTx
}

func makeOverflowHoursSpendTx(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address) coin.Transaction {
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	hours := totalHours / 12

	// These two outputs' hours added up will overflow
	spendTx.PushOutput(toAddr, totalCoins/2, 18446744073709551615)
	spendTx.PushOutput(toAddr, totalCoins-totalCoins/2, hours)

	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()
	return spendTx
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

	_require.PanicsWithLogMessage(t, "precision must be <= droplet.Exponent", func() {
		calculateDivisor(7)
	})
}

func makeTestData(t *testing.T, n int) ([]historydb.Transaction,
	[]coin.SignedBlock, []UnconfirmedTxn, uint64) {
	var txs []historydb.Transaction
	var blocks []coin.SignedBlock
	var uncfmTxs []UnconfirmedTxn
	for i := uint64(0); i < uint64(n); i++ {
		tm := utc.UnixNow() + int64(i)*int64(time.Second)
		txs = append(txs, historydb.Transaction{
			BlockSeq: i,
			Tx: coin.Transaction{
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

		uncfmTxs = append(uncfmTxs, UnconfirmedTxn{
			Txn: coin.Transaction{
				InnerHash: testutil.RandSHA256(t),
			},
			Received: utc.UnixNow() + int64(n)*int64(time.Second),
		})
	}

	return txs, blocks, uncfmTxs, uint64(n)
}

func makeUncfmUxs(txs []UnconfirmedTxn) coin.UxArray {
	var uxs coin.UxArray
	for i := range txs {
		uxs = append(uxs, coin.UxOut{
			Head: coin.UxHead{
				Time: uint64(txs[i].Received),
			},
			Body: coin.UxBody{
				SrcTransaction: txs[i].Hash(),
			},
		})
	}
	return uxs
}

type txsAndUncfmTxs struct {
	Txs      []historydb.Transaction
	UncfmTxs []UnconfirmedTxn
}
type expectTxResult struct {
	txs      []Transaction
	uncfmTxs []Transaction
	err      error
}

func TestGetTransctions(t *testing.T) {
	// Generates test data
	txs, blocks, uncfmTxs, headSeq := makeTestData(t, 10)
	// Generates []Transaction
	var ltxs []Transaction
	for i := range txs {
		height := headSeq - txs[i].BlockSeq + 1
		ltxs = append(ltxs, Transaction{
			Txn:    txs[i].Tx,
			Status: NewConfirmedTransactionStatus(height, txs[i].BlockSeq),
			Time:   blocks[i].Time(),
		})
	}

	// Generate unconfirmed []Transaction
	var luncfmTxs []Transaction
	for i, tx := range uncfmTxs {
		luncfmTxs = append(luncfmTxs, Transaction{
			Txn:    uncfmTxs[i].Txn,
			Status: NewUnconfirmedTransactionStatus(),
			Time:   uint64(nanoToTime(tx.Received).Unix()),
		})
	}

	// Generates addresses
	var addrs []cipher.Address
	for i := 0; i < 10; i++ {
		addrs = append(addrs, testutil.MakeAddress())
	}

	tt := []struct {
		name      string
		addrTxns  map[cipher.Address]txsAndUncfmTxs
		blocks    []coin.SignedBlock
		bcHeadSeq uint64
		filters   []TxFilter
		expect    expectTxResult
	}{
		{
			"addrFilter=1 addr=1 txs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=0 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=0 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=2 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=1 txs=2 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=0 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=0 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=1 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=2 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: nil,
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=2 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=2 unconfirmedTxs=3",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: luncfmTxs[:3],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=3 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: nil,
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:3],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=3 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:3],
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=3 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:3],
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 addr=2 txs=3 unconfirmedTxs=3",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: uncfmTxs[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
			},
			expectTxResult{
				txs:      ltxs[:3],
				uncfmTxs: luncfmTxs[:3],
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false txs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=false confirmedTxs=2 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"confirmedTxFilter=1 confirmed=true confirmedTxs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmedTxFilter=1 confirmed=false addr=1 txs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=1 txs=1 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:1],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txs=2 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2 txs=2 unconfirmedTxs=3",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:3],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2/1 txs=2 unconfirmedTxs=3",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[:2],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=false addr=2/2 txs=2 unconfirmedTxs=3",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: uncfmTxs[2:3],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[1:2]),
				ConfirmedTxFilter(false),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: luncfmTxs[2:3],
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txs=0 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      nil,
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      nil,
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txs=1 unconfirmedTxs=0",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: nil,
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txs=1 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:1],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:1],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txs=2 unconfirmedTxs=1",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=1 txs=2 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2/1 txs=3 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: uncfmTxs[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:1]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:2],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2/2 txs=3 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: uncfmTxs[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[1:2]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[2:3],
				uncfmTxs: nil,
				err:      nil,
			},
		},
		{
			"addrFilter=1 confirmed=true addr=2 txs=3 unconfirmedTxs=2",
			map[cipher.Address]txsAndUncfmTxs{
				addrs[0]: txsAndUncfmTxs{
					Txs:      txs[:2],
					UncfmTxs: uncfmTxs[:1],
				},
				addrs[1]: txsAndUncfmTxs{
					Txs:      txs[2:3],
					UncfmTxs: uncfmTxs[1:2],
				},
			},
			blocks[:],
			headSeq,
			[]TxFilter{
				AddrsFilter(addrs[:2]),
				ConfirmedTxFilter(true),
			},
			expectTxResult{
				txs:      ltxs[:3],
				uncfmTxs: nil,
				err:      nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			his := newHistoryerMock2()
			uncfmTxPool := NewUnconfirmedTxnPoolerMock2()
			for addr, txs := range tc.addrTxns {
				his.On("GetAddrTxns", addr).Return(txs.Txs, nil)
				his.txs = append(his.txs, txs.Txs...)

				uncfmTxPool.On("GetUnspentsOfAddr", addr).Return(makeUncfmUxs(txs.UncfmTxs))
				for i, uncfmTx := range txs.UncfmTxs {
					uncfmTxPool.On("Get", uncfmTx.Hash()).Return(&txs.UncfmTxs[i], true)
				}
				uncfmTxPool.txs = append(uncfmTxPool.txs, txs.UncfmTxs...)
			}

			bc := NewBlockchainerMock()
			for i, b := range tc.blocks {
				bc.On("GetBlockBySeq", b.Seq()).Return(&tc.blocks[i], nil)
			}

			bc.On("HeadSeq").Return(tc.bcHeadSeq)

			v := &Visor{
				history:     his,
				Unconfirmed: uncfmTxPool,
				Blockchain:  bc,
			}

			retTxs, err := v.GetTransactions(tc.filters...)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Len(t, retTxs, len(tc.expect.txs)+len(tc.expect.uncfmTxs))

			// Splits confirmed and unconfirmed txs in returned transactions
			uncfmTxMap := make(map[cipher.SHA256]Transaction)
			txMap := make(map[cipher.SHA256]Transaction)
			for i, tx := range retTxs {
				if retTxs[i].Status.Confirmed {
					txMap[tx.Txn.Hash()] = retTxs[i]
				} else {
					uncfmTxMap[tx.Txn.Hash()] = retTxs[i]
				}
			}

			// Confirms that all expected confirmed transactions must be in the txMap
			for _, tx := range tc.expect.txs {
				retTx, ok := txMap[tx.Txn.Hash()]
				require.True(t, ok)
				require.Equal(t, tx, retTx)
			}

			// Confirms that all expected unconfirmed transactions must be in the uncfmTxMap
			for _, tx := range tc.expect.uncfmTxs {
				retTx, ok := uncfmTxMap[tx.Txn.Hash()]
				require.True(t, ok)
				require.Equal(t, tx, retTx)
			}
		})
	}
}

func TestRefreshUnconfirmed(t *testing.T) {
	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	db, bc, err := loadBlockchain(db, genPublic, false)
	require.NoError(t, err)

	unconfirmed := NewUnconfirmedTxnPool(db)

	cfg := NewVisorConfig()
	cfg.DBPath = db.Path()
	cfg.IsMaster = true
	cfg.BlockchainSeckey = genSecret
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress

	v := &Visor{
		Config:      cfg,
		Unconfirmed: unconfirmed,
		Blockchain:  bc,
		db:          db,
	}

	addGenesisBlock(t, v.Blockchain)
	gb := v.Blockchain.GetGenesisBlock()
	require.NotNil(t, gb)

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	toAddr := testutil.MakeAddress()
	var coins uint64 = 10e6

	// Create a valid transaction that will remain valid
	validTxn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectTransaction(validTxn)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)
	require.Equal(t, 1, unconfirmed.Len())

	// Create a transaction with invalid decimal places
	// It's still injected, because this is considered a soft error
	// This transaction will stay invalid on refresh
	invalidCoins := coins + (maxDropletDivisor / 10)
	alwaysInvalidTxn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, invalidCoins)
	_, softErr, err = v.InjectTransaction(alwaysInvalidTxn)
	require.NoError(t, err)
	testutil.RequireError(t, softErr.Err, errInvalidDecimals.Error())
	require.Equal(t, 2, unconfirmed.Len())

	// Create a transaction that exceeds MaxBlockSize
	// It's still injected, because this is considered a soft error
	// This transaction will become valid on refresh (by increasing MaxBlockSize)
	v.Config.MaxBlockSize = 1
	sometimesInvalidTxn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	_, softErr, err = v.InjectTransaction(sometimesInvalidTxn)
	require.NoError(t, err)
	testutil.RequireError(t, softErr.Err, errTxnExceedsMaxBlockSize.Error())
	require.Equal(t, 3, unconfirmed.Len())

	// The first txn remains valid,
	// the second txn remains invalid,
	// the third txn becomes valid
	v.Config.MaxBlockSize = DefaultMaxBlockSize
	hashes, err := v.RefreshUnconfirmed()
	require.NoError(t, err)
	require.Equal(t, []cipher.SHA256{sometimesInvalidTxn.Hash()}, hashes)

	// Reduce the max block size to affirm that the valid transaction becomes invalid
	// The first txn becomes invalid,
	// the second txn remains invalid,
	// the third txn becomes invalid again
	v.Config.MaxBlockSize = 1
	hashes, err = v.RefreshUnconfirmed()
	require.NoError(t, err)
	require.Nil(t, hashes)

	// Restore the max block size to affirm the expected transaction validity shifts
	// The first txn was valid, became invalid, and is now valid again
	// The second txn was always invalid
	// The third txn was invalid, became valid, became invalid, and is now valid again
	v.Config.MaxBlockSize = DefaultMaxBlockSize
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
	db, shutdown := testutil.PrepareDB(t)
	defer shutdown()

	db, bc, err := loadBlockchain(db, genPublic, true)
	require.NoError(t, err)

	unconfirmed := NewUnconfirmedTxnPool(db)

	cfg := NewVisorConfig()
	cfg.DBPath = db.Path()
	cfg.IsMaster = true
	cfg.BlockchainPubkey = genPublic
	cfg.GenesisAddress = genAddress
	cfg.BlockchainSeckey = genSecret

	v := &Visor{
		Config:      cfg,
		Unconfirmed: unconfirmed,
		Blockchain:  bc,
		db:          db,
	}

	addGenesisBlock(t, v.Blockchain)
	gb := v.Blockchain.GetGenesisBlock()
	require.NotNil(t, gb)

	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])

	// Create two valid transactions, both spending the same inputs, one with a higher fee
	// Then, create a block from these transactions.
	// The one with the higher fee should be included in the block, and the other should be ignored.
	// A call to RemoveInvalidUnconfirmed will remove the other txn, because it would now be a double spend.

	var coins uint64 = 10e6
	txn1 := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins)
	known, softErr, err := v.InjectTransaction(txn1)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)
	require.Equal(t, 1, unconfirmed.Len())

	var fee uint64 = 1
	txn2 := makeSpendTxWithFee(t, uxs, []cipher.SecKey{genSecret}, genAddress, coins, fee)
	known, softErr, err = v.InjectTransaction(txn2)
	require.False(t, known)
	require.Nil(t, softErr)
	require.NoError(t, err)
	require.Equal(t, 2, unconfirmed.Len())

	// Execute a block, txn2 should be included because it has a higher fee
	sb, err := v.CreateAndExecuteBlock()
	require.NoError(t, err)
	require.Equal(t, 1, len(sb.Body.Transactions))
	require.Equal(t, 2, len(sb.Body.Transactions[0].Out))
	require.Equal(t, 1, unconfirmed.Len())
	require.Equal(t, uint64(2), bc.Len())
	require.Equal(t, txn2.TxIDHex(), sb.Body.Transactions[0].TxIDHex())

	// Call RemoveInvalidUnconfirmed, the first txn will be removed because it is now a double-spend txn
	removed, err := v.RemoveInvalidUnconfirmed()
	require.NoError(t, err)
	require.Equal(t, []cipher.SHA256{txn1.Hash()}, removed)
	require.Equal(t, 0, unconfirmed.Len())
}

// historyerMock2 embeds historyerMock, and rewrite the ForEach method
type historyerMock2 struct {
	historyerMock
	txs []historydb.Transaction
}

func newHistoryerMock2() *historyerMock2 {
	return &historyerMock2{}
}

func (h *historyerMock2) ForEach(f func(tx *historydb.Transaction) error) error {
	for i := range h.txs {
		if err := f(&h.txs[i]); err != nil {
			return err
		}
	}
	return nil
}

// UnconfirmedTxnPoolerMock2 embeds UnconfirmedTxnPoolerMock, and rewrite the GetTxns method
type UnconfirmedTxnPoolerMock2 struct {
	UnconfirmedTxnPoolerMock
	txs []UnconfirmedTxn
}

func NewUnconfirmedTxnPoolerMock2() *UnconfirmedTxnPoolerMock2 {
	return &UnconfirmedTxnPoolerMock2{}
}

func (m *UnconfirmedTxnPoolerMock2) GetTxns(f func(tx UnconfirmedTxn) bool) []UnconfirmedTxn {
	var txs []UnconfirmedTxn
	for i := range m.txs {
		if f(m.txs[i]) {
			txs = append(txs, m.txs[i])
		}
	}
	return txs
}
