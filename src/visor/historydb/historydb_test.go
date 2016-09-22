package historydb

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

var (
	genPublic, genSecret = cipher.GenerateKeyPair()
	genAddress           = cipher.AddressFromPubKey(genPublic)
	testMaxSize          = 1024 * 1024
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

func setup(t *testing.T) (*HistoryDB, func(), error) {
	dbName := fmt.Sprintf("%d.db", rand.Int31n(10000))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), dbName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return nil, teardown, err
	}

	util.DataDir = tmpDir
	hdb := &HistoryDB{}
	hdb.Start()

	teardown = func() {
		hdb.Stop()
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	return nil, teardown, nil
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

func TestProcessBlock(t *testing.T) {
	hisDB, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}
	defer teardown()

	ft := FakeTree{}
	bc := coin.NewBlockchain(&ft, nil)
	gb := bc.CreateGenesisBlock(genAddress, _genCoins, _genTime)

	hisDB.ProcessBlock(&gb)
}
