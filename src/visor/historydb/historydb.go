package historydb

import (
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
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

// ProcessBlockchain process the blocks in the chain.
func ProcessBlockchain(bc *coin.Blockchain) error {
	depth := bc.Head().Seq()
	for i := uint64(0); i <= depth; i++ {
		b := bc.GetBlockInDepth(i)
		if err := ProcessBlock(b); err != nil {
			return err
		}
	}
	return nil
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

// GetTransaction get transaction by hash.
func GetTransaction(hash cipher.SHA256) (*Transaction, error) {
	return gTxns.Get(hash)
}
