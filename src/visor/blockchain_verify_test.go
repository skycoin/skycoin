package visor

import (
	"errors"
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

var (
	GenesisPublic, GenesisSecret = cipher.GenerateKeyPair()
	GenesisAddress               = cipher.AddressFromPubKey(GenesisPublic)
)

const (
	TimeIncrement    uint64 = 3600 * 1000
	GenesisTime      uint64 = 1000
	GenesisCoins     uint64 = 1000e6
	GenesisCoinHours uint64 = 1000 * 1000
)

func MakeTransactionForChain(t *testing.T, bc *Blockchain, ux coin.UxOut, sec cipher.SecKey, toAddr cipher.Address, amt, hours, fee uint64) coin.Transaction {
	chrs, err := ux.CoinHours(bc.Time())
	require.NoError(t, err)

	require.Equal(t, cipher.AddressFromPubKey(cipher.PubKeyFromSecKey(sec)), ux.Body.Address)

	knownUx, exists := bc.Unspent().Get(ux.Hash())
	require.True(t, exists)
	require.Equal(t, knownUx, ux)

	tx := coin.Transaction{}
	tx.PushInput(ux.Hash())

	tx.PushOutput(toAddr, amt, hours)

	// Change output
	coinsOut := ux.Body.Coins - amt
	if coinsOut > 0 {
		tx.PushOutput(GenesisAddress, coinsOut, chrs-hours-fee)
	}

	tx.SignInputs([]cipher.SecKey{sec})

	require.Equal(t, len(tx.Sigs), 1)

	err = cipher.ChkSig(ux.Body.Address, cipher.AddSHA256(tx.HashInner(), tx.In[0]), tx.Sigs[0])
	require.NoError(t, err)

	tx.UpdateHeader()

	err = tx.Verify()
	require.NoError(t, err)

	err = bc.VerifySingleTxnHardConstraints(tx)
	require.NoError(t, err)

	return tx
}

func MakeBlockchain(t *testing.T, db *bolt.DB, seckey cipher.SecKey) *Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := NewBlockchain(db, pubkey)
	require.NoError(t, err)
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	if err != nil {
		panic(fmt.Errorf("create genesis block failed: %v", err))
	}

	sig := cipher.SignHash(gb.HashHeader(), seckey)
	db.Update(func(tx *bolt.Tx) error {
		return b.ExecuteBlockWithTx(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   sig,
		})
	})
	return b
}

func MakeAddress() (cipher.PubKey, cipher.SecKey, cipher.Address) {
	p, s := cipher.GenerateKeyPair()
	a := cipher.AddressFromPubKey(p)
	return p, s, a
}

func makeLostCoinTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction {
	tx := coin.Transaction{}
	var totalCoins uint64
	var totalHours uint64

	for _, ux := range uxs {
		tx.PushInput(ux.Hash())
		totalCoins += ux.Body.Coins
		totalHours += ux.Body.Hours
	}

	tx.PushOutput(toAddr, coins, totalHours/4)
	changeCoins := totalCoins - coins
	if changeCoins > 0 {
		tx.PushOutput(uxs[0].Body.Address, changeCoins-1, totalHours/4)
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return tx
}

func makeDuplicateUxOutTx(uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins uint64) coin.Transaction {
	tx := coin.Transaction{}
	var totalCoins uint64
	var totalHours uint64

	for _, ux := range uxs {
		tx.PushInput(ux.Hash())
		totalCoins += ux.Body.Coins
		totalHours += ux.Body.Hours
	}

	tx.PushOutput(toAddr, coins, totalHours/8)
	tx.PushOutput(toAddr, coins, totalHours/8)
	changeCoins := totalCoins - coins
	if changeCoins > 0 {
		tx.PushOutput(uxs[0].Body.Address, changeCoins, totalHours/4)
	}

	tx.SignInputs(keys)
	tx.UpdateHeader()
	return tx
}

// makeUnspentsTx creates a transaction that has a configurable number of outputs sent to the same address.
// The genesis block has only one unspent output, so only one transaction can be made from it.
// This is useful for when multiple test transactions need to be made from the same block.
// Coins and hours are distributed equally amongst all new outputs.
func makeUnspentsTx(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, nUnspents int, maxDivisor uint64) coin.Transaction {
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
func makeSpendTxWithHoursBurned(t *testing.T, uxs coin.UxArray, keys []cipher.SecKey, toAddr cipher.Address, coins, hoursBurned uint64) coin.Transaction {
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

func createGenesisSpendTransaction(t *testing.T, bc *Blockchain, toAddr cipher.Address, coins, hours, fee uint64) coin.Transaction {
	uxOuts, err := bc.Unspent().GetAll()
	require.NoError(t, err)
	require.Len(t, uxOuts, 1)

	txn := MakeTransactionForChain(t, bc, uxOuts[0], GenesisSecret, toAddr, coins, hours, fee)
	require.Equal(t, txn.Out[0].Address.String(), toAddr.String())

	if coins == GenesisCoins {
		// No change output
		require.Len(t, txn.Out, 1)
	} else {
		require.Len(t, txn.Out, 2)
		require.Equal(t, txn.Out[1].Address.String(), GenesisAddress.String())
	}

	return txn
}

func executeGenesisSpendTransaction(t *testing.T, db *bolt.DB, bc *Blockchain, txn coin.Transaction) coin.UxOut {
	block, err := bc.NewBlock(coin.Transactions{txn}, GenesisTime+TimeIncrement)
	require.NoError(t, err)

	sig := cipher.SignHash(block.HashHeader(), GenesisSecret)
	sb := coin.SignedBlock{
		Block: *block,
		Sig:   sig,
	}

	err = db.Update(func(tx *bolt.Tx) error {
		err = bc.ExecuteBlockWithTx(tx, &sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	uxOut, err := coin.CreateUnspent(block.Head, txn, 0)
	require.NoError(t, err)

	return uxOut
}

func requireSoftViolation(t *testing.T, msg string, err error) {
	require.Equal(t, NewErrTxnViolatesSoftConstraint(errors.New(msg)), err)
}

func requireHardViolation(t *testing.T, msg string, err error) {
	require.Equal(t, NewErrTxnViolatesHardConstraint(errors.New(msg)), err)
}

func TestVerifyTransactionAllConstraints(t *testing.T) {
	db, closeDB := testutil.PrepareDB(t)
	defer closeDB()

	store, err := blockdb.NewBlockchain(db, DefaultWalker)
	require.NoError(t, err)

	bc := &Blockchain{
		db:    db,
		store: store,
	}

	gb := addGenesisBlock(t, bc)

	toAddr := testutil.MakeAddress()
	coins := uint64(10e6)

	// create normal spending tx
	uxs := coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	tx := makeSpendTx(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins)
	err = bc.VerifySingleTxnAllConstraints(tx, DefaultMaxBlockSize)
	require.NoError(t, err)

	// Transaction size exceeds maxSize
	err = bc.VerifySingleTxnAllConstraints(tx, tx.Size()-1)
	requireSoftViolation(t, "Transaction size bigger than max block size", err)

	// Invalid transaction fee
	uxs = coin.CreateUnspents(gb.Head, gb.Body.Transactions[0])
	hours := uint64(0)
	for _, ux := range uxs {
		hours += ux.Body.Hours
	}
	tx = makeSpendTxWithHoursBurned(t, uxs, []cipher.SecKey{genSecret}, toAddr, coins, 0)
	err = bc.VerifySingleTxnAllConstraints(tx, DefaultMaxBlockSize)
	requireSoftViolation(t, "Transaction has zero coinhour fee", err)

	// Transaction locking is tested by TestVerifyTransactionIsLocked

	// Test invalid header hash
	originInnerHash := tx.InnerHash
	tx.InnerHash = cipher.SHA256{}
	err = bc.VerifySingleTxnAllConstraints(tx, DefaultMaxBlockSize)
	requireHardViolation(t, "Invalid header hash", err)

	// Set back the originInnerHash
	tx.InnerHash = originInnerHash

	// Create new block to spend the coins
	b, err := bc.NewBlock(coin.Transactions{tx}, genTime+100)
	require.NoError(t, err)

	// Add the block to blockchain
	err = bc.db.Update(func(tx *bolt.Tx) error {
		return bc.store.AddBlockWithTx(tx, &coin.SignedBlock{
			Block: *b,
			Sig:   cipher.SignHash(b.HashHeader(), genSecret),
		})
	})
	require.NoError(t, err)

	// A UxOut does not exist, it was already spent
	err = bc.VerifySingleTxnAllConstraints(tx, DefaultMaxBlockSize)
	expectedErr := NewErrTxnViolatesHardConstraint(blockdb.NewErrUnspentNotExist(tx.In[0].Hex()))
	require.Equal(t, expectedErr, err)

	// Check invalid sig
	uxs = coin.CreateUnspents(b.Head, tx)
	_, key := cipher.GenerateKeyPair()
	toAddr2 := testutil.MakeAddress()
	tx2 := makeSpendTx(t, uxs, []cipher.SecKey{key, key}, toAddr2, 5e6)
	err = bc.VerifySingleTxnAllConstraints(tx2, DefaultMaxBlockSize)
	requireHardViolation(t, "Signature not valid for output being spent", err)

	// Create lost coin transaction
	uxs2 := coin.CreateUnspents(b.Head, tx)
	toAddr3 := testutil.MakeAddress()
	lostCoinTx := makeLostCoinTx(coin.UxArray{uxs2[1]}, []cipher.SecKey{genSecret}, toAddr3, 10e5)
	err = bc.VerifySingleTxnAllConstraints(lostCoinTx, DefaultMaxBlockSize)
	requireHardViolation(t, "Transactions may not destroy coins", err)

	// Create transaction with duplicate UxOuts
	uxs = coin.CreateUnspents(b.Head, tx)
	toAddr4 := testutil.MakeAddress()
	dupUxOutTx := makeDuplicateUxOutTx(coin.UxArray{uxs[0]}, []cipher.SecKey{genSecret}, toAddr4, 1e6)
	err = bc.VerifySingleTxnAllConstraints(dupUxOutTx, DefaultMaxBlockSize)
	requireHardViolation(t, "Duplicate output in transaction", err)
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

	db, close := testutil.PrepareDB(t)
	defer close()

	_, s := cipher.GenerateKeyPair()

	// Setup blockchain
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address
	var coins = GenesisCoins
	var hours uint64 = 1e6
	var fee uint64 = 5e8

	txn := createGenesisSpendTransaction(t, bc, addr, coins, hours, fee)
	uxOut := executeGenesisSpendTransaction(t, db, bc, txn)

	// Create a transaction that spends from the locked address
	// The secret key for the locked address is obviously unavailable here,
	// instead, forge an invalid transaction.
	// Transaction.Verify() is called after TransactionIsLocked(),
	// so for this test it doesn't matter if transaction signature is wrong
	_, _, randomAddress := MakeAddress()
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

	uxIn, err := bc.Unspent().GetArray(txn.In)
	require.NoError(t, err)

	head, err := bc.Head()
	require.NoError(t, err)

	err = VerifySingleTxnSoftConstraints(txn, head.Time(), uxIn, DefaultMaxBlockSize)
	if expectedErr == nil {
		require.NoError(t, err)
	} else {
		requireSoftViolation(t, expectedErr.Error(), err)
	}
}
