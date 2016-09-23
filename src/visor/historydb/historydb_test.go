package historydb_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/historydb"
	"github.com/stretchr/testify/assert"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024
	transactionBkt       = []byte("transactions")
	outputBkt            = []byte("outputs")
	addressInBkt         = []byte("address_in")
	addressOutBkt        = []byte("address_out")
)

var _genTime uint64 = 1000
var _incTime uint64 = 3600 * 1000
var _genCoins uint64 = 1000e6
var _genCoinHours uint64 = 1000 * 1000

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (*bolt.DB, func(), error) {
	dbName := fmt.Sprintf("%ddb", rand.Int31n(10000))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), dbName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return nil, teardown, err
	}

	util.DataDir = tmpDir
	db, err := historydb.NewDB()
	if err != nil {
		t.Fatal(err)
	}

	teardown = func() {
		db.Close()
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	return db, teardown, nil
}

type FakeTree struct {
	blocks []*coin.Block
}

func (ft *FakeTree) AddBlock(b *coin.Block) error {
	ft.blocks = append(ft.blocks, b)
	return nil
}

func (ft *FakeTree) RemoveBlock(b *coin.Block) error {
	return nil
}

func (ft *FakeTree) GetBlock(hash cipher.SHA256) *coin.Block {
	for _, b := range ft.blocks {
		if b.HashHeader() == hash {
			return b
		}
	}
	return nil
}

func (ft *FakeTree) GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block {
	if dep >= uint64(len(ft.blocks)) {
		return nil
	}
	return ft.blocks[int(dep)]
}

func TestProcessGenesisBlock(t *testing.T) {
	db, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	ft := FakeTree{}
	bc := coin.NewBlockchain(&ft, nil)
	gb := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)

	hisDB, err := historydb.New(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := hisDB.ProcessBlock(&gb); err != nil {
		t.Fatal(err)
	}

	// check transactions bucket.
	var tx historydb.Transaction
	txHash := gb.Body.Transactions[0].Hash()
	if err := getBucketValue(db, transactionBkt, txHash[:], &tx); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tx.Transaction, gb)
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
