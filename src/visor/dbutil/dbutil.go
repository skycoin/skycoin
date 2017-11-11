package dbutil

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// CreateBucketFailedErr is returned if creating a bolt.DB bucket fails
type CreateBucketFailedErr struct {
	Bucket string
	Err    error
}

func (e CreateBucketFailedErr) Error() string {
	return fmt.Sprintf("Create bucket \"%s\" failed: %v", e.Bucket, e.Err)
}

// NewCreateBucketFailedErr returns an CreateBucketFailedErr
func NewCreateBucketFailedErr(bucket []byte, err error) error {
	return CreateBucketFailedErr{
		Bucket: string(bucket),
		Err:    err,
	}
}

// BucketNotExistErr is returned if a bolt.DB bucket does not exist
type BucketNotExistErr struct {
	Bucket string
}

func (e BucketNotExistErr) Error() string {
	return fmt.Sprintf("Bucket \"%s\" doesn't exist", e.Bucket)
}

// NewBucketNotExistErr returns an BucketNotExistErr
func NewBucketNotExistErr(bucket []byte) error {
	return BucketNotExistErr{
		Bucket: string(bucket),
	}
}

// ObjectNotExistErr is returned if an object specified by "key" is not found in a bolt.DB bucket
type ObjectNotExistErr struct {
	Bucket string
	Key    string
}

func (e ObjectNotExistErr) Error() string {
	return fmt.Sprintf("Object with key \"%s\" not found in bucket \"%s\"", e.Key, e.Bucket)
}

// NewObjectNotExistErr returns an ObjectNotExistErr
func NewObjectNotExistErr(bucket, key []byte) error {
	return ObjectNotExistErr{
		Bucket: string(bucket),
		Key:    string(key),
	}
}

// GetBucketObjectDecoded returns an encoder-serialized value from a bucket, decoded to an object
func GetBucketObjectDecoded(tx *bolt.Tx, bktName, key []byte, obj interface{}) error {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return err
	}

	if err := encoder.DeserializeRaw(v, obj); err != nil {
		return fmt.Errorf("encoder.DeserializeRaw failed: %v", err)
	}

	return nil
}

// GetBucketObjectJSON returns a JSON value from a bucket, unmarshaled to an object
func GetBucketObjectJSON(tx *bolt.Tx, bktName, key []byte, obj interface{}) error {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(v, obj); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	return nil
}

// GetBucketString returns a string value from a bucket
func GetBucketString(tx *bolt.Tx, bktName, key []byte) (string, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return "", err
	}

	return string(v), nil
}

// GetBucketValue returns a []byte value from a bucket
func GetBucketValue(tx *bolt.Tx, bktName, key []byte) ([]byte, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return nil, err
	}

	// Bytes returned from boltdb are not valid outside of the transaction
	// they are called in, make a copy
	w := make([]byte, len(v))
	copy(w[:], v[:])

	return w, nil
}

// getBucketValue returns a value from a bucket. If the value does not exist,
// it returns an error of type BucketNotExistErr
func getBucketValue(tx *bolt.Tx, bktName, key []byte) ([]byte, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return nil, NewBucketNotExistErr(bktName)
	}

	v := bkt.Get(key)
	if v == nil {
		return nil, NewObjectNotExistErr(bktName, key)
	}

	return v, nil
}

// PutBucketValue puts a value into a bucket under key. If the value's type is
// a string, it stores the value as a string. Otherwise, it marshals the value
// to JSON and stores the JSON string.
func PutBucketValue(tx *bolt.Tx, bktName, key []byte, obj interface{}) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewBucketNotExistErr(bktName)
	}

	switch obj.(type) {
	case []byte:
		return bkt.Put(key, obj.([]byte))
	case string:
		return bkt.Put(key, []byte(obj.(string)))
	default:
		v, err := json.Marshal(obj)
		if err != nil {
			return fmt.Errorf("encode value failed: %v", err)
		}
		return bkt.Put(key, v)
	}
}

// BucketHasKey returns true if a bucket has a non-nil value for a key
func BucketHasKey(tx *bolt.Tx, bktName, key []byte) (bool, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return false, NewBucketNotExistErr(bktName)
	}

	v := bkt.Get(key)
	return v != nil, nil
}

// NextSequence returns the NextSequence() from the bucket
func NextSequence(tx *bolt.Tx, bktName []byte) (uint64, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return 0, NewBucketNotExistErr(bktName)
	}

	return bkt.NextSequence()
}

// ForEach calls ForEach on the bucket
func ForEach(tx *bolt.Tx, bktName []byte, f func(k, v []byte) error) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewBucketNotExistErr(bktName)
	}

	return bkt.ForEach(f)
}

// Delete deletes from a bucket
func Delete(tx *bolt.Tx, bktName, key []byte) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewBucketNotExistErr(bktName)
	}

	return bkt.Delete(key)
}

// Len returns the number of keys in a bucket
func Len(tx *bolt.Tx, bktName []byte) (uint64, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return 0, NewBucketNotExistErr(bktName)
	}

	bstats := bkt.Stats()

	if bstats.KeyN < 0 {
		return 0, errors.New("Negative length queried from db stats")
	}

	return uint64(bstats.KeyN), nil
}
