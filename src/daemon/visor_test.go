package daemon

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
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

func MakeTransactionForChain(t *testing.T, db *dbutil.DB, vs *Visor, ux coin.UxOut, sec cipher.SecKey, toAddr cipher.Address, amt, hours, fee uint64) coin.Transaction {
	headTime, err := vs.v.GetHeadBlockTime()
	require.NoError(t, err)

	chrs := ux.CoinHours(headTime)

	require.Equal(t, cipher.AddressFromPubKey(cipher.PubKeyFromSecKey(sec)), ux.Body.Address)

	knownUx, err := vs.v.GetUnspentOutput(ux.Hash())
	require.NoError(t, err)
	require.NotNil(t, knownUx)
	require.Equal(t, *knownUx, ux)

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

	err = db.View(func(tx *bolt.Tx) error {
		err := vs.v.Blockchain.VerifyTransaction(tx, txn)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	return txn
}

func MakeBlockchain(t *testing.T, db *dbutil.DB, seckey cipher.SecKey) *visor.Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := visor.NewBlockchain(db, pubkey, visor.BlockchainOptions{})
	require.NoError(t, err)
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	require.NoError(t, err, "create genesis block failed: %v", err)

	sig := cipher.SignHash(gb.HashHeader(), seckey)
	err = db.Update(func(tx *bolt.Tx) error {
		return b.ExecuteBlock(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   sig,
		})
	})
	require.NoError(t, err)

	return b
}

func MakeAddress() (cipher.PubKey, cipher.SecKey, cipher.Address) {
	p, s := cipher.GenerateKeyPair()
	a := cipher.AddressFromPubKey(p)
	return p, s, a
}

func setupSimpleVisor(t *testing.T, db *dbutil.DB, bc *visor.Blockchain) *Visor {
	visorCfg := NewVisorConfig()
	visorCfg.DisableNetworking = true
	visorCfg.Config.DBPath = db.Path()

	unconfirmed, err := visor.NewUnconfirmedTxnPool(db)
	require.NoError(t, err)

	return &Visor{
		Config: visorCfg,
		v: &visor.Visor{
			Config:      visorCfg.Config,
			Unconfirmed: unconfirmed,
			Blockchain:  bc,
			DB:          db,
		},
	}
}

func createGenesisSpendTransaction(t *testing.T, db *dbutil.DB, vs *Visor, toAddr cipher.Address, coins, hours, fee uint64) coin.Transaction { // nolint: unparam
	uxOuts, err := vs.v.GetUnspentOutputs()
	require.NoError(t, err)
	require.Len(t, uxOuts, 1)

	txn := MakeTransactionForChain(t, db, vs, uxOuts[0], GenesisSecret, toAddr, coins, hours, fee)
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

func executeGenesisSpendTransaction(t *testing.T, db *dbutil.DB, bc *visor.Blockchain, txn coin.Transaction) coin.UxOut {
	var block *coin.Block

	err := db.Update(func(tx *bolt.Tx) error {
		var err error
		block, err = bc.NewBlock(tx, coin.Transactions{txn}, GenesisTime+TimeIncrement)
		require.NoError(t, err)

		sig := cipher.SignHash(block.HashHeader(), GenesisSecret)
		sb := coin.SignedBlock{
			Block: *block,
			Sig:   sig,
		}

		err = bc.ExecuteBlock(tx, &sb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	uxOut, err := coin.CreateUnspent(block.Head, txn, 0)
	require.NoError(t, err)

	return uxOut
}

func TestVerifyTransactionIsLocked(t *testing.T) {
	for _, addr := range visor.GetLockedDistributionAddresses() {
		t.Run(fmt.Sprintf("IsLocked:%s", addr), func(t *testing.T) {
			testVerifyTransactionAddressLocking(t, addr, "Transaction has locked address inputs")
		})
	}
}

func TestVerifyTransactionIsUnlocked(t *testing.T) {
	// Note: we can't sign transactions that spend the outputs of unlocked addresses,
	// because these are real, hardcoded addresses for which we don't have the secret key.
	// Validation is expected to fail, but it will fail on invalid header hash, rather than
	// due to locked address inputs.
	for _, addr := range visor.GetUnlockedDistributionAddresses() {
		t.Run(fmt.Sprintf("IsUnlocked:%s", addr), func(t *testing.T) {
			testVerifyTransactionAddressLocking(t, addr, "Transaction Verification Failed, Invalid header hash")
		})
	}
}

func testVerifyTransactionAddressLocking(t *testing.T, toAddr, errMsg string) {
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

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	txn := createGenesisSpendTransaction(t, db, v, addr, coins, hours, fee)
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

	// Call verifyTransaction
	err = v.verifyTransaction(txn)
	testutil.RequireError(t, err, errMsg)
}

func TestVerifyTransactionInvalidFee(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var fee uint64
	_, _, addr := MakeAddress()

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	txn := createGenesisSpendTransaction(t, db, v, addr, coins, hours, fee)

	// Call verifyTransaction
	err := v.verifyTransaction(txn)
	testutil.RequireError(t, err, "Transaction coinhour fee minimum not met")
}

func TestVerifyTransactionInvalidSignature(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours uint64
	var fee uint64
	_, _, addr := MakeAddress()

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	txn := createGenesisSpendTransaction(t, db, v, addr, coins, hours, fee)

	// Invalidate signatures
	txn.Sigs = nil

	// Call verifyTransaction
	err := v.verifyTransaction(txn)
	testutil.RequireError(t, err, "Transaction Verification Failed, Invalid number of signatures")
}

func TestInjectValidTransaction(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()

	_, s := cipher.GenerateKeyPair()
	// Setup blockchain
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours uint64
	var fee uint64
	_, _, addr := MakeAddress()

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	txn := createGenesisSpendTransaction(t, db, v, addr, coins, hours, fee)

	// The unconfirmed pool should be empty
	err := db.Update(func(tx *bolt.Tx) error {
		txns, err := v.v.Unconfirmed.RawTxns(tx)
		require.NoError(t, err)
		require.Len(t, txns, 0)
		return nil
	})
	require.NoError(t, err)

	err = v.verifyAndInjectTransaction(txn)
	require.NoError(t, err)

	// The transaction should appear in the unconfirmed pool
	err = db.Update(func(tx *bolt.Tx) error {
		txns, err := v.v.Unconfirmed.RawTxns(tx)
		require.NoError(t, err)
		require.Len(t, txns, 1)
		require.Equal(t, txns[0], txn)
		return nil
	})
	require.NoError(t, err)
}

func TestInjectInvalidTransaction(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var fee uint64
	_, _, addr := MakeAddress()

	// Setup a minimal visor
	v := setupSimpleVisor(t, db, bc)

	txn := createGenesisSpendTransaction(t, db, v, addr, coins, hours, fee)

	// The unconfirmed pool should be empty
	err := db.View(func(tx *bolt.Tx) error {
		txns, err := v.v.Unconfirmed.RawTxns(tx)
		require.NoError(t, err)
		require.Len(t, txns, 0)
		return nil
	})
	require.NoError(t, err)

	err = v.verifyAndInjectTransaction(txn)
	testutil.RequireError(t, err, "Transaction coinhour fee minimum not met")

	// The transaction should not appear in the unconfirmed pool
	err = db.View(func(tx *bolt.Tx) error {
		txns, err := v.v.Unconfirmed.RawTxns(tx)
		require.NoError(t, err)
		require.Len(t, txns, 0)
		return nil
	})
	require.NoError(t, err)
}

func TestSplitHashes(t *testing.T) {
	hashes := make([]cipher.SHA256, 10)
	for i := 0; i < 10; i++ {
		hashes[i] = cipher.SumSHA256(cipher.RandByte(512))
	}

	testCases := []struct {
		name  string
		init  []cipher.SHA256
		n     int
		array [][]cipher.SHA256
	}{
		{
			"has one odd",
			hashes[:],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
				[]cipher.SHA256{
					hashes[3],
					hashes[4],
					hashes[5],
				},
				[]cipher.SHA256{
					hashes[6],
					hashes[7],
					hashes[8],
				},
				[]cipher.SHA256{
					hashes[9],
				},
			},
		},
		{
			"only one value",
			hashes[:1],
			1,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
				},
			},
		},
		{
			"empty value",
			hashes[:0],
			0,
			[][]cipher.SHA256{},
		},
		{
			"with 3 value",
			hashes[:3],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
			},
		},
		{
			"with 8 value",
			hashes[:8],
			3,
			[][]cipher.SHA256{
				[]cipher.SHA256{
					hashes[0],
					hashes[1],
					hashes[2],
				},
				[]cipher.SHA256{
					hashes[3],
					hashes[4],
					hashes[5],
				},
				[]cipher.SHA256{
					hashes[6],
					hashes[7],
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rlt := divideHashes(tc.init, tc.n)
			require.Equal(t, tc.array, rlt)
		})
	}
}
