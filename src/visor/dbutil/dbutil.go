package dbutil

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger        = logging.MustGetLogger("dbutil")
	txViewLog     = false
	txViewTrace   = false
	txUpdateLog   = false
	txUpdateTrace = false
)

// DB wraps a bolt.DB to add logging
type DB struct {
	ViewLog     bool
	ViewTrace   bool
	UpdateLog   bool
	UpdateTrace bool
	*bolt.DB
}

// View wraps *bolt.DB.View to add logging
func (db DB) View(f func(*bolt.Tx) error) error {
	if db.ViewLog {
		logger.Debug("db.View starting")
		defer logger.Debug("db.View done")
	}
	if db.ViewTrace {
		debug.PrintStack()
	}
	return db.DB.View(f)
}

// Update wraps *bolt.DB.Update to add logging
func (db DB) Update(f func(*bolt.Tx) error) error {
	if db.UpdateLog {
		logger.Debug("db.Update starting")
		defer logger.Debug("db.Update done")
	}
	if db.UpdateTrace {
		debug.PrintStack()
	}
	return db.DB.Update(f)
}

// WrapDB returns WrapDB
func WrapDB(db *bolt.DB) *DB {
	return &DB{
		ViewLog:     txViewLog,
		UpdateLog:   txUpdateLog,
		ViewTrace:   txViewTrace,
		UpdateTrace: txUpdateTrace,
		DB:          db,
	}
}

// ErrCreateBucketFailed is returned if creating a bolt.DB bucket fails
type ErrCreateBucketFailed struct {
	Bucket string
	Err    error
}

func (e ErrCreateBucketFailed) Error() string {
	return fmt.Sprintf("Create bucket \"%s\" failed: %v", e.Bucket, e.Err)
}

// NewErrCreateBucketFailed returns an ErrCreateBucketFailed
func NewErrCreateBucketFailed(bucket []byte, err error) error {
	return ErrCreateBucketFailed{
		Bucket: string(bucket),
		Err:    err,
	}
}

// ErrBucketNotExist is returned if a bolt.DB bucket does not exist
type ErrBucketNotExist struct {
	Bucket string
}

func (e ErrBucketNotExist) Error() string {
	return fmt.Sprintf("Bucket \"%s\" doesn't exist", e.Bucket)
}

// NewErrBucketNotExist returns an ErrBucketNotExist
func NewErrBucketNotExist(bucket []byte) error {
	return ErrBucketNotExist{
		Bucket: string(bucket),
	}
}

// CreateBuckets creates multiple buckets
func CreateBuckets(tx *bolt.Tx, buckets [][]byte) error {
	for _, b := range buckets {
		if _, err := tx.CreateBucketIfNotExists(b); err != nil {
			return NewErrCreateBucketFailed(b, err)
		}
	}

	return nil
}

// GetBucketObjectDecoded returns an encoder-serialized value from a bucket, decoded to an object
func GetBucketObjectDecoded(tx *bolt.Tx, bktName, key []byte, obj interface{}) (bool, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return false, err
	} else if v == nil {
		return false, nil
	}

	if err := encoder.DeserializeRaw(v, obj); err != nil {
		return false, fmt.Errorf("encoder.DeserializeRaw failed: %v", err)
	}

	return true, nil
}

// GetBucketObjectJSON returns a JSON value from a bucket, unmarshaled to an object
func GetBucketObjectJSON(tx *bolt.Tx, bktName, key []byte, obj interface{}) (bool, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return false, err
	} else if v == nil {
		return false, nil
	}

	if err := json.Unmarshal(v, obj); err != nil {
		return false, fmt.Errorf("json.Unmarshal failed: %v", err)
	}

	return true, nil
}

// GetBucketString returns a string value from a bucket
func GetBucketString(tx *bolt.Tx, bktName, key []byte) (string, bool, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return "", false, err
	} else if v == nil {
		return "", false, nil
	}

	return string(v), true, nil
}

// GetBucketValue returns a []byte value from a bucket
func GetBucketValue(tx *bolt.Tx, bktName, key []byte) ([]byte, error) {
	v, err := getBucketValue(tx, bktName, key)
	if err != nil {
		return nil, err
	} else if v == nil {
		return nil, nil
	}

	// Bytes returned from boltdb are not valid outside of the transaction
	// they are called in, make a copy
	w := make([]byte, len(v))
	copy(w[:], v[:])

	return w, nil
}

// getBucketValue returns a value from a bucket. If the value does not exist,
// it returns an error of type ErrBucketNotExist
func getBucketValue(tx *bolt.Tx, bktName, key []byte) ([]byte, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return nil, NewErrBucketNotExist(bktName)
	}

	v := bkt.Get(key)
	if v == nil {
		return nil, nil
	}

	return v, nil
}

// PutBucketValue puts a value into a bucket under key.
func PutBucketValue(tx *bolt.Tx, bktName, key, val []byte) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.Put(key, val)
}

// BucketHasKey returns true if a bucket has a non-nil value for a key
func BucketHasKey(tx *bolt.Tx, bktName, key []byte) (bool, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return false, NewErrBucketNotExist(bktName)
	}

	v := bkt.Get(key)
	return v != nil, nil
}

// NextSequence returns the NextSequence() from the bucket
func NextSequence(tx *bolt.Tx, bktName []byte) (uint64, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return 0, NewErrBucketNotExist(bktName)
	}

	return bkt.NextSequence()
}

// ForEach calls ForEach on the bucket
func ForEach(tx *bolt.Tx, bktName []byte, f func(k, v []byte) error) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.ForEach(f)
}

// Delete deletes from a bucket
func Delete(tx *bolt.Tx, bktName, key []byte) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.Delete(key)
}

// Len returns the number of keys in a bucket
func Len(tx *bolt.Tx, bktName []byte) (uint64, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return 0, NewErrBucketNotExist(bktName)
	}

	bstats := bkt.Stats()

	if bstats.KeyN < 0 {
		return 0, errors.New("Negative length queried from db stats")
	}

	return uint64(bstats.KeyN), nil
}

// IsEmpty returns true if the bucket is empty
func IsEmpty(tx *bolt.Tx, bktName []byte) (bool, error) {
	length, err := Len(tx, bktName)
	if err != nil {
		return false, err
	}
	return length == 0, nil
}

// Reset resets the bucket
func Reset(tx *bolt.Tx, bktName []byte) error {
	if err := tx.DeleteBucket(bktName); err != nil {
		return err
	}

	return CreateBuckets(tx, [][]byte{bktName})
}

// Itob converts uint64 to bytes
func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Btoi converts bytes to uint64
func Btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}
