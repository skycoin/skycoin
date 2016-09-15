package blockdb

//https://github.com/boltdb/bolt
//https://github.com/abhigupta912/mbuckets
//https://github.com/asdine/storm

import (
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/util"
)

var db *bolt.DB

// Create 3 buckets. One for
// - blocks
// - block signatures
// - unspent output set

// var bucketBlocks *bucket
// var bucketBlockSigs *bucket
// var bucketUtxos *bucket

// Start the blockdb.
func Start() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbFile := filepath.Join(util.DataDir, "my.db")
	var err error
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Stop the blockdb.
func Stop() {
	db.Close()
}

func UpdateTx(fn func(tx *bolt.Tx) error) error {
	return db.Update(fn)
}
