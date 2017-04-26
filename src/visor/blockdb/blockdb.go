package blockdb

//https://github.com/boltdb/bolt
//https://github.com/abhigupta912/mbuckets
//https://github.com/asdine/storm

import (
	"path/filepath"

	"time"

	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/util"
)

// var db *bolt.DB

// Open the blockdb.
func Open() (*bolt.DB, func()) {
	dbFile := filepath.Join(util.DataDir, "data.db")
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{
		Timeout: 500 * time.Millisecond,
	})
	if err != nil {
		panic(fmt.Errorf("Open boltdb failed, err:%v", err))
	}
	return db, func() {
		db.Close()
	}
}

// func UpdateTx(fn func(tx *bolt.Tx) error) error {
// 	return db.Update(fn)
// }
