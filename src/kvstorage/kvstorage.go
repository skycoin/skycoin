package kvstorage

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/skycoin/skycoin/src/util/file"
)

var (
	// ErrNoSuchKey is returned when the specified key does not exist
	// in the storage instance
	ErrNoSuchKey = NewError(errors.New("no such key exists in the storage"))
)

// kvStorage is a key-value storage for storing arbitrary data
type kvStorage struct {
	fn   string
	data map[string]string
	sync.RWMutex
}

// newKVStorage constructs new storage instance using the file with the filename
// to persist data
func newKVStorage(fn string) (*kvStorage, error) {
	storage := kvStorage{
		fn: fn,
	}

	if err := file.LoadJSON(fn, &storage.data); err != nil {
		logger.Warningf("newKVStorage LoadJSON(%s) failed: %v", fn, err)
		cfp, err := makeCorruptFilePath(fn)
		if err != nil {
			return nil, fmt.Errorf("Failed to make corrupt file path: %v", err)
		}
		if err := os.Rename(fn, cfp); err != nil {
			return nil, fmt.Errorf("Rename %s to %s failed: %v", fn, cfp, err)
		}
		logger.Infof("Backup the corrupted file from: %s to %s", fn, cfp)
		if err := initEmptyStorage(fn); err != nil {
			return nil, err
		}
		storage.data = make(map[string]string)
	}

	return &storage, nil
}

// makeCorruptFilePath creates a $FILE.corrupt.$HASH string based on file path,
// where $HASH is truncated SHA1 of $FILE.
func makeCorruptFilePath(path string) (string, error) {
	fileHash, err := shaFileID(path)
	if err != nil {
		return "", err
	}

	dir, file := filepath.Split(path)
	newFile := fmt.Sprintf("%s.corrupt.%s", file, fileHash)
	newPath := filepath.Join(dir, newFile)

	return newPath, nil
}

// shaFileID return the first 8 bytes of the SHA1 hash of the file,
// hex-encoded
func shaFileID(path string) (string, error) {
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fi.Close()

	h := sha256.New()
	if _, err := io.Copy(h, fi); err != nil {
		return "", err
	}

	sum := h.Sum(nil)
	encodedSum := base64.RawURLEncoding.EncodeToString(sum[:8])
	return encodedSum, nil
}

// get gets the value associated with the `key`. Returns `ErrNoSuchKey`
func (s *kvStorage) get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return val, nil
}

// getAll gets the snapshot of the current storage contents
func (s *kvStorage) getAll() map[string]string {
	s.RLock()
	defer s.RUnlock()

	return copyMap(s.data)
}

// add adds the `val` value to the storage with the specified `key`. Replaces the
// original value if `key` already exists
func (s *kvStorage) add(key, val string) error {
	s.Lock()
	defer s.Unlock()

	// save original data
	oldVal, oldOk := s.data[key]

	s.data[key] = val

	// try to persist data, fall back to original data on error
	if err := s.flush(); err != nil {
		if !oldOk {
			delete(s.data, key)
		} else {
			s.data[key] = oldVal
		}

		return err
	}

	return nil
}

// remove removes the value associated with the `key`. Returns `ErrNoSuchKey`
func (s *kvStorage) remove(key string) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.data[key]; !ok {
		return ErrNoSuchKey
	}

	// save original data
	oldVal := s.data[key]

	delete(s.data, key)

	// try to persist data, fall back to original data on error
	if err := s.flush(); err != nil {
		s.data[key] = oldVal

		return err
	}

	return nil
}

// flush persists data to file
func (s *kvStorage) flush() error {
	return file.SaveJSON(s.fn, s.data, 0600)
}
