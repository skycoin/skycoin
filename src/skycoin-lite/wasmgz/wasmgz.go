// Package wasmgz decompresses the gzipped wasm blobs committed under
// src/skycoin-lite.
//
// It exists so the two embed packages do not each carry a copy of the same
// eight lines, and so the decompression has one place to be correct.
package wasmgz

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Decompress returns the contents of a gzip stream.
func Decompress(compressed []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer reader.Close() //nolint:errcheck // read-only, nothing to flush

	return io.ReadAll(reader)
}
