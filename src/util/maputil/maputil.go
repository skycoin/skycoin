// Package maputil provides Go maps related utility methods
package maputil

// Copy copies the map contents to the new map
func Copy(data map[string]string) map[string]string {
	copied := make(map[string]string, len(data))

	for k, v := range data {
		copied[k] = v
	}

	return copied
}
