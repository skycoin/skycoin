package transactiondb

import (
	"encoding/binary"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/aether/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

var db *bolt.DB

// Start the blockdb.
func Start() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbFile := filepath.Join(util.DataDir, "transactions.db")
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

// Transactions the transaction db instance.
type Transactions struct {
	txns  []byte // transactions bucket name
	depth []byte // transaction depth name.
}

// New create a transaction db instance.
func New() *Transactions {
	t := Transactions{
		txns:  []byte("transactions"),
		depth: []byte("transaction_depth"),
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		// create transaction bucket if does not exist.
		if _, err := tx.CreateBucketIfNotExists(t.txns); err != nil {
			return err
		}

		// create transaction depth bucket if does not exist.
		if _, err := tx.CreateBucketIfNotExists(t.depth); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	return &t
}

// Add transaction to the db.
func (td *Transactions) Add(t *coin.Transaction, depth uint64) error {
	return db.Update(func(tx *bolt.Tx) error {
		txnBkt := tx.Bucket(td.txns)
		key := t.Hash()
		bin := encoder.Serialize(t)
		if err := txnBkt.Put(key[:], bin); err != nil {
			return err
		}

		txDepthBkt := tx.Bucket(td.depth)
		if err := txDepthBkt.Put(key[:], itob(depth)); err != nil {
			return err
		}

		return nil
	})
}

// Get get transaction by tx hash.
func (td Transactions) Get(hash cipher.SHA256) (*coin.Transaction, uint64, error) {
	var txnBin []byte
	var dpBin []byte
	if err := db.View(func(tx *bolt.Tx) error {
		txBkt := tx.Bucket(td.txns)
		txnBin = txBkt.Get(hash[:])

		dpBkt := tx.Bucket(td.depth)
		dpBin = dpBkt.Get(hash[:])
		return nil
	}); err != nil {
		return nil, 0, err
	}

	if txnBin == nil || dpBin == nil {
		return nil, 0, nil
	}

	// deserialize tx
	var tx coin.Transaction
	if err := encoder.DeserializeRaw(txnBin, &tx); err != nil {
		return nil, 0, err
	}

	var dp uint64
	if err := encoder.DeserializeRaw(dpBin, &dp); err != nil {
		return nil, 0, err
	}

	return &tx, dp, nil
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
