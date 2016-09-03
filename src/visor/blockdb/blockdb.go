package blockdb

//https://github.com/boltdb/bolt
//https://github.com/abhigupta912/mbuckets
//https://github.com/asdine/storm

import (
	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
	"log"
)

/*
	Create 3 buckets. One for
	- blocks
	- block signatures
	- unspent output set
*/

var BlockchainDB *bolt.DB = nil
var DatabaseDirectory string = "" //use src/util to resolve

func StartDB() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.

	database_file = filepath.Join(util.DataDir, "my.db")

	BlockchainDB, err := bolt.Open(database_file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func StopDB() {
	db.Close()
}

//block save/load
func SaveBlock(Block coin.Block) {
	//write to DB
	hash := block.HashHeader()               //the key
	block_binary := encoder.Serialize(Block) //the value

}

//return nil on not found
func GetBlock(BlockHash cipher.SHA256) (Block *coin.Block) {

}

//signatures for block

type BlockSignature struct {
	BlockHash     cipher.SHA256 //block.HashHeader
	PrevBlockHash cipher.SHA256 //hash of previous block in the chain
	Sig           cipher.Sig    //signature of block creator
	BkSeq         uint64        //depth of block in tree
}

//return nil on not found
func SetBlockSignature(BlockHash cipher.SHA256, PrevBlockHash cipher.SHA256, Sig cipher.Sig, BkSeq uint64) *BlockSignature {
	var BS BlockSignature = BlockSignature{
		BlockHash,
		PrevBlockHash,
		Sig,
		BkSeq,
	}

	hash := BlockHash                           //key
	binary := encoder.Serialize(BlockSignature) //value

}

//nil on failure
func GetBlockSignature(BlockHash cipher.SHA256) *BlockSignature {

	var b []byte //grab key/value here
	//get object
	var d BlockSignature = BlockSignature{}

	err := DeserializeRaw(b, &d) //deserialize
	if err != nil {
		log.Panic("blockdb.GetBlockSignature, deserialization error")
		return nil
	}

	return nil
}

//unspent output set?

//save unspent out set at start of every N blocks?

/*
write transaction
err := db.Update(func(tx *bolt.Tx) error {
    ...
    return nil
})
*/

/*
Read Transaction
err := db.View(func(tx *bolt.Tx) error {
    ...
    return nil
})
*/

/*
Batch write
err := db.Batch(func(tx *bolt.Tx) error {
    ...
    return nil
})
*/

/*
Buckets are collections of key/value pairs within the database. All keys in a bucket must be unique. You can create a bucket using the DB.CreateBucket() function:

db.Update(func(tx *bolt.Tx) error {
    b, err := tx.CreateBucket([]byte("MyBucket"))
    if err != nil {
        return fmt.Errorf("create bucket: %s", err)
    }
    return nil
})
You can also create a bucket only if it doesn't exist by using the Tx.CreateBucketIfNotExists() function.
*/
