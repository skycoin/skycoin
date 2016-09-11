package blockdb_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"github.com/skycoin/skycoin/src/visor/blockdb"
)

// set rand seed.
var _ = func() int64 {
	t := time.Now().Unix()
	rand.Seed(t)
	return t
}()

func setup(t *testing.T) (string, func(), error) {
	dbName := fmt.Sprintf("%d.db", rand.Int31n(100))
	teardown := func() {}
	tmpDir := filepath.Join(os.TempDir(), dbName)
	if err := os.MkdirAll(tmpDir, 0777); err != nil {
		return "", teardown, err
	}

	util.DataDir = tmpDir
	blockdb.Start()

	teardown = func() {
		blockdb.Stop()
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	return tmpDir, teardown, nil
}

func TestSetAndGetBlocks(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()
	hashs := [10]cipher.SHA256{}
	for i := uint64(0); i < 10; i++ {
		b := coin.Block{}
		b.Head.BkSeq = i
		hashs[i] = b.HashHeader()
		if err := blockdb.SetBlock(b); err != nil {
			t.Fatal(err)
		}
	}

	for i := uint64(0); i < 10; i++ {
		b := blockdb.GetBlock(hashs[i])
		if b == nil {
			t.Fatalf("get block in height: %v failed", i)
		}

		if b.Head.BkSeq != i {
			t.Fatalf("wroing block seq")
		}
	}
}

func TestDisable(t *testing.T) {
	blockdb.Disabled = true

	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()

	b := coin.Block{}
	if err = blockdb.SetBlock(b); err != nil {
		t.Fatal("set block must be nil")
	}

	nb := blockdb.GetBlock(b.HashHeader())
	if nb != nil {
		t.Fatal("get block must be nil")
	}
}

func TestFindBlock(t *testing.T) {
	_, teardown, err := setup(t)
	if err != nil {
		t.Fatal(err)
	}

	defer teardown()
	hashs := make([]cipher.SHA256, 10)
	for i := uint64(0); i < 10; i++ {
		b := coin.Block{}
		b.Head.BkSeq = i
		hashs[i] = b.HashHeader()
		if err := blockdb.SetBlock(b); err != nil {
			t.Fatal(err)
		}
	}

	for i := uint64(0); i < 10; i++ {
		block := blockdb.FindBlock(func(value []byte) (bool, error) {
			b := blockdb.Block{}
			if err := encoder.DeserializeRaw(value, &b); err {
				return false, err
			}

			if b.HashHeader() == hashs[i] {
				return true, nil
			}

			return false, nil
		})

		if block == nil {
			t.Fatal("failed to find block")
		}

		if block.Head.BkSeq != i {
			t.Fatal("failed to find block")
		}
	}

}
