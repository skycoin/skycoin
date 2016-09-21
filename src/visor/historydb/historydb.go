package historydb

import (
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

var db *bolt.DB
var gTxns *Transactions

// Start will open a boltdb named history.db
func Start() error {
	dbFile := filepath.Join(util.DataDir, "history.db")
	var err error
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// initialize the transactions instance.
	gTxns, err = newTransactions(db)
	if err != nil {
		return err
	}

	return nil
}

// Stop the historydb.
func Stop() {
	db.Close()
}

// ProcessBlock will index the index the transaction, outputs,etc.
func ProcessBlock(b *coin.Block) error {
	// index the transactions
	for _, t := range b.Body.Transactions {
		tx := Transaction{
			Transaction: t,
			BlockSeq:    b.Seq(),
		}
		if err := gTxns.Add(&tx); err != nil {
			return err
		}
	}
	return nil
}
