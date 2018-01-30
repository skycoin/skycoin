package daemon

import (
	"errors"
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/daemon/strand"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/skycoin/skycoin/src/util/fee"
	"github.com/skycoin/skycoin/src/visor"
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

func MakeTransactionForChain(t *testing.T, bc *visor.Blockchain, ux coin.UxOut, sec cipher.SecKey, toAddr cipher.Address, amt, hours, fee uint64) coin.Transaction {
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

func MakeBlockchain(t *testing.T, db *bolt.DB, seckey cipher.SecKey) *visor.Blockchain {
	pubkey := cipher.PubKeyFromSecKey(seckey)
	b, err := visor.NewBlockchain(db, pubkey)
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

func createGenesisSpendTransaction(t *testing.T, bc *visor.Blockchain, toAddr cipher.Address, coins, hours, fee uint64) coin.Transaction {
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

func setupSimpleVisor(db *bolt.DB, bc *visor.Blockchain) *Visor {
	visorCfg := NewVisorConfig()
	visorCfg.DisableNetworking = true
	visorCfg.Config.DBPath = db.Path()
	return &Visor{
		Config: visorCfg,
		v: &visor.Visor{
			Config:      visorCfg.Config,
			Unconfirmed: visor.NewUnconfirmedTxnPool(db),
			Blockchain:  bc,
		},
		reqC: make(chan strand.Request, 10),
	}
}

func TestVerifyTransactionInvalidFee(t *testing.T) {
	// Test that a soft constraint is enforced
	// Full verification tests are in visor/blockchain_verify_test.go
	db, close := testutil.PrepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var f uint64
	_, _, addr := MakeAddress()

	txn := createGenesisSpendTransaction(t, bc, addr, coins, hours, f)

	// Setup a minimal visor
	v := setupSimpleVisor(db, bc)
	errC := make(chan error)
	go v.processRequests(errC)
	defer func() {
		errC <- errors.New("stop")
	}()

	_, softErr, err := v.InjectTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, visor.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)
}

func TestVerifyTransactionInvalidSignature(t *testing.T) {
	// Test that a hard constraint is enforced
	// Full verification tests are in visor/blockchain_verify_test.go
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

	txn := createGenesisSpendTransaction(t, bc, addr, coins, hours, fee)

	// Invalidate signatures
	txn.Sigs = nil

	// Setup a minimal visor
	v := setupSimpleVisor(db, bc)
	errC := make(chan error)
	go v.processRequests(errC)
	defer func() {
		errC <- errors.New("stop")
	}()

	_, softErr, err := v.InjectTransaction(txn)
	require.Nil(t, softErr)
	testutil.RequireError(t, err, visor.NewErrTxnViolatesHardConstraint(errors.New("Invalid number of signatures")).Error())
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

	txn := createGenesisSpendTransaction(t, bc, addr, coins, hours, fee)

	// Setup a minimal visor
	v := setupSimpleVisor(db, bc)
	errC := make(chan error)
	go v.processRequests(errC)
	defer func() {
		errC <- errors.New("stop")
	}()

	// The unconfirmed pool should be empty
	txns := v.v.Unconfirmed.RawTxns()
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectTransaction(txn)
	require.Nil(t, softErr)
	require.NoError(t, err)

	// The transaction should appear in the unconfirmed pool
	txns = v.v.Unconfirmed.RawTxns()
	require.Len(t, txns, 1)
	require.Equal(t, txns[0], txn)
}

func TestInjectTransactionSoftViolationNoFee(t *testing.T) {
	db, close := testutil.PrepareDB(t)
	defer close()

	// Setup blockchain
	_, s := cipher.GenerateKeyPair()
	bc := MakeBlockchain(t, db, s)

	// Send coins to the initial address, with invalid fee
	var coins = GenesisCoins
	var hours = GenesisCoinHours * 1e3
	var f uint64
	_, _, addr := MakeAddress()

	txn := createGenesisSpendTransaction(t, bc, addr, coins, hours, f)

	// Setup a minimal visor
	v := setupSimpleVisor(db, bc)
	errC := make(chan error)
	go v.processRequests(errC)
	defer func() {
		errC <- errors.New("stop")
	}()

	// The unconfirmed pool should be empty
	txns := v.v.Unconfirmed.RawTxns()
	require.Len(t, txns, 0)

	// Call injectTransaction
	_, softErr, err := v.InjectTransaction(txn)
	require.NoError(t, err)
	require.NotNil(t, softErr)
	require.Equal(t, visor.NewErrTxnViolatesSoftConstraint(fee.ErrTxnNoFee), *softErr)

	// The transaction should appear in the unconfirmed pool
	txns = v.v.Unconfirmed.RawTxns()
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
