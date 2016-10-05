package bucket

import "github.com/boltdb/bolt"

// Bucket used for grouping the key values in boltdb.
// Also wrap some helper functions.
type Bucket struct {
	Name []byte
	db   *bolt.DB
}

// NewBucket create bucket of specific name.
func New(name []byte, db *bolt.DB) (*Bucket, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(name); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Bucket{name, db}, nil
}

// Get value of specific key in the bucket.
func (b Bucket) Get(key []byte) []byte {
	var value []byte
	b.db.View(func(tx *bolt.Tx) error {
		value = tx.Bucket(b.Name).Get(key)
		return nil
	})
	return value
}

// Put key value in the bucket.
func (b Bucket) Put(key []byte, value []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Name).Put(key, value)
	})
}

// Find find value that match the filter in the bucket.
func (b Bucket) Find(filter func(key, value []byte) bool) []byte {
	var value []byte
	b.db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket(b.Name)

		c := bt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if filter(k, v) {
				value = v
				break
			}
		}
		return nil
	})
	return value
}

// Count return the number of key/value pairs
// func (b Bucket) Count() uint64 {
// 	var count uint64
// 	b.db.View(func(tx *bolt.Tx) error {
// 		bt := tx.Bucket(b.Name)
// 		return bt.ForEach(func(k, v []byte) error {
// 			count++
// 			return nil
// 		})
// 	})
// 	return count
// }
