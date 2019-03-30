package kvstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/util/logging"
)

// Type is a type of a key-value storage
type Type string

const (
	// TypeNotes is a type of storage containing transaction notes
	TypeNotes Type = "txid"
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

	logger = logging.MustGetLogger("storagemanager")
)

// Manager is a manager for key-value storage instances
type Manager struct {
	config   Config
	storages map[Type]*kvStorage
	sync.RWMutex
}

// NewManager constructs new manager according to the config
func NewManager(c Config) (*Manager, error) {
	m := &Manager{
		config:   c,
		storages: make(map[Type]*kvStorage),
	}

	if !m.config.EnableStorageAPI {
		logger.Info("Networking is disabled")
		return m, nil
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

	fileName := m.getStorageFilePath(storageType)

	exists, err := file.Exists(fileName)
	if err != nil {
		return err
	}
	if !exists {
		if err := initEmptyStorage(fileName); err != nil {
			return err
		}
	}

	storage, err := newKVStorage(fileName)
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

	m.RLock()
	defer m.RUnlock()

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

	m.RLock()
	defer m.RUnlock()

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

	m.RLock()
	defer m.RUnlock()

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

	m.RLock()
	defer m.RUnlock()

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
	return fmt.Sprintf("%s%s%s", m.config.StorageDir, storageType, storageFileExtension)
}

// isStorageTypeValid validates the given `storageType` against the predefined available types
func isStorageTypeValid(storageType Type) bool {
	switch storageType {
	case TypeNotes, TypeGeneral:
		return true
	}

	return false
}

// initEmptyStorage creates a file to persist data
func initEmptyStorage(fileName string) error {
	emptyData := make(map[string]string)

	return file.SaveJSON(fileName, emptyData, 0644)
}
