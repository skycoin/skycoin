// Package wasmtinygo provides the TinyGo-compiled WebAssembly cipher.
//
// The wasm is committed gzipped, and stays that way through //go:embed.
//
// Not to save space in the repository — git compresses its objects, and a raw
// wasm and its gzip both pack to about 1.8 MB, which is what the note above
// *.gz in .gitignore is getting at. The saving is in the binary: go:embed takes
// the file verbatim, so an uncompressed wasm would add 3.1 MB to every build.
// The server then passes the compressed bytes to the browser untouched rather
// than inflating them per request. Skywire embeds its wasm-visor blobs gzipped
// for the same reason (pkg/wasmhv/wasmbin).
package wasmtinygo

import (
	_ "embed"

	"github.com/skycoin/skycoin/src/skycoin-lite/wasmgz"
)

// WasmFileGz contains the compiled WebAssembly binary, gzipped.
//
//go:embed skycoin-lite.wasm.gz
var WasmFileGz []byte

// WasmExecJS contains the JavaScript WebAssembly executor
//
//go:embed wasm_exec.js
var WasmExecJS []byte

// WasmFile returns the compiled WebAssembly binary, decompressed. Serving
// WasmFileGz with Content-Encoding: gzip is cheaper and is what the web wallet
// does; this is for callers that need the bytes themselves.
func WasmFile() ([]byte, error) {
	return wasmgz.Decompress(WasmFileGz)
}
