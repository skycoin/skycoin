package kvstorage

// Config is a configuration for storage manager
type Config struct {
	StorageDir       string
	EnabledStorages  []Type
	EnableStorageAPI bool
}

// NewConfig creates a default config.
func NewConfig() Config {
	return Config{
		StorageDir: "./data/",
	}
}
