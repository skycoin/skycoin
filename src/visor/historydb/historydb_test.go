package historydb

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/testutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
	"github.com/stretchr/testify/require"
)

func prepareDB(t *testing.T) (*dbutil.DB, func()) {
	db, shutdown := testutil.PrepareDB(t)

	err := db.Update("", func(tx *dbutil.Tx) error {
		return CreateBuckets(tx)
	})
	if err != nil {
		shutdown()
		t.Fatalf("CreateBuckets failed: %v", err)
	}

	return db, shutdown
}

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
)

var genTime uint64 = 1000
var incTime uint64 = 3600 * 1000
var genCoins uint64 = 1000e6

func feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

// Blockchainer interface for isolating the detail of blockchain.
type Blockchainer interface {
	Head() *coin.Block
	GetBlockInDepth(depth uint64) *coin.Block
	ExecuteBlock(b *coin.Block) (coin.UxArray, error)
	CreateGenesisBlock(genAddress cipher.Address, genCoins, timestamp uint64) coin.Block
	VerifyTransaction(tx coin.Transaction) error
	GetBlock(hash cipher.SHA256) *coin.Block
}

type fakeBlockchain struct {
	blocks  []coin.Block
	unspent map[string]coin.UxOut
	uxhash  cipher.SHA256
}

func newBlockchain() *fakeBlockchain {
	return &fakeBlockchain{
		unspent: make(map[string]coin.UxOut),
	}
}

func (fbc fakeBlockchain) GetBlockInDepth(depth uint64) *coin.Block {
	if depth >= uint64(len(fbc.blocks)) {
		panic(fmt.Sprintf("block depth: %d overflow", depth))
	}

	return &fbc.blocks[depth]
}

func (fbc fakeBlockchain) Head() *coin.Block {
	l := len(fbc.blocks)
	if l == 0 {
		return nil
	}

	return &fbc.blocks[l-1]
}

func (fbc *fakeBlockchain) deleteUxOut(uxids []cipher.SHA256) {
	for _, id := range uxids {
		ux := fbc.unspent[id.Hex()]
		fbc.uxhash = fbc.uxhash.Xor(ux.SnapshotHash())
		delete(fbc.unspent, id.Hex())
	}
}

func (fbc *fakeBlockchain) addUxOut(ux coin.UxOut) {
	fbc.uxhash = fbc.uxhash.Xor(ux.SnapshotHash())
	fbc.unspent[ux.Hash().Hex()] = ux
}

func (fbc *fakeBlockchain) ExecuteBlock(b *coin.Block) (coin.UxArray, error) {
	var uxs coin.UxArray
	txns := b.Body.Transactions
	for _, txn := range txns {
		// Remove spent outputs
		for _, id := range txn.In {
			ux := fbc.unspent[id.Hex()]
			fbc.uxhash = fbc.uxhash.Xor(ux.SnapshotHash())
			delete(fbc.unspent, id.Hex())

		}
		fbc.deleteUxOut(txn.In)
		// Create new outputs
		txnUxs := coin.CreateUnspents(b.Head, txn)
		for i := range txnUxs {
			fbc.addUxOut(txnUxs[i])
		}
		uxs = append(uxs, txnUxs...)
	}

	b.Head.PrevHash = fbc.Head().HashHeader()
	fbc.blocks = append(fbc.blocks, *b)

	return uxs, nil
}

func (fbc *fakeBlockchain) CreateGenesisBlock(genesisAddr cipher.Address, genesisCoins, timestamp uint64) coin.Block {
	txn := coin.Transaction{}
	err := txn.PushOutput(genesisAddr, genesisCoins, genesisCoins)
	if err != nil {
		panic(err)
	}
	body := coin.BlockBody{Transactions: coin.Transactions{txn}}
	prevHash := cipher.SHA256{}
	bodyHash := body.Hash()
	head := coin.BlockHeader{
		Time:     timestamp,
		BodyHash: bodyHash,
		PrevHash: prevHash,
		BkSeq:    0,
		Version:  0,
		Fee:      0,
		UxHash:   cipher.SHA256{},
	}
	b := coin.Block{
		Head: head,
		Body: body,
	}
	// b.Body.Transactions[0].UpdateHeader()
	fbc.blocks = append(fbc.blocks, b)
	ux := coin.UxOut{
		Head: coin.UxHead{
			Time:  timestamp,
			BkSeq: 0,
		},
		Body: coin.UxBody{
			SrcTransaction: txn.InnerHash, //user inner hash
			Address:        genesisAddr,
			Coins:          genesisCoins,
			Hours:          genesisCoins, // Allocate 1 coin hour per coin
		},
	}
	fbc.addUxOut(ux)
	return b
}

func (fbc fakeBlockchain) VerifyTransaction(tx coin.Transaction) error {
	return nil
}

func (fbc fakeBlockchain) GetBlock(hash cipher.SHA256) *coin.Block {
	for _, b := range fbc.blocks {
		if b.HashHeader() == hash {
			return &b
		}
	}
	return nil
}

func TestProcessGenesisBlock(t *testing.T) {
	db, teardown := prepareDB(t)
	defer teardown()

	bc := newBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, genCoins, genTime)
	hisDB := New()

	err := db.Update("", func(tx *dbutil.Tx) error {
		err := hisDB.ParseBlock(tx, gb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)

	// check transactions bucket.
	var txn Transaction
	txnHash := gb.Body.Transactions[0].Hash()
	mustGetBucketValue(t, db, TransactionsBkt, txnHash[:], &txn)
	require.Equal(t, txn.Txn, gb.Body.Transactions[0])

	// check address in
	outID := []cipher.SHA256{}
	mustGetBucketValue(t, db, AddressUxBkt, genAddress.Bytes(), &outID)

	ux, ok := bc.unspent[outID[0].Hex()]
	require.True(t, ok)
	require.Equal(t, outID[0], ux.Hash())

	// check outputs
	output := UxOut{}
	mustGetBucketValue(t, db, UxOutsBkt, outID[0][:], &output)

	require.Equal(t, output.Out, ux)
}

type testData struct {
	PreBlockHash cipher.SHA256
	Vin          txIn
	Vouts        []txOut
	AddrInNum    map[string]int
	AddrOutNum   map[string]int
}

type txIn struct {
	SigKey   string
	Addr     string
	TxID     cipher.SHA256
	BlockSeq uint64
}

type txOut struct {
	ToAddr string
	Coins  uint64
	Hours  uint64
}

// findTxnInBlock looks up a coin.Transaction from a coin.Block.
// Returns the Transaction and whether it was found or not
func findTxnInBlock(b *coin.Block, txnHash cipher.SHA256) (coin.Transaction, bool) {
	txns := b.Body.Transactions
	for i := range txns {
		if txns[i].Hash() == txnHash {
			return txns[i], true
		}
	}
	return coin.Transaction{}, false
}

func getUx(bc Blockchainer, seq uint64, txID cipher.SHA256, addr string) (*coin.UxOut, error) {
	b := bc.GetBlockInDepth(seq)
	if b == nil {
		return nil, fmt.Errorf("no block in depth:%v", seq)
	}

	txn, ok := findTxnInBlock(b, txID)
	if !ok {
		return nil, errors.New("found transaction failed")
	}

	uxs := coin.CreateUnspents(b.Head, txn)
	for _, u := range uxs {
		if u.Body.Address.String() == addr {
			return &u, nil
		}
	}
	return nil, nil
}

func TestProcessBlock(t *testing.T) {
	db, teardown := prepareDB(t)
	defer teardown()
	bc := newBlockchain()
	gb := bc.CreateGenesisBlock(genAddress, genCoins, genTime)

	// create
	hisDB := New()

	err := db.Update("", func(tx *dbutil.Tx) error {
		err := hisDB.ParseBlock(tx, gb)
		require.NoError(t, err)
		return nil
	})
	require.NoError(t, err)
	/*

	                   |-2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS
	   genesisAddr  ==>|                                        |-2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS
	                   |-222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm ==>|
	                                                            |-222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm
	*/

	testData := []testData{
		{
			PreBlockHash: gb.HashHeader(),
			Vin: txIn{
				SigKey:   genSecret.Hex(),
				Addr:     genAddress.String(),
				TxID:     gb.Body.Transactions[0].Hash(),
				BlockSeq: 0,
			},
			Vouts: []txOut{
				{
					ToAddr: "2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS",
					Coins:  10e6,
					Hours:  100,
				},
				{
					ToAddr: "222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm",
					Coins:  genCoins - 10e6,
					Hours:  400,
				},
			},
			AddrInNum: map[string]int{
				"2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS": 1,
				"222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm": 1,
			},
			AddrOutNum: map[string]int{
				genAddress.String(): 1,
			},
		},
		{
			Vin: txIn{
				Addr:     "222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm",
				SigKey:   "62f4d675d991c41a2819d908a4fcf4ba44ff0c31564039e80508c9d68197f90c",
				BlockSeq: 1,
			},
			Vouts: []txOut{
				{
					ToAddr: "2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS",
					Coins:  10e6,
					Hours:  100,
				},
				{
					ToAddr: "222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm",
					Coins:  1000e6 - 20e6,
					Hours:  100,
				},
			},
			AddrInNum: map[string]int{
				"2RxP5N26GhDqHrP6SK45ZzEMSmSpeUeWxsS": 2,
				"222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm": 2,
			},
			AddrOutNum: map[string]int{
				"222uMeCeL1PbkJGZJDgAz5sib2uisv9hYUm": 1,
			},
		},
	}

	testEngine(t, testData, bc, hisDB, db)
}

func testEngine(t *testing.T, tds []testData, bc *fakeBlockchain, hdb *HistoryDB, db *dbutil.DB) {
	for i, td := range tds {
		b, txn, err := addBlock(bc, td, incTime*(uint64(i)+1))
		require.NoError(t, err)

		// update the next block test data.
		if i+1 < len(tds) {
			// update UxOut of next test data.
			tds[i+1].Vin.TxID = txn.Hash()
			tds[i+1].PreBlockHash = b.HashHeader()
		}

		err = db.Update("", func(tx *dbutil.Tx) error {
			err := hdb.ParseBlock(tx, *b)
			require.NoError(t, err)
			return nil
		})
		require.NoError(t, err)

		// check txn
		txnInBkt := Transaction{}
		k := txn.Hash()
		mustGetBucketValue(t, db, TransactionsBkt, k[:], &txnInBkt)
		require.Equal(t, &txnInBkt.Txn, txn)

		// check outputs
		for _, o := range td.Vouts {
			ux, err := getUx(bc, uint64(i+1), txn.Hash(), o.ToAddr)
			require.NoError(t, err)

			uxInDB := UxOut{}
			uxKey := ux.Hash()
			mustGetBucketValue(t, db, UxOutsBkt, uxKey[:], &uxInDB)
			require.Equal(t, &uxInDB.Out, ux)
		}

		// check addr in
		for _, o := range td.Vouts {
			addr := cipher.MustDecodeBase58Address(o.ToAddr)
			uxHashes := []cipher.SHA256{}
			mustGetBucketValue(t, db, AddressUxBkt, addr.Bytes(), &uxHashes)
			require.Equal(t, len(uxHashes), td.AddrInNum[o.ToAddr])
		}
	}
}

func addBlock(bc *fakeBlockchain, td testData, tm uint64) (*coin.Block, *coin.Transaction, error) {
	txn := coin.Transaction{}
	// get unspent output
	ux, err := getUx(bc, td.Vin.BlockSeq, td.Vin.TxID, td.Vin.Addr)
	if err != nil {
		return nil, nil, err
	}
	if ux == nil {
		return nil, nil, errors.New("no unspent output")
	}

	if err := txn.PushInput(ux.Hash()); err != nil {
		return nil, nil, err
	}

	for _, o := range td.Vouts {
		addr, err := cipher.DecodeBase58Address(o.ToAddr)
		if err != nil {
			return nil, nil, err
		}
		if err := txn.PushOutput(addr, o.Coins, o.Hours); err != nil {
			return nil, nil, err
		}
	}

	sigKey := cipher.MustSecKeyFromHex(td.Vin.SigKey)
	txn.SignInputs([]cipher.SecKey{sigKey})
	if err := txn.UpdateHeader(); err != nil {
		return nil, nil, err
	}
	if err := bc.VerifyTransaction(txn); err != nil {
		return nil, nil, err
	}
	preBlock := bc.GetBlock(td.PreBlockHash)
	b := newBlock(*preBlock, tm, bc.uxhash, coin.Transactions{txn}, feeCalc)

	// uxs, err := bc.ExecuteBlock(&b)
	_, err = bc.ExecuteBlock(&b)
	if err != nil {
		return nil, nil, err
	}
	return &b, &txn, nil
}

func mustGetBucketValue(t *testing.T, db *dbutil.DB, name []byte, key []byte, value interface{}) {
	err := db.View("", func(tx *dbutil.Tx) error {
		ok, err := dbutil.GetBucketObjectDecoded(tx, name, key, value)
		require.NoError(t, err)
		require.True(t, ok)
		return err
	})
	require.NoError(t, err)
}

func newBlock(prev coin.Block, currentTime uint64, uxHash cipher.SHA256, txns coin.Transactions, calc coin.FeeCalculator) coin.Block {
	if len(txns) == 0 {
		log.Panic("Refusing to create block with no transactions")
	}
	fee, err := txns.Fees(calc)
	if err != nil {
		// This should have been caught earlier
		log.Panicf("Invalid transaction fees: %v", err)
	}
	body := coin.BlockBody{Transactions: txns}
	return coin.Block{
		Head: newBlockHeader(prev.Head, uxHash, currentTime, fee, body),
		Body: body,
	}
}

func newBlockHeader(prev coin.BlockHeader, uxHash cipher.SHA256, currentTime, fee uint64, body coin.BlockBody) coin.BlockHeader {
	prevHash := prev.Hash()
	bodyHash := body.Hash()
	return coin.BlockHeader{
		BodyHash: bodyHash,
		Version:  prev.Version,
		PrevHash: prevHash,
		Time:     currentTime,
		BkSeq:    prev.BkSeq + 1,
		Fee:      fee,
		UxHash:   uxHash,
	}
}
