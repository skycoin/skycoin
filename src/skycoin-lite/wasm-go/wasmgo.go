// Package wasmgo provides WebAssembly binaries for Go.
package wasmgo

import _ "embed"

// WasmFile contains the compiled WebAssembly binary
//
//go:embed skycoin-lite.wasm
var WasmFile []byte

// WasmExecJS contains the JavaScript WebAssembly executor
//
//go:embed wasm_exec.js
var WasmExecJS []byte
