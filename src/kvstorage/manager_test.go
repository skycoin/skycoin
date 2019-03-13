package kvstorage

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/testutil"
)

func TestLoadStorage(t *testing.T) {
	type expect struct {
		expectError bool
		err         error
	}

	tt := []struct {
		name        string
		manager     *Manager
		storageType Type
		expect      expect
	}{
		{
			name: "API disabled",
			manager: &Manager{
				config: Config{
					StorageDir: "./testdata/",
				},
				storages: make(map[Type]*kvStorage),
			},
			storageType: TypeNotes,
			expect: expect{
				expectError: true,
				err:         ErrStorageAPIDisabled,
			},
		},
		{
			name: "unknown storage type",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: make(map[Type]*kvStorage),
			},
			storageType: "unknown",
			expect: expect{
				expectError: true,
				err:         ErrUnknownKVStorageType,
			},
		},
		{
			name: "storage already loaded",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: map[Type]*kvStorage{
					TypeNotes: nil,
				},
			},
			storageType: TypeNotes,
			expect: expect{
				expectError: true,
				err:         ErrStorageAlreadyLoaded,
			},
		},
		{
			name: "OK",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: make(map[Type]*kvStorage),
			},
			storageType: TypeNotes,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.LoadStorage(tc.storageType)
			if tc.expect.expectError {
				require.Equal(t, tc.expect.err, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			testutil.RequireFileExists(t, tc.manager.getStorageFilePath(tc.storageType))
		})
	}

	err := os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}

func TestUnloadStorage(t *testing.T) {
	type expect struct {
		expectError bool
		err         error
	}

	tt := []struct {
		name        string
		manager     *Manager
		storageType Type
		expect      expect
	}{
		{
			name: "API disabled",
			manager: &Manager{
				config: Config{
					StorageDir: "./testdata/",
				},
			},
			storageType: TypeGeneral,
			expect: expect{
				expectError: true,
				err:         ErrStorageAPIDisabled,
			},
		},
		{
			name: "unknown storage type",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: make(map[Type]*kvStorage),
			},
			storageType: "unknown",
			expect: expect{
				expectError: true,
				err:         ErrUnknownKVStorageType,
			},
		},
		{
			name: "no such storage",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: make(map[Type]*kvStorage),
			},
			storageType: TypeGeneral,
			expect: expect{
				expectError: true,
				err:         ErrNoSuchStorage,
			},
		},
		{
			name: "OK",
			manager: &Manager{
				config: Config{
					StorageDir:       "./testdata/",
					EnableStorageAPI: true,
				},
				storages: map[Type]*kvStorage{
					TypeNotes: nil,
				},
			},
			storageType: TypeNotes,
		},
	}

	// init file for tests
	manager := &Manager{
		config: Config{
			EnableStorageAPI: true,
			StorageDir:       "./testdata/",
		},
		storages: make(map[Type]*kvStorage),
	}
	err := manager.LoadStorage(TypeNotes)
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.manager.UnloadStorage(tc.storageType)
			if tc.expect.expectError {
				require.Equal(t, tc.expect.err, err)
			} else {
				require.NoError(t, err)
			}
			if err != nil {
				return
			}

			testutil.RequireFileExists(t, tc.manager.getStorageFilePath(tc.storageType))
		})
	}

	err = os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}

func TestManagerGetStorageValue(t *testing.T) {
	type expect struct {
		val string
		err error
	}

	tt := []struct {
		name              string
		enableAPI         bool
		storageDataDir    string
		loadStorage       bool
		storageTypeToLoad Type
		storageType       Type
		key               string
		expect            expect
	}{
		{
			name:           "API disabled",
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			key:            "key",
			expect: expect{
				err: ErrStorageAPIDisabled,
			},
		},
		{
			name:           "unknown storage type",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    "unknown",
			expect: expect{
				err: ErrUnknownKVStorageType,
			},
		},
		{
			name:           "no such storage",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			expect: expect{
				err: ErrNoSuchStorage,
			},
		},
		{
			name:              "no such key",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "unknown",
			expect: expect{
				err: ErrNoSuchKey,
			},
		},
		{
			name:              "OK - simple string",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "test1",
			expect: expect{
				val: "some value",
			},
		},
		{
			name:              "OK - complex marshaled data",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "test2",
			expect: expect{
				val: "{\"key\":\"val\",\"key2\":2}",
			},
		},
	}

	err := formTestFile(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewManager(NewConfig())
			require.NoError(t, err)
			m.config.EnableStorageAPI = tc.enableAPI
			m.config.StorageDir = tc.storageDataDir

			if tc.loadStorage {
				err := m.LoadStorage(tc.storageTypeToLoad)
				require.NoError(t, err)
			}

			val, err := m.GetStorageValue(tc.storageType, tc.key)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.val, val)
		})
	}

	err = os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}

func TestManagerGetAllStorageValues(t *testing.T) {
	type expect struct {
		data map[string]string
		err  error
	}

	tt := []struct {
		name              string
		enableAPI         bool
		storageDataDir    string
		loadStorage       bool
		storageTypeToLoad Type
		storageType       Type
		expect            expect
	}{
		{
			name:           "API disabled",
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			expect: expect{
				err: ErrStorageAPIDisabled,
			},
		},
		{
			name:           "unknown storage type",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    "unknown",
			expect: expect{
				err: ErrUnknownKVStorageType,
			},
		},
		{
			name:           "no such storage",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			expect: expect{
				err: ErrNoSuchStorage,
			},
		},
		{
			name:              "OK",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			expect: expect{
				data: map[string]string{
					"test1": "some value",
					"test2": "{\"key\":\"val\",\"key2\":2}",
				},
			},
		},
	}

	err := formTestFile(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewManager(NewConfig())
			require.NoError(t, err)
			m.config.EnableStorageAPI = tc.enableAPI
			m.config.StorageDir = tc.storageDataDir

			if tc.loadStorage {
				err := m.LoadStorage(tc.storageTypeToLoad)
				require.NoError(t, err)
			}

			data, err := m.GetAllStorageValues(tc.storageType)
			require.Equal(t, tc.expect.err, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expect.data, data)
		})
	}

	err = os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}

func TestManagerAddStorageValue(t *testing.T) {
	type expect struct {
		expectErr bool
		err       error
	}

	tt := []struct {
		name              string
		enableAPI         bool
		storageDataDir    string
		loadStorage       bool
		storageTypeToLoad Type
		storageType       Type
		key               string
		val               string
		expect            expect
	}{
		{
			name:           "API disabled",
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			key:            "key",
			val:            "val",
			expect: expect{
				expectErr: true,
				err:       ErrStorageAPIDisabled,
			},
		},
		{
			name:           "unknown storage type",
			enableAPI:      true,
			storageDataDir: "./testdata",
			storageType:    "unknown",
			key:            "key",
			val:            "val",
			expect: expect{
				expectErr: true,
				err:       ErrUnknownKVStorageType,
			},
		},
		{
			name:           "no such storage",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			key:            "key",
			val:            "val",
			expect: expect{
				expectErr: true,
				err:       ErrNoSuchStorage,
			},
		},
		{
			name:              "add new value",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "key",
			val:               "val",
		},
		{
			name:              "replace old value",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "test1",
			val:               "oiuy",
		},
	}

	err := formTestFile(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewManager(NewConfig())
			require.NoError(t, err)
			m.config.EnableStorageAPI = tc.enableAPI
			m.config.StorageDir = tc.storageDataDir

			if tc.loadStorage {
				err := m.LoadStorage(tc.storageTypeToLoad)
				require.NoError(t, err)
			}

			err = m.AddStorageValue(tc.storageType, tc.key, tc.val)
			if tc.expect.expectErr {
				require.Equal(t, tc.expect.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	err = os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}

func TestManagerRemoveStorageValue(t *testing.T) {
	type expect struct {
		expectErr bool
		err       error
	}

	tt := []struct {
		name              string
		enableAPI         bool
		storageDataDir    string
		loadStorage       bool
		storageTypeToLoad Type
		storageType       Type
		key               string
		expect            expect
	}{
		{
			name:           "API disabled",
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			key:            "key",
			expect: expect{
				expectErr: true,
				err:       ErrStorageAPIDisabled,
			},
		},
		{
			name:           "unknown storage type",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    "unknown",
			key:            "key",
			expect: expect{
				expectErr: true,
				err:       ErrUnknownKVStorageType,
			},
		},
		{
			name:           "no such storage",
			enableAPI:      true,
			storageDataDir: "./testdata/",
			storageType:    TypeNotes,
			key:            "key",
			expect: expect{
				expectErr: true,
				err:       ErrNoSuchStorage,
			},
		},
		{
			name:              "no such key",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "key",
			expect: expect{
				expectErr: true,
				err:       ErrNoSuchKey,
			},
		},
		{
			name:              "OK",
			enableAPI:         true,
			storageDataDir:    "./testdata/",
			loadStorage:       true,
			storageTypeToLoad: TypeNotes,
			storageType:       TypeNotes,
			key:               "test1",
		},
	}

	err := formTestFile(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewManager(NewConfig())
			require.NoError(t, err)
			m.config.EnableStorageAPI = tc.enableAPI
			m.config.StorageDir = tc.storageDataDir

			if tc.loadStorage {
				err := m.LoadStorage(tc.storageTypeToLoad)
				require.NoError(t, err)
			}

			err = m.RemoveStorageValue(tc.storageType, tc.key)
			if tc.expect.expectErr {
				require.Equal(t, tc.expect.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	err = os.Remove(fmt.Sprintf("%s%s%s", "./testdata/", TypeNotes, storageFileExtension))
	require.NoError(t, err)
}
