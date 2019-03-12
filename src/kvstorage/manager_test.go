package kvstorage

import "testing"

func TestLoadStorage(t *testing.T) {
	type expect struct {
		expectError bool
		err         error
	}

	tt := []struct {
		name        string
		manager *Manager
		storageType KVStorageType
		expect      expect
	}{
		{
			name: "no data file",
		}
	}
}
