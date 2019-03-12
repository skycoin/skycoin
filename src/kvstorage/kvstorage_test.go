package kvstorage

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/util/file"
)

const (
	testDataFilePath  = "./testdata/data" + storageFileExtension
	testEmptyFilePath = "./testdata/empty" + storageFileExtension
)

func formTestFile(fileName string) error {
	data := map[string]string{
		"test1": "some value",
		"test2": "{\"key\":\"val\",\"key2\":2}",
	}

	return file.SaveJSON(fileName, data, 0644)
}

func formEmptyTestFile(fileName string) error {
	data := make(map[string]string)

	return file.SaveJSON(fileName, data, 0644)
}

func TestNewKVStorage(t *testing.T) {
	type expect struct {
		storage     *kvStorage
		expectError bool
	}

	err := formTestFile(testDataFilePath)
	require.NoError(t, err)

	tt := []struct {
		name     string
		fileName string
		expect   expect
	}{
		{
			name:     "no such file",
			fileName: "nofile.json",
			expect: expect{
				storage:     nil,
				expectError: true,
			},
		},
		{
			name:     "file exists",
			fileName: testDataFilePath,
			expect: expect{
				storage: &kvStorage{
					fileName: testDataFilePath,
					data: map[string]string{
						"test1": "some value",
						"test2": "{\"key\":\"val\",\"key2\":2}",
					},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			storage, err := newKVStorage(tc.fileName)
			if tc.expect.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.storage, storage)
		})
	}
}

func TestGet(t *testing.T) {
	type expect struct {
		val string
		err error
	}

	err := formTestFile(testDataFilePath)
	require.NoError(t, err)

	storage, err := newKVStorage(testDataFilePath)
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

func TestGetAll(t *testing.T) {
	type expect struct {
		data map[string]string
	}

	err := formTestFile(testDataFilePath)
	require.NoError(t, err)
	err = formEmptyTestFile(testEmptyFilePath)
	require.NoError(t, err)

	filledStorage, err := newKVStorage(testDataFilePath)
	require.NoError(t, err)
	emptyStorage, err := newKVStorage(testEmptyFilePath)
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

func TestAdd(t *testing.T) {
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
			err := formTestFile(testDataFilePath)
			require.NoError(t, err)

			storage, err := newKVStorage(testDataFilePath)
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
			err = file.SaveJSON(storage.fileName, originalData, 0644)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, modifiedData)
		})
	}
}

func TestRemove(t *testing.T) {
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
			err := formTestFile(testDataFilePath)
			require.NoError(t, err)

			storage, err := newKVStorage(testDataFilePath)
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
			err = file.SaveJSON(storage.fileName, originalData, 0644)
			require.NoError(t, err)

			require.Equal(t, tc.expect.newData, newData)
		})
	}
}
