package visor

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/blockdb"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

const (
	// GenesisTime is the time of the genesis block created in MakeBlockchain
	GenesisTime uint64 = 1000
	// GenesisCoins is the amount of coins in the genesis block created in MakeBlockchain
	GenesisCoins uint64 = 1000e6
	// GenesisCoinHours is the amount of coin hours in the genesis block created in MakeBlockchain
	GenesisCoinHours uint64 = 1000 * 1000
	// TimeIncrement is the default time increment used when creating a block with CreateGenesisSpendTransaction
	TimeIncrement uint64 = 3600 * 1000
)

var (
	// GenesisPublic is the public key used in the genesis block created in MakeBlockchain
	GenesisPublic cipher.PubKey
	// GenesisSecret is the secret key used in the genesis block created in MakeBlockchain
	GenesisSecret cipher.SecKey
	// GenesisAddress is the address used in the genesis block created in MakeBlockchain
	GenesisAddress cipher.Address
)

func init() {
	GenesisPublic, GenesisSecret = cipher.GenerateKeyPair()
	GenesisAddress = cipher.AddressFromPubKey(GenesisPublic)
}

// MakeBlockchain creates a new blockchain with a genesis block
func MakeBlockchain(t *testing.T, db *dbutil.DB, seckey cipher.SecKey) *Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := NewBlockchain(db, BlockchainConfig{
		Pubkey: pubkey,
	})
	require.NoError(t, err)
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	if err != nil {
		panic(fmt.Errorf("create genesis block failed: %v", err))
	}

	sig := cipher.SignHash(gb.HashHeader(), seckey)
	db.Update("", func(tx *dbutil.Tx) error {
		return b.ExecuteBlock(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   sig,
		})
	})
	return b
}

// CreateGenesisSpendTransaction creates the initial post-genesis transaction that moves genesis coins to another address
func CreateGenesisSpendTransaction(t *testing.T, db *dbutil.DB, bc *Blockchain, toAddr cipher.Address, coins, hours, fee uint64) coin.Transaction {
	var txn coin.Transaction
	err := db.View("", func(tx *dbutil.Tx) error {
		uxOuts, err := bc.Unspent().GetAll(tx)
		require.NoError(t, err)
		require.Len(t, uxOuts, 1)

		txn = makeTransactionForChain(t, tx, bc, uxOuts[0], GenesisSecret, toAddr, coins, hours, fee)
		require.Equal(t, txn.Out[0].Address.String(), toAddr.String())

		if coins == GenesisCoins {
			// No change output
			require.Len(t, txn.Out, 1)
		} else {
			require.Len(t, txn.Out, 2)
			require.Equal(t, txn.Out[1].Address.String(), GenesisAddress.String())
		}

		return nil
	})
	require.NoError(t, err)
	return txn
}

// ExecuteGenesisSpendTransaction executes a genesis block created with CreateGenesisSpendTransaction against a blockchain
// created with MakeBlockchain
func ExecuteGenesisSpendTransaction(t *testing.T, db *dbutil.DB, bc *Blockchain, txn coin.Transaction) coin.UxOut {
	var block *coin.Block
	err := db.View("", func(tx *dbutil.Tx) error {
		var err error
		block, err = bc.NewBlock(tx, coin.Transactions{txn}, GenesisTime+TimeIncrement)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, block)

	sig := cipher.SignHash(block.HashHeader(), GenesisSecret)
	sb := coin.SignedBlock{
		Block: *block,
		Sig:   sig,
	}

	err = db.Update("", func(tx *dbutil.Tx) error {
		err = bc.ExecuteBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	uxOut, err := coin.CreateUnspent(block.Head, txn, 0)
	require.NoError(t, err)

	return uxOut
}

func makeTransactionForChain(t *testing.T, tx *dbutil.Tx, bc *Blockchain, ux coin.UxOut, sec cipher.SecKey, toAddr cipher.Address, amt, hours, fee uint64) coin.Transaction {
	tim, err := bc.Time(tx)
	require.NoError(t, err)

	chrs, err := ux.CoinHours(tim)
	require.NoError(t, err)

	require.Equal(t, cipher.AddressFromPubKey(cipher.PubKeyFromSecKey(sec)), ux.Body.Address)

	knownUx, err := bc.Unspent().Get(tx, ux.Hash())
	require.NoError(t, err)
	require.NotNil(t, knownUx)
	require.Equal(t, knownUx, &ux)

	txn := coin.Transaction{}
	txn.PushInput(ux.Hash())

	txn.PushOutput(toAddr, amt, hours)

	// Change output
	coinsOut := ux.Body.Coins - amt
	if coinsOut > 0 {
		txn.PushOutput(GenesisAddress, coinsOut, chrs-hours-fee)
	}

	txn.SignInputs([]cipher.SecKey{sec})

	require.Equal(t, len(txn.Sigs), 1)

	err = cipher.ChkSig(ux.Body.Address, cipher.AddSHA256(txn.HashInner(), txn.In[0]), txn.Sigs[0])
	require.NoError(t, err)

	txn.UpdateHeader()

	err = txn.Verify()
	require.NoError(t, err)

	err = bc.VerifySingleTxnHardConstraints(tx, txn)
	require.NoError(t, err)

	return txn
}

func makeLostCoinTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction { // nolint: unparam
	txn := coin.Transaction{}
	var totalCoins uint64
	var totalHours uint64

	for _, ux := range uxs {
		txn.PushInput(ux.Hash())
		totalCoins += ux.Body.Coins
		totalHours += ux.Body.Hours
	}

	txn.PushOutput(toAddr, coins, totalHours/4)
	changeCoins := totalCoins - coins
	if changeCoins > 0 {
		txn.PushOutput(uxs[0].Body.Address, changeCoins-1, totalHours/4)
	}

	txn.SignInputs(keys)
	txn.UpdateHeader()
	return txn
}

func makeDuplicateUxOutTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction { // nolint: unparam
	txn := coin.Transaction{}
	var totalCoins uint64
	var totalHours uint64

	for _, ux := range uxs {
		txn.PushInput(ux.Hash())
		totalCoins += ux.Body.Coins
		totalHours += ux.Body.Hours
	}

	txn.PushOutput(toAddr, coins, totalHours/8)
	txn.PushOutput(toAddr, coins, totalHours/8)
	changeCoins := totalCoins - coins
	if changeCoins > 0 {
		txn.PushOutput(uxs[0].Body.Address, changeCoins, totalHours/4)
	}

	txn.SignInputs(keys)
	txn.UpdateHeader()
	return txn
}

// makeUnspentsTx creates a transaction that has a configurable number of outputs sent to the same address.
// The genesis block has only one unspent output, so only one transaction can be made from it.
// This is useful for when multiple test transactions need to be made from the same block.
// Coins and hours are distributed equally amongst all new outputs.
func makeUnspentsTx(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, nUnspents int, maxDivisor uint64) coin.Transaction { // nolint: unparam
	// Add inputs to the transaction
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		var err error
		totalHours, err = coin.AddUint64(totalHours, ux.Body.Hours)
		require.NoError(t, err)

		totalCoins, err = coin.AddUint64(totalCoins, ux.Body.Coins)
		require.NoError(t, err)
	}

	// Distribute coins and hours equally to all of the new outputs
	coins := totalCoins / uint64(nUnspents)
	coins = (coins / maxDivisor) * maxDivisor
	t.Logf("Assigning %d coins to each of %d outputs", coins, nUnspents)
	changeCoins := totalCoins - (coins * uint64(nUnspents))
	t.Logf("Change coins: %d", changeCoins)

	hours := (totalHours / 2) / uint64(nUnspents)
	changeHours := (totalHours / 2) - (hours * uint64(nUnspents))

	// Create the new outputs
	require.True(t, uint64(nUnspents) < hours)
	for i := 0; i < nUnspents; i++ {
		// Subtract index from hours so that the outputs are not all the same,
		// otherwise the output hashes will be duplicated and the transaction
		// will be invalid
		spendHours := hours - uint64(i)
		spendTx.PushOutput(toAddr, coins, spendHours)
	}

	// Add change output, if necessary
	if changeCoins != 0 {
		spendTx.PushOutput(uxs[0].Body.Address, changeCoins, changeHours)
	}

	// Sign the transaction
	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()

	return spendTx
}

// makeSpendTxWithFee creates a txn specified with the extra number of hours to burn in addition to the minimum required burn
func makeSpendTxWithFee(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins, fee uint64) coin.Transaction {
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	require.True(t, coins <= totalCoins)
	require.True(t, fee <= totalHours/2, "Fee must be <= half of total hours")

	spendHours := totalHours/2 - fee

	spendTx.PushOutput(toAddr, coins, spendHours)
	if totalCoins != coins {
		spendTx.PushOutput(uxs[0].Body.Address, totalCoins-coins, 0)
	}
	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()
	return spendTx
}

// makeSpendTxWithHoursBurned creates a txn specified with the total number of hours to burn
func makeSpendTxWithHoursBurned(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins, hoursBurned uint64) coin.Transaction { // nolint: unparam
	spendTx := coin.Transaction{}
	var totalHours uint64
	var totalCoins uint64
	for _, ux := range uxs {
		spendTx.PushInput(ux.Hash())
		totalHours += ux.Body.Hours
		totalCoins += ux.Body.Coins
	}

	require.True(t, coins <= totalCoins)
	require.True(t, hoursBurned <= totalHours, "hoursBurned must be <= totalHours")

	spendHours := totalHours - hoursBurned

	spendTx.PushOutput(toAddr, coins, spendHours)
	if totalCoins != coins {
		spendTx.PushOutput(uxs[0].Body.Address, totalCoins-coins, 0)
	}
	spendTx.SignInputs(keys)
	spendTx.UpdateHeader()
	return spendTx
}

func requireSoftViolation(t *testing.T, msg string, err error) {
	require.Equal(t, NewErrTxnViolatesSoftConstraint(errors.New(msg)), err)
}

func requireHardViolation(t *testing.T, msg string, err error) {
	require.Equal(t, NewErrTxnViolatesHardConstraint(errors.New(msg)), err)
}

func TestVerifyTransactionSoftHardConstraints(t *testing.T) {
	db, closeDB := prepareDB(t)
	defer closeDB()

	err := CreateBuckets(db)
	require.NoError(t, err)

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlockToBlockchain(t, bc)

	toAddr := testutil.MakeAddress()
	coins := uint64(10e6)

	verifySingleTxnSoftHardConstraints := func(txn coin.Transaction, maxBlockSize int) error {
		return db.View("", func(tx *dbutil.Tx) error {
			return bc.VerifySingleTxnSoftHardConstraints(tx, txn, maxBlockSize)
		})
	}

	// create normal spending txn
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	txn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	err = verifySingleTxnSoftHardConstraints(txn, DefaultMaxBlockSize)
	require.NoError(t, err)

	// Transaction size exceeds maxSize
	err = verifySingleTxnSoftHardConstraints(txn, txn.Size()-1)
	requireSoftViolation(t, "Transaction size bigger than max block size", err)

	// Invalid transaction fee
	uxs = coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	hours := uint64(0)
	for _, ux := range uxs {
		hours += ux.Body.Hours
	}
	txn = makeSpendTxWithHoursBurned(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins, 0)
	err = verifySingleTxnSoftHardConstraints(txn, DefaultMaxBlockSize)
	requireSoftViolation(t, "Transaction has zero coinhour fee", err)

	// Invalid transaction fee, part 2
	txn = makeSpendTxWithHoursBurned(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins, 1)
	err = verifySingleTxnSoftHardConstraints(txn, DefaultMaxBlockSize)
	requireSoftViolation(t, "Transaction coinhour fee minimum not met", err)

	// Transaction locking is tested by TestVerifyTransactionIsLocked

	// Test invalid header hash
	originInnerHash := txn.InnerHash
	txn.InnerHash = cipher.SHA256{}
	err = verifySingleTxnSoftHardConstraints(txn, DefaultMaxBlockSize)
	requireHardViolation(t, "Invalid header hash", err)

	// Set back the originInnerHash
	txn.InnerHash = originInnerHash

	// Create new block to spend the coins
	var b *coin.Block
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		b, err = bc.NewBlock(tx, coin.Transactions{txn}, genTime+100)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	require.NotNil(t, b)

	// Add the block to blockchain
	err = bc.db.Update("", func(tx *dbutil.Tx) error {
		return bc.store.AddBlock(tx, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
	})
	require.NoError(t, err)

	// A UxOut does not exist, it was already spent
	err = verifySingleTxnSoftHardConstraints(txn, DefaultMaxBlockSize)
	expectedErr := NewErrTxnViolatesHardConstraint(blockdb.NewErrUnspentNotExist(txn.In[0].Hex()))
	require.Equal(t, expectedErr, err)

	// Check invalid sig
	uxs = coin.CreateUnspents(b.Head, txn)
	_, key := cipher.GenerateKeyPair()
	toAddr2 := testutil.MakeAddress()
	tx2 := makeSpendTx(t, uxs, []cipher.SecKey{key, key}, toAddr2, 5e6)
	err = verifySingleTxnSoftHardConstraints(tx2, DefaultMaxBlockSize)
	requireHardViolation(t, "Signature not valid for output being spent", err)

	// Create lost coin transaction
	uxs2 := coin.CreateUnspents(b.Head, txn)
	toAddr3 := testutil.MakeAddress()
	lostCoinTx := makeLostCoinTx(coin.UxArray{uxs2[1]}, []cipher.SecKey{genSecret}, toAddr3, 10e5)
	err = verifySingleTxnSoftHardConstraints(lostCoinTx, DefaultMaxBlockSize)
	requireHardViolation(t, "Transactions may not destroy coins", err)

	// Create transaction with duplicate UxOuts
	uxs = coin.CreateUnspents(b.Head, txn)
	toAddr4 := testutil.MakeAddress()
	dupUxOutTx := makeDuplicateUxOutTx(coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr4, 1e6)
	err = verifySingleTxnSoftHardConstraints(dupUxOutTx, DefaultMaxBlockSize)
	requireHardViolation(t, "Duplicate output in transaction", err)
}

func TestVerifyTxnFeeCoinHoursAdditionFails(t *testing.T) {
	// Test that VerifySingleTxnSoftConstraints fails if a uxIn.CoinHours() call fails.
	// This is a separate test on its own, because it's not possible to reach the line
	// that is being tested through the blockchain verify API wrappers
	db, closeDB := prepareDB(t)
	defer closeDB()

	err := CreateBuckets(db)
	require.NoError(t, err)

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlockToBlockchain(t, bc)

	toAddr := testutil.MakeAddress()
	coins := uint64(10e6)

	// create normal spending txn
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	txn := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)

	var uxIn coin.UxArray
	var head *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		uxIn, err = bc.Unspent().GetArray(tx, txn.In)
		require.NoError(t, err)
		require.NotEmpty(t, uxIn)

		head, err = bc.Head(tx)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	// Set the uxIn's hours high, so that uxIn.CoinHours() returns an error
	uxIn[0].Body.Hours = math.MaxUint64
	_, coinHoursErr := uxIn[0].CoinHours(head.Time() + 1e6)
	testutil.RequireError(t, coinHoursErr, "UxOut.CoinHours addition of earned coin hours overflow")

	// VerifySingleTxnSoftConstraints should fail on this, when trying to calculate the TransactionFee
	err = VerifySingleTxnSoftConstraints(txn, head.Time()+1e6, uxIn, DefaultMaxBlockSize)
	testutil.RequireError(t, err, NewErrTxnViolatesSoftConstraint(coinHoursErr).Error())

	// VerifySingleTxnHardConstraints should fail on this, when performing the extra check of
	// uxIn.CoinHours() errors, which is ignored by VerifyTransactionHoursSpending if the error
	// is because of the earned hours addition overflow
	head.Block.Head.Time += 1e6
	err = VerifySingleTxnHardConstraints(txn, head, uxIn)
	testutil.RequireError(t, err, NewErrTxnViolatesHardConstraint(coinHoursErr).Error())
}

func TestVerifyTransactionIsLocked(t *testing.T) {
	for _, addr := range GetLockedDistributionAddresses() {
		t.Run(fmt.Sprintf("IsLocked: %s", addr), func(t *testing.T) {
			testVerifyTransactionAddressLocking(t, addr, errors.New("Transaction has locked address inputs"))
		})
	}
}

func TestVerifyTransactionIsUnlocked(t *testing.T) {
	for _, addr := range GetUnlockedDistributionAddresses() {
		t.Run(fmt.Sprintf("IsUnlocked: %s", addr), func(t *testing.T) {
			testVerifyTransactionAddressLocking(t, addr, nil)
		})
	}
}

func testVerifyTransactionAddressLocking(t *testing.T, toAddr string, expectedErr error) {
	addr, err := cipher.DecodeBase58Address(toAddr)
	require.NoError(t, err)

	db, close := prepareDB(t)
	defer close()

	_, s := cipher.GenerateKeyPair()

	// Setup blockchain
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address
	var coins = GenesisCoins
	var hours uint64 = 1e6
	var fee uint64 = 5e8

	txn := CreateGenesisSpendTransaction(t, db, bc, addr, coins, hours, fee)
	uxOut := ExecuteGenesisSpendTransaction(t, db, bc, txn)

	// Create a transaction that spends from the locked address
	// The secret key for the locked address is obviously unavailable here,
	// instead, forge an invalid transaction.
	// Transaction.Verify() is called after TransactionIsLocked(),
	// so for this test it doesn't matter if transaction signature is wrong
	randomAddress := testutil.MakeAddress()
	txn = coin.Transaction{
		In: []cipher.SHA256{uxOut.Hash()},
		Out: []coin.TransactionOutput{
			{
				Address: randomAddress,
				Coins:   uxOut.Body.Coins,
				Hours:   uxOut.Body.Hours / 2,
			},
		},
	}

	var uxIn coin.UxArray
	var head *coin.SignedBlock
	err = db.View("", func(tx *dbutil.Tx) error {
		var err error
		uxIn, err = bc.Unspent().GetArray(tx, txn.In)
		require.NoError(t, err)
		require.NotEmpty(t, uxIn)

		head, err = bc.Head(tx)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	err = VerifySingleTxnSoftConstraints(txn, head.Time(), uxIn, DefaultMaxBlockSize)
	if expectedErr == nil {
		require.NoError(t, err)
	} else {
		requireSoftViolation(t, expectedErr.Error(), err)
	}
}
