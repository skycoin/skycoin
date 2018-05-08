package daemon

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
	"github.com/skycoin/skycoin/src/visor/dbutil"
)

// TODO -- most of these tests should be merged into visor/visor_test.go,
// daemon.Visor is only a thin wrapper around visor.Visor

func prepareDB(t *testing.T) (*dbutil.DB, func()) {
	db, shutdown := testutil.PrepareDB(t)

	err := visor.CreateBuckets(db)
	if err != nil {
		shutdown()
		t.Fatalf("CreateBuckets failed: %v", err)
	}

	return db, shutdown
}

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
func MakeBlockchain(t *testing.T, db *dbutil.DB, seckey cipher.SecKey) *visor.Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := visor.NewBlockchain(db, visor.BlockchainConfig{
		Pubkey: pubkey,
	})
	require.NoError(t, err)
	gb, err := coin.NewGenesisBlock(GenesisAddress, GenesisCoins, GenesisTime)
	if err != nil {
		panic(fmt.Errorf("create genesis block failed: %v", err))
	}

	sig := cipher.SignHash(gb.HashHeader(), seckey)
	db.Update(func(tx *dbutil.Tx) error {
		return b.ExecuteBlock(tx, &coin.SignedBlock{
			Block: *gb,
			Sig:   sig,
		})
	})
	return b
}

// CreateGenesisSpendTransaction creates the initial post-genesis transaction that moves genesis coins to another address
func CreateGenesisSpendTransaction(t *testing.T, db *dbutil.DB, bc *visor.Blockchain, toAddr cipher.Address, coins, hours, fee uint64) coin.Transaction {
	var txn coin.Transaction
	err := db.View(func(tx *dbutil.Tx) error {
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

func makeTransactionForChain(t *testing.T, tx *dbutil.Tx, bc *visor.Blockchain, ux coin.UxOut, sec cipher.SecKey, toAddr cipher.Address, amt, hours, fee uint64) coin.Transaction {
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

func setupSimpleVisor(t *testing.T, db *dbutil.DB, bc *visor.Blockchain) *Visor {
	visorCfg := NewVisorConfig()
	visorCfg.DisableNetworking = true
	visorCfg.Config.DBPath = db.Path()

	pool, err := visor.NewUnconfirmedTxnPool(db)
	require.NoError(t, err)

	return &Visor{
		Config: visorCfg,
		v: &visor.Visor{
			Config:      visorCfg.Config,
			Unconfirmed: pool,
			Blockchain:  bc,
			DB:          db,
		},
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

	_, softErr, err := v.InjectTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, visor.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)
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

	_, softErr, err := v.InjectTransaction(txn)
	require.Nil(t, softErr)
	testutil.RequireError(t, err, visor.NewErrTxnViolatesHardConstraint(errors.New("Invalid number of signatures")).Error())
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
	txns, err := v.v.GetAllUnconfirmedTxns()
	require.NoError(t, err)
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectTransaction(txn)
	require.Nil(t, softErr)
	require.NoError(t, err)

	// The transaction should appear in the unconfirmed pool
	txns, err = v.v.GetAllUnconfirmedTxns()
	require.NoError(t, err)
	require.Len(t, txns, 1)
	require.Equal(t, txns[0].Txn, txn)
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
	txns, err := v.v.GetAllUnconfirmedTxns()
	require.NoError(t, err)
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, visor.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)

	// The transaction should appear in the unconfirmed pool
	txns, err = v.v.GetAllUnconfirmedTxns()
	require.NoError(t, err)
	require.Len(t, txns, 1)
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
