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

	// // create blocks bucket.
	// bucketBlocks, err = newBucket([]byte("blocks"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create block signature bucket.
	// bucketBlockSigs, err = newBucket([]byte("blocksigs"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // create unspent output bucket.
	// bucketUtxos, err = newBucket([]byte("utxos"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// Stop the blockdb.
func Stop() {
	db.Close()
}

// Block in blockdb, includes NextHash, for forward walk the block chain.
// type Block struct {
// 	coin.Block
// 	NextHash cipher.SHA256
// }

// SetBlock save coin block.
// func SetBlock(key cipher.SHA256, value []byte) error {
// 	//write to DB
// 	// bin := encoder.Serialize(block) //the value
// 	return bucketBlocks.Set(key[:], value)
// }

// GetBlock by block hash, return nil on not found.
// func GetBlock(hash cipher.SHA256) []byte {
// 	return bucketBlocks.Get(hash[:])
// 	// block := Block{}
// 	// if err := encoder.DeserializeRaw(bin, &block); err != nil {
// 	// 	return nil
// 	// }
// 	// return &block
// }

// FindBlock return block that match the filter.
// func FindBlock(filter func(value []byte) (bool, error)) []byte {
// 	return bucketBlocks.Find(filter)
// 	// b := Block{}
// 	// if err := encoder.DeserializeRaw(bin, &b); err != nil {
// 	// 	return nil
// 	// }
// 	// return &b
// }

// BlockSignature signatures for block
// type BlockSignature struct {
// 	BlockHash     cipher.SHA256 //block.HashHeader
// 	PrevBlockHash cipher.SHA256 //hash of previous block in the chain
// 	Sig           cipher.Sig    //signature of block creator
// 	BkSeq         uint64        //depth of block in tree
// }

// // SetBlockSignature save block signature
// func SetBlockSignature(hash cipher.SHA256, preHash cipher.SHA256, Sig cipher.Sig, BkSeq uint64) error {
// 	var bs = BlockSignature{
// 		hash,
// 		preHash,
// 		Sig,
// 		BkSeq,
// 	}

// 	bin := encoder.Serialize(bs) //value
// 	return bucketBlockSigs.Set(hash[:], bin)
// }

// // GetBlockSignature return nil on not found
// func GetBlockSignature(hash cipher.SHA256) *BlockSignature {
// 	var bs = BlockSignature{}
// 	//get object
// 	bin := bucketBlockSigs.Get(hash[:])
// 	if err := encoder.DeserializeRaw(bin, &bs); err != nil {
// 		return nil
// 	}

// 	return &bs
// }

// SetUnspentOuts save unspent output
// func SetUnspentOuts(hash cipher.SHA256, utxos coin.UxArray) error {
// 	bin := encoder.Serialize(utxos)
// 	return bucketUtxos.Set(hash[:], bin)
// }

// // GetUnspentOuts get unspent outs by hash, return nil on not found.
// func GetUnspentOuts(hash cipher.SHA256) *coin.UxArray {
// 	bin := bucketUtxos.Get(hash[:])
// 	// deserialize uxout
// 	uxs := coin.UxArray{}
// 	if err := encoder.DeserializeRaw(bin, &uxs); err != nil {
// 		return nil
// 	}
// 	return &uxs
// }

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
