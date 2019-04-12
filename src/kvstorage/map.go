package kvstorage

// Copy copies the map contents to the new map
func copyMap(data map[string]string) map[string]string {
	copied := make(map[string]string, len(data))

	for k, v := range data {
		copied[k] = v
	}

	return copied
}
