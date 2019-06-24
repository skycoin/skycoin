/*
Package dbutil provides boltdb utility methods
*/
package dbutil

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/logging"
)

var (
	logger                       = logging.MustGetLogger("dbutil")
	txViewLog                    = false
	txViewTrace                  = false
	txUpdateLog                  = false
	txUpdateTrace                = false
	txDurationLog                = true
	txDurationReportingThreshold = time.Millisecond * 100

	// StopForEach can be returned by a ForEach caller to interrupt ForEach,
	// this is provided for the convenience of the caller, it has no special purpose
	// other than being a sentinel value
	StopForEach = errors.New("Stopped ForEach loop")
)

// Tx wraps a Tx
type Tx struct {
	*bolt.Tx
}

// String is implemented to prevent a panic when mocking methods with *Tx arguments.
// The mock library forces arguments to be printed with %s which causes Tx to panic.
// See https://github.com/stretchr/testify/pull/596
func (tx *Tx) String() string {
	return fmt.Sprintf("%v", tx.Tx)
}

// DB wraps a bolt.DB to add logging
type DB struct {
	ViewLog                    bool
	ViewTrace                  bool
	UpdateLog                  bool
	UpdateTrace                bool
	DurationLog                bool
	DurationReportingThreshold time.Duration

	*bolt.DB

	// shutdownLock is added to prevent closing the database while a View transaction is in progress
	// bolt.DB will block for Update transactions but not for View transactions, and if
	// the database is closed while in a View transaction, it will panic
	// This will be fixed in coreos's bbolt after this PR is merged:
	// https://github.com/coreos/bbolt/pull/91
	// When coreos has this feature, we can switch to coreos's bbolt and remove this lock
	shutdownLock sync.RWMutex
}

// WrapDB returns WrapDB
func WrapDB(db *bolt.DB) *DB {
	return &DB{
		ViewLog:                    txViewLog,
		UpdateLog:                  txUpdateLog,
		ViewTrace:                  txViewTrace,
		UpdateTrace:                txUpdateTrace,
		DurationLog:                txDurationLog,
		DurationReportingThreshold: txDurationReportingThreshold,
		DB:                         db,
	}
}

// View wraps *bolt.DB.View to add logging
func (db *DB) View(name string, f func(*Tx) error) error {
	db.shutdownLock.RLock()
	defer db.shutdownLock.RUnlock()

	if db.ViewLog {
		logger.Debug("db.View [%s] starting", name)
		defer logger.Debug("db.View [%s] done", name)
	}
	if db.ViewTrace {
		debug.PrintStack()
	}

	t0 := time.Now()

	err := db.DB.View(func(tx *bolt.Tx) error {
		return f(&Tx{tx})
	})

	t1 := time.Now()
	delta := t1.Sub(t0)
	if db.DurationLog && delta > db.DurationReportingThreshold {
		logger.Debugf("db.View [%s] elapsed %s", name, delta)
	}

	return err
}

// Update wraps *bolt.DB.Update to add logging
func (db *DB) Update(name string, f func(*Tx) error) error {
	db.shutdownLock.RLock()
	defer db.shutdownLock.RUnlock()

	if db.UpdateLog {
		logger.Debug("db.Update [%s] starting", name)
		defer logger.Debug("db.Update [%s] done", name)
	}
	if db.UpdateTrace {
		debug.PrintStack()
	}

	t0 := time.Now()

	err := db.DB.Update(func(tx *bolt.Tx) error {
		return f(&Tx{tx})
	})

	t1 := time.Now()
	delta := t1.Sub(t0)
	if db.DurationLog && delta > db.DurationReportingThreshold {
		logger.Debugf("db.Update [%s] elapsed %s", name, delta)
	}

	return err
}

// Close closes the underlying *bolt.DB
func (db *DB) Close() error {
	db.shutdownLock.Lock()
	defer db.shutdownLock.Unlock()

	return db.DB.Close()
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
func CreateBuckets(tx *Tx, buckets [][]byte) error {
	for _, b := range buckets {
		if _, err := tx.CreateBucketIfNotExists(b); err != nil {
			return NewErrCreateBucketFailed(b, err)
		}
	}

	return nil
}

// GetBucketObjectDecoded returns an encoder-serialized value from a bucket, decoded to an object
func GetBucketObjectDecoded(tx *Tx, bktName, key []byte, obj interface{}) (bool, error) {
	v, err := GetBucketValueNoCopy(tx, bktName, key)
	if err != nil {
		return false, err
	} else if v == nil {
		return false, nil
	}

	if err := encoder.DeserializeRawExact(v, obj); err != nil {
		return false, fmt.Errorf("encoder.DeserializeRawExact failed: %v", err)
	}

	return true, nil
}

// GetBucketObjectJSON returns a JSON value from a bucket, unmarshaled to an object
func GetBucketObjectJSON(tx *Tx, bktName, key []byte, obj interface{}) (bool, error) {
	v, err := GetBucketValueNoCopy(tx, bktName, key)
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
func GetBucketString(tx *Tx, bktName, key []byte) (string, bool, error) {
	v, err := GetBucketValueNoCopy(tx, bktName, key)
	if err != nil {
		return "", false, err
	} else if v == nil {
		return "", false, nil
	}

	return string(v), true, nil
}

// GetBucketValue returns a []byte value from a bucket. If the bucket does not exist,
// it returns an error of type ErrBucketNotExist
func GetBucketValue(tx *Tx, bktName, key []byte) ([]byte, error) {
	v, err := GetBucketValueNoCopy(tx, bktName, key)
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

// GetBucketValueNoCopy returns a value from a bucket. If the bucket does not exist,
// it returns an error of type ErrBucketNotExist. The byte value is not copied so is not valid
// outside of the database transaction
func GetBucketValueNoCopy(tx *Tx, bktName, key []byte) ([]byte, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return nil, NewErrBucketNotExist(bktName)
	}

	return bkt.Get(key), nil
}

// PutBucketValue puts a value into a bucket under key.
func PutBucketValue(tx *Tx, bktName, key, val []byte) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.Put(key, val)
}

// BucketHasKey returns true if a bucket has a non-nil value for a key
func BucketHasKey(tx *Tx, bktName, key []byte) (bool, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return false, NewErrBucketNotExist(bktName)
	}

	v := bkt.Get(key)
	return v != nil, nil
}

// NextSequence returns the NextSequence() from the bucket
func NextSequence(tx *Tx, bktName []byte) (uint64, error) {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return 0, NewErrBucketNotExist(bktName)
	}

	return bkt.NextSequence()
}

// ForEach calls ForEach on the bucket
func ForEach(tx *Tx, bktName []byte, f func(k, v []byte) error) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.ForEach(f)
}

// Delete deletes from a bucket
func Delete(tx *Tx, bktName, key []byte) error {
	bkt := tx.Bucket(bktName)
	if bkt == nil {
		return NewErrBucketNotExist(bktName)
	}

	return bkt.Delete(key)
}

// Len returns the number of keys in a bucket
func Len(tx *Tx, bktName []byte) (uint64, error) {
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
func IsEmpty(tx *Tx, bktName []byte) (bool, error) {
	length, err := Len(tx, bktName)
	if err != nil {
		return false, err
	}
	return length == 0, nil
}

// Exists returns true if the bucket exists
func Exists(tx *Tx, bktName []byte) bool {
	return tx.Bucket(bktName) != nil
}

// Reset resets the bucket
func Reset(tx *Tx, bktName []byte) error {
	if err := tx.DeleteBucket(bktName); err != nil {
		return err
	}

	_, err := tx.CreateBucket(bktName)
	return err
}

// Itob converts uint64 to bytes
func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

// Btoi converts bytes to uint64
func Btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}
