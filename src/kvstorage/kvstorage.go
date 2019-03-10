package kvstorage

import (
	"errors"

	"github.com/skycoin/skycoin/src/util/file"
)

var (
	ErrNoSuchKey = errors.New("no such key exists in the storage")
)

// Get gets the value associated with the `key` stored in the file
// with the `fileName`
func Get(fileName, key string) (string, error) {
	var data map[string]string
	if err := file.LoadJSON(fileName, &data); err != nil {
		return "", err
	}

	val, ok := data[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return val, nil
}

// GetAll get all the values stored in the file with the `fileName`
func GetAll(fileName string) (map[string]string, error) {
	var data map[string]string
	if err := file.LoadJSON(fileName, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Add adds the value with the associated `key` and save to the
// file with the `fileName`. Replaces the value if such key already exists
func Add(fileName, key, val string) error {
	var data map[string]string
	if err := file.LoadJSON(fileName, &data); err != nil {
		// either the file doesn't exist or the data is corrupt
		data = make(map[string]string)
	}

	data[key] = val

	return file.SaveJSON(fileName, data, 0644)
}

// Removes the `key` from the data stored in the file with the `fileName`
func Remove(fileName, key string) error {
	var data map[string]string
	if err := file.LoadJSON(fileName, &data); err != nil {
		return err
	}

	if _, ok := data[key]; !ok {
		return ErrNoSuchKey
	}

	delete(data, key)

	return file.SaveJSON(fileName, data, 0644)
}
