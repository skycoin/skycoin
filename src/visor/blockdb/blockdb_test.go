package blockdb

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/util"
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
	Start()

	teardown = func() {
		Stop()
		if err := os.RemoveAll(tmpDir); err != nil {
			panic(err)
		}
	}
	return tmpDir, teardown, nil
}
