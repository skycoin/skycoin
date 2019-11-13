package kvstorage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/util/file"
)

const (
	testDataFilename    = "data" + storageFileExtension
	testEmptyFilename   = "empty" + storageFileExtension
	testCorruptFilename = "corrupt" + storageFileExtension
)

func setupTmpDir(t *testing.T) (string, func()) {
	tmpDir, err := ioutil.TempDir("", "kvstoragetest")
	require.NoError(t, err)

	if err != nil {
		return "", func() {}
	}

	return tmpDir, func() {
		_ = os.RemoveAll(tmpDir) //nolint:errcheck
	}
}

func setupTestFile(t *testing.T, fn string) {
	data := map[string]string{
		"test1": "some value",
		"test2": "{\"key\":\"val\",\"key2\":2}",
	}

	err := file.SaveJSON(fn, data, 0600)
	require.NoError(t, err)
}

func setupEmptyTestFile(t *testing.T, fn string) {
	data := make(map[string]string)

	err := file.SaveJSON(fn, data, 0600)
	require.NoError(t, err)
}

func setupCorruptedTestFile(t *testing.T, fn string) {
	f, err := os.Create(fn)
	require.NoError(t, err)
	defer f.Close()
	_, err = f.Write([]byte("corrupt json file"))
	require.NoError(t, err)
}

func TestNewKVStorage(t *testing.T) {
	type expect struct {
		storage           *kvStorage
		expectError       bool
		expectCorruptFile string
	}

	tmpDir, cleanup := setupTmpDir(t)
	defer cleanup()

	// Setup test data
	dataFilename := filepath.Join(tmpDir, testDataFilename)
	setupTestFile(t, dataFilename)
	// Setup corrupt data file
	corruptDataFilename := filepath.Join(tmpDir, testCorruptFilename)
	setupCorruptedTestFile(t, corruptDataFilename)

	tt := []struct {
		name   string
		fn     string
		expect expect
	}{
		{
			name: "no such file",
			fn:   "nofile.json",
			expect: expect{
				storage:     nil,
				expectError: true,
			},
		},
		{
			name: "file exists",
			fn:   dataFilename,
			expect: expect{
				storage: &kvStorage{
					fn: dataFilename,
					data: map[string]string{
						"test1": "some value",
						"test2": "{\"key\":\"val\",\"key2\":2}",
					},
				},
			},
		},
		{
			name: "corrupted file",
			fn:   corruptDataFilename,
			expect: expect{
				storage: &kvStorage{
					fn:   corruptDataFilename,
					data: map[string]string{}, // an empty file will be when a corrupted file is detected
				},
				expectCorruptFile: corruptDataFilename + ".corrupt.9NGyOAcMBB4",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			storage, err := newKVStorage(tc.fn)
			if tc.expect.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.storage, storage)

			if tc.expect.expectCorruptFile != "" {
				// Check if the file does exist
				_, err := os.Stat(tc.expect.expectCorruptFile)
				require.False(t, os.IsNotExist(err))
			}
		})
	}
}

func TestKVStorageGet(t *testing.T) {
	type expect struct {
		val string
		err error
	}

	tmpDir, cleanup := setupTmpDir(t)
	defer cleanup()

	dataFilename := filepath.Join(tmpDir, testDataFilename)
	setupTestFile(t, dataFilename)

	storage, err := newKVStorage(dataFilename)
	require.NoError(t, err)

	tt := []struct {
		name    string
		storage *kvStorage
		key     string
		expect  expect
	}{
		{
			name:    "no such key",
			storage: storage,
			key:     "key",
			expect: expect{
				err: ErrNoSuchKey,
			},
		},
		{
			name:    "simple string value",
			storage: storage,
			key:     "test1",
			expect: expect{
				val: "some value",
				err: nil,
			},
		},
		{
			name:    "complex marshaled data",
			storage: storage,
			key:     "test2",
			expect: expect{
				val: "{\"key\":\"val\",\"key2\":2}",
				err: nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.storage.get(tc.key)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.val, val)
		})
	}
}

func TestKVStorageGetAll(t *testing.T) {
	type expect struct {
		data map[string]string
	}

	tmpDir, cleanup := setupTmpDir(t)
	defer cleanup()

	dataFilename := filepath.Join(tmpDir, testDataFilename)
	emptyFilename := filepath.Join(tmpDir, testEmptyFilename)

	setupTestFile(t, dataFilename)
	setupEmptyTestFile(t, emptyFilename)

	filledStorage, err := newKVStorage(dataFilename)
	require.NoError(t, err)
	emptyStorage, err := newKVStorage(emptyFilename)
	require.NoError(t, err)

	tt := []struct {
		name    string
		storage *kvStorage
		expect  expect
	}{
		{
			name:    "filled storage",
			storage: filledStorage,
			expect: expect{
				data: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
		{
			name:    "empty storage",
			storage: emptyStorage,
			expect: expect{
				data: map[string]string{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			data := tc.storage.getAll()
			require.Equal(t, tc.expect.data, data)
		})
	}
}

func TestKVStorageAdd(t *testing.T) {
	type expect struct {
		newData     map[string]string
		expectError bool
	}

	tt := []struct {
		name   string
		key    string
		val    string
		expect expect
	}{
		{
			name: "add new value",
			key:  "new key",
			val:  "new value",
			expect: expect{
				newData: map[string]string{
					"test1":   "some value",
					"test2":   "{\"key\":\"val\",\"key2\":2}",
					"new key": "new value",
				},
			},
		},
		{
			name: "replace old value",
			key:  "test1",
			val:  "oiuy",
			expect: expect{
				newData: map[string]string{
					"test1": "oiuy",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, cleanup := setupTmpDir(t)
			defer cleanup()

			dataFilename := filepath.Join(tmpDir, testDataFilename)
			setupTestFile(t, dataFilename)

			storage, err := newKVStorage(dataFilename)
			require.NoError(t, err)

			// acquire the original data
			originalData := storage.getAll()

			err = storage.add(tc.key, tc.val)
			if tc.expect.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			modifiedData := storage.getAll()

			// resave the original data back to file
			err = file.SaveJSON(storage.fn, originalData, 0600)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, modifiedData)
		})
	}
}

func TestKVStorageRemove(t *testing.T) {
	type expect struct {
		newData     map[string]string
		expectError bool
		err         error
	}

	tt := []struct {
		name   string
		key    string
		expect expect
	}{
		{
			name: "no such key",
			key:  "no key",
			expect: expect{
				newData: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
				expectError: true,
				err:         ErrNoSuchKey,
			},
		},
		{
			name: "removed",
			key:  "test1",
			expect: expect{
				newData: map[string]string{
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir, cleanup := setupTmpDir(t)
			defer cleanup()

			dataFilename := filepath.Join(tmpDir, testDataFilename)
			setupTestFile(t, dataFilename)

			storage, err := newKVStorage(dataFilename)
			require.NoError(t, err)

			// acquire the original data
			originalData := storage.getAll()

			err = storage.remove(tc.key)
			if tc.expect.expectError {
				require.Equal(t, tc.expect.err, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			newData := storage.getAll()

			// resave the original data back to file
			err = file.SaveJSON(storage.fn, originalData, 0600)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, newData)
		})
	}
}
