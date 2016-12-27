package blockdb

//https://github.com/boltdb/bolt
//https://github.com/abhigupta912/mbuckets
//https://github.com/asdine/storm

import (
	"log"
	"path/filepath"

	"time"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/util"
)

var db *bolt.DB

// Start the blockdb.
func Start() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	ec := make(chan error)
	go func() {
		dbFile := filepath.Join(util.DataDir, "block.db")
		var err error
		db, err = bolt.Open(dbFile, 0600, nil)
		if err != nil {
			ec <- err
			return
		}
		ec <- nil
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		log.Panic("open boltdb time out")
		return
	case err := <-ec:
		if err != nil {
			log.Panic(err)
		}
	}
}

// Stop the blockdb.
func Stop() {
	db.Close()
}

func UpdateTx(fn func(tx *bolt.Tx) error) error {
	return db.Update(fn)
}
