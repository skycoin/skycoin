package historydb

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024
	blockBkt             = []byte("blocks")
	transactionBkt       = []byte("transactions")
	outputBkt            = []byte("uxouts")
	addressInBkt         = []byte("address_in")
	addressOutBkt        = []byte("address_out")
)

var _genTime uint64 = 1000
var _incTime uint64 = 3600 * 1000
var _genCoins uint64 = 1000e6
var _genCoinHours uint64 = 1000 * 1000

func _feeCalc(t *coin.Transaction) (uint64, error) {
	return 0, nil
}

func _makeFeeCalc(fee uint64) coin.FeeCalculator {
	return func(t *coin.Transaction) (uint64, error) {
		return fee, nil
	}
}

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

type fakeBlockchain struct {
	blocks  []coin.Block
	unspent map[string]coin.UxOut
	uxhash  cipher.SHA256
}

func newBlockchain(db *bolt.DB) *fakeBlockchain {
	return &fakeBlockchain{
		unspent: make(map[string]coin.UxOut),
	}
}

func (fbc fakeBlockchain) GetBlockInDepth(dep uint64) *coin.Block {
	if dep >= uint64(len(fbc.blocks)) {
		panic(fmt.Sprintf("block depth: %d overflow", dep))
	}

	return &fbc.blocks[dep]
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
	for _, tx := range txns {
		// Remove spent outputs
		for _, id := range tx.In {
			ux := fbc.unspent[id.Hex()]
			fbc.uxhash = fbc.uxhash.Xor(ux.SnapshotHash())
			delete(fbc.unspent, id.Hex())

		}
		fbc.deleteUxOut(tx.In)
		// Create new outputs
		txUxs := coin.CreateUnspents(b.Head, tx)
		for i := range txUxs {
			fbc.addUxOut(txUxs[i])
		}
		uxs = append(uxs, txUxs...)
	}

	b.Head.PrevHash = fbc.Head().HashHeader()
	fbc.blocks = append(fbc.blocks, *b)

	return uxs, nil
}

func (fbc *fakeBlockchain) CreateGenesisBlock(genesisAddr cipher.Address, genesisCoins, timestamp uint64) coin.Block {
	txn := coin.Transaction{}
	txn.PushOutput(genesisAddr, genesisCoins, genesisCoins)
	body := coin.BlockBody{Transactions: coin.Transactions{txn}}
	prevHash := cipher.SHA256{}
	head := coin.BlockHeader{
		Time:     timestamp,
		BodyHash: body.Hash(),
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
	db, teardown := testutil.PrepareDB(t)
	defer teardown()

	bc := newBlockchain(db)
	gb := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)
	hisDB, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := hisDB.ProcessBlock(&gb); err != nil {
		t.Fatal(err)
	}

	// check transactions bucket.
	var tx Transaction
	txHash := gb.Body.Transactions[0].Hash()
	if err := getBucketValue(db, transactionBkt, txHash[:], &tx); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tx.Tx, gb.Body.Transactions[0])

	// check address in
	outID := []cipher.SHA256{}
	if err := getBucketValue(db, addressInBkt, genAddress.Bytes(), &outID); err != nil {
		t.Fatal(err)
	}

	ux, ok := bc.unspent[outID[0].Hex()]
	require.True(t, ok)
	require.Equal(t, outID[0], ux.Hash())

	// check outputs
	output := UxOut{}
	err = getBucketValue(db, outputBkt, outID[0][:], &output)
	require.Nil(t, err)

	assert.Equal(t, output.Out, ux)
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

func getUx(bc Blockchainer, seq uint64, txID cipher.SHA256, addr string) (*coin.UxOut, error) {
	b := bc.GetBlockInDepth(seq)
	if b == nil {
		return nil, fmt.Errorf("no block in depth:%v", seq)
	}
	tx, ok := b.GetTransaction(txID)
	if !ok {
		return nil, errors.New("found transaction failed")
	}
	uxs := coin.CreateUnspents(b.Head, tx)
	for _, u := range uxs {
		if u.Body.Address.String() == addr {
			return &u, nil
		}
	}
	return nil, nil
}

func TestProcessBlock(t *testing.T) {
	db, teardown := testutil.PrepareDB(t)
	defer teardown()
	bc := newBlockchain(db)
	gb := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)

	// create
	hisDB, err := New(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := hisDB.ProcessBlock(&gb); err != nil {
		t.Fatal(err)
	}
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
					Coins:  _genCoins - 10e6,
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

func testEngine(t *testing.T, tds []testData, bc *fakeBlockchain, hdb *HistoryDB, db *bolt.DB) {
	for i, td := range tds {
		b, tx, err := addBlock(bc, td, _incTime*(uint64(i)+1))
		if err != nil {
			t.Fatal(err)
		}
		// update the next block test data.
		if i+1 < len(tds) {
			// update UxOut of next test data.
			tds[i+1].Vin.TxID = tx.Hash()
			tds[i+1].PreBlockHash = b.HashHeader()
		}

		if err := hdb.ProcessBlock(b); err != nil {
			t.Fatal(err)
		}

		// check tx
		txInBkt := Transaction{}
		k := tx.Hash()
		if err := getBucketValue(db, transactionBkt, k[:], &txInBkt); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, &txInBkt.Tx, tx)

		// check outputs
		for _, o := range td.Vouts {
			ux, err := getUx(bc, uint64(i+1), tx.Hash(), o.ToAddr)
			if err != nil {
				t.Fatal(err)
			}

			uxInDB := UxOut{}
			uxKey := ux.Hash()
			if err = getBucketValue(db, outputBkt, uxKey[:], &uxInDB); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, &uxInDB.Out, ux)
		}

		// check addr in
		for _, o := range td.Vouts {
			addr := cipher.MustDecodeBase58Address(o.ToAddr)
			uxHashes := []cipher.SHA256{}
			if err := getBucketValue(db, addressInBkt, addr.Bytes(), &uxHashes); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, len(uxHashes), td.AddrInNum[o.ToAddr])
		}
	}
}

func addBlock(bc *fakeBlockchain, td testData, tm uint64) (*coin.Block, *coin.Transaction, error) {
	tx := coin.Transaction{}
	// get unspent output
	ux, err := getUx(bc, td.Vin.BlockSeq, td.Vin.TxID, td.Vin.Addr)
	if err != nil {
		return nil, nil, err
	}
	if ux == nil {
		return nil, nil, errors.New("no unspent output")
	}

	tx.PushInput(ux.Hash())
	for _, o := range td.Vouts {
		addr, err := cipher.DecodeBase58Address(o.ToAddr)
		if err != nil {
			return nil, nil, err
		}
		tx.PushOutput(addr, o.Coins, o.Hours)
	}

	sigKey := cipher.MustSecKeyFromHex(td.Vin.SigKey)
	tx.SignInputs([]cipher.SecKey{sigKey})
	tx.UpdateHeader()
	if err := bc.VerifyTransaction(tx); err != nil {
		return nil, nil, err
	}
	preBlock := bc.GetBlock(td.PreBlockHash)
	b := newBlock(*preBlock, tm, bc.uxhash, coin.Transactions{tx}, _feeCalc)

	// uxs, err := bc.ExecuteBlock(&b)
	_, err = bc.ExecuteBlock(&b)
	if err != nil {
		return nil, nil, err
	}
	return &b, &tx, nil
}

func getBucketValue(db *bolt.DB, name []byte, key []byte, value interface{}) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(name)
		bin := b.Get(key)
		if bin == nil {
			value = nil
			return nil
		}
		return encoder.DeserializeRaw(bin, value)
	})
}

func newBlock(prev coin.Block, currentTime uint64, uxHash cipher.SHA256,
	txns coin.Transactions, calc coin.FeeCalculator) coin.Block {
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

func newBlockHeader(prev coin.BlockHeader, uxHash cipher.SHA256, currentTime,
	fee uint64, body coin.BlockBody) coin.BlockHeader {
	prevHash := prev.Hash()
	return coin.BlockHeader{
		BodyHash: body.Hash(),
		Version:  prev.Version,
		PrevHash: prevHash,
		Time:     currentTime,
		BkSeq:    prev.BkSeq + 1,
		Fee:      fee,
		UxHash:   uxHash,
	}
}
