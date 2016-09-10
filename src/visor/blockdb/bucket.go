package blockdb

import "github.com/boltdb/bolt"

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

// MatchFunc callback function for checking if the value does match.
type MatchFunc func(value []byte) bool

// Find the value that matching, return nil on not found.
func (b bucket) Find(match MatchFunc) []byte {
	var value []byte
	db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket(b.Name)

		c := bt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if match(v) {
				value = v
				break
			}
		}
		return nil
	})
	return value
}
