package kvstorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/skycoin/skycoin/src/util/file"
)

// KVStorageType is a type of a key-value storage
type KVStorageType string

const (
	// KVStorageTypeNotes is a type of storage containing transaction notes
	KVStorageTypeNotes KVStorageType = "notes"
	// KVStorageTypeGeneral is a type of storage for general user data
	KVStorageTypeGeneral KVStorageType = "general"
)

const storageFileExtension = ".json"

var (
	// ErrStorageAPIDisabled is returned while trying to do storage actions while
	// the EnableStorageAPI option is false
	ErrStorageAPIDisabled = NewError(errors.New("storage api is disabled"))
	// ErrNoSuchStorage is returned if no storage with the specified storage type exists
	ErrNoSuchStorage = NewError(errors.New("storage with such type does not exist or is not loaded"))
	// ErrStorageAlreadyLoaded is returned while trying to load already loaded storage
	ErrStorageAlreadyLoaded = NewError(errors.New("storage with such type is already loaded"))
	// ErrUnknownKVStorageType is returned while trying to access the storage of the unknown type
	ErrUnknownKVStorageType = NewError(errors.New("unknown storage type"))
)

// Manager is a manager for key-value storage instances
type Manager struct {
	config   Config
	storages map[KVStorageType]*kvStorage
	sync.RWMutex
}

// NewManager constructs new manager according to the config
func NewManager(c Config) *Manager {
	return &Manager{
		config:   c,
		storages: make(map[KVStorageType]*kvStorage),
	}
}

// LoadStorage loads a new storage instance for the `storageType`
// into the manager. Returns `ErrStorageAlreadyLoaded`, `ErrStorageAPIDisabled`,
// `ErrUnknownKVStorageType`
func (m *Manager) LoadStorage(storageType KVStorageType) error {
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

	fileName := fmt.Sprintf("%s%s%s", m.config.StorageDir,
		storageType, storageFileExtension)

	if !file.Exists(fileName) {
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
func (m *Manager) UnloadStorage(storageType KVStorageType) error {
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

// Get gets the value associated with the `key` from the storage of `storageType.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) Get(storageType KVStorageType, key string) (string, error) {
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

// GetAll gets the snapshot of the current contents from storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) GetAll(storageType KVStorageType) (map[string]string, error) {
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

// Add adds the `val` with the associated `key` to the storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) Add(storageType KVStorageType, key, val string) error {
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

// Remove removes the value with the associated `key` from the storage of `storageType`.
// Returns `ErrNoSuchStorage`, `ErrStorageAPIDisabled`, `ErrUnknownKVStorageType`
func (m *Manager) Remove(storageType KVStorageType, key string) error {
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
func (m *Manager) storageExists(storageType KVStorageType) bool {
	_, ok := m.storages[storageType]

	return ok
}

// isStorageTypeValid validates the given `storageType` against the predefined available types
func isStorageTypeValid(storageType KVStorageType) bool {
	switch storageType {
	case KVStorageTypeNotes, KVStorageTypeGeneral:
		return true
	}

	return false
}

// initEmptyStorage creates a file to persist data
func initEmptyStorage(fileName string) error {
	emptyData := make(map[string]string)

	return file.SaveJSON(fileName, emptyData, 0644)
}
