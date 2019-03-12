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
		storageType KVStorageType
		expect      expect
	}{
		{
			name: "API disabled",
			manager: &Manager{
				config: Config{
					StorageDir: "./testdata/",
				},
				storages: make(map[KVStorageType]*kvStorage),
			},
			storageType: KVStorageTypeNotes,
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
				storages: make(map[KVStorageType]*kvStorage),
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
				storages: map[KVStorageType]*kvStorage{
					KVStorageTypeNotes: nil,
				},
			},
			storageType: KVStorageTypeNotes,
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
				storages: make(map[KVStorageType]*kvStorage),
			},
			storageType: KVStorageTypeNotes,
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

			testutil.RequireFileExists(t, fmt.Sprintf("%s%s%s", tc.manager.config.StorageDir,
				tc.storageType, storageFileExtension))
		})
	}

	err := os.Remove("./testdata/notes.json")
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
		storageType KVStorageType
		expect      expect
	}{
		{
			name:        "API disabled",
			manager:     NewManager(NewConfig()),
			storageType: KVStorageTypeGeneral,
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
				storages: make(map[KVStorageType]*kvStorage),
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
				storages: make(map[KVStorageType]*kvStorage),
			},
			storageType: KVStorageTypeGeneral,
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
				storages: map[KVStorageType]*kvStorage{
					KVStorageTypeNotes: nil,
				},
			},
			storageType: KVStorageTypeNotes,
		},
	}

	// init file for tests
	manager := &Manager{
		config: Config{
			EnableStorageAPI: true,
			StorageDir:       "./testdata/",
		},
		storages: make(map[KVStorageType]*kvStorage),
	}
	err := manager.LoadStorage(KVStorageTypeNotes)
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

			testutil.RequireFileExists(t, fmt.Sprintf("%s%s%s", tc.manager.config.StorageDir,
				tc.storageType, storageFileExtension))
		})
	}

	err = os.Remove("./testdata/notes.json")
	require.NoError(t, err)
}
