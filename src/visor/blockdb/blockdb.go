package blockdb

//https://github.com/boltdb/bolt
//https://github.com/abhigupta912/mbuckets
//https://github.com/asdine/storm

import (
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util"
)

var db *bolt.DB

// Disabled flag used to determine whether using the boltdb.
var Disabled = false

// Create 3 buckets. One for
// - blocks
// - block signatures
// - unspent output set

var bucketBlocks *bucket
var bucketBlockSigs *bucket
var bucketUtxos *bucket

// Start the blockdb.
func Start() {
	if Disabled {
		return
	}
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbFile := filepath.Join(util.DataDir, "my.db")
	var err error
	db, err = bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	// create blocks bucket.
	bucketBlocks, err = newBucket([]byte("blocks"))
	if err != nil {
		log.Fatal(err)
	}

	// create block signature bucket.
	bucketBlockSigs, err = newBucket([]byte("blocksigs"))
	if err != nil {
		log.Fatal(err)
	}

	// create unspent output bucket.
	bucketUtxos, err = newBucket([]byte("utxos"))
	if err != nil {
		log.Fatal(err)
	}
}

// Stop the blockdb.
func Stop() {
	if Disabled {
		return
	}
	db.Close()
}

type bucket struct {
	Name []byte
}

func newBucket(name []byte) (*bucket, error) {
	if !Disabled {
		err := db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return &bucket{name}, nil
}

func (b bucket) Get(key []byte) []byte {
	if Disabled {
		return nil
	}
	var value []byte
	db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(b.Name).Get(key)
		return nil
	})
	return value
}

func (b bucket) Set(key []byte, value []byte) error {
	if Disabled {
		return nil
	}
	return db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Name).Put(key, value)
	})
}

// SetBlock save coin block.
func SetBlock(block coin.Block) error {
	//write to DB
	hash := block.HashHeader()      //the key
	bin := encoder.Serialize(block) //the value
	return bucketBlocks.Set(hash[:], bin)
}

// GetBlock by block hash, return nil on not found.
func GetBlock(hash cipher.SHA256) *coin.Block {
	bin := bucketBlocks.Get(hash[:])
	block := coin.Block{}
	if err := encoder.DeserializeRaw(bin, &block); err != nil {
		return nil
	}
	return &block
}

// BlockSignature signatures for block
type BlockSignature struct {
	BlockHash     cipher.SHA256 //block.HashHeader
	PrevBlockHash cipher.SHA256 //hash of previous block in the chain
	Sig           cipher.Sig    //signature of block creator
	BkSeq         uint64        //depth of block in tree
}

// SetBlockSignature save block signature
func SetBlockSignature(hash cipher.SHA256, preHash cipher.SHA256, Sig cipher.Sig, BkSeq uint64) error {
	var bs = BlockSignature{
		hash,
		preHash,
		Sig,
		BkSeq,
	}

	bin := encoder.Serialize(bs) //value
	return bucketBlockSigs.Set(hash[:], bin)
}

// GetBlockSignature return nil on not found
func GetBlockSignature(hash cipher.SHA256) *BlockSignature {
	var bs = BlockSignature{}
	//get object
	bin := bucketBlockSigs.Get(hash[:])
	if err := encoder.DeserializeRaw(bin, &bs); err != nil {
		return nil
	}

	return &bs
}

// SetUnspentOuts save unspent output
func SetUnspentOuts(hash cipher.SHA256, utxos coin.UxArray) error {
	bin := encoder.Serialize(utxos)
	return bucketUtxos.Set(hash[:], bin)
}

// GetUnspentOuts get unspent outs by hash, return nil on not found.
func GetUnspentOuts(hash cipher.SHA256) *coin.UxArray {
	bin := bucketUtxos.Get(hash[:])
	// deserialize uxout
	uxs := coin.UxArray{}
	if err := encoder.DeserializeRaw(bin, &uxs); err != nil {
		return nil
	}
	return &uxs
}

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
