package bucket

import (
	"errors"

	"github.com/boltdb/bolt"
)

// Bucket used for grouping the key values in boltdb.
// Also wrap some helper functions.
type Bucket struct {
	Name []byte
	db   *bolt.DB
}

// New create bucket of specific name.
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

// Update use callback func to update the value of given key
func (b *Bucket) Update(key []byte, f func([]byte) ([]byte, error)) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		// get the value of given key
		bkt := tx.Bucket(b.Name)
		if v := bkt.Get(key); v != nil {
			var err error
			v, err = f(v)
			if err != nil {
				return err
			}
			return bkt.Put(key, v)
		}
		return errors.New("not exist in bucket")
	})
}

// Delete removes value of given key
func (b *Bucket) Delete(key []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Name).Delete(key)
	})
}

// GetAll returns all values
func (b *Bucket) GetAll() map[interface{}][]byte {
	values := map[interface{}][]byte{}
	b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b.Name)
		c := bkt.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			values[string(k)] = v
		}
		return nil
	})
	return values
}

// RangeUpdate updates range of the values
func (b *Bucket) RangeUpdate(f func(k, v []byte) ([]byte, error)) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b.Name)
		c := bkt.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			v, err := f(k, v)
			if err != nil {
				return err
			}

			if err := bkt.Put(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}

// IsExist check if the value exist of the given key
func (b *Bucket) IsExist(k []byte) bool {
	var exist bool
	b.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(b.Name).Get(k)
		if v != nil {
			exist = true
		}
		return nil
	})
	return exist
}

// ForEach iterate the whole bucket
func (b *Bucket) ForEach(f func(k, v []byte) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(b.Name).ForEach(f)
	})
}

// Len returns the number of key value pairs
func (b *Bucket) Len() (len int) {
	b.db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(b.Name).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			len++
		}
		return nil
	})
	return
}
