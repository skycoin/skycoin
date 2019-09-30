package kvstorage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/SkycoinProject/skycoin/src/util/file"
	"github.com/SkycoinProject/skycoin/src/util/logging"
)

// Type is a type of a key-value storage
type Type string

const (
	// TypeTxIDNotes is a type of storage containing transaction notes
	TypeTxIDNotes Type = "txid"
	// TypeGeneral is a type of storage for general user data
	TypeGeneral Type = "client"
)

const storageFileExtension = ".json"

var (
	// ErrStorageAPIDisabled is returned while trying to do storage actions while
	// the EnableStorageAPI option is false
	ErrStorageAPIDisabled = NewError(errors.New("Storage API is disabled"))
	// ErrNoSuchStorage is returned if no storage with the specified storage type loaded
	ErrNoSuchStorage = NewError(errors.New("Storage with such type is not loaded"))
	// ErrStorageAlreadyLoaded is returned while trying to load already loaded storage
	ErrStorageAlreadyLoaded = NewError(errors.New("Storage with such type is already loaded"))
	// ErrUnknownKVStorageType is returned while trying to access the storage of the unknown type
	ErrUnknownKVStorageType = NewError(errors.New("Unknown storage type"))

	logger = logging.MustGetLogger("kvstorage")
)

// Manager is a manager for key-value storage instances
type Manager struct {
	config   Config
	storages map[Type]*kvStorage
	sync.Mutex
}

// NewManager constructs new manager according to the config
func NewManager(c Config) (*Manager, error) {
	logger.Info("Creating new KVStorage manager")

	m := &Manager{
		config:   c,
		storages: make(map[Type]*kvStorage),
	}

	if !strings.HasSuffix(m.config.StorageDir, "/") {
		m.config.StorageDir += "/"
	}

	if !m.config.EnableStorageAPI {
		logger.Info("KVStorage is disabled")
		return m, nil
	}

	if err := os.MkdirAll(m.config.StorageDir, os.FileMode(0700)); err != nil {
		return nil, fmt.Errorf("failed to create kvstorage directory %s: %v", m.config.StorageDir, err)
	}

	for _, t := range m.config.EnabledStorages {
		if err := m.LoadStorage(t); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// LoadStorage loads a new storage instance for the `storageType`
// into the manager. Returns `ErrStorageAlreadyLoaded`, `ErrStorageAPIDisabled`,
// `ErrUnknownKVStorageType`
func (m *Manager) LoadStorage(storageType Type) error {
	if !isStorageTypeValid(storageType) {
		return ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return ErrStorageAPIDisabled
	}

	if m.storageExists(storageType) {
		return ErrStorageAlreadyLoaded
	}

	fn := m.getStorageFilePath(storageType)

	exists, err := file.Exists(fn)
	if err != nil {
		return fmt.Errorf("Manager.LoadStorage file.Exists failed: %v", err)
	}
	if !exists {
		if err := initEmptyStorage(fn); err != nil {
			return fmt.Errorf("Manager.LoadStorage initEmptyStorage failed: %v", err)
		}
	}

	storage, err := newKVStorage(fn)
	if err != nil {
		return err
	}

	m.storages[storageType] = storage

	return nil
}

// UnloadStorage unloads the storage instance for the given `storageType` from the manager.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) UnloadStorage(storageType Type) error {
	if !isStorageTypeValid(storageType) {
		return ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return ErrStorageAPIDisabled
	}

	if !m.storageExists(storageType) {
		return ErrNoSuchStorage
	}

	delete(m.storages, storageType)

	return nil
}

// GetStorageValue gets the value associated with the `key` from the storage of `storageType.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) GetStorageValue(storageType Type, key string) (string, error) {
	if !isStorageTypeValid(storageType) {
		return "", ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return "", ErrStorageAPIDisabled
	}

	if !m.storageExists(storageType) {
		return "", ErrNoSuchStorage
	}

	return m.storages[storageType].get(key)
}

// GetAllStorageValues gets the snapshot of the current contents from storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) GetAllStorageValues(storageType Type) (map[string]string, error) {
	if !isStorageTypeValid(storageType) {
		return nil, ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return nil, ErrStorageAPIDisabled
	}

	if !m.storageExists(storageType) {
		return nil, ErrNoSuchStorage
	}

	return m.storages[storageType].getAll(), nil
}

// AddStorageValue adds the `val` with the associated `key` to the storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) AddStorageValue(storageType Type, key, val string) error {
	if !isStorageTypeValid(storageType) {
		return ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return ErrStorageAPIDisabled
	}

	if !m.storageExists(storageType) {
		return ErrNoSuchStorage
	}

	return m.storages[storageType].add(key, val)
}

// RemoveStorageValue removes the value with the associated `key` from the storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) RemoveStorageValue(storageType Type, key string) error {
	if !isStorageTypeValid(storageType) {
		return ErrUnknownKVStorageType
	}

	m.Lock()
	defer m.Unlock()

	if !m.config.EnableStorageAPI {
		return ErrStorageAPIDisabled
	}

	if !m.storageExists(storageType) {
		return ErrNoSuchStorage
	}

	return m.storages[storageType].remove(key)
}

// storageExists checks whether the storage of `storageType` exists in the manager
func (m *Manager) storageExists(storageType Type) bool {
	_, ok := m.storages[storageType]

	return ok
}

// getStorageFilePath creates the path to the storage of `storageType` in file system
func (m *Manager) getStorageFilePath(storageType Type) string {
	return filepath.Join(m.config.StorageDir, fmt.Sprintf("%s%s", storageType, storageFileExtension))
}

// isStorageTypeValid validates the given `storageType` against the predefined available types
func isStorageTypeValid(storageType Type) bool {
	switch storageType {
	case TypeTxIDNotes, TypeGeneral:
		return true
	}

	return false
}

// initEmptyStorage creates a file to persist data
func initEmptyStorage(fn string) error {
	return file.SaveJSON(fn, map[string]string{}, 0600)
}
